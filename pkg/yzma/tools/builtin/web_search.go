package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/internal/search"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
)

// RegisterWebSearch registers the web search tool
func RegisterWebSearch(registry *tools.ToolRegistry) error {
	searchService := search.NewService()

	tool := &types.Tool{
		Type:     fantasy.ToolTypeFunction,
		Parallel: true, // Safe to run in parallel - read-only external API call
		Definition: types.ToolDefinition{
			Name:        "web_search",
			Description: "Search the web for current information using Brave Search. Returns real-time search results with titles, URLs, and descriptions.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
					"max_results": map[string]interface{}{
						"type":        "number",
						"description": "Maximum number of results (default: 10)",
					},
				},
				"required": []string{"query"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			query, ok := args["query"]
			if !ok || query == "" {
				return "", fmt.Errorf("query parameter is required")
			}

			maxResults := 10
			if maxStr, ok := args["max_results"]; ok && maxStr != "" {
				// Parse max_results if provided
				var mr int
				if _, err := fmt.Sscanf(maxStr, "%d", &mr); err == nil && mr > 0 {
					maxResults = mr
				}
			}

			// Use real Brave Search API
			searchQuery := search.SearchQuery{
				Query:            query,
				SearchCategories: []string{"general"},
				SearchEngines:    []string{},
				SearchTimeRange:  "anytime",
			}

			response, err := searchService.WebSearch(searchQuery)
			if err != nil {
				log.Printf("⚠️  Web search failed: %v", err)
				return "", fmt.Errorf("search failed: %w", err)
			}

			// Format results for LLM
			results := make([]map[string]interface{}, 0, len(response.Results))
			for i, result := range response.Results {
				if i >= maxResults {
					break
				}
				results = append(results, map[string]interface{}{
					"title":   result.Title,
					"url":     result.URL,
					"snippet": result.Content,
				})
			}

			resultData := map[string]interface{}{
				"query":       query,
				"results":     results,
				"count":       len(results),
				"max_results": maxResults,
			}

			resultJSON, err := json.Marshal(resultData)
			if err != nil {
				return "", fmt.Errorf("failed to marshal results: %w", err)
			}

			return string(resultJSON), nil
		},
		Enabled: true,
	}

	return registry.Register(tool)
}
