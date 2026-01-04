package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/cashbackdistributor"
	"github.com/kawai-network/veridium/pkg/store"
)

// CashbackSettlement handles weekly cashback settlement
type CashbackSettlement struct {
	client      *ethclient.Client
	distributor *cashbackdistributor.DepositCashbackDistributor
	kvStore     *store.KVStore
	privateKey  string
}

// NewCashbackSettlement creates a new cashback settlement service
func NewCashbackSettlement(kvStore *store.KVStore, privateKey string) (*CashbackSettlement, error) {
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	distributorAddr := common.HexToAddress(constant.CashbackDistributorAddress)
	distributor, err := cashbackdistributor.NewDepositCashbackDistributor(distributorAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load distributor contract: %w", err)
	}

	return &CashbackSettlement{
		client:      client,
		distributor: distributor,
		kvStore:     kvStore,
		privateKey:  privateKey,
	}, nil
}

// CashbackLeaf represents a leaf in the Merkle tree
type CashbackLeaf struct {
	Period      uint64
	UserAddress common.Address
	Amount      *big.Int
}

// SettleCashback performs weekly cashback settlement
func (cs *CashbackSettlement) SettleCashback(ctx context.Context, period uint64) error {
	log.Printf("🔄 [CashbackSettlement] Starting settlement for period %d", period)

	// 1. Collect all pending cashback records
	leaves, err := cs.collectPendingCashback(ctx, period)
	if err != nil {
		return fmt.Errorf("failed to collect pending cashback: %w", err)
	}

	if len(leaves) == 0 {
		log.Printf("⚠️  [CashbackSettlement] No pending cashback for period %d", period)
		return nil
	}

	log.Printf("📊 [CashbackSettlement] Collected %d cashback records", len(leaves))

	// 2. Generate Merkle tree
	merkleRoot, proofs, err := cs.generateMerkleTree(leaves)
	if err != nil {
		return fmt.Errorf("failed to generate Merkle tree: %w", err)
	}

	log.Printf("🌳 [CashbackSettlement] Merkle root: %s", common.Bytes2Hex(merkleRoot[:]))

	// 3. Store proofs in KV
	if err := cs.storeProofs(ctx, period, proofs); err != nil {
		return fmt.Errorf("failed to store proofs: %w", err)
	}

	// 4. Set Merkle root on-chain
	if err := cs.setMerkleRoot(ctx, period, merkleRoot); err != nil {
		return fmt.Errorf("failed to set Merkle root: %w", err)
	}

	log.Printf("✅ [CashbackSettlement] Settlement complete for period %d", period)
	return nil
}

// collectPendingCashback collects all pending cashback records for a period
func (cs *CashbackSettlement) collectPendingCashback(ctx context.Context, period uint64) ([]CashbackLeaf, error) {
	// List all cashback keys from KV
	// Note: This is a simplified implementation. In production, you'd want to:
	// 1. Use KV list API with pagination
	// 2. Filter by period
	// 3. Only include unclaimed records

	// For now, we'll use a marker key to track which records belong to which period
	periodKey := fmt.Sprintf("cashback_period:%d:users", period)
	data, err := cs.kvStore.GetMarketplaceData(ctx, periodKey)
	if err != nil {
		// Distinguish between "key not found" vs real errors
		if err.Error() == "key not found" || err.Error() == "not found" {
			// No records for this period yet - this is expected
			return []CashbackLeaf{}, nil
		}
		// Real error (network, permission, etc.)
		return nil, fmt.Errorf("failed to fetch period data: %w", err)
	}

	var userAddresses []string
	if err := json.Unmarshal(data, &userAddresses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user addresses: %w", err)
	}

	leaves := make([]CashbackLeaf, 0, len(userAddresses))
	for _, userAddr := range userAddresses {
		// Get user's pending cashback for this period
		stats, err := cs.kvStore.GetCashbackStats(ctx, userAddr)
		if err != nil {
			log.Printf("⚠️  [CashbackSettlement] Failed to get stats for %s: %v", userAddr, err)
			continue
		}

		// Only include if there's pending cashback
		pendingAmount := new(big.Int)
		if stats.PendingCashback != "" && stats.PendingCashback != "0" {
			if _, ok := pendingAmount.SetString(stats.PendingCashback, 10); !ok {
				log.Printf("⚠️  [CashbackSettlement] Invalid PendingCashback for %s: %s", userAddr, stats.PendingCashback)
				continue
			}
		}

		if pendingAmount.Cmp(big.NewInt(0)) > 0 {
			leaves = append(leaves, CashbackLeaf{
				Period:      period,
				UserAddress: common.HexToAddress(userAddr),
				Amount:      pendingAmount,
			})
		}
	}

	return leaves, nil
}

// generateMerkleTree generates a Merkle tree and proofs for all leaves
func (cs *CashbackSettlement) generateMerkleTree(leaves []CashbackLeaf) ([32]byte, map[string][][]byte, error) {
	if len(leaves) == 0 {
		return [32]byte{}, nil, fmt.Errorf("no leaves to process")
	}

	// 1. Create leaf hashes
	leafHashes := make([][]byte, len(leaves))
	for i, leaf := range leaves {
		// Hash: keccak256(abi.encodePacked(period, user, amount))
		hash := crypto.Keccak256(
			common.LeftPadBytes(big.NewInt(int64(leaf.Period)).Bytes(), 32),
			leaf.UserAddress.Bytes(),
			common.LeftPadBytes(leaf.Amount.Bytes(), 32),
		)
		leafHashes[i] = hash
	}

	// 2. Build Merkle tree
	tree := buildMerkleTree(leafHashes)
	if len(tree) == 0 {
		return [32]byte{}, nil, fmt.Errorf("failed to build Merkle tree")
	}

	// Root is the last element
	var root [32]byte
	copy(root[:], tree[len(tree)-1])

	// 3. Generate proofs for each leaf
	proofs := make(map[string][][]byte)
	for i, leaf := range leaves {
		proof := generateMerkleProof(tree, i, len(leafHashes))
		proofs[leaf.UserAddress.Hex()] = proof
	}

	return root, proofs, nil
}

// buildMerkleTree builds a Merkle tree from leaf hashes
func buildMerkleTree(leaves [][]byte) [][]byte {
	if len(leaves) == 0 {
		return nil
	}

	// Sort leaves for deterministic tree
	sortedLeaves := make([][]byte, len(leaves))
	copy(sortedLeaves, leaves)
	sort.Slice(sortedLeaves, func(i, j int) bool {
		return string(sortedLeaves[i]) < string(sortedLeaves[j])
	})

	tree := make([][]byte, 0, len(sortedLeaves)*2)
	tree = append(tree, sortedLeaves...)

	// Build tree level by level
	for len(sortedLeaves) > 1 {
		nextLevel := make([][]byte, 0, (len(sortedLeaves)+1)/2)

		for i := 0; i < len(sortedLeaves); i += 2 {
			if i+1 < len(sortedLeaves) {
				// Pair exists - sort before hashing (required by OpenZeppelin)
				left, right := sortedLeaves[i], sortedLeaves[i+1]
				if string(left) > string(right) {
					left, right = right, left
				}
				hash := crypto.Keccak256(left, right)
				nextLevel = append(nextLevel, hash)
			} else {
				// Odd node - promote to next level
				nextLevel = append(nextLevel, sortedLeaves[i])
			}
		}

		tree = append(tree, nextLevel...)
		sortedLeaves = nextLevel
	}

	return tree
}

// generateMerkleProof generates a Merkle proof for a leaf at given index
func generateMerkleProof(tree [][]byte, leafIndex int, leafCount int) [][]byte {
	proof := make([][]byte, 0)
	index := leafIndex

	for levelSize := leafCount; levelSize > 1; {
		levelStart := len(tree) - levelSize - (levelSize+1)/2
		if levelStart < 0 {
			levelStart = 0
		}

		// Find sibling
		var sibling []byte
		if index%2 == 0 {
			// Left node - sibling is right
			if index+1 < levelSize {
				sibling = tree[levelStart+index+1]
			}
		} else {
			// Right node - sibling is left
			sibling = tree[levelStart+index-1]
		}

		if sibling != nil {
			proof = append(proof, sibling)
		}

		index /= 2
		levelSize = (levelSize + 1) / 2
	}

	return proof
}

// storeProofs stores Merkle proofs in KV for user claims
func (cs *CashbackSettlement) storeProofs(ctx context.Context, period uint64, proofs map[string][][]byte) error {
	for userAddr, proof := range proofs {
		key := fmt.Sprintf("cashback_proof:%d:%s", period, userAddr)
		data, err := json.Marshal(proof)
		if err != nil {
			return fmt.Errorf("failed to marshal proof for %s: %w", userAddr, err)
		}

		if err := cs.kvStore.StoreMarketplaceData(ctx, key, data); err != nil {
			return fmt.Errorf("failed to store proof for %s: %w", userAddr, err)
		}
	}

	log.Printf("✅ [CashbackSettlement] Stored %d proofs", len(proofs))
	return nil
}

// setMerkleRoot sets the Merkle root on-chain
func (cs *CashbackSettlement) setMerkleRoot(ctx context.Context, period uint64, merkleRoot [32]byte) error {
	// Parse private key
	privateKey, err := crypto.HexToECDSA(cs.privateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get chain ID
	chainID, err := cs.client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Check current period
	currentPeriod, err := cs.distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period: %w", err)
	}

	var tx *types.Transaction
	if period > currentPeriod.Uint64() {
		// Advance to new period
		tx, err = cs.distributor.AdvancePeriod(auth, merkleRoot)
		if err != nil {
			return fmt.Errorf("failed to advance period: %w", err)
		}
		log.Printf("📝 [CashbackSettlement] Advance period tx: %s", tx.Hash().Hex())
	} else {
		// Update existing period
		tx, err = cs.distributor.SetMerkleRoot(auth, merkleRoot)
		if err != nil {
			return fmt.Errorf("failed to set Merkle root: %w", err)
		}
		log.Printf("📝 [CashbackSettlement] Set Merkle root tx: %s", tx.Hash().Hex())
	}

	// Wait for confirmation
	receipt, err := bind.WaitMined(ctx, cs.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for tx confirmation: %w", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	log.Printf("✅ [CashbackSettlement] Merkle root set on-chain (block %d)", receipt.BlockNumber.Uint64())
	return nil
}

// TrackUserForPeriod adds a user to the period's user list for settlement
func (cs *CashbackSettlement) TrackUserForPeriod(ctx context.Context, period uint64, userAddress string) error {
	periodKey := fmt.Sprintf("cashback_period:%d:users", period)
	
	// Get existing users
	data, err := cs.kvStore.GetMarketplaceData(ctx, periodKey)
	var users []string
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &users); err != nil {
			return fmt.Errorf("failed to unmarshal users: %w", err)
		}
	}
	
	// Check if user already tracked
	for _, u := range users {
		if u == userAddress {
			return nil // Already tracked
		}
	}
	
	// Add user
	users = append(users, userAddress)
	data, err = json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}
	
	if err := cs.kvStore.StoreMarketplaceData(ctx, periodKey, data); err != nil {
		return fmt.Errorf("failed to store users: %w", err)
	}
	
	return nil
}

