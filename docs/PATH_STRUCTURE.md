# Path Structure Documentation

## Overview

Veridium uses a centralized path management system through `internal/paths` to ensure consistent data storage across development and production environments.

## Directory Structure

### Development Mode (running from terminal)
```
./data/
├── veridium.db                          # Main SQLite database
├── jarvis/                              # Blockchain/Wallet data
│   ├── keystores/                       # Encrypted wallet files
│   ├── cache.json
│   └── secrets.json
├── libraries/                           # Shared libraries
│   ├── llama-cli
│   ├── llama-server
│   ├── libllama.{so|dylib|dll}
│   └── stable-diffusion/                # SD library subdirectory
│       ├── bin/
│       ├── checksums/
│       └── metadata/
├── models/                              # All AI models (unified flat structure)
│   ├── Qwen/
│   │   └── Qwen3-8B-GGUF/
│   │       └── Qwen3-8B-Q8_0.gguf
│   ├── ggerganov/
│   │   └── whisper.cpp/
│   │       ├── ggml-base.bin
│   │       └── ggml-small.bin
│   ├── CompVis/
│   │   └── stable-diffusion-v-1-4-original/
│   │       └── sd-v1-4.ckpt
│   ├── runwayml/
│   │   └── stable-diffusion-v1-5/
│   │       └── v1-5-pruned-emaonly.safetensors
│   ├── stabilityai/
│   │   └── stable-diffusion-xl-base-1.0/
│   │       └── sd_xl_base_1.0.safetensors
│   └── .index.yaml                      # Unified model index (all types)
├── outputs/                             # Generated content
│   └── stable-diffusion/                # SD generated images
│       ├── image_001.png
│       └── image_002.png
├── catalogs/                            # Model catalogs
│   └── models.json
└── templates/                           # Chat templates
    └── llama-3.2.json
```

### Production Mode (packaged app)

**macOS:**
```
~/Library/Application Support/Kawai/
├── veridium.db
├── jarvis/
│   ├── keystores/
│   ├── cache.json
│   └── secrets.json
├── libraries/
│   ├── llama-cli
│   ├── llama-server
│   ├── libllama.dylib
│   └── stable-diffusion/
│       ├── bin/
│       ├── checksums/
│       └── metadata/
├── models/
│   ├── ggerganov/
│   │   └── whisper.cpp/
│   │       ├── ggml-base.bin
│   │       └── ggml-small.bin
│   ├── CompVis/
│   │   └── stable-diffusion-v-1-4-original/
│   │       └── sd-v1-4.ckpt
│   ├── runwayml/
│   │   └── stable-diffusion-v1-5/
│   │       └── v1-5-pruned-emaonly.safetensors
│   ├── stabilityai/
│   │   └── stable-diffusion-xl-base-1.0/
│   │       └── sd_xl_base_1.0.safetensors
│   └── Qwen/
│       └── Qwen3-8B-GGUF/
│           └── Qwen3-8B-Q8_0.gguf
├── outputs/
│   └── stable-diffusion/
│       └── image_001.png
├── catalogs/
│   └── models.json
└── templates/
    └── llama-3.2.json
```

**Windows:**
```
%APPDATA%\Kawai\
├── veridium.db
├── jarvis\
│   ├── keystores\
│   ├── cache.json
│   └── secrets.json
├── libraries\
│   ├── llama-cli.exe
│   ├── llama-server.exe
│   ├── llama.dll
│   └── stable-diffusion\
│       ├── bin\
│       ├── checksums\
│       └── metadata\
├── models\
│   ├── ggerganov\
│   │   └── whisper.cpp\
│   │       ├── ggml-base.bin
│   │       └── ggml-small.bin
│   ├── CompVis\
│   │   └── stable-diffusion-v-1-4-original\
│   │       └── sd-v1-4.ckpt
│   ├── runwayml\
│   │   └── stable-diffusion-v1-5\
│   │       └── v1-5-pruned-emaonly.safetensors
│   ├── stabilityai\
│   │   └── stable-diffusion-xl-base-1.0\
│   │       └── sd_xl_base_1.0.safetensors
│   └── Qwen\
│       └── Qwen3-8B-GGUF\
│           └── Qwen3-8B-Q8_0.gguf
├── outputs\
│   └── stable-diffusion\
│       └── image_001.png
├── catalogs\
│   └── models.json
└── templates\
    └── llama-3.2.json
```

**Linux:**
```
~/.config/Kawai/
├── veridium.db
├── jarvis/
│   ├── keystores/
│   ├── cache.json
│   └── secrets.json
├── libraries/
│   ├── llama-cli
│   ├── llama-server
│   ├── libllama.so
│   └── stable-diffusion/
│       ├── bin/
│       ├── checksums/
│       └── metadata/
├── models/
│   ├── ggerganov/
│   │   └── whisper.cpp/
│   │       ├── ggml-base.bin
│   │       └── ggml-small.bin
│   ├── CompVis/
│   │   └── stable-diffusion-v-1-4-original/
│   │       └── sd-v1-4.ckpt
│   ├── runwayml/
│   │   └── stable-diffusion-v1-5/
│   │       └── v1-5-pruned-emaonly.safetensors
│   ├── stabilityai/
│   │   └── stable-diffusion-xl-base-1.0/
│   │       └── sd_xl_base_1.0.safetensors
│   └── Qwen/
│       └── Qwen3-8B-GGUF/
│           └── Qwen3-8B-Q8_0.gguf
├── outputs/
│   └── stable-diffusion/
│       └── image_001.png
├── catalogs/
│   └── models.json
└── templates/
    └── llama-3.2.json
```

## Path Functions

### Core Paths

```go
import "github.com/kawai-network/veridium/internal/paths"

// Base directory (platform-specific)
paths.Base()              // ./data/ or ~/Library/Application Support/Kawai/

// Database paths
paths.Database()          // {Base}/veridium.db
paths.DuckDB()           // {Base}/duckdb.db

// File storage
paths.FileBase()         // {Base}/files/
paths.KBAssets()         // {Base}/kb-assets/
```

### Blockchain/Wallet Paths (Jarvis)

```go
paths.Jarvis()                  // {Base}/jarvis/
paths.JarvisKeystores()         // {Base}/jarvis/keystores/
paths.JarvisNetworks()          // {Base}/jarvis/networks/
paths.JarvisAddressBookDB()     // {Base}/jarvis/addressbook.duckdb
paths.JarvisCache()             // {Base}/jarvis/cache.json
paths.JarvisSecrets()           // {Base}/jarvis/secrets.json
```

### AI/ML Paths

```go
// Unified model storage (all models use {author}/{repo}/ structure)
paths.Models()                    // {Base}/models/
paths.ModelPath(huggingfaceURL)  // Extract author/repo from URL

// Shared libraries (llama.cpp, stable-diffusion)
paths.Libraries()                // {Base}/libraries/

// Stable Diffusion specific paths
paths.StableDiffusionOutputs()   // {Base}/outputs/stable-diffusion/
paths.StableDiffusionBin()       // {Base}/libraries/stable-diffusion/bin/
paths.StableDiffusionChecksums() // {Base}/libraries/stable-diffusion/checksums/
paths.StableDiffusionMetadata()  // {Base}/libraries/stable-diffusion/metadata/

// Supporting data
paths.Catalogs()                 // {Base}/catalogs/
paths.Templates()                // {Base}/templates/
```

## Migration from Old Structure

### Old Structure (Deprecated)
```
{Base}/
└── node/
    ├── libraries/
    │   ├── llama-cli
    │   └── libllama.{so|dylib|dll}
    ├── models/
    │   └── llama-3.2-1b-instruct.gguf
    ├── whisper-models/              # Separate directory
    │   └── ggml-base.bin
    ├── catalogs/
    │   └── models.json
    └── templates/
        └── llama-3.2.json
```

**Also deprecated (type-specific subdirectories):**
```
{Base}/
├── models/
│   ├── llm/                         # Old: type-specific subdirectory
│   │   └── Qwen3-8B-Q8_0.gguf
│   ├── diffusion/                   # Old: type-specific subdirectory
│   │   └── sd_v1.5.safetensors
│   └── audio/                       # Old: type-specific subdirectory
│       └── ggml-base.bin
```

**Also deprecated (hardcoded paths):**
```
~/.stable-diffusion/                 # Hardcoded path
├── bin/
├── checksums/
├── metadata/
├── models/
└── outputs/
```

### New Structure (Unified {author}/{repo}/)
```
{Base}/
├── libraries/                       # Moved up one level
│   ├── llama-cli
│   ├── libllama.{so|dylib|dll}
│   └── stable-diffusion/            # SD library organized
│       ├── bin/
│       ├── checksums/
│       └── metadata/
├── models/                          # Unified flat structure for ALL models
│   ├── ggerganov/                   # Whisper models
│   │   └── whisper.cpp/
│   │       └── ggml-base.bin
│   ├── CompVis/                     # SD models
│   │   └── stable-diffusion-v-1-4-original/
│   │       └── sd-v1-4.ckpt
│   ├── runwayml/                    # SD models
│   │   └── stable-diffusion-v1-5/
│   │       └── v1-5-pruned-emaonly.safetensors
│   ├── stabilityai/                 # SD models
│   │   └── stable-diffusion-xl-base-1.0/
│   │       └── sd_xl_base_1.0.safetensors
│   ├── Qwen/                        # LLM models
│   │   └── Qwen3-8B-GGUF/
│   │       └── Qwen3-8B-Q8_0.gguf
│   └── .index.yaml                  # Single unified index with type metadata
├── outputs/                         # Generated content
│   └── stable-diffusion/
├── catalogs/                        # Moved up one level
│   └── models.json
└── templates/                       # Moved up one level
    └── llama-3.2.json
```

### Deprecated Functions (Backward Compatible)

These functions still work but are deprecated:

```go
// Deprecated: Use Base() instead
paths.Node()

// Deprecated: Use Models() instead
paths.NodeModels()

// Deprecated: Use Libraries() instead
paths.NodeLibraries()

// Deprecated: Use Catalogs() instead
paths.NodeCatalogs()

// Deprecated: Use Templates() instead
paths.NodeTemplates()
```

**Note:** Type-specific model path functions (`ModelsWhisper()`, `ModelsStableDiffusion()`) have been removed in favor of the unified `paths.Models()` + `paths.ModelPath(url)` approach.

## Benefits of New Structure

1. **Simpler hierarchy**: No unnecessary `node/` nesting
2. **Unified flat models**: All AI models use consistent {author}/{repo}/ structure without type prefixes
3. **Automatic organization**: Models organized by HuggingFace URL structure
4. **Type detection**: Model type automatically detected from filename patterns and stored in index
5. **No naming conflicts**: Different model types can have same filenames
6. **Easier navigation**: Clear author/repo hierarchy
7. **Better scalability**: Supports unlimited models without manual categorization
8. **Single index**: One `.index.yaml` with type metadata for all models

## Usage Examples

### Setting up libraries
```go
import (
    "github.com/kawai-network/veridium/internal/paths"
    "github.com/kawai-network/veridium/pkg/tools/libs"
)

libMgr, err := libs.New(
    libs.WithBasePath(paths.Libraries()),  // Use centralized path
    libs.WithArch(arch),
    libs.WithOS(opSys),
    libs.WithProcessor(processor),
)
```

### Setting up models
```go
import (
    "github.com/kawai-network/veridium/internal/paths"
    "path/filepath"
)

// All models use unified structure - automatically organized by URL
modelsDir := paths.Models()

// Extract author/repo from HuggingFace URL
whisperURL := "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin"
whisperPath, _ := paths.ModelPath(whisperURL)
// Result: {Base}/models/ggerganov/whisper.cpp/

sdURL := "https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main/v1-5-pruned-emaonly.safetensors"
sdPath, _ := paths.ModelPath(sdURL)
// Result: {Base}/models/runwayml/stable-diffusion-v1-5/

llmURL := "https://huggingface.co/Qwen/Qwen3-8B-GGUF/resolve/main/Qwen3-8B-Q8_0.gguf"
llmPath, _ := paths.ModelPath(llmURL)
// Result: {Base}/models/Qwen/Qwen3-8B-GGUF/
```

### Stable Diffusion paths
```go
import "github.com/kawai-network/veridium/internal/paths"

// Models (automatically organized by author/repo from URL)
modelsPath := paths.Models()
// Example: {Base}/models/runwayml/stable-diffusion-v1-5/

// Generated images output
outputPath := paths.StableDiffusionOutputs()

// Binary and metadata
binPath := paths.StableDiffusionBin()
checksumsPath := paths.StableDiffusionChecksums()
metadataPath := paths.StableDiffusionMetadata()
```

### Custom data directory (development)
```go
import "github.com/kawai-network/veridium/internal/paths"

func main() {
    // Set custom data directory before any path access
    paths.SetDataDir("./custom-data")
    
    // Now all paths will use ./custom-data/ as base
    dbPath := paths.Database()              // ./custom-data/veridium.db
    modelsPath := paths.Models()            // ./custom-data/models/
}
```

## Platform Detection

The path system automatically detects if the app is running in:

1. **Development mode**: Uses `./data/` relative to working directory
2. **Packaged mode**: Uses platform-specific user data directories

Detection logic:
- **macOS**: Checks for `.app` bundle structure
- **Windows**: Checks for `resources/` directory or Program Files
- **Linux**: Checks for `resources/` or `/usr/`, `/opt/` paths

## Best Practices

1. **Always use `internal/paths`**: Never hardcode paths or use `os.UserHomeDir()` directly
2. **Use appropriate functions**: Choose the right path function for your use case
3. **Create directories**: Always ensure directories exist before writing files
4. **Cross-platform**: Test path behavior on all target platforms
5. **Migration**: Update old code to use new path functions

## Testing

```go
import (
    "testing"
    "github.com/kawai-network/veridium/internal/paths"
)

func TestPaths(t *testing.T) {
    // Use temporary directory for tests
    tempDir := t.TempDir()
    paths.SetDataDir(tempDir)
    
    // Now all paths use temp directory
    modelsDir := paths.Models()
    // modelsDir = {tempDir}/models/
}
```
