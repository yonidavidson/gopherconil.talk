package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	endpoint = "https://api.openai.com/v1/chat/completions"
)

type OpenAIProvider struct {
	APIKey string
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestPayload struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	N           int       `json:"n"`
	Stop        *string   `json:"stop"`
}

type choice struct {
	Message message `json:"message"`
}

type responsePayload struct {
	Choices []choice `json:"choices"`
}

// ChatCompletion sends a request to the OpenAI API and returns the response as a byte slice.
func (p OpenAIProvider) ChatCompletion() ([]byte, error) {
	// Define the payload with a system talk
	payload := requestPayload{
		Model: "gpt-4o-mini-2024-07-18",
		Messages: []message{
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
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	var responsePayload responsePayload
	if err := json.Unmarshal(body, &responsePayload); err != nil {
		return nil, err
	}

	return []byte(responsePayload.Choices[0].Message.Content), nil
}
