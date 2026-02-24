package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kawai-network/x/store"
	"github.com/kawai-network/y/types"
)

func main() {
	contributorAddress := "0xaB48220e6721754b906C30463142Dc0A8F5eBba2"

	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	ctx := context.Background()

	fmt.Println("💉 Injecting mining rewards for:", contributorAddress)
	fmt.Println("")

	// Job 1: 100 KAWAI
	record1 := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: contributorAddress,
		UserAddress:        "0x1111111111111111111111111111111111111111", // Valid address
		ReferrerAddress:    "",
		DeveloperAddress:   contributorAddress,
		ContributorAmount:  "100000000000000000000",
		DeveloperAmount:    "5000000000000000000",
		UserAmount:         "5000000000000000000",
		AffiliatorAmount:   "0",
		TokenUsage:         1000,
		RewardType:         types.RewardTypeMining,
		HasReferrer:        false,
		IsSettled:          false,
	}
	if err := kv.SaveJobReward(ctx, record1); err != nil {
		log.Fatal("Failed to save job 1:", err)
	}
	fmt.Println("✅ Job 1: 100 KAWAI")

	time.Sleep(1 * time.Second)
	record2 := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: contributorAddress,
		UserAddress:        "0x2222222222222222222222222222222222222222", // Valid address
		ReferrerAddress:    "",
		DeveloperAddress:   contributorAddress,
		ContributorAmount:  "150000000000000000000",
		DeveloperAmount:    "5000000000000000000",
		UserAmount:         "5000000000000000000",
		AffiliatorAmount:   "0",
		TokenUsage:         1500,
		RewardType:         types.RewardTypeMining,
		HasReferrer:        false,
		IsSettled:          false,
	}
	if err := kv.SaveJobReward(ctx, record2); err != nil {
		log.Fatal("Failed to save job 2:", err)
	}
	fmt.Println("✅ Job 2: 150 KAWAI")

	time.Sleep(1 * time.Second)
	record3 := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: contributorAddress,
		UserAddress:        "0x3333333333333333333333333333333333333333", // Valid address
		ReferrerAddress:    "",
		DeveloperAddress:   contributorAddress,
		ContributorAmount:  "200000000000000000000",
		DeveloperAmount:    "5000000000000000000",
		UserAmount:         "5000000000000000000",
		AffiliatorAmount:   "0",
		TokenUsage:         2000,
		RewardType:         types.RewardTypeMining,
		HasReferrer:        false,
		IsSettled:          false,
	}
	if err := kv.SaveJobReward(ctx, record3); err != nil {
		log.Fatal("Failed to save job 3:", err)
	}
	fmt.Println("✅ Job 3: 200 KAWAI")

	fmt.Println("")
	fmt.Println("📊 Total: 450 KAWAI")
	fmt.Println("✅ Data injected!")
	fmt.Println("")
	fmt.Println("Next: make settle-mining")
}
