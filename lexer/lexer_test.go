package lexer

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/token"
)

func TestLexer(t *testing.T) {
	input :=
		`
			CREATE TABLE 	tbl( one TEXT , two INT );
			INSERT INTO tbl VALUES( 'hello',	 10 );
			INSERT 	INTO tbl  VALUES( 'goodbye', 20 );
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "tbl"},
		{token.LPAREN, "("},
		{token.IDENT, "one"},
		{token.TEXT, "TEXT"},
		{token.COMMA, ","},
		{token.IDENT, "two"},
		{token.INT, "INT"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "hello"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "goodbye"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := RunLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerWithNumbersMixedInLitterals(t *testing.T) {
	input :=
		`
			CREATE TABLE 	tbl2( one TEXT , two INT );
			INSERT INTO tbl2 VALUES( 'hello1',	 10 );
			INSERT 	INTO tbl2  VALUES( 'good123bye', 20 );
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "tbl2"},
		{token.LPAREN, "("},
		{token.IDENT, "one"},
		{token.TEXT, "TEXT"},
		{token.COMMA, ","},
		{token.IDENT, "two"},
		{token.INT, "INT"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl2"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "hello1"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl2"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "good123bye"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := RunLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
