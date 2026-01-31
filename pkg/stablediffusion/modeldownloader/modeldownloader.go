// Package modeldownloader provides functionality to download Stable Diffusion models.
// This package handles model discovery, download with progress tracking, and validation.
package modeldownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/grab"
)

const (
	// DefaultModelURL is the default Stable Diffusion model URL (SD 1.4 ~4GB)
	DefaultModelURL = "https://huggingface.co/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt"
	// DefaultModelFilename is the default filename for the downloaded model
	DefaultModelFilename = "sd-v1-4.ckpt"
	// DownloadTimeout is the timeout for model download (~4GB file)
	DownloadTimeout = 30 * time.Minute
)

// ProgressCallback is called during download to report progress
type ProgressCallback func(bytesComplete, totalBytes int64, mbps float64, done bool)

// ModelDownloader handles model discovery and download
type ModelDownloader struct {
	modelsPath string
	client     *http.Client
}

// New creates a new ModelDownloader instance
func New(modelsPath string) *ModelDownloader {
	return &ModelDownloader{
		modelsPath: modelsPath,
		client: &http.Client{
			Timeout: DownloadTimeout,
		},
	}
}

// NewWithClient creates a new ModelDownloader with a custom HTTP client
func NewWithClient(modelsPath string, client *http.Client) *ModelDownloader {
	return &ModelDownloader{
		modelsPath: modelsPath,
		client:     client,
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
// Uses grab package for resume support and progress tracking
func (m *ModelDownloader) DownloadModel(ctx context.Context, url string, progress ProgressCallback) (string, error) {
	if err := os.MkdirAll(m.modelsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create models directory: %w", err)
	}

	modelDest := filepath.Join(m.modelsPath, DefaultModelFilename)

	// Use grab for download with resume support
	req, err := grab.NewRequest(modelDest, url)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	req = req.WithContext(ctx)

	client := grab.NewClient()
	resp := client.Do(req)

	// Monitor progress
	if progress != nil {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ticker.C:
					if resp.IsComplete() {
						return
					}
					progress(resp.BytesComplete(), resp.Size(), resp.BytesPerSecond()/(1024*1024), false)
				case <-resp.Done:
					return
				}
			}
		}()
	}

	// Wait for download to complete
	if err := resp.Err(); err != nil {
		os.Remove(modelDest) // Cleanup partial file
		return "", fmt.Errorf("download failed: %w", err)
	}

	if progress != nil {
		progress(resp.BytesComplete(), resp.Size(), resp.BytesPerSecond()/(1024*1024), true)
	}

	return modelDest, nil
}

// DownloadModelSimple downloads a model using standard http.Client
// Use this when grab is not needed (simpler cases without resume)
func (m *ModelDownloader) DownloadModelSimple(ctx context.Context, url string) (string, error) {
	if err := os.MkdirAll(m.modelsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create models directory: %w", err)
	}

	modelDest := filepath.Join(m.modelsPath, DefaultModelFilename)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %s", resp.Status)
	}

	out, err := os.Create(modelDest)
	if err != nil {
		return "", fmt.Errorf("failed to create model file: %w", err)
	}

	_, err = io.Copy(out, resp.Body)
	if closeErr := out.Close(); closeErr != nil {
		os.Remove(modelDest) // Cleanup partial file
		return "", fmt.Errorf("failed to close model file: %w", closeErr)
	}
	if err != nil {
		os.Remove(modelDest) // Cleanup partial file
		return "", fmt.Errorf("failed to save model file: %w", err)
	}

	return modelDest, nil
}

// GetModelsPath returns the models directory path
func (m *ModelDownloader) GetModelsPath() string {
	return m.modelsPath
}
