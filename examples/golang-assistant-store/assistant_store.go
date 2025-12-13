package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// Default base URL for agents index
	DefaultAgentsIndexURL = "https://registry.npmmirror.com/@lobehub/agents-index/v1/files/public"
	
	// Default locale
	DefaultLocale = "en-US"
	
	// Cache revalidation time
	CacheRevalidateList    = 3600 * time.Second // 1 hour
	CacheRevalidateDetails = 86400 * time.Second // 24 hours
)

// AssistantStore handles fetching assistant data from NPM registry
type AssistantStore struct {
	baseURL    string
	httpClient *http.Client
	cache      *Cache
}

// AgentMeta contains metadata about an agent
type AgentMeta struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Avatar      string   `json:"avatar"`
}

// AgentIndexItem represents an agent in the index list
type AgentIndexItem struct {
	Identifier string    `json:"identifier"`
	Category   string    `json:"category"`
	Author     string    `json:"author"`
	Meta       AgentMeta `json:"meta"`
	CreatedAt  string    `json:"createdAt,omitempty"`
	Homepage   string    `json:"homepage,omitempty"`
}

// AgentIndexResponse is the response from index.{locale}.json
type AgentIndexResponse struct {
	Agents []AgentIndexItem `json:"agents"`
}

// AgentDetail contains full agent information
type AgentDetail struct {
	Identifier string                 `json:"identifier"`
	Author     string                 `json:"author"`
	SystemRole string                 `json:"systemRole"`
	Meta       AgentMeta              `json:"meta"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Plugins    []string               `json:"plugins,omitempty"`
	CreatedAt  string                 `json:"createdAt,omitempty"`
	Homepage   string                 `json:"homepage,omitempty"`
}

// FilterOptions for filtering agents
type FilterOptions struct {
	Whitelist []string
	Blacklist []string
}

// NewAssistantStore creates a new AssistantStore instance
func NewAssistantStore(baseURL string) *AssistantStore {
	if baseURL == "" {
		baseURL = DefaultAgentsIndexURL
	}

	return &AssistantStore{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: NewCache(),
	}
}

// GetAgentIndexURL returns the URL for agent index based on locale
func (s *AssistantStore) GetAgentIndexURL(locale string) string {
	if locale == "" {
		locale = DefaultLocale
	}
	
	normalizedLocale := normalizeLocale(locale)
	return fmt.Sprintf("%s/index.%s.json", s.baseURL, normalizedLocale)
}

// GetAgentURL returns the URL for a specific agent based on identifier and locale
func (s *AssistantStore) GetAgentURL(identifier, locale string) string {
	if locale == "" {
		locale = DefaultLocale
	}
	
	normalizedLocale := normalizeLocale(locale)
	return fmt.Sprintf("%s/%s.%s.json", s.baseURL, identifier, normalizedLocale)
}

// GetAgentIndex fetches the list of all agents for a given locale
func (s *AssistantStore) GetAgentIndex(locale string) ([]AgentIndexItem, error) {
	if locale == "" {
		locale = DefaultLocale
	}

	// Check cache first
	cacheKey := fmt.Sprintf("index:%s", locale)
	if cached, found := s.cache.Get(cacheKey); found {
		if agents, ok := cached.([]AgentIndexItem); ok {
			return agents, nil
		}
	}

	url := s.GetAgentIndexURL(locale)
	
	// Try to fetch with specified locale
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent index: %w", err)
	}
	defer resp.Body.Close()

	// If 404, fallback to default locale
	if resp.StatusCode == http.StatusNotFound && locale != DefaultLocale {
		url = s.GetAgentIndexURL(DefaultLocale)
		resp, err = s.httpClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch agent index (fallback): %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch agent index error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var indexResponse AgentIndexResponse
	if err := json.Unmarshal(body, &indexResponse); err != nil {
		return nil, fmt.Errorf("failed to parse agent index: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, indexResponse.Agents, CacheRevalidateList)

	return indexResponse.Agents, nil
}

// GetAgentIndexWithFilter fetches agents with whitelist/blacklist filtering
func (s *AssistantStore) GetAgentIndexWithFilter(locale string, filter *FilterOptions) ([]AgentIndexItem, error) {
	agents, err := s.GetAgentIndex(locale)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		return agents, nil
	}

	// Apply whitelist first (if provided)
	if len(filter.Whitelist) > 0 {
		filtered := make([]AgentIndexItem, 0)
		whitelistMap := make(map[string]bool)
		for _, id := range filter.Whitelist {
			whitelistMap[id] = true
		}

		for _, agent := range agents {
			if whitelistMap[agent.Identifier] {
				filtered = append(filtered, agent)
			}
		}
		return filtered, nil
	}

	// Apply blacklist (if no whitelist)
	if len(filter.Blacklist) > 0 {
		filtered := make([]AgentIndexItem, 0)
		blacklistMap := make(map[string]bool)
		for _, id := range filter.Blacklist {
			blacklistMap[id] = true
		}

		for _, agent := range agents {
			if !blacklistMap[agent.Identifier] {
				filtered = append(filtered, agent)
			}
		}
		return filtered, nil
	}

	return agents, nil
}

// GetAgent fetches detailed information for a specific agent
func (s *AssistantStore) GetAgent(identifier, locale string) (*AgentDetail, error) {
	if locale == "" {
		locale = DefaultLocale
	}

	// Check cache first
	cacheKey := fmt.Sprintf("agent:%s:%s", identifier, locale)
	if cached, found := s.cache.Get(cacheKey); found {
		if detail, ok := cached.(*AgentDetail); ok {
			return detail, nil
		}
	}

	url := s.GetAgentURL(identifier, locale)
	
	// Try to fetch with specified locale
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent detail: %w", err)
	}
	defer resp.Body.Close()

	// If 404, fallback to default locale
	if resp.StatusCode == http.StatusNotFound && locale != DefaultLocale {
		url = s.GetAgentURL(identifier, DefaultLocale)
		resp, err = s.httpClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch agent detail (fallback): %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch agent detail error: status=%d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var detail AgentDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("failed to parse agent detail: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, &detail, CacheRevalidateDetails)

	return &detail, nil
}

// SearchAgents searches agents by query string (searches in title, description, tags, author)
func (s *AssistantStore) SearchAgents(locale, query string) ([]AgentIndexItem, error) {
	agents, err := s.GetAgentIndex(locale)
	if err != nil {
		return nil, err
	}

	if query == "" {
		return agents, nil
	}

	queryLower := strings.ToLower(query)
	filtered := make([]AgentIndexItem, 0)

	for _, agent := range agents {
		// Search in multiple fields
		searchText := strings.ToLower(fmt.Sprintf("%s %s %s %s %s",
			agent.Author,
			agent.Meta.Title,
			agent.Meta.Description,
			strings.Join(agent.Meta.Tags, " "),
			agent.Category,
		))

		if strings.Contains(searchText, queryLower) {
			filtered = append(filtered, agent)
		}
	}

	return filtered, nil
}

// FilterByCategory filters agents by category
func (s *AssistantStore) FilterByCategory(locale, category string) ([]AgentIndexItem, error) {
	agents, err := s.GetAgentIndex(locale)
	if err != nil {
		return nil, err
	}

	if category == "" {
		return agents, nil
	}

	filtered := make([]AgentIndexItem, 0)
	for _, agent := range agents {
		if agent.Category == category {
			filtered = append(filtered, agent)
		}
	}

	return filtered, nil
}

// GetCategories returns all unique categories with their counts
func (s *AssistantStore) GetCategories(locale string) (map[string]int, error) {
	agents, err := s.GetAgentIndex(locale)
	if err != nil {
		return nil, err
	}

	categories := make(map[string]int)
	for _, agent := range agents {
		if agent.Category != "" {
			categories[agent.Category]++
		}
	}

	return categories, nil
}

// normalizeLocale converts locale to the format used in file names
func normalizeLocale(locale string) string {
	// Convert locale format (e.g., "en-US", "zh-CN")
	// Already in correct format for most cases
	return locale
}

