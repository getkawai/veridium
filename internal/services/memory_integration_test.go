package services

import "testing"

func TestMemoryIntegrationUsesMuninnBackendFlag(t *testing.T) {
	integration, err := NewMemoryIntegration(&MemoryIntegrationConfig{
		MuninnBackend: &MuninnMemoryBackend{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !integration.UsesMuninnBackend() {
		t.Fatalf("expected muninn backend flag to be true")
	}
}

func TestMemoryIntegrationRegisterMemoryToolDisabled(t *testing.T) {
	integration, err := NewMemoryIntegration(&MemoryIntegrationConfig{
		MuninnBackend: &MuninnMemoryBackend{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := integration.RegisterMemoryTool(nil); err != nil {
		t.Fatalf("expected disabled registration to return nil, got %v", err)
	}
}

func TestMemoryIntegrationRequiresMuninnBackend(t *testing.T) {
	if _, err := NewMemoryIntegration(&MemoryIntegrationConfig{}); err == nil {
		t.Fatalf("expected error when muninn backend is missing")
	}
}
