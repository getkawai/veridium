package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// LocalFileStorage handles local file storage operations
type LocalFileStorage struct {
	BaseDir string // Base file directory
}

// NewLocalFileStorage creates new local file storage
func NewLocalFileStorage(baseDir string) *LocalFileStorage {
	return &LocalFileStorage{BaseDir: baseDir}
}

// CreatePreSignedUrl creates pre-signed upload URL (for local files, just returns path)
func (s *LocalFileStorage) CreatePreSignedUrl(key string) string {
	return filepath.Join(s.BaseDir, key)
}

// CreatePreSignedUrlForPreview creates pre-signed preview URL
func (s *LocalFileStorage) CreatePreSignedUrlForPreview(key string) string {
	return filepath.Join(s.BaseDir, key)
}

// DeleteFile deletes a single file
func (s *LocalFileStorage) DeleteFile(key string) error {
	return s.DeleteFiles([]string{key})
}

// DeleteFiles batch deletes files
func (s *LocalFileStorage) DeleteFiles(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	var errorsEncountered []error
	for _, key := range keys {
		if !s.isValidKey(key) {
			errorsEncountered = append(errorsEncountered, fmt.Errorf("invalid local file path: %s", key))
			continue
		}

		localPath := filepath.Join(s.BaseDir, key)
		if err := os.Remove(localPath); err != nil {
			log.Printf("Failed to delete file %s: %v", localPath, err)
			errorsEncountered = append(errorsEncountered, err)
		}
	}

	if len(errorsEncountered) > 0 {
		return fmt.Errorf("failed to delete %d files: %v", len(errorsEncountered), errorsEncountered[0])
	}

	return nil
}

// isValidKey validates if key is valid (prevent path traversal attack)
func (s *LocalFileStorage) isValidKey(key string) bool {
	cleanKey := filepath.Clean(key)
	if strings.Contains(cleanKey, "..") {
		return false
	}
	return filepath.IsLocal(key)
}

// GetFileByteArray gets file byte array
func (s *LocalFileStorage) GetFileByteArray(key string) ([]byte, error) {
	localPath := filepath.Join(s.BaseDir, key)
	return os.ReadFile(localPath)
}

// GetFileContent gets file content as string
func (s *LocalFileStorage) GetFileContent(key string) (string, error) {
	localPath := filepath.Join(s.BaseDir, key)
	content, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetFullFileUrl gets complete file URL (for local files, returns the key)
func (s *LocalFileStorage) GetFullFileUrl(fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	return s.GetKeyFromFullUrl(fileUrl)
}

// GetKeyFromFullUrl extracts key from complete URL
func (s *LocalFileStorage) GetKeyFromFullUrl(fullUrl string) string {
	// If it's already a file path, return directly
	if !strings.Contains(fullUrl, "://") {
		return fullUrl
	}

	// Parse URL
	parsedURL, err := url.Parse(fullUrl)
	if err != nil {
		log.Printf("Failed to parse URL %s: %v", fullUrl, err)
		return ""
	}

	// Extract path part and remove beginning slash
	filePath := parsedURL.Path
	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	return filePath
}

// UploadContent uploads content (not implemented)
func (s *LocalFileStorage) UploadContent(filePath, content string) error {
	log.Printf("UploadContent not implemented for path: %s", filePath)
	return errors.New("uploadContent not implemented")
}

// ReadFileFromAbsolutePath reads file from absolute path (for drag & drop)
func (s *LocalFileStorage) ReadFileFromAbsolutePath(absolutePath string) ([]byte, error) {
	// Security check: ensure path is absolute and exists
	if !filepath.IsAbs(absolutePath) {
		return nil, fmt.Errorf("path must be absolute: %s", absolutePath)
	}

	// Check if file exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", absolutePath)
	}

	return os.ReadFile(absolutePath)
}

// UploadMedia uploads media file
// buffer parameter can be either:
// - []byte: raw binary data
// - string: base64-encoded data (for Wails binding compatibility)
func (s *LocalFileStorage) UploadMedia(key string, buffer interface{}) (string, error) {
	filename := filepath.Base(key)

	// Convert buffer to []byte
	var data []byte
	switch v := buffer.(type) {
	case []byte:
		data = v
	case string:
		// Decode base64 string
		decoded, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 data: %w", err)
		}
		data = decoded
	default:
		return "", fmt.Errorf("unsupported buffer type: %T", buffer)
	}

	// Calculate file's SHA256 hash
	hash := sha256.Sum256(data)
	hashString := fmt.Sprintf("%x", hash)

	// Construct local file path
	localPath := filepath.Join(s.BaseDir, key)

	// Ensure directory exists
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", localPath, err)
	}

	log.Printf("File uploaded successfully: %s (size: %d bytes, hash: %s)", filename, len(data), hashString)
	return key, nil
}
