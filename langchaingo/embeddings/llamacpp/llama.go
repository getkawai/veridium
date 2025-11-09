package llamacpp

import (
	"fmt"
	"sync"
)

// Vectorizer represents a loaded LLM/Embedding model.
type Vectorizer struct {
	client *LlamaClient
	mutex  sync.RWMutex
}

// NewVectorizer creates a new vectorizer model from the given model file.
func NewVectorizer(modelPath string, gpuLayers int) (*Vectorizer, error) {
	client, err := NewLlamaClient(modelPath, gpuLayers)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize llama client: %w", err)
	}

	return &Vectorizer{
		client: client,
	}, nil
}

// Close closes the model and releases any resources associated with it.
func (m *Vectorizer) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client != nil {
		err := m.client.Close()
		m.client = nil
		return err
	}
	return nil
}

// EmbedText embeds the given text using the model.
func (m *Vectorizer) EmbedText(text string) ([]float32, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.client == nil {
		return nil, fmt.Errorf("vectorizer is closed")
	}

	return m.client.EmbedText(text)
}
