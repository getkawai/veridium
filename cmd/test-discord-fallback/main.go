package main

import (
	"fmt"
	"time"

	"github.com/kawai-network/veridium/pkg/alert"
	"github.com/kawai-network/veridium/pkg/types"
)

func main() {
	fmt.Println("🧪 Testing Discord Fallback for Telegram...")

	// Create Telegram alert with Discord fallback
	telegramAlert := alert.NewTelegramAlert()

	// Test 1: Send Job Reward Log
	fmt.Println("\n1️⃣ Testing Job Reward Log...")
	jobRecord := &types.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: "0x1234567890123456789012345678901234567890",
		UserAddress:        "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		ReferrerAddress:    "0x9876543210987654321098765432109876543210",
		DeveloperAddress:   "0xfedcbafedcbafedcbafedcbafedcbafedcbafed",
		ContributorAmount:  "850000000000000000000", // 850 KAWAI
		DeveloperAmount:    "50000000000000000000",  // 50 KAWAI
		UserAmount:         "50000000000000000000",  // 50 KAWAI
		AffiliatorAmount:   "50000000000000000000",  // 50 KAWAI
		TokenUsage:         1000,
		RewardType:         "kawai",
		HasReferrer:        true,
	}
	telegramAlert.SendJobRewardLog(jobRecord)
	fmt.Println("✅ Job Reward Log sent (check Telegram/Discord)")

	// Test 2: Send Cashback Log
	fmt.Println("\n2️⃣ Testing Cashback Log...")
	cashbackRecord := &types.CashbackRecord{
		UserAddress:    "0x1234567890123456789012345678901234567890",
		TxHash:         "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		DepositAmount:  "1000000000",              // 1000 USDT
		CashbackAmount: "20000000000000000000000", // 20K KAWAI
		RateBPS:        1000,                      // 10%
		Tier:           2,
		IsFirstTime:    true,
		Period:         5,
		Timestamp:      time.Now(),
	}
	telegramAlert.SendCashbackLog(cashbackRecord)
	fmt.Println("✅ Cashback Log sent (check Telegram/Discord)")

	// Test 3: Send Referral Trial Log
	fmt.Println("\n3️⃣ Testing Referral Trial Log...")
	referralRecord := &types.ReferralTrialRecord{
		UserAddress:     "0x1234567890123456789012345678901234567890",
		ReferrerAddress: "0x9876543210987654321098765432109876543210",
		ReferralCode:    "TEST123",
		TrialUSDT:       "10000000",              // 10 USDT
		TrialKAWAI:      "200000000000000000000", // 200 KAWAI
		ReferrerBonus:   "5000000",               // 5 USDT
		MachineID:       "test-machine-id-12345",
		IsReferral:      true,
		Timestamp:       time.Now(),
	}
	telegramAlert.SendReferralTrialLog(referralRecord)
	fmt.Println("✅ Referral Trial Log sent (check Telegram/Discord)")

	// Test 4: Send Alert
	fmt.Println("\n4️⃣ Testing Alert...")
	telegramAlert.SendAlert("SUCCESS", "TestService", "Discord fallback test completed successfully!")
	fmt.Println("✅ Alert sent (check Telegram/Discord)")

	fmt.Println("\n✅ All tests completed! Check your Telegram channel and Discord webhook.")
	fmt.Println("⏳ Waiting 5 seconds for async messages to send...")
	time.Sleep(5 * time.Second)
	fmt.Println("✅ Done!")
}
