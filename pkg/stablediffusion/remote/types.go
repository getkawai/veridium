package remote

import "context"

// GenerationOptions defines options for remote image generation
type GenerationOptions struct {
	Prompt         string
	NegativePrompt string
	Model          string   // Model name for remote API
	OutputPath     string   // Where to save the generated image
	ImageUrl       *string  // Input image URL for img2img
	ImageUrls      []string // Multiple input images
	Width          int
	Height         int
	Size           string // Size preset (e.g., "1024x1024")
	AspectRatio    string // Aspect ratio (e.g., "16:9", "1:1")
	Steps          int
	Cfg            float64
	Strength       float64 // For img2img
	Seed           *int64
	Quality        string // Quality preset
	SamplerName    string // Sampler method
	Scheduler      string // Scheduler type
	OutputFormat   string // Output format (png, jpg, etc.)
}

// Generator interface for remote image generation
type Generator interface {
	// Generate generates an image using remote API
	Generate(ctx context.Context, opts GenerationOptions) error

	// GetAvailableModels returns list of available models
	GetAvailableModels() []string

	// IsAvailable checks if the generator is available (has API keys, etc.)
	IsAvailable() bool
}
