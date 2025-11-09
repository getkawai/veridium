/*
 * Copyright 2024 CloudWeGo Authors
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
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

const (
	MetaKeyTitle   = "_title"
	MetaKeyDesc    = "_description"
	MetaKeyLang    = "_language"
	MetaKeyCharset = "_charset"
	MetaKeySource  = "_source"
)

// HtmlParserConfig holds configuration for HtmlParser
type HtmlParserConfig struct {
	// CSS selector to extract specific content (e.g., "body", "#content", ".article")
	// If nil, extracts entire document
	Selector *string
}

// HtmlParser parses HTML files using goquery
// Extracts text content and metadata (title, description, language, charset)
type HtmlParser struct {
	selector *string
}

// NewHtmlParser creates a new HTML parser
func NewHtmlParser(ctx context.Context, config *HtmlParserConfig) (*HtmlParser, error) {
	if config == nil {
		config = &HtmlParserConfig{}
	}
	return &HtmlParser{
		selector: config.Selector,
	}, nil
}

// Parse implements the parser.Parser interface
func (p *HtmlParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	commonOpts := parser.GetCommonOptions(nil, opts...)

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	// Select content using CSS selector if provided
	var contentSel *goquery.Selection
	if p.selector != nil {
		contentSel = doc.Find(*p.selector).Contents()
	} else {
		contentSel = doc.Contents()
	}

	// Extract and sanitize content (removes XSS, dangerous HTML)
	sanitized := bluemonday.UGCPolicy().Sanitize(contentSel.Text())
	content := strings.TrimSpace(sanitized)

	// Extract metadata from HTML
	meta := p.extractMetadata(doc)

	// Merge with user-provided metadata
	if commonOpts.ExtraMeta != nil {
		for k, v := range commonOpts.ExtraMeta {
			meta[k] = v
		}
	}

	// Create Eino document with UUID
	docs := []*schema.Document{
		{
			ID:       uuid.New().String(),
			Content:  content,
			MetaData: meta,
		},
	}

	return docs, nil
}

// extractMetadata extracts metadata from HTML document
func (p *HtmlParser) extractMetadata(doc *goquery.Document) map[string]any {
	meta := map[string]any{}

	// Extract title
	if title := doc.Find("title").Text(); title != "" {
		meta[MetaKeyTitle] = strings.TrimSpace(title)
	}

	// Extract description
	if desc := doc.Find("meta[name=description]").AttrOr("content", ""); desc != "" {
		meta[MetaKeyDesc] = desc
	}

	// Extract language
	if lang := doc.Find("html").AttrOr("lang", ""); lang != "" {
		meta[MetaKeyLang] = lang
	}

	// Extract charset
	if charset := doc.Find("meta[charset]").AttrOr("charset", ""); charset != "" {
		meta[MetaKeyCharset] = charset
	}

	return meta
}

// GetType returns the parser type
func (p *HtmlParser) GetType() string {
	return "HtmlParser"
}
