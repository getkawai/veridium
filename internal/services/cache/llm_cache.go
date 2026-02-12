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

	"log/slog"

	llamaembed "github.com/kawai-network/veridium/pkg/fantasy/providers/llama-embed"
)

// LLMCacheEntry represents a cached LLM response
type LLMCacheEntry struct {
	ID          string
	Query       string
	Response    string
	ContextHash string    // Hash of system prompt + context
	Model       string    // Model used for generation
	Embedding   []float32 // Query embedding for semantic search
	CreatedAt   time.Time
	HitCount    int64
	Metadata    map[string]string
}

// LLMCacheStats holds cache statistics
type LLMCacheStats struct {
	Hits           int64
	Misses         int64
	Size           int
	TotalSavedCost float64       // Estimated cost saved (based on avg cost per query)
	TotalSavedTime time.Duration // Estimated time saved
}

// LLMCache provides semantic caching for LLM responses
type LLMCache struct {
	embedder llamaembed.Embedder
	entries  map[string]*LLMCacheEntry // ID → Entry
	embedIdx map[string]string         // embedding hash → ID (for exact match)
	mu       sync.RWMutex
	config   *LLMCacheConfig
	stats    LLMCacheStats
	statsMu  sync.Mutex

	// For semantic search, we maintain a simple in-memory index
	// In production, this could use DuckDB or a dedicated vector store
	allEmbeddings []cachedEmbedding
	embMu         sync.RWMutex
}

type cachedEmbedding struct {
	ID        string
	Embedding []float32
}

// LLMCacheConfig holds configuration for the LLM cache
type LLMCacheConfig struct {
	MaxSize             int           // Maximum number of entries (default: 1000)
	TTL                 time.Duration // Time to live (default: 1h, 0 = no expiry)
	SimilarityThreshold float32       // Threshold for semantic match (default: 0.95)
	AvgQueryCost        float64       // Estimated cost per query in USD (for stats)
	AvgQueryTimeMs      int64         // Estimated time per query in ms (for stats)
}

// DefaultLLMCacheConfig returns default configuration
func DefaultLLMCacheConfig() *LLMCacheConfig {
	return &LLMCacheConfig{
		MaxSize:             1000,
		TTL:                 1 * time.Hour,
		SimilarityThreshold: 0.95,
		AvgQueryCost:        0.001, // $0.001 per query estimate
		AvgQueryTimeMs:      2000,  // 2 seconds estimate
	}
}

// NewLLMCache creates a new semantic LLM cache
func NewLLMCache(embedder llamaembed.Embedder, config *LLMCacheConfig) *LLMCache {
	if config == nil {
		config = DefaultLLMCacheConfig()
	}

	cache := &LLMCache{
		embedder:      embedder,
		entries:       make(map[string]*LLMCacheEntry),
		embedIdx:      make(map[string]string),
		config:        config,
		allEmbeddings: make([]cachedEmbedding, 0),
	}

	slog.Info("LLMCache initialized",
		"maxSize", config.MaxSize,
		"ttl", config.TTL,
		"similarityThreshold", config.SimilarityThreshold)
	return cache
}

// generateID creates a unique ID for a cache entry
func generateID() string {
	h := sha256.New()
	h.Write([]byte(time.Now().String()))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// hashContext creates a hash of the context (system prompt, etc.)
func hashContext(context string) string {
	h := sha256.New()
	h.Write([]byte(context))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// Get attempts to find a cached response for the query
// Returns (response, found)
func (c *LLMCache) Get(ctx context.Context, query string, contextHash string, model string) (string, bool) {
	if c.embedder == nil {
		return "", false
	}

	// Generate embedding for query
	embeddings, err := c.embedder.Embed(ctx, []string{query})
	if err != nil || len(embeddings) == 0 {
		slog.WarnContext(ctx, "LLMCache: failed to embed query", "error", err)
		c.recordMiss()
		return "", false
	}
	queryEmb := embeddings[0]

	// Search for similar cached queries
	c.embMu.RLock()
	bestMatch := ""
	bestSimilarity := float32(0)

	for _, cached := range c.allEmbeddings {
		similarity := cosineSimilarity(queryEmb, cached.Embedding)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestMatch = cached.ID
		}
	}
	c.embMu.RUnlock()

	// Check if similarity meets threshold
	if bestSimilarity < c.config.SimilarityThreshold {
		c.recordMiss()
		return "", false
	}

	// Get the entry
	c.mu.RLock()
	entry, ok := c.entries[bestMatch]
	c.mu.RUnlock()

	if !ok {
		c.recordMiss()
		return "", false
	}

	// Check context match (must be same context/system prompt)
	if entry.ContextHash != contextHash {
		c.recordMiss()
		return "", false
	}

	// Check model match (optional: could allow cross-model caching)
	if entry.Model != model && model != "" {
		c.recordMiss()
		return "", false
	}

	// Check TTL
	if c.config.TTL > 0 && time.Since(entry.CreatedAt) > c.config.TTL {
		c.recordMiss()
		return "", false
	}

	// Cache hit!
	c.mu.Lock()
	entry.HitCount++
	c.mu.Unlock()

	c.recordHit()
	slog.DebugContext(ctx, "LLMCache hit",
		"query", truncate(query, 50),
		"similarity", bestSimilarity,
		"hitCount", entry.HitCount)

	return entry.Response, true
}

// Set stores a query-response pair in the cache
func (c *LLMCache) Set(ctx context.Context, query string, response string, contextHash string, model string, metadata map[string]string) error {
	if c.embedder == nil {
		return nil
	}

	// Generate embedding for query
	embeddings, err := c.embedder.Embed(ctx, []string{query})
	if err != nil || len(embeddings) == 0 {
		slog.WarnContext(ctx, "LLMCache: failed to embed query for caching", "error", err)
		return err
	}
	queryEmb := embeddings[0]

	id := generateID()
	entry := &LLMCacheEntry{
		ID:          id,
		Query:       query,
		Response:    response,
		ContextHash: contextHash,
		Model:       model,
		Embedding:   queryEmb,
		CreatedAt:   time.Now(),
		HitCount:    0,
		Metadata:    metadata,
	}

	c.mu.Lock()
	c.entries[id] = entry

	// Evict if over max size
	if len(c.entries) > c.config.MaxSize {
		c.evictLRU()
	}
	c.mu.Unlock()

	// Add to embedding index
	c.embMu.Lock()
	c.allEmbeddings = append(c.allEmbeddings, cachedEmbedding{
		ID:        id,
		Embedding: queryEmb,
	})
	c.embMu.Unlock()

	slog.DebugContext(ctx, "LLMCache: stored entry",
		"id", id,
		"query", truncate(query, 50))

	return nil
}

// evictLRU removes least recently used entries (must be called with lock held)
func (c *LLMCache) evictLRU() {
	// Find entry with lowest hit count and oldest
	var oldestID string
	var oldestTime time.Time
	var lowestHits int64 = -1

	for id, entry := range c.entries {
		if lowestHits == -1 || entry.HitCount < lowestHits ||
			(entry.HitCount == lowestHits && entry.CreatedAt.Before(oldestTime)) {
			oldestID = id
			oldestTime = entry.CreatedAt
			lowestHits = entry.HitCount
		}
	}

	if oldestID != "" {
		delete(c.entries, oldestID)

		// Remove from embedding index
		c.embMu.Lock()
		newEmbeddings := make([]cachedEmbedding, 0, len(c.allEmbeddings)-1)
		for _, e := range c.allEmbeddings {
			if e.ID != oldestID {
				newEmbeddings = append(newEmbeddings, e)
			}
		}
		c.allEmbeddings = newEmbeddings
		c.embMu.Unlock()

		slog.Debug("LLMCache: evicted entry", "id", oldestID)
	}
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (sqrt32(normA) * sqrt32(normB))
}

func sqrt32(x float32) float32 {
	if x <= 0 {
		return 0
	}
	// Newton's method for square root
	z := x / 2
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

// truncate truncates a string to max length
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// recordHit increments hit counter
func (c *LLMCache) recordHit() {
	c.statsMu.Lock()
	c.stats.Hits++
	c.stats.TotalSavedCost += c.config.AvgQueryCost
	c.stats.TotalSavedTime += time.Duration(c.config.AvgQueryTimeMs) * time.Millisecond
	c.statsMu.Unlock()
}

// recordMiss increments miss counter
func (c *LLMCache) recordMiss() {
	c.statsMu.Lock()
	c.stats.Misses++
	c.statsMu.Unlock()
}

// GetStats returns cache statistics
func (c *LLMCache) GetStats() LLMCacheStats {
	c.statsMu.Lock()
	defer c.statsMu.Unlock()

	c.mu.RLock()
	c.stats.Size = len(c.entries)
	c.mu.RUnlock()

	return c.stats
}

// HitRate returns the cache hit rate (0.0 - 1.0)
func (c *LLMCache) HitRate() float64 {
	c.statsMu.Lock()
	defer c.statsMu.Unlock()

	total := c.stats.Hits + c.stats.Misses
	if total == 0 {
		return 0
	}
	return float64(c.stats.Hits) / float64(total)
}

// Clear removes all entries from the cache
func (c *LLMCache) Clear() {
	c.mu.Lock()
	c.entries = make(map[string]*LLMCacheEntry)
	c.embedIdx = make(map[string]string)
	c.mu.Unlock()

	c.embMu.Lock()
	c.allEmbeddings = make([]cachedEmbedding, 0)
	c.embMu.Unlock()

	c.statsMu.Lock()
	c.stats = LLMCacheStats{}
	c.statsMu.Unlock()

	slog.Info("LLMCache cleared")
}

// InvalidateByContext removes all entries with a specific context hash
// Useful when system prompt or knowledge base changes
func (c *LLMCache) InvalidateByContext(contextHash string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	toDelete := make([]string, 0)
	for id, entry := range c.entries {
		if entry.ContextHash == contextHash {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		delete(c.entries, id)
	}

	// Update embedding index
	if len(toDelete) > 0 {
		c.embMu.Lock()
		deleteSet := make(map[string]bool)
		for _, id := range toDelete {
			deleteSet[id] = true
		}
		newEmbeddings := make([]cachedEmbedding, 0)
		for _, e := range c.allEmbeddings {
			if !deleteSet[e.ID] {
				newEmbeddings = append(newEmbeddings, e)
			}
		}
		c.allEmbeddings = newEmbeddings
		c.embMu.Unlock()
	}

	if len(toDelete) > 0 {
		slog.Info("LLMCache: invalidated entries by context", "count", len(toDelete))
	}

	return len(toDelete)
}

// InvalidateByModel removes all entries for a specific model
// Useful when model is updated
func (c *LLMCache) InvalidateByModel(model string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	toDelete := make([]string, 0)
	for id, entry := range c.entries {
		if entry.Model == model {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		delete(c.entries, id)
	}

	// Update embedding index
	if len(toDelete) > 0 {
		c.embMu.Lock()
		deleteSet := make(map[string]bool)
		for _, id := range toDelete {
			deleteSet[id] = true
		}
		newEmbeddings := make([]cachedEmbedding, 0)
		for _, e := range c.allEmbeddings {
			if !deleteSet[e.ID] {
				newEmbeddings = append(newEmbeddings, e)
			}
		}
		c.allEmbeddings = newEmbeddings
		c.embMu.Unlock()
	}

	if len(toDelete) > 0 {
		slog.Info("LLMCache: invalidated entries by model", "model", model, "count", len(toDelete))
	}

	return len(toDelete)
}

// HashContext is a helper to create context hash from system prompt
func HashContext(systemPrompt string, additionalContext ...string) string {
	combined := systemPrompt
	for _, ctx := range additionalContext {
		combined += ctx
	}
	return hashContext(combined)
}
