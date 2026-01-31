package main

import (
	"fmt"
	"log"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
)

// InitializeStableDiffusion should be called during application startup.
// It ensures the library is downloaded and ready to use.
func InitializeStableDiffusion() error {
	// Ensure library is installed (auto-download if needed)
	if err := stablediffusion.EnsureLibrary(); err != nil {
		return fmt.Errorf("failed to initialize Stable Diffusion: %w", err)
	}
	return nil
}

// Example: Application startup
func main() {
	fmt.Println("Starting Application...")
	fmt.Println()

	// Initialize Stable Diffusion (one-time setup)
	if err := InitializeStableDiffusion(); err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}

	fmt.Println("✅ Application ready!")
	fmt.Println()

	// Now you can use Stable Diffusion in your app
	// Example: Create service, handle requests, etc.

	// Check if we have models
	// If not, prompt user to download or download automatically

	fmt.Println("Application is running...")
	fmt.Println("Stable Diffusion is ready for image generation")
}
