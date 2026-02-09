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
	"github.com/kawai-network/veridium/internal/generate/abi/cashbackdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/x/constant"
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

// GetCurrentPeriod returns the current period from the contract
func (cs *CashbackSettlement) GetCurrentPeriod(ctx context.Context) (uint64, error) {
	period, err := cs.distributor.CurrentPeriod(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get current period: %w", err)
	}
	return period.Uint64(), nil
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

	// 3. Store merkle root in KV for queries
	rootKey := fmt.Sprintf("cashback_period:%d:merkle_root", period)
	rootJSON, err := json.Marshal("0x" + common.Bytes2Hex(merkleRoot[:]))
	if err != nil {
		return fmt.Errorf("failed to marshal merkle root: %w", err)
	}
	if err := cs.kvStore.StoreCashbackData(ctx, rootKey, rootJSON); err != nil {
		return fmt.Errorf("failed to store merkle root: %w", err)
	}

	// 4. Store proofs in KV
	if err := cs.storeProofs(ctx, period, proofs); err != nil {
		return fmt.Errorf("failed to store proofs: %w", err)
	}

	// 5. Set Merkle root on-chain
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
	data, err := cs.kvStore.GetCashbackData(ctx, periodKey)
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
	leafAddresses := make([]string, len(leaves))
	for i, leaf := range leaves {
		// Hash: keccak256(abi.encodePacked(period, user, amount))
		// IMPORTANT: Solidity's abi.encodePacked for uint256 uses FULL 32 bytes (left-padded with zeros)
		// This is different from Go's big.Int.Bytes() which returns minimal representation
		periodBytes := common.LeftPadBytes(big.NewInt(int64(leaf.Period)).Bytes(), 32) // 32 bytes
		addressBytes := leaf.UserAddress.Bytes()                                       // 20 bytes
		amountBytes := common.LeftPadBytes(leaf.Amount.Bytes(), 32)                    // 32 bytes

		hash := crypto.Keccak256(periodBytes, addressBytes, amountBytes)
		leafHashes[i] = hash
		leafAddresses[i] = leaf.UserAddress.Hex()
	}

	// 2. Build Merkle tree (this will sort leaves internally)
	tree, sortedIndices := buildMerkleTreeWithIndices(leafHashes)
	if len(tree) == 0 {
		return [32]byte{}, nil, fmt.Errorf("failed to build Merkle tree")
	}

	// Root is the last element
	var root [32]byte
	copy(root[:], tree[len(tree)-1])

	// 3. Generate proofs for each leaf using sorted indices
	proofs := make(map[string][][]byte)
	for originalIdx, sortedIdx := range sortedIndices {
		proof := generateMerkleProof(tree, sortedIdx, len(leafHashes))
		proofs[leafAddresses[originalIdx]] = proof
	}

	return root, proofs, nil
}

// buildMerkleTree builds a Merkle tree from leaf hashes
func buildMerkleTree(leaves [][]byte) [][]byte {
	tree, _ := buildMerkleTreeWithIndices(leaves)
	return tree
}

// buildMerkleTreeWithIndices builds a Merkle tree and returns mapping of original index -> sorted index
func buildMerkleTreeWithIndices(leaves [][]byte) ([][]byte, map[int]int) {
	if len(leaves) == 0 {
		return nil, nil
	}

	// Create index mapping before sorting
	type indexedLeaf struct {
		hash          []byte
		originalIndex int
	}

	indexed := make([]indexedLeaf, len(leaves))
	for i, leaf := range leaves {
		indexed[i] = indexedLeaf{hash: leaf, originalIndex: i}
	}

	// Sort leaves for deterministic tree
	sort.Slice(indexed, func(i, j int) bool {
		return string(indexed[i].hash) < string(indexed[j].hash)
	})

	// Build index mapping: originalIndex -> sortedIndex
	indexMap := make(map[int]int)
	sortedLeaves := make([][]byte, len(indexed))
	for sortedIdx, item := range indexed {
		sortedLeaves[sortedIdx] = item.hash
		indexMap[item.originalIndex] = sortedIdx
	}

	tree := make([][]byte, 0, len(sortedLeaves)*2)
	tree = append(tree, sortedLeaves...)

	// Build tree level by level
	currentLevel := sortedLeaves
	for len(currentLevel) > 1 {
		nextLevel := make([][]byte, 0, (len(currentLevel)+1)/2)

		for i := 0; i < len(currentLevel); i += 2 {
			if i+1 < len(currentLevel) {
				// Pair exists - sort before hashing (required by OpenZeppelin)
				left, right := currentLevel[i], currentLevel[i+1]
				if string(left) > string(right) {
					left, right = right, left
				}
				hash := crypto.Keccak256(left, right)
				nextLevel = append(nextLevel, hash)
			} else {
				// Odd node - promote to next level
				nextLevel = append(nextLevel, currentLevel[i])
			}
		}

		tree = append(tree, nextLevel...)
		currentLevel = nextLevel
	}

	return tree, indexMap
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
	// Get user amounts from pending cashback
	leaves, err := cs.collectPendingCashback(ctx, period)
	if err != nil {
		return fmt.Errorf("failed to get user amounts: %w", err)
	}

	// Create map of user address -> amount
	userAmounts := make(map[string]string)
	for _, leaf := range leaves {
		userAmounts[leaf.UserAddress.Hex()] = leaf.Amount.String()
	}

	for userAddr, proof := range proofs {
		// Convert proof bytes to hex strings for JSON
		proofHex := make([]string, len(proof))
		for i, p := range proof {
			proofHex[i] = "0x" + common.Bytes2Hex(p)
		}

		// Store proof with amount and claimed status
		proofRecord := struct {
			Proof   []string `json:"proof"`
			Amount  string   `json:"amount"`
			Claimed bool     `json:"claimed"`
		}{
			Proof:   proofHex,
			Amount:  userAmounts[userAddr],
			Claimed: false,
		}

		key := fmt.Sprintf("cashback_proof:%d:%s", period, userAddr)
		data, err := json.Marshal(proofRecord)
		if err != nil {
			return fmt.Errorf("failed to marshal proof for %s: %w", userAddr, err)
		}

		if err := cs.kvStore.StoreCashbackData(ctx, key, data); err != nil {
			return fmt.Errorf("failed to store proof for %s: %w", userAddr, err)
		}
	}

	log.Printf("✅ [CashbackSettlement] Stored %d proofs", len(proofs))
	return nil
}

// setMerkleRoot sets the Merkle root on-chain
func (cs *CashbackSettlement) setMerkleRoot(ctx context.Context, period uint64, merkleRoot [32]byte) error {
	// Parse private key (strip 0x prefix if present)
	privateKeyHex := cs.privateKey
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
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
		// Set root for specific period (allows retroactive settlements)
		tx, err = cs.distributor.SetPeriodMerkleRoot(auth, big.NewInt(int64(period)), merkleRoot)
		if err != nil {
			return fmt.Errorf("failed to set period Merkle root: %w", err)
		}
		log.Printf("📝 [CashbackSettlement] Set period %d Merkle root tx: %s", period, tx.Hash().Hex())
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

	// Add period to settled periods index for optimized queries
	if err := cs.kvStore.AddSettledCashbackPeriod(ctx, period); err != nil {
		log.Printf("⚠️  [CashbackSettlement] Failed to add period to settled index: %v", err)
		// Don't fail the whole settlement if index update fails
	}

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
