package main

import (
	"fmt"
	"os"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
)

func main() {
	fmt.Println("Stable Diffusion Simple Test")
	fmt.Println("=============================")

	modelPath := "models/stable-diffusion/sd-v1-4-q8.gguf"

	// Check if model exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		fmt.Printf("Error: Model file not found at %s\n", modelPath)
		fmt.Println("Please download the model first.")
		os.Exit(1)
	}

	fmt.Printf("Using model: %s\n", modelPath)
	fmt.Println("Creating Stable Diffusion instance...")

	// Create instance with minimal config
	sd, err := stablediffusion.NewStableDiffusion(&stablediffusion.ContextParams{
		ModelPath:          modelPath,
		NThreads:           4,
		OffloadParamsToCPU: true,
	})

	if err != nil {
		fmt.Printf("Failed to create instance: %v\n", err)
		os.Exit(1)
	}
	defer sd.Free()

	fmt.Println("✅ Instance created successfully!")
	fmt.Println("\nGenerating image...")

	// Generate simple image
	err = sd.GenerateImage(&stablediffusion.ImgGenParams{
		Prompt:      "a beautiful sunset over mountains",
		Width:       512,
		Height:      512,
		SampleSteps: 20,
		CfgScale:    7.0,
		Seed:        42,
	}, "output_test.png")

	if err != nil {
		fmt.Printf("Failed to generate image: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Image generated successfully!")
	fmt.Println("Output saved to: output_test.png")
}
