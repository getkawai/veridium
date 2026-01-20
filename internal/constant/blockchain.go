package constant

const (
	// Monad Testnet Configuration
	MonadRpcUrl = "https://testnet-rpc.monad.xyz"

	// Contract Addresses (Monad Testnet - Fresh Deployment 2026-01-13)
	KawaiTokenAddress = "0xf68910e8d19047A309f989FFB515E44FBca5D31A"
	OTCMarketAddress  = "0x3E597D76B40004c3fC517C404037fD6F16C8fc34"

	// UsdtTokenAddress: Stablecoin address for deposits and payments
	// - Testnet: MockUSDT (0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc)
	// - Mainnet: USDC (0x754704bc059f8c67012fed69bc8a327a5aafb603)
	// Variable name kept for backward compatibility, but represents any stablecoin
	UsdtTokenAddress = "0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc"

	PaymentVaultAddress         = "0xDA94C8ac2a61eafBd47853EE22702BDCd45B6d93"
	KawaiDistributorAddr        = "0x2B11e8385A859Ea75C77E05Bc0D9756A87017E92"
	USDTDistributorAddr         = "0x896fB97f81ECBEfDBe29DCc3550aC984704932bF"
	CashbackDistributorAddress  = "0x3d5Bfe788782A90ac124096296B45eaFFc43C79B"
	MiningRewardDistributorAddr = "0x1f78c7c472205F1720aAb66a565981561b5EBac0"
	ReferralDistributorAddress  = "0x1c218602218745B20CE201948CaE836f8E94E111"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Fresh deployment: Reset to 0 for clean start
	// - Mainnet: Set to token deployment block to optimize performance
	HolderScanStartBlock = 0
)
