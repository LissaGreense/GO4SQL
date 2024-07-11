package engine

import "strconv"

// TableAlreadyExistsError - error thrown when user tries to create table using name that already
// exists in database
type TableAlreadyExistsError struct {
	tableName string
}

func (m *TableAlreadyExistsError) Error() string {
	return "table with the name of " + m.tableName + " already exists"
}

// TableDoesNotExistError - error thrown when user tries to make operation on un-existing table
type TableDoesNotExistError struct {
	tableName string
}

func (m *TableDoesNotExistError) Error() string {
	return "table with the name of " + m.tableName + " doesn't exist"
}

// ColumnDoesNotExistError - error thrown when user tries to make operation on un-existing column
type ColumnDoesNotExistError struct {
	tableName  string
	columnName string
}

func (m *ColumnDoesNotExistError) Error() string {
	return "column with the name of " + m.columnName + " doesn't exist in table " + m.tableName
}

// InvalidNumberOfParametersError - error thrown when user provides invalid number of expected parameters
// (ex. fewer values in insert than defined )
type InvalidNumberOfParametersError struct {
	expectedNumber int
	actualNumber   int
	commandName    string
}

func (m *InvalidNumberOfParametersError) Error() string {
	return "invalid number of parameters in " + m.commandName + " command, should be: " + strconv.Itoa(m.expectedNumber) + ", but got: " + strconv.Itoa(m.actualNumber)
}

// InvalidValueTypeError - error thrown when user provides value of different type than expected
type InvalidValueTypeError struct {
	expectedType string
	actualType   string
	commandName  string
}

func (m *InvalidValueTypeError) Error() string {
	return "invalid value type provided in " + m.commandName + " command, expecting: " + m.expectedType + ", got: " + m.actualType
}

// UnsupportedValueType - error thrown when engine found unsupported data type to be stored inside
// the columns
type UnsupportedValueType struct {
	variable string
}

func (m *UnsupportedValueType) Error() string {
	return "couldn't map interface to any implementation of it: " + m.variable
}

// UnsupportedOperationTokenError - error thrown when engine found unsupported operation token
// (supported are: AND, OR)
type UnsupportedOperationTokenError struct {
	variable string
}

func (m *UnsupportedOperationTokenError) Error() string {
	return "unsupported operation token has been used: " + m.variable
}

// UnsupportedConditionalTokenError - error thrown when engine found unsupported conditional token
// inside expression (supported are: EQUAL, NOT)
type UnsupportedConditionalTokenError struct {
	variable    string
	commandName string
}

func (m *UnsupportedConditionalTokenError) Error() string {
	return "operation '" + m.variable + "' provided in " + m.commandName + " command isn't allowed"
}

// UnsupportedExpressionTypeError - error thrown when engine found unsupported expression type
type UnsupportedExpressionTypeError struct {
	variable    string
	commandName string
}

func (m *UnsupportedExpressionTypeError) Error() string {
	return "unsupported expression has been used in " + m.commandName + "command: " + m.variable
}

// UnsupportedCommandTypeFromParserError - error thrown when engine found unsupported command
// from parser
type UnsupportedCommandTypeFromParserError struct {
	variable string
}

func (m *UnsupportedCommandTypeFromParserError) Error() string {
	return "unsupported Command detected: " + m.variable
}
