package types

// SearchParams represents optional search parameters
type SearchParams struct {
	SearchCategories []string `json:"searchCategories,omitempty"`
	SearchEngines    []string `json:"searchEngines,omitempty"`
	SearchTimeRange  string   `json:"searchTimeRange,omitempty"`
}

// UniformSearchResult represents a single search result
type UniformSearchResult struct {
	Category      string   `json:"category,omitempty"`
	Content       string   `json:"content"`
	Engines       []string `json:"engines"`
	IframeSrc     string   `json:"iframeSrc,omitempty"`
	ImgSrc        string   `json:"imgSrc,omitempty"`
	ParsedUrl     string   `json:"parsedUrl"`
	PublishedDate string   `json:"publishedDate,omitempty"`
	Score         float64  `json:"score"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	Title         string   `json:"title"`
	URL           string   `json:"url"`
}

// UniformSearchResponse represents the response from a search query
type UniformSearchResponse struct {
	CostTime      int64                 `json:"costTime"` // milliseconds
	Query         string                `json:"query"`
	ResultNumbers int                   `json:"resultNumbers"`
	Results       []UniformSearchResult `json:"results"`
}

