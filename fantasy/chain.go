package fantasy

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kawai-network/veridium/pkg/xlog"
)

// ChainLanguageModel wraps multiple LanguageModels and tries them in order.
// If the primary model fails, it automatically falls back to the next model in the chain.
// This enables resilient LLM calls with graceful degradation.
type ChainLanguageModel struct {
	models []LanguageModel
	name   string

	// Circuit breaker state per model
	mu               sync.RWMutex
	failures         map[int]int       // failure count per model index
	lastFailure      map[int]time.Time // last failure time per model index
	circuitOpen      map[int]bool      // whether circuit is open (skip model)
	failureThreshold int               // failures before opening circuit
	resetTimeout     time.Duration     // time before trying again
}

// ChainOption configures a ChainLanguageModel.
type ChainOption func(*ChainLanguageModel)

// WithCircuitBreaker enables circuit breaker pattern.
// After failureThreshold consecutive failures, the model will be skipped for resetTimeout duration.
func WithCircuitBreaker(failureThreshold int, resetTimeout time.Duration) ChainOption {
	return func(c *ChainLanguageModel) {
		c.failureThreshold = failureThreshold
		c.resetTimeout = resetTimeout
	}
}

// WithChainName sets a descriptive name for the chain (for logging).
func WithChainName(name string) ChainOption {
	return func(c *ChainLanguageModel) {
		c.name = name
	}
}

// NewChain creates a new ChainLanguageModel that tries models in order until one succeeds.
// At least one model must be provided.
func NewChain(models []LanguageModel, opts ...ChainOption) (*ChainLanguageModel, error) {
	if len(models) == 0 {
		return nil, fmt.Errorf("at least one model is required for chain")
	}

	// Filter out nil models
	validModels := make([]LanguageModel, 0, len(models))
	for _, m := range models {
		if m != nil {
			validModels = append(validModels, m)
		}
	}

	if len(validModels) == 0 {
		return nil, fmt.Errorf("at least one non-nil model is required for chain")
	}

	// Build name from model names
	names := make([]string, len(validModels))
	for i, m := range validModels {
		names[i] = fmt.Sprintf("%s/%s", m.Provider(), m.Model())
	}

	chain := &ChainLanguageModel{
		models:           validModels,
		name:             strings.Join(names, " -> "),
		failures:         make(map[int]int),
		lastFailure:      make(map[int]time.Time),
		circuitOpen:      make(map[int]bool),
		failureThreshold: 0, // disabled by default
		resetTimeout:     5 * time.Minute,
	}

	for _, opt := range opts {
		opt(chain)
	}

	return chain, nil
}

// isCircuitOpen checks if a model's circuit breaker is open (should skip).
func (c *ChainLanguageModel) isCircuitOpen(idx int) bool {
	if c.failureThreshold == 0 {
		return false // circuit breaker disabled
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.circuitOpen[idx] {
		return false
	}

	// If resetTimeout is 0, circuit stays open until app restart (for daily rate limits)
	if c.resetTimeout == 0 {
		return true
	}

	// Check if reset timeout has passed
	if time.Since(c.lastFailure[idx]) > c.resetTimeout {
		return false // allow retry (will be reset on success or re-opened on failure)
	}

	return true
}

// recordFailure records a failure for a model.
func (c *ChainLanguageModel) recordFailure(idx int) {
	if c.failureThreshold == 0 {
		return // circuit breaker disabled
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.failures[idx]++
	c.lastFailure[idx] = time.Now()

	if c.failures[idx] >= c.failureThreshold {
		c.circuitOpen[idx] = true
		xlog.Warn(fmt.Sprintf("⚡ Chain[%s]: Circuit opened for model %d (%s/%s) after %d failures",
			c.name, idx, c.models[idx].Provider(), c.models[idx].Model(), c.failures[idx]))
	}
}

// recordSuccess records a success for a model (resets circuit breaker).
func (c *ChainLanguageModel) recordSuccess(idx int) {
	if c.failureThreshold == 0 {
		return // circuit breaker disabled
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.circuitOpen[idx] {
		xlog.Info(fmt.Sprintf("⚡ Chain[%s]: Circuit closed for model %d (%s/%s)",
			c.name, idx, c.models[idx].Provider(), c.models[idx].Model()))
	}

	c.failures[idx] = 0
	c.circuitOpen[idx] = false
}

// Generate tries each model in order until one succeeds.
// Timeout errors trigger fallback to next model with fresh context.
func (c *ChainLanguageModel) Generate(ctx context.Context, call Call) (*Response, error) {
	var lastErr error
	var attemptedModels []string

	for i, model := range c.models {
		// Check circuit breaker
		if c.isCircuitOpen(i) {
			xlog.Info(fmt.Sprintf("🔀 Chain[%s]: Skipping model %d (%s/%s) - circuit open",
				c.name, i, model.Provider(), model.Model()))
			continue
		}

		modelName := fmt.Sprintf("%s/%s", model.Provider(), model.Model())
		attemptedModels = append(attemptedModels, modelName)

		xlog.Info(fmt.Sprintf("🔀 Chain[%s]: Trying model %d/%d (%s)",
			c.name, i+1, len(c.models), modelName))

		// Create fresh context for each model to allow fallback on timeout
		modelCtx := ctx
		if ctx.Err() != nil {
			// Parent context already cancelled/expired, create fresh background context
			// This allows fallback models to still attempt generation
			xlog.Info(fmt.Sprintf("🔄 Chain[%s]: Parent context expired, using fresh context for fallback", c.name))
			modelCtx = context.Background()
		}

		resp, err := model.Generate(modelCtx, call)
		if err == nil {
			c.recordSuccess(i)
			if i > 0 {
				xlog.Info(fmt.Sprintf("✅ Chain[%s]: Fallback succeeded with model %s", c.name, modelName))
			}
			return resp, nil
		}

		c.recordFailure(i)
		lastErr = err
		xlog.Warn(fmt.Sprintf("⚠️  Chain[%s]: Model %s failed: %v", c.name, modelName, err))

		// Only stop if user explicitly cancelled (not timeout)
		// Timeout should trigger fallback to local model
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		}
		// For DeadlineExceeded (timeout), continue to next model
	}

	return nil, fmt.Errorf("all models failed (tried: %s): %w",
		strings.Join(attemptedModels, ", "), lastErr)
}

// Stream tries each model in order.
// If a model fails during setup or *during* streaming, it seamlessly falls back to the next model.
func (c *ChainLanguageModel) Stream(ctx context.Context, call Call) (StreamResponse, error) {
	// We return a single stream iterator that internally manages switching between models
	return func(yield func(StreamPart) bool) {
		var lastErr error
		var attemptedModels []string

		// Setup context for the chain - we might need to recreate context if one model fails/timeouts
		// to give the next model a fresh start, but usually we respect the parent context.
		// If dealing with timeouts per model, complex context handling is needed.
		// For now, we assume the parent context is sufficient or the models handle their own timeouts.

		for i, model := range c.models {
			// Check circuit breaker
			if c.isCircuitOpen(i) {
				continue
			}

			modelName := fmt.Sprintf("%s/%s", model.Provider(), model.Model())
			attemptedModels = append(attemptedModels, modelName)

			xlog.Info(fmt.Sprintf("🔀 Chain[%s]: Trying stream with model %d/%d (%s)", c.name, i+1, len(c.models), modelName))

			// Attempt to start streaming
			stream, err := model.Stream(ctx, call)
			if err != nil {
				c.recordFailure(i)
				lastErr = err
				xlog.Warn(fmt.Sprintf("⚠️  Chain[%s]: Setup failed for %s: %v", c.name, modelName, err))
				continue // Try next model
			}

			// Stream started successfully. Now we consume it.
			// Ideally we want to pass through success parts, but catch error parts.

			streamFailed := false
			for part := range stream {
				// Check for errors in the stream
				if part.Type == StreamPartTypeError || part.Error != nil {
					// Mid-stream failure detected!
					c.recordFailure(i)
					lastErr = part.Error
					if lastErr == nil {
						lastErr = fmt.Errorf("unknown stream error")
					}
					xlog.Warn(fmt.Sprintf("⚠️  Chain[%s]: Mid-stream error from %s: %v. Switching to next model...", c.name, modelName, lastErr))

					streamFailed = true
					break // Stop consuming this stream, move to next model
				}

				// Yield successful part
				if !yield(part) {
					return // Consumer stopped
				}
			}

			if !streamFailed {
				// Model finished successfully
				c.recordSuccess(i)
				return // We are done
			}

			// If we are here, streamFailed was true.
			// We loop back to try the next model.
		}

		// If we exhausted all models, yield the final error
		finalErr := fmt.Errorf("all chain models failed (tried: %s): %w", strings.Join(attemptedModels, ", "), lastErr)
		yield(StreamPart{
			Type:  StreamPartTypeError,
			Error: finalErr,
		})
	}, nil
}

// GenerateObject tries each model in order until one succeeds.
func (c *ChainLanguageModel) GenerateObject(ctx context.Context, call ObjectCall) (*ObjectResponse, error) {
	var lastErr error

	for i, model := range c.models {
		if c.isCircuitOpen(i) {
			continue
		}

		resp, err := model.GenerateObject(ctx, call)
		if err == nil {
			c.recordSuccess(i)
			return resp, nil
		}

		c.recordFailure(i)
		lastErr = err

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("all models failed: %w", lastErr)
}

// StreamObject tries each model in order until one succeeds.
func (c *ChainLanguageModel) StreamObject(ctx context.Context, call ObjectCall) (ObjectStreamResponse, error) {
	var lastErr error

	for i, model := range c.models {
		if c.isCircuitOpen(i) {
			continue
		}

		stream, err := model.StreamObject(ctx, call)
		if err == nil {
			c.recordSuccess(i)
			return stream, nil
		}

		c.recordFailure(i)
		lastErr = err

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("all models failed: %w", lastErr)
}

// Provider returns "chain" as the provider name.
func (c *ChainLanguageModel) Provider() string {
	return "chain"
}

// Model returns a descriptive name of the chain.
func (c *ChainLanguageModel) Model() string {
	return c.name
}

// PrimaryModel returns the first (primary) model in the chain.
func (c *ChainLanguageModel) PrimaryModel() LanguageModel {
	if len(c.models) > 0 {
		return c.models[0]
	}
	return nil
}

// Models returns all models in the chain.
func (c *ChainLanguageModel) Models() []LanguageModel {
	return c.models
}
