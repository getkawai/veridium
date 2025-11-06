package search

import (
	"context"

	"github.com/kawai-network/veridium/internal/search/providers/types"
)

// Provider defines the interface that all search providers must implement
type Provider interface {
	// Query performs a search query and returns results
	Query(ctx context.Context, query string, params *types.SearchParams) (*types.UniformSearchResponse, error)

	// Name returns the provider name
	Name() string
}

// ProviderType represents the type of search provider
type ProviderType string

const (
	ProviderBrave   ProviderType = "brave"
	ProviderTavily  ProviderType = "tavily"
	ProviderSearXNG ProviderType = "searxng"
)
