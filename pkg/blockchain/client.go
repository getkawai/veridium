package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/escrow"
	"github.com/kawai-network/veridium/internal/generate/abi/kawaitoken"
	"github.com/kawai-network/veridium/internal/generate/abi/usdt"
	"github.com/kawai-network/veridium/pkg/config"
)

type Config struct {
	RPCUrl        string
	TokenAddress  string
	EscrowAddress string
	USDTAddress   string
}

type Client struct {
	EthClient *ethclient.Client
	Token     *kawaitoken.KawaiToken
	Escrow    *escrow.OTCMarket
	USDT      *usdt.MockUSDT
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

	usdtAddress := common.HexToAddress(cfg.USDTAddress)
	usdtInstance, err := usdt.NewMockUSDT(usdtAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load USDT: %w", err)
	}

	return &Client{
		EthClient: client,
		Token:     tokenInstance,
		Escrow:    escrowInstance,
		USDT:      usdtInstance,
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

// Marketplace Operations
// Requirements: 6.2, 6.3, 6.4

// MarketplaceCreateOrder creates a new sell order on the OTC marketplace
// Requirements: 6.2 - Smart contract parameter validation and createOrder function call
func (c *Client) MarketplaceCreateOrder(ctx context.Context, transactOpts *bind.TransactOpts, tokenAmount, usdtPrice *big.Int) (*types.Transaction, error) {
	// Validate parameters
	if tokenAmount == nil || tokenAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("token amount must be greater than zero")
	}
	if usdtPrice == nil || usdtPrice.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("USDT price must be greater than zero")
	}
	if transactOpts == nil {
		return nil, fmt.Errorf("transaction options cannot be nil")
	}

	log.Printf("Creating marketplace order: %s KAWAI tokens for %s USDT", tokenAmount.String(), usdtPrice.String())

	// Call the smart contract's createOrder function
	tx, err := c.Escrow.CreateOrder(transactOpts, tokenAmount, usdtPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to create order on smart contract: %w", err)
	}

	log.Printf("✅ Order creation transaction submitted: %s", tx.Hash().Hex())
	return tx, nil
}

// MarketplaceBuyOrder executes a buy order for the specified order ID
// Requirements: 6.3 - Atomic execution via buyOrder function
func (c *Client) MarketplaceBuyOrder(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int) (*types.Transaction, error) {
	// Validate parameters
	if orderID == nil || orderID.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("order ID must be a valid non-negative number")
	}
	if transactOpts == nil {
		return nil, fmt.Errorf("transaction options cannot be nil")
	}

	log.Printf("Buying marketplace order ID: %s", orderID.String())

	// Call the smart contract's buyOrder function
	tx, err := c.Escrow.BuyOrder(transactOpts, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to buy order on smart contract: %w", err)
	}

	log.Printf("✅ Buy order transaction submitted: %s", tx.Hash().Hex())
	return tx, nil
}

// MarketplaceCancelOrder cancels an existing order and returns tokens to seller
// Requirements: 6.4 - Order cancellation and token return verification
func (c *Client) MarketplaceCancelOrder(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int) (*types.Transaction, error) {
	// Validate parameters
	if orderID == nil || orderID.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("order ID must be a valid non-negative number")
	}
	if transactOpts == nil {
		return nil, fmt.Errorf("transaction options cannot be nil")
	}

	log.Printf("Cancelling marketplace order ID: %s", orderID.String())

	// Call the smart contract's cancelOrder function
	tx, err := c.Escrow.CancelOrder(transactOpts, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order on smart contract: %w", err)
	}

	log.Printf("✅ Cancel order transaction submitted: %s", tx.Hash().Hex())
	return tx, nil
}

// MarketplaceGetOrder retrieves order information from the smart contract
func (c *Client) MarketplaceGetOrder(ctx context.Context, orderID *big.Int) (struct {
	Id          *big.Int
	Seller      common.Address
	TokenAmount *big.Int
	PriceInUSDT *big.Int
	IsActive    bool
}, error) {
	// Validate parameters
	if orderID == nil || orderID.Cmp(big.NewInt(0)) < 0 {
		return struct {
			Id          *big.Int
			Seller      common.Address
			TokenAmount *big.Int
			PriceInUSDT *big.Int
			IsActive    bool
		}{}, fmt.Errorf("order ID must be a valid non-negative number")
	}

	// Call the smart contract's orders function
	order, err := c.Escrow.Orders(nil, orderID)
	if err != nil {
		return struct {
			Id          *big.Int
			Seller      common.Address
			TokenAmount *big.Int
			PriceInUSDT *big.Int
			IsActive    bool
		}{}, fmt.Errorf("failed to get order from smart contract: %w", err)
	}

	return order, nil
}

// MarketplaceGetOrdersCount returns the total number of orders in the contract
func (c *Client) MarketplaceGetOrdersCount(ctx context.Context) (*big.Int, error) {
	count, err := c.Escrow.GetOrdersCount(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders count from smart contract: %w", err)
	}
	return count, nil
}

// GetUSDTBalance returns the USDT token balance for a given address
func (c *Client) GetUSDTBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	balance, err := c.USDT.BalanceOf(nil, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get USDT token balance: %w", err)
	}
	return balance, nil
}

// ValidateTradeBalance validates that the buyer has sufficient USDT balance for the trade
func (c *Client) ValidateTradeBalance(ctx context.Context, buyer common.Address, usdtAmount *big.Int) error {
	balance, err := c.GetUSDTBalance(ctx, buyer)
	if err != nil {
		return fmt.Errorf("failed to check buyer USDT balance: %w", err)
	}

	if balance.Cmp(usdtAmount) < 0 {
		return fmt.Errorf("insufficient USDT balance: has %s, needs %s", balance.String(), usdtAmount.String())
	}

	return nil
}

// GetKawaiTokenBalance returns the KAWAI token balance for a given address
func (c *Client) GetKawaiTokenBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	balance, err := c.Token.BalanceOf(nil, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get KAWAI token balance: %w", err)
	}
	return balance, nil
}

// ValidateOrderCreationBalance validates that the seller has sufficient KAWAI token balance
func (c *Client) ValidateOrderCreationBalance(ctx context.Context, seller common.Address, tokenAmount *big.Int) error {
	balance, err := c.GetKawaiTokenBalance(ctx, seller)
	if err != nil {
		return fmt.Errorf("failed to check seller balance: %w", err)
	}

	if balance.Cmp(tokenAmount) < 0 {
		return fmt.Errorf("insufficient KAWAI token balance: has %s, needs %s", balance.String(), tokenAmount.String())
	}

	return nil
}
