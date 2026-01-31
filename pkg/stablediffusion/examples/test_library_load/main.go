package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/veridium/pkg/stablediffusion/sd"
)

func main() {
	fmt.Println("Stable Diffusion Library Load Test")
	fmt.Println("===================================")

	// Test basic context params initialization
	// This will trigger library loading automatically
	fmt.Println("\n1. Testing library load and context params initialization...")
	var ctxParams sd.SDContextParams
	sd.ContextParamsInit(&ctxParams)
	fmt.Println("✅ Library loaded and context params initialized!")

	// Get system info
	fmt.Println("\n2. Getting system info...")
	sysInfo := sd.GetSystemInfo()
	fmt.Printf("System Info: %s\n", sysInfo)

	// Get number of physical cores
	fmt.Println("\n3. Getting CPU cores...")
	cores := sd.GetNumPhysicalCores()
	fmt.Printf("Physical CPU Cores: %d\n", cores)

	fmt.Println("\n✅ All tests passed! Library is working correctly.")
	fmt.Println("\nNote: To generate images, you need to:")
	fmt.Println("  1. Download model files (.gguf format)")
	fmt.Println("  2. Update the model paths in txt2img.go example")
	fmt.Println("  3. Run: go run pkg/stablediffusion/examples/txt2img/txt2img.go")

	os.Exit(0)
}
