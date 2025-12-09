package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kawai-network/veridium/internal/database"
)

// MockEmbedder implements llamaembed.Embedder for testing
type MockEmbedder struct {
	embeddings [][]float32
	dim        int
}

func NewMockEmbedder() *MockEmbedder {
	return &MockEmbedder{dim: 1024}
}

func (m *MockEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if m.embeddings != nil {
		return m.embeddings, nil
	}
	// Return mock embeddings
	result := make([][]float32, len(texts))
	for i := range texts {
		// Create a simple mock embedding (1024 dimensions)
		embedding := make([]float32, m.dim)
		for j := range embedding {
			embedding[j] = float32(i+j) * 0.001
		}
		result[i] = embedding
	}
	return result, nil
}

func (m *MockEmbedder) Dimensions() int {
	return m.dim
}

func (m *MockEmbedder) Close() error {
	return nil
}

func setupTestDB(t *testing.T) (*database.Service, func()) {
	t.Helper()

	// Create temp directory for test database
	tempDir, err := os.MkdirTemp("", "memory_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	dbService, err := database.NewServiceWithPath(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create database: %v", err)
	}

	cleanup := func() {
		dbService.Close()
		os.RemoveAll(tempDir)
	}

	return dbService, cleanup
}

// defaultTestUserID is the user created by database initialization
const defaultTestUserID = "DEFAULT_LOBE_CHAT_USER"

func TestMemoryService_CreateAndGet(t *testing.T) {
	dbService, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	userID := defaultTestUserID

	// Create memory service with mock embedder
	memService, err := NewMemoryService(dbService, &MemoryServiceConfig{
		Embedder:     NewMockEmbedder(),
		EmbeddingDim: 1024,
	})
	if err != nil {
		t.Fatalf("Failed to create memory service: %v", err)
	}

	// Test creating a memory
	memory := &Memory{
		Category: MemoryCategoryFact,
		Layer:    MemoryLayerArchived,
		Type:     "user_profile",
		Title:    "User Name",
		Summary:  "User's name is John Doe",
		Details:  "The user introduced themselves as John Doe from Jakarta",
	}

	created, err := memService.CreateMemory(ctx, userID, memory)
	if err != nil {
		t.Fatalf("Failed to create memory: %v", err)
	}

	if created.ID == "" {
		t.Error("Expected memory ID to be set")
	}
	if created.Title != "User Name" {
		t.Errorf("Expected title 'User Name', got '%s'", created.Title)
	}
	if created.Category != MemoryCategoryFact {
		t.Errorf("Expected category 'fact', got '%s'", created.Category)
	}

	// Test getting the memory
	retrieved, err := memService.GetMemory(ctx, userID, created.ID)
	if err != nil {
		t.Fatalf("Failed to get memory: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID '%s', got '%s'", created.ID, retrieved.ID)
	}
	if retrieved.Summary != "User's name is John Doe" {
		t.Errorf("Expected summary 'User's name is John Doe', got '%s'", retrieved.Summary)
	}
}

func TestMemoryService_ListMemories(t *testing.T) {
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

	// Create multiple memories
	memories := []*Memory{
		{Category: MemoryCategoryFact, Title: "Fact 1", Summary: "Summary 1"},
		{Category: MemoryCategoryPreference, Title: "Preference 1", Summary: "User likes blue"},
		{Category: MemoryCategoryTask, Title: "Task 1", Summary: "Working on project X"},
	}

	for _, m := range memories {
		_, err := memService.CreateMemory(ctx, userID, m)
		if err != nil {
			t.Fatalf("Failed to create memory: %v", err)
		}
	}

	// Test listing memories
	listed, err := memService.ListMemories(ctx, userID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list memories: %v", err)
	}

	if len(listed) != 3 {
		t.Errorf("Expected 3 memories, got %d", len(listed))
	}
}

func TestMemoryService_DeleteMemory(t *testing.T) {
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

	// Create a memory
	memory := &Memory{
		Category: MemoryCategoryFact,
		Title:    "To Delete",
		Summary:  "This will be deleted",
	}

	created, err := memService.CreateMemory(ctx, userID, memory)
	if err != nil {
		t.Fatalf("Failed to create memory: %v", err)
	}

	// Delete the memory
	err = memService.DeleteMemory(ctx, userID, created.ID)
	if err != nil {
		t.Fatalf("Failed to delete memory: %v", err)
	}

	// Verify it's deleted
	_, err = memService.GetMemory(ctx, userID, created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted memory")
	}
}

func TestMemoryService_TextSearch(t *testing.T) {
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

	// Create memories with searchable titles
	memories := []*Memory{
		{Category: MemoryCategoryFact, Title: "User Profile Info", Summary: "Name is John"},
		{Category: MemoryCategoryPreference, Title: "Color Preference", Summary: "Likes blue color"},
		{Category: MemoryCategoryTask, Title: "Project Alpha", Summary: "Working on Alpha"},
	}

	for _, m := range memories {
		_, err := memService.CreateMemory(ctx, userID, m)
		if err != nil {
			t.Fatalf("Failed to create memory: %v", err)
		}
	}

	// Test text search (fallback when no DuckDB)
	results, err := memService.textSearch(ctx, userID, "Profile", 10)
	if err != nil {
		t.Fatalf("Text search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'Profile', got %d", len(results))
	}

	if len(results) > 0 && results[0].Memory.Title != "User Profile Info" {
		t.Errorf("Expected 'User Profile Info', got '%s'", results[0].Memory.Title)
	}
}

func TestMemoryService_FormatForLLM(t *testing.T) {
	dbService, cleanup := setupTestDB(t)
	defer cleanup()

	memService, err := NewMemoryService(dbService, &MemoryServiceConfig{
		Embedder:     NewMockEmbedder(),
		EmbeddingDim: 1024,
	})
	if err != nil {
		t.Fatalf("Failed to create memory service: %v", err)
	}

	memories := []*MemorySearchResult{
		{
			Memory: &Memory{
				Category: MemoryCategoryFact,
				Title:    "User Name",
				Summary:  "User is John Doe",
			},
			Similarity: 0.95,
		},
		{
			Memory: &Memory{
				Category: MemoryCategoryPreference,
				Title:    "Color Preference",
				Summary:  "User prefers blue",
			},
			Similarity: 0.85,
		},
	}

	formatted := memService.FormatForLLM(memories)

	if formatted == "" {
		t.Error("Expected non-empty formatted string")
	}

	// Check that it contains key information
	if !contains(formatted, "User Name") {
		t.Error("Formatted output should contain 'User Name'")
	}
	if !contains(formatted, "0.95") {
		t.Error("Formatted output should contain similarity score")
	}
}

func TestMemoryService_GetRecentMemories(t *testing.T) {
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

	// Create memories
	for i := 0; i < 5; i++ {
		_, err := memService.CreateMemory(ctx, userID, &Memory{
			Category: MemoryCategoryFact,
			Title:    "Memory " + string(rune('A'+i)),
			Summary:  "Summary " + string(rune('A'+i)),
		})
		if err != nil {
			t.Fatalf("Failed to create memory: %v", err)
		}
	}

	// Get recent memories (limit 3)
	recent, err := memService.GetRecentMemories(ctx, userID, 3)
	if err != nil {
		t.Fatalf("Failed to get recent memories: %v", err)
	}

	if len(recent) != 3 {
		t.Errorf("Expected 3 recent memories, got %d", len(recent))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
