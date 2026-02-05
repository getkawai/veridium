package image

import (
	"context"

	"github.com/kawai-network/veridium/pkg/stablediffusion/local"
)

// LocalGenerator handles local Stable Diffusion binary execution
// DEPRECATED: Use pkg/stablediffusion/local.Generator instead
type LocalGenerator struct {
	engine    *StableDiffusion
	generator *local.SDGenerator
}

// NewLocalGenerator creates a new local image generator
// DEPRECATED: Use local.NewGenerator() instead
func NewLocalGenerator(engine *StableDiffusion) *LocalGenerator {
	binaryPath := engine.getBinaryPath()
	modelsPath := engine.GetModelsPath()

	// Use NewGeneratorWithExecutor to preserve process tracking for cleanup
	// This ensures engine.Cleanup() can kill stranded processes
	return &LocalGenerator{
		engine:    engine,
		generator: local.NewGeneratorWithExecutor(binaryPath, modelsPath, engine.Executor),
	}
}

// Generate generates an image using local Stable Diffusion binary
// DEPRECATED: Use generator.Generate() directly
func (lg *LocalGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	// Convert internal GenerationOptions to local.GenerationOptions
	localOpts := local.GenerationOptions{
		Prompt:         opts.Prompt,
		NegativePrompt: opts.NegativePrompt,
		ModelPath:      opts.ModelPath,
		OutputPath:     opts.OutputPath,
		ImageUrl:       opts.ImageUrl,
		Width:          opts.Width,
		Height:         opts.Height,
		Steps:          opts.Steps,
		Cfg:            opts.Cfg,
		Strength:       opts.Strength,
		Seed:           opts.Seed,
		SamplerName:    opts.SamplerName,
		Scheduler:      opts.Scheduler,
	}

	return lg.generator.Generate(ctx, localOpts)
}

// IsAvailable checks if local SD binary is available
// DEPRECATED: Use generator.IsAvailable() directly
func (lg *LocalGenerator) IsAvailable() bool {
	return lg.generator.IsAvailable()
}

// GetFirstAvailableModel returns the first available SD model
// DEPRECATED: Use generator.GetFirstAvailableModel() directly
func (lg *LocalGenerator) GetFirstAvailableModel() string {
	return lg.generator.GetFirstAvailableModel()
}
