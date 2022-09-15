package engine

import (
	"fmt"
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
)

func TestCreateCommand(t *testing.T) {
	input :=
		`
		CREATE TABLE 	tbl( one TEXT , two INT );
		INSERT INTO tbl VALUES( 'hello',	 10 );
		INSERT 	INTO tbl  VALUES( 'goodbye', 20 );
		INSERT 	INTO tbl  VALUES( 'byebye', 3333 );
		`

	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
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

	if len(engine.Tables["tbl"]) != 2 {
		t.Error()
	}

	if engine.Tables["tbl"][0].Type.Type != "TEXT" {
		t.Error()
	}
	if engine.Tables["tbl"][1].Type.Type != "INT" {
		t.Error()
	}
	if fmt.Sprintf("%v", engine.Tables["tbl"][0].Values[0]) != "hello" {
		t.Error()
	}
	if fmt.Sprintf("%v", engine.Tables["tbl"][0].Values[1]) != "goodbye" {
		t.Error()
	}
	if fmt.Sprintf("%v", engine.Tables["tbl"][0].Values[2]) != "byebye" {
		t.Error()
	}

	if fmt.Sprintf("%v", engine.Tables["tbl"][1].Values[0]) != "10" {
		t.Error()
	}
	if fmt.Sprintf("%v", engine.Tables["tbl"][1].Values[1]) != "20" {
		t.Error()
	}
	if fmt.Sprintf("%v", engine.Tables["tbl"][1].Values[2]) != "3333" {
		t.Error()
	}
}

func TestSelectCommand(t *testing.T) {

	input :=
		`
		CREATE TABLE 	tbl( one TEXT , two INT, three INT, four TEXT );
		INSERT INTO tbl 	VALUES( 'hello',	1, 	11, 'q'  );
		INSERT 	INTO tbl  	VALUES( 'goodbye', 	2, 	22, 'w'  );
		INSERT 	INTO tbl  	VALUES( 'byebye', 	3, 	33,	'e'  );
		SELECT * FROM tbl;
		`

	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
	sequences := parserInstance.ParseSequence()

	if len(sequences.Commands) != 5 {
		t.Fatalf("sequences does not contain 5 statements. got=%d", len(sequences.Commands))
	}

	engine := New()
	engine.CreateTable((sequences.Commands[0]).(*ast.CreateCommand))
	engine.InsertIntoTable(sequences.Commands[1].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[2].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[3].(*ast.InsertCommand))

	result := engine.SelectFromTable(sequences.Commands[4].(*ast.SelectCommand)).ToString()

	expectedResult := "one|two|three|four" + "\n" + "'hello'|1|11|'q'" + "\n" + "'goodbye'|2|22|'w'" + "\n" + "'byebye'|3|33|'e'"

	if result != expectedResult {
		t.Error(result)
	}
}

func TestSelectWithColumnNamesCommand(t *testing.T) {
	input :=
		`
		CREATE TABLE 	tbl( one TEXT , two INT, three INT, four TEXT );
		INSERT INTO tbl 	VALUES( 'hello',	1, 	11, 'q'  );
		INSERT 	INTO tbl  	VALUES( 'goodbye', 	2, 	22, 'w'  );
		INSERT 	INTO tbl  	VALUES( 'byebye', 	3, 	33,	'e'  );
		SELECT one, two FROM tbl;
		SELECT two, one FROM tbl;
		SELECT one, two, three, four FROM tbl;
		`

	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
	sequences := parserInstance.ParseSequence()

	if len(sequences.Commands) != 7 {
		t.Fatalf("sequences does not contain 7 statements. got=%d", len(sequences.Commands))
	}

	engine := New()
	engine.CreateTable((sequences.Commands[0]).(*ast.CreateCommand))
	engine.InsertIntoTable(sequences.Commands[1].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[2].(*ast.InsertCommand))
	engine.InsertIntoTable(sequences.Commands[3].(*ast.InsertCommand))

	result := engine.SelectFromTable(sequences.Commands[4].(*ast.SelectCommand)).ToString()

	expectedResult := "one|two" + "\n" + "'hello'|1" + "\n" + "'goodbye'|2" + "\n" + "'byebye'|3"

	if result != expectedResult {
		t.Error(result)
	}

	result = engine.SelectFromTable(sequences.Commands[5].(*ast.SelectCommand)).ToString()

	expectedResult = "two|one" + "\n" + "1|'hello'" + "\n" + "2|'goodbye'" + "\n" + "3|'byebye'"

	if result != expectedResult {
		t.Error(result)
	}

	result = engine.SelectFromTable(sequences.Commands[6].(*ast.SelectCommand)).ToString()

	expectedResult = "one|two|three|four" + "\n" + "'hello'|1|11|'q'" + "\n" + "'goodbye'|2|22|'w'" + "\n" + "'byebye'|3|33|'e'"

	if result != expectedResult {
		t.Error(result)
	}
}
