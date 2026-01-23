package constant

const (
	// Monad Mainnet Configuration
	MonadRpcUrl = "https://rpc.monad.xyz"

	// Contract Addresses (Monad Mainnet - Deployment 2026-01-23)
	KawaiTokenAddress           = "0x9cbdb316b31fd2efa469c57dcf57be0af630f64c"
	OTCMarketAddress            = "0xe7d73b901b7202b4f686166420ee76cfe860d28d"
	StablecoinAddress           = "0x754704bc059f8c67012fed69bc8a327a5aafb603"
	PaymentVaultAddress         = "0xffdf0fb715bec64db41307c26abf545295d31e44"
	RevenueDistributorAddr      = "0x7454495f1a7e2854e4215a4d797e0abd7e14bbe4"
	CashbackDistributorAddress  = "0x1feff071f37a5cb8833e227d8dddea43aa374449"
	MiningRewardDistributorAddr = "0xc58d3f5d04e5748fc1806980e26c1eb487045442"
	ReferralDistributorAddress  = "0xfbbe8b96d1b5eff919ce09da28737c667faa7957"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Fresh deployment: Reset to 0 for clean start
	// - Mainnet: Set to token deployment block to optimize performance
	HolderScanStartBlock = 0
)
