package engine

import (
	"log"
	"strconv"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

type DbEngine struct {
	Tables map[string]map[string]*Column
}

type Column struct {
	Type   token.Token
	Values []interface{}
}

// New Return new DbEngine struct
func New() *DbEngine {
	engine := &DbEngine{}
	engine.Tables = make(map[string]map[string]*Column)
	return engine
}

func (engine *DbEngine) CreateTable(command *ast.CreateCommand) {
	_, exist := engine.Tables[command.Name.Token.Literal]

	if exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " already exist!")
	}

	engine.Tables[command.Name.Token.Literal] = make(map[string]*Column)

	for i, columnName := range command.ColumnNames {
		column := Column{Type: command.ColumnTypes[i]}
		column.Values = make([]interface{}, 0)
		engine.Tables[command.Name.Token.Literal][columnName] = &column
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
	for columnName, columnStruct := range table {
		if columnStruct.Type.Type != engine.Tables[command.Name.Token.Literal][columnName].Type.Type {
			log.Fatal("Invalid Token Type in Insert Command, expecting: " + engine.Tables[command.Name.Token.Literal][columnName].Type.Type + ", got: " + columnStruct.Type.Type)
		}
		engine.Tables[command.Name.Token.Literal][columnName].Values = append(engine.Tables[command.Name.Token.Literal][columnName].Values, command.Values[i].Literal)
		i++
	}
}
