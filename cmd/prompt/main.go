package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	endpoint = "https://api.openai.com/v1/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	N           int       `json:"n"`
	Stop        *string   `json:"stop"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponsePayload struct {
	Choices []Choice `json:"choices"`
}

func main() {
	// Retrieve the API key from the environment variable
	apiKey := os.Getenv("PRIVATE_OPENAI_KEY")
	if apiKey == "" {
		fmt.Println("Error: PRIVATE_OPENAI_KEY environment variable not set")
		return
	}

	// Define the payload with a system prompt
	payload := RequestPayload{
		Model: "gpt-4o-mini-2024-07-18",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant that provides concise and accurate information."},
			{Role: "user", Content: "Translate the following English text to French: 'Hello, how are you?'"},
		},
		MaxTokens:   60,
		Temperature: 0.5,
		TopP:        1.0,
		N:           1,
		Stop:        nil,
	}

	// Marshal the payload into JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Check if the request was successful
	if resp.StatusCode == http.StatusOK {
		var responsePayload ResponsePayload
		if err := json.Unmarshal(body, &responsePayload); err != nil {
			fmt.Printf("Error unmarshaling response JSON: %v\n", err)
			return
		}

		// Print the generated text
		fmt.Println("Response:", responsePayload.Choices[0].Message.Content)
	} else {
		fmt.Printf("Error: %v - %s\n", resp.StatusCode, body)
	}
}
