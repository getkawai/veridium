# Migration from eino-ext markdown splitter

This document describes the migration from `github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown` to the local `pkg/mdsplitter` package.

## Why Migrate?

- **Remove external dependency**: Eliminates the need for the eino-ext dependency
- **Simpler architecture**: Focused, lightweight implementation without the eino framework overhead
- **Easier maintenance**: Local code that can be modified as needed
- **No breaking changes**: The core splitting logic remains the same

## What Changed?

### Before (eino-ext)

```go
import (
    "context"
    "github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
    "github.com/cloudwego/eino/schema"
)

// Configure splitter
config := &markdown.HeaderConfig{
    Headers: map[string]string{
        "##":  "h2",
        "###": "h3",
    },
    TrimHeaders: false,
}

splitter, err := markdown.NewHeaderSplitter(ctx, config)
if err != nil {
    // handle error
}

// Create document
doc := &Document{
    ID:      "doc",
    Content: content,
}

// Split
docs, err := splitter.Transform(ctx, []*Document{doc})
if err != nil {
    // handle error
}

// Access results
for _, doc := range docs {
    fmt.Println(doc.Content)
    fmt.Println(doc.MetaData["h2"])
    fmt.Println(doc.MetaData["h3"])
}
```

### After (mdsplitter)

```go
import (
    "github.com/kawai-network/veridium/pkg/mdsplitter"
)

// Configure splitter
config := &mdsplitter.Config{
    Headers: map[string]string{
        "##":  "h2",
        "###": "h3",
    },
    TrimHeaders: false,
}

splitter, err := mdsplitter.New(config)
if err != nil {
    // handle error
}

// Split (no need for Document wrapper)
chunks := splitter.Split(content)

// Access results
for _, chunk := range chunks {
    fmt.Println(chunk.Content)
    fmt.Println(chunk.Metadata["h2"])
    fmt.Println(chunk.Metadata["h3"])
}
```

## Key Differences

1. **No context required**: The new splitter doesn't need a `context.Context` parameter
2. **Direct string input**: No need to wrap content in `Document`
3. **Simpler output**: Returns `[]Chunk` instead of `[]*Document`
4. **Metadata access**: Use `chunk.Metadata` (map[string]string) instead of `doc.MetaData` (map[string]any)

## Files Modified

1. **Created**: `pkg/mdsplitter/splitter.go` - Main implementation
2. **Created**: `pkg/mdsplitter/splitter_test.go` - Unit tests
3. **Created**: `pkg/mdsplitter/README.md` - Package documentation
4. **Modified**: `internal/services/file_chunking.go` - Updated to use mdsplitter
5. **Modified**: `go.mod` - Removed eino-ext markdown dependency

## Testing

All tests pass:

```bash
cd pkg/mdsplitter
go test -v
```

The implementation maintains the same splitting logic:
- Splits by markdown headers (##, ###, etc.)
- Preserves code blocks (doesn't split on headers inside code blocks)
- Tracks header hierarchy in metadata
- Supports trimming or keeping headers in output

## Cleanup

You can now safely delete the eino-ext markdown dependency:

```bash
# The dependency has already been removed from go.mod
# If you want to clean up the local copy:
rm -rf cloudwego/eino-ext/components/document/transformer/splitter/markdown
```

## Performance

The new implementation has similar performance characteristics:
- O(n) time complexity where n is the number of lines
- Minimal memory overhead
- No external dependencies to load

