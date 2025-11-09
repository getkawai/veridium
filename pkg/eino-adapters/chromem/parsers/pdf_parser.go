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
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/dslipak/pdf"
	"github.com/google/uuid"
)

// PdfParser parses PDF files using dslipak/pdf
type PdfParser struct{}

// NewPdfParser creates a new PDF parser
func NewPdfParser(ctx context.Context) (*PdfParser, error) {
	return &PdfParser{}, nil
}

// Parse implements the parser.Parser interface
func (p *PdfParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	commonOpts := parser.GetCommonOptions(nil, opts...)

	// Read all data
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	// Create in-memory reader (no temp file needed!)
	readerAt := bytes.NewReader(data)
	f, err := pdf.NewReader(readerAt, int64(readerAt.Len()))
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}

	// Extract text from all pages with font caching
	var buf bytes.Buffer
	fonts := make(map[string]*pdf.Font)

	for i := 1; i <= f.NumPage(); i++ {
		page := f.Page(i)

		// Cache fonts to avoid repeated parsing
		for _, name := range page.Fonts() {
			if _, ok := fonts[name]; !ok {
				font := page.Font(name)
				fonts[name] = &font
			}
		}

		// Extract text from page
		text, err := page.GetPlainText(fonts)
		if err != nil {
			return nil, fmt.Errorf("failed to extract page %d: %w", i, err)
		}

		buf.WriteString(text)
		buf.WriteString("\n")
	}

	// Create Eino document
	docs := []*schema.Document{
		{
			ID:       uuid.New().String(),
			Content:  strings.TrimSpace(buf.String()),
			MetaData: commonOpts.ExtraMeta,
		},
	}

	return docs, nil
}

// GetType returns the parser type
func (p *PdfParser) GetType() string {
	return "PdfParser"
}
