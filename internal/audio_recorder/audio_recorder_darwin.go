//go:build darwin

package audio_recorder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// startPlatformRecording starts recording audio using platform-specific tools
func startPlatformRecording(outputPath string) (*exec.Cmd, error) {
	// Use sox on macOS for recording
	// Format: WAV, 16-bit signed integer, 16kHz, mono (optimal for Whisper)
	// Using 'sox -d' instead of 'rec' for better control
	cmd := exec.Command("sox", "-d", "-r", "16000", "-c", "1", "-b", "16", "-e", "signed-integer", outputPath)

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

	// On macOS, send SIGINT (Ctrl+C) for graceful shutdown
	// This allows sox to properly finalize the WAV file
	if err := proc.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Failed to send SIGINT: %v, killing process\n", err)
		return proc.Process.Kill()
	}

	return nil
}

// checkPlatformRecordingTool checks if the recording tool is available
func checkPlatformRecordingTool() (string, error) {
	tool := "sox"
	cmd := exec.Command("which", "sox")

	if err := cmd.Run(); err != nil {
		return tool, fmt.Errorf("%s not found", tool)
	}

	return tool, nil
}

// installPlatformRecordingTool installs the recording tool
func installPlatformRecordingTool() error {
	log.Println("🔧 sox not found, attempting auto-installation...")

	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("homebrew is not installed. Please install Homebrew first: https://brew.sh")
	}

	// Install sox
	log.Println("   Installing sox via Homebrew...")
	cmd := exec.Command("brew", "install", "sox")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install sox: %w\nOutput: %s", err, string(output))
	}

	log.Println("✅ sox installed successfully")
	return nil
}
