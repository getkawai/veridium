package image

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RuntimeImageGenParams matches frontend model-bank/standard-parameters RuntimeImageGenParams
type RuntimeImageGenParams struct {
	Prompt      string   `json:"prompt"`
	ImageUrl    *string  `json:"imageUrl,omitempty"`
	ImageUrls   []string `json:"imageUrls,omitempty"`
	Width       *int     `json:"width,omitempty"`
	Height      *int     `json:"height,omitempty"`
	Size        string   `json:"size,omitempty"`
	AspectRatio string   `json:"aspectRatio,omitempty"`
	Cfg         *float64 `json:"cfg,omitempty"`
	Strength    *float64 `json:"strength,omitempty"`
	Steps       *int     `json:"steps,omitempty"`
	Quality     string   `json:"quality,omitempty"`
	Seed        *int64   `json:"seed,omitempty"`
	SamplerName string   `json:"samplerName,omitempty"`
	Scheduler   string   `json:"scheduler,omitempty"`
}

// CreateImageRequest matches frontend createImage action parameters
type CreateImageRequest struct {
	GenerationTopicId string                `json:"generationTopicId"`
	Provider          string                `json:"provider"`
	Model             string                `json:"model"`
	ImageNum          int                   `json:"imageNum"`
	Params            RuntimeImageGenParams `json:"params"`
}

// GenerationOptions defines internal options for SD binary execution
type GenerationOptions struct {
	Prompt         string
	NegativePrompt string
	ModelPath      string
	OutputPath     string
	ImageUrl       *string
	ImageUrls      []string
	Width          int
	Height         int
	Size           string
	AspectRatio    string
	Steps          int
	Cfg            float64
	Strength       float64
	Seed           *int64
	Quality        string
	SamplerName    string
	Scheduler      string
	OutputFormat   string
	Model          string // Model name for remote API
}

// CreateImage handles frontend CreateImageRequest and generates images asynchronously
// MOVED TO internal/image/service.go

// GetFirstAvailableModel returns the first available SD model
func (sdrm *StableDiffusion) GetFirstAvailableModel() string {
	modelsPath := sdrm.GetModelsPath()
	files, err := os.ReadDir(modelsPath)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			// Check for all supported model formats including GGUF
			if strings.HasSuffix(name, ".ckpt") ||
				strings.HasSuffix(name, ".safetensors") ||
				strings.HasSuffix(name, ".pt") ||
				strings.HasSuffix(name, ".bin") ||
				strings.HasSuffix(name, ".gguf") {
				return filepath.Join(modelsPath, name)
			}
		}
	}
	return ""
}

// GetOutputsPath returns the path for generated images
// Note: This matches the default but we override it in CreateImage to use files/uploads
func (sdrm *StableDiffusion) GetOutputsPath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + "/.stable-diffusion/outputs"
}

// CreateImageWithOptions generates an image using GenerationOptions directly
// Used by internal services like image_designer
func (sdrm *StableDiffusion) CreateImageWithOptions(opts GenerationOptions) error {
	return sdrm.createImageInternal(opts)
}

// generateImageRemote generates an image using Pollinations AI remote API
func (sdrm *StableDiffusion) generateImageRemote(opts GenerationOptions) error {
	// Build the API URL
	baseURL := "https://image.pollinations.ai/prompt/"
	encodedPrompt := url.PathEscape(opts.Prompt)
	apiURL := baseURL + encodedPrompt

	// Build query parameters
	params := url.Values{}
	if opts.Width > 0 {
		params.Add("width", strconv.Itoa(opts.Width))
	}
	if opts.Height > 0 {
		params.Add("height", strconv.Itoa(opts.Height))
	}
	// Use specified model or default to flux
	model := opts.Model
	if model == "" {
		model = "flux"
	}
	params.Add("model", model)
	params.Add("enhance", "false")
	params.Add("nologo", "true")

	if len(params) > 0 {
		apiURL += "$" + params.Encode()
	}

	log.Printf("[Remote] Requesting image from: %s", apiURL)

	// Make HTTP request
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("remote API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote API returned status %d", resp.StatusCode)
	}

	// Create output file
	outFile, err := os.Create(opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Download image
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}

	log.Printf("[Remote] Image downloaded successfully to: %s", opts.OutputPath)
	return nil
}

// createImageInternal executes the Stable Diffusion binary to generate an image
func (sdrm *StableDiffusion) createImageInternal(opts GenerationOptions) error {
	binaryPath := sdrm.getBinaryPath()

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("stable diffusion binary not found at %s", binaryPath)
	}

	// Default values
	if opts.Width == 0 {
		opts.Width = 1024
	}
	if opts.Height == 0 {
		opts.Height = 1024
	}
	if opts.Steps == 0 {
		opts.Steps = 20
	}
	if opts.Cfg == 0 {
		opts.Cfg = 7.0
	}

	// Prepare arguments
	args := []string{
		"-m", opts.ModelPath,
		"-p", opts.Prompt,
		"-o", opts.OutputPath,
		"--width", strconv.Itoa(opts.Width),
		"--height", strconv.Itoa(opts.Height),
		"--steps", strconv.Itoa(opts.Steps),
		"--cfg-scale", strconv.FormatFloat(opts.Cfg, 'f', -1, 64),
	}

	// Add seed if specified
	if opts.Seed != nil {
		args = append(args, "--seed", strconv.FormatInt(*opts.Seed, 10))
	}

	if opts.NegativePrompt != "" {
		args = append(args, "-n", opts.NegativePrompt)
	}

	// Add sampler if specified
	if opts.SamplerName != "" {
		args = append(args, "--sampling-method", opts.SamplerName)
	}

	// Add scheduler if specified
	if opts.Scheduler != "" {
		args = append(args, "--schedule", opts.Scheduler)
	}

	// Add strength for img2img
	if opts.Strength > 0 {
		args = append(args, "--strength", strconv.FormatFloat(opts.Strength, 'f', -1, 64))
	}

	// Add input image for img2img
	if opts.ImageUrl != nil && *opts.ImageUrl != "" {
		args = append(args, "-i", *opts.ImageUrl)
	}

	// Log the command for debugging
	log.Printf("[SD] Executing: %s %v", binaryPath, args)

	// Execute command via the injected executor
	if err := sdrm.Executor.Run(sdrm.ctx, binaryPath, args...); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Verify output exists
	if _, err := os.Stat(opts.OutputPath); err != nil {
		return fmt.Errorf("output file was not created: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
