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

/*
Package parsers provides document parsers for various file formats.

This package implements Eino-compatible parsers that extract text content from
different file formats. All parsers implement the eino/components/document/parser.Parser
interface.

Supported Formats:

  - DOCX: Microsoft Word documents (using gooxml)
    • Section-based extraction (headers, main, footers) similar to Eino-Ext
    • Markdown conversion for main content with tables
    • Plain text extraction for headers and footers
    • Preserves document structure and formatting
  - XLSX: Microsoft Excel spreadsheets (using gooxml)
    • Multiple sheet handling
    • Cell value extraction
  - PDF: Portable Document Format (using dslipak/pdf)
    • Plain text extraction
  - HTML: HyperText Markup Language (using golang.org/x/net/html)
    • Structure preservation
    • Script/style tag filtering
  - TXT/MD: Plain text and Markdown files
    • Direct text extraction

Example Usage:

	// DOCX Parser - Converts to markdown with tables
	docxParser, _ := parsers.NewDocxParser(ctx)
	docs, _ := docxParser.Parse(ctx, reader)
	// Output: Markdown formatted text with tables, headings, lists

	// XLSX Parser (sheets as separate documents)
	xlsxParser, _ := parsers.NewXlsxParser(ctx, &parsers.XlsxParserConfig{
	    SheetsAsDocuments: true,
	})
	docs, _ := xlsxParser.Parse(ctx, reader)

	// HTML Parser (preserve structure)
	htmlParser, _ := parsers.NewHtmlParser(ctx, &parsers.HtmlParserConfig{
	    PreserveStructure: true,
	})
	docs, _ := htmlParser.Parse(ctx, reader)

Document ID Generation:

All parsers automatically generate unique UUIDs for each document's ID field,
following Eino-Ext best practices. This ensures proper indexing, deduplication,
and tracking in vector stores.

	doc := &schema.Document{
	    ID:       "550e8400-e29b-41d4-a716-446655440000",  // Auto-generated UUID
	    Content:  "...",
	    MetaData: {...},
	}

All parsers return []*schema.Document which can be further processed by
Eino transformers (splitters) and indexers.
*/
package parsers

