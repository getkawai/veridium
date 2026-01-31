package image

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// LocalGenerator handles local Stable Diffusion binary execution
type LocalGenerator struct {
	engine *StableDiffusion
}

// NewLocalGenerator creates a new local image generator
func NewLocalGenerator(engine *StableDiffusion) *LocalGenerator {
	return &LocalGenerator{
		engine: engine,
	}
}

// Generate generates an image using local Stable Diffusion binary
func (lg *LocalGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	binaryPath := lg.engine.getBinaryPath()

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("stable diffusion binary not found at %s", binaryPath)
	}

	// Apply defaults
	opts = lg.applyDefaults(opts)

	// Prepare arguments
	args := lg.buildArgs(opts)

	// Log the command for debugging
	log.Printf("[LocalSD] Executing: %s %v", binaryPath, args)

	// Execute command via the injected executor
	if err := lg.engine.Executor.Run(ctx, binaryPath, args...); err != nil {
		return fmt.Errorf("local generation failed: %w", err)
	}

	// Verify output exists
	if _, err := os.Stat(opts.OutputPath); err != nil {
		return fmt.Errorf("output file was not created: %w", err)
	}

	log.Printf("[LocalSD] Image generated successfully: %s", opts.OutputPath)
	return nil
}

// applyDefaults applies default values to generation options
func (lg *LocalGenerator) applyDefaults(opts GenerationOptions) GenerationOptions {
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
	return opts
}

// buildArgs constructs command-line arguments for SD binary
func (lg *LocalGenerator) buildArgs(opts GenerationOptions) []string {
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

	// Add negative prompt
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

	return args
}

// IsAvailable checks if local SD binary is available
func (lg *LocalGenerator) IsAvailable() bool {
	binaryPath := lg.engine.getBinaryPath()
	_, err := os.Stat(binaryPath)
	return err == nil
}

// GetFirstAvailableModel returns the first available SD model
func (lg *LocalGenerator) GetFirstAvailableModel() string {
	modelsPath := lg.engine.GetModelsPath()
	files, err := os.ReadDir(modelsPath)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			// Check for all supported model formats including GGUF
			if lg.isSupportedModelFormat(name) {
				return filepath.Join(modelsPath, name)
			}
		}
	}
	return ""
}

// isSupportedModelFormat checks if file extension is a supported model format
func (lg *LocalGenerator) isSupportedModelFormat(filename string) bool {
	supportedExts := []string{".ckpt", ".safetensors", ".pt", ".bin", ".gguf"}
	for _, ext := range supportedExts {
		if len(filename) >= len(ext) && filename[len(filename)-len(ext):] == ext {
			return true
		}
	}
	return false
}
