package main

import (
	"fmt"
	"os"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
)

func main() {
	fmt.Println("Stable Diffusion Auto-Setup Example")
	fmt.Println("====================================")
	fmt.Println()

	// Check if library is installed
	if stablediffusion.IsLibraryInstalled() {
		fmt.Println("✅ Library already installed")
		fmt.Printf("Location: %s\n", stablediffusion.GetLibraryPath())
		fmt.Printf("Version: %s\n", stablediffusion.GetLibraryVersion())
	} else {
		fmt.Println("⚠️  Library not found")
		fmt.Println("Starting automatic download...")
		fmt.Println()

		// Auto-download library
		if err := stablediffusion.EnsureLibrary(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println()
	fmt.Println("Library is ready to use!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Download a model file (see SETUP.md)")
	fmt.Println("2. Use NewStableDiffusion() to create instance")
	fmt.Println("3. Generate images with GenerateImage()")
}
