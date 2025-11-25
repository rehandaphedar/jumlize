package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"google.golang.org/genai"
)

var model *string
var config *genai.GenerateContentConfig
var client *genai.Client

func runSegmentCommand() {
	segmentFlagSet := flag.NewFlagSet("segment", flag.ExitOnError)

	apiKey := segmentFlagSet.String("api_key", "", "Gemini API key")
	model = segmentFlagSet.String("model", "gemini-2.5-flash", "Model")
	temperature := segmentFlagSet.Float64("temperature", 0, "Temperature")
	thinkingBudget := segmentFlagSet.Int64("thinking_budget", 0, "Thinking budget")
	promptPath := segmentFlagSet.String("template", "prompt.md", "Path to prompt template")
	translationPath = segmentFlagSet.String("translation", "en-sahih-international-simple.json", "Path to JSON file to read from")
	metadataPath = segmentFlagSet.String("metadata", "quran-metadata-ayah.json", "Path to JSON file to read from")
	verseKeys := segmentFlagSet.String("verse_keys", "all", "Comma separated list of verse keys to process. Set to 'all' to process all verses.")
	overwrite := segmentFlagSet.Bool("overwrite", false, "Whether to overwrite existing segments")

	segmentFlagSet.Parse(os.Args[2:])

	printConfig(segmentFlagSet)
	loadTranslation()
	loadMetadata()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		save()
		os.Exit(0)
	}()

	log.Println("Creating API client...")

	var err error
	client, err = genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  *apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("error creating API client: %v\n", err)
	}

	log.Println("Parsing template file ...")

	prompt, err := os.ReadFile(*promptPath)
	if err != nil {
		log.Fatalf("error while parsing template file: %v\n", err)
	}

	thinkingBudgetInt32 := int32(*thinkingBudget)
	temperatureFloat32 := float32(*temperature)
	config = &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudgetInt32,
		},
		Temperature:       &temperatureFloat32,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: string(prompt)}}},
	}

	var verseKeysToProcess []string
	if *verseKeys == "all" {
		for k := range translation {
			verseKeysToProcess = append(verseKeysToProcess, k)
		}
	} else {
		verseKeysToProcess = strings.Split(*verseKeys, ",")
	}

	for _, verseKey := range verseKeysToProcess {
		verse := translation[verseKey]

		if sanityCheck(verseKey, verse.Text, verse.Segments) && !*overwrite {
			continue
		}

		log.Printf("Processing Verse %s...\n", verseKey)

		strs, err := segmentText(verse.Text)
		if err != nil {
			log.Printf("Error while segmenting text: %v\n", err)
			log.Println("Assuming rate limit has been hit.")
			break
		}

		verse.Segments = stringsToSegments(strs)

		if !sanityCheck(verseKey, verse.Text, verse.Segments) {
			log.Println("Sanity check failed...")
		}

		translation[verseKey] = verse
	}

	save()
}

func segmentText(text string) ([]string, error) {
	b, _ := json.Marshal(text)
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{
					Text: string(b),
				},
			},
		},
	}

	result, err := client.Models.GenerateContent(context.Background(), *model, contents, config)
	if err != nil {
		return []string{}, err
	}

	var strs []string
	err = json.Unmarshal([]byte(result.Text()), &strs)
	if err != nil {
		return []string{}, err
	}

	return strs, nil
}
