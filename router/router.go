package router

import (
	"net/http"
	"simple-rag/handlers"
	"simple-rag/models"
)

// Routes configured:
// /health  → HealthHandler
// /ingest  → IngestHandler
// /query   → QueryHandler
// /        → NotFoundHandler (catch-all)

type Router struct {
	mux *http.ServeMux
}

// Dependency Injection: Takes models.RAGService interface
func NewRouter(ragService models.RAGService) *Router {
	mux := http.NewServeMux()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	ingestHandler := handlers.NewIngestHandler(ragService)
	queryHandler := handlers.NewQueryHandler(ragService)
	notFoundHandler := handlers.NewNotFoundHandler()

	// Register routes
	mux.Handle("/health", healthHandler)
	mux.Handle("/ingest", ingestHandler)
	mux.Handle("/query", queryHandler)
	mux.Handle("/", notFoundHandler) // Catch-all

	return &Router{mux: mux}
}

// Accessor method: Provides the underlying http.Handler interface
// Links to server package: This method is called in server.go
func (r *Router) GetHandler() http.Handler {
	return r.mux
}
