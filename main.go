package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

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

	//println(string(bb))

	return string(bb), doc, nil
}

func postResponse() {
	data, _, err := importData("data/flights.json")
	if err != nil {
		log.Fatal(err)
	}

	prompt := fmt.Sprintf("give me the total number of canceled, scheduled, active, landed, incident, diverted, flighst"+
		"with an extra explanation with the respective flight_date and airline for each canceled flight with the format "+
		"'canceled': 10"+
		" %s", data)

	//url := "http://localhost:11434/api/chat"
	//payload := map[string]interface{}{
	//	"model": "llama3",
	//	"messages": []map[string]string{
	//		{
	//			"role":    "user",
	//			"content": prompt,
	//		},
	//	},
	//	"format": "json",
	//	"stream": true,
	//}
	url := "http://localhost:11434/api/generate"
	payload := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}
	startTime := time.Now()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	println(string(body))
	duration := time.Since(startTime)
	valueFromDuration := slog.DurationValue(duration)
	slog.Info("Ollama job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
}
func main() {
	postResponse()
}
