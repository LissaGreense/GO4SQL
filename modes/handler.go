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
func HandleFileMode(filePath string, engine *engine.DbEngine) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	sequences := bytesToSequences(content)
	fmt.Print(evaluateInEngine(sequences, engine))
}

// HandleStreamMode - Handle GO4SQL use case where client sends input via stdin
func HandleStreamMode(engine *engine.DbEngine) {
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		sequences := bytesToSequences(reader.Bytes())
		fmt.Print(evaluateInEngine(sequences, engine))
	}
	err := reader.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// HandleSocketMode - Handle GO4SQL use case where client sends input via socket protocol
func HandleSocketMode(port int, engine *engine.DbEngine) {
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

		go handleSocketClient(conn, engine)
	}
}
func evaluateInEngine(sequences *ast.Sequence, engineSQL *engine.DbEngine) string {
	commands := sequences.Commands

	result := ""
	for commandIndex, command := range commands {

		switch mappedCommand := command.(type) {
		case *ast.WhereCommand:
			continue
		case *ast.OrderByCommand:
			continue
		case *ast.CreateCommand:
			engineSQL.CreateTable(mappedCommand)
			result += "Table '" + mappedCommand.Name.GetToken().Literal + "' has been created\n"
			continue
		case *ast.InsertCommand:
			engineSQL.InsertIntoTable(mappedCommand)
			result += "Data Inserted\n"
			continue
		case *ast.SelectCommand:
			result += getSelectResponse(commandIndex, &commands, engineSQL, mappedCommand) + "\n"
			continue
		case *ast.DeleteCommand:
			nextCommandIndex := commandIndex + 1
			if nextCommandIndex != len(commands) {
				whereCommand, whereCommandIsValid := commands[nextCommandIndex].(*ast.WhereCommand)

				if whereCommandIsValid {
					engineSQL.DeleteFromTable(mappedCommand, whereCommand)
				}
			}
			result += "Data from '" + mappedCommand.Name.GetToken().Literal + "' has been deleted\n"
			continue
		default:
			log.Fatalf("Unsupported Command detected: %v", command)
		}
	}

	return result
}

func getSelectResponse(commandIndex int, commands *[]ast.Command, engineSQL *engine.DbEngine, selectCommand *ast.SelectCommand) string {
	// TODO: this function should be a method of ast.SelectCommand
	nextCommandIndex := commandIndex + 1

	if nextCommandIndex != len(*commands) {
		whereCommand, whereCommandIsValid := (*commands)[nextCommandIndex].(*ast.WhereCommand)

		// TODO: It cannot be like that. Have to be refactored to tree structure.
		if whereCommandIsValid {
			if nextCommandIndex+1 < len(*commands) {
				orderByCommand, orderByCommandIsValid := (*commands)[nextCommandIndex+1].(*ast.OrderByCommand)

				if orderByCommandIsValid {
					return engineSQL.SelectFromTableWithWhereAndOrderBy(selectCommand, whereCommand, orderByCommand).ToString()
				}
			}

			return engineSQL.SelectFromTableWithWhere(selectCommand, whereCommand).ToString()
		}

		orderByCommand, orderByCommandIsValid := (*commands)[nextCommandIndex].(*ast.OrderByCommand)

		if orderByCommandIsValid {
			return engineSQL.SelectFromTableWithOrderBy(selectCommand, orderByCommand).ToString()
		}
	}

	return engineSQL.SelectFromTable(selectCommand).ToString()
}

func bytesToSequences(content []byte) *ast.Sequence {
	lex := lexer.RunLexer(string(content))
	parserInstance := parser.New(lex)
	sequences := parserInstance.ParseSequence()

	return sequences
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
			log.Fatal("Error:", err)
		}
		sequences := bytesToSequences(buffer)
		commandResult := evaluateInEngine(sequences, engine)

		if len(commandResult) > 0 {
			_, err = conn.Write([]byte(commandResult))
		}

		if err != nil {
			log.Fatal("Error:", err)
		}

		fmt.Printf("Received: %s\n", buffer[:n])
	}
}
