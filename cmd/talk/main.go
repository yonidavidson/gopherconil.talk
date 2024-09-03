package main

import (
	"fmt"
	"github.com/yonidavidson/gophercon-israel-2024/prompt"
	"github.com/yonidavidson/gophercon-israel-2024/provider"
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
	messages := []prompt.Message{
		{
			Role:    prompt.RoleSystem,
			Content: "You are a helpful assistant that provides concise and accurate information.",
		},
		{
			Role:    prompt.RoleUser,
			Content: "Translate the following English text to French: 'Hello, how are you",
		},
	}
	r, err := p.ChatCompletion(messages)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(r))
}
