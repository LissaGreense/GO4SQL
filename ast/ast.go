package ast

import "github.com/LissaGreense/GO4SQL/token"

// Sequence - Sequence of operations commands
// Example:
// Command[0] = SELECT * FROM Customers
// Command[1] = WHERE City LIKE '%es%';
type Sequence struct {
	Commands []Command
}

// Node is connector between commands and expressions
type Node interface {
	TokenLiteral() string
}

// Command - Part of sequence - represent single static command
// Example:
// SELECT * FROM Customers
type Command interface {
	Node
	CommandNode()
}

// Expression - Mathematical expression
// Example:
// CustomerID<5
type Expression interface {
	ExpressionNode()
	GetIdentifiers() []*Identifier
}

type Tifier interface {
	IsIdentifier() bool
	GetToken() token.Token
}

func (p *Sequence) TokenLiteral() string {
	if len(p.Commands) > 0 {
		return p.Commands[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token // the token.IDENT token
}

func (ls Identifier) IsIdentifier() bool    { return true }
func (ls Identifier) GetToken() token.Token { return ls.Token }

type Anonymitifier struct {
	Token token.Token // the token.IDENT token
}

func (ls Anonymitifier) IsIdentifier() bool    { return false }
func (ls Anonymitifier) GetToken() token.Token { return ls.Token }

type BooleanExpresion struct {
	Boolean token.Token // example: token.TRUE
}

func (ls BooleanExpresion) ExpressionNode() {}
func (ls BooleanExpresion) GetIdentifiers() []*Identifier {
	var identifiers []*Identifier
	return identifiers
}

type ConditionExpresion struct {
	Left      Tifier      // name of column
	Right     Tifier      // value which column should have
	Condition token.Token // example: token.EQUAL
}

func (ls ConditionExpresion) ExpressionNode() {}
func (ls ConditionExpresion) GetIdentifiers() []*Identifier {
	var identifiers []*Identifier

	if ls.Left.IsIdentifier() {
		identifiers = append(identifiers, &Identifier{ls.Left.GetToken()})
	}

	if ls.Right.IsIdentifier() {
		identifiers = append(identifiers, &Identifier{ls.Right.GetToken()})
	}

	return identifiers
}

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

type CreateCommand struct {
	Token       token.Token
	Name        *Identifier // name of the table
	ColumnNames []string
	ColumnTypes []token.Token
}

func (ls CreateCommand) CommandNode()         {}
func (ls CreateCommand) TokenLiteral() string { return ls.Token.Literal }

type InsertCommand struct {
	Token  token.Token
	Name   *Identifier // name of the table
	Values []token.Token
}

func (ls InsertCommand) CommandNode()         {}
func (ls InsertCommand) TokenLiteral() string { return ls.Token.Literal }

type SelectCommand struct {
	Token token.Token
	Name  *Identifier
	Space []token.Token // ex. column names
}

func (ls SelectCommand) CommandNode()         {}
func (ls SelectCommand) TokenLiteral() string { return ls.Token.Literal }

type WhereCommand struct {
	Token      token.Token
	Expression Expression
}

func (ls WhereCommand) CommandNode()         {}
func (ls WhereCommand) TokenLiteral() string { return ls.Token.Literal }
