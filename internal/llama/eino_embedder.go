/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package llama

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/kawai-network/veridium/pkg/yzma/llama"
)

const (
	defaultContextSize = 2048
	typ                = "Llama"
)

// EmbeddingConfig holds configuration for Llama embedding
type EmbeddingConfig struct {
	// ModelPath specifies the path to the embedding model GGUF file
	// Required
	ModelPath string `json:"model_path"`

	// SkipLibraryInit indicates whether to skip llama.Load() and llama.Init()
	// Set to true when library is already loaded elsewhere (e.g., by LibraryService)
	// Optional. Default: false
	SkipLibraryInit bool `json:"skip_library_init"`

	// LibPath specifies the directory containing llama.cpp libraries
	// Only used if SkipLibraryInit is false
	// Optional. Default: ""
	LibPath string `json:"lib_path"`

	// ContextSize specifies the maximum context size for tokenization
	// Optional. Default: 2048
	ContextSize int `json:"context_size"`
}

var _ embedding.Embedder = (*Embedder)(nil)

// Embedder implements Eino embedding interface using llama.cpp
type Embedder struct {
	conf       *EmbeddingConfig
	model      llama.Model
	context    llama.Context
	vocab      llama.Vocab
	dimensions int32
	mutex      sync.Mutex
	initOnce   sync.Once
	initErr    error
}

// NewEmbedder creates a new Llama embedder
func NewEmbedder(ctx context.Context, config *EmbeddingConfig) (*Embedder, error) {
	if config == nil {
		return nil, fmt.Errorf("embedding config must not be nil")
	}

	if config.ModelPath == "" {
		return nil, fmt.Errorf("model_path is required")
	}

	if config.ContextSize <= 0 {
		config.ContextSize = defaultContextSize
	}

	embedder := &Embedder{
		conf: config,
	}

	// Initialize eagerly to catch errors early
	if err := embedder.initialize(); err != nil {
		return nil, err
	}

	return embedder, nil
}

// initialize loads the library and model
func (e *Embedder) initialize() error {
	e.initOnce.Do(func() {
		// Load llama.cpp library only if not already loaded
		if !e.conf.SkipLibraryInit && e.conf.LibPath != "" {
			expandedLibPath := expandPath(e.conf.LibPath)

			if err := llama.Load(expandedLibPath); err != nil {
				// If Load fails, it might be already loaded - log but continue
				log.Printf("Warning: llama.Load() failed (may be already loaded): %v", err)
			}

			llama.Init()
			log.Printf("✅ Llama.cpp library initialized from: %s", expandedLibPath)
		}

		// Expand home directory in model path
		expandedModelPath := expandPath(e.conf.ModelPath)

		// Validate model file exists
		if _, err := os.Stat(expandedModelPath); err != nil {
			e.initErr = fmt.Errorf("model file not found: %s (%w)", expandedModelPath, err)
			return
		}

		// Load model
		var err error
		e.model, err = llama.ModelLoadFromFile(expandedModelPath, llama.ModelDefaultParams())
		if err != nil || e.model == 0 {
			e.initErr = fmt.Errorf("failed to load embedding model from: %s (%w)", expandedModelPath, err)
			return
		}

		// Get embedding dimensions
		e.dimensions = llama.ModelNEmbd(e.model)
		log.Printf("✅ Embedding model loaded: %s (dimensions: %d)", filepath.Base(expandedModelPath), e.dimensions)

		// Create context for embeddings
		ctxParams := llama.ContextDefaultParams()
		ctxParams.NCtx = uint32(e.conf.ContextSize)
		ctxParams.NBatch = uint32(e.conf.ContextSize)
		ctxParams.NUbatch = uint32(e.conf.ContextSize)
		ctxParams.Embeddings = 1

		e.context, err = llama.InitFromModel(e.model, ctxParams)
		if err != nil || e.context == 0 {
			llama.ModelFree(e.model)
			e.model = 0
			e.initErr = fmt.Errorf("failed to create context for embeddings: %w", err)
			return
		}

		// Get vocabulary
		e.vocab = llama.ModelGetVocab(e.model)
	})

	return e.initErr
}

// EmbedStrings generates embeddings for multiple texts
func (e *Embedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) (
	embeddings [][]float32, err error) {
	defer func() {
		if err != nil {
			callbacks.OnError(ctx, err)
		}
	}()

	if len(texts) == 0 {
		return nil, errors.New("no texts provided")
	}

	// Ensure initialized
	if err := e.initialize(); err != nil {
		return nil, fmt.Errorf("initialization failed: %w", err)
	}

	options := embedding.GetCommonOptions(&embedding.Options{
		Model: &e.conf.ModelPath,
	}, opts...)

	conf := &embedding.Config{
		Model: *options.Model,
	}

	ctx = callbacks.EnsureRunInfo(ctx, e.GetType(), components.ComponentOfEmbedding)
	ctx = callbacks.OnStart(ctx, &embedding.CallbackInput{
		Texts:  texts,
		Config: conf,
	})

	// Generate embeddings for each text
	result := make([][]float32, len(texts))
	for i, text := range texts {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		emb, err := e.generateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}

		// Use float32 directly (no conversion needed)
		result[i] = emb
	}

	callbacks.OnEnd(ctx, &embedding.CallbackOutput{
		Embeddings: result,
		Config:     conf,
	})

	return result, nil
}

// generateEmbedding generates a single embedding (internal method)
func (e *Embedder) generateEmbedding(ctx context.Context, text string) ([]float32, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if text == "" {
		return nil, errors.New("empty text provided")
	}

	// Tokenize text
	tokens := llama.Tokenize(e.vocab, text, true, true)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("failed to tokenize text: %q", text)
	}

	// Truncate tokens if they exceed context size
	maxTokens := e.conf.ContextSize
	if len(tokens) > maxTokens {
		log.Printf("⚠️  Text has %d tokens, truncating to %d tokens", len(tokens), maxTokens)
		tokens = tokens[:maxTokens]
	}

	// Get model's native context size
	nCtxTrain := int(llama.ModelNCtxTrain(e.model))
	if nCtxTrain <= 0 {
		nCtxTrain = defaultContextSize
	}

	// Strict truncation to context size
	if len(tokens) > nCtxTrain {
		log.Printf("⚠️  Input text too long (%d tokens), truncating to model context size (%d)", len(tokens), nCtxTrain)
		tokens = tokens[:nCtxTrain]
	}

	// Create batch
	batch := llama.BatchGetOne(tokens)

	// Decode tokens
	errCode, err := llama.Decode(e.context, batch)
	if err != nil || errCode != 0 {
		return nil, fmt.Errorf("failed to decode tokens: error code %d (%w)", errCode, err)
	}

	// Get embeddings
	vec, err := llama.GetEmbeddingsSeq(e.context, 0, int32(e.dimensions))
	if err != nil {
		return nil, fmt.Errorf("failed to get embeddings from context: %w", err)
	}

	// Normalize embeddings using L2 norm
	return normalizeVector(vec), nil
}

// GetType returns the embedder type
func (e *Embedder) GetType() string {
	return typ
}

// IsCallbacksEnabled returns whether callbacks are enabled
func (e *Embedder) IsCallbacksEnabled() bool {
	return true
}

// Cleanup releases all resources
func (e *Embedder) Cleanup() {
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
	if !e.conf.SkipLibraryInit {
		llama.BackendFree()
		log.Println("✅ Llama.cpp backend freed")
	}
}

// Helper functions

// expandPath expands ~ to home directory in paths
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[1:])
		}
	}
	return path
}

// normalizeVector normalizes a vector using L2 norm
func normalizeVector(v []float32) []float32 {
	var norm float32
	for _, val := range v {
		norm += val * val
	}
	if norm == 0 {
		return v
	}

	norm = float32(1.0) / float32(sqrt(float64(norm)))
	result := make([]float32, len(v))
	for i, val := range v {
		result[i] = val * norm
	}
	return result
}

// sqrt is a helper for square root
func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
