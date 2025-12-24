package admin

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/store"
)

type AdminManager struct {
	Chain *blockchain.Client
	Store store.Store
}

func NewAdminManager(chain *blockchain.Client, kv store.Store) *AdminManager {
	return &AdminManager{
		Chain: chain,
		Store: kv,
	}
}

// AuditWorkers prints the status of all registered workers
// Note: This requires the KV store to support listing keys,
// but for now, we'll implement a simple version that logs status.
func (a *AdminManager) AuditWorkers(ctx context.Context) error {
	log.Println("--- Worker Audit Report ---")

	workers, err := a.Store.ListWorkers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list workers: %w", err)
	}

	if len(workers) == 0 {
		log.Println("No workers registered.")
		return nil
	}

	fmt.Printf("%-42s | %-15s | %-10s | %s\n", "Wallet Address", "Last Seen", "Status", "Specs")
	fmt.Println("------------------------------------------------------------------------------------------------")
	for _, w := range workers {
		lastSeen := w.LastSeen.Format("2006-01-02 15:04")
		fmt.Printf("%-42s | %-15s | %-10s | %s\n", w.WalletAddress, lastSeen, w.Status, w.HardwareSpecs)
	}

	return nil
}

// CalculateDividends (Placeholder)
func (a *AdminManager) CalculateDividends(ctx context.Context) error {
	log.Println("--- Dividend Calculation ---")

	if a.Chain == nil {
		return fmt.Errorf("blockchain client not initialized")
	}

	// 1. Get USDT Balance from PaymentVault (Placeholder address for now)
	vaultAddr := "0x5ba96c283530acfe7e6051d9e20348df" // zoneID from user as proxy? (usually vault)
	log.Printf("Fetching balance from PaymentVault: %s", vaultAddr)

	// TODO: Use a.Chain to call balanceOf on USDT contract for vaultAddr
	balance := big.NewInt(1000000000) // Dummy 1000 USDT (assuming 6 decimals)
	log.Printf("Total Dividends Available: %f USDT", float64(balance.Int64())/1e6)

	// 2. Fetch Holders via Event Scanning
	log.Println("Analyzing Token Holders via Transfer Event Scanning...")

	// Holders map: wallet -> balance
	holders := make(map[string]*big.Int)

	// TODO: Use a.Chain.Client.FilterLogs to get Transfer(address indexed from, address indexed to, uint256 value)
	// Example logic:
	// for _, log := range logs {
	//    from := common.HexToAddress(log.Topics[1].Hex()).Hex()
	//    to := common.HexToAddress(log.Topics[2].Hex()).Hex()
	//    val := new(big.Int).SetBytes(log.Data)
	//    holders[from].Sub(holders[from], val)
	//    holders[to].Add(holders[to], val)
	// }

	// Dummy data for demonstration
	holders["0x123..."] = big.NewInt(5000)
	holders["0xabc..."] = big.NewInt(25000)

	log.Println("Holder | Share (%) | Dividend (USDT)")
	log.Println("---------------------------------------")

	totalSupply := big.NewInt(1000000) // Dummy
	for addr, amt := range holders {
		share := new(big.Float).Quo(new(big.Float).SetInt(amt), new(big.Float).SetInt(totalSupply))
		dividend := new(big.Float).Mul(share, new(big.Float).SetInt(balance))

		// Convert dividend (6 decimals) to float for display
		divFloat, _ := dividend.Float64()
		shareFloat, _ := share.Mul(share, big.NewFloat(100)).Float64()

		fmt.Printf("%s | %.2f%% | %.4f\n", addr, shareFloat, divFloat/1e6)
	}

	return nil
}
