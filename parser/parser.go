package parser

import (
	"fmt"
	"log"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

// Return new Parser struct
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func syntaxError(expected token.TokenType, actual token.TokenType) {
	log.Fatal("Syntax error, expecting: ", expected, ", got: ", actual)
}

func syntaxError2(expected token.TokenType, secondExpected token.TokenType, actual token.TokenType) {
	log.Fatal("Syntax error, expecting: ", expected, ", or: ", secondExpected, "got: ", actual)
}

// create table tbl( one TEXT , two INT );
func (parser *Parser) parseCreateCommand() ast.Command { // TODO make it return the pointer
	createCommand := &ast.CreateCommand{Token: parser.currentToken}

	// Ignore token.CREATE
	parser.nextToken()

	// Check if next token is table
	if parser.currentToken.Type != token.TABLE {
		syntaxError(token.TABLE, parser.currentToken.Type)
	}

	// Ignore token.TABLE
	parser.nextToken()

	// Check if next token is IDENT (variable name)
	if parser.currentToken.Type != token.IDENT {
		syntaxError(token.IDENT, parser.currentToken.Type)
	}
	createCommand.Name = &ast.Identifier{Token: parser.currentToken}

	// Ignore token.IDENT
	parser.nextToken()

	// Check if next token is LPAREN
	if parser.currentToken.Type != token.LPAREN {
		syntaxError(token.LPAREN, parser.currentToken.Type)
	}

	// Ignore token.LPAREN
	parser.nextToken()

	// Begin od insisde of Paren
	for parser.currentToken.Type == token.IDENT {
		if parser.peekToken.Type != token.TEXT && parser.peekToken.Type != token.INT {
			syntaxError2(token.TEXT, token.INT, parser.peekToken.Type)
		}
		createCommand.ColumnNames = append(createCommand.ColumnNames, parser.currentToken.Literal)
		createCommand.ColumnTypes = append(createCommand.ColumnTypes, parser.peekToken)

		// Ignore token.IDENT
		parser.nextToken()
		// Ignore token.TEXT or token.INT
		parser.nextToken()

		if parser.currentToken.Type != token.COMMA {
			break
		}

		// Ignore token.COMMA
		parser.nextToken()
	}
	// End of insisde of Paren

	// Check if next token is RPAREN
	if parser.currentToken.Type != token.RPAREN {
		syntaxError(token.RPAREN, parser.currentToken.Type)
	}

	parser.nextToken()

	// Check if next token is SEMICOLON
	if parser.currentToken.Type != token.SEMICOLON {
		syntaxError(token.SEMICOLON, parser.currentToken.Type)
	}
	// Ignore token.SEMICOLON
	parser.nextToken()

	return createCommand
}

func (parser *Parser) ParseSequence() *ast.Sequence {

	// Create variable holding sequnce/commands
	sequence := &ast.Sequence{}

	for parser.currentToken.Type != token.EOF {
		var statement ast.Command
		switch parser.currentToken.Type {
		case token.CREATE:
			statement = parser.parseCreateCommand()
		case token.INSERT:
			fmt.Println("INSERT")
		case token.SELECT:
			fmt.Println("SELECT")
		default:
			fmt.Println("Syntax error")
		}

		// Add statment to the list of parsed commands
		if statement != nil {
			sequence.Commands = append(sequence.Commands, statement)
		}
	}

	return sequence
}
