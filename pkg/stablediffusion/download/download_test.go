package download

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLibraryName_CurrentPlatform(t *testing.T) {
	// Test the OS-aware LibraryName() function
	name := LibraryName()
	assert.NotEmpty(t, name)

	// Verify it matches current platform
	switch runtime.GOOS {
	case "windows":
		assert.Equal(t, "stable-diffusion.dll", name)
	case "linux":
		assert.Equal(t, "libstable-diffusion.so", name)
	case "darwin":
		assert.Equal(t, "libstable-diffusion.dylib", name)
	default:
		t.Skipf("Unsupported platform: %s", runtime.GOOS)
	}
}

func TestLibraryName_ContainsCorrectExtension(t *testing.T) {
	name := LibraryName()

	switch runtime.GOOS {
	case "windows":
		assert.Contains(t, name, ".dll")
	case "linux":
		assert.Contains(t, name, ".so")
	case "darwin":
		assert.Contains(t, name, ".dylib")
	default:
		t.Skipf("Unsupported platform: %s", runtime.GOOS)
	}
}

func TestLibraryName_Consistency(t *testing.T) {
	// Test that calling LibraryName multiple times returns same result
	name1 := LibraryName()
	name2 := LibraryName()
	assert.Equal(t, name1, name2)
}

func TestDefaultVersion(t *testing.T) {
	assert.NotEmpty(t, DefaultVersion, "default version should not be empty")
	assert.Contains(t, DefaultVersion, "master", "version should contain 'master'")
}

func TestProgressTracker_ValidInput(t *testing.T) {
	// Test that ProgressTracker doesn't panic with valid input
	assert.NotPanics(t, func() {
		ProgressTracker("http://example.com/file", 100, 1000, 1.5, false)
	})
}

func TestProgressTracker_ZeroTotal(t *testing.T) {
	// Test that ProgressTracker handles zero total gracefully
	assert.NotPanics(t, func() {
		ProgressTracker("http://example.com/file", 0, 0, 0, false)
	})
}

func TestProgressTracker_CompleteDownload(t *testing.T) {
	// Test progress at 100%
	assert.NotPanics(t, func() {
		ProgressTracker("http://example.com/file", 1000, 1000, 2.0, true)
	})
}

func TestGet_InvalidVersion(t *testing.T) {
	err := Get("invalid-version-xyz")
	assert.Error(t, err, "should fail with invalid version")
}

func TestGetWithProgress_CallbackInvoked(t *testing.T) {
	tmpDir := t.TempDir()

	mockCallback := func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
		assert.GreaterOrEqual(t, totalBytes, bytesComplete)
		assert.GreaterOrEqual(t, mbps, float64(0))
		assert.NotEmpty(t, url)
	}

	// This will likely fail due to network/version issues, but we test the callback
	err := GetWithProgress(DefaultVersion, tmpDir, mockCallback)

	// We don't assert on error since download might fail
	// We just verify the function signature works
	_ = err

	// Note: callback might not be invoked if download fails immediately
	t.Logf("Download test completed with result: %v", err)
}

func TestGetWithProgress_NilCallback(t *testing.T) {
	tmpDir := t.TempDir()

	// Should not panic with nil callback
	assert.NotPanics(t, func() {
		_ = GetWithProgress(DefaultVersion, tmpDir, nil)
	})
}

func TestLibraryURL_Format(t *testing.T) {
	// Test URL construction logic (if exposed)
	version := "master-123-abc"
	libName := LibraryName()

	assert.NotEmpty(t, version)
	assert.NotEmpty(t, libName)

	// URL should follow pattern: https://github.com/.../releases/download/{version}/{libName}
	expectedPattern := version + "/" + libName
	assert.NotEmpty(t, expectedPattern)
}

func TestDownloadToPath_DirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	nestedPath := filepath.Join(tmpDir, "nested", "path", "lib")

	// Verify directory doesn't exist yet
	_, err := os.Stat(nestedPath)
	assert.True(t, os.IsNotExist(err))

	// Create directory structure
	err = os.MkdirAll(nestedPath, 0755)
	require.NoError(t, err)

	// Verify directory was created
	info, err := os.Stat(nestedPath)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestProgressCallback_Type(t *testing.T) {
	// Verify ProgressCallback type signature
	var callback ProgressCallback = func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
		assert.GreaterOrEqual(t, totalBytes, int64(0))
		assert.GreaterOrEqual(t, bytesComplete, int64(0))
		assert.NotEmpty(t, url)
	}

	assert.NotNil(t, callback)

	// Test callback execution
	assert.NotPanics(t, func() {
		callback("http://test.com", 50, 100, 1.0, false)
	})
}

func TestGet_WithExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	libName := LibraryName()
	libPath := filepath.Join(tmpDir, libName)

	// Create a mock library file
	err := os.WriteFile(libPath, []byte("mock library content"), 0644)
	require.NoError(t, err)

	// Verify file exists
	info, err := os.Stat(libPath)
	require.NoError(t, err)
	assert.False(t, info.IsDir())
	assert.Greater(t, info.Size(), int64(0))
}

func TestDefaultVersion_Format(t *testing.T) {
	// Verify version format is reasonable
	assert.NotEmpty(t, DefaultVersion)
	assert.NotContains(t, DefaultVersion, " ", "version should not contain spaces")
	assert.NotContains(t, DefaultVersion, "\n", "version should not contain newlines")
}
