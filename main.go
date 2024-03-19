package main

import (
	"flag"
	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/engine"
	"github.com/LissaGreense/GO4SQL/modes"
	"log"
)

func main() {
	filePath := flag.String("file", "", "Provide a path to the .sql file")
	streamMode := flag.Bool("stream", false, "Use to redirect stdin to stdout")
	socketMode := flag.Bool("socket", false, "Use to start socket server")
	port := flag.Int("port", 1433, "States on which port socket server will listen")

	flag.Parse()
	engineSQL := engine.New()

	if len(*filePath) > 0 {
		modes.HandleFileMode(*filePath, engineSQL, evaluateInEngine)
	} else if *streamMode {
		modes.HandleStreamMode(engineSQL, evaluateInEngine)
	} else if *socketMode {
		modes.HandleSocketMode(*port, engineSQL, evaluateInEngine)
	} else {
		log.Println("No mode has been providing. Exiting.")
	}
}

func evaluateInEngine(sequences *ast.Sequence, engineSQL *engine.DbEngine) string {
	commands := sequences.Commands

	result := ""
	for commandIndex, command := range commands {

		// TODO: Check if those statements are necessary
		_, whereCommandIsValid := command.(*ast.WhereCommand)
		if whereCommandIsValid {
			continue
		}

		_, orderByCommandIsValid := command.(*ast.OrderByCommand)
		if orderByCommandIsValid {
			continue
		}

		createCommand, createCommandIsValid := command.(*ast.CreateCommand)
		if createCommandIsValid {
			engineSQL.CreateTable(createCommand)
			result += "Table '" + createCommand.Name.GetToken().Literal + "' has been created\n"
			continue
		}

		insertCommand, insertCommandIsValid := command.(*ast.InsertCommand)
		if insertCommandIsValid {
			engineSQL.InsertIntoTable(insertCommand)
			result += "Data Inserted\n"
			continue
		}

		selectCommand, selectCommandIsValid := command.(*ast.SelectCommand)
		if selectCommandIsValid {
			result += getSelectResponse(commandIndex, commands, engineSQL, selectCommand) + "\n"
			continue
		}

		deleteCommand, deleteCommandIsValid := command.(*ast.DeleteCommand)
		if deleteCommandIsValid {
			nextCommandIndex := commandIndex + 1

			if nextCommandIndex != len(commands) {
				whereCommand, whereCommandIsValid := commands[nextCommandIndex].(*ast.WhereCommand)

				if whereCommandIsValid {
					engineSQL.DeleteFromTable(deleteCommand, whereCommand)
				}
			}
			result += "Data from '" + deleteCommand.Name.GetToken().Literal + "' has been deleted\n"
			continue
		}

	}

	return result
}

func getSelectResponse(commandIndex int, commands []ast.Command, engineSQL *engine.DbEngine, selectCommand *ast.SelectCommand) string {
	nextCommandIndex := commandIndex + 1

	if nextCommandIndex != len(commands) {
		whereCommand, whereCommandIsValid := commands[nextCommandIndex].(*ast.WhereCommand)

		// TODO: It cannot be like that. Have to be refactored to tree structure.
		if whereCommandIsValid {
			if nextCommandIndex+1 < len(commands) {
				orderByCommand, orderByCommandIsValid := commands[nextCommandIndex+1].(*ast.OrderByCommand)

				if orderByCommandIsValid {
					return engineSQL.SelectFromTableWithWhereAndOrderBy(selectCommand, whereCommand, orderByCommand).ToString()
				}
			}

			return engineSQL.SelectFromTableWithWhere(selectCommand, whereCommand).ToString()
		}

		orderByCommand, orderByCommandIsValid := commands[nextCommandIndex].(*ast.OrderByCommand)

		if orderByCommandIsValid {
			return engineSQL.SelectFromTableWithOrderBy(selectCommand, orderByCommand).ToString()
		}
	}

	return engineSQL.SelectFromTable(selectCommand).ToString()
}
