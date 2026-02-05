package remote

import (
	"context"
	"fmt"
	"log"
)

// UnifiedGenerator provides a unified interface for all remote generators
type UnifiedGenerator struct {
	gemini       *GeminiGenerator
	pollinations *PollinationsGenerator
}

// NewGenerator creates a new unified remote generator
func NewGenerator() *UnifiedGenerator {
	return &UnifiedGenerator{
		gemini:       NewGeminiGenerator(),
		pollinations: NewPollinationsGenerator(),
	}
}

// Generate generates an image using the best available remote API
func (u *UnifiedGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	// Determine which generator to use based on model
	if opts.Model == "" || opts.Model == "gemini-2.5-flash" || opts.Model == "gemini-2.5-flash-image" {
		// Try Gemini first (if available)
		if u.gemini.IsAvailable() {
			log.Printf("[RemoteGen] Using Gemini API")
			return u.gemini.Generate(ctx, opts)
		}
		// Fallback to Pollinations
		log.Printf("[RemoteGen] Gemini unavailable, falling back to Pollinations")
		return u.pollinations.Generate(ctx, opts)
	}

	// For other models, use Pollinations
	log.Printf("[RemoteGen] Using Pollinations API for model: %s", opts.Model)
	return u.pollinations.Generate(ctx, opts)
}

// GenerateWithFallback generates an image with automatic fallback
func (u *UnifiedGenerator) GenerateWithFallback(ctx context.Context, opts GenerationOptions) error {
	// Try Gemini first
	if u.gemini.IsAvailable() {
		err := u.gemini.Generate(ctx, opts)
		if err == nil {
			return nil
		}
		log.Printf("[RemoteGen] Gemini failed: %v, trying Pollinations", err)
	}

	// Fallback to Pollinations
	return u.pollinations.Generate(ctx, opts)
}

// GetAvailableModels returns all available models from all generators
func (u *UnifiedGenerator) GetAvailableModels() []string {
	models := make([]string, 0)

	if u.gemini.IsAvailable() {
		models = append(models, u.gemini.GetAvailableModels()...)
	}

	models = append(models, u.pollinations.GetAvailableModels()...)

	return models
}

// IsAvailable checks if any remote generator is available
func (u *UnifiedGenerator) IsAvailable() bool {
	return u.gemini.IsAvailable() || u.pollinations.IsAvailable()
}

// GetGenerator returns a specific generator by name
func (u *UnifiedGenerator) GetGenerator(name string) (Generator, error) {
	switch name {
	case "gemini":
		if !u.gemini.IsAvailable() {
			return nil, fmt.Errorf("Gemini generator not available")
		}
		return u.gemini, nil
	case "pollinations":
		return u.pollinations, nil
	default:
		return nil, fmt.Errorf("unknown generator: %s", name)
	}
}
