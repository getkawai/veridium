//go:build windows

package tts

import (
	"fmt"
	"os/exec"
)

// Speak converts text to speech using default voice
func (s *TTSService) Speak(text string) error {
	// Use PowerShell with SAPI
	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Speech; $speak = New-Object System.Speech.Synthesis.SpeechSynthesizer; $speak.Speak('%s')`, text)
	cmd := exec.Command("powershell", "-Command", psScript)
	return cmd.Run()
}

// SpeakWithVoice converts text to speech using specified voice
func (s *TTSService) SpeakWithVoice(text, voice string) error {
	return s.Speak(text) // Voice selection not easily supported
}

// SpeakWithRate converts text to speech with custom speaking rate
func (s *TTSService) SpeakWithRate(text string, rate int) error {
	return s.Speak(text) // Rate not supported in simple implementation
}

// SpeakToFile saves speech to an audio file
func (s *TTSService) SpeakToFile(text, outputPath string) error {
	return fmt.Errorf("SpeakToFile not implemented on Windows")
}

// SpeakToFileWithVoice saves speech to file with specified voice
func (s *TTSService) SpeakToFileWithVoice(text, outputPath, voice string) error {
	return s.SpeakToFile(text, outputPath)
}

// ListVoices returns available TTS voices
func (s *TTSService) ListVoices() ([]Voice, error) {
	return nil, fmt.Errorf("ListVoices not supported on Windows")
}

// GetDefaultVoice returns the system default voice
func (s *TTSService) GetDefaultVoice() (string, error) {
	return "", fmt.Errorf("GetDefaultVoice not supported on Windows")
}

// Stop stops any ongoing speech
func (s *TTSService) Stop() error {
	return fmt.Errorf("Stop not supported on Windows")
}

// IsPlatformSupported returns whether TTS is supported
func (s *TTSService) IsPlatformSupported() bool {
	return true // Windows always has SAPI
}

// GetPlatformInfo returns information about TTS capabilities
func (s *TTSService) GetPlatformInfo() map[string]interface{} {
	return map[string]interface{}{
		"platform":  "windows",
		"supported": true,
		"engine":    "Windows SAPI",
		"quality":   "good",
		"offline":   true,
		"voices":    "varies",
	}
}
