package main

import (
	"flag"
	"fmt"
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
	var err error

	if len(*filePath) > 0 {
		err = modes.HandleFileMode(*filePath, engineSQL)
	} else if *streamMode {
		err = modes.HandleStreamMode(engineSQL)
	} else if *socketMode {
		modes.HandleSocketMode(*port, engineSQL)
	} else {
		err = fmt.Errorf("no mode has been providing, exiting")
	}

	if err != nil {
		log.Fatal(err)
	}
}
