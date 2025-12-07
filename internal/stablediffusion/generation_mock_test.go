package stablediffusion

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// MockCommandExecutor allows testing without running actual binary
type MockCommandExecutor struct {
	LastCmd    string
	LastArgs   []string
	ShouldFail bool
}

func (m *MockCommandExecutor) Run(name string, args ...string) error {
	m.LastCmd = name
	m.LastArgs = args

	if m.ShouldFail {
		return os.ErrNotExist // Just a dummy error
	}

	// Determine output path from args to simulate file creation
	var outputPath string
	for i, arg := range args {
		if arg == "-o" && i+1 < len(args) {
			outputPath = args[i+1]
			break
		}
	}

	if outputPath != "" {
		// Create dummy output file
		return os.WriteFile(outputPath, []byte("dummy image data"), 0644)
	}

	return nil
}

func TestGenerateImage_Mock(t *testing.T) {
	mockExecutor := &MockCommandExecutor{}

	manager := NewStableDiffusionReleaseManager()
	manager.Executor = mockExecutor

	// Create temp output file path
	tmpOutput := "test_output.png"
	defer os.Remove(tmpOutput)

	opts := GenerationOptions{
		Prompt:     "A cute cat",
		ModelPath:  "models/sd-turbo.safetensors",
		OutputPath: tmpOutput,
		Width:      512,
		Height:     512,
		Steps:      10,
	}

	// Mock binary existence check (GenerateImage checks os.Stat)
	// We need to ensure manager.GetBinaryPath() points to something that exists
	// or mock that check too?
	// The current logic checks os.Stat(binaryPath) BEFORE calling Executor.Run.
	// So we must have a dummy binary file.

	// Create dummy binary
	dummyBinDir := t.TempDir()
	manager.BinaryPath = dummyBinDir

	// Create the binary file that GetBinaryPath() expects ("sd" or "sd.exe")
	binName := "sd"
	if runtime.GOOS == "windows" {
		binName = "sd.exe"
	}
	realBinPath := filepath.Join(dummyBinDir, binName)
	os.WriteFile(realBinPath, []byte("dummy binary"), 0755)

	err := manager.GenerateImage(opts)
	if err != nil {
		t.Fatalf("GenerateImage failed with mock: %v", err)
	}

	// Verify args
	// GenerateImage calls GetBinaryPath() which joins BinaryPath with "sd"/"sd.exe"
	if mockExecutor.LastCmd != realBinPath {
		t.Errorf("Expected command %s, got %s", realBinPath, mockExecutor.LastCmd)
	}

	foundPrompt := false
	for _, arg := range mockExecutor.LastArgs {
		if arg == opts.Prompt {
			foundPrompt = true
			break
		}
	}

	if !foundPrompt {
		t.Errorf("Prompt not found in args: %v", mockExecutor.LastArgs)
	}

	// Verify output file creation simulated by mock
	if _, err := os.Stat(tmpOutput); os.IsNotExist(err) {
		t.Errorf("Output file not created by mock")
	}
}
