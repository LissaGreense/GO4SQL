package lexer

import (
	"strings"

	"github.com/LissaGreense/GO4SQL/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	character    byte // current char under examination
}

func RunLexer(input string) *Lexer {
	input = strings.ToUpper(input) // map everything to uppercase
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.character {
	case '*':
		tok = newToken(token.ASTERISK, lexer.character)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.character)
	case ',':
		tok = newToken(token.COMMA, lexer.character)
	case '(':
		tok = newToken(token.LPAREN, lexer.character)
	case ')':
		tok = newToken(token.RPAREN, lexer.character)
	case '\'':
		tok = newToken(token.APOSTROPHE, lexer.character)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lexer.character) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lexer.character) {
			tok.Type = token.LITERAL
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			// unsupported stuff
			tok = newToken(token.ILLEGAL, lexer.character)
		}
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.character == ' ' || lexer.character == '\t' || lexer.character == '\n' || lexer.character == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.character = 0
	} else {
		lexer.character = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.character) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.character) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
