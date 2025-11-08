package llama

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

// TestAutoDownloadIntegration tests the llama-cli integration for auto-downloading models
func TestAutoDownloadIntegration(t *testing.T) {
	log.Println("🧪 Testing llama-cli Auto-Download Integration")
	log.Println("================================================")

	// Initialize service
	service, err := NewService()
	if err != nil {
		t.Fatalf("❌ Failed to create service: %v", err)
	}

	// Get models directory
	modelsDir := service.manager.GetModelsDirectory()
	log.Printf("📁 Models directory: %s", modelsDir)

	// List existing models
	models, err := service.GetAvailableModels()
	if err != nil {
		t.Logf("⚠️  Failed to check models: %v", err)
	} else {
		log.Printf("📦 Found %d existing model(s)", len(models))
		for _, modelName := range models {
			log.Printf("   - %s", modelName)
		}
	}

	// Test hardware detection
	log.Println("\n🔍 Detecting hardware specs...")
	specs := DetectHardwareSpecs()
	log.Printf("   CPU Cores: %d", specs.CPUCores)
	log.Printf("   Total RAM: %d GB", specs.TotalRAM)
	log.Printf("   Available RAM: %d GB", specs.AvailableRAM)
	log.Printf("   GPU: %s", specs.GPUModel)
	log.Printf("   VRAM: %d GB", specs.GPUMemory)

	// Test model selection
	log.Println("\n📦 Selecting optimal model...")
	modelSpec := SelectOptimalQwenModel(specs.AvailableRAM)
	log.Printf("   Selected: %s", modelSpec.Repo)
	log.Printf("   Quantization: %s", modelSpec.Quantization)
	log.Printf("   Parameters: %s", modelSpec.Parameters)
	log.Printf("   Min RAM: %d GB", modelSpec.MinRAM)

	// Test auto-download (only if no models exist)
	if len(models) == 0 {
		log.Println("\n📥 Testing auto-download with llama-cli...")
		if err := service.AutoDownloadRecommendedModel(); err != nil {
			t.Fatalf("❌ Auto-download failed: %v", err)
		}

		// Verify download
		models, err = service.GetAvailableModels()
		if err != nil {
			t.Fatalf("❌ Failed to check models after download: %v", err)
		}

		if len(models) == 0 {
			t.Fatal("❌ No models found after download!")
		}

		log.Printf("✅ Download successful! Found %d model(s):", len(models))
		for _, modelName := range models {
			modelPath := filepath.Join(modelsDir, modelName)
			info, err := os.Stat(modelPath)
			if err == nil {
				sizeMB := float64(info.Size()) / (1024 * 1024)
				log.Printf("   - %s (%.1f MB)", modelName, sizeMB)
			} else {
				log.Printf("   - %s", modelName)
			}
		}
	} else {
		log.Println("\n✅ Models already exist, skipping download test")
	}

	log.Println("\n🎉 All tests passed!")
}

