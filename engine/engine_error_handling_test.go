package engine

import (
	"fmt"
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

// TestEngineUpdateValueTypeErrorHandling - a wrong-type SET value is decided by the statement and
// the schema alone, so it errors the same way on an empty, a no-match and a matching table.
func TestEngineUpdateValueTypeErrorHandling(t *testing.T) {
	intSetup := "CREATE TABLE t(id INT, value INT);"
	textSetup := "CREATE TABLE t(id INT, value TEXT);"
	twoRows := "INSERT INTO t VALUES(1, NULL);INSERT INTO t VALUES(2, NULL);"

	tests := []errorHandlingTestSuite{
		{intSetup + "UPDATE t SET value TO '5';", "invalid value type provided in UPDATE command, expecting: INT, got: '5'"},
		{intSetup + twoRows + "UPDATE t SET value TO '5' WHERE id EQUAL 3;", "invalid value type provided in UPDATE command, expecting: INT, got: '5'"},
		{intSetup + twoRows + "UPDATE t SET value TO '5' WHERE id EQUAL 2;", "invalid value type provided in UPDATE command, expecting: INT, got: '5'"},
		{textSetup + "UPDATE t SET value TO 7;", "invalid value type provided in UPDATE command, expecting: TEXT, got: 7"},
		{textSetup + twoRows + "UPDATE t SET value TO 7 WHERE id EQUAL 3;", "invalid value type provided in UPDATE command, expecting: TEXT, got: 7"},
		{textSetup + twoRows + "UPDATE t SET value TO 7 WHERE id EQUAL 2;", "invalid value type provided in UPDATE command, expecting: TEXT, got: 7"},
		{intSetup + "UPDATE t SET value TO 3.14;", "invalid value type provided in UPDATE command, expecting: INT, got: 3.14"},
		{intSetup + twoRows + "UPDATE t SET value TO 3.14 WHERE id EQUAL 3;", "invalid value type provided in UPDATE command, expecting: INT, got: 3.14"},
		{intSetup + twoRows + "UPDATE t SET value TO 3.14 WHERE id EQUAL 2;", "invalid value type provided in UPDATE command, expecting: INT, got: 3.14"},
		{intSetup + "UPDATE t SET value TO 99999999999999999999;", "invalid value type provided in UPDATE command, expecting: INT, got: 99999999999999999999"},
		{intSetup + twoRows + "UPDATE t SET value TO 99999999999999999999 WHERE id EQUAL 3;", "invalid value type provided in UPDATE command, expecting: INT, got: 99999999999999999999"},
		{intSetup + twoRows + "UPDATE t SET value TO 99999999999999999999 WHERE id EQUAL 2;", "invalid value type provided in UPDATE command, expecting: INT, got: 99999999999999999999"},
		{"CREATE TABLE t(id INT, value INT);INSERT INTO t VALUES(1,3);INSERT INTO t VALUES(2,4);UPDATE t SET value TO 'abc' WHERE id EQUAL 2;", "invalid value type provided in UPDATE command, expecting: INT, got: 'abc'"},
		{"CREATE TABLE t(id INT, value INT);INSERT INTO t VALUES(1,3);INSERT INTO t VALUES(2,4);UPDATE t SET value TO 'abc' WHERE id EQUAL 1;", "invalid value type provided in UPDATE command, expecting: INT, got: 'abc'"},
		{"CREATE TABLE table1( id INT, value INT);INSERT INTO table1 VALUES(1, 3);INSERT INTO table1 VALUES(2, 4);UPDATE table1 SET value TO '10' WHERE id EQUAL 2;", "invalid value type provided in UPDATE command, expecting: INT, got: '10'"},
		{"CREATE TABLE table1( id INT, value INT);INSERT INTO table1 VALUES(1, NULL);INSERT INTO table1 VALUES(2, 3);INSERT INTO table1 VALUES(3, 4);UPDATE table1 SET value TO '10' WHERE id EQUAL 2;", "invalid value type provided in UPDATE command, expecting: INT, got: '10'"},
		{"CREATE TABLE table1( id INT, value TEXT);INSERT INTO table1 VALUES(1, 'Value1');INSERT INTO table1 VALUES(2, 'Value2');UPDATE table1 SET value TO 5 WHERE id EQUAL 1;", "invalid value type provided in UPDATE command, expecting: TEXT, got: 5"},
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
		{"CREATE TABLE tbl(id INT NOT NULL, nickname TEXT); UPDATE tbl SET id TO NULL;", updateNull.Error()},
		{"CREATE TABLE users( id INT NOT NULL, nickname TEXT ); INSERT INTO users VALUES( 1, 'kate' ); UPDATE users SET id TO NULL WHERE id EQUAL 3;", updateNull.Error()},
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

// TestEngineUpdateTypeErrorPrecedesNotNullErrorAcrossMapOrder - every SET value is type-checked
// before any is NOT NULL-checked, so the error category doesn't depend on map iteration order.
func TestEngineUpdateTypeErrorPrecedesNotNullErrorAcrossMapOrder(t *testing.T) {
	want := "invalid value type provided in UPDATE command, expecting: TEXT, got: 5"
	setup := "CREATE TABLE tbl(id INT NOT NULL, nickname TEXT); INSERT INTO tbl VALUES(1, 'a');"

	for i := 0; i < 50; i++ {
		errorMsg := getErrorMessage(t, setup+"UPDATE tbl SET nickname TO 5, id TO NULL;", i)
		if errorMsg != want {
			t.Fatalf("[%d] Was expecting error: \n\t{%s},\n\tbut it was:\n\t{%s}", i, want, errorMsg)
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

// TestEngineFilterErrorsAreDecidedBeforeRowsAreRead - a filter error depends only on the statement
// and the schema, so it must surface whatever rows the scan would see: none, only NULLs, or rows
// that let OR/AND skip the bad side.
func TestEngineFilterErrorsAreDecidedBeforeRowsAreRead(t *testing.T) {
	setup := "CREATE TABLE tbl(one TEXT, two INT);"
	rowsByContext := []struct {
		name string
		rows string
		wrap string
	}{
		{"empty table", "", "%s"},
		{"all-NULL column", "INSERT INTO tbl VALUES(NULL, NULL);", "%s"},
		{"skipped OR side", "INSERT INTO tbl VALUES('a', 1);", "one EQUAL 'a' OR %s"},
		{"skipped AND side", "INSERT INTO tbl VALUES('a', 1);", "one EQUAL 'b' AND %s"},
	}
	badFilters := []struct {
		filter        string
		expectedError string
	}{
		{"two EQUAL 'abc'", "invalid value type provided in WHERE command, expecting: INT, got: 'abc'"},
		{"two NOT 3.14", "invalid value type provided in WHERE command, expecting: INT, got: 3.14"},
		{"one NOT 5", "invalid value type provided in WHERE command, expecting: TEXT, got: 5"},
		{"'5' EQUAL two", "invalid value type provided in WHERE command, expecting: INT, got: '5'"},
		{"one EQUAL two", "invalid value type provided in WHERE command, expecting: TEXT, got: INT"},
		{"'x' EQUAL 5", "invalid value type provided in WHERE command, expecting: TEXT, got: 5"},
		{"two IN (1, 'x')", "invalid value type provided in WHERE command, expecting: INT, got: 'x'"},
		{"two NOTIN (3.14)", "invalid value type provided in WHERE command, expecting: INT, got: 3.14"},
		{"one IN ('a', 5)", "invalid value type provided in WHERE command, expecting: TEXT, got: 5"},
		{"two LIKE '1%'", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"two NOT LIKE '1%'", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"two LIKE NULL", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"'5' LIKE two", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"one LIKE 5", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"nope EQUAL 1", "column with the name of nope doesn't exist in table tbl"},
	}

	var tests []errorHandlingTestSuite
	for _, context := range rowsByContext {
		for _, bad := range badFilters {
			filter := fmt.Sprintf(context.wrap, bad.filter)
			tests = append(tests,
				errorHandlingTestSuite{setup + context.rows + "SELECT * FROM tbl WHERE " + filter + ";", bad.expectedError},
				errorHandlingTestSuite{setup + context.rows + "DELETE FROM tbl WHERE " + filter + ";", bad.expectedError},
			)
		}
	}

	runEngineErrorHandlingSuite(t, tests)
}

// TestEngineJoinWhereErrorsAreDecidedBeforeRowsAreRead - WHERE on a join is checked against the
// joined table's prefixed columns, and names the SELECT's table on a missing column.
func TestEngineJoinWhereErrorsAreDecidedBeforeRowsAreRead(t *testing.T) {
	setup := "CREATE TABLE a(id INT); CREATE TABLE b(id INT, name TEXT);"
	join := "SELECT * FROM a JOIN b ON a.id EQUAL b.id WHERE "

	tests := []errorHandlingTestSuite{
		{setup + join + "b.name EQUAL 5;", "invalid value type provided in WHERE command, expecting: TEXT, got: 5"},
		{setup + join + "a.id LIKE 'x%';", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{setup + join + "a.id EQUAL b.name;", "invalid value type provided in WHERE command, expecting: INT, got: TEXT"},
		{setup + join + "b.nope EQUAL 1;", "column with the name of b.nope doesn't exist in table a"},
	}

	runEngineErrorHandlingSuite(t, tests)
}

// TestEngineFilterMissingColumnPrecedesTypeError - missing columns are reported before any type
// error, so a misspelt column never surfaces as a confusing type mismatch.
func TestEngineFilterMissingColumnPrecedesTypeError(t *testing.T) {
	setup := "CREATE TABLE tbl(one TEXT, two INT);"
	missing := "column with the name of nope doesn't exist in table tbl"

	tests := []errorHandlingTestSuite{
		{setup + "SELECT * FROM tbl WHERE two EQUAL 'abc' OR nope EQUAL 1;", missing},
		{setup + "DELETE FROM tbl WHERE two LIKE '1%' AND nope EQUAL 1;", missing},
	}

	runEngineErrorHandlingSuite(t, tests)
}

// TestEngineWhereOnProjectedColumnIsMissingWhenRowsExist - pins an accepted exception:
// ORDER BY projects before WHERE runs, so a valid column outside the projection is
// reported missing once rows reach the WHERE (#65).
func TestEngineWhereOnProjectedColumnIsMissingWhenRowsExist(t *testing.T) {
	tests := []errorHandlingTestSuite{
		{"CREATE TABLE tbl(one TEXT, two INT); INSERT INTO tbl VALUES('a', 1); SELECT one FROM tbl WHERE two EQUAL 1 ORDER BY one ASC;", "column with the name of two doesn't exist in table tbl"},
	}

	runEngineErrorHandlingSuite(t, tests)
}

// TestEngineUpdateFilterErrorsAreDecidedBeforeRowsAreRead - UPDATE's WHERE error depends only on
// the statement and the schema, so it must surface whether the scan would see no rows, no
// matching rows, or rows that let OR/AND skip the bad side.
func TestEngineUpdateFilterErrorsAreDecidedBeforeRowsAreRead(t *testing.T) {
	setup := "CREATE TABLE tbl(one TEXT, two INT);"
	rowsByContext := []struct {
		name string
		rows string
		wrap string
	}{
		{"empty table", "", "%s"},
		{"no-match table", "INSERT INTO tbl VALUES('a', 1);", "one EQUAL 'b' AND %s"},
		{"skipped OR side", "INSERT INTO tbl VALUES('a', 1);", "one EQUAL 'a' OR %s"},
	}
	badFilters := []struct {
		filter        string
		expectedError string
	}{
		{"two EQUAL 'abc'", "invalid value type provided in WHERE command, expecting: INT, got: 'abc'"},
		{"two NOTIN (3.14)", "invalid value type provided in WHERE command, expecting: INT, got: 3.14"},
		{"one EQUAL two", "invalid value type provided in WHERE command, expecting: TEXT, got: INT"},
		{"two LIKE '1%'", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"one NOT LIKE 5", "invalid operand type provided in WHERE command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"nope EQUAL 1", "column with the name of nope doesn't exist in table tbl"},
	}

	var tests []errorHandlingTestSuite
	for _, context := range rowsByContext {
		for _, bad := range badFilters {
			filter := fmt.Sprintf(context.wrap, bad.filter)
			tests = append(tests, errorHandlingTestSuite{setup + context.rows + "UPDATE tbl SET one TO 'z' WHERE " + filter + ";", bad.expectedError})
		}
	}

	runEngineErrorHandlingSuite(t, tests)
}

// TestEngineJoinOnErrorsAreDecidedBeforeRowsAreRead - JOIN ON is checked against the joined
// table's prefixed columns, so it errors even when either side has no rows. Every join kind
// reports JOIN, and a missing column keeps the empty table name pinned by the JOIN suite above.
func TestEngineJoinOnErrorsAreDecidedBeforeRowsAreRead(t *testing.T) {
	create := "CREATE TABLE a(id INT); CREATE TABLE b(id INT, name TEXT);"
	rowsByContext := []struct {
		name string
		rows string
	}{
		{"both empty", ""},
		{"left empty", "INSERT INTO b VALUES(1, 'x');"},
		{"right empty", "INSERT INTO a VALUES(1);"},
		{"both populated", "INSERT INTO a VALUES(1); INSERT INTO b VALUES(1, 'x');"},
	}
	joinKinds := []string{"JOIN", "INNER JOIN", "LEFT JOIN", "RIGHT JOIN", "FULL JOIN"}
	badFilters := []struct {
		filter        string
		expectedError string
	}{
		{"a.id EQUAL b.name", "invalid value type provided in JOIN command, expecting: INT, got: TEXT"},
		{"a.id LIKE b.name", "invalid operand type provided in JOIN command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"b.name LIKE a.id", "invalid operand type provided in JOIN command for LIKE, expecting: TEXT or NULL, got: INT"},
		{"a.nope EQUAL b.id", "column with the name of a.nope doesn't exist in table "},
	}

	var tests []errorHandlingTestSuite
	for _, context := range rowsByContext {
		for _, joinKind := range joinKinds {
			for _, bad := range badFilters {
				tests = append(tests, errorHandlingTestSuite{create + context.rows + "SELECT * FROM a " + joinKind + " b ON " + bad.filter + ";", bad.expectedError})
			}
		}
	}

	runEngineErrorHandlingSuite(t, tests)
}
