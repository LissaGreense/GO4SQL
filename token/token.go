package token

type Type string

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
	CREATE = "CREATE"
	TABLE  = "TABLE"
	INSERT = "INSERT"
	INTO   = "INTO"
	VALUES = "VALUES"
	SELECT = "SELECT"
	FROM   = "FROM"
	WHERE  = "WHERE"

	// EQUAL - Logical operations
	EQUAL = "EQUAL"
	NOT   = "NOT"

	// TEXT - Data types
	TEXT = "TEXT"
	INT  = "INT"

	// ILLEGAL - System
	ILLEGAL = "ILLEGAL"
)

type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"TEXT":   TEXT,
	"INT":    INT,
	"CREATE": CREATE,
	"TABLE":  TABLE,
	"INSERT": INSERT,
	"INTO":   INTO,
	"SELECT": SELECT,
	"FROM":   FROM,
	"VALUES": VALUES,
	"WHERE":  WHERE,
	"EQUAL":  EQUAL,
	"NOT":    NOT,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
