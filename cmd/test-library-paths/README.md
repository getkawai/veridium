# Library Path Management Test

This test program verifies that library paths are correctly managed after downloading llama.cpp.

## What It Tests

1. ✅ **Installer Creation** - Creates `LlamaCppInstaller` and verifies paths
2. ✅ **Installation Check** - Checks if llama.cpp is already installed
3. ✅ **Download** - Downloads llama.cpp if not installed
4. ✅ **Library Verification** - Verifies all 3 required libraries exist:
   - `libggml.dylib` (or `.so`/`.dll`)
   - `libggml-base.dylib` (or `.so`/`.dll`)
   - `libllama.dylib` (or `.so`/`.dll`)
5. ✅ **Path Consistency** - Verifies paths match `download.RequiredLibraries()`
6. ✅ **Library Loading** - Tests loading libraries with `yzma.Load()`
7. ✅ **Backend Init** - Tests initializing llama.cpp backend
8. ✅ **Service Integration** - Tests `LibraryService` uses correct paths
9. ✅ **Path Matching** - Verifies service and installer use same paths

## Running the Test

### Option 1: Using the shell script (Recommended)

```bash
chmod +x test-library-paths.sh
./test-library-paths.sh
```

### Option 2: Direct go run

```bash
go run ./cmd/test-library-paths/main.go
```

### Option 3: Build and run

```bash
go build -o test-library-paths ./cmd/test-library-paths/main.go
./test-library-paths
```

## Expected Output

```
🧪 Testing Library Path Management
===================================

📦 Step 1: Creating installer...
   Binary Path: /Users/username/.llama-cpp/bin
   Metadata Path: /Users/username/.llama-cpp/metadata
   Models Dir: /Users/username/.llama-cpp/models

🔍 Step 2: Checking installation status...
   Installed: true
   Installed version: b1234

🔍 Step 3: Verifying all required libraries...
   All libraries present: true

📂 Step 4: Getting library paths...
   Library Directory: /Users/username/.llama-cpp/bin
   Main Library File: /Users/username/.llama-cpp/bin/libllama.dylib
   Required Libraries (3):
     1. ✅ /Users/username/.llama-cpp/bin/libggml.dylib
     2. ✅ /Users/username/.llama-cpp/bin/libggml-base.dylib
     3. ✅ /Users/username/.llama-cpp/bin/libllama.dylib

🔍 Step 5: Verifying paths match download package...
   Expected libraries from download package: [libggml.dylib libggml-base.dylib libllama.dylib]
   ✅ Match: libggml.dylib
   ✅ Match: libggml-base.dylib
   ✅ Match: libllama.dylib

🔧 Step 6: Testing library loading with yzma...
   Loading from: /Users/username/.llama-cpp/bin
   ✅ Libraries loaded successfully!

🚀 Step 7: Initializing llama.cpp backend...
   ✅ Backend initialized successfully!

🧪 Step 8: Testing with LibraryService...
   ✅ LibraryService created successfully!

🔍 Step 9: Verifying service uses correct paths...
   ✅ Models directory matches: /Users/username/.llama-cpp/models

==================================================
✅ ALL TESTS PASSED!
==================================================

📊 Summary:
   Platform: darwin/arm64
   Library Path: /Users/username/.llama-cpp/bin
   Required Libraries: 3
   All Libraries Present: ✅
   Library Loading: ✅
   Backend Initialization: ✅
   LibraryService Integration: ✅

🎉 Library path management is working correctly!
```

## What This Verifies

### 1. No Environment Variables Needed
The test runs without setting `YZMA_LIB` environment variable, proving that the programmatic path management works.

### 2. Correct Library Detection
Verifies that `download.RequiredLibraries()` returns the correct libraries for the current platform:
- **macOS**: `.dylib` files
- **Linux**: `.so` files
- **Windows**: `.dll` files

### 3. Path Consistency
Ensures that:
- `installer.GetLibraryPath()` returns the correct directory
- `installer.GetRequiredLibraryPaths()` returns correct full paths
- All paths use the same base directory
- Paths match what `download.RequiredLibraries()` expects

### 4. Library Loading Works
Tests that `yzma.Load()` can successfully load libraries from the installer-provided path.

### 5. Service Integration
Verifies that `LibraryService` correctly uses paths from the installer.

## Troubleshooting

### Test Fails at Step 2 (Not Installed)
```
❌ Failed to download: ...
```

**Solution:** Check network connectivity or try manually downloading:
```bash
go run ./examples/installer/main.go
```

### Test Fails at Step 3 (Missing Libraries)
```
❌ Not all libraries found!
Error: missing required libraries: libggml.dylib
```

**Solution:** Re-download llama.cpp:
```bash
rm -rf ~/.llama-cpp/bin
go run ./cmd/test-library-paths/main.go
```

### Test Fails at Step 5 (Path Mismatch)
```
❌ Path mismatch at index 0:
   Expected: /path/to/libggml.dylib
   Actual:   /other/path/libggml.dylib
```

**Solution:** This indicates a bug in the code. Check:
- `installer.GetRequiredLibraryPaths()` implementation
- `download.RequiredLibraries()` implementation

### Test Fails at Step 6 (Library Loading)
```
❌ Failed to load library: could not load "ggml": ...
```

**Solution:** Verify library files are not corrupted:
```bash
ls -lh ~/.llama-cpp/bin/
file ~/.llama-cpp/bin/libggml.dylib
```

## Platform-Specific Notes

### macOS
- Requires Xcode Command Line Tools
- Libraries use `.dylib` extension
- May require security permissions for downloaded libraries

### Linux
- Libraries use `.so` extension
- May require `libstdc++` or other system libraries

### Windows
- Libraries use `.dll` extension
- May require Visual C++ Redistributable

## See Also

- [LIBRARY_PATHS.md](../../pkg/yzma/download/LIBRARY_PATHS.md) - Library path documentation
- [MIGRATION_NO_ENV_VAR.md](../../pkg/yzma/MIGRATION_NO_ENV_VAR.md) - Migration guide
- [installer.go](../../internal/llama/installer.go) - Installer implementation

