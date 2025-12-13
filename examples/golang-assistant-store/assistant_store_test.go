package main

import (
	"testing"
	"time"
)

func TestNewAssistantStore(t *testing.T) {
	// Test with default URL
	store := NewAssistantStore("")
	if store.baseURL != DefaultAgentsIndexURL {
		t.Errorf("Expected default URL %s, got %s", DefaultAgentsIndexURL, store.baseURL)
	}

	// Test with custom URL
	customURL := "https://custom.com/agents"
	store = NewAssistantStore(customURL)
	if store.baseURL != customURL {
		t.Errorf("Expected custom URL %s, got %s", customURL, store.baseURL)
	}
}

func TestGetAgentIndexURL(t *testing.T) {
	store := NewAssistantStore("")

	tests := []struct {
		locale   string
		expected string
	}{
		{"en-US", DefaultAgentsIndexURL + "/index.en-US.json"},
		{"zh-CN", DefaultAgentsIndexURL + "/index.zh-CN.json"},
		{"", DefaultAgentsIndexURL + "/index.en-US.json"},
	}

	for _, tt := range tests {
		result := store.GetAgentIndexURL(tt.locale)
		if result != tt.expected {
			t.Errorf("GetAgentIndexURL(%s) = %s, want %s", tt.locale, result, tt.expected)
		}
	}
}

func TestGetAgentURL(t *testing.T) {
	store := NewAssistantStore("")

	tests := []struct {
		identifier string
		locale     string
		expected   string
	}{
		{"web-dev", "en-US", DefaultAgentsIndexURL + "/web-dev.en-US.json"},
		{"api-design", "zh-CN", DefaultAgentsIndexURL + "/api-design.zh-CN.json"},
		{"test", "", DefaultAgentsIndexURL + "/test.en-US.json"},
	}

	for _, tt := range tests {
		result := store.GetAgentURL(tt.identifier, tt.locale)
		if result != tt.expected {
			t.Errorf("GetAgentURL(%s, %s) = %s, want %s",
				tt.identifier, tt.locale, result, tt.expected)
		}
	}
}

func TestCache(t *testing.T) {
	cache := NewCache()

	// Test Set and Get
	key := "test-key"
	value := "test-value"
	cache.Set(key, value, 1*time.Second)

	retrieved, found := cache.Get(key)
	if !found {
		t.Error("Expected to find cached value")
	}
	if retrieved != value {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}

	// Test expiration
	time.Sleep(2 * time.Second)
	_, found = cache.Get(key)
	if found {
		t.Error("Expected cache to be expired")
	}

	// Test Delete
	cache.Set(key, value, 10*time.Second)
	cache.Delete(key)
	_, found = cache.Get(key)
	if found {
		t.Error("Expected cache to be deleted")
	}

	// Test Clear
	cache.Set("key1", "value1", 10*time.Second)
	cache.Set("key2", "value2", 10*time.Second)
	cache.Clear()
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	if found1 || found2 {
		t.Error("Expected all cache to be cleared")
	}
}

func TestFilterOptions(t *testing.T) {
	// Mock data
	mockAgents := []AgentIndexItem{
		{Identifier: "agent1", Meta: AgentMeta{Title: "Agent 1"}},
		{Identifier: "agent2", Meta: AgentMeta{Title: "Agent 2"}},
		{Identifier: "agent3", Meta: AgentMeta{Title: "Agent 3"}},
	}

	// Test whitelist
	filter := &FilterOptions{
		Whitelist: []string{"agent1", "agent3"},
	}

	// Note: This is a conceptual test - actual implementation would need mock HTTP responses
	t.Log("Filter options created successfully")

	if len(filter.Whitelist) != 2 {
		t.Errorf("Expected whitelist length 2, got %d", len(filter.Whitelist))
	}

	// Test blacklist
	filter = &FilterOptions{
		Blacklist: []string{"agent2"},
	}

	if len(filter.Blacklist) != 1 {
		t.Errorf("Expected blacklist length 1, got %d", len(filter.Blacklist))
	}

	t.Logf("Mock agents: %d", len(mockAgents))
}

func TestNormalizeLocale(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"en-US", "en-US"},
		{"zh-CN", "zh-CN"},
		{"ja-JP", "ja-JP"},
	}

	for _, tt := range tests {
		result := normalizeLocale(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeLocale(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

// Benchmark tests
func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", 1*time.Hour)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache()
	cache.Set("key", "value", 1*time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}
