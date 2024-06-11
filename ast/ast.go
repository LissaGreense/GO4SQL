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
// GetIdentifiers - Return array for all Identifiers within expression
type Expression interface {
	GetIdentifiers() []Identifier
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

// BooleanExpression - TokenType of Expression that represent single boolean value
//
// Example:
// TRUE
type BooleanExpression struct {
	Boolean token.Token // example: token.TRUE
}

func (ls BooleanExpression) GetIdentifiers() []Identifier {
	var identifiers []Identifier
	return identifiers
}

// ConditionExpression - TokenType of Expression that represent condition that is comparing value from column to static one
//
// Example:
// column1 EQUAL 123
type ConditionExpression struct {
	Left      Tifier      // name of column
	Right     Tifier      // value which column should have
	Condition token.Token // example: token.EQUAL
}

func (ls ConditionExpression) GetIdentifiers() []Identifier {
	var identifiers []Identifier

	if ls.Left.IsIdentifier() {
		identifiers = append(identifiers, Identifier{ls.Left.GetToken()})
	}

	if ls.Right.IsIdentifier() {
		identifiers = append(identifiers, Identifier{ls.Right.GetToken()})
	}

	return identifiers
}

// OperationExpression - TokenType of Expression that represent 2 other Expressions and conditional operation
//
// Example:
// TRUE OR FALSE
type OperationExpression struct {
	Left      Expression  // another operation, condition or boolean
	Right     Expression  // another operation, condition or boolean
	Operation token.Token // example: token.AND
}

func (ls OperationExpression) GetIdentifiers() []Identifier {
	var identifiers []Identifier

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
	Name        Identifier // name of the table
	ColumnNames []string
	ColumnTypes []token.Token
}

func (ls CreateCommand) CommandNode()         {}
func (ls CreateCommand) TokenLiteral() string { return ls.Token.Literal }

// InsertCommand - Part of Command that represent insertion of values into columns
//
// Example:
// INSERT INTO table1 VALUES('hello', 1);
type InsertCommand struct {
	Token  token.Token
	Name   Identifier // name of the table
	Values []token.Token
}

func (ls InsertCommand) CommandNode()         {}
func (ls InsertCommand) TokenLiteral() string { return ls.Token.Literal }

// SelectCommand - Part of Command that represent selecting values from tables
//
// Example:
// SELECT one, two FROM table1;
type SelectCommand struct {
	Token          token.Token
	Name           Identifier      // ex. name of table
	Space          []token.Token   // ex. column names
	WhereCommand   *WhereCommand   // optional
	OrderByCommand *OrderByCommand // optional
	LimitCommand   *LimitCommand   // optional
	OffsetCommand  *OffsetCommand  // optional
}

func (ls SelectCommand) CommandNode()         {}
func (ls SelectCommand) TokenLiteral() string { return ls.Token.Literal }

// HasWhereCommand - returns true if optional HasWhereCommand is present in SelectCommand
//
// Example:
// SELECT * FROM table WHERE column1 NOT 'hi';
// Returns true
//
// SELECT * FROM table;
// Returns false
func (ls SelectCommand) HasWhereCommand() bool {
	if ls.WhereCommand == nil {
		return false
	}
	return true
}

// HasOrderByCommand - returns true if optional OrderByCommand is present in SelectCommand
//
// Example:
// SELECT * FROM table ORDER BY column1 ASC;
// Returns true
//
// SELECT * FROM table;
// Returns false
func (ls SelectCommand) HasOrderByCommand() bool {
	if ls.OrderByCommand == nil {
		return false
	}
	return true
}

// HasLimitCommand - returns true if optional HasLimitCommand is present in SelectCommand
//
// Example:
// SELECT * FROM table LIMIT 5;
// Returns true
//
// SELECT * FROM table;
// Returns false
func (ls SelectCommand) HasLimitCommand() bool {
	if ls.LimitCommand == nil {
		return false
	}
	return true
}

// HasOffsetCommand - returns true if optional HasOffsetCommand is present in SelectCommand
//
// Example:
// SELECT * FROM table OFFSET 100;
// Returns true
//
// SELECT * FROM table LIMIT 10;
// Returns false
func (ls SelectCommand) HasOffsetCommand() bool {
	if ls.OffsetCommand == nil {
		return false
	}
	return true
}

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
	Token        token.Token
	Name         Identifier    // name of the table
	WhereCommand *WhereCommand // optional
}

func (ls DeleteCommand) CommandNode()         {}
func (ls DeleteCommand) TokenLiteral() string { return ls.Token.Literal }

// DropCommand - Part of Command that represent dropping table
//
// Example:
// DROP TABLE table;
type DropCommand struct {
	Token token.Token
	Name  Identifier // name of the table
}

func (ls DropCommand) CommandNode()         {}
func (ls DropCommand) TokenLiteral() string { return ls.Token.Literal }

// HasWhereCommand - returns true if optional HasWhereCommand is present in SelectCommand
//
// Example:
// SELECT * FROM table WHERE column1 NOT 'hi';
// Returns true
//
// SELECT * FROM table;
// Returns false
func (ls DeleteCommand) HasWhereCommand() bool {
	if ls.WhereCommand == nil {
		return false
	}
	return true
}

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

// LimitCommand - Part of Command that limits results from SelectCommand
type LimitCommand struct {
	Token token.Token
	Count int
}

func (ls LimitCommand) CommandNode()         {}
func (ls LimitCommand) TokenLiteral() string { return ls.Token.Literal }

// OffsetCommand - Part of Command that skip begging rows from SelectCommand
type OffsetCommand struct {
	Token token.Token
	Count int
}

func (ls OffsetCommand) CommandNode()         {}
func (ls OffsetCommand) TokenLiteral() string { return ls.Token.Literal }
