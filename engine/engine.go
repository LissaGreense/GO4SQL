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

	if len(command.Values) != len(engine.Tables[command.Name.Token.Literal]) {
		log.Fatal("Invalid number of parameters in insert, should be: " + strconv.Itoa(len(engine.Tables[command.Name.Token.Literal])) + ", but got: " + strconv.Itoa(len(command.Values)))
	}

	i := 0
	for columnName, columnType := range table {
		if columnType.Type.Type != engine.Tables[command.Name.Token.Literal][columnName].Type.Type {
			log.Fatal("Invalid Token Type in Insert Command, expecting: " + engine.Tables[command.Name.Token.Literal][columnName].Type.Type + ", got: " + columnType.Type.Type)
		}
		engine.Tables[command.Name.Token.Literal][columnName].Values = append(engine.Tables[command.Name.Token.Literal][columnName].Values, command.Values[i].Literal)
		i++
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
