package services

import (
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/pkg/mdsplitter"
	"github.com/kawai-network/veridium/types"
)

// ChunkDocument chunks a FileDocument based on its type and configuration
func (l *FileLoader) ChunkDocument(doc *types.FileDocument, config types.ChunkingConfig) []types.DocumentChunk {
	if !config.Enabled || doc == nil {
		return nil
	}

	// Choose chunking strategy based on file type
	switch types.SupportedFileType(doc.FileType) {
	case types.FileTypePDF:
		return l.chunkPDFByPages(doc.Pages, config)
	case types.FileTypeDOCX, types.FileTypePPTX, types.FileTypeXLSX, types.FileTypeTXT, types.FileTypeMarkdown:
		return l.chunkMarkdownWithEino(doc.Content, config)
	default:
		// Fallback to recursive splitter for all other text files
		return l.chunkByRecursiveSplit(doc.Content, config)
	}
}

// chunkMarkdownWithEino uses local markdown header splitter with size enforcement
func (l *FileLoader) chunkMarkdownWithEino(content string, config types.ChunkingConfig) []types.DocumentChunk {
	// Configure markdown splitter
	// Split by ## and ### headers
	splitterConfig := &mdsplitter.Config{
		Headers: map[string]string{
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false, // Keep headers in chunks for context
	}

	splitter, err := mdsplitter.New(splitterConfig)
	if err != nil {
		// Fallback to recursive split if splitter creation fails
		return l.chunkByRecursiveSplit(content, config)
	}

	// Split by headers
	mdChunks := splitter.Split(content)

	// Convert markdown chunks to DocumentChunks with size enforcement
	var chunks []types.DocumentChunk
	chunkID := 1

	for _, mdChunk := range mdChunks {
		chunkContent := mdChunk.Content

		// If chunk exceeds size limit, apply recursive splitting
		if len(chunkContent) > config.ChunkSize {
			subChunks := l.chunkText(chunkContent, config)
			for _, subContent := range subChunks {
				chunks = append(chunks, types.DocumentChunk{
					ID:      fmt.Sprintf("chunk-%d", chunkID),
					Content: subContent,
					Metadata: map[string]interface{}{
						"type":   "markdown-header",
						"h2":     mdChunk.Metadata["h2"],
						"h3":     mdChunk.Metadata["h3"],
						"source": "mdsplitter-split",
					},
				})
				chunkID++
			}
		} else {
			chunks = append(chunks, types.DocumentChunk{
				ID:      fmt.Sprintf("chunk-%d", chunkID),
				Content: chunkContent,
				Metadata: map[string]interface{}{
					"type":   "markdown-header",
					"h2":     mdChunk.Metadata["h2"],
					"h3":     mdChunk.Metadata["h3"],
					"source": "mdsplitter",
				},
			})
			chunkID++
		}
	}

	return chunks
}

// chunkPDFByPages chunks PDF by merging pages until chunkSize is reached
func (l *FileLoader) chunkPDFByPages(pages []types.DocumentPage, config types.ChunkingConfig) []types.DocumentChunk {
	var chunks []types.DocumentChunk
	var currentChunk strings.Builder
	var currentPages []int
	chunkID := 1

	for i, page := range pages {
		pageNum := i + 1

		// Check if adding this page would exceed chunk size
		if currentChunk.Len() > 0 && currentChunk.Len()+page.CharCount > config.ChunkSize {
			// Save current chunk
			chunks = append(chunks, types.DocumentChunk{
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
		chunks = append(chunks, types.DocumentChunk{
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
func (l *FileLoader) chunkByRecursiveSplit(content string, config types.ChunkingConfig) []types.DocumentChunk {
	textChunks := l.chunkText(content, config)

	var chunks []types.DocumentChunk
	for i, text := range textChunks {
		chunks = append(chunks, types.DocumentChunk{
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

// chunkText splits text into chunks using specified configuration
func (l *FileLoader) chunkText(text string, config types.ChunkingConfig) []string {
	if text == "" || !config.Enabled {
		return []string{}
	}

	// Use defaults if not specified
	chunkSize := config.ChunkSize
	overlapSize := config.OverlapSize
	if chunkSize <= 0 {
		chunkSize = 1000 // Default chunk size
	}
	if overlapSize < 0 {
		overlapSize = 200 // Default overlap size
	}

	// Define separators in order of preference (coarse to fine)
	// Try to keep semantic units together as much as possible
	separators := []string{
		"\n\n", // Paragraph breaks (highest priority)
		"\n",   // Line breaks
		". ",   // Sentences
		"? ",   // Questions
		"! ",   // Exclamations
		"; ",   // Semicolons
		", ",   // Commas
		" ",    // Words (last resort)
	}

	return l.recursiveSplit(text, separators, 0, chunkSize, overlapSize)
}

// recursiveSplit implements recursive text splitting with multiple separators
func (l *FileLoader) recursiveSplit(text string, separators []string, depth int, chunkSize, overlapSize int) []string {
	// Base case: text fits within chunk size
	if len(text) <= chunkSize {
		return []string{text}
	}

	// If we've exhausted all separators, force split by character count
	if depth >= len(separators) {
		return l.forceSplitBySize(text, chunkSize, overlapSize)
	}

	separator := separators[depth]

	// Check if separator exists in text
	if !strings.Contains(text, separator) {
		// Try next separator
		return l.recursiveSplit(text, separators, depth+1, chunkSize, overlapSize)
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

		if len(part) > chunkSize {
			// Part is too large, need to process accumulated good parts first
			if len(goodParts) > 0 {
				merged := l.mergeParts(goodParts, separator, chunkSize, overlapSize)
				finalChunks = append(finalChunks, merged...)
				goodParts = nil
			}

			// Recursively split large part with next separator
			subChunks := l.recursiveSplit(part, separators, depth+1, chunkSize, overlapSize)
			finalChunks = append(finalChunks, subChunks...)
		} else {
			// Part is small enough, accumulate it
			goodParts = append(goodParts, part)
		}
	}

	// Process remaining good parts
	if len(goodParts) > 0 {
		merged := l.mergeParts(goodParts, separator, chunkSize, overlapSize)
		finalChunks = append(finalChunks, merged...)
	}

	return finalChunks
}

// mergeParts merges small parts into chunks with overlap support
func (l *FileLoader) mergeParts(parts []string, separator string, chunkSize, overlapSize int) []string {
	var chunks []string
	var currentChunk strings.Builder

	for _, part := range parts {
		partLen := len(part)
		sepLen := 0
		if currentChunk.Len() > 0 {
			sepLen = len(separator)
		}

		// Check if adding this part would exceed chunk size
		if currentChunk.Len() > 0 && currentChunk.Len()+sepLen+partLen > chunkSize {
			// Save current chunk
			chunks = append(chunks, currentChunk.String())

			// Start new chunk with overlap
			currentChunk.Reset()
			if overlapSize > 0 && len(chunks) > 0 {
				prevChunk := chunks[len(chunks)-1]
				overlapStart := len(prevChunk) - overlapSize
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

// forceSplitBySize splits text by character count when no separator works
// This is a last resort to ensure we never exceed chunk size
func (l *FileLoader) forceSplitBySize(text string, chunkSize, overlapSize int) []string {
	var chunks []string

	for len(text) > 0 {
		if len(text) <= chunkSize {
			chunks = append(chunks, text)
			break
		}

		// Take chunk size worth of text
		chunk := text[:chunkSize]
		chunks = append(chunks, chunk)

		// Move forward with overlap
		step := chunkSize - overlapSize
		if step <= 0 {
			step = chunkSize
		}
		text = text[step:]
	}

	return chunks
}
