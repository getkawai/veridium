# Chromem Eino Adapter

Eino-compatible adapters for [chromem](../../chromem), allowing chromem to be used seamlessly in Eino workflows and graphs.

## Features

- ✅ **Local Embeddings**: No external API calls required
- ✅ **Zero Dependencies**: Pure Go vector database
- ✅ **In-Memory + Persistence**: Fast in-memory with optional disk persistence
- ✅ **Eino Integration**: Native support for Eino graphs and workflows
- ✅ **Type-Safe**: Full Go type safety with generics
- ✅ **Backward Compatible**: Works with existing chromem collections
- ✅ **File Management**: Automatic parsing and indexing of DOCX, XLSX, PDF, HTML, TXT, MD
- ✅ **Custom Parsers**: Extensible parser system using gooxml and other libraries
- ✅ **Auto-Chunking**: Intelligent text splitting with Eino recursive splitter

## Installation

```bash
go get github.com/kawai-network/veridium/pkg/eino-adapters/chromem
```

## Quick Start

### Basic Usage

```go
import (
    "context"
    "github.com/cloudwego/eino/schema"
    chromemAdapter "github.com/kawai-network/veridium/pkg/eino-adapters/chromem"
    "github.com/kawai-network/veridium/pkg/chromem"
)

// 1. Create chromem database and collection
db := chromem.NewDB()
collection, _ := db.CreateCollection(
    "my_docs",
    nil,
    chromem.NewEmbeddingFuncDefault(),
)

// 2. Wrap with Eino adapters
indexer := chromemAdapter.NewIndexer(collection)
retriever, _ := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
    Collection: collection,
    TopK:       5,
})

// 3. Use Eino interfaces
ctx := context.Background()

// Add documents
docs := []*schema.Document{
    {Content: "Hello world", MetaData: map[string]any{"source": "test"}},
}
ids, _ := indexer.Store(ctx, docs)

// Query documents
results, _ := retriever.Retrieve(ctx, "hello")
```

### Use in Eino Graph

```go
import (
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/components/retriever"
)

// Create graph
type Input struct {
    Query string
}
type Output struct {
    Documents []*schema.Document
}

graph := compose.NewGraph[Input, Output]()

// Add retriever node
graph.AddRetrieverNode("retriever", retriever,
    compose.WithRetrieverNodeInputKey("Query"),
)

// Add processing node
graph.AddLambdaNode("process",
    compose.InvokableLambda(func(ctx context.Context, docs []*schema.Document) (Output, error) {
        return Output{Documents: docs}, nil
    }),
).AddInput("retriever")

// Compile and run
compiled, _ := graph.Compile(ctx)
result, _ := compiled.Invoke(ctx, Input{Query: "search query"})
```

## File Management

### FileManager

The `FileManager` provides automatic document parsing, chunking, and indexing for various file formats.

```go
type FileManager struct {
    // contains filtered or unexported fields
}

type FileManagerConfig struct {
    Indexer       *Indexer              // Required: Eino indexer
    AssetDir      string                // Optional: Directory for file copies (default: "./assets")
    ChunkSize     int                   // Optional: Max chunk size (default: 1000)
    OverlapSize   int                   // Optional: Chunk overlap (default: 200)
    CustomParsers map[string]parser.Parser // Optional: Custom parsers by extension
}

func NewFileManager(ctx context.Context, config *FileManagerConfig) (*FileManager, error)

// Store a file with automatic parsing and chunking
func (fm *FileManager) StoreFile(ctx context.Context, filePath string, metadata map[string]any) error

// Remove a file and its chunks
func (fm *FileManager) RemoveFile(ctx context.Context, filename string) error

// List all tracked files
func (fm *FileManager) ListFiles() []string

// Check if a file is tracked
func (fm *FileManager) FileExists(filename string) bool

// Get supported file extensions
func (fm *FileManager) GetSupportedExtensions() []string
```

**Supported File Formats:**

| Format | Extension | Parser | Library | Features |
|--------|-----------|--------|---------|----------|
| Microsoft Word | `.docx` | DocxParser | gooxml/document | Markdown conversion, tables, headers, footers |
| Microsoft Excel | `.xlsx` | XlsxParser | gooxml/spreadsheet | Multiple sheets, cell values |
| PDF | `.pdf` | PdfParser | dslipak/pdf | Plain text extraction |
| HTML | `.html`, `.htm` | HtmlParser | golang.org/x/net/html | Structure preservation |
| Text | `.txt`, `.md` | TextParser | built-in | Direct text |

**Example Usage:**

```go
// 1. Setup chromem and Eino adapter
db, _ := chromem.NewPersistentDB("./vectors", true)
collection, _ := db.CreateCollection("docs", nil, chromem.NewEmbeddingFuncDefault())
indexer := chromemAdapter.NewIndexer(collection)

// 2. Create file manager
fileManager, _ := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
    Indexer:     indexer,
    AssetDir:    "./assets",
    ChunkSize:   1500,
    OverlapSize: 300,
})

// 3. Store files (auto-parses and chunks)
fileManager.StoreFile(ctx, "/path/to/manual.pdf", map[string]any{
    "category": "documentation",
    "version":  "2.0",
})

fileManager.StoreFile(ctx, "/path/to/report.docx", map[string]any{
    "category": "report",
    "year":     2024,
})

fileManager.StoreFile(ctx, "/path/to/data.xlsx", map[string]any{
    "category": "data",
})

// 4. List stored files
files := fileManager.ListFiles()
fmt.Printf("Stored %d files\n", len(files))

// 5. Search across all files
retriever, _ := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
    Collection: collection,
    TopK:       5,
})

docs, _ := retriever.Retrieve(ctx, "installation instructions")
for _, doc := range docs {
    fmt.Printf("Found in: %s\n", doc.MetaData["source_file"])
    fmt.Printf("Content: %s\n", doc.Content[:200])
}
```

### Custom Parsers

You can add custom parsers for additional file formats:

```go
// Implement parser.Parser interface
type MyCustomParser struct{}

func (p *MyCustomParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
    // Your parsing logic
    return docs, nil
}

func (p *MyCustomParser) GetType() string {
    return "MyCustomParser"
}

// Add to file manager
fileManager, _ := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
    Indexer: indexer,
    CustomParsers: map[string]parser.Parser{
        ".custom": &MyCustomParser{},
    },
})

// Now you can store .custom files
fileManager.StoreFile(ctx, "/path/to/file.custom", metadata)
```

### Parser Configuration

Each parser can be configured:

```go
// DOCX Parser - automatically converts to markdown
docxParser, _ := parsers.NewDocxParser(ctx)

// XLSX Parser - each sheet as separate document
xlsxParser, _ := parsers.NewXlsxParser(ctx, &parsers.XlsxParserConfig{
    SheetsAsDocuments: true,
})

// HTML Parser - preserve document structure
htmlParser, _ := parsers.NewHtmlParser(ctx, &parsers.HtmlParserConfig{
    PreserveStructure: true,
})

// Use in file manager
fileManager, _ := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
    Indexer: indexer,
    CustomParsers: map[string]parser.Parser{
        ".docx": docxParser,
        ".xlsx": xlsxParser,
        ".html": htmlParser,
    },
})
```

**DOCX Parser Features:**

The DOCX parser extracts content in sections, similar to Eino-Ext:

- ✅ **Section-Based Extraction**: Separates MAIN CONTENT, HEADERS, and FOOTERS
- ✅ **Markdown Conversion**: Main content converted to markdown with tables
- ✅ **Structure Preservation**: Maintains headings, lists, and tables in body
- ✅ **Complete Content**: Includes all document sections
- ✅ **No Configuration Needed**: Works out of the box

```go
// Simple usage - no configuration needed
docxParser, _ := parsers.NewDocxParser(ctx)
docs, _ := docxParser.Parse(ctx, reader)

// Output format (similar to Eino-Ext):
// === HEADERS ===
// [Header text from all headers]
//
// === MAIN CONTENT ===
// [Markdown formatted body with tables]
//
// === FOOTERS ===
// [Footer text from all footers]
```

## API Reference

### Indexer

The `Indexer` wraps `chromem.Collection` and implements `indexer.Indexer` from Eino.

```go
type Indexer struct {
    // contains filtered or unexported fields
}

func NewIndexer(collection *chromem.Collection) *Indexer

// Store implements indexer.Indexer
func (i *Indexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error)

// GetCollection returns the underlying chromem collection
func (i *Indexer) GetCollection() *chromem.Collection
```

**Features:**
- Converts Eino documents to chromem documents automatically
- Generates IDs if not provided
- Handles metadata conversion (any → string)
- Uses all CPU cores for parallel embedding generation

### Retriever

The `Retriever` wraps `chromem.Collection.Query()` and implements `retriever.Retriever` from Eino.

```go
type Retriever struct {
    // contains filtered or unexported fields
}

type RetrieverConfig struct {
    Collection *chromem.Collection // Required
    TopK       int                 // Optional, default: 5
}

func NewRetriever(config *RetrieverConfig) (*Retriever, error)

// Retrieve implements retriever.Retriever
func (r *Retriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error)

// GetCollection returns the underlying chromem collection
func (r *Retriever) GetCollection() *chromem.Collection
```

**Supported Options:**
- `retriever.WithTopK(n)` - Number of results to return
- `retriever.WithScoreThreshold(t)` - Minimum similarity threshold
- `retriever.WithFilters(map[string]string)` - Metadata filters

## Advanced Usage

### Persistent Storage

```go
// Create persistent database
db, _ := chromem.NewPersistentDB("./data/chromem", true)

collection, _ := db.GetOrCreateCollection(
    "persistent_docs",
    nil,
    chromem.NewEmbeddingFuncDefault(),
)

// Use with Eino adapters
indexer := chromemAdapter.NewIndexer(collection)
retriever, _ := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
    Collection: collection,
})

// Documents are automatically persisted to disk
```

### Custom Embedding Functions

```go
// Use Ollama for local embeddings
embedFunc := chromem.NewEmbeddingFuncOllama(
    "http://localhost:11434/api/embeddings",
    "nomic-embed-text",
)

collection, _ := db.CreateCollection("docs", nil, embedFunc)
indexer := chromemAdapter.NewIndexer(collection)
```

### Metadata Filtering

```go
// Add documents with metadata
docs := []*schema.Document{
    {
        Content: "Document about Go",
        MetaData: map[string]any{
            "language": "go",
            "category": "programming",
        },
    },
}
indexer.Store(ctx, docs)

// Query with filters
results, _ := retriever.Retrieve(ctx, "programming",
    retriever.WithFilters(map[string]string{
        "language": "go",
    }),
)
```

### Access Chromem-Specific Features

```go
// Get underlying collection for chromem-specific operations
collection := retriever.GetCollection()

// Use chromem's negative query
results, _ := collection.Query(ctx, chromem.QueryOptions{
    QueryText: "machine learning",
    NResults:  10,
    Negative: chromem.NegativeQueryOptions{
        Text: "deep learning",
        Mode: chromem.NEGATIVE_MODE_FILTER,
    },
})
```

## Integration with Existing Code

This adapter allows you to use your existing chromem collections in Eino workflows without migration:

```go
// Your existing chromem setup
db := getExistingChromemDB()
collection := getExistingCollection()

// Wrap for Eino
einoIndexer := chromemAdapter.NewIndexer(collection)
einoRetriever, _ := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
    Collection: collection,
})

// Now use in Eino graph
graph.AddRetrieverNode("my_retriever", einoRetriever)
```

## Comparison with LangChain Go

| Feature | LangChain Go VectorStore | Chromem Eino Adapter |
|---------|-------------------------|----------------------|
| Interface | Single `VectorStore` | Split: `Indexer` + `Retriever` |
| Graph Support | ❌ No | ✅ Native Eino graph |
| Local Embeddings | ⚠️ Limited | ✅ Full support |
| Persistence | ⚠️ Varies | ✅ Built-in |
| Modularity | ⚠️ Coupled | ✅ Separated concerns |
| Composition | ⚠️ Manual | ✅ Graph-based |

## Examples

See example files for complete runnable examples:

### Basic Examples ([example_test.go](./example_test.go))
- Basic indexing and retrieval
- Using in Eino graphs
- Retriever options
- Persistent storage
- Custom embedding functions

### File Management Examples ([file_manager_example_test.go](./file_manager_example_test.go))
- Automatic file parsing and indexing
- Multi-format document processing
- Custom parser integration
- Knowledge base building
- Persistent file management

## Performance

Chromem is designed for local, in-memory operations:

- **Fast**: In-memory vector search with cosine similarity
- **Efficient**: Parallel embedding generation using all CPU cores
- **Compact**: Optional gzip compression for persistence
- **Scalable**: Suitable for collections up to millions of documents

## Limitations

- **Metadata**: Only string values supported (converted automatically)
- **Embeddings**: Must use same embedding function for all documents in a collection
- **Filters**: Only exact match on metadata (no range queries)
- **File Deletion**: Chromem doesn't support document deletion yet (tracked in file index only)

## Contributing

Contributions are welcome! This adapter is part of the Veridium project.

## License

Part of the Veridium project.

## See Also

- [Chromem Documentation](../../chromem/README.md)
- [Eino Documentation](https://github.com/cloudwego/eino)
- [Eino Components](https://github.com/cloudwego/eino/tree/main/components)
- [Gooxml Documentation](../../gooxml/README.md)
- [Parsers Package](./parsers/doc.go)

