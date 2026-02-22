package models

import (
	"testing"
)

func TestSelectOptimalModel(t *testing.T) {
	tests := []struct {
		name         string
		specs        *HardwareSpecs
		expectedType string
		description  string
	}{
		{
			name: "high_end_system",
			specs: &HardwareSpecs{
				AvailableRAM: 32,
				GPUMemory:    12,
			},
			expectedType: "Qwen-Image-2512",
			description:  "Should select Qwen-Image-2512 for high-end system",
		},
		{
			name: "mid_range_system",
			specs: &HardwareSpecs{
				AvailableRAM: 16,
				GPUMemory:    8,
			},
			expectedType: "Qwen-Image-2512",
			description:  "Should select Qwen-Image-2512 q4 for mid-range system",
		},
		{
			name: "low_end_system",
			specs: &HardwareSpecs{
				AvailableRAM: 8,
				GPUMemory:    4,
			},
			expectedType: "Qwen-Image-2512",
			description:  "Should fallback to first catalog entry when no model fits",
		},
		{
			name: "minimal_system",
			specs: &HardwareSpecs{
				AvailableRAM: 4,
				GPUMemory:    2,
			},
			expectedType: "Qwen-Image-2512",
			description:  "Should fallback to first catalog entry for minimal system",
		},
		{
			name: "cpu_only_system",
			specs: &HardwareSpecs{
				AvailableRAM: 16,
				GPUMemory:    0,
			},
			expectedType: "Qwen-Image-2512",
			description:  "Should select CPU-compatible Qwen-Image variant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := SelectOptimalModel(tt.specs)

			if model.Name == "" {
				t.Errorf("SelectOptimalModel() returned empty model name")
			}

			if model.ModelType != tt.expectedType {
				t.Errorf("SelectOptimalModel() = %v, want %v (%s)",
					model.ModelType, tt.expectedType, tt.description)
			}

			// If model exceeds specs, selector should have fallen back to first model.
			if model.MinRAM > tt.specs.AvailableRAM || (tt.specs.GPUMemory > 0 && model.MinVRAM > tt.specs.GPUMemory) {
				first := GetAvailableModels()[0]
				if model.Name != first.Name {
					t.Errorf("SelectOptimalModel() expected fallback %s, got %s", first.Name, model.Name)
				}
			}
		})
	}
}

func TestGetAvailableModels(t *testing.T) {
	models := GetAvailableModels()

	if len(models) == 0 {
		t.Fatal("GetAvailableModels() returned no models")
	}

	// Verify all models have required fields
	for _, model := range models {
		if model.Name == "" {
			t.Error("Model has empty name")
		}
		if model.URL == "" {
			t.Error("Model has empty URL")
		}
		if model.Filename == "" {
			t.Error("Model has empty filename")
		}
		if model.Size <= 0 {
			t.Errorf("Model %s has invalid size: %d", model.Name, model.Size)
		}
		if model.MinRAM <= 0 {
			t.Errorf("Model %s has invalid MinRAM: %d", model.Name, model.MinRAM)
		}
		if model.ModelType == "" {
			t.Errorf("Model %s has empty ModelType", model.Name)
		}
	}

	// Verify models are sorted by resource requirements (smallest first)
	for i := 1; i < len(models); i++ {
		prev := models[i-1]
		curr := models[i]

		// Check if models are generally sorted by requirements
		// (not strict, but smallest should be first)
		if i == 1 && prev.MinRAM > curr.MinRAM {
			t.Errorf("Models not sorted by requirements: %s (%dGB) before %s (%dGB)",
				prev.Name, prev.MinRAM, curr.Name, curr.MinRAM)
		}
	}
}

func TestModelSpecFields(t *testing.T) {
	models := GetAvailableModels()

	// Test specific known model
	var qwenQ4 *ModelSpec
	for i := range models {
		if models[i].Name == "qwen-image-2512-q4_k_m" {
			qwenQ4 = &models[i]
			break
		}
	}

	if qwenQ4 == nil {
		t.Fatal("Qwen-Image q4 model not found in catalog")
	}

	// Verify Qwen-Image specs
	if qwenQ4.ModelType != "Qwen-Image-2512" {
		t.Errorf("Qwen-Image ModelType = %v, want Qwen-Image-2512", qwenQ4.ModelType)
	}

	if qwenQ4.Quantization != "q4_k_m" {
		t.Errorf("Qwen-Image Quantization = %v, want q4_k_m", qwenQ4.Quantization)
	}

	if qwenQ4.LLMURL == "" || qwenQ4.LLMFilename == "" {
		t.Error("Qwen-Image model missing LLM component")
	}

	if qwenQ4.VAEURL == "" || qwenQ4.VAEFilename == "" {
		t.Error("Qwen-Image model missing VAE component")
	}

	if qwenQ4.EditModelURL == "" || qwenQ4.EditModelFile == "" {
		t.Error("Qwen-Image model missing edit component")
	}
}
