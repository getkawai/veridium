# Stable Diffusion Migration Guide

Branch: `refactor/migrate-stablediffusion-to-pkg`

## Overview

Migrating Stable Diffusion from `internal/image` to `pkg/stablediffusion` for better separation of concerns.

## Phase 1: Remote Generation ✅

**Status:** Complete (7 files, +827/-229 lines, 10/10 tests passing)

```
pkg/stablediffusion/remote/
├── types.go, gemini.go, pollinations.go
├── generator.go, generator_test.go
└── README.md
```

## Phase 2: Local Generation ✅

**Status:** Complete (6 files, +766/-202 lines, 13/13 tests passing)

```
pkg/stablediffusion/local/
├── types.go, executor.go, generator.go
├── generator_test.go
└── README.md
```

## Usage

### Remote
```go
// OLD: import "github.com/kawai-network/veridium/internal/image"
// NEW:
import "github.com/kawai-network/veridium/pkg/stablediffusion/remote"
gen := remote.NewGenerator()
```

### Local
```go
// OLD: gen := image.NewLocalGenerator(engine)
// NEW:
import "github.com/kawai-network/veridium/pkg/stablediffusion/local"
gen := local.NewGenerator(binaryPath, modelsPath)
```

## Phase 3: Model Management ✅

**Status:** Complete (7 files, +363/-558 lines, 3/3 tests passing)

**Quality Fixes:** 8 P1/P2/P3 violations resolved (commit 32791e78)

### Analysis: Kronk vs SD Model Management

**Question**: Can we reuse `pkg/kronk/model/` for Stable Diffusion model management?

**Answer**: No, but we borrowed patterns.

| Aspect | Kronk | SD Models |
|--------|-------|-----------|
| Purpose | LLM runtime | Model catalog |
| Scope | Loading + inference | Specs + selection |
| Dependencies | llama.cpp | Pure Go |
| Focus | Execution | Acquisition |

**Patterns Borrowed**:
- Hardware-aware selection logic
- Model specification structure
- Clean separation of concerns

### Implementation

```
pkg/stablediffusion/models/
├── types.go (ModelSpec, HardwareSpecs, ModelType, ModelFormat, ModelInfo)
├── catalog.go (6 available models)
├── selector.go (hardware-aware selection)
├── selector_test.go (3/3 passing)
├── detector.go (existing - auto-detect models)
├── downloader.go (existing - download from HF/Civitai)
└── README.md
```

**Key Features**:
- Hardware-aware model selection (RAM + VRAM)
- Score-based selection (prefers best model that fits)
- Backward compatibility via `internal/image/models.go` wrapper

### Test Results

```bash
go test -v ./pkg/stablediffusion/models
# ✅ 3/3 tests passing
# - TestSelectOptimalModel (5 scenarios)
# - TestGetAvailableModels
# - TestModelSpecFields
```

### Usage

```go
import "github.com/kawai-network/veridium/pkg/stablediffusion/models"

specs := &models.HardwareSpecs{
    AvailableRAM: 16,
    GPUMemory:    8,
}
model := models.SelectOptimalModel(specs)
// Returns: sdxl-base-f16 (best model that fits)
```

## Next Phases

- **Phase 4:** Service layer refactoring ⏳
- **Phase 5:** Cleanup ⏳

## Quality Fixes (Commit 32791e78)

Fixed 8 violations across packages:

**Critical (P1)**:
- `executor.go`: Environment variables wiped (PATH/HOME/TEMP) - now inherits os.Environ()
- `executor.go`: Windows PATH expansion broken - now uses os.PathListSeparator

**Process Management (P2)**:
- `local.go`: Broken cleanup - now uses NewGeneratorWithExecutor for process tracking

**Library Code Quality (P2/P3)**:
- `selector.go`: Removed global log statements
- `generator.go`: Removed global log statements
- `generator.go`: Case-sensitive extension check - now handles .CKPT/.ckpt

**Model Selection (P2)**:
- `catalog.go`: Reordered to prioritize SD-Turbo over SD1.4 for low-end systems
- `catalog.go`: Fixed 3 incorrect URLs (q4_0/q8_0 models pointed to full-precision files)

## Testing

```bash
go test ./pkg/stablediffusion/remote/... -v  # 10/10 passing
go test ./pkg/stablediffusion/local/... -v   # 13/13 passing
go test ./pkg/stablediffusion/models/... -v  # 3/3 passing
go build ./internal/image/...                # Backward compatible
```

## Files Status

**New Packages:**
- `pkg/stablediffusion/remote/*` (Phase 1) ✅
- `pkg/stablediffusion/local/*` (Phase 2) ✅
- `pkg/stablediffusion/models/*` (Phase 3) ✅

**Modified (Backward Compatible):**
- `internal/image/remote.go` - wraps remote package
- `internal/image/local.go` - wraps local package
- `internal/image/models.go` - wraps models package

**To Migrate:**
- `internal/image/manager.go` - binary management (Phase 4)
- `internal/image/generation.go` - types consolidation (Phase 4)

**Keep:**
- `internal/image/service.go` - business logic & database
