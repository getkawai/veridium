# Stable Diffusion Model Management

This package provides model catalog and hardware-aware model selection for Stable Diffusion.

## Architecture

Inspired by `pkg/kronk/model/` patterns but adapted for Stable Diffusion's needs:

- **types.go**: Core data structures (ModelSpec, HardwareSpecs)
- **catalog.go**: Available model definitions
- **selector.go**: Hardware-aware model selection logic

## Key Differences from Kronk

While this package borrows patterns from kronk's model management, it serves a different purpose:

| Aspect | Kronk (`pkg/kronk/model/`) | SD Models (`pkg/stablediffusion/models/`) |
|--------|---------------------------|------------------------------------------|
| Purpose | LLM runtime management | Model catalog & selection |
| Scope | Model loading + inference | Model specs + hardware matching |
| Dependencies | llama.cpp bindings | None (pure Go) |
| Runtime | Manages loaded models | Provides model metadata |

## Usage

```go
import "github.com/kawai-network/veridium/pkg/stablediffusion/models"

// Get hardware specs (from pkg/hardware or similar)
specs := &models.HardwareSpecs{
    AvailableRAM: 16,
    GPUMemory:    8,
}

// Select optimal model
model := models.SelectOptimalModel(specs)

// Get all available models
allModels := models.GetAvailableModels()
```

## Design Rationale

**Why not reuse kronk directly?**

1. **Different lifecycle**: Kronk manages runtime model loading; SD needs model acquisition
2. **Different dependencies**: Kronk requires llama.cpp; SD models are pure metadata
3. **Different patterns**: Kronk focuses on inference; SD focuses on model selection

**What patterns were borrowed?**

1. Hardware-aware selection logic
2. Model specification structure
3. Clean separation of concerns (types, catalog, selector)

## Integration

This package is used by:
- `pkg/stablediffusion/local/` - For model path resolution
- `internal/image/models.go` - Backward compatibility wrapper
- Future: `internal/image/manager.go` - For model download decisions
