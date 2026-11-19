package engine

import (
	"strconv"

	"github.com/LissaGreense/GO4SQL/token"
)

// coerce converts a value token to the target type. ok=true means the value is of the
// target type or NULL; NOT NULL is the caller's concern. It never panics or builds
// errors, so every caller can word its own error.
func coerce(t token.Token, target SupportedTypes) (ValueInterface, bool) {
	if t.Type == token.NULL {
		return NullValue{}, true
	}

	switch target {
	case IntType:
		if t.Type != token.LITERAL {
			return nil, false
		}
		integer, err := strconv.Atoi(t.Literal)
		if err != nil {
			return nil, false
		}
		return IntegerValue{Value: integer}, true
	case StringType:
		if t.Type == token.LITERAL {
			return nil, false
		}
		return StringValue{Value: t.Literal}, true
	default:
		return nil, false
	}
}

// declaredTarget maps a column's declared type to the type its values must coerce to.
// Unknown types fail closed: only NULL fits them.
func declaredTarget(columnType token.Type) SupportedTypes {
	switch columnType {
	case token.TEXT:
		return StringType
	case token.INT:
		return IntType
	default:
		return NullType
	}
}

// displayValue renders a value token for error messages, re-adding the quotes the lexer
// strips from strings so 'x' and x stay distinguishable from the literal 5.
func displayValue(t token.Token) string {
	if t.Type == token.LITERAL {
		return t.Literal
	}
	return "'" + t.Literal + "'"
}

// inferTarget is the type a filter literal has when there is no column to take a type from.
// It mirrors how row evaluation converts literals, so the check and the scan agree.
func inferTarget(t token.Token) SupportedTypes {
	switch t.Type {
	case token.LITERAL:
		return IntType
	case token.NULL:
		return NullType
	default:
		return StringType
	}
}

// literalValue converts a filter literal for row evaluation. Inferring its type is safe only
// because checkFilter already proved it matches what it is compared against. Trip-wire: 5 against
// a future FLOAT column would pass the check, infer Int, and silently never match.
func literalValue(t token.Token) (ValueInterface, bool) {
	return coerce(t, inferTarget(t))
}
