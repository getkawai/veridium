package whisperapp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSetupOptions(t *testing.T) {
	opts := DefaultSetupOptions()

	assert.NotEmpty(t, opts.ModelsDir)
	assert.NotEmpty(t, opts.LibDir)
	assert.Empty(t, opts.DownloadModels)
	assert.False(t, opts.AutoSelect)
	assert.True(t, opts.Interactive)
}

func TestSetupOptions_CustomValues(t *testing.T) {
	opts := &SetupOptions{
		ModelsDir:      "/custom/models",
		LibDir:         "/custom/lib",
		DownloadModels: []string{"base", "small"},
		AutoSelect:     true,
		Interactive:    false,
	}

	assert.Equal(t, "/custom/models", opts.ModelsDir)
	assert.Equal(t, "/custom/lib", opts.LibDir)
	assert.Equal(t, []string{"base", "small"}, opts.DownloadModels)
	assert.True(t, opts.AutoSelect)
	assert.False(t, opts.Interactive)
}

func TestQuickSetup(t *testing.T) {
	// This is an integration test that would download models
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Use temp directory to avoid polluting real data dir
	tempDir := t.TempDir()

	opts := DefaultSetupOptions()
	opts.ModelsDir = tempDir
	opts.LibDir = filepath.Join(tempDir, "lib")
	opts.DownloadModels = []string{"tiny"}
	opts.Interactive = false

	ctx := context.Background()

	// Note: This will try to download the model, which may fail in tests
	// due to network constraints. We just verify it doesn't panic.
	_ = SetupWhisper(ctx, opts)
}

func TestFullSetup(t *testing.T) {
	// Skip in short mode - would require user input
	if testing.Short() {
		t.Skip("Skipping interactive test in short mode")
	}

	tempDir := t.TempDir()

	opts := DefaultSetupOptions()
	opts.ModelsDir = tempDir
	opts.LibDir = filepath.Join(tempDir, "lib")
	opts.Interactive = true

	ctx := context.Background()

	// Would require user input, so we can't fully test it
	_ = SetupWhisper(ctx, opts)
}

func TestSetupForProduction(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()

	ctx := context.Background()
	err := SetupForProduction(ctx, func(opts *SetupOptions) {
		opts.ModelsDir = tempDir
		opts.LibDir = filepath.Join(tempDir, "lib")
		opts.DownloadModels = []string{"tiny"}
	})

	// May fail due to network, but should not panic
	_ = err
}

func TestSetupWithModels(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()

	opts := DefaultSetupOptions()
	opts.ModelsDir = tempDir
	opts.LibDir = filepath.Join(tempDir, "lib")
	opts.DownloadModels = []string{"tiny", "base"}
	opts.Interactive = false

	ctx := context.Background()

	// May fail due to network, but should not panic
	_ = SetupWhisper(ctx, opts)
}

func TestSetupDirectories(t *testing.T) {
	tempDir := t.TempDir()

	opts := &SetupOptions{
		ModelsDir: filepath.Join(tempDir, "models"),
		LibDir:    filepath.Join(tempDir, "lib"),
	}

	err := setupDirectories(opts)
	require.NoError(t, err)

	// Verify directories were created
	_, err = os.Stat(opts.ModelsDir)
	assert.NoError(t, err)

	_, err = os.Stat(opts.LibDir)
	assert.NoError(t, err)
}

func TestVerifySetup(t *testing.T) {
	tempDir := t.TempDir()

	// Create models directory
	modelsDir := filepath.Join(tempDir, "models")
	err := os.MkdirAll(modelsDir, 0755)
	require.NoError(t, err)

	// Create lib directory
	libDir := filepath.Join(tempDir, "lib")
	err = os.MkdirAll(libDir, 0755)
	require.NoError(t, err)

	opts := &SetupOptions{
		ModelsDir: modelsDir,
		LibDir:    libDir,
	}

	err = verifySetup(opts)
	assert.NoError(t, err)
}

func TestVerifySetup_MissingModelsDir(t *testing.T) {
	tempDir := t.TempDir()

	opts := &SetupOptions{
		ModelsDir: filepath.Join(tempDir, "nonexistent", "models"),
		LibDir:    filepath.Join(tempDir, "lib"),
	}

	err := verifySetup(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "models directory does not exist")
}

func TestVerifySetup_WithModels(t *testing.T) {
	tempDir := t.TempDir()

	// Create models directory with a fake model
	modelsDir := filepath.Join(tempDir, "models")
	err := os.MkdirAll(modelsDir, 0755)
	require.NoError(t, err)

	// Create a fake model file
	modelPath := filepath.Join(modelsDir, "tiny.bin")
	f, err := os.Create(modelPath)
	require.NoError(t, err)
	f.Close()

	// Create lib directory
	libDir := filepath.Join(tempDir, "lib")
	err = os.MkdirAll(libDir, 0755)
	require.NoError(t, err)

	opts := &SetupOptions{
		ModelsDir: modelsDir,
		LibDir:    libDir,
	}

	err = verifySetup(opts)
	assert.NoError(t, err)
}

func TestContains(t *testing.T) {
	tests := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{[]string{"a"}, "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPrintSetupInstructions(t *testing.T) {
	// Should not panic
	PrintSetupInstructions()
}

func TestDiagnoseSetup(t *testing.T) {
	tempDir := t.TempDir()

	// Create directories
	modelsDir := filepath.Join(tempDir, "models")
	libDir := filepath.Join(tempDir, "lib")
	err := os.MkdirAll(modelsDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(libDir, 0755)
	require.NoError(t, err)

	opts := &SetupOptions{
		ModelsDir: modelsDir,
		LibDir:    libDir,
	}

	err = DiagnoseSetup(opts)
	// May have warnings but should not error
	assert.NoError(t, err)
}

func TestDiagnoseSetup_NoFFmpeg(t *testing.T) {
	tempDir := t.TempDir()

	// Create directories
	modelsDir := filepath.Join(tempDir, "models")
	libDir := filepath.Join(tempDir, "lib")
	err := os.MkdirAll(modelsDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(libDir, 0755)
	require.NoError(t, err)

	opts := &SetupOptions{
		ModelsDir: modelsDir,
		LibDir:    libDir,
	}

	// This will likely report FFmpeg not found, but should not panic
	err = DiagnoseSetup(opts)
	// Should complete without error even if FFmpeg is missing
	assert.NoError(t, err)
}
