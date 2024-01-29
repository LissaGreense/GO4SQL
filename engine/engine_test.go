package engine

import (
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

func TestSelectWithWhereLogicalOperationAnd(t *testing.T) {

	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 3, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1 WHERE one EQUAL 'goodbye' AND two NOT 2;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"goodbye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLogicalOperationOR(t *testing.T) {

	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 3, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1 WHERE one NOT 'goodbye' OR two EQUAL 3;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLogicalOperationOROperationAND(t *testing.T) {

	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 3, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1 WHERE one NOT 'goodbye' OR two EQUAL 3 AND four EQUAL 'e';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereEqualToTrue(t *testing.T) {

	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 3, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1 WHERE TRUE;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "2", "22", "w"},
			{"goodbye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereEqualToFalse(t *testing.T) {

	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 3, 	33,	'e'  );",
		},
		selectInput:    "SELECT * FROM tb1 WHERE FALSE;",
		expectedOutput: [][]string{},
	}

	engineTestSuite.runTestSuite(t)
}

func TestDelete(t *testing.T) {

	engineTestSuite := engineTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
			"DELETE FROM tb1 WHERE two EQUAL 3;",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"byebye", "3", "33", "e"}, // TODO DELETE THAT LINE LATER ON
			{"goodbye", "2", "22", "w"},
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
	expectedSequencesNumber := 0
	for inputIndex := 0; inputIndex < len(engineTestSuite.createInputs); inputIndex++ {
		input += engineTestSuite.createInputs[inputIndex] + "\n"
	}
	for inputIndex := 0; inputIndex < len(engineTestSuite.insertInputs); inputIndex++ {
		if strings.HasPrefix(engineTestSuite.insertInputs[inputIndex], "DELETE") {
			expectedSequencesNumber++
		}
		input += engineTestSuite.insertInputs[inputIndex] + "\n"
	}
	input += engineTestSuite.selectInput

	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
	sequences := parserInstance.ParseSequence()

	expectedSequencesNumber += len(engineTestSuite.createInputs) + len(engineTestSuite.insertInputs) + 1

	var actualTable *Table

	if strings.Contains(engineTestSuite.selectInput, " WHERE ") {

		// WHERE CONDITION

		expectedSequencesNumber++
		if len(sequences.Commands) != expectedSequencesNumber {
			t.Fatalf("sequences does not contain %d statements. got=%d", expectedSequencesNumber, len(sequences.Commands))
		}

		engine := engineTestSuite.getEngineWithInsertedValues(sequences)

		actualTable = engine.SelectFromTableWithWhere(sequences.Commands[len(sequences.Commands)-2].(*ast.SelectCommand), sequences.Commands[len(sequences.Commands)-1].(*ast.WhereCommand))
	} else {

		// NO WHERE CONDITION

		if len(sequences.Commands) != expectedSequencesNumber {
			t.Fatalf("sequences does not contain %d statements. got=%d", expectedSequencesNumber, len(sequences.Commands))
		}

		engine := engineTestSuite.getEngineWithInsertedValues(sequences)

		actualTable = engine.SelectFromTable(sequences.Commands[len(sequences.Commands)-1].(*ast.SelectCommand))
	}

	if len(engineTestSuite.expectedOutput) == 0 {
		if len(actualTable.Columns[0].Values) != 0 {
			t.Fatalf("Number of rows is incorrect, should be 0, got %d", len(actualTable.Columns))
		}
	} else {
		if len(actualTable.Columns) != len(engineTestSuite.expectedOutput[0]) {
			t.Fatalf("Number of columns is incorrect, expecting %d, got %d", len(engineTestSuite.expectedOutput[0]), len(actualTable.Columns))
		}

		if len(actualTable.Columns[0].Values) != len(engineTestSuite.expectedOutput)-1 {
			t.Fatalf("Number of rows is incorrect, expecting %d, got %d", len(engineTestSuite.expectedOutput)-1, len(actualTable.Columns[0].Values))
		}

		for iColumn := 0; iColumn < len(actualTable.Columns); iColumn++ {
			for iRow := 0; iRow < len(actualTable.Columns[0].Values); iRow++ {
				if engineTestSuite.expectedOutput[iRow+1][iColumn] != actualTable.Columns[iColumn].Values[iRow].ToString() {
					t.Fatalf("Value doesn't match, expected: %s, got: %s", engineTestSuite.expectedOutput[iRow+1][iColumn], actualTable.Columns[iColumn].Values[iRow].ToString())
				}
			}
		}
	}

}

func (engineTestSuite *engineTestSuite) getEngineWithInsertedValues(sequences *ast.Sequence) *DbEngine {
	engine := New()
	for commandIndex := 0; commandIndex < len(sequences.Commands); commandIndex++ {
		if createCommand, ok := sequences.Commands[commandIndex].(*ast.CreateCommand); ok {
			engine.CreateTable(createCommand)
		}
		if insertCommand, ok := sequences.Commands[commandIndex].(*ast.InsertCommand); ok {
			engine.InsertIntoTable(insertCommand)
		}
		if deleteCommand, ok := sequences.Commands[commandIndex].(*ast.DeleteCommand); ok {
			whereCommand := sequences.Commands[commandIndex+1].(*ast.WhereCommand)
			engine.DeleteFromTable(deleteCommand, whereCommand)
		}
	}
	return engine
}
