package local

import "context"

// GenerationOptions defines options for local SD binary execution
type GenerationOptions struct {
	Prompt         string
	NegativePrompt string
	ModelPath      string  // Path to SD model file
	OutputPath     string  // Where to save generated image
	ImageUrl       *string // Input image for img2img
	Width          int
	Height         int
	Steps          int     // Sampling steps
	Cfg            float64 // CFG scale
	Strength       float64 // Strength for img2img
	Seed           *int64  // Random seed
	SamplerName    string  // Sampler method
	Scheduler      string  // Scheduler type
}

// CommandExecutor interface for executing SD binary commands
type CommandExecutor interface {
	Run(ctx context.Context, name string, args ...string) error
}

// Generator interface for local SD generation
type Generator interface {
	// Generate generates an image using local SD binary
	Generate(ctx context.Context, opts GenerationOptions) error

	// IsAvailable checks if local SD binary is available
	IsAvailable() bool

	// GetFirstAvailableModel returns the first available SD model
	GetFirstAvailableModel() string
}
