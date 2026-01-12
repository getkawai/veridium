package constant

const (
	// Monad Testnet Configuration
	MonadRpcUrl = "https://testnet-rpc.monad.xyz"

	// Contract Addresses (Monad Testnet - Fresh Deployment 2026-01-13)
	KawaiTokenAddress           = "0x80c70F17C4dD23c5C214271293b638197232ab01"
	KawaiEscrowAddress          = "0xd3952Dd84A1Acaf72e3ABD17922e2e850D5E0b58"
	MockUsdtAddress             = "0xDC6f16b4f551638b21C3754F6F93Ea9BbD856298"
	PaymentVaultAddress         = "0x4287dA438FE7FB677D9beB8a7d4A5A09E3C1aC1D"
	KawaiDistributorAddr        = "0x41E2a735aA3D9c9D6cA4dD53b93023ec99FDD7Ef"
	USDTDistributorAddr         = "0x15908521Bd992083F56dCaf7D703b26acFFad742"
	CashbackDistributorAddress  = "0x576564788277b2d8F7475d4B593e08190a2236D6"
	MiningRewardDistributorAddr = "0x86b11B1A7e4e40D181ac06070a0e98648dBc7859"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Fresh deployment: Reset to 0 for clean start
	// - Mainnet: Set to token deployment block to optimize performance
	HolderScanStartBlock = 0
)
