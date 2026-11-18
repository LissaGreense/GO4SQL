package engine

import (
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
	"github.com/LissaGreense/GO4SQL/token"
	"testing"
)

type errorHandlingTestSuite struct {
	input         string
	expectedError string
}

func TestEngineCreateCommandErrorHandling(t *testing.T) {
	duplicateTableNameError := TableAlreadyExistsError{"table1"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE table1( one TEXT , two INT);CREATE TABLE table1(two INT);", duplicateTableNameError.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineDropCommandErrorHandling(t *testing.T) {
	tableDoNotExistError := TableDoesNotExistError{"table1"}
	missingMidSequenceError := TableDoesNotExistError{"nope"}

	tests := []errorHandlingTestSuite{
		{"DROP TABLE table1;", tableDoNotExistError.Error()},
		{"CREATE TABLE a( x INT ); DROP TABLE nope; CREATE TABLE b( x INT );", missingMidSequenceError.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineInsertCommandErrorHandling(t *testing.T) {
	tableDoNotExistError := TableDoesNotExistError{"table1"}
	invalidNumberOfParametersError := InvalidNumberOfParametersError{expectedNumber: 2, actualNumber: 1, commandName: token.INSERT}
	invalidParametersTypeError := InvalidValueTypeError{expectedType: token.TEXT, actualValue: "1", commandName: token.INSERT}
	tests := []errorHandlingTestSuite{
		{"INSERT INTO table1 VALUES( 'hello', 1);", tableDoNotExistError.Error()},
		{"CREATE TABLE table1( one TEXT , two INT); INSERT INTO table1 VALUES(1);", invalidNumberOfParametersError.Error()},
		{"CREATE TABLE table1( one TEXT , two INT); INSERT INTO table1 VALUES(1, 1 );", invalidParametersTypeError.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineInsertValueTypeErrorWording(t *testing.T) {
	tests := []errorHandlingTestSuite{
		{"CREATE TABLE t( a INT ); INSERT INTO t VALUES( 3.14 );", "invalid value type provided in INSERT command, expecting: INT, got: 3.14"},
		{"CREATE TABLE t( a INT ); INSERT INTO t VALUES( 99999999999999999999 );", "invalid value type provided in INSERT command, expecting: INT, got: 99999999999999999999"},
		{"CREATE TABLE t( a INT ); INSERT INTO t VALUES( 'x' );", "invalid value type provided in INSERT command, expecting: INT, got: 'x'"},
		{"CREATE TABLE t( a TEXT ); INSERT INTO t VALUES( 5 );", "invalid value type provided in INSERT command, expecting: TEXT, got: 5"},
		{"CREATE TABLE t( a INT NOT NULL, b INT ); INSERT INTO t VALUES( NULL, 3.14 );", "NULL value isn't allowed in INSERT command for NOT NULL column: a"},
		{"CREATE TABLE t( a INT, b INT NOT NULL ); INSERT INTO t VALUES( 3.14, NULL );", "invalid value type provided in INSERT command, expecting: INT, got: 3.14"},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineCreateDuplicateColumnErrorHandling(t *testing.T) {
	tests := []errorHandlingTestSuite{
		{"CREATE TABLE t( a INT, a TEXT );", "column with the name of a is declared more than once in table t"},
		{"CREATE TABLE t( a INT NOT NULL, a INT );", "column with the name of a is declared more than once in table t"},
		{"CREATE TABLE IF NOT EXISTS t2( a INT, a TEXT );", "column with the name of a is declared more than once in table t2"},
		{"CREATE TABLE t( a INT ); CREATE TABLE t( a INT, a TEXT );", "table with the name of t already exists"},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineSelectCommandErrorHandling(t *testing.T) {
	noTableDoesNotExist := TableDoesNotExistError{"tb1"}
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); SELECT * FROM tb1;", noTableDoesNotExist.Error()},
		{"CREATE TABLE tbl(one TEXT); SELECT two FROM tbl;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineDeleteCommandErrorHandling(t *testing.T) {
	noTableDoesNotExist := TableDoesNotExistError{"tb1"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); DELETE FROM tb1 WHERE one EQUAL 3;", noTableDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineWhereCommandErrorHandling(t *testing.T) {
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); INSERT INTO tbl VALUES('hello'); SELECT * FROM tbl WHERE two EQUAL 3;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineLikeCommandErrorHandling(t *testing.T) {
	invalidLikeOperand := InvalidLikeOperandError{actualType: token.INT, commandName: token.WHERE}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT, two INT); INSERT INTO tbl VALUES('hello', 5); SELECT * FROM tbl WHERE two LIKE '5%';", invalidLikeOperand.Error()},
		{"CREATE TABLE tbl(one TEXT, two INT); INSERT INTO tbl VALUES('hello', 5); SELECT * FROM tbl WHERE one LIKE 5;", invalidLikeOperand.Error()},
		{"CREATE TABLE tbl(one TEXT, two INT); INSERT INTO tbl VALUES('hello', 5); SELECT * FROM tbl WHERE two NOT LIKE '5%';", invalidLikeOperand.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineUpdateCommandErrorHandling(t *testing.T) {
	noTableDoesNotExist := TableDoesNotExistError{"tb1"}
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); UPDATE tb1 SET one TO 2;", noTableDoesNotExist.Error()},
		{"CREATE TABLE tbl(one TEXT);UPDATE tbl SET two TO 2;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineNotNullConstraintErrorHandling(t *testing.T) {
	insertNull := NullValueNotAllowedError{columnName: "id", commandName: "INSERT"}
	updateNull := NullValueNotAllowedError{columnName: "id", commandName: "UPDATE"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(id INT NOT NULL, nickname TEXT); INSERT INTO tbl VALUES(NULL, 'a');", insertNull.Error()},
		{"CREATE TABLE tbl(nickname TEXT, id INT NOT NULL); INSERT INTO tbl VALUES('a', NULL);", insertNull.Error()},
		{"CREATE TABLE tbl(id INT NOT NULL, nickname TEXT); INSERT INTO tbl VALUES(1, 'a'); UPDATE tbl SET id TO NULL;", updateNull.Error()},
		{"CREATE TABLE tbl(id INT NOT NULL, nickname TEXT); INSERT INTO tbl VALUES(1, 'a'); UPDATE tbl SET nickname TO 'b', id TO NULL WHERE nickname EQUAL 'a';", updateNull.Error()},
		// DELETE rebuilds the table via getCopyOfTableWithoutRows; NOT NULL must survive even a no-op DELETE.
		{"CREATE TABLE tbl(id INT NOT NULL, nickname TEXT); INSERT INTO tbl VALUES(1, 'a'); DELETE FROM tbl WHERE id EQUAL 99; UPDATE tbl SET id TO NULL;", updateNull.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

// TestEngineUpdateSetColumnErrorPrecedesWhereColumnError - a nonexistent SET column must be
// reported before a nonexistent WHERE column, matching the precedence that existed before NOT
// NULL validation was added (SET columns are fully resolved before WHERE is evaluated).
func TestEngineUpdateSetColumnErrorPrecedesWhereColumnError(t *testing.T) {
	setColumnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "nosuch"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(id INT); UPDATE tbl SET nosuch TO 1 WHERE alsomissing EQUAL 1;", setColumnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

// TestEngineUpdateSetColumnErrorIsDeterministicAcrossMapOrder - command.Changes is a map, so
// iteration order is randomized per run. A nonexistent SET column must always be reported, never
// masked by an unrelated NOT NULL violation elsewhere in the same SET clause, regardless of which
// entry the map visits first. Runs many times since the bug this guards against is intermittent.
func TestEngineUpdateSetColumnErrorIsDeterministicAcrossMapOrder(t *testing.T) {
	setColumnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "nosuch"}
	setup := "CREATE TABLE tbl(id INT NOT NULL, nickname TEXT); INSERT INTO tbl VALUES(1, 'a');"

	for i := 0; i < 50; i++ {
		errorMsg := getErrorMessage(t, setup+"UPDATE tbl SET nosuch TO 1, id TO NULL;", i)
		if errorMsg != setColumnDoesNotExist.Error() {
			t.Fatalf("[%d] Was expecting error: \n\t{%s},\n\tbut it was:\n\t{%s}", i, setColumnDoesNotExist.Error(), errorMsg)
		}
	}
}

func TestEngineOrderByCommandErrorHandling(t *testing.T) {
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "tbl", columnName: "two"}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT); SELECT * FROM tbl ORDER BY two ASC;", columnDoesNotExist.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineFullJoinErrorHandling(t *testing.T) {
	leftTableNotExist := TableDoesNotExistError{tableName: "leftTable"}
	rightTableNotExist := TableDoesNotExistError{tableName: "rightTable"}
	columnDoesNotExist := ColumnDoesNotExistError{tableName: "", columnName: "leftTable.two"}
	invalidLikeOperandInJoin := InvalidLikeOperandError{actualType: token.INT, commandName: token.JOIN}

	tests := []errorHandlingTestSuite{
		{"CREATE TABLE rightTable(one TEXT); SELECT leftTable.one, rightTable.one FROM leftTable JOIN rightTable ON leftTable.one EQUAL rightTable.one;", leftTableNotExist.Error()},
		{"CREATE TABLE leftTable(one TEXT); SELECT leftTable.one, rightTable.one FROM leftTable JOIN rightTable ON leftTable.one EQUAL rightTable.one;", rightTableNotExist.Error()},
		{"CREATE TABLE leftTable(one TEXT); CREATE TABLE rightTable(one TEXT); INSERT INTO leftTable VALUES('hi'); INSERT INTO rightTable VALUES('hi'); SELECT * FROM leftTable JOIN rightTable ON leftTable.two EQUAL rightTable.one;", columnDoesNotExist.Error()},
		{"CREATE TABLE leftTable(id INT); CREATE TABLE rightTable(id INT); INSERT INTO leftTable VALUES(1); INSERT INTO rightTable VALUES(1); SELECT * FROM leftTable JOIN rightTable ON leftTable.id LIKE rightTable.id;", invalidLikeOperandInJoin.Error()},
	}

	runEngineErrorHandlingSuite(t, tests)
}

func TestEngineAggregateOnIntColumnWithUpdatedNonNumericString(t *testing.T) {
	setup := "CREATE TABLE t(id INT, value INT);INSERT INTO t VALUES(1,3);INSERT INTO t VALUES(2,4);UPDATE t SET value TO 'abc' WHERE id EQUAL 2;"
	want := `strconv.Atoi: parsing "abc": invalid syntax`
	runEngineErrorHandlingSuite(t, []errorHandlingTestSuite{
		{setup + "SELECT SUM(value) FROM t;", want},
		{setup + "SELECT AVG(value) FROM t;", want},
		{"CREATE TABLE t(id INT, value INT);INSERT INTO t VALUES(1,3);INSERT INTO t VALUES(2,4);UPDATE t SET value TO 'abc' WHERE id EQUAL 1;SELECT SUM(value) FROM t;", want},
	})
}

func runEngineErrorHandlingSuite(t *testing.T, suite []errorHandlingTestSuite) {
	for i, test := range suite {
		errorMsg := getErrorMessage(t, test.input, i)

		if errorMsg != test.expectedError {
			t.Fatalf("[%v]Was expecting error: \n\t{%s},\n\tbut it was:\n\t{%s}", i, test.expectedError, errorMsg)
		}
	}
}

func getErrorMessage(t *testing.T, input string, testIndex int) string {
	lexerInstance := lexer.RunLexer(input)
	parserInstance := parser.New(lexerInstance)
	sequences, parserError := parserInstance.ParseSequence()
	if parserError != nil {
		t.Fatalf("[%d] Error has occurred in parser not in engine, error: %s", testIndex, parserError.Error())
	}

	engine := New()
	_, engineError := engine.Evaluate(sequences)
	if engineError == nil {
		t.Fatalf("[%d] Was expecting error from engine but there was none", testIndex)
	}

	return engineError.Error()
}
