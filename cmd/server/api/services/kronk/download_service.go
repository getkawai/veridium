// Package kronk provides setup and download functionality for the DeAI server.
package kronk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kawai-network/veridium/cmd/server/app/domain/ttsapp"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/pkg/stablediffusion/modeldownloader"
	sdmodels "github.com/kawai-network/veridium/pkg/stablediffusion/models"
	"github.com/kawai-network/veridium/pkg/tools/downloader"
	"github.com/kawai-network/veridium/pkg/tools/libs"
)

// Download progress callback type
type ProgressCallback func(completed, total int64, percent float64, mbps float64)

// DownloadResult represents the result of a download operation
type DownloadResult struct {
	Success  bool
	FilePath string
	Bytes    int64
	Duration time.Duration
	Error    error
	Resumed  bool // True if download was resumed from partial file
}

// DownloadService handles all download operations with resume and retry support
type DownloadService struct {
	basePath     string
	modelsPath   string
	maxRetries   int
	retryDelay   time.Duration
	timeout      time.Duration
	progressCb   ProgressCallback
	mu           sync.Mutex
	lastProgress time.Time
	lastBytes    int64
}

// DownloadServiceOption configures the DownloadService
type DownloadServiceOption func(*DownloadService)

// WithBasePath sets the base path for downloads
func WithBasePath(path string) DownloadServiceOption {
	return func(s *DownloadService) {
		s.basePath = path
	}
}

// WithModelsPath sets the models path
func WithModelsPath(path string) DownloadServiceOption {
	return func(s *DownloadService) {
		s.modelsPath = path
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(count int) DownloadServiceOption {
	return func(s *DownloadService) {
		s.maxRetries = count
	}
}

// WithRetryDelay sets the delay between retries
func WithRetryDelay(delay time.Duration) DownloadServiceOption {
	return func(s *DownloadService) {
		s.retryDelay = delay
	}
}

// WithTimeout sets the download timeout
func WithTimeout(timeout time.Duration) DownloadServiceOption {
	return func(s *DownloadService) {
		s.timeout = timeout
	}
}

// WithProgressCallback sets the progress callback function
func WithProgressCallback(cb ProgressCallback) DownloadServiceOption {
	return func(s *DownloadService) {
		s.progressCb = cb
	}
}

// NewDownloadService creates a new DownloadService instance
func NewDownloadService(opts ...DownloadServiceOption) *DownloadService {
	s := &DownloadService{
		basePath:     paths.Base(),
		modelsPath:   paths.Models(),
		maxRetries:   3,
		retryDelay:   5 * time.Second,
		timeout:      30 * time.Minute,
		lastProgress: time.Now(),
		lastBytes:    0,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SetProgressCallback updates the progress callback function
func (s *DownloadService) SetProgressCallback(cb ProgressCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.progressCb = cb
}

// DownloadWithRetry downloads a file with retry support and resume capability
func (s *DownloadService) DownloadWithRetry(ctx context.Context, url, dest string) *DownloadResult {
	var result *DownloadResult
	var lastErr error

	startTime := time.Now()

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-time.After(s.retryDelay):
			case <-ctx.Done():
				return &DownloadResult{
					Success:  false,
					Error:    ctx.Err(),
					Duration: time.Since(startTime),
				}
			}
		}

		// Create context with timeout for this attempt
		downloadCtx, cancel := context.WithTimeout(ctx, s.timeout)

		result = s.downloadInternal(downloadCtx, url, dest, startTime)
		cancel()

		if result.Success {
			return result
		}

		lastErr = result.Error

		// Don't retry on context cancellation
		if ctx.Err() != nil {
			return result
		}
	}

	return &DownloadResult{
		Success:  false,
		Error:    fmt.Errorf("download failed after %d attempts: %w", s.maxRetries+1, lastErr),
		Duration: time.Since(startTime),
	}
}

// downloadInternal performs a single download attempt
func (s *DownloadService) downloadInternal(ctx context.Context, url, dest string, startTime time.Time) *DownloadResult {
	// Create progress wrapper
	var progressCb downloader.ProgressFunc
	if s.progressCb != nil {
		// Signal a fresh transfer so UI progress bars can reset between files/attempts.
		s.progressCb(0, 0, 0, 0)

		cb := s.progressCb
		progressCb = func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
			percent := 0.0
			if totalSize > 0 {
				percent = float64(currentSize) / float64(totalSize) * 100
			}
			cb(currentSize, totalSize, percent, mibPerSec)
		}
	}

	// Download with resume support (via grab package)
	downloaded, err := downloader.Download(ctx, url, dest, progressCb, downloader.SizeIntervalMIB)

	result := &DownloadResult{
		FilePath: dest,
		Duration: time.Since(startTime),
	}

	if err != nil {
		result.Error = err
		result.Success = false
		return result
	}

	if !downloaded {
		// grab may report zero transferred bytes when the destination file is already
		// fully present (e.g., resumed/validated from previous run). Treat this as success
		// if the file exists and has content.
		if info, statErr := os.Stat(dest); statErr == nil && info.Size() > 0 {
			result.Bytes = info.Size()
			result.Success = true
			result.Resumed = true
			return result
		}

		result.Error = fmt.Errorf("download completed but no data transferred")
		result.Success = false
		return result
	}

	// Get file size
	info, err := os.Stat(dest)
	if err == nil {
		result.Bytes = info.Size()
	}

	result.Success = true
	return result
}

// DownloadLibrary downloads a library file
func (s *DownloadService) DownloadLibrary(ctx context.Context, url, dest string) *DownloadResult {
	return s.DownloadWithRetry(ctx, url, dest)
}

// DownloadWhisperModel downloads the Whisper speech-to-text model
// Uses grab package for automatic resume support
func (s *DownloadService) DownloadWhisperModel(ctx context.Context, modelName string) *DownloadResult {
	startTime := time.Now()

	// Whisper model URL and path
	modelURL := fmt.Sprintf("https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-%s.bin", modelName)
	modelPath := filepath.Join(paths.Models(), "ggerganov", "whisper.cpp", fmt.Sprintf("ggml-%s.bin", modelName))

	// Create directory
	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return &DownloadResult{
			Success:  false,
			Error:    fmt.Errorf("failed to create model directory: %w", err),
			Duration: time.Since(startTime),
		}
	}

	result := s.DownloadWithRetry(ctx, modelURL, modelPath)
	result.Duration = time.Since(startTime)
	return result
}

// DownloadStableDiffusionModel downloads the Stable Diffusion image generation model
// Uses grab package for automatic resume support
func (s *DownloadService) DownloadStableDiffusionModel(ctx context.Context) *DownloadResult {
	startTime := time.Now()

	modelURL := modeldownloader.DefaultModelURL
	modelPath := filepath.Join(paths.Models(), "CompVis", "stable-diffusion-v-1-4-original", "sd-v1-4.ckpt")

	// Create directory
	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return &DownloadResult{
			Success:  false,
			Error:    fmt.Errorf("failed to create model directory: %w", err),
			Duration: time.Since(startTime),
		}
	}

	result := s.DownloadWithRetry(ctx, modelURL, modelPath)
	result.Duration = time.Since(startTime)
	return result
}

// DownloadStableDiffusionModelSmart downloads a hardware-selected Stable Diffusion model.
// Uses HuggingFace URL to determine path: models/{author}/{repo}/{filename}
func (s *DownloadService) DownloadStableDiffusionModelSmart(ctx context.Context, spec sdmodels.ModelSpec) *DownloadResult {
	startTime := time.Now()

	if spec.URL == "" || spec.Filename == "" {
		return &DownloadResult{
			Success:  false,
			Error:    fmt.Errorf("invalid model spec: missing URL or filename"),
			Duration: time.Since(startTime),
		}
	}

	totalBytes := int64(0)

	downloadOne := func(component, url, filename string) *DownloadResult {
		modelDir, err := paths.ModelPath(url)
		if err != nil {
			return &DownloadResult{
				Success:  false,
				Error:    fmt.Errorf("failed to parse %s model URL (%s): %w", component, url, err),
				Duration: time.Since(startTime),
			}
		}

		modelPath := filepath.Join(modelDir, filename)
		if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
			return &DownloadResult{
				Success:  false,
				Error:    fmt.Errorf("failed to create model directory: %w", err),
				Duration: time.Since(startTime),
			}
		}

		result := s.DownloadWithRetry(ctx, url, modelPath)
		if result.Success {
			totalBytes += result.Bytes
		} else if result.Error != nil {
			result.Error = fmt.Errorf("%s model download failed (%s): %w", component, url, result.Error)
		}
		return result
	}

	primary := downloadOne("diffusion", spec.URL, spec.Filename)
	if !primary.Success {
		primary.Duration = time.Since(startTime)
		return primary
	}

	if spec.LLMURL != "" && spec.LLMFilename != "" {
		fmt.Printf("\n→ Downloading SD component: llm (%s)\n", spec.LLMFilename)
		llm := downloadOne("llm", spec.LLMURL, spec.LLMFilename)
		if !llm.Success {
			llm.Duration = time.Since(startTime)
			return llm
		}
	}

	if spec.VAEURL != "" && spec.VAEFilename != "" {
		fmt.Printf("\n→ Downloading SD component: vae (%s)\n", spec.VAEFilename)
		vae := downloadOne("vae", spec.VAEURL, spec.VAEFilename)
		if !vae.Success {
			vae.Duration = time.Since(startTime)
			return vae
		}
	}

	if spec.EditModelURL != "" && spec.EditModelFile != "" {
		fmt.Printf("\n→ Downloading SD component: edit (%s)\n", spec.EditModelFile)
		edit := downloadOne("edit", spec.EditModelURL, spec.EditModelFile)
		if !edit.Success {
			if spec.EditFallbackURL != "" && spec.EditFallbackFile != "" {
				fmt.Printf("\n→ Primary edit model failed, trying fallback: %s\n", spec.EditFallbackFile)
				fallback := downloadOne("edit_fallback", spec.EditFallbackURL, spec.EditFallbackFile)
				if !fallback.Success {
					fallback.Duration = time.Since(startTime)
					fallback.Error = fmt.Errorf("failed to download edit model (primary and fallback): primary=%v fallback=%w", edit.Error, fallback.Error)
					return fallback
				}
			} else {
				edit.Duration = time.Since(startTime)
				return edit
			}
		}
	}

	return &DownloadResult{
		Success:  true,
		FilePath: primary.FilePath,
		Bytes:    totalBytes,
		Duration: time.Since(startTime),
	}
}

// DownloadTTSModel downloads the Text-to-Speech model
// Uses grab package for automatic resume support
func (s *DownloadService) DownloadTTSModel(ctx context.Context) *DownloadResult {
	startTime := time.Now()

	modelURL := ttsapp.DefaultTTSModelURL
	modelPath := filepath.Join(paths.Models(), "mmwillet2", "Kokoro_GGUF", "Kokoro_no_espeak_Q4.gguf")

	// Create directory
	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return &DownloadResult{
			Success:  false,
			Error:    fmt.Errorf("failed to create model directory: %w", err),
			Duration: time.Since(startTime),
		}
	}

	result := s.DownloadWithRetry(ctx, modelURL, modelPath)
	result.Duration = time.Since(startTime)
	return result
}

// DownloadLLMModel downloads the Large Language Model
// Uses grab package for automatic resume support
func (s *DownloadService) DownloadLLMModel(ctx context.Context, org, repo, filename string) *DownloadResult {
	startTime := time.Now()

	modelPath := filepath.Join(s.modelsPath, org, repo, filename)
	modelURL := fmt.Sprintf("https://huggingface.co/%s/%s/resolve/main/%s", org, repo, filename)

	// Create directory
	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return &DownloadResult{
			Success:  false,
			Error:    fmt.Errorf("failed to create model directory: %w", err),
			Duration: time.Since(startTime),
		}
	}

	result := s.DownloadWithRetry(ctx, modelURL, modelPath)
	result.Duration = time.Since(startTime)
	return result
}

// DownloadLLMModelDefault downloads the default LLM model (Nemotron 3 Nano)
func (s *DownloadService) DownloadLLMModelDefault(ctx context.Context) *DownloadResult {
	return s.DownloadLLMModel(ctx, DefaultLLMOrg, DefaultLLMRepo, DefaultLLMFile)
}

// DownloadLibraries downloads all required libraries (llama.cpp, whisper.cpp, SD, TTS)
func (s *DownloadService) DownloadLibraries(ctx context.Context) ([]DownloadResult, error) {
	results := make([]DownloadResult, 0)

	// List of libraries to download
	libraries := []libs.LibraryType{
		libs.LibraryLlama,
		libs.LibraryWhisper,
		libs.LibraryStableDiffusion,
		libs.LibraryTTS,
	}

	for _, libType := range libraries {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		// Create lib manager for specific library
		libMgr, err := libs.New(
			libs.WithBasePath(paths.Base()),
			libs.WithAllowUpgrade(true), // Setup should install missing libraries
			libs.WithLibraryType(libType),
		)
		if err != nil {
			results = append(results, DownloadResult{
				Success: false,
				Error:   fmt.Errorf("failed to init %s: %w", libType, err),
			})
			continue
		}

		// Download with progress
		downloadCtx, cancel := context.WithTimeout(ctx, LibraryDownloadTimeout)

		var libResult DownloadResult
		startTime := time.Now()

		_, err = libMgr.DownloadWithProgress(downloadCtx, func(ctx context.Context, msg string, args ...any) {
			// Log callback
		}, func(bytesComplete, totalBytes int64, mbps float64, done bool) {
			if s.progressCb != nil && totalBytes > 0 {
				percent := float64(bytesComplete) / float64(totalBytes) * 100
				s.progressCb(bytesComplete, totalBytes, percent, mbps)
			}
		})

		cancel()

		libResult.Duration = time.Since(startTime)

		if err != nil {
			libResult.Success = false
			libResult.Error = fmt.Errorf("failed to download %s: %w", libType, err)
		} else {
			libResult.Success = true
		}

		results = append(results, libResult)
	}

	return results, nil
}
