package audio_recorder

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AudioRecorderService provides native audio recording capabilities
type AudioRecorderService struct {
	app            *application.App
	recording      bool
	recordingProc  *exec.Cmd
	outputPath     string
	availableTools []string // Cached list of available recording tools
	mu             sync.Mutex
}

// NewAudioRecorderService creates a new audio recorder service
// Automatically installs recording tool if not found
func NewAudioRecorderService(app *application.App) *AudioRecorderService {
	service := &AudioRecorderService{
		app:       app,
		recording: false,
	}

	// Start background initialization
	go service.initializeInBackground()

	return service
}

// initializeInBackground handles recording tool installation
func (s *AudioRecorderService) initializeInBackground() {
	// Check which recording tools are available
	tools := checkAvailableRecordingTools()

	s.mu.Lock()
	s.availableTools = tools
	s.mu.Unlock()

	if len(tools) > 0 {
		log.Printf("✅ Audio recording ready with tools: %v", tools)
		return
	}

	// No tools found, attempt auto-installation
	log.Printf("⚠️  No recording tools found, attempting auto-installation...")
	if installErr := installPlatformRecordingTool(); installErr != nil {
		log.Printf("⚠️  Failed to auto-install recording tool: %v", installErr)
		log.Printf("   Audio recording will not be available until tool is installed")
		return
	}

	// Verify installation
	tools = checkAvailableRecordingTools()
	s.mu.Lock()
	s.availableTools = tools
	s.mu.Unlock()

	if len(tools) == 0 {
		log.Printf("⚠️  Recording tool installation verification failed")
		return
	}

	log.Printf("✅ Audio recording ready with tools: %v", tools)
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

	// Check if any recording tools are available
	if len(s.availableTools) == 0 {
		return "", fmt.Errorf("no recording tools available. Please install one of the supported tools")
	}

	// Create temp file for output
	tempDir := os.TempDir()
	outputPath := filepath.Join(tempDir, fmt.Sprintf("recording_%d.wav", os.Getpid()))

	// Try each available tool in order until one works
	var cmd *exec.Cmd
	var lastErr error

	for _, tool := range s.availableTools {
		cmd, lastErr = startPlatformRecording(tool, outputPath)
		if lastErr == nil {
			log.Printf("Started recording with %s", tool)
			break
		}
		log.Printf("Failed to start recording with %s: %v", tool, lastErr)
	}

	if cmd == nil {
		return "", fmt.Errorf("failed to start recording with any available tool: %w", lastErr)
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
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ [PANIC] Recording process wait panic recovered: %v", r)
					done <- fmt.Errorf("panic: %v", r)
				}
			}()
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
	s.mu.Lock()
	tools := s.availableTools
	s.mu.Unlock()

	result := map[string]interface{}{
		"supported": len(tools) > 0,
		"tools":     tools,
		"error":     "",
	}

	if len(tools) == 0 {
		result["error"] = "no recording tools available"
	}

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
