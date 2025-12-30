package services

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// TestAddressChecksumBug demonstrates the actual bug that was fixed
func TestAddressChecksumBug(t *testing.T) {
	// Address dari keystore (lowercase, no 0x prefix)
	addressFromKeystore := "22cb1309794b475c5709dddb640b79533c38d924"

	// OLD METHOD (BUG): Simple concatenation - produces lowercase
	oldMethod := "0x" + addressFromKeystore
	
	// NEW METHOD (FIX): Proper checksumming - produces mixed case per EIP-55
	newMethod := common.HexToAddress(addressFromKeystore).Hex()

	t.Logf("Address from keystore: %s", addressFromKeystore)
	t.Logf("Old method (buggy):    %s", oldMethod)
	t.Logf("New method (fixed):    %s", newMethod)

	// The bug: old method produces lowercase, not checksummed
	assert.Equal(t, "0x22cb1309794b475c5709dddb640b79533c38d924", oldMethod,
		"Old method produces lowercase (NOT checksummed)")
	
	// The fix: new method produces proper EIP-55 checksum
	assert.Equal(t, "0x22CB1309794B475C5709Dddb640B79533C38D924", newMethod,
		"New method produces proper EIP-55 checksum")
	
	// They are NOT equal!
	assert.NotEqual(t, oldMethod, newMethod,
		"This proves the bug: old method != new method")
}

// TestChecksumValidation shows why checksums matter
func TestChecksumValidation(t *testing.T) {
	correctChecksum := "0x22CB1309794B475C5709Dddb640B79533C38D924"
	incorrectChecksum := "0x22cb1309794b475c5709dddb640b79533c38d924" // all lowercase
	
	// Ethereum clients can validate checksums
	// If you send to wrong checksum, some wallets will reject it
	
	t.Run("correct_checksum", func(t *testing.T) {
		addr := common.HexToAddress(correctChecksum)
		recomputed := addr.Hex()
		
		assert.Equal(t, correctChecksum, recomputed,
			"Correct checksum should match after recomputation")
	})
	
	t.Run("incorrect_checksum_gets_fixed", func(t *testing.T) {
		addr := common.HexToAddress(incorrectChecksum)
		recomputed := addr.Hex()
		
		// The incorrect (lowercase) gets converted to correct checksum
		assert.NotEqual(t, incorrectChecksum, recomputed,
			"Incorrect checksum gets fixed")
		assert.Equal(t, correctChecksum, recomputed,
			"Gets converted to correct checksum")
	})
}

// TestRealWorldScenario demonstrates the real-world impact
func TestRealWorldScenario(t *testing.T) {
	t.Run("scenario_1_import_from_metamask", func(t *testing.T) {
		// User exports keystore from MetaMask
		// MetaMask stores address in lowercase (no 0x)
		keystoreAddress := "995c8e149b803009c2e8dc696d5d5d3429a30a05"
		
		// OLD BUG: Would store as lowercase
		buggyAddress := "0x" + keystoreAddress
		
		// NEW FIX: Stores as checksummed
		fixedAddress := common.HexToAddress(keystoreAddress).Hex()
		
		t.Logf("Keystore address: %s", keystoreAddress)
		t.Logf("Buggy storage:    %s", buggyAddress)
		t.Logf("Fixed storage:    %s", fixedAddress)
		
		// Problem: If we store buggy address, when user checks on Etherscan
		// or other wallets, they see different format (checksummed)
		// This causes confusion and potential security issues
		
		assert.NotEqual(t, buggyAddress, fixedAddress,
			"Bug causes inconsistent address format")
	})
	
	t.Run("scenario_2_address_comparison", func(t *testing.T) {
		// Two users import same wallet
		// One gets lowercase, one gets checksummed
		// String comparison fails even though it's same address!
		
		addr1 := "0x22cb1309794b475c5709dddb640b79533c38d924" // lowercase
		addr2 := "0x22CB1309794B475C5709Dddb640B79533C38D924" // checksummed
		
		// String comparison fails
		assert.NotEqual(t, addr1, addr2,
			"String comparison fails for same address!")
		
		// But they're the same address
		assert.True(t, common.HexToAddress(addr1) == common.HexToAddress(addr2),
			"They represent the same address")
		
		// Solution: Always store in checksummed format
		normalized1 := common.HexToAddress(addr1).Hex()
		normalized2 := common.HexToAddress(addr2).Hex()
		
		assert.Equal(t, normalized1, normalized2,
			"Checksummed format ensures consistency")
	})
	
	t.Run("scenario_3_security_typo_detection", func(t *testing.T) {
		// EIP-55 checksum can detect typos
		correctAddr := "0x22CB1309794B475C5709Dddb640B79533C38D924"
		
		// Attacker changes one character but keeps checksum pattern
		// (this is a hypothetical typo/attack)
		typoAddr := "0x22CB1309794B475C5709Dddb640B79533C38D925" // changed last digit
		
		// Recompute checksums
		correctRecomputed := common.HexToAddress(correctAddr).Hex()
		typoRecomputed := common.HexToAddress(typoAddr).Hex()
		
		// Original correct address maintains checksum
		assert.Equal(t, correctAddr, correctRecomputed,
			"Correct address maintains checksum")
		
		// Typo address gets different checksum
		// (in real scenario, wallet would warn user about invalid checksum)
		t.Logf("Typo address input:      %s", typoAddr)
		t.Logf("Typo address recomputed: %s", typoRecomputed)
		
		// The addresses are different
		assert.NotEqual(t, correctAddr, typoAddr,
			"Checksum helps detect address changes")
	})
}

// TestFileSystemConsistency shows filesystem-related issues
func TestFileSystemConsistency(t *testing.T) {
	// On case-insensitive filesystems (macOS, Windows default)
	// These would be treated as same file:
	// - 0x22cb1309794b475c5709dddb640b79533c38d924.json
	// - 0x22CB1309794B475C5709Dddb640B79533C38D924.json
	
	// But in our metadata and database, they might be different strings!
	// This causes bugs where:
	// 1. File exists but lookup fails (case mismatch)
	// 2. Duplicate detection fails
	// 3. Database queries fail
	
	addr1 := "0x22cb1309794b475c5709dddb640b79533c38d924"
	addr2 := "0x22CB1309794B475C5709Dddb640B79533C38D924"
	
	// Without checksumming, these are different strings
	assert.NotEqual(t, addr1, addr2,
		"Different strings in database/code")
	
	// But same file on case-insensitive FS
	// This causes bugs!
	
	// Solution: Always use checksummed format
	normalized1 := common.HexToAddress(addr1).Hex()
	normalized2 := common.HexToAddress(addr2).Hex()
	
	assert.Equal(t, normalized1, normalized2,
		"Checksummed format prevents filesystem bugs")
}

