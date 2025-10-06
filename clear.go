package main

import (
	"flag"
	"os"
)

func runClearCommand() {
	clearFlagSet := flag.NewFlagSet("clear", flag.ExitOnError)
	translationPath = clearFlagSet.String("translation", "translation.json", "Path to JSON file to read from")
	clearFlagSet.Parse(os.Args[2:])
	printConfig(clearFlagSet)
	loadTranslation()

	for verseKey, verse := range translation {
		verse.Segments = []Segment{}
		translation[verseKey] = verse
	}

	save()
}
