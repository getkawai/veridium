package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/escrow"
	"github.com/kawai-network/veridium/internal/generate/abi/kawaitoken"
	"github.com/kawai-network/veridium/pkg/config"
)

type Config struct {
	RPCUrl        string
	TokenAddress  string
	EscrowAddress string
}

type Client struct {
	EthClient *ethclient.Client
	Token     *kawaitoken.KawaiToken
	Escrow    *escrow.OTCMarket
	ChainID   *big.Int
}

func NewClient(cfg Config) (*Client, error) {
	client, err := ethclient.Dial(cfg.RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	// Fetch ChainID to ensure connection is valid
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	log.Printf("Connected to Chain ID: %s", chainID.String())

	// Load Contracts
	tokenAddress := common.HexToAddress(cfg.TokenAddress)
	tokenInstance, err := kawaitoken.NewKawaiToken(tokenAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load KawaiToken: %w", err)
	}

	escrowAddress := common.HexToAddress(cfg.EscrowAddress)
	escrowInstance, err := escrow.NewOTCMarket(escrowAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load Escrow: %w", err)
	}

	return &Client{
		EthClient: client,
		Token:     tokenInstance,
		Escrow:    escrowInstance,
		ChainID:   chainID,
	}, nil
}

// GetRewardMode determines the current economic phase by checking totalSupply vs MAX_SUPPLY.
// Returns ModeMining if supply < MAX_SUPPLY (Phase 1), ModeUSDT if supply >= MAX_SUPPLY (Phase 2).
func (c *Client) GetRewardMode(ctx context.Context) (config.RewardMode, error) {
	// Get current total supply
	totalSupply, err := c.Token.TotalSupply(nil)
	if err != nil {
		return config.ModeMining, fmt.Errorf("failed to get total supply: %w", err)
	}

	// Get max supply constant from contract
	maxSupply, err := c.Token.MAXSUPPLY(nil)
	if err != nil {
		return config.ModeMining, fmt.Errorf("failed to get max supply: %w", err)
	}

	log.Printf("[Phase Check] Total Supply: %s / Max Supply: %s", totalSupply.String(), maxSupply.String())

	// Compare: if totalSupply >= maxSupply, we're in Phase 2 (USDT mode)
	if totalSupply.Cmp(maxSupply) >= 0 {
		log.Println("[Phase] MAX_SUPPLY reached. Entering Phase 2 (USDT Mode).")
		return config.ModeUSDT, nil
	}

	return config.ModeMining, nil
}
