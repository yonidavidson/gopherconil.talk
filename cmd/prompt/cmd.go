package main

import (
	"fmt"
	"html/template"
	"strings"
)

const promptTemplate = `System: {{.SystemPrompt}}

Chat History:
{{limitTokens .ChatHistory (multiply .MaxTokens 0.3)}}

Context: {{limitTokens .RAGContext (multiply .MaxTokens 0.1)}}

User: {{limitTokens .UserQuery (multiply .MaxTokens 0.2)}}`

type PromptData struct {
	MaxTokens    float64
	RAGContext   string
	UserQuery    string
	ChatHistory  string
	SystemPrompt string
}

func generatePrompt(maxTokens int, ragContext, userQuery, chatHistory, systemPrompt string) (string, error) {
	tmpl, err := template.New("talk").Funcs(template.FuncMap{
		"limitTokens": limitTokens,
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	}).Parse(promptTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	data := PromptData{
		MaxTokens:    float64(maxTokens),
		RAGContext:   ragContext,
		UserQuery:    userQuery,
		ChatHistory:  chatHistory,
		SystemPrompt: systemPrompt,
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return result.String(), nil
}

func limitTokens(s string, maxTokens float64) string {
	const avgTokenLength = 4 // Average token length heuristic
	maxChars := int(maxTokens * avgTokenLength)

	if len(s) <= maxChars {
		return s
	}
	return s[:maxChars] + "..."
}

func main() {
	maxTokens := 100
	ragContext := "Paris, the capital of France, is a major European city and a global center for art, fashion, gastronomy, and culture. Its 19th-century cityscape is crisscrossed by wide boulevards and the River Seine. Beyond such landmarks as the Eiffel Tower and the 12th-century, Gothic Notre-Dame cathedral, the city is known for its cafe culture and designer boutiques along the Rue du Faubourg Saint-Honoré."
	userQuery := "Can you tell me about the history and main attractions of Paris? Also, what`s the best time to visit and are there any local customs I should be aware of?"
	chatHistory := "User: I`m planning a trip to Europe.\nAssistant: That`s exciting! Europe has many wonderful destinations. Do you have any specific countries or cities in mind?\nUser: I'm thinking about visiting France.\nAssistant: France is a great choice! It offers a rich history, beautiful landscapes, and world-renowned cuisine. Are you interested in visiting Paris or exploring other regions as well?"
	systemPrompt := "You are a knowledgeable and helpful travel assistant. Provide accurate and concise information about destinations, attractions, local customs, and travel tips. When appropriate, suggest off-the-beaten-path experiences that tourists might not typically know about. Always prioritize the safety and cultural sensitivity of the traveler."

	prompt, err := generatePrompt(maxTokens, ragContext, userQuery, chatHistory, systemPrompt)
	if err != nil {
		fmt.Printf("Error generating talk: %v\n", err)
		return
	}

	fmt.Println(prompt)
}
