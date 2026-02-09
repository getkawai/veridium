package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_ValidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `version: "1.0"
models:
  diffusion_model: "models/model.gguf"
  vae_model: "models/vae.safetensors"
generation:
  default_prompt: "masterpiece, best quality"
  default_width: 512
  default_height: 512
  default_sample_steps: 20
output:
  output_dir: "./outputs"
  naming_pattern: "sd_{timestamp}_{seed}.png"
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "1.0", cfg.Version)
	assert.Equal(t, "models/model.gguf", cfg.Models.DiffusionModelPath)
	assert.Equal(t, int32(512), cfg.Generation.DefaultWidth)
}

func TestLoadConfig_ValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	jsonContent := `{
  "version": "1.0",
  "models": {
    "diffusion_model": "models/model.gguf"
  },
  "generation": {
    "default_width": 512,
    "default_height": 512
  }
}`

	err := os.WriteFile(configPath, []byte(jsonContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "1.0", cfg.Version)
	assert.Equal(t, int32(512), cfg.Generation.DefaultWidth)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	invalidContent := `
version: "1.0"
models:
  - invalid yaml structure
    without proper indentation
`

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestSaveConfig_YAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "output.yaml")

	cfg := &Config{
		Version: "1.0",
		Models: ModelConfig{
			DiffusionModelPath: "test/model.gguf",
		},
		Generation: GenerationConfig{
			DefaultWidth:  512,
			DefaultHeight: 512,
		},
	}

	err := SaveConfig(cfg, configPath)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Load and verify content
	loadedCfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, cfg.Version, loadedCfg.Version)
	assert.Equal(t, cfg.Models.DiffusionModelPath, loadedCfg.Models.DiffusionModelPath)
}

func TestSaveConfig_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "output.json")

	cfg := &Config{
		Version: "1.0",
		Models: ModelConfig{
			DiffusionModelPath: "test/model.gguf",
		},
	}

	err := SaveConfig(cfg, configPath)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Load and verify content
	loadedCfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, cfg.Version, loadedCfg.Version)
}

func TestSaveConfig_InvalidPath(t *testing.T) {
	cfg := &Config{
		Version: "1.0",
	}

	// Use a path with invalid characters that will fail on all platforms
	// Windows: < > : " | ? * are invalid
	// Unix: null byte is invalid
	invalidPath := "invalid\x00path/config.yaml"
	err := SaveConfig(cfg, invalidPath)
	assert.Error(t, err)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Version)
	assert.NotEmpty(t, cfg.Output.OutputDir)
	assert.Greater(t, cfg.Generation.DefaultWidth, int32(0))
	assert.Greater(t, cfg.Generation.DefaultHeight, int32(0))
}

func TestLoadConfig_UnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.txt")

	err := os.WriteFile(configPath, []byte("some content"), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "unsupported")
}

func TestConfig_CompleteStructure(t *testing.T) {
	cfg := &Config{
		Version: "1.0",
		Models: ModelConfig{
			DiffusionModelPath: "model.gguf",
			VAEPath:            "vae.safetensors",
			ClipLPath:          "clip_l.safetensors",
			ClipGPath:          "clip_g.safetensors",
		},
		Generation: GenerationConfig{
			DefaultPrompt:         "test prompt",
			DefaultNegativePrompt: "bad quality",
			DefaultWidth:          1024,
			DefaultHeight:         1024,
			DefaultSampleSteps:    30,
			DefaultCfgScale:       7.5,
		},
		Output: OutputConfig{
			OutputDir:     "./outputs",
			NamingPattern: "img_{seed}.png",
		},
	}

	assert.Equal(t, "1.0", cfg.Version)
	assert.Equal(t, "model.gguf", cfg.Models.DiffusionModelPath)
	assert.Equal(t, int32(1024), cfg.Generation.DefaultWidth)
}
