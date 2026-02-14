// Package whisperapp provides helper functions for managing whisper models.
package whisperapp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/tools/downloader"
)

// WhisperModelSpec represents a Whisper model specification for standalone whisper
type WhisperModelSpec struct {
	Name        string
	URL         string
	Size        int64
	Parameters  string
	MinRAM      int64
	EnglishOnly bool
	Description string
}

// Model specifications for standalone whisper (gowhisper)
// Note: These models use the same .bin format as whisper.cpp but with different naming convention
var modelSpecs = []WhisperModelSpec{
	{
		Name:        "tiny",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin",
		Size:        77691776, // ~74MB
		Parameters:  "39M",
		MinRAM:      1,
		EnglishOnly: false,
		Description: "Tiny model - fastest, lowest quality, multilingual",
	},
	{
		Name:        "tiny.en",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.en.bin",
		Size:        77691776, // ~74MB
		Parameters:  "39M",
		MinRAM:      1,
		EnglishOnly: true,
		Description: "Tiny model - fastest, lowest quality, English-only",
	},
	{
		Name:        "base",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin",
		Size:        147964544, // ~141MB
		Parameters:  "74M",
		MinRAM:      2,
		EnglishOnly: false,
		Description: "Base model - fast, good quality, multilingual",
	},
	{
		Name:        "base.en",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin",
		Size:        147964544, // ~141MB
		Parameters:  "74M",
		MinRAM:      2,
		EnglishOnly: true,
		Description: "Base model - fast, good quality, English-only",
	},
	{
		Name:        "small",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin",
		Size:        487620096, // ~465MB
		Parameters:  "244M",
		MinRAM:      4,
		EnglishOnly: false,
		Description: "Small model - balanced speed/quality, multilingual",
	},
	{
		Name:        "small.en",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.en.bin",
		Size:        487620096, // ~465MB
		Parameters:  "244M",
		MinRAM:      4,
		EnglishOnly: true,
		Description: "Small model - balanced speed/quality, English-only",
	},
	{
		Name:        "medium",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.bin",
		Size:        1533121536, // ~1.5GB
		Parameters:  "769M",
		MinRAM:      8,
		EnglishOnly: false,
		Description: "Medium model - slower, better quality, multilingual",
	},
	{
		Name:        "medium.en",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.en.bin",
		Size:        1533121536, // ~1.5GB
		Parameters:  "769M",
		MinRAM:      8,
		EnglishOnly: true,
		Description: "Medium model - slower, better quality, English-only",
	},
	{
		Name:        "large-v1",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v1.bin",
		Size:        3094623232, // ~2.9GB
		Parameters:  "1550M",
		MinRAM:      16,
		EnglishOnly: false,
		Description: "Large v1 model - slowest, best quality, multilingual",
	},
	{
		Name:        "large-v2",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v2.bin",
		Size:        3094623232, // ~2.9GB
		Parameters:  "1550M",
		MinRAM:      16,
		EnglishOnly: false,
		Description: "Large v2 model - slowest, best quality, multilingual",
	},
	{
		Name:        "large-v3",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v3.bin",
		Size:        3094623232, // ~2.9GB
		Parameters:  "1550M",
		MinRAM:      16,
		EnglishOnly: false,
		Description: "Large v3 model - slowest, best quality, multilingual",
	},
	{
		Name:        "large-v3-turbo",
		URL:         "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v3-turbo.bin",
		Size:        1623451264, // ~1.5GB
		Parameters:  "809M",
		MinRAM:      8,
		EnglishOnly: false,
		Description: "Large v3 turbo - faster than large-v3, good quality, multilingual",
	},
}

// ProgressCallback is called during download to report progress
type ProgressCallback func(currentBytes, totalBytes int64)

// DownloadModel downloads a whisper model to the specified directory
// Note: Standalone whisper expects models as {name}.bin (not ggml-{name}.bin)
func DownloadModel(ctx context.Context, modelName, modelsDir string, progress ProgressCallback) error {
	spec, exists := GetModelSpec(modelName)
	if !exists {
		return fmt.Errorf("unknown model: %s", modelName)
	}

	// Ensure models directory exists
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Destination file path (standalone whisper uses {name}.bin format)
	destPath := filepath.Join(modelsDir, fmt.Sprintf("%s.bin", modelName))

	// Check if model already exists
	if info, err := os.Stat(destPath); err == nil {
		// Verify file size matches expected size
		if info.Size() == spec.Size {
			return nil // Model already downloaded and correct size
		}
		// File exists but wrong size, re-download
		os.Remove(destPath)
	}

	// Download using centralized downloader with progress tracking
	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		if progress != nil {
			progress(currentSize, totalSize)
		}
	}

	downloaded, err := downloader.Download(ctx, spec.URL, destPath, progressFunc, downloader.SizeIntervalMIB10)
	if err != nil {
		os.Remove(destPath) // Cleanup on failure
		return fmt.Errorf("failed to download model: %w", err)
	}

	if !downloaded {
		os.Remove(destPath) // Cleanup empty file
		return fmt.Errorf("download completed but no data was transferred")
	}

	return nil
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

// SelectOptimalModel selects the best model based on available RAM
// Returns the largest model that fits within available RAM
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

// GetModelPath returns the path to a model file
// For standalone whisper, this is {modelsDir}/{name}.bin
func GetModelPath(modelsDir, modelName string) string {
	return filepath.Join(modelsDir, fmt.Sprintf("%s.bin", modelName))
}

// IsModelDownloaded checks if a model is already downloaded
func IsModelDownloaded(modelsDir, modelName string) bool {
	path := GetModelPath(modelsDir, modelName)
	spec, exists := GetModelSpec(modelName)
	if !exists {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Verify file size matches expected size
	return info.Size() == spec.Size
}

// ListDownloadedModels returns a list of downloaded model names
func ListDownloadedModels(modelsDir string) ([]string, error) {
	if _, err := os.Stat(modelsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	var models []string

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, ".bin") {
			// Extract model name from "{name}.bin"
			modelName := strings.TrimSuffix(name, ".bin")
			// Verify it's a valid model
			if _, exists := GetModelSpec(modelName); exists {
				models = append(models, modelName)
			}
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

// DownloadProgressLogger is a helper that logs download progress
type DownloadProgressLogger struct {
	LastLog time.Time
}

// Log logs progress at reasonable intervals
func (l *DownloadProgressLogger) Log(current, total int64) {
	now := time.Now()
	if now.Sub(l.LastLog) < 2*time.Second && current != total {
		return
	}
	l.LastLog = now

	if total > 0 {
		percent := float64(current) / float64(total) * 100
		fmt.Printf("Download progress: %s / %s (%.1f%%)\n",
			HumanSize(current), HumanSize(total), percent)
	} else {
		fmt.Printf("Downloaded: %s\n", HumanSize(current))
	}
}

// DownloadModelWithLogger downloads a model with automatic progress logging
func DownloadModelWithLogger(ctx context.Context, modelName, modelsDir string) error {
	logger := &DownloadProgressLogger{}

	// Get model spec for initial info
	spec, exists := GetModelSpec(modelName)
	if !exists {
		return fmt.Errorf("unknown model: %s", modelName)
	}

	fmt.Printf("Downloading model '%s' (%s)...\n", modelName, spec.Description)
	fmt.Printf("Model size: %s\n", HumanSize(spec.Size))

	return DownloadModel(ctx, modelName, modelsDir, logger.Log)
}
