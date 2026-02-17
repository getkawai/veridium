package ttsapp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/pkg/tools/downloader"
)

// DefaultTTSModelURL is the default TTS model (Kokoro) from HuggingFace
const DefaultTTSModelURL = "https://huggingface.co/mmwillet2/Kokoro_GGUF/resolve/main/kokoro-v1.0-82M-Q4_K_M.gguf"

// ModelDownloader handles TTS model downloads
type ModelDownloader struct {
	modelsPath string
}

// NewModelDownloader creates a new TTS model downloader
func NewModelDownloader(modelsPath string) *ModelDownloader {
	return &ModelDownloader{
		modelsPath: modelsPath,
	}
}

// DiscoverModel checks if a TTS model already exists
func (d *ModelDownloader) DiscoverModel() (string, error) {
	// Check for any .gguf files in the models directory
	entries, err := os.ReadDir(d.modelsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read models directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gguf" {
			return filepath.Join(d.modelsPath, entry.Name()), nil
		}
	}

	return "", nil
}

// DownloadModel downloads a TTS model from the given URL
func (d *ModelDownloader) DownloadModel(ctx context.Context, url string, progressCb func(bytesComplete, totalBytes int64, mbps float64, done bool)) (string, error) {
	if url == "" {
		url = DefaultTTSModelURL
	}

	// Extract filename from URL
	filename := filepath.Base(url)
	if filename == "" || filename == "." {
		filename = DefaultTTSModelName
	}

	// Create models directory if it doesn't exist
	if err := os.MkdirAll(d.modelsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create models directory: %w", err)
	}

	modelPath := filepath.Join(d.modelsPath, filename)

	// Check if model already exists
	if _, err := os.Stat(modelPath); err == nil {
		return modelPath, nil
	}

	// Download the model
	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		if progressCb != nil {
			progressCb(currentSize, totalSize, mibPerSec, complete)
		}
	}

	_, err := downloader.Download(ctx, url, modelPath, progressFunc, downloader.SizeIntervalMIB)
	if err != nil {
		return "", fmt.Errorf("failed to download TTS model: %w", err)
	}

	return modelPath, nil
}

// GetDefaultModelPath returns the path to the default TTS model
func (d *ModelDownloader) GetDefaultModelPath() string {
	return filepath.Join(d.modelsPath, DefaultTTSModelName)
}

// ListDownloadedModels returns a list of downloaded TTS models
func ListDownloadedModels() ([]string, error) {
	modelsPath := paths.Models()
	ttsModelsPath := filepath.Join(modelsPath, "tts")

	entries, err := os.ReadDir(ttsModelsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read TTS models directory: %w", err)
	}

	var models []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gguf" {
			models = append(models, entry.Name())
		}
	}

	return models, nil
}

// DownloadModelWithProgress downloads a TTS model with progress callback
func DownloadModelWithProgress(ctx context.Context, modelURL string, progressCb func(currentBytes, totalBytes int64)) error {
	modelsPath := paths.Models()
	ttsModelsPath := filepath.Join(modelsPath, "tts")

	downloader := NewModelDownloader(ttsModelsPath)

	progressWrapper := func(currentBytes, totalBytes int64, mbps float64, done bool) {
		if progressCb != nil {
			progressCb(currentBytes, totalBytes)
		}
	}

	_, err := downloader.DownloadModel(ctx, modelURL, progressWrapper)
	return err
}
