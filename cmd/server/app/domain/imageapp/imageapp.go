// Package imageapp provides the image generation api endpoints.
package imageapp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	sd "github.com/kawai-network/veridium/pkg/stablediffusion"
)

type app struct {
	log    *logger.Logger
	engine *sd.StableDiffusion
}

func newApp(cfg Config) *app {
	return &app{
		log:    cfg.Log,
		engine: cfg.Engine,
	}
}

// ImageGenerationRequest represents an OpenAI-compatible image generation request
// Reused from pkg/gateway/image_types.go
type ImageGenerationRequest struct {
	Prompt         string `json:"prompt" binding:"required"`
	Model          string `json:"model"`
	N              int    `json:"n,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Size           string `json:"size,omitempty"`
	Style          string `json:"style,omitempty"`
	User           string `json:"user,omitempty"`
}

// ImageGenerationResponse represents an OpenAI-compatible image generation response
type ImageGenerationResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

// Encode implements web.Encoder.
func (r ImageGenerationResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}

// ImageData represents a single generated image
type ImageData struct {
	B64JSON       string `json:"b64_json,omitempty"`
	URL           string `json:"url,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

func (a *app) generations(ctx context.Context, r *http.Request) web.Encoder {
	if a.engine == nil {
		return errs.Errorf(errs.Unimplemented, "image generation service not available")
	}

	var req ImageGenerationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if req.Prompt == "" {
		return errs.Errorf(errs.InvalidArgument, "prompt is required")
	}

	// OpenAI spec: n defaults to 1 if not provided
	if req.N == 0 {
		req.N = 1
	}

	// Validate N upper bound (OpenAI DALL-E 3 caps at 10)
	if req.N > 10 {
		return errs.Errorf(errs.InvalidArgument, "n must be between 1 and 10")
	}
	if req.N < 1 {
		return errs.Errorf(errs.InvalidArgument, "n must be at least 1")
	}

	a.log.Info(ctx, "image generation", "prompt", req.Prompt, "n", req.N, "size", req.Size)

	// Generate images (reused logic from pkg/gateway/image_executor.go)
	data, err := a.generateImages(ctx, req)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return &ImageGenerationResponse{
		Created: time.Now().Unix(),
		Data:    data,
	}
}

// generateImages implements the image generation logic
// Reused and adapted from pkg/gateway/image_executor.go
func (a *app) generateImages(ctx context.Context, req ImageGenerationRequest) ([]ImageData, error) {
	// Use default output directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	outputDir := filepath.Join(homeDir, ".stable-diffusion", "outputs")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Parse dimensions
	width, height := 1024, 1024
	if req.Size != "" {
		var w, h int
		n, err := fmt.Sscanf(req.Size, "%dx%d", &w, &h)
		if err != nil || n != 2 || w <= 0 || h <= 0 {
			return nil, fmt.Errorf("invalid size format, expected WIDTHxHEIGHT (e.g., 1024x1024)")
		}
		width, height = w, h
	}

	results := make([]ImageData, 0, req.N)

	for i := 0; i < req.N; i++ {
		imageID := uuid.New().String()
		fileName := fmt.Sprintf("gen_%s.png", imageID)
		outputPath := filepath.Join(outputDir, fileName)

		// Map to pkg/stablediffusion Params
		params := &sd.ImgGenParams{
			Prompt:      req.Prompt,
			Width:       int32(width),
			Height:      int32(height),
			SampleSteps: 20,
			CfgScale:    7.0,
			Seed:        0, // Random
		}

		if req.Quality == "hd" {
			params.SampleSteps = 30
		}

		a.log.Info(ctx, "generating image", "id", imageID, "params", params)

		// Generate using the library backend
		if err := a.engine.GenerateImage(params, outputPath); err != nil {
			return nil, fmt.Errorf("generation failed for image %d: %w", i, err)
		}

		imageData := ImageData{}

		if req.ResponseFormat == "b64_json" {
			bytes, err := os.ReadFile(outputPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read generated image: %w", err)
			}
			imageData.B64JSON = base64.StdEncoding.EncodeToString(bytes)
		} else {
			imageData.URL = fmt.Sprintf("/files/%s", fileName)
		}

		imageData.RevisedPrompt = req.Prompt
		results = append(results, imageData)
	}

	return results, nil
}
