package chromem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const defaultBaseURLLlama = "http://localhost:8080"

// llamaEmbeddingRequest represents the request body for llama.cpp embedding API
type llamaEmbeddingRequest struct {
	Content string `json:"content"`
}

// llamaEmbeddingResponse represents the response from llama.cpp embedding API
type llamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// NewEmbeddingFuncLlama returns a function that creates embeddings for a text
// using llama.cpp's embedding API (via llama-server).
//
// The llama-server must be running with an embedding model loaded.
// Example command to start llama-server with an embedding model:
//
//	llama-server -m models/nomic-embed-text-v1.5.Q4_K_M.gguf --port 8080 --embedding
//
// Or use the llama.Service from internal/llama which manages llama-server automatically.
//
// baseURLLlama is the base URL of the llama-server API. If it's empty,
// "http://localhost:8080" is used.
//
// Good embedding models for llama.cpp (as of 2024):
//   - nomic-embed-text (768 dimensions, multilingual)
//   - all-MiniLM-L6-v2 (384 dimensions, English)
//   - bge-small-en-v1.5 (384 dimensions, English)
//   - bge-base-en-v1.5 (768 dimensions, English)
//
// Download from Hugging Face:
//
//	https://huggingface.co/nomic-ai/nomic-embed-text-v1.5-GGUF
//	https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2
func NewEmbeddingFuncLlama(baseURLLlama string) EmbeddingFunc {
	if baseURLLlama == "" {
		baseURLLlama = defaultBaseURLLlama
	}

	// We don't set a default timeout here, although it's usually a good idea.
	// In our case though, the library user can set the timeout on the context,
	// and it might have to be a long timeout, depending on the text length.
	client := &http.Client{}

	var checkedNormalized bool
	checkNormalized := sync.Once{}

	return func(ctx context.Context, text string) ([]float32, error) {
		// Prepare the request body.
		reqBody, err := json.Marshal(llamaEmbeddingRequest{
			Content: text,
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't marshal request body: %w", err)
		}

		// Create the request. Creating it with context is important for a timeout
		// to be possible, because the client is configured without a timeout.
		req, err := http.NewRequestWithContext(ctx, "POST", baseURLLlama+"/embedding", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("couldn't create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send the request.
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("couldn't send request: %w", err)
		}
		defer resp.Body.Close()

		// Check the response status.
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("error response from the embedding API: %s, body: %s", resp.Status, string(body))
		}

		// Read and decode the response body.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("couldn't read response body: %w", err)
		}
		var embeddingResponse llamaEmbeddingResponse
		err = json.Unmarshal(body, &embeddingResponse)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshal response body: %w", err)
		}

		// Check if the response contains embeddings.
		if len(embeddingResponse.Embedding) == 0 {
			return nil, errors.New("no embeddings found in the response")
		}

		v := embeddingResponse.Embedding
		checkNormalized.Do(func() {
			if isNormalized(v) {
				checkedNormalized = true
			} else {
				checkedNormalized = false
			}
		})
		if !checkedNormalized {
			v = normalizeVector(v)
		}

		return v, nil
	}
}
