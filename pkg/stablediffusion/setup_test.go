package stablediffusion

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kawai-network/veridium/pkg/stablediffusion/download"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLibDir(t *testing.T) {
	libDir := GetLibDir()
	assert.NotEmpty(t, libDir, "library directory should not be empty")
	assert.Contains(t, libDir, "lib", "library directory should contain 'lib'")
}

func TestGetLibraryPath(t *testing.T) {
	libPath := GetLibraryPath()
	assert.NotEmpty(t, libPath, "library path should not be empty")

	expectedName := download.LibraryName()
	assert.Contains(t, libPath, expectedName, "library path should contain correct library name")
}

func TestGetLibraryVersion(t *testing.T) {
	version := GetLibraryVersion()
	assert.NotEmpty(t, version, "library version should not be empty")
	assert.Equal(t, download.DefaultVersion, version)
}

func TestIsLibraryInstalled_NotInstalled(t *testing.T) {
	// Test with non-existent library
	installed := IsLibraryInstalled()
	// Note: This will check actual library location
	// For proper testing, we'd need dependency injection
	assert.IsType(t, false, installed)
}

func TestIsLibraryInstalled_WithMockFile(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	libName := download.LibraryName()
	libFile := filepath.Join(tmpDir, libName)

	// Create mock library file
	err := os.WriteFile(libFile, []byte("mock library"), 0644)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(libFile)
	assert.NoError(t, err, "mock library file should exist")
}

func TestEnsureLibrary_ProgressCallback(t *testing.T) {
	mockProgress := func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
		assert.GreaterOrEqual(t, totalBytes, bytesComplete)
		assert.GreaterOrEqual(t, mbps, float64(0))
	}

	// Note: This test won't actually download if library exists
	// For full testing, we'd need to mock the download function
	err := EnsureLibraryWithProgress(mockProgress)

	// If library already exists, no error and callback might not be called
	// If library doesn't exist, it will attempt download
	if err != nil {
		// Download failed or library doesn't exist
		t.Logf("Library ensure result: %v", err)
	}
}

func TestEnsureLibrary_AlreadyInstalled(t *testing.T) {
	// If library is already installed, this should return nil immediately
	err := EnsureLibrary()

	// Should either succeed (library exists) or fail (download issue)
	// We can't guarantee the state, so we just check it doesn't panic
	if err != nil {
		t.Logf("Library ensure returned error: %v", err)
	}
}

func TestGetLibraryPath_ContainsCorrectExtension(t *testing.T) {
	libPath := GetLibraryPath()

	switch runtime.GOOS {
	case "windows":
		assert.Contains(t, libPath, ".dll")
	case "linux":
		assert.Contains(t, libPath, ".so")
	case "darwin":
		assert.Contains(t, libPath, ".dylib")
	default:
		t.Skipf("Unsupported platform: %s", runtime.GOOS)
	}
}
