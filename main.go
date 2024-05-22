package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// dynamic arguments through path or CLI later

func importData(path string) (string, map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	bb, err := io.ReadAll(f)
	if err != nil {
		return "", nil, err
	}

	doc := make(map[string]interface{})
	if err := json.Unmarshal(bb, &doc); err != nil {
		return "", nil, err
	}

	if flightDate, contains := doc["flight_date"]; contains {
		log.Printf("Happy end we have a flight_date: %q\n", flightDate)
	} else {
		log.Println("Json document doesn't have a flightDate field.")
	}

	// test
	beautifiedJSON, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", nil, err
	}

	fmt.Println(string(beautifiedJSON))

	return string(bb), doc, nil
}

// add to CLI later
// importData(os.Args[0])

func main() {
	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Fatal(err)
	}
	flightsData, docs, err := importData("data/flights.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	embs, err := llm.CreateEmbedding(ctx, []string{flightsData})
	fmt.Printf("Got %d embeddings:\n", len(embs))
	for i, emb := range embs {
		fmt.Printf("%d: len=%d; first few=%v\n", i, len(emb), emb[:4])
	}
	prompt := fmt.Sprintf("give me the total number of canceled, scheduled, active, landed, incident, diverted, flighst"+
		"with an extra explanation with the respective flight_date and airline for each canceled flight with the format "+
		"'canceled': 10"+
		" %s", flightsData)

	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
		llms.WithMetadata(docs),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = completion

	// change later for multiple imports

}
