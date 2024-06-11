package engine

import (
	"strconv"

	"github.com/LissaGreense/GO4SQL/token"
)

func getInterfaceValue(t token.Token) (ValueInterface, error) {
	switch t.Type {
	case token.LITERAL:
		castedInteger, err := strconv.Atoi(t.Literal)
		if err != nil {
			return nil, err
		}
		return IntegerValue{Value: castedInteger}, nil
	default:
		return StringValue{Value: t.Literal}, nil
	}
}

func tokenMapper(inputToken token.Type) token.Type {
	switch inputToken {
	case token.TEXT:
		return token.IDENT
	case token.INT:
		return token.LITERAL
	default:
		return inputToken
	}
}

func unique(arr []string) *[]string {
	occurred := map[string]bool{}
	var result []string

	for e := range arr {
		if !occurred[arr[e]] {
			occurred[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return &result
}
