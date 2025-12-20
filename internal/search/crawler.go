package search

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kawai-network/veridium/pkg/htmltomarkdown"
	"github.com/kawai-network/veridium/pkg/htmltomarkdown/converter"
)

// Crawler handles web page crawling
type Crawler struct {
	httpClient *http.Client
	impls      []CrawlImplType
}

// NewCrawler creates a new crawler instance
func NewCrawler(impls []CrawlImplType) *Crawler {
	if len(impls) == 0 {
		// Default: try Naive only
		impls = []CrawlImplType{CrawlImplNaive}
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
	// We only have one implementation now
	result, err := c.crawlNaive(urlStr)
	if err == nil {
		return result
	}

	log.Printf("❌ All crawlers failed for %s: %v", urlStr, err)

	return CrawlResult{
		Error: &CrawlErrorResult{
			ErrorMessage: fmt.Sprintf("crawler failed: %v", err),
			ErrorType:    "CRAWLER_FAILED",
			URL:          urlStr,
		},
	}
}

// crawlNaive performs naive HTTP crawling using htmltomarkdown.ConvertURL
func (c *Crawler) crawlNaive(urlStr string) (CrawlResult, error) {
	log.Printf("🔍 [crawlNaive] URL: %s", urlStr)

	result, err := htmltomarkdown.ConvertURL(urlStr, c.httpClient,
		converter.WithCrawler("kawai"),
		converter.WithMainContentOnly(),
	)
	if err != nil {
		log.Printf("❌ [crawlNaive] Failed: %v", err)
		return CrawlResult{}, err
	}

	if result.Title == "" {
		log.Printf("⚠️  [crawlNaive] Empty title for %s, using URL as fallback", urlStr)
		// Use URL as title if empty
		result.Title = urlStr
	}

	log.Printf("✅ [crawlNaive] Title: %s, Website: %s, Content: %d chars", result.Title, result.Website, len(result.Content))
	return CrawlResult{
		Success: result,
	}, nil
}
