package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/engine"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
)

func main() {
	filePath := flag.String("file", "", "Provide a path to the .sql file")
	streamMode := flag.Bool("stream", false, "Use to redirect stdin to stdout")
	flag.Parse()
	engineSQL := engine.New()

	if len(*filePath) > 0 {
		content, err := ioutil.ReadFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}

		sequences := bytesToSequences(content)
		evaluateInEngine(sequences, engineSQL)
	} else if *streamMode {

		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			sequences := bytesToSequences(reader.Bytes())
			evaluateInEngine(sequences, engineSQL)
		}
		err := reader.Err()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("No mode has been providing. Exiting.")
	}
}

func bytesToSequences(content []byte) *ast.Sequence {
	lex := lexer.RunLexer(string(content))
	parserInstance := parser.New(lex)
	sequences := parserInstance.ParseSequence()

	return sequences
}

func evaluateInEngine(sequences *ast.Sequence, engineSQL *engine.DbEngine) {
	commands := sequences.Commands
	for commandIndex, command := range commands {

		_, whereCommandIsValid := command.(*ast.WhereCommand)
		if whereCommandIsValid {
			continue
		}

		createCommand, createCommandIsValid := command.(*ast.CreateCommand)
		if createCommandIsValid {
			engineSQL.CreateTable(createCommand)
			continue
		}

		insertCommand, insertCommandIsValid := command.(*ast.InsertCommand)
		if insertCommandIsValid {
			engineSQL.InsertIntoTable(insertCommand)
			continue
		}

		selectCommand, selectCommandIsValid := command.(*ast.SelectCommand)
		if selectCommandIsValid {
			result := getSelectResponse(commandIndex, commands, engineSQL, selectCommand)
			fmt.Println(result)
			continue
		}

	}
}

func getSelectResponse(commandIndex int, commands []ast.Command, engineSQL *engine.DbEngine, selectCommand *ast.SelectCommand) string {
	nextCommandIndex := commandIndex + 1

	if nextCommandIndex != len(commands) {
		whereCommand, whereCommandIsValid := commands[nextCommandIndex].(*ast.WhereCommand)

		if whereCommandIsValid {
			return engineSQL.SelectFromTableWithWhere(selectCommand, whereCommand).ToString()
		}
	}

	return engineSQL.SelectFromTable(selectCommand).ToString()
}
