package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple-rag/models"
)

// When you POST to /query with a question, it:
// 1. Converts question to vector
// 2. Searches Pinecone for similar documents
// 3. Generates answer using found documents
// 4. Returns answer with sources
type QueryHandler struct {
	ragService models.RAGService
}

func NewQueryHandler(ragService models.RAGService) *QueryHandler {
	return &QueryHandler{ragService: ragService}
}

func (h *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.ragService.Query(request)
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Printf("✅ Processed query: %s\n", request.Question)
}
