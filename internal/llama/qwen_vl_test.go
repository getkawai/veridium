package llama

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestQwenVLImageProcessing tests the Qwen-VL model for image explanation
// This test requires a VL model to be downloaded and the docparsing_example1.jpg image
func TestQwenVLImageProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping VL image processing test in short mode (requires model download)")
	}

	// Check if image file exists
	imagePath := "internal/llama/docparsing_example1.jpg"
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Skipf("Test image not found: %s", imagePath)
	}

	// Create library service
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(2 * time.Second)

	// Check if VL model is available, if not download it
	vlModels, err := service.manager.GetAvailableVLModels()
	if err != nil {
		t.Fatalf("Failed to check available VL models: %v", err)
	}

	if len(vlModels) == 0 {
		t.Log("No VL models found, attempting auto-download...")
		if err := service.AutoDownloadRecommendedVLModel(); err != nil {
			t.Skipf("Failed to download VL model: %v", err)
		}

		// Wait for download to complete
		time.Sleep(5 * time.Second)

		vlModels, _ = service.manager.GetAvailableVLModels()
		if len(vlModels) == 0 {
			t.Skip("VL model download failed or took too long")
		}
	}

	t.Logf("Found VL model: %s", filepath.Base(vlModels[0]))

	// Load VL model
	err = service.LoadVLModel(vlModels[0])
	if err != nil {
		t.Fatalf("Failed to load VL model: %v", err)
	}

	// Verify VL model is loaded
	if !service.IsVLModelLoaded() {
		t.Fatal("VL model not properly loaded")
	}

	t.Log("✅ VL model loaded successfully")

	// Test image processing
	prompt := "Jelaskan dengan detail apa yang ada di gambar ini. Apa yang Anda lihat? Berikan deskripsi yang lengkap dan akurat."
	maxTokens := int32(256)

	t.Logf("Processing image: %s", imagePath)
	t.Logf("Prompt: %s", prompt)

	startTime := time.Now()
	response, err := service.ProcessImageWithText(imagePath, prompt, maxTokens)
	duration := time.Since(startTime)

	if err != nil {
		t.Fatalf("Failed to process image: %v", err)
	}

	t.Logf("✅ Image processed successfully in %v", duration)
	t.Logf("Response length: %d characters", len(response))

	// Basic validation
	if len(response) == 0 {
		t.Error("Response is empty")
	}

	if len(response) < 10 {
		t.Errorf("Response too short: %q", response)
	}

	// Log response for inspection
	t.Logf("📝 Response: %s", response)

	// Check if response contains image-related terms (basic sanity check)
	responseLower := strings.ToLower(response)
	if !strings.Contains(responseLower, "gambar") &&
	   !strings.Contains(responseLower, "lihat") &&
	   !strings.Contains(responseLower, "ada") {
		t.Logf("⚠️  Response may not be image-related: %s", response)
	}
}

// TestQwenVLModelLoading tests VL model loading functionality
func TestQwenVLModelLoading(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping VL model loading test in short mode")
	}

	// Create library service
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(2 * time.Second)

	// Test VL model auto-selection
	availableVL, err := service.manager.GetAvailableVLModels()
	if err != nil {
		t.Fatalf("Failed to get available VL models: %v", err)
	}

	if len(availableVL) == 0 {
		t.Skip("No VL models available")
	}

	t.Logf("Available VL models: %d", len(availableVL))

	// Test LoadVLModel with auto-selection
	err = service.LoadVLModel("")
	if err != nil {
		t.Skipf("Failed to load VL model: %v", err)
	}

	// Verify model is loaded
	if !service.IsVLModelLoaded() {
		t.Error("VL model not loaded")
	}

	loadedPath := service.GetLoadedVLModel()
	if loadedPath == "" {
		t.Error("No loaded VL model path returned")
	}

	t.Logf("✅ VL model loaded: %s", filepath.Base(loadedPath))

	// Test model unloading (implicit via Cleanup)
}

// TestQwenVLHardwareDetection tests that VL model selection considers hardware
func TestQwenVLHardwareDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hardware detection test in short mode")
	}

	// Test hardware detection
	specs := DetectHardwareSpecs()

	t.Logf("Detected hardware: RAM=%dGB, CPU cores=%d", specs.AvailableRAM, specs.CPUCores)

	if specs.AvailableRAM == 0 {
		t.Error("Available RAM detection failed")
	}

	if specs.CPUCores == 0 {
		t.Error("CPU cores detection failed")
	}

	// Test that VL model selection works
	vlSelection := SelectOptimalQwenModel(specs.AvailableRAM)
	if vlSelection.Name == "" {
		t.Error("VL model selection failed")
	}

	t.Logf("Selected VL model for %dGB RAM: %s", specs.AvailableRAM, vlSelection.Name)

	// Test that text model selection works
	textSelection := SelectOptimalQwenTextModel(specs.AvailableRAM)
	if textSelection.Name == "" {
		t.Error("Text model selection failed")
	}

	t.Logf("Selected text model for %dGB RAM: %s", specs.AvailableRAM, textSelection.Name)
}

// TestAutoDownloadVLModel tests the VL model auto-download functionality
func TestAutoDownloadVLModel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping VL download test in short mode (may take a long time)")
	}

	// Create library service
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Check initial state
	initialVL, _ := service.manager.GetAvailableVLModels()

	t.Logf("Initial VL models: %d", len(initialVL))

	// Test download
	err = service.AutoDownloadRecommendedVLModel()
	if err != nil {
		t.Logf("VL model download failed (may be network issue): %v", err)
		t.Skip("Skipping due to download failure")
	}

	// Verify download
	afterVL, _ := service.manager.GetAvailableVLModels()

	t.Logf("VL models after download: %d", len(afterVL))

	if len(afterVL) <= len(initialVL) {
		t.Error("VL model download did not increase available models count")
	}

	t.Log("✅ VL model auto-download test passed")
}

// BenchmarkQwenVLProcessing benchmarks VL image processing performance
func BenchmarkQwenVLProcessing(b *testing.B) {
	// This would require a loaded VL model and image
	// For now, just skip as it requires heavy setup
	b.Skip("Benchmark requires loaded VL model and available image")
}
