package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

type Translation map[string]Verse

type Verse struct {
	Text     string    `json:"t"`
	Segments []Segment `json:"segments,omitempty"`
}

type Segment struct {
	Text  string `json:"t"`
	Words []int  `json:"words,omitempty"`
}

type Metadata map[string]VerseMetadata

type VerseMetadata struct {
	Id          int    `json:"id"`
	SurahNumber int    `json:"surah_number"`
	AyahNumber  int    `json:"ayah_number"`
	VerseKey    string `json:"verse_key"`
	WordsCount  int    `json:"words_count"`
	Text        string `json:"text"`
}

var translationPath *string
var translation Translation
var metadataPath *string
var metadata Metadata

var wordsCountByVerseKey = make(map[string]int)

func printConfig(flagSet *flag.FlagSet) {
	flags := make(map[string]string)
	flagSet.VisitAll(func(f *flag.Flag) {
		flags[f.Name] = f.Value.String()
	})
	log.Printf("Running subcommand %s with config %+v...\n", flagSet.Name(), flags)

}

func loadTranslation() {
	log.Printf("Loading translation from %s...\n", *translationPath)

	translationFile, err := os.ReadFile(*translationPath)
	if err != nil {
		log.Fatalf("error reading translation file: %v\n", err)
	}

	err = json.Unmarshal(translationFile, &translation)
	if err != nil {
		log.Fatalf("error unmarshaling translation JSON: %v\n", err)
	}
}

func loadMetadata() {
	log.Printf("Loading metadata from %s...\n", *metadataPath)

	metadataFile, err := os.ReadFile(*metadataPath)
	if err != nil {
		log.Fatalf("error reading metadata file: %v\n", err)
	}

	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		log.Fatalf("error unmarshaling metadata JSON: %v\n", err)
	}

	for _, verse := range metadata {
		wordsCountByVerseKey[verse.VerseKey] = verse.WordsCount
	}
}

func save() {
	log.Printf("Saving translation to %s...\n", *translationPath)
	translationJSON, err := json.Marshal(translation)
	if err != nil {
		log.Fatalf("error while marshaling translation to JSON: %v", err)
	}

	err = os.WriteFile(*translationPath, translationJSON, 0666)
	if err != nil {
		log.Fatalf("error while writing translation JSON to file: %v", err)
	}
}

func sanityCheck(verseKey, text string, segments []Segment) bool {
	if text != strings.Join(segmentsToStrings(segments), " ") {
		return false
	}

	if len(segments) > wordsCountByVerseKey[verseKey] {
		return false
	}

	re := regexp.MustCompile(`^\[[^\]]+\],?$`)

	for _, segment := range segments {
		if re.MatchString(segment.Text) {
			return false
		}
	}

	return true
}

func segmentsToStrings(segments []Segment) []string {
	return slices.Collect(func(yield func(string) bool) {
		for _, seg := range segments {
			if !yield(seg.Text) {
				return
			}
		}
	})
}

func stringsToSegments(strs []string) []Segment {
	return slices.Collect(func(yield func(Segment) bool) {
		for _, s := range strs {
			if !yield(Segment{Text: s}) {
				return
			}
		}
	})
}
