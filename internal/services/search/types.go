package search

import "github.com/kawai-network/veridium/internal/services/search/providers/types"

// Re-export provider types for convenience
type SearchParams = types.SearchParams
type UniformSearchResult = types.UniformSearchResult
type UniformSearchResponse = types.UniformSearchResponse

// SearchQuery combines query string with search parameters
type SearchQuery struct {
	Query            string   `json:"query"`
	SearchCategories []string `json:"searchCategories,omitempty"`
	SearchEngines    []string `json:"searchEngines,omitempty"`
	SearchTimeRange  string   `json:"searchTimeRange,omitempty"`
}

// CrawlImplType represents the type of crawler implementation
type CrawlImplType string

const (
	CrawlImplJina        CrawlImplType = "jina"
	CrawlImplNaive       CrawlImplType = "naive"
	CrawlImplBrowserless CrawlImplType = "browserless"
)

// CrawlSuccessResult represents a successful crawl result
type CrawlSuccessResult struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
	Website string `json:"website"`
}

// CrawlErrorResult represents a failed crawl result
type CrawlErrorResult struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorType    string `json:"errorType"`
	URL          string `json:"url"`
}

// CrawlResult is a union type for crawl results
type CrawlResult struct {
	Success *CrawlSuccessResult `json:"success,omitempty"`
	Error   *CrawlErrorResult   `json:"error,omitempty"`
}

// CrawlPagesRequest represents a request to crawl multiple pages
type CrawlPagesRequest struct {
	URLs  []string        `json:"urls"`
	Impls []CrawlImplType `json:"impls,omitempty"`
}

// CrawlPagesResponse represents the response from crawling multiple pages
type CrawlPagesResponse struct {
	Results []CrawlResult `json:"results"`
}

