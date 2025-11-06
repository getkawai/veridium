//go:build linux

package tts

import (
	"fmt"
	"os/exec"
)

// Speak converts text to speech using default voice
func (s *TTSService) Speak(text string) error {
	// Try espeak first, then festival
	if _, err := exec.LookPath("espeak"); err == nil {
		cmd := exec.Command("espeak", text)
		return cmd.Run()
	}

	if _, err := exec.LookPath("festival"); err == nil {
		cmd := exec.Command("festival", "--tts")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		if _, err := stdin.Write([]byte(text)); err != nil {
			return err
		}
		stdin.Close()

		return cmd.Wait()
	}

	return fmt.Errorf("no TTS engine found (install espeak or festival)")
}

// SpeakWithVoice converts text to speech using specified voice
func (s *TTSService) SpeakWithVoice(text, voice string) error {
	return s.Speak(text) // Linux TTS engines don't support voice selection easily
}

// SpeakWithRate converts text to speech with custom speaking rate
func (s *TTSService) SpeakWithRate(text string, rate int) error {
	return s.Speak(text) // Rate not supported
}

// SpeakToFile saves speech to an audio file
func (s *TTSService) SpeakToFile(text, outputPath string) error {
	return fmt.Errorf("SpeakToFile not implemented on Linux")
}

// SpeakToFileWithVoice saves speech to file with specified voice
func (s *TTSService) SpeakToFileWithVoice(text, outputPath, voice string) error {
	return s.SpeakToFile(text, outputPath)
}

// ListVoices returns available TTS voices
func (s *TTSService) ListVoices() ([]Voice, error) {
	return nil, fmt.Errorf("ListVoices not supported on Linux")
}

// GetDefaultVoice returns the system default voice
func (s *TTSService) GetDefaultVoice() (string, error) {
	return "", fmt.Errorf("GetDefaultVoice not supported on Linux")
}

// Stop stops any ongoing speech
func (s *TTSService) Stop() error {
	return fmt.Errorf("Stop not supported on Linux")
}

// IsPlatformSupported returns whether TTS is supported
func (s *TTSService) IsPlatformSupported() bool {
	// Check if espeak or festival is available
	_, errEspeak := exec.LookPath("espeak")
	_, errFestival := exec.LookPath("festival")
	return errEspeak == nil || errFestival == nil
}

// GetPlatformInfo returns information about TTS capabilities
func (s *TTSService) GetPlatformInfo() map[string]interface{} {
	return map[string]interface{}{
		"platform":  "linux",
		"supported": s.IsPlatformSupported(),
		"engine":    "espeak/festival",
		"quality":   "basic",
		"offline":   true,
		"voices":    "varies",
	}
}
