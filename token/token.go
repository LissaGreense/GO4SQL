package token

type TokenType string

const (
	// Operators
	ASTERISK = "*"

	// Identifiers + literals
	IDENT   = "IDENT"   // tab, car, apple...
	LITERAL = "LITERAL" // 1343456

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Special tokens
	EOF        = ""
	APOSTROPHE = "'"

	// Paren
	LPAREN = "("
	RPAREN = ")"

	// Keywords
	CREATE = "CREATE"
	TABLE  = "TABLE"
	INSERT = "INSERT"
	INTO   = "INTO"
	VALUES = "VALUES"
	SELECT = "SELECT"
	FROM   = "FROM"

	// Data types
	TEXT   = "TEXT"
	INT    = "INT"

	// System
	ILLEGAL = "ILLEGAL"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"TEXT":   TEXT,
	"INT":    INT,
	"CREATE": CREATE,
	"TABLE":  TABLE,
	"INSERT": INSERT,
	"INTO":   INTO,
	"SELECT": SELECT,
	"FROM":   SELECT,
	"VALUES": VALUES,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
