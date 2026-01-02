package services

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/kawai-network/veridium/internal/database/generated"
	llamaembed "github.com/kawai-network/veridium/pkg/fantasy/providers/llama-embed"
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
	embedder llamaembed.Embedder
}

// NewVectorSearchService creates a new vector search service
func NewVectorSearchService(
	database *sql.DB,
	duckDB *DuckDBStore,
	embedder llamaembed.Embedder,
) (*VectorSearchService, error) {
	if embedder == nil {
		return nil, fmt.Errorf("embedder is required")
	}

	return &VectorSearchService{
		queries:  db.New(database),
		duckDB:   duckDB,
		embedder: embedder,
	}, nil
}

// SemanticSearch performs semantic search on chunks using DuckDB + SQLite
func (s *VectorSearchService) SemanticSearch(ctx context.Context, query string, fileIDs []string, limit int) ([]SearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if limit <= 0 {
		limit = 30
	}

	// 1. Generate embedding for the query
	embeddings, err := s.embedder.Embed(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return nil, fmt.Errorf("empty embedding generated")
	}

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
	chunks, err := s.queries.GetChunksByIDs(ctx, ids)
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
			file, err := s.queries.GetFile(ctx, chunk.FileID.String)
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
func (s *VectorSearchService) SemanticSearchMultipleFiles(ctx context.Context, query string, fileIDs []string, limit int) ([]SearchResult, error) {
	// Just call SemanticSearch with fileIDs filter
	return s.SemanticSearch(ctx, query, fileIDs, limit)
}

// GetEmbedder returns the embedder
func (s *VectorSearchService) GetEmbedder() llamaembed.Embedder {
	return s.embedder
}

// Close closes the vector search service
func (s *VectorSearchService) Close() error {
	// Nothing to close for now
	return nil
}
