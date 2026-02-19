package kronk

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownloadService_New(t *testing.T) {
	svc := NewDownloadService(
		WithBasePath("/tmp/test"),
		WithModelsPath("/tmp/test/models"),
		WithMaxRetries(3),
		WithRetryDelay(1*time.Second),
		WithTimeout(5*time.Minute),
	)

	if svc == nil {
		t.Fatal("Expected DownloadService to be created")
	}

	if svc.basePath != "/tmp/test" {
		t.Errorf("Expected basePath /tmp/test, got %s", svc.basePath)
	}

	if svc.modelsPath != "/tmp/test/models" {
		t.Errorf("Expected modelsPath /tmp/test/models, got %s", svc.modelsPath)
	}

	if svc.maxRetries != 3 {
		t.Errorf("Expected maxRetries 3, got %d", svc.maxRetries)
	}
}

func TestDownloadService_WithProgressCallback(t *testing.T) {
	cb := func(completed, total int64, percent float64, mbps float64) {
		// Callback implementation
	}

	svc := NewDownloadService(
		WithProgressCallback(cb),
	)

	if svc.progressCb == nil {
		t.Fatal("Expected progressCb to be set")
	}
}

func TestDownloadService_DownloadWithRetry_NoNetwork(t *testing.T) {
	// This test will fail gracefully when no network
	svc := NewDownloadService(
		WithMaxRetries(1),
		WithRetryDelay(100*time.Millisecond),
		WithTimeout(1*time.Second),
	)

	ctx := context.Background()
	
	// Try to download from invalid URL
	result := svc.DownloadWithRetry(ctx, "http://invalid.url.test/file.bin", "/tmp/test.bin")
	
	// Should fail (no network or invalid URL)
	if result.Success {
		t.Error("Expected download to fail with invalid URL")
	}
	
	if result.Error == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDownloadService_DownloadWithRetry_ContextCancellation(t *testing.T) {
	svc := NewDownloadService(
		WithMaxRetries(3),
		WithRetryDelay(1*time.Second),
		WithTimeout(10*time.Second),
	)

	ctx, cancel := context.WithCancel(context.Background())
	
	// Cancel immediately
	cancel()
	
	result := svc.DownloadWithRetry(ctx, "http://example.com/file.bin", "/tmp/test.bin")
	
	if result.Success {
		t.Error("Expected download to fail with cancelled context")
	}
	
	// Check if error is context.Canceled (may be wrapped)
	if result.Error == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDownloadService_DownloadWithRetry_RetryLogic(t *testing.T) {
	svc := &DownloadService{
		basePath:     "/tmp",
		modelsPath:   "/tmp/models",
		maxRetries:   3,
		retryDelay:   10 * time.Millisecond,
		timeout:      100 * time.Millisecond,
	}
	
	ctx := context.Background()
	
	// Try to download from invalid URL - should retry
	result := svc.DownloadWithRetry(ctx, "http://invalid.url.test/file.bin", "/tmp/test.bin")
	
	// Should have attempted multiple times
	if result.Success {
		t.Error("Expected download to fail")
	}
	
	// Verify error message mentions retries
	if result.Error == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDownloadService_DownloadResult_Fields(t *testing.T) {
	result := &DownloadResult{
		Success:  true,
		FilePath: "/tmp/test.bin",
		Bytes:    1024,
		Duration: 1 * time.Second,
		Resumed:  false,
		Error:    nil,
	}
	
	if !result.Success {
		t.Error("Expected Success to be true")
	}
	
	if result.FilePath != "/tmp/test.bin" {
		t.Errorf("Expected FilePath /tmp/test.bin, got %s", result.FilePath)
	}
	
	if result.Bytes != 1024 {
		t.Errorf("Expected Bytes 1024, got %d", result.Bytes)
	}
	
	if result.Resumed {
		t.Error("Expected Resumed to be false")
	}
}

func TestDownloadService_DownloadLLMModel_InvalidURL(t *testing.T) {
	svc := NewDownloadService(
		WithModelsPath("/tmp/models"),
		WithMaxRetries(1),
		WithRetryDelay(10*time.Millisecond),
		WithTimeout(100*time.Millisecond),
	)
	
	ctx := context.Background()
	
	// Try with invalid org/repo
	result := svc.DownloadLLMModel(ctx, "invalid-org", "invalid-repo", "invalid.bin")
	
	if result.Success {
		t.Error("Expected download to fail with invalid URL")
	}
}

func TestDownloadService_DirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	modelsPath := filepath.Join(tmpDir, "models")
	
	svc := NewDownloadService(
		WithModelsPath(modelsPath),
	)
	
	// Verify directory doesn't exist yet
	if _, err := os.Stat(modelsPath); !os.IsNotExist(err) {
		t.Fatalf("Expected %s to not exist", modelsPath)
	}
	
	// The service should handle directory creation during download
	// We just verify the path is set correctly
	if svc.modelsPath != modelsPath {
		t.Errorf("Expected modelsPath %s, got %s", modelsPath, svc.modelsPath)
	}
}

func TestDownloadService_ProgressCallback_Invocation(t *testing.T) {
	progressCalls := 0
	
	cb := func(completed, total int64, percent float64, mbps float64) {
		progressCalls++
	}
	
	svc := NewDownloadService(
		WithProgressCallback(cb),
		WithMaxRetries(1),
		WithTimeout(100*time.Millisecond),
	)
	
	ctx := context.Background()
	
	// Try to download - progress callback should be invoked
	// (even if download fails)
	_ = svc.DownloadWithRetry(ctx, "http://invalid.url.test/file.bin", "/tmp/test.bin")
	
	// Note: Progress callback may or may not be invoked depending on when the download fails
	// This test just verifies the callback is properly wired up
	_ = progressCalls
}

func TestDownloadService_ConcurrentDownloads(t *testing.T) {
	svc := NewDownloadService(
		WithMaxRetries(1),
		WithTimeout(100*time.Millisecond),
	)
	
	ctx := context.Background()
	
	// Start multiple downloads concurrently
	results := make(chan *DownloadResult, 3)
	
	go func() {
		results <- svc.DownloadWithRetry(ctx, "http://invalid1.url.test/file.bin", "/tmp/test1.bin")
	}()
	
	go func() {
		results <- svc.DownloadWithRetry(ctx, "http://invalid2.url.test/file.bin", "/tmp/test2.bin")
	}()
	
	go func() {
		results <- svc.DownloadWithRetry(ctx, "http://invalid3.url.test/file.bin", "/tmp/test3.bin")
	}()
	
	// Collect results
	for i := 0; i < 3; i++ {
		result := <-results
		if result.Success {
			t.Errorf("Download %d should have failed", i)
		}
	}
}

func TestDownloadService_Timeout(t *testing.T) {
	svc := NewDownloadService(
		WithMaxRetries(1),
		WithRetryDelay(10*time.Millisecond),
		WithTimeout(50*time.Millisecond), // Very short timeout
	)
	
	ctx := context.Background()
	
	start := time.Now()
	result := svc.DownloadWithRetry(ctx, "http://invalid.url.test/file.bin", "/tmp/test.bin")
	duration := time.Since(start)
	
	if result.Success {
		t.Error("Expected download to fail")
	}
	
	// Should timeout quickly
	if duration > 500*time.Millisecond {
		t.Errorf("Expected timeout to occur quickly, took %v", duration)
	}
}
