package engine

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
)

func TestCreate(t *testing.T) {
	simpleCreateCase := engineDBContentTestSuite{
		inputs:             []string{"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );"},
		expectedTableNames: []string{"tb1"},
	}

	simpleCreateCase.runTestSuite(t)

	multiplyCreationCase := engineDBContentTestSuite{
		inputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
			"CREATE TABLE tb2( one TEXT, two INT, three INT, four TEXT );",
		},
		expectedTableNames: []string{"tb1", "tb2"},
	}

	multiplyCreationCase.runTestSuite(t)

}

func TestDrop(t *testing.T) {
	simpleDropCase := engineDBContentTestSuite{
		inputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
			"DROP TABLE tb1;",
		},
		expectedTableNames: []string{},
	}
	simpleDropCase.runTestSuite(t)
}

func TestSelectCommand(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
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

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
			"DELETE FROM tb1 WHERE two EQUAL 3;",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "2", "22", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestOrderBy(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 ORDER BY two ASC;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"byebye", "1", "33", "e"},
			{"goodbye", "2", "22", "w"},
			{"hello", "3", "11", "q"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestOrderByWithWhere(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'Ahello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'Bgoodbye', 3, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1 WHERE one NOT 'goodbye' OR two EQUAL 3 AND four EQUAL 'e' ORDER BY one DESC;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"Bgoodbye", "3", "33", "e"},
			{"Ahello", "1", "11", "q"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestOrderByWithMultipleSorts(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'aa'  );",
			"INSERT INTO tb1 VALUES( 'sorry',   2, 	55, 'ba'  );",
		},
		selectInput: "SELECT one FROM tb1 WHERE TRUE ORDER BY two ASC, four DESC;",
		expectedOutput: [][]string{
			{"one"},
			{"byebye"},
			{"sorry"},
			{"goodbye"},
			{"hello"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

type engineDBContentTestSuite struct {
	inputs             []string
	expectedTableNames []string
}

func (engineTestSuite *engineDBContentTestSuite) runTestSuite(t *testing.T) {
	sequences := getSequences(inputsToString(engineTestSuite.inputs))
	engine := New()
	engine.Evaluate(sequences)

	if len(engine.Tables) != len(engineTestSuite.expectedTableNames) {
		t.Fatalf("Number of tables is incorrect, should be %d, got %d", len(engineTestSuite.expectedTableNames), len(engine.Tables))
	}

	for _, tableName := range engineTestSuite.expectedTableNames {
		if engine.Tables[tableName] == nil {
			t.Fatalf("Expected table '%s' does not exist", tableName)
		}
	}
}

type engineTableContentTestSuite struct {
	createInputs          []string
	insertAndDeleteInputs []string
	selectInput           string
	expectedOutput        [][]string
}

func (engineTestSuite *engineTableContentTestSuite) runTestSuite(t *testing.T) {
	expectedSequencesNumber := 0

	input := inputsToString(engineTestSuite.createInputs) + inputsToString(engineTestSuite.insertAndDeleteInputs)

	sequencesWithoutSelect := getSequences(input)
	selectCommand := getSequences(engineTestSuite.selectInput)

	expectedSequencesNumber += len(engineTestSuite.createInputs) + len(engineTestSuite.insertAndDeleteInputs) + 1

	if len(sequencesWithoutSelect.Commands)+len(selectCommand.Commands) != expectedSequencesNumber {
		t.Fatalf("sequences does not contain %d statements. got=%d", expectedSequencesNumber, len(sequencesWithoutSelect.Commands))
	}

	engine := New()
	engine.Evaluate(sequencesWithoutSelect)
	actualTable := engine.getSelectResponse(selectCommand.Commands[0].(*ast.SelectCommand))

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

func inputsToString(inputs []string) string {
	input := ""

	for inputIndex := 0; inputIndex < len(inputs); inputIndex++ {
		input += inputs[inputIndex] + "\n"
	}

	return input
}

func getSequences(input string) *ast.Sequence {
	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
	sequences := parserInstance.ParseSequence()

	return sequences
}
