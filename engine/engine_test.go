package engine

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
	"github.com/LissaGreense/GO4SQL/token"
)

const (
	tableName = "tbl"
)

func getCreateAndInsertCommands() string {
	return `
	CREATE  TABLE  ` + tableName + `( one TEXT , two INT, three INT, four TEXT );
	INSERT  INTO   ` + tableName + ` VALUES( 'hello',	1, 	11, 'q'  );
	INSERT 	INTO   ` + tableName + ` VALUES( 'goodbye', 	2, 	22, 'w'  );
	INSERT 	INTO   ` + tableName + ` VALUES( 'byebye', 	3, 	33,	'e'  );
	`
}

func getOneColumn() *Column {
	column := &Column{
		Name:   "one",
		Type:   token.Token{Type: token.TEXT, Literal: "TEXT"},
		Values: make([]ValueInterface, 0),
	}
	column.Values = append(column.Values, StringValue{Value: "hello"})
	column.Values = append(column.Values, StringValue{Value: "goodbye"})
	column.Values = append(column.Values, StringValue{Value: "byebye"})
	return column
}

func getTwoColumn() *Column {
	column := &Column{
		Name:   "two",
		Type:   token.Token{Type: token.INT, Literal: "INT"},
		Values: make([]ValueInterface, 0),
	}
	column.Values = append(column.Values, IntegerValue{Value: 1})
	column.Values = append(column.Values, IntegerValue{Value: 2})
	column.Values = append(column.Values, IntegerValue{Value: 3})
	return column
}

func getThreeColumn() *Column {
	column := &Column{
		Name:   "three",
		Type:   token.Token{Type: token.INT, Literal: "INT"},
		Values: make([]ValueInterface, 0),
	}
	column.Values = append(column.Values, IntegerValue{Value: 11})
	column.Values = append(column.Values, IntegerValue{Value: 22})
	column.Values = append(column.Values, IntegerValue{Value: 33})
	return column
}

func getFourColumn() *Column {
	column := &Column{
		Name:   "four",
		Type:   token.Token{Type: token.TEXT, Literal: "TEXT"},
		Values: make([]ValueInterface, 0),
	}
	column.Values = append(column.Values, StringValue{Value: "q"})
	column.Values = append(column.Values, StringValue{Value: "w"})
	column.Values = append(column.Values, StringValue{Value: "e"})
	return column
}

func getExpectedTable() *Table {
	expectedTable := &Table{Columns: make([]*Column, 0)}
	expectedTable.Columns = append(expectedTable.Columns, getOneColumn())
	expectedTable.Columns = append(expectedTable.Columns, getTwoColumn())
	expectedTable.Columns = append(expectedTable.Columns, getThreeColumn())
	expectedTable.Columns = append(expectedTable.Columns, getFourColumn())

	return expectedTable
}

func TestCreateCommand(t *testing.T) {
	input := getCreateAndInsertCommands()

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

	tableFromEngineStruct := engine.Tables[tableName]
	expectedTable := getExpectedTable()
	if tableFromEngineStruct.isEqual(expectedTable) == false {
		t.Error("\n" + tableFromEngineStruct.ToString() + "\n not euqal to: \n" + expectedTable.ToString())
	}
}

func TestSelectCommand(t *testing.T) {

	input := getCreateAndInsertCommands() + "\n SELECT * FROM " + tableName + ";"

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

	actualTable := engine.SelectFromTable(sequences.Commands[4].(*ast.SelectCommand))
	expectedTable := getExpectedTable()

	if actualTable.isEqual(expectedTable) == false {
		t.Error("\n" + actualTable.ToString() + "\n not euqal to: \n" + expectedTable.ToString())
	}
}

func TestSelectWithColumnNamesCommand(t *testing.T) {
	input := getCreateAndInsertCommands() +
		`
		SELECT one, two FROM ` + tableName + `;
		SELECT two, one FROM ` + tableName + `;
		SELECT one, two, three, four FROM ` + tableName + `;
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

	actualTable := engine.SelectFromTable(sequences.Commands[4].(*ast.SelectCommand))
	expectedTable := &Table{Columns: make([]*Column, 0)}
	expectedTable.Columns = append(expectedTable.Columns, getOneColumn())
	expectedTable.Columns = append(expectedTable.Columns, getTwoColumn())
	if actualTable.isEqual(expectedTable) == false {
		t.Error("\n" + actualTable.ToString() + "\n not euqal to: \n" + expectedTable.ToString())
	}

	actualTable = engine.SelectFromTable(sequences.Commands[5].(*ast.SelectCommand))
	expectedTable = &Table{Columns: make([]*Column, 0)}
	expectedTable.Columns = append(expectedTable.Columns, getTwoColumn())
	expectedTable.Columns = append(expectedTable.Columns, getOneColumn())
	if actualTable.isEqual(expectedTable) == false {
		t.Error("\n" + actualTable.ToString() + "\n not euqal to: \n" + expectedTable.ToString())
	}

	actualTable = engine.SelectFromTable(sequences.Commands[6].(*ast.SelectCommand))
	expectedResult := getExpectedTable()
	if actualTable.isEqual(expectedResult) == false {
		t.Error("\n" + actualTable.ToString() + "\n not euqal to: \n" + expectedResult.ToString())
	}
}

//
//func TestSelectWithWhereEqual(t *testing.T) {
//	input := getCreateAndInsertCommands() +
//		`
//		SELECT one, three FROM ` + tableName + ` WHERE one EQUAL 'hello';
//		`
//
//	lexerInstance := lexer.RunLexer(input)
//	parserInstance := parser.New(lexerInstance)
//	sequences := parserInstance.ParseSequence()
//
//	if len(sequences.Commands) != 5 {
//		t.Fatalf("sequences does not contain 5 statements. got=%d", len(sequences.Commands))
//	}
//
//	engine := New()
//	engine.CreateTable((sequences.Commands[0]).(*ast.CreateCommand))
//	engine.InsertIntoTable(sequences.Commands[1].(*ast.InsertCommand))
//	engine.InsertIntoTable(sequences.Commands[2].(*ast.InsertCommand))
//	engine.InsertIntoTable(sequences.Commands[3].(*ast.InsertCommand))
//
//	actualTable := engine.SelectFromTable(sequences.Commands[4].(*ast.SelectCommand))
//	expectedTable := &Table{Columns: make([]*Column, 0)}
//	column := &Column{
//		Name:   "one",
//		Type:   token.Token{Type: token.TEXT, Literal: "TEXT"},
//		Values: make([]ValueInterface, 0),
//	}
//	column.Values = append(column.Values, StringValue{Value: "hello"})
//	expectedTable.Columns = append(expectedTable.Columns, column)
//
//	column = &Column{
//		Name:   "three",
//		Type:   token.Token{Type: token.INT, Literal: "INT"},
//		Values: make([]ValueInterface, 0),
//	}
//	column.Values = append(column.Values, IntegerValue{Value: 11})
//	expectedTable.Columns = append(expectedTable.Columns, column)
//
//	if actualTable.isEqual(expectedTable) == false {
//		t.Error("\n" + actualTable.ToString() + "\n not equal to: \n" + expectedTable.ToString())
//	}
//}

//func TestSelectWithWhereNotEqual(t *testing.T) {
//	input := getCreateAndInsertCommands() +
//		`
//		SELECT one, three FROM ` + tableName + ` WHERE two NOT EQUAL 22;
//		`
//
//	lexerInstance := lexer.RunLexer(input)
//	parserInstance := parser.New(lexerInstance)
//	sequences := parserInstance.ParseSequence()
//
//	if len(sequences.Commands) != 5 {
//		t.Fatalf("sequences does not contain 5 statements. got=%d", len(sequences.Commands))
//	}
//
//	engine := New()
//	engine.CreateTable((sequences.Commands[0]).(*ast.CreateCommand))
//	engine.InsertIntoTable(sequences.Commands[1].(*ast.InsertCommand))
//	engine.InsertIntoTable(sequences.Commands[2].(*ast.InsertCommand))
//	engine.InsertIntoTable(sequences.Commands[3].(*ast.InsertCommand))
//
//	actualTable := engine.SelectFromTable(sequences.Commands[4].(*ast.SelectCommand))
//	expectedTable := &Table{Columns: make([]*Column, 0)}
//	column := &Column{
//		Name:   "one",
//		Type:   token.Token{Type: token.TEXT, Literal: "TEXT"},
//		Values: make([]ValueInterface, 0),
//	}
//	column.Values = append(column.Values, StringValue{Value: "hello"})
//	column.Values = append(column.Values, StringValue{Value: "byebye"})
//	expectedTable.Columns = append(expectedTable.Columns, column)
//
//	column = &Column{
//		Name:   "three",
//		Type:   token.Token{Type: token.INT, Literal: "INT"},
//		Values: make([]ValueInterface, 0),
//	}
//	column.Values = append(column.Values, IntegerValue{Value: 11})
//	column.Values = append(column.Values, IntegerValue{Value: 33})
//	expectedTable.Columns = append(expectedTable.Columns, column)
//
//	if actualTable.isEqual(expectedTable) == false {
//		t.Error("\n" + actualTable.ToString() + "\n not equal to: \n" + expectedTable.ToString())
//	}
//}
