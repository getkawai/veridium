package download

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/grab"
)

var (
	ErrUnknownArch    = errors.New("unknown architecture")
	ErrUnknownOS      = errors.New("unknown OS")
	ErrInvalidVersion = errors.New("invalid version")
)

var (
	// DefaultVersion is the default stable-diffusion.cpp version to use
	DefaultVersion = "master-487-43e829f"
	// RetryCount is how many times the package will retry to obtain the latest stable-diffusion.cpp version.
	RetryCount = 3
	// RetryDelay is the delay between retries when obtaining the latest stable-diffusion.cpp version.
	RetryDelay = 3 * time.Second
	// apiURL is the GitHub API URL for fetching the latest stable-diffusion.cpp version.
	apiURL = "https://api.github.com/repos/leejet/stable-diffusion.cpp/releases/latest"
)

// Arch represents the CPU architecture
type Arch int

const (
	AMD64 Arch = iota
	ARM64
)

// OS represents the operating system
type OS int

const (
	Linux OS = iota
	Darwin
	Windows
)

// ParseArch parses a string into an Arch type
func ParseArch(arch string) (Arch, error) {
	switch strings.ToLower(arch) {
	case "amd64", "x86_64", "x64":
		return AMD64, nil
	case "arm64", "aarch64":
		return ARM64, nil
	default:
		return 0, ErrUnknownArch
	}
}

// ParseOS parses a string into an OS type
func ParseOS(os string) (OS, error) {
	switch strings.ToLower(os) {
	case "linux":
		return Linux, nil
	case "darwin", "macos":
		return Darwin, nil
	case "windows":
		return Windows, nil
	default:
		return 0, ErrUnknownOS
	}
}

// ProgressCallback is called during download to report progress
type ProgressCallback func(url string, bytesComplete, totalBytes int64, mbps float64, done bool)

// ProgressTracker is a default progress callback that prints to stdout
var ProgressTracker ProgressCallback = func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
	if done {
		fmt.Printf("\n✅ Download complete: %s\n", filepath.Base(url))
		return
	}

	// Handle case when totalBytes is unknown (server doesn't send Content-Length)
	if totalBytes <= 0 {
		fmt.Printf("\r⬇️  Downloading: %.2f MB (%.2f MB/s)", float64(bytesComplete)/(1024*1024), mbps)
		return
	}

	percent := float64(bytesComplete) / float64(totalBytes) * 100
	fmt.Printf("\r⬇️  Downloading: %.1f%% (%.2f MB/s)", percent, mbps)
}

// SDLatestVersion fetches the latest release tag of stable-diffusion.cpp from the GitHub API.
func SDLatestVersion() (string, error) {
	var version string
	var err error
	for range RetryCount {
		version, err = getLatestVersion()
		if err == nil {
			return version, nil
		}
		time.Sleep(RetryDelay)
	}

	return "", errors.New("unable to fetch latest version")
}

func getLatestVersion() (string, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	// Set required headers for GitHub API
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("received status code %d from GitHub API: %s", resp.StatusCode, string(body))
	}

	var result struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.TagName, nil
}

// getDownloadLocationAndFilename returns the download location and filename for the given parameters.
func getDownloadLocationAndFilename(arch Arch, os OS, version string) (location, filename string, err error) {
	location = fmt.Sprintf("https://github.com/leejet/stable-diffusion.cpp/releases/download/%s", version)

	// Extract commit hash from version (e.g., "master-487-43e829f" -> "43e829f")
	parts := strings.Split(version, "-")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid version format: %s", version)
	}
	commitHash := parts[len(parts)-1] // Get last part (commit hash)

	switch os {
	case Linux:
		if arch == ARM64 {
			return "", "", errors.New("precompiled binaries for Linux ARM64 are not available")
		}
		filename = fmt.Sprintf("sd-master-%s-bin-ubuntu-x64.zip", commitHash)

	case Darwin:
		if arch == ARM64 {
			filename = fmt.Sprintf("sd-master-%s-bin-Darwin-macOS-15.7.3-arm64.zip", commitHash)
		} else {
			filename = fmt.Sprintf("sd-master-%s-bin-Darwin-macOS-x64.zip", commitHash)
		}

	case Windows:
		if arch == ARM64 {
			return "", "", errors.New("precompiled binaries for Windows ARM64 are not available")
		}
		filename = fmt.Sprintf("sd-master-%s-bin-win-x64.zip", commitHash)

	default:
		return "", "", ErrUnknownOS
	}

	return location, filename, nil
}

// Get downloads the stable-diffusion.cpp precompiled binaries for the current system.
// version should be the desired version tag (e.g., "master-487-43e829f").
// If version is empty, it will use DefaultVersion (master-487-43e829f).
// You can use [SDLatestVersion] to obtain the latest release.
// If dest is empty, it will use the default lib directory.
func Get(version string) error {
	if version == "" {
		version = DefaultVersion
	}
	return GetWithProgress(version, "", ProgressTracker)
}

// GetWithProgress downloads the stable-diffusion.cpp precompiled binaries with progress callback.
func GetWithProgress(version string, dest string, progress ProgressCallback) error {
	return GetWithContext(context.Background(), version, dest, progress)
}

// GetWithContext downloads the stable-diffusion.cpp precompiled binaries using the provided context.
// If version is empty, it will use DefaultVersion (master-487-43e829f).
func GetWithContext(ctx context.Context, version string, dest string, progress ProgressCallback) error {
	arch, err := ParseArch(runtime.GOARCH)
	if err != nil {
		return ErrUnknownArch
	}

	os, err := ParseOS(runtime.GOOS)
	if err != nil {
		return ErrUnknownOS
	}

	if version == "" {
		version = DefaultVersion
	}

	location, filename, err := getDownloadLocationAndFilename(arch, os, version)
	if err != nil {
		return err
	}

	// Use default destination if not provided
	if dest == "" {
		dest = "lib"
	}

	url := fmt.Sprintf("%s/%s", location, filename)
	return downloadAndExtractZip(ctx, url, dest, progress)
}

// downloadAndExtractZip downloads a .zip file and extracts it to the destination directory.
func downloadAndExtractZip(ctx context.Context, url, dest string, progress ProgressCallback) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	downloadFile := filepath.Join(dest, filepath.Base(url))

	// Download using grab with resume support
	req, err := grab.NewRequest(downloadFile, url)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	req = req.WithContext(ctx)

	client := grab.NewClient()
	resp := client.Do(req)

	// Monitor progress
	if progress != nil {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ticker.C:
					if resp.IsComplete() {
						return
					}
					progress(url, resp.BytesComplete(), resp.Size(), resp.BytesPerSecond()/(1024*1024), false)
				case <-resp.Done:
					return
				}
			}
		}()
	}

	// Wait for download to complete
	if err := resp.Err(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if progress != nil {
		progress(url, resp.BytesComplete(), resp.Size(), resp.BytesPerSecond()/(1024*1024), true)
	}

	defer func() {
		_ = os.Remove(downloadFile)
	}()

	// Open the downloaded zip file
	zipReader, err := zip.OpenReader(downloadFile)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		if err := zipReader.Close(); err != nil {
			log.Printf("failed to close zip reader: %v", err)
		}
	}()

	// Extract files
	for _, file := range zipReader.File {
		// Get the file path
		filePath := filepath.Join(dest, file.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", filePath)
		}

		if file.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Open the file in the zip
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		// Create the destination file
		dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			_ = srcFile.Close()
			return fmt.Errorf("failed to create file: %w", err)
		}

		// Copy contents
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			_ = srcFile.Close()
			_ = dstFile.Close()
			return fmt.Errorf("failed to write file: %w", err)
		}

		// Close files
		if err := srcFile.Close(); err != nil {
			_ = dstFile.Close()
			return fmt.Errorf("failed to close source file: %w", err)
		}
		if err := dstFile.Close(); err != nil {
			return fmt.Errorf("failed to close destination file: %w", err)
		}
	}

	return nil
}
