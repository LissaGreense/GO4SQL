package greetings

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/token"
)

func TestAbc(t *testing.T) {
	input :=
		`
			create table 	tbl( one TEXT , two INT );
			insert into tbl values( 'hello',	 10 );
			insert 	into tbl  values( 'goodbye', 20 );
			`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "TBL"},
		{token.LPAREN, "("},
		{token.IDENT, "ONE"},
		{token.TEXT, "TEXT"},
		{token.COMMA, ","},
		{token.IDENT, "TWO"},
		{token.INT, "INT"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "TBL"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "HELLO"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "TBL"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "GOODBYE"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := lexer(input)

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
