package download

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kawai-network/grab"
)

var (
	ErrUnknownArch    = errors.New("unknown architecture")
	ErrUnknownOS      = errors.New("unknown OS")
	ErrInvalidVersion = errors.New("invalid version")
)

// DefaultVersion is the default gosd release version to use.
var DefaultVersion = "v0.1.4"

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

// getDownloadLocationAndFilename returns the download location and filename for the given parameters.
func getDownloadLocationAndFilename(arch Arch, os OS, version string) (location, filename string, err error) {
	if version == "" || !strings.HasPrefix(version, "v") {
		return "", "", fmt.Errorf("%w: expected semantic version tag like v0.1.4, got %q", ErrInvalidVersion, version)
	}
	location = fmt.Sprintf("https://github.com/getkawai/stablediffusion/releases/download/%s", version)

	switch os {
	case Linux:
		if arch == ARM64 {
			return "", "", errors.New("precompiled binaries for Linux ARM64 are not available")
		}
		filename = "libgosd-linux.tar.gz"

	case Darwin:
		filename = "libgosd-macos.tar.gz"

	case Windows:
		if arch == ARM64 {
			return "", "", errors.New("precompiled binaries for Windows ARM64 are not available")
		}
		filename = "libgosd-windows.zip"

	default:
		return "", "", ErrUnknownOS
	}

	return location, filename, nil
}

// Get downloads gosd precompiled binaries for the current system.
// version should be the desired release tag (e.g., "v0.1.4").
// If version is empty, it will use DefaultVersion.
// If dest is empty, it will use the default lib directory.
func Get(version string) error {
	if version == "" {
		version = DefaultVersion
	}
	return GetWithProgress(version, "", ProgressTracker)
}

// GetWithProgress downloads gosd precompiled binaries with progress callback.
func GetWithProgress(version string, dest string, progress ProgressCallback) error {
	return GetWithContext(context.Background(), version, dest, progress)
}

// GetWithContext downloads gosd precompiled binaries using the provided context.
// If version is empty, it will use DefaultVersion.
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
	return downloadAndExtractArchive(ctx, url, dest, progress)
}

// downloadAndExtractArchive downloads an archive (.zip or .tar.gz) and extracts it.
func downloadAndExtractArchive(ctx context.Context, url, dest string, progress ProgressCallback) error {
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

	// Check if resume occurred and report initial progress
	if resp.DidResume && progress != nil {
		progress(url, resp.BytesComplete(), resp.Size(), resp.BytesPerSecond()/(1024*1024), false)
	}

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

	switch {
	case strings.HasSuffix(downloadFile, ".zip"):
		return extractZip(downloadFile, dest)
	case strings.HasSuffix(downloadFile, ".tar.gz"), strings.HasSuffix(downloadFile, ".tgz"):
		return extractTarGz(downloadFile, dest)
	default:
		return fmt.Errorf("unsupported archive format: %s", filepath.Base(downloadFile))
	}
}

func extractZip(downloadFile, dest string) error {
	zipReader, err := zip.OpenReader(downloadFile)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		if err := zipReader.Close(); err != nil {
			log.Printf("failed to close zip reader: %v", err)
		}
	}()

	for _, file := range zipReader.File {
		filePath := filepath.Join(dest, file.Name)
		if !isSafeExtractPath(dest, filePath) {
			return fmt.Errorf("illegal file path: %s", filePath)
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}
		dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			_ = srcFile.Close()
			return fmt.Errorf("failed to create file: %w", err)
		}
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			_ = srcFile.Close()
			_ = dstFile.Close()
			return fmt.Errorf("failed to write file: %w", err)
		}
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

func extractTarGz(downloadFile, dest string) error {
	f, err := os.Open(downloadFile)
	if err != nil {
		return fmt.Errorf("failed to open tar.gz file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close tar.gz file: %v", err)
		}
	}()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		if err := gzr.Close(); err != nil {
			log.Printf("failed to close gzip reader: %v", err)
		}
	}()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed reading tar entry: %w", err)
		}

		filePath := filepath.Join(dest, header.Name)
		if !isSafeExtractPath(dest, filePath) {
			return fmt.Errorf("illegal file path: %s", filePath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
			dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(dstFile, tr); err != nil {
				_ = dstFile.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			if err := dstFile.Close(); err != nil {
				return fmt.Errorf("failed to close destination file: %w", err)
			}
		}
	}

	return nil
}

func isSafeExtractPath(dest, filePath string) bool {
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)
	return strings.HasPrefix(filePath, cleanDest)
}
