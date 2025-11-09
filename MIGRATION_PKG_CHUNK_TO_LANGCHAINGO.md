# Migration: pkg/chunk → langchaingo/textsplitter

## Summary

Successfully migrated from custom `pkg/chunk` package to using `langchaingo/textsplitter` directly, eliminating code duplication and leveraging the full feature set of LangChain Go.

## Why Migrate?

### 1. **Code Duplication**
`pkg/chunk/chunking.go` and `langchaingo/textsplitter/paragraph.go` contained **identical code**:

```go
// Both had the exact same function:
func SplitParagraphIntoChunks(paragraph string, maxChunkSize int) []string
```

### 2. **Limited Features**
- **pkg/chunk**: Only 1 simple function (47 lines)
- **langchaingo/textsplitter**: Complete suite of text splitters

### 3. **Single Usage Point**
Only `langchaingo/vectorstores/chromem/persistency.go` was using `pkg/chunk`, making migration straightforward.

## What Changed

### Files Modified

#### 1. `langchaingo/vectorstores/chromem/persistency.go`

**Before:**
```go
import (
    "github.com/kawai-network/veridium/pkg/chunk"
    // ...
)

// Usage:
chunk.SplitParagraphIntoChunks(text, maxchunksize)
```

**After:**
```go
import (
    "github.com/kawai-network/veridium/langchaingo/textsplitter"
    // ...
)

// Usage:
textsplitter.SplitParagraphIntoChunks(text, maxchunksize)
```

### Files Deleted

- ❌ `pkg/chunk/chunking.go` (47 lines)
- ❌ Entire `pkg/chunk/` folder

## Benefits

### 1. **No Code Duplication**
Single source of truth for text splitting functionality.

### 2. **Access to Full Feature Set**
Now you can use all LangChain Go text splitters:

```go
// Recursive Character Splitter (recommended for RAG)
splitter := textsplitter.NewRecursiveCharacter(
    textsplitter.WithChunkSize(500),
    textsplitter.WithChunkOverlap(50),
)
chunks, _ := splitter.SplitText(text)

// Markdown Splitter
mdSplitter := textsplitter.NewMarkdownTextSplitter(
    textsplitter.WithChunkSize(1000),
)

// Token Splitter
tokenSplitter := textsplitter.NewTokenSplitter(
    textsplitter.WithModelName("gpt-4"),
    textsplitter.WithChunkSize(500),
)

// Paragraph Splitter (same as old pkg/chunk)
paragraphSplitter := textsplitter.NewParagraph(
    textsplitter.WithChunkSize(500),
)
```

### 3. **Better Maintenance**
- Updates to LangChain Go automatically benefit your code
- No need to maintain custom chunking logic
- Community-tested and optimized

### 4. **Consistent API**
All splitters implement the same interface:

```go
type TextSplitter interface {
    SplitText(text string) ([]string, error)
}
```

## Available Text Splitters in LangChain Go

### 1. **RecursiveCharacter** (Recommended for RAG)
Hierarchical splitting: paragraphs → lines → sentences → words

```go
splitter := textsplitter.NewRecursiveCharacter(
    textsplitter.WithChunkSize(500),
    textsplitter.WithChunkOverlap(50),
)
```

**Best for**: General documents, articles, books

### 2. **Markdown**
Markdown-aware splitting that preserves structure

```go
splitter := textsplitter.NewMarkdownTextSplitter(
    textsplitter.WithChunkSize(1000),
    textsplitter.WithHeadingHierarchy(true),
)
```

**Best for**: Markdown documentation, README files

### 3. **Token**
Token-based splitting using tiktoken

```go
splitter := textsplitter.NewTokenSplitter(
    textsplitter.WithModelName("gpt-4"),
    textsplitter.WithChunkSize(500),
)
```

**Best for**: LLM context window management

### 4. **Paragraph**
Simple word-based splitting (equivalent to old pkg/chunk)

```go
splitter := textsplitter.NewParagraph(
    textsplitter.WithChunkSize(500),
)
```

**Best for**: Simple text, backward compatibility

## Migration Guide for Other Code

If you have other code using `pkg/chunk`, migrate as follows:

### Before:
```go
import "github.com/kawai-network/veridium/pkg/chunk"

chunks := chunk.SplitParagraphIntoChunks(text, 500)
```

### After (Option 1 - Direct function):
```go
import "github.com/kawai-network/veridium/langchaingo/textsplitter"

chunks := textsplitter.SplitParagraphIntoChunks(text, 500)
```

### After (Option 2 - Using splitter interface):
```go
import "github.com/kawai-network/veridium/langchaingo/textsplitter"

splitter := textsplitter.NewParagraph(
    textsplitter.WithChunkSize(500),
)
chunks, err := splitter.SplitText(text)
```

### After (Option 3 - Upgrade to RecursiveCharacter):
```go
import "github.com/kawai-network/veridium/langchaingo/textsplitter"

splitter := textsplitter.NewRecursiveCharacter(
    textsplitter.WithChunkSize(500),
    textsplitter.WithChunkOverlap(50), // Add overlap for better RAG
)
chunks, err := splitter.SplitText(text)
```

## Verification

### Build Check
```bash
✅ go build ./langchaingo/vectorstores/chromem/...
```

### No Remaining References
```bash
✅ grep -r "pkg/chunk" . 
# No matches found
```

### Linter Check
```bash
✅ No linter errors
```

## Recommendations

### For RAG Applications
Use `RecursiveCharacter` with overlap:

```go
splitter := textsplitter.NewRecursiveCharacter(
    textsplitter.WithChunkSize(500),      // Optimal for embeddings
    textsplitter.WithChunkOverlap(50),    // 10% overlap
)
```

### For Markdown Documentation
Use `MarkdownTextSplitter`:

```go
splitter := textsplitter.NewMarkdownTextSplitter(
    textsplitter.WithChunkSize(1000),
    textsplitter.WithHeadingHierarchy(true), // Preserve header context
)
```

### For LLM Context Management
Use `TokenSplitter`:

```go
splitter := textsplitter.NewTokenSplitter(
    textsplitter.WithModelName("gpt-4"),
    textsplitter.WithChunkSize(4000),     // Leave room for prompt
)
```

## Related Files

- **Modified**: `langchaingo/vectorstores/chromem/persistency.go`
- **Deleted**: `pkg/chunk/chunking.go`
- **Available**: `langchaingo/textsplitter/*.go`

## Documentation

For more information on LangChain Go text splitters:
- Source: `langchaingo/textsplitter/`
- Examples: `langchaingo/examples/`
- Tests: `langchaingo/textsplitter/*_test.go`

## Conclusion

✅ **Migration Complete**
- Eliminated code duplication
- Gained access to full LangChain Go feature set
- Simplified maintenance
- No breaking changes (API compatible)
- All builds passing
- No linter errors

The codebase now uses the standard LangChain Go text splitting library, providing better features and easier maintenance.

