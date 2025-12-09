package search

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	htmltomd "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
)

// Crawler handles web page crawling
type Crawler struct {
	httpClient *http.Client
	impls      []CrawlImplType
}

// NewCrawler creates a new crawler instance
func NewCrawler(impls []CrawlImplType) *Crawler {
	if len(impls) == 0 {
		// Default: try Jina first, then fallback to Naive
		impls = []CrawlImplType{CrawlImplJina, CrawlImplNaive}
	}

	return &Crawler{
		httpClient: &http.Client{
			Timeout: 45 * time.Second, // Increased timeout for slow sites
			Transport: &http.Transport{
				MaxIdleConns:          10,
				IdleConnTimeout:       30 * time.Second,
				DisableCompression:    false, // Let Go auto-decompress gzip/deflate
				DisableKeepAlives:     false,
				MaxConnsPerHost:       5,
				ResponseHeaderTimeout: 30 * time.Second,
				ForceAttemptHTTP2:     false, // Stick to HTTP/1.1 for better compatibility
			},
		},
		impls: impls,
	}
}

// CrawlPages crawls multiple URLs concurrently
func (c *Crawler) CrawlPages(urls []string, impls []CrawlImplType) []CrawlResult {
	if len(impls) == 0 {
		impls = c.impls
	}

	results := make([]CrawlResult, len(urls))
	var wg sync.WaitGroup

	// Limit concurrency to 3
	semaphore := make(chan struct{}, 3)

	for i, urlStr := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()

			semaphore <- struct{}{}        // acquire
			defer func() { <-semaphore }() // release

			results[idx] = c.crawlSingle(u, impls)
		}(i, urlStr)
	}

	wg.Wait()
	return results
}

// crawlSingle crawls a single URL with retry logic
func (c *Crawler) crawlSingle(urlStr string, impls []CrawlImplType) CrawlResult {
	var lastErrors []string

	// Try each implementation in order
	for _, impl := range impls {
		switch impl {
		case CrawlImplJina:
			result, err := c.crawlWithJina(urlStr)
			if err == nil {
				return result
			}
			lastErrors = append(lastErrors, fmt.Sprintf("jina: %v", err))
			log.Printf("⚠️ Crawler[jina] failed for %s: %v", urlStr, err)

		case CrawlImplNaive:
			result, err := c.crawlNaive(urlStr)
			if err == nil {
				return result
			}
			lastErrors = append(lastErrors, fmt.Sprintf("naive: %v", err))
			log.Printf("⚠️ Crawler[naive] failed for %s: %v", urlStr, err)
		}
	}

	// If Naive wasn't in impls, try it as last resort
	if !containsImpl(impls, CrawlImplNaive) {
		result, err := c.crawlNaive(urlStr)
		if err == nil {
			log.Printf("✅ Crawler[naive-fallback] succeeded for %s", urlStr)
			return result
		}
		lastErrors = append(lastErrors, fmt.Sprintf("naive-fallback: %v", err))
	}

	// All implementations failed
	errorMsg := strings.Join(lastErrors, "; ")
	log.Printf("❌ All crawlers failed for %s: %s", urlStr, errorMsg)

	return CrawlResult{
		Error: &CrawlErrorResult{
			ErrorMessage: fmt.Sprintf("all crawlers failed: %s", errorMsg),
			ErrorType:    "CRAWLER_FAILED",
			URL:          urlStr,
		},
	}
}

// containsImpl checks if impl is in the list
func containsImpl(impls []CrawlImplType, impl CrawlImplType) bool {
	for _, i := range impls {
		if i == impl {
			return true
		}
	}
	return false
}

// crawlWithJina crawls using Jina Reader API
func (c *Crawler) crawlWithJina(urlStr string) (CrawlResult, error) {
	// Jina Reader API: https://r.jina.ai/{url}
	jinaURL := fmt.Sprintf("https://r.jina.ai/%s", urlStr)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", jinaURL, nil)
	if err != nil {
		return CrawlResult{}, err
	}

	req.Header.Set("Accept", "text/plain")
	req.Header.Set("X-Respond-With", "markdown") // Request markdown format

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CrawlResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return CrawlResult{}, fmt.Errorf("jina API returned status %d: %s", resp.StatusCode, string(body))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return CrawlResult{}, err
	}

	// Extract title (Jina usually puts it at the top)
	title := extractTitle(string(content), urlStr)
	website := extractWebsite(urlStr)

	return CrawlResult{
		Success: &CrawlSuccessResult{
			Title:   title,
			Content: string(content),
			URL:     urlStr,
			Website: website,
			Crawler: "jina",
		},
	}, nil
}

// crawlNaive performs naive HTTP crawling with better browser simulation
// Returns content in Markdown format (same as Jina)
func (c *Crawler) crawlNaive(urlStr string) (CrawlResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return CrawlResult{}, err
	}

	// Better browser simulation headers
	// Note: Only gzip/deflate - Go auto-decompresses these. Avoid brotli (br) as it needs manual handling.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate") // Go HTTP client auto-decompresses gzip/deflate
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CrawlResult{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Accept 2xx and 3xx status codes
	if resp.StatusCode >= 400 {
		return CrawlResult{}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read body with limit to prevent memory issues
	limitedReader := io.LimitReader(resp.Body, 5*1024*1024) // 5MB max
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return CrawlResult{}, fmt.Errorf("read body failed: %w", err)
	}

	// Check for binary/corrupted content (likely wrong encoding)
	if !isValidUTF8Content(body) {
		return CrawlResult{}, fmt.Errorf("response contains invalid/binary content")
	}

	// Parse HTML for title extraction
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return CrawlResult{}, fmt.Errorf("parse HTML failed: %w", err)
	}

	title := extractHTMLTitle(doc)
	website := extractWebsite(urlStr)

	// Convert HTML to Markdown using html-to-markdown library
	markdown, err := htmltomd.ConvertString(string(body), converter.WithDomain(website))
	if err != nil {
		// Fallback to plain text extraction if markdown conversion fails
		log.Printf("⚠️ Markdown conversion failed for %s: %v, falling back to text", urlStr, err)
		markdown = extractTextContent(doc)
	}

	// Clean up the markdown content
	markdown = cleanMarkdownContent(markdown)

	// Validate we got meaningful content
	if len(markdown) < 50 {
		return CrawlResult{}, fmt.Errorf("content too short (%d chars)", len(markdown))
	}

	return CrawlResult{
		Success: &CrawlSuccessResult{
			Title:   title,
			Content: markdown,
			URL:     urlStr,
			Website: website,
			Crawler: "kawai",
		},
	}, nil
}

// cleanMarkdownContent cleans up markdown content
func cleanMarkdownContent(content string) string {
	lines := strings.Split(content, "\n")
	var cleaned []string
	emptyCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip excessive empty lines (max 2 consecutive)
		if trimmed == "" {
			emptyCount++
			if emptyCount <= 2 {
				cleaned = append(cleaned, "")
			}
			continue
		}
		emptyCount = 0

		// Skip common noise patterns
		if isNoisePattern(trimmed) {
			continue
		}

		cleaned = append(cleaned, line)
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

// isNoisePattern checks if a line is common web noise
func isNoisePattern(line string) bool {
	lower := strings.ToLower(line)

	noisePatterns := []string{
		"skip to content",
		"skip to main",
		"cookie",
		"we use cookies",
		"accept all",
		"reject all",
		"privacy policy",
		"terms of service",
		"subscribe to",
		"sign up for",
		"newsletter",
		"advertisement",
		"sponsored",
		"loading...",
		"please wait",
		"javascript is required",
		"enable javascript",
	}

	for _, pattern := range noisePatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}

// isValidUTF8Content checks if content is valid UTF-8 text (not binary garbage)
func isValidUTF8Content(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Check first 1KB for binary indicators
	checkLen := len(data)
	if checkLen > 1024 {
		checkLen = 1024
	}

	nullCount := 0
	nonPrintable := 0
	for i := 0; i < checkLen; i++ {
		b := data[i]
		if b == 0 {
			nullCount++
		}
		// Count non-printable non-whitespace ASCII control chars
		if b < 32 && b != '\t' && b != '\n' && b != '\r' {
			nonPrintable++
		}
	}

	// If >5% null bytes or >20% non-printable, likely binary
	if float64(nullCount)/float64(checkLen) > 0.05 {
		return false
	}
	if float64(nonPrintable)/float64(checkLen) > 0.20 {
		return false
	}

	return true
}

// extractTitle extracts title from content
func extractTitle(content string, fallbackURL string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			// First non-empty, non-header line might be the title
			if len(line) < 200 {
				return line
			}
		}
		// If it starts with "# ", it's likely the title
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return fallbackURL
}

// extractHTMLTitle extracts title from HTML document
func extractHTMLTitle(n *html.Node) string {
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.TrimSpace(title)
}

// extractTextContent extracts text content from HTML document
func extractTextContent(n *html.Node) string {
	var buf strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				buf.WriteString(text)
				buf.WriteString(" ")
			}
		}
		// Skip script and style tags
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.TrimSpace(buf.String())
}

// extractWebsite extracts the website hostname from URL
func extractWebsite(urlStr string) string {
	if urlStr == "" {
		return ""
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	hostname := u.Hostname()
	if hostname == "" {
		return urlStr
	}
	return hostname
}
