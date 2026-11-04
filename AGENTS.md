# GO4SQL

An in-memory SQL engine written in Go using only the standard library. It runs
in file, stream or socket mode. See `README.md` for usage and supported SQL.

## Structure

SQL flows through the packages in this order: `lexer` → `parser` → `engine`.

- `main.go`: parses flags and starts a mode
- `modes/`: file, stream and socket input handling
- `token/`: token types and keywords
- `lexer/`: turns input text into tokens
- `ast/`: node types for parsed statements
- `parser/`: turns tokens into an AST
- `engine/`: executes statements against in-memory tables
- `e2e/`: SQL files in `test_files/` and their `expected_outputs/`
- `helm/`, `Dockerfile`: deployment

## Commands

Run from the repository root.

```shell
go build                 # creates the GO4SQL binary
go test ./...            # unit tests
./e2e/e2e_test.sh        # e2e tests; run go build first
```
