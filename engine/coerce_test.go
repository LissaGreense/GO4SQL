package engine

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/token"
)

func TestCoerce(t *testing.T) {
	null := token.Token{Type: token.NULL, Literal: "NULL"}
	integer := token.Token{Type: token.LITERAL, Literal: "5"}
	decimal := token.Token{Type: token.LITERAL, Literal: "3.14"}
	overflow := token.Token{Type: token.LITERAL, Literal: "99999999999999999999"}
	str := token.Token{Type: token.IDENT, Literal: "x"}
	asterisk := token.Token{Type: token.ASTERISK, Literal: "*"}

	tests := []struct {
		input     token.Token
		target    SupportedTypes
		wantValue ValueInterface
		wantOk    bool
	}{
		{null, IntType, NullValue{}, true},
		{integer, IntType, IntegerValue{Value: 5}, true},
		{decimal, IntType, nil, false},
		{overflow, IntType, nil, false},
		{str, IntType, nil, false},
		{asterisk, IntType, nil, false},

		{null, StringType, NullValue{}, true},
		{integer, StringType, nil, false},
		{decimal, StringType, nil, false},
		{overflow, StringType, nil, false},
		{str, StringType, StringValue{Value: "x"}, true},
		{asterisk, StringType, StringValue{Value: "*"}, true},

		{null, NullType, NullValue{}, true},
		{integer, NullType, nil, false},
		{decimal, NullType, nil, false},
		{overflow, NullType, nil, false},
		{str, NullType, nil, false},
		{asterisk, NullType, nil, false},

		{null, SupportedTypes(99), NullValue{}, true},
		{integer, SupportedTypes(99), nil, false},
	}

	for _, test := range tests {
		gotValue, gotOk := coerce(test.input, test.target)
		if gotOk != test.wantOk || gotValue != test.wantValue {
			t.Errorf("coerce(%s %q, %s) = (%v, %t), want (%v, %t)",
				test.input.Type, test.input.Literal, test.target, gotValue, gotOk, test.wantValue, test.wantOk)
		}
	}
}
