package parser

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/token"
)

func TestParserCreateCommand(t *testing.T) {
	tests := []struct {
		input               string
		expectedTableName   string
		expectedColumnNames []string
		expectedColumTypes  []token.Token
	}{
		{"CREATE TABLE 	TBL( ONE TEXT );", "TBL", []string{"ONE"}, []token.Token{{Type: token.TEXT, Literal: "TEXT"}}},
		{"CREATE TABLE 	TBL( ONE TEXT,  TWO TEXT, THREE INT);", "TBL", []string{"ONE", "TWO", "THREE"}, []token.Token{{Type: token.TEXT, Literal: "TEXT"}, {Type: token.TEXT, Literal: "TEXT"}, {Type: token.INT, Literal: "INT"}}},
		{"CREATE TABLE 	TBL(  );", "TBL", []string{}, []token.Token{}},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 1 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !testCreateStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedColumnNames, tt.expectedColumTypes) {
			return
		}
	}
}

func testCreateStatement(t *testing.T, command ast.Command, expectedTableName string, expectedColumnNames []string, expectedColumTypes []token.Token) bool {
	if command.TokenLiteral() != "CREATE" {
		t.Errorf("command.TokenLiteral() not 'CREATE'. got=%q", command.TokenLiteral())
		return false
	}

	actualCreateCommand, ok := command.(*ast.CreateCommand)
	if !ok {
		t.Errorf("actualCreateCommand is not %T. got=%T", &ast.CreateCommand{}, command)
		return false
	}

	if actualCreateCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualCreateCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !stringArrayEquals(actualCreateCommand.ColumnNames, expectedColumnNames) {
		t.Errorf("")
		return false
	}

	if !tokenArrayEquals(actualCreateCommand.ColumnTypes, expectedColumTypes) {
		t.Errorf("")
		return false
	}

	return true
}

func TestParseInsertCommand(t *testing.T) {
	tests := []struct {
		input                string
		expectedTableName    string
		expectedValuesTokens []token.Token
	}{
		{"INSERT INTO TBL VALUES();", "TBL", []token.Token{}},
		{"INSERT INTO TBL VALUES( 'HELLO' );", "TBL", []token.Token{{Type: token.IDENT, Literal: "HELLO"}}},
		{"INSERT INTO TBL VALUES( 'HELLO',	 10 , 'LOL');", "TBL", []token.Token{{Type: token.IDENT, Literal: "HELLO"}, {Type: token.LITERAL, Literal: "10"}, {Type: token.IDENT, Literal: "LOL"}}},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 1 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !testInsertStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedValuesTokens) {
			return
		}
	}
}

func testInsertStatement(t *testing.T, command ast.Command, expectedTableName string, expectedValuesTokens []token.Token) bool {
	if command.TokenLiteral() != "INSERT" {
		t.Errorf("command.TokenLiteral() not 'INSERT'. got=%q", command.TokenLiteral())
		return false
	}

	actualInsertCommand, ok := command.(*ast.InsertCommand)
	if !ok {
		t.Errorf("actualInsertCommand is not %T. got=%T", &ast.InsertCommand{}, command)
		return false
	}

	if actualInsertCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualInsertCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !tokenArrayEquals(actualInsertCommand.Values, expectedValuesTokens) {
		t.Errorf("")
		return false
	}

	return true
}

func TestParseSelectCommand(t *testing.T) {
	tests := []struct {
		input             string
		expectedTableName string
		expectedColumns   []token.Token
	}{
		{"SELECT * FROM TBL;", "TBL", []token.Token{{Type: token.ASTERISK, Literal: "*"}}},
		{"SELECT ONE, TWO, THREE FROM TBL;", "TBL", []token.Token{{Type: token.IDENT, Literal: "ONE"}, {Type: token.IDENT, Literal: "TWO"}, {Type: token.IDENT, Literal: "THREE"}}},
		{"SELECT FROM TBL;", "TBL", []token.Token{}},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 1 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !testSelectStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedColumns) {
			return
		}
	}
}

func TestParseWhereCommand(t *testing.T) {
	firstExpression := ast.ConditionExpression{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "fda"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	secondExpression := ast.ConditionExpression{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	tests := []struct {
		input              string
		expectedExpression ast.Expression
	}{
		{
			input:              "SELECT * FROM TBL WHERE colName1 EQUAL 'fda';",
			expectedExpression: firstExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName2 EQUAL 6462389;",
			expectedExpression: secondExpression,
		},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 2 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !whereStatementIsValid(t, sequences.Commands[1], tt.expectedExpression) {
			return
		}
	}
}

func TestParseDeleteCommand(t *testing.T) {
	input := "DELETE FROM colName1 WHERE colName2 EQUAL 6462389;"
	expectedDeleteCommand := ast.DeleteCommand{
		Token: token.Token{Type: token.DELETE, Literal: "DELETE"},
		Name:  &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
	}
	expectedWhereCommand := ast.ConditionExpression{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences := parserInstance.ParseSequence()

	if len(sequences.Commands) != 2 {
		t.Fatalf("sequences does not contain 2 statements. got=%d", len(sequences.Commands))
	}

	actualDeleteCommand, ok := sequences.Commands[0].(*ast.DeleteCommand)
	if !ok {
		t.Errorf("actualDeleteCommand is not %T. got=%T", &ast.DeleteCommand{}, sequences.Commands[0])
	}

	if expectedDeleteCommand.TokenLiteral() != actualDeleteCommand.TokenLiteral() {
		t.Errorf("TokenLiteral of DeleteCommand is not %s. got=%s", expectedDeleteCommand.TokenLiteral(), actualDeleteCommand.TokenLiteral())
	}

	if expectedDeleteCommand.Name.GetToken().Literal != actualDeleteCommand.Name.GetToken().Literal {
		t.Errorf("Table name of DeleteCommand is not %s. got=%s", expectedDeleteCommand.Name.GetToken().Literal, actualDeleteCommand.Name.GetToken().Literal)
	}

	if !whereStatementIsValid(t, sequences.Commands[1], expectedWhereCommand) {
		return
	}
}

func TestSelectWithOrderByCommand(t *testing.T) {
	input := "SELECT * FROM tableName ORDER BY colName1 DESC;"
	expectedSortPattern := ast.SortPattern{
		ColumnName: token.Token{Type: token.IDENT, Literal: "colName1"},
		Order:      token.Token{Type: token.DESC, Literal: "DESC"},
	}
	expectedOrderByCommand := ast.OrderByCommand{
		Token:        token.Token{Type: token.ORDER, Literal: "ORDER"},
		SortPatterns: []ast.SortPattern{expectedSortPattern},
	}
	expectedTableName := "tableName"
	expectedColumnName := []token.Token{{Type: token.ASTERISK, Literal: "*"}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences := parserInstance.ParseSequence()

	if len(sequences.Commands) != 2 {
		t.Fatalf("sequences does not contain 2 statements. got=%d", len(sequences.Commands))
	}

	if !testSelectStatement(t, sequences.Commands[0], expectedTableName, expectedColumnName) {
		return
	}

	actualOrderByCommand, orderByCommandIsOk := sequences.Commands[1].(*ast.OrderByCommand)
	if !orderByCommandIsOk {
		t.Errorf("actualDeleteCommand is not %T. got=%T", &ast.OrderByCommand{}, sequences.Commands[0])
	}

	testOrderByCommands(t, expectedOrderByCommand, actualOrderByCommand)
}

func TestParseLogicOperatorsInCommand(t *testing.T) {

	firstExpression := ast.OperationExpression{
		Left: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "fda"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}},
		Right: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "123"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "NOT"}},
		Operation: token.Token{Type: token.AND, Literal: "AND"},
	}

	secondExpression := ast.OperationExpression{
		Left: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
			Condition: token.Token{Type: token.NOT, Literal: "NOT"}},
		Right: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "qwe"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}},
		Operation: token.Token{Type: token.OR, Literal: "OR"},
	}

	thirdExpression := ast.BooleanExpression{
		Boolean: token.Token{Type: token.TRUE, Literal: "TRUE"},
	}

	tests := []struct {
		input              string
		expectedExpression ast.Expression
	}{
		{
			input:              "SELECT * FROM TBL WHERE colName1 EQUAL 'fda' AND colName2 NOT 123;",
			expectedExpression: firstExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName2 NOT 6462389 OR colName1 EQUAL 'qwe';",
			expectedExpression: secondExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE TRUE;",
			expectedExpression: thirdExpression,
		},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 2 {
			t.Fatalf("sequences does not contain 2 statements. got=%d", len(sequences.Commands))
		}

		if !whereStatementIsValid(t, sequences.Commands[1], tt.expectedExpression) {
			t.Fatalf("Actual expression and expected one are different")
		}
	}
}

func testSelectStatement(t *testing.T, command ast.Command, expectedTableName string, expectedColumnsTokens []token.Token) bool {
	if command.TokenLiteral() != "SELECT" {
		t.Errorf("command.TokenLiteral() not 'SELECT'. got=%q", command.TokenLiteral())
		return false
	}

	actualSelectCommand, ok := command.(*ast.SelectCommand)
	if !ok {
		t.Errorf("actualSelectCommand is not %T. got=%T", &ast.SelectCommand{}, command)
		return false
	}

	if actualSelectCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualSelectCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !tokenArrayEquals(actualSelectCommand.Space, expectedColumnsTokens) {
		t.Errorf("")
		return false
	}

	return true
}

func whereStatementIsValid(t *testing.T, command ast.Command, expectedExpression ast.Expression) bool {
	if command.TokenLiteral() != "WHERE" {
		t.Errorf("command.TokenLiteral() not 'WHERE'. got=%q", command.TokenLiteral())
		return false
	}

	actualWhereCommand, ok := command.(*ast.WhereCommand)
	if !ok {
		t.Errorf("actualWhereCommand is not %T. got=%T", &ast.WhereCommand{}, command)
		return false
	}

	if !expressionsAreEqual(actualWhereCommand.Expression, expectedExpression) {
		t.Errorf("Actual expression is not equal to expected one.\nActual: %#v\nExpected: %#v", actualWhereCommand.Expression, expectedExpression)
		return false
	}

	return true
}

func stringArrayEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func tokenArrayEquals(a []token.Token, b []token.Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Literal != b[i].Literal {
			return false
		}
	}
	return true
}

func testOrderByCommands(t *testing.T, expectedOrderByCommand ast.OrderByCommand, actualOrderByCommand *ast.OrderByCommand) {

	if expectedOrderByCommand.Token.Type != actualOrderByCommand.Token.Type {
		t.Errorf("Expecting Token TokenType: %q, got: %q", expectedOrderByCommand.Token.Type, actualOrderByCommand.Token.Type)
	}
	if expectedOrderByCommand.Token.Literal != actualOrderByCommand.Token.Literal {
		t.Errorf("Expecting Token Literal: %s, got: %s", expectedOrderByCommand.Token.Literal, actualOrderByCommand.Token.Literal)
	}
	if len(expectedOrderByCommand.SortPatterns) != len(actualOrderByCommand.SortPatterns) {
		t.Errorf("Expecting Sorting Pattern Array to have: %d elements, got: %d", len(expectedOrderByCommand.SortPatterns), len(actualOrderByCommand.SortPatterns))
	}

	for i, expectedSortPattern := range expectedOrderByCommand.SortPatterns {
		if expectedSortPattern.Order.Literal != actualOrderByCommand.SortPatterns[i].Order.Literal {
			t.Errorf("Expecting Order: %s, got: %s", expectedSortPattern.Order.Literal, actualOrderByCommand.SortPatterns[i].Order.Literal)
		}
		if expectedSortPattern.ColumnName.Literal != expectedOrderByCommand.SortPatterns[i].ColumnName.Literal {
			t.Errorf("Expecting Column Name: %s, got: %s", expectedSortPattern.ColumnName.Literal, actualOrderByCommand.SortPatterns[i].ColumnName.Literal)
		}
	}

}

func expressionsAreEqual(first ast.Expression, second ast.Expression) bool {

	booleanExpression, booleanExpressionIsValid := first.(*ast.BooleanExpression)
	if booleanExpressionIsValid {
		return validateBooleanExpressions(second, booleanExpression)
	}

	conditionExpression, conditionExpressionIsValid := first.(*ast.ConditionExpression)
	if conditionExpressionIsValid {
		return validateConditionExpression(second, conditionExpression)
	}

	operationExpression, operationExpressionIsValid := first.(*ast.OperationExpression)
	if operationExpressionIsValid {
		return validateOperationExpression(second, operationExpression)
	}

	return false
}

func validateOperationExpression(second ast.Expression, operationExpression *ast.OperationExpression) bool {
	secondOperationExpression, secondOperationExpressionIsValid := second.(ast.OperationExpression)

	if !secondOperationExpressionIsValid {
		return false
	}

	if operationExpression.Operation.Literal != secondOperationExpression.Operation.Literal {
		return false
	}

	return expressionsAreEqual(operationExpression.Left, secondOperationExpression.Left) && expressionsAreEqual(operationExpression.Right, secondOperationExpression.Right)
}

func validateConditionExpression(second ast.Expression, conditionExpression *ast.ConditionExpression) bool {
	secondConditionExpression, secondConditionExpressionIsValid := second.(ast.ConditionExpression)

	if !secondConditionExpressionIsValid {
		return false
	}

	if conditionExpression.Left.GetToken().Literal != secondConditionExpression.Left.GetToken().Literal &&
		conditionExpression.Left.IsIdentifier() == secondConditionExpression.Left.IsIdentifier() {
		return false
	}

	if conditionExpression.Right.GetToken().Literal != secondConditionExpression.Right.GetToken().Literal &&
		conditionExpression.Right.IsIdentifier() == secondConditionExpression.Right.IsIdentifier() {
		return false
	}

	if conditionExpression.Condition.Literal != secondConditionExpression.Condition.Literal {
		return false
	}

	return true
}

func validateBooleanExpressions(second ast.Expression, booleanExpression *ast.BooleanExpression) bool {
	secondBooleanExpresion, secondBooleanExpresionIsValid := second.(ast.BooleanExpression)

	if !secondBooleanExpresionIsValid {
		return false
	}

	if booleanExpression.Boolean.Literal != secondBooleanExpresion.Boolean.Literal {
		return false
	}

	return true
}
