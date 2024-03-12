package modes

import (
	"bufio"
	"fmt"
	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/engine"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/parser"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

// HandleFileMode - Handle GO4SQL use case where client sends input via text file
func HandleFileMode(filePath string, engine *engine.DbEngine, evaluate func(sequences *ast.Sequence, engineSQL *engine.DbEngine) string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	sequences := bytesToSequences(content)
	fmt.Println(evaluate(sequences, engine))
}

// HandleStreamMode - Handle GO4SQL use case where client sends input via stdin
func HandleStreamMode(engine *engine.DbEngine, evaluate func(sequences *ast.Sequence, engineSQL *engine.DbEngine) string) {
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		sequences := bytesToSequences(reader.Bytes())
		fmt.Println(evaluate(sequences, engine))
	}
	err := reader.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// HandleSocketMode - Handle GO4SQL use case where client sends input via socket protocol
func HandleSocketMode(port int, engine *engine.DbEngine, evaluate func(sequences *ast.Sequence, engineSQL *engine.DbEngine) string) {
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	log.Printf("Starting Socket Server on %d port\n", port)

	if err != nil {
		log.Fatal("Error:", err)
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

		go handleSocketClient(conn, engine, evaluate)
	}
}

func bytesToSequences(content []byte) *ast.Sequence {
	lex := lexer.RunLexer(string(content))
	parserInstance := parser.New(lex)
	sequences := parserInstance.ParseSequence()

	return sequences
}

func handleSocketClient(conn net.Conn, engine *engine.DbEngine, evaluate func(sequences *ast.Sequence, engineSQL *engine.DbEngine) string) {
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
			log.Fatal("Error:", err)
		}
		sequences := bytesToSequences(buffer)
		commandResult := evaluate(sequences, engine)

		if len(commandResult) > 0 {
			_, err = conn.Write([]byte(commandResult))
		}

		if err != nil {
			log.Fatal("Error:", err)
		}

		fmt.Printf("Received: %s\n", buffer[:n])
	}
}
