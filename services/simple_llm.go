package services

import (
	"fmt"
	"simple-rag/models"
	"strings"
)

type SimpleLLM struct{}

func NewSimpleLLM() *SimpleLLM {
	return &SimpleLLM{}
}

func (s *SimpleLLM) GenerateResponse(question string, documents []models.Document) string {
	if len(documents) == 0 {
		return "I couldn't find any relevant information to answer your question based on the documents I have access to."
	}

	// Build a proper answer using the actual document content
	answer := s.buildAnswerFromDocuments(question, documents)
	return answer
}

func (s *SimpleLLM) buildAnswerFromDocuments(question string, documents []models.Document) string {
	var answer strings.Builder

	answer.WriteString("Based on the documents I found, here's the answer to your question:\n\n")
	answer.WriteString(fmt.Sprintf("**Question:** %s\n\n", question))
	answer.WriteString("**Answer:** ")

	// Extract key information from the most relevant documents
	for i, doc := range documents {
		if i == 0 {
			// Use the most relevant document as the main answer
			answer.WriteString(s.extractKeyInformation(doc.Content))
		} else if i < 3 {
			// Add supporting information from other relevant documents
			answer.WriteString(" ")
			answer.WriteString(s.extractSupportingInfo(doc.Content))
		}
	}

	answer.WriteString("\n\n**Sources:**\n")
	for i, doc := range documents {
		answer.WriteString(fmt.Sprintf("%d. %s\n", i+1, doc.ID))
	}

	return answer.String()
}

func (s *SimpleLLM) extractKeyInformation(content string) string {
	// Extract the most relevant sentence or key information
	sentences := strings.Split(content, ".")
	if len(sentences) > 0 {
		return strings.TrimSpace(sentences[0]) + "."
	}
	return content
}

func (s *SimpleLLM) extractSupportingInfo(content string) string {
	// Extract additional supporting information
	sentences := strings.Split(content, ".")
	if len(sentences) > 1 {
		return strings.TrimSpace(sentences[1]) + "."
	}
	return ""
}
