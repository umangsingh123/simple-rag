package main

import (
	"fmt"
	"log"
	"simple-rag/router"
	"simple-rag/server"
	"simple-rag/services"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// Steps:
// 1. Initialize all services
// 2. Setup router
// 3. Start server
// 4. Wait for shutdown
func main() {
	fmt.Println(":::::::::: Starting Simple RAG Server...::::::::::")
	fmt.Println("=================================")

	// 1. Initialize Pinecone
	fmt.Println("1. :::::::::: Setting up Pinecone...::::::::::")
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: "",
	})
	if err != nil {
		log.Fatalf(":::::::::: Failed to create Pinecone client: %v", err)
	}

	// 2. Initialize Services
	fmt.Println("2. ::::::::::  Initializing services...")

	pineconeService := services.NewVectorStore(pc, "rag-demo-ach4dab.svc.aped-4627-b74a.pinecone.io")
	fmt.Println("-->> Connected to Pinecone index <<--")
	//embedder := services.NewLlamaEmbedder()
	embedder := services.NewOpenAIEmbedder()
	llm := services.NewSimpleLLM()
	ragService := services.NewRAGService(embedder, pineconeService, llm)

	// 3. Setup Router
	fmt.Println("3. ::::::::::  Setting up routes...::::::::::")
	appRouter := router.NewRouter(ragService)

	// 4. Start Server
	fmt.Println("4. ::::::::::: Starting server...")
	appServer := server.NewServer(":8080", appRouter.GetHandler())

	if err := appServer.Start(); err != nil {
		log.Fatalf(":::::::::: Failed to start server: %v", err)
	}

	// 5. Display startup info
	fmt.Println("=================================")
	fmt.Println(":::::::::: SERVER IS READY!::::::::::")
	fmt.Println(" URL: http://localhost:8080")
	fmt.Println(" Available endpoints:")
	fmt.Println("   GET  /health  - Health check")
	fmt.Println("   POST /ingest  - Add documents")
	fmt.Println("   POST /query   - Ask questions")
	fmt.Println("=================================")
	fmt.Println("  Press Ctrl+C to shutdown gracefully")
	fmt.Println("=================================")

	// 6. Wait for shutdown signal
	appServer.WaitForShutdown()
}
