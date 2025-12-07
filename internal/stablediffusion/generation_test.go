package stablediffusion

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper to mock exec.Command
// Since we can't easily mock exec.Command in Go without changing the internal code structure significantly
// (e.g. using an interface for Command execution), for this test we will focus on
// 1. Validating the arguments can be prepared
// 2. Integration test style: skipping if binary not found

func TestGenerateImage_Integration(t *testing.T) {
	// This test requires the binary to be present, otherwise it skips

	// Create a temp directory for outputs
	tmpDir, err := os.MkdirTemp("", "sd-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewStableDiffusionReleaseManager()

	// Check if binary exists, if not skip
	if !manager.IsStableDiffusionInstalled() {
		t.Skip("Stable Diffusion binary not installed, skipping integration test")
	}

	// Check if models exist
	models, err := manager.CheckInstalledModels()
	if err != nil || len(models) == 0 {
		t.Skip("No Stable Diffusion models installed, skipping integration test")
	}

	t.Logf("Found models: %v", models)

	// We need full path for the model usually, let's try to find it
	modelPath := ""
	modelsPath := manager.GetModelsPath()
	entries, _ := os.ReadDir(modelsPath)
	for _, e := range entries {
		if !e.IsDir() { // Simple check, assuming the first model found is valid
			modelPath = filepath.Join(modelsPath, e.Name())
			break
		}
	}

	if modelPath == "" {
		t.Skip("Could not determine model path")
	}

	outputPath := filepath.Join(tmpDir, "test_output.png")

	opts := GenerationOptions{
		Prompt:     "a small red cube",
		ModelPath:  modelPath,
		OutputPath: outputPath,
		Width:      64, // Small for speed if possible, though SD usually wants 512+
		Height:     64,
		Steps:      1, // Minimal steps for speed
		Seed:       12345,
	}

	err = manager.GenerateImage(opts)
	if err != nil {
		t.Fatalf("GenerateImage failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created at %s", outputPath)
	}
}

// TestGenerateImage_Structure validates that we can call the method and it checks for binary existence properly
func TestGenerateImage_BinaryCheck(t *testing.T) {
	// Setup a dummy manager pointing to non-existent binary
	manager := &StableDiffusionReleaseManager{
		BinaryPath: "/path/to/non/existent/binary",
	}

	opts := GenerationOptions{
		Prompt:     "test",
		ModelPath:  "test.safetensors",
		OutputPath: "out.png",
	}

	err := manager.GenerateImage(opts)
	if err == nil {
		t.Error("Expected error when binary is missing, got nil")
	}
}
