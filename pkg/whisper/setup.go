package whisper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/pkg/whisper/whisper"
	whisperpkg "github.com/kawai-network/whisper"
)

// GetLibDir returns the directory where whisper library should be stored
func GetLibDir() string {
	return paths.WhisperLib()
}

func init() {
	// Set library directory for whisper package
	whisper.SetLibraryDirectory(GetLibDir())
}

// EnsureLibrary checks if the whisper library exists, and downloads it if not
func EnsureLibrary() error {
	return EnsureLibraryWithProgress(nil)
}

// EnsureLibraryWithProgress checks if the whisper library exists with custom progress callback
func EnsureLibraryWithProgress(progress func(string, int64, int64, float64, bool)) error {
	libPath := GetLibDir()
	libName := whisperpkg.LibraryName(runtime.GOOS)
	libFile := filepath.Join(libPath, libName)

	// Check if library already exists
	if _, err := os.Stat(libFile); err == nil {
		// Library exists, try to load it
		if err := whisper.Load(libPath); err == nil {
			return nil
		}
		// If load fails, re-download
	}

	// Download library using external package
	downloader := whisperpkg.NewLibraryDownloader(libPath)

	_, err := downloader.DownloadLatest()
	if err != nil {
		return fmt.Errorf("failed to download library: %w", err)
	}

	// Load the library
	if err := whisper.Load(libPath); err != nil {
		return fmt.Errorf("failed to load library: %w", err)
	}

	return nil
}

// IsLibraryInstalled checks if the whisper library is installed
func IsLibraryInstalled() bool {
	libPath := GetLibDir()
	libName := whisperpkg.LibraryName(runtime.GOOS)
	libFile := filepath.Join(libPath, libName)

	_, err := os.Stat(libFile)
	return err == nil
}

// GetLibraryPath returns the path to the whisper library
func GetLibraryPath() string {
	libPath := GetLibDir()
	libName := whisperpkg.LibraryName(runtime.GOOS)
	return filepath.Join(libPath, libName)
}

// GetLibraryVersion returns the latest library version from GitHub
func GetLibraryVersion() string {
	downloader := whisperpkg.NewLibraryDownloader("")
	release, err := downloader.GetLatestRelease()
	if err != nil {
		return "unknown"
	}
	return release.TagName
}
