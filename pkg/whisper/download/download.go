package download

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	ErrUnknownOS = errors.New("unknown OS")
)

// DefaultVersion is the default whisper.cpp version to use
// This matches the version in third_party/whisper.cpp (CMakeLists.txt: project("whisper.cpp" VERSION 1.8.3))
const DefaultVersion = "v1.8.3"

// ProgressCallback is called during download to report progress
type ProgressCallback func(url string, bytesComplete, totalBytes int64, mbps float64, done bool)

// ProgressTracker is a default progress callback that prints to stdout
var ProgressTracker ProgressCallback = func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
	if done {
		fmt.Printf("\n✅ Download complete: %s\n", filepath.Base(url))
		return
	}

	if totalBytes <= 0 {
		fmt.Printf("\r⬇️  Downloading: %.2f MB (%.2f MB/s)", float64(bytesComplete)/(1024*1024), mbps)
		return
	}

	percent := float64(bytesComplete) / float64(totalBytes) * 100
	fmt.Printf("\r⬇️  Downloading: %.1f%% (%.2f MB/s)", percent, mbps)
}

// LibraryName returns the platform-specific library name
func LibraryName(goos string) string {
	switch goos {
	case "darwin":
		return "libwhisper.dylib"
	case "windows":
		return "whisper.dll"
	case "linux":
		return "libwhisper.so"
	default:
		return "unknown"
	}
}

// GetDownloadURL returns the download URL for whisper.cpp library
// Uses pre-built libraries from whisper.cpp releases
func GetDownloadURL(version, goos, arch string) (string, error) {
	// whisper.cpp releases provide different binaries per platform
	// Format: https://github.com/ggerganov/whisper.cpp/releases/download/v1.7.4/whisper-bin-osx-arm64.zip
	
	baseURL := fmt.Sprintf("https://github.com/ggerganov/whisper.cpp/releases/download/%s", version)
	
	switch goos {
	case "darwin":
		if arch == "arm64" {
			return fmt.Sprintf("%s/whisper-bin-osx-arm64.zip", baseURL), nil
		}
		return fmt.Sprintf("%s/whisper-bin-osx-x64.zip", baseURL), nil
	case "linux":
		// Linux binaries are usually x64
		return fmt.Sprintf("%s/whisper-bin-ubuntu-x64.zip", baseURL), nil
	case "windows":
		return fmt.Sprintf("%s/whisper-bin-win-x64.zip", baseURL), nil
	default:
		return "", ErrUnknownOS
	}
}

// Get downloads the whisper.cpp library to the specified directory
func Get(version, destDir string) error {
	return GetWithProgress(version, destDir, nil)
}

// GetWithProgress downloads with progress callback
func GetWithProgress(version, destDir string, progress ProgressCallback) error {
	if version == "" {
		version = DefaultVersion
	}

	goos := runtime.GOOS
	arch := runtime.GOARCH

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Get download URL
	url, err := GetDownloadURL(version, goos, arch)
	if err != nil {
		return err
	}

	// Download to temp file
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("whisper-%s-%s.zip", goos, arch))
	
	if err := downloadFile(url, tempFile, progress); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(tempFile)

	// Extract library
	libName := LibraryName(goos)
	destPath := filepath.Join(destDir, libName)

	if err := extractLibrary(tempFile, destDir, libName); err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}

	// Make executable (Unix)
	if goos != "windows" {
		os.Chmod(destPath, 0755)
	}

	return nil
}

// downloadFile downloads a file with progress tracking
func downloadFile(url, dest string, progress ProgressCallback) error {
	client := &http.Client{Timeout: 30 * time.Minute}
	
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Track progress
	var downloaded int64
	start := time.Now()
	total := resp.ContentLength
	
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			out.Write(buf[:n])
			downloaded += int64(n)
			
			if progress != nil {
				elapsed := time.Since(start).Seconds()
				mbps := float64(downloaded) / (1024 * 1024) / elapsed
				progress(url, downloaded, total, mbps, false)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	if progress != nil {
		progress(url, downloaded, total, 0, true)
	}

	return nil
}

// extractLibrary extracts the specific library file from archive
func extractLibrary(archivePath, destDir, libName string) error {
	// Check if file already exists
	destPath := filepath.Join(destDir, libName)
	if _, err := os.Stat(destPath); err == nil {
		// Already exists
		return nil
	}

	// Open zip archive
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		// Try to copy directly if not a zip (might be single file)
		data, err := os.ReadFile(archivePath)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, 0755)
	}
	defer reader.Close()

	// Find and extract library file
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		
		name := filepath.Base(file.Name)
		if strings.HasSuffix(name, ".dylib") || strings.HasSuffix(name, ".so") || strings.HasSuffix(name, ".dll") {
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			out, err := os.Create(destPath)
			if err != nil {
				return err
			}
			
			_, err = io.Copy(out, rc)
			out.Close()
			
			if err != nil {
				return err
			}
			
			return os.Chmod(destPath, 0755)
		}
	}

	return fmt.Errorf("library file not found in archive")
}
