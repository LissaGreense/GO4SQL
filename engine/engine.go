package engine

import (
	"fmt"
	"sort"

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
	for _, command := range commands {

		switch mappedCommand := command.(type) {
		case *ast.WhereCommand:
			continue
		case *ast.OrderByCommand:
			continue
		case *ast.LimitCommand:
			continue
		case *ast.OffsetCommand:
			continue
		case *ast.CreateCommand:
			err := engine.createTable(mappedCommand)
			if err != nil {
				return "", err
			}
			result += "Table '" + mappedCommand.Name.GetToken().Literal + "' has been created\n"
			continue
		case *ast.InsertCommand:
			err := engine.insertIntoTable(mappedCommand)
			if err != nil {
				return "", err
			}
			result += "Data Inserted\n"
			continue
		case *ast.SelectCommand:
			selectOutput, err := engine.getSelectResponse(mappedCommand)
			if err != nil {
				return "", err
			}
			result += selectOutput.ToString() + "\n"
			continue
		case *ast.DeleteCommand:
			deleteCommand := command.(*ast.DeleteCommand)
			if deleteCommand.HasWhereCommand() {
				err := engine.deleteFromTable(mappedCommand, deleteCommand.WhereCommand)
				if err != nil {
					return "", err
				}
			}
			result += "Data from '" + mappedCommand.Name.GetToken().Literal + "' has been deleted\n"
			continue
		case *ast.DropCommand:
			engine.dropTable(mappedCommand)
			result += "Table: '" + mappedCommand.Name.GetToken().Literal + "' has been dropped\n"
			continue
		default:
			return "", fmt.Errorf("unsupported Command detected: %v", command)
		}
	}

	return result, nil
}

// getSelectResponse - Returns Select response basing on ast.OrderByCommand and ast.WhereCommand included in this Select
func (engine *DbEngine) getSelectResponse(selectCommand *ast.SelectCommand) (*Table, error) {
	var table *Table
	var err error

	if selectCommand.HasWhereCommand() {
		whereCommand := selectCommand.WhereCommand
		if selectCommand.HasOrderByCommand() {
			orderByCommand := selectCommand.OrderByCommand
			table, err = engine.selectFromTableWithWhereAndOrderBy(selectCommand, whereCommand, orderByCommand)
			if err != nil {
				return nil, err
			}
		} else {
			table, err = engine.selectFromTableWithWhere(selectCommand, whereCommand)
			if err != nil {
				return nil, err
			}
		}
	} else if selectCommand.HasOrderByCommand() {
		table, err = engine.selectFromTableWithOrderBy(selectCommand, selectCommand.OrderByCommand)
		if err != nil {
			return nil, err
		}
	}

	if table == nil {
		table, err = engine.selectFromTable(selectCommand)
		if err != nil {
			return nil, err
		}
	}

	if selectCommand.HasLimitCommand() || selectCommand.HasOffsetCommand() {
		table.applyOffsetAndLimit(selectCommand)
	}

	return table, nil
}

// createTable - initialize new table in engine with specified name
func (engine *DbEngine) createTable(command *ast.CreateCommand) error {
	_, exist := engine.Tables[command.Name.Token.Literal]

	if exist {
		return fmt.Errorf("table with the name of %s already exist", command.Name.Token.Literal)
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

// insertIntoTable - Insert row of values into the table
func (engine *DbEngine) insertIntoTable(command *ast.InsertCommand) error {
	table, exist := engine.Tables[command.Name.Token.Literal]
	if !exist {
		return fmt.Errorf("table with the name of %s doesn't exist", command.Name.Token.Literal)
	}

	columns := table.Columns

	if len(command.Values) != len(columns) {
		return fmt.Errorf("invalid number of parameters in insert, should be: %d, but got: %d", len(columns), len(columns))
	}

	for i := range columns {
		expectedToken := tokenMapper(columns[i].Type.Type)
		if expectedToken != command.Values[i].Type {
			return fmt.Errorf("invalid Token TokenType in Insert Command, expecting: %s, got: %s", expectedToken, command.Values[i].Type)
		}
		interfaceValue, err := getInterfaceValue(command.Values[i])
		if err != nil {
			return err
		}
		columns[i].Values = append(columns[i].Values, interfaceValue)
	}
	return nil
}

// selectFromTable - Return Table containing all values requested by SelectCommand
func (engine *DbEngine) selectFromTable(command *ast.SelectCommand) (*Table, error) {
	table, exist := engine.Tables[command.Name.Token.Literal]

	if !exist {
		return nil, fmt.Errorf("table with the name of %s doesn't exist", command.Name.Token.Literal)
	}

	return engine.selectFromProvidedTable(command, table)
}

func (engine *DbEngine) selectFromProvidedTable(command *ast.SelectCommand, table *Table) (*Table, error) {
	columns := table.Columns

	wantedColumnNames := make([]string, 0)
	if command.Space[0].Type == token.ASTERISK {
		for i := 0; i < len(columns); i++ {
			wantedColumnNames = append(wantedColumnNames, columns[i].Name)
		}
		return extractColumnContent(columns, &wantedColumnNames)
	} else {
		for i := 0; i < len(command.Space); i++ {
			wantedColumnNames = append(wantedColumnNames, command.Space[i].Literal)
		}
		return extractColumnContent(columns, unique(wantedColumnNames))
	}
}

// deleteFromTable - Delete all rows of data from table that match given condition
func (engine *DbEngine) deleteFromTable(deleteCommand *ast.DeleteCommand, whereCommand *ast.WhereCommand) error {
	table, exist := engine.Tables[deleteCommand.Name.Token.Literal]

	if !exist {
		return fmt.Errorf("table with the name of %s doesn't exist", deleteCommand.Name.Token.Literal)
	}

	newTable, err := engine.getFilteredTable(table, whereCommand, true)

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

// selectFromTableWithWhere - Return Table containing all values requested by SelectCommand and filtered by WhereCommand
func (engine *DbEngine) selectFromTableWithWhere(selectCommand *ast.SelectCommand, whereCommand *ast.WhereCommand) (*Table, error) {
	table, exist := engine.Tables[selectCommand.Name.Token.Literal]

	if !exist {
		return nil, fmt.Errorf("table with the name of %s doesn't exist", selectCommand.Name.Token.Literal)
	}

	if len(table.Columns) == 0 || len(table.Columns[0].Values) == 0 {
		return engine.selectFromProvidedTable(selectCommand, &Table{Columns: []*Column{}})
	}

	filteredTable, err := engine.getFilteredTable(table, whereCommand, false)

	if err != nil {
		return nil, err
	}

	return engine.selectFromProvidedTable(selectCommand, filteredTable)
}

// selectFromTableWithWhereAndOrderBy - Return Table containing all values requested by SelectCommand,
// filtered by WhereCommand and sorted by OrderByCommand
func (engine *DbEngine) selectFromTableWithWhereAndOrderBy(selectCommand *ast.SelectCommand, whereCommand *ast.WhereCommand, orderByCommand *ast.OrderByCommand) (*Table, error) {
	table, exist := engine.Tables[selectCommand.Name.Token.Literal]

	if !exist {
		return nil, fmt.Errorf("table with the name of %s doesn't exist", selectCommand.Name.Token.Literal)
	}

	filteredTable, err := engine.getFilteredTable(table, whereCommand, false)

	if err != nil {
		return nil, err
	}

	emptyTable := getCopyOfTableWithoutRows(table)

	return engine.selectFromProvidedTable(selectCommand, engine.getSortedTable(orderByCommand, filteredTable, emptyTable))
}

// selectFromTableWithOrderBy - Return Table containing all values requested by SelectCommand and sorted by OrderByCommand
func (engine *DbEngine) selectFromTableWithOrderBy(selectCommand *ast.SelectCommand, orderByCommand *ast.OrderByCommand) (*Table, error) {
	table, exist := engine.Tables[selectCommand.Name.Token.Literal]

	if !exist {
		return nil, fmt.Errorf("table with the name of %s doesn't exist", selectCommand.Name.Token.Literal)
	}

	emptyTable := getCopyOfTableWithoutRows(table)

	sortedTable := engine.getSortedTable(orderByCommand, table, emptyTable)

	return engine.selectFromProvidedTable(selectCommand, sortedTable)
}

func (engine *DbEngine) getSortedTable(orderByCommand *ast.OrderByCommand, filteredTable *Table, copyOfTable *Table) *Table {
	sortPatterns := orderByCommand.SortPatterns

	rows := MapTableToRows(filteredTable).rows

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
	return copyOfTable
}

func (engine *DbEngine) getFilteredTable(table *Table, whereCommand *ast.WhereCommand, negation bool) (*Table, error) {
	filteredTable := getCopyOfTableWithoutRows(table)

	for _, row := range MapTableToRows(table).rows {
		fulfilledFilters, err := isFulfillingFilters(row, whereCommand.Expression)
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

func (table *Table) applyOffsetAndLimit(command *ast.SelectCommand) {
	var offset = 0
	var limitRaw = -1

	if command.HasLimitCommand() {
		limitRaw = command.LimitCommand.Count
	}
	if command.HasOffsetCommand() {
		offset = command.OffsetCommand.Count
	}

	for _, column := range table.Columns {
		var limit int

		if limitRaw == -1 || limitRaw+offset > len(column.Values) {
			limit = len(column.Values)
		} else {
			limit = limitRaw + offset
		}

		if offset > len(column.Values) || limit == 0 {
			column.Values = make([]ValueInterface, 0)
		} else {
			column.Values = column.Values[offset:limit]
		}
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

func isFulfillingFilters(row map[string]ValueInterface, expressionTree ast.Expression) (bool, error) {
	switch mappedExpression := expressionTree.(type) {
	case *ast.OperationExpression:
		return processOperationExpression(row, mappedExpression)
	case *ast.BooleanExpression:
		return processBooleanExpression(mappedExpression)
	case *ast.ConditionExpression:
		return processConditionExpression(row, mappedExpression)
	default:
		return false, fmt.Errorf("unsupported expression has been used in WHERE command: %v", expressionTree.GetIdentifiers())
	}
}

func processConditionExpression(row map[string]ValueInterface, conditionExpression *ast.ConditionExpression) (bool, error) {
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
		return false, fmt.Errorf("operation '%s' provided in WHERE command isn't allowed", conditionExpression.Condition.Literal)
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

	return false, fmt.Errorf("unsupported operation token has been used: %s", operationExpression.Operation.Literal)
}

func processBooleanExpression(booleanExpression *ast.BooleanExpression) (bool, error) {
	if booleanExpression.Boolean.Literal == token.TRUE {
		return true, nil
	}
	return false, nil
}

func getTifierValue(tifier ast.Tifier, row map[string]ValueInterface) (ValueInterface, error) {
	switch mappedTifier := tifier.(type) {
	case ast.Identifier:
		return row[mappedTifier.GetToken().Literal], nil
	case ast.Anonymitifier:
		return getInterfaceValue(mappedTifier.GetToken())
	default:
		return nil, fmt.Errorf("couldn't map interface to any implementation of it: %s", tifier.GetToken().Literal)
	}
}
