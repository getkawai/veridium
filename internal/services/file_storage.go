package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// CopyFileFromAbsolutePath copies file from absolute path to local storage (for drag & drop)
// Returns the relative key path that can be used to access the file via fileserver
func (s *LocalFileStorage) CopyFileFromAbsolutePath(absolutePath string) (string, error) {
	// Security check: ensure path is absolute and exists
	if !filepath.IsAbs(absolutePath) {
		return "", fmt.Errorf("path must be absolute: %s", absolutePath)
	}

	// Check if file exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", absolutePath)
	}

	// Extract filename
	filename := filepath.Base(absolutePath)

	// Generate unique filename with timestamp
	timestamp := time.Now().UnixMilli()
	uniqueFileName := fmt.Sprintf("%d-%s", timestamp, filename)
	relativeKey := filepath.Join("uploads", uniqueFileName)

	// Construct destination path
	destPath := filepath.Join(s.BaseDir, relativeKey)

	// Ensure directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Copy file
	sourceFile, err := os.Open(absolutePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy content
	written, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	log.Printf("File copied successfully: %s -> %s (%d bytes)", filename, relativeKey, written)

	return relativeKey, nil
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

// FileItem represents a file database record
type FileItem struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int    `json:"size,omitempty"`
}

// FileModelInterface defines file model operations
type FileModelInterface interface {
	FindById(fileId string) (*FileItem, error)
	Delete(fileId string, globalFile bool) error
}

// FileService wraps LocalFileStorage to provide high-level file operations
type FileService struct {
	UserId    string
	FileModel FileModelInterface
	Storage   *LocalFileStorage
}

// NewFileService creates new file service
func NewFileService(userId string, storage *LocalFileStorage) *FileService {
	return &FileService{
		UserId:  userId,
		Storage: storage,
	}
}

// DeleteFile deletes file
func (fs *FileService) DeleteFile(key string) error {
	return fs.Storage.DeleteFile(key)
}

// DeleteFiles batch deletes files
func (fs *FileService) DeleteFiles(keys []string) error {
	return fs.Storage.DeleteFiles(keys)
}

// GetFileContent gets file content
func (fs *FileService) GetFileContent(key string) (string, error) {
	return fs.Storage.GetFileContent(key)
}

// GetFileByteArray gets file byte array
func (fs *FileService) GetFileByteArray(key string) ([]byte, error) {
	return fs.Storage.GetFileByteArray(key)
}

// CreatePreSignedUrl creates pre-signed upload URL
func (fs *FileService) CreatePreSignedUrl(key string) string {
	return fs.Storage.CreatePreSignedUrl(key)
}

// CreatePreSignedUrlForPreview creates pre-signed preview URL
func (fs *FileService) CreatePreSignedUrlForPreview(key string) string {
	return fs.Storage.CreatePreSignedUrlForPreview(key)
}

// UploadContent uploads content
func (fs *FileService) UploadContent(path, content string) error {
	return fs.Storage.UploadContent(path, content)
}

// GetFullFileUrl gets complete file URL
func (fs *FileService) GetFullFileUrl(url string) string {
	return fs.Storage.GetFullFileUrl(url)
}

// GetKeyFromFullUrl extracts key from full URL
func (fs *FileService) GetKeyFromFullUrl(url string) string {
	return fs.Storage.GetKeyFromFullUrl(url)
}

// UploadMedia uploads media file
func (fs *FileService) UploadMedia(key string, buffer interface{}) (string, error) {
	return fs.Storage.UploadMedia(key, buffer)
}

// ReadFileFromAbsolutePath reads file from absolute path (for drag & drop)
func (fs *FileService) ReadFileFromAbsolutePath(absolutePath string) ([]byte, error) {
	return fs.Storage.ReadFileFromAbsolutePath(absolutePath)
}

// CopyFileFromAbsolutePath copies file from absolute path to local storage (for drag & drop)
func (fs *FileService) CopyFileFromAbsolutePath(absolutePath string) (string, error) {
	return fs.Storage.CopyFileFromAbsolutePath(absolutePath)
}

// DownloadFileToLocal downloads file to local temp storage
func (fs *FileService) DownloadFileToLocal(fileId string) (cleanup func(), file *FileItem, filePath string, err error) {
	// Find file by ID
	file, err = fs.FileModel.FindById(fileId)
	if err != nil {
		return nil, nil, "", fmt.Errorf("file not found: %w", err)
	}
	if file == nil {
		return nil, nil, "", errors.New("file not found")
	}

	// Get file content bytes
	content, err := fs.GetFileByteArray(file.URL)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get file content: %w", err)
	}

	if len(content) == 0 {
		return nil, nil, "", errors.New("file content is empty")
	}

	// Convert to base64 string for Wails temp file binding
	dataStr := base64.StdEncoding.EncodeToString(content)
	filePath, err = WriteTempFile(dataStr, file.Name)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create temp file: %w", err)
	}

	// Return cleanup function, file info, and temp file path
	cleanup = func() {
		Cleanup()
	}

	return cleanup, file, filePath, nil
}

// WriteTempFile writes data to temp file (placeholder)
func WriteTempFile(data, name string) (string, error) {
	log.Printf("WriteTempFile not implemented: data length %d, name %s", len(data), name)
	// In actual Wails app, this would call bindings/github.com/kawai-network/veridium/tempfileservice.WriteTempFile
	return "/tmp/placeholder-" + name, nil
}

// Cleanup cleans up temp files (placeholder)
func Cleanup() {
	log.Println("Cleanup not implemented")
	// In actual Wails app, this would call bindings/github.com/kawai-network/veridium/tempfileservice.Cleanup
}
