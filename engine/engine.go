package engine

import (
	"log"
	"strconv"
	"strings"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

type DbEngine struct {
	Tables map[string]map[int]*Column
}

type Column struct {
	Name   string
	Type   token.Token
	Values []ValueInterface
}

// New Return new DbEngine struct
func New() *DbEngine {
	engine := &DbEngine{}
	engine.Tables = make(map[string]map[int]*Column)
	return engine
}

func (engine *DbEngine) CreateTable(command *ast.CreateCommand) {
	_, exist := engine.Tables[command.Name.Token.Literal]

	if exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " already exist!")
	}

	engine.Tables[command.Name.Token.Literal] = make(map[int]*Column)

	for i, columnName := range command.ColumnNames {
		column := Column{Type: command.ColumnTypes[i]}
		column.Values = make([]ValueInterface, 0)
		column.Name = columnName
		engine.Tables[command.Name.Token.Literal][i] = &column
	}
}

func (engine *DbEngine) InsertIntoTable(command *ast.InsertCommand) {
	table, exist := engine.Tables[command.Name.Token.Literal]

	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	if len(command.Values) != len(table) {
		log.Fatal("Invalid number of parameters in insert, should be: " + strconv.Itoa(len(table)) + ", but got: " + strconv.Itoa(len(command.Values)))
	}

	for i := 0; i < len(table); i++ {
		expectedToken := tokenMapper(table[i].Type.Type)
		if expectedToken != command.Values[i].Type {
			log.Fatal("Invalid Token Type in Insert Command, expecting: " + expectedToken + ", got: " + command.Values[i].Type)
		}
		table[i].Values = append(table[i].Values, getInterfaceValue(command.Values[i]))
	}
}

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

func (engine *DbEngine) SelectFromTable(command *ast.SelectCommand) string {
	table, exist := engine.Tables[command.Name.Token.Literal]

	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	wantedColumnNames := make([]string, 0)
	if command.Space[0].Type == token.ASTERISK {
		for i := 0; i < len(table); i++ {
			wantedColumnNames = append(wantedColumnNames, table[i].Name)
		}
		return extractColumnContent(table, wantedColumnNames)
	} else {
		for i := 0; i < len(command.Space); i++ {
			wantedColumnNames = append(wantedColumnNames, command.Space[i].Literal)
		}
		return extractColumnContent(table, unique(wantedColumnNames))
	}
}

func extractColumnContent(table map[int]*Column, wantedColumnNames []string) string {
	mappedIndexes := make([]int, 0)
	for wantedColumnIndex := 0; wantedColumnIndex < len(wantedColumnNames); wantedColumnIndex++ {
		for columnNameIndex := 0; columnNameIndex < len(table); columnNameIndex++ {
			if table[columnNameIndex].Name == wantedColumnNames[wantedColumnIndex] {
				mappedIndexes = append(mappedIndexes, columnNameIndex)
				break
			}
			if columnNameIndex == len(table)-1 {
				log.Fatal("Provided column name: " + wantedColumnNames[wantedColumnIndex] + "doesn't exist")
			}
		}
	}
	result := ""
	for i := 0; i < len(mappedIndexes); i++ {
		result += table[mappedIndexes[i]].Name
		result += "|"
	}
	result = strings.TrimSuffix(result, "|")
	result += "\n"

	rowsCount := len(table[0].Values)

	for iRow := 0; iRow < rowsCount; iRow++ {
		for iColumn := 0; iColumn < len(mappedIndexes); iColumn++ {
			if table[mappedIndexes[iColumn]].Type.Literal == token.TEXT {
				result += "'" + table[mappedIndexes[iColumn]].Values[iRow].ToString() + "'"
			} else {
				result += table[mappedIndexes[iColumn]].Values[iRow].ToString()
			}
			result += "|"
		}
		result = strings.TrimSuffix(result, "|")
		result += "\n"
	}
	result = strings.TrimSuffix(result, "\n")
	return result
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
