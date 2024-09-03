package agent

import (
	"fmt"
	"github.com/yonidavidson/gophercon-israel-2024/prompt"
	"github.com/yonidavidson/gophercon-israel-2024/provider"
	"github.com/yonidavidson/gophercon-israel-2024/rag"
)

// Agent represents a conversational agent that uses a language model and retrieval-augmented generation (RAG) to answer questions.
type Agent struct {
	p provider.OpenAIProvider
	r *rag.Rag
	e []rag.Embedding
}

type promptData struct {
	RAGContext   string
	UserQuery    string
	SystemPrompt string
}

// New creates a new instance of Agent with the provided OpenAI provider, RAG instance, and embeddings
func New(p provider.OpenAIProvider, r *rag.Rag, e []rag.Embedding) *Agent {
	return &Agent{
		p: p,
		r: r,
		e: e,
	}
}

// HandleUserQuery takes a user query, retrieves relevant context using RAG, generates a prompt,
func (a Agent) HandleUserQuery(promptTemplate, systemPrompt, userQuery string) ([]byte, error) {
	var ragContext string
	if a.r != nil && a.e != nil {
		rc, err := a.r.Search(userQuery, a.e)
		if err != nil {
			return nil, fmt.Errorf("error searching text: %v", err)
		}
		ragContext = string(rc)
	}
	m, err := prompt.ParseMessages(promptTemplate, promptData{
		RAGContext:   ragContext,
		UserQuery:    userQuery,
		SystemPrompt: systemPrompt,
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing messages: %v", err)
	}
	c, err := a.p.ChatCompletion(m)
	if err != nil {
		return nil, fmt.Errorf("error getting chat completion: %v", err)
	}
	return c, nil
}
