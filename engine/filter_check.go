package engine

import (
	"fmt"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

// checkFilter rejects a filter that the statement and the declared schema alone make invalid,
// before any row is read (ADR 0002). Row evaluation can't be trusted with this: it never runs on
// an empty table, NULLs compare unequal without erroring, and AND/OR skip a side.
//
// Missing columns are reported before type errors. Unknown conditions, operations and expression
// types fail closed with the same errors row evaluation would raise.
func (engine *DbEngine) checkFilter(expression ast.Expression, table *Table, tableName string, commandName string) error {
	columnNames := make([]string, 0)
	for _, identifier := range expression.GetIdentifiers() {
		columnNames = append(columnNames, identifier.Token.Literal)
	}
	missingColumnName := engine.getMissingColumnName(columnNames, table)
	if missingColumnName != "" {
		return &ColumnDoesNotExistError{tableName: tableName, columnName: missingColumnName}
	}

	return checkTypes(expression, table, commandName)
}

func checkTypes(expression ast.Expression, table *Table, commandName string) error {
	switch mappedExpression := expression.(type) {
	case *ast.OperationExpression:
		if mappedExpression.Operation.Type != token.AND && mappedExpression.Operation.Type != token.OR {
			return &UnsupportedOperationTokenError{mappedExpression.Operation.Literal}
		}
		if err := checkTypes(mappedExpression.Left, table, commandName); err != nil {
			return err
		}
		return checkTypes(mappedExpression.Right, table, commandName)
	case *ast.BooleanExpression:
		return nil
	case *ast.ConditionExpression:
		return checkCondition(mappedExpression, table, commandName)
	case *ast.ContainExpression:
		// The parser only builds IN/NOTIN with a column on the left.
		column := findColumn(mappedExpression.Left.Token.Literal, table)
		for _, item := range mappedExpression.Right {
			if _, ok := coerce(item.Token, declaredTarget(column.Type.Type)); !ok {
				return &InvalidValueTypeError{expectedType: string(column.Type.Type), actualValue: displayValue(item.Token), commandName: commandName}
			}
		}
		return nil
	default:
		return &UnsupportedExpressionTypeError{commandName: commandName, variable: fmt.Sprintf("%s", mappedExpression)}
	}
}

// checkCondition switches on Condition.Type, not Literal: NOT LIKE arrives as Type=LIKE with
// Negated set, and its operand rule doesn't depend on the negation.
func checkCondition(condition *ast.ConditionExpression, table *Table, commandName string) error {
	switch condition.Condition.Type {
	case token.EQUAL, token.NOT:
		return checkSameType(condition.Left, condition.Right, table, commandName)
	case token.LIKE:
		if err := checkLikeOperand(condition.Left, table, commandName); err != nil {
			return err
		}
		return checkLikeOperand(condition.Right, table, commandName)
	default:
		return &UnsupportedConditionalTokenError{variable: condition.Condition.Literal, commandName: commandName}
	}
}

func checkSameType(left ast.Tifier, right ast.Tifier, table *Table, commandName string) error {
	leftColumn := tifierColumn(left, table)
	rightColumn := tifierColumn(right, table)

	switch {
	case leftColumn != nil && rightColumn != nil:
		// Compared directly rather than via declaredTarget, so two unmapped types can't both
		// collapse to NullType and pass.
		if leftColumn.Type.Type != rightColumn.Type.Type {
			return &InvalidValueTypeError{expectedType: string(leftColumn.Type.Type), actualValue: string(rightColumn.Type.Type), commandName: commandName}
		}
		return nil
	case leftColumn != nil:
		return checkLiteralFitsColumn(right.GetToken(), leftColumn, commandName)
	case rightColumn != nil:
		return checkLiteralFitsColumn(left.GetToken(), rightColumn, commandName)
	default:
		leftToken, rightToken := left.GetToken(), right.GetToken()
		target := inferTarget(leftToken)
		if _, ok := coerce(rightToken, target); !ok {
			return &InvalidValueTypeError{expectedType: target.String(), actualValue: displayValue(rightToken), commandName: commandName}
		}
		return nil
	}
}

func checkLiteralFitsColumn(literal token.Token, column *Column, commandName string) error {
	if _, ok := coerce(literal, declaredTarget(column.Type.Type)); !ok {
		return &InvalidValueTypeError{expectedType: string(column.Type.Type), actualValue: displayValue(literal), commandName: commandName}
	}
	return nil
}

// checkLikeOperand reads the column's declared type directly, not via declaredTarget, which maps
// an unknown type to NullType and would let it pass.
func checkLikeOperand(operand ast.Tifier, table *Table, commandName string) error {
	if column := tifierColumn(operand, table); column != nil {
		if column.Type.Type != token.TEXT {
			return &InvalidLikeOperandError{actualType: string(column.Type.Type), commandName: commandName}
		}
		return nil
	}

	target := inferTarget(operand.GetToken())
	if target != StringType && target != NullType {
		return &InvalidLikeOperandError{actualType: target.String(), commandName: commandName}
	}
	return nil
}

// tifierColumn returns the column an identifier names, or nil for a literal. checkFilter has
// already rejected missing columns, so an identifier always resolves.
func tifierColumn(tifier ast.Tifier, table *Table) *Column {
	if !tifier.IsIdentifier() {
		return nil
	}
	return findColumn(tifier.GetToken().Literal, table)
}

func findColumn(name string, table *Table) *Column {
	for _, column := range table.Columns {
		if column.Name == name {
			return column
		}
	}
	return nil
}
