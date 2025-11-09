package scraper

// EnhancedResult represents a structured web scraping result with content extraction
type EnhancedResult struct {
	// Content is the main text content in Markdown format
	Content string `json:"content"`

	// Title is the page title
	Title string `json:"title"`

	// Description is the page description or excerpt
	Description string `json:"description"`

	// SiteName is the name of the website
	SiteName string `json:"siteName,omitempty"`

	// Author is the article author (if available)
	Author string `json:"author,omitempty"`

	// URL is the original URL that was scraped
	URL string `json:"url"`

	// Length is the content length in characters
	Length int `json:"length"`

	// PublishedTime is the published time (if available)
	PublishedTime string `json:"publishedTime,omitempty"`

	// Lang is the content language (if detected)
	Lang string `json:"lang,omitempty"`
}
