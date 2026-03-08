package app

import (
	"testing"

	"github.com/kawai-network/veridium/internal/services"
)

func TestBuildMuninnMemoryIntegrationRequiresBackend(t *testing.T) {
	if _, err := buildMuninnMemoryIntegration(nil); err == nil {
		t.Fatalf("expected error when backend is nil")
	}
}

func TestBuildMuninnMemoryIntegrationUsesMuninnMode(t *testing.T) {
	integration, err := buildMuninnMemoryIntegration(&services.MuninnMemoryBackend{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !integration.UsesMuninnBackend() {
		t.Fatalf("expected integration to use muninn backend")
	}
}
