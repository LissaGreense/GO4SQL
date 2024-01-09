package parser

import (
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

// New Return new Parser struct
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func validateTokenAndSkip(parser *Parser, expectedTokens []token.Type) {
	validateToken(parser.currentToken.Type, expectedTokens)

	// Ignore validated token
	parser.nextToken()
}

func validateToken(tokenType token.Type, expectedTokens []token.Type) {
	var contains = false
	var tokensPrintMessage = ""
	for i, x := range expectedTokens {

		if i == 0 {
			tokensPrintMessage += string(x)
		} else {
			tokensPrintMessage += ", or: " + string(x)
		}

		if x == tokenType {
			contains = true
			break
		}
	}
	if !contains {
		log.Fatal("Syntax error, expecting: ", tokensPrintMessage, ", got: ", tokenType)
	}
}

// create table tbl( one TEXT , two INT );
func (parser *Parser) parseCreateCommand() ast.Command { // TODO make it return the pointer
	// token.CREATE already at current position in parser
	createCommand := &ast.CreateCommand{Token: parser.currentToken}

	// Skip token.CREATE
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.TABLE})

	validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	createCommand.Name = &ast.Identifier{Token: parser.currentToken}

	// Skip token.IDENT
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.LPAREN})

	// Begin of inside Paren
	for parser.currentToken.Type == token.IDENT {
		validateToken(parser.peekToken.Type, []token.Type{token.TEXT, token.INT})
		createCommand.ColumnNames = append(createCommand.ColumnNames, parser.currentToken.Literal)
		createCommand.ColumnTypes = append(createCommand.ColumnTypes, parser.peekToken)

		// Skip token.IDENT
		parser.nextToken()
		// Skip token.TEXT or token.INT
		parser.nextToken()

		if parser.currentToken.Type != token.COMMA {
			break
		}

		// Skip token.COMMA
		parser.nextToken()
	}
	// End of inside Paren

	validateTokenAndSkip(parser, []token.Type{token.RPAREN})
	validateTokenAndSkip(parser, []token.Type{token.SEMICOLON})

	return createCommand
}

func (parser *Parser) skipApostrophe() {
	if parser.currentToken.Type == token.APOSTROPHE {
		parser.nextToken()
	}
}

// insert into tbl values( 'hello',	 10 );
func (parser *Parser) parseInsertCommand() ast.Command {
	// token.INSERT already at current position in parser
	insertCommand := &ast.InsertCommand{Token: parser.currentToken}

	// Ignore token.INSERT
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.INTO})

	validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	insertCommand.Name = &ast.Identifier{Token: parser.currentToken}
	// Ignore token.INDENT
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.VALUES})
	validateTokenAndSkip(parser, []token.Type{token.LPAREN})

	for parser.currentToken.Type == token.IDENT || parser.currentToken.Type == token.LITERAL || parser.currentToken.Type == token.APOSTROPHE {
		// TODO: Add apostrophe validation
		parser.skipApostrophe()

		validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.LITERAL})
		insertCommand.Values = append(insertCommand.Values, parser.currentToken)
		// Ignore token.IDENT or token.LITERAL
		parser.nextToken()
		parser.skipApostrophe()

		if parser.currentToken.Type != token.COMMA {
			break
		}

		// Ignore token.COMMA
		parser.nextToken()
	}

	validateTokenAndSkip(parser, []token.Type{token.RPAREN})
	validateTokenAndSkip(parser, []token.Type{token.SEMICOLON})
	return insertCommand
}

// SELECT col1, col2, col3 FROM tbl;
func (parser *Parser) parseSelectCommand() ast.Command {
	// token.SELECT already at current position in parser
	selectCommand := &ast.SelectCommand{Token: parser.currentToken}

	// Ignore token.SELECT
	parser.nextToken()

	if parser.currentToken.Type == token.ASTERISK {
		selectCommand.Space = append(selectCommand.Space, parser.currentToken)
		parser.nextToken()

	} else {
		for parser.currentToken.Type == token.IDENT {
			// Get column name
			validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
			selectCommand.Space = append(selectCommand.Space, parser.currentToken)
			parser.nextToken()

			if parser.currentToken.Type != token.COMMA {
				break
			}
			// Ignore token.COMMA
			parser.nextToken()
		}
	}

	validateTokenAndSkip(parser, []token.Type{token.FROM})

	selectCommand.Name = &ast.Identifier{Token: parser.currentToken}
	// Ignore token.INDENT
	parser.nextToken()

	// expect SEMICOLON or WHERE
	validateToken(parser.currentToken.Type, []token.Type{token.SEMICOLON, token.WHERE})

	if parser.currentToken.Type == token.SEMICOLON {
		parser.nextToken()
	}

	return selectCommand
}

// WHERE colName EQUAL 'potato'
func (parser *Parser) parseWhereCommand() ast.Command {
	// token.WHERE already at current position in parser
	whereCommand := &ast.WhereCommand{Token: parser.currentToken}
	expressionIsValid := false

	// Ignore token.WHERE
	parser.nextToken()

	expressionIsValid, whereCommand.Expression = parser.getExpression()

	if !expressionIsValid {
		log.Fatal("Expression withing Where statment couldn't be parsed correctly")
	}

	validateTokenAndSkip(parser, []token.Type{token.SEMICOLON})

	return whereCommand
}

func (parser *Parser) getExpression() (bool, ast.Expression) {
	booleanExpressionExists, booleanExpression := parser.getBooleanExpression()

	conditionalExpressionExists, conditionalExpression := parser.getConditionalExpression()

	operationExpressionExists, operationExpression := parser.getOperationExpression(booleanExpressionExists, conditionalExpressionExists, booleanExpression, conditionalExpression)

	if operationExpressionExists {
		return true, operationExpression
	}

	if conditionalExpressionExists {
		return true, conditionalExpression
	}

	if booleanExpressionExists {
		return true, booleanExpression
	}

	return false, nil
}

func (parser *Parser) getOperationExpression(booleanExpressionExists bool, conditionalExpressionExists bool, booleanExpression *ast.BooleanExpresion, conditionalExpression *ast.ConditionExpresion) (bool, *ast.OperationExpression) {
	operationExpression := &ast.OperationExpression{}

	if (booleanExpressionExists || conditionalExpressionExists) && (parser.currentToken.Type == token.OR || parser.currentToken.Type == token.AND) {
		if booleanExpressionExists {
			operationExpression.Left = booleanExpression
		}

		if conditionalExpressionExists {
			operationExpression.Left = conditionalExpression
		}

		operationExpression.Operation = parser.currentToken
		parser.nextToken()

		expressionIsValid, expression := parser.getExpression()

		if !expressionIsValid {
			log.Fatal("Couldn't parse right side of the OperationExpression after ", operationExpression.Operation.Literal, " token.")
		}

		operationExpression.Right = expression

		return true, operationExpression
	}

	return false, operationExpression
}

func (parser *Parser) getBooleanExpression() (bool, *ast.BooleanExpresion) {
	booleanExpression := &ast.BooleanExpresion{}
	isValid := false

	if parser.currentToken.Type == token.TRUE || parser.currentToken.Type == token.FALSE {
		booleanExpression.Boolean = parser.currentToken
		parser.nextToken()
		isValid = true
	}

	return isValid, booleanExpression
}

func (parser *Parser) getConditionalExpression() (bool, *ast.ConditionExpresion) {
	// TODO REFACTOR THIS
	conditionalExpression := &ast.ConditionExpresion{}

	if parser.currentToken.Type == token.IDENT {
		conditionalExpression.Left = ast.Identifier{
			Token: parser.currentToken,
		}
		parser.nextToken()

	} else if parser.currentToken.Type == token.APOSTROPHE {
		parser.skipApostrophe()

		conditionalExpression.Left = ast.Anonymitifier{
			Token: parser.currentToken,
		}

		parser.nextToken()

		validateTokenAndSkip(parser, []token.Type{token.APOSTROPHE})
	} else if parser.currentToken.Type == token.LITERAL {
		conditionalExpression.Left = ast.Anonymitifier{
			Token: parser.currentToken,
		}
		parser.nextToken()

	} else {
		return false, conditionalExpression
	}

	validateToken(parser.currentToken.Type, []token.Type{token.EQUAL, token.NOT})
	conditionalExpression.Condition = parser.currentToken
	parser.nextToken()

	if parser.currentToken.Type == token.IDENT {
		conditionalExpression.Right = ast.Identifier{
			Token: parser.currentToken,
		}
		parser.nextToken()

	} else if parser.currentToken.Type == token.APOSTROPHE {
		parser.skipApostrophe()

		conditionalExpression.Right = ast.Anonymitifier{
			Token: parser.currentToken,
		}

		parser.nextToken()

		validateTokenAndSkip(parser, []token.Type{token.APOSTROPHE})

	} else if parser.currentToken.Type == token.LITERAL {
		conditionalExpression.Right = ast.Anonymitifier{
			Token: parser.currentToken,
		}
		parser.nextToken()

	} else {
		log.Fatal("Syntax error, expecting: ", token.APOSTROPHE, ",", token.IDENT, ",", token.LITERAL, ", got: ", parser.currentToken.Literal)
	}

	return true, conditionalExpression
}

func (parser *Parser) ParseSequence() *ast.Sequence {
	// Create variable holding sequence/commands
	sequence := &ast.Sequence{}

	for parser.currentToken.Type != token.EOF {
		var command ast.Command
		switch parser.currentToken.Type {
		case token.CREATE:
			command = parser.parseCreateCommand()
		case token.INSERT:
			command = parser.parseInsertCommand()
		case token.SELECT:
			command = parser.parseSelectCommand()
		case token.WHERE:
			if len(sequence.Commands) == 0 || sequence.Commands[len(sequence.Commands)-1].TokenLiteral() != token.SELECT {
				log.Fatal("Syntax error, WHERE command needs SELECT command before")
			}
			command = parser.parseWhereCommand()
		default:
			log.Fatal("Syntax error, invalid command found")
		}

		// Add command to the list of parsed commands
		if command != nil {
			sequence.Commands = append(sequence.Commands, command)
		}
	}

	return sequence
}
