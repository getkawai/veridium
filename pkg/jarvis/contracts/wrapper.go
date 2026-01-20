package contracts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/internal/generate/abi/cashbackdistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/distributor"
	"github.com/kawai-network/veridium/internal/generate/abi/escrow"
	"github.com/kawai-network/veridium/internal/generate/abi/kawaitoken"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/referraldistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/vault"
	"github.com/kawai-network/veridium/pkg/jarvis/util"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
)

// ResolveAddress resolves a string (hex or name) to a common.Address using Jarvis
func ResolveAddress(addrStr string) (common.Address, error) {
	addr, _, err := util.GetAddressFromString(addrStr)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(addr), nil
}

// KawaiToken wraps the generated KawaiToken binding
func KawaiToken(addrStr string, r *reader.EthReader) (*kawaitoken.KawaiToken, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return kawaitoken.NewKawaiToken(addr, backend)
}

// Escrow wraps the generated OTCMarket binding
func Escrow(addrStr string, r *reader.EthReader) (*escrow.OTCMarket, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return escrow.NewOTCMarket(addr, backend)
}

// Vault wraps the generated PaymentVault binding
func Vault(addrStr string, r *reader.EthReader) (*vault.PaymentVault, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return vault.NewPaymentVault(addr, backend)
}

// MerkleDistributor wraps the generated MerkleDistributor binding
func MerkleDistributor(addrStr string, r *reader.EthReader) (*distributor.MerkleDistributor, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return distributor.NewMerkleDistributor(addr, backend)
}

// MiningRewardDistributor wraps the generated MiningRewardDistributor binding
// This contract supports referral-based mining rewards with flexible developer addresses
func MiningRewardDistributor(addrStr string, r *reader.EthReader) (*miningdistributor.MiningRewardDistributor, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return miningdistributor.NewMiningRewardDistributor(addr, backend)
}

// CashbackDistributor wraps the generated DepositCashbackDistributor binding
// This contract distributes KAWAI cashback rewards for USDT deposits
func CashbackDistributor(addrStr string, r *reader.EthReader) (*cashbackdistributor.DepositCashbackDistributor, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return cashbackdistributor.NewDepositCashbackDistributor(addr, backend)
}

// ReferralRewardDistributor wraps the generated ReferralRewardDistributor binding
// This contract distributes KAWAI referral commission rewards
func ReferralRewardDistributor(addrStr string, r *reader.EthReader) (*referraldistributor.ReferralRewardDistributor, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return referraldistributor.NewReferralRewardDistributor(addr, backend)
}

// Stablecoin wraps any ERC-20 stablecoin token (MockUSDT on testnet, USDC on mainnet)
// Uses KawaiToken binding since it implements standard ERC-20 interface
func Stablecoin(addrStr string, r *reader.EthReader) (*kawaitoken.KawaiToken, error) {
	addr, err := ResolveAddress(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve stablecoin address %s: %w", addrStr, err)
	}
	backend := NewJarvisBackend(r)
	return kawaitoken.NewKawaiToken(addr, backend)
}
