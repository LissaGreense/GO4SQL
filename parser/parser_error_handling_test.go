package parser

import (
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/token"
	"testing"
)

type errorHandlingTestSuite struct {
	input         string
	expectedError string
}

func TestParseCreateCommandErrorHandling(t *testing.T) {
	noTableKeyword := SyntaxError{[]string{token.TABLE}, token.IDENT}
	noTableName := SyntaxError{[]string{token.IDENT}, token.LPAREN}
	noLeftParen := SyntaxError{[]string{token.LPAREN}, token.IDENT}
	noRightParen := SyntaxError{[]string{token.RPAREN}, token.SEMICOLON}
	noColumDefinition := SyntaxError{[]string{token.IDENT}, token.RPAREN}
	noColumnName := SyntaxError{[]string{token.IDENT}, token.TEXT}
	noColumnType := SyntaxError{[]string{token.TEXT, token.INT}, token.COMMA}
	noSemicolon := SyntaxError{[]string{token.SEMICOLON}, ""}

	tests := []errorHandlingTestSuite{
		{"CREATE tbl(one TEXT);", noTableKeyword.Error()},
		{"CREATE TABLE (one TEXT);", noTableName.Error()},
		{"CREATE TABLE tbl one TEXT);", noLeftParen.Error()},
		{"CREATE TABLE tbl (one TEXT;", noRightParen.Error()},
		{"CREATE TABLE tbl ();", noColumDefinition.Error()},
		{"CREATE TABLE tbl (TEXT, two INT);", noColumnName.Error()},
		{"CREATE TABLE tbl (one , two INT);", noColumnType.Error()},
		{"CREATE TABLE tbl (one TEXT, two INT)", noSemicolon.Error()},
	}

	runParserErrorHandlingSuite(t, tests)

}

func TestParseDropCommandErrorHandling(t *testing.T) {
	missingTableKeywordError := SyntaxError{expecting: []string{token.TABLE}, got: token.IDENT}
	missingDropKeywordError := SyntaxInvalidCommandError{token.TABLE}
	missingSemicolonError := &SyntaxError{expecting: []string{token.SEMICOLON}, got: ""}
	invalidIdentError := &SyntaxError{expecting: []string{token.IDENT}, got: token.LITERAL}
	tests := []errorHandlingTestSuite{
		{input: "DROP table;", expectedError: missingTableKeywordError.Error()},
		{input: "TABLE table;", expectedError: missingDropKeywordError.Error()},
		{input: "DROP TABLE table", expectedError: missingSemicolonError.Error()},
		{input: "DROP TABLE 2;", expectedError: invalidIdentError.Error()},
	}

	runParserErrorHandlingSuite(t, tests)
}

func TestParseInsertCommandErrorHandling(t *testing.T) {
	noIntoKeyword := SyntaxError{[]string{token.INTO}, token.IDENT}
	noTableName := SyntaxError{[]string{token.IDENT}, token.VALUES}
	noLeftParen := SyntaxError{[]string{token.LPAREN}, token.APOSTROPHE}
	noValue := SyntaxError{[]string{token.IDENT, token.LITERAL}, token.APOSTROPHE}
	noRightParen := SyntaxError{[]string{token.RPAREN}, token.SEMICOLON}
	noSemicolon := SyntaxError{[]string{token.SEMICOLON}, ""}

	tests := []errorHandlingTestSuite{
		{"INSERT tbl VALUES( 'hello', 10);", noIntoKeyword.Error()},
		{"INSERT INTO VALUES( 'hello', 10);", noTableName.Error()},
		{"INSERT INTO tl VALUES 'hello', 10);", noLeftParen.Error()},
		{"INSERT INTO tl VALUES ('', 10);", noValue.Error()},
		{"INSERT INTO tl VALUES ('hello', 10;", noRightParen.Error()},
		{"INSERT INTO tl VALUES ('hello', 10)", noSemicolon.Error()},
	}

	runParserErrorHandlingSuite(t, tests)

}

func TestParseUpdateCommandErrorHandling(t *testing.T) {
	notableName := SyntaxError{expecting: []string{token.IDENT}, got: token.SEMICOLON}
	noSetKeyword := SyntaxError{expecting: []string{token.SET}, got: token.SEMICOLON}
	noColumnName := SyntaxError{expecting: []string{token.IDENT}, got: token.LITERAL}
	noToKeyword := SyntaxError{expecting: []string{token.TO}, got: token.SEMICOLON}
	noSecondIdentOrLiteralForValue := SyntaxError{expecting: []string{token.IDENT, token.LITERAL}, got: token.SEMICOLON}
	//noCommaBetweenValues := SyntaxInvalidCommandError{"column_name_2"}
	noWhereOrSemicolon := SyntaxError{expecting: []string{token.SEMICOLON, token.WHERE}, got: token.SELECT}

	// UPDATE table SET col1 TO 'value' WHERE col2 EQUAL 10;
	tests := []errorHandlingTestSuite{
		{"UPDATE;", notableName.Error()},
		{"UPDATE table;", noSetKeyword.Error()},
		{"UPDATE table SET 2;", noColumnName.Error()},
		{"UPDATE table SET column_name_1;", noToKeyword.Error()},
		{"UPDATE table SET column_name_1 TO;", noSecondIdentOrLiteralForValue.Error()},
		//TODO: ADD no comma handling or not ???
		//{"UPDATE table SET column_name_1 TO new_value_1 column_name_2 TO new_value_2;", noCommaBetweenValues.Error()},
		{"UPDATE table SET column_name_1 TO 'new_value_1' SELECT;", noWhereOrSemicolon.Error()},
	}

	runParserErrorHandlingSuite(t, tests)

}

func TestParseSelectCommandErrorHandling(t *testing.T) {
	noFromKeyword := SyntaxError{[]string{token.FROM}, token.IDENT}
	noColumns := SyntaxError{[]string{token.ASTERISK, token.IDENT}, token.FROM}
	noTableName := SyntaxError{[]string{token.IDENT}, token.SEMICOLON}
	noSemicolon := SyntaxError{[]string{token.SEMICOLON, token.WHERE, token.ORDER, token.LIMIT, token.OFFSET}, ""}

	tests := []errorHandlingTestSuite{
		{"SELECT column1, column2 tbl;", noFromKeyword.Error()},
		{"SELECT FROM table;", noColumns.Error()},
		{"SELECT column1, column2 FROM ;", noTableName.Error()},
		{"SELECT column1, column2 FROM table", noSemicolon.Error()},
	}

	runParserErrorHandlingSuite(t, tests)
}

func TestParseDeleteCommandErrorHandling(t *testing.T) {
	noFromKeyword := SyntaxError{[]string{token.FROM}, token.IDENT}
	noTableName := SyntaxError{[]string{token.IDENT}, token.WHERE}
	noWhereCommand := SyntaxError{[]string{token.WHERE}, ";"}

	tests := []errorHandlingTestSuite{
		{"DELETE table WHERE TRUE", noFromKeyword.Error()},
		{"DELETE FROM WHERE TRUE;", noTableName.Error()},
		{"DELETE FROM table;", noWhereCommand.Error()},
	}

	runParserErrorHandlingSuite(t, tests)
}

func runParserErrorHandlingSuite(t *testing.T, suite []errorHandlingTestSuite) {
	for i, test := range suite {
		errorMsg := getErrorMessage(t, test.input, i)

		if errorMsg != test.expectedError {
			t.Fatalf("[%v]Was expecting error: \n\t{%s},\n\tbut it was:\n\t{%s}", i, test.expectedError, errorMsg)
		}
	}
}

func getErrorMessage(t *testing.T, input string, testIndex int) string {
	lexerInstance := lexer.RunLexer(input)
	parserInstance := New(lexerInstance)
	_, err := parserInstance.ParseSequence()

	if err == nil {
		t.Fatalf("[%v]Was expecting error from parser but there was none", testIndex)
	}

	return err.Error()
}
