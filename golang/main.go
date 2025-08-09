package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"git.sr.ht/~rehandaphedar/lafzize/v3/pkg/api"
	"google.golang.org/genai"
)

type Output struct {
	Segments map[string][]string `json:"segments"`
	Errors   map[string][]string `json:"errors"`
}

type Model struct {
	Slug    string
	RPM     int
	Timeout time.Duration
}

func main() {
	if os.Args[1] == "check" {
		check()
		os.Exit(0)
	}
	apiKey := flag.String("gemini_api_key", "", "Gemini API key")
	dataPath := flag.String("data", "data.json", "Path to JSON file to read from")
	outputPath := flag.String("output", "output.json", "Path to JSON file to write to")

	flag.Parse()

	models := []Model{
		// {
		// 	Slug: "gemini-2.5-flash-lite",
		// 	RPM:  15,
		// },
		// {
		// 	Slug: "gemini-2.0-flash-lite",
		// 	RPM:  30,
		// },
		// {
		// 	Slug: "gemini-2.0-flash",
		// 	RPM:  15,
		// },
		// {
		// 	Slug: "gemini-2.5-flash",
		// 	RPM:  10,
		// },
		{
			Slug: "gemma-3-12b-it",
			RPM:  30,
		},
		{
			Slug: "gemma-3-4b-it",
			RPM:  30,
		},
		{
			Slug: "gemma-3-1b-it",
			RPM:  30,
		},
		{
			Slug: "gemma-3n-e4b-it",
			RPM:  30,
		},
		{
			Slug: "gemma-3n-e2b-it",
			RPM:  30,
		},
		{
			Slug: "gemma-3-27b-it",
			RPM:  30,
		},
	}

	for modelIdx, model := range models {
		models[modelIdx].Timeout = time.Duration(int(time.Second) * 64 / model.RPM)
	}

	modelIdx := 0
	model := models[modelIdx]

	log.Println("Creating API client...")

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  *apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("error creating API client: %v\n", err)
	}

	log.Println("Reading data and output files...")

	dataFile, err := os.ReadFile(*dataPath)
	if err != nil {
		log.Fatalf("error reading data file: %v\n", err)
	}

	var data api.API
	err = json.Unmarshal(dataFile, &data)
	if err != nil {
		log.Fatalf("error unmarshaling data JSON: %v\n", err)
	}

	outputFile, err := os.ReadFile(*outputPath)
	if err != nil {
		log.Fatalf("error reading output file: %v\n", err)
	}

	var output Output
	err = json.Unmarshal(outputFile, &output)
	if err != nil {
		log.Fatalf("error unmarshaling output JSON: %v\n", err)
	}

	verseKeys := api.GetVerseKeys(data)

	processedSegments := slices.Collect(maps.Keys(output.Segments))
	processedErrors := slices.Collect(maps.Keys(output.Errors))

	log.Printf("Using model %s with timeout %s...\n", model.Slug, model.Timeout)

	curModelRequests := 0

	for _, verseKey := range verseKeys {
		if slices.Contains(processedSegments, verseKey) {
			// log.Printf("Verse %s already processed successfully, skipping...\n", verseKey)
			continue
		}
		if slices.Contains(processedErrors, verseKey) {
			// log.Printf("Verse %s already processed unsuccessfully, skipping...\n", verseKey)
			continue
		}

		log.Printf("Processing Verse %s...\n", verseKey)

		translation := data.Verses[verseKey].Translations[0].Text
		segments, err, isAPIError := splitTranslation(translation, client, model.Slug)

		if err != nil {
			log.Printf("Error while splitting translation: %v\n", err)

			if isAPIError {
				log.Println("Assuming rate limit has been hit.")
				modelIdx += 1

				if modelIdx == len(models) {
					log.Println("All models exhausted. Writing computed results to output file.")
					break
					// modelIdx = -1
				}

				model = models[modelIdx]
				log.Printf("Switching model to %s with timeout %s\n", model.Slug, model.Timeout)
				continue
			}

			continue
		}

		joined := strings.Join(segments, " ")

		if joined != translation {
			log.Println("Segments do not return original string on joining...")
			output.Errors[verseKey] = segments
		}

		if joined == translation {
			output.Segments[verseKey] = segments
		}

		// time.Sleep(model.Timeout)
		curModelRequests += 1
		if curModelRequests+1 == model.RPM {
			curModelRequests = 0
			modelIdx += 1

			if modelIdx == len(models) {
				model = models[modelIdx]
				log.Printf("Switching model to %s with timeout %s\n", model.Slug, model.Timeout)
			}
		}

		outputJSON, err := json.Marshal(output)
		if err != nil {
			log.Fatalf("error while marshaling output to JSON: %v", err)
		}

		err = os.WriteFile(*outputPath, outputJSON, 0666)
		if err != nil {
			log.Fatalf("error while writing output JSON to file: %v", err)
		}
	}
}

func splitTranslation(translation string, client *genai.Client, modelSlug string) ([]string, error, bool) {
	thinkingBudgetVal := int32(0)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudgetVal,
		},
	}

	prompt := fmt.Sprintf(`Split the following text into individual sentences.
When all strings are joined, they must recreate the original text EXACTLY.
Sentences can end with periods, question marks, exclamation marks, or other punctuation.

Text to split:
%s`, translation)

	if strings.HasPrefix(modelSlug, "gemma") {
		config = nil
		prompt = fmt.Sprintf(`Split the following text into individual sentences.
When all strings are joined, they must recreate the original text EXACTLY.
Sentences can end with periods, question marks, exclamation marks, or other punctuation.
Output as a JSON list. Do not wrap in a code block.

Example Input:
Hi Tom, what are you doing? Let's go have some pizza. I like the shop down the block.

Example Output:
[
"Hi Tom, what are you doing?",
"Let's go have some pizza.",
"I like the shop down the block."
]

Text to split:
%s`, translation)
	}

	parts := []*genai.Part{
		{Text: prompt},
	}

	result, err := client.Models.GenerateContent(context.Background(), modelSlug, []*genai.Content{{Parts: parts}}, config)
	if err != nil {
		return []string{}, err, true
	}

	var sentences []string
	err = json.Unmarshal([]byte(result.Text()), &sentences)
	if err != nil {
		return []string{}, err, false
	}

	return sentences, nil, false
}

func check() {
	dataPath := flag.String("data", "data.json", "Path to JSON file to read from")
	outputPath := flag.String("output", "sentences.json", "Path to JSON file to write to")

	flag.Parse()

	log.Println("Reading data and output files...")

	dataFile, err := os.ReadFile(*dataPath)
	if err != nil {
		log.Fatalf("error reading data file: %v\n", err)
	}

	var data api.API
	err = json.Unmarshal(dataFile, &data)
	if err != nil {
		log.Fatalf("error unmarshaling data JSON: %v\n", err)
	}

	outputFile, err := os.ReadFile(*outputPath)
	if err != nil {
		log.Fatalf("error reading output file: %v\n", err)
	}

	var output Output
	err = json.Unmarshal(outputFile, &output)
	if err != nil {
		log.Fatalf("error unmarshaling output JSON: %v\n", err)
	}

	output.Errors = make(map[string][]string)
	verseKeys := api.GetVerseKeys(data)
	for _, verseKey := range verseKeys {
		if strings.Join(output.Segments[verseKey], " ") != data.Verses[verseKey].Translations[0].Text {
			output.Errors[verseKey] = output.Segments[verseKey]
		}
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("error while marshaling output to JSON: %v", err)
	}

	err = os.WriteFile(*outputPath, outputJSON, 0666)
	if err != nil {
		log.Fatalf("error while writing output JSON to file: %v", err)
	}
}
