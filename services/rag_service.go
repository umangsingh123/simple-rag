package services

import (
	"fmt"
	"simple-rag/models"
)

// RAG Pipeline:
// 1. Question → Vector (Llama)
// 2. Vector → Similar Documents (Pinecone)
// 3. Documents → Answer (SimpleLLM)
type RAGService struct {
	//Embedder *LlamaEmbedder
	Embedder *OpenAIEmbedder
	Store    *VectorStore
	LLM      *SimpleLLM
}

func NewRAGService(embedder *OpenAIEmbedder, store *VectorStore, llm *SimpleLLM) *RAGService {
	return &RAGService{
		Embedder: embedder,
		Store:    store,
		LLM:      llm,
	}
}

func (r *RAGService) Query(request models.QueryRequest) (*models.QueryResponse, error) {
	fmt.Printf(">>> Processing question: %s\n", request.Question)

	embedding, err := r.Embedder.CreateEmbedding(request.Question)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %v", err)
	}

	topK := request.TopK
	if topK == 0 {
		topK = 3
	}

	documents, err := r.Store.Search(embedding, topK)
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	fmt.Printf(">>>>> Found %d relevant documents\n", len(documents))

	answer := r.LLM.GenerateResponse(request.Question, documents)

	return &models.QueryResponse{
		Answer:  answer,
		Sources: documents,
	}, nil
}

func (r *RAGService) Ingest(request models.IngestionRequest) error {
	fmt.Printf(">>>>>>> Ingesting %d documents...\n", len(request.Documents))

	for i := range request.Documents {

		embedding, err := r.Embedder.CreateEmbedding(request.Documents[i].Content)

		if err != nil {
			return fmt.Errorf("failed to embed document %s: %v", request.Documents[i].ID, err)
		}
		request.Documents[i].Embedding = embedding
	}

	if err := r.Store.Upsert(request.Documents); err != nil {
		return fmt.Errorf("failed to store documents: %v", err)
	}

	fmt.Println(">>>>> Documents ingested successfully!")
	return nil
}
