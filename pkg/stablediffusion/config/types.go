// Package config provides configuration management for stable-diffusion-go
package config

import (
	"fmt"
	"time"
)

// Version is the current config version
const Version = "1.0"

// Config represents the full configuration for stable-diffusion-go
type Config struct {
	Version     string            `yaml:"version" json:"version"`
	Models      ModelConfig       `yaml:"models" json:"models"`
	Generation  GenerationConfig  `yaml:"generation" json:"generation"`
	Output      OutputConfig      `yaml:"output" json:"output"`
	Performance PerformanceConfig `yaml:"performance" json:"performance"`
	Logging     LoggingConfig     `yaml:"logging" json:"logging"`
}

// ModelConfig contains model paths and settings
type ModelConfig struct {
	DiffusionModelPath          string            `yaml:"diffusion_model" json:"diffusion_model"`
	LLMPath                     string            `yaml:"llm_model" json:"llm_model"`
	VAEPath                     string            `yaml:"vae_model" json:"vae_model"`
	T5XXLPath                   string            `yaml:"t5xxl_model" json:"t5xxl_model"`
	ClipLPath                   string            `yaml:"clipl_model" json:"clipl_model"`
	ClipGPath                   string            `yaml:"clipg_model" json:"clipg_model"`
	ClipVisionPath              string            `yaml:"clip_vision_model" json:"clip_vision_model"`
	LLMVisionPath               string            `yaml:"llm_vision_model" json:"llm_vision_model"`
	HighNoiseDiffusionModelPath string            `yaml:"high_noise_diffusion_model" json:"high_noise_diffusion_model"`
	TAESDPath                   string            `yaml:"taesd_model" json:"taesd_model"`
	ControlNetPath              string            `yaml:"controlnet_model" json:"controlnet_model"`
	PhotoMakerPath              string            `yaml:"photomaker_model" json:"photomaker_model"`
	Embeddings                  []EmbeddingConfig `yaml:"embeddings" json:"embeddings"`
}

// EmbeddingConfig represents a single embedding configuration
type EmbeddingConfig struct {
	Name string `yaml:"name" json:"name"`
	Path string `yaml:"path" json:"path"`
}

// GenerationConfig contains default generation parameters
type GenerationConfig struct {
	DefaultPrompt           string  `yaml:"default_prompt" json:"default_prompt"`
	DefaultNegativePrompt   string  `yaml:"default_negative_prompt" json:"default_negative_prompt"`
	DefaultWidth            int32   `yaml:"default_width" json:"default_width"`
	DefaultHeight           int32   `yaml:"default_height" json:"default_height"`
	DefaultSampleSteps      int32   `yaml:"default_sample_steps" json:"default_sample_steps"`
	DefaultCfgScale         float32 `yaml:"default_cfg_scale" json:"default_cfg_scale"`
	DefaultImageCfgScale    float32 `yaml:"default_image_cfg_scale" json:"default_image_cfg_scale"`
	DefaultDistilledGuidance float32 `yaml:"default_distilled_guidance" json:"default_distilled_guidance"`
	DefaultSampleMethod     string  `yaml:"default_sample_method" json:"default_sample_method"`
	DefaultScheduler        string  `yaml:"default_scheduler" json:"default_scheduler"`
	DefaultPrediction       string  `yaml:"default_prediction" json:"default_prediction"`
	DefaultRNGType          string  `yaml:"default_rng_type" json:"default_rng_type"`
	DefaultWType            string  `yaml:"default_wtype" json:"default_wtype"`
	DefaultClipSkip         int32   `yaml:"default_clip_skip" json:"default_clip_skip"`
	DefaultStrength         float32 `yaml:"default_strength" json:"default_strength"`
	DefaultSeed             int64   `yaml:"default_seed" json:"default_seed"`
	DefaultBatchCount       int32   `yaml:"default_batch_count" json:"default_batch_count"`
}

// OutputConfig contains output settings
type OutputConfig struct {
	OutputDir       string `yaml:"output_dir" json:"output_dir"`
	NamingPattern   string `yaml:"naming_pattern" json:"naming_pattern"`
	AutoTimestamp   bool   `yaml:"auto_timestamp" json:"auto_timestamp"`
	CreateSubdirs   bool   `yaml:"create_subdirs" json:"create_subdirs"`
	SubdirPattern   string `yaml:"subdir_pattern" json:"subdir_pattern"`
}

// PerformanceConfig contains performance-related settings
type PerformanceConfig struct {
	NThreads           int32  `yaml:"n_threads" json:"n_threads"`
	OffloadParamsToCPU bool   `yaml:"offload_to_cpu" json:"offload_to_cpu"`
	EnableMmap         bool   `yaml:"enable_mmap" json:"enable_mmap"`
	DiffusionFlashAttn bool   `yaml:"flash_attention" json:"flash_attention"`
	KeepClipOnCPU      bool   `yaml:"keep_clip_on_cpu" json:"keep_clip_on_cpu"`
	KeepControlNetOnCPU bool  `yaml:"keep_controlnet_on_cpu" json:"keep_controlnet_on_cpu"`
	KeepVAEOnCPU       bool   `yaml:"keep_vae_on_cpu" json:"keep_vae_on_cpu"`
	DiffusionConvDirect bool  `yaml:"diffusion_conv_direct" json:"diffusion_conv_direct"`
	VAEConvDirect      bool   `yaml:"vae_conv_direct" json:"vae_conv_direct"`
	TAEPreviewOnly     bool   `yaml:"tae_preview_only" json:"tae_preview_only"`
	FreeParamsImmediately bool `yaml:"free_params_immediately" json:"free_params_immediately"`
	VAEDecodeOnly      bool   `yaml:"vae_decode_only" json:"vae_decode_only"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level     string `yaml:"level" json:"level"`
	LogToFile bool   `yaml:"log_to_file" json:"log_to_file"`
	LogFile   string `yaml:"log_file" json:"log_file"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Version == "" {
		c.Version = Version
	}

	// Validate models
	if err := c.Models.Validate(); err != nil {
		return fmt.Errorf("models validation failed: %w", err)
	}

	// Validate generation config
	if err := c.Generation.Validate(); err != nil {
		return fmt.Errorf("generation validation failed: %w", err)
	}

	// Validate output config
	if err := c.Output.Validate(); err != nil {
		return fmt.Errorf("output validation failed: %w", err)
	}

	// Validate performance config
	if err := c.Performance.Validate(); err != nil {
		return fmt.Errorf("performance validation failed: %w", err)
	}

	return nil
}

// Validate checks if model configuration is valid
// At least one model path must be specified for the configuration to be usable
func (m *ModelConfig) Validate() error {
	// At least one diffusion model should be specified
	if m.DiffusionModelPath == "" {
		return fmt.Errorf("diffusion model path must be specified (models.diffusion_model in config)")
	}
	return nil
}

// GetPrimaryModelPath returns the primary model path (diffusion or full model)
func (m *ModelConfig) GetPrimaryModelPath() string {
	if m.DiffusionModelPath != "" {
		return m.DiffusionModelPath
	}
	return "" // Will be handled by ContextParams
}

// Validate checks if generation configuration is valid
func (g *GenerationConfig) Validate() error {
	if g.DefaultWidth <= 0 {
		g.DefaultWidth = 512
	}
	if g.DefaultHeight <= 0 {
		g.DefaultHeight = 512
	}
	if g.DefaultSampleSteps <= 0 {
		g.DefaultSampleSteps = 20
	}
	if g.DefaultCfgScale <= 0 {
		g.DefaultCfgScale = 5.0
	}
	if g.DefaultSeed < 0 {
		g.DefaultSeed = 42
	}
	if g.DefaultBatchCount <= 0 {
		g.DefaultBatchCount = 1
	}
	return nil
}

// Validate checks if output configuration is valid
func (o *OutputConfig) Validate() error {
	if o.OutputDir == "" {
		o.OutputDir = "./outputs"
	}
	if o.NamingPattern == "" {
		o.NamingPattern = "sd_{timestamp}_{seed}.png"
	}
	return nil
}

// Validate checks if performance configuration is valid
func (p *PerformanceConfig) Validate() error {
	if p.NThreads <= 0 {
		p.NThreads = -1 // Use default
	}
	return nil
}

// GenerateOutputPath generates an output path based on the naming pattern
func (o *OutputConfig) GenerateOutputPath(seed int64, extension string) string {
	pattern := o.NamingPattern
	
	if o.AutoTimestamp {
		timestamp := time.Now().Format("20060102_150405")
		pattern = replacePlaceholder(pattern, "timestamp", timestamp)
	}
	
	pattern = replacePlaceholder(pattern, "seed", fmt.Sprintf("%d", seed))
	pattern = replacePlaceholder(pattern, "date", time.Now().Format("20060102"))
	pattern = replacePlaceholder(pattern, "time", time.Now().Format("150405"))
	
	// Ensure correct extension
	if extension != "" && !hasExtension(pattern, extension) {
		pattern = pattern + "." + extension
	}
	
	return pattern
}

// replacePlaceholder replaces a placeholder in the pattern with a value
func replacePlaceholder(pattern, placeholder, value string) string {
	return fmt.Sprintf(pattern, value) // Simple implementation
}

// hasExtension checks if the filename has the given extension
func hasExtension(filename, ext string) bool {
	if len(filename) < len(ext)+1 {
		return false
	}
	return filename[len(filename)-len(ext)-1:] == "."+ext
}