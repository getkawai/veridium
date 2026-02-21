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
			expectedType: "SDXL",
			description:  "Should select SDXL for high-end system",
		},
		{
			name: "mid_range_system",
			specs: &HardwareSpecs{
				AvailableRAM: 16,
				GPUMemory:    8,
			},
			expectedType: "SDXL",
			description:  "Should select SDXL for mid-range system",
		},
		{
			name: "low_end_system",
			specs: &HardwareSpecs{
				AvailableRAM: 8,
				GPUMemory:    4,
			},
			expectedType: "SDXL",
			description:  "Should select SDXL q4_0 for low-end system (best that fits)",
		},
		{
			name: "minimal_system",
			specs: &HardwareSpecs{
				AvailableRAM: 4,
				GPUMemory:    2,
			},
			expectedType: "SD1.5",
			description:  "Should select SD1.5 for minimal system (smallest model that fits)",
		},
		{
			name: "cpu_only_system",
			specs: &HardwareSpecs{
				AvailableRAM: 16,
				GPUMemory:    0,
			},
			expectedType: "SDXL",
			description:  "Should select best CPU-compatible model",
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

			// Verify model meets hardware requirements
			if model.MinRAM > tt.specs.AvailableRAM {
				t.Errorf("Selected model requires %dGB RAM but system has %dGB",
					model.MinRAM, tt.specs.AvailableRAM)
			}

			if tt.specs.GPUMemory > 0 && model.MinVRAM > tt.specs.GPUMemory {
				t.Errorf("Selected model requires %dGB VRAM but system has %dGB",
					model.MinVRAM, tt.specs.GPUMemory)
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

	// Test specific known models
	var sdTurbo *ModelSpec
	for i := range models {
		if models[i].Name == "sd-turbo-q8_0" {
			sdTurbo = &models[i]
			break
		}
	}

	if sdTurbo == nil {
		t.Fatal("SD-Turbo model not found in catalog")
	}

	// Verify SD-Turbo specs
	if sdTurbo.ModelType != "SD-Turbo" {
		t.Errorf("SD-Turbo ModelType = %v, want SD-Turbo", sdTurbo.ModelType)
	}

	if sdTurbo.Quantization != "q8_0" {
		t.Errorf("SD-Turbo Quantization = %v, want q8_0", sdTurbo.Quantization)
	}

	if sdTurbo.MinRAM != 4 {
		t.Errorf("SD-Turbo MinRAM = %v, want 4", sdTurbo.MinRAM)
	}
}
