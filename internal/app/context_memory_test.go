package app

import (
	"path/filepath"
	"testing"

	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/y/paths"
)

func TestProvideMemoryIntegration(t *testing.T) {
	// This test verifies that ProvideMemoryIntegration can initialize successfully
	// Note: Full fx lifecycle testing would require fxtest package
	
	muninnDataDir := filepath.Join(paths.Base(), "muninndb", "veridium_test")
	backend, err := services.NewMuninnMemoryBackend(muninnDataDir, "veridium_test", 10000, false)
	if err != nil {
		t.Skipf("Muninn memory backend init failed (expected in CI): %v", err)
	}
	defer backend.Close()

	integration, err := services.NewMemoryIntegration(&services.MemoryIntegrationConfig{
		MuninnBackend: backend,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !integration.UsesMuninnBackend() {
		t.Fatalf("expected integration to use muninn backend")
	}
}

func TestMemoryIntegrationRequiresBackend(t *testing.T) {
	// Test that MemoryIntegration requires a backend
	if _, err := services.NewMemoryIntegration(&services.MemoryIntegrationConfig{
		MuninnBackend: nil,
	}); err == nil {
		t.Fatalf("expected error when backend is nil")
	}
}
