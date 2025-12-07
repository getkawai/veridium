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
	"embed"
	"encoding/json"
	"log"
	"strings"
	"sync"
)

//go:embed configs/*.json
var configsFS embed.FS

// ModelInfo contains metadata about a specific model
type ModelInfo struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	CostPerMillionIn   float64        `json:"cost_per_1m_in"`
	CostPerMillionOut  float64        `json:"cost_per_1m_out"`
	CostPerMillionInCached  float64   `json:"cost_per_1m_in_cached"`
	CostPerMillionOutCached float64   `json:"cost_per_1m_out_cached"`
	ContextWindow      int            `json:"context_window"`
	DefaultMaxTokens   int            `json:"default_max_tokens"`
	CanReason          bool           `json:"can_reason"`
	ReasoningLevels    []string       `json:"reasoning_levels,omitempty"`
	DefaultReasoningEffort string     `json:"default_reasoning_effort,omitempty"`
	SupportsAttachments bool          `json:"supports_attachments"`
	Options            map[string]any `json:"options,omitempty"`
}

// ProviderCatalog contains all models for a provider
type ProviderCatalog struct {
	Name              string      `json:"name"`
	ID                string      `json:"id"`
	APIKey            string      `json:"api_key"`
	APIEndpoint       string      `json:"api_endpoint"`
	Type              string      `json:"type"`
	DefaultLargeModel string      `json:"default_large_model_id"`
	DefaultSmallModel string      `json:"default_small_model_id"`
	Models            []ModelInfo `json:"models"`
	modelIndex        map[string]*ModelInfo
}

// ModelCatalog provides access to all provider model catalogs
type ModelCatalog struct {
	providers map[string]*ProviderCatalog
	mu        sync.RWMutex
}

var (
	globalCatalog *ModelCatalog
	catalogOnce   sync.Once
)

// GetCatalog returns the global model catalog (singleton)
func GetCatalog() *ModelCatalog {
	catalogOnce.Do(func() {
		globalCatalog = &ModelCatalog{
			providers: make(map[string]*ProviderCatalog),
		}
		globalCatalog.loadEmbedded()
	})
	return globalCatalog
}

// loadEmbedded loads all JSON configs from embedded filesystem
func (c *ModelCatalog) loadEmbedded() {
	files := []string{
		"configs/openrouter.json",
		"configs/zai.json",
	}

	for _, file := range files {
		data, err := configsFS.ReadFile(file)
		if err != nil {
			log.Printf("⚠️  ModelCatalog: Failed to read %s: %v", file, err)
			continue
		}

		var catalog ProviderCatalog
		if err := json.Unmarshal(data, &catalog); err != nil {
			log.Printf("⚠️  ModelCatalog: Failed to parse %s: %v", file, err)
			continue
		}

		// Build model index for fast lookup
		// For OpenRouter, only index free models (ID ends with ":free")
		catalog.modelIndex = make(map[string]*ModelInfo)
		freeCount := 0
		for i := range catalog.Models {
			model := &catalog.Models[i]
			if catalog.ID == "openrouter" {
				// Only include free models for OpenRouter
				if strings.HasSuffix(model.ID, ":free") {
					catalog.modelIndex[model.ID] = model
					freeCount++
				}
			} else {
				catalog.modelIndex[model.ID] = model
			}
		}

		c.providers[catalog.ID] = &catalog
		if catalog.ID == "openrouter" {
			log.Printf("📚 ModelCatalog: Loaded %s with %d free models (filtered from %d)", catalog.Name, freeCount, len(catalog.Models))
		} else {
			log.Printf("📚 ModelCatalog: Loaded %s with %d models", catalog.Name, len(catalog.Models))
		}
	}
}

// GetProvider returns a provider catalog by ID
func (c *ModelCatalog) GetProvider(providerID string) *ProviderCatalog {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.providers[providerID]
}

// GetModel returns model info by provider and model ID
func (c *ModelCatalog) GetModel(providerID, modelID string) *ModelInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	provider := c.providers[providerID]
	if provider == nil {
		return nil
	}
	return provider.modelIndex[modelID]
}

// GetModelInfo returns model info, searching across all providers if needed
func (c *ModelCatalog) GetModelInfo(modelID string) *ModelInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try each provider
	for _, provider := range c.providers {
		if info := provider.modelIndex[modelID]; info != nil {
			return info
		}
	}
	return nil
}

// ListProviders returns all loaded provider IDs
func (c *ModelCatalog) ListProviders() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ids := make([]string, 0, len(c.providers))
	for id := range c.providers {
		ids = append(ids, id)
	}
	return ids
}

// GetEndpoint returns the API endpoint for a provider
func (c *ModelCatalog) GetEndpoint(providerID string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if p := c.providers[providerID]; p != nil {
		return p.APIEndpoint
	}
	return ""
}

// GetDefaultModel returns the default model ID for a provider (large or small)
func (c *ModelCatalog) GetDefaultModel(providerID string, small bool) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if p := c.providers[providerID]; p != nil {
		if small {
			return p.DefaultSmallModel
		}
		return p.DefaultLargeModel
	}
	return ""
}
