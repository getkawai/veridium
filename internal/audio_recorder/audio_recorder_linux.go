//go:build linux

package audio_recorder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// startPlatformRecording starts recording audio using platform-specific tools
func startPlatformRecording(outputPath string) (*exec.Cmd, error) {
	// Use arecord on Linux
	cmd := exec.Command("arecord", "-f", "S16_LE", "-r", "16000", "-c", "1", "-t", "wav", outputPath)
	
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start recording: %w", err)
	}
	
	return cmd, nil
}

// stopPlatformRecording stops the recording process gracefully
func stopPlatformRecording(proc *exec.Cmd) error {
	if proc == nil || proc.Process == nil {
		return nil
	}
	
	// On Linux, send SIGINT (Ctrl+C) for graceful shutdown
	// This allows arecord to properly finalize the WAV file
	if err := proc.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Failed to send SIGINT: %v, killing process\n", err)
		return proc.Process.Kill()
	}
	
	return nil
}

// checkPlatformRecordingTool checks if the recording tool is available
func checkPlatformRecordingTool() (string, error) {
	tool := "arecord"
	cmd := exec.Command("which", "arecord")
	
	if err := cmd.Run(); err != nil {
		return tool, fmt.Errorf("%s not found", tool)
	}
	
	return tool, nil
}

// installPlatformRecordingTool installs the recording tool
func installPlatformRecordingTool() error {
	log.Println("🔧 arecord not found, attempting auto-installation...")
	
	// Try to detect package manager and install
	// Try apt-get (Debian/Ubuntu)
	if _, err := exec.LookPath("apt-get"); err == nil {
		log.Println("   Installing alsa-utils via apt-get...")
		cmd := exec.Command("sudo", "apt-get", "install", "-y", "alsa-utils")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install alsa-utils: %w\nOutput: %s", err, string(output))
		}
		log.Println("✅ alsa-utils installed successfully")
		return nil
	}
	
	// Try yum (RHEL/CentOS/Fedora)
	if _, err := exec.LookPath("yum"); err == nil {
		log.Println("   Installing alsa-utils via yum...")
		cmd := exec.Command("sudo", "yum", "install", "-y", "alsa-utils")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install alsa-utils: %w\nOutput: %s", err, string(output))
		}
		log.Println("✅ alsa-utils installed successfully")
		return nil
	}
	
	// Try pacman (Arch Linux)
	if _, err := exec.LookPath("pacman"); err == nil {
		log.Println("   Installing alsa-utils via pacman...")
		cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "alsa-utils")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install alsa-utils: %w\nOutput: %s", err, string(output))
		}
		log.Println("✅ alsa-utils installed successfully")
		return nil
	}
	
	return fmt.Errorf("no supported package manager found. Please install alsa-utils manually")
}

