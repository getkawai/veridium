package kronk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateStableDiffusionModelBundle_RequiredPathEmpty(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	diffusion := filepath.Join(tmpDir, "diffusion.gguf")
	vae := filepath.Join(tmpDir, "vae.safetensors")
	for _, file := range []string{diffusion, vae} {
		if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
			t.Fatalf("write model file %s: %v", file, err)
		}
	}

	err := validateStableDiffusionModelBundle(&stableDiffusionModelBundle{
		diffusionPath: diffusion,
		llmPath:       "",
		vaePath:       vae,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "path is empty (llm_model)") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateStableDiffusionModelBundle_FileNotFound(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	diffusion := filepath.Join(tmpDir, "diffusion.gguf")
	llm := filepath.Join(tmpDir, "llm.gguf")

	if err := os.WriteFile(diffusion, []byte("x"), 0644); err != nil {
		t.Fatalf("write diffusion file: %v", err)
	}
	if err := os.WriteFile(llm, []byte("x"), 0644); err != nil {
		t.Fatalf("write llm file: %v", err)
	}

	err := validateStableDiffusionModelBundle(&stableDiffusionModelBundle{
		diffusionPath: diffusion,
		llmPath:       llm,
		vaePath:       filepath.Join(tmpDir, "missing-vae.safetensors"),
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "file not found (vae_model)") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateStableDiffusionModelBundle_Valid(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	diffusion := filepath.Join(tmpDir, "diffusion.gguf")
	llm := filepath.Join(tmpDir, "llm.gguf")
	vae := filepath.Join(tmpDir, "vae.safetensors")
	edit := filepath.Join(tmpDir, "edit.gguf")

	for _, file := range []string{diffusion, llm, vae, edit} {
		if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
			t.Fatalf("write model file %s: %v", file, err)
		}
	}

	err := validateStableDiffusionModelBundle(&stableDiffusionModelBundle{
		diffusionPath: diffusion,
		llmPath:       llm,
		vaePath:       vae,
		editModelPath: edit,
	})
	if err != nil {
		t.Fatalf("expected validation success, got: %v", err)
	}
}
