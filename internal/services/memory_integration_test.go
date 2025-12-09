package services

import (
	"context"
	"testing"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/fantasy/llamalib/tools"
)

func TestMemoryIntegration_RegisterMemoryTool(t *testing.T) {
	dbService, cleanup := setupTestDB(t)
	defer cleanup()

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
	})
	if err != nil {
		t.Fatalf("Failed to create enrichment service: %v", err)
	}

	integration, err := NewMemoryIntegration(&MemoryIntegrationConfig{
		MemoryService:     memService,
		EnrichmentService: enrichService,
	})
	if err != nil {
		t.Fatalf("Failed to create memory integration: %v", err)
	}

	// Create tool registry
	registry := tools.NewToolRegistry()

	// Register memory tool
	err = integration.RegisterMemoryTool(registry, userID)
	if err != nil {
		t.Fatalf("Failed to register memory tool: %v", err)
	}

	// Verify tool was registered
	tool, exists := registry.Get("search_memory")
	if !exists {
		t.Error("Expected search_memory tool to be registered")
	}
	if tool == nil {
		t.Error("Expected tool to not be nil")
	}
}

func TestMemoryIntegration_ProcessSessionBuffer(t *testing.T) {
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
	})
	if err != nil {
		t.Fatalf("Failed to create enrichment service: %v", err)
	}

	// Create integration with low threshold for testing
	integration, err := NewMemoryIntegration(&MemoryIntegrationConfig{
		MemoryService:     memService,
		EnrichmentService: enrichService,
		BufferConfig: &BufferConfig{
			MaxBufferSize:    10,
			ArchiveBatchSize: 2,
			ArchiveThreshold: 4,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create memory integration: %v", err)
	}

	// Create messages exceeding threshold
	messages := make([]fantasy.Message, 6)
	for i := 0; i < 6; i++ {
		messages[i] = fantasy.Message{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "Saya suka belajar AI"}},
		}
	}

	remaining, err := integration.ProcessSessionBuffer(ctx, userID, messages)
	if err != nil {
		t.Fatalf("ProcessSessionBuffer failed: %v", err)
	}

	// Should have archived some messages
	if len(remaining) >= len(messages) {
		t.Errorf("Expected buffer to be reduced, got %d (original: %d)", len(remaining), len(messages))
	}
}

func TestMemoryIntegration_GetRelevantMemories(t *testing.T) {
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

	integration, err := NewMemoryIntegration(&MemoryIntegrationConfig{
		MemoryService: memService,
	})
	if err != nil {
		t.Fatalf("Failed to create memory integration: %v", err)
	}

	// First create some memories
	_, err = memService.CreateMemory(ctx, userID, &Memory{
		Category: MemoryCategoryFact,
		Title:    "User Profile",
		Summary:  "User name is Alice from Jakarta",
	})
	if err != nil {
		t.Fatalf("Failed to create memory: %v", err)
	}

	// Get relevant memories
	formattedMemories, err := integration.GetRelevantMemories(ctx, userID, "What is my name?", 5)
	if err != nil {
		t.Fatalf("GetRelevantMemories failed: %v", err)
	}

	// Should return formatted string (may be empty if semantic search not available)
	// At minimum, text search fallback should work
	t.Logf("Formatted memories: %s", formattedMemories)
}

func TestMemoryIntegration_BuildHybridContext(t *testing.T) {
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

	integration, err := NewMemoryIntegration(&MemoryIntegrationConfig{
		MemoryService: memService,
	})
	if err != nil {
		t.Fatalf("Failed to create memory integration: %v", err)
	}

	// Create some memories first
	_, _ = memService.CreateMemory(ctx, userID, &Memory{
		Category: MemoryCategoryPreference,
		Title:    "Color Preference",
		Summary:  "User prefers blue color for UI",
	})

	shortTermMessages := []fantasy.Message{
		{
			Role:    fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: "What color should I use?"}},
		},
	}

	hybridContext, err := integration.BuildHybridContext(ctx, userID, "color preference", shortTermMessages)
	if err != nil {
		t.Fatalf("BuildHybridContext failed: %v", err)
	}

	// Hybrid context is for adding to system prompt
	t.Logf("Hybrid context: %s", hybridContext)
}

func TestMemoryIntegration_NilServices(t *testing.T) {
	// Test with nil memory service (should handle gracefully)
	integration, err := NewMemoryIntegration(&MemoryIntegrationConfig{
		MemoryService: nil,
	})
	if err != nil {
		t.Fatalf("Should allow nil memory service: %v", err)
	}

	ctx := context.Background()

	// These should not panic with nil services
	_, err = integration.GetRelevantMemories(ctx, "user", "query", 5)
	if err != nil {
		t.Errorf("GetRelevantMemories should handle nil service: %v", err)
	}

	registry := tools.NewToolRegistry()
	err = integration.RegisterMemoryTool(registry, "user")
	if err != nil {
		t.Errorf("RegisterMemoryTool should handle nil service: %v", err)
	}
}
