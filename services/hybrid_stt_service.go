package services

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// HybridSTTService provides speech-to-text using both native and Whisper engines
// It automatically selects the best engine based on platform and availability
type HybridSTTService struct {
	nativeSTT  *NativeSTTService
	whisperSTT *WhisperService
	useNative  bool
	locale     string
}

// STTEngine represents the engine used for transcription
type STTEngine string

const (
	EngineNative  STTEngine = "native"  // macOS Speech Framework
	EngineWhisper STTEngine = "whisper" // Whisper.cpp
	EngineAuto    STTEngine = "auto"    // Automatic selection
)

// TranscriptionOptions contains options for transcription
type TranscriptionOptions struct {
	Engine       STTEngine // Preferred engine (auto, native, whisper)
	Locale       string    // Language code (e.g., "en-US", "id-ID")
	WhisperModel string    // Whisper model to use (if engine is whisper)
	Timeout      time.Duration
}

// NewHybridSTTService creates a new hybrid STT service
func NewHybridSTTService(locale string) (*HybridSTTService, error) {
	if locale == "" {
		locale = "en-US"
	}
	
	service := &HybridSTTService{
		locale: locale,
	}
	
	// Try to initialize native STT (macOS only)
	if runtime.GOOS == "darwin" {
		native, err := NewNativeSTTService(locale)
		if err == nil && native.IsAvailable() {
			service.nativeSTT = native
			service.useNative = true
		}
	}
	
	// Always initialize Whisper as fallback
	whisper, err := NewWhisperService()
	if err != nil {
		// If both native and whisper fail, return error
		if service.nativeSTT == nil {
			return nil, fmt.Errorf("failed to initialize any STT engine: %w", err)
		}
		// Native is available, continue without Whisper
	} else {
		service.whisperSTT = whisper
	}
	
	return service, nil
}

// Transcribe transcribes audio file using the best available engine
func (s *HybridSTTService) Transcribe(audioPath string) (string, error) {
	return s.TranscribeWithOptions(audioPath, TranscriptionOptions{
		Engine:       EngineAuto,
		Locale:       s.locale,
		WhisperModel: "ggml-base",
		Timeout:      60 * time.Second,
	})
}

// TranscribeWithOptions transcribes audio with specific options
func (s *HybridSTTService) TranscribeWithOptions(audioPath string, opts TranscriptionOptions) (string, error) {
	// Determine which engine to use
	engine := opts.Engine
	if engine == EngineAuto {
		engine = s.selectBestEngine()
	}
	
	// Try selected engine first
	switch engine {
	case EngineNative:
		if s.nativeSTT != nil {
			text, err := s.transcribeWithNative(audioPath)
			if err == nil {
				return text, nil
			}
			// Fallback to Whisper if native fails
			if s.whisperSTT != nil {
				return s.transcribeWithWhisper(audioPath, opts.WhisperModel, opts.Timeout)
			}
			return "", err
		}
		// Native not available, try Whisper
		if s.whisperSTT != nil {
			return s.transcribeWithWhisper(audioPath, opts.WhisperModel, opts.Timeout)
		}
		return "", fmt.Errorf("no STT engine available")
		
	case EngineWhisper:
		if s.whisperSTT != nil {
			return s.transcribeWithWhisper(audioPath, opts.WhisperModel, opts.Timeout)
		}
		// Whisper not available, try native
		if s.nativeSTT != nil {
			return s.transcribeWithNative(audioPath)
		}
		return "", fmt.Errorf("whisper engine not available")
		
	default:
		return "", fmt.Errorf("unknown engine: %s", engine)
	}
}

// transcribeWithNative uses native macOS Speech Framework
func (s *HybridSTTService) transcribeWithNative(audioPath string) (string, error) {
	if s.nativeSTT == nil {
		return "", fmt.Errorf("native STT not available")
	}
	
	return s.nativeSTT.TranscribeFile(audioPath)
}

// transcribeWithWhisper uses Whisper.cpp
func (s *HybridSTTService) transcribeWithWhisper(audioPath, model string, timeout time.Duration) (string, error) {
	if s.whisperSTT == nil {
		return "", fmt.Errorf("whisper STT not available")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.whisperSTT.Transcribe(ctx, model, audioPath)
}

// selectBestEngine automatically selects the best engine
func (s *HybridSTTService) selectBestEngine() STTEngine {
	// Prefer native on macOS for speed and quality
	if s.useNative && s.nativeSTT != nil && s.nativeSTT.IsAvailable() {
		return EngineNative
	}
	
	// Fallback to Whisper
	if s.whisperSTT != nil {
		return EngineWhisper
	}
	
	return EngineNative // Will fail gracefully if not available
}

// GetAvailableEngines returns list of available engines
func (s *HybridSTTService) GetAvailableEngines() []STTEngine {
	engines := make([]STTEngine, 0, 2)
	
	if s.nativeSTT != nil && s.nativeSTT.IsAvailable() {
		engines = append(engines, EngineNative)
	}
	
	if s.whisperSTT != nil {
		engines = append(engines, EngineWhisper)
	}
	
	return engines
}

// GetEngineInfo returns information about available engines
func (s *HybridSTTService) GetEngineInfo() map[string]interface{} {
	info := map[string]interface{}{
		"platform": runtime.GOOS,
		"locale":   s.locale,
		"engines":  make(map[string]interface{}),
	}
	
	engines := info["engines"].(map[string]interface{})
	
	if s.nativeSTT != nil {
		engines["native"] = map[string]interface{}{
			"available": s.nativeSTT.IsAvailable(),
			"engine":    "macOS Speech Framework",
			"quality":   "excellent",
			"speed":     "fast",
			"offline":   "partial",
			"realtime":  true,
		}
	}
	
	if s.whisperSTT != nil {
		engines["whisper"] = map[string]interface{}{
			"available":     true,
			"engine":        "Whisper.cpp",
			"quality":       "excellent",
			"speed":         "medium",
			"offline":       "full",
			"realtime":      false,
			"models_dir":    s.whisperSTT.GetModelsDirectory(),
			"models_count":  len(s.whisperSTT.ListModels()),
		}
	}
	
	return info
}

// GetSupportedLocales returns supported locales for all engines
func (s *HybridSTTService) GetSupportedLocales() map[string][]string {
	locales := make(map[string][]string)
	
	if s.nativeSTT != nil {
		nativeLocales, err := s.nativeSTT.GetSupportedLocales()
		if err == nil {
			locales["native"] = nativeLocales
		}
	}
	
	// Whisper supports 99 languages
	if s.whisperSTT != nil {
		locales["whisper"] = []string{
			"en", "zh", "de", "es", "ru", "ko", "fr", "ja", "pt", "tr",
			"pl", "ca", "nl", "ar", "sv", "it", "id", "hi", "fi", "vi",
			"he", "uk", "el", "ms", "cs", "ro", "da", "hu", "ta", "no",
			"th", "ur", "hr", "bg", "lt", "la", "mi", "ml", "cy", "sk",
			"te", "fa", "lv", "bn", "sr", "az", "sl", "kn", "et", "mk",
			"br", "eu", "is", "hy", "ne", "mn", "bs", "kk", "sq", "sw",
			"gl", "mr", "pa", "si", "km", "sn", "yo", "so", "af", "oc",
			"ka", "be", "tg", "sd", "gu", "am", "yi", "lo", "uz", "fo",
			"ht", "ps", "tk", "nn", "mt", "sa", "lb", "my", "bo", "tl",
			"mg", "as", "tt", "haw", "ln", "ha", "ba", "jw", "su",
		}
	}
	
	return locales
}

// SetLocale changes the locale for native STT
func (s *HybridSTTService) SetLocale(locale string) error {
	if s.nativeSTT != nil {
		// Need to recreate native STT with new locale
		s.nativeSTT.Close()
		
		native, err := NewNativeSTTService(locale)
		if err != nil {
			return err
		}
		
		s.nativeSTT = native
	}
	
	s.locale = locale
	return nil
}

// GetCurrentEngine returns the engine that would be used for transcription
func (s *HybridSTTService) GetCurrentEngine() STTEngine {
	return s.selectBestEngine()
}

// Close releases all resources
func (s *HybridSTTService) Close() error {
	var err error
	
	if s.nativeSTT != nil {
		if closeErr := s.nativeSTT.Close(); closeErr != nil {
			err = closeErr
		}
	}
	
	if s.whisperSTT != nil {
		if closeErr := s.whisperSTT.Close(); closeErr != nil {
			err = closeErr
		}
	}
	
	return err
}

// Benchmark compares performance of both engines
func (s *HybridSTTService) Benchmark(audioPath string) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	
	// Benchmark native
	if s.nativeSTT != nil && s.nativeSTT.IsAvailable() {
		start := time.Now()
		text, err := s.transcribeWithNative(audioPath)
		duration := time.Since(start)
		
		results["native"] = map[string]interface{}{
			"duration": duration.Seconds(),
			"success":  err == nil,
			"text":     text,
			"error":    err,
		}
	}
	
	// Benchmark whisper
	if s.whisperSTT != nil {
		start := time.Now()
		text, err := s.transcribeWithWhisper(audioPath, "ggml-base", 60*time.Second)
		duration := time.Since(start)
		
		results["whisper"] = map[string]interface{}{
			"duration": duration.Seconds(),
			"success":  err == nil,
			"text":     text,
			"error":    err,
		}
	}
	
	return results, nil
}

