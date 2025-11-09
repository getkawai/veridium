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
	"strings"

	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

// TextParser parses plain text files (TXT, MD, etc.)
type TextParser struct{}

// NewTextParser creates a new text parser
func NewTextParser(ctx context.Context) (*TextParser, error) {
	return &TextParser{}, nil
}

// Parse implements the parser.Parser interface
func (p *TextParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	commonOpts := parser.GetCommonOptions(nil, opts...)

	// Read all content
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read text: %w", err)
	}

	// Create Eino document
	docs := []*schema.Document{
		{
			ID:       uuid.New().String(),
			Content:  strings.TrimSpace(string(data)),
			MetaData: commonOpts.ExtraMeta,
		},
	}

	return docs, nil
}

// GetType returns the parser type
func (p *TextParser) GetType() string {
	return "TextParser"
}

