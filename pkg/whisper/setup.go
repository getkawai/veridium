package whisper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/pkg/whisper/download"
	"github.com/kawai-network/veridium/pkg/whisper/whisper"
)

// GetLibDir returns the directory where whisper library should be stored
func GetLibDir() string {
	return filepath.Join(paths.Base(), "lib", "whisper")
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
func EnsureLibraryWithProgress(progress download.ProgressCallback) error {
	libPath := GetLibDir()
	libName := download.LibraryName(runtime.GOOS)
	libFile := filepath.Join(libPath, libName)

	// Check if library already exists
	if _, err := os.Stat(libFile); err == nil {
		// Library exists, try to load it
		if err := whisper.Load(libPath); err == nil {
			return nil
		}
		// If load fails, re-download
	}

	// Download library
	if err := download.GetWithProgress(download.DefaultVersion, libPath, progress); err != nil {
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
	libName := download.LibraryName(runtime.GOOS)
	libFile := filepath.Join(libPath, libName)

	_, err := os.Stat(libFile)
	return err == nil
}

// GetLibraryPath returns the path to the whisper library
func GetLibraryPath() string {
	libPath := GetLibDir()
	libName := download.LibraryName(runtime.GOOS)
	return filepath.Join(libPath, libName)
}

// GetLibraryVersion returns the default library version
func GetLibraryVersion() string {
	return download.DefaultVersion
}
