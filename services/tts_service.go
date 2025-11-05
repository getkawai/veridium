package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// TTSService provides text-to-speech using native OS capabilities
type TTSService struct {
	platform string
}

// Voice represents a TTS voice
type Voice struct {
	Name     string
	Language string
	Gender   string
}

// NewTTSService creates a new TTS service instance
func NewTTSService() (*TTSService, error) {
	return &TTSService{
		platform: runtime.GOOS,
	}, nil
}

// Speak converts text to speech using default voice
func (s *TTSService) Speak(text string) error {
	switch s.platform {
	case "darwin":
		return s.speakMacOS(text, "")
	case "windows":
		return s.speakWindows(text)
	case "linux":
		return s.speakLinux(text)
	default:
		return fmt.Errorf("unsupported platform: %s", s.platform)
	}
}

// SpeakWithVoice converts text to speech using specified voice
func (s *TTSService) SpeakWithVoice(text, voice string) error {
	switch s.platform {
	case "darwin":
		return s.speakMacOS(text, voice)
	default:
		return s.Speak(text) // Fallback to default
	}
}

// SpeakWithRate converts text to speech with custom speaking rate
// rate: words per minute (default is ~175 for macOS)
func (s *TTSService) SpeakWithRate(text string, rate int) error {
	if s.platform != "darwin" {
		return s.Speak(text)
	}

	cmd := exec.Command("say", "-r", fmt.Sprintf("%d", rate), text)
	return cmd.Run()
}

// SpeakToFile saves speech to an audio file
// Supported formats: AIFF (macOS default), WAV, etc.
func (s *TTSService) SpeakToFile(text, outputPath string) error {
	switch s.platform {
	case "darwin":
		cmd := exec.Command("say", "-o", outputPath, text)
		return cmd.Run()
	default:
		return fmt.Errorf("SpeakToFile not supported on %s", s.platform)
	}
}

// SpeakToFileWithVoice saves speech to file with specified voice
func (s *TTSService) SpeakToFileWithVoice(text, outputPath, voice string) error {
	if s.platform != "darwin" {
		return s.SpeakToFile(text, outputPath)
	}

	cmd := exec.Command("say", "-v", voice, "-o", outputPath, text)
	return cmd.Run()
}

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

// ListVoices returns available TTS voices
func (s *TTSService) ListVoices() ([]Voice, error) {
	switch s.platform {
	case "darwin":
		return s.listVoicesMacOS()
	default:
		return nil, fmt.Errorf("ListVoices not supported on %s", s.platform)
	}
}

// GetDefaultVoice returns the system default voice
func (s *TTSService) GetDefaultVoice() (string, error) {
	if s.platform != "darwin" {
		return "", fmt.Errorf("GetDefaultVoice not supported on %s", s.platform)
	}

	// Get default voice from system preferences
	cmd := exec.Command("defaults", "read", "com.apple.speech.voice.prefs", "SelectedVoiceName")
	output, err := cmd.Output()
	if err != nil {
		return "Samantha", nil // Fallback to Samantha
	}

	return strings.TrimSpace(string(output)), nil
}

// Stop stops any ongoing speech (macOS only)
func (s *TTSService) Stop() error {
	if s.platform != "darwin" {
		return fmt.Errorf("Stop not supported on %s", s.platform)
	}

	// Kill all running 'say' processes
	cmd := exec.Command("killall", "say")
	return cmd.Run()
}

// Platform-specific implementations

func (s *TTSService) speakMacOS(text, voice string) error {
	args := []string{}

	if voice != "" {
		args = append(args, "-v", voice)
	}

	args = append(args, text)

	cmd := exec.Command("say", args...)
	return cmd.Run()
}

func (s *TTSService) speakWindows(text string) error {
	// Use PowerShell with SAPI
	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Speech; $speak = New-Object System.Speech.Synthesis.SpeechSynthesizer; $speak.Speak('%s')`, text)
	cmd := exec.Command("powershell", "-Command", psScript)
	return cmd.Run()
}

func (s *TTSService) speakLinux(text string) error {
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

func (s *TTSService) listVoicesMacOS() ([]Voice, error) {
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

// IsPlatformSupported returns whether TTS is supported on current platform
func (s *TTSService) IsPlatformSupported() bool {
	switch s.platform {
	case "darwin", "windows", "linux":
		return true
	default:
		return false
	}
}

// GetPlatformInfo returns information about TTS capabilities on current platform
func (s *TTSService) GetPlatformInfo() map[string]interface{} {
	info := map[string]interface{}{
		"platform":  s.platform,
		"supported": s.IsPlatformSupported(),
	}

	switch s.platform {
	case "darwin":
		info["engine"] = "macOS Speech Synthesis (say)"
		info["quality"] = "excellent"
		info["offline"] = true
		info["voices"] = "70+"
	case "windows":
		info["engine"] = "Windows SAPI"
		info["quality"] = "good"
		info["offline"] = true
		info["voices"] = "varies"
	case "linux":
		info["engine"] = "espeak/festival"
		info["quality"] = "basic"
		info["offline"] = true
		info["voices"] = "varies"
	}

	return info
}
