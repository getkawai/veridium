package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dslipak/pdf"
	"github.com/kawai-network/veridium/gooxml/document"
	"github.com/kawai-network/veridium/gooxml/presentation"
	"github.com/kawai-network/veridium/gooxml/spreadsheet"
)

// FileLoader provides file loading functionality
type FileLoader struct{}

// NewFileLoader creates a new FileLoader
func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// LoadFile loads a file and returns a FileDocument with markdown content
func (l *FileLoader) LoadFile(filePath string, fileMetadata *FileMetadata) (*FileDocument, error) {
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
func (l *FileLoader) detectFileType(filePath string) (SupportedFileType, error) {
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
		// Check if it's an image
		if l.isImageFile(ext) {
			return FileTypeImage, nil
		}
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

// isTextReadableFile checks if file extension indicates text-readable content
func (l *FileLoader) isTextReadableFile(ext string) bool {
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

// isImageFile checks if file extension indicates an image
func (l *FileLoader) isImageFile(ext string) bool {
	// Remove leading dot
	ext = strings.TrimPrefix(ext, ".")
	ext = strings.ToLower(ext)

	imageExtensions := []string{
		"jpg", "jpeg", "png", "gif", "webp", "svg", "bmp", "tiff",
	}

	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// CanChunkForRAG checks if a MIME type can be chunked for RAG processing
// Returns true if the file type can be parsed into text and chunked
// Returns false for images, videos, audio, and binary formats
func (l *FileLoader) CanChunkForRAG(mimeType string) bool {
	// Media files cannot be chunked into text (except images which we process with VL model)
	if len(mimeType) >= 6 {
		prefix := mimeType[:6]
		if prefix == "video/" || prefix == "audio/" {
			return false
		}
		// Images ARE supported now via VL model description generation
		if prefix == "image/" {
			return true
		}
	}

	// Binary/archive formats cannot be chunked
	unsupportedTypes := []string{
		"application/octet-stream",
		"application/zip",
		"application/x-rar",
		"application/x-7z-compressed",
		"application/x-tar",
		"application/gzip",
		"application/x-bzip2",
		"application/x-xz",
	}

	for _, unsupported := range unsupportedTypes {
		if mimeType == unsupported {
			return false
		}
	}

	// All other types (documents, text files) can be chunked
	return true
}

// loadContent loads content based on file type and converts to markdown
func (l *FileLoader) loadContent(filePath string, fileType SupportedFileType) ([]DocumentPage, string, string, error) {
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
	case FileTypeImage:
		return l.loadImageFile(filePath)
	default:
		return nil, "", "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// loadTextFile loads text files
func (l *FileLoader) loadTextFile(filePath string) ([]DocumentPage, string, string, error) {
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
func (l *FileLoader) loadPDFFile(filePath string) ([]DocumentPage, string, string, error) {
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

// loadDOCXFile loads DOCX files using gooxml/document
func (l *FileLoader) loadDOCXFile(filePath string) ([]DocumentPage, string, string, error) {
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

// extractDOCXContent extracts content from DOCX file using gooxml/document
func (l *FileLoader) extractDOCXContent(filePath string) (string, error) {
	// Open the DOCX document
	doc, err := document.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX document: %w", err)
	}

	// Convert to markdown with images served via URLs
	markdown, err := doc.ToMarkdownWithImageURLs("/files")
	if err != nil {
		return "", fmt.Errorf("failed to convert DOCX to markdown: %w", err)
	}

	if markdown == "" {
		return "# DOCX Document\n\n*No content found in document*", nil
	}

	return markdown, nil
}

// loadExcelFile loads Excel files using gooxml/spreadsheet
func (l *FileLoader) loadExcelFile(filePath string) ([]DocumentPage, string, string, error) {
	markdown, err := l.extractXLSXContent(filePath)
	if err != nil {
		return nil, "", fmt.Sprintf("failed to convert XLSX to markdown: %v", err), err
	}

	// Create a single page for XLSX
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

// extractXLSXContent extracts content from XLSX file using gooxml/spreadsheet
func (l *FileLoader) extractXLSXContent(filePath string) (string, error) {
	// Open the XLSX workbook
	wb, err := spreadsheet.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open XLSX workbook: %w", err)
	}
	defer wb.Close()

	// Convert to markdown with images served via URLs
	markdown, err := wb.ToMarkdownWithImageURLs("/files")
	if err != nil {
		return "", fmt.Errorf("failed to convert XLSX to markdown: %w", err)
	}

	if markdown == "" {
		return "# Excel Workbook\n\n*No content found in workbook*", nil
	}

	return markdown, nil
}

// loadPPTXFile loads PPTX files using gooxml/presentation
func (l *FileLoader) loadPPTXFile(filePath string) ([]DocumentPage, string, string, error) {
	markdown, err := l.extractPPTXContent(filePath)
	if err != nil {
		return nil, "", fmt.Sprintf("failed to convert PPTX to markdown: %v", err), err
	}

	// Create a single page for PPTX
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

// loadImageFile loads image files
func (l *FileLoader) loadImageFile(filePath string) ([]DocumentPage, string, string, error) {
	// For images, we just create a markdown reference
	// The actual description will be generated by the FileProcessorService using VL model
	// and appended to this content
	filename := filepath.Base(filePath)

	// We use the /files/ route which is served by the fileserver service
	markdown := fmt.Sprintf("![%s](/files/%s)\n\n", filename, filename)

	page := DocumentPage{
		CharCount:   len(markdown),
		LineCount:   1,
		Metadata:    map[string]interface{}{"type": "image"},
		PageContent: markdown,
	}

	pages := []DocumentPage{page}

	return pages, markdown, "", nil
}

// extractPPTXContent extracts content from PPTX file using gooxml/presentation
func (l *FileLoader) extractPPTXContent(filePath string) (string, error) {
	// Open the PPTX presentation
	doc, err := presentation.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PPTX presentation: %w", err)
	}

	// Convert to markdown with images served via URLs
	markdown, err := doc.ToMarkdownWithImageURLs("/files")
	if err != nil {
		return "", fmt.Errorf("failed to convert PPTX to markdown: %w", err)
	}

	if markdown == "" {
		return "# PowerPoint Presentation\n\n*No content found in presentation*", nil
	}

	return markdown, nil
}

// createErrorDocument creates a FileDocument with error information
func (l *FileLoader) createErrorDocument(filePath string, fileMetadata *FileMetadata, errorMsg string) (*FileDocument, error) {
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
