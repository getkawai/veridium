package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/veridium/pkg/merkle"
	"github.com/kawai-network/veridium/pkg/store"
)

// ReferralSettlement handles referral commission settlement
type ReferralSettlement struct {
	kvStore *store.KVStore
}

// ReferralCommissionRecord represents a referrer's accumulated commission for a period
type ReferralCommissionRecord struct {
	ReferrerAddress string `json:"referrer_address"`
	Period          uint64 `json:"period"`
	TotalKawai      string `json:"total_kawai"` // Total commission earned (wei)
	JobCount        int    `json:"job_count"`   // Number of jobs that generated commission
}

// NewReferralSettlement creates a new referral settlement instance
func NewReferralSettlement(kvStore *store.KVStore, privateKey string) *ReferralSettlement {
	return &ReferralSettlement{
		kvStore: kvStore,
	}
}

// SettleReferral generates Merkle tree for referral commissions and stores proofs
// Period is weekly, same as mining settlement
func (rs *ReferralSettlement) SettleReferral(ctx context.Context, period uint64) error {
	log.Printf("🤝 Starting referral commission settlement for period %d", period)
	log.Println("")

	// 1. Collect all referral commissions for this period
	commissions, err := rs.collectReferralCommissions(ctx, period)
	if err != nil {
		return fmt.Errorf("failed to collect referral commissions: %w", err)
	}

	if len(commissions) == 0 {
		log.Println("⚠️  No referral commissions found for this period")
		return nil
	}

	log.Printf("Found %d referrers with commissions", len(commissions))
	log.Println("")

	// 2. Generate Merkle tree (3-field: period, account, amount)
	leaves := make([][]byte, len(commissions))
	referrers := make([]string, len(commissions))

	for i, comm := range commissions {
		// Create leaf: keccak256(period, account, amount)
		leaf := rs.createLeaf(period, comm.ReferrerAddress, comm.TotalKawai)
		leaves[i] = leaf
		referrers[i] = comm.ReferrerAddress

		log.Printf("  %d. %s: %s KAWAI (%d jobs)",
			i+1, comm.ReferrerAddress, formatKawai(comm.TotalKawai), comm.JobCount)
	}

	log.Println("")
	log.Println("🌳 Generating Merkle tree...")

	tree := merkle.NewMerkleTree(leaves)

	merkleRoot := common.BytesToHash(tree.Root).Hex()
	log.Printf("✅ Merkle Root: %s", merkleRoot)
	log.Println("")

	// 3. Store proofs in KV
	log.Println("💾 Storing Merkle proofs...")
	for i, comm := range commissions {
		proof, found := tree.GetProof(leaves[i])
		if !found {
			return fmt.Errorf("proof not found for %s", comm.ReferrerAddress)
		}

		proofHex := make([]string, len(proof))
		for j, p := range proof {
			proofHex[j] = common.BytesToHash(p).Hex()
		}

		// Store proof
		proofData := map[string]interface{}{
			"proof":  proofHex,
			"amount": comm.TotalKawai,
		}

		data, err := json.Marshal(proofData)
		if err != nil {
			return fmt.Errorf("failed to marshal proof for %s: %w", comm.ReferrerAddress, err)
		}

		// Store in KV: referral_proof:period:address
		proofKey := fmt.Sprintf("referral_proof:%d:%s", period, comm.ReferrerAddress)
		if err := rs.kvStore.StoreMarketplaceData(ctx, proofKey, data); err != nil {
			return fmt.Errorf("failed to store proof for %s: %w", comm.ReferrerAddress, err)
		}

		log.Printf("  ✅ Stored proof for %s", comm.ReferrerAddress)
	}

	log.Println("")

	// 4. Store Merkle root
	log.Println("💾 Storing Merkle root...")
	rootKey := fmt.Sprintf("referral_period:%d:merkle_root", period)
	rootData, err := json.Marshal(merkleRoot)
	if err != nil {
		return fmt.Errorf("failed to marshal merkle root: %w", err)
	}

	if err := rs.kvStore.StoreMarketplaceData(ctx, rootKey, rootData); err != nil {
		return fmt.Errorf("failed to store merkle root: %w", err)
	}

	log.Println("✅ Merkle root stored")
	log.Println("")

	// Summary
	totalCommission := big.NewInt(0)
	for _, comm := range commissions {
		amount := new(big.Int)
		amount.SetString(comm.TotalKawai, 10)
		totalCommission.Add(totalCommission, amount)
	}

	log.Println("📊 Settlement Summary:")
	log.Printf("   Period:         %d", period)
	log.Printf("   Referrers:      %d", len(commissions))
	log.Printf("   Total Commission: %s KAWAI", formatKawai(totalCommission.String()))
	log.Printf("   Merkle Root:    %s", merkleRoot)
	log.Println("")
	log.Println("✅ Referral settlement completed!")
	log.Println("")
	log.Println("📝 Next: Upload Merkle root to ReferralRewardDistributor contract")

	return nil
}

// collectReferralCommissions collects all referral commissions for a period
// by scanning job rewards and aggregating by referrer
func (rs *ReferralSettlement) collectReferralCommissions(ctx context.Context, period uint64) ([]ReferralCommissionRecord, error) {
	// Get all unsettled job rewards (this includes referrer info)
	jobRewards, err := rs.kvStore.GetAllUnsettledJobRewards(ctx, "kawai")
	if err != nil {
		return nil, fmt.Errorf("failed to get job rewards: %w", err)
	}

	// Aggregate by referrer
	commissionMap := make(map[string]*ReferralCommissionRecord)

	for _, jobs := range jobRewards {
		for _, job := range jobs {
			// Skip if no referrer
			if job.ReferrerAddress == "" {
				continue
			}

			// Get or create commission record
			if _, exists := commissionMap[job.ReferrerAddress]; !exists {
				commissionMap[job.ReferrerAddress] = &ReferralCommissionRecord{
					ReferrerAddress: job.ReferrerAddress,
					Period:          period,
					TotalKawai:      "0",
					JobCount:        0,
				}
			}

			// Add affiliator amount to total
			current := new(big.Int)
			current.SetString(commissionMap[job.ReferrerAddress].TotalKawai, 10)

			affiliatorAmount := new(big.Int)
			affiliatorAmount.SetString(job.AffiliatorAmount, 10)

			total := new(big.Int).Add(current, affiliatorAmount)
			commissionMap[job.ReferrerAddress].TotalKawai = total.String()
			commissionMap[job.ReferrerAddress].JobCount++
		}
	}

	// Convert map to slice
	commissions := make([]ReferralCommissionRecord, 0, len(commissionMap))
	for _, comm := range commissionMap {
		// Only include if commission > 0
		amount := new(big.Int)
		amount.SetString(comm.TotalKawai, 10)
		if amount.Cmp(big.NewInt(0)) > 0 {
			commissions = append(commissions, *comm)
		}
	}

	return commissions, nil
}

// createLeaf creates a Merkle leaf for referral commission
// Leaf = keccak256(period, account, amount)
func (rs *ReferralSettlement) createLeaf(period uint64, account string, amount string) []byte {
	// Convert period to bytes32
	periodBytes := make([]byte, 32)
	periodBig := new(big.Int).SetUint64(period)
	periodBig.FillBytes(periodBytes)

	// Convert address to bytes20
	accountBytes := common.HexToAddress(account).Bytes()

	// Convert amount to bytes32
	amountBytes := make([]byte, 32)
	amountBig := new(big.Int)
	amountBig.SetString(amount, 10)
	amountBig.FillBytes(amountBytes)

	// Concatenate and hash
	data := append(periodBytes, accountBytes...)
	data = append(data, amountBytes...)

	return crypto.Keccak256(data)
}

// formatKawai formats wei amount to KAWAI (divide by 1e18)
func formatKawai(weiStr string) string {
	wei := new(big.Int)
	wei.SetString(weiStr, 10)

	// Divide by 1e18
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	kawai := new(big.Int).Div(wei, divisor)

	return kawai.String()
}

