// Model Management Example
// This example demonstrates how to use the model management feature
package main

import (
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/stablediffusion/models"
)

func main() {
	fmt.Println("Stable Diffusion Go - Model Management Example")
	fmt.Println("===============================================")

	// Example 1: Create and load model registry
	fmt.Println("\n1. Creating and loading model registry...")
	registry := models.NewRegistry("./model_registry.json")
	if err := registry.Load(); err != nil {
		fmt.Printf("Note: Starting with empty registry: %v\n", err)
	}

	// Example 2: Manually register models
	fmt.Println("\n2. Registering models manually...")

	diffusionModel := &models.ModelInfo{
		ID:          "z-image-turbo",
		Name:        "Z-Image Turbo",
		Type:        models.ModelTypeDiffusion,
		Path:        "D:\\hf-mirror\\Z-Image-Turbo-GGUF\\z_image_turbo-Q4_K_M.gguf",
		Format:      models.FormatGGUF,
		Size:        4 * 1024 * 1024 * 1024, // 4GB
		Tags:        []string{"turbo", "fast", "q4"},
		Description: "Fast diffusion model for quick image generation",
		Source:      "https://huggingface.co/model",
	}

	if err := registry.Register(diffusionModel); err != nil {
		log.Printf("Failed to register diffusion model: %v", err)
	} else {
		fmt.Printf("✓ Registered: %s (%s)\n", diffusionModel.Name, diffusionModel.HumanSize())
	}

	llmModel := &models.ModelInfo{
		ID:          "qwen3-4b",
		Name:        "Qwen3 4B",
		Type:        models.ModelTypeLLM,
		Path:        "D:\\hf-mirror\\Z-Image-Turbo-GGUF\\Qwen3-4B-Instruct-2507-Q4_K_M.gguf",
		Format:      models.FormatGGUF,
		Size:        2 * 1024 * 1024 * 1024, // 2GB
		Tags:        []string{"llm", "qwen", "4b", "q4"},
		Description: "Qwen3 4B language model",
		Source:      "https://huggingface.co/qwen",
	}

	if err := registry.Register(llmModel); err != nil {
		log.Printf("Failed to register LLM model: %v", err)
	} else {
		fmt.Printf("✓ Registered: %s (%s)\n", llmModel.Name, llmModel.HumanSize())
	}

	vaeModel := &models.ModelInfo{
		ID:          "sd-vae",
		Name:        "SD VAE",
		Type:        models.ModelTypeVAE,
		Path:        "D:\\hf-mirror\\Z-Image-Turbo-GGUF\\diffusion_pytorch_model.safetensors",
		Format:      models.FormatSafetensors,
		Size:        300 * 1024 * 1024, // 300MB
		Tags:        []string{"vae", "sd"},
		Description: "VAE for Stable Diffusion",
	}

	if err := registry.Register(vaeModel); err != nil {
		log.Printf("Failed to register VAE model: %v", err)
	} else {
		fmt.Printf("✓ Registered: %s (%s)\n", vaeModel.Name, vaeModel.HumanSize())
	}

	// Example 3: List all models
	fmt.Println("\n3. Listing all registered models...")
	allModels := registry.List(models.ModelTypeUnknown)
	fmt.Printf("Total models: %d\n", len(allModels))
	for _, m := range allModels {
		fmt.Printf("  - [%s] %s (%s) - %s\n", m.Type, m.Name, m.HumanSize(), m.Path)
	}

	// Example 4: List models by type
	fmt.Println("\n4. Listing diffusion models only...")
	diffusionModels := registry.List(models.ModelTypeDiffusion)
	for _, m := range diffusionModels {
		fmt.Printf("  - %s (%s)\n", m.Name, m.HumanSize())
	}

	// Example 5: Search models
	fmt.Println("\n5. Searching for models with 'qwen'...")
	searchResults := registry.Search("qwen")
	for _, m := range searchResults {
		fmt.Printf("  - Found: %s (%s)\n", m.Name, m.Path)
	}

	// Example 6: Filter by tag
	fmt.Println("\n6. Filtering models by tag 'q4'...")
	q4Models := registry.FilterByTag("q4")
	for _, m := range q4Models {
		fmt.Printf("  - %s (%s)\n", m.Name, m.HumanSize())
	}

	// Example 7: Get model by ID
	fmt.Println("\n7. Getting model by ID...")
	model, err := registry.Get("z-image-turbo")
	if err != nil {
		log.Printf("Failed to get model: %v", err)
	} else {
		fmt.Printf("Found: %s\n", model.Name)
		fmt.Printf("  Path: %s\n", model.Path)
		fmt.Printf("  Size: %s\n", model.HumanSize())
		fmt.Printf("  Tags: %v\n", model.Tags)
	}

	// Example 8: Get registry statistics
	fmt.Println("\n8. Registry statistics...")
	stats := registry.GetStats()
	fmt.Printf("Total models: %d\n", stats.TotalModels)
	fmt.Printf("Total size: %s\n", formatBytes(stats.TotalSize))
	fmt.Printf("Average size: %s\n", formatBytes(stats.AverageSize))
	fmt.Println("By type:")
	for t, count := range stats.ByType {
		fmt.Printf("  - %s: %d\n", t, count)
	}

	// Example 9: Validate registry
	fmt.Println("\n9. Validating registry...")
	validationErrors := registry.Validate()
	if len(validationErrors) > 0 {
		fmt.Println("Validation errors:")
		for _, err := range validationErrors {
			fmt.Printf("  - %v\n", err)
		}
	} else {
		fmt.Println("✓ All models validated successfully!")
	}

	// Example 10: Auto-detect models in directory
	fmt.Println("\n10. Auto-detecting models in directory...")
	// Note: This would scan a directory for model files
	// For this example, we'll skip the actual scanning
	fmt.Println("Note: Auto-detection would scan a directory for model files")
	fmt.Println("Usage: models.AutoRegister(registry, \"./models\", true)")

	// Example 11: Update model usage
	fmt.Println("\n11. Updating model usage...")
	if model != nil {
		model.UpdateLastUsed()
		fmt.Printf("Model '%s' use count: %d\n", model.Name, model.UseCount)
		fmt.Printf("Last used: %v\n", model.LastUsed)
	}

	// Example 12: Model download (example URLs)
	fmt.Println("\n12. Model download examples...")
	fmt.Println("HuggingFace:")
	fmt.Println("  hf := models.NewHuggingFaceDownloader(\"user/repo\", \"./models\")")
	fmt.Println("  model, err := hf.Download(\"model.gguf\")")
	fmt.Println("\nCivitai:")
	fmt.Println("  civ := models.NewCivitaiDownloader(\"./models\")")
	fmt.Println("  model, err := civ.DownloadByID(12345, \"model.safetensors\")")

	fmt.Println("\n✅ Model management example completed!")
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
