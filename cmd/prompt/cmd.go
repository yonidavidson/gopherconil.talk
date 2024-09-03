package main

import (
	"fmt"
	"github.com/yonidavidson/gophercon-israel-2024/prompt"
	"github.com/yonidavidson/gophercon-israel-2024/provider"
)

const promptTemplate = `<system>{{.SystemPrompt}}</system>
{{limitTokens .ChatHistory (multiply .MaxTokens 1)}}
<user>
Context: {{limitTokens .RAGContext (multiply .MaxTokens 1)}}
User Query: {{limitTokens .UserQuery (multiply .MaxTokens 1)}}</user>`

type promptData struct {
	MaxTokens    float64
	RAGContext   string
	UserQuery    string
	ChatHistory  string
	SystemPrompt string
}

func main() {
	maxTokens := 1000
	ragContext := "Paris, the capital of France, is a major European city and a global center for art, fashion, gastronomy, and culture. Its 19th-century cityscape is crisscrossed by wide boulevards and the River Seine. Beyond such landmarks as the Eiffel Tower and the 12th-century, Gothic Notre-Dame cathedral, the city is known for its cafe culture and designer boutiques along the Rue du Faubourg Saint-Honoré."
	userQuery := "Can you tell me about the history and main attractions of Paris? Also, what`s the best time to visit and are there any local customs I should be aware of?"
	chatHistory := "<user>I`m planning a trip to Europe.</user>\n<assistant>That`s exciting! Europe has many wonderful destinations. Do you have any specific countries or cities in mind?</assistant>\n<user>I am thinking about visiting France.</user>\n<assistant>France is a great choice! It offers a rich history, beautiful landscapes, and world-renowned cuisine. Are you interested in visiting Paris or exploring other regions as well?</assistant>"
	systemPrompt := "You are a knowledgeable and helpful travel assistant. Provide accurate and concise information about destinations, attractions, local customs, and travel tips. When appropriate, suggest off-the-beaten-path experiences that tourists might not typically know about. Always prioritize the safety and cultural sensitivity of the traveler."
	data := promptData{
		MaxTokens:    float64(maxTokens),
		RAGContext:   ragContext,
		UserQuery:    userQuery,
		ChatHistory:  chatHistory,
		SystemPrompt: systemPrompt,
	}
	m, err := prompt.ParseMessages(promptTemplate, data)
	if err != nil {
		fmt.Printf("Error parsing messages: %v\n", err)
		return
	}
	p, err := provider.NewOpenAIProvider()
	if err != nil {
		fmt.Printf("Error creating OpenAI provider: %v\n", err)
		return
	}
	r, err := p.ChatCompletion(m)
	if err != nil {
		fmt.Printf("Error getting chat completion: %v\n", err)
		return
	}
	fmt.Println("\n\n\n\n" + string(r))
}
