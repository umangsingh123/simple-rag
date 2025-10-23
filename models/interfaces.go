package models

// RAGService interface defines the contract for RAG operations
type RAGService interface {
	Query(request QueryRequest) (*QueryResponse, error)
	Ingest(request IngestionRequest) error
}
