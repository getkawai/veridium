# Library Path Management

This document explains how to properly manage llama.cpp library paths for library-based usage (via yzma).

## Overview

llama.cpp requires **three separate library files** to function:

1. **libggml** - Core GGML library
2. **libggml-base** - Base GGML functionality
3. **libllama** - Main llama.cpp library

All three libraries must be present in the same directory for llama.cpp to work correctly.

## Platform-Specific Library Names

### Linux / FreeBSD
```
libggml.so
libggml-base.so
libllama.so
```

### macOS (Darwin)
```
libggml.dylib
libggml-base.dylib
libllama.dylib
```

### Windows
```
ggml.dll
ggml-base.dll
llama.dll
```

## API Functions

### `RequiredLibraries(os string) []string`

Returns all required library filenames for a given OS.

**Example:**
```go
libs := download.RequiredLibraries("darwin")
// Returns: ["libggml.dylib", "libggml-base.dylib", "libllama.dylib"]
```

### `LibraryName(os string) string`

Returns only the main llama library filename.

**Example:**
```go
mainLib := download.LibraryName("linux")
// Returns: "libllama.so"
```

### `GetLibraryExtension(os string) string`

Returns the library file extension for a given OS.

**Example:**
```go
ext := download.GetLibraryExtension("windows")
// Returns: ".dll"
```

## Usage in LlamaCppInstaller

The `LlamaCppInstaller` provides high-level methods for working with library paths:

### `GetLibraryPath() string`

Returns the **directory** containing all libraries. This is what you should pass to `llama.Load()`.

**Example:**
```go
installer := llama.NewLlamaCppInstaller()
libPath := installer.GetLibraryPath()
// Returns: "/Users/username/.llama-cpp/bin"

// Use with yzma:
llama.Load(libPath)
```

### `GetLibraryFilePath() string`

Returns the full path to the main llama library file.

**Example:**
```go
mainLibPath := installer.GetLibraryFilePath()
// macOS: "/Users/username/.llama-cpp/bin/libllama.dylib"
// Linux: "/home/username/.llama-cpp/bin/libllama.so"
// Windows: "C:\Users\username\.llama-cpp\bin\llama.dll"
```

### `GetRequiredLibraryPaths() []string`

Returns full paths to all required library files.

**Example:**
```go
allPaths := installer.GetRequiredLibraryPaths()
// macOS returns:
// [
//   "/Users/username/.llama-cpp/bin/libggml.dylib",
//   "/Users/username/.llama-cpp/bin/libggml-base.dylib",
//   "/Users/username/.llama-cpp/bin/libllama.dylib"
// ]
```

### `VerifyAllLibrariesExist() bool`

Checks if all required libraries are present.

**Example:**
```go
if installer.VerifyAllLibrariesExist() {
    fmt.Println("✅ All libraries present")
} else {
    fmt.Println("❌ Missing libraries")
}
```

## Complete Example: LibraryService Initialization

```go
package llama

import (
    "fmt"
    "github.com/kawai-network/veridium/pkg/yzma/llama"
)

func (s *LibraryService) InitializeLibrary() error {
    // Get library path from installer
    libPath := s.manager.GetLibraryPath()
    
    // Verify all required libraries exist
    if !s.manager.VerifyAllLibrariesExist() {
        return fmt.Errorf("llama.cpp libraries not found in %s", libPath)
    }
    
    // Log which libraries were found
    requiredPaths := s.manager.GetRequiredLibraryPaths()
    for _, path := range requiredPaths {
        log.Printf("  ✓ Found: %s", filepath.Base(path))
    }
    
    // Load the library (pass directory, not file path!)
    if err := llama.Load(libPath); err != nil {
        return fmt.Errorf("failed to load llama.cpp: %w", err)
    }
    
    // Initialize llama.cpp backend
    llama.Init()
    
    return nil
}
```

## Common Mistakes to Avoid

### ❌ Wrong: Passing file path to llama.Load()
```go
// DON'T DO THIS
mainLib := installer.GetLibraryFilePath()
llama.Load(mainLib) // ERROR: expects directory, not file!
```

### ✅ Correct: Passing directory path to llama.Load()
```go
// DO THIS
libDir := installer.GetLibraryPath()
llama.Load(libDir) // Correct: directory containing all libraries
```

### ❌ Wrong: Only checking main library
```go
// DON'T DO THIS
mainLib := filepath.Join(libPath, "libllama.dylib")
if _, err := os.Stat(mainLib); err == nil {
    // Missing libggml and libggml-base!
}
```

### ✅ Correct: Checking all required libraries
```go
// DO THIS
if installer.VerifyAllLibrariesExist() {
    // All three libraries are present
}
```

## Directory Structure

After installation, your directory should look like this:

```
~/.llama-cpp/
├── bin/                      ← GetLibraryPath() returns this
│   ├── libggml.dylib        ← Required library 1
│   ├── libggml-base.dylib   ← Required library 2
│   └── libllama.dylib       ← Required library 3 (main)
├── models/                   ← GetModelsDirectory() returns this
│   ├── qwen2.5-0.5b-instruct-q4_k_m.gguf
│   └── all-MiniLM-L6-v2-Q4_K_M.gguf
└── metadata/
    └── installed-version.json
```

## Testing

Run tests to verify library path functions:

```bash
# Test download package functions
go test ./pkg/yzma/download -run TestRequiredLibraries
go test ./pkg/yzma/download -run TestGetLibraryExtension

# Test installer functions
go test ./internal/llama -run TestGetLibraryPath
go test ./internal/llama -run TestVerifyAllLibrariesExist
```

## Troubleshooting

### "Failed to load llama.cpp library"

**Cause:** Missing one or more required libraries.

**Solution:**
```go
// Check which libraries are missing
if err := installer.VerifyInstalledBinary(); err != nil {
    log.Printf("Missing libraries: %v", err)
    // Re-download
    installer.DownloadRelease("", nil)
}
```

### "Library loaded but crashes on Init()"

**Cause:** Library version mismatch or corrupted files.

**Solution:**
```go
// Clean up and re-download
installer.CleanupPartialDownloads()
installer.DownloadRelease("", nil)
```

## See Also

- [DOWNLOAD_FEATURES.md](../../../internal/llama/DOWNLOAD_FEATURES.md) - Download strategy
- [LIBRARY_USAGE.md](../../../internal/llama/LIBRARY_USAGE.md) - Library-based usage guide
- [yzma documentation](https://github.com/hybridgroup/yzma) - Low-level library bindings

