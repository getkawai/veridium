package remote

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// UnifiedGenerator provides a unified interface for all remote generators
type UnifiedGenerator struct {
	gemini     *GeminiGenerator
	cloudflare *CloudflareGenerator
}

// NewGenerator creates a new unified remote generator
func NewGenerator() *UnifiedGenerator {
	return &UnifiedGenerator{
		gemini:     NewGeminiGenerator(),
		cloudflare: NewCloudflareGenerator(),
	}
}

// Generate generates an image using the best available remote API
func (u *UnifiedGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	// Prioritize Cloudflare for specific models or if explicitly requested
	if u.cloudflare.IsAvailable() && (opts.Model == "@cf/black-forest-labs/flux-2-klein-9b" || strings.HasPrefix(opts.Model, "@cf/")) {
		log.Printf("[RemoteGen] Using Cloudflare API for model: %s", opts.Model)
		return u.cloudflare.Generate(ctx, opts)
	}

	// Determine which generator to use based on model
	if opts.Model == "" || opts.Model == "gemini-2.5-flash" || opts.Model == "gemini-2.5-flash-image" {
		// Try Gemini first (if available)
		if u.gemini.IsAvailable() {
			log.Printf("[RemoteGen] Using Gemini API")
			return u.gemini.Generate(ctx, opts)
		}

		// Fallback to Cloudflare if Gemini is down and model is generic
		if u.cloudflare.IsAvailable() {
			log.Printf("[RemoteGen] Gemini unavailable, falling back to Cloudflare")
			return u.cloudflare.Generate(ctx, opts)
		}

		return fmt.Errorf("no remote generators available")
	}

	// For Cloudflare models
	if u.cloudflare.IsAvailable() {
		log.Printf("[RemoteGen] Using Cloudflare API for model: %s", opts.Model)
		return u.cloudflare.Generate(ctx, opts)
	}

	return fmt.Errorf("no suitable generator available for model: %s", opts.Model)
}

// GenerateWithFallback generates an image with automatic fallback
func (u *UnifiedGenerator) GenerateWithFallback(ctx context.Context, opts GenerationOptions) error {
	// Try Gemini first
	if u.gemini.IsAvailable() {
		err := u.gemini.Generate(ctx, opts)
		if err == nil {
			return nil
		}
		log.Printf("[RemoteGen] Gemini failed: %v, trying Cloudflare", err)
	}

	// Try Cloudflare
	if u.cloudflare.IsAvailable() {
		err := u.cloudflare.Generate(ctx, opts)
		if err == nil {
			return nil
		}
		log.Printf("[RemoteGen] Cloudflare failed: %v", err)
		return err
	}

	return fmt.Errorf("all remote generators failed or unavailable")
}

// GetAvailableModels returns all available models from all generators
func (u *UnifiedGenerator) GetAvailableModels() []string {
	models := make([]string, 0)

	if u.cloudflare.IsAvailable() {
		models = append(models, u.cloudflare.GetAvailableModels()...)
	}

	if u.gemini.IsAvailable() {
		models = append(models, u.gemini.GetAvailableModels()...)
	}

	return models
}

// IsAvailable checks if any remote generator is available
func (u *UnifiedGenerator) IsAvailable() bool {
	return u.gemini.IsAvailable() || u.cloudflare.IsAvailable()
}

// GetGenerator returns a specific generator by name
func (u *UnifiedGenerator) GetGenerator(name string) (Generator, error) {
	switch name {
	case "gemini":
		if !u.gemini.IsAvailable() {
			return nil, fmt.Errorf("gemini generator not available")
		}
		return u.gemini, nil
	case "cloudflare":
		if !u.cloudflare.IsAvailable() {
			return nil, fmt.Errorf("cloudflare generator not available")
		}
		return u.cloudflare, nil
	default:
		return nil, fmt.Errorf("unknown generator: %s", name)
	}
}
