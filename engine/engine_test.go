package engine

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
)

func TestParserCreateCommand(t *testing.T) {
	input :=
		`
		create table 	tbl( one TEXT , two INT );
		insert into tbl values( 'hello',	 10 );
		insert 	into tbl  values( 'goodbye', 20 );
		insert 	into tbl  values( 'byebye', 3333 );
		`

	lexer := lexer.RunLexer(input)
	parserInstance := parser.New(lexer)
	sequences := parserInstance.ParseSequence()

	if len(sequences.Commands) != 4 {
		t.Fatalf("sequences does not contain 4 statements. got=%d", len(sequences.Commands))
	}

	engine := New()
	engine.CreateTable((sequences.Commands[0]).(*ast.CreateCommand))
	engine.InsertIntoTable(sequences.Commands[1].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[2].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[3].(*ast.InsertCommand))

	if engine.Tables == nil {
		t.Error()
	}

	if len(engine.Tables["TBL"]) != 2 {
		t.Error()
	}

	if engine.Tables["TBL"]["ONE"].Type.Type != "TEXT" {
		t.Error()
	}
	if engine.Tables["TBL"]["TWO"].Type.Type != "INT" {
		t.Error()
	}
	if engine.Tables["TBL"]["ONE"].Values[0] != "HELLO" {
		t.Error()
	}
	if engine.Tables["TBL"]["ONE"].Values[1] != "GOODBYE" {
		t.Error()
	}
	if engine.Tables["TBL"]["ONE"].Values[2] != "BYEBYE" {
		t.Error()
	}

	if engine.Tables["TBL"]["TWO"].Values[0] != "10" {
		t.Error()
	}
	if engine.Tables["TBL"]["TWO"].Values[1] != "20" {
		t.Error()
	}
	if engine.Tables["TBL"]["TWO"].Values[2] != "3333" {
		t.Error()
	}
}
