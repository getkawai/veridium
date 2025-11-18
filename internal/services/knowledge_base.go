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

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/pkg/chromem"
	chromemAdapter "github.com/kawai-network/veridium/pkg/eino-adapters/chromem"
)

// KnowledgeBaseService manages knowledge bases with RAG capabilities
type KnowledgeBaseService struct {
	db              *database.Service
	chromemDB       *chromem.DB
	embedFunc       chromem.EmbeddingFunc
	collections     map[string]*chromem.Collection
	indexers        map[string]*chromemAdapter.Indexer
	retrievers      map[string]*chromemAdapter.Retriever
	fileManagers    map[string]*chromemAdapter.FileManager
	defaultAssetDir string
}

// KnowledgeBaseConfig holds configuration for KB service
type KnowledgeBaseConfig struct {
	ChromemPath   string                // Path for vector DB persistence
	EmbeddingFunc chromem.EmbeddingFunc // Embedding function for vector generation
	AssetDir      string                // Base directory for file copies
}

// NewKnowledgeBaseService creates a new knowledge base service
func NewKnowledgeBaseService(db *database.Service, config *KnowledgeBaseConfig) (*KnowledgeBaseService, error) {
	// Create persistent chromem DB
	chromemDB, err := chromem.NewPersistentDB(config.ChromemPath, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create chromem DB: %w", err)
	}

	// Ensure asset directory exists
	if err := os.MkdirAll(config.AssetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create asset directory: %w", err)
	}

	service := &KnowledgeBaseService{
		db:              db,
		chromemDB:       chromemDB,
		embedFunc:       config.EmbeddingFunc,
		collections:     make(map[string]*chromem.Collection),
		indexers:        make(map[string]*chromemAdapter.Indexer),
		retrievers:      make(map[string]*chromemAdapter.Retriever),
		fileManagers:    make(map[string]*chromemAdapter.FileManager),
		defaultAssetDir: config.AssetDir,
	}

	// Load existing knowledge bases from DB
	if err := service.loadKnowledgeBases(context.Background()); err != nil {
		log.Printf("⚠️  Warning: Failed to load knowledge bases: %v", err)
	}

	return service, nil
}

// CreateKnowledgeBase creates a new knowledge base
func (s *KnowledgeBaseService) CreateKnowledgeBase(ctx context.Context, name, description string, userID string) (string, error) {
	// Generate unique ID
	kbID := uuid.New().String()

	// 1. Create KB record in SQLite
	now := time.Now().Unix() * 1000 // timestamp in milliseconds
	kb, err := s.db.Queries().CreateKnowledgeBase(ctx, db.CreateKnowledgeBaseParams{
		ID:          kbID,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		Avatar:      sql.NullString{},
		Type:        sql.NullString{},
		UserID:      userID,
		ClientID:    sql.NullString{},
		IsPublic:    0,
		Settings:    sql.NullString{String: "{}", Valid: true},
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create knowledge base: %w", err)
	}

	// 2. Create chromem collection
	collection, err := s.chromemDB.GetOrCreateCollection(
		kb.ID, // Use KB ID as collection name
		nil,   // No metadata
		s.embedFunc,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create collection: %w", err)
	}

	// 3. Create Eino adapters
	indexer := chromemAdapter.NewIndexer(collection)
	retriever, err := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       10,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create retriever: %w", err)
	}

	// 4. Create file manager
	assetPath := filepath.Join(s.defaultAssetDir, kb.ID)
	fileManager, err := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
		Indexer:     indexer,
		AssetDir:    assetPath,
		ChunkSize:   1500,
		OverlapSize: 300,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create file manager: %w", err)
	}

	// 5. Cache in memory
	s.collections[kb.ID] = collection
	s.indexers[kb.ID] = indexer
	s.retrievers[kb.ID] = retriever
	s.fileManagers[kb.ID] = fileManager

	log.Printf("✅ Knowledge base created: %s (%s)", name, kb.ID)
	return kb.ID, nil
}

// GetKnowledgeBase gets a knowledge base by ID
func (s *KnowledgeBaseService) GetKnowledgeBase(ctx context.Context, kbID, userID string) (db.KnowledgeBasis, error) {
	return s.db.Queries().GetKnowledgeBase(ctx, db.GetKnowledgeBaseParams{
		ID:     kbID,
		UserID: userID,
	})
}

// ListKnowledgeBases lists all knowledge bases for a user
func (s *KnowledgeBaseService) ListKnowledgeBases(ctx context.Context, userID string) ([]db.KnowledgeBasis, error) {
	return s.db.Queries().ListKnowledgeBases(ctx, userID)
}

// UpdateKnowledgeBase updates a knowledge base
func (s *KnowledgeBaseService) UpdateKnowledgeBase(ctx context.Context, kbID, name, description, userID string) error {
	now := time.Now().Unix() * 1000
	_, err := s.db.Queries().UpdateKnowledgeBase(ctx, db.UpdateKnowledgeBaseParams{
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
	// Delete from database
	if err := s.db.Queries().DeleteKnowledgeBase(ctx, db.DeleteKnowledgeBaseParams{
		ID:     kbID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("failed to delete knowledge base: %w", err)
	}

	// Remove from memory
	delete(s.collections, kbID)
	delete(s.indexers, kbID)
	delete(s.retrievers, kbID)
	delete(s.fileManagers, kbID)

	// TODO: Delete chromem collection and asset files

	log.Printf("🗑️  Knowledge base deleted: %s", kbID)
	return nil
}

// AddFileToKnowledgeBase adds a file to a knowledge base
func (s *KnowledgeBaseService) AddFileToKnowledgeBase(ctx context.Context, kbID, filePath string, metadata map[string]any, userID string) error {
	// 1. Get file manager
	fm, ok := s.fileManagers[kbID]
	if !ok {
		// Try to load the KB
		if err := s.loadKnowledgeBase(ctx, kbID, userID); err != nil {
			return fmt.Errorf("knowledge base not found: %s", kbID)
		}
		fm = s.fileManagers[kbID]
	}

	// 2. Store file (auto-parses and indexes)
	if err := fm.StoreFile(ctx, filePath, metadata); err != nil {
		return fmt.Errorf("failed to store file: %w", err)
	}

	// 3. Create file record in SQLite
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	now := time.Now().Unix() * 1000
	fileID := uuid.New().String()

	file, err := s.db.Queries().CreateFile(ctx, db.CreateFileParams{
		ID:        fileID,
		UserID:    userID,
		FileType:  fileExt,
		FileHash:  sql.NullString{},
		Name:      fileName,
		Size:      fileInfo.Size(),
		Url:       filePath,
		Source:    sql.NullString{},
		ClientID:  sql.NullString{},
		Metadata:  sql.NullString{String: "{}", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("failed to create file record: %w", err)
	}

	// 4. Link file to knowledge base
	if err := s.db.Queries().LinkKnowledgeBaseToFile(ctx, db.LinkKnowledgeBaseToFileParams{
		KnowledgeBaseID: kbID,
		FileID:          file.ID,
		UserID:          userID,
		CreatedAt:       now,
	}); err != nil {
		return fmt.Errorf("failed to link file to KB: %w", err)
	}

	log.Printf("✅ File added to KB: %s -> %s", fileName, kbID)
	return nil
}

// RemoveFileFromKnowledgeBase removes a file from a knowledge base
func (s *KnowledgeBaseService) RemoveFileFromKnowledgeBase(ctx context.Context, kbID, fileID, userID string) error {
	// Unlink from KB
	if err := s.db.Queries().UnlinkKnowledgeBaseFromFile(ctx, db.UnlinkKnowledgeBaseFromFileParams{
		KnowledgeBaseID: kbID,
		FileID:          fileID,
		UserID:          userID,
	}); err != nil {
		return fmt.Errorf("failed to unlink file: %w", err)
	}

	// TODO: Remove file chunks from chromem

	return nil
}

// QueryKnowledgeBase performs semantic search on a knowledge base
func (s *KnowledgeBaseService) QueryKnowledgeBase(ctx context.Context, kbID, query string, topK int, userID string) ([]*schema.Document, error) {
	// Get retriever
	retriever, ok := s.retrievers[kbID]
	if !ok {
		// Try to load the KB
		if err := s.loadKnowledgeBase(ctx, kbID, userID); err != nil {
			return nil, fmt.Errorf("knowledge base not found: %s", kbID)
		}
		retriever = s.retrievers[kbID]
	}

	// Set default topK
	if topK <= 0 {
		topK = 5
	}

	// Retrieve documents
	docs, err := retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	// Limit to topK
	if len(docs) > topK {
		docs = docs[:topK]
	}

	return docs, nil
}

// GetRetriever returns the Eino retriever for a knowledge base
// This is used in RAG workflows and agent tools
func (s *KnowledgeBaseService) GetRetriever(ctx context.Context, kbID, userID string) (*chromemAdapter.Retriever, error) {
	retriever, ok := s.retrievers[kbID]
	if !ok {
		// Try to load the KB
		if err := s.loadKnowledgeBase(ctx, kbID, userID); err != nil {
			return nil, fmt.Errorf("knowledge base not found: %s", kbID)
		}
		retriever = s.retrievers[kbID]
	}
	return retriever, nil
}

// GetIndexer returns the Eino indexer for a knowledge base
func (s *KnowledgeBaseService) GetIndexer(ctx context.Context, kbID, userID string) (*chromemAdapter.Indexer, error) {
	indexer, ok := s.indexers[kbID]
	if !ok {
		// Try to load the KB
		if err := s.loadKnowledgeBase(ctx, kbID, userID); err != nil {
			return nil, fmt.Errorf("knowledge base not found: %s", kbID)
		}
		indexer = s.indexers[kbID]
	}
	return indexer, nil
}

// GetFileManager returns the file manager for a knowledge base
func (s *KnowledgeBaseService) GetFileManager(ctx context.Context, kbID, userID string) (*chromemAdapter.FileManager, error) {
	fm, ok := s.fileManagers[kbID]
	if !ok {
		// Try to load the KB
		if err := s.loadKnowledgeBase(ctx, kbID, userID); err != nil {
			return nil, fmt.Errorf("knowledge base not found: %s", kbID)
		}
		fm = s.fileManagers[kbID]
	}
	return fm, nil
}

// loadKnowledgeBases loads existing KBs from database on startup
func (s *KnowledgeBaseService) loadKnowledgeBases(ctx context.Context) error {
	// TODO: Load all KBs for all users
	// For now, this is a no-op as KBs are loaded on-demand
	log.Println("📚 Knowledge bases will be loaded on-demand")
	return nil
}

// loadKnowledgeBase loads a specific KB into memory
func (s *KnowledgeBaseService) loadKnowledgeBase(ctx context.Context, kbID, userID string) error {
	// Get KB from database
	kb, err := s.db.Queries().GetKnowledgeBase(ctx, db.GetKnowledgeBaseParams{
		ID:     kbID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("failed to get knowledge base: %w", err)
	}

	// Get or create chromem collection
	collection, err := s.chromemDB.GetOrCreateCollection(
		kb.ID,
		nil,
		s.embedFunc,
	)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	// Create Eino adapters
	indexer := chromemAdapter.NewIndexer(collection)
	retriever, err := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       10,
	})
	if err != nil {
		return fmt.Errorf("failed to create retriever: %w", err)
	}

	// Create file manager
	assetPath := filepath.Join(s.defaultAssetDir, kb.ID)
	fileManager, err := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
		Indexer:     indexer,
		AssetDir:    assetPath,
		ChunkSize:   1500,
		OverlapSize: 300,
	})
	if err != nil {
		return fmt.Errorf("failed to create file manager: %w", err)
	}

	// Cache in memory
	s.collections[kb.ID] = collection
	s.indexers[kb.ID] = indexer
	s.retrievers[kb.ID] = retriever
	s.fileManagers[kb.ID] = fileManager

	log.Printf("📚 Knowledge base loaded: %s (%s)", kb.Name, kb.ID)
	return nil
}
