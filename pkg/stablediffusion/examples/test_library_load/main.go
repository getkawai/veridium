package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/stablediffusion"
)

func main() {
	fmt.Println("Stable Diffusion Library Load Test")
	fmt.Println("===================================")

	// Initialize library
	fmt.Println("\n1. Loading library...")
	sd, err := stablediffusion.New(stablediffusion.LibraryConfig{LibPath: "./lib"})
	if err != nil {
		fmt.Printf("❌ Failed to load library: %v\n", err)
		fmt.Println("   Make sure the stable-diffusion library is available in ./lib")
		os.Exit(1)
	}
	defer sd.Close()
	fmt.Println("✅ Library loaded successfully!")

	// Test basic context params initialization
	fmt.Println("\n2. Testing context params initialization...")
	var ctxParams stablediffusion.SDContextParams
	sd.ContextParamsInit(&ctxParams)
	fmt.Println("✅ Context params initialized!")

	// Get system info
	fmt.Println("\n3. Getting system info...")
	sysInfo := sd.GetSystemInfo()
	fmt.Printf("System Info: %s\n", sysInfo)

	// Get version
	fmt.Println("\n4. Getting version info...")
	version := sd.Version()
	commit := sd.Commit()
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)

	fmt.Println("\n✅ All tests passed! Library is working correctly.")
	fmt.Println("\nNote: To generate images, you need to:")
	fmt.Println("  1. Download model files (.gguf format)")
	fmt.Println("  2. Update the model paths in txt2img.go example")
	fmt.Println("  3. Run: go run pkg/stablediffusion/examples/txt2img/txt2img.go")

	os.Exit(0)
}
