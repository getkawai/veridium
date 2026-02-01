package whisper

import (
	"testing"

	"github.com/kawai-network/veridium/pkg/whisper/download"
	"github.com/kawai-network/veridium/pkg/whisper/whisper"
	"github.com/stretchr/testify/assert"
)

func TestCString(t *testing.T) {
	tests := []string{
		"hello",
		"world",
		"test with spaces",
		"",
	}

	for _, tc := range tests {
		// Test CString conversion
		cstr := whisper.CString(tc)
		if tc != "" {
			assert.NotNil(t, cstr)
			// Verify null termination by reading back
			result := whisper.GoString(cstr)
			assert.Equal(t, tc, result)
		}
	}
}

func TestGoString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.input == "" {
				assert.Equal(t, "", whisper.GoString(nil))
			} else {
				cstr := whisper.CString(tc.input)
				result := whisper.GoString(cstr)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestGetLibDir(t *testing.T) {
	dir := GetLibDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, "whisper")
}

func TestLibraryName(t *testing.T) {
	tests := []struct {
		goos     string
		expected string
	}{
		{"darwin", "libwhisper.dylib"},
		{"linux", "libwhisper.so"},
		{"windows", "whisper.dll"},
		{"freebsd", "unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.goos, func(t *testing.T) {
			result := download.LibraryName(tc.goos)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsLibraryInstalled(t *testing.T) {
	// This will return false unless library is actually installed
	result := IsLibraryInstalled()
	// Just verify it doesn't panic
	assert.False(t, result) // Should be false in test environment
}

func TestGetLibraryVersion(t *testing.T) {
	version := GetLibraryVersion()
	assert.NotEmpty(t, version)
}

func TestGetDownloadURL(t *testing.T) {
	tests := []struct {
		version string
		goos    string
		arch    string
		wantErr bool
	}{
		{"v1.7.4", "darwin", "arm64", false},
		{"v1.7.4", "darwin", "amd64", false},
		{"v1.7.4", "linux", "amd64", false},
		{"v1.7.4", "windows", "amd64", false},
		{"v1.7.4", "freebsd", "amd64", true},
	}

	for _, tc := range tests {
		t.Run(tc.goos+"_"+tc.arch, func(t *testing.T) {
			url, err := download.GetDownloadURL(tc.version, tc.goos, tc.arch)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, url, "github.com")
				assert.Contains(t, url, tc.version)
			}
		})
	}
}
