package download

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/grab"
	"github.com/kawai-network/veridium/pkg/xlog"
	"golang.org/x/time/rate"
)

var (
	// ErrUnknownOS is returned when an unknown operating system is specified
	ErrUnknownOS = errors.New("unknown OS")
	// ErrUnknownProcessor is returned when an unknown processor is specified
	ErrUnknownProcessor = errors.New("unknown processor")
	// ErrInvalidVersion is returned when an invalid version string is provided
	ErrInvalidVersion = errors.New("invalid version")
)

// RetryCount is how many times the package will retry to obtain the latest llama.cpp version.
var RetryCount = 3

// FallbackVersion is used when GitHub API is unavailable (rate limit, network issues, etc.)
// This should be updated periodically to a known stable version
const FallbackVersion = "b7248"

// LlamaLatestVersion fetches the latest release tag of llama.cpp from the GitHub API.
// Falls back to a hardcoded version if GitHub API is unavailable.
func LlamaLatestVersion() (string, error) {
	var version string
	var err error
	for range RetryCount {
		version, err = getLatestVersion()
		if err == nil {
			return version, nil
		}
		time.Sleep(3 * time.Second)
	}

	// If all retries failed, use fallback version
	xlog.Warn("⚠️  Failed to fetch version from GitHub API", "error", err)
	xlog.Warn("📦 Using fallback version", "version", FallbackVersion)
	return FallbackVersion, nil
}

func getLatestVersion() (string, error) {
	const apiURL = "https://api.github.com/repos/ggml-org/llama.cpp/releases/latest"

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code %d from GitHub API", resp.StatusCode)
	}

	var result struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode GitHub API response: %w", err)
	}

	if result.TagName == "" {
		return "", fmt.Errorf("empty version tag from GitHub API")
	}

	xlog.Info("📦 Latest llama.cpp version from GitHub", "version", result.TagName)
	return result.TagName, nil
}

// LlamaAvailableVersions returns a list of available llama.cpp versions from GitHub releases.
// Returns versions in descending order (newest first).
// limit specifies how many releases to fetch (default: 10, max: 100).
func LlamaAvailableVersions(limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10 // Default to 10 releases
	}
	if limit > 100 {
		limit = 100 // GitHub API max per page
	}

	url := fmt.Sprintf("https://api.github.com/repos/ggml-org/llama.cpp/releases?per_page=%d", limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var releases []struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode releases: %w", err)
	}

	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		if release.TagName != "" {
			versions = append(versions, release.TagName)
		}
	}

	return versions, nil
}

// Get downloads the llama.cpp precompiled binaries for the desired OS/processor.
// os can be one of the following values: "linux", "darwin", "windows".
// processor can be one of the following values: "cpu", "cuda", "vulkan", "metal".
// version should be the desired `b1234` formatted llama.cpp version. You can use the
// [LlamaLatestVersion] function to obtain the latest release.
// dest in the destination directory for the downloaded binaries.
func Get(os string, processor string, version string, dest string) error {
	if err := VersionIsValid(version); err != nil {
		return err
	}

	var location, filename string
	location = fmt.Sprintf("https://github.com/ggml-org/llama.cpp/releases/download/%s", version)

	switch os {
	case "linux":
		switch processor {
		case "cpu":
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.zip//build/bin", version)
		case "cuda":
			location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-x64.zip", version)
		case "vulkan":
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.zip//build/bin", version)
		default:
			return ErrUnknownProcessor
		}
	case "darwin":
		switch processor {
		case "cpu", "metal":
			filename = fmt.Sprintf("llama-%s-bin-macos-arm64.zip//build/bin", version)
		default:
			return ErrUnknownProcessor
		}

	case "windows":
		switch processor {
		case "cpu":
			filename = fmt.Sprintf("llama-%s-bin-win-cpu-x64.zip//build/bin", version)
		case "cuda":
			filename = fmt.Sprintf("llama-%s-bin-win-cuda-12.4-x64.zip//build/bin", version)
		case "vulkan":
			filename = fmt.Sprintf("llama-%s-bin-win-vulkan-x64.zip//build/bin", version)
		default:
			return ErrUnknownProcessor
		}

	default:
		return ErrUnknownOS
	}

	// Extract the actual filename (before //) for URL construction
	actualFilename := filename
	if strings.Contains(filename, "//") {
		actualFilename = strings.SplitN(filename, "//", 2)[0]
	}

	url := fmt.Sprintf("%s/%s", location, actualFilename)
	return get(url, filename, dest)
}

// get downloads a file using grab and optionally extracts it if it's a ZIP
// Implements retry logic with exponential backoff for GitHub CDN propagation delays
func get(url, filename, dest string) error {
	const maxRetries = 3
	const initialBackoff = 2 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := initialBackoff * time.Duration(1<<uint(attempt-1)) // 2s, 4s, 8s
			xlog.Info("⏳ Retry download (GitHub CDN may be propagating)...", "attempt", attempt+1, "max_retries", maxRetries, "backoff", backoff)
			time.Sleep(backoff)
		}

		err := downloadAndExtract(url, filename, dest)
		if err == nil {
			return nil // Success!
		}

		lastErr = err

		// Check if it's a 404 error (release assets not ready yet)
		if strings.Contains(err.Error(), "404") {
			// For 404, retry might help if it's CDN propagation delay
			if attempt < maxRetries-1 {
				xlog.Warn("⚠️  404 error, retrying (may be CDN delay)...")
				continue
			}
			// After all retries, return specific 404 message
			return fmt.Errorf("download failed: release assets not available yet (404). The release tag exists but binaries are still being built. Please try again in a few minutes or use a previous version")
		}

		// For other errors, don't retry
		return err
	}

	return fmt.Errorf("download failed after %d attempts: %w", maxRetries, lastErr)
}

// downloadAndExtract performs a single download attempt with optional ZIP extraction
func downloadAndExtract(url, filename, dest string) error {
	// Create temp directory for download
	tempDir, err := os.MkdirTemp("", "llama-download-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Determine if we need to extract from ZIP
	// go-getter syntax: "file.zip//path/inside" means extract path/inside from file.zip
	var zipFile, extractPath string
	if strings.Contains(filename, "//") {
		parts := strings.SplitN(filename, "//", 2)
		zipFile = parts[0]
		extractPath = parts[1]
	} else {
		zipFile = filename
	}

	// Download ZIP file
	zipPath := filepath.Join(tempDir, filepath.Base(zipFile))
	req, err := grab.NewRequest(zipPath, url)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	client := grab.NewClient()
	resp := client.Do(req)

	// Wait for download to complete
	<-resp.Done
	if err := resp.Err(); err != nil {
		// Check if it's a 404 error
		if resp.HTTPResponse != nil && resp.HTTPResponse.StatusCode == 404 {
			return fmt.Errorf("404 error")
		}
		return fmt.Errorf("download failed: %w", err)
	}

	// If no extraction needed, just move the file
	if extractPath == "" {
		return os.Rename(zipPath, filepath.Join(dest, filepath.Base(zipFile)))
	}

	// Extract specific path from ZIP
	return extractFromZip(zipPath, extractPath, dest)
}

// extractFromZip extracts a specific path from a ZIP file to destination
func extractFromZip(zipPath, extractPath, dest string) error {
	// Open ZIP file
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}

	// Extract files that match the extractPath
	extractedCount := 0
	for _, f := range r.File {
		// Check if file is in the extract path
		if !strings.HasPrefix(f.Name, extractPath) {
			continue
		}

		// Calculate relative path
		relPath := strings.TrimPrefix(f.Name, extractPath)
		relPath = strings.TrimPrefix(relPath, "/")
		if relPath == "" {
			continue // Skip the directory itself
		}

		targetPath := filepath.Join(dest, relPath)

		if f.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Extract file
		if err := extractFile(f, targetPath); err != nil {
			return fmt.Errorf("failed to extract %s: %w", f.Name, err)
		}
		extractedCount++
	}

	if extractedCount == 0 {
		return fmt.Errorf("no files found in path %s", extractPath)
	}

	return nil
}

// extractFile extracts a single file from ZIP
func extractFile(f *zip.File, targetPath string) error {
	// Open file in ZIP
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create target file
	outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy content
	if _, err := io.Copy(outFile, rc); err != nil {
		return err
	}

	return nil
}

// VersionIsValid checks if the provided version string is valid.
func VersionIsValid(version string) error {
	if version == "" {
		return fmt.Errorf("%w: version string is empty", ErrInvalidVersion)
	}
	if !strings.HasPrefix(version, "b") {
		return fmt.Errorf("%w: version must start with 'b', got: %s", ErrInvalidVersion, version)
	}

	return nil
}

// LibraryName returns the name for the llama.cpp library for any given OS.
func LibraryName(os string) string {
	switch os {
	case "linux", "freebsd":
		return "libllama.so"
	case "windows":
		return "llama.dll"
	case "darwin":
		return "libllama.dylib"
	default:
		return "unknown"
	}
}

// RequiredLibraries returns all required library files for llama.cpp on a given OS
// llama.cpp requires multiple libraries: libggml, libggml-base, and libllama
// This is critical for library-based usage (via yzma)
func RequiredLibraries(os string) []string {
	switch os {
	case "linux", "freebsd":
		return []string{
			"libggml.so",
			"libggml-base.so",
			"libllama.so",
		}
	case "windows":
		return []string{
			"ggml.dll",
			"ggml-base.dll",
			"llama.dll",
		}
	case "darwin":
		return []string{
			"libggml.dylib",
			"libggml-base.dylib",
			"libllama.dylib",
		}
	default:
		return []string{}
	}
}

// GetLibraryExtension returns the library file extension for a given OS
func GetLibraryExtension(os string) string {
	switch os {
	case "linux", "freebsd":
		return ".so"
	case "windows":
		return ".dll"
	case "darwin":
		return ".dylib"
	default:
		return ""
	}
}

// ============================================================================
// Advanced Download Functions with Progress Tracking and Rate Limiting
// ============================================================================

// DownloadOptions configures download behavior
type DownloadOptions struct {
	// MaxRetries specifies how many times to retry failed downloads
	MaxRetries int

	// ShowProgress enables real-time progress logging
	ShowProgress bool

	// ResumeIfPossible enables automatic resume of partial downloads
	ResumeIfPossible bool

	// ProgressInterval specifies how often to log progress updates
	ProgressInterval time.Duration

	// RateLimitMBps limits download speed in MB/s (0 = unlimited)
	RateLimitMBps int
}

// DefaultDownloadOptions returns sensible default options
func DefaultDownloadOptions() DownloadOptions {
	return DownloadOptions{
		MaxRetries:       3,
		ShowProgress:     true,
		ResumeIfPossible: true,
		ProgressInterval: 2 * time.Second,
		RateLimitMBps:    0, // Unlimited by default
	}
}

// WithRateLimit returns download options with specified rate limit in MB/s
func WithRateLimit(mbps int) DownloadOptions {
	opts := DefaultDownloadOptions()
	opts.RateLimitMBps = mbps
	return opts
}

// GetWithProgress downloads a file using grab with retry logic, progress tracking,
// and automatic resume support. This is the recommended way to download large files.
//
// Features:
// - Automatic retry with exponential backoff (2s, 4s, 6s)
// - Resume support for interrupted downloads
// - Real-time progress tracking with speed and ETA
// - Optional rate limiting
// - Validates HTTP status codes
// - Thread-safe and context-aware
//
// The file is downloaded to destPath. If the download is interrupted,
// it can be resumed on the next attempt automatically.
func GetWithProgress(url, destPath string, opts DownloadOptions) error {
	client := grab.NewClient()

	var lastErr error
	for attempt := 1; attempt <= opts.MaxRetries; attempt++ {
		if attempt > 1 {
			backoff := time.Duration(attempt) * 2 * time.Second
			xlog.Info("🔄 Retry attempt (download)", "attempt", attempt, "max", opts.MaxRetries, "backoff", backoff)
			time.Sleep(backoff)
		}

		// Create request
		req, err := grab.NewRequest(destPath, url)
		if err != nil {
			lastErr = fmt.Errorf("failed to create download request: %w", err)
			xlog.Warn("⚠️  Download attempt failed", "attempt", attempt, "error", lastErr)
			continue
		}

		// Configure resume behavior
		req.NoResume = !opts.ResumeIfPossible

		// Configure rate limiting if specified
		if opts.RateLimitMBps > 0 {
			bytesPerSecond := opts.RateLimitMBps * 1024 * 1024
			req.RateLimiter = rate.NewLimiter(rate.Limit(bytesPerSecond), bytesPerSecond)
			xlog.Debug("🔧 Download rate limit set", "limit_mbps", opts.RateLimitMBps)
		}

		// Start download
		resp := client.Do(req)

		// Track progress
		if opts.ShowProgress {
			if err := trackProgress(resp, opts.ProgressInterval); err != nil {
				lastErr = err
				xlog.Warn("⚠️  Download attempt failed", "attempt", attempt, "error", lastErr)
				continue
			}
		} else {
			// Wait for completion without progress tracking
			<-resp.Done
			if err := resp.Err(); err != nil {
				lastErr = err
				xlog.Warn("⚠️  Download attempt failed", "attempt", attempt, "error", lastErr)
				continue
			}
		}

		// Success!
		sizeMB := float64(resp.Size()) / (1024 * 1024)
		xlog.Info("✅ Download complete", "filename", resp.Filename, "size_mb", sizeMB)

		if resp.DidResume {
			xlog.Info("📦 Download resumed successfully")
		}

		return nil
	}

	return fmt.Errorf("download failed after %d attempts: %w", opts.MaxRetries, lastErr)
}

// trackProgress monitors download progress and logs updates at regular intervals
func trackProgress(resp *grab.Response, interval time.Duration) error {
	t := time.NewTicker(interval)
	defer t.Stop()

	lastProgress := float64(0)
	stuckCount := 0
	maxStuckCount := 5 // Consider download stuck after 5 intervals with no progress

	for {
		select {
		case <-t.C:
			progress := resp.Progress() * 100
			speed := resp.BytesPerSecond() / (1024 * 1024) // MB/s
			eta := resp.ETA().Round(time.Second)

			// Check if download is stuck
			if progress == lastProgress && progress < 100 {
				stuckCount++
				if stuckCount >= maxStuckCount {
					return fmt.Errorf("download appears stuck at %.1f%% for %v", progress, interval*time.Duration(maxStuckCount))
				}
			} else {
				stuckCount = 0
			}
			lastProgress = progress

			if speed > 0 {
				xlog.Info("📥 Download progress", "file", filepath.Base(resp.Filename), "progress_percent", fmt.Sprintf("%.1f", progress), "speed_mbps", fmt.Sprintf("%.2f", speed), "eta", eta)
			} else {
				xlog.Info("📥 Download progress (starting...)", "file", filepath.Base(resp.Filename), "progress_percent", fmt.Sprintf("%.1f", progress))
			}

		case <-resp.Done:
			// Download completed or failed
			if err := resp.Err(); err != nil {
				return fmt.Errorf("download error: %w", err)
			}
			return nil
		}
	}
}

// GetBatch downloads multiple files concurrently using grab's batch feature.
// This is useful for downloading multiple files at once.
//
// Returns a channel that receives responses for each download as they complete.
// The channel is closed when all downloads are complete.
func GetBatch(workers int, destDir string, urls ...string) (<-chan *grab.Response, error) {
	if workers <= 0 {
		workers = 3 // Default to 3 concurrent downloads
	}

	return grab.GetBatch(workers, destDir, urls...)
}
