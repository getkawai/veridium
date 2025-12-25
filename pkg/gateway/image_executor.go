package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/image"
)

// ImageExecutor defines the interface for image generation execution
type ImageExecutor interface {
	GenerateImage(ctx context.Context, req ImageGenerationRequest) ([]ImageData, error)
}

// SDLocalExecutor implements ImageExecutor using local Stable Diffusion
type SDLocalExecutor struct {
	engine *image.StableDiffusion
}

// NewSDLocalExecutor creates a new SDLocalExecutor
func NewSDLocalExecutor(engine *image.StableDiffusion) *SDLocalExecutor {
	return &SDLocalExecutor{
		engine: engine,
	}
}

// GenerateImage generates images using local SD engine
func (e *SDLocalExecutor) GenerateImage(ctx context.Context, req ImageGenerationRequest) ([]ImageData, error) {
	// 1. Prepare options
	outputDir := e.engine.GetOutputsPath()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Resolve model
	modelPath := e.engine.GetFirstAvailableModel()
	if req.Model != "" {
		// TODO: Implement specific model selection logic if needed
		// For now, we trust the user or strictly use available one
		// Ideally we check if req.Model matches an available file
	}
	if modelPath == "" {
		return nil, fmt.Errorf("no stable diffusion model available")
	}

	// Parse dimensions
	width, height := 1024, 1024
	if req.Size != "" {
		// Simple parsing, can be more robust
		fmt.Sscanf(req.Size, "%dx%d", &width, &height)
	}

	n := 1
	if req.N > 0 {
		n = req.N
	}

	results := make([]ImageData, 0, n)

	for i := 0; i < n; i++ {
		// Generate unique filename
		imageID := uuid.New().String()
		fileName := fmt.Sprintf("gen_%s.png", imageID)
		outputPath := filepath.Join(outputDir, fileName)

		// Set seed
		seed := time.Now().UnixNano()

		// Build options
		opts := image.GenerationOptions{
			Prompt:     req.Prompt,
			ModelPath:  modelPath,
			OutputPath: outputPath,
			Width:      width,
			Height:     height,
			Steps:      20, // Default steps
			Cfg:        7.0,
			Seed:       &seed,
		}

		if req.Quality == "hd" {
			opts.Steps = 30
		}

		// Execute generation
		if err := e.engine.CreateImageWithOptions(opts); err != nil {
			return nil, fmt.Errorf("generation failed for image %d: %w", i, err)
		}

		// Handle Response Format
		imageData := ImageData{}

		if req.ResponseFormat == "b64_json" {
			bytes, err := os.ReadFile(outputPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read generated image: %w", err)
			}
			imageData.B64JSON = base64.StdEncoding.EncodeToString(bytes)

			// Clean up file if only b64_json is requested?
			// OpenAI usually keeps it ephemeral, but we keep it in outputs for now or clean up.
			// Let's keep it in outputs for cache/debugging.
		} else {
			// URL format - using the file server route
			imageData.URL = fmt.Sprintf("/files/%s", fileName)
		}

		imageData.RevisedPrompt = req.Prompt // We mimic revised prompt
		results = append(results, imageData)
	}

	return results, nil
}
