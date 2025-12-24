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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_usdt\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"usdt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051610fe3380380610fe3833981810160405281019061003191906102d0565b805f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610099919061031d565b60405180910390fd5b6100b18161017f60201b60201c565b5060016100d06100c561024060201b60201c565b61026960201b60201c565b5f01819055505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610144576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161013b90610390565b60405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505050506103ae565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61029f82610276565b9050919050565b6102af81610295565b81146102b9575f5ffd5b50565b5f815190506102ca816102a6565b92915050565b5f5f604083850312156102e6576102e5610272565b5b5f6102f3858286016102bc565b9250506020610304858286016102bc565b9150509250929050565b61031781610295565b82525050565b5f6020820190506103305f83018461030e565b92915050565b5f82825260208201905092915050565b7f496e76616c6964205553445420616464726573730000000000000000000000005f82015250565b5f61037a601483610336565b915061038582610346565b602082019050919050565b5f6020820190508181035f8301526103a78161036e565b9050919050565b608051610c086103db5f395f8181610100015281816101ae01528181610350015261042f0152610c085ff3fe608060405234801561000f575f5ffd5b5060043610610060575f3560e01c80632f48ab7d14610064578063715018a6146100825780638da5cb5b1461008c578063b6b55f25146100aa578063f2fde38b146100c6578063f3fef3a3146100e2575b5f5ffd5b61006c6100fe565b60405161007991906108dc565b60405180910390f35b61008a610122565b005b610094610135565b6040516100a19190610915565b60405180910390f35b6100c460048036038101906100bf9190610965565b61015c565b005b6100e060048036038101906100db91906109ba565b61024c565b005b6100fc60048036038101906100f791906109e5565b6102d0565b005b7f000000000000000000000000000000000000000000000000000000000000000081565b61012a6104cd565b6101335f610554565b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b610164610615565b5f81116101a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161019d90610a7d565b60405180910390fd5b6101f33330837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610637909392919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff167f2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4826040516102399190610aaa565b60405180910390a261024961068c565b50565b6102546104cd565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036102c4575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016102bb9190610915565b60405180910390fd5b6102cd81610554565b50565b6102d86104cd565b6102e0610615565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361034e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161034590610b0d565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016103a79190610915565b602060405180830381865afa1580156103c2573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103e69190610b3f565b811115610428576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161041f90610bb4565b60405180910390fd5b61047382827f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166106a69092919063ffffffff16565b8173ffffffffffffffffffffffffffffffffffffffff167f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5826040516104b99190610aaa565b60405180910390a26104c961068c565b5050565b6104d56106f9565b73ffffffffffffffffffffffffffffffffffffffff166104f3610135565b73ffffffffffffffffffffffffffffffffffffffff1614610552576105166106f9565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016105499190610915565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b61061d610700565b600261062f61062a610741565b61076a565b5f0181905550565b610645848484846001610773565b61068657836040517f5274afe700000000000000000000000000000000000000000000000000000000815260040161067d9190610915565b60405180910390fd5b50505050565b600161069e610699610741565b61076a565b5f0181905550565b6106b383838360016107e4565b6106f457826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016106eb9190610915565b60405180910390fd5b505050565b5f33905090565b610708610846565b1561073f576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5f6323b872dd60e01b9050604051815f525f1960601c87166004525f1960601c86166024528460445260205f60645f5f8c5af1925060015f511483166107d15783831516156107c5573d5f823e3d81fd5b5f883b113d1516831692505b806040525f606052505095945050505050565b5f5f63a9059cbb60e01b9050604051815f525f1960601c86166004528460245260205f60445f5f8b5af1925060015f5114831661083857838315161561082c573d5f823e3d81fd5b5f873b113d1516831692505b806040525050949350505050565b5f6002610859610854610741565b61076a565b5f015414905090565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f6108a461089f61089a84610862565b610881565b610862565b9050919050565b5f6108b58261088a565b9050919050565b5f6108c6826108ab565b9050919050565b6108d6816108bc565b82525050565b5f6020820190506108ef5f8301846108cd565b92915050565b5f6108ff82610862565b9050919050565b61090f816108f5565b82525050565b5f6020820190506109285f830184610906565b92915050565b5f5ffd5b5f819050919050565b61094481610932565b811461094e575f5ffd5b50565b5f8135905061095f8161093b565b92915050565b5f6020828403121561097a5761097961092e565b5b5f61098784828501610951565b91505092915050565b610999816108f5565b81146109a3575f5ffd5b50565b5f813590506109b481610990565b92915050565b5f602082840312156109cf576109ce61092e565b5b5f6109dc848285016109a6565b91505092915050565b5f5f604083850312156109fb576109fa61092e565b5b5f610a08858286016109a6565b9250506020610a1985828601610951565b9150509250929050565b5f82825260208201905092915050565b7f416d6f756e74206d757374206265203e203000000000000000000000000000005f82015250565b5f610a67601283610a23565b9150610a7282610a33565b602082019050919050565b5f6020820190508181035f830152610a9481610a5b565b9050919050565b610aa481610932565b82525050565b5f602082019050610abd5f830184610a9b565b92915050565b7f496e76616c696420726563697069656e740000000000000000000000000000005f82015250565b5f610af7601183610a23565b9150610b0282610ac3565b602082019050919050565b5f6020820190508181035f830152610b2481610aeb565b9050919050565b5f81519050610b398161093b565b92915050565b5f60208284031215610b5457610b5361092e565b5b5f610b6184828501610b2b565b91505092915050565b7f496e73756666696369656e742062616c616e63650000000000000000000000005f82015250565b5f610b9e601483610a23565b9150610ba982610b6a565b602082019050919050565b5f6020820190508181035f830152610bcb81610b92565b905091905056fea26469706673582212202aaaba38229b9914839716e43f7c07a535df5552af5ecad51a19975b93a48fd364736f6c634300081e0033",
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
