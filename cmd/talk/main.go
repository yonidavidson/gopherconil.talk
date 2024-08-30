package main

import (
	"fmt"
	"github/yonidavidson/gophercon2024/provider"
	"os"
)

func main() {
	// Retrieve the API key from the environment variable
	apiKey := os.Getenv("PRIVATE_OPENAI_KEY")
	if apiKey == "" {
		fmt.Println("Error: PRIVATE_OPENAI_KEY environment variable not set")
		return
	}
	p := provider.OpenAIProvider{APIKey: apiKey}
	p.Response()
}
