package lifecycle

import (
	"sync"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	lm := NewManager()

	if lm == nil {
		t.Fatal("NewManager returned nil")
	}

	if lm.Count() != 0 {
		t.Errorf("Expected 0 cleanups, got %d", lm.Count())
	}

	if lm.IsShutdown() {
		t.Error("Expected IsShutdown to be false")
	}
}

func TestRegisterCleanup(t *testing.T) {
	lm := NewManager()

	executed := false
	lm.RegisterCleanup("test", func() {
		executed = true
	})

	if lm.Count() != 1 {
		t.Errorf("Expected 1 cleanup, got %d", lm.Count())
	}

	lm.Shutdown()

	if !executed {
		t.Error("Cleanup function was not executed")
	}

	if !lm.IsShutdown() {
		t.Error("Expected IsShutdown to be true")
	}
}

func TestLIFOOrder(t *testing.T) {
	lm := NewManager()

	var order []int
	var mu sync.Mutex

	// Register 3 cleanup functions
	lm.RegisterCleanup("first", func() {
		mu.Lock()
		order = append(order, 1)
		mu.Unlock()
	})

	lm.RegisterCleanup("second", func() {
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
	})

	lm.RegisterCleanup("third", func() {
		mu.Lock()
		order = append(order, 3)
		mu.Unlock()
	})

	lm.Shutdown()

	// Should execute in reverse order: 3, 2, 1
	expected := []int{3, 2, 1}
	if len(order) != len(expected) {
		t.Fatalf("Expected %d executions, got %d", len(expected), len(order))
	}

	for i, v := range expected {
		if order[i] != v {
			t.Errorf("Expected order[%d] = %d, got %d", i, v, order[i])
		}
	}
}

func TestPanicRecovery(t *testing.T) {
	lm := NewManager()

	executed := false

	// Register a cleanup that panics
	lm.RegisterCleanup("panic", func() {
		panic("test panic")
	})

	// Register another cleanup that should still execute
	lm.RegisterCleanup("after-panic", func() {
		executed = true
	})

	// Should not panic
	lm.Shutdown()

	if !executed {
		t.Error("Cleanup after panic was not executed")
	}
}

func TestRegisterAfterShutdown(t *testing.T) {
	lm := NewManager()

	lm.Shutdown()

	// Should not panic, but should log warning
	lm.RegisterCleanup("late", func() {})

	if lm.Count() != 0 {
		t.Errorf("Expected 0 cleanups after shutdown, got %d", lm.Count())
	}
}

func TestMultipleShutdowns(t *testing.T) {
	lm := NewManager()

	executionCount := 0
	lm.RegisterCleanup("test", func() {
		executionCount++
	})

	lm.Shutdown()
	lm.Shutdown() // Second shutdown should be no-op

	if executionCount != 1 {
		t.Errorf("Expected cleanup to execute once, got %d times", executionCount)
	}
}

func TestGetRegisteredCleanups(t *testing.T) {
	lm := NewManager()

	lm.RegisterCleanup("first", func() {})
	lm.RegisterCleanup("second", func() {})
	lm.RegisterCleanup("third", func() {})

	names := lm.GetRegisteredCleanups()

	if len(names) != 3 {
		t.Fatalf("Expected 3 names, got %d", len(names))
	}

	expected := []string{"first", "second", "third"}
	for i, name := range expected {
		if names[i] != name {
			t.Errorf("Expected names[%d] = %s, got %s", i, name, names[i])
		}
	}
}

func TestRegisterCleanupWithTimeout(t *testing.T) {
	lm := NewManager()

	executed := false

	// Fast cleanup
	lm.RegisterCleanupWithTimeout("fast", func() {
		executed = true
	}, 1*time.Second)

	lm.Shutdown()

	if !executed {
		t.Error("Fast cleanup was not executed")
	}
}

func TestRegisterCleanupWithTimeoutExceeded(t *testing.T) {
	lm := NewManager()

	// Slow cleanup that exceeds timeout
	lm.RegisterCleanupWithTimeout("slow", func() {
		time.Sleep(2 * time.Second)
	}, 100*time.Millisecond)

	start := time.Now()
	lm.Shutdown()
	duration := time.Since(start)

	// Should timeout after ~100ms, not wait for 2 seconds
	if duration > 500*time.Millisecond {
		t.Errorf("Timeout not working, took %v", duration)
	}
}

func TestConcurrentRegistration(t *testing.T) {
	lm := NewManager()

	var wg sync.WaitGroup
	numGoroutines := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			lm.RegisterCleanup("cleanup", func() {})
		}(i)
	}

	wg.Wait()

	if lm.Count() != numGoroutines {
		t.Errorf("Expected %d cleanups, got %d", numGoroutines, lm.Count())
	}
}

func TestString(t *testing.T) {
	lm := NewManager()

	lm.RegisterCleanup("test1", func() {})
	lm.RegisterCleanup("test2", func() {})

	str := lm.String()
	if str == "" {
		t.Error("String() returned empty string")
	}

	// Should contain count and shutdown status
	if !contains(str, "2") || !contains(str, "false") {
		t.Errorf("String() output doesn't contain expected values: %s", str)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
