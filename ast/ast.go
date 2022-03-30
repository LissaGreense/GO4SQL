package ast

import "github.com/LissaGreense/GO4SQL/token"

// Sequence of operations commands
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

// Part of sequence - represent single static command
// Example:
// SELECT * FROM Customers
type Command interface {
	Node
	CommandNode()
}

// Mathematical expression
// Example:
// CustomerID<5
type Expression interface {
	Node
	expressionNode()
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

type CreateCommand struct {
	Token       token.Token
	Name        *Identifier // name of the table
	ColumnNames []string
	ColumnTypes []token.Token
}

func (ls *CreateCommand) CommandNode()         {}
func (ls *CreateCommand) Node()                {}
func (ls *CreateCommand) TokenLiteral() string { return ls.Token.Literal }

type InsertCommand struct {
	Token  token.Token
	Name   *Identifier // name od the table
	Values []token.Token
}

type SelectCommand struct {
	Token token.Token
	Name  *Identifier
	Space []token.Token // ex. column names
}
