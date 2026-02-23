package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/veridium/pkg/stablediffusion"
)

func main() {
	fmt.Println("Stable Diffusion Library Load Test")
	fmt.Println("===================================")

	// Initialize library
	fmt.Println("\n1. Loading library...")
	libPath := stablediffusion.GetLibraryPath()
	if err := stablediffusion.InitLibrary(libPath); err != nil {
		fmt.Printf("❌ Failed to load library: %v\n", err)
		fmt.Printf("   Make sure the stable-diffusion library is available at %s\n", libPath)
		os.Exit(1)
	}
	fmt.Println("✅ Library loaded successfully!")

	fmt.Println("\n2. Getting version info...")
	version := stablediffusion.GetLibraryVersion()
	commit := "n/a"
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)

	fmt.Println("\n✅ All tests passed! Library is working correctly.")
	fmt.Println("\nNote: To generate images, you need to:")
	fmt.Println("  1. Download model files (.gguf format)")
	fmt.Println("  2. Update the model paths in txt2img.go example")
	fmt.Println("  3. Run: go run pkg/stablediffusion/examples/txt2img/txt2img.go")

	os.Exit(0)
}
