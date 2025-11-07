//go:build darwin

package llama

import (
	"fmt"
	"log"
	"os/exec"
)

// InstallLlamaCpp installs llama.cpp using Homebrew on macOS
func (lcm *LlamaCppReleaseManager) InstallLlamaCpp() error {
	if lcm.IsLlamaCppInstalled() {
		log.Println("llama.cpp is already installed")
		return nil
	}

	log.Println("Installing llama.cpp via Homebrew...")

	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew is not installed. Please install Homebrew first: https://brew.sh")
	}

	// Install llama.cpp
	cmd := exec.Command("brew", "install", "llama.cpp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install llama.cpp: %w\nOutput: %s", err, string(output))
	}

	log.Println("llama.cpp installed successfully via Homebrew")
	return nil
}
