package token

type Type string

// Token - contains
type Token struct {
	Type    Type
	Literal string
}

const (
	// ASTERISK - Operators
	ASTERISK = "*"

	// IDENT - Identifiers + literals
	IDENT   = "IDENT"   // tab, car, apple...
	LITERAL = "LITERAL" // 1343456

	// COMMA - Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// EOF - Special tokens
	EOF        = ""
	APOSTROPHE = "'"

	// LPAREN - Paren
	LPAREN = "("
	RPAREN = ")"

	// CREATE - Keywords
	CREATE   = "CREATE"
	DROP     = "DROP"
	TABLE    = "TABLE"
	INSERT   = "INSERT"
	INTO     = "INTO"
	VALUES   = "VALUES"
	SELECT   = "SELECT"
	FROM     = "FROM"
	WHERE    = "WHERE"
	DELETE   = "DELETE"
	ORDER    = "ORDER"
	BY       = "BY"
	ASC      = "ASC"
	DESC     = "DESC"
	LIMIT    = "LIMIT"
	OFFSET   = "OFFSET"
	UPDATE   = "UPDATE"
	SET      = "SET"
	DISTINCT = "DISTINCT"
	JOIN     = "JOIN"
	INNER    = "INNER"
	FULL     = "FULL"
	LEFT     = "LEFT"
	RIGHT    = "RIGHT"
	ON       = "ON"

	TO = "TO"

	// EQUAL - Logical operations
	EQUAL = "EQUAL"
	NOT   = "NOT"
	AND   = "AND"
	OR    = "OR"
	TRUE  = "TRUE"
	FALSE = "FALSE"

	// TEXT - Data types
	TEXT = "TEXT"
	INT  = "INT"

	// ILLEGAL - System
	ILLEGAL = "ILLEGAL"
)

var keywords = map[string]Type{
	"TEXT":     TEXT,
	"INT":      INT,
	"CREATE":   CREATE,
	"DROP":     DROP,
	"TABLE":    TABLE,
	"INSERT":   INSERT,
	"INTO":     INTO,
	"SELECT":   SELECT,
	"FROM":     FROM,
	"DELETE":   DELETE,
	"ORDER":    ORDER,
	"BY":       BY,
	"ASC":      ASC,
	"DESC":     DESC,
	"LIMIT":    LIMIT,
	"OFFSET":   OFFSET,
	"UPDATE":   UPDATE,
	"SET":      SET,
	"DISTINCT": DISTINCT,
	"INNER":    INNER,
	"FULL":     FULL,
	"LEFT":     LEFT,
	"RIGHT":    RIGHT,
	"JOIN":     JOIN,
	"ON":       ON,
	"TO":       TO,
	"VALUES":   VALUES,
	"WHERE":    WHERE,
	"EQUAL":    EQUAL,
	"NOT":      NOT,
	"AND":      AND,
	"OR":       OR,
	"TRUE":     TRUE,
	"FALSE":    FALSE,
}

// LookupIdent - Return keyword type from defined list if exists, otherwise it returns IDENT type
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
