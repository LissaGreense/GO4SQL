package lexer

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/token"
)

func TestLexerWithInsertCommand(t *testing.T) {
	input :=
		`
			CREATE TABLE 	1tbl( one TEXT , two INT );
			INSERT INTO tbl VALUES( 'CREATE',	 10 );
			INSERT 	INTO tbl  VALUES( 'goodbye', 20 );
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "1tbl"},
		{token.LPAREN, "("},
		{token.IDENT, "one"},
		{token.TEXT, "TEXT"},
		{token.COMMA, ","},
		{token.IDENT, "two"},
		{token.INT, "INT"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "CREATE"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "goodbye"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestLexerWithUpdateCommand(t *testing.T) {
	input :=
		`
			UPDATE table1
			SET column_name_1 TO 'UPDATE', column_name_2 TO 42
			WHERE column_name_3 EQUAL 1;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.UPDATE, "UPDATE"},
		{token.IDENT, "table1"},
		{token.SET, "SET"},
		{token.IDENT, "column_name_1"},
		{token.TO, "TO"},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "UPDATE"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.IDENT, "column_name_2"},
		{token.TO, "TO"},
		{token.LITERAL, "42"},
		{token.WHERE, "WHERE"},
		{token.IDENT, "column_name_3"},
		{token.EQUAL, "EQUAL"},
		{token.LITERAL, "1"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestLexerWithNumbersMixedInLiterals(t *testing.T) {
	input :=
		`
			CREATE TABLE 	tbl2( one TEXT , two INT );
			INSERT INTO tbl2 VALUES( 'hello1',	 10 );
			INSERT 	INTO tbl2  VALUES( 'good123bye', 20 );
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "tbl2"},
		{token.LPAREN, "("},
		{token.IDENT, "one"},
		{token.TEXT, "TEXT"},
		{token.COMMA, ","},
		{token.IDENT, "two"},
		{token.INT, "INT"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl2"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "hello1"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl2"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "good123bye"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestLexerWithNumbersWithWhitespacesIdentifier(t *testing.T) {
	input :=
		`
			CREATE TABLE 	tbl2( one TEXT , two INT );
			INSERT INTO tbl2 VALUES( ' hello1',	 10 );
			INSERT 	INTO tbl2  VALUES( 'Hello	 Dear', 20 );
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "tbl2"},
		{token.LPAREN, "("},
		{token.IDENT, "one"},
		{token.TEXT, "TEXT"},
		{token.COMMA, ","},
		{token.IDENT, "two"},
		{token.INT, "INT"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl2"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, " hello1"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "10"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "tbl2"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "Hello\t Dear"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.LITERAL, "20"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestLogicalStatements(t *testing.T) {
	input :=
		`
			WHERE FALSE AND three EQUAL 33;
			WHERE two NOT 11 OR TRUE;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.WHERE, "WHERE"},
		{token.FALSE, "FALSE"},
		{token.AND, "AND"},
		{token.IDENT, "three"},
		{token.EQUAL, "EQUAL"},
		{token.LITERAL, "33"},
		{token.SEMICOLON, ";"},
		{token.WHERE, "WHERE"},
		{token.IDENT, "two"},
		{token.NOT, "NOT"},
		{token.LITERAL, "11"},
		{token.OR, "OR"},
		{token.TRUE, "TRUE"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestInStatement(t *testing.T) {
	input :=
		`
			WHERE two IN (1, 2) AND
			WHERE three NOTIN ('one', 'two');
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.WHERE, "WHERE"},
		{token.IDENT, "two"},
		{token.IN, "IN"},
		{token.LPAREN, "("},
		{token.LITERAL, "1"},
		{token.COMMA, ","},
		{token.LITERAL, "2"},
		{token.RPAREN, ")"},
		{token.AND, "AND"},
		{token.WHERE, "WHERE"},
		{token.IDENT, "three"},
		{token.NOTIN, "NOTIN"},
		{token.LPAREN, "("},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "one"},
		{token.APOSTROPHE, "'"},
		{token.COMMA, ","},
		{token.APOSTROPHE, "'"},
		{token.IDENT, "two"},
		{token.APOSTROPHE, "'"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestDeleteStatement(t *testing.T) {
	input := `DELETE FROM table WHERE two NOT 11 OR TRUE;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DELETE, "DELETE"},
		{token.FROM, "FROM"},
		{token.IDENT, "table"},
		{token.WHERE, "WHERE"},
		{token.IDENT, "two"},
		{token.NOT, "NOT"},
		{token.LITERAL, "11"},
		{token.OR, "OR"},
		{token.TRUE, "TRUE"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestOrderByStatement(t *testing.T) {
	input := `SELECT * FROM table ORDER BY something ASC;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.ASTERISK, "*"},
		{token.FROM, "FROM"},
		{token.IDENT, "table"},
		{token.ORDER, "ORDER"},
		{token.BY, "BY"},
		{token.IDENT, "something"},
		{token.ASC, "ASC"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestDropStatement(t *testing.T) {
	input := `DROP TABLE table;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DROP, "DROP"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "table"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestLimitAndOffsetStatement(t *testing.T) {
	input := `LIMIT 5 OFFSET 6;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LIMIT, "LIMIT"},
		{token.LITERAL, "5"},
		{token.OFFSET, "OFFSET"},
		{token.LITERAL, "6"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestAggregateFunctions(t *testing.T) {
	input := `SELECT MIN(colOne), MAX(colOne), COUNT(colOne), SUM(colOne), AVG(colOne) FROM tbl;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.MIN, "MIN"},
		{token.LPAREN, "("},
		{token.IDENT, "colOne"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.MAX, "MAX"},
		{token.LPAREN, "("},
		{token.IDENT, "colOne"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.COUNT, "COUNT"},
		{token.LPAREN, "("},
		{token.IDENT, "colOne"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.SUM, "SUM"},
		{token.LPAREN, "("},
		{token.IDENT, "colOne"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.AVG, "AVG"},
		{token.LPAREN, "("},
		{token.IDENT, "colOne"},
		{token.RPAREN, ")"},
		{token.FROM, "FROM"},
		{token.IDENT, "tbl"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestSelectWithDistinct(t *testing.T) {
	input := `SELECT DISTINCT * FROM table;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.DISTINCT, "DISTINCT"},
		{token.ASTERISK, "*"},
		{token.FROM, "FROM"},
		{token.IDENT, "table"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestDefaultJoin(t *testing.T) {
	input := `	SELECT title FROM books
    			JOIN authors ON
        		books.author_id EQUAL authors.author_id;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.IDENT, "title"},
		{token.FROM, "FROM"},
		{token.IDENT, "books"},
		{token.JOIN, "JOIN"},
		{token.IDENT, "authors"},
		{token.ON, "ON"},
		{token.IDENT, "books.author_id"},
		{token.EQUAL, "EQUAL"},
		{token.IDENT, "authors.author_id"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestInnerJoin(t *testing.T) {
	input := `	SELECT title FROM books
    			INNER JOIN authors ON
        		books.author_id EQUAL authors.author_id;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.IDENT, "title"},
		{token.FROM, "FROM"},
		{token.IDENT, "books"},
		{token.INNER, "INNER"},
		{token.JOIN, "JOIN"},
		{token.IDENT, "authors"},
		{token.ON, "ON"},
		{token.IDENT, "books.author_id"},
		{token.EQUAL, "EQUAL"},
		{token.IDENT, "authors.author_id"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestLeftJoin(t *testing.T) {
	input := `	SELECT title FROM books
    			LEFT JOIN authors ON
        		books.author_id EQUAL authors.author_id;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.IDENT, "title"},
		{token.FROM, "FROM"},
		{token.IDENT, "books"},
		{token.LEFT, "LEFT"},
		{token.JOIN, "JOIN"},
		{token.IDENT, "authors"},
		{token.ON, "ON"},
		{token.IDENT, "books.author_id"},
		{token.EQUAL, "EQUAL"},
		{token.IDENT, "authors.author_id"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestRightJoin(t *testing.T) {
	input := `	SELECT title FROM books
    			RIGHT JOIN authors ON
        		books.author_id EQUAL authors.author_id;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.IDENT, "title"},
		{token.FROM, "FROM"},
		{token.IDENT, "books"},
		{token.RIGHT, "RIGHT"},
		{token.JOIN, "JOIN"},
		{token.IDENT, "authors"},
		{token.ON, "ON"},
		{token.IDENT, "books.author_id"},
		{token.EQUAL, "EQUAL"},
		{token.IDENT, "authors.author_id"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func TestFullJoin(t *testing.T) {
	input := `	SELECT title FROM books
    			FULL JOIN authors ON
        		books.author_id EQUAL authors.author_id;
			`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.IDENT, "title"},
		{token.FROM, "FROM"},
		{token.IDENT, "books"},
		{token.FULL, "FULL"},
		{token.JOIN, "JOIN"},
		{token.IDENT, "authors"},
		{token.ON, "ON"},
		{token.IDENT, "books.author_id"},
		{token.EQUAL, "EQUAL"},
		{token.IDENT, "authors.author_id"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runLexerTestSuite(t, input, tests)
}

func runLexerTestSuite(t *testing.T, input string, tests []struct {
	expectedType    token.Type
	expectedLiteral string
}) {
	l := RunLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
