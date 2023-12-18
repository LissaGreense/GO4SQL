package engine

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

type DbEngine struct {
	Tables map[string]*Table
}

// New Return new DbEngine struct
func New() *DbEngine {
	engine := &DbEngine{}
	engine.Tables = make(map[string]*Table)
	return engine
}

func (engine *DbEngine) CreateTable(command *ast.CreateCommand) {
	_, exist := engine.Tables[command.Name.Token.Literal]

	if exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " already exist!")
	}

	engine.Tables[command.Name.Token.Literal] = &Table{Columns: []*Column{}}
	for i, columnName := range command.ColumnNames {
		engine.Tables[command.Name.Token.Literal].Columns = append(engine.Tables[command.Name.Token.Literal].Columns,
			&Column{
				Type:   command.ColumnTypes[i],
				Values: make([]ValueInterface, 0),
				Name:   columnName,
			})
	}
}

func (engine *DbEngine) InsertIntoTable(command *ast.InsertCommand) {
	table, exist := engine.Tables[command.Name.Token.Literal]
	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	columns := table.Columns

	if len(command.Values) != len(columns) {
		log.Fatal("Invalid number of parameters in insert, should be: " + strconv.Itoa(len(columns)) + ", but got: " + strconv.Itoa(len(columns)))
	}

	for i := 0; i < len(columns); i++ {
		expectedToken := tokenMapper(columns[i].Type.Type)
		if expectedToken != command.Values[i].Type {
			log.Fatal("Invalid Token Type in Insert Command, expecting: " + expectedToken + ", got: " + command.Values[i].Type)
		}
		columns[i].Values = append(columns[i].Values, getInterfaceValue(command.Values[i]))
	}
}

func (engine *DbEngine) SelectFromTable(command *ast.SelectCommand) *Table {
	table, exist := engine.Tables[command.Name.Token.Literal]

	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	return engine.selectFromProvidedTable(command, table)
}

func (engine *DbEngine) selectFromProvidedTable(command *ast.SelectCommand, table *Table) *Table {
	columns := table.Columns

	wantedColumnNames := make([]string, 0)
	if command.Space[0].Type == token.ASTERISK {
		for i := 0; i < len(columns); i++ {
			wantedColumnNames = append(wantedColumnNames, columns[i].Name)
		}
		return extractColumnContent(columns, wantedColumnNames)
	} else {
		for i := 0; i < len(command.Space); i++ {
			wantedColumnNames = append(wantedColumnNames, command.Space[i].Literal)
		}
		return extractColumnContent(columns, unique(wantedColumnNames))
	}
}

func (engine *DbEngine) SelectFromTableWithWhere(selectCommand *ast.SelectCommand, whereCommand *ast.WhereCommand) *Table {
	table, exist := engine.Tables[selectCommand.Name.Token.Literal]

	if !exist {
		log.Fatal("Table with the name of " + selectCommand.Name.Token.Literal + " doesn't exist!")
	}

	filteredTable := getCopyOfTableWithoutRows(table)

	if len(table.Columns) == 0 || len(table.Columns[0].Values) == 0 {
		return engine.selectFromProvidedTable(selectCommand, filteredTable)
	}

	//TODO: maybe rows should have separate structure, so it would would have it's on methods
	rows := mapTableToRows(table)

	for _, row := range rows {
		fulfilledFilters, err := isFulfillingFilters(row, whereCommand.Expression)
		if err != nil {
			log.Fatal(err.Error())
		}

		if fulfilledFilters {
			for _, filteredColumn := range filteredTable.Columns {
				value := row[filteredColumn.Name]
				filteredColumn.Values = append(filteredColumn.Values, value)
			}
		}
	}

	return engine.selectFromProvidedTable(selectCommand, filteredTable)
}

func getCopyOfTableWithoutRows(table *Table) *Table {
	filteredTable := &Table{Columns: []*Column{}}

	for _, column := range table.Columns {
		filteredTable.Columns = append(filteredTable.Columns,
			&Column{
				Type:   column.Type,
				Values: make([]ValueInterface, 0),
				Name:   column.Name,
			})
	}
	return filteredTable
}

func mapTableToRows(table *Table) []map[string]ValueInterface {
	rows := make([]map[string]ValueInterface, 0)

	numberOfRows := len(table.Columns[0].Values)

	for rowIndex := 0; rowIndex < numberOfRows; rowIndex++ {
		row := make(map[string]ValueInterface)
		for _, column := range table.Columns {
			row[column.Name] = column.Values[rowIndex]
		}
		rows = append(rows, row)
	}
	return rows
}

func isFulfillingFilters(row map[string]ValueInterface, expressionTree ast.Expression) (bool, error) {
	operationExpression, operationExpressionIsValid := expressionTree.(*ast.OperationExpression)
	if operationExpressionIsValid {
		return processOperationExpression(row, operationExpression)
	}

	booleanExpression, booleanExpressionIsValid := expressionTree.(*ast.BooleanExpresion)
	if booleanExpressionIsValid {
		return processBooleanExpression(booleanExpression)
	}

	conditionExpression, conditionExpressionIsValid := expressionTree.(*ast.ConditionExpresion)
	if conditionExpressionIsValid {
		return processConditionExpression(row, conditionExpression)
	}

	return false, fmt.Errorf("unsupported expression has been used in WHERE command: %v", expressionTree.GetIdentifiers())
}

func processConditionExpression(row map[string]ValueInterface, conditionExpression *ast.ConditionExpresion) (bool, error) {
	valueLeft, isValueLeftValid := getTifierValue(conditionExpression.Left, row)
	if isValueLeftValid != nil {
		log.Fatal(isValueLeftValid.Error())
	}

	valueRight, isValueRightValid := getTifierValue(conditionExpression.Right, row)
	if isValueLeftValid != nil {
		log.Fatal(isValueRightValid.Error())
	}

	switch conditionExpression.Condition.Type {
	case token.EQUAL:
		return valueLeft.IsEqual(valueRight), nil
	case token.NOT:
		return !(valueLeft.IsEqual(valueRight)), nil
	default:
		return false, errors.New("Operation '" + conditionExpression.Condition.Literal + "' provided in WHERE command isn't allowed!")
	}
}

func processOperationExpression(row map[string]ValueInterface, operationExpression *ast.OperationExpression) (bool, error) {
	if operationExpression.Operation.Type == token.AND {
		left, err := isFulfillingFilters(row, operationExpression.Left)
		if !left {
			return left, err
		}
		right, err := isFulfillingFilters(row, operationExpression.Right)

		return left && right, err
	}

	if operationExpression.Operation.Type == token.OR {
		left, err := isFulfillingFilters(row, operationExpression.Left)
		if left {
			return left, err
		}
		right, err := isFulfillingFilters(row, operationExpression.Right)

		return left || right, err
	}

	return false, errors.New("unsupported operation token has been used: " + operationExpression.Operation.Literal)
}

func processBooleanExpression(booleanExpression *ast.BooleanExpresion) (bool, error) {
	if booleanExpression.Boolean.Literal == token.TRUE {
		return true, nil
	}
	return false, nil
}

func getTifierValue(tifier ast.Tifier, row map[string]ValueInterface) (ValueInterface, error) {
	identifier, identifierIsValid := tifier.(ast.Identifier)

	if identifierIsValid {
		return row[identifier.GetToken().Literal], nil
	}

	anonymitifier, anonymitifierIsValid := tifier.(ast.Anonymitifier)
	if anonymitifierIsValid {
		return getInterfaceValue(anonymitifier.GetToken()), nil
	}

	// TODO: Maybe information in which table this column doesn't exist is needed
	return nil, errors.New("Column name:'" + tifier.GetToken().Literal + "' doesn't exist!")
}
