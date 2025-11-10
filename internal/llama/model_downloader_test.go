package llama

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// setupTestService creates a test service with a temporary models directory
func setupTestService(t *testing.T) (*Service, string, func()) {
	tempDir := t.TempDir()

	// Override HOME to use temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)

	// Create the expected models directory
	modelsDir := filepath.Join(tempDir, ".llama-cpp", "models")
	os.MkdirAll(modelsDir, 0755)

	manager := NewLlamaCppReleaseManager()
	service := &Service{manager: manager}

	cleanup := func() {
		os.Setenv("HOME", oldHome)
	}

	return service, modelsDir, cleanup
}

// Test model selection based on RAM
func TestSelectOptimalQwenModel(t *testing.T) {
	tests := []struct {
		name          string
		availableRAM  int64
		expectedModel string
	}{
		{
			name:          "Low RAM (1GB) - selects smallest model",
			availableRAM:  1,
			expectedModel: "qwen2.5-0.5b-instruct-q4_k_m",
		},
		{
			name:          "Medium RAM (4GB) - selects 1.5B model",
			availableRAM:  4,
			expectedModel: "qwen2.5-1.5b-instruct-q4_k_m",
		},
		{
			name:          "High RAM (10GB) - selects 7B model",
			availableRAM:  10,
			expectedModel: "qwen2.5-7b-instruct-q4_k_m",
		},
		{
			name:          "Very High RAM (32GB) - selects largest available",
			availableRAM:  32,
			expectedModel: "qwen2.5-7b-instruct-q4_k_m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := SelectOptimalQwenModel(tt.availableRAM)
			if model.Name != tt.expectedModel {
				t.Errorf("SelectOptimalQwenModel() = %v, want %v", model.Name, tt.expectedModel)
			}
		})
	}
}

// Test successful download
func TestDownloadModel_Success(t *testing.T) {
	// Create test server that serves a valid GGUF file
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write GGUF magic number
		w.Write([]byte("GGUF"))
		// Write dummy data to reach expected size
		dummyData := make([]byte, 1024*1024) // 1MB
		w.Write(dummyData)
	}))
	defer server.Close()

	service, modelsDir, cleanup := setupTestService(t)
	defer cleanup()

	// Create test model spec
	modelSpec := QwenModelSpec{
		Name:   "test-model",
		URL:    server.URL,
		Size:   1024 * 1024, // 1MB
		SHA256: "",          // Skip checksum for this test
	}

	// Test download
	err := service.DownloadModel(modelSpec)
	if err != nil {
		t.Fatalf("DownloadModel() error = %v", err)
	}

	// Verify file exists
	modelPath := filepath.Join(modelsDir, "test-model.gguf")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Errorf("Downloaded model file not found at %s", modelPath)
	}

	// Verify file content starts with GGUF
	content, err := os.ReadFile(modelPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if !strings.HasPrefix(string(content), "GGUF") {
		t.Errorf("Downloaded file doesn't start with GGUF magic number")
	}
}

// Test download with network failure
func TestDownloadModel_NetworkFailure(t *testing.T) {
	// Create test server that fails
	failCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failCount++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")
	service := &Service{manager: manager}

	modelSpec := QwenModelSpec{
		Name: "test-model-fail",
		URL:  server.URL,
		Size: 1024,
	}

	err := service.DownloadModel(modelSpec)
	if err == nil {
		t.Error("DownloadModel() expected error for network failure, got nil")
	}

	// Verify no partial file left behind
	modelPath := filepath.Join(tempDir, "test-model-fail.gguf")
	if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
		t.Error("Partial file should be cleaned up after failure")
	}

	// Verify temp file is also cleaned up
	tempPath := modelPath + ".tmp"
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		t.Error("Temporary file should be cleaned up after failure")
	}
}

// Test download with size mismatch
func TestDownloadModel_SizeMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write GGUF header but wrong size
		w.Write([]byte("GGUF"))
		w.Write(make([]byte, 100)) // Only 104 bytes instead of expected 1MB
	}))
	defer server.Close()

	tempDir := t.TempDir()
	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")
	service := &Service{manager: manager}

	modelSpec := QwenModelSpec{
		Name: "test-model-size",
		URL:  server.URL,
		Size: 1024 * 1024, // Expect 1MB
	}

	err := service.DownloadModel(modelSpec)
	if err == nil {
		t.Error("DownloadModel() expected error for size mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "size mismatch") {
		t.Errorf("DownloadModel() error should mention size mismatch, got: %v", err)
	}

	// Verify cleanup
	modelPath := filepath.Join(tempDir, "test-model-size.gguf")
	if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
		t.Error("File with wrong size should be cleaned up")
	}
}

// Test download with invalid GGUF file (corrupt file)
func TestDownloadModel_InvalidGGUF(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid magic number
		w.Write([]byte("FAKE"))
		w.Write(make([]byte, 1024*1024))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")
	service := &Service{manager: manager}

	modelSpec := QwenModelSpec{
		Name: "test-model-corrupt",
		URL:  server.URL,
		Size: 1024 * 1024,
	}

	err := service.DownloadModel(modelSpec)
	if err == nil {
		t.Error("DownloadModel() expected error for invalid GGUF, got nil")
	}
	if !strings.Contains(err.Error(), "invalid GGUF") {
		t.Errorf("DownloadModel() error should mention invalid GGUF, got: %v", err)
	}

	// Verify cleanup
	modelPath := filepath.Join(tempDir, "test-model-corrupt.gguf")
	if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
		t.Error("Corrupt file should be cleaned up")
	}
}

// Test download with checksum verification
func TestDownloadModel_ChecksumMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GGUF"))
		w.Write([]byte("test data"))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")
	service := &Service{manager: manager}

	modelSpec := QwenModelSpec{
		Name:   "test-model-checksum",
		URL:    server.URL,
		Size:   13,                                                                 // 4 (GGUF) + 9 (test data)
		SHA256: "0000000000000000000000000000000000000000000000000000000000000000", // Wrong checksum
	}

	err := service.DownloadModel(modelSpec)
	if err == nil {
		t.Error("DownloadModel() expected error for checksum mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "integrity check failed") {
		t.Errorf("DownloadModel() error should mention integrity check, got: %v", err)
	}
}

// Test download with slow connection (timeout simulation)
func TestDownloadModel_SlowConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow connection test in short mode")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow connection
		w.Write([]byte("GGUF"))
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			w.Write([]byte{byte(i)})
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")
	service := &Service{manager: manager}

	modelSpec := QwenModelSpec{
		Name: "test-model-slow",
		URL:  server.URL,
		Size: 14, // 4 (GGUF) + 10 bytes
	}

	// Should succeed even with slow connection
	err := service.DownloadModel(modelSpec)
	if err != nil {
		t.Errorf("DownloadModel() should handle slow connections, got error: %v", err)
	}
}

// Test existing model skip
func TestDownloadModel_ExistingModel(t *testing.T) {
	tempDir := t.TempDir()

	// Create existing valid model file
	existingPath := filepath.Join(tempDir, "existing-model.gguf")
	err := os.WriteFile(existingPath, []byte("GGUFexisting data"), 0644)
	if err != nil {
		t.Fatalf("Failed to create existing model file: %v", err)
	}

	// Track if server was called
	serverCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		w.Write([]byte("GGUFnew data"))
	}))
	defer server.Close()

	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")
	service := &Service{manager: manager}

	modelSpec := QwenModelSpec{
		Name: "existing-model",
		URL:  server.URL,
		Size: 12,
	}

	err = service.DownloadModel(modelSpec)
	if err != nil {
		t.Errorf("DownloadModel() should skip existing model, got error: %v", err)
	}

	if serverCalled {
		t.Error("Server should not be called when model already exists")
	}

	// Verify original file unchanged
	content, _ := os.ReadFile(existingPath)
	if string(content) != "GGUFexisting data" {
		t.Error("Existing model file should not be modified")
	}
}

// Test validateGGUFFile
func TestValidateGGUFFile(t *testing.T) {
	tests := []struct {
		name      string
		content   []byte
		wantError bool
	}{
		{
			name:      "Valid GGUF file",
			content:   []byte("GGUFsome model data here"),
			wantError: false,
		},
		{
			name:      "Invalid magic number",
			content:   []byte("FAKEsome model data here"),
			wantError: true,
		},
		{
			name:      "Too short file",
			content:   []byte("GG"),
			wantError: true,
		},
		{
			name:      "Empty file",
			content:   []byte(""),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := filepath.Join(t.TempDir(), "test.gguf")
			err := os.WriteFile(tempFile, tt.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			manager := NewLlamaCppReleaseManager()
			service := &Service{manager: manager}

			err = service.validateGGUFFile(tempFile)
			if (err != nil) != tt.wantError {
				t.Errorf("validateGGUFFile() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// Test verifyModelChecksum
func TestVerifyModelChecksum(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.gguf")
	testData := []byte("test content for checksum")

	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate actual checksum
	actualChecksum := "87de2c7cc9d3a7d03005c32f08a3d433c5c5c4b91b00baaee8b27b1f1f0e2e3d"

	manager := NewLlamaCppReleaseManager()
	service := &Service{manager: manager}

	tests := []struct {
		name      string
		checksum  string
		wantError bool
	}{
		{
			name:      "Valid checksum",
			checksum:  actualChecksum,
			wantError: false,
		},
		{
			name:      "Invalid checksum",
			checksum:  "0000000000000000000000000000000000000000000000000000000000000000",
			wantError: true,
		},
		{
			name:      "Uppercase checksum (should work)",
			checksum:  strings.ToUpper(actualChecksum),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.verifyModelChecksum(testFile, tt.checksum)
			if (err != nil) != tt.wantError {
				t.Errorf("verifyModelChecksum() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// Test partial download cleanup
func TestDownloadModel_PartialDownloadCleanup(t *testing.T) {
	// Server that closes connection mid-transfer
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GGUF"))
		// Simulate partial write
		w.Write(make([]byte, 100))
		// Force close connection
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")
	service := &Service{manager: manager}

	modelSpec := QwenModelSpec{
		Name: "test-model-partial",
		URL:  server.URL,
		Size: 1024 * 1024, // Expect much larger
	}

	err := service.DownloadModel(modelSpec)
	if err == nil {
		t.Error("DownloadModel() expected error for partial download")
	}

	// Verify all files are cleaned up
	modelPath := filepath.Join(tempDir, "test-model-partial.gguf")
	tempPath := modelPath + ".tmp"

	if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
		t.Error("Partial model file should be cleaned up")
	}
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		t.Error("Partial temp file should be cleaned up")
	}
}

// Test interrupted download with stale temp file
func TestDownloadModel_InterruptedDownloadResume(t *testing.T) {
	service, modelsDir, cleanup := setupTestService(t)
	defer cleanup()

	// Simulate a stale temp file from previous interrupted download
	staleTempFile := filepath.Join(modelsDir, "test-model-resume.gguf.tmp")
	staleContent := []byte("INCOMPLETE DOWNLOAD DATA FROM PREVIOUS RUN")
	err := os.WriteFile(staleTempFile, staleContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create stale temp file: %v", err)
	}

	// Verify stale file exists
	if _, err := os.Stat(staleTempFile); os.IsNotExist(err) {
		t.Fatal("Stale temp file should exist before test")
	}

	// Create server for successful download
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GGUF"))
		w.Write(make([]byte, 1024*100)) // 100KB
	}))
	defer server.Close()

	modelSpec := QwenModelSpec{
		Name: "test-model-resume",
		URL:  server.URL,
		Size: 1024 * 100,
	}

	// Download should clean up stale file and succeed
	err = service.DownloadModel(modelSpec)
	if err != nil {
		t.Fatalf("DownloadModel() should clean up stale file and succeed: %v", err)
	}

	// Verify final model exists
	finalPath := filepath.Join(modelsDir, "test-model-resume.gguf")
	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		t.Error("Final model file should exist after download")
	}

	// Verify stale temp file was cleaned up
	if _, err := os.Stat(staleTempFile); !os.IsNotExist(err) {
		t.Error("Stale temp file should be cleaned up")
	}

	// Verify final file has correct content (not stale content)
	content, _ := os.ReadFile(finalPath)
	if strings.Contains(string(content), "INCOMPLETE DOWNLOAD") {
		t.Error("Final file should not contain stale download data")
	}
}

// Test CleanupStaleTempFiles function
func TestCleanupStaleTempFiles(t *testing.T) {
	service, modelsDir, cleanup := setupTestService(t)
	defer cleanup()

	// Create multiple stale temp files
	staleFiles := []string{
		"model1.gguf.tmp",
		"model2.gguf.tmp",
		"model3.gguf.tmp",
	}

	for _, filename := range staleFiles {
		path := filepath.Join(modelsDir, filename)
		content := []byte("STALE DATA " + filename)
		err := os.WriteFile(path, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create stale file %s: %v", filename, err)
		}
	}

	// Also create a valid model file (should NOT be deleted)
	validModelPath := filepath.Join(modelsDir, "valid-model.gguf")
	err := os.WriteFile(validModelPath, []byte("GGUF valid model data"), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid model: %v", err)
	}

	// Run cleanup
	err = service.CleanupStaleTempFiles()
	if err != nil {
		t.Fatalf("CleanupStaleTempFiles() error = %v", err)
	}

	// Verify all temp files were deleted
	for _, filename := range staleFiles {
		path := filepath.Join(modelsDir, filename)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("Stale temp file %s should be deleted", filename)
		}
	}

	// Verify valid model file still exists
	if _, err := os.Stat(validModelPath); os.IsNotExist(err) {
		t.Error("Valid model file should not be deleted")
	}
}

// Test application crash during download simulation
func TestDownloadModel_ApplicationCrashDuringDownload(t *testing.T) {
	service, modelsDir, cleanup := setupTestService(t)
	defer cleanup()

	// Simulate partial download by creating an incomplete temp file
	partialTempFile := filepath.Join(modelsDir, "crashed-download.gguf.tmp")
	partialContent := []byte("GGUF") // Valid header but incomplete
	for i := 0; i < 1000; i++ {
		partialContent = append(partialContent, byte(i%256))
	}
	err := os.WriteFile(partialTempFile, partialContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create partial temp file: %v", err)
	}

	log.Printf("📝 Test: Simulated crash with partial file: %d bytes", len(partialContent))

	// Now try to download the model again (simulating app restart after crash)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GGUF"))
		w.Write(make([]byte, 1024*200)) // 200KB complete file
	}))
	defer server.Close()

	modelSpec := QwenModelSpec{
		Name: "crashed-download",
		URL:  server.URL,
		Size: 1024 * 200,
	}

	// Download should succeed and replace partial file
	err = service.DownloadModel(modelSpec)
	if err != nil {
		t.Fatalf("DownloadModel() should succeed after crash cleanup: %v", err)
	}

	// Verify final file has correct size (not partial)
	finalPath := filepath.Join(modelsDir, "crashed-download.gguf")
	info, err := os.Stat(finalPath)
	if err != nil {
		t.Fatalf("Final model file should exist: %v", err)
	}

	// Size should be close to expected (200KB + header)
	if info.Size() < 1024*200 {
		t.Errorf("Final file seems incomplete: got %d bytes, expected ~%d bytes", info.Size(), 1024*200)
	}

	// Verify temp file was cleaned up
	if _, err := os.Stat(partialTempFile); !os.IsNotExist(err) {
		t.Error("Partial temp file should be cleaned up")
	}
}

// Test multiple interrupted downloads
func TestCleanupStaleTempFiles_MultipleInterrupted(t *testing.T) {
	service, modelsDir, cleanup := setupTestService(t)
	defer cleanup()

	// Simulate multiple interrupted downloads of different sizes
	interruptedDownloads := []struct {
		name string
		size int
	}{
		{"tiny-model.gguf.tmp", 1024},              // 1KB
		{"small-model.gguf.tmp", 1024 * 100},       // 100KB
		{"medium-model.gguf.tmp", 1024 * 1024},     // 1MB
		{"large-model.gguf.tmp", 1024 * 1024 * 10}, // 10MB
	}

	totalSize := int64(0)
	for _, dl := range interruptedDownloads {
		path := filepath.Join(modelsDir, dl.name)
		content := make([]byte, dl.size)
		// Fill with some pattern
		for i := range content {
			content[i] = byte(i % 256)
		}
		err := os.WriteFile(path, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create interrupted download %s: %v", dl.name, err)
		}
		totalSize += int64(dl.size)
	}

	log.Printf("📊 Created %d interrupted downloads, total size: %.2f MB",
		len(interruptedDownloads), float64(totalSize)/(1024*1024))

	// Run cleanup
	err := service.CleanupStaleTempFiles()
	if err != nil {
		t.Fatalf("CleanupStaleTempFiles() error = %v", err)
	}

	// Verify all were cleaned
	for _, dl := range interruptedDownloads {
		path := filepath.Join(modelsDir, dl.name)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("Interrupted download %s should be cleaned up", dl.name)
		}
	}

	// Verify directory still exists and is empty of .tmp files
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		t.Fatalf("Failed to read models directory: %v", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			t.Errorf("Found remaining .tmp file: %s", entry.Name())
		}
	}
}

// Benchmark model download (mocked)
func BenchmarkDownloadModel(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GGUF"))
		w.Write(make([]byte, 1024*100)) // 100KB
	}))
	defer server.Close()

	tempDir := b.TempDir()
	manager := NewLlamaCppReleaseManager()
	manager.BinaryPath = filepath.Join(tempDir, "bin")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service := &Service{manager: manager}
		modelSpec := QwenModelSpec{
			Name: fmt.Sprintf("bench-model-%d", i),
			URL:  server.URL,
			Size: 1024 * 100,
		}
		service.DownloadModel(modelSpec)
	}
}
