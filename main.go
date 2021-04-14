package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := flag.String("file", "", "Provide a path to the .sql file")
	streamMode := flag.Bool("stream", false, "Use to redirect stdin to stdout")

	flag.Parse()

	fmt.Println("Provided file: " + *filePath)

	if *streamMode {
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			fmt.Println(reader.Text())
		}

		err := reader.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

}
