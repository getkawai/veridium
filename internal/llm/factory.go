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
	"fmt"
	"log"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/fantasy/providers/openai"
	"github.com/kawai-network/veridium/fantasy/providers/openaicompat"
	"github.com/kawai-network/veridium/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
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
func (f *ProviderFactory) CreateProvider(config types.ProviderConfig) (Provider, error) {
	log.Printf("🏭 Creating provider: type=%s, model=%s", config.Type, config.Model)

	// Set default model if not specified
	if config.Model == "" {
		if defaultModel, ok := types.DefaultModels[config.Type]; ok {
			config.Model = defaultModel
		}
	}

	var fantasyProvider fantasy.Provider
	var err error

	switch config.Type {
	case types.ProviderOpenRouter:
		opts := []openrouter.Option{
			openrouter.WithAPIKey(config.APIKey),
		}
		if config.Name != "" {
			opts = append(opts, openrouter.WithName(config.Name))
		}
		fantasyProvider, err = openrouter.New(opts...)

	case types.ProviderZhipuAI:
		baseURL := config.BaseURL
		if baseURL == "" {
			baseURL = types.ProviderEndpoints[types.ProviderZhipuAI]
		}
		opts := []openaicompat.Option{
			openaicompat.WithBaseURL(baseURL),
			openaicompat.WithAPIKey(config.APIKey),
			openaicompat.WithName("zhipu"),
		}
		fantasyProvider, err = openaicompat.New(opts...)

	case types.ProviderOpenAI:
		opts := []openai.Option{
			openai.WithAPIKey(config.APIKey),
		}
		if config.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(config.BaseURL))
		}
		fantasyProvider, err = openai.New(opts...)

	case types.ProviderLlama:
		return nil, fmt.Errorf("llama provider must be created with LibraryService, use NewLlamaProviderAdapter instead")

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create fantasy provider: %w", err)
	}

	adapter, err := NewFantasyProviderAdapter(fantasyProvider, config.Model, f.toolRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	return adapter, nil
}

// CreateOpenRouterProvider creates an OpenRouter provider with convenient defaults
func (f *ProviderFactory) CreateOpenRouterProvider(apiKey, model string) (Provider, error) {
	config := types.ProviderConfig{
		Type:   types.ProviderOpenRouter,
		Name:   "OpenRouter",
		APIKey: apiKey,
		Model:  model,
	}
	return f.CreateProvider(config)
}

// CreateZhipuProvider creates a Zhipu GLM provider with convenient defaults
func (f *ProviderFactory) CreateZhipuProvider(apiKey, model string) (Provider, error) {
	if model == "" {
		model = "glm-4-flash"
	}
	config := types.ProviderConfig{
		Type:   types.ProviderZhipuAI,
		Name:   "Zhipu GLM",
		APIKey: apiKey,
		Model:  model,
	}
	return f.CreateProvider(config)
}

// ============================================================================
// Provider Registry (for managing multiple configured providers)
// ============================================================================

// ProviderRegistry manages multiple LLM provider configurations
type ProviderRegistry struct {
	factory   *ProviderFactory
	providers map[string]Provider
	configs   map[string]types.ProviderConfig
	active    string
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
	if provider, exists := r.providers[name]; exists {
		return provider, nil
	}

	config, exists := r.configs[name]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", name)
	}

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
