package models

// Document represents a piece of text with its vector embedding
type Document struct {
	ID        string    `json:"id"`        // Unique ID for the document
	Content   string    `json:"content"`   // The actual text content
	Embedding []float32 `json:"embedding"` // Vector representation (1536 numbers from OpenAI)
}

// QueryRequest is what users send when asking questions
type QueryRequest struct {
	Question string `json:"question"` // User's question
	TopK     int    `json:"top_k"`    // How many results to return
}

// QueryResponse is what we send back to users
type QueryResponse struct {
	Answer  string     `json:"answer"`  // Generated answer
	Sources []Document `json:"sources"` // Documents used for answer
}

// IngestionRequest is for adding documents to the system
type IngestionRequest struct {
	Documents []Document `json:"documents"` // List of documents to add
}
