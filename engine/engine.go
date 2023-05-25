package engine

import (
	"log"
	"strconv"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

type DbEngine struct {
	Tables map[string]*Table
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

func (engine *DbEngine) SelectFromTable(command *ast.SelectCommand) *Table {
	table, exist := engine.Tables[command.Name.Token.Literal]

	if !exist {
		log.Fatal("Table with the name of " + command.Name.Token.Literal + " doesn't exist!")
	}

	return engine.selectFromProvidedTable(command, table)
}

func (engine *DbEngine) selectFromProvidedTable(command *ast.SelectCommand, table *Table) *Table {
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

func (engine *DbEngine) SelectFromTableWithWhere(selectcommand *ast.SelectCommand, whereCommand *ast.WhereCommand) *Table {
	table, exist := engine.Tables[selectcommand.Name.Token.Literal]

	if !exist {
		log.Fatal("Table with the name of " + selectcommand.Name.Token.Literal + " doesn't exist!")
	}

	columns := table.Columns

	conditionalColumnName := whereCommand.Expression.Left
	conditionalOperation := whereCommand.Expression.OperationToken
	conditionalValue := whereCommand.Expression.Right

	filteredTable := &Table{Columns: []*Column{}}

	conditionalColumnIndex := -1

	// use create table after decorator implementation
	for i := range columns {
		filteredTable.Columns = append(filteredTable.Columns,
			&Column{
				Type:   columns[i].Type,
				Values: make([]ValueInterface, 0),
				Name:   columns[i].Name,
			})

		if conditionalColumnName.Literal == columns[i].Name {
			conditionalColumnIndex = i
		}
	}

	if conditionalColumnIndex == -1 {
		log.Fatal("In table" + selectcommand.Name.Token.Literal + ", column with the name of" + conditionalColumnName.Literal + " doesn't exist!")
	}

	for rowIndex, value := range columns[conditionalColumnIndex].Values {
		switch conditionalOperation.Type {
		case token.EQUAL:
			if value == getInterfaceValue(conditionalValue) {
				filteredTable.appendRow(columns, rowIndex)
			}
		case token.NOT:
			if value != getInterfaceValue(conditionalValue) {
				filteredTable.appendRow(columns, rowIndex)
			}
		default:
			log.Fatal("Operation '" + conditionalOperation.Literal + "' provided in WHERE command isn't allowed!")
		}
	}

	return engine.selectFromProvidedTable(selectcommand, filteredTable)
}
