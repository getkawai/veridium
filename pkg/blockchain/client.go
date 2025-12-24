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
		return nil, fmt.Errorf("failed to connect to BSC RPC: %w", err)
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
