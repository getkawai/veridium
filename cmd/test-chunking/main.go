package main

import (
	"fmt"
	"strings"
)

// Simplified version of RAGProcessor for testing
type testProcessor struct {
	chunkSize   int
	overlapSize int
}

func main() {
	processor := &testProcessor{
		chunkSize:   100,
		overlapSize: 20,
	}

	// Test cases
	tests := []struct {
		name string
		text string
	}{
		{
			name: "Normal text with periods",
			text: "This is a sentence. This is another sentence. And here is a third one. Finally, a fourth sentence to test chunking.",
		},
		{
			name: "Long text without periods",
			text: strings.Repeat("word ", 50), // 250 chars, no periods
		},
		{
			name: "Code-like text",
			text: `func main() {
	fmt.Println("Hello")
	for i := 0; i < 10; i++ {
		process(i)
	}
}`,
		},
		{
			name: "Mixed content",
			text: "Paragraph one with some text.\n\nParagraph two with more content.\n\nParagraph three is here.",
		},
	}

	for _, tt := range tests {
		fmt.Printf("\n=== Test: %s ===\n", tt.name)
		fmt.Printf("Input length: %d chars\n", len(tt.text))

		chunks := processor.chunkText(tt.text)

		fmt.Printf("Output: %d chunks\n", len(chunks))
		for i, chunk := range chunks {
			fmt.Printf("  Chunk %d (%d chars): %q\n", i+1, len(chunk), chunk)
		}
	}
}

// chunkText - copy of the new implementation
func (r *testProcessor) chunkText(text string) []string {
	if text == "" {
		return []string{}
	}

	separators := []string{
		"\n\n", "\n", ". ", "? ", "! ", "; ", ", ", " ",
	}

	return r.recursiveSplit(text, separators, 0)
}

func (r *testProcessor) recursiveSplit(text string, separators []string, depth int) []string {
	if len(text) <= r.chunkSize {
		return []string{text}
	}

	if depth >= len(separators) {
		return r.forceSplitBySize(text)
	}

	separator := separators[depth]

	if !strings.Contains(text, separator) {
		return r.recursiveSplit(text, separators, depth+1)
	}

	parts := strings.Split(text, separator)

	var finalChunks []string
	var goodParts []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if len(part) > r.chunkSize {
			if len(goodParts) > 0 {
				merged := r.mergeParts(goodParts, separator)
				finalChunks = append(finalChunks, merged...)
				goodParts = nil
			}

			subChunks := r.recursiveSplit(part, separators, depth+1)
			finalChunks = append(finalChunks, subChunks...)
		} else {
			goodParts = append(goodParts, part)
		}
	}

	if len(goodParts) > 0 {
		merged := r.mergeParts(goodParts, separator)
		finalChunks = append(finalChunks, merged...)
	}

	return finalChunks
}

func (r *testProcessor) mergeParts(parts []string, separator string) []string {
	var chunks []string
	var currentChunk strings.Builder

	for _, part := range parts {
		partLen := len(part)
		sepLen := 0
		if currentChunk.Len() > 0 {
			sepLen = len(separator)
		}

		if currentChunk.Len() > 0 && currentChunk.Len()+sepLen+partLen > r.chunkSize {
			chunks = append(chunks, currentChunk.String())

			currentChunk.Reset()
			if r.overlapSize > 0 && len(chunks) > 0 {
				prevChunk := chunks[len(chunks)-1]
				overlapStart := len(prevChunk) - r.overlapSize
				if overlapStart < 0 {
					overlapStart = 0
				}
				currentChunk.WriteString(prevChunk[overlapStart:])
				currentChunk.WriteString(separator)
			}
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(separator)
		}
		currentChunk.WriteString(part)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

func (r *testProcessor) forceSplitBySize(text string) []string {
	var chunks []string

	for len(text) > 0 {
		if len(text) <= r.chunkSize {
			chunks = append(chunks, text)
			break
		}

		chunk := text[:r.chunkSize]
		chunks = append(chunks, chunk)

		step := r.chunkSize - r.overlapSize
		if step <= 0 {
			step = r.chunkSize
		}
		text = text[step:]
	}

	return chunks
}
