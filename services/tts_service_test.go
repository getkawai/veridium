package services

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestTTSService_Basic(t *testing.T) {
	service, err := NewTTSService()
	if err != nil {
		t.Fatalf("Failed to create TTS service: %v", err)
	}
	
	t.Run("PlatformSupported", func(t *testing.T) {
		if !service.IsPlatformSupported() {
			t.Skip("TTS not supported on this platform")
		}
	})
	
	t.Run("GetPlatformInfo", func(t *testing.T) {
		info := service.GetPlatformInfo()
		if info == nil {
			t.Error("Expected platform info")
		}
		
		if platform, ok := info["platform"]; !ok || platform != runtime.GOOS {
			t.Errorf("Expected platform %s, got %v", runtime.GOOS, platform)
		}
	})
}

func TestTTSService_Speak(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping speak test on non-macOS platform")
	}
	
	service, err := NewTTSService()
	if err != nil {
		t.Fatalf("Failed to create TTS service: %v", err)
	}
	
	t.Run("SpeakSimple", func(t *testing.T) {
		err := service.Speak("Hello, this is a test")
		if err != nil {
			t.Errorf("Failed to speak: %v", err)
		}
	})
	
	t.Run("SpeakWithVoice", func(t *testing.T) {
		err := service.SpeakWithVoice("Testing with Samantha voice", "Samantha")
		if err != nil {
			t.Errorf("Failed to speak with voice: %v", err)
		}
	})
	
	t.Run("SpeakWithRate", func(t *testing.T) {
		err := service.SpeakWithRate("Fast speech test", 250)
		if err != nil {
			t.Errorf("Failed to speak with rate: %v", err)
		}
	})
}

func TestTTSService_SpeakToFile(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping file output test on non-macOS platform")
	}
	
	service, err := NewTTSService()
	if err != nil {
		t.Fatalf("Failed to create TTS service: %v", err)
	}
	
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_output.aiff")
	
	t.Run("SaveToFile", func(t *testing.T) {
		err := service.SpeakToFile("This is a test for file output", outputPath)
		if err != nil {
			t.Fatalf("Failed to save speech to file: %v", err)
		}
		
		// Verify file exists
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Error("Output file was not created")
		}
		
		// Verify file has content
		info, err := os.Stat(outputPath)
		if err != nil {
			t.Fatalf("Failed to stat output file: %v", err)
		}
		
		if info.Size() == 0 {
			t.Error("Output file is empty")
		}
		
		t.Logf("Generated audio file: %s (size: %d bytes)", outputPath, info.Size())
	})
	
	t.Run("SaveToFileWithVoice", func(t *testing.T) {
		outputPath2 := filepath.Join(tmpDir, "test_output_voice.aiff")
		err := service.SpeakToFileWithVoice("Testing with specific voice", outputPath2, "Alex")
		if err != nil {
			t.Fatalf("Failed to save speech with voice to file: %v", err)
		}
		
		// Verify file exists
		if _, err := os.Stat(outputPath2); os.IsNotExist(err) {
			t.Error("Output file was not created")
		}
	})
}

func TestTTSService_ListVoices(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping voice listing test on non-macOS platform")
	}
	
	service, err := NewTTSService()
	if err != nil {
		t.Fatalf("Failed to create TTS service: %v", err)
	}
	
	t.Run("ListAllVoices", func(t *testing.T) {
		voices, err := service.ListVoices()
		if err != nil {
			t.Fatalf("Failed to list voices: %v", err)
		}
		
		if len(voices) == 0 {
			t.Error("Expected at least one voice")
		}
		
		t.Logf("Found %d voices", len(voices))
		
		// Log first few voices
		for i, voice := range voices {
			if i >= 5 {
				break
			}
			t.Logf("Voice %d: %s (%s) - %s", i+1, voice.Name, voice.Language, voice.Gender)
		}
	})
	
	t.Run("GetVoicesByLanguage", func(t *testing.T) {
		voices, err := service.GetVoicesByLanguage("en")
		if err != nil {
			t.Fatalf("Failed to get English voices: %v", err)
		}
		
		if len(voices) == 0 {
			t.Error("Expected at least one English voice")
		}
		
		t.Logf("Found %d English voices", len(voices))
	})
	
	t.Run("GetRecommendedVoices", func(t *testing.T) {
		recommended := service.GetRecommendedVoices()
		if len(recommended) == 0 {
			t.Error("Expected recommended voices")
		}
		
		// Check for common languages
		if _, ok := recommended["en-US"]; !ok {
			t.Error("Expected en-US in recommended voices")
		}
		
		t.Logf("Recommended voices: %v", recommended)
	})
}

func TestTTSService_DefaultVoice(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping default voice test on non-macOS platform")
	}
	
	service, err := NewTTSService()
	if err != nil {
		t.Fatalf("Failed to create TTS service: %v", err)
	}
	
	t.Run("GetDefaultVoice", func(t *testing.T) {
		voice, err := service.GetDefaultVoice()
		if err != nil {
			t.Fatalf("Failed to get default voice: %v", err)
		}
		
		if voice == "" {
			t.Error("Expected non-empty default voice")
		}
		
		t.Logf("Default voice: %s", voice)
	})
}

func TestTTSService_MultiLanguage(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping multi-language test on non-macOS platform")
	}
	
	if testing.Short() {
		t.Skip("Skipping multi-language test in short mode")
	}
	
	service, err := NewTTSService()
	if err != nil {
		t.Fatalf("Failed to create TTS service: %v", err)
	}
	
	tests := []struct {
		name  string
		text  string
		voice string
	}{
		{"English", "Hello, how are you?", "Samantha"},
		{"Indonesian", "Halo, apa kabar?", "Damayanti"},
		{"Japanese", "こんにちは", "Kyoko"},
		{"Chinese", "你好", "Ting-Ting"},
		{"Spanish", "Hola, ¿cómo estás?", "Monica"},
	}
	
	tmpDir := t.TempDir()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(tmpDir, tt.name+".aiff")
			err := service.SpeakToFileWithVoice(tt.text, outputPath, tt.voice)
			if err != nil {
				t.Logf("Warning: Failed to generate %s speech: %v", tt.name, err)
				// Don't fail - voice might not be installed
				return
			}
			
			// Verify file exists
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("%s output file was not created", tt.name)
			} else {
				t.Logf("✅ Generated %s audio: %s", tt.name, outputPath)
			}
		})
	}
}

func TestTTSService_Stop(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping stop test on non-macOS platform")
	}
	
	service, err := NewTTSService()
	if err != nil {
		t.Fatalf("Failed to create TTS service: %v", err)
	}
	
	t.Run("StopSpeech", func(t *testing.T) {
		// Start long speech in background
		go service.Speak("This is a very long text that will take some time to speak completely")
		
		// Stop it immediately
		err := service.Stop()
		if err != nil {
			t.Logf("Warning: Failed to stop speech: %v", err)
			// Don't fail - might not be speaking
		}
	})
}

