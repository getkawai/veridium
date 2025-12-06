// Package llama provides an implementation of the fantasy AI SDK for local llama.cpp models.
package llama

import (
	"context"
	"fmt"

	"github.com/kawai-network/veridium/fantasy"
	internalllama "github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

const (
	// Name is the name of the llama provider.
	Name = "llama"
)

type provider struct {
	options options
}

type options struct {
	name         string
	libService   *internalllama.LibraryService
	toolRegistry *tools.ToolRegistry
	modelPath    string
}

// Option defines a function that configures llama provider options.
type Option = func(*options)

// New creates a new llama provider with the given options.
func New(opts ...Option) (fantasy.Provider, error) {
	providerOptions := options{
		name: Name,
	}
	for _, o := range opts {
		o(&providerOptions)
	}

	// Create LibraryService if not provided
	if providerOptions.libService == nil {
		libService, err := internalllama.NewLibraryService()
		if err != nil {
			return nil, fmt.Errorf("failed to create library service: %w", err)
		}
		providerOptions.libService = libService
	}

	return &provider{options: providerOptions}, nil
}

// WithName sets the name for the llama provider.
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithLibraryService sets a pre-configured LibraryService.
func WithLibraryService(libService *internalllama.LibraryService) Option {
	return func(o *options) {
		o.libService = libService
	}
}

// WithToolRegistry sets the tool registry for tool calling support.
func WithToolRegistry(registry *tools.ToolRegistry) Option {
	return func(o *options) {
		o.toolRegistry = registry
	}
}

// WithModelPath sets a specific model path to load.
func WithModelPath(modelPath string) Option {
	return func(o *options) {
		o.modelPath = modelPath
	}
}

// LanguageModel implements fantasy.Provider.
func (p *provider) LanguageModel(ctx context.Context, modelID string) (fantasy.LanguageModel, error) {
	// Wait for library initialization
	if err := p.options.libService.WaitForInitialization(ctx); err != nil {
		return nil, fmt.Errorf("library initialization failed: %w", err)
	}

	// Load model if not loaded or if a specific model is requested
	modelPath := p.options.modelPath
	if modelID != "" {
		modelPath = modelID
	}

	if !p.options.libService.IsChatModelLoaded() || modelPath != "" {
		if err := p.options.libService.LoadChatModel(modelPath); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	return newLanguageModel(
		p.options.libService.GetLoadedChatModel(),
		p.options.name,
		p.options.libService,
		p.options.toolRegistry,
	), nil
}

// Name implements fantasy.Provider.
func (p *provider) Name() string {
	return p.options.name
}

// GetLibraryService returns the underlying LibraryService for advanced usage.
func (p *provider) GetLibraryService() *internalllama.LibraryService {
	return p.options.libService
}

// Cleanup releases all resources held by the provider.
func (p *provider) Cleanup() {
	if p.options.libService != nil {
		p.options.libService.Cleanup()
	}
}
