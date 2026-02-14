// Package whisperapp provides utilities for managing whisper models.
// Note: Model conversion utilities are deprecated as whisper is not yet in production.
package whisperapp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Deprecated: ConvertModel is deprecated as whisper is not yet in production.
// Old format: ggml-{name}.bin (whisper.cpp format)
// New format: {name}.bin (gowhisper format)
// Kept for future reference if migration is needed.
// Example: ggml-base.bin -> base.bin
func ConvertModel(modelsDir, modelName string) error {
	oldPath := filepath.Join(modelsDir, fmt.Sprintf("ggml-%s.bin", modelName))
	newPath := filepath.Join(modelsDir, fmt.Sprintf("%s.bin", modelName))

	// Check if old file exists
	oldInfo, err := os.Stat(oldPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("old model file not found: %s", oldPath)
		}
		return fmt.Errorf("failed to stat old model: %w", err)
	}

	// Check if new file already exists
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("new model file already exists: %s (use --force to overwrite)", newPath)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check new model file: %w", err)
	}

	// Copy file from old to new format
	if err := copyFile(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to copy model file: %w", err)
	}

	// Verify new file size matches old file size
	newInfo, err := os.Stat(newPath)
	if err != nil {
		os.Remove(newPath) // Cleanup on error
		return fmt.Errorf("failed to stat new model: %w", err)
	}

	if oldInfo.Size() != newInfo.Size() {
		os.Remove(newPath) // Cleanup on error
		return fmt.Errorf("file size mismatch after conversion: old=%d, new=%d", oldInfo.Size(), newInfo.Size())
	}

	return nil
}

// Deprecated: ConvertModelWithForce is deprecated as whisper is not yet in production.
// Converts a model, overwriting the target if it exists.
func ConvertModelWithForce(modelsDir, modelName string) error {
	newPath := filepath.Join(modelsDir, fmt.Sprintf("%s.bin", modelName))

	// Remove new file if it exists
	if _, err := os.Stat(newPath); err == nil {
		if err := os.Remove(newPath); err != nil {
			return fmt.Errorf("failed to remove existing new model: %w", err)
		}
	}

	return ConvertModel(modelsDir, modelName)
}

// Deprecated: ConvertAllModels is deprecated as whisper is not yet in production.
// Converts all models from old format to new format in a directory.
// Returns list of converted models and any errors encountered.
func ConvertAllModels(modelsDir string, force bool) ([]string, []error) {
	var converted []string
	var errors []error

	// Read directory
	files, err := os.ReadDir(modelsDir)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to read models directory: %w", err))
		return converted, errors
	}

	// Find all old format models
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasPrefix(name, "ggml-") && strings.HasSuffix(name, ".bin") {
			// Extract model name from "ggml-{name}.bin"
			modelName := strings.TrimPrefix(name, "ggml-")
			modelName = strings.TrimSuffix(modelName, ".bin")

			// Convert model
			var err error
			if force {
				err = ConvertModelWithForce(modelsDir, modelName)
			} else {
				err = ConvertModel(modelsDir, modelName)
			}

			if err != nil {
				errors = append(errors, fmt.Errorf("failed to convert %s: %w", modelName, err))
			} else {
				converted = append(converted, modelName)
				fmt.Printf("✓ Converted: %s -> %s\n", name, modelName+".bin")
			}
		}
	}

	return converted, errors
}

// Deprecated: DeleteOldModels is deprecated as whisper is not yet in production.
// Deletes old format model files after successful conversion.
// Use with caution!
func DeleteOldModels(modelsDir string) ([]string, []error) {
	var deleted []string
	var errors []error

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to read models directory: %w", err))
		return deleted, errors
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasPrefix(name, "ggml-") && strings.HasSuffix(name, ".bin") {
			path := filepath.Join(modelsDir, name)

			if err := os.Remove(path); err != nil {
				errors = append(errors, fmt.Errorf("failed to delete %s: %w", name, err))
			} else {
				deleted = append(deleted, name)
				fmt.Printf("✓ Deleted: %s\n", name)
			}
		}
	}

	return deleted, errors
}

// Deprecated: VerifyConversion is deprecated as whisper is not yet in production.
// Verifies that models have been successfully converted.
// Checks that old format models are gone and new format models exist.
func VerifyConversion(modelsDir string, modelNames []string) ([]string, []error) {
	var verified []string
	var errors []error

	for _, modelName := range modelNames {
		oldPath := filepath.Join(modelsDir, fmt.Sprintf("ggml-%s.bin", modelName))
		newPath := filepath.Join(modelsDir, fmt.Sprintf("%s.bin", modelName))

		// Check old file doesn't exist
		if _, err := os.Stat(oldPath); err == nil {
			errors = append(errors, fmt.Errorf("old model still exists: %s", modelName))
			continue
		}

		// Check new file exists
		if _, err := os.Stat(newPath); err != nil {
			if os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("new model not found: %s", modelName))
			} else {
				errors = append(errors, fmt.Errorf("failed to check new model %s: %w", modelName, err))
			}
			continue
		}

		verified = append(verified, modelName)
		fmt.Printf("✓ Verified: %s\n", modelName)
	}

	return verified, errors
}

// Deprecated: GetOldModels is deprecated as whisper is not yet in production.
// Returns a list of old format models in a directory.
func GetOldModels(modelsDir string) ([]string, error) {
	var models []string

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasPrefix(name, "ggml-") && strings.HasSuffix(name, ".bin") {
			// Extract model name from "ggml-{name}.bin"
			modelName := strings.TrimPrefix(name, "ggml-")
			modelName = strings.TrimSuffix(modelName, ".bin")
			models = append(models, modelName)
		}
	}

	return models, nil
}

// Deprecated: GetNewModels is deprecated as whisper is not yet in production.
// Returns a list of new format models in a directory.
func GetNewModels(modelsDir string) ([]string, error) {
	var models []string

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		// Check if it's a valid model (not old format)
		if strings.HasSuffix(name, ".bin") && !strings.HasPrefix(name, "ggml-") {
			// Extract model name from "{name}.bin"
			modelName := strings.TrimSuffix(name, ".bin")

			// Verify it's a valid model by checking specs
			if _, exists := GetModelSpec(modelName); exists {
				models = append(models, modelName)
			}
		}
	}

	return models, nil
}

// Deprecated: ConversionStats is deprecated as whisper is not yet in production.
// Provides statistics about the conversion process.
type ConversionStats struct {
	TotalOldModels   int
	TotalNewModels   int
	ConvertedModels  int
	FailedModels     int
	DeletedOldModels int
	VerifiedModels   int
}

// Deprecated: GetConversionStats is deprecated as whisper is not yet in production.
// Returns statistics about model conversion status.
func GetConversionStats(modelsDir string) (*ConversionStats, error) {
	stats := &ConversionStats{}

	oldModels, err := GetOldModels(modelsDir)
	if err != nil {
		return nil, err
	}
	stats.TotalOldModels = len(oldModels)

	newModels, err := GetNewModels(modelsDir)
	if err != nil {
		return nil, err
	}
	stats.TotalNewModels = len(newModels)

	// Count converted models (models that exist in new format but not old)
	for _, newModel := range newModels {
		oldPath := filepath.Join(modelsDir, fmt.Sprintf("ggml-%s.bin", newModel))
		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			stats.ConvertedModels++
		}
	}

	return stats, nil
}

// Deprecated: PrintConversionSummary is deprecated as whisper is not yet in production.
// Prints a summary of the conversion process.
func PrintConversionSummary(stats *ConversionStats, errors []error) {
	fmt.Println("\n=== Conversion Summary ===")
	fmt.Printf("Old format models found: %d\n", stats.TotalOldModels)
	fmt.Printf("New format models found: %d\n", stats.TotalNewModels)
	fmt.Printf("Successfully converted: %d\n", stats.ConvertedModels)
	fmt.Printf("Failed conversions: %d\n", stats.FailedModels)
	if stats.DeletedOldModels > 0 {
		fmt.Printf("Deleted old models: %d\n", stats.DeletedOldModels)
	}
	if stats.VerifiedModels > 0 {
		fmt.Printf("Verified conversions: %d\n", stats.VerifiedModels)
	}

	if len(errors) > 0 {
		fmt.Printf("\n=== Errors (%d) ===\n", len(errors))
		for i, err := range errors {
			fmt.Printf("%d. %v\n", i+1, err)
		}
	}

	fmt.Println("========================")
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
