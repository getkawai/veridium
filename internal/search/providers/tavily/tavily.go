package tavily

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kawai-network/veridium/internal/search/providers/types"
)

// Provider implements the Tavily search provider
type Provider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// TavilyResult represents a single search result from Tavily
type TavilyResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

// TavilyRequest represents the request to Tavily Search API
type TavilyRequest struct {
	Query                    string `json:"query"`
	SearchDepth              string `json:"search_depth,omitempty"`
	IncludeAnswer            bool   `json:"include_answer"`
	IncludeImages            bool   `json:"include_images"`
	IncludeImageDescriptions bool   `json:"include_image_descriptions"`
	IncludeRawContent        bool   `json:"include_raw_content"`
	MaxResults               int    `json:"max_results"`
	Topic                    string `json:"topic,omitempty"`
	TimeRange                string `json:"time_range,omitempty"`
}

// TavilyResponse represents the response from Tavily Search API
type TavilyResponse struct {
	Results []TavilyResult `json:"results"`
}

// NewProvider creates a new Tavily search provider
func NewProvider() *Provider {
	apiKey := os.Getenv("TAVILY_API_KEY")
	searchDepth := os.Getenv("TAVILY_SEARCH_DEPTH")
	if searchDepth == "" {
		searchDepth = "basic"
	}

	return &Provider{
		apiKey:  apiKey,
		baseURL: "https://api.tavily.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "tavily"
}

// Query performs a search query using Tavily Search API
func (p *Provider) Query(ctx context.Context, query string, params *types.SearchParams) (*types.UniformSearchResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEY environment variable not set")
	}

	endpoint := fmt.Sprintf("%s/search", p.baseURL)

	// Build request body
	req := TavilyRequest{
		Query:                    query,
		SearchDepth:              os.Getenv("TAVILY_SEARCH_DEPTH"),
		IncludeAnswer:            false,
		IncludeImages:            false,
		IncludeImageDescriptions: true,
		IncludeRawContent:        false,
		MaxResults:               15,
	}

	if req.SearchDepth == "" {
		req.SearchDepth = "basic"
	}

	// Add time range if specified
	if params != nil && params.SearchTimeRange != "" && params.SearchTimeRange != "anytime" {
		req.TimeRange = params.SearchTimeRange
	}

	// Add topic (Tavily only supports 'news' and 'general')
	if params != nil && len(params.SearchCategories) > 0 {
		for _, cat := range params.SearchCategories {
			if cat == "news" || cat == "general" {
				req.Topic = cat
				break
			}
		}
	}
	if req.Topic == "" {
		req.Topic = "general"
	}

	// Marshal request body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	// Execute request
	startTime := time.Now()
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	costTime := time.Since(startTime).Milliseconds()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tavily API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tavilyResp TavilyResponse
	if err := json.NewDecoder(resp.Body).Decode(&tavilyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to uniform format
	results := make([]types.UniformSearchResult, 0, len(tavilyResp.Results))
	for _, result := range tavilyResp.Results {
		parsedURL := ""
		if u, err := url.Parse(result.URL); err == nil {
			parsedURL = u.Hostname()
		}

		results = append(results, types.UniformSearchResult{
			Category:  req.Topic,
			Content:   result.Content,
			Engines:   []string{"tavily"},
			ParsedUrl: parsedURL,
			Score:     result.Score,
			Title:     result.Title,
			URL:       result.URL,
		})
	}

	return &types.UniformSearchResponse{
		CostTime:      costTime,
		Query:         query,
		ResultNumbers: len(results),
		Results:       results,
	}, nil
}
