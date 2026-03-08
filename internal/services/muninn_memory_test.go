package services

import (
	"context"
	"strings"
	"testing"

	"github.com/scrypster/muninndb/pkg/embedded"
)

func TestMuninnMemoryBackend_StorePersistsConversation(t *testing.T) {
	ctx := context.Background()
	dataDir := t.TempDir()
	backend, err := NewMuninnMemoryBackend(dataDir, "test_vault", 1024, false)
	if err != nil {
		t.Fatalf("failed to create backend: %v", err)
	}

	const token = "muninn-e2e-token-123"
	user := "Please remember this token: " + token
	assistant := "Stored token " + token + " for future reference."

	if err := backend.StoreConversationMemory(ctx, user, assistant); err != nil {
		t.Fatalf("failed to store conversation memory: %v", err)
	}

	memories, err := backend.GetRelevantMemories(ctx, token, 5)
	if err != nil {
		t.Fatalf("failed to recall memories: %v", err)
	}
	if memories != "" && !strings.Contains(memories, token) {
		t.Fatalf("expected recalled memories to contain token %q when non-empty, got: %s", token, memories)
	}

	if err := backend.Close(); err != nil {
		t.Fatalf("failed to close backend: %v", err)
	}

	audit := embedded.NewService()
	if _, err := audit.Attach(embedded.AttachOptions{
		Name:         "audit",
		DataDir:      dataDir,
		DefaultVault: "test_vault",
		CacheSize:    1024,
		NoSync:       false,
	}); err != nil {
		t.Fatalf("failed to attach audit service: %v", err)
	}
	t.Cleanup(func() {
		if detachErr := audit.Detach("audit"); detachErr != nil {
			t.Fatalf("failed to detach audit service: %v", detachErr)
		}
	})

	stat, err := audit.Status(ctx, embedded.StatusInput{
		Connection: "audit",
		Vault:      "test_vault",
	})
	if err != nil {
		t.Fatalf("failed to read status: %v", err)
	}
	if stat.EngramCount < 1 {
		t.Fatalf("expected at least 1 engram, got %d", stat.EngramCount)
	}
}

func TestMuninnMemoryBackend_StoreAndRecallConversation(t *testing.T) {
	ctx := context.Background()
	backend, err := NewMuninnMemoryBackend(t.TempDir(), "test_vault", 1024, true)
	if err != nil {
		t.Fatalf("failed to create backend: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := backend.Close(); closeErr != nil {
			t.Fatalf("failed to close backend: %v", closeErr)
		}
	})

	const token = "muninn-recall-token-456"
	user := "Remember token " + token
	assistant := "Token " + token + " acknowledged."
	if err := backend.StoreConversationMemory(ctx, user, assistant); err != nil {
		t.Fatalf("failed to store conversation memory: %v", err)
	}

	memories, err := backend.GetRelevantMemories(ctx, token, 5)
	if err != nil {
		t.Fatalf("failed to recall memories: %v", err)
	}
	if memories != "" && !strings.Contains(memories, token) {
		t.Fatalf("expected recalled memories to contain token %q, got: %s", token, memories)
	}
}

func TestMuninnMemoryBackend_GetRelevantMemoriesEmptyQuery(t *testing.T) {
	backend, err := NewMuninnMemoryBackend(t.TempDir(), "test_vault", 1024, true)
	if err != nil {
		t.Fatalf("failed to create backend: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := backend.Close(); closeErr != nil {
			t.Fatalf("failed to close backend: %v", closeErr)
		}
	})

	memories, err := backend.GetRelevantMemories(context.Background(), "   ", 5)
	if err != nil {
		t.Fatalf("unexpected error for empty query: %v", err)
	}
	if memories != "" {
		t.Fatalf("expected empty memories for empty query, got: %q", memories)
	}
}
