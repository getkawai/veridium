package image

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kawai-network/veridium/pkg/stablediffusion/remote"
)

// RuntimeImageGenParams matches frontend model-bank/standard-parameters RuntimeImageGenParams
type RuntimeImageGenParams struct {
	Prompt      string   `json:"prompt"`
	ImageUrl    *string  `json:"imageUrl,omitempty"`
	ImageUrls   []string `json:"imageUrls,omitempty"`
	Width       *int     `json:"width,omitempty"`
	Height      *int     `json:"height,omitempty"`
	Size        string   `json:"size,omitempty"`
	AspectRatio string   `json:"aspectRatio,omitempty"`
	Cfg         *float64 `json:"cfg,omitempty"`
	Strength    *float64 `json:"strength,omitempty"`
	Steps       *int     `json:"steps,omitempty"`
	Quality     string   `json:"quality,omitempty"`
	Seed        *int64   `json:"seed,omitempty"`
	SamplerName string   `json:"samplerName,omitempty"`
	Scheduler   string   `json:"scheduler,omitempty"`
}

// CreateImageRequest matches frontend createImage action parameters
type CreateImageRequest struct {
	GenerationTopicId string                `json:"generationTopicId"`
	Provider          string                `json:"provider"`
	Model             string                `json:"model"`
	ImageNum          int                   `json:"imageNum"`
	Params            RuntimeImageGenParams `json:"params"`
}

// GenerationOptions defines internal options for SD binary execution
type GenerationOptions struct {
	Prompt         string
	NegativePrompt string
	ModelPath      string
	OutputPath     string
	ImageUrl       *string
	ImageUrls      []string
	Width          int
	Height         int
	Size           string
	AspectRatio    string
	Steps          int
	Cfg            float64
	Strength       float64
	Seed           *int64
	Quality        string
	SamplerName    string
	Scheduler      string
	OutputFormat   string
	Model          string // Model name for remote API
}

// CreateImage handles frontend CreateImageRequest and generates images asynchronously
// MOVED TO internal/image/service.go

// CreateImageWithOptions generates an image using GenerationOptions directly
// Used by internal services like image_designer
// This method now delegates to the appropriate generator (local or remote)
func (sdrm *StableDiffusion) CreateImageWithOptions(opts GenerationOptions) error {
	ctx := context.Background()

	// Use remote generation by default (Gemini API)
	remoteGen := remote.NewGenerator()

	// Convert to remote.GenerationOptions
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

	return remoteGen.Generate(ctx, remoteOpts)
}

// generateImageRemote generates an image using remote APIs (wrapper for backward compatibility)
func (sdrm *StableDiffusion) generateImageRemote(opts GenerationOptions) error {
	return sdrm.CreateImageWithOptions(opts)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
