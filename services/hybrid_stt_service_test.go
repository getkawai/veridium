package services

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestHybridSTTService_Basic(t *testing.T) {
	service, err := NewHybridSTTService("en-US")
	if err != nil {
		t.Fatalf("Failed to create Hybrid STT service: %v", err)
	}
	defer service.Close()

	t.Run("GetAvailableEngines", func(t *testing.T) {
		engines := service.GetAvailableEngines()
		if len(engines) == 0 {
			t.Error("Expected at least one engine")
		}

		t.Logf("Available engines: %v", engines)
	})

	t.Run("GetCurrentEngine", func(t *testing.T) {
		engine := service.GetCurrentEngine()
		if engine == "" {
			t.Error("Expected non-empty engine")
		}

		t.Logf("Current engine: %s", engine)
	})

	t.Run("GetEngineInfo", func(t *testing.T) {
		info := service.GetEngineInfo()
		if info == nil {
			t.Error("Expected engine info")
		}

		t.Logf("Engine info: %+v", info)
	})
}

func TestHybridSTTService_SupportedLocales(t *testing.T) {
	service, err := NewHybridSTTService("en-US")
	if err != nil {
		t.Fatalf("Failed to create Hybrid STT service: %v", err)
	}
	defer service.Close()

	t.Run("GetSupportedLocales", func(t *testing.T) {
		locales := service.GetSupportedLocales()
		if len(locales) == 0 {
			t.Error("Expected supported locales")
		}

		for engine, langs := range locales {
			t.Logf("Engine %s supports %d languages", engine, len(langs))
		}
	})
}

func TestHybridSTTService_Transcribe(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transcription test in short mode")
	}

	service, err := NewHybridSTTService("en-US")
	if err != nil {
		t.Fatalf("Failed to create Hybrid STT service: %v", err)
	}
	defer service.Close()

	// Create test audio file
	tmpDir := t.TempDir()
	testAudioPath := filepath.Join(tmpDir, "test.wav")
	if err := createTestWavFile(testAudioPath); err != nil {
		t.Fatalf("Failed to create test audio: %v", err)
	}

	t.Run("TranscribeAuto", func(t *testing.T) {
		text, err := service.Transcribe(testAudioPath)
		if err != nil {
			t.Fatalf("Failed to transcribe: %v", err)
		}

		t.Logf("Transcription result: %s", text)
	})
}

func TestHybridSTTService_TranscribeWithOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transcription with options test in short mode")
	}

	service, err := NewHybridSTTService("en-US")
	if err != nil {
		t.Fatalf("Failed to create Hybrid STT service: %v", err)
	}
	defer service.Close()

	// Create test audio file
	tmpDir := t.TempDir()
	testAudioPath := filepath.Join(tmpDir, "test.wav")
	if err := createTestWavFile(testAudioPath); err != nil {
		t.Fatalf("Failed to create test audio: %v", err)
	}

	engines := service.GetAvailableEngines()

	for _, engine := range engines {
		t.Run(string(engine), func(t *testing.T) {
			opts := TranscriptionOptions{
				Engine:       engine,
				Locale:       "en-US",
				WhisperModel: "ggml-tiny",
				Timeout:      60 * time.Second,
			}

			text, err := service.TranscribeWithOptions(testAudioPath, opts)
			if err != nil {
				t.Logf("Warning: Failed to transcribe with %s: %v", engine, err)
				return
			}

			t.Logf("✅ %s transcription: %s", engine, text)
		})
	}
}

func TestHybridSTTService_Benchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping benchmark test in short mode")
	}

	service, err := NewHybridSTTService("en-US")
	if err != nil {
		t.Fatalf("Failed to create Hybrid STT service: %v", err)
	}
	defer service.Close()

	// Create test audio file
	tmpDir := t.TempDir()
	testAudioPath := filepath.Join(tmpDir, "test.wav")
	if err := createTestWavFile(testAudioPath); err != nil {
		t.Fatalf("Failed to create test audio: %v", err)
	}

	t.Run("BenchmarkEngines", func(t *testing.T) {
		// Download Whisper model first if needed
		if service.whisperSTT != nil {
			models := service.whisperSTT.ListModels()
			if len(models) == 0 {
				t.Log("Downloading Whisper tiny model...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				defer cancel()

				if err := service.whisperSTT.DownloadModel(ctx, "ggml-tiny.bin"); err != nil {
					t.Logf("Warning: Failed to download model: %v", err)
					return
				}
			}
		}

		results, err := service.Benchmark(testAudioPath)
		if err != nil {
			t.Fatalf("Failed to benchmark: %v", err)
		}

		t.Log("Benchmark Results:")
		for engine, result := range results {
			resultMap := result.(map[string]interface{})
			t.Logf("  %s:", engine)
			t.Logf("    Duration: %.2fs", resultMap["duration"])
			t.Logf("    Success: %v", resultMap["success"])
			if resultMap["success"].(bool) {
				t.Logf("    Text: %s", resultMap["text"])
			} else if resultMap["error"] != nil {
				t.Logf("    Error: %v", resultMap["error"])
			}
		}
	})
}

func TestHybridSTTService_SetLocale(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping locale change test on non-macOS platform")
	}

	service, err := NewHybridSTTService("en-US")
	if err != nil {
		t.Fatalf("Failed to create Hybrid STT service: %v", err)
	}
	defer service.Close()

	t.Run("ChangeLocale", func(t *testing.T) {
		err := service.SetLocale("ja-JP")
		if err != nil {
			t.Fatalf("Failed to set locale: %v", err)
		}

		if service.locale != "ja-JP" {
			t.Errorf("Expected locale ja-JP, got %s", service.locale)
		}

		t.Log("✅ Locale changed to ja-JP")
	})
}

func TestHybridSTTService_MultiLanguage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping multi-language test in short mode")
	}

	if runtime.GOOS != "darwin" {
		t.Skip("Skipping multi-language test on non-macOS platform")
	}

	// Test different languages
	languages := []struct {
		code string
		name string
	}{
		{"en-US", "English"},
		{"ja-JP", "Japanese"},
		{"zh-CN", "Chinese"},
		{"id-ID", "Indonesian"},
	}

	for _, lang := range languages {
		t.Run(lang.name, func(t *testing.T) {
			service, err := NewHybridSTTService(lang.code)
			if err != nil {
				t.Logf("Warning: Failed to create service for %s: %v", lang.name, err)
				return
			}
			defer service.Close()

			engines := service.GetAvailableEngines()
			t.Logf("✅ %s (%s): Available engines: %v", lang.name, lang.code, engines)
		})
	}
}

// Note: createTestWavFile is defined in whisper_service_test.go
