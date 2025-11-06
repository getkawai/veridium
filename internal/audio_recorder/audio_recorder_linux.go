//go:build linux

package audio_recorder

import (
	"fmt"
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
		return tool, fmt.Errorf("%s not found. Please install %s", tool, tool)
	}
	
	return tool, nil
}

