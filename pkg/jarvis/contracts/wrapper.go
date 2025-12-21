package contracts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/internal/generate/abi/escrow"
	"github.com/kawai-network/veridium/internal/generate/abi/kawaitoken"
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
