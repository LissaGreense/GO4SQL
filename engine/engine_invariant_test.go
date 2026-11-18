package engine

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
	"github.com/LissaGreense/GO4SQL/token"
)

// The Feared-Failure Invariant: after any statement run serially, including a failed one,
// every table's columns have equal length and every stored value is the column's declared
// type, or NULL where the column allows NULL. A failed statement also changes nothing.

// storedTypeOracle is deliberately independent of declaredTarget, so the walker doesn't
// share the code it is checking.
var storedTypeOracle = map[token.Type]SupportedTypes{
	token.INT:  IntType,
	token.TEXT: StringType,
}

func walkTables(t *testing.T, engine *DbEngine, context string) {
	t.Helper()
	for tableName, table := range engine.Tables {
		if len(table.Columns) == 0 {
			continue
		}
		rowCount := len(table.Columns[0].Values)
		for _, column := range table.Columns {
			if len(column.Values) != rowCount {
				t.Fatalf("%s: table %s is ragged: column %s has %d values, column %s has %d",
					context, tableName, column.Name, len(column.Values), table.Columns[0].Name, rowCount)
			}
			wantType, known := storedTypeOracle[column.Type.Type]
			if !known {
				t.Fatalf("%s: table %s column %s has unknown declared type %s", context, tableName, column.Name, column.Type.Type)
			}
			for rowIndex, value := range column.Values {
				if value.GetType() == NullType {
					if column.IsNotNull {
						t.Fatalf("%s: table %s row %d stores NULL in NOT NULL column %s", context, tableName, rowIndex, column.Name)
					}
					continue
				}
				if value.GetType() != wantType {
					t.Fatalf("%s: table %s row %d column %s stores %s %q, declared %s",
						context, tableName, rowIndex, column.Name, value.GetType(), value.ToString(), column.Type.Type)
				}
			}
		}
	}
}

func snapshotTables(engine *DbEngine) map[string][]Column {
	snapshot := make(map[string][]Column, len(engine.Tables))
	for tableName, table := range engine.Tables {
		columns := make([]Column, len(table.Columns))
		for i, column := range table.Columns {
			columns[i] = *column
			columns[i].Values = append([]ValueInterface{}, column.Values...)
		}
		snapshot[tableName] = columns
	}
	return snapshot
}

func assertTablesEqualSnapshot(t *testing.T, engine *DbEngine, snapshot map[string][]Column, context string) {
	t.Helper()
	if after := snapshotTables(engine); !reflect.DeepEqual(after, snapshot) {
		t.Fatalf("%s: tables changed:\n\tbefore: %+v\n\tafter:  %+v", context, snapshot, after)
	}
}

// expectedExitCode mirrors the e2e runner's convention (PR #75): a fixture may ship
// <name>_expected_exit next to its expected output; without it the fixture must exit 0.
func expectedExitCode(t *testing.T, testFile string) int {
	t.Helper()
	name := strings.TrimSuffix(filepath.Base(testFile), "_test")
	content, err := os.ReadFile(filepath.Join("..", "e2e", "expected_outputs", name+"_expected_exit"))
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatalf("reading expected exit for %s: %s", testFile, err)
	}
	code, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		t.Fatalf("expected exit for %s isn't a number: %s", testFile, err)
	}
	return code
}

func TestInvariantHoldsAcrossE2EScripts(t *testing.T) {
	testFiles, err := filepath.Glob(filepath.Join("..", "e2e", "test_files", "*_test"))
	if err != nil {
		t.Fatal(err)
	}
	if len(testFiles) == 0 {
		t.Fatal("no e2e scripts found")
	}

	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatal(err)
			}
			mayFail := expectedExitCode(t, testFile) != 0

			// One parse over the whole file: script 12 has a string literal spanning lines.
			sequence, err := parser.New(lexer.RunLexer(string(content))).ParseSequence()
			if err != nil {
				if mayFail {
					return
				}
				t.Fatalf("parser error: %s", err)
			}

			engine := New()
			for i, command := range sequence.Commands {
				context := "command " + strconv.Itoa(i) + " (" + command.TokenLiteral() + ")"
				snapshot := snapshotTables(engine)
				_, err := engine.Evaluate(&ast.Sequence{Commands: []ast.Command{command}})
				walkTables(t, engine, context)
				if err != nil {
					if !mayFail {
						t.Fatalf("%s: unexpected engine error: %s", context, err)
					}
					assertTablesEqualSnapshot(t, engine, snapshot, context)
					return
				}
			}
		})
	}
}

func TestFailedStatementsLeaveTablesUnchanged(t *testing.T) {
	tests := []struct {
		name  string
		setup []string
		bad   string
	}{
		{
			name:  "string into INT",
			setup: []string{"CREATE TABLE t( a INT, b INT );", "INSERT INTO t VALUES( 7, 8 );"},
			bad:   "INSERT INTO t VALUES( 1, 'x' );",
		},
		{
			name:  "decimal into INT",
			setup: []string{"CREATE TABLE t( a INT, b TEXT );", "INSERT INTO t VALUES( 7, 'y' );"},
			bad:   "INSERT INTO t VALUES( 3.14, 'x' );",
		},
		{
			name:  "overflow into INT",
			setup: []string{"CREATE TABLE t( a INT, b TEXT );", "INSERT INTO t VALUES( 7, 'y' );"},
			bad:   "INSERT INTO t VALUES( 99999999999999999999, 'x' );",
		},
		{
			name:  "NULL into NOT NULL",
			setup: []string{"CREATE TABLE t( a INT NOT NULL, b TEXT );", "INSERT INTO t VALUES( 7, 'y' );"},
			bad:   "INSERT INTO t VALUES( NULL, 'x' );",
		},
		{
			name: "duplicate column of different type",
			bad:  "CREATE TABLE t( a INT, a TEXT );",
		},
		{
			name: "duplicate column with NOT NULL",
			bad:  "CREATE TABLE t( a INT NOT NULL, a INT );",
		},
		{
			name:  "IF NOT EXISTS duplicate column on new table",
			setup: []string{"CREATE TABLE t( a INT );"},
			bad:   "CREATE TABLE IF NOT EXISTS t2( a INT, a TEXT );",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine := New()
			evaluateInputs(t, engine, test.setup...)
			snapshot := snapshotTables(engine)

			sequence, err := parser.New(lexer.RunLexer(test.bad)).ParseSequence()
			if err != nil {
				t.Fatalf("parser error for %q: %s", test.bad, err)
			}
			if _, err := engine.Evaluate(sequence); err == nil {
				t.Fatalf("expected error for %q, got nil", test.bad)
			}
			walkTables(t, engine, test.bad)
			assertTablesEqualSnapshot(t, engine, snapshot, test.bad)
		})
	}
}
