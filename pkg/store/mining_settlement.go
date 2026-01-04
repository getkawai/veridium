package store

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/veridium/pkg/merkle"
)

// GenerateMiningSettlement creates a Merkle tree for mining rewards
// with referral-based splits (85/5/5/5 or 90/5/5)
// This generates 9-field Merkle leaves for MiningRewardDistributor contract
func (s *KVStore) GenerateMiningSettlement(ctx context.Context, rewardType string) (*SettlementPeriod, error) {
	slog.Info("Generating mining settlement", "reward_type", rewardType)

	// 1. Get all unsettled job rewards grouped by contributor
	jobRewardsByContributor, err := s.GetAllUnsettledJobRewards(ctx, rewardType)
	if err != nil {
		return nil, fmt.Errorf("failed to get unsettled job rewards: %w", err)
	}

	if len(jobRewardsByContributor) == 0 {
		return nil, fmt.Errorf("no unsettled job rewards found")
	}

	slog.Info("Found unsettled job rewards", "contributors", len(jobRewardsByContributor))

	// 2. Aggregate rewards per contributor and generate Merkle leaves
	var leaves [][]byte
	proofs := make(map[string]*MerkleProofData)
	period := uint64(time.Now().Unix())
	
	totalAmount := big.NewInt(0)
	currentIndex := uint64(0)

	for contributorAddr, jobs := range jobRewardsByContributor {
		if len(jobs) == 0 {
			continue
		}

		// Aggregate amounts from all jobs
		totalContrib := big.NewInt(0)
		totalDev := big.NewInt(0)
		totalUser := big.NewInt(0)
		totalAff := big.NewInt(0)

		// Use addresses from the most recent job
		// (In practice, developer address varies but user/affiliator should be consistent)
		latestJob := jobs[len(jobs)-1]
		developerAddr := latestJob.DeveloperAddress
		userAddr := latestJob.UserAddress
		affiliatorAddr := latestJob.ReferrerAddress

		for _, job := range jobs {
			amt := new(big.Int)
			amt.SetString(job.ContributorAmount, 10)
			totalContrib.Add(totalContrib, amt)

			amt = new(big.Int)
			amt.SetString(job.DeveloperAmount, 10)
			totalDev.Add(totalDev, amt)

			amt = new(big.Int)
			amt.SetString(job.UserAmount, 10)
			totalUser.Add(totalUser, amt)

			amt = new(big.Int)
			amt.SetString(job.AffiliatorAmount, 10)
			totalAff.Add(totalAff, amt)
		}

		// Generate 9-field Merkle leaf
		leaf := generateMiningMerkleLeaf(
			period,
			common.HexToAddress(contributorAddr),
			totalContrib, totalDev, totalUser, totalAff,
			common.HexToAddress(developerAddr),
			common.HexToAddress(userAddr),
			common.HexToAddress(affiliatorAddr),
		)
		leaves = append(leaves, leaf)

		// Store proof data with mining-specific fields
		proofs[contributorAddr] = &MerkleProofData{
			Index:              currentIndex,
			Amount:             totalContrib.String(), // For backward compatibility
			PeriodID:           int64(period),
			RewardType:         rewardType,
			ContributorAmount:  totalContrib.String(),
			DeveloperAmount:    totalDev.String(),
			UserAmount:         totalUser.String(),
			AffiliatorAmount:   totalAff.String(),
			DeveloperAddress:   developerAddr,
			UserAddress:        userAddr,
			AffiliatorAddress:  affiliatorAddr,
			ClaimStatus:        ClaimStatusUnclaimed,
			CreatedAt:          time.Now(),
		}

		// Add to total
		totalAmount.Add(totalAmount, totalContrib)
		currentIndex++
	}

	if len(leaves) == 0 {
		return nil, fmt.Errorf("no valid leaves generated")
	}

	slog.Info("Generated Merkle leaves", "count", len(leaves))

	// 3. Build Merkle tree
	tree := merkle.NewMerkleTree(leaves)
	root := tree.Root

	slog.Info("Merkle tree built", "root", fmt.Sprintf("0x%x", root))

	// 4. Generate proofs for each leaf
	for i, leaf := range leaves {
		proof, ok := tree.GetProof(leaf)
		if !ok {
			slog.Warn("Failed to generate proof", "index", i)
			continue
		}

		var proofHex []string
		for _, p := range proof {
			proofHex = append(proofHex, fmt.Sprintf("0x%x", p))
		}

		// Find corresponding contributor address
		for contributorAddr, proofData := range proofs {
			if proofData.Index == uint64(i) {
				proofData.Proof = proofHex
				proofData.MerkleRoot = fmt.Sprintf("0x%x", root)
				slog.Debug("Generated proof", "contributor", contributorAddr, "index", i)
				break
			}
		}
	}

	// 5. Save settlement period using existing parallel settlement logic
	settlement, err := s.PerformSettlementParallel(
		ctx,
		int64(period),
		fmt.Sprintf("0x%x", root),
		rewardType,
		proofs,
		10, // 10 workers
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save settlement: %w", err)
	}

	// 6. Mark all job rewards as settled
	for contributorAddr := range jobRewardsByContributor {
		if err := s.MarkJobRewardsAsSettled(ctx, contributorAddr, int64(period)); err != nil {
			slog.Warn("Failed to mark job rewards as settled", 
				"contributor", contributorAddr, "error", err)
		}
	}

	slog.Info("Mining settlement completed", 
		"period", period,
		"contributors", len(proofs),
		"total_amount", totalAmount.String())

	return settlement, nil
}

// generateMiningMerkleLeaf creates a 9-field Merkle leaf for MiningRewardDistributor
// Matches the Solidity keccak256(abi.encodePacked(...)) format
func generateMiningMerkleLeaf(
	period uint64,
	contributor common.Address,
	contributorAmt, developerAmt, userAmt, affiliatorAmt *big.Int,
	developer, user, affiliator common.Address,
) []byte {
	// Solidity abi.encodePacked packs values tightly without padding
	// For uint256, it uses 32 bytes; for address, it uses 20 bytes
	return crypto.Keccak256(
		common.LeftPadBytes(big.NewInt(int64(period)).Bytes(), 32),  // uint256
		contributor.Bytes(),                                          // address (20 bytes)
		common.LeftPadBytes(contributorAmt.Bytes(), 32),             // uint256
		common.LeftPadBytes(developerAmt.Bytes(), 32),               // uint256
		common.LeftPadBytes(userAmt.Bytes(), 32),                    // uint256
		common.LeftPadBytes(affiliatorAmt.Bytes(), 32),              // uint256
		developer.Bytes(),                                            // address (20 bytes)
		user.Bytes(),                                                 // address (20 bytes)
		affiliator.Bytes(),                                           // address (20 bytes)
	)
}

