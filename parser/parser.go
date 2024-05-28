package parser

import (
	"log"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/token"
)

// Parser - Contain token that is currently analyzed by parser and the next one. Lexer is used to tokenize the client
// text input.
type Parser struct {
	lexer        lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

// New - Return new Parser struct
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: *lexer}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken - Move pointer to the next token
func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

// validateTokenAndSkip - Check if current token type is appearing in provided expectedTokens array then move to the next token
func validateTokenAndSkip(parser *Parser, expectedTokens []token.Type) {
	validateToken(parser.currentToken.Type, expectedTokens)

	// Ignore validated token
	parser.nextToken()
}

// validateToken - Check if current token type is appearing in provided expectedTokens array
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

// parseCreateCommand - Return ast.CreateCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.CreateCommand:
// create table tbl( one TEXT , two INT );
func (parser *Parser) parseCreateCommand() ast.Command {
	// token.CREATE already at current position in parser
	createCommand := &ast.CreateCommand{Token: parser.currentToken}

	// Skip token.CREATE
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.TABLE})

	validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	createCommand.Name = ast.Identifier{Token: parser.currentToken}

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

func (parser *Parser) skipIfCurrentTokenIsApostrophe() {
	if parser.currentToken.Type == token.APOSTROPHE {
		parser.nextToken()
	}
}

func (parser *Parser) skipIfCurrentTokenIsSemicolon() {
	if parser.currentToken.Type == token.SEMICOLON {
		parser.nextToken()
	}
}

// parseInsertCommand - Return ast.InsertCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.InsertCommand:
// insert into tbl values( 'hello',	 10 );
func (parser *Parser) parseInsertCommand() ast.Command {
	// token.INSERT already at current position in parser
	insertCommand := &ast.InsertCommand{Token: parser.currentToken}

	// Ignore token.INSERT
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.INTO})

	validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	insertCommand.Name = ast.Identifier{Token: parser.currentToken}
	// Ignore token.INDENT
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.VALUES})
	validateTokenAndSkip(parser, []token.Type{token.LPAREN})

	for parser.currentToken.Type == token.IDENT || parser.currentToken.Type == token.LITERAL || parser.currentToken.Type == token.APOSTROPHE {
		parser.skipIfCurrentTokenIsApostrophe()

		validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.LITERAL})
		insertCommand.Values = append(insertCommand.Values, parser.currentToken)
		// Ignore token.IDENT or token.LITERAL
		parser.nextToken()

		parser.skipIfCurrentTokenIsApostrophe()

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

// parseSelectCommand - Return ast.SelectCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.SelectCommand:
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

	selectCommand.Name = ast.Identifier{Token: parser.currentToken}
	// Ignore token.INDENT
	parser.nextToken()

	// expect SEMICOLON or WHERE
	validateToken(parser.currentToken.Type, []token.Type{token.SEMICOLON, token.WHERE, token.ORDER})

	if parser.currentToken.Type == token.SEMICOLON {
		parser.nextToken()
	}

	return selectCommand
}

// parseWhereCommand - Return ast.WhereCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.WhereCommand:
// WHERE colName EQUAL 'potato'
func (parser *Parser) parseWhereCommand() ast.Command {
	// token.WHERE already at current position in parser
	whereCommand := &ast.WhereCommand{Token: parser.currentToken}
	expressionIsValid := false

	// Ignore token.WHERE
	parser.nextToken()

	expressionIsValid, whereCommand.Expression = parser.getExpression()

	if !expressionIsValid {
		log.Fatal("Expression withing Where statement couldn't be parsed correctly")
	}

	validateToken(parser.currentToken.Type, []token.Type{token.SEMICOLON, token.ORDER})

	parser.skipIfCurrentTokenIsSemicolon()

	return whereCommand
}

// parseDeleteCommand - Return ast.DeleteCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.DeleteCommand:
// DELETE FROM table;
func (parser *Parser) parseDeleteCommand() ast.Command {
	// token.DELETE already at current position in parser
	deleteCommand := &ast.DeleteCommand{Token: parser.currentToken}

	// token.DELETE no longer needed
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.FROM})

	validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	deleteCommand.Name = ast.Identifier{Token: parser.currentToken}

	// token.IDENT no longer needed
	parser.nextToken()

	// expect WHERE
	validateToken(parser.currentToken.Type, []token.Type{token.WHERE})

	return deleteCommand
}

// parseDropCommand - Return ast.DropCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.DropCommand:
// DROP TABLE table;
func (parser *Parser) parseDropCommand() ast.Command {
	// token.DROP already at current position in parser
	dropCommand := &ast.DropCommand{Token: parser.currentToken}

	// token.DROP no longer needed
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.TABLE})

	validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	dropCommand.Name = ast.Identifier{Token: parser.currentToken}

	// token.IDENT no longer needed
	parser.nextToken()

	parser.skipIfCurrentTokenIsSemicolon()

	return dropCommand
}

// parseOrderByCommand - Return ast.OrderByCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.OrderByCommand:
// ORDER BY colName ASC
func (parser *Parser) parseOrderByCommand() ast.Command {
	// token.ORDER already at current position in parser
	orderCommand := &ast.OrderByCommand{Token: parser.currentToken}

	// token.ORDER no longer needed
	parser.nextToken()

	validateTokenAndSkip(parser, []token.Type{token.BY})

	// ensure that loop below will execute at least once
	validateToken(parser.currentToken.Type, []token.Type{token.IDENT})

	// array of SortPattern
	for parser.currentToken.Type == token.IDENT {
		// Get column name
		validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
		columnName := parser.currentToken
		parser.nextToken()

		// Get ASC or DESC
		validateToken(parser.currentToken.Type, []token.Type{token.ASC, token.DESC})
		order := parser.currentToken
		parser.nextToken()

		// append sortPattern
		orderCommand.SortPatterns = append(orderCommand.SortPatterns, ast.SortPattern{ColumnName: columnName, Order: order})

		if parser.currentToken.Type != token.COMMA {
			break
		}
		// Ignore token.COMMA
		parser.nextToken()
	}

	validateTokenAndSkip(parser, []token.Type{token.SEMICOLON})

	return orderCommand
}

// getExpression - Return proper structure of ast.Expression and validate the syntax
//
// Available expressions:
// - ast.OperationExpression
// - ast.BooleanExpression
// - ast.ConditionExpression
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

// getOperationExpression - Return ast.OperationExpression created from tokens and validate the syntax
func (parser *Parser) getOperationExpression(booleanExpressionExists bool, conditionalExpressionExists bool, booleanExpression *ast.BooleanExpression, conditionalExpression *ast.ConditionExpression) (bool, *ast.OperationExpression) {
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

// getBooleanExpression - Return ast.BooleanExpression created from tokens and validate the syntax
func (parser *Parser) getBooleanExpression() (bool, *ast.BooleanExpression) {
	booleanExpression := &ast.BooleanExpression{}
	isValid := false

	if parser.currentToken.Type == token.TRUE || parser.currentToken.Type == token.FALSE {
		booleanExpression.Boolean = parser.currentToken
		parser.nextToken()
		isValid = true
	}

	return isValid, booleanExpression
}

// getConditionalExpression - Return ast.ConditionExpression created from tokens and validate the syntax
func (parser *Parser) getConditionalExpression() (bool, *ast.ConditionExpression) {
	conditionalExpression := &ast.ConditionExpression{}

	switch parser.currentToken.Type {
	case token.IDENT:
		conditionalExpression.Left = ast.Identifier{Token: parser.currentToken}
		parser.nextToken()
	case token.APOSTROPHE:
		parser.skipIfCurrentTokenIsApostrophe()
		conditionalExpression.Left = ast.Anonymitifier{Token: parser.currentToken}
		parser.nextToken()
		validateTokenAndSkip(parser, []token.Type{token.APOSTROPHE})
	case token.LITERAL:
		conditionalExpression.Left = ast.Anonymitifier{Token: parser.currentToken}
		parser.nextToken()
	default:
		return false, conditionalExpression
	}

	validateToken(parser.currentToken.Type, []token.Type{token.EQUAL, token.NOT})
	conditionalExpression.Condition = parser.currentToken
	parser.nextToken()

	switch parser.currentToken.Type {
	case token.IDENT:
		conditionalExpression.Right = ast.Identifier{Token: parser.currentToken}
		parser.nextToken()
	case token.APOSTROPHE:
		parser.skipIfCurrentTokenIsApostrophe()
		conditionalExpression.Right = ast.Anonymitifier{Token: parser.currentToken}
		parser.nextToken()
		validateTokenAndSkip(parser, []token.Type{token.APOSTROPHE})
	case token.LITERAL:
		conditionalExpression.Right = ast.Anonymitifier{Token: parser.currentToken}
		parser.nextToken()
	default:
		log.Fatal("Syntax error, expecting: ", token.APOSTROPHE, ",", token.IDENT, ",", token.LITERAL, ", got: ", parser.currentToken.Literal)
	}

	return true, conditionalExpression
}

// ParseSequence - Return ast.Sequence (sequence of commands) created from client input after tokenization
//
// Parse tokens returned by lexer to structures defines in ast package, and it's responsible for syntax validation.
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
		case token.DELETE:
			command = parser.parseDeleteCommand()
		case token.DROP:
			command = parser.parseDropCommand()
		case token.WHERE:
			lastCommand := parser.getLastCommand(sequence)

			if lastCommand.TokenLiteral() == token.SELECT {
				lastCommand.(*ast.SelectCommand).WhereCommand = parser.parseWhereCommand().(*ast.WhereCommand)
			} else if lastCommand.TokenLiteral() == token.DELETE {
				lastCommand.(*ast.DeleteCommand).WhereCommand = parser.parseWhereCommand().(*ast.WhereCommand)
			} else {
				log.Fatal("Syntax error, WHERE command needs SELECT or DELETE command before")
			}
		case token.ORDER:
			lastCommand := parser.getLastCommand(sequence)

			if lastCommand.TokenLiteral() != token.SELECT {
				log.Fatal("Syntax error, ORDER BY command needs SELECT command before")
			}

			selectCommand := lastCommand.(*ast.SelectCommand)
			selectCommand.OrderByCommand = parser.parseOrderByCommand().(*ast.OrderByCommand)
		default:
			log.Fatal("Syntax error, invalid command found: ", parser.currentToken.Type)
		}

		// Add command to the list of parsed commands
		if command != nil {
			sequence.Commands = append(sequence.Commands, command)
		}
	}

	return sequence
}

func (parser *Parser) getLastCommand(sequence *ast.Sequence) ast.Command {
	if len(sequence.Commands) == 0 {
		log.Fatal("Syntax error, Where Command can't be used without predecessor")
	}
	lastCommand := sequence.Commands[len(sequence.Commands)-1]
	return lastCommand
}
