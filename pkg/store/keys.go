package store

import (
	"fmt"
	"strconv"
	"strings"
)

// Key generation functions for each namespace
// These functions ensure consistent key formatting across the codebase

// =============================================================================
// CONTRIBUTORS NAMESPACE
// =============================================================================
// Key format: {address}
// Example: 0x742d35Cc6634C0532925a3b844Bc454e4438f44e

// ContributorKey generates a key for contributor data
// Uses lowercase address for consistency
func ContributorKey(address string) string {
	return strings.ToLower(address)
}

// =============================================================================
// PROOFS NAMESPACE
// =============================================================================
// Key format: {address}:{periodID}
// Example: 0x742d35cc6634c0532925a3b844bc454e4438f44e:1704067200000000000

// ProofKey generates a key for Merkle proof data
// Uses lowercase address and periodID for consistent ordering
func ProofKey(address string, periodID int64) string {
	return fmt.Sprintf("%s:%d", strings.ToLower(address), periodID)
}

// ProofPrefixForAddress generates a prefix to list all proofs for an address
func ProofPrefixForAddress(address string) string {
	return fmt.Sprintf("%s:", strings.ToLower(address))
}

// ParseProofKey extracts address and periodID from a proof key
func ParseProofKey(key string) (address string, periodID int64, err error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid proof key format: %s", key)
	}

	address = parts[0]
	periodID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid period ID in key: %s", key)
	}

	return address, periodID, nil
}

// =============================================================================
// SETTLEMENTS NAMESPACE
// =============================================================================
// Key format: {periodID}
// Example: 1704067200000000000

// SettlementKey generates a key for settlement period data
func SettlementKey(periodID int64) string {
	return strconv.FormatInt(periodID, 10)
}

// ParseSettlementKey extracts periodID from a settlement key
func ParseSettlementKey(key string) (int64, error) {
	return strconv.ParseInt(key, 10, 64)
}

