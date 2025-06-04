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

func (parser *Parser) skipIfCurrentTokenIsApostrophe() bool {
	if parser.currentToken.Type == token.APOSTROPHE {
		parser.nextToken()
		return true
	}
	return false
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

	for parser.currentToken.Type == token.IDENT || parser.currentToken.Type == token.LITERAL || parser.currentToken.Type == token.NULL || parser.currentToken.Type == token.APOSTROPHE {
		startedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()

		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.LITERAL, token.NULL})
		if err != nil {
			return nil, err
		}
		value := parser.currentToken
		insertCommand.Values = append(insertCommand.Values, value)
		// Ignore token.IDENT, token.LITERAL or token.NULL
		parser.nextToken()

		finishedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()

		err = validateApostropheWrapping(startedWithApostrophe, finishedWithApostrophe, value)

		if err != nil {
			return nil, err
		}

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

func validateApostropheWrapping(startedWithApostrophe bool, finishedWithApostrophe bool, value token.Token) error {
	if startedWithApostrophe && !finishedWithApostrophe {
		return &NoApostropheOnRightParserError{ident: value.Literal}
	} else if !startedWithApostrophe && finishedWithApostrophe {
		return &NoApostropheOnLeftParserError{ident: value.Literal}
	}
	return nil
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
		command, err := parser.parseSelectSpace(selectCommand)
		if err != nil {
			return command, err
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

func (parser *Parser) parseSelectSpace(selectCommand *ast.SelectCommand) (ast.Command, error) {
	for parser.currentToken.Type == token.IDENT || isAggregateFunction(parser.currentToken.Type) {
		if parser.currentToken.Type != token.IDENT {
			err := parser.parseAggregateFunction(selectCommand)
			if err != nil {
				return nil, err
			}
		} else {
			// Get column name
			err := validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
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
	return nil, nil
}

func (parser *Parser) parseAggregateFunction(selectCommand *ast.SelectCommand) error {
	aggregateFunction := parser.currentToken
	parser.nextToken()
	err := validateTokenAndSkip(parser, []token.Type{token.LPAREN})
	if err != nil {
		return err
	}
	if aggregateFunction.Type == token.COUNT {
		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.ASTERISK})
	} else {
		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT})
	}
	if err != nil {
		return err
	}
	selectCommand.Space = append(selectCommand.Space, ast.Space{ColumnName: parser.currentToken, AggregateFunc: &aggregateFunction})
	parser.nextToken()

	err = validateTokenAndSkip(parser, []token.Type{token.RPAREN})
	return err
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

	// Ignore token.WHERE
	parser.nextToken()
	var err error
	whereCommand.Expression, err = parser.getExpression()
	if err != nil {
		return nil, err
	}

	if whereCommand.Expression == nil {
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

	joinCommand.Expression, err = parser.getExpression()
	if err != nil {
		return nil, err
	}

	if joinCommand.Expression == nil {
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

		startedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()
		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.LITERAL, token.NULL})
		if err != nil {
			return nil, err
		}
		updateCommand.Changes[colKey] = ast.Anonymitifier{Token: parser.currentToken}

		// skip token.IDENT, token.LITERAL or token.NULL
		parser.nextToken()
		finishedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()

		err = validateApostropheWrapping(startedWithApostrophe, finishedWithApostrophe, updateCommand.Changes[colKey].GetToken())

		if err != nil {
			return nil, err
		}

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

// getExpression is the entry point for parsing logical expressions.
// It handles OR operations (lowest precedence).
func (parser *Parser) getExpression() (ast.Expression, error) {
	leftExpr, err := parser.getTerm()
	if err != nil {
		return nil, err // Propagate error from getTerm
	}
	if leftExpr == nil { // If getTerm returned (nil,nil) convert to error
		return nil, &LogicalExpressionParsingError{customMessage: "Missing left operand for OR expression", afterToken: &parser.currentToken.Literal}
	}

	for parser.currentToken.Type == token.OR {
		operatorToken := parser.currentToken // Store OR token
		parser.nextToken()                   // Consume OR

		rightExpr, err := parser.getTerm()
		if err != nil {
			return nil, err // Propagate error from getTerm (for RHS)
		}
		if rightExpr == nil { // If getTerm for RHS returned (nil,nil) convert to error
			return nil, &LogicalExpressionParsingError{customMessage: "Missing right operand for OR expression", afterToken: &operatorToken.Literal}
		}

		leftExpr = &ast.OperationExpression{Left: leftExpr, Operation: operatorToken, Right: rightExpr}
	}

	return leftExpr, nil
}

// getTerm handles AND operations (medium precedence).
func (parser *Parser) getTerm() (ast.Expression, error) {
	leftExpr, err := parser.getFactor()
	if err != nil {
		return nil, err // Propagate error from getFactor
	}
	if leftExpr == nil { // If getFactor returned (nil,nil) convert to error
		return nil, &LogicalExpressionParsingError{customMessage: "Missing left operand for AND expression", afterToken: &parser.currentToken.Literal}
	}

	for parser.currentToken.Type == token.AND {
		operatorToken := parser.currentToken // Store AND token
		parser.nextToken()                   // Consume AND

		rightExpr, err := parser.getFactor()
		if err != nil {
			return nil, err // Propagate error from getFactor (for RHS)
		}
		if rightExpr == nil { // If getFactor for RHS returned (nil,nil) convert to error
			return nil, &LogicalExpressionParsingError{customMessage: "Missing right operand for AND expression", afterToken: &operatorToken.Literal}
		}

		leftExpr = &ast.OperationExpression{Left: leftExpr, Operation: operatorToken, Right: rightExpr}
	}

	return leftExpr, nil
}

// getFactor handles individual conditions, literals, and parenthesized expressions (highest precedence).
func (parser *Parser) getFactor() (ast.Expression, error) {
	if parser.currentToken.Type == token.LPAREN {
		parser.nextToken() // Consume LPAREN
		expr, err := parser.getExpression() // Recursive call to parse the inner expression
		if err != nil {
			return nil, err
		}
		if expr == nil { // If inner expression is (nil,nil) convert to error
			return nil, &LogicalExpressionParsingError{customMessage: "Empty or invalid parenthesized expression"}
		}
		if parser.currentToken.Type != token.RPAREN {
			return nil, &SyntaxError{expecting: []string{string(token.RPAREN)}, got: string(parser.currentToken.Type)}
		}
		parser.nextToken() // Consume RPAREN
		return expr, nil
	}

	// Otherwise, parse a simple condition or boolean
	leftSideValue, isAnonymitifier, err := parser.getExpressionLeftSideValue()
	if err != nil { // Error from getExpressionLeftSideValue (e.g. unexpected token like AND when factor expected)
		return nil, err
	}
	// If getExpressionLeftSideValue succeeded, leftSideValue is a valid token.

	// Check for specific condition types based on the *next* token
	// Current token is the operator (e.g. EQUAL, ISNULL, IN)
	if parser.currentToken.Type == token.EQUAL ||
		parser.currentToken.Type == token.NOT || // Assuming NOT here implies a condition like "NOT EQUAL" or part of "IS NOT NULL"
		parser.currentToken.Type == token.ISNULL ||
		parser.currentToken.Type == token.ISNOTNULL {
		expr, err := parser.getConditionalExpression(leftSideValue, isAnonymitifier)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, &LogicalExpressionParsingError{customMessage: "Conditional expression parsing failed", afterToken: &leftSideValue.Literal}
		}
		return expr, nil
	} else if parser.currentToken.Type == token.IN || parser.currentToken.Type == token.NOTIN {
		expr, err := parser.getContainExpression(leftSideValue, isAnonymitifier)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, &LogicalExpressionParsingError{customMessage: "Contain expression parsing failed", afterToken: &leftSideValue.Literal}
		}
		return expr, nil
	} else if leftSideValue.Type == token.TRUE || leftSideValue.Type == token.FALSE {
		// If no operator follows, but the leftSide itself is TRUE/FALSE, it's a valid factor.
		return &ast.BooleanExpression{Boolean: leftSideValue}, nil
	}

	// If leftSideValue was parsed, but no valid operator followed, and it's not TRUE/FALSE by itself
	return nil, &LogicalExpressionParsingError{customMessage: "Identifier or literal not followed by a valid operator", afterToken: &leftSideValue.Literal}
}

func (parser *Parser) getExpressionLeftSideValue() (token.Token, bool, error) {
	// Check if the current token can start an expression leaf (factor)
	if parser.currentToken.Type == token.IDENT ||
		parser.currentToken.Type == token.LITERAL ||
		parser.currentToken.Type == token.NULL ||
		parser.currentToken.Type == token.APOSTROPHE ||
		parser.currentToken.Type == token.TRUE ||
		parser.currentToken.Type == token.FALSE {

		var leftSide token.Token
		isAnonymitifier := false
		startedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()

		if startedWithApostrophe {
			isAnonymitifier = true
			value := ""
			for parser.currentToken.Type != token.EOF && parser.currentToken.Type != token.APOSTROPHE {
				value += parser.currentToken.Literal
				parser.nextToken()
			}
			leftSide = token.Token{Type: token.IDENT, Literal: value}
			finishedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()
			err := validateApostropheWrapping(startedWithApostrophe, finishedWithApostrophe, leftSide)
			if err != nil {
				return token.Token{}, isAnonymitifier, err
			}
		} else {
			leftSide = parser.currentToken
			parser.nextToken() // Consume the token
		}
		return leftSide, isAnonymitifier, nil
	}
	// If the token is not one of the above, it cannot start a simple factor.
	return token.Token{}, false, &SyntaxError{expecting: []string{"IDENT", "LITERAL", "NULL", "TRUE", "FALSE", "'"}, got: string(parser.currentToken.Type)}
}

// getConditionalExpression - Return ast.ConditionExpression created from tokens and validate the syntax
// This function assumes parser.currentToken is the operator (EQUAL, NOT, ISNULL, ISNOTNULL)
// when it's called by getFactor.
func (parser *Parser) getConditionalExpression(leftSide token.Token, isAnonymitifier bool) (*ast.ConditionExpression, error) {
	conditionalExpression := &ast.ConditionExpression{Condition: parser.currentToken}

	if isAnonymitifier {
		conditionalExpression.Left = ast.Anonymitifier{Token: leftSide}
	} else {
		conditionalExpression.Left = ast.Identifier{Token: leftSide}
	}

	// skip EQUAL or NOT
	parser.nextToken()

	if parser.currentToken.Type == token.IDENT || parser.currentToken.Type == token.LITERAL ||
		parser.currentToken.Type == token.NULL || parser.currentToken.Type == token.APOSTROPHE {
		startedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()

		if !startedWithApostrophe && parser.currentToken.Type == token.IDENT {
			conditionalExpression.Right = ast.Identifier{Token: parser.currentToken}
		} else {
			conditionalExpression.Right = ast.Anonymitifier{Token: parser.currentToken}
		}
		parser.nextToken()

		finishedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()
		err := validateApostropheWrapping(startedWithApostrophe, finishedWithApostrophe, conditionalExpression.Right.GetToken())
		if err != nil {
			return nil, err
		}
	} else {
		return nil, &SyntaxError{expecting: []string{token.APOSTROPHE, token.IDENT, token.LITERAL, token.NULL}, got: parser.currentToken.Literal}
	}

	return conditionalExpression, nil
}

// getContainExpression - Return ast.ContainExpression created from tokens and validate the syntax
func (parser *Parser) getContainExpression(leftSide token.Token, isAnonymitifier bool) (*ast.ContainExpression, error) {
	containExpression := &ast.ContainExpression{}

	if isAnonymitifier {
		return nil, &SyntaxError{expecting: []string{token.IDENT}, got: "'" + leftSide.Literal + "'"}
	}

	containExpression.Left = ast.Identifier{Token: leftSide}

	if parser.currentToken.Type == token.IN {
		containExpression.Contains = true
	} else {
		containExpression.Contains = false
	}

	// skip IN or NOTIN
	parser.nextToken()

	err := validateTokenAndSkip(parser, []token.Type{token.LPAREN})
	if err != nil {
		return nil, err
	}

	for parser.currentToken.Type == token.IDENT || parser.currentToken.Type == token.LITERAL || parser.currentToken.Type == token.NULL || parser.currentToken.Type == token.APOSTROPHE {
		startedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()

		err = validateToken(parser.currentToken.Type, []token.Type{token.IDENT, token.LITERAL, token.NULL})
		if err != nil {
			return nil, err
		}
		currentAnonymitifier := ast.Anonymitifier{Token: parser.currentToken}
		containExpression.Right = append(containExpression.Right, currentAnonymitifier)
		// Ignore token.IDENT, token.LITERAL or token.NULL
		parser.nextToken()

		finishedWithApostrophe := parser.skipIfCurrentTokenIsApostrophe()

		err = validateApostropheWrapping(startedWithApostrophe, finishedWithApostrophe, currentAnonymitifier.GetToken())
		if err != nil {
			return nil, err
		}

		if parser.currentToken.Type != token.COMMA {
			if parser.currentToken.Type != token.RPAREN {
				return nil, &SyntaxError{expecting: []string{token.COMMA, token.RPAREN}, got: string(parser.currentToken.Type)}
			}
			break
		}

		// Ignore token.COMMA
		parser.nextToken()
	}

	err = validateTokenAndSkip(parser, []token.Type{token.RPAREN})
	if err != nil {
		return nil, err
	}

	return containExpression, err
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
			err = parser.updateLastCommandWithWhereConstraints(sequence)
		case token.ORDER:
			err = parser.updateSelectCommandWithOrderByConstraints(sequence)
		case token.LIMIT:
			err = parser.updateSelectCommandWithLimitConstraints(sequence)
		case token.OFFSET:
			err = parser.updateSelectCommandWithOffsetConstraints(sequence)
		case token.JOIN, token.LEFT, token.RIGHT, token.INNER, token.FULL:
			err = parser.updateSelectCommandWithJoinConstraints(sequence)
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

func (parser *Parser) updateSelectCommandWithJoinConstraints(sequence *ast.Sequence) error {
	lastCommand, parserError := parser.getLastCommand(sequence, token.JOIN)
	if parserError != nil {
		return parserError
	}
	if lastCommand.TokenLiteral() != token.SELECT {
		return &SyntaxCommandExpectedError{command: "JOIN", neededCommands: []string{"SELECT"}}
	}
	selectCommand := lastCommand.(*ast.SelectCommand)
	newCommand, err := parser.parseJoinCommand()
	if err != nil {
		return err
	}
	selectCommand.JoinCommand = newCommand.(*ast.JoinCommand)
	return nil
}

func (parser *Parser) updateSelectCommandWithOffsetConstraints(sequence *ast.Sequence) error {
	lastCommand, parserError := parser.getLastCommand(sequence, token.OFFSET)
	if parserError != nil {
		return parserError
	}
	if lastCommand.TokenLiteral() != token.SELECT {
		return &SyntaxCommandExpectedError{command: "OFFSET", neededCommands: []string{"SELECT"}}
	}
	selectCommand := lastCommand.(*ast.SelectCommand)
	newCommand, err := parser.parseOffsetCommand()
	if err != nil {
		return err
	}
	selectCommand.OffsetCommand = newCommand.(*ast.OffsetCommand)
	return nil
}

func (parser *Parser) updateSelectCommandWithLimitConstraints(sequence *ast.Sequence) error {
	lastCommand, parserError := parser.getLastCommand(sequence, token.LIMIT)
	if parserError != nil {
		return parserError
	}
	if lastCommand.TokenLiteral() != token.SELECT {
		return &SyntaxCommandExpectedError{command: "LIMIT", neededCommands: []string{"SELECT"}}
	}
	selectCommand := lastCommand.(*ast.SelectCommand)
	newCommand, err := parser.parseLimitCommand()
	if err != nil {
		return err
	}
	selectCommand.LimitCommand = newCommand.(*ast.LimitCommand)
	return nil
}

func (parser *Parser) updateSelectCommandWithOrderByConstraints(sequence *ast.Sequence) error {
	lastCommand, parserError := parser.getLastCommand(sequence, token.ORDER)
	if parserError != nil {
		return parserError
	}

	if lastCommand.TokenLiteral() != token.SELECT {
		return &SyntaxCommandExpectedError{command: "ORDER BY", neededCommands: []string{"SELECT"}}
	}

	selectCommand := lastCommand.(*ast.SelectCommand)
	newCommand, err := parser.parseOrderByCommand()
	if err != nil {
		return err
	}
	selectCommand.OrderByCommand = newCommand.(*ast.OrderByCommand)
	return nil
}

func (parser *Parser) updateLastCommandWithWhereConstraints(sequence *ast.Sequence) error {
	lastCommand, parserError := parser.getLastCommand(sequence, token.WHERE)
	if parserError != nil {
		return parserError
	}

	if lastCommand.TokenLiteral() == token.SELECT {
		newCommand, err := parser.parseWhereCommand()
		if err != nil {
			return err
		}
		lastCommand.(*ast.SelectCommand).WhereCommand = newCommand.(*ast.WhereCommand)
	} else if lastCommand.TokenLiteral() == token.DELETE {
		newCommand, err := parser.parseWhereCommand()
		if err != nil {
			return err
		}
		lastCommand.(*ast.DeleteCommand).WhereCommand = newCommand.(*ast.WhereCommand)
	} else if lastCommand.TokenLiteral() == token.UPDATE {
		newCommand, err := parser.parseWhereCommand()
		if err != nil {
			return err
		}
		lastCommand.(*ast.UpdateCommand).WhereCommand = newCommand.(*ast.WhereCommand)
	} else {
		return &SyntaxCommandExpectedError{command: "WHERE", neededCommands: []string{"SELECT", "DELETE", "UPDATE"}}
	}
	return nil
}

func (parser *Parser) getLastCommand(sequence *ast.Sequence, currentToken string) (ast.Command, error) {
	if len(sequence.Commands) == 0 {
		return nil, &NoPredecessorParserError{command: currentToken}
	}
	lastCommand := sequence.Commands[len(sequence.Commands)-1]
	return lastCommand, nil
}
