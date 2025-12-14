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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/pkg/xlog"
)

// MemoryEnrichmentService converts raw chat messages into enriched facts
// This implements the "Summary-Based RAG" strategy for infinite memory
type MemoryEnrichmentService struct {
	memoryService *MemoryService
	llm           fantasy.LanguageModel // For summarization
}

// MemoryEnrichmentConfig holds configuration for enrichment service
type MemoryEnrichmentConfig struct {
	MemoryService *MemoryService
	LLM           fantasy.LanguageModel
}

// NewMemoryEnrichmentService creates a new memory enrichment service
func NewMemoryEnrichmentService(config *MemoryEnrichmentConfig) (*MemoryEnrichmentService, error) {
	if config.MemoryService == nil {
		return nil, fmt.Errorf("memory service is required")
	}

	return &MemoryEnrichmentService{
		memoryService: config.MemoryService,
		llm:           config.LLM,
	}, nil
}

// BufferConfig defines the buffer configuration
type BufferConfig struct {
	MaxBufferSize    int // Max messages in context window (default: 20)
	ArchiveBatchSize int // Messages to archive at once (default: 5)
	ArchiveThreshold int // Trigger archive when buffer > threshold (default: 15)
}

// DefaultBufferConfig returns default buffer configuration
func DefaultBufferConfig() BufferConfig {
	return BufferConfig{
		MaxBufferSize:    20,
		ArchiveBatchSize: 5,
		ArchiveThreshold: 15,
	}
}

// EnrichmentResult represents the result of enrichment
type EnrichmentResult struct {
	Memories     []*Memory `json:"memories"`
	FactCount    int       `json:"fact_count"`
	SourceMsgIDs []string  `json:"source_message_ids"`
}

// ExtractedFact represents a fact extracted from messages
type ExtractedFact struct {
	Type    string `json:"type"`    // user_profile, preference, task, context
	Title   string `json:"title"`   // Short title
	Summary string `json:"summary"` // Enriched summary
	Details string `json:"details"` // Original content
}

// EnrichMessages extracts facts from messages and stores as memories
func (s *MemoryEnrichmentService) EnrichMessages(ctx context.Context, messages []fantasy.Message) (*EnrichmentResult, error) {
	if len(messages) == 0 {
		return &EnrichmentResult{}, nil
	}

	// 1. Format messages for LLM analysis
	formattedMessages := s.formatMessagesForAnalysis(messages)

	// 2. Extract facts using LLM (if available) or rule-based extraction
	var facts []ExtractedFact
	var err error

	if s.llm != nil {
		facts, err = s.extractFactsWithLLM(ctx, formattedMessages)
		if err != nil {
			xlog.Warn("⚠️  LLM extraction failed, using rule-based", "error", err)
			facts = s.extractFactsRuleBased(messages)
		}
	} else {
		facts = s.extractFactsRuleBased(messages)
	}

	// 3. Store facts as memories
	memories := make([]*Memory, 0, len(facts))
	sourceIDs := s.extractMessageIDs(messages)

	for _, fact := range facts {
		memory := &Memory{
			Category: s.factTypeToCategory(fact.Type),
			Layer:    MemoryLayerArchived,
			Type:     fact.Type,
			Title:    fact.Title,
			Summary:  fact.Summary,
			Details:  fact.Details,
			Status:   "active",
		}

		created, err := s.memoryService.CreateMemory(ctx, memory)
		if err != nil {
			xlog.Warn("⚠️  Failed to store memory", "error", err)
			continue
		}
		memories = append(memories, created)
	}

	return &EnrichmentResult{
		Memories:     memories,
		FactCount:    len(memories),
		SourceMsgIDs: sourceIDs,
	}, nil
}

// AutoArchive checks buffer size and archives old messages if needed
func (s *MemoryEnrichmentService) AutoArchive(ctx context.Context, messages []fantasy.Message, config BufferConfig) ([]fantasy.Message, error) {
	if len(messages) <= config.ArchiveThreshold {
		return messages, nil
	}

	// Get messages to archive (oldest ones)
	toArchive := messages[:config.ArchiveBatchSize]
	remaining := messages[config.ArchiveBatchSize:]

	// Enrich and store
	result, err := s.EnrichMessages(ctx, toArchive)
	if err != nil {
		return messages, fmt.Errorf("failed to archive messages: %w", err)
	}

	xlog.Info("✅ Auto-archived messages", "message_count", len(toArchive), "fact_count", result.FactCount)

	return remaining, nil
}

// formatMessagesForAnalysis formats messages for LLM analysis
func (s *MemoryEnrichmentService) formatMessagesForAnalysis(messages []fantasy.Message) string {
	var sb strings.Builder
	sb.WriteString("<conversation>\n")

	for _, msg := range messages {
		role := "unknown"
		switch msg.Role {
		case fantasy.MessageRoleUser:
			role = "user"
		case fantasy.MessageRoleAssistant:
			role = "assistant"
		case fantasy.MessageRoleSystem:
			role = "system"
		}

		content := ""
		for _, part := range msg.Content {
			if part.GetType() == fantasy.ContentTypeText {
				if textPart, ok := fantasy.AsContentType[fantasy.TextPart](part); ok {
					content = textPart.Text
					break
				}
			}
		}

		sb.WriteString(fmt.Sprintf("<%s>%s</%s>\n", role, content, role))
	}

	sb.WriteString("</conversation>")
	return sb.String()
}

// extractFactsWithLLM uses LLM to extract facts from conversation
func (s *MemoryEnrichmentService) extractFactsWithLLM(ctx context.Context, conversation string) ([]ExtractedFact, error) {
	systemPrompt := `You are a memory extraction assistant. Your task is to extract important facts from conversations that should be remembered for future reference.

Extract facts in the following categories:
- user_profile: Information about the user (name, location, job, etc.)
- preference: User preferences and likes/dislikes
- task: Tasks, projects, or ongoing work
- context: Important contextual information

For each fact, provide:
1. type: The category of the fact
2. title: A short title (max 50 chars)
3. summary: An enriched summary that resolves pronouns and provides full context
4. details: The original content that led to this fact

IMPORTANT: 
- Resolve all pronouns ("it", "that", "this") to their actual referents
- Make summaries self-contained and understandable without context
- Only extract facts that would be useful to remember

Output as JSON array:
[{"type": "...", "title": "...", "summary": "...", "details": "..."}]

If no important facts found, return: []`

	userPrompt := fmt.Sprintf("Extract important facts from this conversation:\n\n%s", conversation)

	// Build call for LLM
	call := fantasy.Call{
		Prompt: fantasy.Prompt{
			fantasy.NewSystemMessage(systemPrompt),
			fantasy.NewUserMessage(userPrompt),
		},
	}

	// Generate response
	response, err := s.llm.Generate(ctx, call)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// Parse response - get text from ResponseContent
	responseText := response.Content.Text()

	// Extract JSON from response
	responseText = s.extractJSON(responseText)

	var facts []ExtractedFact
	if err := json.Unmarshal([]byte(responseText), &facts); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return facts, nil
}

// extractFactsRuleBased extracts facts using simple rules (fallback)
func (s *MemoryEnrichmentService) extractFactsRuleBased(messages []fantasy.Message) []ExtractedFact {
	var facts []ExtractedFact

	for _, msg := range messages {
		if msg.Role != fantasy.MessageRoleUser {
			continue
		}

		content := ""
		for _, part := range msg.Content {
			if part.GetType() == fantasy.ContentTypeText {
				if textPart, ok := fantasy.AsContentType[fantasy.TextPart](part); ok {
					content = textPart.Text
					break
				}
			}
		}

		lowerContent := strings.ToLower(content)

		// Simple pattern matching for common fact types
		if strings.Contains(lowerContent, "nama saya") || strings.Contains(lowerContent, "my name is") {
			facts = append(facts, ExtractedFact{
				Type:    "user_profile",
				Title:   "User Name",
				Summary: fmt.Sprintf("User introduced themselves: %s", content),
				Details: content,
			})
		}

		if strings.Contains(lowerContent, "saya suka") || strings.Contains(lowerContent, "i like") ||
			strings.Contains(lowerContent, "saya prefer") || strings.Contains(lowerContent, "i prefer") {
			facts = append(facts, ExtractedFact{
				Type:    "preference",
				Title:   "User Preference",
				Summary: fmt.Sprintf("User expressed preference: %s", content),
				Details: content,
			})
		}

		if strings.Contains(lowerContent, "kerja") || strings.Contains(lowerContent, "work") ||
			strings.Contains(lowerContent, "project") || strings.Contains(lowerContent, "proyek") {
			facts = append(facts, ExtractedFact{
				Type:    "task",
				Title:   "Work/Project Info",
				Summary: fmt.Sprintf("User mentioned work/project: %s", content),
				Details: content,
			})
		}
	}

	return facts
}

// extractJSON extracts JSON from potentially markdown-wrapped response
func (s *MemoryEnrichmentService) extractJSON(text string) string {
	// Try to find JSON array in the text
	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")

	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}

	return "[]"
}

// factTypeToCategory converts fact type string to MemoryCategory
func (s *MemoryEnrichmentService) factTypeToCategory(factType string) MemoryCategory {
	switch factType {
	case "user_profile":
		return MemoryCategoryFact
	case "preference":
		return MemoryCategoryPreference
	case "task":
		return MemoryCategoryTask
	case "context":
		return MemoryCategoryContext
	default:
		return MemoryCategoryConversation
	}
}

// extractMessageIDs extracts message IDs from messages (if available)
func (s *MemoryEnrichmentService) extractMessageIDs(messages []fantasy.Message) []string {
	// Note: fantasy.Message doesn't have ID field by default
	// This would need to be extended or use metadata
	return []string{}
}
