/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package chromem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/eino-adapters/chromem/parsers"
)

// FileManager manages file-based document indexing with automatic parsing and chunking
type FileManager struct {
	parsers   map[string]parser.Parser // extension -> parser
	splitter  document.Transformer
	indexer   *Indexer
	assetDir  string
	fileIndex map[string][]string // filename -> doc IDs
}

// FileManagerConfig holds configuration for FileManager
type FileManagerConfig struct {
	// Indexer is the Eino-compatible indexer to use for storing documents
	// Required.
	Indexer *Indexer

	// AssetDir is the directory where original files are copied
	// Optional, defaults to "./assets"
	AssetDir string

	// ChunkSize is the maximum size of each text chunk
	// Optional, defaults to 1000
	ChunkSize int

	// OverlapSize is the overlap between consecutive chunks
	// Optional, defaults to 200
	OverlapSize int

	// CustomParsers allows adding custom parsers for specific file extensions
	// Optional. Key is file extension (e.g., ".custom"), value is the parser
	CustomParsers map[string]parser.Parser
}

// NewFileManager creates a new file manager with automatic document parsing
func NewFileManager(ctx context.Context, config *FileManagerConfig) (*FileManager, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if config.Indexer == nil {
		return nil, fmt.Errorf("indexer is required")
	}
	if config.AssetDir == "" {
		config.AssetDir = "./assets"
	}
	if config.ChunkSize <= 0 {
		config.ChunkSize = 1000
	}
	if config.OverlapSize < 0 {
		config.OverlapSize = 200
	}

	// Create asset directory if it doesn't exist
	if err := os.MkdirAll(config.AssetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create asset directory: %w", err)
	}

	// Create custom parsers using gooxml and other libraries
	docxParser, _ := parsers.NewDocxParser(ctx)
	xlsxParser, _ := parsers.NewXlsxParser(ctx, &parsers.XlsxParserConfig{
		SheetsAsDocuments: false, // All sheets in one document
	})
	pdfParser, _ := parsers.NewPdfParser(ctx)
	htmlParser, _ := parsers.NewHtmlParser(ctx, &parsers.HtmlParserConfig{
		PreserveStructure: true,
	})
	textParser, _ := parsers.NewTextParser(ctx)

	// Map extensions to parsers
	parserMap := map[string]parser.Parser{
		".docx": docxParser,
		".xlsx": xlsxParser,
		".pdf":  pdfParser,
		".html": htmlParser,
		".htm":  htmlParser,
		".txt":  textParser,
		".md":   textParser,
	}

	// Add custom parsers
	if config.CustomParsers != nil {
		for ext, p := range config.CustomParsers {
			parserMap[ext] = p
		}
	}

	// Create simple text splitter (using langchaingo textsplitter)
	// Note: For production, you may want to use Eino-Ext's recursive splitter
	// For now, we'll use a simple implementation
	splitter := &simpleTextSplitter{
		chunkSize:   config.ChunkSize,
		overlapSize: config.OverlapSize,
	}

	return &FileManager{
		parsers:   parserMap,
		splitter:  splitter,
		indexer:   config.Indexer,
		assetDir:  config.AssetDir,
		fileIndex: make(map[string][]string),
	}, nil
}

// StoreFile loads, parses, splits, and indexes a file
func (fm *FileManager) StoreFile(ctx context.Context, filePath string, metadata map[string]any) error {
	// 1. Determine file type
	ext := strings.ToLower(filepath.Ext(filePath))
	p, ok := fm.parsers[ext]
	if !ok {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	// 2. Open file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 3. Parse file
	docs, err := p.Parse(ctx, file)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// 4. Split documents into chunks
	chunks, err := fm.splitter.Transform(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to split documents: %w", err)
	}

	// 5. Add metadata to all chunks
	fileName := filepath.Base(filePath)
	for _, chunk := range chunks {
		if chunk.MetaData == nil {
			chunk.MetaData = make(map[string]any)
		}
		// Add user metadata
		for k, v := range metadata {
			chunk.MetaData[k] = v
		}
		// Add system metadata
		chunk.MetaData["source_file"] = fileName
		chunk.MetaData["file_type"] = ext
		chunk.MetaData["file_path"] = filePath
	}

	// 6. Index using chromem (via Eino indexer)
	ids, err := fm.indexer.Store(ctx, chunks)
	if err != nil {
		return fmt.Errorf("failed to index documents: %w", err)
	}

	// 7. Track in file index
	fm.fileIndex[fileName] = ids

	// 8. Copy to asset directory
	return fm.copyFile(filePath)
}

// RemoveFile removes a file and its chunks from the index
func (fm *FileManager) RemoveFile(ctx context.Context, filename string) error {
	ids, ok := fm.fileIndex[filename]
	if !ok {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Delete from chromem
	// Note: chromem doesn't have a Delete method yet, so we just remove from tracking
	// TODO: Implement Delete in chromem collection
	_ = ids // Will be used when Delete is implemented

	// Remove from index
	delete(fm.fileIndex, filename)

	// Remove from asset directory
	assetPath := filepath.Join(fm.assetDir, filename)
	if err := os.Remove(assetPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove asset file: %w", err)
	}

	return nil
}

// ListFiles returns all tracked files
func (fm *FileManager) ListFiles() []string {
	files := make([]string, 0, len(fm.fileIndex))
	for file := range fm.fileIndex {
		files = append(files, file)
	}
	return files
}

// FileExists checks if a file is tracked
func (fm *FileManager) FileExists(filename string) bool {
	_, ok := fm.fileIndex[filename]
	return ok
}

// GetFileChunks returns the chunk IDs for a specific file
func (fm *FileManager) GetFileChunks(filename string) ([]string, error) {
	ids, ok := fm.fileIndex[filename]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", filename)
	}
	return ids, nil
}

// GetSupportedExtensions returns a list of supported file extensions
func (fm *FileManager) GetSupportedExtensions() []string {
	exts := make([]string, 0, len(fm.parsers))
	for ext := range fm.parsers {
		exts = append(exts, ext)
	}
	return exts
}

// copyFile copies a file to the asset directory
func (fm *FileManager) copyFile(src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	dst := filepath.Join(fm.assetDir, filepath.Base(src))
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("failed to write asset file: %w", err)
	}

	return nil
}

// simpleTextSplitter is a basic text splitter implementation
type simpleTextSplitter struct {
	chunkSize   int
	overlapSize int
}

// Transform implements document.Transformer
func (s *simpleTextSplitter) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	result := []*schema.Document{}
	
	for _, doc := range docs {
		content := doc.Content
		if len(content) <= s.chunkSize {
			result = append(result, doc)
			continue
		}
		
		// Split into chunks with overlap
		for i := 0; i < len(content); i += s.chunkSize - s.overlapSize {
			end := i + s.chunkSize
			if end > len(content) {
				end = len(content)
			}
			
			chunk := &schema.Document{
				ID:       fmt.Sprintf("%s-chunk-%d", doc.ID, len(result)),
				Content:  content[i:end],
				MetaData: make(map[string]any),
			}
			
			// Copy metadata
			for k, v := range doc.MetaData {
				chunk.MetaData[k] = v
			}
			
			result = append(result, chunk)
			
			if end == len(content) {
				break
			}
		}
	}
	
	return result, nil
}

// GetType returns the transformer type
func (s *simpleTextSplitter) GetType() string {
	return "SimpleTextSplitter"
}

