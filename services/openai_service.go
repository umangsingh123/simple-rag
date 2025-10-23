package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIEmbedder struct {
	BaseURL string
	Client  *http.Client
	APIKey  string
}

func NewOpenAIEmbedder() *OpenAIEmbedder {
	return &OpenAIEmbedder{
		BaseURL: "https://api.openai.com/v1",
		Client:  &http.Client{Timeout: 30 * time.Second},
		APIKey:  "",
	}
}

func (e *OpenAIEmbedder) CreateEmbedding(text string) ([]float32, error) {
	reqBody := map[string]interface{}{
		"model": "text-embedding-3-small", // Lowest cost so used that
		"input": text,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", e.BaseURL+"/embeddings", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+e.APIKey)

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// ✅ Debug line to see actual response if something fails
	fmt.Println("Embedding API response:", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding request failed: %s", string(body))
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
		return nil, fmt.Errorf("no embedding received: %s", string(body))
	}

	return result.Data[0].Embedding, nil
}
