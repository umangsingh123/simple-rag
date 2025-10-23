package handlers

// When you POST to /ingest with documents, it:
// 1. Converts documents to vectors using Llama
// 2. Stores them in Pinecone
// 3. Returns success message
import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple-rag/models"
)

type IngestHandler struct {
	ragService models.RAGService
}

func NewIngestHandler(ragService models.RAGService) *IngestHandler {
	return &IngestHandler{ragService: ragService}
}

func (h *IngestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.IngestionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.ragService.Ingest(request); err != nil {
		http.Error(w, "Ingestion failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":        "Documents added successfully",
		"document_count": len(request.Documents),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Printf("✅ Ingested %d documents\n", len(request.Documents))
}
