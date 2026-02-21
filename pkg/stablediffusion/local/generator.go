package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// SDGenerator handles local Stable Diffusion binary execution
type SDGenerator struct {
	binaryPath string
	modelsPath string
	executor   CommandExecutor
}

// NewGenerator creates a new local SD generator
func NewGenerator(binaryPath, modelsPath string) *SDGenerator {
	binDir := filepath.Dir(binaryPath)
	return &SDGenerator{
		binaryPath: binaryPath,
		modelsPath: modelsPath,
		executor:   NewDefaultExecutor(binDir),
	}
}

// NewGeneratorWithExecutor creates a new local SD generator with custom executor
func NewGeneratorWithExecutor(binaryPath, modelsPath string, executor CommandExecutor) *SDGenerator {
	return &SDGenerator{
		binaryPath: binaryPath,
		modelsPath: modelsPath,
		executor:   executor,
	}
}

// Generate generates an image using local Stable Diffusion binary
func (g *SDGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	// Check if binary exists
	if _, err := os.Stat(g.binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("stable diffusion binary not found at %s", g.binaryPath)
	}

	// Apply defaults
	opts = g.applyDefaults(opts)

	// Prepare arguments
	args := g.buildArgs(opts)

	// Execute command (logging removed - library code should not log)
	if err := g.executor.Run(ctx, g.binaryPath, args...); err != nil {
		return fmt.Errorf("local generation failed: %w", err)
	}

	// Verify output exists
	if _, err := os.Stat(opts.OutputPath); err != nil {
		return fmt.Errorf("output file was not created: %w", err)
	}

	return nil
}

// IsAvailable checks if local SD binary is available
func (g *SDGenerator) IsAvailable() bool {
	_, err := os.Stat(g.binaryPath)
	return err == nil
}

// GetFirstAvailableModel returns the first available SD model
// Scans models directory recursively to find supported model formats
func (g *SDGenerator) GetFirstAvailableModel() string {
	return g.findFirstModel(g.modelsPath)
}

// findFirstModel recursively searches for the first supported model file
func (g *SDGenerator) findFirstModel(dir string) string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if model := g.findFirstModel(fullPath); model != "" {
				return model
			}
		} else if g.isSupportedModelFormat(file.Name()) {
			return fullPath
		}
	}
	return ""
}

// GetBinaryPath returns the path to the SD binary
func (g *SDGenerator) GetBinaryPath() string {
	return g.binaryPath
}

// GetModelsPath returns the path to the models directory
func (g *SDGenerator) GetModelsPath() string {
	return g.modelsPath
}

// applyDefaults applies default values to generation options
func (g *SDGenerator) applyDefaults(opts GenerationOptions) GenerationOptions {
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
func (g *SDGenerator) buildArgs(opts GenerationOptions) []string {
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

// isSupportedModelFormat checks if file extension is a supported model format
// Case-insensitive check to support uppercase extensions on Windows/macOS
func (g *SDGenerator) isSupportedModelFormat(filename string) bool {
	supportedExts := []string{".ckpt", ".safetensors", ".pt", ".bin", ".gguf"}
	filenameLower := strings.ToLower(filename)
	for _, ext := range supportedExts {
		if strings.HasSuffix(filenameLower, ext) {
			return true
		}
	}
	return false
}

// GetBinaryName returns the appropriate binary name for the current platform
func GetBinaryName() string {
	if runtime.GOOS == "windows" {
		return "sd.exe"
	}
	return "sd"
}
