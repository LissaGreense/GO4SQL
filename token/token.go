package token

type Type string

// Token - contains
type Token struct {
	Type    Type
	Literal string
}

const (
	// Operators
	ASTERISK = "*"

	// Identifiers & Literals
	IDENT   = "IDENT"   // e.g., table, column names
	LITERAL = "LITERAL" // e.g., numeric or string literals

	// Delimiters
	COMMA      = ","
	SEMICOLON  = ";"
	APOSTROPHE = "'"

	// Parentheses
	LPAREN = "("
	RPAREN = ")"

	// Special Tokens
	EOF     = ""
	ILLEGAL = "ILLEGAL"

	// Commands
	CREATE = "CREATE"
	DROP   = "DROP"
	TABLE  = "TABLE"
	INSERT = "INSERT"
	INTO   = "INTO"
	VALUES = "VALUES"
	SELECT = "SELECT"
	DELETE = "DELETE"
	UPDATE = "UPDATE"

	// Clauses
	FROM     = "FROM"
	WHERE    = "WHERE"
	ORDER    = "ORDER"
	BY       = "BY"
	ASC      = "ASC"
	DESC     = "DESC"
	LIMIT    = "LIMIT"
	OFFSET   = "OFFSET"
	SET      = "SET"
	DISTINCT = "DISTINCT"
	TO       = "TO"

	// Joins
	JOIN  = "JOIN"
	INNER = "INNER"
	FULL  = "FULL"
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
	ON    = "ON"

	// Aggregates
	MIN   = "MIN"
	MAX   = "MAX"
	COUNT = "COUNT"
	SUM   = "SUM"
	AVG   = "AVG"

	// Logical
	EQUAL = "EQUAL"
	NOT   = "NOT"
	AND   = "AND"
	OR    = "OR"
	TRUE  = "TRUE"
	FALSE = "FALSE"
	IN    = "IN"
	NOTIN = "NOTIN"
	NULL  = "NULL"

	// Data Types
	TEXT = "TEXT"
	INT  = "INT"
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
	"MIN":      MIN,
	"MAX":      MAX,
	"COUNT":    COUNT,
	"SUM":      SUM,
	"AVG":      AVG,
	"IN":       IN,
	"NOTIN":    NOTIN,
	"TO":       TO,
	"VALUES":   VALUES,
	"WHERE":    WHERE,
	"EQUAL":    EQUAL,
	"NOT":      NOT,
	"AND":      AND,
	"OR":       OR,
	"TRUE":     TRUE,
	"FALSE":    FALSE,
	"NULL":     NULL,
}

// LookupIdent - Return keyword type from defined list if exists, otherwise it returns IDENT type
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
