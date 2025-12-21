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

package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/kawai-network/veridium/pkg/xlog"
	llamaembed "github.com/kawai-network/veridium/pkg/fantasy/providers/llama-embed"
)

// EmbeddingCacheEntry represents a cached embedding with metadata
type EmbeddingCacheEntry struct {
	Embedding []float32
	CreatedAt time.Time
	HitCount  int64
}

// EmbeddingCacheStats holds cache statistics
type EmbeddingCacheStats struct {
	Hits       int64
	Misses     int64
	Size       int
	TotalSaved time.Duration
}

// EmbeddingCache wraps an embedder with caching capability
type EmbeddingCache struct {
	embedder   llamaembed.Embedder
	cache      map[string]*EmbeddingCacheEntry
	mu         sync.RWMutex
	maxSize    int
	ttl        time.Duration
	stats      EmbeddingCacheStats
	statsMu    sync.Mutex
	avgEmbedMs int64 // average embedding time in ms for stats
}

// EmbeddingCacheConfig holds configuration for the cache
type EmbeddingCacheConfig struct {
	MaxSize    int           // Maximum number of entries (default: 10000)
	TTL        time.Duration // Time to live for entries (default: 24h, 0 = no expiry)
	AvgEmbedMs int64         // Estimated average embedding time in ms (for stats)
}

// DefaultEmbeddingCacheConfig returns default cache configuration
func DefaultEmbeddingCacheConfig() *EmbeddingCacheConfig {
	return &EmbeddingCacheConfig{
		MaxSize:    10000,
		TTL:        24 * time.Hour,
		AvgEmbedMs: 100,
	}
}

// NewEmbeddingCache creates a new embedding cache wrapper
func NewEmbeddingCache(embedder llamaembed.Embedder, config *EmbeddingCacheConfig) *EmbeddingCache {
	if config == nil {
		config = DefaultEmbeddingCacheConfig()
	}

	cache := &EmbeddingCache{
		embedder:   embedder,
		cache:      make(map[string]*EmbeddingCacheEntry),
		maxSize:    config.MaxSize,
		ttl:        config.TTL,
		avgEmbedMs: config.AvgEmbedMs,
	}

	xlog.Info("EmbeddingCache initialized", "maxSize", config.MaxSize, "ttl", config.TTL)
	return cache
}

// hashText generates a SHA256 hash of the input text
func hashText(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// Embed generates embeddings with caching
func (c *EmbeddingCache) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	results := make([][]float32, len(texts))
	toEmbed := make([]string, 0)
	toEmbedIdx := make([]int, 0)
	hashes := make([]string, len(texts))

	// Check cache for each text
	c.mu.RLock()
	for i, text := range texts {
		hash := hashText(text)
		hashes[i] = hash

		if entry, ok := c.cache[hash]; ok {
			// Check TTL
			if c.ttl == 0 || time.Since(entry.CreatedAt) < c.ttl {
				results[i] = entry.Embedding
				entry.HitCount++
				c.recordHit()
				continue
			}
		}
		// Cache miss - need to embed
		toEmbed = append(toEmbed, text)
		toEmbedIdx = append(toEmbedIdx, i)
		c.recordMiss()
	}
	c.mu.RUnlock()

	// Generate embeddings for cache misses
	if len(toEmbed) > 0 {
		embeddings, err := c.embedder.Embed(ctx, toEmbed)
		if err != nil {
			return nil, err
		}

		// Store results and update cache
		c.mu.Lock()
		for j, emb := range embeddings {
			idx := toEmbedIdx[j]
			results[idx] = emb

			// Store in cache
			hash := hashes[idx]
			c.cache[hash] = &EmbeddingCacheEntry{
				Embedding: emb,
				CreatedAt: time.Now(),
				HitCount:  0,
			}
		}

		// Evict if over max size (simple LRU-like: remove oldest)
		if len(c.cache) > c.maxSize {
			c.evictOldest(len(c.cache) - c.maxSize)
		}
		c.mu.Unlock()

		xlog.Debug("EmbeddingCache: generated new embeddings", "count", len(toEmbed))
	}

	return results, nil
}

// evictOldest removes the oldest entries (must be called with lock held)
func (c *EmbeddingCache) evictOldest(count int) {
	if count <= 0 {
		return
	}

	// Find oldest entries
	type entry struct {
		hash      string
		createdAt time.Time
	}
	entries := make([]entry, 0, len(c.cache))
	for hash, e := range c.cache {
		entries = append(entries, entry{hash: hash, createdAt: e.CreatedAt})
	}

	// Sort by creation time (oldest first)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].createdAt.Before(entries[i].createdAt) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Remove oldest
	for i := 0; i < count && i < len(entries); i++ {
		delete(c.cache, entries[i].hash)
	}

	xlog.Debug("EmbeddingCache: evicted entries", "count", count)
}

// recordHit increments hit counter
func (c *EmbeddingCache) recordHit() {
	c.statsMu.Lock()
	c.stats.Hits++
	c.stats.TotalSaved += time.Duration(c.avgEmbedMs) * time.Millisecond
	c.statsMu.Unlock()
}

// recordMiss increments miss counter
func (c *EmbeddingCache) recordMiss() {
	c.statsMu.Lock()
	c.stats.Misses++
	c.statsMu.Unlock()
}

// Dimensions returns the embedding dimension size
func (c *EmbeddingCache) Dimensions() int {
	return c.embedder.Dimensions()
}

// Close releases resources
func (c *EmbeddingCache) Close() error {
	c.mu.Lock()
	c.cache = nil
	c.mu.Unlock()
	return c.embedder.Close()
}

// GetStats returns cache statistics
func (c *EmbeddingCache) GetStats() EmbeddingCacheStats {
	c.statsMu.Lock()
	defer c.statsMu.Unlock()

	c.mu.RLock()
	c.stats.Size = len(c.cache)
	c.mu.RUnlock()

	return c.stats
}

// HitRate returns the cache hit rate (0.0 - 1.0)
func (c *EmbeddingCache) HitRate() float64 {
	c.statsMu.Lock()
	defer c.statsMu.Unlock()

	total := c.stats.Hits + c.stats.Misses
	if total == 0 {
		return 0
	}
	return float64(c.stats.Hits) / float64(total)
}

// Clear removes all entries from the cache
func (c *EmbeddingCache) Clear() {
	c.mu.Lock()
	c.cache = make(map[string]*EmbeddingCacheEntry)
	c.mu.Unlock()

	c.statsMu.Lock()
	c.stats = EmbeddingCacheStats{}
	c.statsMu.Unlock()

	xlog.Info("EmbeddingCache cleared")
}

// Warmup preloads embeddings into cache
func (c *EmbeddingCache) Warmup(ctx context.Context, texts []string) error {
	if len(texts) == 0 {
		return nil
	}

	xlog.Info("EmbeddingCache warmup started", "count", len(texts))
	_, err := c.Embed(ctx, texts)
	if err != nil {
		return err
	}
	xlog.Info("EmbeddingCache warmup completed")
	return nil
}

// Ensure EmbeddingCache implements llamaembed.Embedder
var _ llamaembed.Embedder = (*EmbeddingCache)(nil)
