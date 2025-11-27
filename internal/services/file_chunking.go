package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
)

// ChunkDocument chunks a FileDocument based on its type and configuration
func (l *FileLoader) ChunkDocument(doc *FileDocument, config ChunkingConfig) []DocumentChunk {
	if !config.Enabled || doc == nil {
		return nil
	}

	// Choose chunking strategy based on file type
	switch doc.FileType {
	case "pdf":
		return l.chunkPDFByPages(doc.Pages, config)
	case "docx", "pptx", "xlsx", "md", "markdown":
		return l.chunkMarkdownWithEino(doc.Content, config)
	case "go", "py", "js", "ts", "java", "cpp", "c", "rs":
		// For code files, use recursive splitter
		// TODO: Implement AST-based chunking for better semantic splitting
		return l.chunkByRecursiveSplit(doc.Content, config)
	default:
		// Fallback to recursive splitter for all other text files
		return l.chunkByRecursiveSplit(doc.Content, config)
	}
}

// chunkMarkdownWithEino uses Eino's markdown header splitter with size enforcement
func (l *FileLoader) chunkMarkdownWithEino(content string, config ChunkingConfig) []DocumentChunk {
	ctx := context.Background()

	// Configure Eino markdown splitter
	// Split by ## and ### headers
	einoConfig := &markdown.HeaderConfig{
		Headers: map[string]string{
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false, // Keep headers in chunks for context
	}

	splitter, err := markdown.NewHeaderSplitter(ctx, einoConfig)
	if err != nil {
		// Fallback to recursive split if Eino fails
		return l.chunkByRecursiveSplit(content, config)
	}

	// Create Eino document
	einoDoc := &schema.Document{
		ID:      "doc",
		Content: content,
	}

	// Split by headers
	einoDocs, err := splitter.Transform(ctx, []*schema.Document{einoDoc})
	if err != nil {
		// Fallback to recursive split
		return l.chunkByRecursiveSplit(content, config)
	}

	// Convert Eino chunks to DocumentChunks with size enforcement
	var chunks []DocumentChunk
	chunkID := 1

	for _, einoChunk := range einoDocs {
		chunkContent := einoChunk.Content

		// If chunk exceeds size limit, apply recursive splitting
		if len(chunkContent) > config.ChunkSize {
			subChunks := l.recursiveSplitText(chunkContent, config)
			for _, subContent := range subChunks {
				chunks = append(chunks, DocumentChunk{
					ID:      fmt.Sprintf("chunk-%d", chunkID),
					Content: subContent,
					Metadata: map[string]interface{}{
						"type":   "markdown-header",
						"h2":     einoChunk.MetaData["h2"],
						"h3":     einoChunk.MetaData["h3"],
						"source": "eino-split",
					},
				})
				chunkID++
			}
		} else {
			chunks = append(chunks, DocumentChunk{
				ID:      fmt.Sprintf("chunk-%d", chunkID),
				Content: chunkContent,
				Metadata: map[string]interface{}{
					"type":   "markdown-header",
					"h2":     einoChunk.MetaData["h2"],
					"h3":     einoChunk.MetaData["h3"],
					"source": "eino",
				},
			})
			chunkID++
		}
	}

	return chunks
}

// chunkPDFByPages chunks PDF by merging pages until chunkSize is reached
func (l *FileLoader) chunkPDFByPages(pages []DocumentPage, config ChunkingConfig) []DocumentChunk {
	var chunks []DocumentChunk
	var currentChunk strings.Builder
	var currentPages []int
	chunkID := 1

	for i, page := range pages {
		pageNum := i + 1

		// Check if adding this page would exceed chunk size
		if currentChunk.Len() > 0 && currentChunk.Len()+page.CharCount > config.ChunkSize {
			// Save current chunk
			chunks = append(chunks, DocumentChunk{
				ID:      fmt.Sprintf("chunk-%d", chunkID),
				Content: currentChunk.String(),
				Metadata: map[string]interface{}{
					"type":  "pdf-pages",
					"pages": currentPages,
				},
			})
			chunkID++

			// Start new chunk with overlap (last page content)
			currentChunk.Reset()
			currentPages = nil

			if config.OverlapSize > 0 && i > 0 {
				prevPage := pages[i-1].PageContent
				overlapStart := len(prevPage) - config.OverlapSize
				if overlapStart < 0 {
					overlapStart = 0
				}
				currentChunk.WriteString(prevPage[overlapStart:])
				currentChunk.WriteString("\n\n")
			}
		}

		// Add page to current chunk
		currentChunk.WriteString(page.PageContent)
		currentChunk.WriteString("\n\n")
		currentPages = append(currentPages, pageNum)
	}

	// Add final chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, DocumentChunk{
			ID:      fmt.Sprintf("chunk-%d", chunkID),
			Content: currentChunk.String(),
			Metadata: map[string]interface{}{
				"type":  "pdf-pages",
				"pages": currentPages,
			},
		})
	}

	return chunks
}

// chunkByRecursiveSplit uses recursive splitting strategy (fallback)
func (l *FileLoader) chunkByRecursiveSplit(content string, config ChunkingConfig) []DocumentChunk {
	textChunks := l.recursiveSplitText(content, config)

	var chunks []DocumentChunk
	for i, text := range textChunks {
		chunks = append(chunks, DocumentChunk{
			ID:      fmt.Sprintf("chunk-%d", i+1),
			Content: text,
			Metadata: map[string]interface{}{
				"type":   "recursive-split",
				"source": "fallback",
			},
		})
	}

	return chunks
}

// recursiveSplitText implements recursive text splitting (reused from RAGProcessor logic)
func (l *FileLoader) recursiveSplitText(text string, config ChunkingConfig) []string {
	separators := []string{
		"\n\n", // Paragraph breaks
		"\n",   // Line breaks
		". ",   // Sentences
		"? ",   // Questions
		"! ",   // Exclamations
		"; ",   // Semicolons
		", ",   // Commas
		" ",    // Words
	}

	return l.recursiveSplit(text, separators, 0, config)
}

// recursiveSplit implements the recursive splitting algorithm
func (l *FileLoader) recursiveSplit(text string, separators []string, depth int, config ChunkingConfig) []string {
	// Base case: text fits within chunk size
	if len(text) <= config.ChunkSize {
		return []string{text}
	}

	// If exhausted all separators, force split by size
	if depth >= len(separators) {
		return l.forceSplitBySize(text, config)
	}

	separator := separators[depth]

	// Check if separator exists in text
	if !strings.Contains(text, separator) {
		// Try next separator
		return l.recursiveSplit(text, separators, depth+1, config)
	}

	// Split by current separator
	parts := strings.Split(text, separator)

	var finalChunks []string
	var goodParts []string

	// Process each part
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if len(part) > config.ChunkSize {
			// Part too large, process accumulated good parts first
			if len(goodParts) > 0 {
				merged := l.mergeParts(goodParts, separator, config)
				finalChunks = append(finalChunks, merged...)
				goodParts = nil
			}

			// Recursively split large part
			subChunks := l.recursiveSplit(part, separators, depth+1, config)
			finalChunks = append(finalChunks, subChunks...)
		} else {
			// Part is small enough, accumulate it
			goodParts = append(goodParts, part)
		}
	}

	// Process remaining good parts
	if len(goodParts) > 0 {
		merged := l.mergeParts(goodParts, separator, config)
		finalChunks = append(finalChunks, merged...)
	}

	return finalChunks
}

// mergeParts merges small parts into chunks with overlap
func (l *FileLoader) mergeParts(parts []string, separator string, config ChunkingConfig) []string {
	var chunks []string
	var currentChunk strings.Builder

	for _, part := range parts {
		partLen := len(part)
		sepLen := 0
		if currentChunk.Len() > 0 {
			sepLen = len(separator)
		}

		// Check if adding this part would exceed chunk size
		if currentChunk.Len() > 0 && currentChunk.Len()+sepLen+partLen > config.ChunkSize {
			// Save current chunk
			chunks = append(chunks, currentChunk.String())

			// Start new chunk with overlap
			currentChunk.Reset()
			if config.OverlapSize > 0 && len(chunks) > 0 {
				prevChunk := chunks[len(chunks)-1]
				overlapStart := len(prevChunk) - config.OverlapSize
				if overlapStart < 0 {
					overlapStart = 0
				}
				currentChunk.WriteString(prevChunk[overlapStart:])
				currentChunk.WriteString(separator)
			}
		}

		// Add part to current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(separator)
		}
		currentChunk.WriteString(part)
	}

	// Add final chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// forceSplitBySize splits text by character count (last resort)
func (l *FileLoader) forceSplitBySize(text string, config ChunkingConfig) []string {
	var chunks []string

	for len(text) > 0 {
		if len(text) <= config.ChunkSize {
			chunks = append(chunks, text)
			break
		}

		// Take chunk size worth of text
		chunk := text[:config.ChunkSize]
		chunks = append(chunks, chunk)

		// Move forward with overlap
		step := config.ChunkSize - config.OverlapSize
		if step <= 0 {
			step = config.ChunkSize
		}
		text = text[step:]
	}

	return chunks
}
