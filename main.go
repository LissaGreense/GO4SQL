package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
)

func main() {
	filePath := flag.String("file", "", "Provide a path to the .sql file")
	streamMode := flag.Bool("stream", false, "Use to redirect stdin to stdout")
	debugMode := flag.Bool("debug", false, "Use to enable debug mode")

	// file := "test_file"
	// filePath = &file
	flag.Parse()

	if len(*filePath) > 0 {
		log.Println("Reading file: ", *filePath)

		content, err := ioutil.ReadFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}

		lex := lexer.RunLexer(string(content))

		// FIXME: printing lexer breaks parser input - using NextToken() doesn't let to start reading from the beginning again
		if *debugMode {
			log.Println("Lexer output:")
			for {
				token := lex.NextToken()
				if len(token.Literal) == 0 {
					break
				}
				log.Println(token.Type, " : ", token.Literal)
			}
		}

		parserInstance := parser.New(lex)
		sequences := parserInstance.ParseSequence()

		// TODO: Print only when debug mode is turned on - can be done after lexer print fix
		log.Println("Parser output:")
		for _, command := range sequences.Commands {
			log.Println(command)
		}
	} else if *streamMode {
		log.Println("Reading from stream")

		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			fmt.Println(reader.Text())
		}

		err := reader.Err()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("No mode has been providing. Exiting.")
	}
}
