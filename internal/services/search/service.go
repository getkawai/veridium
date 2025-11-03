package search

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/kawai-network/veridium/internal/services/search/providers/brave"
	"github.com/kawai-network/veridium/internal/services/search/providers/searxng"
	"github.com/kawai-network/veridium/internal/services/search/providers/tavily"
)

// Service provides search and crawl functionality
type Service struct {
	provider Provider
	crawler  *Crawler
	mu       sync.RWMutex
}

// NewService creates a new search service
func NewService() *Service {
	// Get provider from environment or use default
	providerType := getProviderFromEnv()
	provider := createProvider(providerType)

	// Get crawler implementations from environment
	crawlerImpls := getCrawlerImplsFromEnv()
	crawler := NewCrawler(crawlerImpls)

	return &Service{
		provider: provider,
		crawler:  crawler,
	}
}

// Query performs a search query using the configured provider
func (s *Service) Query(query string, params *SearchParams) (*UniformSearchResponse, error) {
	ctx := context.Background()
	return s.provider.Query(ctx, query, params)
}

// WebSearch performs a web search with retry logic
func (s *Service) WebSearch(query SearchQuery) (*UniformSearchResponse, error) {
	params := &SearchParams{
		SearchCategories: query.SearchCategories,
		SearchEngines:    query.SearchEngines,
		SearchTimeRange:  query.SearchTimeRange,
	}

	// First attempt with all parameters
	data, err := s.Query(query.Query, params)
	if err != nil {
		return nil, err
	}

	// First retry: remove search engine restrictions if no results found
	if len(data.Results) == 0 && len(query.SearchEngines) > 0 {
		paramsExcludeSearchEngines := &SearchParams{
			SearchCategories: query.SearchCategories,
			SearchEngines:    nil,
			SearchTimeRange:  query.SearchTimeRange,
		}
		data, err = s.Query(query.Query, paramsExcludeSearchEngines)
		if err != nil {
			return nil, err
		}
	}

	// Second retry: remove all restrictions if still no results found
	if len(data.Results) == 0 {
		data, err = s.Query(query.Query, nil)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// CrawlPages crawls multiple URLs concurrently
func (s *Service) CrawlPages(req CrawlPagesRequest) (*CrawlPagesResponse, error) {
	results := s.crawler.CrawlPages(req.URLs, req.Impls)
	return &CrawlPagesResponse{Results: results}, nil
}

// getProviderFromEnv reads the search provider from environment variables
func getProviderFromEnv() ProviderType {
	envStr := os.Getenv("SEARCH_PROVIDERS")
	if envStr == "" {
		return ProviderSearXNG // default
	}

	// Parse comma-separated list and get first provider
	providers := strings.Split(strings.ReplaceAll(envStr, "，", ","), ",")
	if len(providers) > 0 {
		switch strings.TrimSpace(providers[0]) {
		case "brave":
			return ProviderBrave
		case "tavily":
			return ProviderTavily
		case "searxng":
			return ProviderSearXNG
		}
	}

	return ProviderSearXNG
}

// getCrawlerImplsFromEnv reads crawler implementations from environment
func getCrawlerImplsFromEnv() []CrawlImplType {
	envStr := os.Getenv("CRAWLER_IMPLS")
	if envStr == "" {
		return []CrawlImplType{CrawlImplJina} // default to Jina
	}

	// Parse comma-separated list
	implStrs := strings.Split(strings.ReplaceAll(envStr, "，", ","), ",")
	var impls []CrawlImplType
	for _, s := range implStrs {
		s = strings.TrimSpace(s)
		if s != "" {
			impls = append(impls, CrawlImplType(s))
		}
	}

	if len(impls) == 0 {
		return []CrawlImplType{CrawlImplJina}
	}

	return impls
}

// createProvider creates a provider instance based on the type
func createProvider(providerType ProviderType) Provider {
	switch providerType {
	case ProviderBrave:
		return brave.NewProvider()
	case ProviderTavily:
		return tavily.NewProvider()
	case ProviderSearXNG:
		return searxng.NewProvider()
	default:
		// Fallback to SearXNG
		fmt.Printf("Unknown provider type: %s, falling back to SearXNG\n", providerType)
		return searxng.NewProvider()
	}
}
