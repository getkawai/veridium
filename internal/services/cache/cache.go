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

// Package cache provides caching layers for embeddings and LLM responses
// to improve performance and reduce costs.
//
// Usage:
//
//	// Wrap existing embedder with cache
//	cachedEmbedder := cache.NewEmbeddingCache(embedder, nil)
//
//	// Use cachedEmbedder instead of embedder
//	vectorSearchService, _ := services.NewVectorSearchService(db, duckDB, cachedEmbedder)
//
//	// For LLM cache
//	llmCache := cache.NewLLMCache(embedder, nil)
//	response, found := llmCache.Get(ctx, query, contextHash, model)
//	if !found {
//	    response = callLLM(query)
//	    llmCache.Set(ctx, query, response, contextHash, model, nil)
//	}
package cache

import (
	"log/slog"

	llamaembed "github.com/getkawai/unillm/providers/llama-embed"
)

// CacheManager manages all cache instances
type CacheManager struct {
	embeddingCache *EmbeddingCache
	llmCache       *LLMCache
}

// CacheConfig holds configuration for all caches
type CacheConfig struct {
	EmbeddingConfig *EmbeddingCacheConfig
	LLMConfig       *LLMCacheConfig
	Enabled         bool
}

// DefaultCacheConfig returns default configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		EmbeddingConfig: DefaultEmbeddingCacheConfig(),
		LLMConfig:       DefaultLLMCacheConfig(),
		Enabled:         true,
	}
}

// NewCacheManager creates a new cache manager
func NewCacheManager(embedder llamaembed.Embedder, config *CacheConfig) *CacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}

	if !config.Enabled {
		slog.Info("Cache disabled")
		return &CacheManager{}
	}

	manager := &CacheManager{}

	if embedder != nil {
		manager.embeddingCache = NewEmbeddingCache(embedder, config.EmbeddingConfig)
		manager.llmCache = NewLLMCache(embedder, config.LLMConfig)
	}

	return manager
}

// GetEmbeddingCache returns the embedding cache (implements llamaembed.Embedder)
func (m *CacheManager) GetEmbeddingCache() *EmbeddingCache {
	return m.embeddingCache
}

// GetLLMCache returns the LLM cache
func (m *CacheManager) GetLLMCache() *LLMCache {
	return m.llmCache
}

// GetCachedEmbedder returns the cached embedder if available, otherwise returns the original
func (m *CacheManager) GetCachedEmbedder(original llamaembed.Embedder) llamaembed.Embedder {
	if m.embeddingCache != nil {
		return m.embeddingCache
	}
	return original
}

// PrintStats logs cache statistics
func (m *CacheManager) PrintStats() {
	if m.embeddingCache != nil {
		stats := m.embeddingCache.GetStats()
		slog.Info("EmbeddingCache stats",
			"hits", stats.Hits,
			"misses", stats.Misses,
			"size", stats.Size,
			"hitRate", m.embeddingCache.HitRate(),
			"timeSaved", stats.TotalSaved)
	}

	if m.llmCache != nil {
		stats := m.llmCache.GetStats()
		slog.Info("LLMCache stats",
			"hits", stats.Hits,
			"misses", stats.Misses,
			"size", stats.Size,
			"hitRate", m.llmCache.HitRate(),
			"costSaved", stats.TotalSavedCost,
			"timeSaved", stats.TotalSavedTime)
	}
}

// Clear clears all caches
func (m *CacheManager) Clear() {
	if m.embeddingCache != nil {
		m.embeddingCache.Clear()
	}
	if m.llmCache != nil {
		m.llmCache.Clear()
	}
}
