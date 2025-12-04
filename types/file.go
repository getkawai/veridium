package types

import (
	"time"
)

// SupportedFileType represents supported file types
type SupportedFileType string

const (
	FileTypeTXT      SupportedFileType = "txt"
	FileTypePDF      SupportedFileType = "pdf"
	FileTypeDOC      SupportedFileType = "doc"
	FileTypeDOCX     SupportedFileType = "docx"
	FileTypeXLS      SupportedFileType = "xls"
	FileTypeXLSX     SupportedFileType = "xlsx"
	FileTypePPTX     SupportedFileType = "pptx"
	FileTypeImage    SupportedFileType = "image"
	FileTypeMarkdown SupportedFileType = "markdown"
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
	Content        string          `json:"content"`
	CreatedTime    time.Time       `json:"createdTime"`
	FileType       string          `json:"fileType"`
	Filename       string          `json:"filename"`
	Metadata       FileMetadata    `json:"metadata"`
	ModifiedTime   time.Time       `json:"modifiedTime"`
	Pages          []DocumentPage  `json:"pages"`
	Source         string          `json:"source"`
	TotalCharCount int             `json:"totalCharCount"`
	TotalLineCount int             `json:"totalLineCount"`
	Chunks         []DocumentChunk `json:"chunks,omitempty"` // NEW: Pre-chunked content for RAG
}

// ChunkingConfig configures how documents are chunked
type ChunkingConfig struct {
	Enabled     bool // Whether to enable chunking
	ChunkSize   int  // Maximum chunk size in characters
	OverlapSize int  // Overlap between chunks in characters
}

// DocumentChunk represents a chunk of document content
type DocumentChunk = Document

// Document is a piece of text with metadata.
type Document struct {
	// ID is the unique identifier of the document.
	ID string `json:"id"`
	// Content is the content of the document.
	Content string `json:"content"`
	// Metadata is the metadata of the document, can be used to store extra information.
	Metadata map[string]interface{} `json:"metadata"`
}
