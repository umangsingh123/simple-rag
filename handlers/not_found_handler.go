package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"error":     "Endpoint not found",
		"path":      r.URL.Path,
		"method":    r.Method,
		"timestamp": time.Now().Format(time.RFC3339),
		"available_endpoints": []string{
			"GET /health",
			"POST /ingest",
			"POST /query",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}
