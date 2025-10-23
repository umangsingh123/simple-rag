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

func (r *Router) GetHandler() http.Handler {
	return r.mux
}
