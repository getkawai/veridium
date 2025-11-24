package main

import (
	"fmt"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
)

func main() {
	fmt.Println("=== Hardware Validation Test for Reasoning Mode ===\n")

	// Detect hardware
	installer := llama.NewLlamaCppInstaller()
	specs := installer.HardwareSpecs

	fmt.Printf("Detected Hardware:\n")
	fmt.Printf("  Total RAM:      %d GB\n", specs.TotalRAM)
	fmt.Printf("  Available RAM:  %d GB\n", specs.AvailableRAM)
	fmt.Printf("  CPU:            %s\n", specs.CPU)
	fmt.Printf("  CPU Cores:      %d\n", specs.CPUCores)
	fmt.Printf("  GPU:            %s\n", specs.GPUModel)
	fmt.Printf("  GPU Memory:     %d GB\n\n", specs.GPUMemory)

	// Test each reasoning mode
	modes := []services.ReasoningMode{
		services.ReasoningDisabled,
		services.ReasoningEnabled,
		services.ReasoningVerbose,
	}

	fmt.Println("Testing Reasoning Modes:")
	fmt.Println("=" + "==============================================")

	for _, mode := range modes {
		config := services.ReasoningConfig{Mode: mode}
		
		fmt.Printf("\n%s Mode:\n", mode)
		fmt.Printf("  %s\n", config.GetModeDescription())
		
		// Get requirements
		req := config.GetHardwareRequirements()
		fmt.Printf("  Requirements: %s\n", req.Description)
		fmt.Printf("    - Min RAM:       %d GB\n", req.MinRAM)
		fmt.Printf("    - Min CPU Cores: %d\n", req.MinCPUCores)
		fmt.Printf("    - GPU:           %v\n", map[bool]string{true: "Recommended", false: "Not required"}[req.RecommendGPU])
		
		// Validate hardware
		valid, reason := config.ValidateHardware(specs)
		if valid {
			fmt.Printf("  ✅ Hardware is SUFFICIENT\n")
		} else {
			fmt.Printf("  ❌ Hardware is INSUFFICIENT\n")
			fmt.Printf("  Reason: %s\n", reason)
		}
		
		// Expected performance
		perf := config.GetExpectedPerformance()
		fmt.Printf("  Expected Performance:\n")
		fmt.Printf("    - Speed:            %s\n", perf["speed"])
		fmt.Printf("    - Token efficiency: %s\n", perf["token_efficiency"])
		fmt.Printf("    - Max turns:        %s\n", perf["max_turns"])
	}

	// Suggest best mode for this hardware
	fmt.Println("\n" + "=" + "==============================================")
	suggested := services.SuggestModeForHardware(specs)
	fmt.Printf("\n💡 Suggested Mode for Your Hardware: %s\n", suggested)
	
	suggestedConfig := services.ReasoningConfig{Mode: suggested}
	fmt.Printf("   %s\n", suggestedConfig.GetModeDescription())
	
	fmt.Println("\n=== Test Complete ===")
}

