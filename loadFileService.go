package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dslipak/pdf"
	"github.com/eino-contrib/docx2md"
	"github.com/unidoc/unioffice/presentation"
	"github.com/xuri/excelize/v2"
)

// LoadFileService provides file loading functionality as a Wails service
type LoadFileService struct{}

// SupportedFileType represents supported file types
type SupportedFileType string

const (
	FileTypeTXT  SupportedFileType = "txt"
	FileTypePDF  SupportedFileType = "pdf"
	FileTypeDOC  SupportedFileType = "doc"
	FileTypeDOCX SupportedFileType = "docx"
	FileTypeXLS  SupportedFileType = "xls"
	FileTypeXLSX SupportedFileType = "xlsx"
	FileTypePPTX SupportedFileType = "pptx"
)

// FileMetadata represents file metadata
type FileMetadata struct {
	Source       string    `json:"source,omitempty"`
	Filename     string    `json:"filename,omitempty"`
	FileType     string    `json:"fileType,omitempty"`
	CreatedTime  time.Time `json:"createdTime,omitempty"`
	ModifiedTime time.Time `json:"modifiedTime,omitempty"`
	Error        string    `json:"error,omitempty"`
}

// DocumentPage represents a page in a document
type DocumentPage struct {
	CharCount   int                    `json:"charCount"`
	LineCount   int                    `json:"lineCount"`
	Metadata    map[string]interface{} `json:"metadata"`
	PageContent string                 `json:"pageContent"`
}

// FileDocument represents the loaded file document
type FileDocument struct {
	Content        string         `json:"content"`
	CreatedTime    time.Time      `json:"createdTime"`
	FileType       string         `json:"fileType"`
	Filename       string         `json:"filename"`
	Metadata       FileMetadata   `json:"metadata"`
	ModifiedTime   time.Time      `json:"modifiedTime"`
	Pages          []DocumentPage `json:"pages"`
	Source         string         `json:"source"`
	TotalCharCount int            `json:"totalCharCount"`
	TotalLineCount int            `json:"totalLineCount"`
}

// LoadFile loads a file and returns a FileDocument with markdown content
func (l *LoadFileService) LoadFile(filePath string, fileMetadata *FileMetadata) (*FileDocument, error) {
	// Get file stats
	stats, err := os.Stat(filePath)
	if err != nil {
		return l.createErrorDocument(filePath, fileMetadata, fmt.Sprintf("Failed to access file stats: %v", err))
	}

	// Determine file type
	fileType, err := l.detectFileType(filePath)
	if err != nil {
		return l.createErrorDocument(filePath, fileMetadata, fmt.Sprintf("Failed to detect file type: %v", err))
	}

	// Get base file info
	ext := filepath.Ext(filePath)
	baseFilename := filepath.Base(filePath)
	fileExtension := strings.ToLower(strings.TrimPrefix(ext, "."))

	// Apply metadata overrides or use defaults
	source := filePath
	filename := baseFilename
	if fileMetadata != nil {
		if fileMetadata.Source != "" {
			source = fileMetadata.Source
		}
		if fileMetadata.Filename != "" {
			filename = fileMetadata.Filename
		}
		if fileMetadata.FileType != "" {
			fileExtension = fileMetadata.FileType
		}
	}

	createdTime := stats.ModTime()
	modifiedTime := stats.ModTime()
	if fileMetadata != nil {
		if !fileMetadata.CreatedTime.IsZero() {
			createdTime = fileMetadata.CreatedTime
		}
		if !fileMetadata.ModifiedTime.IsZero() {
			modifiedTime = fileMetadata.ModifiedTime
		}
	}

	// Load content based on file type
	pages, aggregatedContent, loaderError, err := l.loadContent(filePath, SupportedFileType(fileType))
	if err != nil {
		return l.createErrorDocument(filePath, fileMetadata, fmt.Sprintf("Failed to load content: %v", err))
	}

	// Calculate totals
	totalCharCount := 0
	totalLineCount := 0
	for _, page := range pages {
		totalCharCount += page.CharCount
		totalLineCount += page.LineCount
	}

	// Create metadata
	metadata := FileMetadata{
		Source:       source,
		Filename:     filename,
		FileType:     fileExtension,
		CreatedTime:  createdTime,
		ModifiedTime: modifiedTime,
	}
	if loaderError != "" {
		metadata.Error = loaderError
	}

	// Create document
	doc := &FileDocument{
		Content:        aggregatedContent,
		CreatedTime:    createdTime,
		FileType:       fileExtension,
		Filename:       filename,
		Metadata:       metadata,
		ModifiedTime:   modifiedTime,
		Pages:          pages,
		Source:         source,
		TotalCharCount: totalCharCount,
		TotalLineCount: totalLineCount,
	}

	return doc, nil
}

// detectFileType determines file type based on extension
func (l *LoadFileService) detectFileType(filePath string) (SupportedFileType, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return FileTypePDF, nil
	case ".doc":
		return FileTypeDOC, nil
	case ".docx":
		return FileTypeDOCX, nil
	case ".xlsx", ".xls":
		return FileTypeXLSX, nil
	case ".pptx":
		return FileTypePPTX, nil
	case ".txt", "":
		return FileTypeTXT, nil
	default:
		// Check if it's text readable
		if l.isTextReadableFile(ext) {
			return FileTypeTXT, nil
		}
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

// isTextReadableFile checks if file extension indicates text-readable content
func (l *LoadFileService) isTextReadableFile(ext string) bool {
	// Remove leading dot
	ext = strings.TrimPrefix(ext, ".")

	textExtensions := []string{
		"txt", "md", "markdown", "json", "xml", "html", "htm", "css", "js", "ts",
		"py", "java", "cpp", "c", "h", "hpp", "cs", "php", "rb", "go", "rs", "sh",
		"yml", "yaml", "toml", "ini", "cfg", "conf", "log", "csv", "tsv",
	}

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

// loadContent loads content based on file type and converts to markdown
func (l *LoadFileService) loadContent(filePath string, fileType SupportedFileType) ([]DocumentPage, string, string, error) {
	switch fileType {
	case FileTypeTXT:
		return l.loadTextFile(filePath)
	case FileTypePDF:
		return l.loadPDFFile(filePath)
	case FileTypeDOCX:
		return l.loadDOCXFile(filePath)
	case FileTypeXLSX:
		return l.loadExcelFile(filePath)
	case FileTypePPTX:
		return l.loadPPTXFile(filePath)
	default:
		return nil, "", "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// loadTextFile loads text files
func (l *LoadFileService) loadTextFile(filePath string) ([]DocumentPage, string, string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Sprintf("Failed to read text file: %v", err), err
	}

	textContent := string(content)
	lines := strings.Split(textContent, "\n")
	lineCount := len(lines)
	charCount := len(textContent)

	page := DocumentPage{
		CharCount:   charCount,
		LineCount:   lineCount,
		Metadata:    map[string]interface{}{"lineNumberEnd": lineCount, "lineNumberStart": 1},
		PageContent: textContent,
	}

	pages := []DocumentPage{page}
	// For text files, content is already in readable format, wrap in markdown code block
	aggregatedContent := fmt.Sprintf("```\n%s\n```", textContent)

	return pages, aggregatedContent, "", nil
}

// loadPDFFile loads PDF files using github.com/dslipak/pdf
func (l *LoadFileService) loadPDFFile(filePath string) ([]DocumentPage, string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Sprintf("Failed to open PDF file: %v", err), err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, "", fmt.Sprintf("Failed to get file stat: %v", err), err
	}

	reader, err := pdf.NewReader(file, stat.Size())
	if err != nil {
		return nil, "", fmt.Sprintf("Failed to create PDF reader: %v", err), err
	}

	var pages []DocumentPage
	var markdownContent strings.Builder

	markdownContent.WriteString("# PDF Document\n\n")

	numPages := reader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		// Try to extract text from the page
		var pageContent string
		content := page.Content()
		if len(content.Text) > 0 {
			// Convert []pdf.Text to string
			var textBuilder strings.Builder
			for _, textItem := range content.Text {
				textBuilder.WriteString(textItem.S)
			}
			pageContent = textBuilder.String()
		} else {
			pageContent = "[Unable to extract text from this page]"
		}

		lines := strings.Split(pageContent, "\n")
		charCount := len(pageContent)
		lineCount := len(lines)

		docPage := DocumentPage{
			CharCount:   charCount,
			LineCount:   lineCount,
			Metadata:    map[string]interface{}{"pageNumber": i},
			PageContent: pageContent,
		}
		pages = append(pages, docPage)

		markdownContent.WriteString(fmt.Sprintf("## Page %d\n\n%s\n\n", i, pageContent))
	}

	return pages, markdownContent.String(), "", nil
}

// loadDOCXFile loads DOCX files using github.com/eino-contrib/docx2md
func (l *LoadFileService) loadDOCXFile(filePath string) ([]DocumentPage, string, string, error) {
	// For now, use a simple text extraction approach
	// The docx2md library might have different API
	markdown, err := l.extractDOCXContent(filePath)
	if err != nil {
		return nil, "", fmt.Sprintf("failed to convert DOCX to markdown: %v", err), err
	}

	// Create a single page for DOCX
	lines := strings.Split(markdown, "\n")
	charCount := len(markdown)
	lineCount := len(lines)

	page := DocumentPage{
		CharCount:   charCount,
		LineCount:   lineCount,
		Metadata:    map[string]interface{}{},
		PageContent: markdown,
	}

	pages := []DocumentPage{page}

	return pages, markdown, "", nil
}

// extractDOCXContent extracts content from DOCX file using docx2md
func (l *LoadFileService) extractDOCXContent(filePath string) (string, error) {
	// Configure docx2md to include headers, footers and tables
	config := &docx2md.Config{
		IncludeHeaders: true,
		IncludeFooters: true,
		IncludeTables:  true,
	}

	// Convert DOCX to markdown sections
	sections, err := docx2md.DocxConvert(filePath, config)
	if err != nil {
		return "", fmt.Errorf("failed to convert DOCX to markdown: %w", err)
	}

	// Combine all sections into one markdown document
	var contentBuilder strings.Builder

	// Section titles mapping (similar to CloudWeGo example)
	sectionTitles := map[string]string{
		"main": "MAIN CONTENT",
	}

	getSectionTitle := func(key string) string {
		if title, ok := sectionTitles[key]; ok {
			return title
		}
		return strings.ToUpper(key)
	}

	// Build content from sections
	for key, section := range sections {
		trimmed := strings.TrimSpace(section)
		if trimmed != "" {
			contentBuilder.WriteString(fmt.Sprintf("=== %s ===\n\n", getSectionTitle(key)))
			contentBuilder.WriteString(trimmed)
			contentBuilder.WriteString("\n\n")
		}
	}

	content := contentBuilder.String()
	if content == "" {
		return "# DOCX Document\n\n*No content found in document*", nil
	}

	return content, nil
}

// loadExcelFile loads Excel files using github.com/xuri/excelize/v2
func (l *LoadFileService) loadExcelFile(filePath string) ([]DocumentPage, string, string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, "", fmt.Sprintf("failed to open Excel file: %v", err), err
	}
	defer f.Close()

	var pages []DocumentPage
	var markdownContent strings.Builder

	markdownContent.WriteString("# Excel Document\n\n")

	sheetList := f.GetSheetList()
	for _, sheetName := range sheetList {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		markdownContent.WriteString(fmt.Sprintf("## Sheet: %s\n\n", sheetName))

		if len(rows) > 0 {
			// Create table header
			if len(rows[0]) > 0 {
				markdownContent.WriteString("| ")
				for i, header := range rows[0] {
					markdownContent.WriteString(header)
					if i < len(rows[0])-1 {
						markdownContent.WriteString(" | ")
					}
				}
				markdownContent.WriteString(" |\n| ")
				for i := range rows[0] {
					markdownContent.WriteString("---")
					if i < len(rows[0])-1 {
						markdownContent.WriteString(" | ")
					}
				}
				markdownContent.WriteString(" |\n")
			}

			// Create table rows
			for i, row := range rows {
				if i == 0 { // Skip header
					continue
				}
				markdownContent.WriteString("| ")
				for j, cell := range row {
					markdownContent.WriteString(cell)
					if j < len(row)-1 {
						markdownContent.WriteString(" | ")
					}
				}
				markdownContent.WriteString(" |\n")
			}
		}

		markdownContent.WriteString("\n")

		// Create a page for this sheet
		sheetContent := l.rowsToString(rows)
		lines := strings.Split(sheetContent, "\n")
		charCount := len(sheetContent)
		lineCount := len(lines)

		page := DocumentPage{
			CharCount:   charCount,
			LineCount:   lineCount,
			Metadata:    map[string]interface{}{"sheetName": sheetName},
			PageContent: sheetContent,
		}
		pages = append(pages, page)
	}

	return pages, markdownContent.String(), "", nil
}

// loadPPTXFile loads PPTX files using github.com/unidoc/unioffice
func (l *LoadFileService) loadPPTXFile(filePath string) ([]DocumentPage, string, string, error) {
	doc, err := presentation.Open(filePath)
	if err != nil {
		return nil, "", fmt.Sprintf("failed to open PPTX file: %v", err), err
	}
	defer doc.Close()

	var pages []DocumentPage
	var markdownContent strings.Builder

	markdownContent.WriteString("# PowerPoint Presentation\n\n")

	slides := doc.Slides()
	for i, slide := range slides {
		// Extract text from slide using unioffice API
		slideText := slide.ExtractText().Text()

		// Clean up the text
		slideText = strings.TrimSpace(slideText)
		if slideText == "" {
			slideText = "[Empty slide]"
		}

		lines := strings.Split(slideText, "\n")
		charCount := len(slideText)
		lineCount := len(lines)

		page := DocumentPage{
			CharCount:   charCount,
			LineCount:   lineCount,
			Metadata:    map[string]interface{}{"slideNumber": i + 1},
			PageContent: slideText,
		}
		pages = append(pages, page)

		markdownContent.WriteString(fmt.Sprintf("## Slide %d\n\n%s\n\n", i+1, slideText))
	}

	return pages, markdownContent.String(), "", nil
}

// rowsToString converts Excel rows to string representation
func (l *LoadFileService) rowsToString(rows [][]string) string {
	var content strings.Builder
	for i, row := range rows {
		for j, cell := range row {
			content.WriteString(cell)
			if j < len(row)-1 {
				content.WriteString("\t")
			}
		}
		if i < len(rows)-1 {
			content.WriteString("\n")
		}
	}
	return content.String()
}

// createErrorDocument creates a FileDocument with error information
func (l *LoadFileService) createErrorDocument(filePath string, fileMetadata *FileMetadata, errorMsg string) (*FileDocument, error) {
	baseFilename := filepath.Base(filePath)

	filename := baseFilename
	source := filePath
	if fileMetadata != nil {
		if fileMetadata.Filename != "" {
			filename = fileMetadata.Filename
		}
		if fileMetadata.Source != "" {
			source = fileMetadata.Source
		}
	}

	errorPage := DocumentPage{
		CharCount:   0,
		LineCount:   0,
		Metadata:    map[string]interface{}{"error": errorMsg},
		PageContent: "",
	}

	doc := &FileDocument{
		Content:     "",
		CreatedTime: time.Now(),
		FileType:    "",
		Filename:    filename,
		Metadata: FileMetadata{
			Source: source,
			Error:  errorMsg,
		},
		ModifiedTime:   time.Now(),
		Pages:          []DocumentPage{errorPage},
		Source:         source,
		TotalCharCount: 0,
		TotalLineCount: 0,
	}

	return doc, nil
}
