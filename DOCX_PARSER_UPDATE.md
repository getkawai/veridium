# DOCX Parser Update - Markdown Conversion & Table Extraction

## Summary

Updated DOCX parser to support markdown conversion and table extraction, matching the capabilities of Eino-Ext's DOCX parser while using gooxml's native features.

## Changes Made

### 1. Added Configuration Options

```go
type DocxParserConfig struct {
    ToMarkdown    bool  // Convert to markdown format
    IncludeTables bool  // Extract table content
}
```

### 2. Markdown Conversion Support

The parser now supports gooxml's built-in `ToMarkdown()` method:

```go
if p.ToMarkdown {
    markdown := doc.ToMarkdown()
    content = markdown
}
```

**Benefits:**
- ✅ Preserves document structure (headings, lists, tables)
- ✅ Better for semantic search
- ✅ Better for chunking (structure-aware)
- ✅ Native gooxml implementation (no external dependencies)

### 3. Table Extraction

When `ToMarkdown` is disabled, tables can still be extracted as plain text:

```go
if p.IncludeTables && doc.X().Body != nil {
    // Navigate through document structure
    for _, ble := range doc.X().Body.EG_BlockLevelElts {
        for _, c := range ble.EG_ContentBlockContent {
            for _, tbl := range c.Tbl {
                // Extract table content
                table := document.NewTable(doc, tbl)
                // Extract text from each cell
                for _, rowContent := range table.X().EG_ContentRowContent {
                    for _, tr := range rowContent.Tr {
                        for _, tc := range tr.EG_ContentCellContent {
                            for _, cell := range tc.Tc {
                                // Extract cell paragraphs
                                for _, cellBle := range cell.EG_BlockLevelElts {
                                    for _, cellContent := range cellBle.EG_ContentBlockContent {
                                        for _, p := range cellContent.P {
                                            // Extract runs and text
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
```

### 4. Default Configuration in FileManager

FileManager now uses markdown mode by default:

```go
docxParser, _ := parsers.NewDocxParser(ctx, &parsers.DocxParserConfig{
    ToMarkdown:    true,  // Enable markdown conversion
    IncludeTables: true,  // Include table content
})
```

## Comparison with Eino-Ext

| Feature | Eino-Ext DOCX Parser | Our DOCX Parser |
|---------|---------------------|-----------------|
| **Library** | docx2md (external) | gooxml (built-in) |
| **Markdown Conversion** | ✅ Yes | ✅ Yes |
| **Table Support** | ✅ Yes | ✅ Yes |
| **Headers/Footers** | ✅ Configurable | ✅ Always included |
| **Section Splitting** | ✅ Yes | ⚠️ Not yet |
| **Comments** | ❌ Not supported | ❌ Not supported |
| **Images** | ⚠️ Via docx2md | ⚠️ Via ToMarkdownWithImages() |
| **Dependencies** | External (docx2md) | Built-in (gooxml) |

## Usage Examples

### Example 1: Markdown Mode (Recommended)

```go
parser, _ := parsers.NewDocxParser(ctx, &parsers.DocxParserConfig{
    ToMarkdown:    true,
    IncludeTables: true,
})

docs, _ := parser.Parse(ctx, reader)
fmt.Println(docs[0].Content)
// Output: Markdown formatted text with tables
```

### Example 2: Plain Text Mode

```go
parser, _ := parsers.NewDocxParser(ctx, &parsers.DocxParserConfig{
    ToMarkdown:    false,
    IncludeTables: true,
})

docs, _ := parser.Parse(ctx, reader)
fmt.Println(docs[0].Content)
// Output: Plain text with tab-separated table cells
```

### Example 3: Custom Parser in FileManager

```go
// Override default parser
customDocxParser, _ := parsers.NewDocxParser(ctx, &parsers.DocxParserConfig{
    ToMarkdown:    false,  // Use plain text mode
    IncludeTables: false,  // Skip tables
})

fileManager, _ := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
    Indexer: indexer,
    CustomParsers: map[string]parser.Parser{
        ".docx": customDocxParser,
    },
})
```

## Gooxml ToMarkdown Features

The gooxml library provides three markdown conversion methods:

1. **`ToMarkdown()`** - Basic markdown conversion (used by default)
2. **`ToMarkdownWithImages(imageDir)`** - Saves images to local directory
3. **`ToMarkdownWithImageURLs(baseURL)`** - References images via URLs

**Current Implementation:** Uses `ToMarkdown()` for simplicity and no external dependencies.

**Future Enhancement:** Could support image extraction by using `ToMarkdownWithImages()` or `ToMarkdownWithImageURLs()`.

## Benefits of This Approach

### 1. No External Dependencies
- Uses gooxml which is already in the project
- No need for docx2md or other external libraries
- Reduces dependency management complexity

### 2. Better Structure Preservation
- Markdown format preserves document structure
- Headings, lists, and tables are properly formatted
- Better for semantic search and RAG applications

### 3. Flexible Configuration
- Can choose between markdown and plain text
- Can enable/disable table extraction
- Easy to extend with more options

### 4. Performance
- Native gooxml implementation
- No external process calls
- Efficient memory usage

## Testing

### Manual Testing Required
1. Test with DOCX files containing tables
2. Test with DOCX files with complex formatting
3. Test with DOCX files with headers/footers
4. Compare markdown output quality
5. Verify table extraction accuracy

### Example Test Files
- Simple document with paragraphs
- Document with tables
- Document with headers and footers
- Document with complex formatting (bold, italic, lists)
- Document with images (for future enhancement)

## Future Enhancements

### 1. Section Splitting
Similar to Eino-Ext, could split document into sections:
- Main content
- Headers
- Footers
- Tables

### 2. Image Support
Use `ToMarkdownWithImages()` to extract and save images:
```go
if p.IncludeImages {
    markdown, err := doc.ToMarkdownWithImages(imageDir)
    // Handle images
}
```

### 3. Advanced Table Formatting
Could enhance table extraction with:
- Column headers detection
- Cell alignment
- Merged cells handling

### 4. Metadata Extraction
Extract document properties:
- Author
- Title
- Creation date
- Last modified date

## Migration Guide

### From Old Parser (No Config)

**Before:**
```go
parser, _ := parsers.NewDocxParser(ctx)
```

**After:**
```go
parser, _ := parsers.NewDocxParser(ctx, &parsers.DocxParserConfig{
    ToMarkdown:    true,  // New default
    IncludeTables: true,  // New default
})
```

### FileManager (Automatic)

FileManager automatically uses the new configuration:
```go
fileManager, _ := chromemAdapter.NewFileManager(ctx, config)
// Already uses markdown mode with tables
```

## Conclusion

The updated DOCX parser now provides:
- ✅ Markdown conversion for better structure preservation
- ✅ Table extraction in both markdown and plain text modes
- ✅ No external dependencies (uses built-in gooxml)
- ✅ Flexible configuration options
- ✅ Compatible with Eino-Ext's approach but using native features

This brings our DOCX parser to feature parity with Eino-Ext while maintaining the benefits of using gooxml, which is already integrated into the project.

