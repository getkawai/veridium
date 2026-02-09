package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
)

// ToContextParams converts Config to ContextParams
func (c *Config) ToContextParams() (*stablediffusion.ContextParams, error) {
	params := &stablediffusion.ContextParams{}

	// Model paths
	params.DiffusionModelPath = c.Models.DiffusionModelPath
	params.LLMPath = c.Models.LLMPath
	params.VAEPath = c.Models.VAEPath
	params.T5XXLPath = c.Models.T5XXLPath
	params.ClipLPath = c.Models.ClipLPath
	params.ClipGPath = c.Models.ClipGPath
	params.ClipVisionPath = c.Models.ClipVisionPath
	params.LLMVisionPath = c.Models.LLMVisionPath
	params.HighNoiseDiffusionModelPath = c.Models.HighNoiseDiffusionModelPath
	params.TAESDPath = c.Models.TAESDPath
	params.ControlNetPath = c.Models.ControlNetPath
	params.PhotoMakerPath = c.Models.PhotoMakerPath

	// Performance settings
	params.NThreads = c.Performance.NThreads
	params.OffloadParamsToCPU = c.Performance.OffloadParamsToCPU
	params.EnableMmap = c.Performance.EnableMmap
	params.DiffusionFlashAttn = c.Performance.DiffusionFlashAttn
	params.KeepClipOnCPU = c.Performance.KeepClipOnCPU
	params.KeepControlNetOnCPU = c.Performance.KeepControlNetOnCPU
	params.KeepVAEOnCPU = c.Performance.KeepVAEOnCPU
	params.DiffusionConvDirect = c.Performance.DiffusionConvDirect
	params.VAEConvDirect = c.Performance.VAEConvDirect
	params.TAEPreviewOnly = c.Performance.TAEPreviewOnly
	params.FreeParamsImmediately = c.Performance.FreeParamsImmediately
	params.VAEDecodeOnly = c.Performance.VAEDecodeOnly

	// String-based settings
	params.WType = c.Generation.DefaultWType
	params.RNGType = c.Generation.DefaultRNGType
	params.SamplerRNGType = c.Generation.DefaultRNGType
	params.Prediction = c.Generation.DefaultPrediction
	params.LoraApplyMode = "auto"

	// Embeddings
	if len(c.Models.Embeddings) > 0 {
		params.Embeddings = &stablediffusion.Embedding{
			Name: c.Models.Embeddings[0].Name,
			Path: c.Models.Embeddings[0].Path,
		}
		params.EmbeddingCount = uint32(len(c.Models.Embeddings))
	}

	return params, nil
}

// ToImgGenParams converts Config to ImgGenParams with optional overrides
func (c *Config) ToImgGenParams(overrides *stablediffusion.ImgGenParams) *stablediffusion.ImgGenParams {
	params := &stablediffusion.ImgGenParams{
		Prompt:            c.Generation.DefaultPrompt,
		NegativePrompt:    c.Generation.DefaultNegativePrompt,
		Width:             c.Generation.DefaultWidth,
		Height:            c.Generation.DefaultHeight,
		SampleSteps:       c.Generation.DefaultSampleSteps,
		CfgScale:          c.Generation.DefaultCfgScale,
		ImageCfgScale:     c.Generation.DefaultImageCfgScale,
		DistilledGuidance: c.Generation.DefaultDistilledGuidance,
		SampleMethod:      c.Generation.DefaultSampleMethod,
		Scheduler:         c.Generation.DefaultScheduler,
		ClipSkip:          c.Generation.DefaultClipSkip,
		Strength:          c.Generation.DefaultStrength,
		Seed:              c.Generation.DefaultSeed,
		BatchCount:        c.Generation.DefaultBatchCount,
	}

	// Apply overrides if provided
	if overrides != nil {
		if overrides.Prompt != "" {
			params.Prompt = overrides.Prompt
		}
		if overrides.NegativePrompt != "" {
			params.NegativePrompt = overrides.NegativePrompt
		}
		if overrides.Width > 0 {
			params.Width = overrides.Width
		}
		if overrides.Height > 0 {
			params.Height = overrides.Height
		}
		if overrides.SampleSteps > 0 {
			params.SampleSteps = overrides.SampleSteps
		}
		if overrides.CfgScale > 0 {
			params.CfgScale = overrides.CfgScale
		}
		if overrides.SampleMethod != "" {
			params.SampleMethod = overrides.SampleMethod
		}
		if overrides.Scheduler != "" {
			params.Scheduler = overrides.Scheduler
		}
		if overrides.Seed != 0 {
			params.Seed = overrides.Seed
		}
		if overrides.InitImagePath != "" {
			params.InitImagePath = overrides.InitImagePath
		}
		if overrides.MaskImagePath != "" {
			params.MaskImagePath = overrides.MaskImagePath
		}
	}

	return params
}

// ToVidGenParams converts Config to VidGenParams with optional overrides
func (c *Config) ToVidGenParams(overrides *stablediffusion.VidGenParams) *stablediffusion.VidGenParams {
	params := &stablediffusion.VidGenParams{
		Prompt:            c.Generation.DefaultPrompt,
		NegativePrompt:    c.Generation.DefaultNegativePrompt,
		Width:             c.Generation.DefaultWidth,
		Height:            c.Generation.DefaultHeight,
		SampleSteps:       c.Generation.DefaultSampleSteps,
		CfgScale:          c.Generation.DefaultCfgScale,
		ImageCfgScale:     c.Generation.DefaultImageCfgScale,
		DistilledGuidance: c.Generation.DefaultDistilledGuidance,
		SampleMethod:      c.Generation.DefaultSampleMethod,
		Scheduler:         c.Generation.DefaultScheduler,
		ClipSkip:          c.Generation.DefaultClipSkip,
		Strength:          c.Generation.DefaultStrength,
		Seed:              c.Generation.DefaultSeed,
		VideoFrames:       33, // Default video frames
	}

	// Apply overrides if provided
	if overrides != nil {
		if overrides.Prompt != "" {
			params.Prompt = overrides.Prompt
		}
		if overrides.Width > 0 {
			params.Width = overrides.Width
		}
		if overrides.Height > 0 {
			params.Height = overrides.Height
		}
		if overrides.SampleSteps > 0 {
			params.SampleSteps = overrides.SampleSteps
		}
		if overrides.VideoFrames > 0 {
			params.VideoFrames = overrides.VideoFrames
		}
		if overrides.InitImagePath != "" {
			params.InitImagePath = overrides.InitImagePath
		}
		if overrides.EndImagePath != "" {
			params.EndImagePath = overrides.EndImagePath
		}
	}

	return params
}

// GenerateOutputPath generates a complete output path based on config
func (c *Config) GenerateOutputPath(seed int64, extension string) string {
	// Create output directory if needed
	outputDir := c.Output.OutputDir
	if c.Output.CreateSubdirs && c.Output.SubdirPattern != "" {
		subdir := generateSubdirName(c.Output.SubdirPattern)
		outputDir = filepath.Join(outputDir, subdir)
	}

	// Ensure directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		// Log error but don't fail - return path anyway
		log.Printf("warning: failed to create output directory: %v", err)
	}

	// Generate filename
	filename := c.Output.GenerateOutputPath(seed, extension)

	return filepath.Join(outputDir, filename)
}

// generateSubdirName generates a subdirectory name based on pattern
func generateSubdirName(pattern string) string {
	now := time.Now()

	replacements := map[string]string{
		"{date}":   now.Format("20060102"),
		"{year}":   now.Format("2006"),
		"{month}":  now.Format("01"),
		"{day}":    now.Format("02"),
		"{time}":   now.Format("150405"),
		"{hour}":   now.Format("15"),
		"{minute}": now.Format("04"),
	}

	result := pattern
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// ValidateModelPaths checks if all configured model paths exist
func (c *Config) ValidateModelPaths() []error {
	var errors []error

	paths := map[string]string{
		"diffusion_model": c.Models.DiffusionModelPath,
		"llm_model":       c.Models.LLMPath,
		"vae_model":       c.Models.VAEPath,
		"t5xxl_model":     c.Models.T5XXLPath,
	}

	for name, path := range paths {
		if path == "" {
			continue // Skip empty paths
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("%s not found: %s", name, path))
		}
	}

	// Validate embeddings
	for _, emb := range c.Models.Embeddings {
		if emb.Path != "" {
			if _, err := os.Stat(emb.Path); os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("embedding not found: %s (%s)", emb.Name, emb.Path))
			}
		}
	}

	return errors
}

// GetEffectivePrompt returns the effective prompt (combining default and specific)
func (c *Config) GetEffectivePrompt(specificPrompt string) string {
	if specificPrompt == "" {
		return c.Generation.DefaultPrompt
	}
	if c.Generation.DefaultPrompt == "" {
		return specificPrompt
	}
	// Combine default and specific prompt
	return c.Generation.DefaultPrompt + ", " + specificPrompt
}

// GetEffectiveNegativePrompt returns the effective negative prompt
func (c *Config) GetEffectiveNegativePrompt(specificNegPrompt string) string {
	if specificNegPrompt == "" {
		return c.Generation.DefaultNegativePrompt
	}
	if c.Generation.DefaultNegativePrompt == "" {
		return specificNegPrompt
	}
	// Combine default and specific negative prompt
	return c.Generation.DefaultNegativePrompt + ", " + specificNegPrompt
}
