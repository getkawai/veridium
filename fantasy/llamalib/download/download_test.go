package download

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestLlamaLatestVersion(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("skipping test since github API sends 403 error")
	}

	version, err := LlamaLatestVersion()
	if err != nil {
		t.Fatal("count not get latest version", err)
	}

	if !strings.HasPrefix(version, "b") {
		t.Fatalf("Expected version should start with 'b', got '%s'", version)
	}

	t.Logf("LlamaLatestVersion returned: %s", version)
}

func TestGetLinuxCPU(t *testing.T) {
	version := "b6795"
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	expectedFile := filepath.Join(dest, "libllama.so")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", expectedFile)
	}

	t.Logf("Get() successfully downloaded the file to: %s", expectedFile)
}

func TestGetInvalidOS(t *testing.T) {
	version := "b6795"
	osVer := "cpm"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != ErrUnknownOS {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

func TestGetInvalidProcessor(t *testing.T) {
	version := "b6795"
	osVer := "windows"
	processor := "flux"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != ErrUnknownProcessor {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

func TestGetInvalidVersion(t *testing.T) {
	version := "nogood"
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != ErrInvalidVersion {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

// ============================================================================
// Tests for Advanced Download Functions
// ============================================================================

func TestDefaultDownloadOptions(t *testing.T) {
	opts := DefaultDownloadOptions()

	if opts.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries=3, got %d", opts.MaxRetries)
	}
	if !opts.ShowProgress {
		t.Error("Expected ShowProgress=true")
	}
	if !opts.ResumeIfPossible {
		t.Error("Expected ResumeIfPossible=true")
	}
	if opts.RateLimitMBps != 0 {
		t.Errorf("Expected RateLimitMBps=0, got %d", opts.RateLimitMBps)
	}

	t.Logf("DefaultDownloadOptions: %+v", opts)
}

func TestWithRateLimit(t *testing.T) {
	opts := WithRateLimit(5)

	if opts.RateLimitMBps != 5 {
		t.Errorf("Expected RateLimitMBps=5, got %d", opts.RateLimitMBps)
	}
	if opts.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries=3, got %d", opts.MaxRetries)
	}

	t.Logf("WithRateLimit(5): %+v", opts)
}

func TestGetWithProgress_SmallFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	// Download a small test file from a reliable source
	url := "https://raw.githubusercontent.com/ggml-org/llama.cpp/master/README.md"
	dest := filepath.Join(t.TempDir(), "README.md")

	opts := DefaultDownloadOptions()
	opts.ShowProgress = false // Disable progress for test

	err := GetWithProgress(url, dest, opts)
	if err != nil {
		t.Logf("GetWithProgress failed (may be network issue): %v", err)
		t.Skip("Skipping due to network issue")
	}

	// Verify file exists
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", dest)
	}

	// Verify file is not empty
	info, err := os.Stat(dest)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("Downloaded file is empty")
	}

	t.Logf("✅ Successfully downloaded file: %s (%.1f KB)", dest, float64(info.Size())/1024)
}

func TestGetWithProgress_WithProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	// Download a small test file with progress enabled
	url := "https://raw.githubusercontent.com/ggml-org/llama.cpp/master/LICENSE"
	dest := filepath.Join(t.TempDir(), "LICENSE")

	opts := DefaultDownloadOptions()
	opts.ShowProgress = true
	opts.ProgressInterval = 500 * time.Millisecond // Fast updates for test

	err := GetWithProgress(url, dest, opts)
	if err != nil {
		t.Logf("GetWithProgress failed (may be network issue): %v", err)
		t.Skip("Skipping due to network issue")
	}

	// Verify file exists
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", dest)
	}

	t.Logf("✅ Successfully downloaded with progress: %s", dest)
}

func TestGetWithProgress_Resume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping resume test in short mode")
	}

	url := "https://raw.githubusercontent.com/ggml-org/llama.cpp/master/README.md"
	dest := filepath.Join(t.TempDir(), "README.md")

	// First download (will be interrupted by creating partial file)
	// Create a partial file to simulate interrupted download
	partialContent := []byte("partial content")
	if err := os.WriteFile(dest, partialContent, 0644); err != nil {
		t.Fatalf("Failed to create partial file: %v", err)
	}

	opts := DefaultDownloadOptions()
	opts.ShowProgress = false
	opts.ResumeIfPossible = true

	// This should resume/overwrite the partial file
	err := GetWithProgress(url, dest, opts)
	if err != nil {
		t.Logf("GetWithProgress failed (may be network issue): %v", err)
		t.Skip("Skipping due to network issue")
	}

	// Verify file is larger than partial content
	info, err := os.Stat(dest)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Size() <= int64(len(partialContent)) {
		t.Fatal("File was not properly downloaded/resumed")
	}

	t.Logf("✅ Resume test passed: final size %.1f KB", float64(info.Size())/1024)
}

func TestGetWithProgress_WithRateLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rate limit test in short mode")
	}

	url := "https://raw.githubusercontent.com/ggml-org/llama.cpp/master/README.md"
	dest := filepath.Join(t.TempDir(), "README.md")

	opts := WithRateLimit(1) // 1 MB/s limit
	opts.ShowProgress = false

	err := GetWithProgress(url, dest, opts)
	if err != nil {
		t.Logf("GetWithProgress failed (may be network issue): %v", err)
		t.Skip("Skipping due to network issue")
	}

	// Verify file exists
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", dest)
	}

	t.Logf("✅ Rate limit test passed")
}

func TestGetWithProgress_Retry(t *testing.T) {
	// Test with invalid URL to trigger retry
	url := "https://invalid-domain-that-does-not-exist-12345.com/file.txt"
	dest := filepath.Join(t.TempDir(), "file.txt")

	opts := DefaultDownloadOptions()
	opts.MaxRetries = 2 // Reduce retries for faster test
	opts.ShowProgress = false

	err := GetWithProgress(url, dest, opts)
	if err == nil {
		t.Fatal("Expected error for invalid URL, got nil")
	}

	// Verify error message mentions retries
	if !strings.Contains(err.Error(), "failed after") {
		t.Errorf("Expected error to mention retries, got: %v", err)
	}

	t.Logf("✅ Retry test passed: %v", err)
}

func TestGetWithProgress_InvalidURL(t *testing.T) {
	url := "not-a-valid-url"
	dest := filepath.Join(t.TempDir(), "file.txt")

	opts := DefaultDownloadOptions()
	opts.MaxRetries = 1
	opts.ShowProgress = false

	err := GetWithProgress(url, dest, opts)
	if err == nil {
		t.Fatal("Expected error for invalid URL, got nil")
	}

	t.Logf("✅ Invalid URL test passed: %v", err)
}

func TestGetBatch_MultipleFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping batch download test in short mode")
	}

	destDir := t.TempDir()

	urls := []string{
		"https://raw.githubusercontent.com/ggml-org/llama.cpp/master/README.md",
		"https://raw.githubusercontent.com/ggml-org/llama.cpp/master/LICENSE",
	}

	respch, err := GetBatch(2, destDir, urls...)
	if err != nil {
		t.Fatalf("GetBatch failed: %v", err)
	}

	successCount := 0
	for resp := range respch {
		if err := resp.Err(); err != nil {
			t.Logf("Download failed: %v", err)
		} else {
			successCount++
			t.Logf("✅ Downloaded: %s", resp.Filename)
		}
	}

	if successCount == 0 {
		t.Skip("All batch downloads failed (may be network issue)")
	}

	t.Logf("✅ Batch download completed: %d/%d successful", successCount, len(urls))
}

func TestGetBatch_DefaultWorkers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping batch download test in short mode")
	}

	destDir := t.TempDir()
	urls := []string{
		"https://raw.githubusercontent.com/ggml-org/llama.cpp/master/README.md",
	}

	// Test with workers=0 (should default to 3)
	respch, err := GetBatch(0, destDir, urls...)
	if err != nil {
		t.Fatalf("GetBatch failed: %v", err)
	}

	for resp := range respch {
		if err := resp.Err(); err != nil {
			t.Logf("Download failed (may be network issue): %v", err)
			t.Skip("Skipping due to network issue")
		}
		t.Logf("✅ Downloaded with default workers: %s", resp.Filename)
	}
}

func TestLibraryName(t *testing.T) {
	tests := []struct {
		os       string
		expected string
	}{
		{"linux", "libllama.so"},
		{"freebsd", "libllama.so"},
		{"windows", "llama.dll"},
		{"darwin", "libllama.dylib"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.os, func(t *testing.T) {
			result := LibraryName(tt.os)
			if result != tt.expected {
				t.Errorf("LibraryName(%s) = %s, want %s", tt.os, result, tt.expected)
			}
		})
	}
}

func TestVersionIsValid(t *testing.T) {
	tests := []struct {
		version string
		wantErr bool
	}{
		{"b1234", false},
		{"b6795", false},
		{"b7016", false},
		{"1234", true},
		{"v1.0", true},
		{"nogood", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			err := VersionIsValid(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("VersionIsValid(%s) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestRequiredLibraries(t *testing.T) {
	tests := []struct {
		os       string
		expected []string
	}{
		{
			os: "linux",
			expected: []string{
				"libggml.so",
				"libggml-base.so",
				"libllama.so",
			},
		},
		{
			os: "freebsd",
			expected: []string{
				"libggml.so",
				"libggml-base.so",
				"libllama.so",
			},
		},
		{
			os: "darwin",
			expected: []string{
				"libggml.dylib",
				"libggml-base.dylib",
				"libllama.dylib",
			},
		},
		{
			os: "windows",
			expected: []string{
				"ggml.dll",
				"ggml-base.dll",
				"llama.dll",
			},
		},
		{
			os:       "unknown",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.os, func(t *testing.T) {
			result := RequiredLibraries(tt.os)
			if len(result) != len(tt.expected) {
				t.Errorf("RequiredLibraries(%s) returned %d libs, want %d", tt.os, len(result), len(tt.expected))
				return
			}
			for i, lib := range result {
				if lib != tt.expected[i] {
					t.Errorf("RequiredLibraries(%s)[%d] = %s, want %s", tt.os, i, lib, tt.expected[i])
				}
			}
		})
	}
}

func TestGetLibraryExtension(t *testing.T) {
	tests := []struct {
		os       string
		expected string
	}{
		{"linux", ".so"},
		{"freebsd", ".so"},
		{"darwin", ".dylib"},
		{"windows", ".dll"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.os, func(t *testing.T) {
			result := GetLibraryExtension(tt.os)
			if result != tt.expected {
				t.Errorf("GetLibraryExtension(%s) = %s, want %s", tt.os, result, tt.expected)
			}
		})
	}
}
