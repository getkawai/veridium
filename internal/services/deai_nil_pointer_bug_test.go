package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNilPointerBug demonstrates Bug #1: Nil Pointer Dereference Risk
// After fix, these should return errors instead of panicking
func TestNilPointerBug(t *testing.T) {
	t.Run("GetVaultBalance_with_locked_wallet", func(t *testing.T) {
		// Setup: Create DeAIService with locked wallet
		wallet := &WalletService{
			currentAccount: nil, // Wallet is LOCKED
			address:        "",
		}

		deai := &DeAIService{
			wallet: wallet,
			reader: nil, // Not needed for this test
			kv:     nil,
		}

		// AFTER FIX: Should return error, not panic
		result, err := deai.GetVaultBalance()

		assert.Empty(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ FIX VERIFIED: GetVaultBalance returns error when wallet is locked")
	})

	t.Run("DepositToVault_with_locked_wallet", func(t *testing.T) {
		wallet := &WalletService{
			currentAccount: nil, // Wallet is LOCKED
		}

		deai := &DeAIService{
			wallet: wallet,
		}

		// AFTER FIX: Should return error, not panic
		result, err := deai.DepositToVault("100")

		assert.Empty(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ FIX VERIFIED: DepositToVault returns error when wallet is locked")
	})

	t.Run("GetClaimableRewards_with_locked_wallet", func(t *testing.T) {
		wallet := &WalletService{
			currentAccount: nil, // Wallet is LOCKED
		}

		deai := &DeAIService{
			wallet: wallet,
		}

		// GOOD: GetClaimableRewards DOES check if wallet is locked
		result, err := deai.GetClaimableRewards()

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no wallet connected")

		t.Log("✅ GOOD: GetClaimableRewards properly checks for nil wallet")
	})
}

// TestInconsistentNilChecks shows the inconsistency across methods
func TestInconsistentNilChecks(t *testing.T) {
	wallet := &WalletService{
		currentAccount: nil, // Wallet is LOCKED
	}

	deai := &DeAIService{
		wallet: wallet,
	}

	t.Run("methods_with_nil_check", func(t *testing.T) {
		// These methods PROPERLY check for nil
		methods := []struct {
			name string
			fn   func() error
		}{
			{
				name: "GetClaimableRewards",
				fn: func() error {
					_, err := deai.GetClaimableRewards()
					return err
				},
			},
			{
				name: "ConfirmRewardClaim",
				fn: func() error {
					return deai.ConfirmRewardClaim(1)
				},
			},
			{
				name: "MarkClaimFailed",
				fn: func() error {
					return deai.MarkClaimFailed(1, "test")
				},
			},
		}

		for _, method := range methods {
			t.Run(method.name, func(t *testing.T) {
				err := method.fn()
				assert.Error(t, err, "Should return error, not panic")
				assert.Contains(t, err.Error(), "no wallet connected",
					"Should have proper error message")
				t.Logf("✅ %s properly checks for nil wallet", method.name)
			})
		}
	})

	t.Run("methods_NOW_WITH_nil_check", func(t *testing.T) {
		// These methods NOW check for nil after fix
		methods := []struct {
			name string
			fn   func() error
		}{
			{
				name: "GetVaultBalance",
				fn: func() error {
					_, err := deai.GetVaultBalance()
					return err
				},
			},
			{
				name: "DepositToVault",
				fn: func() error {
					_, err := deai.DepositToVault("100")
					return err
				},
			},
		}

		for _, method := range methods {
			t.Run(method.name, func(t *testing.T) {
				err := method.fn()
				assert.Error(t, err, "Should return error, not panic")
				assert.Contains(t, err.Error(), "wallet is locked",
					"Should have proper error message")
				t.Logf("✅ %s now properly checks for nil wallet", method.name)
			})
		}
	})
}

// TestRealWorldScenario demonstrates when this bug would occur
func TestNilPointerRealWorldScenario(t *testing.T) {
	t.Run("scenario_1_user_locks_wallet_then_navigates", func(t *testing.T) {
		// User flow:
		// 1. User unlocks wallet
		// 2. User navigates to wallet page - sees balance
		// 3. User locks wallet (for security)
		// 4. User navigates to another page that calls GetVaultBalance
		// 5. APP CRASHES! (nil pointer panic)

		wallet := &WalletService{}
		deai := &DeAIService{wallet: wallet}

		// Step 1-2: Wallet is unlocked (simulated)
		// wallet.currentAccount = someAccount
		// balance, _ := deai.GetVaultBalance() // Works fine

		// Step 3: User locks wallet
		wallet.LockWallet()

		// Step 4: App tries to get balance (e.g., background refresh)
		_, err := deai.GetVaultBalance()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ FIX VERIFIED: User locks wallet → App shows error message (no crash)")
	})

	t.Run("scenario_2_app_starts_without_unlock", func(t *testing.T) {
		// User flow:
		// 1. User opens app
		// 2. Wallet is locked by default
		// 3. App tries to load balance on startup
		// 4. APP CRASHES!

		wallet := &WalletService{
			currentAccount: nil, // Default state
		}
		deai := &DeAIService{wallet: wallet}

		// App tries to load balance before user unlocks
		_, err := deai.GetVaultBalance()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ FIX VERIFIED: App prompts to unlock instead of crashing")
	})

	t.Run("scenario_3_session_timeout", func(t *testing.T) {
		// User flow:
		// 1. User unlocks wallet
		// 2. User leaves app open
		// 3. Session times out → wallet auto-locks
		// 4. Background task tries to refresh balance
		// 5. APP CRASHES!

		wallet := &WalletService{}
		deai := &DeAIService{wallet: wallet}

		// Simulate session timeout
		wallet.LockWallet()

		// Background task tries to refresh
		_, err := deai.GetVaultBalance()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ FIX VERIFIED: Session timeout → Background task handles gracefully")
	})
}

// TestProperNilCheckPattern shows the correct pattern
func TestProperNilCheckPattern(t *testing.T) {
	t.Run("correct_pattern_example", func(t *testing.T) {
		wallet := &WalletService{
			currentAccount: nil,
		}
		deai := &DeAIService{wallet: wallet}

		// CORRECT PATTERN (used in GetClaimableRewards):
		// if s.wallet.currentAccount == nil {
		//     return nil, fmt.Errorf("no wallet connected")
		// }

		_, err := deai.GetClaimableRewards()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no wallet connected")

		t.Log("✅ Correct pattern: Check for nil, return error (no panic)")
	})

	t.Run("now_uses_correct_pattern", func(t *testing.T) {
		wallet := &WalletService{
			currentAccount: nil,
		}
		deai := &DeAIService{wallet: wallet}

		// NOW USES CORRECT PATTERN:
		// if s.wallet.currentAccount == nil {
		//     return "", fmt.Errorf("wallet is locked")
		// }

		_, err := deai.GetVaultBalance()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ Now uses correct pattern: Check for nil, return error")
	})
}

// TestComparisonWithWalletService shows WalletService does it correctly
func TestComparisonWithWalletService(t *testing.T) {
	t.Run("WalletService_SignMessage_checks_nil", func(t *testing.T) {
		wallet := &WalletService{
			currentAccount: nil,
		}

		// WalletService.SignMessage PROPERLY checks for nil
		_, err := wallet.SignMessage("test")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ WalletService.SignMessage properly checks for nil")
	})

	t.Run("WalletService_getTransactOpts_checks_nil", func(t *testing.T) {
		wallet := &WalletService{
			currentAccount: nil,
		}

		// WalletService.getTransactOpts PROPERLY checks for nil
		_, err := wallet.getTransactOpts(nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is locked")

		t.Log("✅ WalletService.getTransactOpts properly checks for nil")
	})

	t.Run("comparison", func(t *testing.T) {
		t.Log("📊 COMPARISON:")
		t.Log("   WalletService methods: Check for nil ✅")
		t.Log("   DeAIService methods:    Some check ✅, some don't ❌")
		t.Log("   Result: INCONSISTENT behavior across services")
	})
}
