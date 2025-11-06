//go:build darwin

package tts

import (
	"fmt"
	"os/exec"
	"strings"
)

// Speak converts text to speech using default voice
func (s *TTSService) Speak(text string) error {
	return s.SpeakWithVoice(text, "")
}

// SpeakWithVoice converts text to speech using specified voice
func (s *TTSService) SpeakWithVoice(text, voice string) error {
	args := []string{}

	if voice != "" {
		args = append(args, "-v", voice)
	}

	args = append(args, text)

	cmd := exec.Command("say", args...)
	return cmd.Run()
}

// SpeakWithRate converts text to speech with custom speaking rate
func (s *TTSService) SpeakWithRate(text string, rate int) error {
	cmd := exec.Command("say", "-r", fmt.Sprintf("%d", rate), text)
	return cmd.Run()
}

// SpeakToFile saves speech to an audio file
func (s *TTSService) SpeakToFile(text, outputPath string) error {
	return s.SpeakToFileWithVoice(text, outputPath, "")
}

// SpeakToFileWithVoice saves speech to file with specified voice
func (s *TTSService) SpeakToFileWithVoice(text, outputPath, voice string) error {
	args := []string{"-o", outputPath}

	if voice != "" {
		args = append(args, "-v", voice)
	}

	args = append(args, text)

	cmd := exec.Command("say", args...)
	return cmd.Run()
}

// ListVoices returns available TTS voices
func (s *TTSService) ListVoices() ([]Voice, error) {
	cmd := exec.Command("say", "-v", "?")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	voices := make([]Voice, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse format: "VoiceName    language    # description"
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		voice := Voice{
			Name:     parts[0],
			Language: parts[1],
		}

		// Try to determine gender from description
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "female") || strings.Contains(lowerLine, "woman") {
			voice.Gender = "female"
		} else if strings.Contains(lowerLine, "male") || strings.Contains(lowerLine, "man") {
			voice.Gender = "male"
		} else {
			voice.Gender = "unknown"
		}

		voices = append(voices, voice)
	}

	return voices, nil
}

// GetDefaultVoice returns the system default voice
func (s *TTSService) GetDefaultVoice() (string, error) {
	cmd := exec.Command("defaults", "read", "com.apple.speech.voice.prefs", "SelectedVoiceName")
	output, err := cmd.Output()
	if err != nil {
		return "Samantha", nil // Fallback to Samantha
	}

	return strings.TrimSpace(string(output)), nil
}

// Stop stops any ongoing speech
func (s *TTSService) Stop() error {
	cmd := exec.Command("killall", "say")
	return cmd.Run()
}

// IsPlatformSupported returns whether TTS is supported on current platform
func (s *TTSService) IsPlatformSupported() bool {
	return true
}

// GetPlatformInfo returns information about TTS capabilities
func (s *TTSService) GetPlatformInfo() map[string]interface{} {
	return map[string]interface{}{
		"platform":  "darwin",
		"supported": true,
		"engine":    "macOS Speech Synthesis (say)",
		"quality":   "excellent",
		"offline":   true,
		"voices":    "70+",
	}
}
