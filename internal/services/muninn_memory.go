package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/scrypster/muninndb/pkg/embedded"
)

// MuninnMemoryBackend is the single memory backend for Veridium.
type MuninnMemoryBackend struct {
	service      *embedded.Service
	connection   string
	defaultVault string
}

func NewMuninnMemoryBackend(dataDir string, defaultVault string, cacheSize int, noSync bool) (*MuninnMemoryBackend, error) {
	if strings.TrimSpace(dataDir) == "" {
		return nil, fmt.Errorf("muninn data dir is required")
	}
	if strings.TrimSpace(defaultVault) == "" {
		defaultVault = "default"
	}

	service := embedded.NewService()
	connection := "veridium-memory"
	if _, err := service.Attach(embedded.AttachOptions{
		Name:         connection,
		DataDir:      dataDir,
		DefaultVault: defaultVault,
		CacheSize:    cacheSize,
		NoSync:       noSync,
	}); err != nil {
		return nil, fmt.Errorf("attach muninn backend: %w", err)
	}

	return &MuninnMemoryBackend{
		service:      service,
		connection:   connection,
		defaultVault: defaultVault,
	}, nil
}

func (m *MuninnMemoryBackend) Close() error {
	if m.service == nil {
		return nil
	}
	return m.service.Detach(m.connection)
}

func (m *MuninnMemoryBackend) StoreConversationMemory(ctx context.Context, userMessage, assistantResponse string) error {
	return m.StoreConversationMemoryForScope(ctx, "", userMessage, assistantResponse)
}

func (m *MuninnMemoryBackend) StoreConversationMemoryForScope(ctx context.Context, scopeKey, userMessage, assistantResponse string) error {
	userMessage = strings.TrimSpace(userMessage)
	assistantResponse = strings.TrimSpace(assistantResponse)
	if userMessage == "" && assistantResponse == "" {
		return nil
	}

	concept := truncateForMemory(userMessage, 120)
	if concept == "" {
		concept = truncateForMemory(assistantResponse, 120)
	}
	if concept == "" {
		concept = "conversation memory"
	}

	content := fmt.Sprintf("user: %s\nassistant: %s", userMessage, assistantResponse)
	summary := truncateForMemory(assistantResponse, 400)

	_, err := m.service.Remember(ctx, embedded.RememberInput{
		Connection: m.connection,
		Vault:      m.vaultForScope(scopeKey),
		Concept:    concept,
		Content:    content,
		Summary:    summary,
		Tags:       []string{"chat", "conversation"},
		TypeLabel:  "conversation",
	})
	return err
}

func (m *MuninnMemoryBackend) GetRelevantMemories(ctx context.Context, query string, limit int) (string, error) {
	return m.GetRelevantMemoriesForScope(ctx, "", query, limit)
}

func (m *MuninnMemoryBackend) GetRelevantMemoriesForScope(ctx context.Context, scopeKey, query string, limit int) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", nil
	}
	if limit <= 0 {
		limit = 5
	}

	resp, err := m.service.Recall(ctx, embedded.RecallInput{
		Connection: m.connection,
		Vault:      m.vaultForScope(scopeKey),
		Context:    []string{query},
		MaxResults: limit,
		MaxHops:    1,
		IncludeWhy: false,
	})
	if err != nil {
		return "", err
	}
	if resp == nil || len(resp.Activations) == 0 {
		return "", nil
	}

	var b strings.Builder
	b.WriteString("Relevant memories from past conversations:\n\n")
	for i, item := range resp.Activations {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Concept))
		if item.Content != "" {
			b.WriteString(fmt.Sprintf("   %s\n", truncateForMemory(item.Content, 240)))
		}
		b.WriteString(fmt.Sprintf("   (Relevance: %.2f)\n\n", item.Score))
	}
	return b.String(), nil
}

func truncateForMemory(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}

var nonVaultRune = regexp.MustCompile(`[^a-z0-9_-]+`)

func (m *MuninnMemoryBackend) vaultForScope(scopeKey string) string {
	scopeKey = strings.TrimSpace(scopeKey)
	if scopeKey == "" {
		return m.defaultVault
	}

	scoped := strings.ToLower(scopeKey)
	scoped = nonVaultRune.ReplaceAllString(scoped, "_")
	scoped = strings.Trim(scoped, "_")
	if scoped == "" {
		return m.defaultVault
	}
	if len(scoped) > 48 {
		scoped = scoped[:48]
	}
	return "user_" + scoped
}
