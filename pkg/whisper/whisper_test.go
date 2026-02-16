package whisper

import (
	"testing"

	"github.com/kawai-network/veridium/pkg/whisper/whisper"
	whisperpkg "github.com/kawai-network/whisper"
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
		{"darwin", "libgowhisper.dylib"},
		{"linux", "libgowhisper.so"},
		{"windows", "gowhisper.dll"},
	}

	for _, tc := range tests {
		t.Run(tc.goos, func(t *testing.T) {
			result := whisperpkg.LibraryName(tc.goos)
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
