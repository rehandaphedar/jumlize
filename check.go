package main

import (
	"flag"
	"fmt"
	"os"
)

func runCheckCommand() {
	checkFlagSet := flag.NewFlagSet("check", flag.ExitOnError)
	translationPath = checkFlagSet.String("translation", "en-sahih-international-simple.json", "Path to JSON file to read from")
	metadataPath = checkFlagSet.String("metadata", "quran-metadata-ayah.json", "Path to JSON file to read from")
	checkFlagSet.Parse(os.Args[2:])
	printConfig(checkFlagSet)

	loadTranslation()
	loadMetadata()

	for verseKey, verse := range translation {
		if !sanityCheck(verseKey, verse.Text, verse.Segments) {
			fmt.Println(verseKey)
		}
	}
}
