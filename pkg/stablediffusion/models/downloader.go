package models

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Downloader handles downloading models from various sources
type Downloader struct {
	client      *http.Client
	progressCb  func(downloaded, total int64)
	destination string
}

// DownloaderOption configures the downloader
type DownloaderOption func(*Downloader)

// WithProgressCallback sets a progress callback
func WithProgressCallback(cb func(downloaded, total int64)) DownloaderOption {
	return func(d *Downloader) {
		d.progressCb = cb
	}
}

// WithTimeout sets the HTTP timeout
func WithTimeout(timeout time.Duration) DownloaderOption {
	return func(d *Downloader) {
		d.client.Timeout = timeout
	}
}

// NewDownloader creates a new model downloader
func NewDownloader(destination string, opts ...DownloaderOption) *Downloader {
	d := &Downloader{
		client: &http.Client{
			Timeout: 30 * time.Minute, // Long timeout for large files
		},
		destination: destination,
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

// Download downloads a model from a URL
func (d *Downloader) Download(sourceURL, filename string) (*ModelInfo, error) {
	// Parse URL
	u, err := url.Parse(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(d.destination, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Determine filename
	if filename == "" {
		filename = filepath.Base(u.Path)
	}

	destinationPath := filepath.Join(d.destination, filename)

	// Check if file already exists
	if _, err := os.Stat(destinationPath); err == nil {
		return nil, fmt.Errorf("file already exists: %s", destinationPath)
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "stable-diffusion-go/1.0")

	// Perform download
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Get content length
	contentLength := resp.ContentLength

	// Create destination file
	out, err := os.Create(destinationPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Printf("failed to close output file: %v", err)
		}
	}()

	// Copy with progress
	var downloaded int64
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return nil, fmt.Errorf("failed to write file: %w", writeErr)
			}
			downloaded += int64(n)

			// Report progress
			if d.progressCb != nil {
				d.progressCb(downloaded, contentLength)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("download error: %w", err)
		}
	}

	// Get file info
	fileInfo, err := out.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Detect model info
	detector := NewDetector()
	model, err := detector.Detect(destinationPath)
	if err != nil {
		// Return basic info if detection fails
		model = &ModelInfo{
			ID:       generateModelIDFromPath(destinationPath),
			Name:     strings.TrimSuffix(filename, filepath.Ext(filename)),
			Type:     ModelTypeUnknown,
			Path:     destinationPath,
			Format:   detectFormat(filepath.Ext(filename)),
			Size:     fileInfo.Size(),
			Tags:     []string{},
			Metadata: map[string]string{"source": sourceURL},
			AddedAt:  time.Now(),
		}
	}

	model.Source = sourceURL

	return model, nil
}

// HuggingFaceDownloader downloads models from HuggingFace
type HuggingFaceDownloader struct {
	*Downloader
	repo string
}

// NewHuggingFaceDownloader creates a new HuggingFace downloader
func NewHuggingFaceDownloader(repo, destination string, opts ...DownloaderOption) *HuggingFaceDownloader {
	return &HuggingFaceDownloader{
		Downloader: NewDownloader(destination, opts...),
		repo:       repo,
	}
}

// Download downloads a file from HuggingFace
func (h *HuggingFaceDownloader) Download(filename string) (*ModelInfo, error) {
	url := fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", h.repo, filename)
	return h.Downloader.Download(url, filename)
}

// DownloadWithRevision downloads a specific revision
func (h *HuggingFaceDownloader) DownloadWithRevision(filename, revision string) (*ModelInfo, error) {
	url := fmt.Sprintf("https://huggingface.co/%s/resolve/%s/%s", h.repo, revision, filename)
	return h.Downloader.Download(url, filename)
}

// CivitaiDownloader downloads models from Civitai
type CivitaiDownloader struct {
	*Downloader
}

// NewCivitaiDownloader creates a new Civitai downloader
func NewCivitaiDownloader(destination string, opts ...DownloaderOption) *CivitaiDownloader {
	return &CivitaiDownloader{
		Downloader: NewDownloader(destination, opts...),
	}
}

// DownloadByID downloads a model by Civitai model ID
func (c *CivitaiDownloader) DownloadByID(modelID int, filename string) (*ModelInfo, error) {
	url := fmt.Sprintf("https://civitai.com/api/download/models/%d", modelID)
	return c.Download(url, filename)
}

// DownloadByVersionID downloads a specific model version
func (c *CivitaiDownloader) DownloadByVersionID(versionID int, filename string) (*ModelInfo, error) {
	url := fmt.Sprintf("https://civitai.com/api/download/models/%d", versionID)
	return c.Download(url, filename)
}

// DownloadSource represents a download source configuration
type DownloadSource struct {
	Type        string            `json:"type" yaml:"type"` // huggingface, civitai, direct
	URL         string            `json:"url" yaml:"url"`
	Repo        string            `json:"repo" yaml:"repo"`         // For HuggingFace
	ModelID     int               `json:"model_id" yaml:"model_id"` // For Civitai
	Filename    string            `json:"filename" yaml:"filename"`
	Headers     map[string]string `json:"headers" yaml:"headers"`
	Description string            `json:"description" yaml:"description"`
}

// DownloadManager manages multiple downloads
type DownloadManager struct {
	downloader *Downloader
	queue      []DownloadTask
	results    []DownloadResult
}

// DownloadTask represents a download task
type DownloadTask struct {
	Source   DownloadSource
	Priority int
	Status   DownloadStatus
	Error    error
	Model    *ModelInfo
}

// DownloadStatus represents the status of a download
type DownloadStatus int

const (
	StatusPending DownloadStatus = iota
	StatusDownloading
	StatusCompleted
	StatusFailed
	StatusCancelled
)

// DownloadResult represents the result of a download
type DownloadResult struct {
	Task  DownloadTask
	Model *ModelInfo
	Error error
	Time  time.Duration
}

// NewDownloadManager creates a new download manager
func NewDownloadManager(destination string) *DownloadManager {
	return &DownloadManager{
		downloader: NewDownloader(destination),
		queue:      make([]DownloadTask, 0),
		results:    make([]DownloadResult, 0),
	}
}

// Add adds a download task to the queue
func (dm *DownloadManager) Add(source DownloadSource, priority int) {
	dm.queue = append(dm.queue, DownloadTask{
		Source:   source,
		Priority: priority,
		Status:   StatusPending,
	})
}

// Execute executes all pending downloads
func (dm *DownloadManager) Execute() []DownloadResult {
	results := make([]DownloadResult, 0, len(dm.queue))

	for i := range dm.queue {
		task := &dm.queue[i]
		task.Status = StatusDownloading

		start := time.Now()
		var model *ModelInfo
		var err error

		switch task.Source.Type {
		case "direct":
			model, err = dm.downloader.Download(task.Source.URL, task.Source.Filename)
		case "huggingface":
			hf := NewHuggingFaceDownloader(task.Source.Repo, dm.downloader.destination)
			model, err = hf.Download(task.Source.Filename)
		case "civitai":
			civ := NewCivitaiDownloader(dm.downloader.destination)
			model, err = civ.DownloadByID(task.Source.ModelID, task.Source.Filename)
		default:
			err = fmt.Errorf("unknown source type: %s", task.Source.Type)
		}

		duration := time.Since(start)

		if err != nil {
			task.Status = StatusFailed
			task.Error = err
		} else {
			task.Status = StatusCompleted
			task.Model = model
		}

		results = append(results, DownloadResult{
			Task:  *task,
			Model: model,
			Error: err,
			Time:  duration,
		})
	}

	dm.results = append(dm.results, results...)
	return results
}

// GetResults returns all download results
func (dm *DownloadManager) GetResults() []DownloadResult {
	return dm.results
}

// SimpleProgressReporter returns a simple progress callback
func SimpleProgressReporter() func(downloaded, total int64) {
	return func(downloaded, total int64) {
		if total > 0 {
			percent := float64(downloaded) / float64(total) * 100
			fmt.Printf("\rDownloading: %.1f%% (%s / %s)",
				percent,
				formatBytes(downloaded),
				formatBytes(total),
			)
		} else {
			fmt.Printf("\rDownloading: %s", formatBytes(downloaded))
		}
	}
}

// formatBytes formats bytes to human-readable string
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ParseHuggingFaceURL parses a HuggingFace URL to extract repo and filename
func ParseHuggingFaceURL(hfURL string) (repo, filename string, err error) {
	u, err := url.Parse(hfURL)
	if err != nil {
		return "", "", err
	}

	// Expected format: https://huggingface.co/{repo}/resolve/main/{filename}
	parts := strings.Split(u.Path, "/")
	if len(parts) < 5 {
		return "", "", fmt.Errorf("invalid HuggingFace URL format")
	}

	// Find "resolve" in path
	for i, part := range parts {
		if part == "resolve" && i+2 < len(parts) {
			repo = strings.Join(parts[1:i], "/")
			filename = strings.Join(parts[i+2:], "/")
			return repo, filename, nil
		}
	}

	return "", "", fmt.Errorf("could not parse HuggingFace URL")
}

// ParseCivitaiURL parses a Civitai URL to extract model ID
func ParseCivitaiURL(civitaiURL string) (modelID int, err error) {
	u, err := url.Parse(civitaiURL)
	if err != nil {
		return 0, err
	}

	// Check for models/{id} pattern
	if strings.Contains(u.Path, "/models/") {
		parts := strings.Split(u.Path, "/")
		for i, part := range parts {
			if part == "models" && i+1 < len(parts) {
				id, err := strconv.Atoi(parts[i+1])
				if err == nil {
					return id, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("could not parse Civitai URL")
}
