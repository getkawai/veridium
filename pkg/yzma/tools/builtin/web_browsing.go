package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kawai-network/veridium/internal/search"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// WebBrowsingService wraps search.Service to provide lobe-web-browsing compatible responses
type WebBrowsingService struct {
	searchService *search.Service
}

// NewWebBrowsingService creates a new web browsing service
func NewWebBrowsingService() *WebBrowsingService {
	return &WebBrowsingService{
		searchService: search.NewService(),
	}
}

// ============================================================================
// Response Types (matching frontend expected format)
// ============================================================================

// Note: UniformSearchResult and UniformSearchResponse are imported from internal/search
// to avoid duplication. They match the frontend UniformSearchResult interface.

// CrawlData matches frontend expected crawl data format
type CrawlData struct {
	Content     string `json:"content"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// CrawlResult matches frontend CrawlResult interface
type CrawlResult struct {
	OriginalUrl string    `json:"originalUrl"`
	Crawler     string    `json:"crawler"`
	Data        CrawlData `json:"data"`
}

// CrawlPluginState matches frontend CrawlPluginState interface
type CrawlPluginState struct {
	Results []CrawlResult `json:"results"`
}

// ============================================================================
// Service Methods
// ============================================================================

// Search performs web search and returns frontend-compatible response
func (s *WebBrowsingService) Search(query string, categories, engines []string, timeRange string) (*search.UniformSearchResponse, error) {
	startTime := time.Now()

	searchQuery := search.SearchQuery{
		Query:            query,
		SearchCategories: categories,
		SearchEngines:    engines,
		SearchTimeRange:  timeRange,
	}

	response, err := s.searchService.WebSearch(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// No transformation needed - response is already in the correct format
	// Just update the cost time
	response.CostTime = time.Since(startTime).Milliseconds()

	return response, nil
}

// CrawlSinglePage crawls a single URL and returns frontend-compatible response
func (s *WebBrowsingService) CrawlSinglePage(url string) (*CrawlPluginState, error) {
	return s.CrawlMultiPages([]string{url})
}

// CrawlMultiPages crawls multiple URLs and returns frontend-compatible response
func (s *WebBrowsingService) CrawlMultiPages(urls []string) (*CrawlPluginState, error) {
	response, err := s.searchService.CrawlPages(search.CrawlPagesRequest{
		URLs: urls,
	})
	if err != nil {
		return nil, fmt.Errorf("crawl failed: %w", err)
	}

	// Transform to frontend format
	results := make([]CrawlResult, 0, len(response.Results))
	for i, r := range response.Results {
		originalUrl := ""
		if i < len(urls) {
			originalUrl = urls[i]
		}

		if r.Error != nil {
			// Include error in data
			results = append(results, CrawlResult{
				OriginalUrl: originalUrl,
				Crawler:     "jina",
				Data: CrawlData{
					Content:     fmt.Sprintf("Error: %s", r.Error.ErrorMessage),
					URL:         r.Error.URL,
					Title:       "Crawl Failed",
					Description: r.Error.ErrorType,
				},
			})
		} else if r.Success != nil {
			results = append(results, CrawlResult{
				OriginalUrl: originalUrl,
				Crawler:     "jina",
				Data: CrawlData{
					Content:     r.Success.Content,
					URL:         r.Success.URL,
					Title:       r.Success.Title,
					Description: "",
				},
			})
		}
	}

	return &CrawlPluginState{
		Results: results,
	}, nil
}

// ============================================================================
// Tool Registration
// ============================================================================

// RegisterWebBrowsing registers the lobe-web-browsing tools (search, crawlSinglePage, crawlMultiPages)
func RegisterWebBrowsing(registry *tools.ToolRegistry) error {
	service := NewWebBrowsingService()

	// Tool 1: search
	searchTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-web-browsing__search",
			Description: "Search the web for information. Returns a list of search results with title, content, and URL.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query string",
					},
					"searchCategories": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Search categories: general, images, news, science, videos",
					},
					"searchEngines": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Search engines: google, bing, duckduckgo, brave, wikipedia, github, arxiv",
					},
					"searchTimeRange": map[string]interface{}{
						"type":        "string",
						"description": "Time range filter: anytime, day, week, month, year",
					},
				},
				"required": []string{"query"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			query := args["query"]
			if query == "" {
				return "", fmt.Errorf("query is required")
			}

			var categories, engines []string
			timeRange := "anytime"

			if catStr := args["searchCategories"]; catStr != "" {
				json.Unmarshal([]byte(catStr), &categories)
			}
			if engStr := args["searchEngines"]; engStr != "" {
				json.Unmarshal([]byte(engStr), &engines)
			}
			if tr := args["searchTimeRange"]; tr != "" {
				timeRange = tr
			}

			response, err := service.Search(query, categories, engines, timeRange)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(response)
			if err != nil {
				return "", fmt.Errorf("failed to marshal response: %w", err)
			}

			log.Printf("🔍 Web search completed: query=%s, results=%d", query, len(response.Results))
			return string(resultJSON), nil
		},
		Enabled: true,
	}

	if err := registry.Register(searchTool); err != nil {
		return fmt.Errorf("failed to register search tool: %w", err)
	}

	// Tool 2: crawlSinglePage
	crawlSingleTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-web-browsing__crawlSinglePage",
			Description: "Retrieve content from a specific webpage. Returns the page title, content, URL and website.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL of the webpage to crawl",
					},
				},
				"required": []string{"url"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			url := args["url"]
			if url == "" {
				return "", fmt.Errorf("url is required")
			}

			response, err := service.CrawlSinglePage(url)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(response)
			if err != nil {
				return "", fmt.Errorf("failed to marshal response: %w", err)
			}

			log.Printf("🌐 Crawled single page: url=%s", url)
			return string(resultJSON), nil
		},
		Enabled: true,
	}

	if err := registry.Register(crawlSingleTool); err != nil {
		return fmt.Errorf("failed to register crawlSinglePage tool: %w", err)
	}

	// Tool 3: crawlMultiPages
	crawlMultiTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-web-browsing__crawlMultiPages",
			Description: "Retrieve content from multiple webpages simultaneously. Returns an array of page results.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"urls": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "The URLs of the webpages to crawl",
					},
				},
				"required": []string{"urls"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			urlsStr := args["urls"]
			if urlsStr == "" {
				return "", fmt.Errorf("urls is required")
			}

			var urls []string
			if err := json.Unmarshal([]byte(urlsStr), &urls); err != nil {
				return "", fmt.Errorf("failed to parse urls: %w", err)
			}

			if len(urls) == 0 {
				return "", fmt.Errorf("at least one URL is required")
			}

			response, err := service.CrawlMultiPages(urls)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(response)
			if err != nil {
				return "", fmt.Errorf("failed to marshal response: %w", err)
			}

			log.Printf("🌐 Crawled %d pages", len(urls))
			return string(resultJSON), nil
		},
		Enabled: true,
	}

	if err := registry.Register(crawlMultiTool); err != nil {
		return fmt.Errorf("failed to register crawlMultiPages tool: %w", err)
	}

	return nil
}
