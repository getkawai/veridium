package searxng

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kawai-network/veridium/internal/search/providers/types"
)

// Provider implements the SearXNG search provider
type Provider struct {
	baseURL    string
	httpClient *http.Client
}

// SearXNGResult represents a single search result from SearXNG
type SearXNGResult struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	Content       string   `json:"content"`
	Engines       []string `json:"engines"`
	Category      string   `json:"category"`
	PublishedDate string   `json:"publisheddate,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	ImgSrc        string   `json:"img_src,omitempty"`
	IframeSrc     string   `json:"iframe_src,omitempty"`
}

// SearXNGResponse represents the response from SearXNG API
type SearXNGResponse struct {
	Results []SearXNGResult `json:"results"`
	Query   string          `json:"query"`
}

// NewProvider creates a new SearXNG search provider
func NewProvider() *Provider {
	baseURL := os.Getenv("SEARXNG_BASE_URL")
	if baseURL == "" {
		baseURL = "https://searx.be" // default public instance
	}

	return &Provider{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "searxng"
}

// Query performs a search query using SearXNG API
func (p *Provider) Query(ctx context.Context, query string, params *types.SearchParams) (*types.UniformSearchResponse, error) {
	endpoint := fmt.Sprintf("%s/search", p.baseURL)

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("q", query)
	queryParams.Set("format", "json")
	queryParams.Set("pageno", "1")

	// Add categories if specified
	if params != nil && len(params.SearchCategories) > 0 {
		queryParams.Set("categories", strings.Join(params.SearchCategories, ","))
	}

	// Add time range if specified
	if params != nil && params.SearchTimeRange != "" && params.SearchTimeRange != "anytime" {
		queryParams.Set("time_range", params.SearchTimeRange)
	}

	// Add search engines if specified
	if params != nil && len(params.SearchEngines) > 0 {
		queryParams.Set("engines", strings.Join(params.SearchEngines, ","))
	}

	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	// Execute request
	startTime := time.Now()
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	costTime := time.Since(startTime).Milliseconds()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SearXNG API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var searxngResp SearXNGResponse
	if err := json.NewDecoder(resp.Body).Decode(&searxngResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to uniform format
	results := make([]types.UniformSearchResult, 0, len(searxngResp.Results))
	for _, result := range searxngResp.Results {
		parsedURL := ""
		if u, err := url.Parse(result.URL); err == nil {
			parsedURL = u.Hostname()
		}

		results = append(results, types.UniformSearchResult{
			Category:      result.Category,
			Content:       result.Content,
			Engines:       result.Engines,
			ParsedUrl:     parsedURL,
			Score:         1.0, // SearXNG doesn't provide scores
			Title:         result.Title,
			URL:           result.URL,
			PublishedDate: result.PublishedDate,
			Thumbnail:     result.Thumbnail,
			ImgSrc:        result.ImgSrc,
			IframeSrc:     result.IframeSrc,
		})
	}

	return &types.UniformSearchResponse{
		CostTime:      costTime,
		Query:         query,
		ResultNumbers: len(results),
		Results:       results,
	}, nil
}
