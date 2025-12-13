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
	"log"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/fantasy/llamalib/tools"
	"github.com/kawai-network/veridium/fantasy/llamalib/tools/builtin"
)

// MemoryIntegration provides integration between memory services and chat
type MemoryIntegration struct {
	memoryService     *MemoryService
	enrichmentService *MemoryEnrichmentService
	bufferConfig      BufferConfig
}

// MemoryIntegrationConfig holds configuration for memory integration
type MemoryIntegrationConfig struct {
	MemoryService     *MemoryService
	EnrichmentService *MemoryEnrichmentService
	BufferConfig      *BufferConfig
}

// NewMemoryIntegration creates a new memory integration
func NewMemoryIntegration(config *MemoryIntegrationConfig) (*MemoryIntegration, error) {
	bufferConfig := DefaultBufferConfig()
	if config.BufferConfig != nil {
		bufferConfig = *config.BufferConfig
	}

	return &MemoryIntegration{
		memoryService:     config.MemoryService,
		enrichmentService: config.EnrichmentService,
		bufferConfig:      bufferConfig,
	}, nil
}

// RegisterMemoryTool registers the search_memory tool with the given registry
func (m *MemoryIntegration) RegisterMemoryTool(registry *tools.ToolRegistry) error {
	if m.memoryService == nil {
		log.Println("⚠️  Memory service not available, skipping memory tool registration")
		return nil
	}

	// Create adapter that converts MemoryService to MemorySearcher interface
	adapter := builtin.NewMemoryServiceAdapter(
		// SemanticSearch adapter
		func(ctx context.Context, query string, limit int) ([]builtin.MemorySearchResult, error) {
			results, err := m.memoryService.SemanticSearch(ctx, query, limit)
			if err != nil {
				return nil, err
			}

			searchResults := make([]builtin.MemorySearchResult, len(results))
			for i, r := range results {
				searchResults[i] = builtin.MemorySearchResult{
					ID:         r.Memory.ID,
					Category:   string(r.Memory.Category),
					Title:      r.Memory.Title,
					Summary:    r.Memory.Summary,
					Similarity: r.Similarity,
				}
			}
			return searchResults, nil
		},
		// SemanticSearchByCategory adapter (uses same function for now)
		func(ctx context.Context, query, category string, limit int) ([]builtin.MemorySearchResult, error) {
			results, err := m.memoryService.SemanticSearch(ctx, query, limit)
			if err != nil {
				return nil, err
			}

			// Filter by category
			var filtered []*MemorySearchResult
			for _, r := range results {
				if category == "" || string(r.Memory.Category) == category {
					filtered = append(filtered, r)
				}
			}

			searchResults := make([]builtin.MemorySearchResult, len(filtered))
			for i, r := range filtered {
				searchResults[i] = builtin.MemorySearchResult{
					ID:         r.Memory.ID,
					Category:   string(r.Memory.Category),
					Title:      r.Memory.Title,
					Summary:    r.Memory.Summary,
					Similarity: r.Similarity,
				}
			}
			return searchResults, nil
		},
	)

	return builtin.RegisterMemorySearch(registry, adapter)
}

// ProcessSessionBuffer processes the session buffer for auto-archiving
// Call this before processing new messages to ensure buffer doesn't overflow
func (m *MemoryIntegration) ProcessSessionBuffer(ctx context.Context, messages []fantasy.Message) ([]fantasy.Message, error) {
	if m.enrichmentService == nil {
		return messages, nil
	}

	return m.enrichmentService.AutoArchive(ctx, messages, m.bufferConfig)
}

// EnrichAndStoreMessages manually enriches messages and stores as memories
func (m *MemoryIntegration) EnrichAndStoreMessages(ctx context.Context, messages []fantasy.Message) (*EnrichmentResult, error) {
	if m.enrichmentService == nil {
		return &EnrichmentResult{}, nil
	}

	return m.enrichmentService.EnrichMessages(ctx, messages)
}

// GetRelevantMemories retrieves memories relevant to a query
func (m *MemoryIntegration) GetRelevantMemories(ctx context.Context, query string, limit int) (string, error) {
	if m.memoryService == nil {
		return "", nil
	}

	results, err := m.memoryService.SemanticSearch(ctx, query, limit)
	if err != nil {
		return "", err
	}

	return m.memoryService.FormatForLLM(results), nil
}

// BuildHybridContext builds context combining short-term buffer and long-term memory
// This implements the "RAM vs Hard Disk" analogy from MemGPT
func (m *MemoryIntegration) BuildHybridContext(ctx context.Context, currentQuery string, shortTermMessages []fantasy.Message) (string, error) {
	if m.memoryService == nil {
		return "", nil
	}

	// 1. Get relevant long-term memories based on current query
	relevantMemories, err := m.GetRelevantMemories(ctx, currentQuery, 5)
	if err != nil {
		log.Printf("⚠️  Failed to retrieve memories: %v", err)
		relevantMemories = ""
	}

	// 2. Short-term buffer is already in messages, no need to add here
	// The relevantMemories will be added to system context

	return relevantMemories, nil
}

// ArchiveOldMemories archives memories that haven't been accessed recently
func (m *MemoryIntegration) ArchiveOldMemories(ctx context.Context, olderThanDays int) error {
	if m.memoryService == nil {
		return nil
	}

	return m.memoryService.ArchiveOldMemories(ctx, olderThanDays)
}

// GetMemoryService returns the underlying memory service
func (m *MemoryIntegration) GetMemoryService() *MemoryService {
	return m.memoryService
}

// GetEnrichmentService returns the underlying enrichment service
func (m *MemoryIntegration) GetEnrichmentService() *MemoryEnrichmentService {
	return m.enrichmentService
}

// StoreConversationMemory stores a conversation exchange as memory
// This is called automatically after each chat response
func (m *MemoryIntegration) StoreConversationMemory(ctx context.Context, userMessage, assistantResponse string) error {
	if m.enrichmentService == nil {
		return nil
	}

	// Create messages from the conversation
	messages := []fantasy.Message{
		{Role: fantasy.MessageRoleUser, Content: []fantasy.MessagePart{fantasy.TextPart{Text: userMessage}}},
		{Role: fantasy.MessageRoleAssistant, Content: []fantasy.MessagePart{fantasy.TextPart{Text: assistantResponse}}},
	}

	// Enrich and store
	result, err := m.enrichmentService.EnrichMessages(ctx, messages)
	if err != nil {
		return err
	}

	if result.FactCount > 0 {
		log.Printf("🧠 [Memory] Stored %d facts from conversation", result.FactCount)
	}

	return nil
}
