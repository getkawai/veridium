package image

import (
	"context"

	"github.com/kawai-network/veridium/pkg/stablediffusion/remote"
)

// RemoteGenerator handles remote API-based image generation
// DEPRECATED: Use pkg/stablediffusion/remote.Generator instead
type RemoteGenerator struct {
	generator *remote.UnifiedGenerator
}

// NewRemoteGenerator creates a new remote image generator
// DEPRECATED: Use remote.NewGenerator() instead
func NewRemoteGenerator() *RemoteGenerator {
	return &RemoteGenerator{
		generator: remote.NewGenerator(),
	}
}

// Generate generates an image using remote APIs (Gemini, Pollinations, etc.)
// DEPRECATED: Use generator.Generate() directly
func (rg *RemoteGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	// Convert internal GenerationOptions to remote.GenerationOptions
	remoteOpts := remote.GenerationOptions{
		Prompt:         opts.Prompt,
		NegativePrompt: opts.NegativePrompt,
		Model:          opts.Model,
		OutputPath:     opts.OutputPath,
		ImageUrl:       opts.ImageUrl,
		ImageUrls:      opts.ImageUrls,
		Width:          opts.Width,
		Height:         opts.Height,
		Size:           opts.Size,
		AspectRatio:    opts.AspectRatio,
		Steps:          opts.Steps,
		Cfg:            opts.Cfg,
		Strength:       opts.Strength,
		Seed:           opts.Seed,
		Quality:        opts.Quality,
		SamplerName:    opts.SamplerName,
		Scheduler:      opts.Scheduler,
		OutputFormat:   opts.OutputFormat,
	}

	return rg.generator.Generate(ctx, remoteOpts)
}

// GetAvailableModels returns list of available remote models
// DEPRECATED: Use generator.GetAvailableModels() directly
func (rg *RemoteGenerator) GetAvailableModels() []string {
	return rg.generator.GetAvailableModels()
}

// IsAvailable checks if remote generation is available (has API keys)
// DEPRECATED: Use generator.IsAvailable() directly
func (rg *RemoteGenerator) IsAvailable() bool {
	return rg.generator.IsAvailable()
}
