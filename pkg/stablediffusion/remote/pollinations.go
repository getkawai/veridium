package remote

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
)

// PollinationsGenerator handles image generation using Pollinations AI
type PollinationsGenerator struct {
	baseURL string
}

// NewPollinationsGenerator creates a new Pollinations image generator
func NewPollinationsGenerator() *PollinationsGenerator {
	return &PollinationsGenerator{
		baseURL: "https://image.pollinations.ai/prompt/",
	}
}

// Generate generates an image using Pollinations AI API
func (p *PollinationsGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	// Create context with timeout to prevent indefinite hangs
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Build the API URL
	encodedPrompt := url.PathEscape(opts.Prompt)
	apiURL := p.baseURL + encodedPrompt

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
		apiURL += "?" + params.Encode()
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

// GetAvailableModels returns list of available Pollinations models
func (p *PollinationsGenerator) GetAvailableModels() []string {
	return []string{
		"flux",
		"flux-realism",
		"flux-anime",
		"flux-3d",
		"any-dark",
		"turbo",
	}
}

// IsAvailable checks if Pollinations generation is available
func (p *PollinationsGenerator) IsAvailable() bool {
	// Pollinations is always available (no API key required)
	return true
}
