package engine

import (
	"log"
	"strconv"
	"strings"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

type DbEngine struct {
	Tables map[string]*Table
}

type Table struct {
	Columns []*Column
}

type Column struct {
	Name   string
	Type   token.Token
	Values []ValueInterface
}

// New Return new DbEngine struct
func New() *DbEngine {
	engine := &DbEngine{}
	engine.Tables = make(map[string]*Table)
	return engine
}

func (engine *DbEngine) CreateTable(command *ast.CreateCommand) {
	_, exist := engine.Tables[command.Name.Token.Literal]

	if exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " already exist!")
	}

	engine.Tables[command.Name.Token.Literal] = &Table{Columns: []*Column{}}
	for i, columnName := range command.ColumnNames {
		engine.Tables[command.Name.Token.Literal].Columns = append(engine.Tables[command.Name.Token.Literal].Columns,
			&Column{
				Type:   command.ColumnTypes[i],
				Values: make([]ValueInterface, 0),
				Name:   columnName,
			})
	}
}

func (engine *DbEngine) InsertIntoTable(command *ast.InsertCommand) {
	table, exist := engine.Tables[command.Name.Token.Literal]
	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	columns := table.Columns

	if len(command.Values) != len(columns) {
		log.Fatal("Invalid number of parameters in insert, should be: " + strconv.Itoa(len(columns)) + ", but got: " + strconv.Itoa(len(columns)))
	}

	for i := 0; i < len(columns); i++ {
		expectedToken := tokenMapper(columns[i].Type.Type)
		if expectedToken != command.Values[i].Type {
			log.Fatal("Invalid Token Type in Insert Command, expecting: " + expectedToken + ", got: " + command.Values[i].Type)
		}
		columns[i].Values = append(columns[i].Values, getInterfaceValue(command.Values[i]))
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

func (engine *DbEngine) SelectFromTable(command *ast.SelectCommand) *Table {
	table, exist := engine.Tables[command.Name.Token.Literal]

	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	columns := table.Columns

	wantedColumnNames := make([]string, 0)
	if command.Space[0].Type == token.ASTERISK {
		for i := 0; i < len(columns); i++ {
			wantedColumnNames = append(wantedColumnNames, columns[i].Name)
		}
		return extractColumnContent(columns, wantedColumnNames)
	} else {
		for i := 0; i < len(command.Space); i++ {
			wantedColumnNames = append(wantedColumnNames, command.Space[i].Literal)
		}
		return extractColumnContent(columns, unique(wantedColumnNames))
	}
}

func extractColumnContent(columns []*Column, wantedColumnNames []string) *Table {
	selectedTable := &Table{Columns: make([]*Column, 0)}
	mappedIndexes := make([]int, 0)
	for wantedColumnIndex := 0; wantedColumnIndex < len(wantedColumnNames); wantedColumnIndex++ {
		for columnNameIndex := 0; columnNameIndex < len(columns); columnNameIndex++ {
			if columns[columnNameIndex].Name == wantedColumnNames[wantedColumnIndex] {
				mappedIndexes = append(mappedIndexes, columnNameIndex)
				break
			}
			if columnNameIndex == len(columns)-1 {
				log.Fatal("Provided column name: " + wantedColumnNames[wantedColumnIndex] + "doesn't exist")
			}
		}
	}

	for i := 0; i < len(mappedIndexes); i++ {
		selectedTable.Columns = append(selectedTable.Columns, &Column{
			Name:   columns[mappedIndexes[i]].Name,
			Type:   columns[mappedIndexes[i]].Type,
			Values: make([]ValueInterface, 0),
		})
	}
	rowsCount := len(columns[0].Values)

	for iRow := 0; iRow < rowsCount; iRow++ {
		for iColumn := 0; iColumn < len(mappedIndexes); iColumn++ {
			selectedTable.Columns[iColumn].Values = append(selectedTable.Columns[iColumn].Values, columns[mappedIndexes[iColumn]].Values[iRow])
		}
	}
	return selectedTable
}

func (table *Table) IsEqual(secondTable *Table) bool {
	if len(table.Columns) != len(secondTable.Columns) {
		return false
	}

	for i := 0; i < len(table.Columns); i++ {
		if table.Columns[i].Name != secondTable.Columns[i].Name {
			return false
		}
		if table.Columns[i].Type.Literal != secondTable.Columns[i].Type.Literal {
			return false
		}
		if string(table.Columns[i].Type.Type) != string(secondTable.Columns[i].Type.Type) {
			return false
		}
		if len(table.Columns[i].Values) != len(secondTable.Columns[i].Values) {
			return false
		}
		for j := 0; j < len(table.Columns[i].Values); j++ {
			if table.Columns[i].Values[j].ToString() != secondTable.Columns[i].Values[j].ToString() {
				return false
			}
		}
	}

	return true
}

func (table *Table) ToString() string {
	result := ""

	for i := 0; i < len(table.Columns); i++ {
		result += table.Columns[i].Name
		result += "|"
	}
	result = strings.TrimSuffix(result, "|")
	result += "\n"

	rowsCount := len(table.Columns[0].Values)

	for iRow := 0; iRow < rowsCount; iRow++ {

		for iColumn := 0; iColumn < len(table.Columns); iColumn++ {
			if table.Columns[iColumn].Type.Literal == token.TEXT {
				result += "'" + table.Columns[iColumn].Values[iRow].ToString() + "'"
			} else {
				result += table.Columns[iColumn].Values[iRow].ToString()
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
