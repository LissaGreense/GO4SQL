package engine

import (
	"log"
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
	"github.com/LissaGreense/GO4SQL/token"
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
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, NULL  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	NULL, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "2", "22", "NULL"},
			{"byebye", "NULL", "33", "e"},
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
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	NULL, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE three NOT NULL;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"byebye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereContains(t *testing.T) {

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	NULL, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE three IN (11, NULL, 67);",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "2", "NULL", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereNotContains(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one NOTIN ('hello', 'byebye', 'youAreTheBest');",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"goodbye", "2", "22", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereContainsButResponseIsEmpty(t *testing.T) {

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one IN ('I', 'dont', 'exist', 'anymore');",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereNotContainsButResponseIsEmpty(t *testing.T) {

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE two NOTIN (1, 2, 3, 4);",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikePrefix(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE 'he%';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeSuffix(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE '%bye';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"goodbye", "2", "22", "w"},
			{"byebye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeContains(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE '%ood%';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"goodbye", "2", "22", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeSingleCharWildcard(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'hallo', 	2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE 'h_llo';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"hallo", "2", "22", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeIsCaseSensitive(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'HELLO',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'hello', 	2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE 'hello';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "2", "22", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeOnNullValueDoesNotMatch(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( NULL, 		2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE '%';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereNotLike(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one NOT LIKE '%bye';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereNotLikeOnNullValueMatches(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( NULL, 		2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one NOT LIKE '%';",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"NULL", "2", "22", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeCombinedWithAndOr(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE '%bye' AND two EQUAL 3;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"byebye", "3", "33", "e"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeAgainstNullPattern(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE NULL;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestSelectWithWhereLikeAgainstColumnPattern(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'hel%'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'by%'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'he%'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 WHERE one LIKE four;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "hel%"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestJoinOnLikeCondition(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE books( title TEXT);",
			"CREATE TABLE authors( title_pattern TEXT, name TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO books VALUES('Fire and Ash');",
			"INSERT INTO books VALUES('Earthbound');",
			"INSERT INTO authors VALUES('Fire%', 'Alissa Ireneus');",
			"INSERT INTO authors VALUES('Earth%', 'Reynold Boyka');",
		},
		selectInput: "SELECT books.title, authors.name FROM books JOIN authors ON books.title LIKE authors.title_pattern;",
		expectedOutput: [][]string{
			{"books.title", "authors.name"},
			{"Fire and Ash", "Alissa Ireneus"},
			{"Earthbound", "Reynold Boyka"},
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
			"INSERT INTO tb1 VALUES( 'goodbye', NULL, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 3, 	33,	'e'  );",
		},
		selectInput: "SELECT * FROM tb1 WHERE one EQUAL 'goodbye' AND two NOT NULL;",
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
			"INSERT INTO tb1 VALUES( 'goodbye', 3, 	33,	NULL  );",
		},
		selectInput: "SELECT * FROM tb1 WHERE one NOT 'goodbye' OR two IN (3) AND four EQUAL NULL;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "3", "33", "NULL"},
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

func TestDistinctSelect(t *testing.T) {

	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', NULL, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', NULL, 	22, 'w'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT DISTINCT * FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"goodbye", "NULL", "22", "w"},
			{"goodbye", "2", "22", "w"},
		},
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

func TestUpdateWithWhere(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
			"UPDATE tb1 SET one TO 'hi hello', three TO NULL WHERE two EQUAL 3;",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hello", "1", "11", "q"},
			{"hi hello", "3", "NULL", "e"},
			{"goodbye", "2", "22", "w"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestCreateWithNotNullColumn(t *testing.T) {
	engine := New()
	evaluateInputs(t, engine, "CREATE TABLE users( id INT NOT NULL, nickname TEXT );")

	columns := engine.Tables["users"].Columns
	if !columns[0].IsNotNull {
		t.Errorf("expected column %s to be NOT NULL", columns[0].Name)
	}
	if columns[1].IsNotNull {
		t.Errorf("expected column %s to be nullable", columns[1].Name)
	}
}

func TestNotNullColumnAcceptsValuesAndNullableColumnAcceptsNull(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE users( id INT NOT NULL, nickname TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO users VALUES( 1, 'kate' );",
			"INSERT INTO users VALUES( 2, NULL );",
			"UPDATE users SET nickname TO NULL WHERE id EQUAL 1;",
			"UPDATE users SET id TO NULL WHERE id EQUAL 3;",
			"UPDATE users SET id TO 5 WHERE id EQUAL 2;",
		},
		selectInput: "SELECT * FROM users;",
		expectedOutput: [][]string{
			{"id", "nickname"},
			{"1", "NULL"},
			{"5", "NULL"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestRejectedNotNullCommandsLeaveTableUnchanged(t *testing.T) {
	engine := New()
	evaluateInputs(t, engine,
		"CREATE TABLE users( nickname TEXT, id INT NOT NULL );",
		"INSERT INTO users VALUES( 'kate', 1 );",
	)

	for _, input := range []string{
		"INSERT INTO users VALUES( 'bob', NULL );",
		"UPDATE users SET nickname TO 'changed', id TO NULL;",
	} {
		sequences, err := parser.New(lexer.RunLexer(input)).ParseSequence()
		if err != nil {
			t.Fatalf("parser error: %s", err)
		}
		if _, err := engine.Evaluate(sequences); err == nil {
			t.Fatalf("expected error for %q", input)
		}
	}

	columns := engine.Tables["users"].Columns
	if len(columns[0].Values) != 1 || len(columns[1].Values) != 1 {
		t.Fatalf("expected 1 row in each column, got %d and %d", len(columns[0].Values), len(columns[1].Values))
	}
	if columns[0].Values[0].ToString() != "kate" || columns[1].Values[0].ToString() != "1" {
		t.Fatalf("expected row (kate, 1), got (%s, %s)", columns[0].Values[0].ToString(), columns[1].Values[0].ToString())
	}
}

func evaluateInputs(t *testing.T, engine *DbEngine, inputs ...string) {
	for _, input := range inputs {
		sequences, err := parser.New(lexer.RunLexer(input)).ParseSequence()
		if err != nil {
			t.Fatalf("parser error for %q: %s", input, err)
		}
		if _, err := engine.Evaluate(sequences); err != nil {
			t.Fatalf("engine error for %q: %s", input, err)
		}
	}
}

func TestUpdateWithoutWhere(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	1, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	3, 	33,	'e'  );",
			"UPDATE tb1 SET one TO 'hi hello', three TO 5;",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"hi hello", "1", "5", "q"},
			{"hi hello", "3", "5", "e"},
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
			"INSERT INTO tb1 VALUES( 'byebye', 	NULL, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 2, 	22, 'w'  );",
		},
		selectInput: "SELECT one, two, three, four FROM tb1 ORDER BY two ASC;",
		expectedOutput: [][]string{
			{"one", "two", "three", "four"},
			{"byebye", "NULL", "33", "e"},
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

func TestLimit(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',		3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 		1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 	2, 	22, 'aa'  );",
			"INSERT INTO tb1 VALUES( 'sorry',		2, 	55, 'ba'  );",
			"INSERT INTO tb1 VALUES( 'welcome',		2, 	66, 'bb'  );",
			"INSERT INTO tb1 VALUES( 'seeYouLater', 2, 	95, 'ab'  );",
		},
		selectInput: "SELECT one FROM tb1 LIMIT 2;",
		expectedOutput: [][]string{
			{"one"},
			{"hello"},
			{"byebye"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestLimitEqualToZero(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',		3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 		1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 	2, 	22, 'aa'  );",
			"INSERT INTO tb1 VALUES( 'sorry',		2, 	55, 'ba'  );",
			"INSERT INTO tb1 VALUES( 'welcome',		2, 	66, 'bb'  );",
			"INSERT INTO tb1 VALUES( 'seeYouLater', 2, 	95, 'ab'  );",
		},
		selectInput: "SELECT one FROM tb1 LIMIT 0;",
		expectedOutput: [][]string{
			{"one"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestLimitThatIsMoreThanSize(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',		3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 		1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 	2, 	22, 'aa'  );",
		},
		selectInput: "SELECT one FROM tb1 LIMIT 666;",
		expectedOutput: [][]string{
			{"one"},
			{"hello"},
			{"byebye"},
			{"goodbye"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestOffset(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 4, 	22, 'aa'  );",
			"INSERT INTO tb1 VALUES( 'sorry',   2, 	55, 'ba'  );",
		},
		selectInput: "SELECT one FROM tb1 OFFSET 3;",
		expectedOutput: [][]string{
			{"one"},
			{"sorry"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestOffsetThatOverExceedSize(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 4, 	22, 'aa'  );",
			"INSERT INTO tb1 VALUES( 'sorry',   2, 	55, 'ba'  );",
		},
		selectInput: "SELECT one FROM tb1 WHERE TRUE ORDER BY two ASC OFFSET 4;",
		expectedOutput: [][]string{
			{"one"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestOffsetEqualToZero(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE tb1( one TEXT, two INT, three INT, four TEXT );",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO tb1 VALUES( 'hello',	3, 	11, 'q'  );",
			"INSERT INTO tb1 VALUES( 'byebye', 	1, 	33,	'e'  );",
			"INSERT INTO tb1 VALUES( 'goodbye', 4, 	22, 'aa'  );",
			"INSERT INTO tb1 VALUES( 'sorry',   2, 	55, 'ba'  );",
		},
		selectInput: "SELECT one FROM tb1 OFFSET 0;",
		expectedOutput: [][]string{
			{"one"},
			{"hello"},
			{"byebye"},
			{"goodbye"},
			{"sorry"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestLimitAndOffset(t *testing.T) {
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
		selectInput: "SELECT one FROM tb1 WHERE TRUE ORDER BY two ASC, four DESC LIMIT 2 OFFSET 2;",
		expectedOutput: [][]string{
			{"one"},
			{"goodbye"},
			{"hello"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestDefaultJoinToBehaveLikeInnerJoin(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE books( author_id INT, title TEXT);",
			"CREATE TABLE authors( author_id INT, name TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO books VALUES(2, 'Fire');",
			"INSERT INTO books VALUES(1, 'Earth');",
			"INSERT INTO books VALUES(1, 'Air');",
			"INSERT INTO books VALUES(3, 'Smoke');",
			"INSERT INTO authors VALUES( 1, 'Reynold Boyka'  );",
			"INSERT INTO authors VALUES( 2, 'Alissa Ireneus'  );",
			"INSERT INTO authors VALUES( 3, NULL  );",
		},
		selectInput: "SELECT books.title, authors.name FROM books JOIN authors ON books.author_id EQUAL authors.author_id;",
		expectedOutput: [][]string{
			{"books.title", "authors.name"},
			{"Fire", "Alissa Ireneus"},
			{"Earth", "Reynold Boyka"},
			{"Air", "Reynold Boyka"},
			{"Smoke", "NULL"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestInnerJoinOnMultipleMatches(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE books( author_id INT, title TEXT);",
			"CREATE TABLE authors( author_id INT, name TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO books VALUES(1, 'Book One');",
			"INSERT INTO books VALUES(1, 'Book Two');",
			"INSERT INTO authors VALUES(1, 'Author One');",
			"INSERT INTO authors VALUES(1, 'Author Two');",
		},
		selectInput: "SELECT books.title, authors.name FROM books JOIN authors ON books.author_id EQUAL authors.author_id;",
		expectedOutput: [][]string{
			{"books.title", "authors.name"},
			{"Book One", "Author One"},
			{"Book One", "Author Two"},
			{"Book Two", "Author One"},
			{"Book Two", "Author Two"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestFullJoinOnIdenticalTables(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
			"CREATE TABLE table2( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table2 VALUES(2, 'Value2');",
			"INSERT INTO table2 VALUES(3, NULL);",
		},
		selectInput: "SELECT table1.value, table2.value FROM table1 FULL JOIN table2 ON table1.id EQUAL table2.id;",
		expectedOutput: [][]string{
			{"table1.value", "table2.value"},
			{"Value1", "NULL"},
			{"Value2", "Value2"},
			{"NULL", "NULL"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestInnerJoinWithSpecifiedKeywordOnIdenticalTables(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
			"CREATE TABLE table2( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(NULL, NULL);",
			"INSERT INTO table2 VALUES(2, 'Value2');",
			"INSERT INTO table2 VALUES(3, 'Value3');",
			"INSERT INTO table2 VALUES(NULL, 'Value4');",
		},
		selectInput: "SELECT table1.value, table2.value FROM table1 INNER JOIN table2 ON table1.id EQUAL table2.id;",
		expectedOutput: [][]string{
			{"table1.value", "table2.value"},
			{"Value2", "Value2"},
			{"NULL", "Value4"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestLeftJoinOnIdenticalTables(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
			"CREATE TABLE table2( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(NULL, 'Value4');",
			"INSERT INTO table2 VALUES(2, 'Value2');",
			"INSERT INTO table2 VALUES(3, 'Value3');",
			"INSERT INTO table2 VALUES(NULL, NULL);",
		},
		selectInput: "SELECT table1.value, table2.value FROM table1 LEFT JOIN table2 ON table1.id EQUAL table2.id;",
		expectedOutput: [][]string{
			{"table1.value", "table2.value"},
			{"Value1", "NULL"},
			{"Value2", "Value2"},
			{"Value4", "NULL"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestRightJoinOnIdenticalTables(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
			"CREATE TABLE table2( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(NULL, NULL);",
			"INSERT INTO table2 VALUES(2, 'Value2');",
			"INSERT INTO table2 VALUES(3, 'Value3');",
			"INSERT INTO table2 VALUES(NULL, 'Value4');",
		},
		selectInput: "SELECT table1.value, table2.value FROM table1 RIGHT JOIN table2 ON table1.id EQUAL table2.id;",
		expectedOutput: [][]string{
			{"table1.value", "table2.value"},
			{"Value2", "Value2"},
			{"NULL", "Value3"},
			{"NULL", "Value4"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionMax(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
		},
		selectInput: "SELECT MAX(id), MAX(value) FROM table1;",
		expectedOutput: [][]string{
			{"MAX(id)", "MAX(value)"},
			{"2", "Value2"},
		},
		expectedColumnTypes: []token.Token{intColumnType, textColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionMin(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(3, NULL);",
		},
		selectInput: "SELECT MIN(value), MIN(id) FROM table1 WHERE value NOT NULL;",
		expectedOutput: [][]string{
			{"MIN(value)", "MIN(id)"},
			{"Value1", "1"},
		},
		expectedColumnTypes: []token.Token{textColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionMinWithNull(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(3, NULL);",
		},
		selectInput: "SELECT MIN(value), MIN(id) FROM table1;",
		expectedOutput: [][]string{
			{"MIN(value)", "MIN(id)"},
			{"NULL", "1"},
		},
		expectedColumnTypes: []token.Token{textColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionCount(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(3, 'Value3');",
			"INSERT INTO table1 VALUES(NULL, NULL);",
		},
		selectInput: "SELECT COUNT(*), COUNT(id), COUNT(value) FROM table1;",
		expectedOutput: [][]string{
			{"COUNT(*)", "COUNT(id)", "COUNT(value)"},
			{"4", "3", "3"},
		},
		expectedColumnTypes: []token.Token{intColumnType, intColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSum(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(3, 'Value3');",
			"INSERT INTO table1 VALUES(NULL, 'Value4');",
		},
		selectInput: "SELECT SUM(id), SUM(value) FROM table1;",
		expectedOutput: [][]string{
			{"SUM(id)", "SUM(value)"},
			{"6", "0"},
		},
		expectedColumnTypes: []token.Token{intColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSumWithAllNulls(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( score INT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(NULL);",
			"INSERT INTO table1 VALUES(NULL);",
		},
		selectInput: "SELECT SUM(score) FROM table1;",
		expectedOutput: [][]string{
			{"SUM(score)"},
			{"NULL"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionAvg(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(3, 'Value3');",
			"INSERT INTO table1 VALUES(10, NULL);",
		},
		selectInput: "SELECT AVG(id), AVG(value) FROM table1;",
		expectedOutput: [][]string{
			{"AVG(id)", "AVG(value)"},
			{"4", "0"},
		},
		expectedColumnTypes: []token.Token{intColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionAvgWithAllNulls(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( score INT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(NULL);",
			"INSERT INTO table1 VALUES(NULL);",
		},
		selectInput: "SELECT AVG(score) FROM table1;",
		expectedOutput: [][]string{
			{"AVG(score)"},
			{"NULL"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionAvgWithPartialNullsInIntColumn(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE s( name TEXT, score INT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO s VALUES('ann', 10);",
			"INSERT INTO s VALUES('bob', 20);",
			"INSERT INTO s VALUES('cat', NULL);",
			"INSERT INTO s VALUES('dan', NULL);",
		},
		selectInput: "SELECT AVG(score) FROM s;",
		expectedOutput: [][]string{
			{"AVG(score)"},
			{"15"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSumOnTextColumnWithNullInFirstRow(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(NULL);",
			"INSERT INTO table1 VALUES('Value1');",
		},
		selectInput: "SELECT SUM(value) FROM table1;",
		expectedOutput: [][]string{
			{"SUM(value)"},
			{"0"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionAvgOnTextColumnWithNullInFirstRow(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(NULL);",
			"INSERT INTO table1 VALUES('Value1');",
		},
		selectInput: "SELECT AVG(value) FROM table1;",
		expectedOutput: [][]string{
			{"AVG(value)"},
			{"0"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSumOnIntColumnWithNullInFirstRow(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( value INT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(NULL);",
			"INSERT INTO table1 VALUES(4);",
		},
		selectInput: "SELECT SUM(value) FROM table1;",
		expectedOutput: [][]string{
			{"SUM(value)"},
			{"4"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSumOnIntColumnWithUpdatedStringInLaterRow(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value INT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 3);",
			"INSERT INTO table1 VALUES(2, 4);",
			"UPDATE table1 SET value TO '10' WHERE id EQUAL 2;",
		},
		selectInput: "SELECT SUM(value) FROM table1;",
		expectedOutput: [][]string{
			{"SUM(value)"},
			{"13"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSumOnIntColumnWithUpdatedStringInFirstNonNullRow(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value INT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, NULL);",
			"INSERT INTO table1 VALUES(2, 3);",
			"INSERT INTO table1 VALUES(3, 4);",
			"UPDATE table1 SET value TO '10' WHERE id EQUAL 2;",
		},
		selectInput: "SELECT SUM(value) FROM table1;",
		expectedOutput: [][]string{
			{"SUM(value)"},
			{"14"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSumOnTextColumnWithUpdatedNumberInFirstRow(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"UPDATE table1 SET value TO 5 WHERE id EQUAL 1;",
		},
		selectInput: "SELECT SUM(value) FROM table1;",
		expectedOutput: [][]string{
			{"SUM(value)"},
			{"0"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionSumOnTextColumnWithWhereAndOrderBy(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
		},
		selectInput: "SELECT SUM(value), value FROM table1 WHERE TRUE ORDER BY id ASC;",
		expectedOutput: [][]string{
			{"SUM(value)", "value"},
			{"0", "Value1"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionWithColumnSelection(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(3, NULL);",
			"INSERT INTO table1 VALUES(6, 'Value3');",
		},
		selectInput: "SELECT AVG(id), id FROM table1;",
		expectedOutput: [][]string{
			{"AVG(id)", "id"},
			{"3", "1"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionWithColumnSelectionAndOrderBy(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
			"INSERT INTO table1 VALUES(3, 'Value3');",
			"INSERT INTO table1 VALUES(4, NULL);",
		},
		selectInput: "SELECT MAX(id), id FROM table1 ORDER BY id DESC;",
		expectedOutput: [][]string{
			{"MAX(id)", "id"},
			{"4", "4"},
		},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionWithTextColumnSelection(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
			"INSERT INTO table1 VALUES(2, 'Value2');",
		},
		selectInput: "SELECT MAX(id), value FROM table1;",
		expectedOutput: [][]string{
			{"MAX(id)", "value"},
			{"2", "Value1"},
		},
		expectedColumnTypes: []token.Token{intColumnType, textColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionMinAndMaxWithAllNulls(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(NULL, NULL);",
			"INSERT INTO table1 VALUES(NULL, NULL);",
		},
		selectInput: "SELECT MIN(value), MAX(value), MIN(id), MAX(id) FROM table1;",
		expectedOutput: [][]string{
			{"MIN(value)", "MAX(value)", "MIN(id)", "MAX(id)"},
			{"NULL", "NULL", "NULL", "NULL"},
		},
		expectedColumnTypes: []token.Token{textColumnType, textColumnType, intColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionMinAndMaxOnEmptyTable(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		selectInput: "SELECT MIN(value), MAX(value), MIN(id), MAX(id) FROM table1;",
		expectedOutput: [][]string{
			{"MIN(value)", "MAX(value)", "MIN(id)", "MAX(id)"},
			{"NULL", "NULL", "NULL", "NULL"},
		},
		expectedColumnTypes: []token.Token{textColumnType, textColumnType, intColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionMinAndMaxWhenWhereMatchesNoRows(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		insertAndDeleteInputs: []string{
			"INSERT INTO table1 VALUES(1, 'Value1');",
		},
		selectInput: "SELECT MIN(value), MAX(id) FROM table1 WHERE id EQUAL 999;",
		expectedOutput: [][]string{
			{"MIN(value)", "MAX(id)"},
			{"NULL", "NULL"},
		},
		expectedColumnTypes: []token.Token{textColumnType, intColumnType},
	}

	engineTestSuite.runTestSuite(t)
}

func TestAggregateFunctionCountOnEmptyTable(t *testing.T) {
	engineTestSuite := engineTableContentTestSuite{
		createInputs: []string{
			"CREATE TABLE table1( id INT, value TEXT);",
		},
		selectInput: "SELECT COUNT(*) FROM table1;",
		expectedOutput: [][]string{
			{"COUNT(*)"},
			{"0"},
		},
		expectedColumnTypes: []token.Token{intColumnType},
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

var (
	textColumnType = token.Token{Type: token.TEXT, Literal: "TEXT"}
	intColumnType  = token.Token{Type: token.INT, Literal: "INT"}
)

type engineTableContentTestSuite struct {
	createInputs          []string
	insertAndDeleteInputs []string
	selectInput           string
	expectedOutput        [][]string
	expectedColumnTypes   []token.Token
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
	_, err := engine.Evaluate(sequencesWithoutSelect)
	if err != nil {
		log.Fatal(err)
	}
	actualTable, err := engine.getSelectResponse(selectCommand.Commands[0].(*ast.SelectCommand))
	if err != nil {
		log.Fatal(err)
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

	if engineTestSuite.expectedColumnTypes == nil {
		return
	}
	if len(actualTable.Columns) != len(engineTestSuite.expectedColumnTypes) {
		t.Fatalf("Number of column types is incorrect, expecting %d, got %d", len(engineTestSuite.expectedColumnTypes), len(actualTable.Columns))
	}
	for i, column := range actualTable.Columns {
		if column.Type != engineTestSuite.expectedColumnTypes[i] {
			t.Errorf("Column %s has wrong type, expected: %v, got: %v", column.Name, engineTestSuite.expectedColumnTypes[i], column.Type)
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
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		log.Fatal(err)
	}
	return sequences
}
