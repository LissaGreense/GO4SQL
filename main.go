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

	flag.Parse()

	if len(*filePath) > 0 {
		log.Println("Reading file: ", *filePath)

		content, err := ioutil.ReadFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}

		lex := lexer.RunLexer(string(content))

		if *debugMode {
			log.Println("Lexer output:")
			for {
				token := (lex.NextToken())
				if len(token.Literal) == 0 {
					break
				}
				log.Println(token.Type, " : ", token.Literal)
			}
		}

		parser := parser.New(lex)
		sequences := parser.ParseSequence()
		log.Println(sequences.Commands[0])
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
