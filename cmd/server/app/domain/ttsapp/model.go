package ttsapp

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/pkg/tools/downloader"
)

// DefaultTTSModelURL is the default TTS model (Kokoro) from HuggingFace
const DefaultTTSModelURL = "https://huggingface.co/mmwillet2/Kokoro_GGUF/resolve/main/Kokoro_no_espeak_Q4.gguf"

// DefaultTTSModelName is the default filename for the TTS model
const DefaultTTSModelName = "Kokoro_no_espeak_Q4.gguf"

// DefaultTTSModelOrg is the default organization for TTS model
const DefaultTTSModelOrg = "mmwillet2"

// DefaultTTSModelRepo is the default repository for TTS model
const DefaultTTSModelRepo = "Kokoro_GGUF"

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

// modelFilePathAndName extracts the standard model path and filename from a HuggingFace URL
// Following the standard: models/{org}/{repo}/{filename}
func modelFilePathAndName(modelFileURL string) (string, string, error) {
	mURL, err := url.Parse(modelFileURL)
	if err != nil {
		return "", "", fmt.Errorf("unable to parse fileURL: %w", err)
	}

	parts := strings.Split(mURL.Path, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid huggingface url: %q", mURL.Path)
	}

	fileName, err := extractFileName(modelFileURL)
	if err != nil {
		return "", "", fmt.Errorf("unable to extract file name: %w", err)
	}

	// Standard path: models/{org}/{repo}/{filename}
	// parts: ["", "mmwillet2", "Kokoro_GGUF", "resolve", "main", "Kokoro_no_espeak_Q4.gguf"]
	modelFilePath := filepath.Join(paths.Models(), parts[1], parts[2])
	modelFileName := filepath.Join(modelFilePath, fileName)

	return modelFilePath, modelFileName, nil
}

// extractFileName extracts the filename from a URL
func extractFileName(modelFileURL string) (string, error) {
	u, err := url.Parse(modelFileURL)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	name := path.Base(u.Path)
	return name, nil
}

// DiscoverModel checks if a TTS model already exists
// Uses the standard path structure: models/{org}/{repo}/{filename}
func (d *ModelDownloader) DiscoverModel() (string, error) {
	// Check standard path: models/{org}/{repo}/
	standardPath := filepath.Join(paths.Models(), DefaultTTSModelOrg, DefaultTTSModelRepo)
	entries, err := os.ReadDir(standardPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read models directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gguf" {
			return filepath.Join(standardPath, entry.Name()), nil
		}
	}

	return "", nil
}

// DownloadModel downloads a TTS model from the given URL
// Uses the standard path structure: models/{org}/{repo}/{filename}
func (d *ModelDownloader) DownloadModel(ctx context.Context, modelURL string, progressCb func(bytesComplete, totalBytes int64, mbps float64, done bool)) (string, error) {
	if modelURL == "" {
		modelURL = DefaultTTSModelURL
	}

	// Extract standard path from URL
	modelFilePath, modelFileName, err := modelFilePathAndName(modelURL)
	if err != nil {
		return "", fmt.Errorf("failed to extract model path: %w", err)
	}

	// Create models directory if it doesn't exist
	if err := os.MkdirAll(modelFilePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create models directory: %w", err)
	}

	// Check if model already exists
	if _, err := os.Stat(modelFileName); err == nil {
		return modelFileName, nil
	}

	// Download the model
	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		if progressCb != nil {
			progressCb(currentSize, totalSize, mibPerSec, complete)
		}
	}

	// Pass the full file path (modelFileName) to downloader.Download
	_, err = downloader.Download(ctx, modelURL, modelFileName, progressFunc, downloader.SizeIntervalMIB)
	if err != nil {
		return "", fmt.Errorf("failed to download TTS model: %w", err)
	}

	return modelFileName, nil
}

// GetDefaultModelPath returns the path to the default TTS model
// Uses the standard path structure: models/{org}/{repo}/{filename}
func (d *ModelDownloader) GetDefaultModelPath() string {
	return filepath.Join(paths.Models(), DefaultTTSModelOrg, DefaultTTSModelRepo, DefaultTTSModelName)
}

// ListDownloadedModels returns a list of downloaded TTS models
// Uses the standard path structure: models/{org}/{repo}/
func ListDownloadedModels() ([]string, error) {
	standardPath := filepath.Join(paths.Models(), DefaultTTSModelOrg, DefaultTTSModelRepo)
	entries, err := os.ReadDir(standardPath)
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
	downloader := NewModelDownloader("")

	progressWrapper := func(currentBytes, totalBytes int64, mbps float64, done bool) {
		if progressCb != nil {
			progressCb(currentBytes, totalBytes)
		}
	}

	_, err := downloader.DownloadModel(ctx, modelURL, progressWrapper)
	return err
}
