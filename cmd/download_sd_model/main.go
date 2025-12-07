package main

import (
	"log"

	"github.com/kawai-network/veridium/internal/stablediffusion"
)

func main() {
	manager := stablediffusion.NewStableDiffusionReleaseManager()

	// Get sd-turbo-q8_0 model spec
	models := stablediffusion.GetAvailableModels()
	var targetModel stablediffusion.ModelSpec
	for _, m := range models {
		if m.Name == "sd-turbo-q8_0" {
			targetModel = m
			break
		}
	}

	if targetModel.Name == "" {
		log.Fatal("Model sd-turbo-q8_0 not found in available models")
	}

	log.Printf("Downloading model: %s (%s)", targetModel.Name, targetModel.Filename)
	log.Printf("URL: %s", targetModel.URL)
	log.Printf("Expected size: %d MB", targetModel.Size)

	modelSpec := map[string]interface{}{
		"name":     targetModel.Name,
		"url":      targetModel.URL,
		"filename": targetModel.Filename,
		"size":     targetModel.Size,
	}

	err := manager.DownloadModel(modelSpec, func(progress float64) {
		if int(progress)%10 == 0 {
			log.Printf("Progress: %.1f%%", progress)
		}
	})

	if err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	log.Println("Model downloaded successfully!")

	// Verify
	models2, err := manager.CheckInstalledModels()
	if err != nil {
		log.Fatalf("Failed to check installed models: %v", err)
	}
	log.Printf("Installed models: %v", models2)
}
