package stablediffusion

import (
	"fmt"
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
}

// CreateImage handles frontend CreateImageRequest and generates images
func (sdrm *StableDiffusion) CreateImage(req CreateImageRequest) error {
	// Get first available model if not specified
	modelPath := sdrm.GetFirstAvailableModel()
	if modelPath == "" {
		return fmt.Errorf("no SD model found")
	}

	// Generate output path
	outputDir := sdrm.GetOutputsPath()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Convert RuntimeImageGenParams to GenerationOptions
	opts := GenerationOptions{
		Prompt:      req.Params.Prompt,
		ModelPath:   modelPath,
		ImageUrl:    req.Params.ImageUrl,
		ImageUrls:   req.Params.ImageUrls,
		Size:        req.Params.Size,
		AspectRatio: req.Params.AspectRatio,
		Quality:     req.Params.Quality,
		SamplerName: req.Params.SamplerName,
		Scheduler:   req.Params.Scheduler,
		Seed:        req.Params.Seed,
	}

	// Handle optional numeric params
	if req.Params.Width != nil {
		opts.Width = *req.Params.Width
	}
	if req.Params.Height != nil {
		opts.Height = *req.Params.Height
	}
	if req.Params.Steps != nil {
		opts.Steps = *req.Params.Steps
	}
	if req.Params.Cfg != nil {
		opts.Cfg = *req.Params.Cfg
	}
	if req.Params.Strength != nil {
		opts.Strength = *req.Params.Strength
	}

	// Generate multiple images based on imageNum
	for i := 0; i < req.ImageNum; i++ {
		outputPath := fmt.Sprintf("%s/%s_%d.png", outputDir, req.GenerationTopicId, i)
		opts.OutputPath = outputPath

		if err := sdrm.createImageInternal(opts); err != nil {
			return err
		}
	}

	return nil
}

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
func (sdrm *StableDiffusion) GetOutputsPath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + "/.stable-diffusion/outputs"
}

// CreateImageWithOptions generates an image using GenerationOptions directly
// Used by internal services like image_designer
func (sdrm *StableDiffusion) CreateImageWithOptions(opts GenerationOptions) error {
	return sdrm.createImageInternal(opts)
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
