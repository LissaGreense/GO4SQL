package parser

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/token"
)

func TestParserCreateCommand(t *testing.T) {
	tests := []struct {
		input               string
		expectedTableName   string
		expectedColumnNames []string
		expectedColumTypes  []token.Token
	}{
		{"CREATE TABLE 	TBL( ONE TEXT );", "TBL", []string{"ONE"}, []token.Token{{Type: token.TEXT, Literal: "TEXT"}}},
		{"CREATE TABLE 	TBL( ONE TEXT,  TWO TEXT, THREE INT);", "TBL", []string{"ONE", "TWO", "THREE"}, []token.Token{{Type: token.TEXT, Literal: "TEXT"}, {Type: token.TEXT, Literal: "TEXT"}, {Type: token.INT, Literal: "INT"}}},
		{"CREATE TABLE 	TBL(  );", "TBL", []string{}, []token.Token{}},
	}

	for testIndex, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences, err := parserInstance.ParseSequence()
		if err != nil {
			t.Fatalf("[%d] Got error from parser: %s", testIndex, err)
		}

		if len(sequences.Commands) != 1 {
			t.Fatalf("[%d] sequences does not contain 1 statements. got=%d", testIndex, len(sequences.Commands))
		}

		if !testCreateStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedColumnNames, tt.expectedColumTypes) {
			return
		}
	}
}

func testCreateStatement(t *testing.T, command ast.Command, expectedTableName string, expectedColumnNames []string, expectedColumTypes []token.Token) bool {
	if command.TokenLiteral() != "CREATE" {
		t.Errorf("command.TokenLiteral() not 'CREATE'. got=%q", command.TokenLiteral())
		return false
	}

	actualCreateCommand, ok := command.(*ast.CreateCommand)
	if !ok {
		t.Errorf("actualCreateCommand is not %T. got=%T", &ast.CreateCommand{}, command)
		return false
	}

	if actualCreateCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualCreateCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !stringArrayEquals(actualCreateCommand.ColumnNames, expectedColumnNames) {
		t.Errorf("")
		return false
	}

	if !tokenArrayEquals(actualCreateCommand.ColumnTypes, expectedColumTypes) {
		t.Errorf("")
		return false
	}

	return true
}

func TestParseInsertCommand(t *testing.T) {
	tests := []struct {
		input                string
		expectedTableName    string
		expectedValuesTokens []token.Token
	}{
		{"INSERT INTO TBL VALUES();", "TBL", []token.Token{}},
		{"INSERT INTO TBL VALUES( 'HELLO' );", "TBL", []token.Token{{Type: token.IDENT, Literal: "HELLO"}}},
		{"INSERT INTO TBL VALUES( 'HELLO',	 10 , 'LOL');", "TBL", []token.Token{{Type: token.IDENT, Literal: "HELLO"}, {Type: token.LITERAL, Literal: "10"}, {Type: token.IDENT, Literal: "LOL"}}},
		{"INSERT INTO TBL VALUES(NULL, 'NULL', null);", "TBL", []token.Token{{Type: token.NULL, Literal: "NULL"}, {Type: token.IDENT, Literal: "NULL"}, {Type: token.IDENT, Literal: "null"}}},
	}

	for testIndex, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences, err := parserInstance.ParseSequence()
		if err != nil {
			t.Fatalf("[%d] Got error from parser: %s", testIndex, err)
		}

		if len(sequences.Commands) != 1 {
			t.Fatalf("[%d] sequences does not contain 1 statements. got=%d", testIndex, len(sequences.Commands))
		}

		if !testInsertStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedValuesTokens) {
			return
		}
	}
}

func testInsertStatement(t *testing.T, command ast.Command, expectedTableName string, expectedValuesTokens []token.Token) bool {
	if command.TokenLiteral() != "INSERT" {
		t.Errorf("command.TokenLiteral() not 'INSERT'. got=%q", command.TokenLiteral())
		return false
	}

	actualInsertCommand, ok := command.(*ast.InsertCommand)
	if !ok {
		t.Errorf("actualInsertCommand is not %T. got=%T", &ast.InsertCommand{}, command)
		return false
	}

	if actualInsertCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualInsertCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !tokenArrayEquals(actualInsertCommand.Values, expectedValuesTokens) {
		t.Errorf("")
		return false
	}

	return true
}

func TestParseSelectCommand(t *testing.T) {
	tests := []struct {
		input             string
		expectedTableName string
		expectedSpaces    []ast.Space
		expectedDistinct  bool
	}{
		{"SELECT * FROM TBL;", "TBL", []ast.Space{{ColumnName: token.Token{Type: token.ASTERISK, Literal: "*"}}}, false},
		{"SELECT ONE, TWO, THREE FROM TBL;", "TBL", []ast.Space{{ColumnName: token.Token{Type: token.IDENT, Literal: "ONE"}}, {ColumnName: token.Token{Type: token.IDENT, Literal: "TWO"}}, {ColumnName: token.Token{Type: token.IDENT, Literal: "THREE"}}}, false},
		{"SELECT DISTINCT * FROM TBL;", "TBL", []ast.Space{{ColumnName: token.Token{Type: token.ASTERISK, Literal: "*"}}}, true},
	}

	for testIndex, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences, err := parserInstance.ParseSequence()
		if err != nil {
			t.Fatalf("[%d] Got error from parser: %s", testIndex, err)
		}

		if len(sequences.Commands) != 1 {
			t.Fatalf("[%d] sequences does not contain 1 statements. got=%d", testIndex, len(sequences.Commands))
		}

		if !testSelectStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedSpaces, tt.expectedDistinct) {
			return
		}
	}
}

func TestParseWhereCommand(t *testing.T) {
	firstExpression := ast.ConditionExpression{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "fda"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	secondExpression := ast.ConditionExpression{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	thirdExpression := ast.ContainExpression{
		Left: ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName3"}},
		Right: []ast.Anonymitifier{
			{Token: token.Token{Type: token.LITERAL, Literal: "1"}},
			{Token: token.Token{Type: token.LITERAL, Literal: "2"}},
		},
		Contains: true,
	}

	fourthExpression := ast.ContainExpression{
		Left: ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName4"}},
		Right: []ast.Anonymitifier{
			{Token: token.Token{Type: token.IDENT, Literal: "one"}},
			{Token: token.Token{Type: token.IDENT, Literal: "two"}},
		},
		Contains: false,
	}

	fifthExpression := ast.ConditionExpression{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName5"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.NULL, Literal: "NULL"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	tests := []struct {
		input              string
		expectedExpression ast.Expression
	}{
		{
			input:              "SELECT * FROM TBL WHERE colName1 EQUAL 'fda';",
			expectedExpression: firstExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName2 EQUAL 6462389;",
			expectedExpression: secondExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName3 IN (1, 2);",
			expectedExpression: thirdExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName4 NOTIN ('one', 'two');",
			expectedExpression: fourthExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName5 EQUAL NULL;",
			expectedExpression: fifthExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName5 EQUAL NULL;",
			expectedExpression: fifthExpression,
		},
	}

	for testIndex, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences, err := parserInstance.ParseSequence()
		if err != nil {
			t.Fatalf("[%d] Got error from parser: %s", testIndex, err)
		}

		if len(sequences.Commands) != 1 {
			t.Fatalf("[%d] sequences does not contain 1 statements, got=%d", testIndex, len(sequences.Commands))
		}

		selectCommand := sequences.Commands[0].(*ast.SelectCommand)
		if !selectCommand.HasWhereCommand() {
			t.Fatalf("[%d] sequences does not contain where command", testIndex)
		}

		if !whereStatementIsValid(t, selectCommand.WhereCommand, tt.expectedExpression) {
			return
		}
	}
}

func TestParseDeleteCommand(t *testing.T) {
	input := "DELETE FROM colName1 WHERE colName2 EQUAL 6462389;"
	expectedDeleteCommand := ast.DeleteCommand{
		Token: token.Token{Type: token.DELETE, Literal: "DELETE"},
		Name:  ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
	}
	expectedWhereCommand := ast.ConditionExpression{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	actualDeleteCommand, ok := sequences.Commands[0].(*ast.DeleteCommand)
	if !ok {
		t.Errorf("actualDeleteCommand is not %T. got=%T", &ast.DeleteCommand{}, sequences.Commands[0])
	}

	if expectedDeleteCommand.TokenLiteral() != actualDeleteCommand.TokenLiteral() {
		t.Errorf("TokenLiteral of DeleteCommand is not %s. got=%s", expectedDeleteCommand.TokenLiteral(), actualDeleteCommand.TokenLiteral())
	}

	if expectedDeleteCommand.Name.GetToken().Literal != actualDeleteCommand.Name.GetToken().Literal {
		t.Errorf("Table name of DeleteCommand is not %s. got=%s", expectedDeleteCommand.Name.GetToken().Literal, actualDeleteCommand.Name.GetToken().Literal)
	}

	if !actualDeleteCommand.HasWhereCommand() {
		t.Fatalf("sequences does not contain where command")
	}

	if !whereStatementIsValid(t, actualDeleteCommand.WhereCommand, expectedWhereCommand) {
		return
	}
}

func TestParseDropCommand(t *testing.T) {
	input := "DROP TABLE table;"
	expectedDropCommand := ast.DropCommand{
		Token: token.Token{Type: token.DROP, Literal: "DROP"},
		Name:  ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "table"}},
	}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	actualDropCommand, ok := sequences.Commands[0].(*ast.DropCommand)
	if !ok {
		t.Errorf("actualDropCommand is not %T. got=%T", &ast.DropCommand{}, sequences.Commands[0])
	}

	if expectedDropCommand.TokenLiteral() != actualDropCommand.TokenLiteral() {
		t.Errorf("TokenLiteral of DropCommand is not %s. got=%s", expectedDropCommand.TokenLiteral(), actualDropCommand.TokenLiteral())
	}

	if expectedDropCommand.Name.GetToken().Literal != actualDropCommand.Name.GetToken().Literal {
		t.Errorf("Table name of DropCommand is not %s. got=%s", expectedDropCommand.Name.GetToken().Literal, actualDropCommand.Name.GetToken().Literal)
	}
}

func TestSelectWithOrderByCommand(t *testing.T) {
	input := "SELECT * FROM tableName ORDER BY colName1 DESC;"
	expectedSortPattern := ast.SortPattern{
		ColumnName: token.Token{Type: token.IDENT, Literal: "colName1"},
		Order:      token.Token{Type: token.DESC, Literal: "DESC"},
	}
	expectedOrderByCommand := ast.OrderByCommand{
		Token:        token.Token{Type: token.ORDER, Literal: "ORDER"},
		SortPatterns: []ast.SortPattern{expectedSortPattern},
	}
	expectedTableName := "tableName"
	expectedSpaces := []ast.Space{{ColumnName: token.Token{Type: token.ASTERISK, Literal: "*"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}

	if !selectCommand.HasOrderByCommand() {
		t.Fatalf("sequences does not contain where command")
	}

	testOrderByCommands(t, expectedOrderByCommand, selectCommand.OrderByCommand)
}

func TestSelectWithLimitCommand(t *testing.T) {
	input := "SELECT * FROM tableName LIMIT 5;"
	expectedLimitCommand := ast.LimitCommand{
		Token: token.Token{Type: token.LIMIT, Literal: "LIMIT"},
		Count: 5,
	}
	expectedTableName := "tableName"
	expectedSpaces := []ast.Space{{ColumnName: token.Token{Type: token.ASTERISK, Literal: "*"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}

	if !selectCommand.HasLimitCommand() {
		t.Fatalf("sequences does not contain where command")
	}

	testLimitCommands(t, expectedLimitCommand, selectCommand.LimitCommand)
}

func TestSelectWithOffsetCommand(t *testing.T) {
	input := "SELECT * FROM tableName OFFSET 5;"
	expectedOffsetCommand := ast.OffsetCommand{
		Token: token.Token{Type: token.OFFSET, Literal: "OFFSET"},
		Count: 5,
	}

	expectedTableName := "tableName"
	expectedSpaces := []ast.Space{{ColumnName: token.Token{Type: token.ASTERISK, Literal: "*"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}

	if !selectCommand.HasOffsetCommand() {
		t.Fatalf("select command should have offset command")
	}
	testOffsetCommands(t, expectedOffsetCommand, selectCommand.OffsetCommand)
}

func TestSelectWithLimitAndOffsetCommand(t *testing.T) {
	input := "SELECT * FROM tableName ORDER BY colName1 DESC LIMIT 2 OFFSET 13;"
	expectedLimitCommand := ast.LimitCommand{
		Token: token.Token{Type: token.LIMIT, Literal: "LIMIT"},
		Count: 2,
	}
	expectedOffsetCommand := ast.OffsetCommand{
		Token: token.Token{Type: token.OFFSET, Literal: "OFFSET"},
		Count: 13,
	}
	expectedTableName := "tableName"
	expectedSpaces := []ast.Space{{ColumnName: token.Token{Type: token.ASTERISK, Literal: "*"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}

	if !selectCommand.HasLimitCommand() {
		t.Fatalf("select command should have limit command")
	}
	if !selectCommand.HasOffsetCommand() {
		t.Fatalf("select command should have offset command")
	}

	testLimitCommands(t, expectedLimitCommand, selectCommand.LimitCommand)
	testOffsetCommands(t, expectedOffsetCommand, selectCommand.OffsetCommand)
}

func TestSelectWithDefaultInnerJoinCommand(t *testing.T) {
	input := "SELECT tbl.one, tbl2.two FROM tbl JOIN tbl2 ON tbl.one EQUAL tbl2.one;"
	expectedJoinCommand := ast.JoinCommand{
		Token:    token.Token{Type: token.JOIN, Literal: "JOIN"},
		Name:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2"}},
		JoinType: token.Token{Type: token.INNER, Literal: "INNER"},
		Expression: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl.one"}},
			Right:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2.one"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
		},
	}
	expectedTableName := "tbl"
	expectedSpace := []ast.Space{{ColumnName: token.Token{Type: token.IDENT, Literal: "tbl.one"}}, {ColumnName: token.Token{Type: token.IDENT, Literal: "tbl2.two"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpace, false) {
		return
	}

	if !selectCommand.HasJoinCommand() {
		t.Fatalf("select command should have join command")
	}

	testJoinCommands(t, expectedJoinCommand, *selectCommand.JoinCommand)
}

func TestSelectWithInnerJoinCommand(t *testing.T) {
	input := "SELECT tbl.one, tbl2.two FROM tbl INNER JOIN tbl2 ON tbl.one EQUAL tbl2.one;"
	expectedJoinCommand := ast.JoinCommand{
		Token:    token.Token{Type: token.JOIN, Literal: "JOIN"},
		Name:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2"}},
		JoinType: token.Token{Type: token.INNER, Literal: "INNER"},
		Expression: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl.one"}},
			Right:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2.one"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
		},
	}
	expectedTableName := "tbl"
	expectedSpace := []ast.Space{{ColumnName: token.Token{Type: token.IDENT, Literal: "tbl.one"}}, {ColumnName: token.Token{Type: token.IDENT, Literal: "tbl2.two"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpace, false) {
		return
	}

	if !selectCommand.HasJoinCommand() {
		t.Fatalf("select command should have join command")
	}

	testJoinCommands(t, expectedJoinCommand, *selectCommand.JoinCommand)
}

func TestSelectWithLeftJoinCommand(t *testing.T) {
	input := "SELECT tbl.one, tbl2.two FROM tbl LEFT JOIN tbl2 ON tbl.one EQUAL tbl2.one;"
	expectedJoinCommand := ast.JoinCommand{
		Token:    token.Token{Type: token.JOIN, Literal: "JOIN"},
		Name:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2"}},
		JoinType: token.Token{Type: token.LEFT, Literal: "LEFT"},
		Expression: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl.one"}},
			Right:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2.one"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
		},
	}
	expectedTableName := "tbl"
	expectedSpaces := []ast.Space{{ColumnName: token.Token{Type: token.IDENT, Literal: "tbl.one"}}, {ColumnName: token.Token{Type: token.IDENT, Literal: "tbl2.two"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}

	if !selectCommand.HasJoinCommand() {
		t.Fatalf("select command should have join command")
	}

	testJoinCommands(t, expectedJoinCommand, *selectCommand.JoinCommand)
}

func TestSelectWithRightJoinCommand(t *testing.T) {
	input := "SELECT tbl.one, tbl2.two FROM tbl RIGHT JOIN tbl2 ON tbl.one EQUAL tbl2.one;"
	expectedJoinCommand := ast.JoinCommand{
		Token:    token.Token{Type: token.JOIN, Literal: "JOIN"},
		Name:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2"}},
		JoinType: token.Token{Type: token.RIGHT, Literal: "RIGHT"},
		Expression: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl.one"}},
			Right:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2.one"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
		},
	}
	expectedTableName := "tbl"
	expectedSpaces := []ast.Space{{ColumnName: token.Token{Type: token.IDENT, Literal: "tbl.one"}}, {ColumnName: token.Token{Type: token.IDENT, Literal: "tbl2.two"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}

	if !selectCommand.HasJoinCommand() {
		t.Fatalf("select command should have join command")
	}

	testJoinCommands(t, expectedJoinCommand, *selectCommand.JoinCommand)
}

func TestSelectWithFullJoinCommand(t *testing.T) {
	input := "SELECT tbl.one, tbl2.two FROM tbl FULL JOIN tbl2 ON tbl.one EQUAL tbl2.one;"
	expectedJoinCommand := ast.JoinCommand{
		Token:    token.Token{Type: token.JOIN, Literal: "JOIN"},
		Name:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2"}},
		JoinType: token.Token{Type: token.FULL, Literal: "FULL"},
		Expression: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl.one"}},
			Right:     ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "tbl2.one"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
		},
	}
	expectedTableName := "tbl"
	expectedSpaces := []ast.Space{{ColumnName: token.Token{Type: token.IDENT, Literal: "tbl.one"}}, {ColumnName: token.Token{Type: token.IDENT, Literal: "tbl2.two"}}}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}

	if !selectCommand.HasJoinCommand() {
		t.Fatalf("select command should have join command")
	}

	testJoinCommands(t, expectedJoinCommand, *selectCommand.JoinCommand)
}

func TestSelectWithAggregateFunctions(t *testing.T) {
	input := "SELECT MIN(colOne), MAX(colOne), COUNT(*), COUNT(colOne), SUM(colOne), AVG(colOne) FROM tbl;"

	expectedTableName := "tbl"
	expectedSpaces := []ast.Space{
		{
			ColumnName:    token.Token{Type: token.IDENT, Literal: "colOne"},
			AggregateFunc: &token.Token{Type: token.MIN, Literal: "MIN"},
		},
		{
			ColumnName:    token.Token{Type: token.ASTERISK, Literal: "colOne"},
			AggregateFunc: &token.Token{Type: token.MAX, Literal: "MAX"},
		},
		{
			ColumnName:    token.Token{Type: token.IDENT, Literal: "*"},
			AggregateFunc: &token.Token{Type: token.COUNT, Literal: "COUNT"},
		},
		{
			ColumnName:    token.Token{Type: token.IDENT, Literal: "colOne"},
			AggregateFunc: &token.Token{Type: token.COUNT, Literal: "COUNT"},
		},
		{
			ColumnName:    token.Token{Type: token.IDENT, Literal: "colOne"},
			AggregateFunc: &token.Token{Type: token.SUM, Literal: "SUM"},
		},
		{
			ColumnName:    token.Token{Type: token.IDENT, Literal: "colOne"},
			AggregateFunc: &token.Token{Type: token.AVG, Literal: "AVG"},
		},
	}

	lexer := lexer.RunLexer(input)
	parserInstance := New(lexer)
	sequences, err := parserInstance.ParseSequence()
	if err != nil {
		t.Fatalf("Got error from parser: %s", err)
	}

	if len(sequences.Commands) != 1 {
		t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
	}

	selectCommand := sequences.Commands[0].(*ast.SelectCommand)

	if !testSelectStatement(t, selectCommand, expectedTableName, expectedSpaces, false) {
		return
	}
}

func TestParseUpdateCommand(t *testing.T) {
	tests := []struct {
		input             string
		expectedTableName string
		expectedChanges   map[token.Token]ast.Anonymitifier
	}{
		{
			input: "UPDATE tbl SET colName TO 5;", expectedTableName: "tbl", expectedChanges: map[token.Token]ast.Anonymitifier{
				{Type: token.IDENT, Literal: "colName"}: {Token: token.Token{Type: token.LITERAL, Literal: "5"}},
			},
		},
		{
			input: "UPDATE tbl1 SET colName1 TO 'hi hello', colName2 TO 5;", expectedTableName: "tbl1", expectedChanges: map[token.Token]ast.Anonymitifier{
				{Type: token.IDENT, Literal: "colName1"}: {Token: token.Token{Type: token.IDENT, Literal: "hi hello"}},
				{Type: token.IDENT, Literal: "colName2"}: {Token: token.Token{Type: token.LITERAL, Literal: "5"}},
			},
		},
		{
			input: "UPDATE tbl1 SET colName1 TO NULL, colName2 TO 'NULL';", expectedTableName: "tbl1", expectedChanges: map[token.Token]ast.Anonymitifier{
				{Type: token.IDENT, Literal: "colName1"}: {Token: token.Token{Type: token.NULL, Literal: "NULL"}},
				{Type: token.IDENT, Literal: "colName2"}: {Token: token.Token{Type: token.LITERAL, Literal: "NULL"}},
			},
		},
	}

	for testIndex, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences, err := parserInstance.ParseSequence()
		if err != nil {
			t.Fatalf("[%d] Got error from parser: %s", testIndex, err)
		}

		if len(sequences.Commands) != 1 {
			t.Fatalf("[%d] sequences does not contain 1 statements. got=%d", testIndex, len(sequences.Commands))
		}

		if !testUpdateStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedChanges) {
			return
		}
	}
}

func TestParseUpdateCommandWithWhere(t *testing.T) {
	tests := []struct {
		input                string
		expectedTableName    string
		expectedChanges      map[token.Token]ast.Anonymitifier
		expectedWhereCommand ast.Expression
	}{
		{
			input:             "UPDATE tbl SET colName TO 5 WHERE id EQUAL 3;",
			expectedTableName: "tbl",
			expectedChanges: map[token.Token]ast.Anonymitifier{
				{Type: token.IDENT, Literal: "colName"}: {Token: token.Token{Type: token.LITERAL, Literal: "5"}},
			},
			expectedWhereCommand: ast.ConditionExpression{
				Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "id"}},
				Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "3"}},
				Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
			},
		},
	}

	for testIndex, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences, err := parserInstance.ParseSequence()
		if err != nil {
			t.Fatalf("Got error from parser: %s", err)
		}

		if len(sequences.Commands) != 1 {
			t.Fatalf("[%d] sequences does not contain 1 statements. got=%d", testIndex, len(sequences.Commands))
		}

		actualUpdateCommand, ok := sequences.Commands[0].(*ast.UpdateCommand)

		if !ok {
			t.Errorf("[%d] actualUpdateCommand is not %T. got=%T", testIndex, &ast.UpdateCommand{}, sequences.Commands[0])
		}

		if !testUpdateStatement(t, actualUpdateCommand, tt.expectedTableName, tt.expectedChanges) {
			return
		}

		if !actualUpdateCommand.HasWhereCommand() {
			t.Errorf("[%d] actualUpdateCommand should have where command", testIndex)
		}

		if !whereStatementIsValid(t, actualUpdateCommand.WhereCommand, tt.expectedWhereCommand) {
			return
		}
	}
}

func TestParseLogicOperatorsInCommand(t *testing.T) {

	firstExpression := ast.OperationExpression{
		Left: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "fda"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}},
		Right: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.NULL, Literal: "NULL"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "NOT"}},
		Operation: token.Token{Type: token.AND, Literal: "AND"},
	}

	secondExpression := ast.OperationExpression{
		Left: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
			Condition: token.Token{Type: token.NOT, Literal: "NOT"}},
		Right: ast.ConditionExpression{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "NULL"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}},
		Operation: token.Token{Type: token.OR, Literal: "OR"},
	}

	thirdExpression := ast.BooleanExpression{
		Boolean: token.Token{Type: token.TRUE, Literal: "TRUE"},
	}

	fourthExpression := ast.ConditionExpression{
		Left:      ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "colName1 EQUAL;"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "colName1 EQUAL;"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}}

	tests := []struct {
		input              string
		expectedExpression ast.Expression
	}{
		{
			input:              "SELECT * FROM TBL WHERE colName1 EQUAL 'fda' AND colName2 NOT NULL;",
			expectedExpression: firstExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName2 NOT 6462389 OR colName1 EQUAL 'NULL';",
			expectedExpression: secondExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE TRUE;",
			expectedExpression: thirdExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE 'colName1 EQUAL;' EQUAL 'colName1 EQUAL;';",
			expectedExpression: fourthExpression,
		},
	}

	for testIndex, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences, err := parserInstance.ParseSequence()
		if err != nil {
			t.Fatalf("Got error from parser: %s", err)
		}

		if len(sequences.Commands) != 1 {
			t.Fatalf("[%d] sequences does not contain 1 statements. got=%d", testIndex, len(sequences.Commands))
		}

		selectCommand := sequences.Commands[0].(*ast.SelectCommand)

		if !selectCommand.HasWhereCommand() {
			t.Fatalf("[%d] sequences does not contain where command", testIndex)
		}

		if !whereStatementIsValid(t, selectCommand.WhereCommand, tt.expectedExpression) {
			t.Fatalf("[%d] Actual expression and expected one are different", testIndex)
		}
	}
}

func testSelectStatement(t *testing.T, command ast.Command, expectedTableName string, expectedSpaces []ast.Space, expectedDistinct bool) bool {
	if command.TokenLiteral() != "SELECT" {
		t.Errorf("command.TokenLiteral() not 'SELECT'. got=%q", command.TokenLiteral())
		return false
	}

	actualSelectCommand, ok := command.(*ast.SelectCommand)
	if !ok {
		t.Errorf("actualSelectCommand is not %T. got=%T", &ast.SelectCommand{}, command)
		return false
	}

	if actualSelectCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualSelectCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if actualSelectCommand.HasDistinct != expectedDistinct {
		t.Errorf("HasDistinct should be set to %t, got=%t", expectedDistinct, actualSelectCommand.HasDistinct)
		return false
	}

	if !spaceArrayEquals(actualSelectCommand.Space, expectedSpaces) {
		t.Errorf("actualSelectCommand has diffrent space than expected. %+v != %+v", actualSelectCommand.Space, expectedSpaces)
		return false
	}

	return true
}

func testUpdateStatement(t *testing.T, command ast.Command, expectedTableName string, expectedChanges map[token.Token]ast.Anonymitifier) bool {
	if command.TokenLiteral() != "UPDATE" {
		t.Errorf("command.TokenLiteral() not 'UPDATE'. got=%q", command.TokenLiteral())
		return false
	}
	actualUpdateCommand, ok := command.(*ast.UpdateCommand)

	if !ok {
		t.Errorf("actualUpdateCommand is not %T. got=%T", &ast.UpdateCommand{}, command)
		return false
	}
	if actualUpdateCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualUpdateCommand.TokenLiteral(), expectedTableName)
		return false
	}
	if !tokenMapEquals(actualUpdateCommand.Changes, expectedChanges) {
		t.Errorf("")
		return false
	}

	return true
}

func whereStatementIsValid(t *testing.T, command ast.Command, expectedExpression ast.Expression) bool {
	if command.TokenLiteral() != "WHERE" {
		t.Errorf("command.TokenLiteral() not 'WHERE'. got=%q", command.TokenLiteral())
		return false
	}

	actualWhereCommand, ok := command.(*ast.WhereCommand)
	if !ok {
		t.Errorf("actualWhereCommand is not %T. got=%T", &ast.WhereCommand{}, command)
		return false
	}

	if !expressionsAreEqual(actualWhereCommand.Expression, expectedExpression) {
		t.Errorf("Actual expression is not equal to expected one.\nActual: %#v\nExpected: %#v", actualWhereCommand.Expression, expectedExpression)
		return false
	}

	return true
}

func stringArrayEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func tokenArrayEquals(a []token.Token, b []token.Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Literal != b[i].Literal {
			return false
		}
		if v.Type != b[i].Type {
			return false
		}
	}
	return true
}

func spaceArrayEquals(a []ast.Space, b []ast.Space) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.ColumnName.Literal != b[i].ColumnName.Literal {
			return false
		}
		if v.ContainsAggregateFunc() != b[i].ContainsAggregateFunc() {
			return false
		}
		if v.ContainsAggregateFunc() && b[i].ContainsAggregateFunc() && v.AggregateFunc.Literal != b[i].AggregateFunc.Literal {
			return false
		}
	}
	return true
}

func tokenMapEquals(a map[token.Token]ast.Anonymitifier, b map[token.Token]ast.Anonymitifier) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if v.GetToken().Literal != b[k].GetToken().Literal {
			return false
		}
	}
	return true
}

func testOrderByCommands(t *testing.T, expectedOrderByCommand ast.OrderByCommand, actualOrderByCommand *ast.OrderByCommand) {

	if expectedOrderByCommand.Token.Type != actualOrderByCommand.Token.Type {
		t.Errorf("Expecting Token TokenType: %q, got: %q", expectedOrderByCommand.Token.Type, actualOrderByCommand.Token.Type)
	}
	if expectedOrderByCommand.Token.Literal != actualOrderByCommand.Token.Literal {
		t.Errorf("Expecting Token Literal: %s, got: %s", expectedOrderByCommand.Token.Literal, actualOrderByCommand.Token.Literal)
	}
	if len(expectedOrderByCommand.SortPatterns) != len(actualOrderByCommand.SortPatterns) {
		t.Errorf("Expecting Sorting Pattern Array to have: %d elements, got: %d", len(expectedOrderByCommand.SortPatterns), len(actualOrderByCommand.SortPatterns))
	}

	for i, expectedSortPattern := range expectedOrderByCommand.SortPatterns {
		if expectedSortPattern.Order.Literal != actualOrderByCommand.SortPatterns[i].Order.Literal {
			t.Errorf("Expecting Order: %s, got: %s", expectedSortPattern.Order.Literal, actualOrderByCommand.SortPatterns[i].Order.Literal)
		}
		if expectedSortPattern.ColumnName.Literal != expectedOrderByCommand.SortPatterns[i].ColumnName.Literal {
			t.Errorf("Expecting Column Name: %s, got: %s", expectedSortPattern.ColumnName.Literal, actualOrderByCommand.SortPatterns[i].ColumnName.Literal)
		}
	}
}

func testLimitCommands(t *testing.T, expectedLimitCommand ast.LimitCommand, actualLimitCommand *ast.LimitCommand) {

	if expectedLimitCommand.Token.Type != actualLimitCommand.Token.Type {
		t.Errorf("Expecting Token TokenType: %q, got: %q", expectedLimitCommand.Token.Type, actualLimitCommand.Token.Type)
	}
	if expectedLimitCommand.Token.Literal != actualLimitCommand.Token.Literal {
		t.Errorf("Expecting Token Literal: %s, got: %s", expectedLimitCommand.Token.Literal, actualLimitCommand.Token.Literal)
	}
	if expectedLimitCommand.Count != actualLimitCommand.Count {
		t.Errorf("Expecting Count to have value: %d, got: %d", expectedLimitCommand.Count, actualLimitCommand.Count)
	}
}

func testOffsetCommands(t *testing.T, expectedOffsetCommand ast.OffsetCommand, actualOffsetCommand *ast.OffsetCommand) {

	if expectedOffsetCommand.Token.Type != actualOffsetCommand.Token.Type {
		t.Errorf("Expecting Token TokenType: %q, got: %q", expectedOffsetCommand.Token.Type, actualOffsetCommand.Token.Type)
	}
	if expectedOffsetCommand.Token.Literal != actualOffsetCommand.Token.Literal {
		t.Errorf("Expecting Token Literal: %s, got: %s", expectedOffsetCommand.Token.Literal, actualOffsetCommand.Token.Literal)
	}
	if expectedOffsetCommand.Count != actualOffsetCommand.Count {
		t.Errorf("Expecting Count to have value: %d, got: %d", expectedOffsetCommand.Count, actualOffsetCommand.Count)
	}
}

func testJoinCommands(t *testing.T, expectedJoinCommand ast.JoinCommand, actualJoinCommand ast.JoinCommand) {

	if expectedJoinCommand.Token.Type != actualJoinCommand.Token.Type {
		t.Errorf("Expecting Token TokenType: %q, got: %q", expectedJoinCommand.Token.Type, actualJoinCommand.Token.Type)
	}
	if expectedJoinCommand.Token.Literal != actualJoinCommand.Token.Literal {
		t.Errorf("Expecting Token Literal: %s, got: %s", expectedJoinCommand.Token.Literal, actualJoinCommand.Token.Literal)
	}
	if expectedJoinCommand.Name != actualJoinCommand.Name {
		t.Errorf("Expecting Name to has a value: %s, got: %s", expectedJoinCommand.Name, actualJoinCommand.Name)
	}
	if !expressionsAreEqual(actualJoinCommand.Expression, expectedJoinCommand.Expression) {
		t.Errorf("Actual expression is not equal to expected one.\nActual: %#v\nExpected: %#v", actualJoinCommand.Expression, expectedJoinCommand.Expression)
	}
}
func expressionsAreEqual(first ast.Expression, second ast.Expression) bool {

	booleanExpression, booleanExpressionIsValid := first.(*ast.BooleanExpression)
	if booleanExpressionIsValid {
		return validateBooleanExpressions(second, booleanExpression)
	}

	conditionExpression, conditionExpressionIsValid := first.(*ast.ConditionExpression)
	if conditionExpressionIsValid {
		return validateConditionExpression(second, conditionExpression)
	}

	operationExpression, operationExpressionIsValid := first.(*ast.OperationExpression)
	if operationExpressionIsValid {
		return validateOperationExpression(second, operationExpression)
	}

	containExpression, containExpressionIsValid := first.(*ast.ContainExpression)
	if containExpressionIsValid {
		return validateContainExpression(second, containExpression)
	}

	return false
}

func validateOperationExpression(second ast.Expression, operationExpression *ast.OperationExpression) bool {
	secondOperationExpression, secondOperationExpressionIsValid := second.(ast.OperationExpression)

	if !secondOperationExpressionIsValid {
		return false
	}

	if operationExpression.Operation.Literal != secondOperationExpression.Operation.Literal {
		return false
	}

	return expressionsAreEqual(operationExpression.Left, secondOperationExpression.Left) && expressionsAreEqual(operationExpression.Right, secondOperationExpression.Right)
}

func validateContainExpression(expression ast.Expression, expectedContainExpression *ast.ContainExpression) bool {
	actualContainExpression, actualContainExpressionIsValid := expression.(ast.ContainExpression)

	if !actualContainExpressionIsValid {
		return false
	}

	if expectedContainExpression.Contains != actualContainExpression.Contains {
		return false
	}

	if actualContainExpression.Left.GetToken().Literal != expectedContainExpression.Left.GetToken().Literal &&
		actualContainExpression.Left.IsIdentifier() == expectedContainExpression.Left.IsIdentifier() {
		return false
	}

	if len(expectedContainExpression.Right) != len(actualContainExpression.Right) {
		return false
	}

	for i := 0; i < len(expectedContainExpression.Right); i++ {
		if expectedContainExpression.Right[i] != actualContainExpression.Right[i] {
			return false
		}
	}

	return true
}

func validateConditionExpression(second ast.Expression, conditionExpression *ast.ConditionExpression) bool {
	secondConditionExpression, secondConditionExpressionIsValid := second.(ast.ConditionExpression)

	if !secondConditionExpressionIsValid {
		return false
	}

	if conditionExpression.Left.GetToken().Literal != secondConditionExpression.Left.GetToken().Literal &&
		conditionExpression.Left.IsIdentifier() == secondConditionExpression.Left.IsIdentifier() {
		return false
	}

	if conditionExpression.Right.GetToken().Literal != secondConditionExpression.Right.GetToken().Literal &&
		conditionExpression.Right.IsIdentifier() == secondConditionExpression.Right.IsIdentifier() {
		return false
	}

	if conditionExpression.Condition.Literal != secondConditionExpression.Condition.Literal {
		return false
	}

	return true
}

func validateBooleanExpressions(second ast.Expression, booleanExpression *ast.BooleanExpression) bool {
	secondBooleanExpression, secondBooleanExpressionIsValid := second.(ast.BooleanExpression)

	if !secondBooleanExpressionIsValid {
		return false
	}

	if booleanExpression.Boolean.Literal != secondBooleanExpression.Boolean.Literal {
		return false
	}

	return true
}
