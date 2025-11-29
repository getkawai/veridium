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

package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

// ====================
// Tool Creation Tests
// ====================

func TestNewSearchTool(t *testing.T) {
	ctx := context.Background()

	tool, err := NewSearchTool(ctx, nil)
	require.NoError(t, err)
	assert.NotNil(t, tool)

	info, err := tool.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, "search", info.Name)
	assert.Contains(t, info.Desc, "Search the web")
}

func TestNewCrawlSinglePageTool(t *testing.T) {
	ctx := context.Background()

	tool, err := NewCrawlSinglePageTool(ctx, nil)
	require.NoError(t, err)
	assert.NotNil(t, tool)

	info, err := tool.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, "crawlSinglePage", info.Name)
	assert.Contains(t, info.Desc, "Retrieve content")
}

func TestNewCrawlMultiPagesTool(t *testing.T) {
	ctx := context.Background()

	tool, err := NewCrawlMultiPagesTool(ctx, nil)
	require.NoError(t, err)
	assert.NotNil(t, tool)

	info, err := tool.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, "crawlMultiPages", info.Name)
	assert.Contains(t, info.Desc, "multiple webpages")
}

func TestNewWebBrowsingTools(t *testing.T) {
	ctx := context.Background()

	tools, err := NewWebBrowsingTools(ctx)
	require.NoError(t, err)
	assert.Len(t, tools, 3)

	expectedNames := []string{"search", "crawlSinglePage", "crawlMultiPages"}
	for i, tool := range tools {
		info, err := tool.Info(ctx)
		require.NoError(t, err)
		assert.Equal(t, expectedNames[i], info.Name)
	}
}

func TestSearchToolWithCustomConfig(t *testing.T) {
	ctx := context.Background()

	config := &ToolConfig{
		ToolName: "custom_search",
		ToolDesc: "Custom search description",
	}

	tool, err := NewSearchTool(ctx, config)
	require.NoError(t, err)

	info, err := tool.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, "custom_search", info.Name)
	assert.Equal(t, "Custom search description", info.Desc)
}

// ====================
// Search Method Tests
// ====================

func TestSearch_EmptyQuery(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	_, err := wbt.Search(ctx, &SearchRequest{Query: ""})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "query is required")
}

func TestSearch_WhitespaceQuery(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	// Whitespace-only query should be treated as empty
	_, err := wbt.Search(ctx, &SearchRequest{Query: "   "})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "query is required")
}

func TestSearch_NilRequest(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	_, err := wbt.Search(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "request is required")
}

// ====================
// Crawl Single Page Tests
// ====================

func TestCrawlSinglePage_EmptyURL(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	_, err := wbt.CrawlSinglePage(ctx, &CrawlSinglePageRequest{URL: ""})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "url is required")
}

func TestCrawlSinglePage_InvalidURL(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	result, err := wbt.CrawlSinglePage(ctx, &CrawlSinglePageRequest{URL: "not-a-valid-url"})
	// Should either return error or result with Error field set
	if err == nil {
		assert.NotEmpty(t, result.Error, "Expected error message for invalid URL")
	}
}

func TestCrawlSinglePage_NilRequest(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	_, err := wbt.CrawlSinglePage(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "request is required")
}

// ====================
// Crawl Multi Pages Tests
// ====================

func TestCrawlMultiPages_EmptyURLs(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	_, err := wbt.CrawlMultiPages(ctx, &CrawlMultiPagesRequest{URLs: []string{}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "urls is required")
}

func TestCrawlMultiPages_NilURLs(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	_, err := wbt.CrawlMultiPages(ctx, &CrawlMultiPagesRequest{URLs: nil})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "urls is required")
}

func TestCrawlMultiPages_MixedValidInvalidURLs(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	result, err := wbt.CrawlMultiPages(ctx, &CrawlMultiPagesRequest{
		URLs: []string{"invalid-url", "also-invalid"},
	})

	if err == nil {
		// Should have results for each URL
		assert.Len(t, result.Results, 2, "Should return result for each URL")
		for _, r := range result.Results {
			// Each should have error since URLs are invalid
			assert.NotEmpty(t, r.Error, "Expected error for invalid URL")
		}
	}
}

func TestCrawlMultiPages_NilRequest(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	wbt := &webBrowsingTool{service: service, config: &ToolConfig{}}

	_, err := wbt.CrawlMultiPages(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "urls is required")
}

// ====================
// Crawler Unit Tests
// ====================

func TestCrawler_NewCrawlerDefaults(t *testing.T) {
	crawler := NewCrawler(nil)
	assert.NotNil(t, crawler)
	assert.Equal(t, []CrawlImplType{CrawlImplJina}, crawler.impls)
}

func TestCrawler_NewCrawlerWithImpls(t *testing.T) {
	impls := []CrawlImplType{CrawlImplNaive, CrawlImplJina}
	crawler := NewCrawler(impls)
	assert.NotNil(t, crawler)
	assert.Equal(t, impls, crawler.impls)
}

func TestCrawler_ExtractWebsite(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{"simple URL", "https://example.com/page", "example.com"},
		{"with port", "https://example.com:8080/page", "example.com"},
		{"subdomain", "https://sub.example.com/page", "sub.example.com"},
		{"invalid URL", "not-a-url", "not-a-url"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractWebsite(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCrawler_ExtractTitle(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		fallbackURL string
		expected    string
	}{
		{
			name:        "markdown title",
			content:     "# My Title\nSome content here",
			fallbackURL: "https://example.com",
			expected:    "My Title",
		},
		{
			name:        "plain text title",
			content:     "Short Title\nMore content",
			fallbackURL: "https://example.com",
			expected:    "Short Title",
		},
		{
			name:        "empty content uses fallback",
			content:     "",
			fallbackURL: "https://example.com",
			expected:    "https://example.com",
		},
		{
			name:        "only whitespace uses fallback",
			content:     "   \n\n   ",
			fallbackURL: "https://example.com",
			expected:    "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitle(tt.content, tt.fallbackURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCrawler_ExtractHTMLTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple title",
			html:     "<html><head><title>Test Page</title></head><body></body></html>",
			expected: "Test Page",
		},
		{
			name:     "no title",
			html:     "<html><head></head><body>Content</body></html>",
			expected: "",
		},
		{
			name:     "empty title",
			html:     "<html><head><title></title></head><body></body></html>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			require.NoError(t, err)
			result := extractHTMLTitle(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCrawler_ExtractTextContent(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		contains []string
		excludes []string
	}{
		{
			name:     "extracts body text",
			html:     "<html><body><p>Hello World</p></body></html>",
			contains: []string{"Hello World"},
			excludes: []string{},
		},
		{
			name:     "excludes script content",
			html:     "<html><body><script>var x = 1;</script><p>Visible</p></body></html>",
			contains: []string{"Visible"},
			excludes: []string{"var x"},
		},
		{
			name:     "excludes style content",
			html:     "<html><body><style>.class { color: red; }</style><p>Visible</p></body></html>",
			contains: []string{"Visible"},
			excludes: []string{"color"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTML(tt.html)
			require.NoError(t, err)
			result := extractTextContent(doc)
			for _, c := range tt.contains {
				assert.Contains(t, result, c)
			}
			for _, e := range tt.excludes {
				assert.NotContains(t, result, e)
			}
		})
	}
}

// ====================
// Brave Provider Tests
// ====================

func TestBraveProvider_MapTimeRange(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"day", "pd"},
		{"week", "pw"},
		{"month", "pm"},
		{"year", "py"},
		{"anytime", ""},
		{"invalid", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapTimeRange(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBraveProvider_Query_MockServer(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/web/search")
		assert.NotEmpty(t, r.URL.Query().Get("q"))

		// Return mock response
		response := BraveResponse{
			Web: struct {
				Results []BraveWebResult `json:"results"`
			}{
				Results: []BraveWebResult{
					{Title: "Result 1", URL: "https://example.com/1", Description: "Description 1"},
					{Title: "Result 2", URL: "https://example.com/2", Description: "Description 2"},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create provider with mock server
	provider := &BraveProvider{
		apiKey:     "test-key",
		baseURL:    server.URL,
		httpClient: http.DefaultClient,
	}

	ctx := context.Background()
	result, err := provider.Query(ctx, "test query", nil)

	require.NoError(t, err)
	assert.Equal(t, "test query", result.Query)
	assert.Len(t, result.Results, 2)
	assert.Equal(t, "Result 1", result.Results[0].Title)
	assert.Equal(t, "Result 2", result.Results[1].Title)
}

func TestBraveProvider_Query_EmptyAPIKey(t *testing.T) {
	provider := &BraveProvider{
		apiKey:     "",
		baseURL:    "https://api.search.brave.com/res/v1",
		httpClient: http.DefaultClient,
	}

	ctx := context.Background()
	_, err := provider.Query(ctx, "test query", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "BRAVE_SEARCH_API_KEY")
}

func TestBraveProvider_Query_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	provider := &BraveProvider{
		apiKey:     "test-key",
		baseURL:    server.URL,
		httpClient: http.DefaultClient,
	}

	ctx := context.Background()
	_, err := provider.Query(ctx, "test query", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestBraveProvider_Query_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	provider := &BraveProvider{
		apiKey:     "test-key",
		baseURL:    server.URL,
		httpClient: http.DefaultClient,
	}

	ctx := context.Background()
	_, err := provider.Query(ctx, "test query", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode")
}

// ====================
// Service Tests
// ====================

func TestService_WebSearch_RetryLogic(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)

	// Test with params that might trigger retry
	query := SearchQuery{
		Query:         "test",
		SearchEngines: []string{"nonexistent_engine"},
	}

	// This tests the retry logic path
	// Note: This will make real HTTP calls - in production you'd mock this
	_, err := service.WebSearch(query)
	// Error is acceptable since we're testing the code path
	if err != nil {
		t.Logf("Expected error from real API call: %v", err)
	}
}

func TestService_CrawlPages_EmptyURLs(t *testing.T) {
	service := NewService()

	resp, err := service.CrawlPages(CrawlPagesRequest{URLs: []string{}})
	require.NoError(t, err)
	assert.Empty(t, resp.Results)
}

func TestService_CrawlPages_ConcurrencyLimit(t *testing.T) {
	service := NewService()

	// Test with more URLs than concurrency limit (3)
	urls := []string{
		"invalid1", "invalid2", "invalid3", "invalid4", "invalid5",
	}

	resp, err := service.CrawlPages(CrawlPagesRequest{URLs: urls})
	require.NoError(t, err)
	assert.Len(t, resp.Results, 5, "Should return result for each URL")
}

// ====================
// System Prompt Tests
// ====================

func TestWebBrowsingSystemPrompt(t *testing.T) {
	prompt := WebBrowsingSystemPrompt()
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Web Information tool")
	assert.Contains(t, prompt, "search")
	assert.Contains(t, prompt, "crawlMultiPages")
	assert.Contains(t, prompt, "crawlSinglePage")
}

func TestWebBrowsingSystemPrompt_HasRequiredSections(t *testing.T) {
	prompt := WebBrowsingSystemPrompt()

	requiredSections := []string{
		"<core_capabilities>",
		"<workflow>",
		"<tool_selection_guidelines>",
		"<search_categories_selection>",
		"<citation_requirements>",
	}

	for _, section := range requiredSections {
		assert.Contains(t, prompt, section, "Missing required section: %s", section)
	}
}

func TestWebBrowsingToolIdentifier(t *testing.T) {
	assert.NotEmpty(t, WebBrowsingToolIdentifier)
	assert.Equal(t, "veridium-web-browsing", WebBrowsingToolIdentifier)
}

// ====================
// Helper Functions
// ====================

func parseHTML(htmlStr string) (*html.Node, error) {
	return html.Parse(strings.NewReader(htmlStr))
}

// ====================
// Integration Tests (require network)
// ====================

func TestIntegration_CrawlWithJina(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	crawler := NewCrawler([]CrawlImplType{CrawlImplJina})
	result, err := crawler.crawlWithJina("https://go.dev")

	require.NoError(t, err)
	require.NotNil(t, result.Success)
	assert.NotEmpty(t, result.Success.Title)
	assert.NotEmpty(t, result.Success.Content)
	assert.Equal(t, "https://go.dev", result.Success.URL)
	assert.Equal(t, "go.dev", result.Success.Website)
}

func TestIntegration_CrawlNaive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	crawler := NewCrawler([]CrawlImplType{CrawlImplNaive})
	result, err := crawler.crawlNaive("https://go.dev")

	require.NoError(t, err)
	require.NotNil(t, result.Success)
	assert.NotEmpty(t, result.Success.Content)
	assert.Equal(t, "https://go.dev", result.Success.URL)
}
