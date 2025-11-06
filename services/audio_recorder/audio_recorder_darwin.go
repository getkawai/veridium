//go:build darwin

package services

import (
	"fmt"
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
	tool := "sox (rec command)"
	cmd := exec.Command("which", "rec")
	
	if err := cmd.Run(); err != nil {
		return tool, fmt.Errorf("%s not found. Please install %s", tool, tool)
	}
	
	return tool, nil
}

