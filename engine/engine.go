package engine

import (
	"fmt"
	"maps"
	"sort"
	"strconv"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

type DbEngine struct {
	Tables Tables
}
type Tables map[string]*Table

// New Return new DbEngine struct
func New() *DbEngine {
	engine := &DbEngine{}
	engine.Tables = make(Tables)

	return engine
}

// Evaluate - it takes sequences, map them to specific implementation and then process it in SQL engine
func (engine *DbEngine) Evaluate(sequences *ast.Sequence) (string, error) {
	commands := sequences.Commands

	result := ""
	var err error
	for _, command := range commands {
		if err != nil {
			return "", err
		}
		switch mappedCommand := command.(type) {
		case *ast.WhereCommand:
			continue
		case *ast.OrderByCommand:
			continue
		case *ast.LimitCommand:
			continue
		case *ast.OffsetCommand:
			continue
		case *ast.JoinCommand:
			continue
		case *ast.CreateCommand:
			err = engine.createTable(mappedCommand)
			result += "Table '" + mappedCommand.Name.GetToken().Literal + "' has been created\n"
			continue
		case *ast.InsertCommand:
			err = engine.insertIntoTable(mappedCommand)
			result += "Data Inserted\n"
			continue
		case *ast.SelectCommand:
			var selectOutput *Table
			selectOutput, err = engine.getSelectResponse(mappedCommand)
			result += selectOutput.ToString() + "\n"
			continue
		case *ast.DeleteCommand:
			deleteCommand := command.(*ast.DeleteCommand)
			if deleteCommand.HasWhereCommand() {
				err = engine.deleteFromTable(mappedCommand, deleteCommand.WhereCommand)
			}
			result += "Data from '" + mappedCommand.Name.GetToken().Literal + "' has been deleted\n"
			continue
		case *ast.DropCommand:
			engine.dropTable(mappedCommand)
			result += "Table: '" + mappedCommand.Name.GetToken().Literal + "' has been dropped\n"
			continue
		case *ast.UpdateCommand:
			err = engine.updateTable(mappedCommand)
			result += "Table: '" + mappedCommand.Name.GetToken().Literal + "' has been updated\n"
			continue
		default:
			return "", &UnsupportedCommandTypeFromParserError{variable: fmt.Sprintf("%s", command)}
		}
	}
	return result, err
}

// getSelectResponse - processes a SELECT query represented by the ast.SelectCommand and applies a pipeline of
// transformations based on options applied to ast.SelectCommand
func (engine *DbEngine) getSelectResponse(selectCommand *ast.SelectCommand) (*Table, error) {
	var table *Table
	var err error

	if selectCommand.HasJoinCommand() {
		table, err = engine.joinTables(selectCommand.JoinCommand, selectCommand.Name.Token.Literal)
	} else {
		var exists bool
		table, exists = engine.Tables[selectCommand.Name.Token.Literal]
		if !exists {
			return nil, &TableDoesNotExistError{selectCommand.Name.Token.Literal}
		}
	}

	if err != nil {
		return nil, err
	}

	processor := NewSelectProcessor(engine, selectCommand)

	// Build the transformation pipeline using the builder pattern
	if selectCommand.HasOrderByCommand() {
		processor.WithOrderByClause()
	}

	if selectCommand.HasWhereCommand() {
		processor.WithWhereClause()
	}

	// If no WHERE or ORDER BY, the vanilla select (projection) is applied first.
	// Otherwise, WHERE/ORDER BY are applied, and then the projection happens within them (handled by their respective transformers).
	if !selectCommand.HasOrderByCommand() && !selectCommand.HasWhereCommand() {
		processor.WithVanillaSelectClause()
	}

	if selectCommand.HasOffsetCommand() || selectCommand.HasLimitCommand() {
		processor.WithOffsetLimitClause()
	}

	if selectCommand.HasDistinct {
		processor.WithDistinctClause()
	}

	return processor.Process(table)
}

// createTable - initialize new table in engine with specified name
func (engine *DbEngine) createTable(command *ast.CreateCommand) error {
	_, exist := engine.Tables[command.Name.Token.Literal]

	if exist {
		return &TableAlreadyExistsError{command.Name.Token.Literal}
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
	return nil
}

func (engine *DbEngine) updateTable(command *ast.UpdateCommand) error {
	table, exists := engine.Tables[command.Name.Token.Literal]
	if !exists {
		return &TableDoesNotExistError{command.Name.Token.Literal}
	}

	columnIndices := make(map[string]int, len(table.Columns))
	for i, col := range table.Columns {
		columnIndices[col.Name] = i
	}

	// Map changes to column indices
	type change struct {
		index int
		value ast.Anonymitifier
	}

	changes := make([]change, 0, len(command.Changes))
	for colToken, newValue := range command.Changes {
		colName := colToken.Literal
		colIndex, ok := columnIndices[colName]
		if !ok {
			return &ColumnDoesNotExistError{
				tableName:  command.Name.Token.Literal,
				columnName: colName,
			}
		}
		changes = append(changes, change{index: colIndex, value: newValue})
	}

	for rowIndex := 0; rowIndex < len(table.Columns[0].Values); rowIndex++ {
		if command.HasWhereCommand() {
			matches, err := isFulfillingFilters(getRow(table, rowIndex), command.WhereCommand.Expression, command.WhereCommand.Token.Literal)
			if err != nil {
				return err
			}
			if !matches {
				continue
			}
		}

		for _, change := range changes {
			val, err := getInterfaceValue(change.value.GetToken())
			if err != nil {
				return err
			}
			table.Columns[change.index].Values[rowIndex] = val
		}
	}

	return nil
}

// insertIntoTable - Insert row of values into the table
func (engine *DbEngine) insertIntoTable(command *ast.InsertCommand) error {
	table, exist := engine.Tables[command.Name.Token.Literal]
	if !exist {
		return &TableDoesNotExistError{command.Name.Token.Literal}
	}

	columns := table.Columns

	if len(command.Values) != len(columns) {
		return &InvalidNumberOfParametersError{
			expectedNumber: len(columns),
			actualNumber:   len(command.Values),
			commandName:    command.Token.Literal,
		}
	}

	for i := range columns {
		expectedToken := tokenMapper(columns[i].Type.Type)
		colValueType := command.Values[i].Type

		if (expectedToken != colValueType) && (colValueType != token.NULL) {
			return &InvalidValueTypeError{
				expectedType: string(expectedToken),
				actualType:   string(colValueType),
				commandName:  command.Token.Literal,
			}
		}
		interfaceValue, err := getInterfaceValue(command.Values[i])
		if err != nil {
			return err
		}
		columns[i].Values = append(columns[i].Values, interfaceValue)
	}
	return nil
}

func (engine *DbEngine) selectFromProvidedTable(command *ast.SelectCommand, table *Table) (*Table, error) {
	columns := table.Columns

	wantedColumnNames := make([]string, 0)
	if command.AggregateFunctionAppears() {
		selectedTable := &Table{Columns: make([]*Column, 0)}

		for i := 0; i < len(command.Space); i++ {
			col := &Column{}
			var columnValues []ValueInterface
			currentSpace := command.Space[i]

			if currentSpace.ColumnName.Type == token.ASTERISK && currentSpace.AggregateFunc.Type == token.COUNT {
				if len(columns) > 0 {
					columnValues = columns[0].Values
				}
			} else {
				var err error
				columnValues, err = getValuesOfColumn(currentSpace.ColumnName.Literal, columns)
				if err != nil {
					return nil, err
				}
			}

			if currentSpace.ContainsAggregateFunc() {
				col.Name = fmt.Sprintf("%s(%s)", currentSpace.AggregateFunc.Literal,
					currentSpace.ColumnName.Literal)
				col.Type = evaluateColumnTypeOfAggregateFunc(currentSpace)
				aggregatedValue, err := aggregateColumnContent(currentSpace, columnValues)
				if err != nil {
					return nil, err
				}
				col.Values = []ValueInterface{aggregatedValue}
			} else {
				col.Name = currentSpace.ColumnName.Literal
				col.Type = currentSpace.ColumnName
				col.Values = []ValueInterface{columnValues[0]}
			}
			selectedTable.Columns = append(selectedTable.Columns, col)
		}
		return selectedTable, nil
	} else if command.Space[0].ColumnName.Type == token.ASTERISK {
		for i := 0; i < len(columns); i++ {
			wantedColumnNames = append(wantedColumnNames, columns[i].Name)
		}
		return extractColumnContent(columns, &wantedColumnNames, command.Name.GetToken().Literal)
	} else {
		for i := 0; i < len(command.Space); i++ {
			wantedColumnNames = append(wantedColumnNames, command.Space[i].ColumnName.Literal)
		}
		return extractColumnContent(columns, unique(wantedColumnNames), command.Name.GetToken().Literal)
	}
}

func getValuesOfColumn(columnName string, columns []*Column) ([]ValueInterface, error) {
	wantedColumnName := []string{columnName}
	columnContent, err := extractColumnContent(columns, &wantedColumnName, "")
	if err != nil {
		return nil, err
	}
	return columnContent.Columns[0].Values, nil
}

func evaluateColumnTypeOfAggregateFunc(space ast.Space) token.Token {
	if space.AggregateFunc.Type == token.MIN ||
		space.AggregateFunc.Type == token.MAX {
		return space.ColumnName
	}
	return token.Token{Type: token.INT, Literal: "INT"}
}

func aggregateColumnContent(space ast.Space, columnValues []ValueInterface) (ValueInterface, error) {
	if space.AggregateFunc.Type == token.COUNT {
		return getCount(space, columnValues)
	}
	if len(columnValues) == 0 {
		return NullValue{}, nil
	}
	switch space.AggregateFunc.Type {
	case token.MAX:
		return getMax(columnValues)
	case token.MIN:
		return getMin(columnValues)
	case token.SUM:
		return getSum(columnValues)
	default:
		return getAvg(columnValues)
	}
}

func getCount(space ast.Space, columnValues []ValueInterface) (*IntegerValue, error) {
	if space.ColumnName.Type == token.ASTERISK {
		return &IntegerValue{Value: len(columnValues)}, nil
	}
	count := 0
	for _, value := range columnValues {
		if value.GetType() != NullType {
			count++
		}
	}
	return &IntegerValue{Value: count}, nil
}

func getAvg(columnValues []ValueInterface) (*IntegerValue, error) {
	sum, err := getSum(columnValues)
	if err != nil {
		return nil, err
	}
	return &IntegerValue{Value: sum.Value / len(columnValues)}, nil
}

func getSum(columnValues []ValueInterface) (*IntegerValue, error) {
	if columnValues[0].GetType() == StringType {
		return &IntegerValue{Value: 0}, nil
	} else {
		sum := 0
		for _, value := range columnValues {
			if value.GetType() != NullType {
				num, err := strconv.Atoi(value.ToString())
				if err != nil {
					return nil, err
				}
				sum += num
			}
		}
		return &IntegerValue{Value: sum}, nil
	}
}

// deleteFromTable - Delete all rows of data from table that match given condition
func (engine *DbEngine) deleteFromTable(deleteCommand *ast.DeleteCommand, whereCommand *ast.WhereCommand) error {
	table, exist := engine.Tables[deleteCommand.Name.Token.Literal]

	if !exist {
		return &TableDoesNotExistError{deleteCommand.Name.Token.Literal}
	}

	newTable, err := engine.getFilteredTable(table, whereCommand, true, deleteCommand.Name.Token.Literal)

	if err != nil {
		return err
	}
	engine.Tables[deleteCommand.Name.Token.Literal] = newTable

	return nil
}

// dropTable - Drop table with given name
func (engine *DbEngine) dropTable(dropCommand *ast.DropCommand) {
	delete(engine.Tables, dropCommand.Name.GetToken().Literal)
}

func (engine *DbEngine) getSortedTable(orderByCommand *ast.OrderByCommand, table *Table, copyOfTable *Table, tableName string) (*Table, error) {
	sortPatterns := orderByCommand.SortPatterns

	columnNames := make([]string, 0)
	for _, sortPattern := range sortPatterns {
		columnNames = append(columnNames, sortPattern.ColumnName.Literal)
	}

	missingColName := engine.getMissingColumnName(columnNames, table)
	if missingColName != "" {
		return nil, &ColumnDoesNotExistError{
			tableName:  tableName,
			columnName: missingColName,
		}
	}

	rows := MapTableToRows(table).rows

	sort.Slice(rows, func(i, j int) bool {
		howDeepWeSort := 0
		sortingType := sortPatterns[howDeepWeSort].Order.Type
		columnToSort := sortPatterns[howDeepWeSort].ColumnName.Literal

		for rows[i][columnToSort].IsEqual(rows[j][columnToSort]) {
			howDeepWeSort++
			sortingType = sortPatterns[howDeepWeSort].Order.Type

			if howDeepWeSort >= len(orderByCommand.SortPatterns) {
				return true
			}
			columnToSort = sortPatterns[howDeepWeSort].ColumnName.Literal
		}

		if sortingType == token.DESC {
			return rows[i][columnToSort].isGreaterThan(rows[j][columnToSort])
		}

		return rows[i][columnToSort].isSmallerThan(rows[j][columnToSort])
	})

	for _, row := range rows {
		for _, newColumn := range copyOfTable.Columns {
			value := row[newColumn.Name]
			newColumn.Values = append(newColumn.Values, value)
		}
	}
	return copyOfTable, nil
}

func (engine *DbEngine) getMissingColumnName(columnNames []string, table *Table) string {
	for _, columnName := range columnNames {
		exists := false
		for _, column := range table.Columns {
			if column.Name == columnName {
				exists = true
				break
			}
		}
		if !exists {
			return columnName
		}
	}
	return ""
}

func (engine *DbEngine) getFilteredTable(table *Table, whereCommand *ast.WhereCommand, negation bool, tableName string) (*Table, error) {
	filteredTable := getCopyOfTableWithoutRows(table)

	identifiers := whereCommand.Expression.GetIdentifiers()
	columnNames := make([]string, 0)
	for _, identifier := range identifiers {
		columnNames = append(columnNames, identifier.Token.Literal)
	}
	missingColumnName := engine.getMissingColumnName(columnNames, table)
	if missingColumnName != "" {
		return nil, &ColumnDoesNotExistError{tableName: tableName, columnName: missingColumnName}
	}

	for _, row := range MapTableToRows(table).rows {
		fulfilledFilters, err := isFulfillingFilters(row, whereCommand.Expression, whereCommand.Token.Literal)
		if err != nil {
			return nil, err
		}

		if xor(fulfilledFilters, negation) {
			for _, filteredColumn := range filteredTable.Columns {
				value := row[filteredColumn.Name]
				filteredColumn.Values = append(filteredColumn.Values, value)
			}
		}
	}
	return filteredTable, nil
}

func (engine *DbEngine) joinTables(joinCommand *ast.JoinCommand, leftTableName string) (*Table, error) {
	leftTable, exist := engine.Tables[leftTableName]
	leftTablePrefix := leftTableName + "."
	if !exist {
		return nil, &TableDoesNotExistError{leftTableName}
	}

	rightTableName := joinCommand.Name.Token.Literal
	rightTablePrefix := rightTableName + "."
	rightTable, exist := engine.Tables[rightTableName]
	if !exist {
		return nil, &TableDoesNotExistError{rightTableName}
	}

	joinedTable := &Table{Columns: []*Column{}}

	addColumnsWithPrefix(joinedTable, leftTable.Columns, leftTablePrefix)
	addColumnsWithPrefix(joinedTable, rightTable.Columns, rightTablePrefix)

	leftTableWithAddedPrefix := leftTable.getTableCopyWithAddedPrefixToColumnNames(leftTablePrefix)
	rightTableWithAddedPrefix := rightTable.getTableCopyWithAddedPrefixToColumnNames(rightTablePrefix)
	var unmatchedRightRows = make(map[int]bool)

	for leftRowIndex := 0; leftRowIndex < len(leftTable.Columns[0].Values); leftRowIndex++ {
		joinedRowLeft := getRow(leftTableWithAddedPrefix, leftRowIndex)
		leftRowMatches := false

		for rightRowIndex := 0; rightRowIndex < len(rightTable.Columns[0].Values); rightRowIndex++ {
			joinedRowRight := getRow(rightTableWithAddedPrefix, rightRowIndex)
			maps.Copy(joinedRowRight, joinedRowLeft)

			fulfilledFilters, err := isFulfillingFilters(joinedRowRight, joinCommand.Expression, joinCommand.Token.Literal)
			if err != nil {
				return nil, err
			}

			isLastLeftRow := leftRowIndex == len(leftTable.Columns[0].Values)-1

			if fulfilledFilters {
				for colIndex, column := range joinedTable.Columns {
					joinedTable.Columns[colIndex].Values = append(joinedTable.Columns[colIndex].Values, joinedRowRight[column.Name])
				}
				leftRowMatches, unmatchedRightRows[rightRowIndex] = true, true
			} else if isLastLeftRow && joinCommand.ShouldTakeRightSide() && !unmatchedRightRows[rightRowIndex] {
				joinedRowRight = getRow(rightTableWithAddedPrefix, rightRowIndex)
				aggregateRowIntoJoinTable(leftTableWithAddedPrefix, joinedRowRight, joinedTable)
			}
		}

		if joinCommand.ShouldTakeLeftSide() && !leftRowMatches {
			aggregateRowIntoJoinTable(rightTableWithAddedPrefix, joinedRowLeft, joinedTable)
		}
	}

	return joinedTable, nil
}

func aggregateRowIntoJoinTable(tableWithAddedPrefix *Table, joinedRow map[string]ValueInterface, joinedTable *Table) {
	joinedEmptyRow := getEmptyRow(tableWithAddedPrefix)
	maps.Copy(joinedRow, joinedEmptyRow)
	for colIndex, column := range joinedTable.Columns {
		joinedTable.Columns[colIndex].Values = append(joinedTable.Columns[colIndex].Values, joinedRow[column.Name])
	}
}

func addColumnsWithPrefix(finalTable *Table, columnsToAdd []*Column, prefix string) {
	for _, column := range columnsToAdd {
		finalTable.Columns = append(finalTable.Columns,
			&Column{
				Type:   column.Type,
				Values: make([]ValueInterface, 0),
				Name:   prefix + column.Name,
			})
	}
}

func xor(fulfilledFilters bool, negation bool) bool {
	return (fulfilledFilters || negation) && !(fulfilledFilters && negation)
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

func isFulfillingFilters(row map[string]ValueInterface, expressionTree ast.Expression, commandName string) (bool, error) {
	switch mappedExpression := expressionTree.(type) {
	case *ast.OperationExpression:
		return processOperationExpression(row, mappedExpression, commandName)
	case *ast.BooleanExpression:
		return processBooleanExpression(mappedExpression), nil
	case *ast.ConditionExpression:
		return processConditionExpression(row, mappedExpression, commandName)
	case *ast.ContainExpression:
		return processContainExpression(row, mappedExpression)

	default:
		return false, &UnsupportedExpressionTypeError{commandName: commandName, variable: fmt.Sprintf("%s", mappedExpression)}
	}
}

func processConditionExpression(row map[string]ValueInterface, conditionExpression *ast.ConditionExpression, commandName string) (bool, error) {
	valueLeft, err := getTifierValue(conditionExpression.Left, row)
	if err != nil {
		return false, err
	}

	valueRight, err := getTifierValue(conditionExpression.Right, row)
	if err != nil {
		return false, err
	}

	switch conditionExpression.Condition.Type {
	case token.EQUAL:
		return valueLeft.IsEqual(valueRight), nil
	case token.NOT:
		return !(valueLeft.IsEqual(valueRight)), nil
	default:
		return false, &UnsupportedConditionalTokenError{variable: conditionExpression.Condition.Literal, commandName: commandName}
	}
}

func processContainExpression(row map[string]ValueInterface, containExpression *ast.ContainExpression) (bool, error) {
	valueLeft, err := getTifierValue(containExpression.Left, row)
	if err != nil {
		return false, err
	}

	result, err := ifValueInterfaceInArray(containExpression.Right, valueLeft)

	if containExpression.Contains {
		return result, err
	}

	return !result, err
}

func ifValueInterfaceInArray(array []ast.Anonymitifier, valueLeft ValueInterface) (bool, error) {
	for _, expectedValue := range array {
		value, err := getInterfaceValue(expectedValue.Token)
		if err != nil {
			return false, err
		}
		if value.IsEqual(valueLeft) {
			return true, nil
		}
	}
	return false, nil
}

func processOperationExpression(row map[string]ValueInterface, operationExpression *ast.OperationExpression, commandName string) (bool, error) {
	if operationExpression.Operation.Type == token.AND {
		left, err := isFulfillingFilters(row, operationExpression.Left, commandName)
		if !left {
			return left, err
		}
		right, err := isFulfillingFilters(row, operationExpression.Right, commandName)

		return right, err
	}

	if operationExpression.Operation.Type == token.OR {
		left, err := isFulfillingFilters(row, operationExpression.Left, commandName)
		if left {
			return left, err
		}
		right, err := isFulfillingFilters(row, operationExpression.Right, commandName)

		return right, err
	}

	return false, &UnsupportedOperationTokenError{operationExpression.Operation.Literal}
}

func processBooleanExpression(booleanExpression *ast.BooleanExpression) bool {
	if booleanExpression.Boolean.Literal == token.TRUE {
		return true
	}
	return false
}

func getTifierValue(tifier ast.Tifier, row map[string]ValueInterface) (ValueInterface, error) {
	switch mappedTifier := tifier.(type) {
	case ast.Identifier:
		value, ok := row[mappedTifier.GetToken().Literal]
		if ok == false {
			return nil, &ColumnDoesNotExistError{tableName: "", columnName: mappedTifier.GetToken().Literal}
		}
		return value, nil
	case ast.Anonymitifier:
		return getInterfaceValue(mappedTifier.GetToken())
	default:
		return nil, &UnsupportedValueType{tifier.GetToken().Literal}
	}
}
