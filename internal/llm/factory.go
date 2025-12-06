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

package llm

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/internal/llm/openai"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
	"github.com/kawai-network/veridium/types/message"
)

// ProviderFactory creates LLM providers based on configuration
type ProviderFactory struct {
	toolRegistry *tools.ToolRegistry
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory(toolRegistry *tools.ToolRegistry) *ProviderFactory {
	return &ProviderFactory{
		toolRegistry: toolRegistry,
	}
}

// CreateProvider creates an LLM provider from configuration
// For local llama provider, use the existing LlamaProviderAdapter instead
func (f *ProviderFactory) CreateProvider(config types.ProviderConfig) (Provider, error) {
	log.Printf("🏭 Creating provider: type=%s, model=%s", config.Type, config.Model)

	switch config.Type {
	case types.ProviderOpenRouter, types.ProviderZhipuAI:
		// OpenAI-compatible providers use the same implementation
		// Wrap in adapter to implement Provider interface
		return NewOpenAIProviderAdapter(openai.NewProvider(config, f.toolRegistry)), nil

	case types.ProviderLlama:
		// Local llama provider should be created via NewLlamaProviderAdapter
		// in agent_chat_service.go which has access to LibraryService
		return nil, fmt.Errorf("llama provider must be created with LibraryService, use NewLlamaProviderAdapter instead")

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}
}

// ============================================================================
// OpenAI Provider Adapter (implements llm.Provider interface)
// ============================================================================

// OpenAIProviderAdapter wraps openai.Provider to implement llm.Provider interface
type OpenAIProviderAdapter struct {
	provider *openai.Provider
}

// NewOpenAIProviderAdapter creates a new adapter
func NewOpenAIProviderAdapter(provider *openai.Provider) *OpenAIProviderAdapter {
	return &OpenAIProviderAdapter{provider: provider}
}

// Generate implements Provider.Generate
func (a *OpenAIProviderAdapter) Generate(ctx context.Context, messages []message.Message) (*types.LLMResponse, error) {
	return a.provider.Generate(ctx, messages)
}

// RunAgentLoop implements Provider.RunAgentLoop
func (a *OpenAIProviderAdapter) RunAgentLoop(ctx context.Context, messages []message.Message, maxIterations int) (*types.LLMResponse, []message.Message, error) {
	return a.provider.RunAgentLoop(ctx, messages, maxIterations)
}

// RunAgentLoopWithStreaming implements Provider.RunAgentLoopWithStreaming
func (a *OpenAIProviderAdapter) RunAgentLoopWithStreaming(ctx context.Context, messages []message.Message, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*types.LLMResponse, []message.Message, error) {
	return a.provider.RunAgentLoopWithStreaming(ctx, messages, maxIterations, streamCallback, toolCallback)
}

// WithTools implements Provider.WithTools
func (a *OpenAIProviderAdapter) WithTools(toolNames []string) Provider {
	return &OpenAIProviderAdapter{provider: a.provider.WithTools(toolNames)}
}

// WithoutTools implements Provider.WithoutTools
func (a *OpenAIProviderAdapter) WithoutTools() Provider {
	return &OpenAIProviderAdapter{provider: a.provider.WithoutTools()}
}

// GetConfig returns the provider configuration
func (a *OpenAIProviderAdapter) GetConfig() types.ProviderConfig {
	return a.provider.GetConfig()
}

// SetModel changes the model
func (a *OpenAIProviderAdapter) SetModel(model string) {
	a.provider.SetModel(model)
}

// CreateOpenRouterProvider creates an OpenRouter provider with convenient defaults
func (f *ProviderFactory) CreateOpenRouterProvider(apiKey, model string) Provider {
	config := types.ProviderConfig{
		Type:      types.ProviderOpenRouter,
		Name:      "OpenRouter",
		APIKey:    apiKey,
		Model:     model,
		MaxTokens: 4096,
		Options: map[string]any{
			"app_name": "Veridium",
		},
	}
	return NewOpenAIProviderAdapter(openai.NewProvider(config, f.toolRegistry))
}

// CreateZhipuProvider creates a Zhipu GLM provider with convenient defaults
func (f *ProviderFactory) CreateZhipuProvider(apiKey, model string) Provider {
	if model == "" {
		model = "glm-4-flash" // Fast and cost-effective
	}
	config := types.ProviderConfig{
		Type:      types.ProviderZhipuAI,
		Name:      "Zhipu GLM",
		APIKey:    apiKey,
		Model:     model,
		MaxTokens: 4096,
	}
	return NewOpenAIProviderAdapter(openai.NewProvider(config, f.toolRegistry))
}

// ============================================================================
// Provider Registry (for managing multiple configured providers)
// ============================================================================

// ProviderRegistry manages multiple LLM provider configurations
type ProviderRegistry struct {
	factory   *ProviderFactory
	providers map[string]Provider             // name -> provider instance
	configs   map[string]types.ProviderConfig // name -> config
	active    string                          // currently active provider name
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry(toolRegistry *tools.ToolRegistry) *ProviderRegistry {
	return &ProviderRegistry{
		factory:   NewProviderFactory(toolRegistry),
		providers: make(map[string]Provider),
		configs:   make(map[string]types.ProviderConfig),
	}
}

// RegisterConfig registers a provider configuration (does not create provider yet)
func (r *ProviderRegistry) RegisterConfig(name string, config types.ProviderConfig) {
	r.configs[name] = config
	log.Printf("📋 Registered provider config: %s (type: %s, model: %s)", name, config.Type, config.Model)
}

// GetProvider gets or creates a provider by name
func (r *ProviderRegistry) GetProvider(name string) (Provider, error) {
	// Check if already instantiated
	if provider, exists := r.providers[name]; exists {
		return provider, nil
	}

	// Check if config exists
	config, exists := r.configs[name]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", name)
	}

	// Create provider
	provider, err := r.factory.CreateProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider %s: %w", name, err)
	}

	r.providers[name] = provider
	return provider, nil
}

// SetActive sets the active provider
func (r *ProviderRegistry) SetActive(name string) error {
	if _, exists := r.configs[name]; !exists {
		return fmt.Errorf("provider not found: %s", name)
	}
	r.active = name
	log.Printf("✅ Active provider set to: %s", name)
	return nil
}

// GetActive returns the currently active provider
func (r *ProviderRegistry) GetActive() (Provider, error) {
	if r.active == "" {
		return nil, fmt.Errorf("no active provider set")
	}
	return r.GetProvider(r.active)
}

// GetActiveName returns the name of the active provider
func (r *ProviderRegistry) GetActiveName() string {
	return r.active
}

// ListProviders returns all registered provider names
func (r *ProviderRegistry) ListProviders() []string {
	names := make([]string, 0, len(r.configs))
	for name := range r.configs {
		names = append(names, name)
	}
	return names
}

// GetConfig returns the configuration for a provider
func (r *ProviderRegistry) GetConfig(name string) (types.ProviderConfig, bool) {
	config, exists := r.configs[name]
	return config, exists
}

// RemoveProvider removes a provider from the registry
func (r *ProviderRegistry) RemoveProvider(name string) {
	delete(r.providers, name)
	delete(r.configs, name)
	if r.active == name {
		r.active = ""
	}
	log.Printf("🗑️  Removed provider: %s", name)
}
