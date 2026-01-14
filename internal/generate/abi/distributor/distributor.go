// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package distributor

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

// MerkleDistributorMetaData contains all meta data concerning the MerkleDistributor contract.
var MerkleDistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"token_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"mintOnClaim_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isClaimed\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mintOnClaim\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60c0346100f257601f6108e938819003918201601f19168301916001600160401b038311848410176100f65780849260409485528339810103126100f25780516001600160a01b0391828216918290036100f257602001519081151582036100f25733156100da575f543360018060a01b03198216175f55604051933391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a360805260a0526107de908161010b82396080518181816092015281816104e001526105f8015260a05181818161028601526104b90152f35b604051631e4fbdf760e01b81525f6004820152602490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f803560e01c80632e7ba6ef146102c95780632eb4a7ab146102ab57806359aae2fe1461026e578063715018a6146102155780637cb64759146101c05780638da5cb5b146101995780639e34070f14610157578063f2fde38b146100c45763fc0c546a1461007d575f80fd5b346100c157806003193601126100c1576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b80fd5b50346100c15760203660031901126100c1576004356001600160a01b0381811691829003610153576100f461077d565b811561013a575f54826bffffffffffffffffffffffff60a01b8216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a380f35b604051631e4fbdf760e01b815260048101849052602490fd5b5f80fd5b50346100c15760203660031901126100c157602061018f6004358060081c5f526002602052600160ff60405f205492161b8091161490565b6040519015158152f35b50346100c157806003193601126100c157546040516001600160a01b039091168152602090f35b50346100c15760203660031901126100c1576004356101dd61077d565b7ffd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f875760406001548151908152836020820152a160015580f35b50346100c157806003193601126100c15761022e61077d565b5f80546001600160a01b0319811682556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b50346100c157806003193601126100c15760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b50346100c157806003193601126100c1576020600154604051908152f35b50346100c15760803660031901126100c157602435906001600160a01b03821682036100c15767ffffffffffffffff606435116100c1573660236064350112156100c15767ffffffffffffffff60643560040135116100c1573660246064356004013560051b6064350101116100c15761035d6004358060081c5f526002602052600160ff60405f205492161b8091161490565b61072757604051916020830160043581526bffffffffffffffffffffffff198260601b16604085015260443560548501526054845283608081011067ffffffffffffffff608086011117610713576080848101604081905285519092206001549592909167ffffffffffffffff6004606435013560051b603f01601f191685019091019081119111176106ff57600460643590810135600581901b603f01601f19168401608090810160405284015260240160a083015b60246064356004013560051b606435010182106106ef5750509183925b60808301518410156104905760a08460051b84010151908181105f14610481578552602052604084205b925f19811461046d5760010192610431565b634e487b7160e01b85526011600452602485fd5b9085526020526040842061045b565b858591036106a05760043560081c81526002602052604081208054600160ff600435161b1790557f00000000000000000000000000000000000000000000000000000000000000005f146105d2577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156105ce576040516340c10f1960e01b81526001600160a01b0384166004820152604480356024830152909183918391829084905af180156105c357610596575b509060607f4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed026915b6040805160043581526001600160a01b03909216602083015260443590820152a180f35b67ffffffffffffffff81116105af57604052606061054b565b634e487b7160e01b82526041600452602482fd5b6040513d84823e3d90fd5b5080fd5b60405163a9059cbb60e01b82526001600160a01b038316600452604480356024529192917f00000000000000000000000000000000000000000000000000000000000000009060209085908180855af16001855114811615610681575b826040521561066357505060607f4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed02691610572565b635274afe760e01b82526001600160a01b0316600482015260249150fd5b600181151661069757813b15153d15161661062f565b823d86823e3d90fd5b60405162461bcd60e51b815260206004820152602160248201527f4d65726b6c654469737472696275746f723a20496e76616c69642070726f6f666044820152601760f91b6064820152608490fd5b8135815260209182019101610414565b634e487b7160e01b84526041600452602484fd5b634e487b7160e01b83526041600452602483fd5b60405162461bcd60e51b815260206004820152602860248201527f4d65726b6c654469737472696275746f723a2044726f7020616c72656164792060448201526731b630b4b6b2b21760c11b6064820152608490fd5b5f546001600160a01b0316330361079057565b60405163118cdaa760e01b8152336004820152602490fdfea2646970667358221220d3ae44b020c04f46a83943c35e5c0ea799948fc68a8a7d457d59d349fb91aa0864736f6c63430008140033",
}

// MerkleDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleDistributorMetaData.ABI instead.
var MerkleDistributorABI = MerkleDistributorMetaData.ABI

// MerkleDistributorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MerkleDistributorMetaData.Bin instead.
var MerkleDistributorBin = MerkleDistributorMetaData.Bin

// DeployMerkleDistributor deploys a new Ethereum contract, binding an instance of MerkleDistributor to it.
func DeployMerkleDistributor(auth *bind.TransactOpts, backend bind.ContractBackend, token_ common.Address, mintOnClaim_ bool) (common.Address, *types.Transaction, *MerkleDistributor, error) {
	parsed, err := MerkleDistributorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MerkleDistributorBin), backend, token_, mintOnClaim_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MerkleDistributor{MerkleDistributorCaller: MerkleDistributorCaller{contract: contract}, MerkleDistributorTransactor: MerkleDistributorTransactor{contract: contract}, MerkleDistributorFilterer: MerkleDistributorFilterer{contract: contract}}, nil
}

// MerkleDistributor is an auto generated Go binding around an Ethereum contract.
type MerkleDistributor struct {
	MerkleDistributorCaller     // Read-only binding to the contract
	MerkleDistributorTransactor // Write-only binding to the contract
	MerkleDistributorFilterer   // Log filterer for contract events
}

// MerkleDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleDistributorSession struct {
	Contract     *MerkleDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MerkleDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleDistributorCallerSession struct {
	Contract *MerkleDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MerkleDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleDistributorTransactorSession struct {
	Contract     *MerkleDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MerkleDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleDistributorRaw struct {
	Contract *MerkleDistributor // Generic contract binding to access the raw methods on
}

// MerkleDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleDistributorCallerRaw struct {
	Contract *MerkleDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleDistributorTransactorRaw struct {
	Contract *MerkleDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleDistributor creates a new instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributor(address common.Address, backend bind.ContractBackend) (*MerkleDistributor, error) {
	contract, err := bindMerkleDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributor{MerkleDistributorCaller: MerkleDistributorCaller{contract: contract}, MerkleDistributorTransactor: MerkleDistributorTransactor{contract: contract}, MerkleDistributorFilterer: MerkleDistributorFilterer{contract: contract}}, nil
}

// NewMerkleDistributorCaller creates a new read-only instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributorCaller(address common.Address, caller bind.ContractCaller) (*MerkleDistributorCaller, error) {
	contract, err := bindMerkleDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorCaller{contract: contract}, nil
}

// NewMerkleDistributorTransactor creates a new write-only instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleDistributorTransactor, error) {
	contract, err := bindMerkleDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorTransactor{contract: contract}, nil
}

// NewMerkleDistributorFilterer creates a new log filterer instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleDistributorFilterer, error) {
	contract, err := bindMerkleDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorFilterer{contract: contract}, nil
}

// bindMerkleDistributor binds a generic wrapper to an already deployed contract.
func bindMerkleDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MerkleDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleDistributor *MerkleDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleDistributor.Contract.MerkleDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleDistributor *MerkleDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.MerkleDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleDistributor *MerkleDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.MerkleDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleDistributor *MerkleDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleDistributor *MerkleDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleDistributor *MerkleDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.contract.Transact(opts, method, params...)
}

// IsClaimed is a free data retrieval call binding the contract method 0x9e34070f.
//
// Solidity: function isClaimed(uint256 index) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCaller) IsClaimed(opts *bind.CallOpts, index *big.Int) (bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "isClaimed", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsClaimed is a free data retrieval call binding the contract method 0x9e34070f.
//
// Solidity: function isClaimed(uint256 index) view returns(bool)
func (_MerkleDistributor *MerkleDistributorSession) IsClaimed(index *big.Int) (bool, error) {
	return _MerkleDistributor.Contract.IsClaimed(&_MerkleDistributor.CallOpts, index)
}

// IsClaimed is a free data retrieval call binding the contract method 0x9e34070f.
//
// Solidity: function isClaimed(uint256 index) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCallerSession) IsClaimed(index *big.Int) (bool, error) {
	return _MerkleDistributor.Contract.IsClaimed(&_MerkleDistributor.CallOpts, index)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorSession) MerkleRoot() ([32]byte, error) {
	return _MerkleDistributor.Contract.MerkleRoot(&_MerkleDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCallerSession) MerkleRoot() ([32]byte, error) {
	return _MerkleDistributor.Contract.MerkleRoot(&_MerkleDistributor.CallOpts)
}

// MintOnClaim is a free data retrieval call binding the contract method 0x59aae2fe.
//
// Solidity: function mintOnClaim() view returns(bool)
func (_MerkleDistributor *MerkleDistributorCaller) MintOnClaim(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "mintOnClaim")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MintOnClaim is a free data retrieval call binding the contract method 0x59aae2fe.
//
// Solidity: function mintOnClaim() view returns(bool)
func (_MerkleDistributor *MerkleDistributorSession) MintOnClaim() (bool, error) {
	return _MerkleDistributor.Contract.MintOnClaim(&_MerkleDistributor.CallOpts)
}

// MintOnClaim is a free data retrieval call binding the contract method 0x59aae2fe.
//
// Solidity: function mintOnClaim() view returns(bool)
func (_MerkleDistributor *MerkleDistributorCallerSession) MintOnClaim() (bool, error) {
	return _MerkleDistributor.Contract.MintOnClaim(&_MerkleDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleDistributor *MerkleDistributorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleDistributor *MerkleDistributorSession) Owner() (common.Address, error) {
	return _MerkleDistributor.Contract.Owner(&_MerkleDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleDistributor *MerkleDistributorCallerSession) Owner() (common.Address, error) {
	return _MerkleDistributor.Contract.Owner(&_MerkleDistributor.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_MerkleDistributor *MerkleDistributorCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_MerkleDistributor *MerkleDistributorSession) Token() (common.Address, error) {
	return _MerkleDistributor.Contract.Token(&_MerkleDistributor.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_MerkleDistributor *MerkleDistributorCallerSession) Token() (common.Address, error) {
	return _MerkleDistributor.Contract.Token(&_MerkleDistributor.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x2e7ba6ef.
//
// Solidity: function claim(uint256 index, address account, uint256 amount, bytes32[] merkleProof) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) Claim(opts *bind.TransactOpts, index *big.Int, account common.Address, amount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "claim", index, account, amount, merkleProof)
}

// Claim is a paid mutator transaction binding the contract method 0x2e7ba6ef.
//
// Solidity: function claim(uint256 index, address account, uint256 amount, bytes32[] merkleProof) returns()
func (_MerkleDistributor *MerkleDistributorSession) Claim(index *big.Int, account common.Address, amount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Claim(&_MerkleDistributor.TransactOpts, index, account, amount, merkleProof)
}

// Claim is a paid mutator transaction binding the contract method 0x2e7ba6ef.
//
// Solidity: function claim(uint256 index, address account, uint256 amount, bytes32[] merkleProof) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) Claim(index *big.Int, account common.Address, amount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Claim(&_MerkleDistributor.TransactOpts, index, account, amount, merkleProof)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleDistributor *MerkleDistributorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleDistributor *MerkleDistributorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MerkleDistributor.Contract.RenounceOwnership(&_MerkleDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MerkleDistributor.Contract.RenounceOwnership(&_MerkleDistributor.TransactOpts)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) SetMerkleRoot(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "setMerkleRoot", _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_MerkleDistributor *MerkleDistributorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetMerkleRoot(&_MerkleDistributor.TransactOpts, _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetMerkleRoot(&_MerkleDistributor.TransactOpts, _merkleRoot)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleDistributor *MerkleDistributorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.TransferOwnership(&_MerkleDistributor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.TransferOwnership(&_MerkleDistributor.TransactOpts, newOwner)
}

// MerkleDistributorClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the MerkleDistributor contract.
type MerkleDistributorClaimedIterator struct {
	Event *MerkleDistributorClaimed // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorClaimed)
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
		it.Event = new(MerkleDistributorClaimed)
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
func (it *MerkleDistributorClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorClaimed represents a Claimed event raised by the MerkleDistributor contract.
type MerkleDistributorClaimed struct {
	Index   *big.Int
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed026.
//
// Solidity: event Claimed(uint256 index, address account, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterClaimed(opts *bind.FilterOpts) (*MerkleDistributorClaimedIterator, error) {

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorClaimedIterator{contract: _MerkleDistributor.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed026.
//
// Solidity: event Claimed(uint256 index, address account, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *MerkleDistributorClaimed) (event.Subscription, error) {

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorClaimed)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseClaimed(log types.Log) (*MerkleDistributorClaimed, error) {
	event := new(MerkleDistributorClaimed)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorMerkleRootUpdatedIterator is returned from FilterMerkleRootUpdated and is used to iterate over the raw logs and unpacked data for MerkleRootUpdated events raised by the MerkleDistributor contract.
type MerkleDistributorMerkleRootUpdatedIterator struct {
	Event *MerkleDistributorMerkleRootUpdated // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorMerkleRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorMerkleRootUpdated)
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
		it.Event = new(MerkleDistributorMerkleRootUpdated)
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
func (it *MerkleDistributorMerkleRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorMerkleRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorMerkleRootUpdated represents a MerkleRootUpdated event raised by the MerkleDistributor contract.
type MerkleDistributorMerkleRootUpdated struct {
	OldRoot [32]byte
	NewRoot [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootUpdated is a free log retrieval operation binding the contract event 0xfd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f8757.
//
// Solidity: event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterMerkleRootUpdated(opts *bind.FilterOpts) (*MerkleDistributorMerkleRootUpdatedIterator, error) {

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "MerkleRootUpdated")
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorMerkleRootUpdatedIterator{contract: _MerkleDistributor.contract, event: "MerkleRootUpdated", logs: logs, sub: sub}, nil
}

// WatchMerkleRootUpdated is a free log subscription operation binding the contract event 0xfd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f8757.
//
// Solidity: event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchMerkleRootUpdated(opts *bind.WatchOpts, sink chan<- *MerkleDistributorMerkleRootUpdated) (event.Subscription, error) {

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "MerkleRootUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorMerkleRootUpdated)
				if err := _MerkleDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseMerkleRootUpdated(log types.Log) (*MerkleDistributorMerkleRootUpdated, error) {
	event := new(MerkleDistributorMerkleRootUpdated)
	if err := _MerkleDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MerkleDistributor contract.
type MerkleDistributorOwnershipTransferredIterator struct {
	Event *MerkleDistributorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorOwnershipTransferred)
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
		it.Event = new(MerkleDistributorOwnershipTransferred)
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
func (it *MerkleDistributorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorOwnershipTransferred represents a OwnershipTransferred event raised by the MerkleDistributor contract.
type MerkleDistributorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MerkleDistributorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorOwnershipTransferredIterator{contract: _MerkleDistributor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MerkleDistributorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorOwnershipTransferred)
				if err := _MerkleDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseOwnershipTransferred(log types.Log) (*MerkleDistributorOwnershipTransferred, error) {
	event := new(MerkleDistributorOwnershipTransferred)
	if err := _MerkleDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
