package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToContextParams_BasicConversion(t *testing.T) {
	cfg := &Config{
		Models: ModelConfig{
			DiffusionModelPath: "models/diffusion.gguf",
			VAEPath:            "models/vae.safetensors",
			ClipLPath:          "models/clip_l.safetensors",
		},
		Performance: PerformanceConfig{
			NThreads:           8,
			DiffusionFlashAttn: true,
			OffloadParamsToCPU: true,
		},
	}

	params, err := cfg.ToContextParams()
	require.NoError(t, err)
	assert.NotNil(t, params)
	assert.Equal(t, "models/diffusion.gguf", params.DiffusionModelPath)
	assert.Equal(t, "models/vae.safetensors", params.VAEPath)
	assert.Equal(t, int32(8), params.NThreads)
	assert.True(t, params.DiffusionFlashAttn)
}

func TestToContextParams_EmptyConfig(t *testing.T) {
	cfg := &Config{}

	params, err := cfg.ToContextParams()
	// Empty config should fail validation because diffusion model path is required
	require.Error(t, err)
	assert.Nil(t, params)
	assert.Contains(t, err.Error(), "diffusion model path must be specified")
}

func TestToImgGenParams_BasicConversion(t *testing.T) {
	cfg := &Config{
		Generation: GenerationConfig{
			DefaultPrompt:         "a beautiful landscape",
			DefaultNegativePrompt: "blurry, low quality",
			DefaultWidth:          512,
			DefaultHeight:         512,
			DefaultSampleSteps:    20,
			DefaultCfgScale:       7.5,
		},
	}

	params := cfg.ToImgGenParams(nil)
	assert.NotNil(t, params)
	assert.Equal(t, "a beautiful landscape", params.Prompt)
	assert.Equal(t, "blurry, low quality", params.NegativePrompt)
	assert.Equal(t, int32(512), params.Width)
	assert.Equal(t, int32(512), params.Height)
	assert.Equal(t, int32(20), params.SampleSteps)
	assert.Equal(t, float32(7.5), params.CfgScale)
}

func TestToImgGenParams_WithOverrides(t *testing.T) {
	cfg := &Config{
		Generation: GenerationConfig{
			DefaultWidth:  512,
			DefaultHeight: 512,
		},
	}

	params := cfg.ToImgGenParams(nil)

	// Override values
	params.Width = 1024
	params.Height = 1024
	params.Prompt = "custom prompt"

	assert.Equal(t, int32(1024), params.Width)
	assert.Equal(t, int32(1024), params.Height)
	assert.Equal(t, "custom prompt", params.Prompt)
}

func TestToVidGenParams_BasicConversion(t *testing.T) {
	cfg := &Config{
		Generation: GenerationConfig{
			DefaultPrompt:      "a video of nature",
			DefaultWidth:       512,
			DefaultHeight:      512,
			DefaultSampleSteps: 30,
		},
	}

	params := cfg.ToVidGenParams(nil)
	assert.NotNil(t, params)
	assert.Equal(t, "a video of nature", params.Prompt)
	assert.Equal(t, int32(512), params.Width)
	assert.Equal(t, int32(512), params.Height)
}

func TestConfig_ToContextParams_AllFields(t *testing.T) {
	cfg := &Config{
		Models: ModelConfig{
			DiffusionModelPath: "diffusion.gguf",
			VAEPath:            "vae.safetensors",
			ClipLPath:          "clip_l.safetensors",
			ClipGPath:          "clip_g.safetensors",
			T5XXLPath:          "t5xxl.gguf",
		},
		Performance: PerformanceConfig{
			NThreads:           16,
			DiffusionFlashAttn: true,
			OffloadParamsToCPU: true,
			KeepClipOnCPU:      true,
			KeepVAEOnCPU:       true,
			EnableMmap:         true,
		},
	}

	params, err := cfg.ToContextParams()
	require.NoError(t, err)
	assert.Equal(t, "diffusion.gguf", params.DiffusionModelPath)
	assert.Equal(t, "vae.safetensors", params.VAEPath)
	assert.Equal(t, "clip_l.safetensors", params.ClipLPath)
	assert.Equal(t, "clip_g.safetensors", params.ClipGPath)
	assert.Equal(t, "t5xxl.gguf", params.T5XXLPath)
	assert.Equal(t, int32(16), params.NThreads)
	assert.True(t, params.DiffusionFlashAttn)
	assert.True(t, params.OffloadParamsToCPU)
	assert.True(t, params.KeepClipOnCPU)
	assert.True(t, params.KeepVAEOnCPU)
	assert.True(t, params.EnableMmap)
}

func TestGenerationConfig_DefaultValues(t *testing.T) {
	cfg := GenerationConfig{}

	assert.Equal(t, int32(0), cfg.DefaultWidth)
	assert.Equal(t, int32(0), cfg.DefaultHeight)
	assert.Equal(t, int32(0), cfg.DefaultSampleSteps)
	assert.Equal(t, float32(0), cfg.DefaultCfgScale)
	assert.Empty(t, cfg.DefaultPrompt)
}

func TestModelConfig_AllPaths(t *testing.T) {
	models := ModelConfig{
		DiffusionModelPath: "path/to/diffusion.gguf",
		VAEPath:            "path/to/vae.safetensors",
		ClipLPath:          "path/to/clip_l.safetensors",
		ClipGPath:          "path/to/clip_g.safetensors",
		T5XXLPath:          "path/to/t5xxl.gguf",
		ControlNetPath:     "path/to/controlnet.safetensors",
	}

	assert.NotEmpty(t, models.DiffusionModelPath)
	assert.NotEmpty(t, models.VAEPath)
	assert.NotEmpty(t, models.ClipLPath)
	assert.NotEmpty(t, models.ClipGPath)
	assert.NotEmpty(t, models.T5XXLPath)
	assert.NotEmpty(t, models.ControlNetPath)
}

func TestPerformanceConfig_Flags(t *testing.T) {
	perf := PerformanceConfig{
		NThreads:           8,
		DiffusionFlashAttn: true,
		OffloadParamsToCPU: true,
		KeepClipOnCPU:      false,
		KeepVAEOnCPU:       false,
		EnableMmap:         true,
	}

	assert.Equal(t, int32(8), perf.NThreads)
	assert.True(t, perf.DiffusionFlashAttn)
	assert.True(t, perf.OffloadParamsToCPU)
	assert.False(t, perf.KeepClipOnCPU)
	assert.False(t, perf.KeepVAEOnCPU)
	assert.True(t, perf.EnableMmap)
}

func TestOutputConfig_Structure(t *testing.T) {
	output := OutputConfig{
		OutputDir:     "./outputs",
		NamingPattern: "img_{timestamp}_{seed}.png",
	}

	assert.Equal(t, "./outputs", output.OutputDir)
	assert.Equal(t, "img_{timestamp}_{seed}.png", output.NamingPattern)
}
