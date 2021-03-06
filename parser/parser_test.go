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
		{"CREATE TABLE 	TBL( ONE TEXT );", "TBL", []string{"ONE"}, []token.Token{{token.TEXT, "TEXT"}}},
		{"CREATE TABLE 	TBL( ONE TEXT,  TWO TEXT, THREE INT);", "TBL", []string{"ONE", "TWO", "THREE"}, []token.Token{{token.TEXT, "TEXT"}, {token.TEXT, "TEXT"}, {token.INT, "INT"}}},
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
		{"INSERT INTO TBL VALUES( 'HELLO' );", "TBL", []token.Token{{token.IDENT, "HELLO"}}},
		{"INSERT INTO TBL VALUES( 'HELLO',	 10 , 'LOL');", "TBL", []token.Token{{token.IDENT, "HELLO"}, {token.LITERAL, "10"}, {token.IDENT, "LOL"}}},
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
		{"SELECT * FROM TBL;", "TBL", []token.Token{{token.ASTERISK, "*"}}},
		{"SELECT ONE, TWO, THREE FROM TBL;", "TBL", []token.Token{{token.IDENT, "ONE"}, {token.IDENT, "TWO"}, {token.IDENT, "THREE"}}},
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
