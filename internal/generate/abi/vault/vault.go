// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vault

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PaymentVaultMetaData contains all meta data concerning the PaymentVault contract.
var PaymentVaultMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_usdt\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdt\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001017380380620010178339818101604052810190620000379190620002a2565b80600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603620000ad5760006040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401620000a49190620002fa565b60405180910390fd5b620000be816200017460201b60201c565b5060018081905550600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000138576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200012f9062000378565b60405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505050506200039a565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200026a826200023d565b9050919050565b6200027c816200025d565b81146200028857600080fd5b50565b6000815190506200029c8162000271565b92915050565b60008060408385031215620002bc57620002bb62000238565b5b6000620002cc858286016200028b565b9250506020620002df858286016200028b565b9150509250929050565b620002f4816200025d565b82525050565b6000602082019050620003116000830184620002e9565b92915050565b600082825260208201905092915050565b7f496e76616c696420555344542061646472657373000000000000000000000000600082015250565b60006200036060148362000317565b91506200036d8262000328565b602082019050919050565b60006020820190508181036000830152620003938162000351565b9050919050565b608051610c4c620003cb60003960008181610103015281816101b50152818161035a015261043b0152610c4c6000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80632f48ab7d14610067578063715018a6146100855780638da5cb5b1461008f578063b6b55f25146100ad578063f2fde38b146100c9578063f3fef3a3146100e5575b600080fd5b61006f610101565b60405161007c919061089d565b60405180910390f35b61008d610125565b005b610097610139565b6040516100a491906108d9565b60405180910390f35b6100c760048036038101906100c2919061092f565b610162565b005b6100e360048036038101906100de9190610988565b610253565b005b6100ff60048036038101906100fa91906109b5565b6102d9565b005b7f000000000000000000000000000000000000000000000000000000000000000081565b61012d6104d9565b6101376000610560565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b61016a610624565b600081116101ad576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101a490610a52565b60405180910390fd5b6101fa3330837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661066a909392919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff167f2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4826040516102409190610a81565b60405180910390a26102506106ec565b50565b61025b6104d9565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036102cd5760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016102c491906108d9565b60405180910390fd5b6102d681610560565b50565b6102e16104d9565b6102e9610624565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610358576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161034f90610ae8565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016103b191906108d9565b602060405180830381865afa1580156103ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103f29190610b1d565b811115610434576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161042b90610b96565b60405180910390fd5b61047f82827f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166106f59092919063ffffffff16565b8173ffffffffffffffffffffffffffffffffffffffff167f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5826040516104c59190610a81565b60405180910390a26104d56106ec565b5050565b6104e1610774565b73ffffffffffffffffffffffffffffffffffffffff166104ff610139565b73ffffffffffffffffffffffffffffffffffffffff161461055e57610522610774565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161055591906108d9565b60405180910390fd5b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600260015403610660576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002600181905550565b6106e6848573ffffffffffffffffffffffffffffffffffffffff166323b872dd86868660405160240161069f93929190610bb6565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061077c565b50505050565b60018081905550565b61076f838473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8585604051602401610728929190610bed565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061077c565b505050565b600033905090565b600080602060008451602086016000885af18061079f576040513d6000823e3d81fd5b3d9250600051915050600082146107ba5760018114156107d6565b60008473ffffffffffffffffffffffffffffffffffffffff163b145b1561081857836040517f5274afe700000000000000000000000000000000000000000000000000000000815260040161080f91906108d9565b60405180910390fd5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600061086361085e6108598461081e565b61083e565b61081e565b9050919050565b600061087582610848565b9050919050565b60006108878261086a565b9050919050565b6108978161087c565b82525050565b60006020820190506108b2600083018461088e565b92915050565b60006108c38261081e565b9050919050565b6108d3816108b8565b82525050565b60006020820190506108ee60008301846108ca565b92915050565b600080fd5b6000819050919050565b61090c816108f9565b811461091757600080fd5b50565b60008135905061092981610903565b92915050565b600060208284031215610945576109446108f4565b5b60006109538482850161091a565b91505092915050565b610965816108b8565b811461097057600080fd5b50565b6000813590506109828161095c565b92915050565b60006020828403121561099e5761099d6108f4565b5b60006109ac84828501610973565b91505092915050565b600080604083850312156109cc576109cb6108f4565b5b60006109da85828601610973565b92505060206109eb8582860161091a565b9150509250929050565b600082825260208201905092915050565b7f416d6f756e74206d757374206265203e20300000000000000000000000000000600082015250565b6000610a3c6012836109f5565b9150610a4782610a06565b602082019050919050565b60006020820190508181036000830152610a6b81610a2f565b9050919050565b610a7b816108f9565b82525050565b6000602082019050610a966000830184610a72565b92915050565b7f496e76616c696420726563697069656e74000000000000000000000000000000600082015250565b6000610ad26011836109f5565b9150610add82610a9c565b602082019050919050565b60006020820190508181036000830152610b0181610ac5565b9050919050565b600081519050610b1781610903565b92915050565b600060208284031215610b3357610b326108f4565b5b6000610b4184828501610b08565b91505092915050565b7f496e73756666696369656e742062616c616e6365000000000000000000000000600082015250565b6000610b806014836109f5565b9150610b8b82610b4a565b602082019050919050565b60006020820190508181036000830152610baf81610b73565b9050919050565b6000606082019050610bcb60008301866108ca565b610bd860208301856108ca565b610be56040830184610a72565b949350505050565b6000604082019050610c0260008301856108ca565b610c0f6020830184610a72565b939250505056fea264697066735822122054a84314d66a35f8f90549311a8d2e2a5f83f26caac4e268abbba565cf0e95e464736f6c63430008140033",
}

// PaymentVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymentVaultMetaData.ABI instead.
var PaymentVaultABI = PaymentVaultMetaData.ABI

// PaymentVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PaymentVaultMetaData.Bin instead.
var PaymentVaultBin = PaymentVaultMetaData.Bin

// DeployPaymentVault deploys a new Ethereum contract, binding an instance of PaymentVault to it.
func DeployPaymentVault(auth *bind.TransactOpts, backend bind.ContractBackend, _usdt common.Address, initialOwner common.Address) (common.Address, *types.Transaction, *PaymentVault, error) {
	parsed, err := PaymentVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymentVaultBin), backend, _usdt, initialOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PaymentVault{PaymentVaultCaller: PaymentVaultCaller{contract: contract}, PaymentVaultTransactor: PaymentVaultTransactor{contract: contract}, PaymentVaultFilterer: PaymentVaultFilterer{contract: contract}}, nil
}

// PaymentVault is an auto generated Go binding around an Ethereum contract.
type PaymentVault struct {
	PaymentVaultCaller     // Read-only binding to the contract
	PaymentVaultTransactor // Write-only binding to the contract
	PaymentVaultFilterer   // Log filterer for contract events
}

// PaymentVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type PaymentVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PaymentVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PaymentVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PaymentVaultSession struct {
	Contract     *PaymentVault     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PaymentVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PaymentVaultCallerSession struct {
	Contract *PaymentVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PaymentVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PaymentVaultTransactorSession struct {
	Contract     *PaymentVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PaymentVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type PaymentVaultRaw struct {
	Contract *PaymentVault // Generic contract binding to access the raw methods on
}

// PaymentVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PaymentVaultCallerRaw struct {
	Contract *PaymentVaultCaller // Generic read-only contract binding to access the raw methods on
}

// PaymentVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PaymentVaultTransactorRaw struct {
	Contract *PaymentVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPaymentVault creates a new instance of PaymentVault, bound to a specific deployed contract.
func NewPaymentVault(address common.Address, backend bind.ContractBackend) (*PaymentVault, error) {
	contract, err := bindPaymentVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PaymentVault{PaymentVaultCaller: PaymentVaultCaller{contract: contract}, PaymentVaultTransactor: PaymentVaultTransactor{contract: contract}, PaymentVaultFilterer: PaymentVaultFilterer{contract: contract}}, nil
}

// NewPaymentVaultCaller creates a new read-only instance of PaymentVault, bound to a specific deployed contract.
func NewPaymentVaultCaller(address common.Address, caller bind.ContractCaller) (*PaymentVaultCaller, error) {
	contract, err := bindPaymentVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentVaultCaller{contract: contract}, nil
}

// NewPaymentVaultTransactor creates a new write-only instance of PaymentVault, bound to a specific deployed contract.
func NewPaymentVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*PaymentVaultTransactor, error) {
	contract, err := bindPaymentVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentVaultTransactor{contract: contract}, nil
}

// NewPaymentVaultFilterer creates a new log filterer instance of PaymentVault, bound to a specific deployed contract.
func NewPaymentVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*PaymentVaultFilterer, error) {
	contract, err := bindPaymentVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymentVaultFilterer{contract: contract}, nil
}

// bindPaymentVault binds a generic wrapper to an already deployed contract.
func bindPaymentVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaymentVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PaymentVault *PaymentVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PaymentVault.Contract.PaymentVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PaymentVault *PaymentVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentVault.Contract.PaymentVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PaymentVault *PaymentVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PaymentVault.Contract.PaymentVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PaymentVault *PaymentVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PaymentVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PaymentVault *PaymentVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PaymentVault *PaymentVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PaymentVault.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PaymentVault *PaymentVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PaymentVault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PaymentVault *PaymentVaultSession) Owner() (common.Address, error) {
	return _PaymentVault.Contract.Owner(&_PaymentVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PaymentVault *PaymentVaultCallerSession) Owner() (common.Address, error) {
	return _PaymentVault.Contract.Owner(&_PaymentVault.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_PaymentVault *PaymentVaultCaller) Usdt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PaymentVault.contract.Call(opts, &out, "usdt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_PaymentVault *PaymentVaultSession) Usdt() (common.Address, error) {
	return _PaymentVault.Contract.Usdt(&_PaymentVault.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_PaymentVault *PaymentVaultCallerSession) Usdt() (common.Address, error) {
	return _PaymentVault.Contract.Usdt(&_PaymentVault.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_PaymentVault *PaymentVaultTransactor) Deposit(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _PaymentVault.contract.Transact(opts, "deposit", _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_PaymentVault *PaymentVaultSession) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _PaymentVault.Contract.Deposit(&_PaymentVault.TransactOpts, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) returns()
func (_PaymentVault *PaymentVaultTransactorSession) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _PaymentVault.Contract.Deposit(&_PaymentVault.TransactOpts, _amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PaymentVault *PaymentVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PaymentVault *PaymentVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _PaymentVault.Contract.RenounceOwnership(&_PaymentVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PaymentVault *PaymentVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PaymentVault.Contract.RenounceOwnership(&_PaymentVault.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PaymentVault *PaymentVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PaymentVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PaymentVault *PaymentVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PaymentVault.Contract.TransferOwnership(&_PaymentVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PaymentVault *PaymentVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PaymentVault.Contract.TransferOwnership(&_PaymentVault.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _to, uint256 _amount) returns()
func (_PaymentVault *PaymentVaultTransactor) Withdraw(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _PaymentVault.contract.Transact(opts, "withdraw", _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _to, uint256 _amount) returns()
func (_PaymentVault *PaymentVaultSession) Withdraw(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _PaymentVault.Contract.Withdraw(&_PaymentVault.TransactOpts, _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _to, uint256 _amount) returns()
func (_PaymentVault *PaymentVaultTransactorSession) Withdraw(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _PaymentVault.Contract.Withdraw(&_PaymentVault.TransactOpts, _to, _amount)
}

// PaymentVaultDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the PaymentVault contract.
type PaymentVaultDepositedIterator struct {
	Event *PaymentVaultDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentVaultDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentVaultDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentVaultDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentVaultDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentVaultDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentVaultDeposited represents a Deposited event raised by the PaymentVault contract.
type PaymentVaultDeposited struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed user, uint256 amount)
func (_PaymentVault *PaymentVaultFilterer) FilterDeposited(opts *bind.FilterOpts, user []common.Address) (*PaymentVaultDepositedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PaymentVault.contract.FilterLogs(opts, "Deposited", userRule)
	if err != nil {
		return nil, err
	}
	return &PaymentVaultDepositedIterator{contract: _PaymentVault.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed user, uint256 amount)
func (_PaymentVault *PaymentVaultFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *PaymentVaultDeposited, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PaymentVault.contract.WatchLogs(opts, "Deposited", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentVaultDeposited)
				if err := _PaymentVault.contract.UnpackLog(event, "Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposited is a log parse operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed user, uint256 amount)
func (_PaymentVault *PaymentVaultFilterer) ParseDeposited(log types.Log) (*PaymentVaultDeposited, error) {
	event := new(PaymentVaultDeposited)
	if err := _PaymentVault.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PaymentVault contract.
type PaymentVaultOwnershipTransferredIterator struct {
	Event *PaymentVaultOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentVaultOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentVaultOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentVaultOwnershipTransferred represents a OwnershipTransferred event raised by the PaymentVault contract.
type PaymentVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PaymentVault *PaymentVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PaymentVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PaymentVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PaymentVaultOwnershipTransferredIterator{contract: _PaymentVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PaymentVault *PaymentVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PaymentVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PaymentVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentVaultOwnershipTransferred)
				if err := _PaymentVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PaymentVault *PaymentVaultFilterer) ParseOwnershipTransferred(log types.Log) (*PaymentVaultOwnershipTransferred, error) {
	event := new(PaymentVaultOwnershipTransferred)
	if err := _PaymentVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentVaultWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the PaymentVault contract.
type PaymentVaultWithdrawnIterator struct {
	Event *PaymentVaultWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentVaultWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentVaultWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentVaultWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentVaultWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentVaultWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentVaultWithdrawn represents a Withdrawn event raised by the PaymentVault contract.
type PaymentVaultWithdrawn struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed to, uint256 amount)
func (_PaymentVault *PaymentVaultFilterer) FilterWithdrawn(opts *bind.FilterOpts, to []common.Address) (*PaymentVaultWithdrawnIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PaymentVault.contract.FilterLogs(opts, "Withdrawn", toRule)
	if err != nil {
		return nil, err
	}
	return &PaymentVaultWithdrawnIterator{contract: _PaymentVault.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed to, uint256 amount)
func (_PaymentVault *PaymentVaultFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *PaymentVaultWithdrawn, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PaymentVault.contract.WatchLogs(opts, "Withdrawn", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentVaultWithdrawn)
				if err := _PaymentVault.contract.UnpackLog(event, "Withdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawn is a log parse operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed to, uint256 amount)
func (_PaymentVault *PaymentVaultFilterer) ParseWithdrawn(log types.Log) (*PaymentVaultWithdrawn, error) {
	event := new(PaymentVaultWithdrawn)
	if err := _PaymentVault.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
