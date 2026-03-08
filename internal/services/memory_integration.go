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

	"github.com/getkawai/tools"
	unillm "github.com/getkawai/unillm"
)

// MemoryIntegration provides integration between memory services and chat
type MemoryIntegration struct {
	memoryService     *MemoryService
	enrichmentService *MemoryEnrichmentService
	muninnBackend     *MuninnMemoryBackend
	bufferConfig      BufferConfig
}

// MemoryIntegrationConfig holds configuration for memory integration
type MemoryIntegrationConfig struct {
	MemoryService     *MemoryService
	EnrichmentService *MemoryEnrichmentService
	MuninnBackend     *MuninnMemoryBackend
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
		muninnBackend:     config.MuninnBackend,
		bufferConfig:      bufferConfig,
	}, nil
}

// RegisterMemoryTool registers the search_memory tool with the given registry
func (m *MemoryIntegration) RegisterMemoryTool(registry *tools.ToolRegistry) error {
	_ = registry
	// Big-bang migration: memory is served by MuninnDB backend and muninn_* tools.
	// Legacy search_memory tool is intentionally disabled.
	log.Println("ℹ️  Legacy search_memory tool disabled (using MuninnDB backend)")
	return nil
}

// ProcessSessionBuffer processes the session buffer for auto-archiving
// Call this before processing new messages to ensure buffer doesn't overflow
func (m *MemoryIntegration) ProcessSessionBuffer(ctx context.Context, messages []unillm.Message) ([]unillm.Message, error) {
	if m.enrichmentService == nil {
		return messages, nil
	}

	return m.enrichmentService.AutoArchive(ctx, messages, m.bufferConfig)
}

// EnrichAndStoreMessages manually enriches messages and stores as memories
func (m *MemoryIntegration) EnrichAndStoreMessages(ctx context.Context, messages []unillm.Message) (*EnrichmentResult, error) {
	if m.enrichmentService == nil {
		return &EnrichmentResult{}, nil
	}

	return m.enrichmentService.EnrichMessages(ctx, messages)
}

// GetRelevantMemories retrieves memories relevant to a query
func (m *MemoryIntegration) GetRelevantMemories(ctx context.Context, query string, limit int) (string, error) {
	if m.muninnBackend != nil {
		return m.muninnBackend.GetRelevantMemories(ctx, query, limit)
	}
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
func (m *MemoryIntegration) BuildHybridContext(ctx context.Context, currentQuery string, shortTermMessages []unillm.Message) (string, error) {
	if m.muninnBackend != nil {
		relevantMemories, err := m.muninnBackend.GetRelevantMemories(ctx, currentQuery, 5)
		if err != nil {
			log.Printf("⚠️  Failed to retrieve Muninn memories: %v", err)
			return "", err
		}
		return relevantMemories, nil
	}
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
	_ = ctx
	_ = olderThanDays
	if m.muninnBackend != nil {
		// Muninn memory lifecycle is managed internally by scoring/activation.
		return nil
	}
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
	if m.muninnBackend != nil {
		if err := m.muninnBackend.StoreConversationMemory(ctx, userMessage, assistantResponse); err != nil {
			return err
		}
		log.Printf("🧠 [Memory] Stored conversation in MuninnDB")
		return nil
	}

	if m.enrichmentService == nil {
		return nil
	}

	// Create messages from the conversation
	messages := []unillm.Message{
		{Role: unillm.MessageRoleUser, Content: []unillm.MessagePart{unillm.TextPart{Text: userMessage}}},
		{Role: unillm.MessageRoleAssistant, Content: []unillm.MessagePart{unillm.TextPart{Text: assistantResponse}}},
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
