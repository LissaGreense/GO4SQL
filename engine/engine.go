package engine

import (
	"fmt"
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
	Values []interface{}
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
		column.Values = make([]interface{}, 0)
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
		table[i].Values = append(table[i].Values, command.Values[i].Literal)
	}
}

func (engine *DbEngine) SelectFromTable(command *ast.SelectCommand) string {
	table, exist := engine.Tables[command.Name.Token.Literal]
	result := ""

	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	if command.Space[0].Type == token.ASTERISK {
		for i := 0; i < len(table); i++ {
			result += table[i].Name
			result += "|"
		}
		result = strings.TrimSuffix(result, "|")
		result += "\n"

		rowsCount := len(table[0].Values)

		for iRow := 0; iRow < rowsCount; iRow++ {
			for iColumn := 0; iColumn < len(table); iColumn++ {
				if table[iColumn].Type.Literal == token.TEXT {
					result += "'" + fmt.Sprintf("%v", table[iColumn].Values[iRow]) + "'"
				} else {
					result += fmt.Sprintf("%v", table[iColumn].Values[iRow])
				}
				result += "|"
			}
			result = strings.TrimSuffix(result, "|")
			result += "\n"
		}
		result = strings.TrimSuffix(result, "\n")
	}

	return result
}

func tokenMapper(inputToken token.TokenType) token.TokenType {
	switch inputToken {
	case token.TEXT:
		return token.IDENT
	case token.INT:
		return token.LITERAL
	default:
		return inputToken
	}
}
