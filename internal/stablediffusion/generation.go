package stablediffusion

import (
	"fmt"
	"os"
	"strconv"
)

// GenerationOptions defines the options for image generation
type GenerationOptions struct {
	Prompt         string
	NegativePrompt string
	ModelPath      string
	OutputPath     string
	Width          int
	Height         int
	Steps          int
	Seed           int64
	OutputFormat   string // "png", "jpg", etc. (default to png if empty)
}

// GenerateImage executes the Stable Diffusion binary to generate an image
func (sdrm *StableDiffusionReleaseManager) GenerateImage(opts GenerationOptions) error {
	binaryPath := sdrm.GetBinaryPath()

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

	// Prepare arguments
	args := []string{
		"-m", opts.ModelPath,
		"-p", opts.Prompt,
		"-o", opts.OutputPath,
		"--width", strconv.Itoa(opts.Width),
		"--height", strconv.Itoa(opts.Height),
		"--steps", strconv.Itoa(opts.Steps),
		"--seed", strconv.FormatInt(opts.Seed, 10),
	}

	if opts.NegativePrompt != "" {
		args = append(args, "-n", opts.NegativePrompt)
	}

	// Add format if specific args are needed, though SD.cpp usually infers from extension or defaults
	// Current sd.cpp might not strictly require format arg if extension is present,
	// but let's keep it simple based on the previous implementation.

	// Helper for command execution to allow mocking in tests if we were injecting an executor
	// For now, we use exec.Command directly.

	// Execute command via the injected executor
	if err := sdrm.Executor.Run(binaryPath, args...); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Verify output exists
	if _, err := os.Stat(opts.OutputPath); err != nil {
		return fmt.Errorf("output file was not created: %w", err)
	}

	return nil
}
