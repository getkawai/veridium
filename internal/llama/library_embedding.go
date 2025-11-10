package llama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// LibraryEmbeddingService provides embedding functionality using llama library
type LibraryEmbeddingService struct {
	libService *LibraryService
}

// NewLibraryEmbeddingService creates a new library-based embedding service
func NewLibraryEmbeddingService(libService *LibraryService) *LibraryEmbeddingService {
	return &LibraryEmbeddingService{
		libService: libService,
	}
}

// EmbeddingRequest represents an OpenAI-compatible embedding request
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"` // Can be string or array of strings
}

// EmbeddingResponse represents the embedding response
type EmbeddingResponse struct {
	Object string              `json:"object"`
	Data   []EmbeddingData     `json:"data"`
	Model  string              `json:"model"`
	Usage  *EmbeddingUsage     `json:"usage"`
}

// EmbeddingData represents a single embedding
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingUsage represents token usage for embeddings
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// CreateEmbedding creates embeddings for the given input texts
func (e *LibraryEmbeddingService) CreateEmbedding(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	// Ensure embedding model is loaded
	if !e.libService.IsEmbeddingModelLoaded() {
		log.Println("📥 Embedding model not loaded, loading now...")
		if err := e.libService.LoadEmbeddingModel(""); err != nil {
			return nil, fmt.Errorf("failed to load embedding model: %w", err)
		}
	}
	
	// Process each input text
	var embeddings []EmbeddingData
	totalTokens := 0
	
	for i, text := range req.Input {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		// Generate embedding
		embedding, err := e.libService.GenerateEmbedding(text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for input %d: %w", i, err)
		}
		
		embeddings = append(embeddings, EmbeddingData{
			Object:    "embedding",
			Embedding: embedding,
			Index:     i,
		})
		
		// Rough token estimate
		totalTokens += len(text) / 4
	}
	
	return &EmbeddingResponse{
		Object: "list",
		Data:   embeddings,
		Model:  e.libService.GetLoadedEmbeddingModel(),
		Usage: &EmbeddingUsage{
			PromptTokens: totalTokens,
			TotalTokens:  totalTokens,
		},
	}, nil
}

// HandleEmbeddingRequest handles an embedding request from proxy
func (e *LibraryEmbeddingService) HandleEmbeddingRequest(ctx context.Context, body string) (*ProxyResponse, error) {
	// Parse request
	var req EmbeddingRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return nil, fmt.Errorf("failed to parse embedding request: %w", err)
	}
	
	// Generate embeddings
	resp, err := e.CreateEmbedding(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// Marshal response
	respData, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	
	return &ProxyResponse{
		Status:     200,
		StatusText: "OK",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(respData),
	}, nil
}

// BatchEmbedding generates embeddings for multiple texts efficiently
func (e *LibraryEmbeddingService) BatchEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	// Ensure embedding model is loaded
	if !e.libService.IsEmbeddingModelLoaded() {
		if err := e.libService.LoadEmbeddingModel(""); err != nil {
			return nil, fmt.Errorf("failed to load embedding model: %w", err)
		}
	}
	
	results := make([][]float32, len(texts))
	
	for i, text := range texts {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		embedding, err := e.libService.GenerateEmbedding(text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}
		
		results[i] = embedding
	}
	
	return results, nil
}

// GetEmbeddingDimension returns the dimension of embeddings from the loaded model
func (e *LibraryEmbeddingService) GetEmbeddingDimension() (int32, error) {
	if !e.libService.IsEmbeddingModelLoaded() {
		return 0, fmt.Errorf("embedding model not loaded")
	}
	
	e.libService.embMutex.Lock()
	defer e.libService.embMutex.Unlock()
	
	// For now, generate a test embedding to get dimension
	testEmbed, err := e.libService.GenerateEmbedding("test")
	if err != nil {
		return 0, err
	}
	
	return int32(len(testEmbed)), nil
}

// LibraryProxyServiceWithEmbedding extends LibraryProxyService with embedding support
func (p *LibraryProxyService) HandleEmbeddingEndpoint(ctx context.Context, request ProxyRequest) (*ProxyResponse, error) {
	embService := NewLibraryEmbeddingService(p.libService)
	
	// Check if this is an embedding request
	if strings.HasSuffix(request.Path, "/embeddings") || strings.HasSuffix(request.Path, "/v1/embeddings") {
		return embService.HandleEmbeddingRequest(ctx, request.Body)
	}
	
	return nil, fmt.Errorf("not an embedding endpoint")
}

// Update the Fetch method to handle embeddings
func (p *LibraryProxyService) FetchWithEmbeddings(ctx context.Context, request ProxyRequest) (*ProxyResponse, error) {
	log.Printf("📍 [LibraryProxyService.FetchWithEmbeddings] method=%s, path=%s", request.Method, request.Path)
	
	// Check if it's an embedding request
	if strings.HasSuffix(request.Path, "/embeddings") || strings.HasSuffix(request.Path, "/v1/embeddings") {
		return p.HandleEmbeddingEndpoint(ctx, request)
	}
	
	// Otherwise, handle as chat
	if strings.HasSuffix(request.Path, "/chat/completions") || strings.HasSuffix(request.Path, "/v1/chat/completions") {
		var req ChatCompletionRequest
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return nil, fmt.Errorf("failed to parse request: %w", err)
		}
		
		req.Stream = false
		
		resp, err := p.chatService.ChatCompletion(ctx, req)
		if err != nil {
			return nil, err
		}
		
		respData, err := json.Marshal(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		
		return &ProxyResponse{
			Status:     200,
			StatusText: "OK",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(respData),
		}, nil
	}
	
	return nil, fmt.Errorf("endpoint not supported: %s", request.Path)
}

