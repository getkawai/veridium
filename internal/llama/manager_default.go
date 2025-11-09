//go:build !darwin

package llama

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

// IsLlamaCppInstalled checks if llama.cpp binaries exist and are valid
// Default implementation for non-macOS platforms
func (lcm *LlamaCppReleaseManager) IsLlamaCppInstalled() bool {
	// Check for main binary
	mainBinaryPath := lcm.GetMainBinaryPath()
	if _, err := os.Stat(mainBinaryPath); os.IsNotExist(err) {
		return false
	}

	// Verify the binary integrity by trying to run it
	if err := lcm.VerifyInstalledBinary(); err != nil {
		log.Printf("Binary integrity check failed: %v", err)
		return false
	}

	return true
}

// VerifyInstalledBinary verifies the installed binary
// Default implementation for non-macOS platforms
func (lcm *LlamaCppReleaseManager) VerifyInstalledBinary() error {
	binaryPath := lcm.GetMainBinaryPath()

	// Check if the binary exists
	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("binary file not found: %w", err)
	}

	// On Unix systems, check if the binary is executable
	if runtime.GOOS != "windows" {
		info, err := os.Stat(binaryPath)
		if err != nil {
			return fmt.Errorf("failed to get binary info: %w", err)
		}
		if info.Mode()&0111 == 0 {
			return fmt.Errorf("binary is not executable")
		}
	}

	// Try to run the binary with --help to verify it works
	cmd := exec.Command(binaryPath, "--help")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("binary failed to execute: %w", err)
	}

	return nil
}

