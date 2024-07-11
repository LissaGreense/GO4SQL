package engine

import (
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
	"github.com/LissaGreense/GO4SQL/token"
	"testing"
)

type errorHandlingTestSuite struct {
	input         string
	expectedError string
}

func TestEngineCreateCommandErrorHandling(t *testing.T) {
	duplicateTableNameError := TableAlreadyExistsError{"table1"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE table1( one TEXT , two INT);CREATE TABLE table1(two INT);", duplicateTableNameError.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineInsertCommandErrorHandling(t *testing.T) {
	tableDoNotExistError := TableDoesNotExistError{"table1"}
	invalidNumberOfParametersError := InvalidNumberOfParametersError{expectedNumber: 2, actualNumber: 1, commandName: token.INSERT}
	invalidParametersTypeError := InvalidValueTypeError{expectedType: token.IDENT, actualType: token.LITERAL, commandName: token.INSERT}
	tests := []errorHandlingTestSuite{
		{"INSERT INTO table1 VALUES( 'hello', 1);", tableDoNotExistError.Error()},
		{"CREATE TABLE table1( one TEXT , two INT); INSERT INTO table1 VALUES(1);", invalidNumberOfParametersError.Error()},
		{"CREATE TABLE table1( one TEXT , two INT); INSERT INTO table1 VALUES(1, 1 );", invalidParametersTypeError.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineSelectCommandErrorHandling(t *testing.T) {
	noTableDoesNotExist := TableDoesNotExistError{"tb1"}
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); SELECT * FROM tb1;", noTableDoesNotExist.Error()},
		{"CREATE TABLE tbl(one TEXT); SELECT two FROM tbl;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineDeleteCommandErrorHandling(t *testing.T) {
	noTableDoesNotExist := TableDoesNotExistError{"tb1"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); DELETE FROM tb1 WHERE one EQUAL 3;", noTableDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineWhereCommandErrorHandling(t *testing.T) {
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); INSERT INTO tbl VALUES('hello'); SELECT * FROM tbl WHERE two EQUAL 3;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineUpdateCommandErrorHandling(t *testing.T) {
	noTableDoesNotExist := TableDoesNotExistError{"tb1"}
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); UPDATE tb1 SET one TO 2;", noTableDoesNotExist.Error()},
		{"CREATE TABLE tbl(one TEXT);UPDATE tbl SET two TO 2;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineOrderByCommandErrorHandling(t *testing.T) {
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); SELECT * FROM tbl ORDER BY two ASC;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func runEngineErrorHandlingSuite(t *testing.T, suite []errorHandlingTestSuite) {
	for i, test := range suite {
		errorMsg := getErrorMessage(t, test.input, i)

		if errorMsg != test.expectedError {
			t.Fatalf("[%v]Was expecting error: \n\t{%s},\n\tbut it was:\n\t{%s}", i, test.expectedError, errorMsg)
		}
	}
}

func getErrorMessage(t *testing.T, input string, testIndex int) string {
	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
	sequences, parserError := parserInstance.ParseSequence()
	if parserError != nil {
		t.Fatalf("[%d] Error has occured in parser not in engine, error: %s", testIndex, parserError.Error())
	}

	engine := New()
	_, engineError := engine.Evaluate(sequences)
	if engineError == nil {
		t.Fatalf("[%d] Was expecting error from engine but there was none", testIndex)
	}

	return engineError.Error()
}
