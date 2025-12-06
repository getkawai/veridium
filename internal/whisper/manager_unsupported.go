//go:build !darwin && !linux && !windows

package whisper

import (
	"fmt"
	"log"
	"runtime"
)

// InstallWhisper is not supported on this platform
func (m *Manager) InstallWhisper() error {
	if m.IsWhisperInstalled() {
		log.Println("whisper-cpp is already installed")
		return nil
	}

	return fmt.Errorf("whisper-cpp installation not supported on %s. Please install manually from: https://github.com/ggml-org/whisper.cpp", runtime.GOOS)
}

// InstallFFmpeg is not supported on this platform
func (m *Manager) InstallFFmpeg() error {
	if m.IsFFmpegInstalled() {
		log.Println("ffmpeg is already installed")
		return nil
	}

	return fmt.Errorf("ffmpeg installation not supported on %s. Please install manually from: https://ffmpeg.org/download.html", runtime.GOOS)
}

// detectAvailableRAM returns default RAM for unsupported platforms
func (m *Manager) detectAvailableRAM() int64 {
	log.Printf("⚠️  RAM detection not supported on %s, defaulting to 8GB", runtime.GOOS)
	return 8 // Default to 8GB
}
