package stablediffusion

import "log"

// ModelSpec represents a Stable Diffusion model specification
type ModelSpec struct {
	Name            string // Model name/identifier
	URL             string // Download URL for the model
	Filename        string // Local filename for the model
	Size            int64  // Model file size in MB
	MinRAM          int64  // Minimum RAM required in GB
	RecommendedRAM  int64  // Recommended RAM in GB
	MinVRAM         int64  // Minimum VRAM required in GB (0 if CPU-only)
	RecommendedVRAM int64  // Recommended VRAM in GB
	ModelType       string // Type of model (SD1.5, SDXL, etc.)
	Description     string // Model description
	Quantization    string // Quantization level (f16, q4_0, q8_0, etc.)
}

// GetAvailableModels returns all available Stable Diffusion models sorted by resource requirements
func GetAvailableModels() []ModelSpec {
	return []ModelSpec{
		{
			Name:            "sd-v1-4-q4_0",
			URL:             "https://huggingface.co/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt",
			Filename:        "sd-v1-4-q4_0.ckpt",
			Size:            2048,
			MinRAM:          4,
			RecommendedRAM:  8,
			MinVRAM:         2,
			RecommendedVRAM: 4,
			ModelType:       "SD1.4",
			Description:     "Stable Diffusion v1.4 quantized (4-bit) - compact and efficient",
			Quantization:    "q4_0",
		},
		{
			Name:            "sd-v1-5-f16",
			URL:             "https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main/v1-5-pruned-emaonly.safetensors",
			Filename:        "sd-v1-5-f16.safetensors",
			Size:            3972,
			MinRAM:          6,
			RecommendedRAM:  12,
			MinVRAM:         4,
			RecommendedVRAM: 6,
			ModelType:       "SD1.5",
			Description:     "Stable Diffusion v1.5 (16-bit) - balanced quality and performance",
			Quantization:    "f16",
		},
		{
			Name:            "sd-v1-5-q8_0",
			URL:             "https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main/v1-5-pruned-emaonly.safetensors",
			Filename:        "sd-v1-5-q8_0.safetensors",
			Size:            2048,
			MinRAM:          4,
			RecommendedRAM:  8,
			MinVRAM:         2,
			RecommendedVRAM: 4,
			ModelType:       "SD1.5",
			Description:     "Stable Diffusion v1.5 quantized (8-bit) - good quality, smaller size",
			Quantization:    "q8_0",
		},
		{
			Name:            "sdxl-base-f16",
			URL:             "https://huggingface.co/stabilityai/stable-diffusion-xl-base-1.0/resolve/main/sd_xl_base_1.0.safetensors",
			Filename:        "sdxl-base-f16.safetensors",
			Size:            6938,
			MinRAM:          12,
			RecommendedRAM:  16,
			MinVRAM:         8,
			RecommendedVRAM: 12,
			ModelType:       "SDXL",
			Description:     "Stable Diffusion XL Base (16-bit) - highest quality, large size",
			Quantization:    "f16",
		},
		{
			Name:            "sdxl-base-q4_0",
			URL:             "https://huggingface.co/stabilityai/stable-diffusion-xl-base-1.0/resolve/main/sd_xl_base_1.0.safetensors",
			Filename:        "sdxl-base-q4_0.safetensors",
			Size:            3500,
			MinRAM:          8,
			RecommendedRAM:  12,
			MinVRAM:         4,
			RecommendedVRAM: 8,
			ModelType:       "SDXL",
			Description:     "Stable Diffusion XL Base quantized (4-bit) - good quality, manageable size",
			Quantization:    "q4_0",
		},
		{
			Name:            "sd-turbo-q8_0",
			URL:             "https://huggingface.co/gpustack/stable-diffusion-v2-1-turbo-GGUF/resolve/main/stable-diffusion-v2-1-turbo-Q8_0.gguf",
			Filename:        "stable-diffusion-v2-1-turbo-Q8_0.gguf",
			Size:            2215,
			MinRAM:          4,
			RecommendedRAM:  6,
			MinVRAM:         2,
			RecommendedVRAM: 4,
			ModelType:       "SD-Turbo",
			Description:     "SD-Turbo v2.1 quantized (8-bit GGUF) - ultra-fast generation, 1-4 steps",
			Quantization:    "q8_0",
		},
	}
}

// HardwareSpecs represents system hardware specifications (compatible with llama package)
type HardwareSpecs struct {
	TotalRAM     int64  // Total RAM in GB
	AvailableRAM int64  // Available RAM in GB
	CPU          string // CPU model
	CPUCores     int    // Number of CPU cores
	GPUMemory    int64  // GPU VRAM in GB (if available)
	GPUModel     string // GPU model
}

// SelectOptimalModel selects the best Stable Diffusion model based on hardware specs
func SelectOptimalModel(specs *HardwareSpecs) ModelSpec {
	models := GetAvailableModels()

	var selectedModel ModelSpec
	for _, model := range models {
		ramOk := model.MinRAM <= specs.AvailableRAM
		vramOk := specs.GPUMemory == 0 || model.MinVRAM <= specs.GPUMemory

		if ramOk && vramOk {
			selectedModel = model
		}
	}

	if selectedModel.Name == "" {
		selectedModel = models[0]
		log.Printf("Warning: System has limited resources (RAM=%dGB, VRAM=%dGB), using smallest model",
			specs.AvailableRAM, specs.GPUMemory)
	}

	log.Printf("Selected Stable Diffusion model: %s (%s) - requires %dGB RAM, %dGB VRAM (system has %dGB RAM, %dGB VRAM available)",
		selectedModel.Name, selectedModel.ModelType, selectedModel.MinRAM, selectedModel.MinVRAM,
		specs.AvailableRAM, specs.GPUMemory)

	return selectedModel
}
