package constant

const (
	// Monad Testnet Configuration
	MonadRpcUrl = "https://testnet-rpc.monad.xyz"

	// Contract Addresses (Monad Testnet - Fresh Deployment 2026-01-12)
	KawaiTokenAddress           = "0xE32660b39D99988Df4bFdc7e4b68A4DC9D654722"
	KawaiEscrowAddress          = "0xd065F9DDb66aa90a1FF62c10868BeF921be2E103"
	MockUsdtAddress             = "0x2cBe796033377352158df11Ab388010ab3097F58"
	PaymentVaultAddress         = "0x9a5A9e31977cB86cD502DC9E0B568d8F17977dAd"
	KawaiDistributorAddr        = "0xaB0DdFbb4bD94d23a32d0C40f9F96d9A61b45463"
	USDTDistributorAddr         = "0x98a7590406a08Cc64dc074D8698B71e4D997a268"
	CashbackDistributorAddress  = "0xdE64f6F5bEe28762c91C76ff762365D553204e35"
	MiningRewardDistributorAddr = "0x8117D77A219EeF5F7869897C3F0973Afb87d8427"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Fresh deployment: Reset to 0 for clean start
	// - Mainnet: Set to token deployment block to optimize performance
	HolderScanStartBlock = 0
)
