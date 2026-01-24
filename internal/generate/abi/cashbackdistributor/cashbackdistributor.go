// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cashbackdistributor

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

// DepositCashbackDistributorMetaData contains all meta data concerning the DepositCashbackDistributor contract.
var DepositCashbackDistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"kawaiToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"advancePeriod\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimCashback\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimMultiplePeriods\",\"inputs\":[{\"name\":\"periods\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"kawaiAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"merkleProofs\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeriodMerkleRoot\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStats\",\"inputs\":[],\"outputs\":[{\"name\":\"_currentPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalKawaiDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalUsers\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedAnyPeriod\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasUserClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kawaiToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"periodMerkleRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeriodMerkleRoot\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalKawaiDistributed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalUsers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CashbackClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodAdvanced\",\"inputs\":[{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a03461014b57601f6110ad38819003918201601f19168301916001600160401b0383118484101761014f5780849260209460405283398101031261014b57516001600160a01b0381169081900361014b573315610138575f8054336001600160a01b03198216811783556040519290916001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005581156100f6575060805260018055604051610f49908161016482396080518181816101c30152818161054801526107b30152f35b62461bcd60e51b815260206004820152601560248201527f496e76616c6964204b41574149206164647265737300000000000000000000006044820152606490fd5b631e4fbdf760e01b5f525f60045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6080806040526004361015610012575f80fd5b5f905f3560e01c9081630604061814610bd45750806307c7a72d1461090e5780630ae6540314610b685780632eb4a7ab14610b4b5780633f4ba83a14610adf5780634a03da8a14610a1d5780635869bc5a14610ac25780635c975abb14610a9e578063715018a614610a47578063727a7c5a14610a1d5780637cb64759146109b75780638456cb5914610957578063873f6f9e1461090e5780638a90e20f146106a35780638da5cb5b1461067c578063adeacbd31461063d578063bff1f9e11461061f578063c40c91bd146105a6578063c59d484714610577578063cb56cd4f14610532578063f2fde38b146104ac5763f75cc2b914610110575f80fd5b3461043c57606036600319011261043c5760043567ffffffffffffffff811161029f57610141903690600401610c1a565b9060243567ffffffffffffffff811161043f57610162903690600401610c1a565b9060443567ffffffffffffffff81116104a857610183903690600401610c1a565b9261018c610e4d565b610194610e85565b8086148061049f575b15610462579493929190869587955b8087106102ae5788886101c0811515610caa565b817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b1561029f576040516340c10f1960e01b8152336004820152602481018490529082908290604490829084905af180156102a35761028a575b505061023390600654610df6565b600655338152600560205260ff60408220541615610260575b60015f516020610ef45f395f51905f525580f35b338152600560205260408120600160ff19825416179055610282600754610c4b565b60075561024c565b8161029491610ced565b61029f578183610225565b5080fd5b6040513d84823e3d90fd5b909192939495966102c0888389610e03565b35976102cd818588610e03565b35898b52600460205260408b2060018060a01b0333165f5260205260ff60405f205416610457576103026001548b1115610c6d565b61030d811515610caa565b60408051602081018c81523360601b6bffffffffffffffffffffffff1916928201929092526054810183905261035081607481015b03601f198101835282610ced565b5190208a8c52600360205260408c20549361036c851515610d23565b8984101561044357908c918460051b890135601e198a36030181121561043f5789019567ffffffffffffffff87351161043f576020870193873560051b3603851361043c578e6103da6103d588966040956103d060019d6104029b36913590610d64565b610ea2565b610dba565b8152600460205220868060a01b0333165f5260205260405f208660ff19825416179055610df6565b996040519182527f81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f60203393a35b019594939291906101ac565b80fd5b8380fd5b634e487b7160e01b8d52603260045260248dfd5b506001919850610430565b60405162461bcd60e51b8152602060048201526015602482015274082e4e4c2f240d8cadccee8d040dad2e6dac2e8c6d605b1b6044820152606490fd5b5083861461019d565b8580fd5b503461043c57602036600319011261043c576104c6610c04565b6104ce610e27565b6001600160a01b0316801561051e5781546001600160a01b03198116821783556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b631e4fbdf760e01b82526004829052602482fd5b503461043c578060031936011261043c576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b503461043c578060031936011261043c5760606001546006546007549060405192835260208301526040820152f35b503461043c57604036600319011261043c576024356004356105c6610e27565b6105d4600154821115610c6d565b8083526003602052807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c6040808620548151908152856020820152a282526003602052604082205580f35b503461043c578060031936011261043c576020600754604051908152f35b503461043c57602036600319011261043c5760209060ff906040906001600160a01b03610668610c04565b168152600584522054166040519015158152f35b503461043c578060031936011261043c57546040516001600160a01b039091168152602090f35b50346108c55760603660031901126108c55760243560043560443567ffffffffffffffff81116108c5576106db903690600401610c1a565b6106e3610e4d565b6106eb610e85565b6106f9600154841115610c6d565b5f83815260046020908152604080832033845290915290205460ff166108c95761078e916103d59161072c861515610caa565b6103d0604051602081019061076a816103428b338c87916054939183526bffffffffffffffffffffffff199060601b16602083015260348201520190565b51902092865f52600360205260405f205492610787841515610d23565b3691610d64565b5f8181526004602090815260408083203384529091529020805460ff191660011790557f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156108c5576040516340c10f1960e01b815233600482015260248101849052905f908290604490829084905af180156108ba576108a5575b5061082282600654610df6565b600655338352600560205260ff6040842054161561087b575b6040519182527f81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f60203393a360015f516020610ef45f395f51905f525580f35b338352600560205260408320600160ff1982541617905561089d600754610c4b565b60075561083b565b6108b29193505f90610ced565b5f915f610815565b6040513d5f823e3d90fd5b5f80fd5b60405162461bcd60e51b815260206004820152601f60248201527f416c726561647920636c61696d656420666f72207468697320706572696f64006044820152606490fd5b346108c55760403660031901126108c557610927610bee565b6004355f52600460205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b346108c5575f3660031901126108c55761096f610e27565b610977610e85565b5f805460ff60a01b1916600160a01b1790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602090a1005b346108c55760203660031901126108c5576004356109d3610e27565b600154807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c60406002548151908152856020820152a2816002555f52600360205260405f20555f80f35b346108c55760203660031901126108c5576004355f526003602052602060405f2054604051908152f35b346108c5575f3660031901126108c557610a5f610e27565b5f80546001600160a01b0319811682556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346108c5575f3660031901126108c557602060ff5f5460a01c166040519015158152f35b346108c5575f3660031901126108c5576020600654604051908152f35b346108c5575f3660031901126108c557610af7610e27565b5f5460ff8160a01c1615610b3c5760ff60a01b19165f556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa90602090a1005b638dfc202b60e01b5f5260045ffd5b346108c5575f3660031901126108c5576020600254604051908152f35b346108c55760203660031901126108c557600435610b84610e27565b610b8f600154610c4b565b80600155816002555f5260036020528060405f20557fa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf6602060015492604051908152a2005b346108c5575f3660031901126108c5576020906001548152f35b602435906001600160a01b03821682036108c557565b600435906001600160a01b03821682036108c557565b9181601f840112156108c55782359167ffffffffffffffff83116108c5576020808501948460051b0101116108c557565b5f198114610c595760010190565b634e487b7160e01b5f52601160045260245ffd5b15610c7457565b60405162461bcd60e51b815260206004820152600e60248201526d125b9d985b1a59081c195c9a5bd960921b6044820152606490fd5b15610cb157565b60405162461bcd60e51b81526020600482015260146024820152734e6f20636173686261636b20746f20636c61696d60601b6044820152606490fd5b90601f8019910116810190811067ffffffffffffffff821117610d0f57604052565b634e487b7160e01b5f52604160045260245ffd5b15610d2a57565b60405162461bcd60e51b815260206004820152601260248201527114195c9a5bd9081b9bdd081cd95d1d1b195960721b6044820152606490fd5b9291909267ffffffffffffffff8411610d0f578360051b906020604051610d8d82850182610ced565b80968152019181019283116108c557905b828210610daa57505050565b8135815260209182019101610d9e565b15610dc157565b60405162461bcd60e51b815260206004820152600d60248201526c24b73b30b634b210383937b7b360991b6044820152606490fd5b91908201809211610c5957565b9190811015610e135760051b0190565b634e487b7160e01b5f52603260045260245ffd5b5f546001600160a01b03163303610e3a57565b63118cdaa760e01b5f523360045260245ffd5b60025f516020610ef45f395f51905f525414610e765760025f516020610ef45f395f51905f5255565b633ee5aeb560e01b5f5260045ffd5b60ff5f5460a01c16610e9357565b63d93c066560e01b5f5260045ffd5b929091905f915b8451831015610eeb5760208360051b86010151908181105f14610eda575f52602052600160405f205b920191610ea9565b905f52602052600160405f20610ed2565b91509250149056fe9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00a2646970667358221220ed062cf2b900546adc91e21c14a0216de96bc6300065289717e5e26fa1aee3ee64736f6c63430008210033",
}

// DepositCashbackDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositCashbackDistributorMetaData.ABI instead.
var DepositCashbackDistributorABI = DepositCashbackDistributorMetaData.ABI

// DepositCashbackDistributorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepositCashbackDistributorMetaData.Bin instead.
var DepositCashbackDistributorBin = DepositCashbackDistributorMetaData.Bin

// DeployDepositCashbackDistributor deploys a new Ethereum contract, binding an instance of DepositCashbackDistributor to it.
func DeployDepositCashbackDistributor(auth *bind.TransactOpts, backend bind.ContractBackend, kawaiToken_ common.Address) (common.Address, *types.Transaction, *DepositCashbackDistributor, error) {
	parsed, err := DepositCashbackDistributorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepositCashbackDistributorBin), backend, kawaiToken_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DepositCashbackDistributor{DepositCashbackDistributorCaller: DepositCashbackDistributorCaller{contract: contract}, DepositCashbackDistributorTransactor: DepositCashbackDistributorTransactor{contract: contract}, DepositCashbackDistributorFilterer: DepositCashbackDistributorFilterer{contract: contract}}, nil
}

// DepositCashbackDistributor is an auto generated Go binding around an Ethereum contract.
type DepositCashbackDistributor struct {
	DepositCashbackDistributorCaller     // Read-only binding to the contract
	DepositCashbackDistributorTransactor // Write-only binding to the contract
	DepositCashbackDistributorFilterer   // Log filterer for contract events
}

// DepositCashbackDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositCashbackDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositCashbackDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositCashbackDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositCashbackDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositCashbackDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositCashbackDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositCashbackDistributorSession struct {
	Contract     *DepositCashbackDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DepositCashbackDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositCashbackDistributorCallerSession struct {
	Contract *DepositCashbackDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// DepositCashbackDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositCashbackDistributorTransactorSession struct {
	Contract     *DepositCashbackDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// DepositCashbackDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositCashbackDistributorRaw struct {
	Contract *DepositCashbackDistributor // Generic contract binding to access the raw methods on
}

// DepositCashbackDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositCashbackDistributorCallerRaw struct {
	Contract *DepositCashbackDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// DepositCashbackDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositCashbackDistributorTransactorRaw struct {
	Contract *DepositCashbackDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositCashbackDistributor creates a new instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributor(address common.Address, backend bind.ContractBackend) (*DepositCashbackDistributor, error) {
	contract, err := bindDepositCashbackDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributor{DepositCashbackDistributorCaller: DepositCashbackDistributorCaller{contract: contract}, DepositCashbackDistributorTransactor: DepositCashbackDistributorTransactor{contract: contract}, DepositCashbackDistributorFilterer: DepositCashbackDistributorFilterer{contract: contract}}, nil
}

// NewDepositCashbackDistributorCaller creates a new read-only instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributorCaller(address common.Address, caller bind.ContractCaller) (*DepositCashbackDistributorCaller, error) {
	contract, err := bindDepositCashbackDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorCaller{contract: contract}, nil
}

// NewDepositCashbackDistributorTransactor creates a new write-only instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositCashbackDistributorTransactor, error) {
	contract, err := bindDepositCashbackDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorTransactor{contract: contract}, nil
}

// NewDepositCashbackDistributorFilterer creates a new log filterer instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositCashbackDistributorFilterer, error) {
	contract, err := bindDepositCashbackDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorFilterer{contract: contract}, nil
}

// bindDepositCashbackDistributor binds a generic wrapper to an already deployed contract.
func bindDepositCashbackDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositCashbackDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositCashbackDistributor *DepositCashbackDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositCashbackDistributor.Contract.DepositCashbackDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositCashbackDistributor *DepositCashbackDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.DepositCashbackDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositCashbackDistributor *DepositCashbackDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.DepositCashbackDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositCashbackDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.contract.Transact(opts, method, params...)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) CurrentPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "currentPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) CurrentPeriod() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.CurrentPeriod(&_DepositCashbackDistributor.CallOpts)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) CurrentPeriod() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.CurrentPeriod(&_DepositCashbackDistributor.CallOpts)
}

// GetPeriodMerkleRoot is a free data retrieval call binding the contract method 0x4a03da8a.
//
// Solidity: function getPeriodMerkleRoot(uint256 period) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) GetPeriodMerkleRoot(opts *bind.CallOpts, period *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "getPeriodMerkleRoot", period)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPeriodMerkleRoot is a free data retrieval call binding the contract method 0x4a03da8a.
//
// Solidity: function getPeriodMerkleRoot(uint256 period) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) GetPeriodMerkleRoot(period *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.GetPeriodMerkleRoot(&_DepositCashbackDistributor.CallOpts, period)
}

// GetPeriodMerkleRoot is a free data retrieval call binding the contract method 0x4a03da8a.
//
// Solidity: function getPeriodMerkleRoot(uint256 period) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) GetPeriodMerkleRoot(period *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.GetPeriodMerkleRoot(&_DepositCashbackDistributor.CallOpts, period)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 _currentPeriod, uint256 _totalKawaiDistributed, uint256 _totalUsers)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) GetStats(opts *bind.CallOpts) (struct {
	CurrentPeriod         *big.Int
	TotalKawaiDistributed *big.Int
	TotalUsers            *big.Int
}, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "getStats")

	outstruct := new(struct {
		CurrentPeriod         *big.Int
		TotalKawaiDistributed *big.Int
		TotalUsers            *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CurrentPeriod = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalKawaiDistributed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalUsers = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 _currentPeriod, uint256 _totalKawaiDistributed, uint256 _totalUsers)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) GetStats() (struct {
	CurrentPeriod         *big.Int
	TotalKawaiDistributed *big.Int
	TotalUsers            *big.Int
}, error) {
	return _DepositCashbackDistributor.Contract.GetStats(&_DepositCashbackDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 _currentPeriod, uint256 _totalKawaiDistributed, uint256 _totalUsers)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) GetStats() (struct {
	CurrentPeriod         *big.Int
	TotalKawaiDistributed *big.Int
	TotalUsers            *big.Int
}, error) {
	return _DepositCashbackDistributor.Contract.GetStats(&_DepositCashbackDistributor.CallOpts)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) HasClaimed(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "hasClaimed", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimed(&_DepositCashbackDistributor.CallOpts, arg0, arg1)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimed(&_DepositCashbackDistributor.CallOpts, arg0, arg1)
}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) HasClaimedAnyPeriod(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "hasClaimedAnyPeriod", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) HasClaimedAnyPeriod(arg0 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimedAnyPeriod(&_DepositCashbackDistributor.CallOpts, arg0)
}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) HasClaimedAnyPeriod(arg0 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimedAnyPeriod(&_DepositCashbackDistributor.CallOpts, arg0)
}

// HasUserClaimed is a free data retrieval call binding the contract method 0x07c7a72d.
//
// Solidity: function hasUserClaimed(uint256 period, address user) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) HasUserClaimed(opts *bind.CallOpts, period *big.Int, user common.Address) (bool, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "hasUserClaimed", period, user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasUserClaimed is a free data retrieval call binding the contract method 0x07c7a72d.
//
// Solidity: function hasUserClaimed(uint256 period, address user) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) HasUserClaimed(period *big.Int, user common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasUserClaimed(&_DepositCashbackDistributor.CallOpts, period, user)
}

// HasUserClaimed is a free data retrieval call binding the contract method 0x07c7a72d.
//
// Solidity: function hasUserClaimed(uint256 period, address user) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) HasUserClaimed(period *big.Int, user common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasUserClaimed(&_DepositCashbackDistributor.CallOpts, period, user)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) KawaiToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "kawaiToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) KawaiToken() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.KawaiToken(&_DepositCashbackDistributor.CallOpts)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) KawaiToken() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.KawaiToken(&_DepositCashbackDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) MerkleRoot() ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.MerkleRoot(&_DepositCashbackDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) MerkleRoot() ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.MerkleRoot(&_DepositCashbackDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) Owner() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.Owner(&_DepositCashbackDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) Owner() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.Owner(&_DepositCashbackDistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) Paused() (bool, error) {
	return _DepositCashbackDistributor.Contract.Paused(&_DepositCashbackDistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) Paused() (bool, error) {
	return _DepositCashbackDistributor.Contract.Paused(&_DepositCashbackDistributor.CallOpts)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) PeriodMerkleRoots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "periodMerkleRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.PeriodMerkleRoots(&_DepositCashbackDistributor.CallOpts, arg0)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.PeriodMerkleRoots(&_DepositCashbackDistributor.CallOpts, arg0)
}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) TotalKawaiDistributed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "totalKawaiDistributed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) TotalKawaiDistributed() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalKawaiDistributed(&_DepositCashbackDistributor.CallOpts)
}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) TotalKawaiDistributed() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalKawaiDistributed(&_DepositCashbackDistributor.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) TotalUsers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "totalUsers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) TotalUsers() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalUsers(&_DepositCashbackDistributor.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) TotalUsers() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalUsers(&_DepositCashbackDistributor.CallOpts)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) AdvancePeriod(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "advancePeriod", _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.AdvancePeriod(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.AdvancePeriod(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// ClaimCashback is a paid mutator transaction binding the contract method 0x8a90e20f.
//
// Solidity: function claimCashback(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) ClaimCashback(opts *bind.TransactOpts, period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "claimCashback", period, kawaiAmount, merkleProof)
}

// ClaimCashback is a paid mutator transaction binding the contract method 0x8a90e20f.
//
// Solidity: function claimCashback(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) ClaimCashback(period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimCashback(&_DepositCashbackDistributor.TransactOpts, period, kawaiAmount, merkleProof)
}

// ClaimCashback is a paid mutator transaction binding the contract method 0x8a90e20f.
//
// Solidity: function claimCashback(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) ClaimCashback(period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimCashback(&_DepositCashbackDistributor.TransactOpts, period, kawaiAmount, merkleProof)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) ClaimMultiplePeriods(opts *bind.TransactOpts, periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "claimMultiplePeriods", periods, kawaiAmounts, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) ClaimMultiplePeriods(periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimMultiplePeriods(&_DepositCashbackDistributor.TransactOpts, periods, kawaiAmounts, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) ClaimMultiplePeriods(periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimMultiplePeriods(&_DepositCashbackDistributor.TransactOpts, periods, kawaiAmounts, merkleProofs)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) Pause() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.Pause(&_DepositCashbackDistributor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) Pause() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.Pause(&_DepositCashbackDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.RenounceOwnership(&_DepositCashbackDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.RenounceOwnership(&_DepositCashbackDistributor.TransactOpts)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) SetMerkleRoot(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "setMerkleRoot", _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetMerkleRoot(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetMerkleRoot(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) SetPeriodMerkleRoot(opts *bind.TransactOpts, period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "setPeriodMerkleRoot", period, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) SetPeriodMerkleRoot(period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetPeriodMerkleRoot(&_DepositCashbackDistributor.TransactOpts, period, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) SetPeriodMerkleRoot(period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetPeriodMerkleRoot(&_DepositCashbackDistributor.TransactOpts, period, _merkleRoot)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.TransferOwnership(&_DepositCashbackDistributor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.TransferOwnership(&_DepositCashbackDistributor.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) Unpause() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.Unpause(&_DepositCashbackDistributor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) Unpause() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.Unpause(&_DepositCashbackDistributor.TransactOpts)
}

// DepositCashbackDistributorCashbackClaimedIterator is returned from FilterCashbackClaimed and is used to iterate over the raw logs and unpacked data for CashbackClaimed events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorCashbackClaimedIterator struct {
	Event *DepositCashbackDistributorCashbackClaimed // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorCashbackClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorCashbackClaimed)
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
		it.Event = new(DepositCashbackDistributorCashbackClaimed)
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
func (it *DepositCashbackDistributorCashbackClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorCashbackClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorCashbackClaimed represents a CashbackClaimed event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorCashbackClaimed struct {
	Period      *big.Int
	User        common.Address
	KawaiAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCashbackClaimed is a free log retrieval operation binding the contract event 0x81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f.
//
// Solidity: event CashbackClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterCashbackClaimed(opts *bind.FilterOpts, period []*big.Int, user []common.Address) (*DepositCashbackDistributorCashbackClaimedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "CashbackClaimed", periodRule, userRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorCashbackClaimedIterator{contract: _DepositCashbackDistributor.contract, event: "CashbackClaimed", logs: logs, sub: sub}, nil
}

// WatchCashbackClaimed is a free log subscription operation binding the contract event 0x81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f.
//
// Solidity: event CashbackClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchCashbackClaimed(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorCashbackClaimed, period []*big.Int, user []common.Address) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "CashbackClaimed", periodRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorCashbackClaimed)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "CashbackClaimed", log); err != nil {
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

// ParseCashbackClaimed is a log parse operation binding the contract event 0x81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f.
//
// Solidity: event CashbackClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParseCashbackClaimed(log types.Log) (*DepositCashbackDistributorCashbackClaimed, error) {
	event := new(DepositCashbackDistributorCashbackClaimed)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "CashbackClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorMerkleRootUpdatedIterator is returned from FilterMerkleRootUpdated and is used to iterate over the raw logs and unpacked data for MerkleRootUpdated events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorMerkleRootUpdatedIterator struct {
	Event *DepositCashbackDistributorMerkleRootUpdated // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorMerkleRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorMerkleRootUpdated)
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
		it.Event = new(DepositCashbackDistributorMerkleRootUpdated)
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
func (it *DepositCashbackDistributorMerkleRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorMerkleRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorMerkleRootUpdated represents a MerkleRootUpdated event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorMerkleRootUpdated struct {
	Period  *big.Int
	OldRoot [32]byte
	NewRoot [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootUpdated is a free log retrieval operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterMerkleRootUpdated(opts *bind.FilterOpts, period []*big.Int) (*DepositCashbackDistributorMerkleRootUpdatedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorMerkleRootUpdatedIterator{contract: _DepositCashbackDistributor.contract, event: "MerkleRootUpdated", logs: logs, sub: sub}, nil
}

// WatchMerkleRootUpdated is a free log subscription operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchMerkleRootUpdated(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorMerkleRootUpdated, period []*big.Int) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorMerkleRootUpdated)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
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

// ParseMerkleRootUpdated is a log parse operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParseMerkleRootUpdated(log types.Log) (*DepositCashbackDistributorMerkleRootUpdated, error) {
	event := new(DepositCashbackDistributorMerkleRootUpdated)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorOwnershipTransferredIterator struct {
	Event *DepositCashbackDistributorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorOwnershipTransferred)
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
		it.Event = new(DepositCashbackDistributorOwnershipTransferred)
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
func (it *DepositCashbackDistributorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorOwnershipTransferred represents a OwnershipTransferred event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DepositCashbackDistributorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorOwnershipTransferredIterator{contract: _DepositCashbackDistributor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorOwnershipTransferred)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParseOwnershipTransferred(log types.Log) (*DepositCashbackDistributorOwnershipTransferred, error) {
	event := new(DepositCashbackDistributorOwnershipTransferred)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorPausedIterator struct {
	Event *DepositCashbackDistributorPaused // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorPaused)
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
		it.Event = new(DepositCashbackDistributorPaused)
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
func (it *DepositCashbackDistributorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorPaused represents a Paused event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterPaused(opts *bind.FilterOpts) (*DepositCashbackDistributorPausedIterator, error) {

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorPausedIterator{contract: _DepositCashbackDistributor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorPaused) (event.Subscription, error) {

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorPaused)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParsePaused(log types.Log) (*DepositCashbackDistributorPaused, error) {
	event := new(DepositCashbackDistributorPaused)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorPeriodAdvancedIterator is returned from FilterPeriodAdvanced and is used to iterate over the raw logs and unpacked data for PeriodAdvanced events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorPeriodAdvancedIterator struct {
	Event *DepositCashbackDistributorPeriodAdvanced // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorPeriodAdvancedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorPeriodAdvanced)
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
		it.Event = new(DepositCashbackDistributorPeriodAdvanced)
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
func (it *DepositCashbackDistributorPeriodAdvancedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorPeriodAdvancedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorPeriodAdvanced represents a PeriodAdvanced event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorPeriodAdvanced struct {
	NewPeriod  *big.Int
	MerkleRoot [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPeriodAdvanced is a free log retrieval operation binding the contract event 0xa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf6.
//
// Solidity: event PeriodAdvanced(uint256 indexed newPeriod, bytes32 merkleRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterPeriodAdvanced(opts *bind.FilterOpts, newPeriod []*big.Int) (*DepositCashbackDistributorPeriodAdvancedIterator, error) {

	var newPeriodRule []interface{}
	for _, newPeriodItem := range newPeriod {
		newPeriodRule = append(newPeriodRule, newPeriodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "PeriodAdvanced", newPeriodRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorPeriodAdvancedIterator{contract: _DepositCashbackDistributor.contract, event: "PeriodAdvanced", logs: logs, sub: sub}, nil
}

// WatchPeriodAdvanced is a free log subscription operation binding the contract event 0xa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf6.
//
// Solidity: event PeriodAdvanced(uint256 indexed newPeriod, bytes32 merkleRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchPeriodAdvanced(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorPeriodAdvanced, newPeriod []*big.Int) (event.Subscription, error) {

	var newPeriodRule []interface{}
	for _, newPeriodItem := range newPeriod {
		newPeriodRule = append(newPeriodRule, newPeriodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "PeriodAdvanced", newPeriodRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorPeriodAdvanced)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
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

// ParsePeriodAdvanced is a log parse operation binding the contract event 0xa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf6.
//
// Solidity: event PeriodAdvanced(uint256 indexed newPeriod, bytes32 merkleRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParsePeriodAdvanced(log types.Log) (*DepositCashbackDistributorPeriodAdvanced, error) {
	event := new(DepositCashbackDistributorPeriodAdvanced)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorUnpausedIterator struct {
	Event *DepositCashbackDistributorUnpaused // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorUnpaused)
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
		it.Event = new(DepositCashbackDistributorUnpaused)
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
func (it *DepositCashbackDistributorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorUnpaused represents a Unpaused event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*DepositCashbackDistributorUnpausedIterator, error) {

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorUnpausedIterator{contract: _DepositCashbackDistributor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorUnpaused) (event.Subscription, error) {

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorUnpaused)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParseUnpaused(log types.Log) (*DepositCashbackDistributorUnpaused, error) {
	event := new(DepositCashbackDistributorUnpaused)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
