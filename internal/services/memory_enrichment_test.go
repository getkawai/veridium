package services

import (
	"context"
	"testing"

	"github.com/kawai-network/veridium/fantasy"
)

func TestMemoryEnrichmentService_ExtractFactsRuleBased(t *testing.T) {
	dbService, cleanup := setupTestDB(t)
	defer cleanup()

	memService, err := NewMemoryService(dbService, &MemoryServiceConfig{
		Embedder:     NewMockEmbedder(),
		EmbeddingDim: 1024,
	})
	if err != nil {
		t.Fatalf("Failed to create memory service: %v", err)
	}

	enrichService, err := NewMemoryEnrichmentService(&MemoryEnrichmentConfig{
		MemoryService: memService,
		LLM:           nil, // No LLM, will use rule-based extraction
	})
	if err != nil {
		t.Fatalf("Failed to create enrichment service: %v", err)
	}

	// Test rule-based extraction
	messages := []fantasy.Message{
		{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Nama saya John Doe"}},
		},
		{
			Role:    fantasy.MessageRoleAssistant,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Senang bertemu John!"}},
		},
		{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Saya suka warna biru"}},
		},
	}

	facts := enrichService.extractFactsRuleBased(messages)

	if len(facts) < 2 {
		t.Errorf("Expected at least 2 facts, got %d", len(facts))
	}

	// Check that we got user_profile and preference facts
	foundProfile := false
	foundPreference := false
	for _, f := range facts {
		if f.Type == "user_profile" {
			foundProfile = true
		}
		if f.Type == "preference" {
			foundPreference = true
		}
	}

	if !foundProfile {
		t.Error("Expected to find user_profile fact")
	}
	if !foundPreference {
		t.Error("Expected to find preference fact")
	}
}

func TestMemoryEnrichmentService_EnrichMessages(t *testing.T) {
	dbService, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := defaultTestUserID

	memService, err := NewMemoryService(dbService, &MemoryServiceConfig{
		Embedder:     NewMockEmbedder(),
		EmbeddingDim: 1024,
	})
	if err != nil {
		t.Fatalf("Failed to create memory service: %v", err)
	}

	enrichService, err := NewMemoryEnrichmentService(&MemoryEnrichmentConfig{
		MemoryService: memService,
		LLM:           nil,
	})
	if err != nil {
		t.Fatalf("Failed to create enrichment service: %v", err)
	}

	messages := []fantasy.Message{
		{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "My name is Alice and I work in Jakarta"}},
		},
		{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "I prefer dark mode for my IDE"}},
		},
	}

	result, err := enrichService.EnrichMessages(ctx, userID, messages)
	if err != nil {
		t.Fatalf("EnrichMessages failed: %v", err)
	}

	if result.FactCount == 0 {
		t.Error("Expected at least some facts to be extracted")
	}

	// Verify memories were stored
	memories, err := memService.ListMemories(ctx, userID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list memories: %v", err)
	}

	if len(memories) == 0 {
		t.Error("Expected memories to be stored in database")
	}
}

func TestMemoryEnrichmentService_AutoArchive(t *testing.T) {
	dbService, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := defaultTestUserID

	memService, err := NewMemoryService(dbService, &MemoryServiceConfig{
		Embedder:     NewMockEmbedder(),
		EmbeddingDim: 1024,
	})
	if err != nil {
		t.Fatalf("Failed to create memory service: %v", err)
	}

	enrichService, err := NewMemoryEnrichmentService(&MemoryEnrichmentConfig{
		MemoryService: memService,
		LLM:           nil,
	})
	if err != nil {
		t.Fatalf("Failed to create enrichment service: %v", err)
	}

	// Create buffer config with low threshold for testing
	config := BufferConfig{
		MaxBufferSize:    10,
		ArchiveBatchSize: 3,
		ArchiveThreshold: 5,
	}

	// Create messages that exceed threshold
	messages := make([]fantasy.Message, 8)
	for i := 0; i < 8; i++ {
		messages[i] = fantasy.Message{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Saya suka programming dengan Go"}},
		}
	}

	// Auto archive should reduce buffer
	remaining, err := enrichService.AutoArchive(ctx, userID, messages, config)
	if err != nil {
		t.Fatalf("AutoArchive failed: %v", err)
	}

	expectedRemaining := 8 - config.ArchiveBatchSize
	if len(remaining) != expectedRemaining {
		t.Errorf("Expected %d remaining messages, got %d", expectedRemaining, len(remaining))
	}
}

func TestMemoryEnrichmentService_FormatMessagesForAnalysis(t *testing.T) {
	dbService, cleanup := setupTestDB(t)
	defer cleanup()

	memService, err := NewMemoryService(dbService, &MemoryServiceConfig{
		Embedder:     NewMockEmbedder(),
		EmbeddingDim: 1024,
	})
	if err != nil {
		t.Fatalf("Failed to create memory service: %v", err)
	}

	enrichService, err := NewMemoryEnrichmentService(&MemoryEnrichmentConfig{
		MemoryService: memService,
	})
	if err != nil {
		t.Fatalf("Failed to create enrichment service: %v", err)
	}

	messages := []fantasy.Message{
		{
			Role:    fantasy.MessageRoleSystem,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "You are a helpful assistant"}},
		},
		{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Hello, my name is Bob"}},
		},
		{
			Role:    fantasy.MessageRoleAssistant,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Hi Bob! How can I help?"}},
		},
	}

	formatted := enrichService.formatMessagesForAnalysis(messages)

	// Check format contains expected tags
	if !contains(formatted, "<conversation>") {
		t.Error("Expected <conversation> tag")
	}
	if !contains(formatted, "<user>") {
		t.Error("Expected <user> tag")
	}
	if !contains(formatted, "<assistant>") {
		t.Error("Expected <assistant> tag")
	}
	if !contains(formatted, "Bob") {
		t.Error("Expected 'Bob' in formatted output")
	}
}

func TestBufferConfig_Default(t *testing.T) {
	config := DefaultBufferConfig()

	if config.MaxBufferSize != 20 {
		t.Errorf("Expected MaxBufferSize 20, got %d", config.MaxBufferSize)
	}
	if config.ArchiveBatchSize != 5 {
		t.Errorf("Expected ArchiveBatchSize 5, got %d", config.ArchiveBatchSize)
	}
	if config.ArchiveThreshold != 15 {
		t.Errorf("Expected ArchiveThreshold 15, got %d", config.ArchiveThreshold)
	}
}
