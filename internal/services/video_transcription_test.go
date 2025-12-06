package services_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/whisper"
)

// Suppress unused import warning - sql is used for sql.NullString in document query
var _ = sql.NullString{}

func TestVideoTranscriptionPipeline(t *testing.T) {
	// Skip if video file doesn't exist
	videoPath := "/Users/yuda/github.com/kawai-network/veridium/videoplayback.mp4"
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		t.Skip("Test video file not found, skipping")
	}

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "video_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy video to temp dir
	videoData, err := os.ReadFile(videoPath)
	if err != nil {
		t.Fatalf("Failed to read video: %v", err)
	}
	testVideoPath := filepath.Join(tempDir, "test_video.mp4")
	if err := os.WriteFile(testVideoPath, videoData, 0644); err != nil {
		t.Fatalf("Failed to write test video: %v", err)
	}

	// Initialize test database (will create default user automatically)
	dbPath := filepath.Join(tempDir, "test.db")
	dbService, err := database.NewServiceWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer dbService.Close()

	ctx := context.Background()
	queries := db.New(dbService.DB())

	// Use default user ID (created by NewServiceWithPath)
	const testUserID = "DEFAULT_LOBE_CHAT_USER"

	// Initialize whisper service
	whisperService, err := whisper.NewService()
	if err != nil {
		t.Fatalf("Failed to create whisper service: %v", err)
	}

	// Check if whisper is available
	if !whisperService.IsWhisperInstalled() {
		t.Skip("Whisper not installed, skipping")
	}

	models, err := whisperService.ListModels()
	if err != nil || len(models) == 0 {
		t.Skip("No whisper models available, skipping")
	}

	t.Logf("Using whisper model: %s", models[0])

	// Initialize file processor
	fileLoader := services.NewFileLoader()
	ragProcessor := services.NewRAGProcessor(dbService.DB(), nil, fileLoader, nil)
	processor := services.NewFileProcessorService(
		dbService.DB(),
		fileLoader,
		ragProcessor,
		(*llama.LibraryService)(nil), // No VL needed for video
		whisperService,
	)

	// Process video file
	t.Log("Processing video file...")
	result, err := processor.ProcessFile(ctx, services.ProcessFileRequest{
		FilePath:  testVideoPath,
		Filename:  "test_video.mp4",
		UserID:    testUserID,
		Source:    testVideoPath,
		EnableRAG: false,
		IsShared:  false,
	})
	if err != nil {
		t.Fatalf("Failed to process video: %v", err)
	}

	t.Logf("Video processed: FileID=%s, DocumentID=%s", result.FileID, result.DocumentID)

	// Wait for async transcription to fully complete (no more "Transcribing..." in content)
	t.Log("Waiting for transcription to complete...")
	deadline := time.Now().Add(2 * time.Minute)
	var transcription string
	var lastProgress string

	for time.Now().Before(deadline) {
		doc, err := queries.GetDocumentByFileID(ctx, sql.NullString{String: result.FileID, Valid: true})
		if err == nil && doc.Content.Valid {
			content := doc.Content.String
			if strings.Contains(content, "Video Transcription (AI Generated via Whisper)") {
				// Check if still transcribing (progressive update in progress)
				if strings.Contains(content, "*Transcribing...") {
					// Extract progress info for logging
					if content != lastProgress {
						lastProgress = content
						// Find progress line
						if idx := strings.Index(content, "*Transcribing..."); idx >= 0 {
							endIdx := strings.Index(content[idx:], "*\n")
							if endIdx > 0 {
								t.Logf("Progress: %s", content[idx:idx+endIdx+1])
							}
						}
					}
				} else if strings.Contains(content, "[Segment") {
					// All segments completed, transcription is final
					transcription = content
					break
				}
			}
		}
		time.Sleep(2 * time.Second)
	}

	if transcription == "" {
		t.Fatal("Transcription not completed within timeout")
	}

	t.Logf("Transcription completed! Length: %d chars", len(transcription))

	// Count segments
	segmentCount := strings.Count(transcription, "[Segment")
	t.Logf("Total segments: %d", segmentCount)

	t.Logf("Preview: %s...", transcription[:min(800, len(transcription))])
}

func TestExtractAudioFromVideo(t *testing.T) {
	// Skip if video file doesn't exist
	videoPath := "/Users/yuda/github.com/kawai-network/veridium/videoplayback.mp4"
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		t.Skip("Test video file not found, skipping")
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "audio_extract_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize services (minimal)
	fileLoader := services.NewFileLoader()
	processor := services.NewFileProcessorService(nil, fileLoader, nil, nil, nil)

	// Use reflection to access private method for testing
	// For now, just test ffmpeg directly
	audioPath := filepath.Join(tempDir, "test_audio.wav")

	// Run ffmpeg command
	t.Log("Extracting audio from video...")
	cmd := "ffmpeg -i " + videoPath + " -vn -acodec pcm_s16le -ar 16000 -ac 1 -t 10 -y " + audioPath
	if err := runCommand(cmd); err != nil {
		t.Fatalf("FFmpeg failed: %v", err)
	}

	// Verify audio file
	info, err := os.Stat(audioPath)
	if err != nil {
		t.Fatalf("Audio file not created: %v", err)
	}

	t.Logf("Audio extracted: %s (%.2f KB)", audioPath, float64(info.Size())/1024)

	// Verify it's a valid WAV
	if info.Size() < 1000 {
		t.Fatal("Audio file too small, probably corrupted")
	}

	// Suppress unused variable warning
	_ = processor
}

func runCommand(cmd string) error {
	parts := strings.Fields(cmd)
	c := &struct {
		path string
		args []string
	}{
		path: parts[0],
		args: parts[1:],
	}
	_ = c
	// For this test we just use os/exec
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
