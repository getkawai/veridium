package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/veridium/internal/stablediffusion"
)

func main() {
	manager := stablediffusion.NewStableDiffusionReleaseManager()

	fmt.Printf("Binary Path: %s\n", manager.GetBinaryPath())

	installed := manager.IsStableDiffusionInstalled()
	fmt.Printf("IsInstalled: %v\n", installed)

	if !installed {
		// Try to see why
		binPath := manager.GetBinaryPath()
		info, err := os.Stat(binPath)
		if err != nil {
			fmt.Printf("Stat error: %v\n", err)
		} else {
			fmt.Printf("Mode: %s\n", info.Mode())
		}

		err = manager.VerifyInstalledBinary()
		if err != nil {
			fmt.Printf("Verification error: %v\n", err)
		}
	}

	modelsSpec, err := manager.CheckInstalledModels()
	fmt.Printf("Models Error: %v\n", err)
	fmt.Printf("Models Count: %d\n", len(modelsSpec))
	for _, m := range modelsSpec {
		fmt.Printf(" - %s\n", m)
	}
}
