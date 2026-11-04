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

## Design rules

Based on Thorsten Ball's *Writing An Interpreter In Go*, adapted for SQL.

- Keep the pipeline one-way, with each stage in its own package.
- Hand-write the lexer and parser. Use only the standard library.
- Build features vertically: token, lexer, AST and parser, engine, tests.
- Prefer readable code to clever code, and tolerate some duplication.
- Write table-driven tests where they fit, and run real SQL through the real
  pipeline with no mocks. New statements also get a numbered e2e pair.

## Differences from the book

- The parser returns the first typed error; it doesn't collect errors.
- There is no Pratt parsing or precedence. WHERE recurses right, so
  `a AND b OR c` means `a AND (b OR c)`.
- `Command` replaces the book's `Statement`, and SELECT runs a transformer
  pipeline rather than a pure `Eval` walk.

## Commands

Run from the repository root.

```shell
go build                 # creates the GO4SQL binary
go test ./...            # unit tests
./e2e/e2e_test.sh        # e2e tests; run go build first
```
