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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/types"
)

// KnowledgeBaseService manages knowledge bases with RAG capabilities
type KnowledgeBaseService struct {
	dbService       *database.Service
	ragProcessor    *RAGProcessor
	vectorSearch    *VectorSearchService
	fileLoader      *FileLoader
	defaultAssetDir string
}

// KnowledgeBaseConfig holds configuration for KB service
type KnowledgeBaseConfig struct {
	RAGProcessor *RAGProcessor
	VectorSearch *VectorSearchService
	FileLoader   *FileLoader
	AssetDir     string // Base directory for file copies
}

// NewKnowledgeBaseService creates a new knowledge base service
func NewKnowledgeBaseService(dbService *database.Service, config *KnowledgeBaseConfig) (*KnowledgeBaseService, error) {
	// Ensure asset directory exists
	if err := os.MkdirAll(config.AssetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create asset directory: %w", err)
	}

	service := &KnowledgeBaseService{
		dbService:       dbService,
		ragProcessor:    config.RAGProcessor,
		vectorSearch:    config.VectorSearch,
		fileLoader:      config.FileLoader,
		defaultAssetDir: config.AssetDir,
	}

	log.Println("✅ Knowledge Base service initialized (DuckDB + SQLite)")
	return service, nil
}

// CreateKnowledgeBase creates a new knowledge base
func (s *KnowledgeBaseService) CreateKnowledgeBase(ctx context.Context, name, description string, userID string) (string, error) {
	// Generate unique ID
	kbID := uuid.New().String()

	// Create KB record in SQLite
	_, err := s.dbService.Queries().CreateKnowledgeBase(ctx, db.CreateKnowledgeBaseParams{
		ID:          kbID,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		Avatar:      sql.NullString{},
		Type:        sql.NullString{},
		UserID:      userID,
		IsPublic:    0,
		Settings:    sql.NullString{String: "{}", Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create knowledge base: %w", err)
	}

	log.Printf("✅ Knowledge base created: %s (%s)", name, kbID)
	return kbID, nil
}

// GetKnowledgeBase gets a knowledge base by ID
func (s *KnowledgeBaseService) GetKnowledgeBase(ctx context.Context, kbID, userID string) (db.KnowledgeBasis, error) {
	return s.dbService.Queries().GetKnowledgeBase(ctx, db.GetKnowledgeBaseParams{
		ID:     kbID,
		UserID: userID,
	})
}

// ListKnowledgeBases lists all knowledge bases for a user
func (s *KnowledgeBaseService) ListKnowledgeBases(ctx context.Context, userID string) ([]db.KnowledgeBasis, error) {
	return s.dbService.Queries().ListKnowledgeBases(ctx, userID)
}

// UpdateKnowledgeBase updates a knowledge base
func (s *KnowledgeBaseService) UpdateKnowledgeBase(ctx context.Context, kbID, name, description, userID string) error {
	now := time.Now().Unix() * 1000
	_, err := s.dbService.Queries().UpdateKnowledgeBase(ctx, db.UpdateKnowledgeBaseParams{
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		Avatar:      sql.NullString{},
		Settings:    sql.NullString{String: "{}", Valid: true},
		UpdatedAt:   now,
		ID:          kbID,
		UserID:      userID,
	})
	return err
}

// DeleteKnowledgeBase deletes a knowledge base
func (s *KnowledgeBaseService) DeleteKnowledgeBase(ctx context.Context, kbID, userID string) error {
	// Delete from database (cascades to files and chunks)
	if err := s.dbService.Queries().DeleteKnowledgeBase(ctx, db.DeleteKnowledgeBaseParams{
		ID:     kbID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("failed to delete knowledge base: %w", err)
	}

	// TODO: Delete vectors from DuckDB
	// TODO: Delete asset files

	log.Printf("🗑️  Knowledge base deleted: %s", kbID)
	return nil
}

// AddFileToKnowledgeBase adds a file to a knowledge base
func (s *KnowledgeBaseService) AddFileToKnowledgeBase(ctx context.Context, kbID, filePath string, metadata map[string]any, userID string) error {
	// 1. Load and parse file
	fileDoc, err := s.fileLoader.LoadFile(filePath, nil)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	// 2. Create file record in SQLite
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	now := time.Now().Unix() * 1000
	fileID := uuid.New().String()

	file, err := s.dbService.Queries().CreateFile(ctx, db.CreateFileParams{
		ID:       fileID,
		UserID:   userID,
		FileType: fileExt,
		FileHash: sql.NullString{},
		Name:     fileName,
		Size:     fileInfo.Size(),
		Url:      filePath,
		Source:   sql.NullString{},
		Metadata: sql.NullString{String: "{}", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create file record: %w", err)
	}

	// 3. Create document record
	documentID := uuid.New().String()
	_, err = s.dbService.Queries().CreateDocument(ctx, db.CreateDocumentParams{
		ID:             documentID,
		Title:          sql.NullString{String: fileName, Valid: true},
		Content:        sql.NullString{String: fileDoc.Content, Valid: true},
		FileType:       fileDoc.FileType,
		Filename:       sql.NullString{String: fileName, Valid: true},
		TotalCharCount: int64(fileDoc.TotalCharCount),
		TotalLineCount: int64(fileDoc.TotalLineCount),
		Metadata:       sql.NullString{String: "{}", Valid: true},
		Pages:          sql.NullString{},
		SourceType:     "file",
		Source:         filePath,
		FileID:         sql.NullString{String: fileID, Valid: true},
		UserID:         userID,
		EditorData:     sql.NullString{},
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	// 4. Process for RAG (chunking + embedding + indexing)
	_, err = s.ragProcessor.ProcessFile(ctx, RAGProcessRequest{
		FilePath:   filePath,
		FileID:     fileID,
		DocumentID: documentID,
		UserID:     userID,
		Filename:   fileName,
	})
	if err != nil {
		return fmt.Errorf("failed to process file for RAG: %w", err)
	}

	// 5. Link file to knowledge base
	if err := s.dbService.Queries().LinkKnowledgeBaseToFile(ctx, db.LinkKnowledgeBaseToFileParams{
		KnowledgeBaseID: kbID,
		FileID:          file.ID,
		UserID:          userID,
	}); err != nil {
		return fmt.Errorf("failed to link file to KB: %w", err)
	}

	log.Printf("✅ File added to KB: %s -> %s", fileName, kbID)
	return nil
}

// RemoveFileFromKnowledgeBase removes a file from a knowledge base
func (s *KnowledgeBaseService) RemoveFileFromKnowledgeBase(ctx context.Context, kbID, fileID, userID string) error {
	// Unlink from KB
	if err := s.dbService.Queries().UnlinkKnowledgeBaseFromFile(ctx, db.UnlinkKnowledgeBaseFromFileParams{
		KnowledgeBaseID: kbID,
		FileID:          fileID,
		UserID:          userID,
	}); err != nil {
		return fmt.Errorf("failed to unlink file: %w", err)
	}

	// TODO: Remove file chunks from DuckDB

	return nil
}

// QueryKnowledgeBase performs semantic search on a knowledge base
func (s *KnowledgeBaseService) QueryKnowledgeBase(ctx context.Context, kbID, query string, topK int, userID string) ([]*types.Document, error) {
	// Set default topK
	if topK <= 0 {
		topK = 5
	}

	// Get files in this KB
	kbFiles, err := s.dbService.Queries().ListKnowledgeBaseFiles(ctx, db.ListKnowledgeBaseFilesParams{
		KnowledgeBaseID: kbID,
		UserID:          userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get KB files: %w", err)
	}

	if len(kbFiles) == 0 {
		return []*types.Document{}, nil
	}

	// Extract file IDs
	fileIDs := make([]string, len(kbFiles))
	for i, f := range kbFiles {
		fileIDs[i] = f.FileID
	}

	// Perform semantic search
	results, err := s.vectorSearch.SemanticSearchMultipleFiles(ctx, userID, query, fileIDs, topK)
	if err != nil {
		return nil, fmt.Errorf("semantic search failed: %w", err)
	}

	// Convert to Eino Document format for compatibility
	docs := make([]*types.Document, len(results))
	for i, r := range results {
		docs[i] = &types.Document{
			ID:      r.ID,
			Content: r.Text,
			Metadata: map[string]any{
				"file_id":    r.FileID,
				"file_name":  r.FileName,
				"type":       r.Type,
				"index":      r.Index,
				"similarity": r.Similarity,
			},
		}
	}

	return docs, nil
}

// GetRetriever returns a simple retriever function for compatibility
// This replaces the Eino chromem adapter
func (s *KnowledgeBaseService) GetRetriever(ctx context.Context, kbID, userID string) (func(context.Context, string) ([]*types.Document, error), error) {
	// Return a closure that captures kbID and userID
	retriever := func(ctx context.Context, query string) ([]*types.Document, error) {
		return s.QueryKnowledgeBase(ctx, kbID, query, 10, userID)
	}
	return retriever, nil
}
