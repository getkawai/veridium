package models

import (
	"testing"
)

func TestDetectModelType(t *testing.T) {
	tests := []struct {
		name     string
		modelID  string
		files    []string
		expected ModelType
	}{
		{
			name:     "SD v1.5 safetensors without sd keyword",
			modelID:  "v1-5-pruned",
			files:    []string{"/path/to/v1-5-pruned.safetensors"},
			expected: ModelTypeDiffusion,
		},
		{
			name:     "RealVisXL safetensors",
			modelID:  "RealVisXL",
			files:    []string{"/path/to/RealVisXL.safetensors"},
			expected: ModelTypeDiffusion,
		},
		{
			name:     "SD with explicit sd keyword",
			modelID:  "sd-v1-5",
			files:    []string{"/path/to/sd-v1-5.safetensors"},
			expected: ModelTypeDiffusion,
		},
		{
			name:     "SDXL model",
			modelID:  "sdxl-base",
			files:    []string{"/path/to/sdxl-base.safetensors"},
			expected: ModelTypeDiffusion,
		},
		{
			name:     "Flux model",
			modelID:  "flux-dev",
			files:    []string{"/path/to/flux-dev.safetensors"},
			expected: ModelTypeDiffusion,
		},
		{
			name:     "CKPT file is diffusion",
			modelID:  "some-model",
			files:    []string{"/path/to/model.ckpt"},
			expected: ModelTypeDiffusion,
		},
		{
			name:     "LLM with GGUF",
			modelID:  "llama-3",
			files:    []string{"/path/to/llama-3.gguf"},
			expected: ModelTypeLLM,
		},
		{
			name:     "LLM safetensors with llama keyword",
			modelID:  "llama-2-7b",
			files:    []string{"/path/to/llama-2-7b.safetensors"},
			expected: ModelTypeLLM,
		},
		{
			name:     "LLM safetensors with qwen keyword",
			modelID:  "qwen-7b",
			files:    []string{"/path/to/qwen-7b.safetensors"},
			expected: ModelTypeLLM,
		},
		{
			name:     "Whisper model",
			modelID:  "whisper-large",
			files:    []string{"/path/to/ggml-large.bin"},
			expected: ModelTypeAudio,
		},
		{
			name:     "Generic safetensors defaults to diffusion",
			modelID:  "custom-model",
			files:    []string{"/path/to/custom-model.safetensors"},
			expected: ModelTypeDiffusion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectModelType(tt.modelID, tt.files)
			if result != tt.expected {
				t.Errorf("detectModelType(%q, %v) = %v, want %v", tt.modelID, tt.files, result, tt.expected)
			}
		})
	}
}
