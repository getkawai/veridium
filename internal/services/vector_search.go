package services

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/kawai-network/veridium/pkg/chromem"
)

// ChunkData represents chunk data for vector storage
type ChunkData struct {
	ID       string            `json:"id"`
	Text     string            `json:"text"`
	FileID   string            `json:"fileId"`
	FileName string            `json:"fileName"`
	Type     string            `json:"type"`
	Index    int               `json:"index"`
	Metadata map[string]string `json:"metadata"`
}

// SearchResult represents a search result from vector database
type SearchResult struct {
	ID         string            `json:"id"`
	Similarity float32           `json:"similarity"`
	Text       string            `json:"text"`
	FileID     string            `json:"fileId"`
	FileName   string            `json:"fileName"`
	Type       string            `json:"type"`
	Index      int               `json:"index"`
	Metadata   map[string]string `json:"metadata"`
}

// VectorSearchService handles vector search operations using chromem
type VectorSearchService struct {
	db          *chromem.DB
	collections map[string]*chromem.Collection
	mu          sync.RWMutex
	embedFunc   chromem.EmbeddingFunc
}

// NewVectorSearchService creates a new vector search service
func NewVectorSearchService(persistPath string, embeddingProvider string, embeddingModel string) (*VectorSearchService, error) {
	if persistPath == "" {
		persistPath = "./data/vector-db"
	}

	// Create persistent DB with compression
	db, err := chromem.NewPersistentDB(persistPath, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create chromem DB: %w", err)
	}

	// Setup embedding function based on provider
	var embedFunc chromem.EmbeddingFunc
	switch embeddingProvider {
	case "ollama":
		if embeddingModel == "" {
			embeddingModel = "nomic-embed-text"
		}
		embedFunc = chromem.NewEmbeddingFuncOllama(
			"http://localhost:11434/api/embeddings",
			embeddingModel,
		)
	case "openai":
		if embeddingModel == "" {
			embeddingModel = "text-embedding-3-small"
		}
		embedFunc = chromem.NewEmbeddingFuncDefault() // Uses OpenAI by default
	default:
		// Default to OpenAI
		embedFunc = chromem.NewEmbeddingFuncDefault()
	}

	return &VectorSearchService{
		db:          db,
		collections: make(map[string]*chromem.Collection),
		embedFunc:   embedFunc,
	}, nil
}

// GetUserCollection gets or creates a collection for a user
func (s *VectorSearchService) GetUserCollection(ctx context.Context, userID string) (*chromem.Collection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	collectionName := fmt.Sprintf("user_%s_chunks", userID)

	// Check if collection already exists in memory
	if col, exists := s.collections[collectionName]; exists {
		return col, nil
	}

	// Get or create collection
	col, err := s.db.GetOrCreateCollection(collectionName, nil, s.embedFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create collection: %w", err)
	}

	s.collections[collectionName] = col
	return col, nil
}

// AddChunks adds chunks with embeddings to the vector database
func (s *VectorSearchService) AddChunks(ctx context.Context, userID string, chunks []ChunkData) error {
	if len(chunks) == 0 {
		return nil
	}

	col, err := s.GetUserCollection(ctx, userID)
	if err != nil {
		return err
	}

	// Convert to chromem documents
	docs := make([]chromem.Document, len(chunks))
	for i, chunk := range chunks {
		metadata := map[string]string{
			"fileId":   chunk.FileID,
			"fileName": chunk.FileName,
			"type":     chunk.Type,
			"index":    fmt.Sprintf("%d", chunk.Index),
			"userId":   userID,
		}

		// Add custom metadata if provided
		for k, v := range chunk.Metadata {
			metadata[k] = v
		}

		docs[i] = chromem.Document{
			ID:       chunk.ID,
			Content:  chunk.Text,
			Metadata: metadata,
		}
	}

	// Add documents concurrently (chromem handles threading)
	err = col.AddDocuments(ctx, docs, runtime.NumCPU())
	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}

	return nil
}

// SemanticSearch performs semantic search on chunks
func (s *VectorSearchService) SemanticSearch(ctx context.Context, userID string, query string, fileIDs []string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 30
	}

	col, err := s.GetUserCollection(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build metadata filter
	var where map[string]string
	if len(fileIDs) > 0 {
		// Note: chromem supports exact match only
		// For multiple files, we'll need to query each and merge results
		// For now, we'll use the first fileID
		where = map[string]string{
			"fileId": fileIDs[0],
		}
	}

	// Perform query
	results, err := col.Query(ctx, query, limit, where, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	// Convert to search results
	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		index := 0
		if idxStr, ok := r.Metadata["index"]; ok {
			fmt.Sscanf(idxStr, "%d", &index)
		}

		searchResults[i] = SearchResult{
			ID:         r.ID,
			Similarity: r.Similarity,
			Text:       r.Content,
			FileID:     r.Metadata["fileId"],
			FileName:   r.Metadata["fileName"],
			Type:       r.Metadata["type"],
			Index:      index,
			Metadata:   r.Metadata,
		}
	}

	return searchResults, nil
}

// SemanticSearchMultipleFiles performs semantic search across multiple files
func (s *VectorSearchService) SemanticSearchMultipleFiles(ctx context.Context, userID string, query string, fileIDs []string, limit int) ([]SearchResult, error) {
	if len(fileIDs) == 0 {
		// Search all files
		return s.SemanticSearch(ctx, userID, query, nil, limit)
	}

	// Query each file and merge results
	allResults := []SearchResult{}
	for _, fileID := range fileIDs {
		results, err := s.SemanticSearch(ctx, userID, query, []string{fileID}, limit)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, results...)
	}

	// Sort by similarity and limit
	// Simple bubble sort for small datasets
	for i := 0; i < len(allResults)-1; i++ {
		for j := 0; j < len(allResults)-i-1; j++ {
			if allResults[j].Similarity < allResults[j+1].Similarity {
				allResults[j], allResults[j+1] = allResults[j+1], allResults[j]
			}
		}
	}

	if len(allResults) > limit {
		allResults = allResults[:limit]
	}

	return allResults, nil
}

// DeleteChunks deletes chunks from the vector database
func (s *VectorSearchService) DeleteChunks(ctx context.Context, userID string, chunkIDs []string) error {
	if len(chunkIDs) == 0 {
		return nil
	}

	col, err := s.GetUserCollection(ctx, userID)
	if err != nil {
		return err
	}

	// Delete documents by IDs
	err = col.Delete(ctx, nil, nil, chunkIDs...)
	if err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	return nil
}

// DeleteChunksByFileID deletes all chunks for a specific file
func (s *VectorSearchService) DeleteChunksByFileID(ctx context.Context, userID string, fileID string) error {
	col, err := s.GetUserCollection(ctx, userID)
	if err != nil {
		return err
	}

	// Query all chunks for this file
	where := map[string]string{
		"fileId": fileID,
	}

	// Get all documents (use high limit)
	results, err := col.Query(ctx, "", 10000, where, nil)
	if err != nil {
		return fmt.Errorf("failed to query chunks for deletion: %w", err)
	}

	// Delete each document
	chunkIDs := make([]string, len(results))
	for i, r := range results {
		chunkIDs[i] = r.ID
	}

	return s.DeleteChunks(ctx, userID, chunkIDs)
}

// GetStats returns statistics about the vector database
func (s *VectorSearchService) GetStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	col, err := s.GetUserCollection(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all documents to count
	results, err := col.Query(ctx, "", 100000, nil, nil)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"totalChunks":    len(results),
		"collectionName": fmt.Sprintf("user_%s_chunks", userID),
	}

	return stats, nil
}

// Close closes the vector database
func (s *VectorSearchService) Close() error {
	// chromem DB doesn't have explicit close, but we can clear collections
	s.mu.Lock()
	defer s.mu.Unlock()
	s.collections = make(map[string]*chromem.Collection)
	return nil
}

