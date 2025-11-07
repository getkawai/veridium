//go:build linux

package llama

import (
	"fmt"
	"log"
)

// InstallLlamaCpp attempts to install llama.cpp on Linux
// Note: Package managers on Linux require sudo which needs user interaction
// So we skip package manager and go directly to GitHub download
func (lcm *LlamaCppReleaseManager) InstallLlamaCpp() error {
	if lcm.IsLlamaCppInstalled() {
		log.Println("llama.cpp is already installed")
		return nil
	}

	// Package managers on Linux require sudo (user interaction)
	// Return error to trigger GitHub download fallback
	return fmt.Errorf("package manager installation requires sudo. Will use GitHub download")
}
