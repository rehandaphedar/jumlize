package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "segment":
		runSegmentCommand()
	case "check":
		runCheckCommand()
	case "clear":
		runClearCommand()
	default:
		printHelp()
	}
}

func printHelp() {
	log.Println("Invalid command")
	fmt.Println(`Usage: jumlize [subcommand] [flags]
Subcommands:
- segment
- check
- clear`)
}
