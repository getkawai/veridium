package constant

const (
	// Monad Testnet Configuration
	MonadRpcUrl = "https://testnet-rpc.monad.xyz"

	// Contract Addresses (Monad Testnet - Fresh Deployment 2026-01-13)
	KawaiTokenAddress           = "0xcBFdd2F2B8a4174711021A0cA8E00426d0c14d0B"
	OTCMarketAddress            = "0xf4CE5EEBdAa162756A58BE495dda968e374fbcc1"
	UsdtTokenAddress            = "0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc"
	PaymentVaultAddress         = "0xD353745dE21B218228C84f26f6c8862CC27F817C"
	KawaiDistributorAddr        = "0xcb39CA020ba056Df8cc746Fe08fbb06A0AC58a50"
	USDTDistributorAddr         = "0x026bC3673D9B5DA393Ad6149ECD68b2A5df1E811"
	CashbackDistributorAddress  = "0x93D6dD33e1e9e735FEF13416a83153964C914581"
	MiningRewardDistributorAddr = "0xA65d51efc7FdcAD26AA5cA5Ca93B760c7dD48de6"
	ReferralDistributorAddress  = "0x3c8b41c1075c26d66Eb7d90Ad218449e169079e9"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Fresh deployment: Reset to 0 for clean start
	// - Mainnet: Set to token deployment block to optimize performance
	HolderScanStartBlock = 0
)
