package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/xlog"
)

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

// VectorSearchService handles vector search operations using DuckDB + SQLite
type VectorSearchService struct {
	queries  *db.Queries
	duckDB   *DuckDBStore
	embedder embedding.Embedder // Eino embedding interface
}

// NewVectorSearchService creates a new vector search service
func NewVectorSearchService(
	database *sql.DB,
	duckDB *DuckDBStore,
	embeddingProvider string,
	embeddingModel string,
	libService *llama.LibraryService,
) (*VectorSearchService, error) {
	// Setup embedding using Eino interface
	var embedder embedding.Embedder
	var err error

	switch embeddingProvider {
	case "llama", "":
		// Default to llama.cpp (local, no API key needed)
		xlog.Info("Using llama.cpp embedding provider (Eino)")

		if embeddingModel == "" {
			embeddingModel = llama.GetRecommendedEmbeddingModel()
		}

		// Get model info from catalog to get correct filename
		model, exists := llama.GetEmbeddingModel(embeddingModel)
		if !exists {
			return nil, fmt.Errorf("embedding model not found in catalog: %s", embeddingModel)
		}

		// Get full model path (library is already loaded by LibraryService)
		installer := llama.NewLlamaCppInstaller()
		modelPath := installer.GetModelsDirectory() + "/" + model.Filename

		// Create Eino Llama embedder
		embedder, err = llama.NewEmbedder(context.Background(), &llama.EmbeddingConfig{
			ModelPath:       modelPath,
			SkipLibraryInit: true, // Library already loaded by LibraryService
			ContextSize:     2048,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create Llama embedder: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s (only 'llama' is supported with Eino)", embeddingProvider)
	}

	return &VectorSearchService{
		queries:  db.New(database),
		duckDB:   duckDB,
		embedder: embedder,
	}, nil
}

// SemanticSearch performs semantic search on chunks using DuckDB + SQLite
func (s *VectorSearchService) SemanticSearch(ctx context.Context, userID string, query string, fileIDs []string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 30
	}

	// 1. Generate embedding for the query using Eino
	embeddings, err := s.embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return nil, fmt.Errorf("empty embedding generated")
	}

	// Use embedding directly ([]float64 from Eino)
	embedding := embeddings[0]

	// 2. Search DuckDB for similar vectors
	if s.duckDB == nil {
		return nil, fmt.Errorf("DuckDB not initialized")
	}

	// Fetch more results than limit to account for potential filtering
	vectorResults, err := s.duckDB.SearchVectors(ctx, embedding, limit*2)
	if err != nil {
		return nil, fmt.Errorf("DuckDB search failed: %w", err)
	}

	if len(vectorResults) == 0 {
		return []SearchResult{}, nil
	}

	// Extract IDs and create similarity map
	ids := make([]string, len(vectorResults))
	similarityMap := make(map[string]float64)
	for i, vr := range vectorResults {
		ids[i] = vr.ID
		similarityMap[vr.ID] = vr.Similarity
	}

	// 3. Fetch full content from SQLite
	chunks, err := s.queries.GetChunksByIDs(ctx, db.GetChunksByIDsParams{
		Ids:    ids,
		UserID: sql.NullString{String: userID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chunks from SQLite: %w", err)
	}

	// 4. Convert to SearchResult and Filter
	results := make([]SearchResult, 0, len(chunks))
	for _, chunk := range chunks {
		// Filter by fileIDs if needed
		if len(fileIDs) > 0 {
			found := false
			if chunk.FileID.Valid {
				for _, fid := range fileIDs {
					if chunk.FileID.String == fid {
						found = true
						break
					}
				}
			}
			if !found {
				continue
			}
		}

		// Parse metadata
		metadata := make(map[string]string)
		if chunk.Metadata.Valid {
			metadata["raw"] = chunk.Metadata.String
		}
		if chunk.FileID.Valid {
			metadata["fileId"] = chunk.FileID.String
		}

		// Get similarity from map
		similarity := float32(similarityMap[chunk.ID])

		// Get filename from files table
		fileName := ""
		if chunk.FileID.Valid {
			file, err := s.queries.GetFile(ctx, db.GetFileParams{
				ID:     chunk.FileID.String,
				UserID: userID,
			})
			if err == nil {
				fileName = file.Name
			}
		}

		results = append(results, SearchResult{
			ID:         chunk.ID,
			Similarity: similarity,
			Text:       chunk.Text.String,
			FileID:     chunk.FileID.String,
			FileName:   fileName,
			Type:       chunk.Type.String,
			Index:      int(chunk.ChunkIndex.Int64),
			Metadata:   metadata,
		})

		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

// SemanticSearchMultipleFiles performs semantic search across multiple files
func (s *VectorSearchService) SemanticSearchMultipleFiles(ctx context.Context, userID string, query string, fileIDs []string, limit int) ([]SearchResult, error) {
	// Just call SemanticSearch with fileIDs filter
	return s.SemanticSearch(ctx, userID, query, fileIDs, limit)
}

// GetEmbedder returns the Eino embedder
func (s *VectorSearchService) GetEmbedder() embedding.Embedder {
	return s.embedder
}

// Close closes the vector search service
func (s *VectorSearchService) Close() error {
	// Nothing to close for now
	return nil
}
