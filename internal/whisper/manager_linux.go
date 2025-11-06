//go:build linux

package whisper

import (
	"fmt"
	"log"
)

// InstallWhisper provides instructions for Linux installation
func (m *Manager) InstallWhisper() error {
	if m.IsWhisperInstalled() {
		log.Println("whisper-cpp is already installed")
		return nil
	}

	return fmt.Errorf("please install whisper-cpp manually on Linux:\n" +
		"1. Install Homebrew for Linux: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"\n" +
		"2. Run: brew install whisper-cpp\n" +
		"Or build from source: https://github.com/ggml-org/whisper.cpp")
}
