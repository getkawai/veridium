package main

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ZipService provides ZIP file extraction functionality as a Wails service
type ZipService struct{}

// ExtractFiles extracts files from a ZIP archive based on a filter pattern
func (z *ZipService) ExtractFiles(zipPath string, filterPattern string) ([]ExtractedFile, error) {
	// Read the ZIP file using the file system service
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer reader.Close()

	var extractedFiles []ExtractedFile

	// Compile the regex pattern
	regex, err := regexp.Compile(filterPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	for _, file := range reader.File {
		// Check if file matches the filter pattern
		if regex.MatchString(file.Name) {
			// Skip directories
			if strings.HasSuffix(file.Name, "/") {
				continue
			}

			// Open the file
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", file.Name, err)
			}

			// Read the content
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", file.Name, err)
			}

			// Convert content to base64 for JavaScript compatibility
			encodedContent := base64.StdEncoding.EncodeToString(content)

			extractedFiles = append(extractedFiles, ExtractedFile{
				Content: encodedContent,
				Path:    file.Name,
			})
		}
	}

	return extractedFiles, nil
}

// ExtractedFile represents an extracted file from a ZIP archive
type ExtractedFile struct {
	Content string `json:"content"`
	Path    string `json:"path"`
}
