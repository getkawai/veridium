package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

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
		impls = []CrawlImplType{CrawlImplJina}
	}

	return &Crawler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
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

// crawlSingle crawls a single URL
func (c *Crawler) crawlSingle(urlStr string, impls []CrawlImplType) CrawlResult {
	// Try each implementation in order
	for _, impl := range impls {
		switch impl {
		case CrawlImplJina:
			result, err := c.crawlWithJina(urlStr)
			if err == nil {
				return result
			}
		case CrawlImplNaive:
			result, err := c.crawlNaive(urlStr)
			if err == nil {
				return result
			}
		}
	}

	// All implementations failed
	return CrawlResult{
		Error: &CrawlErrorResult{
			ErrorMessage: "all crawler implementations failed",
			ErrorType:    "CRAWLER_FAILED",
			URL:          urlStr,
		},
	}
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
		},
	}, nil
}

// crawlNaive performs naive HTTP crawling
func (c *Crawler) crawlNaive(urlStr string) (CrawlResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return CrawlResult{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; VeridiumCrawler/1.0)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CrawlResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CrawlResult{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return CrawlResult{}, err
	}

	title := extractHTMLTitle(doc)
	content := extractTextContent(doc)
	website := extractWebsite(urlStr)

	return CrawlResult{
		Success: &CrawlSuccessResult{
			Title:   title,
			Content: content,
			URL:     urlStr,
			Website: website,
		},
	}, nil
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
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	return u.Hostname()
}
