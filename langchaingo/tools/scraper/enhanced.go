package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	readability "github.com/go-shiori/go-readability"
	"github.com/kawai-network/veridium/langchaingo/httputil"
)

// EnhancedOptions contains options for enhanced scraping
type EnhancedOptions struct {
	// EnableReadability enables Mozilla Readability algorithm for content extraction
	EnableReadability bool

	// PureText strips images and links from the output
	PureText bool

	// Timeout for HTTP request (default: 30 seconds)
	Timeout time.Duration

	// UserAgent for HTTP requests
	UserAgent string
}

// DefaultEnhancedOptions returns default options for enhanced scraping
func DefaultEnhancedOptions() EnhancedOptions {
	return EnhancedOptions{
		EnableReadability: true,
		PureText:          false,
		Timeout:           30 * time.Second,
		UserAgent:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	}
}

// CrawlEnhanced performs enhanced web scraping with content extraction and markdown conversion
func (s *Scraper) CrawlEnhanced(ctx context.Context, targetURL string, opts EnhancedOptions) (*EnhancedResult, error) {
	// Validate URL
	parsedURL, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid URL: %v", ErrScrapingFailed, err)
	}

	// Set defaults if not provided
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.UserAgent == "" {
		opts.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout:   opts.Timeout,
		Transport: httputil.DefaultTransport,
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %v", ErrScrapingFailed, err)
	}

	// Set headers
	req.Header.Set("User-Agent", opts.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Fetch HTML
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch URL: %v", ErrScrapingFailed, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: HTTP %d: %s", ErrScrapingFailed, resp.StatusCode, resp.Status)
	}

	// Read HTML
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response: %v", ErrScrapingFailed, err)
	}
	htmlContent := string(htmlBytes)

	// Extract content with readability if enabled
	var article readability.Article
	var contentHTML string

	if opts.EnableReadability {
		article, err = readability.FromReader(strings.NewReader(htmlContent), parsedURL)
		if err != nil {
			// Fallback to raw HTML if readability fails
			contentHTML = htmlContent
			article.Title = extractTitleFromHTML(htmlContent)
		} else {
			contentHTML = article.Content
		}
	} else {
		contentHTML = htmlContent
		article.Title = extractTitleFromHTML(htmlContent)
	}

	// Convert HTML to Markdown
	converter := md.NewConverter("", true, nil)

	// Configure converter for pure text if needed
	if opts.PureText {
		// For pure text mode, we'll strip images and convert links to plain text
		// This is handled by the converter's default behavior with custom rules
		converter.AddRules(md.Rule{
			Filter: []string{"img"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				empty := ""
				return &empty
			},
		})
		converter.AddRules(md.Rule{
			Filter: []string{"a"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				return &content
			},
		})
	}

	markdown, err := converter.ConvertString(contentHTML)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to convert to markdown: %v", ErrScrapingFailed, err)
	}

	// Clean up markdown
	markdown = strings.TrimSpace(markdown)

	// Build result
	publishedTime := ""
	if article.PublishedTime != nil {
		publishedTime = article.PublishedTime.Format(time.RFC3339)
	}

	result := &EnhancedResult{
		Content:       markdown,
		Title:         article.Title,
		Description:   article.Excerpt,
		SiteName:      article.SiteName,
		Author:        article.Byline,
		URL:           targetURL,
		Length:        len(markdown),
		PublishedTime: publishedTime,
		Lang:          article.Language,
	}

	return result, nil
}

// extractTitleFromHTML extracts title from HTML using simple string matching
func extractTitleFromHTML(html string) string {
	// Simple title extraction
	start := strings.Index(html, "<title>")
	if start == -1 {
		return ""
	}
	start += 7 // len("<title>")

	end := strings.Index(html[start:], "</title>")
	if end == -1 {
		return ""
	}

	title := html[start : start+end]
	return strings.TrimSpace(title)
}
