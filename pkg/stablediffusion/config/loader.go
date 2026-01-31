package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from a YAML or JSON file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	// Apply defaults and validate
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(config *Config, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var data []byte
	var err error

	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %w", err)
		}
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadOrCreate loads config from file or creates default if not exists
func LoadOrCreate(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config
		config := DefaultConfig()
		if err := SaveConfig(config, path); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	return LoadConfig(path)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Version: Version,
		Models: ModelConfig{
			DiffusionModelPath: "",
			LLMPath:            "",
			VAEPath:            "",
			T5XXLPath:          "",
			ClipLPath:          "",
			ClipGPath:          "",
			ClipVisionPath:    "",
			LLMVisionPath:     "",
			Embeddings:        []EmbeddingConfig{},
		},
		Generation: GenerationConfig{
			DefaultPrompt:           "masterpiece, best quality",
			DefaultNegativePrompt:   "blurry, low quality, distorted, ugly",
			DefaultWidth:            512,
			DefaultHeight:           512,
			DefaultSampleSteps:      20,
			DefaultCfgScale:         5.0,
			DefaultImageCfgScale:    1.0,
			DefaultDistilledGuidance: 3.5,
			DefaultSampleMethod:     "euler",
			DefaultScheduler:        "discrete",
			DefaultPrediction:       "default",
			DefaultRNGType:          "cuda",
			DefaultWType:            "default",
			DefaultClipSkip:         -1,
			DefaultStrength:         0.75,
			DefaultSeed:             42,
			DefaultBatchCount:       1,
		},
		Output: OutputConfig{
			OutputDir:       "./outputs",
			NamingPattern:   "sd_{timestamp}_{seed}.png",
			AutoTimestamp:   true,
			CreateSubdirs:   false,
			SubdirPattern:   "{date}",
		},
		Performance: PerformanceConfig{
			NThreads:              -1,
			OffloadParamsToCPU:    true,
			EnableMmap:            false,
			DiffusionFlashAttn:    true,
			KeepClipOnCPU:         false,
			KeepControlNetOnCPU:   false,
			KeepVAEOnCPU:          false,
			DiffusionConvDirect:   false,
			VAEConvDirect:         false,
			TAEPreviewOnly:        false,
			FreeParamsImmediately: false,
			VAEDecodeOnly:         false,
		},
		Logging: LoggingConfig{
			Level:     "info",
			LogToFile: false,
			LogFile:   "sd.log",
		},
	}
}

// MergeConfigs merges two configurations, with override taking precedence
func MergeConfigs(base, override *Config) *Config {
	if override == nil {
		return base
	}

	result := *base

	// Merge models
	if override.Models.DiffusionModelPath != "" {
		result.Models.DiffusionModelPath = override.Models.DiffusionModelPath
	}
	if override.Models.LLMPath != "" {
		result.Models.LLMPath = override.Models.LLMPath
	}
	if override.Models.VAEPath != "" {
		result.Models.VAEPath = override.Models.VAEPath
	}
	if override.Models.T5XXLPath != "" {
		result.Models.T5XXLPath = override.Models.T5XXLPath
	}
	if len(override.Models.Embeddings) > 0 {
		result.Models.Embeddings = override.Models.Embeddings
	}

	// Merge generation settings (only if explicitly set)
	if override.Generation.DefaultPrompt != "" {
		result.Generation.DefaultPrompt = override.Generation.DefaultPrompt
	}
	if override.Generation.DefaultNegativePrompt != "" {
		result.Generation.DefaultNegativePrompt = override.Generation.DefaultNegativePrompt
	}
	if override.Generation.DefaultWidth > 0 {
		result.Generation.DefaultWidth = override.Generation.DefaultWidth
	}
	if override.Generation.DefaultHeight > 0 {
		result.Generation.DefaultHeight = override.Generation.DefaultHeight
	}
	if override.Generation.DefaultSampleSteps > 0 {
		result.Generation.DefaultSampleSteps = override.Generation.DefaultSampleSteps
	}
	if override.Generation.DefaultCfgScale > 0 {
		result.Generation.DefaultCfgScale = override.Generation.DefaultCfgScale
	}

	// Merge output settings
	if override.Output.OutputDir != "" {
		result.Output.OutputDir = override.Output.OutputDir
	}
	if override.Output.NamingPattern != "" {
		result.Output.NamingPattern = override.Output.NamingPattern
	}

	// Merge performance settings
	if override.Performance.NThreads != 0 {
		result.Performance.NThreads = override.Performance.NThreads
	}

	return &result
}

// ConfigManager provides a higher-level interface for config management
type ConfigManager struct {
	config   *Config
	path     string
	autoSave bool
}

// NewConfigManager creates a new config manager
func NewConfigManager(path string, autoSave bool) (*ConfigManager, error) {
	config, err := LoadOrCreate(path)
	if err != nil {
		return nil, err
	}

	return &ConfigManager{
		config:   config,
		path:     path,
		autoSave: autoSave,
	}, nil
}

// Get returns the current configuration
func (cm *ConfigManager) Get() *Config {
	return cm.config
}

// Set updates the configuration
func (cm *ConfigManager) Set(config *Config) error {
	cm.config = config
	if cm.autoSave {
		return cm.Save()
	}
	return nil
}

// Save persists the current configuration
func (cm *ConfigManager) Save() error {
	return SaveConfig(cm.config, cm.path)
}

// Update updates specific fields and optionally saves
func (cm *ConfigManager) Update(updater func(*Config)) error {
	updater(cm.config)
	if err := cm.config.Validate(); err != nil {
		return err
	}
	if cm.autoSave {
		return cm.Save()
	}
	return nil
}