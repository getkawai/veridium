// Package lifecycle provides centralized application lifecycle management
package lifecycle

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// CleanupFunc represents a cleanup function with a name
type CleanupFunc struct {
	Name string
	Fn   func()
}

// Manager handles application lifecycle events and cleanup
type Manager struct {
	cleanupFuncs []CleanupFunc
	mu           sync.Mutex
	isShutdown   bool
}

// NewManager creates a new lifecycle manager
func NewManager() *Manager {
	return &Manager{
		cleanupFuncs: make([]CleanupFunc, 0),
		isShutdown:   false,
	}
}

// RegisterCleanup registers a cleanup function with a descriptive name
// Cleanup functions are executed in LIFO order (last registered, first executed)
func (lm *Manager) RegisterCleanup(name string, fn func()) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lm.isShutdown {
		log.Printf("⚠️  Warning: Attempted to register cleanup '%s' after shutdown", name)
		return
	}

	wrappedFn := func() {
		startTime := time.Now()
		log.Printf("🧹 Cleaning up: %s", name)

		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ Cleanup panic in %s: %v", name, r)
			} else {
				duration := time.Since(startTime)
				log.Printf("✅ Cleanup completed: %s (took %v)", name, duration)
			}
		}()

		fn()
	}

	lm.cleanupFuncs = append(lm.cleanupFuncs, CleanupFunc{
		Name: name,
		Fn:   wrappedFn,
	})

	log.Printf("📝 Registered cleanup: %s (total: %d)", name, len(lm.cleanupFuncs))
}

// RegisterCleanupWithTimeout registers a cleanup function with a timeout
func (lm *Manager) RegisterCleanupWithTimeout(name string, fn func(), timeout time.Duration) {
	lm.RegisterCleanup(name, func() {
		done := make(chan struct{})

		go func() {
			defer close(done)
			// Capture panic within the goroutine to prevent app crash andensure channel close
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ Cleanup panic in %s (timeout wrapper): %v", name, r)
				}
			}()
			fn()
		}()

		select {
		case <-done:
			// Cleanup completed successfully
		case <-time.After(timeout):
			log.Printf("⏱️  Cleanup timeout for %s after %v", name, timeout)
		}
	})
}

// Shutdown executes all registered cleanup functions in LIFO order
// This ensures that resources are cleaned up in the reverse order they were initialized
func (lm *Manager) Shutdown() {
	lm.mu.Lock()
	if lm.isShutdown {
		lm.mu.Unlock()
		log.Printf("⚠️  Warning: Shutdown already called")
		return
	}

	lm.isShutdown = true

	// Create a copy of cleanups to execute outside the lock
	// This prevents deadlocks if a cleanup function calls back into the manager
	cleanups := make([]CleanupFunc, len(lm.cleanupFuncs))
	copy(cleanups, lm.cleanupFuncs)
	lm.mu.Unlock()

	totalFuncs := len(cleanups)
	if totalFuncs == 0 {
		log.Printf("ℹ️  No cleanup functions registered")
		return
	}

	log.Printf("🚀 Starting shutdown sequence (%d cleanup functions)", totalFuncs)
	startTime := time.Now()

	// Execute cleanup functions in LIFO order (reverse)
	for i := len(cleanups) - 1; i >= 0; i-- {
		cleanup := cleanups[i]
		cleanup.Fn()
	}

	duration := time.Since(startTime)
	log.Printf("🎉 Shutdown sequence completed in %v", duration)
}

// GetRegisteredCleanups returns the names of all registered cleanup functions
// Useful for debugging and testing
func (lm *Manager) GetRegisteredCleanups() []string {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	names := make([]string, len(lm.cleanupFuncs))
	for i, cleanup := range lm.cleanupFuncs {
		names[i] = cleanup.Name
	}
	return names
}

// Count returns the number of registered cleanup functions
func (lm *Manager) Count() int {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return len(lm.cleanupFuncs)
}

// IsShutdown returns true if shutdown has been called
func (lm *Manager) IsShutdown() bool {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return lm.isShutdown
}

// String returns a string representation of the lifecycle manager
func (lm *Manager) String() string {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	return fmt.Sprintf("LifecycleManager{cleanups: %d, shutdown: %v}",
		len(lm.cleanupFuncs), lm.isShutdown)
}
