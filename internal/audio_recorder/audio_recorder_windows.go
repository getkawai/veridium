//go:build windows

package audio_recorder

import (
	"fmt"
	"log"
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
		return tool, fmt.Errorf("%s not found", tool)
	}
	
	return tool, nil
}

// installPlatformRecordingTool installs the recording tool
func installPlatformRecordingTool() error {
	log.Println("🔧 ffmpeg not found, attempting auto-installation...")
	
	// Check if winget is available (Windows 10+)
	if _, err := exec.LookPath("winget"); err == nil {
		log.Println("   Installing ffmpeg via winget...")
		cmd := exec.Command("winget", "install", "ffmpeg", "--accept-package-agreements", "--accept-source-agreements")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install ffmpeg: %w\nOutput: %s", err, string(output))
		}
		log.Println("✅ ffmpeg installed successfully")
		return nil
	}
	
	// Check if chocolatey is available
	if _, err := exec.LookPath("choco"); err == nil {
		log.Println("   Installing ffmpeg via chocolatey...")
		cmd := exec.Command("choco", "install", "ffmpeg", "-y")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install ffmpeg: %w\nOutput: %s", err, string(output))
		}
		log.Println("✅ ffmpeg installed successfully")
		return nil
	}
	
	return fmt.Errorf("no supported package manager found (winget or chocolatey). Please install ffmpeg manually from: https://ffmpeg.org/download.html")
}

