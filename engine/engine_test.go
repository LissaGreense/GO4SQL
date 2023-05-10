package engine

import (
	"log"
	"strings"
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
)

func TestSelectCommand(t *testing.T) {
	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "2", "22", "w"},
			{"byebye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithColumnNamesCommand(t *testing.T) {
	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "2", "22", "w"},
			{"byebye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereEqual(t *testing.T) {
	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one EQUAL 'hello';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereNotEqual(t *testing.T) {

	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE three NOT 22;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"byebye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

type engineTestSuite struct {
	createInputs   []string
	insertInputs   []string
	selectInput    string
	expectedOutput [][]string
}

func (engineTestSuite *engineTestSuite) runTestSuite(t *testing.T) {
	input := ""
	for i := 0; i < len(engineTestSuite.createInputs); i++ {
		input += engineTestSuite.createInputs[i] + "\n"
	}
	for i := 0; i < len(engineTestSuite.insertInputs); i++ {
		input += engineTestSuite.insertInputs[i] + "\n"
	}
	input += engineTestSuite.selectInput

	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
	sequences := parserInstance.ParseSequence()

	expectedSequencesNumber := len(engineTestSuite.createInputs) + len(engineTestSuite.insertInputs) + 1

	var actualTable *Table

	if strings.Contains(engineTestSuite.selectInput, " WHERE ") {

		// WHERE CONDITION

		expectedSequencesNumber++
		if len(sequences.Commands) != expectedSequencesNumber {
			t.Fatalf("sequences does not contain %d statements. got=%d", expectedSequencesNumber, len(sequences.Commands))
		}

		engine := New()
		for i := 0; i < len(sequences.Commands)-2; i++ {
			if createCommand, ok := sequences.Commands[i].(*ast.CreateCommand); ok {
				engine.CreateTable(createCommand)
			}
			if insertCommand, ok := sequences.Commands[i].(*ast.InsertCommand); ok {
				engine.InsertIntoTable(insertCommand)
			}
		}

		actualTable = engine.SelectFromTableWithWhere(sequences.Commands[len(sequences.Commands)-2].(*ast.SelectCommand), sequences.Commands[len(sequences.Commands)-1].(*ast.WhereCommand))
	} else {

		// NO WHERE CONDITION

		if len(sequences.Commands) != expectedSequencesNumber {
			t.Fatalf("sequences does not contain %d statements. got=%d", expectedSequencesNumber, len(sequences.Commands))
		}

		engine := New()
		for i := 0; i < len(sequences.Commands)-1; i++ {
			if createCommand, ok := sequences.Commands[i].(*ast.CreateCommand); ok {
				engine.CreateTable(createCommand)
			}
			if insertCommand, ok := sequences.Commands[i].(*ast.InsertCommand); ok {
				engine.InsertIntoTable(insertCommand)
			}
		}

		actualTable = engine.SelectFromTable(sequences.Commands[len(sequences.Commands)-1].(*ast.SelectCommand))
	}

	if len(actualTable.Columns) != len(engineTestSuite.expectedOutput[0]) {
		log.Fatalf("Number of columns is incorrect, expecting %d, got %d", len(engineTestSuite.expectedOutput[0]), len(actualTable.Columns))
	}

	if len(actualTable.Columns[0].Values) != len(engineTestSuite.expectedOutput)-1 {
		log.Fatalf("Number of rows is incorrect, expecting %d, got %d", len(engineTestSuite.expectedOutput)-1, len(actualTable.Columns[0].Values))
	}

	for iColumn := 0; iColumn < len(actualTable.Columns); iColumn++ {
		if actualTable.Columns[iColumn].Name != engineTestSuite.expectedOutput[0][iColumn] {
			t.Fatalf("Column names doesn't match, expected: %s, got: %s", engineTestSuite.expectedOutput[0][iColumn], actualTable.Columns[iColumn].Name)
		}

		for iRow := 0; iRow < len(actualTable.Columns[0].Values); iRow++ {
			if engineTestSuite.expectedOutput[iRow+1][iColumn] != actualTable.Columns[iColumn].Values[iRow].ToString() {
				t.Fatalf("Value doesn't match, expected: %s, got: %s", engineTestSuite.expectedOutput[iRow+1][iColumn], actualTable.Columns[iColumn].Values[iRow].ToString())
			}
		}
	}
}
