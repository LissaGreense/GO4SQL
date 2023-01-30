package engine

import (
	"log"
	"strconv"

	"github.com/LissaGreense/GO4SQL/token"
)

func getInterfaceValue(t token.Token) ValueInterface {
	switch t.Type {
	case token.INT:
		castedInteger, err := strconv.Atoi(t.Literal)
		if err != nil {
			log.Fatal("Cannot cast \"" + t.Literal + "\" to Integer")
		}
		return IntegerValue{Value: castedInteger}
	default:
		return StringValue{Value: t.Literal}
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

func unique(arr []string) []string {
	occurred := map[string]bool{}
	var result []string

	for e := range arr {
		if !occurred[arr[e]] {
			occurred[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}
