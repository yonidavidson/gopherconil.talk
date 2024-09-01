package rag

import (
	"github.com/yonidavidson/gophercon-israel-2024/provider"
	"math"
	"sort"
)

type (
	// Rag is a Retrieval Augmented Generation (RAG) struct
	Rag struct {
		provider provider.OpenAIProvider
	}

	// Embedding represents a text embedding
	Embedding struct {
		text   string
		vector []float64
	}

	// scoredEmbedding represents an Embedding with a score
	scoredEmbedding struct {
		Embedding
		score float64
	}
)

// New creates a new Rag struct
func New(provider provider.OpenAIProvider) *Rag {
	return &Rag{
		provider: provider,
	}
}

// Embed receives a large text and returns a slice embeddings
func (r *Rag) Embed(text string, chunkSize int) ([]Embedding, error) {
	// Split the text into chunks of the specified size
	var chunks []string
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	// Get embeddings for each chunk
	vectors, err := r.provider.TextEmbedding(chunks)
	if err != nil {
		return nil, err
	}

	// Create embeddings slice
	result := make([]Embedding, len(chunks))
	for i, chunk := range chunks {
		result[i] = Embedding{
			text:   chunk,
			vector: vectors[i],
		}
	}

	return result, nil
}

// Search receives a query and a slice of embeddings and returns the most relevant embeddings
func (r *Rag) Search(query string, embeddings []Embedding) ([]byte, error) {
	// Get the Embedding for the query
	queryEmbedding, err := r.provider.TextEmbedding([]string{query})
	if err != nil {
		return nil, err
	}

	// Calculate the similarity between the query Embedding and each text Embedding

	var scoredEmbeddings []scoredEmbedding
	for _, emb := range embeddings {
		score := cosineSimilarity(queryEmbedding[0], emb.vector)
		scoredEmbeddings = append(scoredEmbeddings, scoredEmbedding{emb, score})
	}

	// Sort the embeddings by similarity score in descending order
	sort.Slice(scoredEmbeddings, func(i, j int) bool {
		return scoredEmbeddings[i].score > scoredEmbeddings[j].score
	})

	// Convert scored embeddings back to the original embeddings type
	result := make([]Embedding, len(scoredEmbeddings))
	for i, se := range scoredEmbeddings {
		result[i] = se.Embedding
	}

	return []byte(result[0].text), nil
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
