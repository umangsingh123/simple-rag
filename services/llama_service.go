package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Methods:
// - CreateEmbedding(): Text → Vector conversion
// - Calls your Llama server at localhost:8081
type LlamaEmbedder struct {
	BaseURL string
	Client  *http.Client
}

func NewLlamaEmbedder() *LlamaEmbedder {
	return &LlamaEmbedder{
		BaseURL: "http://localhost:8081",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (l *LlamaEmbedder) CreateEmbedding(text string) ([]float32, error) {
	reqBody := map[string]interface{}{
		"model": "llama-text-embed-v2",
		"input": []string{text},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := l.Client.Post(l.BaseURL+"/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Llama error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding received")
	}

	return result.Data[0].Embedding, nil
}
