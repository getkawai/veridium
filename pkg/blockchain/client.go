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
	"github.com/kawai-network/veridium/internal/generate/abi/kawaitoken"
	"github.com/kawai-network/veridium/internal/generate/abi/mockstablecoin"
	"github.com/kawai-network/veridium/internal/generate/abi/otcmarket"
	"golang.org/x/time/rate"
)

type Config struct {
	RPCUrl           string
	TokenAddress     string
	OTCMarketAddress string
	USDTAddress      string
}

type Client struct {
	EthClient   *ethclient.Client
	Token       *kawaitoken.KawaiToken
	OTCMarket   *otcmarket.OTCMarket
	USDT        *mockstablecoin.MockStablecoin
	ChainID     *big.Int
	rateLimiter *rate.Limiter // ✅ Rate limiter for RPC calls
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

	otcMarketAddress := common.HexToAddress(cfg.OTCMarketAddress)
	otcMarketInstance, err := otcmarket.NewOTCMarket(otcMarketAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load OTCMarket: %w", err)
	}

	usdtAddress := common.HexToAddress(cfg.USDTAddress)
	usdtInstance, err := mockstablecoin.NewMockStablecoin(usdtAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load MockStablecoin: %w", err)
	}

	return &Client{
		EthClient:   client,
		Token:       tokenInstance,
		OTCMarket:   otcMarketInstance,
		USDT:        usdtInstance,
		ChainID:     chainID,
		rateLimiter: rate.NewLimiter(rate.Limit(10), 20), // ✅ 10 RPC calls/sec, burst 20
	}, nil
}

// GetTotalSupply retrieves the current total supply of KAWAI tokens from the smart contract.
func (c *Client) GetTotalSupply(ctx context.Context) (*big.Int, error) {
	return c.Token.TotalSupply(nil)
}

// GetMaxSupply retrieves the cap (max supply) from the smart contract.
func (c *Client) GetMaxSupply(ctx context.Context) (*big.Int, error) {
	return c.Token.Cap(nil)
}

// Marketplace Operations
// Requirements: 6.2, 6.3, 6.4

// MarketplaceCreateOrder creates a new sell order on the OTC marketplace
// Requirements: 6.2 - Smart contract parameter validation and createOrder function call
func (c *Client) MarketplaceCreateOrder(ctx context.Context, transactOpts *bind.TransactOpts, tokenAmount, stablecoinPrice *big.Int) (*types.Transaction, error) {
	// Validate parameters
	if tokenAmount == nil || tokenAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("token amount must be greater than zero")
	}
	if stablecoinPrice == nil || stablecoinPrice.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("stablecoin price must be greater than zero")
	}
	if transactOpts == nil {
		return nil, fmt.Errorf("transaction options cannot be nil")
	}

	log.Printf("Creating marketplace order: %s KAWAI tokens for %s stablecoin", tokenAmount.String(), stablecoinPrice.String())

	// Call the smart contract's createOrder function
	tx, err := c.OTCMarket.CreateOrder(transactOpts, tokenAmount, stablecoinPrice)
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
	tx, err := c.OTCMarket.BuyOrder(transactOpts, orderID)
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
	tx, err := c.OTCMarket.CancelOrder(transactOpts, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order on smart contract: %w", err)
	}

	log.Printf("✅ Cancel order transaction submitted: %s", tx.Hash().Hex())
	return tx, nil
}

// MarketplaceGetOrder retrieves order information from the smart contract
func (c *Client) MarketplaceGetOrder(ctx context.Context, orderID *big.Int) (struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
}, error) {
	// Validate parameters
	if orderID == nil || orderID.Cmp(big.NewInt(0)) < 0 {
		return struct {
			Id              *big.Int
			Seller          common.Address
			TokenAmount     *big.Int
			PriceInUSDT     *big.Int
			RemainingAmount *big.Int
			IsActive        bool
		}{}, fmt.Errorf("order ID must be a valid non-negative number")
	}

	// Call the smart contract's getOrder function (new view function)
	order, err := c.OTCMarket.GetOrder(nil, orderID)
	if err != nil {
		return struct {
			Id              *big.Int
			Seller          common.Address
			TokenAmount     *big.Int
			PriceInUSDT     *big.Int
			RemainingAmount *big.Int
			IsActive        bool
		}{}, fmt.Errorf("failed to get order from smart contract: %w", err)
	}

	return order, nil
}

// MarketplaceGetOrdersCount returns the total number of orders in the contract
func (c *Client) MarketplaceGetOrdersCount(ctx context.Context) (*big.Int, error) {
	count, err := c.OTCMarket.GetOrdersCount(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders count from smart contract: %w", err)
	}
	return count, nil
}

// GetUSDTBalance returns the stablecoin token balance for a given address
// Note: Function name kept for backward compatibility, works with MockStablecoin (testnet) or USDC (mainnet)
func (c *Client) GetUSDTBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	balance, err := c.USDT.BalanceOf(nil, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get stablecoin token balance: %w", err)
	}
	return balance, nil
}

// ValidateTradeBalance validates that the buyer has sufficient stablecoin balance for the trade
func (c *Client) ValidateTradeBalance(ctx context.Context, buyer common.Address, stablecoinAmount *big.Int) error {
	balance, err := c.GetUSDTBalance(ctx, buyer)
	if err != nil {
		return fmt.Errorf("failed to check buyer stablecoin balance: %w", err)
	}

	if balance.Cmp(stablecoinAmount) < 0 {
		return fmt.Errorf("insufficient stablecoin balance: has %s, needs %s", balance.String(), stablecoinAmount.String())
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

// ✅ NEW: Buy partial order
func (c *Client) MarketplaceBuyOrderPartial(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int, amount *big.Int) (*types.Transaction, error) {
	// Validate parameters
	if orderID == nil || orderID.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("order ID must be a valid non-negative number")
	}
	if amount == nil || amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if transactOpts == nil {
		return nil, fmt.Errorf("transaction options cannot be nil")
	}

	// Call smart contract
	tx, err := c.OTCMarket.BuyOrderPartial(transactOpts, orderID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to buy partial order on smart contract: %w", err)
	}

	log.Printf("Partial buy order transaction submitted: %s", tx.Hash().Hex())
	return tx, nil
}

// ✅ NEW: Get orders by seller
func (c *Client) MarketplaceGetOrdersBySeller(ctx context.Context, seller common.Address, offset, limit *big.Int) ([]otcmarket.OTCMarketOrder, error) {
	// Validate parameters
	if offset == nil {
		offset = big.NewInt(0)
	}
	if limit == nil {
		limit = big.NewInt(100)
	}

	// ✅ Wait for rate limit
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// Call smart contract
	orders, err := c.OTCMarket.GetOrdersBySeller(nil, seller, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by seller from smart contract: %w", err)
	}

	return orders, nil
}

// ✅ NEW: Get active orders
func (c *Client) MarketplaceGetActiveOrders(ctx context.Context, offset, limit *big.Int) ([]otcmarket.OTCMarketOrder, error) {
	// Validate parameters
	if offset == nil {
		offset = big.NewInt(0)
	}
	if limit == nil {
		limit = big.NewInt(100)
	}

	// Wait for rate limit
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// Call smart contract
	orders, err := c.OTCMarket.GetActiveOrders(nil, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get active orders from smart contract: %w", err)
	}

	return orders, nil
}

// ✅ NEW: Get multiple orders at once
func (c *Client) MarketplaceGetOrders(ctx context.Context, orderIDs []*big.Int) ([]otcmarket.OTCMarketOrder, error) {
	// Validate parameters
	if len(orderIDs) == 0 {
		return []otcmarket.OTCMarketOrder{}, nil
	}

	// Wait for rate limit
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// Call smart contract
	orders, err := c.OTCMarket.GetOrders(nil, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders from smart contract: %w", err)
	}

	return orders, nil
}
