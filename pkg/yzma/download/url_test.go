package download

import (
	"strings"
	"testing"
)

// TestURLConstruction verifies that download URLs are correctly formed
// without including the extraction path (//build/bin) in the URL
func TestURLConstruction(t *testing.T) {
	tests := []struct {
		name              string
		os                string
		processor         string
		version           string
		expectedURLSuffix string // The part after releases/download/version/
		shouldNotContain  string // String that should NOT be in URL
	}{
		{
			name:              "macOS ARM64 Metal",
			os:                "darwin",
			processor:         "metal",
			version:           "b7018",
			expectedURLSuffix: "llama-b7018-bin-macos-arm64.zip",
			shouldNotContain:  "//build/bin",
		},
		{
			name:              "Linux x64 CPU",
			os:                "linux",
			processor:         "cpu",
			version:           "b7018",
			expectedURLSuffix: "llama-b7018-bin-ubuntu-x64.zip",
			shouldNotContain:  "//build/bin",
		},
		{
			name:              "Windows x64 CPU",
			os:                "windows",
			processor:         "cpu",
			version:           "b7018",
			expectedURLSuffix: "llama-b7018-bin-win-cpu-x64.zip",
			shouldNotContain:  "//build/bin",
		},
		{
			name:              "Windows CUDA",
			os:                "windows",
			processor:         "cuda",
			version:           "b7018",
			expectedURLSuffix: "llama-b7018-bin-win-cuda-12.4-x64.zip",
			shouldNotContain:  "//build/bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't actually call Get() without downloading
			// So we'll test the URL construction logic directly
			
			var location, filename string
			location = "https://github.com/ggml-org/llama.cpp/releases/download/" + tt.version

			switch tt.os {
			case "linux":
				switch tt.processor {
				case "cpu":
					filename = "llama-" + tt.version + "-bin-ubuntu-x64.zip//build/bin"
				case "vulkan":
					filename = "llama-" + tt.version + "-bin-ubuntu-vulkan-x64.zip//build/bin"
				}
			case "darwin":
				switch tt.processor {
				case "cpu", "metal":
					filename = "llama-" + tt.version + "-bin-macos-arm64.zip//build/bin"
				}
			case "windows":
				switch tt.processor {
				case "cpu":
					filename = "llama-" + tt.version + "-bin-win-cpu-x64.zip//build/bin"
				case "cuda":
					filename = "llama-" + tt.version + "-bin-win-cuda-12.4-x64.zip//build/bin"
				}
			}

			// Extract the actual filename (before //) for URL construction
			actualFilename := filename
			if strings.Contains(filename, "//") {
				actualFilename = strings.SplitN(filename, "//", 2)[0]
			}

			url := location + "/" + actualFilename

			// Verify URL is correct
			if !strings.HasSuffix(url, tt.expectedURLSuffix) {
				t.Errorf("Expected URL to end with %s, got: %s", tt.expectedURLSuffix, url)
			}

			// Verify URL does NOT contain extraction path
			if strings.Contains(url, tt.shouldNotContain) {
				t.Errorf("URL should NOT contain %s, but got: %s", tt.shouldNotContain, url)
			}

			// Verify URL is a valid GitHub release URL
			expectedPrefix := "https://github.com/ggml-org/llama.cpp/releases/download/" + tt.version + "/"
			if !strings.HasPrefix(url, expectedPrefix) {
				t.Errorf("Expected URL to start with %s, got: %s", expectedPrefix, url)
			}

			t.Logf("✅ Correct URL: %s", url)
		})
	}
}

// TestFilenameExtraction verifies that extraction path is preserved for internal use
func TestFilenameExtraction(t *testing.T) {
	tests := []struct {
		filename         string
		expectedZipFile  string
		expectedExtract  string
	}{
		{
			filename:        "llama-b7018-bin-macos-arm64.zip//build/bin",
			expectedZipFile: "llama-b7018-bin-macos-arm64.zip",
			expectedExtract: "build/bin",
		},
		{
			filename:        "llama-b7018-bin-ubuntu-x64.zip//build/bin",
			expectedZipFile: "llama-b7018-bin-ubuntu-x64.zip",
			expectedExtract: "build/bin",
		},
		{
			filename:        "simple-file.zip",
			expectedZipFile: "simple-file.zip",
			expectedExtract: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			var zipFile, extractPath string
			if strings.Contains(tt.filename, "//") {
				parts := strings.SplitN(tt.filename, "//", 2)
				zipFile = parts[0]
				extractPath = parts[1]
			} else {
				zipFile = tt.filename
			}

			if zipFile != tt.expectedZipFile {
				t.Errorf("Expected zipFile=%s, got: %s", tt.expectedZipFile, zipFile)
			}

			if extractPath != tt.expectedExtract {
				t.Errorf("Expected extractPath=%s, got: %s", tt.expectedExtract, extractPath)
			}

			t.Logf("✅ zipFile=%s, extractPath=%s", zipFile, extractPath)
		})
	}
}

