package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yonidavidson/gopherconil.talk/prompt"
	"io"
	"net/http"
	"os"
)

type (
	// OpenAIProvider is a provider that uses the OpenAI API
	OpenAIProvider struct {
		APIKey string
	}

	// requestPayload is the JSON payload we send to the OpenAI API
	requestPayload struct {
		Model       string    `json:"model"`
		Messages    []message `json:"messages"`
		MaxTokens   int       `json:"max_tokens"`
		Temperature float64   `json:"temperature"`
		TopP        float64   `json:"top_p"`
		N           int       `json:"n"`
		Stop        *string   `json:"stop"`
	}

	// choice struct represents a single choice from the OpenAI API response
	choice struct {
		Message message `json:"message"`
	}

	// responsePayload is the JSON payload we receive from the OpenAI API
	responsePayload struct {
		Choices []choice `json:"choices"`
	}

	// embeddingRequestPayload is the JSON payload we send to the OpenAI API for embedding
	embeddingRequestPayload struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}

	embedding struct {
		Embedding []float64 `json:"embedding"`
	}

	embeddingResponsePayload struct {
		Data []embedding `json:"data"`
	}

	message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

const (
	endpoint          = "https://api.openai.com/v1/chat/completions"
	embeddingEndpoint = "https://api.openai.com/v1/embeddings"
)

// NewOpenAIProvider creates a new instance of OpenAIProvider with the API key from the environment variable.
func NewOpenAIProvider() (*OpenAIProvider, error) {
	apiKey := os.Getenv("PRIVATE_OPENAI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PRIVATE_OPENAI_KEY environment variable is not set")
	}
	return &OpenAIProvider{APIKey: apiKey}, nil
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
func (p OpenAIProvider) TextEmbedding(input []string) ([][]float64, error) {
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

	// Convert the embedding data to the expected return type
	embeddings := make([][]float64, len(responsePayload.Data))
	for i, emb := range responsePayload.Data {
		embeddings[i] = emb.Embedding
	}

	return embeddings, nil
}
