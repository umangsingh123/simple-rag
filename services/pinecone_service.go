package services

import (
	"context"
	"fmt"
	"simple-rag/models"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"google.golang.org/protobuf/types/known/structpb" //Protocol Buffers for metadata handling (required by Pinecone SDK)
)

type VectorStore struct {
	Client    *pinecone.Client
	IndexHost string
}

func NewVectorStore(client *pinecone.Client, indexHost string) *VectorStore {
	return &VectorStore{
		Client:    client,
		IndexHost: indexHost,
	}
}

func (v *VectorStore) Upsert(documents []models.Document) error {
	ctx := context.Background()

	index, err := v.Client.Index(pinecone.NewIndexConnParams{
		Host: v.IndexHost,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to index: %v", err)
	}

	// Create vectors using the WORKING pinecone.Vector struct
	vectors := make([]pinecone.Vector, 0, len(documents))
	for i, doc := range documents {
		fmt.Printf(">>>>> Document %d: ID=%s, Embedding=%d dimensions\n", i+1, doc.ID, len(doc.Embedding))

		// Create metadata using structpb
		metadata, err := structpb.NewStruct(map[string]interface{}{
			"content": doc.Content,
		})
		if err != nil {
			return fmt.Errorf("failed to create metadata: %v", err)
		}

		// Create a copy of embedding for pointer
		//Pinecone SDK requires *[]float32 for Values field
		embeddingCopy := make([]float32, len(doc.Embedding))
		copy(embeddingCopy, doc.Embedding)

		vectors = append(vectors, pinecone.Vector{
			Id:       doc.ID,
			Values:   &embeddingCopy, // Pointer to []float32
			Metadata: metadata,       // structpb.Struct
		})
	}

	// Convert to pointers for UpsertVectors
	vectorPointers := make([]*pinecone.Vector, len(vectors))
	for i := range vectors {
		vectorPointers[i] = &vectors[i]
	}

	// Use UpsertVectors - the working method!
	_, err = index.UpsertVectors(ctx, vectorPointers)
	if err != nil {
		return fmt.Errorf("failed to upsert vectors: %v", err)
	}

	fmt.Printf("::: Successfully upserted %d vectors\n", len(vectors))
	return nil
}

func (v *VectorStore) Search(embedding []float32, topK int) ([]models.Document, error) {
	if topK == 0 {
		topK = 5
	}

	ctx := context.Background()
	index, err := v.Client.Index(pinecone.NewIndexConnParams{
		Host: v.IndexHost,
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf(">>>>>>>> Searching with %d-dimensional vector <<<<<<<, topK=%d\n", len(embedding), topK)

	// Use SearchRecords with the correct structure
	res, err := index.SearchRecords(ctx, &pinecone.SearchRecordsRequest{
		Query: pinecone.SearchRecordsQuery{
			TopK: int32(topK),
			Vector: &pinecone.SearchRecordsVector{
				Values: &embedding, // This should work since UpsertVectors stores the vectors
			},
		},
		Fields: &[]string{"content"}, // Request the content field back
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf(">>>>> Search found %d matches\n", len(res.Result.Hits))

	documents := make([]models.Document, len(res.Result.Hits))
	for i, hit := range res.Result.Hits {
		content := ""
		if hit.Fields != nil {
			if contentVal, ok := hit.Fields["content"]; ok {
				if contentStr, ok := contentVal.(string); ok {
					content = contentStr
				}
			}
		}

		fmt.Printf("   Match %d: %s (score: %.3f)\n", i+1, hit.Id, hit.Score)
		documents[i] = models.Document{
			ID:      hit.Id,
			Content: content,
		}
	}

	return documents, nil
}
