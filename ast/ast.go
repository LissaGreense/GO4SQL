package ast

import "github.com/LissaGreense/GO4SQL/token"

// Sequence - Sequence of operations commands
//
// Example:
// Commands[0] = SELECT * FROM Customers
// Commands[1] = WHERE City LIKE '%es%';
type Sequence struct {
	Commands []Command
}

// Node is connector between commands and expressions
type Node interface {
	TokenLiteral() string
}

// Command - Part of sequence - represent single static command
//
// Example:
// SELECT * FROM Customers;
type Command interface {
	Node
	CommandNode()
}

// Expression - Mathematical expression that is used to evaluate conditions
//
// Methods:
//
// ExpressionNode: Abstraction needed for creating tree abstraction in order to optimise evaluating
// GetIdentifiers - Return array of pointers for all Identifiers within expression
type Expression interface {
	// ExpressionNode TODO: Check if ExpressionNode is needed
	ExpressionNode()
	GetIdentifiers() []*Identifier
}

// Tifier - Interface that represent Token with string value
//
// Methods:
//
// IsIdentifier - Check if Tifier is Identifier
// GetToken - return token within Tifier
type Tifier interface {
	IsIdentifier() bool
	GetToken() token.Token
}

// TokenLiteral - Return first literal in sequence
func (p *Sequence) TokenLiteral() string {
	if len(p.Commands) > 0 {
		return p.Commands[0].TokenLiteral()
	} else {
		return ""
	}
}

// Identifier - Represent Token with string value that is equal to either column or table name
type Identifier struct {
	Token token.Token // the token.IDENT token
}

func (ls Identifier) IsIdentifier() bool    { return true }
func (ls Identifier) GetToken() token.Token { return ls.Token }

// Anonymitifier - Represent Token with string value that is equal to simple value that is put into columns
type Anonymitifier struct {
	Token token.Token // the token.IDENT token
}

func (ls Anonymitifier) IsIdentifier() bool    { return false }
func (ls Anonymitifier) GetToken() token.Token { return ls.Token }

// BooleanExpression - Type of Expression that represent single boolean value
//
// Example:
// TRUE
type BooleanExpression struct {
	Boolean token.Token // example: token.TRUE
}

func (ls BooleanExpression) ExpressionNode() {}
func (ls BooleanExpression) GetIdentifiers() []*Identifier {
	var identifiers []*Identifier
	return identifiers
}

// ConditionExpression - Type of Expression that represent condition that is comparing value from column to static one
//
// Example:
// column1 EQUAL 123
type ConditionExpression struct {
	Left      Tifier      // name of column
	Right     Tifier      // value which column should have
	Condition token.Token // example: token.EQUAL
}

func (ls ConditionExpression) ExpressionNode() {}
func (ls ConditionExpression) GetIdentifiers() []*Identifier {
	var identifiers []*Identifier

	if ls.Left.IsIdentifier() {
		identifiers = append(identifiers, &Identifier{ls.Left.GetToken()})
	}

	if ls.Right.IsIdentifier() {
		identifiers = append(identifiers, &Identifier{ls.Right.GetToken()})
	}

	return identifiers
}

// OperationExpression - Type of Expression that represent 2 other Expressions and conditional operation
//
// Example:
// TRUE OR FALSE
type OperationExpression struct {
	Left      Expression  // another operation, condition or boolean
	Right     Expression  // another operation, condition or boolean
	Operation token.Token // example: token.AND
}

func (ls OperationExpression) ExpressionNode() {}
func (ls OperationExpression) GetIdentifiers() []*Identifier {
	var identifiers []*Identifier

	identifiers = append(identifiers, ls.Left.GetIdentifiers()...)
	identifiers = append(identifiers, ls.Right.GetIdentifiers()...)

	return identifiers
}

// CreateCommand - Part of Command that represent creation of table
//
// Example:
// CREATE TABLE table1( one TEXT , two INT);
type CreateCommand struct {
	Token       token.Token
	Name        *Identifier // name of the table
	ColumnNames []string
	ColumnTypes []token.Token
}

func (ls CreateCommand) CommandNode()         {}
func (ls CreateCommand) TokenLiteral() string { return ls.Token.Literal }

// InsertCommand - Part of Command that represent insertion of values into columns
//
// Example:
// INSERT INTO table1 VALUES( 'hello', 1);
type InsertCommand struct {
	Token  token.Token
	Name   *Identifier // name of the table
	Values []token.Token
}

func (ls InsertCommand) CommandNode()         {}
func (ls InsertCommand) TokenLiteral() string { return ls.Token.Literal }

// SelectCommand - Part of Command that represent selecting values from tables
//
// Example:
// SELECT one, two FROM table1;
type SelectCommand struct {
	Token token.Token
	Name  *Identifier
	Space []token.Token // ex. column names
}

func (ls SelectCommand) CommandNode()         {}
func (ls SelectCommand) TokenLiteral() string { return ls.Token.Literal }

// WhereCommand - Part of Command that represent Where statement with expression that will qualify values from Select
//
// Example:
// WHERE column1 NOT 'goodbye' OR column2 EQUAL 3;
type WhereCommand struct {
	Token      token.Token
	Expression Expression
}

func (ls WhereCommand) CommandNode()         {}
func (ls WhereCommand) TokenLiteral() string { return ls.Token.Literal }

// DeleteCommand - Part of Command that represent deleting row from table
//
// Example:
// DELETE FROM tb1 WHERE two EQUAL 3;
type DeleteCommand struct {
	Token token.Token
	Name  *Identifier // name of the table
}

func (ls DeleteCommand) CommandNode()         {}
func (ls DeleteCommand) TokenLiteral() string { return ls.Token.Literal }

// OrderByCommand - Part of Command that ordering columns from SelectCommand
//
// Example:
// ORDER BY column1 ASC, column2 DESC;
type OrderByCommand struct {
	Token        token.Token
	SortPatterns []SortPattern // column name and sorting type
}

func (ls OrderByCommand) CommandNode()         {}
func (ls OrderByCommand) TokenLiteral() string { return ls.Token.Literal }

// SortPattern - Represent in which order declared columns should be sorted
type SortPattern struct {
	ColumnName token.Token // column name
	Order      token.Token // ASC or DESC
}
