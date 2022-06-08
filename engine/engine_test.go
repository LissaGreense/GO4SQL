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

	if engine.Tables["TBL"][0].Type.Type != "TEXT" {
		t.Error()
	}
	if engine.Tables["TBL"][1].Type.Type != "INT" {
		t.Error()
	}
	if engine.Tables["TBL"][0].Values[0] != "HELLO" {
		t.Error()
	}
	if engine.Tables["TBL"][0].Values[1] != "GOODBYE" {
		t.Error()
	}
	if engine.Tables["TBL"][0].Values[2] != "BYEBYE" {
		t.Error()
	}

	if engine.Tables["TBL"][1].Values[0] != "10" {
		t.Error()
	}
	if engine.Tables["TBL"][1].Values[1] != "20" {
		t.Error()
	}
	if engine.Tables["TBL"][1].Values[2] != "3333" {
		t.Error()
	}
}

func TestSelectCommand(t *testing.T) {

	input :=
		`
		create table 	tbl( one TEXT , two INT, three INT, four TEXT );
		insert into tbl 	values( 'hello',	1, 	11, 'q'  );
		insert 	into tbl  	values( 'goodbye', 	2, 	22, 'w'  );
		insert 	into tbl  	values( 'byebye', 	3, 	33,	'e'  );
		select * from tbl;
		`

	lexer := lexer.RunLexer(input)
	parserInstance := parser.New(lexer)
	sequences := parserInstance.ParseSequence()

	if len(sequences.Commands) != 5 {
		t.Fatalf("sequences does not contain 5 statements. got=%d", len(sequences.Commands))
	}

	engine := New()
	engine.CreateTable((sequences.Commands[0]).(*ast.CreateCommand))
	engine.InsertIntoTable(sequences.Commands[1].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[2].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[3].(*ast.InsertCommand))

	result := engine.SelectFromTable(sequences.Commands[4].(*ast.SelectCommand))

	expectedResult := "ONE|TWO|THREE|FOUR" + "\n" + "'HELLO'|1|11|'Q'" + "\n" + "'GOODBYE'|2|22|'W'" + "\n" + "'BYEBYE'|3|33|'E'"

	if result != expectedResult {
		t.Error(result)
	}
}
