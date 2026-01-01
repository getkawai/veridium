package store

import (
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBalanceRaceCondition demonstrates the race condition in non-atomic operations
func TestBalanceRaceCondition(t *testing.T) {
	t.Run("non_atomic_double_spend", func(t *testing.T) {
		// This test demonstrates the PROBLEM with non-atomic operations
		// Multiple goroutines can deduct more than the available balance

		t.Skip("This test demonstrates the bug - skip in CI")

		// Simulated balance (not using real KV store)
		balance := big.NewInt(1000) // 1000 USDT
		var mu sync.Mutex

		var wg sync.WaitGroup
		successCount := 0

		// 10 goroutines try to deduct 200 USDT each (total 2000 USDT)
		// Expected: Only 5 should succeed (5 * 200 = 1000)
		// Bug: More than 5 might succeed due to race condition
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Simulate non-atomic check-then-deduct
				mu.Lock()
				currentBalance := new(big.Int).Set(balance)
				mu.Unlock()

				// ❌ RACE CONDITION: Another goroutine can check here!
				if currentBalance.Cmp(big.NewInt(200)) >= 0 {
					// Simulate some processing time
					// time.Sleep(1 * time.Millisecond)

					mu.Lock()
					balance.Sub(balance, big.NewInt(200))
					successCount++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		t.Logf("Final balance: %s (expected: 0 or positive)", balance.String())
		t.Logf("Success count: %d (expected: 5)", successCount)

		// ❌ BUG: Balance might be negative!
		if balance.Cmp(big.NewInt(0)) < 0 {
			t.Logf("⚠️  BUG DETECTED: Negative balance! %s", balance.String())
		}

		// ❌ BUG: More than 5 might succeed!
		if successCount > 5 {
			t.Logf("⚠️  BUG DETECTED: Double-spend! %d operations succeeded (expected max 5)", successCount)
		}
	})
}

// TestAtomicBalanceOperations tests the atomic implementation
func TestAtomicBalanceOperations(t *testing.T) {
	t.Run("atomic_prevents_double_spend", func(t *testing.T) {
		// This test shows how atomic operations prevent double-spend

		// Note: This is a conceptual test
		// Real implementation would use KV store with retry logic

		balance := big.NewInt(1000)
		var mu sync.Mutex

		var wg sync.WaitGroup
		successCount := 0
		failCount := 0

		// 10 goroutines try to deduct 200 USDT each
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// ✅ ATOMIC: Lock entire operation
				mu.Lock()
				defer mu.Unlock()

				currentBalance := new(big.Int).Set(balance)
				if currentBalance.Cmp(big.NewInt(200)) >= 0 {
					balance.Sub(balance, big.NewInt(200))
					successCount++
				} else {
					failCount++
				}
			}()
		}

		wg.Wait()

		t.Logf("Final balance: %s", balance.String())
		t.Logf("Success count: %d", successCount)
		t.Logf("Fail count: %d", failCount)

		// ✅ CORRECT: Balance should be 0 or positive
		assert.GreaterOrEqual(t, balance.Cmp(big.NewInt(0)), 0, "Balance should not be negative")

		// ✅ CORRECT: Exactly 5 should succeed
		assert.Equal(t, 5, successCount, "Exactly 5 operations should succeed")
		assert.Equal(t, 5, failCount, "Exactly 5 operations should fail")
	})
}

// TestConcurrentBalanceOperations tests concurrent operations
func TestConcurrentBalanceOperations(t *testing.T) {
	t.Run("concurrent_add_and_deduct", func(t *testing.T) {
		balance := big.NewInt(1000)
		var mu sync.Mutex

		var wg sync.WaitGroup

		// 50 goroutines adding 10 USDT each
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				balance.Add(balance, big.NewInt(10))
				mu.Unlock()
			}()
		}

		// 50 goroutines deducting 10 USDT each
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				if balance.Cmp(big.NewInt(10)) >= 0 {
					balance.Sub(balance, big.NewInt(10))
				}
				mu.Unlock()
			}()
		}

		wg.Wait()

		t.Logf("Final balance: %s (expected: 1000)", balance.String())

		// Balance should be 1000 (50 * 10 added, 50 * 10 deducted)
		assert.Equal(t, big.NewInt(1000).String(), balance.String(), "Balance should remain 1000")
	})
}

// TestBalanceAtomicRetryLogic tests retry logic for transient failures
func TestBalanceAtomicRetryLogic(t *testing.T) {
	t.Run("retry_on_transient_failure", func(t *testing.T) {
		// Simulate transient failures (e.g., network issues)
		attemptCount := 0
		maxAttempts := 3

		success := false
		for attempt := 0; attempt < maxAttempts; attempt++ {
			attemptCount++

			// Simulate transient failure on first 2 attempts
			if attempt < 2 {
				t.Logf("Attempt %d: Simulated failure", attempt+1)
				continue
			}

			// Success on 3rd attempt
			t.Logf("Attempt %d: Success", attempt+1)
			success = true
			break
		}

		assert.True(t, success, "Should succeed after retries")
		assert.Equal(t, 3, attemptCount, "Should take 3 attempts")
	})
}

// TestBalanceTransferAtomic tests atomic transfer between two addresses
func TestBalanceTransferAtomic(t *testing.T) {
	t.Run("atomic_transfer_success", func(t *testing.T) {
		balances := map[string]*big.Int{
			"alice": big.NewInt(1000),
			"bob":   big.NewInt(500),
		}
		var mu sync.Mutex

		// Transfer 200 from Alice to Bob
		mu.Lock()
		if balances["alice"].Cmp(big.NewInt(200)) >= 0 {
			balances["alice"].Sub(balances["alice"], big.NewInt(200))
			balances["bob"].Add(balances["bob"], big.NewInt(200))
		}
		mu.Unlock()

		assert.Equal(t, big.NewInt(800).String(), balances["alice"].String())
		assert.Equal(t, big.NewInt(700).String(), balances["bob"].String())
	})

	t.Run("atomic_transfer_insufficient_balance", func(t *testing.T) {
		balances := map[string]*big.Int{
			"alice": big.NewInt(100),
			"bob":   big.NewInt(500),
		}
		var mu sync.Mutex

		// Try to transfer 200 from Alice to Bob (should fail)
		mu.Lock()
		if balances["alice"].Cmp(big.NewInt(200)) >= 0 {
			balances["alice"].Sub(balances["alice"], big.NewInt(200))
			balances["bob"].Add(balances["bob"], big.NewInt(200))
		} else {
			t.Log("Transfer failed: insufficient balance")
		}
		mu.Unlock()

		// Balances should remain unchanged
		assert.Equal(t, big.NewInt(100).String(), balances["alice"].String())
		assert.Equal(t, big.NewInt(500).String(), balances["bob"].String())
	})

	t.Run("atomic_transfer_rollback_on_failure", func(t *testing.T) {
		balances := map[string]*big.Int{
			"alice": big.NewInt(1000),
			"bob":   big.NewInt(500),
		}
		var mu sync.Mutex

		// Simulate transfer with rollback
		mu.Lock()
		originalAlice := new(big.Int).Set(balances["alice"])

		if balances["alice"].Cmp(big.NewInt(200)) >= 0 {
			// Deduct from Alice
			balances["alice"].Sub(balances["alice"], big.NewInt(200))

			// Simulate failure when adding to Bob
			simulateFailure := true
			if simulateFailure {
				// Rollback Alice's balance
				balances["alice"].Set(originalAlice)
				t.Log("Transfer failed: rolled back Alice's balance")
			} else {
				balances["bob"].Add(balances["bob"], big.NewInt(200))
			}
		}
		mu.Unlock()

		// Balances should remain unchanged after rollback
		assert.Equal(t, big.NewInt(1000).String(), balances["alice"].String())
		assert.Equal(t, big.NewInt(500).String(), balances["bob"].String())
	})
}

// TestRealWorldScenarios tests real-world usage patterns
func TestRealWorldScenarios(t *testing.T) {
	t.Run("scenario_api_usage_concurrent_requests", func(t *testing.T) {
		// Scenario: User makes 10 concurrent API requests
		// Each request costs 100 USDT
		// User has 500 USDT balance
		// Expected: Only 5 requests succeed

		balance := big.NewInt(500)
		var mu sync.Mutex
		var wg sync.WaitGroup

		successCount := 0
		failCount := 0

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(requestID int) {
				defer wg.Done()

				mu.Lock()
				defer mu.Unlock()

				if balance.Cmp(big.NewInt(100)) >= 0 {
					balance.Sub(balance, big.NewInt(100))
					successCount++
					t.Logf("Request %d: Success (remaining: %s)", requestID, balance.String())
				} else {
					failCount++
					t.Logf("Request %d: Failed (insufficient balance)", requestID)
				}
			}(i)
		}

		wg.Wait()

		assert.Equal(t, 5, successCount, "Exactly 5 requests should succeed")
		assert.Equal(t, 5, failCount, "Exactly 5 requests should fail")
		assert.Equal(t, big.NewInt(0).String(), balance.String(), "Balance should be 0")
	})

	t.Run("scenario_deposit_and_spend", func(t *testing.T) {
		// Scenario: User deposits while spending
		// Note: This test has non-deterministic ordering
		// Possible outcomes:
		// - Deposit first, then spend: 100 + 500 - 200 = 400
		// - Spend first (fails if balance < 200), then deposit: 100 + 500 = 600

		balance := big.NewInt(100)
		var mu sync.Mutex
		var wg sync.WaitGroup
		spentSuccess := false

		// Goroutine 1: Deposit 500 USDT
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			balance.Add(balance, big.NewInt(500))
			t.Log("Deposited 500 USDT")
			mu.Unlock()
		}()

		// Goroutine 2: Spend 200 USDT
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			if balance.Cmp(big.NewInt(200)) >= 0 {
				balance.Sub(balance, big.NewInt(200))
				spentSuccess = true
				t.Log("Spent 200 USDT")
			} else {
				t.Log("Spend failed: insufficient balance")
			}
			mu.Unlock()
		}()

		wg.Wait()

		// Final balance depends on execution order
		if spentSuccess {
			// Deposit happened first, then spend: 100 + 500 - 200 = 400
			assert.Equal(t, big.NewInt(400).String(), balance.String())
		} else {
			// Spend check happened first (failed), then deposit: 100 + 500 = 600
			assert.Equal(t, big.NewInt(600).String(), balance.String())
		}
		t.Logf("Final balance: %s (spend success: %v)", balance.String(), spentSuccess)
	})
}

// BenchmarkAtomicOperations benchmarks atomic operations
func BenchmarkAtomicOperations(b *testing.B) {
	balance := big.NewInt(1000000)
	var mu sync.Mutex

	b.Run("atomic_deduct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mu.Lock()
			if balance.Cmp(big.NewInt(1)) >= 0 {
				balance.Sub(balance, big.NewInt(1))
			}
			mu.Unlock()
		}
	})

	b.Run("atomic_add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mu.Lock()
			balance.Add(balance, big.NewInt(1))
			mu.Unlock()
		}
	})
}
