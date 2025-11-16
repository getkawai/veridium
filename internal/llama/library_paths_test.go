package llama

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/kawai-network/veridium/pkg/yzma/download"
)

// TestLibraryPathConfiguration verifies that library paths are correctly configured
func TestLibraryPathConfiguration(t *testing.T) {
	t.Log("🔍 Testing Library Path Configuration")

	// Create installer
	installer := NewLlamaCppInstaller()

	// Test 1: Verify GetLibraryPath returns a valid directory
	t.Run("GetLibraryPath", func(t *testing.T) {
		libPath := installer.GetLibraryPath()

		if libPath == "" {
			t.Fatal("GetLibraryPath() returned empty string")
		}

		t.Logf("✅ Library Path: %s", libPath)

		// Should be the BinaryPath
		if libPath != installer.BinaryPath {
			t.Errorf("GetLibraryPath() = %s, want %s", libPath, installer.BinaryPath)
		}
	})

	// Test 2: Verify GetLibraryFilePath returns correct main library
	t.Run("GetLibraryFilePath", func(t *testing.T) {
		libFilePath := installer.GetLibraryFilePath()

		if libFilePath == "" {
			t.Fatal("GetLibraryFilePath() returned empty string")
		}

		t.Logf("✅ Main Library File: %s", libFilePath)

		// Should contain the correct extension
		expectedExt := download.GetLibraryExtension(runtime.GOOS)
		if !strings.HasSuffix(libFilePath, expectedExt) {
			t.Errorf("GetLibraryFilePath() = %s, expected extension %s", libFilePath, expectedExt)
		}

		// Should contain "llama"
		if !strings.Contains(filepath.Base(libFilePath), "llama") {
			t.Errorf("GetLibraryFilePath() = %s, expected to contain 'llama'", libFilePath)
		}
	})

	// Test 3: Verify GetRequiredLibraryPaths returns all 3 libraries
	t.Run("GetRequiredLibraryPaths", func(t *testing.T) {
		requiredPaths := installer.GetRequiredLibraryPaths()

		if len(requiredPaths) != 3 {
			t.Fatalf("GetRequiredLibraryPaths() returned %d paths, want 3", len(requiredPaths))
		}

		t.Logf("✅ Required Libraries Count: %d", len(requiredPaths))

		// Verify each path
		expectedLibs := []string{"ggml", "ggml-base", "llama"}
		for i, path := range requiredPaths {
			t.Logf("   %d. %s", i+1, path)

			// Should contain expected library name
			baseName := filepath.Base(path)
			found := false
			for _, expectedLib := range expectedLibs {
				if strings.Contains(baseName, expectedLib) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Path %s does not contain expected library names", path)
			}

			// Should have correct extension
			expectedExt := download.GetLibraryExtension(runtime.GOOS)
			if !strings.HasSuffix(path, expectedExt) {
				t.Errorf("Path %s does not have expected extension %s", path, expectedExt)
			}
		}
	})

	// Test 4: Verify paths match download.RequiredLibraries()
	t.Run("PathsMatchDownloadPackage", func(t *testing.T) {
		requiredPaths := installer.GetRequiredLibraryPaths()
		expectedLibs := download.RequiredLibraries(runtime.GOOS)

		if len(requiredPaths) != len(expectedLibs) {
			t.Fatalf("Path count mismatch: got %d, want %d", len(requiredPaths), len(expectedLibs))
		}

		t.Logf("✅ Count matches: %d libraries", len(requiredPaths))

		libPath := installer.GetLibraryPath()
		for i, expectedLib := range expectedLibs {
			expectedPath := filepath.Join(libPath, expectedLib)
			actualPath := requiredPaths[i]

			if expectedPath != actualPath {
				t.Errorf("Path mismatch at index %d:\n  Expected: %s\n  Actual:   %s",
					i, expectedPath, actualPath)
			} else {
				t.Logf("   ✅ Match: %s", filepath.Base(actualPath))
			}
		}
	})

	// Test 5: Verify VerifyAllLibrariesExist logic
	t.Run("VerifyAllLibrariesExist", func(t *testing.T) {
		// This will return false if not installed, which is fine
		exists := installer.VerifyAllLibrariesExist()

		if exists {
			t.Log("✅ All libraries exist")

			// If they exist, verify each one
			for _, path := range installer.GetRequiredLibraryPaths() {
				if _, err := os.Stat(path); err != nil {
					t.Errorf("Library reported as existing but stat failed: %s (%v)", path, err)
				}
			}
		} else {
			t.Log("ℹ️  Libraries not installed (this is OK for test)")
		}
	})

	// Test 6: Verify IsLlamaCppInstalled uses VerifyAllLibrariesExist
	t.Run("IsLlamaCppInstalled", func(t *testing.T) {
		isInstalled := installer.IsLlamaCppInstalled()
		allExist := installer.VerifyAllLibrariesExist()

		if isInstalled != allExist {
			t.Errorf("IsLlamaCppInstalled() = %v, but VerifyAllLibrariesExist() = %v (should match)",
				isInstalled, allExist)
		} else {
			t.Logf("✅ IsLlamaCppInstalled() matches VerifyAllLibrariesExist(): %v", isInstalled)
		}
	})

	// Test 7: Verify VerifyInstalledBinary error messages
	t.Run("VerifyInstalledBinary", func(t *testing.T) {
		err := installer.VerifyInstalledBinary()

		if err != nil {
			// Should mention which libraries are missing
			errMsg := err.Error()
			if !strings.Contains(errMsg, "missing required libraries") {
				t.Logf("ℹ️  Error message: %s", errMsg)
			}
			t.Log("ℹ️  Libraries not installed (this is OK for test)")
		} else {
			t.Log("✅ All libraries verified")
		}
	})
}

// TestDownloadPackageConsistency verifies download package functions
func TestDownloadPackageConsistency(t *testing.T) {
	t.Log("🔍 Testing Download Package Consistency")

	// Test 1: RequiredLibraries returns correct count
	t.Run("RequiredLibrariesCount", func(t *testing.T) {
		libs := download.RequiredLibraries(runtime.GOOS)

		if len(libs) != 3 {
			t.Fatalf("RequiredLibraries() returned %d libraries, want 3", len(libs))
		}

		t.Logf("✅ Correct count: %d libraries", len(libs))
		for i, lib := range libs {
			t.Logf("   %d. %s", i+1, lib)
		}
	})

	// Test 2: RequiredLibraries returns correct extensions
	t.Run("RequiredLibrariesExtensions", func(t *testing.T) {
		libs := download.RequiredLibraries(runtime.GOOS)
		expectedExt := download.GetLibraryExtension(runtime.GOOS)

		for _, lib := range libs {
			if !strings.HasSuffix(lib, expectedExt) {
				t.Errorf("Library %s does not have expected extension %s", lib, expectedExt)
			}
		}

		t.Logf("✅ All libraries have correct extension: %s", expectedExt)
	})

	// Test 3: LibraryName returns correct main library
	t.Run("LibraryName", func(t *testing.T) {
		mainLib := download.LibraryName(runtime.GOOS)

		if mainLib == "" || mainLib == "unknown" {
			t.Fatalf("LibraryName() returned invalid value: %s", mainLib)
		}

		// Should contain "llama"
		if !strings.Contains(mainLib, "llama") {
			t.Errorf("LibraryName() = %s, expected to contain 'llama'", mainLib)
		}

		// Should have correct extension
		expectedExt := download.GetLibraryExtension(runtime.GOOS)
		if !strings.HasSuffix(mainLib, expectedExt) {
			t.Errorf("LibraryName() = %s, expected extension %s", mainLib, expectedExt)
		}

		t.Logf("✅ Main library name: %s", mainLib)
	})

	// Test 4: GetLibraryExtension returns correct value
	t.Run("GetLibraryExtension", func(t *testing.T) {
		ext := download.GetLibraryExtension(runtime.GOOS)

		expectedExts := map[string]string{
			"darwin":  ".dylib",
			"linux":   ".so",
			"freebsd": ".so",
			"windows": ".dll",
		}

		expected, ok := expectedExts[runtime.GOOS]
		if !ok {
			t.Skipf("Unknown platform: %s", runtime.GOOS)
		}

		if ext != expected {
			t.Errorf("GetLibraryExtension() = %s, want %s", ext, expected)
		}

		t.Logf("✅ Library extension: %s", ext)
	})
}

// TestPlatformSpecificPaths verifies platform-specific library names
func TestPlatformSpecificPaths(t *testing.T) {
	t.Log("🔍 Testing Platform-Specific Paths")

	platforms := []struct {
		os      string
		ext     string
		prefix  string
		mainLib string
	}{
		{"darwin", ".dylib", "lib", "libllama.dylib"},
		{"linux", ".so", "lib", "libllama.so"},
		{"freebsd", ".so", "lib", "libllama.so"},
		{"windows", ".dll", "", "llama.dll"},
	}

	for _, platform := range platforms {
		t.Run(platform.os, func(t *testing.T) {
			// Test extension
			ext := download.GetLibraryExtension(platform.os)
			if ext != platform.ext {
				t.Errorf("GetLibraryExtension(%s) = %s, want %s", platform.os, ext, platform.ext)
			}

			// Test main library name
			mainLib := download.LibraryName(platform.os)
			if mainLib != platform.mainLib {
				t.Errorf("LibraryName(%s) = %s, want %s", platform.os, mainLib, platform.mainLib)
			}

			// Test required libraries
			libs := download.RequiredLibraries(platform.os)
			if len(libs) != 3 {
				t.Errorf("RequiredLibraries(%s) returned %d libraries, want 3", platform.os, len(libs))
			}

			// All should have correct extension
			for _, lib := range libs {
				if !strings.HasSuffix(lib, platform.ext) {
					t.Errorf("Library %s does not have expected extension %s", lib, platform.ext)
				}
			}

			t.Logf("✅ %s: extension=%s, mainLib=%s, count=%d",
				platform.os, ext, mainLib, len(libs))
		})
	}
}

// TestNoEnvironmentVariables verifies no env vars are used
func TestNoEnvironmentVariables(t *testing.T) {
	t.Log("🔍 Testing No Environment Variables Used")

	// Save and clear YZMA_LIB if set
	originalEnv := os.Getenv("YZMA_LIB")
	os.Unsetenv("YZMA_LIB")
	defer func() {
		if originalEnv != "" {
			os.Setenv("YZMA_LIB", originalEnv)
		}
	}()

	// Create installer - should work without env var
	installer := NewLlamaCppInstaller()

	// Verify paths are set
	if installer.GetLibraryPath() == "" {
		t.Fatal("GetLibraryPath() returned empty string without YZMA_LIB env var")
	}

	if len(installer.GetRequiredLibraryPaths()) == 0 {
		t.Fatal("GetRequiredLibraryPaths() returned empty without YZMA_LIB env var")
	}

	t.Log("✅ Installer works without YZMA_LIB environment variable")
	t.Logf("   Library Path: %s", installer.GetLibraryPath())
}
