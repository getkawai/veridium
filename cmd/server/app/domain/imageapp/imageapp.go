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
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/internal/paths"
	sd "github.com/kawai-network/veridium/pkg/stablediffusion"
)

type app struct {
	log        *logger.Logger
	engine     *sd.StableDiffusion
	editEngine *sd.StableDiffusion
	mu         sync.Mutex
}

const (
	// Qwen-Image 2512 stable-diffusion.cpp baseline.
	qwenImageDefaultSteps = 40
	qwenImageHDSteps      = 50
	qwenImageCFGScale     = 2.5
	qwenImageFallback1Dim = 768
	qwenImageFallback2Dim = 512
)

type imageFallbackPreset struct {
	name      string
	steps     int32
	cfgScale  float32
	maxDim    int32
	forceSize bool
}

func newApp(cfg Config) *app {
	return &app{
		log:        cfg.Log,
		engine:     cfg.Engine,
		editEngine: cfg.EditEngine,
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
	outputDir := paths.StableDiffusionOutputs()
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
			SampleSteps: qwenImageDefaultSteps,
			CfgScale:    qwenImageCFGScale,
			Seed:        -1, // Random
		}

		if req.Quality == "hd" {
			params.SampleSteps = qwenImageHDSteps
		}

		a.log.Info(ctx, "generating image", "id", imageID, "params", params)

		// Generate using the library backend
		if err := a.generateImage(ctx, a.engine, params, outputPath); err != nil {
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

// ImageEditRequest represents an OpenAI-compatible image edit request
// POST /v1/images/edits
type ImageEditRequest struct {
	Prompt         string `json:"prompt" binding:"required"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Size           string `json:"size,omitempty"`
	User           string `json:"user,omitempty"`
	// Image is the base64 encoded image data or file path
	Image string `json:"image,omitempty"`
	// Mask is the base64 encoded mask image or file path (optional)
	Mask string `json:"mask,omitempty"`
}

// edits handles POST /v1/images/edits
func (a *app) edits(ctx context.Context, r *http.Request) web.Encoder {
	if a.editEngine == nil {
		return errs.Errorf(errs.Unimplemented, "image editing service not available")
	}

	var req ImageEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if req.Prompt == "" {
		return errs.Errorf(errs.InvalidArgument, "prompt is required")
	}

	if req.Image == "" {
		return errs.Errorf(errs.InvalidArgument, "image is required")
	}

	// OpenAI spec: n defaults to 1 if not provided
	if req.N == 0 {
		req.N = 1
	}

	if req.N > 10 {
		return errs.Errorf(errs.InvalidArgument, "n must be between 1 and 10")
	}
	if req.N < 1 {
		return errs.Errorf(errs.InvalidArgument, "n must be at least 1")
	}

	a.log.Info(ctx, "image edit request", "prompt", req.Prompt, "n", req.N, "size", req.Size, "has_mask", req.Mask != "")

	// Generate edited images
	data, err := a.editImages(ctx, req)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return &ImageGenerationResponse{
		Created: time.Now().Unix(),
		Data:    data,
	}
}

// editImages implements image editing logic
func (a *app) editImages(ctx context.Context, req ImageEditRequest) ([]ImageData, error) {
	outputDir := paths.StableDiffusionOutputs()
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

	// Save image from base64 or use file path
	initImagePath, err := a.saveImageFromBase64(req.Image, "edit_init_")
	if err != nil {
		return nil, fmt.Errorf("failed to save init image: %w", err)
	}
	defer os.Remove(initImagePath) // Cleanup temp file

	var maskImagePath string
	if req.Mask != "" {
		maskImagePath, err = a.saveImageFromBase64(req.Mask, "edit_mask_")
		if err != nil {
			return nil, fmt.Errorf("failed to save mask image: %w", err)
		}
		defer os.Remove(maskImagePath) // Cleanup temp file
	}

	results := make([]ImageData, 0, req.N)

	for i := 0; i < req.N; i++ {
		imageID := uuid.New().String()
		fileName := fmt.Sprintf("edit_%s.png", imageID)
		outputPath := filepath.Join(outputDir, fileName)

		// Map to pkg/stablediffusion Params
		params := &sd.ImgGenParams{
			Prompt:        req.Prompt,
			InitImagePath: initImagePath,
			MaskImagePath: maskImagePath,
			Width:         int32(width),
			Height:        int32(height),
			SampleSteps:   qwenImageDefaultSteps,
			CfgScale:      qwenImageCFGScale,
			ImageCfgScale: 1.0,
			Seed:          -1, // Random
		}

		if req.Quality == "hd" {
			params.SampleSteps = qwenImageHDSteps
		}

		a.log.Info(ctx, "editing image", "id", imageID, "params", params)

		// Generate using the library backend
		if err := a.generateImage(ctx, a.editEngine, params, outputPath); err != nil {
			return nil, fmt.Errorf("edit failed for image %d: %w", i, err)
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

// ImageVariationRequest represents an OpenAI-compatible image variation request
// POST /v1/images/variations
type ImageVariationRequest struct {
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Size           string `json:"size,omitempty"`
	User           string `json:"user,omitempty"`
	// Image is the base64 encoded image data or file path
	Image string `json:"image" binding:"required"`
}

// variations handles POST /v1/images/variations
func (a *app) variations(ctx context.Context, r *http.Request) web.Encoder {
	if a.editEngine == nil {
		return errs.Errorf(errs.Unimplemented, "image variation service not available")
	}

	var req ImageVariationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if req.Image == "" {
		return errs.Errorf(errs.InvalidArgument, "image is required")
	}

	// OpenAI spec: n defaults to 1 if not provided
	if req.N == 0 {
		req.N = 1
	}

	if req.N > 10 {
		return errs.Errorf(errs.InvalidArgument, "n must be between 1 and 10")
	}
	if req.N < 1 {
		return errs.Errorf(errs.InvalidArgument, "n must be at least 1")
	}

	a.log.Info(ctx, "image variation request", "n", req.N, "size", req.Size)

	// Generate image variations
	data, err := a.variateImages(ctx, req)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return &ImageGenerationResponse{
		Created: time.Now().Unix(),
		Data:    data,
	}
}

// variateImages implements image variation logic
func (a *app) variateImages(ctx context.Context, req ImageVariationRequest) ([]ImageData, error) {
	outputDir := paths.StableDiffusionOutputs()
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

	// Save image from base64
	initImagePath, err := a.saveImageFromBase64(req.Image, "var_init_")
	if err != nil {
		return nil, fmt.Errorf("failed to save init image: %w", err)
	}
	defer os.Remove(initImagePath) // Cleanup temp file

	results := make([]ImageData, 0, req.N)

	for i := 0; i < req.N; i++ {
		imageID := uuid.New().String()
		fileName := fmt.Sprintf("var_%s.png", imageID)
		outputPath := filepath.Join(outputDir, fileName)

		// Map to pkg/stablediffusion Params
		params := &sd.ImgGenParams{
			Prompt:        "", // No prompt for pure variation
			InitImagePath: initImagePath,
			Width:         int32(width),
			Height:        int32(height),
			SampleSteps:   qwenImageDefaultSteps,
			CfgScale:      qwenImageCFGScale,
			Strength:      0.75, // Default strength for img2img variation
			Seed:          -1,   // Random seed for variation
		}

		if req.Quality == "hd" {
			params.SampleSteps = qwenImageHDSteps
		}

		a.log.Info(ctx, "generating variation", "id", imageID, "params", params)

		// Generate using the library backend
		if err := a.generateImage(ctx, a.editEngine, params, outputPath); err != nil {
			return nil, fmt.Errorf("variation failed for image %d: %w", i, err)
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

		imageData.RevisedPrompt = "Variation of input image"
		results = append(results, imageData)
	}

	return results, nil
}

// saveImageFromBase64 saves a base64 encoded image to a temp file and returns the path
func (a *app) saveImageFromBase64(base64Data string, prefix string) (string, error) {
	// Check if it's a file path (not base64)
	if len(base64Data) < 100 || !strings.HasPrefix(base64Data, "data:") && !strings.HasPrefix(base64Data, "/") {
		// Assume it's already a file path
		if _, err := os.Stat(base64Data); err == nil {
			return base64Data, nil
		}
	}

	// Decode base64
	imageData := base64Data
	// Remove data URL prefix if present (e.g., "data:image/png;base64,")
	if idx := strings.Index(base64Data, ","); idx > 0 && strings.HasPrefix(base64Data, "data:") {
		imageData = base64Data[idx+1:]
	}

	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", prefix+"*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Write data
	if _, err := tempFile.Write(decoded); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	return tempFile.Name(), nil
}

// generateImage serializes access because StableDiffusion context is not thread-safe.
// If an OOM-like error occurs, it retries with lower-cost presets automatically.
func (a *app) generateImage(ctx context.Context, engine *sd.StableDiffusion, params *sd.ImgGenParams, outputPath string) error {
	if engine == nil {
		return fmt.Errorf("image generation engine is not available")
	}

	presets := []imageFallbackPreset{
		{name: "primary", steps: params.SampleSteps, cfgScale: params.CfgScale},
		{name: "fallback_768", steps: 28, cfgScale: 2.0, maxDim: qwenImageFallback1Dim},
		{name: "fallback_512", steps: 20, cfgScale: 1.8, maxDim: qwenImageFallback2Dim, forceSize: true},
	}

	base := *params
	var lastErr error

	for idx, preset := range presets {
		attemptParams := base
		if idx > 0 {
			attemptParams.SampleSteps = preset.steps
			attemptParams.CfgScale = preset.cfgScale
			if preset.forceSize {
				attemptParams.Width = preset.maxDim
				attemptParams.Height = preset.maxDim
			} else {
				attemptParams.Width, attemptParams.Height = scaleDimensions(base.Width, base.Height, preset.maxDim)
			}
		}

		a.mu.Lock()
		err := engine.GenerateImage(&attemptParams, outputPath)
		a.mu.Unlock()
		if err == nil {
			if idx > 0 {
				a.log.Info(ctx, "image generation fallback succeeded",
					"preset", preset.name,
					"steps", attemptParams.SampleSteps,
					"cfg_scale", attemptParams.CfgScale,
					"size", fmt.Sprintf("%dx%d", attemptParams.Width, attemptParams.Height),
				)
			}
			return nil
		}

		lastErr = err
		if !isOOMError(err) {
			return err
		}

		if idx < len(presets)-1 {
			a.log.Warn(ctx, "image generation OOM, retrying with fallback preset",
				"current_preset", preset.name,
				"error", err,
			)
		}
	}

	return fmt.Errorf("insufficient VRAM: %w", lastErr)
}

func scaleDimensions(width, height, maxDim int32) (int32, int32) {
	if maxDim <= 0 || width <= maxDim && height <= maxDim {
		return width, height
	}

	if width >= height {
		scaledHeight := int32(float64(height) * float64(maxDim) / float64(width))
		if scaledHeight < 64 {
			scaledHeight = 64
		}
		return maxDim, scaledHeight
	}

	scaledWidth := int32(float64(width) * float64(maxDim) / float64(height))
	if scaledWidth < 64 {
		scaledWidth = 64
	}
	return scaledWidth, maxDim
}

func isOOMError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	patterns := []string{
		"out of memory",
		"cuda out of memory",
		"not enough memory",
		"insufficient memory",
		"failed to allocate",
		"vram",
		"memory allocation",
	}

	for _, pattern := range patterns {
		if strings.Contains(msg, pattern) {
			return true
		}
	}

	return false
}
