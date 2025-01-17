package modes

import (
	"bufio"
	"fmt"
	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/engine"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
	"log"
	"net"
	"os"
	"strconv"
)

// HandleFileMode - Handle GO4SQL use case where client sends input via text file
func HandleFileMode(filePath string, engine *engine.DbEngine) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	sequences, err := bytesToSequences(content)
	if err != nil {
		return err
	}
	evaluate, err := engine.Evaluate(sequences)
	if err != nil {
		return err
	}
	fmt.Print(evaluate)
	return nil
}

// HandleStreamMode - Handle GO4SQL use case where client sends input via stdin
func HandleStreamMode(engine *engine.DbEngine) error {
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		sequences, err := bytesToSequences(reader.Bytes())
		if err != nil {
			fmt.Print(err)
		} else {
			evaluate, err := engine.Evaluate(sequences)
			if err != nil {
				fmt.Print(err)
			} else {
				fmt.Print(evaluate)
			}
		}
	}
	return reader.Err()
}

// HandleSocketMode - Handle GO4SQL use case where client sends input via socket protocol
func HandleSocketMode(port int, engine *engine.DbEngine) {
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	log.Printf("Starting Socket Server on %d port\n", port)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal("Error:", err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleSocketClient(conn, engine)
	}
}

func bytesToSequences(content []byte) (*ast.Sequence, error) {
	lex := lexer.RunLexer(string(content))
	parserInstance := parser.New(lex)
	sequences, err := parserInstance.ParseSequence()
	return sequences, err
}

func handleSocketClient(conn net.Conn, engine *engine.DbEngine) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Error:", err)
		}
	}(conn)

	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			log.Fatal(err.Error())
		}
		sequences, err := bytesToSequences(buffer)

		if err != nil {
			log.Fatal(err.Error())
		}

		commandResult, err := engine.Evaluate(sequences)

		if err != nil {
			_, err = conn.Write([]byte(err.Error()))
		} else if len(commandResult) > 0 {
			_, err = conn.Write([]byte(commandResult))
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Received: %s\n", buffer[:n])
	}
}
