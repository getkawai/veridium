package store

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleWeeklySettlement demonstrates a complete weekly settlement cycle with all safety features
func ExampleWeeklySettlement() {
	ctx := context.Background()
	
	// Initialize store
	store, err := NewKVStore("api-token", "account-id", "namespace-id")
	if err != nil {
		log.Fatal(err)
	}
	
	// ===== WEEK 1: Accumulation Phase =====
	log.Println("=== Week 1: Jobs Execution & Reward Accumulation ===")
	
	// Jobs are executed throughout the week, rewards accumulate
	// (This happens automatically via RecordJobReward calls)
	
	// ===== END OF WEEK 1: Settlement Phase =====
	log.Println("\n=== End of Week 1: Settlement ===")
	
	// Step 1: Get snapshots of all balances (SORTED for consistent Merkle tree)
	snapshots, err := store.GetSettlementSnapshots(ctx, "kawai")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Snapshotted %d contributors (sorted by address)", len(snapshots))
	for _, s := range snapshots {
		log.Printf("  - %s: %s KAWAI", s.Address, s.Amount)
	}
	
	// Step 2: Generate Merkle tree (external library - not implemented here)
	// IMPORTANT: Use the same sorted order as GetSettlementSnapshots!
	periodID := GenerateUniquePeriodID() // Use unique ID to prevent collision
	merkleRoot := "0x123abc..." // Generated from tree
	
	// Create proofs map (would come from Merkle tree generation)
	proofs := map[string]*MerkleProofData{
		"0xContributor1": {
			Index:  0,
			Amount: "1000000000000000000000", // 1000 KAWAI
			Proof:  []string{"0xabc...", "0xdef..."},
		},
		"0xContributor2": {
			Index:  1,
			Amount: "500000000000000000000", // 500 KAWAI
			Proof:  []string{"0x123...", "0x456..."},
		},
	}
	
	// Step 3: Perform settlement with rollback support
	period, err := store.PerformSettlement(ctx, periodID, merkleRoot, "kawai", proofs)
	if err != nil {
		log.Fatalf("Settlement failed (rolled back): %v", err)
	}
	
	log.Printf("Settlement completed:")
	log.Printf("  - Period ID: %d", period.PeriodID)
	log.Printf("  - Status: %s", period.Status)
	log.Printf("  - Proofs saved: %d", period.ProofsSaved)
	log.Printf("  - Balances reset: %d", period.BalancesReset)
	log.Printf("  - Total distributed: %s KAWAI", period.TotalAmount)
}

// ExampleClaimWithStatusTracking demonstrates safe claim flow with status tracking
func ExampleClaimWithStatusTracking() {
	ctx := context.Background()
	
	store, err := NewKVStore("api-token", "account-id", "namespace-id")
	if err != nil {
		log.Fatal(err)
	}
	
	address := "0xContributor1"
	periodID := int64(1704067200)
	
	// Step 1: Get claimable rewards
	claimable, err := store.GetClaimableRewards(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	
	unclaimedProofs := claimable["unclaimed_proofs"].([]*MerkleProofData)
	pendingProofs := claimable["pending_proofs"].([]*MerkleProofData)
	
	log.Printf("Contributor %s:", address)
	log.Printf("  - Unclaimed proofs: %d", len(unclaimedProofs))
	log.Printf("  - Pending proofs: %d", len(pendingProofs))
	log.Printf("  - Total KAWAI claimable: %s", claimable["total_kawai_claimable"])
	log.Printf("  - Currently accumulating: %s", claimable["current_kawai_accumulating"])
	
	// Step 2: Get proof for claiming
	proof, err := store.GetMerkleProofForPeriod(ctx, address, periodID)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("\nPreparing to claim period %d:", periodID)
	log.Printf("  - Amount: %s", proof.Amount)
	log.Printf("  - Merkle Root: %s", proof.MerkleRoot)
	
	// Step 3: Submit transaction to smart contract
	txHash := "0xabc123..." // From web3 transaction
	
	// IMPORTANT: Mark as pending BEFORE waiting for confirmation
	err = store.MarkClaimPending(ctx, address, periodID, txHash)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Marked claim as pending, tx: %s", txHash)
	
	// Step 4: Wait for transaction confirmation (simulated)
	log.Println("Waiting for transaction confirmation...")
	txSuccess := true // Simulated - would check blockchain
	
	if txSuccess {
		// Step 5a: Transaction confirmed - mark as claimed
		err = store.ConfirmClaim(ctx, address, periodID)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Claim confirmed successfully!")
	} else {
		// Step 5b: Transaction failed - mark as failed for retry
		err = store.MarkClaimFailed(ctx, address, periodID, "Transaction reverted")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Claim failed, can retry later")
		
		// Later: Retry the claim
		err = store.RetryFailedClaim(ctx, address, periodID)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Claim reset for retry")
	}
}

// ExampleLargeScaleSettlement demonstrates settlement for 1000+ contributors
func ExampleLargeScaleSettlement() {
	ctx := context.Background()
	
	store, err := NewKVStore("api-token", "account-id", "namespace-id")
	if err != nil {
		log.Fatal(err)
	}
	
	// For large settlements, use custom config
	config := &SettlementConfig{
		BatchSize:      100,               // Process 100 at a time
		BatchDelay:     200 * time.Millisecond, // Wait 200ms between batches
		MaxRetries:     5,                 // Retry failed ops 5 times
		EnableRollback: true,              // Enable rollback on failure
	}
	
	periodID := GenerateUniquePeriodID()
	merkleRoot := "0x123abc..."
	
	// Large proofs map (1000+ entries)
	proofs := make(map[string]*MerkleProofData)
	// ... populate proofs
	
	// Option 1: Sequential with batching
	period, err := store.PerformSettlementWithConfig(ctx, periodID, merkleRoot, "kawai", proofs, config)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Sequential settlement completed in %v", period.CompletedAt.Sub(period.StartedAt))
	
	// Option 2: Parallel processing (faster but more resource intensive)
	periodID2 := GenerateUniquePeriodID()
	period2, err := store.PerformSettlementParallel(ctx, periodID2, merkleRoot, "kawai", proofs, 10) // 10 workers
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Parallel settlement completed in %v", period2.CompletedAt.Sub(period2.StartedAt))
}

// ExampleResumeInterruptedSettlement demonstrates resuming a failed settlement
func ExampleResumeInterruptedSettlement() {
	ctx := context.Background()
	
	store, err := NewKVStore("api-token", "account-id", "namespace-id")
	if err != nil {
		log.Fatal(err)
	}
	
	// Assume settlement was interrupted
	periodID := int64(1704067200)
	
	// Check settlement status
	period, err := store.GetSettlementPeriod(ctx, periodID)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Settlement %d status: %s", periodID, period.Status)
	
	if period.Status != SettlementStatusCompleted {
		// Resume from where it left off
		proofs := make(map[string]*MerkleProofData)
		// ... reload proofs if needed
		
		config := DefaultSettlementConfig()
		resumedPeriod, err := store.ResumeSettlement(ctx, periodID, proofs, config)
		if err != nil {
			log.Fatal(err)
		}
		
		log.Printf("Settlement resumed and completed: %s", resumedPeriod.Status)
	}
}

// ExampleSoftDeleteContributor demonstrates soft delete flow
func ExampleSoftDeleteContributor() {
	ctx := context.Background()
	
	store, err := NewKVStore("api-token", "account-id", "namespace-id")
	if err != nil {
		log.Fatal(err)
	}
	
	address := "0xContributor1"
	
	// Register contributor
	contributor, err := store.RegisterContributor(ctx, address, "https://node.example.com", "RTX 4090")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Registered: %s, IsActive: %v", contributor.WalletAddress, contributor.IsActive)
	
	// Contributor wants to leave (but has pending rewards)
	err = store.SoftDeleteContributor(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Soft deleted: %s", address)
	
	// Check - contributor still exists but inactive
	deleted, err := store.GetContributor(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("After delete - IsActive: %v, DeletedAt: %v", deleted.IsActive, deleted.DeletedAt)
	
	// Settlement still includes inactive contributors with balance
	snapshots, err := store.GetSettlementSnapshots(ctx, "kawai")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Snapshots include %d contributors (including inactive with balance)", len(snapshots))
	
	// Contributor can still claim rewards
	claimable, err := store.GetClaimableRewards(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Claimable rewards: %s KAWAI", claimable["total_kawai_claimable"])
	
	// Later: Contributor returns
	err = store.RestoreContributor(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Contributor restored!")
}

// ExamplePendingClaimsMonitoring demonstrates monitoring pending claims
func ExamplePendingClaimsMonitoring() {
	ctx := context.Background()
	
	store, err := NewKVStore("api-token", "account-id", "namespace-id")
	if err != nil {
		log.Fatal(err)
	}
	
	// Get all pending claims (useful for admin monitoring)
	pendingClaims, err := store.GetPendingClaims(ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Found %d pending claims:", len(pendingClaims))
	for _, claim := range pendingClaims {
		log.Printf("  - Period %d: %s, TX: %s, Attempts: %d",
			claim.PeriodID,
			claim.Amount,
			claim.ClaimTxHash,
			claim.ClaimAttempts)
	}
	
	// Check each pending claim on blockchain
	for _, claim := range pendingClaims {
		// Check if TX is confirmed on blockchain
		txConfirmed := checkBlockchainTx(claim.ClaimTxHash)
		
		if txConfirmed {
			// Extract address from claim (would need to store this)
			address := "0x..." // Get from claim data
			store.ConfirmClaim(ctx, address, claim.PeriodID)
		} else {
			// TX still pending or failed
			// Could implement timeout logic here
		}
	}
}

// Helper function (simulated)
func checkBlockchainTx(txHash string) bool {
	// Would check blockchain for TX confirmation
	return true
}

// ExampleCompleteFlow demonstrates the entire flow from registration to claim
func ExampleCompleteFlow() {
	fmt.Println(`
Complete Flow with All Safety Features:
=======================================

1. REGISTRATION
   RegisterContributor(ctx, address, endpoint, specs)
   → Creates contributor with IsActive=true

2. ACCUMULATION (Week 1)
   RecordJobReward(ctx, contributor, tokens, admin, mode)
   → AccumulatedRewards increases
   → EnsureAdminExists() auto-registers admin if needed

3. SETTLEMENT (End of Week 1)
   a. periodID := GenerateUniquePeriodID()  // Unique ID prevents collision
   b. snapshots := GetSettlementSnapshots(ctx, "kawai")  // SORTED by address
   c. merkleRoot, proofs := GenerateMerkleTree(snapshots)  // Use same order!
   d. PerformSettlement(ctx, periodID, merkleRoot, "kawai", proofs)
      → Checks for existing period (collision prevention)
      → Saves all proofs FIRST (with rollback on failure)
      → Resets balances AFTER all proofs saved
      → Saves settlement metadata with status tracking

4. CONTRIBUTOR UNREGISTER (Optional)
   SoftDeleteContributor(ctx, address)
   → IsActive=false, but data preserved
   → Still included in settlements (if has balance)
   → Can still claim rewards

5. CLAIM FLOW
   a. GetClaimableRewards(ctx, address)  // Shows unclaimed + pending
   b. GetMerkleProofForPeriod(ctx, address, periodID)
   c. MarkClaimPending(ctx, address, periodID, txHash)  // BEFORE sending TX
   d. Submit TX to blockchain
   e. Wait for confirmation
   f. ConfirmClaim(ctx, address, periodID)  // On success
      OR MarkClaimFailed(ctx, address, periodID, reason)  // On failure
   g. RetryFailedClaim(ctx, address, periodID)  // To retry later

6. MONITORING
   GetPendingClaims(ctx)  // Admin can monitor stuck claims
   ListSettlementPeriods(ctx)  // View settlement history

7. MAINTENANCE
   CleanupOldProofs(ctx, 90*24*time.Hour)  // Clean old confirmed claims

Safety Features Implemented:
- ✅ Unique Period ID (prevents collision)
- ✅ Sorted snapshots (consistent Merkle tree)
- ✅ Rollback on settlement failure
- ✅ Batch processing (handles 1000+ contributors)
- ✅ Claim status tracking (pending/confirmed/failed)
- ✅ Soft delete (preserves data for claims)
- ✅ Resumable settlement (handles interruptions)
- ✅ Auto-register admin
`)
}
