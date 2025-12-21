/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	llamaembed "github.com/kawai-network/veridium/pkg/fantasy/providers/llama-embed"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
)

// MemoryCategory defines the type of memory
type MemoryCategory string

const (
	MemoryCategoryConversation MemoryCategory = "conversation" // From chat history
	MemoryCategoryFact         MemoryCategory = "fact"         // Extracted facts about user
	MemoryCategoryPreference   MemoryCategory = "preference"   // User preferences
	MemoryCategoryContext      MemoryCategory = "context"      // Contextual information
	MemoryCategoryTask         MemoryCategory = "task"         // Task-related memory
)

// MemoryLayer defines the memory layer (like RAM vs disk)
type MemoryLayer string

const (
	MemoryLayerWorking  MemoryLayer = "working"  // In context window (short-term)
	MemoryLayerArchived MemoryLayer = "archived" // In vector store (long-term)
)

// Memory represents a single memory fact
type Memory struct {
	ID             string         `json:"id"`
	Category       MemoryCategory `json:"category"`
	Layer          MemoryLayer    `json:"layer"`
	Type           string         `json:"type"` // More specific type within category
	Title          string         `json:"title"`
	Summary        string         `json:"summary"` // Enriched/summarized content
	Details        string         `json:"details"` // Original content
	Status         string         `json:"status"`
	AccessCount    int            `json:"access_count"`
	LastAccessedAt int64          `json:"last_accessed_at"`
	CreatedAt      int64          `json:"created_at"`
	UpdatedAt      int64          `json:"updated_at"`

	// Similarity score (filled during search)
	Similarity float64 `json:"similarity,omitempty"`
}

// MemorySearchResult represents a search result with similarity score
type MemorySearchResult struct {
	Memory     *Memory `json:"memory"`
	Similarity float64 `json:"similarity"`
}

// MemoryService manages infinite memory for conversations
type MemoryService struct {
	dbService    *database.Service
	duckDB       *DuckDBStore
	embedder     llamaembed.Embedder
	embeddingDim int
}

// MemoryServiceConfig holds configuration for memory service
type MemoryServiceConfig struct {
	DuckDB       *DuckDBStore
	Embedder     llamaembed.Embedder
	EmbeddingDim int // Default: 1024
}

// NewMemoryService creates a new memory service
func NewMemoryService(dbService *database.Service, config *MemoryServiceConfig) (*MemoryService, error) {
	if config.Embedder == nil {
		return nil, fmt.Errorf("embedder is required for memory service")
	}

	embeddingDim := config.EmbeddingDim
	if embeddingDim == 0 {
		embeddingDim = 1024
	}

	service := &MemoryService{
		dbService:    dbService,
		duckDB:       config.DuckDB,
		embedder:     config.Embedder,
		embeddingDim: embeddingDim,
	}

	log.Println("✅ Memory service initialized")
	return service, nil
}

// CreateMemory creates a new memory with embedding
func (s *MemoryService) CreateMemory(ctx context.Context, memory *Memory) (*Memory, error) {
	if memory.ID == "" {
		memory.ID = uuid.New().String()
	}

	now := time.Now().UnixMilli()
	memory.CreatedAt = now
	memory.UpdatedAt = now
	memory.LastAccessedAt = now

	if memory.Status == "" {
		memory.Status = "active"
	}

	// Generate embedding for summary
	var summaryVector []byte
	var detailsVector []byte

	if memory.Summary != "" {
		embeddings, err := s.embedder.Embed(ctx, []string{memory.Summary})
		if err != nil {
			log.Printf("⚠️  Failed to generate summary embedding: %v", err)
		} else if len(embeddings) > 0 {
			summaryVector = float32SliceToBytes(embeddings[0])
		}
	}

	if memory.Details != "" {
		embeddings, err := s.embedder.Embed(ctx, []string{memory.Details})
		if err != nil {
			log.Printf("⚠️  Failed to generate details embedding: %v", err)
		} else if len(embeddings) > 0 {
			detailsVector = float32SliceToBytes(embeddings[0])
		}
	}

	// Store in SQLite
	result, err := s.dbService.Queries().CreateUserMemory(ctx, db.CreateUserMemoryParams{
		ID:                memory.ID,
		MemoryCategory:    sql.NullString{String: string(memory.Category), Valid: true},
		MemoryLayer:       sql.NullString{String: string(memory.Layer), Valid: true},
		MemoryType:        sql.NullString{String: memory.Type, Valid: memory.Type != ""},
		Title:             sql.NullString{String: memory.Title, Valid: memory.Title != ""},
		Summary:           sql.NullString{String: memory.Summary, Valid: memory.Summary != ""},
		SummaryVector1024: summaryVector,
		Details:           sql.NullString{String: memory.Details, Valid: memory.Details != ""},
		DetailsVector1024: detailsVector,
		Status:            sql.NullString{String: memory.Status, Valid: true},
		AccessedCount:     sql.NullInt64{Int64: 0, Valid: true},
		LastAccessedAt:    now,
		CreatedAt:         now,
		UpdatedAt:         now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create memory: %w", err)
	}

	// Also store in DuckDB for vector search if available
	if s.duckDB != nil && summaryVector != nil {
		embeddings, _ := s.embedder.Embed(ctx, []string{memory.Summary})
		if len(embeddings) > 0 {
			_ = s.duckDB.UpsertVector(ctx, memory.ID, "", embeddings[0])
		}
	}

	log.Printf("✅ Memory created: %s [%s/%s]", memory.ID, memory.Category, memory.Type)

	return dbMemoryToMemory(&result), nil
}

// GetUserMemory retrieves a memory by ID
func (s *MemoryService) GetUserMemory(ctx context.Context, memoryID string) (*Memory, error) {
	result, err := s.dbService.Queries().GetUserMemory(ctx, memoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	// Update access count
	now := time.Now().UnixMilli()
	_ = s.dbService.Queries().UpdateMemoryAccessCount(ctx, db.UpdateMemoryAccessCountParams{
		LastAccessedAt: now,
		UpdatedAt:      now,
		ID:             memoryID,
	})

	return dbMemoryToMemory(&result), nil
}

// ListMemories lists memories with pagination
func (s *MemoryService) ListMemories(ctx context.Context, limit, offset int) ([]*Memory, error) {
	if limit <= 0 {
		limit = 20
	}

	results, err := s.dbService.Queries().ListUserMemories(ctx, db.ListUserMemoriesParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list memories: %w", err)
	}

	memories := make([]*Memory, len(results))
	for i, r := range results {
		memories[i] = dbMemoryToMemory(&r)
	}

	return memories, nil
}

// GetUserMemoriesByIds gets memories by IDs
func (s *MemoryService) GetUserMemoriesByIds(ctx context.Context, ids []string) ([]*Memory, error) {
	results, err := s.dbService.Queries().GetUserMemoriesByIds(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories by ids: %w", err)
	}

	memories := make([]*Memory, len(results))
	for i, r := range results {
		memories[i] = dbMemoryToMemory(&r)
	}

	return memories, nil
}

// SemanticSearch performs semantic search on memories using embeddings
func (s *MemoryService) SemanticSearch(ctx context.Context, query string, limit int) ([]*MemorySearchResult, error) {
	if limit <= 0 {
		limit = 5
	}

	// Generate query embedding
	embeddings, err := s.embedder.Embed(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return nil, fmt.Errorf("empty embedding generated")
	}

	queryEmbedding := embeddings[0]

	// Search in DuckDB if available
	if s.duckDB != nil {
		vectorResults, err := s.duckDB.SearchVectors(ctx, queryEmbedding, limit*2)
		if err != nil {
			log.Printf("⚠️  DuckDB search failed, falling back to SQLite: %v", err)
		} else if len(vectorResults) > 0 {
			// Get memory IDs and fetch from SQLite
			ids := make([]string, len(vectorResults))
			similarityMap := make(map[string]float64)
			for i, vr := range vectorResults {
				ids[i] = vr.ID
				similarityMap[vr.ID] = vr.Similarity
			}

			memories, err := s.dbService.Queries().GetUserMemoriesByIds(ctx, ids)
			if err == nil {
				results := make([]*MemorySearchResult, 0, len(memories))
				for _, m := range memories {
					memory := dbMemoryToMemory(&m)
					results = append(results, &MemorySearchResult{
						Memory:     memory,
						Similarity: similarityMap[memory.ID],
					})
				}
				return results, nil
			}
		}
	}

	// Fallback: text search
	return s.textSearch(ctx, query, limit)
}

// textSearch performs simple text-based search (fallback)
func (s *MemoryService) textSearch(ctx context.Context, query string, limit int) ([]*MemorySearchResult, error) {
	results, err := s.dbService.Queries().SearchMemoriesByTitle(ctx, db.SearchMemoriesByTitleParams{
		Column1: sql.NullString{String: query, Valid: true},
		Limit:   int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("text search failed: %w", err)
	}

	searchResults := make([]*MemorySearchResult, len(results))
	for i, r := range results {
		searchResults[i] = &MemorySearchResult{
			Memory:     dbMemoryToMemory(&r),
			Similarity: 0.5, // Default similarity for text match
		}
	}

	return searchResults, nil
}

// DeleteUserMemory deletes a memory
func (s *MemoryService) DeleteUserMemory(ctx context.Context, memoryID string) error {
	err := s.dbService.Queries().DeleteUserMemory(ctx, memoryID)
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	// Also delete from DuckDB if available
	if s.duckDB != nil {
		_ = s.duckDB.DeleteVector(ctx, memoryID)
	}

	log.Printf("✅ Memory deleted: %s", memoryID)
	return nil
}

// GetRecentMemories gets the most recent memories
func (s *MemoryService) GetRecentMemories(ctx context.Context, limit int) ([]*Memory, error) {
	if limit <= 0 {
		limit = 10
	}

	results, err := s.dbService.Queries().GetRecentMemories(ctx, int64(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get recent memories: %w", err)
	}

	memories := make([]*Memory, len(results))
	for i, r := range results {
		memories[i] = dbMemoryToMemory(&r)
	}

	return memories, nil
}

// ArchiveOldMemories archives memories that haven't been accessed recently
func (s *MemoryService) ArchiveOldMemories(ctx context.Context, olderThanDays int) error {
	cutoff := time.Now().AddDate(0, 0, -olderThanDays).UnixMilli()
	now := time.Now().UnixMilli()

	err := s.dbService.Queries().ArchiveOldMemories(ctx, db.ArchiveOldMemoriesParams{
		UpdatedAt:      now,
		LastAccessedAt: cutoff,
	})
	if err != nil {
		return fmt.Errorf("failed to archive old memories: %w", err)
	}

	log.Printf("✅ Archived memories older than %d days", olderThanDays)
	return nil
}

// FormatForLLM formats memories for inclusion in LLM context
func (s *MemoryService) FormatForLLM(memories []*MemorySearchResult) string {
	if len(memories) == 0 {
		return ""
	}

	var result string
	result = "Relevant memories from past conversations:\n\n"

	for i, m := range memories {
		result += fmt.Sprintf("%d. [%s] %s\n", i+1, m.Memory.Category, m.Memory.Title)
		if m.Memory.Summary != "" {
			result += fmt.Sprintf("   Summary: %s\n", m.Memory.Summary)
		}
		result += fmt.Sprintf("   (Relevance: %.2f)\n\n", m.Similarity)
	}

	return result
}

// Helper functions

func dbMemoryToMemory(m *db.UserMemory) *Memory {
	return &Memory{
		ID:             m.ID,
		Category:       MemoryCategory(m.MemoryCategory.String),
		Layer:          MemoryLayer(m.MemoryLayer.String),
		Type:           m.MemoryType.String,
		Title:          m.Title.String,
		Summary:        m.Summary.String,
		Details:        m.Details.String,
		Status:         m.Status.String,
		AccessCount:    int(m.AccessedCount.Int64),
		LastAccessedAt: m.LastAccessedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func float32SliceToBytes(floats []float32) []byte {
	data, _ := json.Marshal(floats)
	return data
}

func bytesToFloat32Slice(data []byte) []float32 {
	var floats []float32
	_ = json.Unmarshal(data, &floats)
	return floats
}
