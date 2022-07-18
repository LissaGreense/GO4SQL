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
		log.Println("Reading file: ", *filePath)

		content, err := ioutil.ReadFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}

		sequences := bytesToSequences(content)

		log.Println("Parser output:")
		for _, command := range sequences.Commands {
			log.Println(command)
		}
		evaluateInEngine(sequences, engineSQL)
	} else if *streamMode {
		log.Println("Reading from stream")

		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			sequences := bytesToSequences(reader.Bytes())
			for _, command := range sequences.Commands {
				fmt.Println(command)
			}
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
	for _, command := range sequences.Commands {

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
			fmt.Println(engineSQL.SelectFromTable(selectCommand))
		}
	}
}
