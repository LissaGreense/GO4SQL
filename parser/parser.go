package parser

import (
	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/token"
	"strconv"
	"strings"
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
func validateTokenAndSkip(parser *Parser, expectedTokens []token.Type) error {
	err := validateToken(parser.currentToken.Type, expectedTokens)

	if err != nil {
		return err
	}

	// Ignore validated token
	parser.nextToken()
	return nil
}

// validateToken - Check if current token type is appearing in provided expectedTokens array
func validateToken(tokenType token.Type, expectedTokens []token.Type) error {
	var contains = false
	expectedTokensStrings := make([]string, 0)
	for _, x := range expectedTokens {
		expectedTokensStrings = append(expectedTokensStrings, string(x))

		if x == tokenType {
			contains = true
			break
		}
	}
	if !contains {
		return &SyntaxError{expectedTokensStrings, string(tokenType)}
	}
	return nil
}

// parseCreateCommand - Return ast.CreateCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.CreateCommand:
// create table tbl( one TEXT , two INT );
func (parser *Parser) parseCreateCommand() (ast.Command, error) {
	// token.CREATE already at current position in parser
	createCommand := &ast.CreateCommand{Token: parser.currentToken}

	// Skip token.CREATE
	parser.nextToken()

	err := validateTokenAndSkip(parser, []token.Type{token.TABLE})
	if err != nil {
		return nil, err
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}

	if strings.Contains(parser.currentToken.Literal, ".") {
		return nil, &IllegalPeriodInIdentParserError{name: parser.currentToken.Literal}
	}
	createCommand.Name = ast.Identifier{Token: parser.currentToken}

	// Skip token.IDENT
	parser.nextToken()

	err = validateTokenAndSkip(parser, []token.Type{token.LPAREN})
	if err != nil {
		return nil, err
	}

	// Begin of inside Paren
	for parser.currentToken.Type == token.IDENT {
		err = validateToken(parser.peekToken.Type, []token.Type{token.TEXT, token.INT})
		if err != nil {
			return nil, err
		}

		if strings.Contains(parser.currentToken.Literal, ".") {
			return nil, &IllegalPeriodInIdentParserError{name: parser.currentToken.Literal}
		}

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

	err = validateTokenAndSkip(parser, []token.Type{token.RPAREN})
	if err != nil {
		return nil, err
	}
	err = validateTokenAndSkip(parser, []token.Type{token.SEMICOLON})
	if err != nil {
		return nil, err
	}

	return createCommand, nil
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
func (parser *Parser) parseInsertCommand() (ast.Command, error) {
	// token.INSERT already at current position in parser
	insertCommand := &ast.InsertCommand{Token: parser.currentToken}

	// Ignore token.INSERT
	parser.nextToken()

	err := validateTokenAndSkip(parser, []token.Type{token.INTO})
	if err != nil {
		return nil, err
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}
	insertCommand.Name = ast.Identifier{Token: parser.currentToken}
	// Ignore token.INDENT
	parser.nextToken()

	err = validateTokenAndSkip(parser, []token.Type{token.VALUES})
	if err != nil {
		return nil, err
	}
	err = validateTokenAndSkip(parser, []token.Type{token.LPAREN})
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type == token.IDENT || parser.currentToken.Type == token.LITERAL || parser.currentToken.Type == token.APOSTROPHE {
		parser.skipIfCurrentTokenIsApostrophe()

		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.LITERAL})
		if err != nil {
			return nil, err
		}
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

	err = validateTokenAndSkip(parser, []token.Type{token.RPAREN})
	if err != nil {
		return nil, err
	}

	err = validateTokenAndSkip(parser, []token.Type{token.SEMICOLON})
	if err != nil {
		return nil, err
	}

	return insertCommand, nil
}

// parseSelectCommand - Return ast.SelectCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.SelectCommand:
// SELECT col1, col2, col3 FROM tbl;
func (parser *Parser) parseSelectCommand() (ast.Command, error) {
	// token.SELECT already at current position in parser
	selectCommand := &ast.SelectCommand{Token: parser.currentToken}

	// Ignore token.SELECT
	parser.nextToken()

	// optional DISTINCT
	if parser.currentToken.Type == token.DISTINCT {
		selectCommand.HasDistinct = true

		// Ignore token.DISTINCT
		parser.nextToken()
	}

	err := validateToken(parser.currentToken.Type, []token.Type{token.ASTERISK, token.IDENT, token.MAX, token.MIN, token.SUM, token.AVG, token.COUNT})
	if err != nil {
		return nil, err
	}

	if parser.currentToken.Type == token.ASTERISK {
		selectCommand.Space = append(selectCommand.Space, ast.Space{ColumnName: parser.currentToken})
		parser.nextToken()
	} else {
		for parser.currentToken.Type == token.IDENT || isAggregateFunction(parser.currentToken.Type) {
			if parser.currentToken.Type != token.IDENT {
				aggregateFunction := parser.currentToken
				parser.nextToken()
				err := validateTokenAndSkip(parser, []token.Type{token.LPAREN})
				if err != nil {
					return nil, err
				}
				// Get column name
				err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.ASTERISK})
				if err != nil {
					return nil, err
				}
				selectCommand.Space = append(selectCommand.Space, ast.Space{ColumnName: parser.currentToken, AggregateFunc: &aggregateFunction})
				parser.nextToken()

				err = validateTokenAndSkip(parser, []token.Type{token.RPAREN})
				if err != nil {
					return nil, err
				}
			} else {
				// Get column name
				err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
				if err != nil {
					return nil, err
				}
				selectCommand.Space = append(selectCommand.Space, ast.Space{ColumnName: parser.currentToken})
				parser.nextToken()
			}

			if parser.currentToken.Type != token.COMMA {
				break
			}
			// Ignore token.COMMA
			parser.nextToken()
		}
	}

	err = validateTokenAndSkip(parser, []token.Type{token.FROM})
	if err != nil {
		return nil, err
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}

	selectCommand.Name = ast.Identifier{Token: parser.currentToken}
	// Ignore token.IDENT
	parser.nextToken()

	// expect SEMICOLON or other keywords expected in SELECT statement
	err = validateToken(parser.currentToken.Type, []token.Type{token.SEMICOLON, token.WHERE, token.ORDER, token.LIMIT, token.OFFSET, token.JOIN, token.LEFT, token.RIGHT, token.INNER, token.FULL})
	if err != nil {
		return nil, err
	}

	if parser.currentToken.Type == token.SEMICOLON {
		parser.nextToken()
	}

	return selectCommand, nil
}

func (parser *Parser) getColumnName(err error, selectCommand *ast.SelectCommand, aggregateFunction token.Token) error {
	// Get column name
	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.ASTERISK})
	if err != nil {
		return err
	}
	selectCommand.Space = append(selectCommand.Space, ast.Space{ColumnName: parser.currentToken, AggregateFunc: &aggregateFunction})
	parser.nextToken()
	return nil
}

func isAggregateFunction(t token.Type) bool {
	return t == token.MIN || t == token.MAX || t == token.COUNT || t == token.SUM || t == token.AVG
}

// parseWhereCommand - Return ast.WhereCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.WhereCommand:
// WHERE colName EQUAL 'potato'
func (parser *Parser) parseWhereCommand() (ast.Command, error) {
	// token.WHERE already at current position in parser
	whereCommand := &ast.WhereCommand{Token: parser.currentToken}
	expressionIsValid := false

	// Ignore token.WHERE
	parser.nextToken()
	var err error
	expressionIsValid, whereCommand.Expression, err = parser.getExpression()
	if err != nil {
		return nil, err
	}

	if !expressionIsValid {
		return nil, &LogicalExpressionParsingError{}
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.SEMICOLON, token.ORDER})
	if err != nil {
		return nil, err
	}

	parser.skipIfCurrentTokenIsSemicolon()

	return whereCommand, nil
}

// parseDeleteCommand - Return ast.DeleteCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.DeleteCommand:
// DELETE FROM table;
func (parser *Parser) parseDeleteCommand() (ast.Command, error) {
	// token.DELETE already at current position in parser
	deleteCommand := &ast.DeleteCommand{Token: parser.currentToken}

	// token.DELETE no longer needed
	parser.nextToken()

	err := validateTokenAndSkip(parser, []token.Type{token.FROM})
	if err != nil {
		return nil, err
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}
	deleteCommand.Name = ast.Identifier{Token: parser.currentToken}

	// token.IDENT no longer needed
	parser.nextToken()

	// expect WHERE
	err = validateToken(parser.currentToken.Type, []token.Type{token.WHERE})

	return deleteCommand, err
}

// parseDropCommand - Return ast.DropCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.DropCommand:
// DROP TABLE table;
func (parser *Parser) parseDropCommand() (ast.Command, error) {
	// token.DROP already at current position in parser
	dropCommand := &ast.DropCommand{Token: parser.currentToken}

	// token.DROP no longer needed
	parser.nextToken()

	err := validateTokenAndSkip(parser, []token.Type{token.TABLE})
	if err != nil {
		return nil, err
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}
	dropCommand.Name = ast.Identifier{Token: parser.currentToken}

	// token.IDENT no longer needed
	parser.nextToken()

	err = validateTokenAndSkip(parser, []token.Type{token.SEMICOLON})

	return dropCommand, err
}

// parseOrderByCommand - Return ast.OrderByCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.OrderByCommand:
// ORDER BY colName ASC
func (parser *Parser) parseOrderByCommand() (ast.Command, error) {
	// token.ORDER already at current position in parser
	orderCommand := &ast.OrderByCommand{Token: parser.currentToken}

	// token.ORDER no longer needed
	parser.nextToken()

	err := validateTokenAndSkip(parser, []token.Type{token.BY})
	if err != nil {
		return nil, err
	}

	// ensure that loop below will execute at least once
	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}

	// array of SortPattern
	for parser.currentToken.Type == token.IDENT {
		// Get column name
		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
		if err != nil {
			return nil, err
		}
		columnName := parser.currentToken
		parser.nextToken()

		// Get ASC or DESC
		err = validateToken(parser.currentToken.Type, []token.Type{token.ASC, token.DESC})
		if err != nil {
			return nil, err
		}
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

	parser.skipIfCurrentTokenIsSemicolon()

	return orderCommand, nil
}

// parseLimitCommand - Return ast.LimitCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.LimitCommand:
// LIMIT 10
func (parser *Parser) parseLimitCommand() (ast.Command, error) {
	// token.LIMIT already at current position in parser
	limitCommand := &ast.LimitCommand{Token: parser.currentToken}

	// token.LIMIT no longer needed
	parser.nextToken()

	err := validateToken(parser.currentToken.Type, []token.Type{token.LITERAL})
	if err != nil {
		return nil, err
	}

	// convert count number to int
	count, err := strconv.Atoi(parser.currentToken.Literal)
	if err != nil {
		return nil, err
	}

	if count < 0 {
		return nil, &ArithmeticLessThanZeroParserError{variable: "limit"}
	}

	limitCommand.Count = count

	// Skip token.IDENT
	parser.nextToken()

	parser.skipIfCurrentTokenIsSemicolon()

	return limitCommand, nil
}

// parseOffsetCommand - Return ast.OffsetCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.LimitCommand:
// OFFSET 10
func (parser *Parser) parseOffsetCommand() (ast.Command, error) {
	// token.OFFSET already at current position in parser
	offsetCommand := &ast.OffsetCommand{Token: parser.currentToken}

	// token.OFFSET no longer needed
	parser.nextToken()

	err := validateToken(parser.currentToken.Type, []token.Type{token.LITERAL})
	if err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(parser.currentToken.Literal)
	if err != nil {
		return nil, err
	}
	if count < 0 {
		return nil, &ArithmeticLessThanZeroParserError{variable: "offset"}
	}

	offsetCommand.Count = count

	// Skip token.IDENT
	parser.nextToken()

	parser.skipIfCurrentTokenIsSemicolon()

	return offsetCommand, nil
}

// parseJoinCommand - Return ast.JoinCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.JoinCommand:
// JOIN table on table.one EQUAL table2.one;
func (parser *Parser) parseJoinCommand() (ast.Command, error) {
	// parser has either token.JOIN, token.LEFT, token.RIGHT, token.INNER or token.FULL
	var joinCommand *ast.JoinCommand

	if parser.currentToken.Type == token.JOIN {
		joinCommand = &ast.JoinCommand{Token: parser.currentToken}
		joinCommand.JoinType = token.Token{Type: token.INNER, Literal: token.INNER}
	} else {
		joinTypeTokenType := parser.currentToken
		parser.nextToken()
		err := validateToken(parser.currentToken.Type, []token.Type{token.JOIN})
		if err != nil {
			return nil, err
		}
		joinCommand = &ast.JoinCommand{Token: parser.currentToken}
		joinCommand.JoinType = joinTypeTokenType
	}

	// token.JOIN no longer needed
	parser.nextToken()

	err := validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}

	joinCommand.Name = ast.Identifier{Token: parser.currentToken}
	parser.nextToken()

	err = validateTokenAndSkip(parser, []token.Type{token.ON})
	if err != nil {
		return nil, err
	}

	var expressionIsValid bool
	expressionIsValid, joinCommand.Expression, err = parser.getExpression()
	if err != nil {
		return nil, err
	}

	if !expressionIsValid {
		return nil, &LogicalExpressionParsingError{}
	}

	parser.skipIfCurrentTokenIsSemicolon()

	return joinCommand, nil
}

// parseUpdateCommand - Return ast.parseUpdateCommand created from tokens and validate the syntax
//
// Example of input parsable to the ast.parseUpdateCommand:
// UPDATE table SET col1 TO 'value' WHERE col2 EQUAL 10;
func (parser *Parser) parseUpdateCommand() (ast.Command, error) {
	// token.UPDATE already at current position in parser
	updateCommand := &ast.UpdateCommand{Token: parser.currentToken}

	// Ignore token.UPDATE
	parser.nextToken()

	// Get table name
	err := validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}
	updateCommand.Name = ast.Identifier{Token: parser.currentToken}

	// Ignore token.IDENT
	parser.nextToken()

	err = validateTokenAndSkip(parser, []token.Type{token.SET})
	if err != nil {
		return nil, err
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	if err != nil {
		return nil, err
	}

	updateCommand.Changes = make(map[token.Token]ast.Anonymitifier)
	for parser.currentToken.Type == token.IDENT {
		// Get column name
		err := validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
		if err != nil {
			return nil, err
		}
		colKey := parser.currentToken

		// skip column name
		parser.nextToken()

		err = validateToken(parser.currentToken.Type, []token.Type{token.TO})
		if err != nil {
			return nil, err
		}
		// skip token.TO
		parser.nextToken()

		parser.skipIfCurrentTokenIsApostrophe()
		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.LITERAL})
		if err != nil {
			return nil, err
		}
		updateCommand.Changes[colKey] = ast.Anonymitifier{Token: parser.currentToken}

		// skip token.IDENT or token.LITERAL
		parser.nextToken()
		parser.skipIfCurrentTokenIsApostrophe()

		if parser.currentToken.Type != token.COMMA {
			break
		}

		// Skip token.COMMA
		parser.nextToken()
	}

	err = validateToken(parser.currentToken.Type, []token.Type{token.SEMICOLON, token.WHERE})
	if err != nil {
		return nil, err
	}
	parser.skipIfCurrentTokenIsSemicolon()
	return updateCommand, nil
}

// getExpression - Return proper structure of ast.Expression and validate the syntax
//
// Available expressions:
// - ast.OperationExpression
// - ast.BooleanExpression
// - ast.ConditionExpression
func (parser *Parser) getExpression() (bool, ast.Expression, error) {
	booleanExpressionExists, booleanExpression := parser.getBooleanExpression()

	conditionalExpressionExists, conditionalExpression, err := parser.getConditionalExpression()
	if err != nil {
		return false, nil, err
	}

	operationExpressionExists, operationExpression, err := parser.getOperationExpression(booleanExpressionExists, conditionalExpressionExists, booleanExpression, conditionalExpression)
	if err != nil {
		return false, nil, err
	}

	if operationExpressionExists {
		return true, operationExpression, err
	}

	if conditionalExpressionExists {
		return true, conditionalExpression, err
	}

	if booleanExpressionExists {
		return true, booleanExpression, err
	}

	return false, nil, err
}

// getOperationExpression - Return ast.OperationExpression created from tokens and validate the syntax
func (parser *Parser) getOperationExpression(booleanExpressionExists bool, conditionalExpressionExists bool, booleanExpression *ast.BooleanExpression, conditionalExpression *ast.ConditionExpression) (bool, *ast.OperationExpression, error) {
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

		expressionIsValid, expression, err := parser.getExpression()

		if err != nil {
			return false, nil, err
		}
		if !expressionIsValid {
			return false, nil, &LogicalExpressionParsingError{afterToken: &operationExpression.Operation.Literal}
		}

		operationExpression.Right = expression

		return true, operationExpression, nil
	}

	return false, operationExpression, nil
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
func (parser *Parser) getConditionalExpression() (bool, *ast.ConditionExpression, error) {
	conditionalExpression := &ast.ConditionExpression{}

	switch parser.currentToken.Type {
	case token.IDENT:
		conditionalExpression.Left = ast.Identifier{Token: parser.currentToken}
		parser.nextToken()
	case token.APOSTROPHE:
		parser.skipIfCurrentTokenIsApostrophe()
		conditionalExpression.Left = ast.Anonymitifier{Token: parser.currentToken}
		parser.nextToken()
		err := validateTokenAndSkip(parser, []token.Type{token.APOSTROPHE})
		if err != nil {
			return false, nil, err
		}
	case token.LITERAL:
		conditionalExpression.Left = ast.Anonymitifier{Token: parser.currentToken}
		parser.nextToken()
	default:
		return false, conditionalExpression, nil
	}

	err := validateToken(parser.currentToken.Type, []token.Type{token.EQUAL, token.NOT})
	if err != nil {
		return false, nil, err
	}
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
		err := validateTokenAndSkip(parser, []token.Type{token.APOSTROPHE})
		if err != nil {
			return false, nil, err
		}
	case token.LITERAL:
		conditionalExpression.Right = ast.Anonymitifier{Token: parser.currentToken}
		parser.nextToken()
	default:
		return false, nil, &SyntaxError{expecting: []string{token.APOSTROPHE, token.IDENT, token.LITERAL}, got: parser.currentToken.Literal}
	}

	return true, conditionalExpression, nil
}

// ParseSequence - Return ast.Sequence (sequence of commands) created from client input after tokenization
//
// Parse tokens returned by lexer to structures defines in ast package, and it's responsible for syntax validation.
func (parser *Parser) ParseSequence() (*ast.Sequence, error) {
	// Create variable holding sequence/commands
	sequence := &ast.Sequence{}

	for parser.currentToken.Type != token.EOF {
		var command ast.Command
		var err error
		switch parser.currentToken.Type {
		case token.CREATE:
			command, err = parser.parseCreateCommand()
		case token.INSERT:
			command, err = parser.parseInsertCommand()
		case token.UPDATE:
			command, err = parser.parseUpdateCommand()
		case token.SELECT:
			command, err = parser.parseSelectCommand()
		case token.DELETE:
			command, err = parser.parseDeleteCommand()
		case token.DROP:
			command, err = parser.parseDropCommand()
		case token.WHERE:
			lastCommand, parserError := parser.getLastCommand(sequence, token.WHERE)
			if parserError != nil {
				return nil, parserError
			}

			if lastCommand.TokenLiteral() == token.SELECT {
				newCommand, err := parser.parseWhereCommand()
				if err != nil {
					return nil, err
				}
				lastCommand.(*ast.SelectCommand).WhereCommand = newCommand.(*ast.WhereCommand)
			} else if lastCommand.TokenLiteral() == token.DELETE {
				newCommand, err := parser.parseWhereCommand()
				if err != nil {
					return nil, err
				}
				lastCommand.(*ast.DeleteCommand).WhereCommand = newCommand.(*ast.WhereCommand)
			} else if lastCommand.TokenLiteral() == token.UPDATE {
				newCommand, err := parser.parseWhereCommand()
				if err != nil {
					return nil, err
				}
				lastCommand.(*ast.UpdateCommand).WhereCommand = newCommand.(*ast.WhereCommand)
			} else {
				return nil, &SyntaxCommandExpectedError{command: "WHERE", neededCommands: []string{"SELECT", "DELETE", "UPDATE"}}
			}
		case token.ORDER:
			lastCommand, parserError := parser.getLastCommand(sequence, token.ORDER)
			if parserError != nil {
				return nil, parserError
			}

			if lastCommand.TokenLiteral() != token.SELECT {
				return nil, &SyntaxCommandExpectedError{command: "ORDER BY", neededCommands: []string{"SELECT"}}
			}

			selectCommand := lastCommand.(*ast.SelectCommand)
			newCommand, err := parser.parseOrderByCommand()
			if err != nil {
				return nil, err
			}
			selectCommand.OrderByCommand = newCommand.(*ast.OrderByCommand)
		case token.LIMIT:
			lastCommand, parserError := parser.getLastCommand(sequence, token.LIMIT)
			if parserError != nil {
				return nil, parserError
			}
			if lastCommand.TokenLiteral() != token.SELECT {
				return nil, &SyntaxCommandExpectedError{command: "LIMIT", neededCommands: []string{"SELECT"}}
			}
			selectCommand := lastCommand.(*ast.SelectCommand)
			newCommand, err := parser.parseLimitCommand()
			if err != nil {
				return nil, err
			}
			selectCommand.LimitCommand = newCommand.(*ast.LimitCommand)
		case token.OFFSET:
			lastCommand, parserError := parser.getLastCommand(sequence, token.OFFSET)
			if parserError != nil {
				return nil, parserError
			}
			if lastCommand.TokenLiteral() != token.SELECT {
				return nil, &SyntaxCommandExpectedError{command: "OFFSET", neededCommands: []string{"SELECT"}}
			}
			selectCommand := lastCommand.(*ast.SelectCommand)
			newCommand, err := parser.parseOffsetCommand()
			if err != nil {
				return nil, err
			}
			selectCommand.OffsetCommand = newCommand.(*ast.OffsetCommand)
		case token.JOIN, token.LEFT, token.RIGHT, token.INNER, token.FULL:
			lastCommand, parserError := parser.getLastCommand(sequence, token.JOIN)
			if parserError != nil {
				return nil, parserError
			}
			if lastCommand.TokenLiteral() != token.SELECT {
				return nil, &SyntaxCommandExpectedError{command: "JOIN", neededCommands: []string{"SELECT"}}
			}
			selectCommand := lastCommand.(*ast.SelectCommand)
			newCommand, err := parser.parseJoinCommand()
			if err != nil {
				return nil, err
			}
			selectCommand.JoinCommand = newCommand.(*ast.JoinCommand)
		default:
			return nil, &SyntaxInvalidCommandError{invalidCommand: parser.currentToken.Literal}
		}

		if err != nil {
			return nil, err
		}

		// Add command to the list of parsed commands
		if command != nil {
			sequence.Commands = append(sequence.Commands, command)
		}
	}

	return sequence, nil
}

func (parser *Parser) getLastCommand(sequence *ast.Sequence, currentToken string) (ast.Command, error) {
	if len(sequence.Commands) == 0 {
		return nil, &NoPredecessorParserError{command: currentToken}
	}
	lastCommand := sequence.Commands[len(sequence.Commands)-1]
	return lastCommand, nil
}
