package lexer

import (
	"bytes"
	"github.com/LissaGreense/GO4SQL/token"
)

type Lexer struct {
	input             string
	position          int  // current position in input (points to current char)
	readPosition      int  // current reading position in input (after current char)
	character         byte // current char under examination
	insideApostrophes bool // flag which tells if lexer is between apostrophes
}

func RunLexer(input string) *Lexer {
	l := &Lexer{input: input, insideApostrophes: false}
	l.readChar()
	return l
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.character {
	case '*':
		tok = newToken(token.ASTERISK, string(lexer.character))
	case ';':
		tok = newToken(token.SEMICOLON, string(lexer.character))
	case ',':
		tok = newToken(token.COMMA, string(lexer.character))
	case '(':
		tok = newToken(token.LPAREN, string(lexer.character))
	case ')':
		tok = newToken(token.RPAREN, string(lexer.character))
	case '\'':
		lexer.insideApostrophes = !lexer.insideApostrophes
		tok = newToken(token.APOSTROPHE, string(lexer.character))
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if lexer.insideApostrophes {
			characters := lexer.processCharacters([]byte{'\''}, []byte{' ', '\n', '\t', '\r'})
			return characters
		} else {
			characters := lexer.processCharacters([]byte{'\'', '(', ',', ';', '*', ')'}, []byte{})
			return characters
		}
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) skipWhitespace() {
	for isWhitespace(lexer.character) && !lexer.insideApostrophes {
		lexer.readChar()
	}
}

func (lexer *Lexer) readChar() {
	lexer.character = lexer.getNextChar()
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) getNextChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.readPosition]
}

func (lexer *Lexer) processCharacters(blacklist []byte, whitelist []byte) token.Token {
	nextChar := lexer.getNextChar()
	position := lexer.position

	hasDigit := isDigit(lexer.character)
	hasLetter := isLetter(lexer.character)

	for validChar(nextChar, whitelist) && !bytes.ContainsAny(blacklist, string(nextChar)) {
		lexer.readChar()
		nextChar = lexer.getNextChar()

		hasDigit = hasDigit || isDigit(lexer.character)
		hasLetter = hasLetter || isLetter(lexer.character)
	}
	lexer.readChar()

	return lexer.evaluateToken(position, hasLetter, hasDigit)
}

func (lexer *Lexer) evaluateToken(position int, hasLetter bool, hasDigit bool) token.Token {
	characters := lexer.input[position:lexer.position]

	if (hasLetter && hasDigit) || lexer.insideApostrophes {
		return newToken(token.IDENT, characters)
	}
	if hasLetter && !hasDigit {
		return newToken(token.LookupIdent(characters), characters)
	}
	if !hasLetter && hasDigit {
		return newToken(token.LITERAL, characters)
	}

	return newToken(token.ILLEGAL, characters)
}

func validChar(nextChar byte, whitelist []byte) bool {
	return '!' <= nextChar && nextChar <= '~' || bytes.ContainsAny(whitelist, string(nextChar))
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func newToken(tokenType token.Type, characters string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: characters,
	}
}
