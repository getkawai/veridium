# mdsplitter

A simple, dependency-free markdown header-based text splitter for Go.

## Features

- Split markdown text by headers (e.g., `##`, `###`)
- Preserve or trim headers in output chunks
- Track header hierarchy in metadata
- Handle code blocks correctly (don't split on headers inside code blocks)
- Zero external dependencies

## Usage

```go
package main

import (
	"fmt"
	"github.com/kawai-network/veridium/pkg/mdsplitter"
)

func main() {
	// Create a splitter configuration
	config := &mdsplitter.Config{
		Headers: map[string]string{
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false, // Keep headers in chunks
	}

	// Create the splitter
	splitter, err := mdsplitter.New(config)
	if err != nil {
		panic(err)
	}

	// Split markdown text
	text := `## Introduction

This is the introduction section.

### Background

Some background information here.

## Main Content

The main content goes here.`

	chunks := splitter.Split(text)

	// Process chunks
	for i, chunk := range chunks {
		fmt.Printf("Chunk %d:\n", i+1)
		fmt.Printf("Content: %s\n", chunk.Content)
		fmt.Printf("Metadata: %+v\n\n", chunk.Metadata)
	}
}
```

## API

### Config

```go
type Config struct {
	// Headers maps header prefixes to metadata keys
	// Example: map[string]string{"##": "h2", "###": "h3"}
	Headers map[string]string

	// TrimHeaders removes header lines from chunk content if true
	TrimHeaders bool
}
```

### Chunk

```go
type Chunk struct {
	Content  string            // The text content of the chunk
	Metadata map[string]string // Header metadata (e.g., {"h2": "Section Title"})
}
```

### Splitter

```go
// New creates a new markdown header splitter
func New(config *Config) (*Splitter, error)

// Split splits the given text by markdown headers
func (s *Splitter) Split(text string) []Chunk
```

## Design

This package is a simplified version inspired by the cloudwego/eino-ext markdown splitter, but with:
- No external dependencies (no eino/schema imports)
- Simpler API (direct text input/output instead of Document objects)
- Focused functionality (just splitting, no transformation pipeline)
- Easy to integrate into existing codebases

## License

Same as the parent project.

