package audio_recorder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AudioRecorderService provides native audio recording capabilities
type AudioRecorderService struct {
	app           *application.App
	recording     bool
	recordingProc *exec.Cmd
	outputPath    string
	mu            sync.Mutex
}

// NewAudioRecorderService creates a new audio recorder service
func NewAudioRecorderService(app *application.App) *AudioRecorderService {
	return &AudioRecorderService{
		app:       app,
		recording: false,
	}
}

// SetApp sets the application instance (for event emission)
// This can be called after the app is created if the service was initialized before
func (s *AudioRecorderService) SetApp(app *application.App) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.app = app
}

// StartRecording starts recording audio from the microphone
// Returns the path where the audio will be saved
func (s *AudioRecorderService) StartRecording(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.recording {
		return "", fmt.Errorf("already recording")
	}

	// Create temp file for output
	tempDir := os.TempDir()
	outputPath := filepath.Join(tempDir, fmt.Sprintf("recording_%d.wav", os.Getpid()))

	// Start recording process using platform-specific implementation
	cmd, err := startPlatformRecording(outputPath)
	if err != nil {
		return "", err
	}

	s.recordingProc = cmd
	s.outputPath = outputPath
	s.recording = true

	// Emit event to frontend
	if s.app != nil {
		s.app.Event.Emit("audio:recording:started", outputPath)
	}

	return outputPath, nil
}

// StopRecording stops the current recording
// Returns the path to the recorded audio file
func (s *AudioRecorderService) StopRecording() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.recording {
		return "", fmt.Errorf("not recording")
	}

	outputPath := s.outputPath

	// Stop the recording process using platform-specific implementation
	if s.recordingProc != nil && s.recordingProc.Process != nil {
		// Stop recording gracefully
		if err := stopPlatformRecording(s.recordingProc); err != nil {
			fmt.Printf("Error stopping recording: %v\n", err)
		}

		// Wait for process to finish writing the file
		// Use a channel with timeout to avoid blocking forever
		done := make(chan error, 1)
		go func() {
			done <- s.recordingProc.Wait()
		}()

		select {
		case err := <-done:
			// Process finished
			if err != nil {
				fmt.Printf("Recording process finished with error: %v\n", err)
			}
		case <-time.After(3 * time.Second):
			// Timeout - force kill
			fmt.Printf("Recording process timeout, force killing\n")
			s.recordingProc.Process.Kill()
			s.recordingProc.Wait()
		}

		s.recordingProc = nil
	}

	s.recording = false
	s.outputPath = ""

	// Wait for file to be written and have content
	// sox may take a moment to flush and close the file
	maxRetries := 20                     // Increased from 10
	retryDelay := 200 * time.Millisecond // Increased from 100ms

	fmt.Printf("Waiting for recording file: %s\n", outputPath)

	for i := 0; i < maxRetries; i++ {
		if fileInfo, err := os.Stat(outputPath); err == nil && fileInfo.Size() > 0 {
			// File exists and has content - good to go!
			fmt.Printf("Recording file ready: %s (size: %d bytes)\n", outputPath, fileInfo.Size())
			break
		} else if i == maxRetries-1 {
			// Last retry failed - log details
			if err != nil {
				fmt.Printf("File check failed: %v\n", err)
			} else {
				fmt.Printf("File exists but is empty\n")
			}
			return "", fmt.Errorf("recording file not found or empty after %d retries: %s", maxRetries, outputPath)
		}

		// Wait before retry
		time.Sleep(retryDelay)
	}

	// Emit event to frontend
	if s.app != nil {
		s.app.Event.Emit("audio:recording:stopped", outputPath)
	}

	return outputPath, nil
}

// IsRecording returns whether audio is currently being recorded
func (s *AudioRecorderService) IsRecording() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.recording
}

// GetRecordingPath returns the path to the current recording
func (s *AudioRecorderService) GetRecordingPath() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.outputPath
}

// CheckRecordingCapabilities checks if audio recording is supported
func (s *AudioRecorderService) CheckRecordingCapabilities() map[string]interface{} {
	result := map[string]interface{}{
		"supported": false,
		"tool":      "",
		"error":     "",
	}

	tool, err := checkPlatformRecordingTool()
	if err != nil {
		result["error"] = err.Error()
		return result
	}

	result["supported"] = true
	result["tool"] = tool
	return result
}

// CancelRecording cancels the current recording without saving
func (s *AudioRecorderService) CancelRecording() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.recording {
		return fmt.Errorf("not recording")
	}

	// Kill the recording process
	if s.recordingProc != nil && s.recordingProc.Process != nil {
		s.recordingProc.Process.Kill()
		s.recordingProc.Wait()
		s.recordingProc = nil
	}

	// Remove the output file
	if s.outputPath != "" {
		os.Remove(s.outputPath)
		s.outputPath = ""
	}

	s.recording = false

	// Emit event to frontend
	if s.app != nil {
		s.app.Event.Emit("audio:recording:cancelled", nil)
	}

	return nil
}
