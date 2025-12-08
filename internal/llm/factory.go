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

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/fantasy/llamalib/tools"
	"github.com/kawai-network/veridium/fantasy/providers/openai"
	"github.com/kawai-network/veridium/fantasy/providers/openaicompat"
	"github.com/kawai-network/veridium/fantasy/providers/openrouter"
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

// CreateLanguageModel creates a fantasy.LanguageModel from configuration
func (f *ProviderFactory) CreateLanguageModel(config ProviderConfig) (fantasy.LanguageModel, error) {
	log.Printf("🏭 Creating language model: type=%s, model=%s", config.Type, config.Model)

	// Set default model if not specified
	if config.Model == "" {
		if defaultModel, ok := DefaultModels[config.Type]; ok {
			config.Model = defaultModel
		}
	}

	// Enrich config from catalog if available
	catalog := GetCatalog()
	if modelInfo := catalog.GetModelInfo(config.Model); modelInfo != nil {
		if config.MaxTokens == 0 {
			config.MaxTokens = modelInfo.DefaultMaxTokens
		}
		log.Printf("📚 ModelCatalog: %s (context=%d, max_tokens=%d, can_reason=%v)",
			modelInfo.Name, modelInfo.ContextWindow, modelInfo.DefaultMaxTokens, modelInfo.CanReason)
	}

	var fantasyProvider fantasy.Provider
	var err error

	switch config.Type {
	case ProviderOpenRouter:
		opts := []openrouter.Option{
			openrouter.WithAPIKey(config.APIKey),
		}
		if config.Name != "" {
			opts = append(opts, openrouter.WithName(config.Name))
		}
		fantasyProvider, err = openrouter.New(opts...)

	case ProviderZhipuAI:
		baseURL := config.BaseURL
		if baseURL == "" {
			baseURL = ProviderEndpoints[ProviderZhipuAI]
		}
		opts := []openaicompat.Option{
			openaicompat.WithBaseURL(baseURL),
			openaicompat.WithAPIKey(config.APIKey),
			openaicompat.WithName("zhipu"),
		}
		fantasyProvider, err = openaicompat.New(opts...)

	case ProviderOpenAI:
		opts := []openai.Option{
			openai.WithAPIKey(config.APIKey),
		}
		if config.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(config.BaseURL))
		}
		fantasyProvider, err = openai.New(opts...)

	case ProviderLlama:
		return nil, fmt.Errorf("llama model must be created via fantasy/providers/llama provider with LibraryService")

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create fantasy provider: %w", err)
	}

	// Create language model directly from provider
	ctx := context.Background()
	model, err := fantasyProvider.LanguageModel(ctx, config.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to create language model: %w", err)
	}

	return model, nil
}

// CreateOpenRouterModel creates an OpenRouter model with convenient defaults
func (f *ProviderFactory) CreateOpenRouterModel(apiKey, model string) (fantasy.LanguageModel, error) {
	config := ProviderConfig{
		Type:   ProviderOpenRouter,
		Name:   "OpenRouter",
		APIKey: apiKey,
		Model:  model,
	}
	return f.CreateLanguageModel(config)
}

// CreateZhipuModel creates a Zhipu GLM model with convenient defaults
func (f *ProviderFactory) CreateZhipuModel(apiKey, model string) (fantasy.LanguageModel, error) {
	if model == "" {
		model = "glm-4-flash"
	}
	config := ProviderConfig{
		Type:   ProviderZhipuAI,
		Name:   "Zhipu GLM",
		APIKey: apiKey,
		Model:  model,
	}
	return f.CreateLanguageModel(config)
}

// ============================================================================
// Model Registry (for managing multiple configured models)
// ============================================================================

// ModelRegistry manages multiple LLM model configurations
type ModelRegistry struct {
	factory *ProviderFactory
	models  map[string]fantasy.LanguageModel
	configs map[string]ProviderConfig
	active  string
}

// NewModelRegistry creates a new model registry
func NewModelRegistry(toolRegistry *tools.ToolRegistry) *ModelRegistry {
	return &ModelRegistry{
		factory: NewProviderFactory(toolRegistry),
		models:  make(map[string]fantasy.LanguageModel),
		configs: make(map[string]ProviderConfig),
	}
}

// RegisterConfig registers a model configuration (does not create model yet)
func (r *ModelRegistry) RegisterConfig(name string, config ProviderConfig) {
	r.configs[name] = config
	log.Printf("📋 Registered model config: %s (type: %s, model: %s)", name, config.Type, config.Model)
}

// GetModel gets or creates a model by name
func (r *ModelRegistry) GetModel(name string) (fantasy.LanguageModel, error) {
	if model, exists := r.models[name]; exists {
		return model, nil
	}

	config, exists := r.configs[name]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", name)
	}

	model, err := r.factory.CreateLanguageModel(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create model %s: %w", name, err)
	}

	r.models[name] = model
	return model, nil
}

// SetActive sets the active model
func (r *ModelRegistry) SetActive(name string) error {
	if _, exists := r.configs[name]; !exists {
		return fmt.Errorf("model not found: %s", name)
	}
	r.active = name
	log.Printf("✅ Active model set to: %s", name)
	return nil
}

// GetActive returns the currently active model
func (r *ModelRegistry) GetActive() (fantasy.LanguageModel, error) {
	if r.active == "" {
		return nil, fmt.Errorf("no active model set")
	}
	return r.GetModel(r.active)
}

// GetActiveName returns the name of the active model
func (r *ModelRegistry) GetActiveName() string {
	return r.active
}

// ListModels returns all registered model names
func (r *ModelRegistry) ListModels() []string {
	names := make([]string, 0, len(r.configs))
	for name := range r.configs {
		names = append(names, name)
	}
	return names
}

// GetConfig returns the configuration for a model
func (r *ModelRegistry) GetConfig(name string) (ProviderConfig, bool) {
	config, exists := r.configs[name]
	return config, exists
}

// RemoveModel removes a model from the registry
func (r *ModelRegistry) RemoveModel(name string) {
	delete(r.models, name)
	delete(r.configs, name)
	if r.active == name {
		r.active = ""
	}
	log.Printf("🗑️  Removed model: %s", name)
}
