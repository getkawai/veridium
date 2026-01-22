package constant

const (
	// Monad Testnet Configuration
	MonadRpcUrl = "https://testnet-rpc.monad.xyz"

	// Contract Addresses (Monad Testnet - Fresh Deployment 2026-01-13)
	KawaiTokenAddress           = "0xcb0FbFEe8253fD4B797dcF7f90e6fDCbE3a34b02"
	OTCMarketAddress            = "0x69d9c1959Ca499369C420A5Aaa7fa7E3b23b1f31"
	StablecoinAddress           = ""
	PaymentVaultAddress         = "0xAaFc7c1f31a53d42B38028B810393926BfD30479"
	RevenueDistributorAddr      = ""
	CashbackDistributorAddress  = "0xCd3882103BB72608A173cb5BA6C4D32e93501Ad0"
	MiningRewardDistributorAddr = "0x531C8Aca995AA92279f36b2F3121ba0004bab3BC"
	ReferralDistributorAddress  = "0x1A4bEc99eC4ADd78F537D0c9f1f9D5b6f110a7dC"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Fresh deployment: Reset to 0 for clean start
	// - Mainnet: Set to token deployment block to optimize performance
	HolderScanStartBlock = 0
)
