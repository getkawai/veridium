// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package revenuedistributor

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

// RevenueDistributorMetaData contains all meta data concerning the RevenueDistributor contract.
var RevenueDistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"token_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isClaimed\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60a0346100d057601f6106d738819003918201601f19168301916001600160401b038311848410176100d4578084926020946040528339810103126100d057516001600160a01b038116908190036100d05733156100bd575f8054336001600160a01b0319821681178355604051939290916001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a36080526105ee90816100e98239608051818181608601526103f20152f35b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c80632e7ba6ef146102715780632eb4a7ab14610254578063715018a6146101fd5780637cb64759146101aa5780638da5cb5b146101835780639e34070f14610142578063f2fde38b146100b95763fc0c546a14610071575f80fd5b346100b5575f3660031901126100b5576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b5f80fd5b346100b55760203660031901126100b5576004356001600160a01b038116908190036100b5576100e7610592565b801561012f575f80546001600160a01b03198116831782556001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a3005b631e4fbdf760e01b5f525f60045260245ffd5b346100b55760203660031901126100b55760206101796004358060081c5f526002602052600160ff60405f205492161b8091161490565b6040519015158152f35b346100b5575f3660031901126100b5575f546040516001600160a01b039091168152602090f35b346100b55760203660031901126100b5576004356101c6610592565b7ffd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f875760406001548151908152836020820152a1600155005b346100b5575f3660031901126100b557610215610592565b5f80546001600160a01b0319811682556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346100b5575f3660031901126100b5576020600154604051908152f35b346100b55760803660031901126100b5576024356004356001600160a01b0382168083036100b5576064359160443567ffffffffffffffff84116100b557366023850112156100b55783600401359367ffffffffffffffff85116100b5578460051b95602487830101903682116100b557610304858060081c5f526002602052600160ff60405f205492161b8091161490565b6105055760405160208101918683526bffffffffffffffffffffffff199060601b1660408201528460548201526054815261034060748261055c565b519020916001549661035860206040519a018a61055c565b8852602401602088015b8282106104f557505050925f935b86518510156103b25760208560051b88010151908181105f146103a1575f52602052600160405f205b940193610370565b905f52602052600160405f20610399565b85036104a557600883901c5f9081526002602090815260408083208054600160ff89161b1790555163a9059cbb60e01b83526004849052602485905294917f0000000000000000000000000000000000000000000000000000000000000000919060448180855af160015f5114811615610486575b8560405215610466577f4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed026606086868686835260208301526040820152a1005b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b600181151661049c57813b15153d151616610427565b853d5f823e3d90fd5b60405162461bcd60e51b815260206004820152602260248201527f526576656e75654469737472696275746f723a20496e76616c69642070726f6f604482015261331760f11b6064820152608490fd5b8135815260209182019101610362565b60405162461bcd60e51b815260206004820152602960248201527f526576656e75654469737472696275746f723a2044726f7020616c72656164796044820152681031b630b4b6b2b21760b91b6064820152608490fd5b90601f8019910116810190811067ffffffffffffffff82111761057e57604052565b634e487b7160e01b5f52604160045260245ffd5b5f546001600160a01b031633036105a557565b63118cdaa760e01b5f523360045260245ffdfea2646970667358221220039cf01bb4af4b68ecfb7e0f2db7c76dfe2001c3319adf89e4b22f17ab1385ea64736f6c63430008210033",
}

// RevenueDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use RevenueDistributorMetaData.ABI instead.
var RevenueDistributorABI = RevenueDistributorMetaData.ABI

// RevenueDistributorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RevenueDistributorMetaData.Bin instead.
var RevenueDistributorBin = RevenueDistributorMetaData.Bin

// DeployRevenueDistributor deploys a new Ethereum contract, binding an instance of RevenueDistributor to it.
func DeployRevenueDistributor(auth *bind.TransactOpts, backend bind.ContractBackend, token_ common.Address) (common.Address, *types.Transaction, *RevenueDistributor, error) {
	parsed, err := RevenueDistributorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RevenueDistributorBin), backend, token_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RevenueDistributor{RevenueDistributorCaller: RevenueDistributorCaller{contract: contract}, RevenueDistributorTransactor: RevenueDistributorTransactor{contract: contract}, RevenueDistributorFilterer: RevenueDistributorFilterer{contract: contract}}, nil
}

// RevenueDistributor is an auto generated Go binding around an Ethereum contract.
type RevenueDistributor struct {
	RevenueDistributorCaller     // Read-only binding to the contract
	RevenueDistributorTransactor // Write-only binding to the contract
	RevenueDistributorFilterer   // Log filterer for contract events
}

// RevenueDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type RevenueDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevenueDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RevenueDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevenueDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RevenueDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevenueDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RevenueDistributorSession struct {
	Contract     *RevenueDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RevenueDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RevenueDistributorCallerSession struct {
	Contract *RevenueDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// RevenueDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RevenueDistributorTransactorSession struct {
	Contract     *RevenueDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// RevenueDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type RevenueDistributorRaw struct {
	Contract *RevenueDistributor // Generic contract binding to access the raw methods on
}

// RevenueDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RevenueDistributorCallerRaw struct {
	Contract *RevenueDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// RevenueDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RevenueDistributorTransactorRaw struct {
	Contract *RevenueDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRevenueDistributor creates a new instance of RevenueDistributor, bound to a specific deployed contract.
func NewRevenueDistributor(address common.Address, backend bind.ContractBackend) (*RevenueDistributor, error) {
	contract, err := bindRevenueDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RevenueDistributor{RevenueDistributorCaller: RevenueDistributorCaller{contract: contract}, RevenueDistributorTransactor: RevenueDistributorTransactor{contract: contract}, RevenueDistributorFilterer: RevenueDistributorFilterer{contract: contract}}, nil
}

// NewRevenueDistributorCaller creates a new read-only instance of RevenueDistributor, bound to a specific deployed contract.
func NewRevenueDistributorCaller(address common.Address, caller bind.ContractCaller) (*RevenueDistributorCaller, error) {
	contract, err := bindRevenueDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RevenueDistributorCaller{contract: contract}, nil
}

// NewRevenueDistributorTransactor creates a new write-only instance of RevenueDistributor, bound to a specific deployed contract.
func NewRevenueDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*RevenueDistributorTransactor, error) {
	contract, err := bindRevenueDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RevenueDistributorTransactor{contract: contract}, nil
}

// NewRevenueDistributorFilterer creates a new log filterer instance of RevenueDistributor, bound to a specific deployed contract.
func NewRevenueDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*RevenueDistributorFilterer, error) {
	contract, err := bindRevenueDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RevenueDistributorFilterer{contract: contract}, nil
}

// bindRevenueDistributor binds a generic wrapper to an already deployed contract.
func bindRevenueDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RevenueDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RevenueDistributor *RevenueDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RevenueDistributor.Contract.RevenueDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RevenueDistributor *RevenueDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.RevenueDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RevenueDistributor *RevenueDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.RevenueDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RevenueDistributor *RevenueDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RevenueDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RevenueDistributor *RevenueDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RevenueDistributor *RevenueDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.contract.Transact(opts, method, params...)
}

// IsClaimed is a free data retrieval call binding the contract method 0x9e34070f.
//
// Solidity: function isClaimed(uint256 index) view returns(bool)
func (_RevenueDistributor *RevenueDistributorCaller) IsClaimed(opts *bind.CallOpts, index *big.Int) (bool, error) {
	var out []interface{}
	err := _RevenueDistributor.contract.Call(opts, &out, "isClaimed", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsClaimed is a free data retrieval call binding the contract method 0x9e34070f.
//
// Solidity: function isClaimed(uint256 index) view returns(bool)
func (_RevenueDistributor *RevenueDistributorSession) IsClaimed(index *big.Int) (bool, error) {
	return _RevenueDistributor.Contract.IsClaimed(&_RevenueDistributor.CallOpts, index)
}

// IsClaimed is a free data retrieval call binding the contract method 0x9e34070f.
//
// Solidity: function isClaimed(uint256 index) view returns(bool)
func (_RevenueDistributor *RevenueDistributorCallerSession) IsClaimed(index *big.Int) (bool, error) {
	return _RevenueDistributor.Contract.IsClaimed(&_RevenueDistributor.CallOpts, index)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_RevenueDistributor *RevenueDistributorCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RevenueDistributor.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_RevenueDistributor *RevenueDistributorSession) MerkleRoot() ([32]byte, error) {
	return _RevenueDistributor.Contract.MerkleRoot(&_RevenueDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_RevenueDistributor *RevenueDistributorCallerSession) MerkleRoot() ([32]byte, error) {
	return _RevenueDistributor.Contract.MerkleRoot(&_RevenueDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RevenueDistributor *RevenueDistributorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RevenueDistributor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RevenueDistributor *RevenueDistributorSession) Owner() (common.Address, error) {
	return _RevenueDistributor.Contract.Owner(&_RevenueDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RevenueDistributor *RevenueDistributorCallerSession) Owner() (common.Address, error) {
	return _RevenueDistributor.Contract.Owner(&_RevenueDistributor.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_RevenueDistributor *RevenueDistributorCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RevenueDistributor.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_RevenueDistributor *RevenueDistributorSession) Token() (common.Address, error) {
	return _RevenueDistributor.Contract.Token(&_RevenueDistributor.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_RevenueDistributor *RevenueDistributorCallerSession) Token() (common.Address, error) {
	return _RevenueDistributor.Contract.Token(&_RevenueDistributor.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x2e7ba6ef.
//
// Solidity: function claim(uint256 index, address account, uint256 amount, bytes32[] merkleProof) returns()
func (_RevenueDistributor *RevenueDistributorTransactor) Claim(opts *bind.TransactOpts, index *big.Int, account common.Address, amount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _RevenueDistributor.contract.Transact(opts, "claim", index, account, amount, merkleProof)
}

// Claim is a paid mutator transaction binding the contract method 0x2e7ba6ef.
//
// Solidity: function claim(uint256 index, address account, uint256 amount, bytes32[] merkleProof) returns()
func (_RevenueDistributor *RevenueDistributorSession) Claim(index *big.Int, account common.Address, amount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.Claim(&_RevenueDistributor.TransactOpts, index, account, amount, merkleProof)
}

// Claim is a paid mutator transaction binding the contract method 0x2e7ba6ef.
//
// Solidity: function claim(uint256 index, address account, uint256 amount, bytes32[] merkleProof) returns()
func (_RevenueDistributor *RevenueDistributorTransactorSession) Claim(index *big.Int, account common.Address, amount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.Claim(&_RevenueDistributor.TransactOpts, index, account, amount, merkleProof)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RevenueDistributor *RevenueDistributorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RevenueDistributor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RevenueDistributor *RevenueDistributorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RevenueDistributor.Contract.RenounceOwnership(&_RevenueDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RevenueDistributor *RevenueDistributorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RevenueDistributor.Contract.RenounceOwnership(&_RevenueDistributor.TransactOpts)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_RevenueDistributor *RevenueDistributorTransactor) SetMerkleRoot(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _RevenueDistributor.contract.Transact(opts, "setMerkleRoot", _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_RevenueDistributor *RevenueDistributorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.SetMerkleRoot(&_RevenueDistributor.TransactOpts, _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_RevenueDistributor *RevenueDistributorTransactorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.SetMerkleRoot(&_RevenueDistributor.TransactOpts, _merkleRoot)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RevenueDistributor *RevenueDistributorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RevenueDistributor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RevenueDistributor *RevenueDistributorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.TransferOwnership(&_RevenueDistributor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RevenueDistributor *RevenueDistributorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RevenueDistributor.Contract.TransferOwnership(&_RevenueDistributor.TransactOpts, newOwner)
}

// RevenueDistributorClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the RevenueDistributor contract.
type RevenueDistributorClaimedIterator struct {
	Event *RevenueDistributorClaimed // Event containing the contract specifics and raw log

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
func (it *RevenueDistributorClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RevenueDistributorClaimed)
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
		it.Event = new(RevenueDistributorClaimed)
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
func (it *RevenueDistributorClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RevenueDistributorClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RevenueDistributorClaimed represents a Claimed event raised by the RevenueDistributor contract.
type RevenueDistributorClaimed struct {
	Index   *big.Int
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed026.
//
// Solidity: event Claimed(uint256 index, address account, uint256 amount)
func (_RevenueDistributor *RevenueDistributorFilterer) FilterClaimed(opts *bind.FilterOpts) (*RevenueDistributorClaimedIterator, error) {

	logs, sub, err := _RevenueDistributor.contract.FilterLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return &RevenueDistributorClaimedIterator{contract: _RevenueDistributor.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed026.
//
// Solidity: event Claimed(uint256 index, address account, uint256 amount)
func (_RevenueDistributor *RevenueDistributorFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *RevenueDistributorClaimed) (event.Subscription, error) {

	logs, sub, err := _RevenueDistributor.contract.WatchLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RevenueDistributorClaimed)
				if err := _RevenueDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0x4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed026.
//
// Solidity: event Claimed(uint256 index, address account, uint256 amount)
func (_RevenueDistributor *RevenueDistributorFilterer) ParseClaimed(log types.Log) (*RevenueDistributorClaimed, error) {
	event := new(RevenueDistributorClaimed)
	if err := _RevenueDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RevenueDistributorMerkleRootUpdatedIterator is returned from FilterMerkleRootUpdated and is used to iterate over the raw logs and unpacked data for MerkleRootUpdated events raised by the RevenueDistributor contract.
type RevenueDistributorMerkleRootUpdatedIterator struct {
	Event *RevenueDistributorMerkleRootUpdated // Event containing the contract specifics and raw log

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
func (it *RevenueDistributorMerkleRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RevenueDistributorMerkleRootUpdated)
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
		it.Event = new(RevenueDistributorMerkleRootUpdated)
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
func (it *RevenueDistributorMerkleRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RevenueDistributorMerkleRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RevenueDistributorMerkleRootUpdated represents a MerkleRootUpdated event raised by the RevenueDistributor contract.
type RevenueDistributorMerkleRootUpdated struct {
	OldRoot [32]byte
	NewRoot [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootUpdated is a free log retrieval operation binding the contract event 0xfd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f8757.
//
// Solidity: event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot)
func (_RevenueDistributor *RevenueDistributorFilterer) FilterMerkleRootUpdated(opts *bind.FilterOpts) (*RevenueDistributorMerkleRootUpdatedIterator, error) {

	logs, sub, err := _RevenueDistributor.contract.FilterLogs(opts, "MerkleRootUpdated")
	if err != nil {
		return nil, err
	}
	return &RevenueDistributorMerkleRootUpdatedIterator{contract: _RevenueDistributor.contract, event: "MerkleRootUpdated", logs: logs, sub: sub}, nil
}

// WatchMerkleRootUpdated is a free log subscription operation binding the contract event 0xfd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f8757.
//
// Solidity: event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot)
func (_RevenueDistributor *RevenueDistributorFilterer) WatchMerkleRootUpdated(opts *bind.WatchOpts, sink chan<- *RevenueDistributorMerkleRootUpdated) (event.Subscription, error) {

	logs, sub, err := _RevenueDistributor.contract.WatchLogs(opts, "MerkleRootUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RevenueDistributorMerkleRootUpdated)
				if err := _RevenueDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
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

// ParseMerkleRootUpdated is a log parse operation binding the contract event 0xfd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f8757.
//
// Solidity: event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot)
func (_RevenueDistributor *RevenueDistributorFilterer) ParseMerkleRootUpdated(log types.Log) (*RevenueDistributorMerkleRootUpdated, error) {
	event := new(RevenueDistributorMerkleRootUpdated)
	if err := _RevenueDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RevenueDistributorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RevenueDistributor contract.
type RevenueDistributorOwnershipTransferredIterator struct {
	Event *RevenueDistributorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RevenueDistributorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RevenueDistributorOwnershipTransferred)
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
		it.Event = new(RevenueDistributorOwnershipTransferred)
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
func (it *RevenueDistributorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RevenueDistributorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RevenueDistributorOwnershipTransferred represents a OwnershipTransferred event raised by the RevenueDistributor contract.
type RevenueDistributorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RevenueDistributor *RevenueDistributorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RevenueDistributorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RevenueDistributor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RevenueDistributorOwnershipTransferredIterator{contract: _RevenueDistributor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RevenueDistributor *RevenueDistributorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RevenueDistributorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RevenueDistributor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RevenueDistributorOwnershipTransferred)
				if err := _RevenueDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_RevenueDistributor *RevenueDistributorFilterer) ParseOwnershipTransferred(log types.Log) (*RevenueDistributorOwnershipTransferred, error) {
	event := new(RevenueDistributorOwnershipTransferred)
	if err := _RevenueDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
