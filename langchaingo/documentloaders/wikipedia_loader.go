package documentloaders

import (
	"context"
	"fmt"

	"github.com/kawai-network/veridium/langchaingo/schema"
	"github.com/kawai-network/veridium/langchaingo/textsplitter"
)

// Wikipedia loads Wikipedia articles.
// It implements the Loader interface.
type Wikipedia struct {
	query       string
	language    string
	maxArticles int
}

// NewWikipedia creates a new Wikipedia document loader.
func NewWikipedia(query, language string, maxArticles int) *Wikipedia {
	return &Wikipedia{
		query:       query,
		language:    language,
		maxArticles: maxArticles,
	}
}

// Load fetches Wikipedia articles and returns them as documents.
func (w *Wikipedia) Load(ctx context.Context) ([]schema.Document, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Search Wikipedia for articles
	titles, err := SearchWikipedia(w.query, w.language, w.maxArticles)
	if err != nil {
		return nil, fmt.Errorf("wikipedia search failed: %w", err)
	}

	if len(titles) == 0 {
		return nil, fmt.Errorf("no Wikipedia articles found for query: %s", w.query)
	}

	// Load each article
	docs := make([]schema.Document, 0, len(titles))
	for _, title := range titles {
		// Check context cancellation for each article
		select {
		case <-ctx.Done():
			return docs, ctx.Err()
		default:
		}

		content, err := GetWikipediaArticle(title, w.language)
		if err != nil {
			// Log error but continue with other articles
			fmt.Printf("Warning: failed to load article '%s': %v\n", title, err)
			continue
		}

		// Clean the content
		cleanContent := CleanText(content)

		docs = append(docs, schema.Document{
			PageContent: cleanContent,
			Metadata: map[string]any{
				"source":   "wikipedia",
				"title":    title,
				"language": w.language,
				"query":    w.query,
			},
		})
	}

	return docs, nil
}

// LoadAndSplit loads Wikipedia articles and splits them into chunks using a text splitter.
func (w *Wikipedia) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := w.Load(ctx)
	if err != nil {
		return nil, err
	}

	return textsplitter.SplitDocuments(splitter, docs)
}

// Compile-time check that Wikipedia implements Loader
var _ Loader = (*Wikipedia)(nil)
