package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWhisperService_Basic(t *testing.T) {
	// Create temporary directory for test models
	tmpDir := t.TempDir()

	// Create WhisperService with temp directory
	service := &WhisperService{
		modelsDir: tmpDir,
	}

	// Initialize whisper
	var err error
	service.whisper, err = initWhisper(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize whisper: %v", err)
	}
	defer service.Close()

	t.Run("GetModelsDirectory", func(t *testing.T) {
		dir := service.GetModelsDirectory()
		if dir != tmpDir {
			t.Errorf("Expected models directory %s, got %s", tmpDir, dir)
		}
	})

	t.Run("ListModels_Empty", func(t *testing.T) {
		models := service.ListModels()
		if len(models) != 0 {
			t.Errorf("Expected 0 models, got %d", len(models))
		}
	})

	t.Run("GetAvailableModels", func(t *testing.T) {
		available := service.GetAvailableModels()
		if len(available) == 0 {
			t.Error("Expected some available models")
		}

		// Check structure of first model
		if len(available) > 0 {
			model := available[0]
			if _, ok := model["id"]; !ok {
				t.Error("Model should have 'id' field")
			}
			if _, ok := model["name"]; !ok {
				t.Error("Model should have 'name' field")
			}
			if _, ok := model["size"]; !ok {
				t.Error("Model should have 'size' field")
			}
		}
	})
}

func TestWhisperService_DownloadModel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping model download test in short mode")
	}

	// Create temporary directory for test models
	tmpDir := t.TempDir()

	// Create WhisperService
	service := &WhisperService{
		modelsDir: tmpDir,
	}

	var err error
	service.whisper, err = initWhisper(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize whisper: %v", err)
	}
	defer service.Close()

	t.Run("DownloadTinyModel", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Download tiny model (smallest, ~75MB)
		err := service.DownloadModel(ctx, "ggml-tiny.bin")
		if err != nil {
			t.Fatalf("Failed to download model: %v", err)
		}

		// Verify model was downloaded
		models := service.ListModels()
		if len(models) == 0 {
			t.Error("Expected at least 1 model after download")
		}

		// Check if tiny model exists
		found := false
		for _, m := range models {
			if m.Id == "ggml-tiny" || m.Id == "ggml-tiny.bin" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Downloaded model not found in list")
		}
	})
}

func TestWhisperService_Transcribe(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transcription test in short mode")
	}

	// Create temporary directory for test models
	tmpDir := t.TempDir()

	// Create WhisperService
	service := &WhisperService{
		modelsDir: tmpDir,
	}

	var err error
	service.whisper, err = initWhisper(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize whisper: %v", err)
	}
	defer service.Close()

	// Download tiny model first
	ctx := context.Background()
	downloadCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	t.Log("Downloading tiny model...")
	err = service.DownloadModel(downloadCtx, "ggml-tiny.bin")
	if err != nil {
		t.Fatalf("Failed to download model: %v", err)
	}

	// Create a simple test WAV file
	testWavPath := filepath.Join(tmpDir, "test.wav")
	err = createTestWavFile(testWavPath)
	if err != nil {
		t.Fatalf("Failed to create test WAV file: %v", err)
	}

	t.Run("TranscribeSimple", func(t *testing.T) {
		transcribeCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		text, err := service.Transcribe(transcribeCtx, "ggml-tiny", testWavPath)
		if err != nil {
			t.Fatalf("Failed to transcribe: %v", err)
		}

		if text == "" {
			t.Error("Expected non-empty transcription")
		}

		t.Logf("Transcription result: %s", text)
	})

	t.Run("TranscribeWithSegments", func(t *testing.T) {
		transcribeCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		segments, err := service.TranscribeWithSegments(transcribeCtx, "ggml-tiny", testWavPath)
		if err != nil {
			t.Fatalf("Failed to transcribe with segments: %v", err)
		}

		if len(segments) == 0 {
			t.Error("Expected at least one segment")
		}

		for i, seg := range segments {
			t.Logf("Segment %d: [%v - %v] %s", i, seg.Start, seg.End, seg.Text)
		}
	})
}

func TestWhisperService_ModelManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping model management test in short mode")
	}

	// Create temporary directory for test models
	tmpDir := t.TempDir()

	// Create WhisperService
	service := &WhisperService{
		modelsDir: tmpDir,
	}

	var err error
	service.whisper, err = initWhisper(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize whisper: %v", err)
	}
	defer service.Close()

	// Download tiny model
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err = service.DownloadModel(ctx, "ggml-tiny.bin")
	if err != nil {
		t.Fatalf("Failed to download model: %v", err)
	}

	t.Run("GetModel", func(t *testing.T) {
		model := service.GetModel("ggml-tiny")
		if model == nil {
			t.Error("Expected to find ggml-tiny model")
		}
	})

	t.Run("DeleteModel", func(t *testing.T) {
		// Get model ID first
		models := service.ListModels()
		if len(models) == 0 {
			t.Fatal("No models to delete")
		}

		modelId := models[0].Id
		err := service.DeleteModel(modelId)
		if err != nil {
			t.Errorf("Failed to delete model: %v", err)
		}

		// Verify deletion
		model := service.GetModel(modelId)
		if model != nil {
			t.Error("Model should be deleted")
		}
	})
}

// createTestWavFile creates a simple silent WAV file for testing
func createTestWavFile(path string) error {
	// Create a simple 1-second silent WAV file (16kHz, mono, 16-bit)
	sampleRate := 16000
	duration := 1 // seconds
	numSamples := sampleRate * duration

	// WAV header
	header := []byte{
		// RIFF header
		'R', 'I', 'F', 'F',
		0, 0, 0, 0, // File size (will be filled)
		'W', 'A', 'V', 'E',

		// fmt chunk
		'f', 'm', 't', ' ',
		16, 0, 0, 0, // fmt chunk size
		1, 0, // Audio format (1 = PCM)
		1, 0, // Number of channels (1 = mono)
		0x80, 0x3e, 0, 0, // Sample rate (16000 Hz)
		0, 0x7d, 0, 0, // Byte rate (16000 * 1 * 2)
		2, 0, // Block align (1 * 2)
		16, 0, // Bits per sample (16)

		// data chunk
		'd', 'a', 't', 'a',
		0, 0, 0, 0, // Data size (will be filled)
	}

	// Calculate sizes
	dataSize := numSamples * 2             // 2 bytes per sample (16-bit)
	fileSize := len(header) + dataSize - 8 // -8 for RIFF header

	// Fill in sizes
	header[4] = byte(fileSize)
	header[5] = byte(fileSize >> 8)
	header[6] = byte(fileSize >> 16)
	header[7] = byte(fileSize >> 24)

	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	// Create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header
	if _, err := f.Write(header); err != nil {
		return err
	}

	// Write silent audio data (all zeros)
	silence := make([]byte, dataSize)
	if _, err := f.Write(silence); err != nil {
		return err
	}

	return nil
}
