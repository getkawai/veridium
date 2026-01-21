package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/alert"
)

func main() {
	fmt.Println("🧪 Testing Claim Failure Alert System (Telegram + Discord)...")
	fmt.Println()

	telegramAlerter := alert.NewTelegramAlert()
	discordAlerter := &alert.DiscordAlert{
		WebhookURL: constant.GetDiscordClaimFailure(),
		Client:     &http.Client{Timeout: 10 * time.Second},
	}

	sendDualAlert := func(level, source, message string) {
		telegramAlerter.SendAlert(level, source, message)
		discordAlerter.SendAlert(level, source, message)
	}

	// Test 1: Cashback claim failure
	fmt.Println("1️⃣ Testing cashback claim failure alert...")
	sendDualAlert("WARNING", "Claim",
		"⚠️ Cashback claim failed\n\nUser: 0x1234...5678\nPeriod: 5\nAmount: 1000 KAWAI\nError: insufficient balance")
	time.Sleep(2 * time.Second)

	// Test 2: Mining claim failure
	fmt.Println("2️⃣ Testing mining claim failure alert...")
	sendDualAlert("WARNING", "Claim",
		"⚠️ Mining claim failed\n\nClaimer: 0xabcd...ef01\nPeriod: 1704067200\nContributor: 5000 KAWAI\nError: RPC connection timeout")
	time.Sleep(2 * time.Second)

	// Test 3: Transaction revert (cashback)
	fmt.Println("3️⃣ Testing cashback claim revert alert...")
	sendDualAlert("ERROR", "Claim",
		"❌ Cashback claim reverted!\n\nUser: 0x1234...5678\nPeriod: 5\nAmount: 1000 KAWAI\nTx: 0xabc123...\nGas Used: 150000")
	time.Sleep(2 * time.Second)

	// Test 4: Transaction revert (mining)
	fmt.Println("4️⃣ Testing mining claim revert alert...")
	sendDualAlert("ERROR", "Claim",
		"❌ Mining claim reverted!\n\nClaimer: 0xabcd...ef01\nPeriod: 1704067200\nContributor: 5000 KAWAI\nTx: 0xdef456...\nGas Used: 250000")
	time.Sleep(2 * time.Second)

	// Test 5: KAWAI claim failure
	fmt.Println("5️⃣ Testing KAWAI claim failure alert...")
	sendDualAlert("WARNING", "Claim",
		"⚠️ KAWAI claim failed\n\nClaimer: 0x9876...5432\nPeriod: 10\nIndex: 42\nAmount: 2500000000000000000000\nError: merkle proof verification failed")
	time.Sleep(2 * time.Second)

	// Test 6: USDT claim revert
	fmt.Println("6️⃣ Testing USDT claim revert alert...")
	sendDualAlert("ERROR", "Claim",
		"❌ USDT claim reverted!\n\nClaimer: 0x9876...5432\nPeriod: 10\nIndex: 42\nAmount: 100000000\nTx: 0x789abc...\nGas Used: 180000")
	time.Sleep(2 * time.Second)

	fmt.Println()
	fmt.Println("✅ All test alerts sent!")
	fmt.Println("📱 Check your Telegram for 6 alert messages")
	fmt.Println("💬 Check your Discord for 6 alert messages")
	fmt.Println()
	fmt.Println("Expected alerts:")
	fmt.Println("  - 4 WARNING alerts (claim failures)")
	fmt.Println("  - 2 ERROR alerts (transaction reverts)")
	fmt.Println()
	fmt.Println("Channels:")
	fmt.Println("  - Telegram: General alerts")
	fmt.Println("  - Discord: Claim failure webhook")
}
