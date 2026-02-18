package services

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	jarviscommon "github.com/kawai-network/veridium/pkg/jarvis/common"
	"github.com/kawai-network/veridium/pkg/jarvis/db"
	"github.com/kawai-network/veridium/pkg/jarvis/networks"
	"github.com/kawai-network/veridium/pkg/jarvis/txanalyzer"
	"github.com/kawai-network/veridium/pkg/jarvis/util"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
	"github.com/kawai-network/x/constant"
)

// NetworkInfo represents a blockchain network for the frontend
type NetworkInfo struct {
	ID                 uint64 `json:"id"`
	Name               string `json:"name"`
	NativeTokenSymbol  string `json:"nativeTokenSymbol"`
	NativeTokenDecimal uint64 `json:"nativeTokenDecimal"`
	ExplorerURL        string `json:"explorerURL"`
	IsTestnet          bool   `json:"isTestnet"`
	Icon               string `json:"icon"`
	StablecoinSymbol   string `json:"stablecoinSymbol"` // "MOCK" (testnet) or "USDC" (mainnet)
	StablecoinName     string `json:"stablecoinName"`   // Full display name
	StablecoinShort    string `json:"stablecoinShort"`  // "USDT" (testnet) or "USDC" (mainnet) for messages
}

// TokenInfo represents ERC20 token information
type TokenInfo struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint64 `json:"decimals"`
	IsKnown  bool   `json:"isKnown"`
}

// BalanceInfo represents a balance result
type BalanceInfo struct {
	Raw       string `json:"raw"`
	Formatted string `json:"formatted"`
	Decimals  uint64 `json:"decimals"`
}

// GasEstimate represents gas estimation for a network
type GasEstimate struct {
	MaxGasPriceGwei float64 `json:"maxGasPriceGwei"`
	MaxTipGwei      float64 `json:"maxTipGwei"`
	IsDynamicFee    bool    `json:"isDynamicFee"`
}

// TxAnalysis represents a decoded transaction
type TxAnalysis struct {
	Hash        string      `json:"hash"`
	Status      string      `json:"status"`
	From        string      `json:"from"`
	To          string      `json:"to"`
	Value       string      `json:"value"`
	GasUsed     string      `json:"gasUsed"`
	GasCost     string      `json:"gasCost"`
	Method      string      `json:"method"`
	TxType      string      `json:"txType"`
	Params      []ParamInfo `json:"params"`
	Logs        []LogInfo   `json:"logs"`
	BlockNumber uint64      `json:"blockNumber"`
	Error       string      `json:"error,omitempty"`
}

// ParamInfo represents a decoded parameter
type ParamInfo struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// LogInfo represents a decoded event log
type LogInfo struct {
	Name   string      `json:"name"`
	Topics []string    `json:"topics"`
	Data   []ParamInfo `json:"data"`
}

// JarvisService provides blockchain utilities from the jarvis package
type JarvisService struct {
	readers   map[uint64]*reader.EthReader
	readersMu sync.RWMutex // Protect concurrent map access
	wallet    WalletOperations
}

// NewJarvisService creates a new JarvisService instance
func NewJarvisService(wallet WalletOperations) *JarvisService {
	return &JarvisService{
		readers: make(map[uint64]*reader.EthReader),
		wallet:  wallet,
	}
}

// getReader returns or creates an EthReader for a network
func (s *JarvisService) getReader(networkID uint64) (*reader.EthReader, networks.Network, error) {
	network, err := networks.GetNetworkByID(networkID)
	if err != nil {
		return nil, nil, fmt.Errorf("network not supported: %w", err)
	}

	// Check if reader exists (read lock)
	s.readersMu.RLock()
	r, exists := s.readers[networkID]
	s.readersMu.RUnlock()

	if exists {
		return r, network, nil
	}

	// Create new reader (write lock)
	s.readersMu.Lock()
	defer s.readersMu.Unlock()

	// Double-check after acquiring write lock (another goroutine might have created it)
	if r, exists := s.readers[networkID]; exists {
		return r, network, nil
	}

	r = reader.NewEthReaderGeneric(network.GetDefaultNodes(), nil)
	s.readers[networkID] = r
	return r, network, nil
}

// getWeb3IconSlugByChainID maps chain IDs to Web3Icons slugs for accurate matching
func getWeb3IconSlugByChainID(chainID uint64, name string) string {
	// Map known chain IDs to icon slugs
	switch chainID {
	// Ethereum networks
	case 1: // Ethereum Mainnet
		return "ethereum"
	case 3, 4, 42: // Ropsten, Rinkeby, Kovan (deprecated testnets)
		return "ethereum"
	case 10001: // Ethereum PoW
		return "ethereum"

	// BSC
	case 56: // BSC Mainnet
		return "binance-smart-chain"
	case 97: // BSC Testnet
		return "binance-smart-chain"

	// Polygon
	case 137: // Polygon (Matic) Mainnet
		return "polygon"
	case 80001: // Mumbai Testnet
		return "polygon"
	case 1101: // Polygon zkEVM
		return "polygon-zkevm"

	// Avalanche
	case 43114: // Avalanche C-Chain
		return "avalanche"
	case 43113: // Avalanche Fuji Testnet
		return "avalanche"

	// Fantom
	case 250: // Fantom Opera
		return "fantom"
	case 4002: // Fantom Testnet
		return "fantom"

	// Optimism
	case 10: // Optimism Mainnet
		return "optimism"
	case 420: // Optimism Goerli
		return "optimism"

	// Arbitrum
	case 42161: // Arbitrum One
		return "arbitrum"
	case 421613: // Arbitrum Goerli
		return "arbitrum"

	// Base
	case 8453: // Base Mainnet
		return "base"
	case 84531: // Base Goerli
		return "base"

	// Scroll
	case 534352: // Scroll Mainnet
		return "scroll"

	// Linea
	case 59144: // Linea Mainnet
		return "linea"

	// Monad
	case 143, 10143: // Monad Mainnet/Testnet
		return "monad"

	// BTTC - no icon available, fallback to ethereum
	case 199:
		return "ethereum"

	// Bitfi - no icon available, fallback to ethereum
	case 891891:
		return "ethereum"
	}

	// Fallback: try to match by name
	return getWeb3IconSlugByName(name)
}

// getWeb3IconSlugByName maps network names to Web3Icons slugs (fallback)
func getWeb3IconSlugByName(name string) string {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "ethereum") || strings.Contains(lower, "mainnet") {
		return "ethereum"
	} else if strings.Contains(lower, "bnb") || strings.Contains(lower, "bsc") || strings.Contains(lower, "binance") {
		return "binance-smart-chain"
	} else if strings.Contains(lower, "zkevm") {
		return "polygon-zkevm"
	} else if strings.Contains(lower, "matic") || strings.Contains(lower, "polygon") {
		return "polygon"
	} else if strings.Contains(lower, "avalanche") || strings.Contains(lower, "avax") || strings.Contains(lower, "snowtrace") {
		return "avalanche"
	} else if strings.Contains(lower, "fantom") || strings.Contains(lower, "ftm") {
		return "fantom"
	} else if strings.Contains(lower, "arbitrum") {
		return "arbitrum"
	} else if strings.Contains(lower, "optimism") {
		return "optimism"
	} else if strings.Contains(lower, "base") {
		return "base"
	} else if strings.Contains(lower, "scroll") {
		return "scroll"
	} else if strings.Contains(lower, "linea") {
		return "linea"
	} else if strings.Contains(lower, "monad") {
		return "monad"
	}
	// Default fallback to ethereum icon
	return "ethereum"
}

// getStablecoinInfo returns stablecoin information based on network testnet status
func getStablecoinInfo(isTestnet bool) (symbol, name, short string) {
	if isTestnet {
		return "MOCK", "Mock Stablecoin (Testnet)", "MOCK"
	}
	return "USDC", "USD Coin", "USDC"
}

// GetSupportedNetworks returns all supported blockchain networks
func (s *JarvisService) GetSupportedNetworks() []NetworkInfo {
	allNetworks := networks.GetSupportedNetworks()

	// Use map to deduplicate by chain ID (since alternative names cause duplicates)
	seen := make(map[uint64]bool)
	result := make([]NetworkInfo, 0, len(allNetworks))

	for _, n := range allNetworks {
		chainID := n.GetChainID()

		// Skip if we've already added this network (by chain ID)
		if seen[chainID] {
			continue
		}
		seen[chainID] = true

		name := n.GetName()
		isTestnet := strings.Contains(strings.ToLower(name), "testnet") ||
			strings.Contains(strings.ToLower(name), "ropsten") ||
			strings.Contains(strings.ToLower(name), "kovan") ||
			strings.Contains(strings.ToLower(name), "rinkeby") ||
			strings.Contains(strings.ToLower(name), "mumbai")

		// Use chain ID to determine icon for more accurate matching
		icon := getWeb3IconSlugByChainID(chainID, name)

		// Get stablecoin info based on network type
		stablecoinSymbol, stablecoinName, stablecoinShort := getStablecoinInfo(isTestnet)

		result = append(result, NetworkInfo{
			ID:                 chainID,
			Name:               name,
			NativeTokenSymbol:  n.GetNativeTokenSymbol(),
			NativeTokenDecimal: n.GetNativeTokenDecimal(),
			ExplorerURL:        n.GetBlockExplorerAPIURL(),
			IsTestnet:          isTestnet,
			Icon:               icon,
			StablecoinSymbol:   stablecoinSymbol,
			StablecoinName:     stablecoinName,
			StablecoinShort:    stablecoinShort,
		})
	}

	return result
}

// GetNetworkByID returns a specific network by chain ID
func (s *JarvisService) GetNetworkByID(chainID uint64) (*NetworkInfo, error) {
	network, err := networks.GetNetworkByID(chainID)
	if err != nil {
		return nil, err
	}

	name := network.GetName()
	isTestnet := strings.Contains(strings.ToLower(name), "testnet") ||
		strings.Contains(strings.ToLower(name), "ropsten") ||
		strings.Contains(strings.ToLower(name), "kovan") ||
		strings.Contains(strings.ToLower(name), "rinkeby") ||
		strings.Contains(strings.ToLower(name), "mumbai")

	// Get stablecoin info based on network type
	stablecoinSymbol, stablecoinName, stablecoinShort := getStablecoinInfo(isTestnet)

	return &NetworkInfo{
		ID:                 chainID,
		Name:               name,
		NativeTokenSymbol:  network.GetNativeTokenSymbol(),
		NativeTokenDecimal: network.GetNativeTokenDecimal(),
		ExplorerURL:        network.GetBlockExplorerAPIURL(),
		IsTestnet:          isTestnet,
		Icon:               getWeb3IconSlugByChainID(chainID, name),
		StablecoinSymbol:   stablecoinSymbol,
		StablecoinName:     stablecoinName,
		StablecoinShort:    stablecoinShort,
	}, nil
}

// GetTokenInfo fetches ERC20 token information from the blockchain
func (s *JarvisService) GetTokenInfo(tokenAddress string, networkID uint64) (*TokenInfo, error) {
	r, network, err := s.getReader(networkID)
	if err != nil {
		return nil, err
	}

	// Check if it's a known token
	knownName, isKnown := db.TOKENS[tokenAddress]

	// Get symbol from chain
	symbol, err := r.ERC20Symbol(tokenAddress)
	if err != nil {
		// If we have a known token, use that info
		if isKnown {
			return &TokenInfo{
				Address:  tokenAddress,
				Name:     knownName,
				Symbol:   strings.Split(knownName, " ")[0], // Extract symbol from "XXX token"
				Decimals: 18,                               // Default
				IsKnown:  true,
			}, nil
		}
		return nil, fmt.Errorf("failed to get token symbol: %w", err)
	}

	// Get decimals
	decimals, err := util.GetERC20Decimal(tokenAddress, network)
	if err != nil {
		decimals = 18 // Default
	}

	name := symbol + " Token"
	if isKnown {
		name = knownName
	}

	return &TokenInfo{
		Address:  tokenAddress,
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
		IsKnown:  isKnown,
	}, nil
}

// GetNativeBalance returns the native token balance for an address
func (s *JarvisService) GetNativeBalance(address string, networkID uint64) (*BalanceInfo, error) {
	r, network, err := s.getReader(networkID)
	if err != nil {
		return nil, err
	}

	balance, err := r.GetBalance(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	decimals := network.GetNativeTokenDecimal()
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	formatted := new(big.Float).Quo(new(big.Float).SetInt(balance), divisor)

	return &BalanceInfo{
		Raw:       balance.String(),
		Formatted: formatted.Text('f', 6),
		Decimals:  decimals,
	}, nil
}

// GetTokenBalance returns the ERC20 token balance for an address
func (s *JarvisService) GetTokenBalance(tokenAddress string, walletAddress string, networkID uint64) (*BalanceInfo, error) {
	r, network, err := s.getReader(networkID)
	if err != nil {
		return nil, err
	}

	balance, err := r.ERC20Balance(tokenAddress, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	decimals, err := util.GetERC20Decimal(tokenAddress, network)
	if err != nil {
		decimals = 18
	}

	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	formatted := new(big.Float).Quo(new(big.Float).SetInt(balance), divisor)

	return &BalanceInfo{
		Raw:       balance.String(),
		Formatted: formatted.Text('f', 6),
		Decimals:  decimals,
	}, nil
}

// EstimateGas returns gas price estimates for a network
func (s *JarvisService) EstimateGas(networkID uint64) (*GasEstimate, error) {
	r, _, err := s.getReader(networkID)
	if err != nil {
		return nil, err
	}

	maxGasPrice, maxTip, err := r.SuggestedGasSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	isDynamic, _ := r.CheckDynamicFeeTxAvailable()

	return &GasEstimate{
		MaxGasPriceGwei: maxGasPrice,
		MaxTipGwei:      maxTip,
		IsDynamicFee:    isDynamic,
	}, nil
}

// GetCurrentBlock returns the current block number
func (s *JarvisService) GetCurrentBlock(networkID uint64) (uint64, error) {
	r, _, err := s.getReader(networkID)
	if err != nil {
		return 0, err
	}

	return r.CurrentBlock()
}

// AnalyzeTransaction decodes and analyzes a transaction
func (s *JarvisService) AnalyzeTransaction(txHash string, networkID uint64) (*TxAnalysis, error) {
	r, network, err := s.getReader(networkID)
	if err != nil {
		return nil, err
	}

	// Get transaction info
	txInfo, err := r.TxInfoFromHash(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	result := &TxAnalysis{
		Hash:   txHash,
		Status: txInfo.Status,
	}

	if txInfo.Status == "notfound" {
		result.Error = "Transaction not found"
		return result, nil
	}

	if txInfo.Status == "pending" {
		result.TxType = "pending"
		if txInfo.Tx != nil {
			if txInfo.Tx.Extra.From != nil {
				result.From = txInfo.Tx.Extra.From.Hex()
			}
			if txInfo.Tx.To() != nil {
				result.To = txInfo.Tx.To().Hex()
			}
			result.Value = jarviscommon.BigToFloatString(txInfo.Tx.Value(), network.GetNativeTokenDecimal())
		}
		return result, nil
	}

	// Analyze completed transaction
	if txInfo.Tx != nil {
		if txInfo.Tx.Extra.From != nil {
			result.From = txInfo.Tx.Extra.From.Hex()
		}
		if txInfo.Tx.To() != nil {
			result.To = txInfo.Tx.To().Hex()
		}
		result.Value = jarviscommon.BigToFloatString(txInfo.Tx.Value(), network.GetNativeTokenDecimal())

		// Determine if it's a contract call
		if len(txInfo.Tx.Data()) > 0 {
			result.TxType = "contract call"
		} else {
			result.TxType = "transfer"
		}
	}

	if txInfo.Receipt != nil {
		result.GasUsed = fmt.Sprintf("%d", txInfo.Receipt.GasUsed)
		result.GasCost = jarviscommon.BigToFloatString(txInfo.GasCost(), network.GetNativeTokenDecimal())
		result.BlockNumber = txInfo.Receipt.BlockNumber.Uint64()
	}

	// Try to decode the function call
	if result.TxType == "contract call" && txInfo.Tx.To() != nil {
		analyzer, err := txanalyzer.EthAnalyzer(network)
		if err == nil {
			txResult := analyzer.AnalyzeOffline(
				&txInfo,
				util.GetABI,
				nil,
				true,
				network,
			)

			if txResult.FunctionCall != nil {
				result.Method = txResult.FunctionCall.Method

				// Convert params
				for _, p := range txResult.FunctionCall.Params {
					paramValue := ""
					if len(p.Values) > 0 {
						paramValue = p.Values[0].Value
					}
					result.Params = append(result.Params, ParamInfo{
						Name:  p.Name,
						Type:  p.Type,
						Value: paramValue,
					})
				}
			}

			// Convert logs
			for _, log := range txResult.Logs {
				logInfo := LogInfo{
					Name: log.Name,
				}
				for _, topic := range log.Topics {
					if len(topic.Value) > 0 {
						logInfo.Topics = append(logInfo.Topics, topic.Value[0].Value)
					}
				}
				for _, d := range log.Data {
					dataValue := ""
					if len(d.Values) > 0 {
						dataValue = d.Values[0].Value
					}
					logInfo.Data = append(logInfo.Data, ParamInfo{
						Name:  d.Name,
						Type:  d.Type,
						Value: dataValue,
					})
				}
				result.Logs = append(result.Logs, logInfo)
			}
		}
	}

	return result, nil
}

// IsContract checks if an address is a smart contract
func (s *JarvisService) IsContract(address string, networkID uint64) (bool, error) {
	r, _, err := s.getReader(networkID)
	if err != nil {
		return false, err
	}

	code, err := r.GetCode(address)
	if err != nil {
		return false, err
	}

	return len(code) > 0, nil
}

// ValidateAddress checks if a string is a valid EVM address
func (s *JarvisService) ValidateAddress(address string) bool {
	return common.IsHexAddress(address)
}

// LookupKnownToken checks if an address is in the known tokens database
func (s *JarvisService) LookupKnownToken(address string) (string, bool) {
	name, exists := db.TOKENS[address]
	return name, exists
}

// ProjectToken represents a project-specific token
type ProjectToken struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

// GetProjectTokens returns all project-specific tokens (deployed contracts)
func (s *JarvisService) GetProjectTokens() []ProjectToken {
	tokens := make([]ProjectToken, 0, len(constant.PROJECT_TOKENS))
	for addr, name := range constant.PROJECT_TOKENS {
		// Only include actual tokens (not contracts like Escrow, PaymentVault, Distributors)
		if name == "Stablecoin" || name == "KawaiToken" {
			tokens = append(tokens, ProjectToken{
				Address: addr,
				Name:    name,
				Symbol:  getSymbolFromName(name),
			})
		}
	}
	return tokens
}

// getSymbolFromName extracts symbol from token name
func getSymbolFromName(name string) string {
	switch name {
	case "Stablecoin":
		return "MOCK"
	case "KawaiToken":
		return "KAWAI"
	default:
		return name
	}
}

// GetTokenAllowance returns the allowance of a token
func (s *JarvisService) GetTokenAllowance(tokenAddress string, owner string, spender string, networkID uint64) (*BalanceInfo, error) {
	r, network, err := s.getReader(networkID)
	if err != nil {
		return nil, err
	}

	allowance, err := r.ERC20Allowance(tokenAddress, owner, spender)
	if err != nil {
		return nil, fmt.Errorf("failed to get allowance: %w", err)
	}

	decimals, err := util.GetERC20Decimal(tokenAddress, network)
	if err != nil {
		decimals = 18
	}

	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	formatted := new(big.Float).Quo(new(big.Float).SetInt(allowance), divisor)

	return &BalanceInfo{
		Raw:       allowance.String(),
		Formatted: formatted.Text('f', 6),
		Decimals:  decimals,
	}, nil
}

// GetTokenPrice returns the price of a token in USD
func (s *JarvisService) GetTokenPrice(token string, networkID uint64) (float64, error) {
	network, err := networks.GetNetworkByID(networkID)
	if err != nil {
		return 0, err
	}
	return util.GetCoinGeckoRateInUSD(network, token)
}
