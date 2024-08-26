package main

import (
	"fmt"
	"html/template"
	"strings"
)

const promptTemplate = `System: You are a helpful assistant that provides concise and accurate information. Please limit your response to approximately {{.MaxTokens}} tokens.

Context: {{.RAGContext}}

User: {{.UserQuery}}`

type PromptData struct {
	MaxTokens  string
	RAGContext string
	UserQuery  string
}

func generatePrompt(maxTokens int, ragContext, userQuery string) (string, error) {
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	data := PromptData{
		MaxTokens:  fmt.Sprintf("%d", maxTokens),
		RAGContext: limitTokens(ragContext, maxTokens/2),
		UserQuery:  limitTokens(userQuery, maxTokens/4),
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return result.String(), nil
}

func limitTokens(s string, maxTokens int) string {
	words := strings.Fields(s)
	if len(words) <= maxTokens {
		return s
	}
	return strings.Join(words[:maxTokens], " ") + "..."
}

func main() {
	maxTokens := 100
	ragContext := "This is some context information for the LLM."
	userQuery := "What is the capital of France?"

	prompt, err := generatePrompt(maxTokens, ragContext, userQuery)
	if err != nil {
		fmt.Printf("Error generating prompt: %v\n", err)
		return
	}

	fmt.Println(prompt)
}
