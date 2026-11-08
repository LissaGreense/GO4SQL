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

func TestAggregateResultColumnType(t *testing.T) {
	createInput := "CREATE TABLE s( name TEXT, score INT );"
	textType := token.Token{Type: token.TEXT, Literal: "TEXT"}
	intType := token.Token{Type: token.INT, Literal: "INT"}
	rows := []string{
		"INSERT INTO s VALUES( 'dan', 10 );",
		"INSERT INTO s VALUES( 'amy', 30 );",
	}
	nullRows := []string{
		"INSERT INTO s VALUES( NULL, NULL );",
		"INSERT INTO s VALUES( NULL, NULL );",
	}

	testCases := []struct {
		name          string
		inserts       []string
		selectInput   string
		expectedTypes []token.Token
	}{
		{
			name:          "MIN on TEXT",
			inserts:       rows,
			selectInput:   "SELECT MIN(name) FROM s;",
			expectedTypes: []token.Token{textType},
		},
		{
			name:          "MAX on TEXT",
			inserts:       rows,
			selectInput:   "SELECT MAX(name) FROM s;",
			expectedTypes: []token.Token{textType},
		},
		{
			name:          "MIN on INT",
			inserts:       rows,
			selectInput:   "SELECT MIN(score) FROM s;",
			expectedTypes: []token.Token{intType},
		},
		{
			name:          "MAX on INT",
			inserts:       rows,
			selectInput:   "SELECT MAX(score) FROM s;",
			expectedTypes: []token.Token{intType},
		},
		{
			name:          "plain column next to an aggregate",
			inserts:       rows,
			selectInput:   "SELECT MAX(score), name FROM s;",
			expectedTypes: []token.Token{intType, textType},
		},
		{
			name:          "MIN and MAX on all NULL values",
			inserts:       nullRows,
			selectInput:   "SELECT MIN(name), MAX(name), MIN(score), MAX(score) FROM s;",
			expectedTypes: []token.Token{textType, textType, intType, intType},
		},
		{
			name:          "MIN and MAX on an empty table",
			selectInput:   "SELECT MIN(name), MAX(name), MIN(score), MAX(score) FROM s;",
			expectedTypes: []token.Token{textType, textType, intType, intType},
		},
		{
			name:          "MIN and MAX when WHERE matches no rows",
			inserts:       rows,
			selectInput:   "SELECT MIN(name), MAX(score) FROM s WHERE score EQUAL 999;",
			expectedTypes: []token.Token{textType, intType},
		},
		{
			name:          "COUNT, SUM and AVG stay INT",
			inserts:       rows,
			selectInput:   "SELECT COUNT(*), COUNT(name), SUM(name), AVG(name), SUM(score), AVG(score) FROM s;",
			expectedTypes: []token.Token{intType, intType, intType, intType, intType, intType},
		},
		{
			name:          "COUNT(*) on an empty table stays INT",
			selectInput:   "SELECT COUNT(*) FROM s;",
			expectedTypes: []token.Token{intType},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			engine := New()
			setup := inputsToString(append([]string{createInput}, testCase.inserts...))
			if _, err := engine.Evaluate(getSequences(setup)); err != nil {
				t.Fatal(err)
			}

			selectCommand := getSequences(testCase.selectInput).Commands[0].(*ast.SelectCommand)
			actualTable, err := engine.getSelectResponse(selectCommand)
			if err != nil {
				t.Fatal(err)
			}

			if len(actualTable.Columns) != len(testCase.expectedTypes) {
				t.Fatalf("expected %d columns, got %d", len(testCase.expectedTypes), len(actualTable.Columns))
			}
			for i, column := range actualTable.Columns {
				if column.Type != testCase.expectedTypes[i] {
					t.Errorf("column %s: expected type %v, got %v", column.Name, testCase.expectedTypes[i], column.Type)
				}
			}
		})
	}
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
