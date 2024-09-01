package rag

import (
	"github.com/yonidavidson/gophercon-israel-2024/provider"
	"math"
	"sort"
)

type rag struct {
	provider provider.OpenAIProvider
}

type embedding struct {
	text   string
	vector []float64
}

type scoredEmbedding struct {
	embedding
	score float64
}

// Embed receives a large text and returns a slice embeddings
func (r *rag) Embed(text string) ([]embedding, error) {
	// Split the text into chunks (assuming a simple split by sentences for this example)
	chunks := []string{text} // You might want to split the text into smaller chunks

	// Get embeddings for each chunk
	vectors, err := r.provider.TextEmbedding(chunks)
	if err != nil {
		return nil, err
	}

	// Create embeddings slice
	result := make([]embedding, len(chunks))
	for i, chunk := range chunks {
		result[i] = embedding{
			text:   chunk,
			vector: vectors,
		}
	}

	return result, nil
}

// Search receives a query and a slice of embeddings and returns the most relevant embeddings
func (r *rag) Search(query string, embeddings []embedding) ([]embedding, error) {
	// Get the embedding for the query
	queryEmbedding, err := r.provider.TextEmbedding([]string{query})
	if err != nil {
		return nil, err
	}

	// Calculate the similarity between the query embedding and each text embedding

	var scoredEmbeddings []scoredEmbedding
	for _, emb := range embeddings {
		score := cosineSimilarity(queryEmbedding, emb.vector)
		scoredEmbeddings = append(scoredEmbeddings, scoredEmbedding{emb, score})
	}

	// Sort the embeddings by similarity score in descending order
	sort.Slice(scoredEmbeddings, func(i, j int) bool {
		return scoredEmbeddings[i].score > scoredEmbeddings[j].score
	})

	// Convert scored embeddings back to the original embeddings type
	result := make([]embedding, len(scoredEmbeddings))
	for i, se := range scoredEmbeddings {
		result[i] = se.embedding
	}

	return result, nil
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
