package stablediffusion

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/pkg/stablediffusion/download"
)

// GetLibDir returns the directory where stable-diffusion library should be stored.
// Uses paths.Libraries() to ensure consistent location for both development and packaged apps.
func GetLibDir() string {
	return paths.Libraries()
}

// EnsureLibrary checks if the stable-diffusion library exists, and downloads it if not.
// This should be called once during application startup.
func EnsureLibrary() error {
	return EnsureLibraryWithProgress(download.ProgressTracker)
}

// EnsureLibraryWithProgress checks if the stable-diffusion library exists with custom progress callback.
func EnsureLibraryWithProgress(progress download.ProgressCallback) error {
	libPath := GetLibDir()
	libName := download.LibraryName(runtime.GOOS)
	libFile := filepath.Join(libPath, libName)

	// Check if library already exists
	if _, err := os.Stat(libFile); err == nil {
		// Library exists
		return nil
	}

	slog.Info("Stable Diffusion library not found, downloading", "size", "~18MB")

	// Download library
	err := download.GetWithProgress(download.DefaultVersion, libPath, progress)
	if err != nil {
		return fmt.Errorf("failed to download library: %w", err)
	}

	slog.Info("Library setup complete", "path", libFile)
	return nil
}

// IsLibraryInstalled checks if the stable-diffusion library is installed.
func IsLibraryInstalled() bool {
	libPath := GetLibDir()
	libName := download.LibraryName(runtime.GOOS)
	libFile := filepath.Join(libPath, libName)

	_, err := os.Stat(libFile)
	return err == nil
}

// GetLibraryPath returns the path to the stable-diffusion library.
func GetLibraryPath() string {
	libPath := GetLibDir()
	libName := download.LibraryName(runtime.GOOS)
	return filepath.Join(libPath, libName)
}

// GetLibraryVersion returns the default library version.
func GetLibraryVersion() string {
	return download.DefaultVersion
}
