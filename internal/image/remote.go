package image

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/kawai-network/veridium/internal/constant"
	"google.golang.org/genai"
)

// RemoteGenerator handles remote API-based image generation
type RemoteGenerator struct {
	// Future: could add API keys, rate limiters, etc.
}

// NewRemoteGenerator creates a new remote image generator
func NewRemoteGenerator() *RemoteGenerator {
	return &RemoteGenerator{}
}

// Generate generates an image using remote APIs (Gemini, Pollinations, etc.)
func (rg *RemoteGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	// Default to Gemini API
	if opts.Model == "" || opts.Model == "gemini-2.5-flash" || opts.Model == "gemini-2.5-flash-image" {
		return rg.generateWithGemini(ctx, opts)
	}

	// Fallback to Pollinations for other models
	return rg.generateWithPollinations(ctx, opts)
}

// generateWithGemini generates an image using Google Gemini API
// Reference: https://ai.google.dev/gemini-api/docs/image-generation#go
func (rg *RemoteGenerator) generateWithGemini(ctx context.Context, opts GenerationOptions) error {
	// Create context with timeout to prevent indefinite hangs
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Get random Gemini API key from pool (from internal/constant/llm.go)
	apiKey := constant.GetRandomGeminiApiKey()
	if apiKey == "" {
		return fmt.Errorf("no Gemini API key available")
	}

	// Create Gemini client with API key directly (no environment variable needed)
	// This is thread-safe and avoids race conditions in concurrent goroutines
	clientConfig := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	client, err := genai.NewClient(ctx, clientConfig)
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}
	// Note: genai.Client doesn't have Close() method, no cleanup needed

	// Determine model to use
	// gemini-2.5-flash-image (Nano Banana) - fast, 1024px, free tier quota available
	model := "gemini-2.5-flash-image"

	// Priority: explicit model choice
	if opts.Model != "" {
		// User explicitly specified a model - respect their choice
		if opts.Model == "gemini-2.5-flash" || opts.Model == "gemini-2.5-flash-image" {
			model = "gemini-2.5-flash-image"
		}
		// Note: gemini-3-pro-image-preview removed - no free tier quota available
	}

	log.Printf("[Gemini] Using model: %s for prompt: %s", model, opts.Prompt)

	// Build generation config with aspect ratio
	config := &genai.GenerateContentConfig{
		ImageConfig: &genai.ImageConfig{},
	}

	// Map dimensions to aspect ratio
	aspectRatio := rg.calculateAspectRatio(opts)
	config.ImageConfig.AspectRatio = aspectRatio

	log.Printf("[Gemini] Aspect ratio: %s", aspectRatio)

	// Generate content with image
	result, err := client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(opts.Prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("Gemini API generation failed: %w", err)
	}

	// Extract image data from response
	if len(result.Candidates) == 0 {
		return fmt.Errorf("no candidates returned from Gemini API")
	}

	var imageBytes []byte
	foundImage := false

	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			log.Printf("[Gemini] Response text: %s", part.Text)
		} else if part.InlineData != nil {
			imageBytes = part.InlineData.Data
			foundImage = true
			log.Printf("[Gemini] Received image data: %d bytes", len(imageBytes))
			break
		}
	}

	if !foundImage || len(imageBytes) == 0 {
		return fmt.Errorf("no image data returned from Gemini API")
	}

	// Write image to output file
	err = os.WriteFile(opts.OutputPath, imageBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write image file: %w", err)
	}

	log.Printf("[Gemini] Image saved successfully to: %s", opts.OutputPath)
	return nil
}

// calculateAspectRatio determines the aspect ratio from width/height or explicit setting
func (rg *RemoteGenerator) calculateAspectRatio(opts GenerationOptions) string {
	aspectRatio := opts.AspectRatio
	if aspectRatio == "" && opts.Width > 0 && opts.Height > 0 {
		// Calculate aspect ratio from width/height
		ratio := float64(opts.Width) / float64(opts.Height)
		switch {
		case ratio >= 0.95 && ratio <= 1.05:
			aspectRatio = "1:1"
		case ratio >= 0.65 && ratio <= 0.68:
			aspectRatio = "2:3"
		case ratio >= 1.48 && ratio <= 1.52:
			aspectRatio = "3:2"
		case ratio >= 0.72 && ratio <= 0.76:
			aspectRatio = "3:4"
		case ratio >= 1.32 && ratio <= 1.36:
			aspectRatio = "4:3"
		case ratio >= 0.77 && ratio <= 0.79:
			aspectRatio = "4:5"
		case ratio >= 1.27 && ratio <= 1.29:
			aspectRatio = "5:4"
		case ratio >= 0.56 && ratio <= 0.58:
			aspectRatio = "9:16"
		case ratio >= 1.76 && ratio <= 1.79:
			aspectRatio = "16:9"
		case ratio >= 2.32 && ratio <= 2.35:
			aspectRatio = "21:9"
		default:
			aspectRatio = "1:1" // Default square
		}
	}
	if aspectRatio == "" {
		aspectRatio = "1:1"
	}
	return aspectRatio
}

// generateWithPollinations generates an image using Pollinations AI (legacy/fallback)
func (rg *RemoteGenerator) generateWithPollinations(ctx context.Context, opts GenerationOptions) error {
	// Build the API URL
	baseURL := "https://image.pollinations.ai/prompt/"
	encodedPrompt := url.PathEscape(opts.Prompt)
	apiURL := baseURL + encodedPrompt

	// Build query parameters
	params := url.Values{}
	if opts.Width > 0 {
		params.Add("width", strconv.Itoa(opts.Width))
	}
	if opts.Height > 0 {
		params.Add("height", strconv.Itoa(opts.Height))
	}
	// Use specified model or default to flux
	model := opts.Model
	if model == "" {
		model = "flux"
	}
	params.Add("model", model)
	params.Add("enhance", "false")
	params.Add("nologo", "true")

	if len(params) > 0 {
		apiURL += "$" + params.Encode()
	}

	log.Printf("[Pollinations] Requesting image from: %s", apiURL)

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Make HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("remote API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote API returned status %d", resp.StatusCode)
	}

	// Create output file
	outFile, err := os.Create(opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Download image
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}

	log.Printf("[Pollinations] Image downloaded successfully to: %s", opts.OutputPath)
	return nil
}

// GetAvailableModels returns list of available remote models
func (rg *RemoteGenerator) GetAvailableModels() []string {
	return []string{
		"gemini-2.5-flash",       // Fast, 1024px (Nano Banana) - free tier
		"gemini-2.5-flash-image", // Explicit image model
		// Pollinations models (fallback)
		"flux",
		"flux-realism",
		"flux-anime",
		"flux-3d",
		"any-dark",
		"turbo",
	}
}

// IsAvailable checks if remote generation is available (has API keys)
func (rg *RemoteGenerator) IsAvailable() bool {
	// Check if we have Gemini API keys
	return constant.GetRandomGeminiApiKey() != ""
}
