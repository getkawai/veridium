package tts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TTSService provides text-to-speech using native OS capabilities
type TTSService struct{}

// Voice represents a TTS voice
type Voice struct {
	Name     string
	Language string
	Gender   string
}

// NewTTSService creates a new TTS service instance
func NewTTSService() (*TTSService, error) {
	return &TTSService{}, nil
}

// Platform-specific methods are implemented in:
// - darwin.go (macOS)
// - linux.go (Linux)
// - windows.go (Windows)
// Go compiler automatically selects the right file based on GOOS

// SpeakToAudio generates speech and returns audio data as byte array
// This is more convenient than SpeakToFile for frontend integration
func (s *TTSService) SpeakToAudio(text string) ([]byte, error) {
	return s.SpeakToAudioWithVoice(text, "")
}

// SpeakToAudioWithVoice generates speech with specified voice and returns audio data
func (s *TTSService) SpeakToAudioWithVoice(text, voice string) ([]byte, error) {
	// Create temporary file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("tts-%d.aiff", os.Getpid()))

	// Ensure cleanup
	defer os.Remove(tempFile)

	// Generate speech to temp file
	var err error
	if voice != "" {
		err = s.SpeakToFileWithVoice(text, tempFile, voice)
	} else {
		err = s.SpeakToFile(text, tempFile)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate speech: %w", err)
	}

	// Read the audio file
	audioData, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}

	return audioData, nil
}

// GetVoicesByLanguage returns voices filtered by language code
func (s *TTSService) GetVoicesByLanguage(langCode string) ([]Voice, error) {
	allVoices, err := s.ListVoices()
	if err != nil {
		return nil, err
	}

	filtered := make([]Voice, 0)
	for _, voice := range allVoices {
		if strings.HasPrefix(voice.Language, langCode) {
			filtered = append(filtered, voice)
		}
	}

	return filtered, nil
}

// GetRecommendedVoices returns high-quality voices for common languages
func (s *TTSService) GetRecommendedVoices() map[string]string {
	return map[string]string{
		"en-US": "Samantha",  // High quality female voice
		"en-GB": "Daniel",    // British male
		"en-AU": "Karen",     // Australian female
		"en-IN": "Veena",     // Indian female
		"id-ID": "Damayanti", // Indonesian female
		"ja-JP": "Kyoko",     // Japanese female
		"zh-CN": "Ting-Ting", // Chinese female
		"zh-TW": "Mei-Jia",   // Taiwanese female
		"es-ES": "Monica",    // Spanish female
		"fr-FR": "Amelie",    // French female
		"de-DE": "Anna",      // German female
		"it-IT": "Alice",     // Italian female
		"ko-KR": "Yuna",      // Korean female
		"ru-RU": "Milena",    // Russian female
		"ar-SA": "Maged",     // Arabic male
	}
}
