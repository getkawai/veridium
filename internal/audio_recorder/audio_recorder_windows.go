//go:build windows

package audio_recorder

import (
	"fmt"
	"os/exec"
)

// startPlatformRecording starts recording audio using platform-specific tools
func startPlatformRecording(outputPath string) (*exec.Cmd, error) {
	// Use ffmpeg on Windows (requires ffmpeg to be installed)
	cmd := exec.Command("ffmpeg", "-f", "dshow", "-i", "audio=", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)
	
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
	
	// On Windows, kill the process directly
	return proc.Process.Kill()
}

// checkPlatformRecordingTool checks if the recording tool is available
func checkPlatformRecordingTool() (string, error) {
	tool := "ffmpeg"
	cmd := exec.Command("where", "ffmpeg")
	
	if err := cmd.Run(); err != nil {
		return tool, fmt.Errorf("%s not found. Please install %s", tool, tool)
	}
	
	return tool, nil
}

