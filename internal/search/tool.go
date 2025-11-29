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
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// ToolConfig represents the web browsing tool configuration
type ToolConfig struct {
	ToolName string `json:"tool_name"` // default: web_search
	ToolDesc string `json:"tool_desc"` // default: "Search the web for information"
}

// webBrowsingTool wraps the Service to provide Eino tool interface
type webBrowsingTool struct {
	service *Service
	config  *ToolConfig
}

// NewSearchTool creates a new web search tool for Eino agent
func NewSearchTool(ctx context.Context, config *ToolConfig) (tool.InvokableTool, error) {
	if config == nil {
		config = &ToolConfig{}
	}

	if config.ToolName == "" {
		config.ToolName = "search"
	}
	if config.ToolDesc == "" {
		config.ToolDesc = "Search the web for information. Returns a list of search results with title, content, and URL."
	}

	service := NewService()
	wbt := &webBrowsingTool{
		service: service,
		config:  config,
	}

	searchTool, err := utils.InferTool(config.ToolName, config.ToolDesc, wbt.Search)
	if err != nil {
		return nil, fmt.Errorf("failed to infer search tool: %w", err)
	}

	return searchTool, nil
}

// NewCrawlSinglePageTool creates a tool for crawling a single web page
func NewCrawlSinglePageTool(ctx context.Context, config *ToolConfig) (tool.InvokableTool, error) {
	if config == nil {
		config = &ToolConfig{}
	}

	if config.ToolName == "" {
		config.ToolName = "crawlSinglePage"
	}
	if config.ToolDesc == "" {
		config.ToolDesc = "Retrieve content from a specific webpage. Returns the page title, content, URL and website."
	}

	service := NewService()
	wbt := &webBrowsingTool{
		service: service,
		config:  config,
	}

	crawlTool, err := utils.InferTool(config.ToolName, config.ToolDesc, wbt.CrawlSinglePage)
	if err != nil {
		return nil, fmt.Errorf("failed to infer crawl single page tool: %w", err)
	}

	return crawlTool, nil
}

// NewCrawlMultiPagesTool creates a tool for crawling multiple web pages
func NewCrawlMultiPagesTool(ctx context.Context, config *ToolConfig) (tool.InvokableTool, error) {
	if config == nil {
		config = &ToolConfig{}
	}

	if config.ToolName == "" {
		config.ToolName = "crawlMultiPages"
	}
	if config.ToolDesc == "" {
		config.ToolDesc = "Retrieve content from multiple webpages simultaneously. Returns an array of page results."
	}

	service := NewService()
	wbt := &webBrowsingTool{
		service: service,
		config:  config,
	}

	crawlTool, err := utils.InferTool(config.ToolName, config.ToolDesc, wbt.CrawlMultiPages)
	if err != nil {
		return nil, fmt.Errorf("failed to infer crawl multi pages tool: %w", err)
	}

	return crawlTool, nil
}

// NewWebBrowsingTools creates all web browsing tools (search, crawlSinglePage, crawlMultiPages)
func NewWebBrowsingTools(ctx context.Context) ([]tool.InvokableTool, error) {
	searchTool, err := NewSearchTool(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search tool: %w", err)
	}

	crawlSingleTool, err := NewCrawlSinglePageTool(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create crawl single page tool: %w", err)
	}

	crawlMultiTool, err := NewCrawlMultiPagesTool(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create crawl multi pages tool: %w", err)
	}

	return []tool.InvokableTool{searchTool, crawlSingleTool, crawlMultiTool}, nil
}

// SearchRequest represents the search request parameters
type SearchRequest struct {
	Query            string   `json:"query" jsonschema_description:"The search query string"`
	SearchCategories []string `json:"searchCategories,omitempty" jsonschema_description:"Search categories: general, images, news, science, videos"`
	SearchEngines    []string `json:"searchEngines,omitempty" jsonschema_description:"Search engines to use: google, bing, duckduckgo, brave, wikipedia, github, arxiv, etc."`
	SearchTimeRange  string   `json:"searchTimeRange,omitempty" jsonschema_description:"Time range filter: anytime, day, week, month, year"`
}

// SearchResultItem represents a single search result
type SearchResultItem struct {
	Title         string   `json:"title" jsonschema_description:"The title of the search result"`
	Content       string   `json:"content" jsonschema_description:"The content/description of the search result"`
	URL           string   `json:"url" jsonschema_description:"The URL of the search result"`
	Engines       []string `json:"engines,omitempty" jsonschema_description:"Search engines that returned this result"`
	PublishedDate string   `json:"publishedDate,omitempty" jsonschema_description:"Published date if available"`
}

// SearchResponse represents the search response
type SearchResponse struct {
	Query   string              `json:"query" jsonschema_description:"The original search query"`
	Results []*SearchResultItem `json:"results" jsonschema_description:"The search results"`
	Count   int                 `json:"count" jsonschema_description:"Number of results returned"`
}

// Search performs a web search
func (w *webBrowsingTool) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	query := strings.TrimSpace(req.Query)
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	searchQuery := SearchQuery{
		Query:            req.Query,
		SearchCategories: req.SearchCategories,
		SearchEngines:    req.SearchEngines,
		SearchTimeRange:  req.SearchTimeRange,
	}

	result, err := w.service.WebSearch(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	items := make([]*SearchResultItem, 0, len(result.Results))
	for _, r := range result.Results {
		items = append(items, &SearchResultItem{
			Title:         r.Title,
			Content:       r.Content,
			URL:           r.URL,
			Engines:       r.Engines,
			PublishedDate: r.PublishedDate,
		})
	}

	return &SearchResponse{
		Query:   result.Query,
		Results: items,
		Count:   len(items),
	}, nil
}

// CrawlSinglePageRequest represents a request to crawl a single page
type CrawlSinglePageRequest struct {
	URL string `json:"url" jsonschema_description:"The URL of the webpage to crawl"`
}

// CrawlPageResult represents the result of crawling a page
type CrawlPageResult struct {
	Title   string `json:"title" jsonschema_description:"The title of the webpage"`
	Content string `json:"content" jsonschema_description:"The content of the webpage"`
	URL     string `json:"url" jsonschema_description:"The URL of the webpage"`
	Website string `json:"website" jsonschema_description:"The website hostname"`
	Error   string `json:"error,omitempty" jsonschema_description:"Error message if crawl failed"`
}

// CrawlSinglePage crawls a single webpage and returns its content
func (w *webBrowsingTool) CrawlSinglePage(ctx context.Context, req *CrawlSinglePageRequest) (*CrawlPageResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if strings.TrimSpace(req.URL) == "" {
		return nil, fmt.Errorf("url is required")
	}

	resp, err := w.service.CrawlPages(CrawlPagesRequest{
		URLs: []string{req.URL},
	})
	if err != nil {
		return nil, fmt.Errorf("crawl failed: %w", err)
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("no results returned")
	}

	result := resp.Results[0]
	if result.Error != nil {
		return &CrawlPageResult{
			URL:   req.URL,
			Error: result.Error.ErrorMessage,
		}, nil
	}

	return &CrawlPageResult{
		Title:   result.Success.Title,
		Content: result.Success.Content,
		URL:     result.Success.URL,
		Website: result.Success.Website,
	}, nil
}

// CrawlMultiPagesRequest represents a request to crawl multiple pages
type CrawlMultiPagesRequest struct {
	URLs []string `json:"urls" jsonschema_description:"The URLs of the webpages to crawl"`
}

// CrawlMultiPagesResponse represents the response from crawling multiple pages
type CrawlMultiPagesResponse struct {
	Results []*CrawlPageResult `json:"results" jsonschema_description:"The crawl results for each URL"`
}

// CrawlMultiPages crawls multiple webpages and returns their content
func (w *webBrowsingTool) CrawlMultiPages(ctx context.Context, req *CrawlMultiPagesRequest) (*CrawlMultiPagesResponse, error) {
	if req == nil || len(req.URLs) == 0 {
		return nil, fmt.Errorf("urls is required")
	}

	resp, err := w.service.CrawlPages(CrawlPagesRequest{
		URLs: req.URLs,
	})
	if err != nil {
		return nil, fmt.Errorf("crawl failed: %w", err)
	}

	results := make([]*CrawlPageResult, 0, len(resp.Results))
	for i, r := range resp.Results {
		if r.Error != nil {
			results = append(results, &CrawlPageResult{
				URL:   req.URLs[i],
				Error: r.Error.ErrorMessage,
			})
		} else {
			results = append(results, &CrawlPageResult{
				Title:   r.Success.Title,
				Content: r.Success.Content,
				URL:     r.Success.URL,
				Website: r.Success.Website,
			})
		}
	}

	return &CrawlMultiPagesResponse{
		Results: results,
	}, nil
}
