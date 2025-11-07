package llama

import (
	"encoding/json"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// HardwareSpecs represents system hardware specifications
type HardwareSpecs struct {
	TotalRAM     int64  // Total RAM in GB
	AvailableRAM int64  // Available RAM in GB
	CPU          string // CPU model
	CPUCores     int    // Number of CPU cores
	GPUMemory    int64  // GPU VRAM in GB (if available)
	GPUModel     string // GPU model
}

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

// StableDiffusionModelSpec represents a Stable Diffusion model specification
type StableDiffusionModelSpec struct {
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

// DetectHardwareSpecs detects the current system's hardware specifications
func DetectHardwareSpecs() *HardwareSpecs {
	specs := &HardwareSpecs{}

	// Detect CPU cores
	specs.CPUCores = runtime.NumCPU()

	// Detect system memory based on OS
	if runtime.GOOS == "darwin" {
		specs.detectMacOSSpecs()
	} else if runtime.GOOS == "linux" {
		specs.detectLinuxSpecs()
	} else if runtime.GOOS == "windows" {
		specs.detectWindowsSpecs()
	} else {
		log.Printf("Unsupported OS for hardware detection: %s", runtime.GOOS)
		// Fallback values
		specs.TotalRAM = 8
		specs.AvailableRAM = 6
	}

	log.Printf("Detected hardware: RAM=%dGB (available=%dGB), CPU cores=%d, GPU=%s (VRAM=%dGB)",
		specs.TotalRAM, specs.AvailableRAM, specs.CPUCores, specs.GPUModel, specs.GPUMemory)

	return specs
}

// detectMacOSSpecs detects hardware specs on macOS
func (specs *HardwareSpecs) detectMacOSSpecs() {
	// Get total memory
	if out, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
		if memBytes, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64); err == nil {
			specs.TotalRAM = memBytes / (1024 * 1024 * 1024) // Convert to GB
		}
	}

	// Estimate available memory (roughly 80% of total, accounting for OS usage)
	specs.AvailableRAM = int64(float64(specs.TotalRAM) * 0.8)

	// Get CPU model
	if out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
		specs.CPU = strings.TrimSpace(string(out))
	}

	// Try to detect GPU (basic detection)
	if out, err := exec.Command("system_profiler", "SPDisplaysDataType", "-json").Output(); err == nil {
		specs.parseGPUFromSystemProfiler(string(out))
	}
}

// detectLinuxSpecs detects hardware specs on Linux
func (specs *HardwareSpecs) detectLinuxSpecs() {
	// Get memory info from /proc/meminfo
	if out, err := exec.Command("cat", "/proc/meminfo").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					if memKB, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
						specs.TotalRAM = memKB / (1024 * 1024) // Convert KB to GB
					}
				}
			} else if strings.HasPrefix(line, "MemAvailable:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					if memKB, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
						specs.AvailableRAM = memKB / (1024 * 1024) // Convert KB to GB
					}
				}
			}
		}
	}

	// If MemAvailable not found, estimate as 80% of total
	if specs.AvailableRAM == 0 {
		specs.AvailableRAM = int64(float64(specs.TotalRAM) * 0.8)
	}

	// Get CPU info
	if out, err := exec.Command("cat", "/proc/cpuinfo").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					specs.CPU = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}

	// Try to detect NVIDIA GPU
	if out, err := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader,nounits").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) > 0 {
			parts := strings.Split(lines[0], ",")
			if len(parts) >= 2 {
				specs.GPUModel = strings.TrimSpace(parts[0])
				if vramMB, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64); err == nil {
					specs.GPUMemory = vramMB / 1024 // Convert MB to GB
				}
			}
		}
	}
}

// detectWindowsSpecs detects hardware specs on Windows
func (specs *HardwareSpecs) detectWindowsSpecs() {
	// Get total memory using wmic
	if out, err := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/value").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "TotalPhysicalMemory=") {
				memStr := strings.TrimPrefix(line, "TotalPhysicalMemory=")
				memStr = strings.TrimSpace(memStr)
				if memBytes, err := strconv.ParseInt(memStr, 10, 64); err == nil {
					specs.TotalRAM = memBytes / (1024 * 1024 * 1024) // Convert to GB
				}
			}
		}
	}

	// Estimate available memory
	specs.AvailableRAM = int64(float64(specs.TotalRAM) * 0.8)

	// Get CPU info
	if out, err := exec.Command("wmic", "cpu", "get", "name", "/value").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Name=") {
				specs.CPU = strings.TrimPrefix(line, "Name=")
				specs.CPU = strings.TrimSpace(specs.CPU)
				break
			}
		}
	}
}

// parseGPUFromSystemProfiler parses GPU info from macOS system_profiler JSON output
func (specs *HardwareSpecs) parseGPUFromSystemProfiler(jsonStr string) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return
	}

	if displays, ok := data["SPDisplaysDataType"].([]interface{}); ok && len(displays) > 0 {
		if display, ok := displays[0].(map[string]interface{}); ok {
			if name, ok := display["sppci_model"].(string); ok {
				specs.GPUModel = name
			}
			// Note: macOS system_profiler doesn't easily provide VRAM info for all GPUs
			// For Apple Silicon, we can make educated guesses based on model
			if strings.Contains(specs.GPUModel, "Apple") {
				// Apple Silicon GPUs share system memory
				specs.GPUMemory = specs.TotalRAM / 4 // Conservative estimate
			}
		}
	}
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

// GetAvailableStableDiffusionModels returns all available Stable Diffusion models sorted by resource requirements
func GetAvailableStableDiffusionModels() []StableDiffusionModelSpec {
	return []StableDiffusionModelSpec{
		{
			Name:            "sd-v1-4-q4_0",
			URL:             "https://huggingface.co/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt",
			Filename:        "sd-v1-4-q4_0.ckpt",
			Size:            2048, // ~2GB
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
			Size:            3972, // ~4GB
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
			Size:            2048, // ~2GB
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
			Size:            6938, // ~7GB
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
			Size:            3500, // ~3.5GB
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
			URL:             "https://huggingface.co/stabilityai/sd-turbo/resolve/main/sd_turbo.safetensors",
			Filename:        "sd-turbo-q8_0.safetensors",
			Size:            2048, // ~2GB
			MinRAM:          4,
			RecommendedRAM:  6,
			MinVRAM:         2,
			RecommendedVRAM: 4,
			ModelType:       "SD-Turbo",
			Description:     "SD-Turbo quantized (8-bit) - ultra-fast generation, 1-4 steps",
			Quantization:    "q8_0",
		},
	}
}

// SelectOptimalQwenModel selects the best Qwen model based on hardware specs
func SelectOptimalQwenModel(specs *HardwareSpecs) QwenModelSpec {
	models := GetAvailableQwenModels()

	// Find the largest model that fits within the hardware constraints
	var selectedModel QwenModelSpec
	for _, model := range models {
		// Check if the model fits within RAM constraints
		if model.MinRAM <= specs.AvailableRAM {
			selectedModel = model
		} else {
			break // Models are sorted by size, so we can stop here
		}
	}

	// If no model was selected (very low RAM), default to the smallest
	if selectedModel.Name == "" {
		selectedModel = models[0]
		log.Printf("Warning: System has very low RAM (%dGB), using smallest model", specs.AvailableRAM)
	}

	log.Printf("Selected Qwen3 model: %s (%s) - requires %dGB RAM (system has %dGB available)",
		selectedModel.Name, selectedModel.Parameters, selectedModel.MinRAM, specs.AvailableRAM)

	return selectedModel
}

// SelectOptimalStableDiffusionModel selects the best Stable Diffusion model based on hardware specs
func SelectOptimalStableDiffusionModel(specs *HardwareSpecs) StableDiffusionModelSpec {
	models := GetAvailableStableDiffusionModels()

	// Find the largest model that fits within the hardware constraints
	var selectedModel StableDiffusionModelSpec
	for _, model := range models {
		// Check if the model fits within RAM and VRAM constraints
		ramOk := model.MinRAM <= specs.AvailableRAM
		vramOk := specs.GPUMemory == 0 || model.MinVRAM <= specs.GPUMemory

		if ramOk && vramOk {
			selectedModel = model
		}
	}

	// If no model was selected (very low RAM/VRAM), default to the smallest
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
