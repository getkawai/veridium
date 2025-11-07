//go:build windows

package llama

import (
	"fmt"
	"log"
	"os/exec"
)

// InstallLlamaCpp attempts to install llama.cpp on Windows using winget or scoop
func (lcm *LlamaCppReleaseManager) InstallLlamaCpp() error {
	if lcm.IsLlamaCppInstalled() {
		log.Println("llama.cpp is already installed")
		return nil
	}

	log.Println("Installing llama.cpp on Windows...")

	// Try winget first (built-in on Windows 10/11)
	if _, err := exec.LookPath("winget"); err == nil {
		log.Println("Trying winget installation...")
		cmd := exec.Command("winget", "install", "--id", "ggerganov.llama.cpp",
			"--silent",
			"--accept-package-agreements",
			"--accept-source-agreements",
			"--disable-interactivity")
		output, err := cmd.CombinedOutput()
		if err == nil {
			log.Println("llama.cpp installed successfully via winget")
			return nil
		}
		log.Printf("winget installation failed: %v\nOutput: %s", err, string(output))
	}

	// Try scoop as fallback (scoop is non-interactive by default)
	if _, err := exec.LookPath("scoop"); err == nil {
		log.Println("Trying scoop installation...")
		cmd := exec.Command("scoop", "install", "llama.cpp", "--no-cache")
		output, err := cmd.CombinedOutput()
		if err == nil {
			log.Println("llama.cpp installed successfully via scoop")
			return nil
		}
		log.Printf("scoop installation failed: %v\nOutput: %s", err, string(output))
	}

	// If both package managers are not available or failed
	return fmt.Errorf("no package manager found (winget/scoop). Will fallback to GitHub download")
}
