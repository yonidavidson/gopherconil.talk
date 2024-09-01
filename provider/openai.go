package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yonidavidson/gophercon-israel-2024/prompt"
	"io"
	"net/http"
)

const (
	endpoint          = "https://api.openai.com/v1/chat/completions"
	embeddingEndpoint = "https://api.openai.com/v1/embeddings"
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

type embeddingRequestPayload struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embedding struct {
	Embedding []float64 `json:"embedding"`
}

type embeddingResponsePayload struct {
	Data []embedding `json:"data"`
}

// ChatCompletion sends a request to the OpenAI API and returns the response as a byte slice.
func (p OpenAIProvider) ChatCompletion(m []prompt.Message) ([]byte, error) {
	// convert from []prompt.Message to []message
	messages := make([]message, len(m))
	for i, m := range m {
		messages[i] = message{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}
	// Define the payload with a system talk
	payload := requestPayload{
		Model:       "gpt-4o-mini-2024-07-18",
		Messages:    messages,
		MaxTokens:   1000,
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

// TextEmbedding sends a request to the OpenAI API to get text embeddings and returns the response as a slice of float64.
func (p OpenAIProvider) TextEmbedding(input []string) ([]float64, error) {
	// Define the payload
	payload := embeddingRequestPayload{
		Model: "text-embedding-3-small",
		Input: input,
	}

	// Marshal the payload into JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", embeddingEndpoint, bytes.NewBuffer(payloadBytes))
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

	// Unmarshal the response
	var responsePayload embeddingResponsePayload
	if err := json.Unmarshal(body, &responsePayload); err != nil {
		return nil, err
	}

	// Assuming we are interested in the first embedding in the response
	if len(responsePayload.Data) == 0 {
		return nil, fmt.Errorf("no embedding data found in the response")
	}

	return responsePayload.Data[0].Embedding, nil
}
