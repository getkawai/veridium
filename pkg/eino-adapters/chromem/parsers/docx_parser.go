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

package parsers

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/kawai-network/veridium/gooxml/document"
)

// DocxParser parses DOCX files using gooxml
// It converts documents to markdown format with full table support
type DocxParser struct{}

// NewDocxParser creates a new DOCX parser
func NewDocxParser(ctx context.Context) (*DocxParser, error) {
	return &DocxParser{}, nil
}

// Parse implements the parser.Parser interface
func (p *DocxParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	commonOpts := parser.GetCommonOptions(nil, opts...)

	// Save reader to temp file (gooxml needs file path)
	tmpFile, err := os.CreateTemp("", "docx-*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Copy reader to temp file
	if _, err := io.Copy(tmpFile, reader); err != nil {
		return nil, fmt.Errorf("failed to copy to temp file: %w", err)
	}

	// Close temp file before opening with gooxml
	tmpFile.Close()

	// Open with gooxml
	doc, err := document.Open(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open docx: %w", err)
	}

	// Extract sections similar to docx2md approach
	sections := make(map[string]string)

	// 1. Extract headers
	var headerBuilder strings.Builder
	for _, header := range doc.Headers() {
		for _, para := range header.Paragraphs() {
			for _, run := range para.Runs() {
				headerBuilder.WriteString(run.Text())
			}
			headerBuilder.WriteString("\n")
		}
	}
	if trimmed := strings.TrimSpace(headerBuilder.String()); trimmed != "" {
		sections["headers"] = trimmed
	}

	// 2. Extract main content (body) - convert to markdown
	mainContent := doc.ToMarkdown()
	if trimmed := strings.TrimSpace(mainContent); trimmed != "" {
		sections["main"] = trimmed
	}

	// 3. Extract footers
	var footerBuilder strings.Builder
	for _, footer := range doc.Footers() {
		for _, para := range footer.Paragraphs() {
			for _, run := range para.Runs() {
				footerBuilder.WriteString(run.Text())
			}
			footerBuilder.WriteString("\n")
		}
	}
	if trimmed := strings.TrimSpace(footerBuilder.String()); trimmed != "" {
		sections["footers"] = trimmed
	}

	// Build final content with section markers (similar to Eino-Ext)
	var contentBuilder strings.Builder

	// Add sections in order: headers, main, footers
	if headers, ok := sections["headers"]; ok {
		contentBuilder.WriteString("=== HEADERS ===\n")
		contentBuilder.WriteString(headers)
		contentBuilder.WriteString("\n\n")
	}

	if main, ok := sections["main"]; ok {
		contentBuilder.WriteString("=== MAIN CONTENT ===\n")
		contentBuilder.WriteString(main)
		contentBuilder.WriteString("\n\n")
	}

	if footers, ok := sections["footers"]; ok {
		contentBuilder.WriteString("=== FOOTERS ===\n")
		contentBuilder.WriteString(footers)
		contentBuilder.WriteString("\n")
	}

	// Create Eino document
	finalContent := strings.TrimSpace(contentBuilder.String())
	if finalContent == "" {
		return []*schema.Document{}, nil
	}

	docs := []*schema.Document{
		{
			ID:       uuid.New().String(),
			Content:  finalContent,
			MetaData: commonOpts.ExtraMeta,
		},
	}

	return docs, nil
}

// GetType returns the parser type
func (p *DocxParser) GetType() string {
	return "DocxParser"
}
