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
	noColumnName := SyntaxError{[]string{token.RPAREN}, token.TEXT}
	noColumnType := SyntaxError{[]string{token.TEXT, token.INT}, token.COMMA}
	noSemicolon := SyntaxError{[]string{token.SEMICOLON}, ""}

	tests := []errorHandlingTestSuite{
		{"CREATE tbl(one TEXT);", noTableKeyword.Error()},
		{"CREATE TABLE (one TEXT);", noTableName.Error()},
		{"CREATE TABLE tbl one TEXT);", noLeftParen.Error()},
		{"CREATE TABLE tbl (one TEXT;", noRightParen.Error()},
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
	noValue := SyntaxError{[]string{token.IDENT, token.LITERAL, token.NULL}, token.APOSTROPHE}
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
	noSecondIdentOrLiteralForValue := SyntaxError{expecting: []string{token.IDENT, token.LITERAL, token.NULL}, got: token.SEMICOLON}
	noCommaBetweenValues := SyntaxError{expecting: []string{token.SEMICOLON, token.WHERE}, got: token.IDENT}
	noWhereOrSemicolon := SyntaxError{expecting: []string{token.SEMICOLON, token.WHERE}, got: token.SELECT}

	tests := []errorHandlingTestSuite{
		{"UPDATE;", notableName.Error()},
		{"UPDATE table;", noSetKeyword.Error()},
		{"UPDATE table SET 2;", noColumnName.Error()},
		{"UPDATE table SET column_name_1;", noToKeyword.Error()},
		{"UPDATE table SET column_name_1 TO;", noSecondIdentOrLiteralForValue.Error()},
		{"UPDATE table SET column_name_1 TO 2 column_name_1 TO 3;", noCommaBetweenValues.Error()},
		{"UPDATE table SET column_name_1 TO 'new_value_1' SELECT;", noWhereOrSemicolon.Error()},
	}

	runParserErrorHandlingSuite(t, tests)

}

func TestParseSelectCommandErrorHandling(t *testing.T) {
	noFromKeyword := SyntaxError{[]string{token.FROM}, token.IDENT}
	noColumns := SyntaxError{[]string{token.ASTERISK, token.IDENT, token.MAX, token.MIN, token.SUM, token.AVG, token.COUNT}, token.FROM}
	noTableName := SyntaxError{[]string{token.IDENT}, token.SEMICOLON}
	noSemicolon := SyntaxError{[]string{token.SEMICOLON, token.WHERE, token.ORDER, token.LIMIT, token.OFFSET, token.JOIN, token.LEFT, token.RIGHT, token.INNER, token.FULL}, ""}
	noAggregateFunctionParenClosure := SyntaxError{[]string{token.RPAREN}, ","}
	noAggregateFunctionLeftParen := SyntaxError{[]string{token.LPAREN}, token.IDENT}
	noFromAfterAsterisk := SyntaxError{[]string{token.FROM}, ","}
	noAsteriskInsideMaxArgument := SyntaxError{[]string{token.IDENT}, "*"}

	tests := []errorHandlingTestSuite{
		{"SELECT column1, column2 tbl;", noFromKeyword.Error()},
		{"SELECT FROM table;", noColumns.Error()},
		{"SELECT column1, column2 FROM ;", noTableName.Error()},
		{"SELECT column1, column2 FROM table", noSemicolon.Error()},
		{"SELECT SUM(column1, column2 FROM table", noAggregateFunctionParenClosure.Error()},
		{"SELECT SUM column1 FROM table", noAggregateFunctionLeftParen.Error()},
		{"SELECT *, colName FROM table", noFromAfterAsterisk.Error()},
		{"SELECT MAX(*) FROM table", noAsteriskInsideMaxArgument.Error()},
	}

	runParserErrorHandlingSuite(t, tests)
}

func TestParseWhereCommandErrorHandling(t *testing.T) {
	selectCommandPrefix := "SELECT * FROM tbl "
	noPredecessorError := NoPredecessorParserError{command: token.WHERE}
	noColName := LogicalExpressionParsingError{}
	noOperatorInsideWhereStatementException := LogicalExpressionParsingError{}
	valueIsMissing := SyntaxError{expecting: []string{token.APOSTROPHE, token.IDENT, token.LITERAL, token.NULL}, got: token.SEMICOLON}
	tokenAnd := token.AND
	conjunctionIsMissing := SyntaxError{expecting: []string{token.SEMICOLON, token.ORDER}, got: token.IDENT}
	nextLogicalExpressionIsMissing := LogicalExpressionParsingError{afterToken: &tokenAnd}
	noSemicolon := SyntaxError{expecting: []string{token.SEMICOLON, token.ORDER}, got: ""}
	noLeftParGotSemicolon := SyntaxError{expecting: []string{token.LPAREN}, got: ";"}
	noLeftParGotNumber := SyntaxError{expecting: []string{token.LPAREN}, got: token.LITERAL}
	noComma := SyntaxError{expecting: []string{token.COMMA, token.RPAREN}, got: token.LITERAL}
	noInKeywordException := LogicalExpressionParsingError{}

	tests := []errorHandlingTestSuite{
		{"WHERE col1 NOT 'goodbye' OR col2 EQUAL 3;", noPredecessorError.Error()},
		{selectCommandPrefix + "WHERE NOT 'goodbye' OR column2 EQUAL 3;", noColName.Error()},
		{selectCommandPrefix + "WHERE one 'goodbye';", noOperatorInsideWhereStatementException.Error()},
		{selectCommandPrefix + "WHERE one EQUAL;", valueIsMissing.Error()},
		{selectCommandPrefix + "WHERE one EQUAL 5 two NOT 1;", conjunctionIsMissing.Error()},
		{selectCommandPrefix + "WHERE one EQUAL 5 AND;", nextLogicalExpressionIsMissing.Error()},
		{selectCommandPrefix + "WHERE one EQUAL 5 AND two NOT 5", noSemicolon.Error()},
		{selectCommandPrefix + "WHERE one IN ;", noLeftParGotSemicolon.Error()},
		{selectCommandPrefix + "WHERE one IN 5;", noLeftParGotNumber.Error()},
		{selectCommandPrefix + "WHERE one IN (5 6);", noComma.Error()},
		{selectCommandPrefix + "WHERE one (5, 6);", noInKeywordException.Error()},
	}

	runParserErrorHandlingSuite(t, tests)

}

func TestParseOrderByCommandErrorHandling(t *testing.T) {
	selectCommandPrefix := "SELECT * FROM tbl "
	noPredecessorError := NoPredecessorParserError{command: token.ORDER}
	noAscDescError := SyntaxError{expecting: []string{token.ASC, token.DESC}, got: token.SEMICOLON}
	noByKeywordError := SyntaxError{expecting: []string{token.BY}, got: token.IDENT}
	noIdentKeywordError := SyntaxError{expecting: []string{token.IDENT}, got: token.ASC}

	tests := []errorHandlingTestSuite{
		{"ORDER BY column1;", noPredecessorError.Error()},
		{selectCommandPrefix + "ORDER BY column1;", noAscDescError.Error()},
		{selectCommandPrefix + "ORDER  column1 ASC;", noByKeywordError.Error()},
		{selectCommandPrefix + "ORDER BY ASC;", noIdentKeywordError.Error()},
	}

	runParserErrorHandlingSuite(t, tests)
}

func TestParseLimitCommandErrorHandling(t *testing.T) {
	selectCommandPrefix := "SELECT * FROM tbl "
	noPredecessorError := NoPredecessorParserError{command: token.LIMIT}
	noLiteralError := SyntaxError{expecting: []string{token.LITERAL}, got: token.SEMICOLON}
	lessThanZeroError := ArithmeticLessThanZeroParserError{variable: "limit"}

	tests := []errorHandlingTestSuite{
		{"LIMIT 5;", noPredecessorError.Error()},
		{selectCommandPrefix + "LIMIT;", noLiteralError.Error()},
		{selectCommandPrefix + "LIMIT -10;", lessThanZeroError.Error()},
	}

	runParserErrorHandlingSuite(t, tests)
}

func TestParseOffsetCommandErrorHandling(t *testing.T) {
	selectCommandPrefix := "SELECT * FROM tbl "
	noPredecessorError := NoPredecessorParserError{command: token.OFFSET}
	noLiteralError := SyntaxError{expecting: []string{token.LITERAL}, got: token.IDENT}
	lessThanZeroError := ArithmeticLessThanZeroParserError{variable: "offset"}

	tests := []errorHandlingTestSuite{
		{"OFFSET 5;", noPredecessorError.Error()},
		{selectCommandPrefix + "OFFSET hi;", noLiteralError.Error()},
		{selectCommandPrefix + "OFFSET -10;", lessThanZeroError.Error()},
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

func TestPeriodInIdentWhileCreatingTableErrorHandling(t *testing.T) {
	illegalPeriodInTableName := IllegalPeriodInIdentParserError{"tab.le"}
	illegalPeriodInColumnName := IllegalPeriodInIdentParserError{"col.umn"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tab.le( one TEXT , two INT);", illegalPeriodInTableName.Error()},
		{"CREATE TABLE table1( col.umn TEXT , two INT);", illegalPeriodInColumnName.Error()},
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
