// Package modeldownloader provides functionality to download Stable Diffusion models.
// This package handles model discovery, download with progress tracking, and validation.
package modeldownloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kawai-network/veridium/pkg/tools/downloader"
)

const (
	// DefaultModelURL is the default Stable Diffusion model URL (SD 1.4 ~4GB)
	DefaultModelURL = "https://huggingface.co/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt"
	// DefaultModelFilename is the default filename for the downloaded model
	DefaultModelFilename = "sd-v1-4.ckpt"
)

// ProgressCallback is called during download to report progress
type ProgressCallback func(bytesComplete, totalBytes int64, mbps float64, done bool)

// ModelDownloader handles model discovery and download
type ModelDownloader struct {
	modelsPath string
}

// New creates a new ModelDownloader instance
func New(modelsPath string) *ModelDownloader {
	return &ModelDownloader{
		modelsPath: modelsPath,
	}
}

// NewWithClient creates a new ModelDownloader (kept for backward compatibility)
// Note: client parameter is ignored as we now use pkg/tools/downloader
func NewWithClient(modelsPath string, client interface{}) *ModelDownloader {
	return &ModelDownloader{
		modelsPath: modelsPath,
	}
}

// DiscoverModel searches for existing model files in the models directory
// Returns the path to the first found model file, or empty string if none found
func (m *ModelDownloader) DiscoverModel() (string, error) {
	var modelFile string

	err := filepath.Walk(m.modelsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			name := info.Name()
			if strings.HasSuffix(name, ".safetensors") ||
				strings.HasSuffix(name, ".ckpt") ||
				strings.HasSuffix(name, ".gguf") {
				modelFile = path
				return filepath.SkipAll // Found one
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error walking models directory: %w", err)
	}

	return modelFile, nil
}

// DownloadModel downloads a model from the given URL to the models directory
// Organizes models by {author}/{repo}/ structure from HuggingFace URLs
// Uses grab package for resume support and progress tracking
func (m *ModelDownloader) DownloadModel(ctx context.Context, url string, progress ProgressCallback) (string, error) {
	// Extract author/repo from HuggingFace URL and filename
	// URL format: https://huggingface.co/{author}/{repo}/resolve/main/{filename}
	modelDir := m.modelsPath
	filename := DefaultModelFilename

	if strings.Contains(url, "huggingface.co") {
		parts := strings.Split(url, "/")
		for i, part := range parts {
			if strings.Contains(part, "huggingface.co") && i+2 < len(parts) {
				author := parts[i+1]
				repo := parts[i+2]

				// Validate author and repo to prevent path traversal
				if author == "" || repo == "" {
					return "", fmt.Errorf("empty author or repo in URL")
				}
				if strings.Contains(author, "..") || strings.Contains(repo, "..") ||
					strings.ContainsAny(author, "/\\") || strings.ContainsAny(repo, "/\\") {
					return "", fmt.Errorf("invalid author or repo in URL: path traversal attempt detected")
				}

				modelDir = filepath.Join(m.modelsPath, author, repo)

				// Extract filename from URL (last part) and strip query parameters
				if len(parts) > 0 {
					rawFilename := parts[len(parts)-1]
					// Strip query parameters (e.g., ?token=xxx)
					if idx := strings.Index(rawFilename, "?"); idx != -1 {
						rawFilename = rawFilename[:idx]
					}
					// Strip fragment (e.g., #section)
					if idx := strings.Index(rawFilename, "#"); idx != -1 {
						rawFilename = rawFilename[:idx]
					}
					if rawFilename != "" {
						filename = rawFilename
					}
				}
				break
			}
		}
	}

	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create model directory: %w", err)
	}

	modelDest := filepath.Join(modelDir, filename)

	// Use centralized downloader with resume support
	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		if progress != nil {
			progress(currentSize, totalSize, mibPerSec, complete)
		}
	}

	// Use 100MB interval for large SD models (typically 2-7GB)
	downloaded, err := downloader.Download(ctx, url, modelDest, progressFunc, downloader.SizeIntervalMIB100)
	if err != nil {
		os.Remove(modelDest) // Cleanup on failure
		return "", fmt.Errorf("download failed: %w", err)
	}

	// Check if download actually transferred data
	if !downloaded {
		os.Remove(modelDest) // Cleanup empty file
		return "", fmt.Errorf("download completed but no data was transferred")
	}

	return modelDest, nil
}

// DownloadModelSimple is deprecated. Use DownloadModel instead.
// Kept for backward compatibility.
func (m *ModelDownloader) DownloadModelSimple(ctx context.Context, url string) (string, error) {
	return m.DownloadModel(ctx, url, nil)
}

// GetModelsPath returns the models directory path
func (m *ModelDownloader) GetModelsPath() string {
	return m.modelsPath
}
