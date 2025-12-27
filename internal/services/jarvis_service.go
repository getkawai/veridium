package services

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	jarviscommon "github.com/kawai-network/veridium/pkg/jarvis/common"
	"github.com/kawai-network/veridium/pkg/jarvis/db"
	"github.com/kawai-network/veridium/pkg/jarvis/networks"
	"github.com/kawai-network/veridium/pkg/jarvis/txanalyzer"
	"github.com/kawai-network/veridium/pkg/jarvis/util"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
)

// NetworkInfo represents a blockchain network for the frontend
type NetworkInfo struct {
	ID                 uint64 `json:"id"`
	Name               string `json:"name"`
	NativeTokenSymbol  string `json:"nativeTokenSymbol"`
	NativeTokenDecimal uint64 `json:"nativeTokenDecimal"`
	ExplorerURL        string `json:"explorerURL"`
	IsTestnet          bool   `json:"isTestnet"`
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
	readers map[uint64]*reader.EthReader
	wallet  *WalletService
}

// NewJarvisService creates a new JarvisService instance
func NewJarvisService(wallet *WalletService) *JarvisService {
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

	if r, exists := s.readers[networkID]; exists {
		return r, network, nil
	}

	r := reader.NewEthReaderGeneric(network.GetDefaultNodes(), nil)
	s.readers[networkID] = r
	return r, network, nil
}

// GetSupportedNetworks returns all supported blockchain networks
func (s *JarvisService) GetSupportedNetworks() []NetworkInfo {
	allNetworks := networks.GetSupportedNetworks()
	result := make([]NetworkInfo, 0, len(allNetworks))

	for _, n := range allNetworks {
		name := n.GetName()
		isTestnet := strings.Contains(strings.ToLower(name), "testnet") ||
			strings.Contains(strings.ToLower(name), "ropsten") ||
			strings.Contains(strings.ToLower(name), "kovan") ||
			strings.Contains(strings.ToLower(name), "rinkeby") ||
			strings.Contains(strings.ToLower(name), "mumbai")

		result = append(result, NetworkInfo{
			ID:                 n.GetChainID(),
			Name:               name,
			NativeTokenSymbol:  n.GetNativeTokenSymbol(),
			NativeTokenDecimal: n.GetNativeTokenDecimal(),
			ExplorerURL:        n.GetBlockExplorerAPIURL(),
			IsTestnet:          isTestnet,
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

	return &NetworkInfo{
		ID:                 network.GetChainID(),
		Name:               name,
		NativeTokenSymbol:  network.GetNativeTokenSymbol(),
		NativeTokenDecimal: network.GetNativeTokenDecimal(),
		ExplorerURL:        network.GetBlockExplorerAPIURL(),
		IsTestnet:          isTestnet,
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
