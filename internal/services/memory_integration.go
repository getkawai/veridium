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
	"fmt"
	"log"
	"strings"

	"github.com/getkawai/tools"
	unillm "github.com/getkawai/unillm"
)

// MemoryIntegration provides integration between memory services and chat
type MemoryIntegration struct {
	muninnBackend *MuninnMemoryBackend
	bufferConfig  BufferConfig
}

// BufferConfig defines session-buffer behavior before memory archival.
type BufferConfig struct {
	MaxBufferSize    int
	ArchiveBatchSize int
	ArchiveThreshold int
}

// DefaultBufferConfig returns default buffer configuration.
func DefaultBufferConfig() BufferConfig {
	return BufferConfig{
		MaxBufferSize:    20,
		ArchiveBatchSize: 5,
		ArchiveThreshold: 15,
	}
}

// EnrichmentResult summarizes stored memory outcomes.
type EnrichmentResult struct {
	FactCount int `json:"fact_count"`
}

// MemoryIntegrationConfig holds configuration for memory integration
type MemoryIntegrationConfig struct {
	MuninnBackend *MuninnMemoryBackend
	BufferConfig  *BufferConfig
}

// NewMemoryIntegration creates a new memory integration
func NewMemoryIntegration(config *MemoryIntegrationConfig) (*MemoryIntegration, error) {
	if config == nil || config.MuninnBackend == nil {
		return nil, ErrMuninnBackendRequired()
	}

	bufferConfig := DefaultBufferConfig()
	if config.BufferConfig != nil {
		bufferConfig = *config.BufferConfig
	}

	return &MemoryIntegration{
		muninnBackend: config.MuninnBackend,
		bufferConfig:  bufferConfig,
	}, nil
}

func ErrMuninnBackendRequired() error {
	return &memoryIntegrationError{msg: "muninn backend is required"}
}

type memoryIntegrationError struct {
	msg string
}

func (e *memoryIntegrationError) Error() string {
	return e.msg
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
	_ = ctx
	_ = m.bufferConfig
	return messages, nil
}

// EnrichAndStoreMessages manually enriches messages and stores as memories
func (m *MemoryIntegration) EnrichAndStoreMessages(ctx context.Context, messages []unillm.Message) (*EnrichmentResult, error) {
	if m == nil || m.muninnBackend == nil {
		return nil, ErrMuninnBackendRequired()
	}
	if len(messages) == 0 {
		return &EnrichmentResult{}, nil
	}

	stored := 0
	for i := 0; i+1 < len(messages); i++ {
		if messages[i].Role != unillm.MessageRoleUser || messages[i+1].Role != unillm.MessageRoleAssistant {
			continue
		}
		userText := textFromMessage(messages[i])
		assistantText := textFromMessage(messages[i+1])
		if err := m.StoreConversationMemory(ctx, userText, assistantText); err != nil {
			return nil, err
		}
		stored++
	}

	return &EnrichmentResult{FactCount: stored}, nil
}

// GetRelevantMemories retrieves memories relevant to a query
func (m *MemoryIntegration) GetRelevantMemories(ctx context.Context, query string, limit int) (string, error) {
	return m.GetRelevantMemoriesForScope(ctx, "", query, limit)
}

func (m *MemoryIntegration) GetRelevantMemoriesForScope(ctx context.Context, scopeKey, query string, limit int) (string, error) {
	if m == nil || m.muninnBackend == nil {
		return "", ErrMuninnBackendRequired()
	}
	return m.muninnBackend.GetRelevantMemoriesForScope(ctx, scopeKey, query, limit)
}

// BuildHybridContext builds context combining short-term buffer and long-term memory
// This implements the "RAM vs Hard Disk" analogy from MemGPT
func (m *MemoryIntegration) BuildHybridContext(ctx context.Context, currentQuery string, shortTermMessages []unillm.Message) (string, error) {
	return m.BuildHybridContextForScope(ctx, "", currentQuery, shortTermMessages)
}

func (m *MemoryIntegration) BuildHybridContextForScope(ctx context.Context, scopeKey, currentQuery string, shortTermMessages []unillm.Message) (string, error) {
	if m == nil || m.muninnBackend == nil {
		return "", ErrMuninnBackendRequired()
	}

	relevantMemories, err := m.muninnBackend.GetRelevantMemoriesForScope(ctx, scopeKey, currentQuery, 5)
	if err != nil {
		log.Printf("⚠️  Failed to retrieve Muninn memories: %v", err)
		return "", err
	}

	shortTermContext := serializeShortTermMessages(shortTermMessages)
	switch {
	case shortTermContext == "" && relevantMemories == "":
		return "", nil
	case shortTermContext == "":
		return relevantMemories, nil
	case relevantMemories == "":
		return shortTermContext, nil
	default:
		return shortTermContext + "\n\n--- Long-Term Relevant Memories ---\n\n" + relevantMemories, nil
	}
}

// ArchiveOldMemories archives memories that haven't been accessed recently
func (m *MemoryIntegration) ArchiveOldMemories(ctx context.Context, olderThanDays int) error {
	_ = ctx
	_ = olderThanDays
	if m == nil || m.muninnBackend == nil {
		return ErrMuninnBackendRequired()
	}
	// Muninn memory lifecycle is managed internally by scoring/activation.
	return nil
}

// UsesMuninnBackend reports whether this integration is running in Muninn mode.
func (m *MemoryIntegration) UsesMuninnBackend() bool {
	return m != nil && m.muninnBackend != nil
}

// StoreConversationMemory stores a conversation exchange as memory
// This is called automatically after each chat response
func (m *MemoryIntegration) StoreConversationMemory(ctx context.Context, userMessage, assistantResponse string) error {
	return m.StoreConversationMemoryForScope(ctx, "", userMessage, assistantResponse)
}

func (m *MemoryIntegration) StoreConversationMemoryForScope(ctx context.Context, scopeKey, userMessage, assistantResponse string) error {
	if m == nil || m.muninnBackend == nil {
		return ErrMuninnBackendRequired()
	}
	if err := m.muninnBackend.StoreConversationMemoryForScope(ctx, scopeKey, userMessage, assistantResponse); err != nil {
		return err
	}
	log.Printf("🧠 [Memory] Stored conversation in MuninnDB")
	return nil
}

func textFromMessage(msg unillm.Message) string {
	for _, part := range msg.Content {
		if part.GetType() != unillm.ContentTypeText {
			continue
		}
		if textPart, ok := unillm.AsContentType[unillm.TextPart](part); ok {
			return textPart.Text
		}
	}
	return ""
}

func serializeShortTermMessages(messages []unillm.Message) string {
	if len(messages) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Current session context:\n")
	for i := range messages {
		content := strings.TrimSpace(textFromMessage(messages[i]))
		if content == "" {
			continue
		}
		b.WriteString(fmt.Sprintf("- %s: %s\n", roleLabel(messages[i].Role), content))
	}
	return strings.TrimSpace(b.String())
}

func roleLabel(role unillm.MessageRole) string {
	switch role {
	case unillm.MessageRoleUser:
		return "user"
	case unillm.MessageRoleAssistant:
		return "assistant"
	case unillm.MessageRoleSystem:
		return "system"
	default:
		return "unknown"
	}
}
