package constant

const (
	// Monad Testnet Configuration
	MonadRpcUrl = "https://testnet-rpc.monad.xyz"

	// Contract Addresses (Monad Testnet - Deployed 2025-12-31)
	KawaiTokenAddress          = "0xF27c5c43a746B329B1c767CE1b319c9EBfE8012E"
	KawaiEscrowAddress         = "0x5b1235038B2F05aC88b791A23814130710eFaaEa"
	MockUsdtAddress            = "0xb8cD3f468E9299Fa58B2f4210Fe06fe678d1A1B7"
	PaymentVaultAddress        = "0x714238F32A7aE70C0D208D58Cc041D8Dda28e813"
	KawaiDistributorAddr       = "0xf4CCb09208cA77153e1681d256247dae0ff119ba"
	USDTDistributorAddr        = "0xE964B52D496F37749bd0caF287A356afdC10836C"
	CashbackDistributorAddress = "0xcc992d001Bc1963A44212D62F711E502DE162B8E"

	// Holder Scanner Configuration
	// HolderScanStartBlock: Starting block for holder scanning
	// - Testnet: Set to recent block to avoid RPC limits (max 100 blocks per query)
	// - Mainnet: Set to token deployment block to optimize performance
	// Example: If token deployed at block 1000000, set HolderScanStartBlock = 1000000
	// Current: 5437070 (for testing - within 90 blocks to stay under RPC 100 block limit)
	HolderScanStartBlock = 5437070
)
