package libs

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
	"github.com/kawai-network/veridium/pkg/tools/downloader"
)

func init() {
	registerDownloader(LibraryTTS, &ttsDownloader{})
}

// ttsDownloader handles downloading TTS.cpp libraries from GitHub releases.
// TTS.cpp releases are distributed as ZIP archives containing multiple files:
// - libtts_c_api.{so,dylib,dll} (C API library)
// - libtts.{so,dylib,dll} (main library)
// - libggml*.{so,dylib,dll} (ggml dependencies)
// - tts_c_api.h (header file)
type ttsDownloader struct{}

func (d *ttsDownloader) LatestVersion() (string, error) {
	// TTS.cpp uses GitHub releases, fetch latest tag
	// For now, return hardcoded version or implement GitHub API fetch
	// TODO: Implement proper version fetching from GitHub API
	return "v0.1.1", nil
}

func (d *ttsDownloader) Download(ctx context.Context, arch, osName, processor, version, dest string, progress download.ProgressCallback) error {
	if version == "" {
		version = "v0.1.1"
	}

	// Determine asset name based on OS and arch
	assetName, libName := d.getAssetNames(osName, arch)
	if assetName == "" {
		return fmt.Errorf("unsupported platform: %s/%s", osName, arch)
	}

	// Construct download URL
	zipURL := fmt.Sprintf("https://github.com/kawai-network/TTS.cpp/releases/download/%s/%s", version, assetName)
	zipPath := filepath.Join(dest, "tts-download.zip")

	// Download the ZIP file
	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		if progress != nil {
			progress(src, currentSize, totalSize, mibPerSec, complete)
		}
	}

	_, err := downloader.Download(ctx, zipURL, zipPath, progressFunc, downloader.SizeIntervalMIB)
	if err != nil {
		return fmt.Errorf("failed to download TTS library: %w", err)
	}
	defer os.Remove(zipPath)

	// Extract the ZIP file
	if err := d.extractZip(zipPath, dest); err != nil {
		return fmt.Errorf("failed to extract TTS library: %w", err)
	}

	// Verify the main library file exists
	libPath := filepath.Join(dest, libName)
	if _, err := os.Stat(libPath); err != nil {
		return fmt.Errorf("library file not found after extraction: %s", libName)
	}

	return nil
}

func (d *ttsDownloader) LibraryName(os string) string {
	_, libName := d.getAssetNames(os, "")
	return libName
}

// getAssetNames returns the asset name and library name for the given platform.
func (d *ttsDownloader) getAssetNames(os, arch string) (assetName, libName string) {
	switch os {
	case "linux":
		// Default to gcc build for Linux
		return "tts-c-api-linux-gcc.zip", "libtts_c_api.so"
	case "darwin":
		// macOS only has arm64 builds currently
		return "tts-c-api-macos-arm64.zip", "libtts_c_api.dylib"
	case "windows":
		return "tts-c-api-windows-msvc.zip", "tts_c_api.dll"
	default:
		return "", ""
	}
}

// extractZip extracts a ZIP file to the destination directory.
func (d *ttsDownloader) extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, f := range r.File {
		// Skip directories
		if f.FileInfo().IsDir() {
			continue
		}

		// Extract only library files and headers
		name := f.Name
		if !d.isRequiredFile(name) {
			continue
		}

		// Open file inside ZIP
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip %s: %w", name, err)
		}

		// Create destination file
		dstPath := filepath.Join(destDir, filepath.Base(name))
		out, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create file %s: %w", dstPath, err)
		}

		// Copy content
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", name, err)
		}
	}

	return nil
}

// isRequiredFile checks if a file should be extracted from the ZIP.
func (d *ttsDownloader) isRequiredFile(name string) bool {
	// Only extract library files and headers
	// Skip subdirectories in the ZIP (files are in tts-c-api-*/ subfolder)
	base := filepath.Base(name)

	// Check if it's a library file or header
	if filepath.Ext(base) == ".h" {
		return true
	}

	ext := filepath.Ext(base)
	return ext == ".so" || ext == ".dylib" || ext == ".dll"
}
