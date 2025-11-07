package llama

import "log"

// QwenModelSpec represents a Qwen3 model specification
type QwenModelSpec struct {
	Name            string // Model name for Ollama
	Parameters      string // Parameter size (0.6b, 1.7b, etc.)
	MinRAM          int64  // Minimum RAM required in GB
	RecommendedRAM  int64  // Recommended RAM in GB
	MinVRAM         int64  // Minimum VRAM required in GB (0 if CPU-only)
	RecommendedVRAM int64  // Recommended VRAM in GB
	Description     string // Model description
}

// GetAvailableQwenModels returns all available Qwen3 models sorted by resource requirements
func GetAvailableQwenModels() []QwenModelSpec {
	return []QwenModelSpec{
		{
			Name:            "qwen3:0.6b",
			Parameters:      "0.6b",
			MinRAM:          2,
			RecommendedRAM:  4,
			MinVRAM:         0,
			RecommendedVRAM: 2,
			Description:     "Smallest Qwen3 model, suitable for basic tasks on low-end hardware",
		},
		{
			Name:            "qwen3:1.7b",
			Parameters:      "1.7b",
			MinRAM:          4,
			RecommendedRAM:  8,
			MinVRAM:         0,
			RecommendedVRAM: 4,
			Description:     "Lightweight Qwen3 model for general conversation and simple tasks",
		},
		{
			Name:            "qwen3:4b",
			Parameters:      "4b",
			MinRAM:          6,
			RecommendedRAM:  12,
			MinVRAM:         0,
			RecommendedVRAM: 6,
			Description:     "Balanced Qwen3 model for most general-purpose tasks",
		},
		{
			Name:            "qwen3:8b",
			Parameters:      "8b",
			MinRAM:          10,
			RecommendedRAM:  16,
			MinVRAM:         0,
			RecommendedVRAM: 8,
			Description:     "High-quality Qwen3 model for advanced conversations and reasoning",
		},
		{
			Name:            "qwen3:14b",
			Parameters:      "14b",
			MinRAM:          18,
			RecommendedRAM:  32,
			MinVRAM:         0,
			RecommendedVRAM: 16,
			Description:     "Large Qwen3 model for complex tasks and professional use",
		},
		{
			Name:            "qwen3:30b",
			Parameters:      "30b",
			MinRAM:          40,
			RecommendedRAM:  64,
			MinVRAM:         0,
			RecommendedVRAM: 32,
			Description:     "Very large Qwen3 MoE model for demanding applications",
		},
		{
			Name:            "qwen3:32b",
			Parameters:      "32b",
			MinRAM:          42,
			RecommendedRAM:  64,
			MinVRAM:         0,
			RecommendedVRAM: 32,
			Description:     "Large Qwen3 dense model for high-performance tasks",
		},
		{
			Name:            "qwen3:235b",
			Parameters:      "235b",
			MinRAM:          120,
			RecommendedRAM:  256,
			MinVRAM:         0,
			RecommendedVRAM: 80,
			Description:     "Flagship Qwen3 MoE model (235B total, ~22B active) for extreme performance",
		},
	}
}

// SelectOptimalQwenModel selects the best Qwen model based on hardware specs
func SelectOptimalQwenModel(specs *HardwareSpecs) QwenModelSpec {
	models := GetAvailableQwenModels()

	var selectedModel QwenModelSpec
	for _, model := range models {
		if model.MinRAM <= specs.AvailableRAM {
			selectedModel = model
		} else {
			break
		}
	}

	if selectedModel.Name == "" {
		selectedModel = models[0]
		log.Printf("Warning: System has very low RAM (%dGB), using smallest model", specs.AvailableRAM)
	}

	log.Printf("Selected Qwen3 model: %s (%s) - requires %dGB RAM (system has %dGB available)",
		selectedModel.Name, selectedModel.Parameters, selectedModel.MinRAM, specs.AvailableRAM)

	return selectedModel
}

