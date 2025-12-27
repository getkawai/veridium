// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package usdt

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

// MockUSDTMetaData contains all meta data concerning the MockUSDT contract.
var MockUSDTMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f5ffd5b506040518060400160405280600981526020017f4d6f636b205553445400000000000000000000000000000000000000000000008152506040518060400160405280600481526020017f5553445400000000000000000000000000000000000000000000000000000000815250816003908161008b91906105bd565b50806004908161009b91906105bd565b5050506100d5336100b06100da60201b60201c565b600a6100bc91906107f4565b620f42406100ca919061083e565b6100e260201b60201c565b610967565b5f6006905090565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610152575f6040517fec442f0500000000000000000000000000000000000000000000000000000000815260040161014991906108be565b60405180910390fd5b6101635f838361016760201b60201c565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036101b7578060025f8282546101ab91906108d7565b92505081905550610285565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610240578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161023793929190610919565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036102cc578060025f8282540392505081905550610316565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610373919061094e565b60405180910390a3505050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806103fb57607f821691505b60208210810361040e5761040d6103b7565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026104707fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610435565b61047a8683610435565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6104be6104b96104b484610492565b61049b565b610492565b9050919050565b5f819050919050565b6104d7836104a4565b6104eb6104e3826104c5565b848454610441565b825550505050565b5f5f905090565b6105026104f3565b61050d8184846104ce565b505050565b5b81811015610530576105255f826104fa565b600181019050610513565b5050565b601f8211156105755761054681610414565b61054f84610426565b8101602085101561055e578190505b61057261056a85610426565b830182610512565b50505b505050565b5f82821c905092915050565b5f6105955f198460080261057a565b1980831691505092915050565b5f6105ad8383610586565b9150826002028217905092915050565b6105c682610380565b67ffffffffffffffff8111156105df576105de61038a565b5b6105e982546103e4565b6105f4828285610534565b5f60209050601f831160018114610625575f8415610613578287015190505b61061d85826105a2565b865550610684565b601f19841661063386610414565b5f5b8281101561065a57848901518255600182019150602085019450602081019050610635565b868310156106775784890151610673601f891682610586565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f8160011c9050919050565b5f5f8291508390505b600185111561070e578086048111156106ea576106e961068c565b5b60018516156106f95780820291505b8081029050610707856106b9565b94506106ce565b94509492505050565b5f8261072657600190506107e1565b81610733575f90506107e1565b8160018114610749576002811461075357610782565b60019150506107e1565b60ff8411156107655761076461068c565b5b8360020a91508482111561077c5761077b61068c565b5b506107e1565b5060208310610133831016604e8410600b84101617156107b75782820a9050838111156107b2576107b161068c565b5b6107e1565b6107c484848460016106c5565b925090508184048111156107db576107da61068c565b5b81810290505b9392505050565b5f60ff82169050919050565b5f6107fe82610492565b9150610809836107e8565b92506108367fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8484610717565b905092915050565b5f61084882610492565b915061085383610492565b925082820261086181610492565b915082820484148315176108785761087761068c565b5b5092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6108a88261087f565b9050919050565b6108b88161089e565b82525050565b5f6020820190506108d15f8301846108af565b92915050565b5f6108e182610492565b91506108ec83610492565b92508282019050808211156109045761090361068c565b5b92915050565b61091381610492565b82525050565b5f60608201905061092c5f8301866108af565b610939602083018561090a565b610946604083018461090a565b949350505050565b5f6020820190506109615f83018461090a565b92915050565b610e96806109745f395ff3fe608060405234801561000f575f5ffd5b506004361061009c575f3560e01c806340c10f191161006457806340c10f191461015a57806370a082311461017657806395d89b41146101a6578063a9059cbb146101c4578063dd62ed3e146101f45761009c565b806306fdde03146100a0578063095ea7b3146100be57806318160ddd146100ee57806323b872dd1461010c578063313ce5671461013c575b5f5ffd5b6100a8610224565b6040516100b59190610b0f565b60405180910390f35b6100d860048036038101906100d39190610bc0565b6102b4565b6040516100e59190610c18565b60405180910390f35b6100f66102d6565b6040516101039190610c40565b60405180910390f35b61012660048036038101906101219190610c59565b6102df565b6040516101339190610c18565b60405180910390f35b61014461030d565b6040516101519190610cc4565b60405180910390f35b610174600480360381019061016f9190610bc0565b610315565b005b610190600480360381019061018b9190610cdd565b610323565b60405161019d9190610c40565b60405180910390f35b6101ae610368565b6040516101bb9190610b0f565b60405180910390f35b6101de60048036038101906101d99190610bc0565b6103f8565b6040516101eb9190610c18565b60405180910390f35b61020e60048036038101906102099190610d08565b61041a565b60405161021b9190610c40565b60405180910390f35b60606003805461023390610d73565b80601f016020809104026020016040519081016040528092919081815260200182805461025f90610d73565b80156102aa5780601f10610281576101008083540402835291602001916102aa565b820191905f5260205f20905b81548152906001019060200180831161028d57829003601f168201915b5050505050905090565b5f5f6102be61049c565b90506102cb8185856104a3565b600191505092915050565b5f600254905090565b5f5f6102e961049c565b90506102f68582856104b5565b610301858585610548565b60019150509392505050565b5f6006905090565b61031f8282610638565b5050565b5f5f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b60606004805461037790610d73565b80601f01602080910402602001604051908101604052809291908181526020018280546103a390610d73565b80156103ee5780601f106103c5576101008083540402835291602001916103ee565b820191905f5260205f20905b8154815290600101906020018083116103d157829003601f168201915b5050505050905090565b5f5f61040261049c565b905061040f818585610548565b600191505092915050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b5f33905090565b6104b083838360016106b7565b505050565b5f6104c0848461041a565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8110156105425781811015610533578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161052a93929190610db2565b60405180910390fd5b61054184848484035f6106b7565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036105b8575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016105af9190610de7565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610628575f6040517fec442f0500000000000000000000000000000000000000000000000000000000815260040161061f9190610de7565b60405180910390fd5b610633838383610886565b505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036106a8575f6040517fec442f0500000000000000000000000000000000000000000000000000000000815260040161069f9190610de7565b60405180910390fd5b6106b35f8383610886565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610727575f6040517fe602df0500000000000000000000000000000000000000000000000000000000815260040161071e9190610de7565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610797575f6040517f94280d6200000000000000000000000000000000000000000000000000000000815260040161078e9190610de7565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015610880578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516108779190610c40565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036108d6578060025f8282546108ca9190610e2d565b925050819055506109a4565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205490508181101561095f578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161095693929190610db2565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036109eb578060025f8282540392505081905550610a35565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610a929190610c40565b60405180910390a3505050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f610ae182610a9f565b610aeb8185610aa9565b9350610afb818560208601610ab9565b610b0481610ac7565b840191505092915050565b5f6020820190508181035f830152610b278184610ad7565b905092915050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610b5c82610b33565b9050919050565b610b6c81610b52565b8114610b76575f5ffd5b50565b5f81359050610b8781610b63565b92915050565b5f819050919050565b610b9f81610b8d565b8114610ba9575f5ffd5b50565b5f81359050610bba81610b96565b92915050565b5f5f60408385031215610bd657610bd5610b2f565b5b5f610be385828601610b79565b9250506020610bf485828601610bac565b9150509250929050565b5f8115159050919050565b610c1281610bfe565b82525050565b5f602082019050610c2b5f830184610c09565b92915050565b610c3a81610b8d565b82525050565b5f602082019050610c535f830184610c31565b92915050565b5f5f5f60608486031215610c7057610c6f610b2f565b5b5f610c7d86828701610b79565b9350506020610c8e86828701610b79565b9250506040610c9f86828701610bac565b9150509250925092565b5f60ff82169050919050565b610cbe81610ca9565b82525050565b5f602082019050610cd75f830184610cb5565b92915050565b5f60208284031215610cf257610cf1610b2f565b5b5f610cff84828501610b79565b91505092915050565b5f5f60408385031215610d1e57610d1d610b2f565b5b5f610d2b85828601610b79565b9250506020610d3c85828601610b79565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610d8a57607f821691505b602082108103610d9d57610d9c610d46565b5b50919050565b610dac81610b52565b82525050565b5f606082019050610dc55f830186610da3565b610dd26020830185610c31565b610ddf6040830184610c31565b949350505050565b5f602082019050610dfa5f830184610da3565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610e3782610b8d565b9150610e4283610b8d565b9250828201905080821115610e5a57610e59610e00565b5b9291505056fea2646970667358221220c9da84c91efe0a91b71fec5cf27395900a965122f7fbe4cec99be69376086e5064736f6c634300081e0033",
}

// MockUSDTABI is the input ABI used to generate the binding from.
// Deprecated: Use MockUSDTMetaData.ABI instead.
var MockUSDTABI = MockUSDTMetaData.ABI

// MockUSDTBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockUSDTMetaData.Bin instead.
var MockUSDTBin = MockUSDTMetaData.Bin

// DeployMockUSDT deploys a new Ethereum contract, binding an instance of MockUSDT to it.
func DeployMockUSDT(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MockUSDT, error) {
	parsed, err := MockUSDTMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockUSDTBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockUSDT{MockUSDTCaller: MockUSDTCaller{contract: contract}, MockUSDTTransactor: MockUSDTTransactor{contract: contract}, MockUSDTFilterer: MockUSDTFilterer{contract: contract}}, nil
}

// MockUSDT is an auto generated Go binding around an Ethereum contract.
type MockUSDT struct {
	MockUSDTCaller     // Read-only binding to the contract
	MockUSDTTransactor // Write-only binding to the contract
	MockUSDTFilterer   // Log filterer for contract events
}

// MockUSDTCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockUSDTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockUSDTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockUSDTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockUSDTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockUSDTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockUSDTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockUSDTSession struct {
	Contract     *MockUSDT         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockUSDTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockUSDTCallerSession struct {
	Contract *MockUSDTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MockUSDTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockUSDTTransactorSession struct {
	Contract     *MockUSDTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MockUSDTRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockUSDTRaw struct {
	Contract *MockUSDT // Generic contract binding to access the raw methods on
}

// MockUSDTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockUSDTCallerRaw struct {
	Contract *MockUSDTCaller // Generic read-only contract binding to access the raw methods on
}

// MockUSDTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockUSDTTransactorRaw struct {
	Contract *MockUSDTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockUSDT creates a new instance of MockUSDT, bound to a specific deployed contract.
func NewMockUSDT(address common.Address, backend bind.ContractBackend) (*MockUSDT, error) {
	contract, err := bindMockUSDT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockUSDT{MockUSDTCaller: MockUSDTCaller{contract: contract}, MockUSDTTransactor: MockUSDTTransactor{contract: contract}, MockUSDTFilterer: MockUSDTFilterer{contract: contract}}, nil
}

// NewMockUSDTCaller creates a new read-only instance of MockUSDT, bound to a specific deployed contract.
func NewMockUSDTCaller(address common.Address, caller bind.ContractCaller) (*MockUSDTCaller, error) {
	contract, err := bindMockUSDT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockUSDTCaller{contract: contract}, nil
}

// NewMockUSDTTransactor creates a new write-only instance of MockUSDT, bound to a specific deployed contract.
func NewMockUSDTTransactor(address common.Address, transactor bind.ContractTransactor) (*MockUSDTTransactor, error) {
	contract, err := bindMockUSDT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockUSDTTransactor{contract: contract}, nil
}

// NewMockUSDTFilterer creates a new log filterer instance of MockUSDT, bound to a specific deployed contract.
func NewMockUSDTFilterer(address common.Address, filterer bind.ContractFilterer) (*MockUSDTFilterer, error) {
	contract, err := bindMockUSDT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockUSDTFilterer{contract: contract}, nil
}

// bindMockUSDT binds a generic wrapper to an already deployed contract.
func bindMockUSDT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockUSDTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockUSDT *MockUSDTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockUSDT.Contract.MockUSDTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockUSDT *MockUSDTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUSDT.Contract.MockUSDTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockUSDT *MockUSDTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockUSDT.Contract.MockUSDTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockUSDT *MockUSDTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockUSDT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockUSDT *MockUSDTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockUSDT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockUSDT *MockUSDTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockUSDT.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockUSDT *MockUSDTCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockUSDT.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockUSDT *MockUSDTSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockUSDT.Contract.Allowance(&_MockUSDT.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockUSDT *MockUSDTCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockUSDT.Contract.Allowance(&_MockUSDT.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockUSDT *MockUSDTCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockUSDT.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockUSDT *MockUSDTSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockUSDT.Contract.BalanceOf(&_MockUSDT.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockUSDT *MockUSDTCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockUSDT.Contract.BalanceOf(&_MockUSDT.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockUSDT *MockUSDTCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockUSDT.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockUSDT *MockUSDTSession) Decimals() (uint8, error) {
	return _MockUSDT.Contract.Decimals(&_MockUSDT.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockUSDT *MockUSDTCallerSession) Decimals() (uint8, error) {
	return _MockUSDT.Contract.Decimals(&_MockUSDT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockUSDT *MockUSDTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockUSDT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockUSDT *MockUSDTSession) Name() (string, error) {
	return _MockUSDT.Contract.Name(&_MockUSDT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockUSDT *MockUSDTCallerSession) Name() (string, error) {
	return _MockUSDT.Contract.Name(&_MockUSDT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockUSDT *MockUSDTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockUSDT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockUSDT *MockUSDTSession) Symbol() (string, error) {
	return _MockUSDT.Contract.Symbol(&_MockUSDT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockUSDT *MockUSDTCallerSession) Symbol() (string, error) {
	return _MockUSDT.Contract.Symbol(&_MockUSDT.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockUSDT *MockUSDTCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockUSDT.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockUSDT *MockUSDTSession) TotalSupply() (*big.Int, error) {
	return _MockUSDT.Contract.TotalSupply(&_MockUSDT.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockUSDT *MockUSDTCallerSession) TotalSupply() (*big.Int, error) {
	return _MockUSDT.Contract.TotalSupply(&_MockUSDT.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.Approve(&_MockUSDT.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.Approve(&_MockUSDT.TransactOpts, spender, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_MockUSDT *MockUSDTTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MockUSDT.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_MockUSDT *MockUSDTSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.Mint(&_MockUSDT.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_MockUSDT *MockUSDTTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.Mint(&_MockUSDT.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.Transfer(&_MockUSDT.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.Transfer(&_MockUSDT.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.TransferFrom(&_MockUSDT.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockUSDT *MockUSDTTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockUSDT.Contract.TransferFrom(&_MockUSDT.TransactOpts, from, to, value)
}

// MockUSDTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the MockUSDT contract.
type MockUSDTApprovalIterator struct {
	Event *MockUSDTApproval // Event containing the contract specifics and raw log

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
func (it *MockUSDTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUSDTApproval)
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
		it.Event = new(MockUSDTApproval)
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
func (it *MockUSDTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUSDTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUSDTApproval represents a Approval event raised by the MockUSDT contract.
type MockUSDTApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockUSDT *MockUSDTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*MockUSDTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockUSDT.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &MockUSDTApprovalIterator{contract: _MockUSDT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockUSDT *MockUSDTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MockUSDTApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockUSDT.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUSDTApproval)
				if err := _MockUSDT.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockUSDT *MockUSDTFilterer) ParseApproval(log types.Log) (*MockUSDTApproval, error) {
	event := new(MockUSDTApproval)
	if err := _MockUSDT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockUSDTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the MockUSDT contract.
type MockUSDTTransferIterator struct {
	Event *MockUSDTTransfer // Event containing the contract specifics and raw log

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
func (it *MockUSDTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockUSDTTransfer)
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
		it.Event = new(MockUSDTTransfer)
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
func (it *MockUSDTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockUSDTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockUSDTTransfer represents a Transfer event raised by the MockUSDT contract.
type MockUSDTTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockUSDT *MockUSDTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockUSDTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockUSDT.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockUSDTTransferIterator{contract: _MockUSDT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockUSDT *MockUSDTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MockUSDTTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockUSDT.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockUSDTTransfer)
				if err := _MockUSDT.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockUSDT *MockUSDTFilterer) ParseTransfer(log types.Log) (*MockUSDTTransfer, error) {
	event := new(MockUSDTTransfer)
	if err := _MockUSDT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
