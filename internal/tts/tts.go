package tts

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/pemistahl/lingua-go"
)

// TTSService provides text-to-speech using native OS capabilities
type TTSService struct {
	detector     lingua.LanguageDetector
	detectorOnce sync.Once
}

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
// Automatically detects language and selects appropriate voice
func (s *TTSService) SpeakToAudio(text string) ([]byte, error) {
	// Auto-detect language and select voice
	voice := s.SelectVoiceForText(text)
	return s.SpeakToAudioWithVoice(text, voice)
}

// SpeakToAudioWithVoice generates speech with specified voice and returns audio data
func (s *TTSService) SpeakToAudioWithVoice(text, voice string) ([]byte, error) {
	// Create temporary file with unique name to avoid conflicts
	tempFile, err := os.CreateTemp(os.TempDir(), "tts-*.aiff")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close() // Close immediately, we just need the path

	// Ensure cleanup
	defer os.Remove(tempPath)

	// Generate speech to temp file
	if voice != "" {
		err = s.SpeakToFileWithVoice(text, tempPath, voice)
	} else {
		err = s.SpeakToFile(text, tempPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate speech: %w", err)
	}

	// Read the audio file
	audioData, err := os.ReadFile(tempPath)
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

// initDetector initializes the language detector (lazy initialization)
func (s *TTSService) initDetector() {
	s.detectorOnce.Do(func() {
		// Build detector with common languages for better performance
		languages := []lingua.Language{
			lingua.English,
			lingua.Chinese,
			lingua.Indonesian,
			lingua.Japanese,
			lingua.Korean,
			lingua.Spanish,
			lingua.French,
			lingua.German,
			lingua.Italian,
			lingua.Russian,
			lingua.Arabic,
		}
		s.detector = lingua.NewLanguageDetectorBuilder().
			FromLanguages(languages...).
			Build()
	})
}

// DetectLanguage detects the language of the given text
// Returns locale code (e.g., "en-US", "zh-CN") and confidence score
func (s *TTSService) DetectLanguage(text string) (string, float64) {
	s.initDetector()

	if text == "" {
		return "en-US", 0.0
	}

	// Detect language with confidence
	if language, exists := s.detector.DetectLanguageOf(text); exists {
		confidence := s.detector.ComputeLanguageConfidence(text, language)

		// Map lingua.Language to locale code
		langCode := s.mapLanguageToLocale(language)

		fmt.Printf("[TTS] Detected language: %s (confidence: %.2f)\n", langCode, confidence)
		return langCode, confidence
	}

	// Default to English if detection fails
	fmt.Printf("[TTS] Language detection failed, defaulting to en-US\n")
	return "en-US", 0.0
}

// DetectLanguageCode detects the language and returns only the language code
// This is a convenience method for cases where confidence score is not needed
func (s *TTSService) DetectLanguageCode(text string) string {
	langCode, _ := s.DetectLanguage(text)
	return langCode
}

// mapLanguageToLocale maps lingua.Language to locale code
func (s *TTSService) mapLanguageToLocale(lang lingua.Language) string {
	mapping := map[lingua.Language]string{
		lingua.English:    "en-US",
		lingua.Chinese:    "zh-CN",
		lingua.Indonesian: "id-ID",
		lingua.Japanese:   "ja-JP",
		lingua.Korean:     "ko-KR",
		lingua.Spanish:    "es-ES",
		lingua.French:     "fr-FR",
		lingua.German:     "de-DE",
		lingua.Italian:    "it-IT",
		lingua.Russian:    "ru-RU",
		lingua.Arabic:     "ar-SA",
	}

	if locale, ok := mapping[lang]; ok {
		return locale
	}

	return "en-US" // Default fallback
}

// SelectVoiceForText detects language and selects appropriate voice
func (s *TTSService) SelectVoiceForText(text string) string {
	langCode, confidence := s.DetectLanguage(text)

	// Use recommended voice for detected language
	recommendedVoices := s.GetRecommendedVoices()
	if voice, ok := recommendedVoices[langCode]; ok {
		fmt.Printf("[TTS] Selected voice: %s for language: %s (confidence: %.2f)\n",
			voice, langCode, confidence)
		return voice
	}

	// Fallback to default voice
	fmt.Printf("[TTS] No recommended voice for %s, using default\n", langCode)
	return ""
}
