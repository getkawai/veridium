package file

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// FileServiceImpl defines file service interface
type FileServiceImpl interface {
	CreatePreSignedUrl(key string) string
	CreatePreSignedUrlForPreview(key string) string
	DeleteFile(key string) error
	DeleteFiles(keys []string) error
	GetFileByteArray(key string) ([]byte, error)
	GetFileContent(key string) (string, error)
	GetFullFileUrl(fileUrl string) string
	GetKeyFromFullUrl(fullUrl string) string
	UploadContent(path, content string) error
	UploadMedia(key string, buffer []byte) (string, error)
	ReadFileFromAbsolutePath(absolutePath string) ([]byte, error)
}

// DesktopLocalFileImpl desktop local file service implementation
type DesktopLocalFileImpl struct {
	BaseDir string // Base file directory
}

// NewDesktopLocalFileImpl creates new local file service implementation
func NewDesktopLocalFileImpl(baseDir string) *DesktopLocalFileImpl {
	return &DesktopLocalFileImpl{BaseDir: baseDir}
}

// CreatePreSignedUrl creates pre-signed upload URL (local version directly returns file path)
func (d *DesktopLocalFileImpl) CreatePreSignedUrl(key string) string {
	// In local file implementation, no pre-signed URL is needed
	// Directly return file path
	return filepath.Join(d.BaseDir, key)
}

// CreatePreSignedUrlForPreview creates pre-signed preview URL
func (d *DesktopLocalFileImpl) CreatePreSignedUrlForPreview(key string) string {
	// For local files, directly return local file path
	return filepath.Join(d.BaseDir, key)
}

// DeleteFile deletes file
func (d *DesktopLocalFileImpl) DeleteFile(key string) error {
	return d.DeleteFiles([]string{key})
}

// DeleteFiles batch deletes files
func (d *DesktopLocalFileImpl) DeleteFiles(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	var errorsEncountered []error
	for _, key := range keys {
		if !d.isValidKey(key) {
			errorsEncountered = append(errorsEncountered, fmt.Errorf("invalid local file path: %s", key))
			continue
		}

		localPath := filepath.Join(d.BaseDir, key)
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
func (d *DesktopLocalFileImpl) isValidKey(key string) bool {
	// Check if it is a valid local path
	cleanKey := filepath.Clean(key)
	if strings.Contains(cleanKey, "..") {
		return false
	}
	return filepath.IsLocal(key)
}

// GetFileByteArray gets file byte array
func (d *DesktopLocalFileImpl) GetFileByteArray(key string) ([]byte, error) {
	localPath := d.getLocalPath(key)
	return ioutil.ReadFile(localPath)
}

// GetFileContent gets file content
func (d *DesktopLocalFileImpl) GetFileContent(key string) (string, error) {
	localPath := d.getLocalPath(key)
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// getLocalPath gets local path
func (d *DesktopLocalFileImpl) getLocalPath(key string) string {
	return filepath.Join(d.BaseDir, key)
}

// GetFullFileUrl gets complete file URL
func (d *DesktopLocalFileImpl) GetFullFileUrl(fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	// For local files, directly return local path
	return d.GetKeyFromFullUrl(fileUrl)
}

// GetKeyFromFullUrl extracts key from complete URL
func (d *DesktopLocalFileImpl) GetKeyFromFullUrl(fullUrl string) string {
	// If it's already a file path, return directly
	if d.isLocalPath(fullUrl) {
		return fullUrl
	}

	// Parse URL
	parsedURL, err := url.Parse(fullUrl)
	if err != nil {
		log.Printf("Failed to parse URL %s: %v", fullUrl, err)
		return ""
	}

	// Extract path part
	filePath := parsedURL.Path

	// Remove beginning slash
	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	return filePath
}

// isLocalPath determines if it's a local path
func (d *DesktopLocalFileImpl) isLocalPath(path string) bool {
	return !strings.Contains(path, "://")
}

// UploadContent upload content (not implemented)
func (d *DesktopLocalFileImpl) UploadContent(filePath, content string) error {
	// This needs to be implemented according to specific requirements
	log.Printf("UploadContent not implemented for path: %s", filePath)
	return errors.New("uploadContent not implemented")
}

// ReadFileFromAbsolutePath reads file from absolute path (for drag & drop)
func (d *DesktopLocalFileImpl) ReadFileFromAbsolutePath(absolutePath string) ([]byte, error) {
	// Security check: ensure path is absolute and exists
	if !filepath.IsAbs(absolutePath) {
		return nil, fmt.Errorf("path must be absolute: %s", absolutePath)
	}

	// Check if file exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", absolutePath)
	}

	// Read file
	return ioutil.ReadFile(absolutePath)
}

// UploadMedia upload media file
func (d *DesktopLocalFileImpl) UploadMedia(key string, buffer []byte) (string, error) {
	// Extract filename from key
	filename := filepath.Base(key)

	// Calculate file's SHA256 hash
	hash := sha256.Sum256(buffer)
	hashString := fmt.Sprintf("%x", hash)

	// Construct local file path
	localPath := d.getLocalPath(key)

	// Ensure directory exists
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(localPath, buffer, 0644); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", localPath, err)
	}

	log.Printf("File uploaded successfully: %s (hash: %s)", filename, hashString)
	return key, nil
}

// extractKeyFromUrlOrReturnOriginal extracts key from URL or returns original string
// Handles legacy data where full URLs were stored instead of keys
func extractKeyFromUrlOrReturnOriginal(url string, getKeyFromFullUrl func(string) string) string {
	// Only process URLs that start with http:// or https://
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// Extract key from full URL for legacy data compatibility
		return getKeyFromFullUrl(url)
	}
	// Return original input if it's already a key
	return url
}

// S3StaticFileImpl S3 file service implementation
type S3StaticFileImpl struct {
	Bucket          string
	PublicDomain    string
	EnablePathStyle bool
	SetACL          bool
}

// CreateS3StaticFileImpl creates new S3 file service implementation
func NewS3StaticFileImpl(bucket, publicDomain string, enablePathStyle, setACL bool) *S3StaticFileImpl {
	return &S3StaticFileImpl{
		Bucket:          bucket,
		PublicDomain:    publicDomain,
		EnablePathStyle: enablePathStyle,
		SetACL:          setACL,
	}
}

// CreatePreSignedUrl creates pre-signed upload URL for S3
func (s *S3StaticFileImpl) CreatePreSignedUrl(key string) string {
	// Placeholder - in real implementation would generate S3 pre-signed URL
	log.Printf("CreatePreSignedUrl not implemented for key: %s", key)
	return ""
}

// CreatePreSignedUrlForPreview creates pre-signed preview URL for S3
func (s *S3StaticFileImpl) CreatePreSignedUrlForPreview(key string) string {
	// Placeholder - in real implementation would generate S3 pre-signed URL for preview
	log.Printf("CreatePreSignedUrlForPreview not implemented for key: %s", key)
	return ""
}

// DeleteFile deletes S3 file
func (s *S3StaticFileImpl) DeleteFile(key string) error {
	// Placeholder - in real implementation would delete from S3
	log.Printf("DeleteFile not implemented for key: %s", key)
	return errors.New("deleteFile not implemented for S3")
}

// DeleteFiles batch deletes S3 files
func (s *S3StaticFileImpl) DeleteFiles(keys []string) error {
	// Placeholder - in real implementation would batch delete from S3
	log.Printf("DeleteFiles not implemented for keys: %v", keys)
	return errors.New("deleteFiles not implemented for S3")
}

// GetFileByteArray gets S3 file byte array
func (s *S3StaticFileImpl) GetFileByteArray(key string) ([]byte, error) {
	// Placeholder - in real implementation would download from S3
	log.Printf("GetFileByteArray not implemented for key: %s", key)
	return nil, errors.New("getFileByteArray not implemented for S3")
}

// GetFileContent gets S3 file content
func (s *S3StaticFileImpl) GetFileContent(key string) (string, error) {
	// Placeholder - in real implementation would download from S3
	log.Printf("GetFileContent not implemented for key: %s", key)
	return "", errors.New("getFileContent not implemented for S3")
}

// GetFullFileUrl gets complete S3 file URL
func (s *S3StaticFileImpl) GetFullFileUrl(fileUrl string) string {
	if fileUrl == "" {
		return ""
	}

	// Handle legacy data compatibility
	key := extractKeyFromUrlOrReturnOriginal(fileUrl, s.GetKeyFromFullUrl)

	// If bucket is not set public read, preview address would need to be regenerated
	// Placeholder - in real implementation would generate pre-signed URL if needed
	if !s.SetACL {
		log.Printf("Would need to generate pre-signed URL for private bucket")
		// In real implementation: return createPreSignedUrlForPreview(key)
	}

	// Construct public URL
	if s.EnablePathStyle && s.Bucket != "" {
		// Path style: https://domain.com/bucket/key
		return fmt.Sprintf("%s/%s/%s", s.PublicDomain, s.Bucket, key)
	}
	// Virtual-hosted style: https://bucket.domain.com/key or https://domain.com/key
	return fmt.Sprintf("%s/%s", s.PublicDomain, key)
}

// GetKeyFromFullUrl extracts key from full S3 URL
func (s *S3StaticFileImpl) GetKeyFromFullUrl(fullUrl string) string {
	// Parse URL
	parsedURL, err := url.Parse(fullUrl)
	if err != nil {
		log.Printf("Failed to parse URL %s: %v", fullUrl, err)
		return fullUrl
	}

	pathname := parsedURL.Path

	if s.EnablePathStyle && s.Bucket != "" {
		// Path style: /bucket/key -> key
		bucketPrefix := "/" + s.Bucket + "/"
		if strings.HasPrefix(pathname, bucketPrefix) {
			return strings.TrimPrefix(pathname, bucketPrefix)
		}
		// Fallback
		return strings.TrimPrefix(pathname, "/")
	}

	// Virtual-hosted style: /key -> key
	return strings.TrimPrefix(pathname, "/")
}

// UploadContent uploads content to S3
func (s *S3StaticFileImpl) UploadContent(path, content string) error {
	// Placeholder - in real implementation would upload to S3
	log.Printf("UploadContent not implemented for path: %s", path)
	return errors.New("uploadContent not implemented for S3")
}

// UploadMedia uploads media file to S3
func (s *S3StaticFileImpl) UploadMedia(key string, buffer []byte) (string, error) {
	// Placeholder - in real implementation would upload buffer to S3
	log.Printf("UploadMedia not implemented for key: %s", key)
	return key, nil
}

// ReadFileFromAbsolutePath not supported for S3 (only for local files)
func (s *S3StaticFileImpl) ReadFileFromAbsolutePath(absolutePath string) ([]byte, error) {
	return nil, errors.New("ReadFileFromAbsolutePath not supported for S3 implementation")
}

// CreateFileServiceModule creates file service module (factory function)
func CreateFileServiceModule(baseDir, s3Bucket, s3PublicDomain string, s3EnablePathStyle, s3SetACL, isDesktop bool) FileServiceImpl {
	if isDesktop {
		// Use local file implementation for desktop
		return NewDesktopLocalFileImpl(baseDir)
	}
	// Use S3 implementation for server
	return NewS3StaticFileImpl(s3Bucket, s3PublicDomain, s3EnablePathStyle, s3SetACL)
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

// FileService wraps FileServiceImpl to provide high-level file operations
type FileService struct {
	UserId    string
	FileModel FileModelInterface
	Impl      FileServiceImpl
}

// NewFileService creates new file service with database and implementation
func NewFileService(userId string, impl FileServiceImpl) *FileService {
	return &FileService{
		UserId: userId,
		// FileModel would need proper database initialization - placeholder
		Impl: impl,
	}
}

// DeleteFile deletes file via implementation
func (fs *FileService) DeleteFile(key string) error {
	return fs.Impl.DeleteFile(key)
}

// DeleteFiles batch deletes files via implementation
func (fs *FileService) DeleteFiles(keys []string) error {
	return fs.Impl.DeleteFiles(keys)
}

// GetFileContent gets file content via implementation
func (fs *FileService) GetFileContent(key string) (string, error) {
	return fs.Impl.GetFileContent(key)
}

// GetFileByteArray gets file byte array via implementation
func (fs *FileService) GetFileByteArray(key string) ([]byte, error) {
	return fs.Impl.GetFileByteArray(key)
}

// CreatePreSignedUrl creates pre-signed upload URL via implementation
func (fs *FileService) CreatePreSignedUrl(key string) string {
	return fs.Impl.CreatePreSignedUrl(key)
}

// CreatePreSignedUrlForPreview creates pre-signed preview URL via implementation
func (fs *FileService) CreatePreSignedUrlForPreview(key string) string {
	return fs.Impl.CreatePreSignedUrlForPreview(key)
}

// UploadContent uploads content via implementation
func (fs *FileService) UploadContent(path, content string) error {
	return fs.Impl.UploadContent(path, content)
}

// GetFullFileUrl gets complete file URL via implementation
func (fs *FileService) GetFullFileUrl(url string) string {
	return fs.Impl.GetFullFileUrl(url)
}

// GetKeyFromFullUrl extracts key from full URL via implementation
func (fs *FileService) GetKeyFromFullUrl(url string) string {
	return fs.Impl.GetKeyFromFullUrl(url)
}

// UploadMedia uploads media file via implementation
func (fs *FileService) UploadMedia(key string, buffer []byte) (string, error) {
	return fs.Impl.UploadMedia(key, buffer)
}

// ReadFileFromAbsolutePath reads file from absolute path (for drag & drop)
func (fs *FileService) ReadFileFromAbsolutePath(absolutePath string) ([]byte, error) {
	return fs.Impl.ReadFileFromAbsolutePath(absolutePath)
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
		// Handle S3 not found case
		if strings.Contains(err.Error(), "NoSuchKey") {
			// Remove from database if it's gone from storage
			if delErr := fs.FileModel.Delete(fileId, false); delErr != nil {
				log.Printf("Failed to delete missing file from DB: %v", delErr)
			}
			return nil, nil, "", errors.New("file content is empty")
		}
		return nil, nil, "", fmt.Errorf("failed to get file content: %w", err)
	}

	if len(content) == 0 {
		return nil, nil, "", errors.New("file content is empty")
	}

	// Convert Uint8Array to base64 string for Wails temp file binding
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

// WriteTempFile and Cleanup functions - placeholder for external binding
// These would be imported from the Wails binding in actual implementation

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
