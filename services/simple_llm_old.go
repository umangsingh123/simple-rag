package services

import (
	"fmt"
	"simple-rag/models"
)

// Methods:
// - GenerateResponse(): Creates answer from documents
// - Simple template-based (no API costs)
type SimpleLLM_old struct{}

func NewSimpleLLM_old() *SimpleLLM {
	return &SimpleLLM{}
}

func (s *SimpleLLM) GenerateResponse_old(question string, documents []models.Document) string {
	if len(documents) == 0 {
		return "I couldn't find any relevant information to answer your question."
	}

	answer := fmt.Sprintf("Question: %s\n\n", question)
	answer += fmt.Sprintf("I found %d relevant documents:\n\n", len(documents))

	for i, doc := range documents {
		answer += fmt.Sprintf("📄 Document %d: %s\n\n", i+1, doc.Content)
	}

	answer += "These documents contain information relevant to your question."
	return answer
}
