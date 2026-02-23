package stablediffusion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImgGenParams_DefaultValues(t *testing.T) {
	params := &ImgGenParams{
		Prompt: "test prompt",
	}

	assert.Equal(t, "test prompt", params.Prompt)
	assert.Equal(t, int32(0), params.Width)
	assert.Equal(t, int32(0), params.Height)
	assert.Equal(t, float32(0), params.CfgScale)
	assert.Equal(t, int32(0), params.SampleSteps)
}

func TestVidGenParams_DefaultValues(t *testing.T) {
	params := &VidGenParams{
		Prompt: "test video prompt",
	}

	assert.Equal(t, "test video prompt", params.Prompt)
	assert.Equal(t, int32(0), params.Width)
	assert.Equal(t, int32(0), params.Height)
	assert.Equal(t, float32(0), params.CfgScale)
	assert.Equal(t, int32(0), params.VideoFrames)
}

func TestLora_Structure(t *testing.T) {
	lora := &Lora{
		IsHighNoise: true,
		Multiplier:  0.8,
		Path:        "/path/to/lora.safetensors",
	}

	assert.True(t, lora.IsHighNoise)
	assert.Equal(t, float32(0.8), lora.Multiplier)
	assert.Equal(t, "/path/to/lora.safetensors", lora.Path)
}

func TestEmbedding_Structure(t *testing.T) {
	embedding := &Embedding{
		Name: "test_embedding",
		Path: "/path/to/embedding.pt",
	}

	assert.Equal(t, "test_embedding", embedding.Name)
	assert.Equal(t, "/path/to/embedding.pt", embedding.Path)
}

func TestUpscalerParams_DefaultValues(t *testing.T) {
	params := &UpscalerParams{
		EsrganPath: "/path/to/esrgan.pth",
	}

	assert.Equal(t, "/path/to/esrgan.pth", params.EsrganPath)
	assert.Equal(t, 0, params.NThreads)
	assert.Equal(t, 0, params.TileSize)
	assert.False(t, params.OffloadParamsToCPU)
}

// TestDefaultConstants validates that default constants have sensible values
func TestDefaultConstants_ImageGeneration(t *testing.T) {
	// CfgScale should be in reasonable range (1-20)
	assert.InDelta(t, float32(5.0), DefaultCfgScale, 0.01)

	// Sample steps should be positive
	assert.InDelta(t, int32(20), DefaultSampleSteps, 1)

	// Strength should be in valid range [0, 1]
	assert.InDelta(t, float32(0.75), DefaultStrength, 0.01)

	// Seed should be non-negative
	assert.InDelta(t, int64(42), DefaultSeed, 1)

	// Dimensions should be positive and reasonable
	assert.InDelta(t, int32(512), DefaultWidth, 1)
	assert.InDelta(t, int32(512), DefaultHeight, 1)
}

func TestDefaultConstants_VideoGeneration(t *testing.T) {
	// Video frames should be positive
	assert.InDelta(t, int32(33), DefaultVideoFrames, 1)

	// MOE boundary should be in valid range [0, 1]
	assert.InDelta(t, float32(0.875), DefaultMOEBoundary, 0.001)

	// VACE strength should be positive
	assert.InDelta(t, float32(1.0), DefaultVaceStrength, 0.01)

	// High noise cfg scale should be positive
	assert.InDelta(t, float32(6.0), DefaultHighNoiseCfgScale, 0.01)
}

func TestDefaultConstants_Guidance(t *testing.T) {
	// Image cfg scale should be positive
	assert.InDelta(t, float32(1.0), DefaultImageCfgScale, 0.01)

	// Distilled guidance should be positive
	assert.InDelta(t, float32(3.5), DefaultDistilledGuidance, 0.01)

	// SLG start should be in valid range [0, 1]
	assert.InDelta(t, float32(0.01), DefaultSkipLayerStart, 0.001)

	// SLG end should be in valid range [0, 1] and greater than start
	assert.InDelta(t, float32(0.2), DefaultSkipLayerEnd, 0.001)
	assert.Less(t, float32(DefaultSkipLayerStart), float32(DefaultSkipLayerEnd))
}

func TestDefaultConstants_Misc(t *testing.T) {
	// Batch count should be positive
	assert.InDelta(t, int32(1), DefaultBatchCount, 1)

	// Eta should be non-negative
	assert.InDelta(t, float32(1.0), DefaultEta, 0.01)

	// Control strength should be in valid range [0, 1]
	assert.InDelta(t, float32(0.9), DefaultControlStrength, 0.01)

	// Clip skip should be negative (no skip by default)
	assert.InDelta(t, int32(-1), DefaultClipSkip, 1)

	// High noise sample steps: -1 means auto
	assert.InDelta(t, int32(-1), DefaultHighNoiseSampleSteps, 1)
}
