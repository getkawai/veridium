// Model Management Example
// This example demonstrates how to use the model catalog and selection features
package main

import (
	"fmt"

	"github.com/kawai-network/veridium/pkg/stablediffusion/models"
)

func main() {
	fmt.Println("Stable Diffusion Go - Model Management Example")
	fmt.Println("===============================================")

	// Example 1: Get available models from catalog
	fmt.Println("\n1. Getting available models from catalog...")
	availableModels := models.GetAvailableModels()
	fmt.Printf("Total models in catalog: %d\n", len(availableModels))

	// Example 2: List all models with details
	fmt.Println("\n2. Listing all models with details...")
	for i, model := range availableModels {
		fmt.Printf("\n[%d] %s\n", i+1, model.Name)
		fmt.Printf("    Type: %s\n", model.ModelType)
		fmt.Printf("    Size: %d MB\n", model.Size)
		fmt.Printf("    Quantization: %s\n", model.Quantization)
		fmt.Printf("    Min RAM: %d GB (Recommended: %d GB)\n", model.MinRAM, model.RecommendedRAM)
		fmt.Printf("    Min VRAM: %d GB (Recommended: %d GB)\n", model.MinVRAM, model.RecommendedVRAM)
		fmt.Printf("    Description: %s\n", model.Description)
		fmt.Printf("    URL: %s\n", model.URL)
	}

	// Example 3: Select optimal model for high-end system
	fmt.Println("\n3. Selecting optimal model for high-end system...")
	highEndSpecs := &models.HardwareSpecs{
		TotalRAM:     32,
		AvailableRAM: 24,
		CPU:          "AMD Ryzen 9 5950X",
		CPUCores:     16,
		GPUMemory:    12,
		GPUModel:     "NVIDIA RTX 3080",
	}
	selectedModel := models.SelectOptimalModel(highEndSpecs)
	fmt.Printf("Selected: %s\n", selectedModel.Name)
	fmt.Printf("  Type: %s\n", selectedModel.ModelType)
	fmt.Printf("  Size: %d MB\n", selectedModel.Size)
	fmt.Printf("  Quantization: %s\n", selectedModel.Quantization)

	// Example 4: Select optimal model for mid-range system
	fmt.Println("\n4. Selecting optimal model for mid-range system...")
	midRangeSpecs := &models.HardwareSpecs{
		TotalRAM:     16,
		AvailableRAM: 12,
		CPU:          "Intel Core i5-12400",
		CPUCores:     6,
		GPUMemory:    6,
		GPUModel:     "NVIDIA RTX 3060",
	}
	selectedModel = models.SelectOptimalModel(midRangeSpecs)
	fmt.Printf("Selected: %s\n", selectedModel.Name)
	fmt.Printf("  Type: %s\n", selectedModel.ModelType)
	fmt.Printf("  Size: %d MB\n", selectedModel.Size)

	// Example 5: Select optimal model for low-end system
	fmt.Println("\n5. Selecting optimal model for low-end system...")
	lowEndSpecs := &models.HardwareSpecs{
		TotalRAM:     8,
		AvailableRAM: 6,
		CPU:          "Intel Core i3-10100",
		CPUCores:     4,
		GPUMemory:    2,
		GPUModel:     "NVIDIA GTX 1650",
	}
	selectedModel = models.SelectOptimalModel(lowEndSpecs)
	fmt.Printf("Selected: %s\n", selectedModel.Name)
	fmt.Printf("  Type: %s\n", selectedModel.ModelType)
	fmt.Printf("  Size: %d MB\n", selectedModel.Size)

	// Example 6: Select optimal model for CPU-only system
	fmt.Println("\n6. Selecting optimal model for CPU-only system...")
	cpuOnlySpecs := &models.HardwareSpecs{
		TotalRAM:     16,
		AvailableRAM: 12,
		CPU:          "AMD Ryzen 7 5800X",
		CPUCores:     8,
		GPUMemory:    0, // No GPU
		GPUModel:     "",
	}
	selectedModel = models.SelectOptimalModel(cpuOnlySpecs)
	fmt.Printf("Selected: %s\n", selectedModel.Name)
	fmt.Printf("  Type: %s\n", selectedModel.ModelType)
	fmt.Printf("  Size: %d MB\n", selectedModel.Size)

	// Example 7: Filter models by type
	fmt.Println("\n7. Filtering models by type...")
	fmt.Println("\nSD1.5 models:")
	for _, model := range availableModels {
		if model.ModelType == "SD1.5" {
			fmt.Printf("  - %s (%s, %d MB)\n", model.Name, model.Quantization, model.Size)
		}
	}

	fmt.Println("\nSDXL models:")
	for _, model := range availableModels {
		if model.ModelType == "SDXL" {
			fmt.Printf("  - %s (%s, %d MB)\n", model.Name, model.Quantization, model.Size)
		}
	}

	// Example 8: Filter models by quantization
	fmt.Println("\n8. Filtering models by quantization...")
	fmt.Println("\nQ4_0 quantized models (smallest):")
	for _, model := range availableModels {
		if model.Quantization == "q4_0" {
			fmt.Printf("  - %s (%d MB)\n", model.Name, model.Size)
		}
	}

	fmt.Println("\nF16 models (highest quality):")
	for _, model := range availableModels {
		if model.Quantization == "f16" {
			fmt.Printf("  - %s (%d MB)\n", model.Name, model.Size)
		}
	}

	// Example 9: Find models that fit specific constraints
	fmt.Println("\n9. Finding models that fit specific constraints...")
	maxRAM := int64(8)
	maxVRAM := int64(4)
	fmt.Printf("Finding models that fit in %d GB RAM and %d GB VRAM:\n", maxRAM, maxVRAM)
	for _, model := range availableModels {
		if model.MinRAM <= maxRAM && model.MinVRAM <= maxVRAM {
			fmt.Printf("  ✓ %s (RAM: %d GB, VRAM: %d GB)\n", model.Name, model.MinRAM, model.MinVRAM)
		}
	}

	// Example 10: Model download information
	fmt.Println("\n10. Model download information...")
	fmt.Println("To download a model, use the URL from the catalog:")
	exampleModel := availableModels[0]
	fmt.Printf("\nExample: %s\n", exampleModel.Name)
	fmt.Printf("  URL: %s\n", exampleModel.URL)
	fmt.Printf("  Save as: %s\n", exampleModel.Filename)
	fmt.Printf("  Expected size: %d MB\n", exampleModel.Size)

	fmt.Println("\n✅ Model management example completed!")
	fmt.Println("\nNote: This example demonstrates the model catalog and selection features.")
	fmt.Println("For actual model downloading, use the download package or manual download.")
}
