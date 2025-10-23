# Simple RAG System

A production-ready Retrieval-Augmented Generation (RAG) system built in Go that allows you to store documents and ask questions based on their content.

## ::::::::: Architecture Overview ::::::::::

```mermaid



graph TB
    subgraph "Client Layer"
        C[Client Apps]
    end
  
    subgraph "API Layer"
        R[Router]
        H1[Health Handler]
        H2[Ingest Handler]
        H3[Query Handler]
        H4[404 Handler]
    end
  
    subgraph "Business Logic Layer"
        RS[RAG Service]
        LS[Llama Service]
        PS[Pinecone Service]
        LLM[Simple LLM]
    end
  
    subgraph "External Services"
        L[Llama Server<br/>localhost:8081]
        P[Pinecone<br/>Vector Database]
    end
  
    subgraph "Data Layer"
        M[Models]
    end
  
    C --> R
    R --> H1
    R --> H2
    R --> H3
    R --> H4
  
    H2 --> RS
    H3 --> RS
  
    RS --> LS
    RS --> PS
    RS --> LLM
  
    LS --> L
    PS --> P
  
    RS --> M
    LS --> M
    PS --> M

```

## ::::::::: Data Ingestion Flow :::::::::

1. Client POST /ingest with documents
2. IngestHandler validates JSON
3. RAGService processes each document:
   └── OpenAiService converts text → vector
   └── PineconeService stores vector + metadata
4. Returns success response

## **::::::::: Query Flow :::::::::**

1. Client POST /query with question
2. QueryHandler validates JSON
3. RAGService executes RAG pipeline:
   └── OpenAIService: question → vector
   └── PineconeService: vector → similar documents
   └── SimpleLLM: documents → answer
4. Returns answer with sources

## **::::::::: Set Env Variable :::::::::**

```export PINECONE_API_KEY="your_pinecone_api_key"
export PINECONE_INDEX_HOST="your_index_host"
export PINECONE_INDEX_NAME="your_index_name"
```

## **::::::::: Run App :::::::::**

```go run main.go

```

## **::::::::: Ingest Document  :::::::::::**

```curl -X POST http://localhost:8080/ingest -H "Content-Type: application/json" -d '{
"documents": [
{
"id": "doc1",
"content": "Go is an open source programming language..."
}
]
}'
```

## **::::::::: Ask Question  :::::::::::**

```curl -X POST http://localhost:8080/query
-H "Content-Type: application/json"
-d '{
"question": "What is Go programming?",
"top_k": 3
}
```
