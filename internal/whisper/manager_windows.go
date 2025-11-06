//go:build windows

package whisper

import (
	"fmt"
	"log"
)

// InstallWhisper provides instructions for Windows installation
func (m *Manager) InstallWhisper() error {
	if m.IsWhisperInstalled() {
		log.Println("whisper-cpp is already installed")
		return nil
	}

	return fmt.Errorf("please install whisper-cpp manually on Windows:\n" +
		"1. Download pre-built binaries from: https://github.com/ggml-org/whisper.cpp/releases\n" +
		"2. Extract and add to PATH\n" +
		"Or build from source: https://github.com/ggml-org/whisper.cpp")
}
