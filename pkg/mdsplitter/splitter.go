// Package mdsplitter provides a simple markdown header-based text splitter.
// This is a simplified version inspired by cloudwego/eino-ext markdown splitter,
// but with no external dependencies and a simpler API.
package mdsplitter

import (
	"fmt"
	"strings"
)

// Chunk represents a split chunk of text with its metadata
type Chunk struct {
	Content  string
	Metadata map[string]string
}

// Config configures the markdown header splitter
type Config struct {
	// Headers specify the headers to be identified and their names in chunk metadata.
	// Headers can only consist of '#'.
	// Example:
	//   Headers: map[string]string{
	//     "##":  "h2",
	//     "###": "h3",
	//   }
	Headers map[string]string

	// TrimHeaders specifies if results should exclude header lines.
	// If false, headers are included in the chunk content.
	TrimHeaders bool
}

// Splitter splits markdown text by headers
type Splitter struct {
	headers     map[string]string
	trimHeaders bool
}

// New creates a new markdown header splitter
func New(config *Config) (*Splitter, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("no headers specified")
	}

	// Validate that headers only contain '#'
	for header := range config.Headers {
		for _, c := range header {
			if c != '#' {
				return nil, fmt.Errorf("header can only consist of '#': %s", header)
			}
		}
	}

	return &Splitter{
		headers:     config.Headers,
		trimHeaders: config.TrimHeaders,
	}, nil
}

// Split splits the given text by markdown headers
func (s *Splitter) Split(text string) []Chunk {
	return s.splitText(text)
}

const (
	codeFenceBacktick = "```"
	codeFenceTilde    = "~~~"
)

type metaRecord struct {
	name  string
	level int
	data  string
}

func (s *Splitter) splitText(text string) []Chunk {
	var recordedMetaList []metaRecord
	recordedMetaMap := make(map[string]string)
	var currentLines []string
	var inCodeBlock bool
	var openingFence string
	var chunks []Chunk

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(line) == 0 {
			currentLines = append(currentLines, line)
			continue
		}

		trimmedLine := strings.TrimSpace(line)

		// Handle code blocks (don't process headers inside code blocks)
		if !inCodeBlock {
			if strings.HasPrefix(trimmedLine, codeFenceBacktick) && strings.Count(trimmedLine, codeFenceBacktick) == 1 {
				inCodeBlock = true
				openingFence = codeFenceBacktick
			} else if strings.HasPrefix(trimmedLine, codeFenceTilde) {
				inCodeBlock = true
				openingFence = codeFenceTilde
			}
		} else {
			if strings.HasPrefix(trimmedLine, openingFence) {
				inCodeBlock = false
				openingFence = ""
			}
		}

		if inCodeBlock {
			currentLines = append(currentLines, line)
			continue
		}

		// Check if the line starts with any configured header
		isNewHeader := false
		for header, name := range s.headers {
			if strings.HasPrefix(trimmedLine, header) && (len(trimmedLine) == len(header) || trimmedLine[len(header)] == ' ') {
				// Save current chunk if we have accumulated lines
				if len(currentLines) > 0 {
					chunks = append(chunks, Chunk{
						Content:  strings.Join(currentLines, "\n"),
						Metadata: deepCopyMap(recordedMetaMap),
					})
					currentLines = nil
				}

				// Add header to current chunk if not trimming
				if !s.trimHeaders {
					currentLines = append(currentLines, line)
				}

				// Update metadata tracking
				newLevel := len(header)

				// Remove metadata from higher or equal level headers
				for i := len(recordedMetaList) - 1; i >= 0; i-- {
					if recordedMetaList[i].level >= newLevel {
						delete(recordedMetaMap, recordedMetaList[i].name)
						recordedMetaList = recordedMetaList[:i]
					} else {
						break
					}
				}

				// Add new header metadata
				headerData := strings.TrimSpace(trimmedLine[len(header):])
				recordedMetaList = append(recordedMetaList, metaRecord{
					name:  name,
					level: newLevel,
					data:  headerData,
				})
				recordedMetaMap[name] = headerData

				isNewHeader = true
				break
			}
		}

		if !isNewHeader {
			currentLines = append(currentLines, line)
		}
	}

	// Add final chunk
	if len(currentLines) > 0 {
		chunks = append(chunks, Chunk{
			Content:  strings.Join(currentLines, "\n"),
			Metadata: deepCopyMap(recordedMetaMap),
		})
	}

	return chunks
}

func deepCopyMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
