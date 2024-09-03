package main

import (
	"fmt"
	"github.com/yonidavidson/gophercon-israel-2024/prompt"
	"github.com/yonidavidson/gophercon-israel-2024/provider"
)

func main() {
	p, err := provider.NewOpenAIProvider()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
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
