package chromem

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
)

// Constants for embedding configuration
const (
	defaultContextSize = 2048
	defaultBatchSize   = 512
)

// LlamaEmbedder wraps llama.cpp library for generating embeddings.
// It provides thread-safe, lazy-initialized embedding generation with automatic normalization.
type LlamaEmbedder struct {
	libPath         string
	modelPath       string
	model           llama.Model
	context         llama.Context
	vocab           llama.Vocab
	mutex           sync.Mutex
	initOnce        sync.Once
	initErr         error
	skipLibraryInit bool // Skip Load() and Init() if already done elsewhere
	dimensions      int32
}

// NewLlamaEmbedder creates a new llama embedder that uses library directly.
//
// Parameters:
//   - libPath: directory containing llama.cpp libraries (e.g., ~/.llama-cpp/bin)
//   - modelPath: path to the embedding model GGUF file
//
// Recommended embedding models (as of 2024):
//   - granite-embedding-107m-multilingual (384 dimensions, 18 languages) - RECOMMENDED
//   - nomic-embed-text-v1.5 (768 dimensions, multilingual)
//   - all-MiniLM-L6-v2 (384 dimensions, English)
//   - bge-small-en-v1.5 (384 dimensions, English)
//   - bge-base-en-v1.5 (768 dimensions, English)
//
// Download from Hugging Face:
//   - https://huggingface.co/bartowski/granite-embedding-107m-multilingual-GGUF
//   - https://huggingface.co/nomic-ai/nomic-embed-text-v1.5-GGUF
//   - https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2
//
// Example:
//
//	embedder := NewLlamaEmbedder(
//	    "~/.llama-cpp/bin",
//	    "~/.llama-cpp/models/granite-embedding-107m-multilingual-Q6_K_L.gguf",
//	)
//	defer embedder.Cleanup()
//
//	embedding, err := embedder.GenerateEmbedding(ctx, "Hello world")
func NewLlamaEmbedder(libPath, modelPath string) *LlamaEmbedder {
	return &LlamaEmbedder{
		libPath:   libPath,
		modelPath: modelPath,
	}
}

// NewLlamaEmbedderWithPreloadedLibrary creates an embedder assuming library is already loaded.
// Use this when llama.Load() and llama.Init() have already been called elsewhere
// (e.g., by LibraryService) to avoid double-initialization.
//
// This is useful when sharing the llama.cpp library with other components like chat models.
//
// Example:
//
//	// In LibraryService initialization
//	llama.Load(libPath)
//	llama.Init()
//
//	// Create embedder without re-initializing
//	embedder := NewLlamaEmbedderWithPreloadedLibrary(modelPath)
func NewLlamaEmbedderWithPreloadedLibrary(modelPath string) *LlamaEmbedder {
	return &LlamaEmbedder{
		modelPath:       modelPath,
		skipLibraryInit: true,
	}
}

// Initialize loads the library and model.
// This is called automatically on first use via sync.Once for thread-safe lazy initialization.
// It's safe to call this multiple times - it will only initialize once.
func (e *LlamaEmbedder) Initialize() error {
	e.initOnce.Do(func() {
		// Load llama.cpp library only if not already loaded
		if !e.skipLibraryInit && e.libPath != "" {
			// Expand home directory if needed
			expandedLibPath := e.expandPath(e.libPath)

			if err := llama.Load(expandedLibPath); err != nil {
				// If Load fails, it might be already loaded - log but continue
				// (yzma doesn't track if library is already loaded)
				log.Printf("Warning: llama.Load() failed (may be already loaded): %v", err)
			}

			// Initialize backend only if we loaded the library
			llama.Init()
			log.Printf("✅ Llama.cpp library initialized from: %s", expandedLibPath)
		}

		// Expand home directory in model path
		expandedModelPath := e.expandPath(e.modelPath)

		// Validate model file exists
		if _, err := os.Stat(expandedModelPath); err != nil {
			e.initErr = fmt.Errorf("model file not found: %s (%w)", expandedModelPath, err)
			return
		}

		// Load model (this is safe to do even if library was loaded elsewhere)
		e.model = llama.ModelLoadFromFile(expandedModelPath, llama.ModelDefaultParams())
		if e.model == 0 {
			e.initErr = fmt.Errorf("failed to load embedding model from: %s", expandedModelPath)
			return
		}

		// Get embedding dimensions
		e.dimensions = llama.ModelNEmbd(e.model)
		log.Printf("✅ Embedding model loaded: %s (dimensions: %d)", filepath.Base(expandedModelPath), e.dimensions)

		// Create context for embeddings
		ctxParams := llama.ContextDefaultParams()
		ctxParams.NCtx = defaultContextSize           // Context size
		ctxParams.NBatch = defaultBatchSize           // Batch size
		ctxParams.Embeddings = 1                      // Enable embeddings
		ctxParams.PoolingType = llama.PoolingTypeMean // Mean pooling for better quality

		e.context = llama.InitFromModel(e.model, ctxParams)
		if e.context == 0 {
			llama.ModelFree(e.model)
			e.model = 0
			e.initErr = errors.New("failed to create context for embeddings")
			return
		}

		// Get vocabulary
		e.vocab = llama.ModelGetVocab(e.model)
	})

	return e.initErr
}

// expandPath expands ~ to home directory in paths
func (e *LlamaEmbedder) expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[1:])
		}
	}
	return path
}

// GenerateEmbedding generates normalized embeddings for the given text.
// The embeddings are L2-normalized, making them suitable for cosine similarity.
//
// Thread-safe: Multiple goroutines can call this method concurrently.
//
// Parameters:
//   - ctx: context for cancellation (currently not used, but reserved for future)
//   - text: input text to embed
//
// Returns:
//   - []float32: normalized embedding vector
//   - error: any error during generation
//
// Example:
//
//	embedding, err := embedder.GenerateEmbedding(ctx, "Machine learning is fascinating")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Embedding dimensions: %d\n", len(embedding))
func (e *LlamaEmbedder) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Ensure initialized
	if err := e.Initialize(); err != nil {
		return nil, fmt.Errorf("initialization failed: %w", err)
	}

	// Validate input
	if text == "" {
		return nil, errors.New("empty text provided")
	}

	// Tokenize text
	tokens := llama.Tokenize(e.vocab, text, true, true)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("failed to tokenize text: %q", text)
	}

	// Create batch and decode
	batch := llama.BatchGetOne(tokens)
	if err := llama.Decode(e.context, batch); err != 0 {
		return nil, fmt.Errorf("failed to decode tokens: error code %d", err)
	}

	// Get embeddings from sequence
	vec := llama.GetEmbeddingsSeq(e.context, 0, e.dimensions)
	if vec == nil {
		return nil, errors.New("failed to get embeddings from context")
	}

	// Copy to result slice
	result := make([]float32, e.dimensions)
	for i := int32(0); i < e.dimensions; i++ {
		result[i] = vec[i]
	}

	// Normalize embeddings using L2 norm
	result = normalizeVector(result)

	return result, nil
}

// BatchGenerateEmbeddings generates embeddings for multiple texts efficiently.
// This is more efficient than calling GenerateEmbedding multiple times.
//
// Thread-safe: Uses internal locking for concurrent access.
//
// Parameters:
//   - ctx: context for cancellation
//   - texts: slice of input texts to embed
//
// Returns:
//   - [][]float32: slice of normalized embedding vectors
//   - error: any error during generation
//
// Example:
//
//	texts := []string{"Hello", "World", "Machine Learning"}
//	embeddings, err := embedder.BatchGenerateEmbeddings(ctx, texts)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Generated %d embeddings\n", len(embeddings))
func (e *LlamaEmbedder) BatchGenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, errors.New("no texts provided")
	}

	results := make([][]float32, len(texts))

	for i, text := range texts {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		embedding, err := e.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}

		results[i] = embedding
	}

	return results, nil
}

// GetDimensions returns the embedding dimensions of the loaded model.
// Returns 0 if model is not yet initialized.
func (e *LlamaEmbedder) GetDimensions() int32 {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.dimensions
}

// IsInitialized returns whether the embedder has been initialized.
func (e *LlamaEmbedder) IsInitialized() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.model != 0 && e.context != 0
}

// Cleanup releases all resources including model, context, and optionally the backend.
//
// Note: This does NOT call llama.BackendFree() if skipLibraryInit is true,
// because the library might be shared with other components.
//
// It's safe to call this multiple times - subsequent calls will be no-ops.
//
// Example:
//
//	embedder := NewLlamaEmbedder(libPath, modelPath)
//	defer embedder.Cleanup() // Ensure cleanup on exit
func (e *LlamaEmbedder) Cleanup() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.context != 0 {
		llama.Free(e.context)
		e.context = 0
		log.Println("✅ Embedding context freed")
	}

	if e.model != 0 {
		llama.ModelFree(e.model)
		e.model = 0
		log.Println("✅ Embedding model freed")
	}

	// Only free backend if we initialized it
	if !e.skipLibraryInit {
		llama.BackendFree()
		log.Println("✅ Llama.cpp backend freed")
	}
}

// NewEmbeddingFuncLlama returns an EmbeddingFunc that creates embeddings
// using llama.cpp library directly (no HTTP server required).
//
// This function is designed to be used with chromem's embedding system.
// It automatically handles normalization and provides a simple interface.
//
// Parameters:
//   - libPath: directory containing llama.cpp libraries
//   - modelPath: path to the embedding model GGUF file
//
// Returns:
//   - EmbeddingFunc: function compatible with chromem's embedding interface
//
// Example:
//
//	embedFunc := NewEmbeddingFuncLlama(
//	    "~/.llama-cpp/bin",
//	    "~/.llama-cpp/models/granite-embedding-107m-multilingual-Q6_K_L.gguf",
//	)
//
//	// Use with chromem
//	collection, err := db.GetOrCreateCollection("docs", nil, embedFunc)
func NewEmbeddingFuncLlama(libPath, modelPath string) EmbeddingFunc {
	embedder := NewLlamaEmbedder(libPath, modelPath)

	var checkedNormalized bool
	checkNormalized := sync.Once{}

	return func(ctx context.Context, text string) ([]float32, error) {
		v, err := embedder.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("llama embedding generation failed: %w", err)
		}

		// Check if embeddings are normalized (only once for efficiency)
		checkNormalized.Do(func() {
			checkedNormalized = isNormalized(v)
			if checkedNormalized {
				log.Println("✅ Embeddings are already normalized")
			} else {
				log.Println("⚠️  Embeddings not normalized, applying normalization")
			}
		})

		// Normalize if needed (though GenerateEmbedding already normalizes)
		// This is a safety check in case the model doesn't normalize
		if !checkedNormalized {
			v = normalizeVector(v)
		}

		return v, nil
	}
}

// NewEmbeddingFuncLlamaWithPreloadedLibrary creates an EmbeddingFunc using a preloaded library.
// Use this when the llama.cpp library is already loaded by another component.
//
// Parameters:
//   - modelPath: path to the embedding model GGUF file
//
// Returns:
//   - EmbeddingFunc: function compatible with chromem's embedding interface
//
// Example:
//
//	// In your main initialization
//	llama.Load(libPath)
//	llama.Init()
//
//	// Create embedding function without re-initializing library
//	embedFunc := NewEmbeddingFuncLlamaWithPreloadedLibrary(modelPath)
func NewEmbeddingFuncLlamaWithPreloadedLibrary(modelPath string) EmbeddingFunc {
	embedder := NewLlamaEmbedderWithPreloadedLibrary(modelPath)

	var checkedNormalized bool
	checkNormalized := sync.Once{}

	return func(ctx context.Context, text string) ([]float32, error) {
		v, err := embedder.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("llama embedding generation failed: %w", err)
		}

		// Check if embeddings are normalized (only once)
		checkNormalized.Do(func() {
			checkedNormalized = isNormalized(v)
		})

		// Normalize if needed
		if !checkedNormalized {
			v = normalizeVector(v)
		}

		return v, nil
	}
}

// Note: normalizeVector is already defined in vector.go
// We use the existing implementation from the chromem package
