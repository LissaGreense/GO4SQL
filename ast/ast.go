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
	Node
	ExpressionNode()
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

type Condition struct {
	Left           *Identifier // name of column
	Right          token.Token // value which column should have
	OperationToken token.Token // example: token.EQUAL
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
	Expression *Condition
}

func (ls WhereCommand) CommandNode()         {}
func (ls WhereCommand) TokenLiteral() string { return ls.Token.Literal }
