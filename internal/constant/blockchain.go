package constant

const (
	// Monad Testnet Configuration
	MonadRpcUrl = "https://testnet-rpc.monad.xyz"

	// Contract Addresses (Monad Testnet - Fresh Deployment 2026-01-13)
	KawaiTokenAddress           = "0xeBB1bf5Db06e9E858FfA346520e22D25944fbB87"
	KawaiEscrowAddress          = "0xA71306429A0232C8b9EF5F2908e04Fcad7E96B8E"
	MockUsdtAddress             = "0x60B99499aA369b63BE1aC7123E383Ef450227a63"
	PaymentVaultAddress         = "0x44581F99a39eEab898194e54c3B89bb261d3b2c3"
	KawaiDistributorAddr        = "0x578C6E95908399467f98bDf10378beA95b44EF23"
	USDTDistributorAddr         = "0xe024c9a670E0039A887c7f385892333CAb04a275"
	CashbackDistributorAddress  = "0x56Bc3045088C51f329F86AE5Dec3faED59d77664"
	MiningRewardDistributorAddr = "0xFEC16f47BD9DD4B9E05DAaC7BBef8C047f010289"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Fresh deployment: Reset to 0 for clean start
	// - Mainnet: Set to token deployment block to optimize performance
	HolderScanStartBlock = 0
)
