package services

import (
	"context"
	"strings"
	"testing"
	"time"

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

	requireRecallContainsToken(t, backend, token)

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

	requireRecallContainsToken(t, backend, token)
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

func requireRecallContainsToken(t *testing.T, backend *MuninnMemoryBackend, token string) {
	t.Helper()

	ctx := context.Background()
	deadline := time.Now().Add(2 * time.Second)
	var lastMemories string

	for {
		memories, err := backend.GetRelevantMemories(ctx, token, 5)
		if err != nil {
			t.Fatalf("failed to recall memories: %v", err)
		}
		lastMemories = memories
		if strings.Contains(memories, token) {
			return
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if strings.TrimSpace(lastMemories) == "" {
		t.Fatalf("expected recalled memories to be non-empty for token %q", token)
	}
	t.Fatalf("expected recalled memories to contain token %q, got: %s", token, lastMemories)
}
