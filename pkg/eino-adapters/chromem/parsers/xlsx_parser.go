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
	"github.com/kawai-network/veridium/gooxml/spreadsheet"
)

// XlsxParser parses XLSX files using gooxml
type XlsxParser struct {
	// If true, each sheet becomes a separate document
	SheetsAsDocuments bool
}

// XlsxParserConfig holds configuration for XlsxParser
type XlsxParserConfig struct {
	SheetsAsDocuments bool
}

// NewXlsxParser creates a new XLSX parser
func NewXlsxParser(ctx context.Context, config *XlsxParserConfig) (*XlsxParser, error) {
	if config == nil {
		config = &XlsxParserConfig{}
	}
	return &XlsxParser{
		SheetsAsDocuments: config.SheetsAsDocuments,
	}, nil
}

// Parse implements the parser.Parser interface
func (p *XlsxParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	commonOpts := parser.GetCommonOptions(nil, opts...)

	// Save reader to temp file
	tmpFile, err := os.CreateTemp("", "xlsx-*.xlsx")
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
	wb, err := spreadsheet.Open(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open xlsx: %w", err)
	}

	docs := []*schema.Document{}

	if p.SheetsAsDocuments {
		// Each sheet as separate document
		for _, sheet := range wb.Sheets() {
			content := p.extractSheetText(sheet)

			metadata := make(map[string]any)
			for k, v := range commonOpts.ExtraMeta {
				metadata[k] = v
			}
			metadata["sheet_name"] = sheet.Name()

			docs = append(docs, &schema.Document{
				ID:       uuid.New().String(),
				Content:  content,
				MetaData: metadata,
			})
		}
	} else {
		// All sheets in one document
		var content strings.Builder
		for _, sheet := range wb.Sheets() {
			content.WriteString(fmt.Sprintf("=== Sheet: %s ===\n", sheet.Name()))
			content.WriteString(p.extractSheetText(sheet))
			content.WriteString("\n\n")
		}

		docs = append(docs, &schema.Document{
			ID:       uuid.New().String(),
			Content:  strings.TrimSpace(content.String()),
			MetaData: commonOpts.ExtraMeta,
		})
	}

	return docs, nil
}

// extractSheetText extracts text from a sheet
func (p *XlsxParser) extractSheetText(sheet spreadsheet.Sheet) string {
	var content strings.Builder

	for _, row := range sheet.Rows() {
		rowContent := []string{}
		for _, cell := range row.Cells() {
			// Get cell value as string
			value := cell.GetString()
			if value != "" {
				rowContent = append(rowContent, value)
			}
		}
		if len(rowContent) > 0 {
			content.WriteString(strings.Join(rowContent, "\t"))
			content.WriteString("\n")
		}
	}

	return strings.TrimSpace(content.String())
}

// GetType returns the parser type
func (p *XlsxParser) GetType() string {
	return "XlsxParser"
}
