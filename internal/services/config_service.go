package services

import (
	"strings"

	"github.com/kawai-network/contracts"
)

// NetworkEnvironment represents a blockchain network
type NetworkEnvironment struct {
	Name      string `json:"name"`
	ChainId   uint64 `json:"chainId"`
	IsTestnet bool   `json:"isTestnet"`
}

// ContractAddresses represents all smart contract addresses
type ContractAddresses struct {
	// Token Addresses
	Usdt  string `json:"usdt"` // Stablecoin address: MockStablecoin (testnet) or USDC (mainnet)
	Kawai string `json:"kawai"`

	// Payment & Vault
	PaymentVault string `json:"paymentVault"`

	// OTC Market
	OtcMarket string `json:"otcMarket"`

	// Reward Distributors
	MiningDistributor   string `json:"miningDistributor"`
	CashbackDistributor string `json:"cashbackDistributor"`
	ReferralDistributor string `json:"referralDistributor"`

	// Revenue Sharing
	UsdtDistributor string `json:"usdtDistributor"`
}

// BackendConfig represents the complete backend configuration
type BackendConfig struct {
	Environment string             `json:"environment"` // "testnet" | "mainnet"
	Network     NetworkEnvironment `json:"network"`
	Contracts   ContractAddresses  `json:"contracts"`
}

// ConfigService exposes backend configuration to frontend
type ConfigService struct{}

// GetConfig returns current backend configuration
// This reads from internal/constant which is generated from .env
func (s *ConfigService) GetConfig() BackendConfig {
	// Determine environment from explicit chain ID or RPC URL
	// Priority: Use explicit MONAD_CHAIN_ID if set, otherwise infer from RPC URL
	var env string
	var chainId uint64
	var networkName string
	var isTestnet bool

	// Try to determine from RPC URL as fallback
	isTestnetRpc := strings.Contains(contracts.MonadRpcUrl, "testnet")

	if isTestnetRpc {
		env = "testnet"
		chainId = 10143
		networkName = "Monad Testnet"
		isTestnet = true
	} else {
		env = "mainnet"
		chainId = 143 // Monad Mainnet chain ID
		networkName = "Monad Mainnet"
		isTestnet = false
	}

	return BackendConfig{
		Environment: env,
		Network: NetworkEnvironment{
			Name:      networkName,
			ChainId:   chainId,
			IsTestnet: isTestnet,
		},
		Contracts: ContractAddresses{
			// Token Addresses
			Usdt:  contracts.StablecoinAddress,
			Kawai: contracts.KawaiTokenAddress,

			// Payment & Vault
			PaymentVault: contracts.PaymentVaultAddress,

			// OTC Market
			OtcMarket: contracts.OTCMarketAddress,

			// Reward Distributors
			MiningDistributor:   contracts.MiningRewardDistributorAddress,
			CashbackDistributor: contracts.CashbackDistributorAddress,
			ReferralDistributor: contracts.ReferralDistributorAddress,

			// Revenue Sharing
			UsdtDistributor: contracts.RevenueDistributorAddress,
		},
	}
}

// GetEnvironment returns current environment name
func (s *ConfigService) GetEnvironment() string {
	config := s.GetConfig()
	return config.Environment
}

// IsTestnet returns true if running on testnet
func (s *ConfigService) IsTestnet() bool {
	config := s.GetConfig()
	return config.Network.IsTestnet
}

// IsMainnet returns true if running on mainnet
func (s *ConfigService) IsMainnet() bool {
	return !s.IsTestnet()
}

// GetNetwork returns current network configuration
func (s *ConfigService) GetNetwork() NetworkEnvironment {
	config := s.GetConfig()
	return config.Network
}

// GetContracts returns all contract addresses
func (s *ConfigService) GetContracts() ContractAddresses {
	config := s.GetConfig()
	return config.Contracts
}

// GetContractAddress returns specific contract address by name
func (s *ConfigService) GetContractAddress(name string) string {
	config := s.GetConfig()

	switch strings.ToLower(name) {
	case "usdt":
		return config.Contracts.Usdt
	case "kawai":
		return config.Contracts.Kawai
	case "paymentvault", "payment_vault":
		return config.Contracts.PaymentVault
	case "otcmarket", "otc_market", "otc":
		return config.Contracts.OtcMarket
	case "miningdistributor", "mining_distributor":
		return config.Contracts.MiningDistributor
	case "cashbackdistributor", "cashback_distributor":
		return config.Contracts.CashbackDistributor
	case "referraldistributor", "referral_distributor":
		return config.Contracts.ReferralDistributor
	case "usdtdistributor", "usdt_distributor":
		return config.Contracts.UsdtDistributor
	default:
		return ""
	}
}
