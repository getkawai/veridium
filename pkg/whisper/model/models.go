// Package model provides whisper model specifications and downloading capabilities.
//
// This package contains:
//   - Model specifications for all official Whisper models
//   - Model downloader with progress tracking
//   - Auto-selection based on available RAM
//
// Usage:
//
//	// Download a specific model
//	err := model.DownloadModel("base", "./models", nil)
//
//	// Auto-select and download best model
//	spec := model.SelectOptimalModel(8) // 8GB RAM
//	err := model.DownloadModel(spec.Name, "./models", progressCallback)
//
// Available models: tiny, tiny.en, base, base.en, small, small.en, medium, medium.en, large-v1, large-v2, large-v3, large-v3-turbo
package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// WhisperModelSpec represents a Whisper model specification
type WhisperModelSpec struct {
	Name         string // Model name (e.g., "base.en")
	URL          string // Direct download URL
	Size         int64  // Model size in bytes
	Parameters   string // Parameter count (e.g., "74M")
	MinRAM       int64  // Minimum RAM required in GB
	EnglishOnly  bool   // Whether model is English-only
	Description  string // Model description
}

// HuggingFace base URL for whisper models
const huggingFaceBaseURL = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main"

// Model specifications (ordered from smallest to largest)
var modelSpecs = []WhisperModelSpec{
	{
		Name:        "tiny",
		URL:         huggingFaceBaseURL + "/ggml-tiny.bin",
		Size:        77691776, // ~74MB
		Parameters:  "39M",
		MinRAM:      1,
		EnglishOnly: false,
		Description: "Tiny model - fastest, lowest quality, multilingual",
	},
	{
		Name:        "tiny.en",
		URL:         huggingFaceBaseURL + "/ggml-tiny.en.bin",
		Size:        77691776, // ~74MB
		Parameters:  "39M",
		MinRAM:      1,
		EnglishOnly: true,
		Description: "Tiny model - fastest, lowest quality, English-only",
	},
	{
		Name:        "base",
		URL:         huggingFaceBaseURL + "/ggml-base.bin",
		Size:        147964544, // ~141MB
		Parameters:  "74M",
		MinRAM:      2,
		EnglishOnly: false,
		Description: "Base model - fast, good quality, multilingual",
	},
	{
		Name:        "base.en",
		URL:         huggingFaceBaseURL + "/ggml-base.en.bin",
		Size:        147964544, // ~141MB
		Parameters:  "74M",
		MinRAM:      2,
		EnglishOnly: true,
		Description: "Base model - fast, good quality, English-only",
	},
	{
		Name:        "small",
		URL:         huggingFaceBaseURL + "/ggml-small.bin",
		Size:        487620096, // ~465MB
		Parameters:  "244M",
		MinRAM:      4,
		EnglishOnly: false,
		Description: "Small model - balanced speed/quality, multilingual",
	},
	{
		Name:        "small.en",
		URL:         huggingFaceBaseURL + "/ggml-small.en.bin",
		Size:        487620096, // ~465MB
		Parameters:  "244M",
		MinRAM:      4,
		EnglishOnly: true,
		Description: "Small model - balanced speed/quality, English-only",
	},
	{
		Name:        "medium",
		URL:         huggingFaceBaseURL + "/ggml-medium.bin",
		Size:        1533121536, // ~1.5GB
		Parameters:  "769M",
		MinRAM:      8,
		EnglishOnly: false,
		Description: "Medium model - slower, better quality, multilingual",
	},
	{
		Name:        "medium.en",
		URL:         huggingFaceBaseURL + "/ggml-medium.en.bin",
		Size:        1533121536, // ~1.5GB
		Parameters:  "769M",
		MinRAM:      8,
		EnglishOnly: true,
		Description: "Medium model - slower, better quality, English-only",
	},
	{
		Name:        "large-v1",
		URL:         huggingFaceBaseURL + "/ggml-large-v1.bin",
		Size:        3094623232, // ~2.9GB
		Parameters:  "1550M",
		MinRAM:      16,
		EnglishOnly: false,
		Description: "Large v1 model - slowest, best quality, multilingual",
	},
	{
		Name:        "large-v2",
		URL:         huggingFaceBaseURL + "/ggml-large-v2.bin",
		Size:        3094623232, // ~2.9GB
		Parameters:  "1550M",
		MinRAM:      16,
		EnglishOnly: false,
		Description: "Large v2 model - slowest, best quality, multilingual",
	},
	{
		Name:        "large-v3",
		URL:         huggingFaceBaseURL + "/ggml-large-v3.bin",
		Size:        3094623232, // ~2.9GB
		Parameters:  "1550M",
		MinRAM:      16,
		EnglishOnly: false,
		Description: "Large v3 model - slowest, best quality, multilingual",
	},
	{
		Name:        "large-v3-turbo",
		URL:         huggingFaceBaseURL + "/ggml-large-v3-turbo.bin",
		Size:        1623451264, // ~1.5GB
		Parameters:  "809M",
		MinRAM:      8,
		EnglishOnly: false,
		Description: "Large v3 turbo - faster than large-v3, good quality, multilingual",
	},
}

// GetModelSpec returns the model specification for a given model name
func GetModelSpec(name string) (*WhisperModelSpec, bool) {
	for i := range modelSpecs {
		if modelSpecs[i].Name == name {
			return &modelSpecs[i], true
		}
	}
	return nil, false
}

// GetAllModels returns all available model specifications
func GetAllModels() []WhisperModelSpec {
	return modelSpecs
}

// GetAvailableModels returns models that fit within the given RAM (in GB)
func GetAvailableModels(availableRAM int64) []WhisperModelSpec {
	var available []WhisperModelSpec
	for _, spec := range modelSpecs {
		if spec.MinRAM <= availableRAM {
			available = append(available, spec)
		}
	}
	return available
}

// SelectOptimalModel selects the best model based on available RAM.
// Returns the largest model that fits within available RAM.
func SelectOptimalModel(availableRAM int64) *WhisperModelSpec {
	var selected *WhisperModelSpec
	for i := range modelSpecs {
		if modelSpecs[i].MinRAM <= availableRAM {
			selected = &modelSpecs[i]
		}
	}
	// Fallback to tiny if nothing fits
	if selected == nil {
		return &modelSpecs[0]
	}
	return selected
}

// SelectOptimalModelEnglish selects the best English-only model based on available RAM
func SelectOptimalModelEnglish(availableRAM int64) *WhisperModelSpec {
	var selected *WhisperModelSpec
	for i := range modelSpecs {
		if modelSpecs[i].EnglishOnly && modelSpecs[i].MinRAM <= availableRAM {
			selected = &modelSpecs[i]
		}
	}
	// Fallback to first English model
	if selected == nil {
		for i := range modelSpecs {
			if modelSpecs[i].EnglishOnly {
				return &modelSpecs[i]
			}
		}
	}
	return selected
}

// ProgressCallback is called during download to report progress
type ProgressCallback func(bytesComplete, totalBytes int64)

// DownloadModel downloads a whisper model to the specified directory.
// If progress is nil, no progress updates will be sent.
func DownloadModel(modelName, modelsDir string, progress ProgressCallback) error {
	spec, exists := GetModelSpec(modelName)
	if !exists {
		return fmt.Errorf("unknown model: %s", modelName)
	}

	// Ensure models directory exists
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Build file paths
	modelFileName := fmt.Sprintf("ggml-%s.bin", spec.Name)
	destPath := filepath.Join(modelsDir, modelFileName)
	tempPath := destPath + ".tmp"

	// Check if model already exists
	if _, err := os.Stat(destPath); err == nil {
		return nil // Model already exists
	}

	// Download to temp file
	if err := downloadFile(spec.URL, tempPath, spec.Size, progress); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to download model: %w", err)
	}

	// Move temp file to final destination
	if err := os.Rename(tempPath, destPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to move downloaded file: %w", err)
	}

	return nil
}

// downloadFile downloads a file with progress tracking
func downloadFile(url, destPath string, expectedSize int64, progress ProgressCallback) error {
	// Create temp file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get content length
	totalSize := resp.ContentLength
	if totalSize < 0 {
		totalSize = expectedSize
	}

	// Copy with progress tracking
	var downloaded int64
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)
			if progress != nil {
				progress(downloaded, totalSize)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// GetModelPath returns the expected path for a model file
func GetModelPath(modelsDir, modelName string) string {
	return filepath.Join(modelsDir, fmt.Sprintf("ggml-%s.bin", modelName))
}

// IsModelDownloaded checks if a model is already downloaded
func IsModelDownloaded(modelsDir, modelName string) bool {
	path := GetModelPath(modelsDir, modelName)
	_, err := os.Stat(path)
	return err == nil
}

// ListDownloadedModels returns a list of downloaded model names
func ListDownloadedModels(modelsDir string) ([]string, error) {
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var models []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, "ggml-") && strings.HasSuffix(name, ".bin") {
			// Extract model name from "ggml-{name}.bin"
			modelName := strings.TrimPrefix(name, "ggml-")
			modelName = strings.TrimSuffix(modelName, ".bin")
			models = append(models, modelName)
		}
	}
	return models, nil
}

// DeleteModel deletes a downloaded model
func DeleteModel(modelsDir, modelName string) error {
	path := GetModelPath(modelsDir, modelName)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("model not found: %s", modelName)
		}
		return err
	}
	return nil
}

// HumanSize returns a human-readable size string
func HumanSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
