package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func runCheckCommand() {
	checkFlagSet := flag.NewFlagSet("check", flag.ExitOnError)
	translationPath = checkFlagSet.String("translation", "translation.json", "Path to JSON file to read from")
	checkFlagSet.Parse(os.Args[2:])
	printConfig(checkFlagSet)
	loadTranslation()

	for verseKey, verse := range translation {
		if strings.Join(segmentsToStrings(verse.Segments), " ") != verse.Text {
			fmt.Println(verseKey)
		}
	}
}
