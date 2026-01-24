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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_stablecoin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stablecoin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60a03461015d57601f61074e38819003918201601f19168301916001600160401b0383118484101761016157808492604094855283398101031261015d5761004681610175565b906001600160a01b039061005c90602001610175565b1690811561014a575f80546001600160a01b031981168417825560405193916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556001600160a01b031690811561010857506080526040516105c4908161018a823960805181818160ad0152818161031601526103690152f35b62461bcd60e51b815260206004820152601a60248201527f496e76616c696420737461626c65636f696e20616464726573730000000000006044820152606490fd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b038216820361015d5756fe60806040526004361015610011575f80fd5b5f3560e01c8063715018a6146104a35780638da5cb5b1461047c578063b6b55f2514610345578063e9cbd82214610301578063f2fde38b1461027c5763f3fef3a31461005b575f80fd5b34610219576040366003190112610219576100746104fa565b60243590610080610510565b610088610536565b6001600160a01b0316908115610243576040516370a0823160e01b81523060048201527f000000000000000000000000000000000000000000000000000000000000000091906001600160a01b03831690602081602481855afa908115610238575f916101d8575b50821161019c576040519263a9059cbb60e01b5f52846004528260245260205f60448180855af19060015f511482161561017d575b50836040521561016b5750816020917f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d59352a260015f51602061056f5f395f51905f5255005b635274afe760e01b5f5260045260245ffd5b6001821516610193573b15153d1516165f610125565b843d5f823e3d90fd5b60405162461bcd60e51b8152602060048201526014602482015273496e73756666696369656e742062616c616e636560601b6044820152606490fd5b905060203d602011610231575b601f8101601f1916820167ffffffffffffffff81118382101761021d5760209183916040528101031261021957515f6100f0565b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b503d6101e5565b6040513d5f823e3d90fd5b60405162461bcd60e51b8152602060048201526011602482015270125b9d985b1a59081c9958da5c1a595b9d607a1b6044820152606490fd5b34610219576020366003190112610219576102956104fa565b61029d610510565b6001600160a01b031680156102ee575f80546001600160a01b03198116831782556001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a3005b631e4fbdf760e01b5f525f60045260245ffd5b34610219575f366003190112610219576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b3461021957602036600319011261021957600435610361610536565b8015610442577f000000000000000000000000000000000000000000000000000000000000000090604051916323b872dd60e01b5f5233600452306024528160445260205f60648180855af160015f5114811615610423575b836040525f60605215610403575081527f2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c460203392a260015f51602061056f5f395f51905f5255005b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b600181151661043957813b15153d1516166103ba565b833d5f823e3d90fd5b60405162461bcd60e51b81526020600482015260126024820152710416d6f756e74206d757374206265203e20360741b6044820152606490fd5b34610219575f366003190112610219575f546040516001600160a01b039091168152602090f35b34610219575f366003190112610219576104bb610510565b5f80546001600160a01b0319811682556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b600435906001600160a01b038216820361021957565b5f546001600160a01b0316330361052357565b63118cdaa760e01b5f523360045260245ffd5b60025f51602061056f5f395f51905f52541461055f5760025f51602061056f5f395f51905f5255565b633ee5aeb560e01b5f5260045ffdfe9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00a264697066735822122053ed86fb91eca143184d068f7e3744208d11f56684a0b4eb80f97625744eb68b64736f6c63430008210033",
}

// PaymentVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymentVaultMetaData.ABI instead.
var PaymentVaultABI = PaymentVaultMetaData.ABI

// PaymentVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PaymentVaultMetaData.Bin instead.
var PaymentVaultBin = PaymentVaultMetaData.Bin

// DeployPaymentVault deploys a new Ethereum contract, binding an instance of PaymentVault to it.
func DeployPaymentVault(auth *bind.TransactOpts, backend bind.ContractBackend, _stablecoin common.Address, initialOwner common.Address) (common.Address, *types.Transaction, *PaymentVault, error) {
	parsed, err := PaymentVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymentVaultBin), backend, _stablecoin, initialOwner)
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

// Stablecoin is a free data retrieval call binding the contract method 0xe9cbd822.
//
// Solidity: function stablecoin() view returns(address)
func (_PaymentVault *PaymentVaultCaller) Stablecoin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PaymentVault.contract.Call(opts, &out, "stablecoin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Stablecoin is a free data retrieval call binding the contract method 0xe9cbd822.
//
// Solidity: function stablecoin() view returns(address)
func (_PaymentVault *PaymentVaultSession) Stablecoin() (common.Address, error) {
	return _PaymentVault.Contract.Stablecoin(&_PaymentVault.CallOpts)
}

// Stablecoin is a free data retrieval call binding the contract method 0xe9cbd822.
//
// Solidity: function stablecoin() view returns(address)
func (_PaymentVault *PaymentVaultCallerSession) Stablecoin() (common.Address, error) {
	return _PaymentVault.Contract.Stablecoin(&_PaymentVault.CallOpts)
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
