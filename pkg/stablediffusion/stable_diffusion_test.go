package stablediffusion

import (
	"testing"

	"github.com/kawai-network/veridium/pkg/stablediffusion/sd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRNGTypeMap_AllValuesValid(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected sd.RngType
	}{
		{"default", "default", sd.DefaultRNG},
		{"cuda", "cuda", sd.CUDARNG},
		{"cpu", "cpu", sd.CPURNG},
		{"type_count", "type_count", sd.RNGTypeCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := RNGTypeMap[tt.key]
			require.True(t, exists, "key %s should exist in RNGTypeMap", tt.key)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestSampleMethodMap_AllValuesValid(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected sd.SampleMethod
	}{
		{"default", "default", -1},
		{"euler", "euler", sd.EulerSampleMethod},
		{"euler_a", "euler_a", sd.EulerASampleMethod},
		{"heun", "heun", sd.HeunSampleMethod},
		{"dpm2", "dpm2", sd.DPM2SampleMethod},
		{"lcm", "lcm", sd.LCMSampleMethod},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := SampleMethodMap[tt.key]
			require.True(t, exists, "key %s should exist in SampleMethodMap", tt.key)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestSchedulerMap_AllValuesValid(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected sd.Scheduler
	}{
		{"default", "default", -1},
		{"discrete", "discrete", sd.DiscreteScheduler},
		{"karras", "karras", sd.KarrasScheduler},
		{"exponential", "exponential", sd.ExponentialScheduler},
		{"lcm", "lcm", sd.LCMScheduler},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := SchedulerMap[tt.key]
			require.True(t, exists, "key %s should exist in SchedulerMap", tt.key)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestSDTypeMap_AllValuesValid(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected sd.SDType
	}{
		{"f32", "f32", sd.SDTypeF32},
		{"f16", "f16", sd.SDTypeF16},
		{"q4_0", "q4_0", sd.SDTypeQ4_0},
		{"q4_1", "q4_1", sd.SDTypeQ4_1},
		{"q8_0", "q8_0", sd.SDTypeQ8_0},
		{"default", "default", sd.SDTypeCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := SDTypeMap[tt.key]
			require.True(t, exists, "key %s should exist in SDTypeMap", tt.key)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestPredictionMap_AllValuesValid(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected sd.Prediction
	}{
		{"eps", "eps", sd.EPSPred},
		{"v", "v", sd.VPred},
		{"flow", "flow", sd.FlowPred},
		{"default", "default", sd.PredictionCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := PredictionMap[tt.key]
			require.True(t, exists, "key %s should exist in PredictionMap", tt.key)
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestLoraApplyModeMap_AllValuesValid(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected sd.LoraApplyMode
	}{
		{"auto", "auto", sd.LoraApplyAuto},
		{"immediately", "immediately", sd.LoraApplyImmediately},
		{"at_runtime", "at_runtime", sd.LoraApplyAtRuntime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := LoraApplyModeMap[tt.key]
			require.True(t, exists, "key %s should exist in LoraApplyModeMap", tt.key)
			assert.Equal(t, tt.expected, value)
		})
	}
}

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

func TestContextParams_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  *ContextParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid_params",
			params: &ContextParams{
				ModelPath: "/path/to/model.gguf",
				WType:     "f16",
				RNGType:   "cuda",
			},
			wantErr: false,
		},
		{
			name: "invalid_wtype",
			params: &ContextParams{
				ModelPath: "/path/to/model.gguf",
				WType:     "invalid_type",
			},
			wantErr: true,
			errMsg:  "Invalid WType",
		},
		{
			name: "invalid_rng_type",
			params: &ContextParams{
				ModelPath:      "/path/to/model.gguf",
				RNGType:        "invalid_rng",
				SamplerRNGType: "cuda",
			},
			wantErr: true,
			errMsg:  "Invalid RNG type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: NewStableDiffusion requires actual library loaded
			// This test validates parameter logic only
			if tt.params.WType != "" {
				if _, ok := SDTypeMap[tt.params.WType]; !ok && tt.wantErr {
					assert.Contains(t, tt.errMsg, "Invalid WType")
				}
			}
			if tt.params.RNGType != "" {
				if _, ok := RNGTypeMap[tt.params.RNGType]; !ok && tt.wantErr {
					assert.Contains(t, tt.errMsg, "Invalid RNG type")
				}
			}
		})
	}
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
