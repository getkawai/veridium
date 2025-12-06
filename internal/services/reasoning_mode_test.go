package services

import (
	"testing"

	"github.com/kawai-network/veridium/internal/llama"
)

func TestHardwareValidation(t *testing.T) {
	tests := []struct {
		name          string
		mode          ReasoningMode
		specs         *llama.HardwareSpecs
		expectValid   bool
		expectContain string
	}{
		{
			name: "Disabled mode with low-end hardware - should pass",
			mode: ReasoningDisabled,
			specs: &llama.HardwareSpecs{
				AvailableRAM: 4,
				CPUCores:     2,
				GPUMemory:    0,
			},
			expectValid: true,
		},
		{
			name: "Enabled mode with sufficient hardware - should pass",
			mode: ReasoningEnabled,
			specs: &llama.HardwareSpecs{
				AvailableRAM: 8,
				CPUCores:     4,
				GPUMemory:    0,
			},
			expectValid: true,
		},
		{
			name: "Enabled mode with insufficient RAM - should fail",
			mode: ReasoningEnabled,
			specs: &llama.HardwareSpecs{
				AvailableRAM: 4,
				CPUCores:     4,
				GPUMemory:    0,
			},
			expectValid:   false,
			expectContain: "Insufficient RAM",
		},
		{
			name: "Enabled mode with insufficient CPU - should fail",
			mode: ReasoningEnabled,
			specs: &llama.HardwareSpecs{
				AvailableRAM: 8,
				CPUCores:     2,
				GPUMemory:    0,
			},
			expectValid:   false,
			expectContain: "Insufficient CPU cores",
		},
		{
			name: "Verbose mode with high-end hardware - should pass",
			mode: ReasoningVerbose,
			specs: &llama.HardwareSpecs{
				AvailableRAM: 16,
				CPUCores:     8,
				GPUMemory:    8,
				GPUModel:     "NVIDIA RTX 4090",
			},
			expectValid: true,
		},
		{
			name: "Verbose mode with insufficient RAM - should fail",
			mode: ReasoningVerbose,
			specs: &llama.HardwareSpecs{
				AvailableRAM: 8,
				CPUCores:     8,
				GPUMemory:    0,
			},
			expectValid:   false,
			expectContain: "Insufficient RAM",
		},
		{
			name:        "Nil specs - should allow (with warning)",
			mode:        ReasoningVerbose,
			specs:       nil,
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ReasoningConfig{Mode: tt.mode}
			valid, reason := config.ValidateHardware(tt.specs)

			if valid != tt.expectValid {
				t.Errorf("ValidateHardware() valid = %v, want %v", valid, tt.expectValid)
			}

			if !tt.expectValid && tt.expectContain != "" {
				if reason == "" {
					t.Errorf("Expected reason to contain '%s', but got empty reason", tt.expectContain)
				}
			}
		})
	}
}

func TestSuggestModeForHardware(t *testing.T) {
	tests := []struct {
		name         string
		specs        *llama.HardwareSpecs
		expectedMode ReasoningMode
	}{
		{
			name: "High-end system - suggest verbose",
			specs: &llama.HardwareSpecs{
				AvailableRAM: 32,
				CPUCores:     16,
				GPUMemory:    16,
			},
			expectedMode: ReasoningVerbose,
		},
		{
			name: "Mid-range system - suggest enabled",
			specs: &llama.HardwareSpecs{
				AvailableRAM: 8,
				CPUCores:     4,
				GPUMemory:    0,
			},
			expectedMode: ReasoningEnabled,
		},
		{
			name: "Low-end system - suggest disabled",
			specs: &llama.HardwareSpecs{
				AvailableRAM: 4,
				CPUCores:     2,
				GPUMemory:    0,
			},
			expectedMode: ReasoningDisabled,
		},
		{
			name: "Just enough for enabled",
			specs: &llama.HardwareSpecs{
				AvailableRAM: 8,
				CPUCores:     4,
				GPUMemory:    0,
			},
			expectedMode: ReasoningEnabled,
		},
		{
			name: "Just enough for verbose",
			specs: &llama.HardwareSpecs{
				AvailableRAM: 16,
				CPUCores:     6,
				GPUMemory:    0,
			},
			expectedMode: ReasoningVerbose,
		},
		{
			name:         "Nil specs - suggest disabled (safest)",
			specs:        nil,
			expectedMode: ReasoningDisabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := SuggestModeForHardware(tt.specs)
			if mode != tt.expectedMode {
				t.Errorf("SuggestModeForHardware() = %v, want %v", mode, tt.expectedMode)
			}
		})
	}
}

func TestGetHardwareRequirements(t *testing.T) {
	tests := []struct {
		name           string
		mode           ReasoningMode
		expectedMinRAM int64
		expectedCPU    int
	}{
		{
			name:           "Disabled mode requirements",
			mode:           ReasoningDisabled,
			expectedMinRAM: 4,
			expectedCPU:    2,
		},
		{
			name:           "Enabled mode requirements",
			mode:           ReasoningEnabled,
			expectedMinRAM: 8,
			expectedCPU:    4,
		},
		{
			name:           "Verbose mode requirements",
			mode:           ReasoningVerbose,
			expectedMinRAM: 16,
			expectedCPU:    6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ReasoningConfig{Mode: tt.mode}
			req := config.GetHardwareRequirements()

			if req.MinRAM != tt.expectedMinRAM {
				t.Errorf("MinRAM = %d, want %d", req.MinRAM, tt.expectedMinRAM)
			}
			if req.MinCPUCores != tt.expectedCPU {
				t.Errorf("MinCPUCores = %d, want %d", req.MinCPUCores, tt.expectedCPU)
			}
			if req.Description == "" {
				t.Error("Description should not be empty")
			}
		})
	}
}
