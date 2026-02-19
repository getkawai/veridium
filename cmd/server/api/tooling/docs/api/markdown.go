// Package api provides markdown documentation generation for kawai-website/docs.
package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getMarkdownOutputDir returns the absolute path to the markdown output directory
func getMarkdownOutputDir() string {
	// Output to kawai-website/docs-source/api directory
	return "/Users/yuda/github.com/kawai-network/kawai-website/docs-source/api"
}

// RunMarkdown generates markdown files for kawai-website/docs
func RunMarkdown() error {
	docs := []apiDoc{
		chatDoc(),
		messagesDoc(),
		responsesDoc(),
		embeddingsDoc(),
		rerankDoc(),
		toolsDoc(),
		speechDoc(),
		transcriptionsDoc(),
		imageDoc(),
	}

	outputDir := getMarkdownOutputDir()

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Generate overview page
	overviewMD := generateOverviewMarkdown(docs)
	overviewPath := filepath.Join(outputDir, "overview.md")
	if err := os.WriteFile(overviewPath, []byte(overviewMD), 0644); err != nil {
		return fmt.Errorf("writing overview: %w", err)
	}
	fmt.Printf("Generated %s\n", overviewPath)

	// Generate individual API pages
	for _, doc := range docs {
		md := generateMarkdown(doc)

		filename := doc.Filename
		filename = strings.TrimSuffix(filename, ".tsx")
		filename = strings.TrimPrefix(filename, "DocsAPI")
		filename = toKebabCase(filename) + ".md"

		outputPath := filepath.Join(outputDir, filename)
		if err := os.WriteFile(outputPath, []byte(md), 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Printf("Generated %s\n", outputPath)
	}

	return nil
}

func generateOverviewMarkdown(docs []apiDoc) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString("title: \"API Reference\"\n")
	b.WriteString("description: \"Complete API reference for Kawai DeAI Network\"\n")
	b.WriteString("date: \"2026-02-17\"\n")
	b.WriteString("categories: [\"API\"]\n")
	b.WriteString("---\n\n")

	b.WriteString("# API Reference\n\n")
	b.WriteString("Complete API reference for the Kawai DeAI Network. All endpoints are prefixed with `/v1`.\n\n")

	b.WriteString("## Base URL\n\n")
	b.WriteString("```\n")
	b.WriteString("https://api.getkawai.com/v1\n")
	b.WriteString("```\n\n")

	b.WriteString("## Authentication\n\n")
	b.WriteString("When authentication is enabled, include your API token in the Authorization header:\n\n")
	b.WriteString("```\n")
	b.WriteString("Authorization: Bearer API_KEY\n")
	b.WriteString("```\n\n")

	b.WriteString("## Available APIs\n\n")

	for _, doc := range docs {
		filename := toKebabCase(strings.TrimPrefix(doc.Filename, "DocsAPI"))
		filename = strings.TrimSuffix(filename, ".tsx")
		b.WriteString(fmt.Sprintf("### [%s](/docs/api/%s)\n", doc.Name, filename))
		b.WriteString(fmt.Sprintf("%s\n\n", doc.Description))
	}

	b.WriteString("## Quick Reference\n\n")
	b.WriteString("| Endpoint | Method | Description |\n")
	b.WriteString("|----------|--------|-------------|\n")
	b.WriteString("| `/chat/completions` | POST | Chat completions (OpenAI compatible) |\n")
	b.WriteString("| `/messages` | POST | Messages (Anthropic compatible) |\n")
	b.WriteString("| `/responses` | POST | Responses API |\n")
	b.WriteString("| `/embeddings` | POST | Text embeddings |\n")
	b.WriteString("| `/rerank` | POST | Document reranking |\n")
	b.WriteString("| `/images/generations` | POST | Image generation |\n")
	b.WriteString("| `/images/edits` | POST | Image editing |\n")
	b.WriteString("| `/images/variations` | POST | Image variations |\n")
	b.WriteString("| `/audio/transcriptions` | POST | Speech-to-text |\n")
	b.WriteString("| `/audio/translations` | POST | Audio translation |\n")
	b.WriteString("| `/audio/speech` | POST | Text-to-speech |\n")
	b.WriteString("| `/models` | GET | List models |\n")
	b.WriteString("| `/catalog` | GET | Browse model catalog |\n")
	b.WriteString("| `/libs` | GET | Library information |\n")

	return b.String()
}

func generateMarkdown(doc apiDoc) string {
	var b strings.Builder

	// Front matter
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("title: \"%s\"\n", doc.Name))
	b.WriteString(fmt.Sprintf("description: \"%s\"\n", escapeMarkdown(doc.Description)))
	b.WriteString("date: \"2026-02-17\"\n")
	b.WriteString("categories: [\"API\"]\n")
	b.WriteString("---\n\n")

	// Title
	b.WriteString(fmt.Sprintf("# %s\n\n", doc.Name))
	b.WriteString(fmt.Sprintf("%s\n\n", doc.Description))

	// Base URL info
	b.WriteString("## Base URL\n\n")
	b.WriteString("```\n")
	b.WriteString("https://api.getkawai.com/v1\n")
	b.WriteString("```\n\n")

	// Authentication
	b.WriteString("## Authentication\n\n")
	b.WriteString("When authentication is enabled, include your token in the Authorization header:\n\n")
	b.WriteString("```\n")
	b.WriteString("Authorization: Bearer API_KEY\n")
	b.WriteString("```\n\n")

	// Generate content for each group
	for _, group := range doc.Groups {
		b.WriteString(fmt.Sprintf("## %s\n\n", group.Name))
		b.WriteString(fmt.Sprintf("%s\n\n", group.Description))

		for _, ep := range group.Endpoints {
			generateEndpointMarkdown(&b, ep, group.Name)
		}
	}

	return b.String()
}

func generateEndpointMarkdown(b *strings.Builder, ep endpoint, groupName string) {
	// Endpoint header
	if ep.Method != "" {
		b.WriteString(fmt.Sprintf("### `%s %s`\n\n", ep.Method, ep.Path))
	} else {
		b.WriteString(fmt.Sprintf("### %s\n\n", ep.Path))
	}

	b.WriteString(fmt.Sprintf("%s\n\n", ep.Description))

	// Authentication info
	if ep.Auth != "" {
		b.WriteString(fmt.Sprintf("**Authentication:** %s\n\n", ep.Auth))
	}

	// Headers
	if len(ep.Headers) > 0 {
		b.WriteString("#### Headers\n\n")
		b.WriteString("| Header | Required | Description |\n")
		b.WriteString("|--------|----------|-------------|\n")
		for _, h := range ep.Headers {
			required := "No"
			if h.Required {
				required = "Yes"
			}
			b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", h.Name, required, h.Description))
		}
		b.WriteString("\n")
	}

	// Request Body
	if ep.RequestBody != nil && len(ep.RequestBody.Fields) > 0 {
		b.WriteString("#### Request Body\n\n")
		b.WriteString(fmt.Sprintf("Content-Type: `%s`\n\n", ep.RequestBody.ContentType))
		b.WriteString("| Field | Type | Required | Description |\n")
		b.WriteString("|-------|------|----------|-------------|\n")
		for _, f := range ep.RequestBody.Fields {
			required := "No"
			if f.Required {
				required = "Yes"
			}
			b.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | %s |\n", f.Name, f.Type, required, f.Description))
		}
		b.WriteString("\n")
	}

	// Response
	if ep.Response != nil {
		b.WriteString("#### Response\n\n")
		b.WriteString(fmt.Sprintf("%s\n\n", ep.Response.Description))
		if ep.Response.ContentType != "" {
			b.WriteString(fmt.Sprintf("Content-Type: `%s`\n\n", ep.Response.ContentType))
		}
	}

	// Examples
	if len(ep.Examples) > 0 {
		b.WriteString("#### Examples\n\n")
		for i, ex := range ep.Examples {
			if ex.Description != "" {
				b.WriteString(fmt.Sprintf("%s\n\n", ex.Description))
			}
			// Detect language from content
			lang := detectLanguage(ex.Code)
			b.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, ex.Code))
			if i < len(ep.Examples)-1 {
				b.WriteString("---\n\n")
			}
		}
	}

	b.WriteString("---\n\n")
}

func detectLanguage(code string) string {
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "curl") {
		return "bash"
	}
	if strings.HasPrefix(code, "{") || strings.HasPrefix(code, "[") {
		return "json"
	}
	if strings.HasPrefix(code, "event:") {
		return "text"
	}
	return ""
}

func toKebabCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && isUpper(r) && !isUpper(rune(s[i-1])) {
			result = append(result, '-')
		}
		result = append(result, []rune(strings.ToLower(string(r)))...)
	}
	return string(result)
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func escapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
