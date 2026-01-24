// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mockstablecoin

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

// MockStablecoinMetaData contains all meta data concerning the MockStablecoin contract.
var MockStablecoinMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x6080604052346103b957604080519081016001600160401b038111828210176102c1576040908152600f82526e26b7b1b59029ba30b13632b1b7b4b760891b602083015280519081016001600160401b038111828210176102c15760405260048152634d4f434b60e01b602082015281516001600160401b0381116102c157600354600181811c911680156103af575b60208210146102a357601f8111610341575b50602092601f82116001146102e057928192935f926102d5575b50508160011b915f199060031b1c1916176003555b80516001600160401b0381116102c157600454600181811c911680156102b7575b60208210146102a357601f8111610235575b50602091601f82116001146101d5579181925f926101ca575b50508160011b915f199060031b1c1916176004555b33156101b75760025464e8d4a5100081018091116101a357600255335f525f60205260405f2064e8d4a51000815401905560405164e8d4a5100081525f7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60203393a36040516106a190816103be8239f35b634e487b7160e01b5f52601160045260245ffd5b63ec442f0560e01b5f525f60045260245ffd5b015190505f8061011c565b601f1982169260045f52805f20915f5b85811061021d57508360019510610205575b505050811b01600455610131565b01515f1960f88460031b161c191690555f80806101f7565b919260206001819286850151815501940192016101e5565b818111156101035760045f52601f820160051c7f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b6020841061029b575b81601f9101920160051c03905f5b82811061028e575050610103565b5f82820155600101610280565b5f9150610272565b634e487b7160e01b5f52602260045260245ffd5b90607f16906100f1565b634e487b7160e01b5f52604160045260245ffd5b015190505f806100bb565b601f1982169360035f52805f20915f5b8681106103295750836001959610610311575b505050811b016003556100d0565b01515f1960f88460031b161c191690555f8080610303565b919260206001819286850151815501940192016102f0565b818111156100a15760035f52601f820160051c7fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b602084106103a7575b81601f9101920160051c03905f5b82811061039a5750506100a1565b5f8282015560010161038c565b5f915061037e565b90607f169061008f565b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c90816306fdde031461049c57508063095ea7b31461041a57806318160ddd146103fd57806323b872dd1461031e578063313ce5671461030357806340c10f191461026157806370a082311461022a57806395d89b411461010f578063a9059cbb146100de5763dd62ed3e1461008a575f80fd5b346100da5760403660031901126100da576100a3610595565b6100ab6105ab565b6001600160a01b039182165f908152600160209081526040808320949093168252928352819020549051908152f35b5f80fd5b346100da5760403660031901126100da576101046100fa610595565b60243590336105c1565b602060405160018152f35b346100da575f3660031901126100da576040515f6004548060011c90600181168015610220575b60208310811461020c578285529081156101f0575060011461019b575b50819003601f01601f191681019067ffffffffffffffff821181831017610187576101838291826040528261056b565b0390f35b634e487b7160e01b5f52604160045260245ffd5b905060045f527f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b5f905b8282106101da57506020915082010182610153565b60018160209254838588010152019101906101c5565b90506020925060ff191682840152151560051b82010182610153565b634e487b7160e01b5f52602260045260245ffd5b91607f1691610136565b346100da5760203660031901126100da576001600160a01b0361024b610595565b165f525f602052602060405f2054604051908152f35b346100da5760403660031901126100da5761027a610595565b6001600160a01b031660243581156102f057600254908082018092116102dc5760207fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef915f9360025584845283825260408420818154019055604051908152a3005b634e487b7160e01b5f52601160045260245ffd5b63ec442f0560e01b5f525f60045260245ffd5b346100da575f3660031901126100da57602060405160068152f35b346100da5760603660031901126100da57610337610595565b61033f6105ab565b6001600160a01b0382165f818152600160209081526040808320338452909152902054909260443592915f19811061037d575b5061010493506105c1565b8381106103e25784156103cf5733156103bc57610104945f52600160205260405f2060018060a01b0333165f526020528360405f209103905584610372565b634a1406b160e11b5f525f60045260245ffd5b63e602df0560e01b5f525f60045260245ffd5b8390637dc7a0d960e11b5f523360045260245260445260645ffd5b346100da575f3660031901126100da576020600254604051908152f35b346100da5760403660031901126100da57610433610595565b6024359033156103cf576001600160a01b03169081156103bc57335f52600160205260405f20825f526020528060405f20556040519081527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560203392a3602060405160018152f35b346100da575f3660031901126100da575f6003548060011c90600181168015610561575b60208310811461020c578285529081156101f0575060011461050c5750819003601f01601f191681019067ffffffffffffffff821181831017610187576101838291826040528261056b565b905060035f527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b5f905b82821061054b57506020915082010182610153565b6001816020925483858801015201910190610536565b91607f16916104c0565b602060409281835280519182918282860152018484015e5f828201840152601f01601f1916010190565b600435906001600160a01b03821682036100da57565b602435906001600160a01b03821682036100da57565b6001600160a01b0316908115610658576001600160a01b03169182156102f057815f525f60205260405f205481811061063f57817fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92602092855f525f84520360405f2055845f525f825260405f20818154019055604051908152a3565b8263391434e360e21b5f5260045260245260445260645ffd5b634b637e8f60e11b5f525f60045260245ffdfea26469706673582212209bd4386ae7ddbb983964f1e9b58606e44b14cd6801119f7d9eaed6fdeb8875d764736f6c63430008210033",
}

// MockStablecoinABI is the input ABI used to generate the binding from.
// Deprecated: Use MockStablecoinMetaData.ABI instead.
var MockStablecoinABI = MockStablecoinMetaData.ABI

// MockStablecoinBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockStablecoinMetaData.Bin instead.
var MockStablecoinBin = MockStablecoinMetaData.Bin

// DeployMockStablecoin deploys a new Ethereum contract, binding an instance of MockStablecoin to it.
func DeployMockStablecoin(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MockStablecoin, error) {
	parsed, err := MockStablecoinMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockStablecoinBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockStablecoin{MockStablecoinCaller: MockStablecoinCaller{contract: contract}, MockStablecoinTransactor: MockStablecoinTransactor{contract: contract}, MockStablecoinFilterer: MockStablecoinFilterer{contract: contract}}, nil
}

// MockStablecoin is an auto generated Go binding around an Ethereum contract.
type MockStablecoin struct {
	MockStablecoinCaller     // Read-only binding to the contract
	MockStablecoinTransactor // Write-only binding to the contract
	MockStablecoinFilterer   // Log filterer for contract events
}

// MockStablecoinCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockStablecoinCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockStablecoinTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockStablecoinTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockStablecoinFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockStablecoinFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockStablecoinSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockStablecoinSession struct {
	Contract     *MockStablecoin   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockStablecoinCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockStablecoinCallerSession struct {
	Contract *MockStablecoinCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MockStablecoinTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockStablecoinTransactorSession struct {
	Contract     *MockStablecoinTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MockStablecoinRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockStablecoinRaw struct {
	Contract *MockStablecoin // Generic contract binding to access the raw methods on
}

// MockStablecoinCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockStablecoinCallerRaw struct {
	Contract *MockStablecoinCaller // Generic read-only contract binding to access the raw methods on
}

// MockStablecoinTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockStablecoinTransactorRaw struct {
	Contract *MockStablecoinTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockStablecoin creates a new instance of MockStablecoin, bound to a specific deployed contract.
func NewMockStablecoin(address common.Address, backend bind.ContractBackend) (*MockStablecoin, error) {
	contract, err := bindMockStablecoin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockStablecoin{MockStablecoinCaller: MockStablecoinCaller{contract: contract}, MockStablecoinTransactor: MockStablecoinTransactor{contract: contract}, MockStablecoinFilterer: MockStablecoinFilterer{contract: contract}}, nil
}

// NewMockStablecoinCaller creates a new read-only instance of MockStablecoin, bound to a specific deployed contract.
func NewMockStablecoinCaller(address common.Address, caller bind.ContractCaller) (*MockStablecoinCaller, error) {
	contract, err := bindMockStablecoin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockStablecoinCaller{contract: contract}, nil
}

// NewMockStablecoinTransactor creates a new write-only instance of MockStablecoin, bound to a specific deployed contract.
func NewMockStablecoinTransactor(address common.Address, transactor bind.ContractTransactor) (*MockStablecoinTransactor, error) {
	contract, err := bindMockStablecoin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockStablecoinTransactor{contract: contract}, nil
}

// NewMockStablecoinFilterer creates a new log filterer instance of MockStablecoin, bound to a specific deployed contract.
func NewMockStablecoinFilterer(address common.Address, filterer bind.ContractFilterer) (*MockStablecoinFilterer, error) {
	contract, err := bindMockStablecoin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockStablecoinFilterer{contract: contract}, nil
}

// bindMockStablecoin binds a generic wrapper to an already deployed contract.
func bindMockStablecoin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockStablecoinMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockStablecoin *MockStablecoinRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockStablecoin.Contract.MockStablecoinCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockStablecoin *MockStablecoinRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockStablecoin.Contract.MockStablecoinTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockStablecoin *MockStablecoinRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockStablecoin.Contract.MockStablecoinTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockStablecoin *MockStablecoinCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockStablecoin.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockStablecoin *MockStablecoinTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockStablecoin.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockStablecoin *MockStablecoinTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockStablecoin.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockStablecoin *MockStablecoinCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockStablecoin.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockStablecoin *MockStablecoinSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockStablecoin.Contract.Allowance(&_MockStablecoin.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockStablecoin *MockStablecoinCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockStablecoin.Contract.Allowance(&_MockStablecoin.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockStablecoin *MockStablecoinCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockStablecoin.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockStablecoin *MockStablecoinSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockStablecoin.Contract.BalanceOf(&_MockStablecoin.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockStablecoin *MockStablecoinCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockStablecoin.Contract.BalanceOf(&_MockStablecoin.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockStablecoin *MockStablecoinCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockStablecoin.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockStablecoin *MockStablecoinSession) Decimals() (uint8, error) {
	return _MockStablecoin.Contract.Decimals(&_MockStablecoin.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_MockStablecoin *MockStablecoinCallerSession) Decimals() (uint8, error) {
	return _MockStablecoin.Contract.Decimals(&_MockStablecoin.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockStablecoin *MockStablecoinCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockStablecoin.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockStablecoin *MockStablecoinSession) Name() (string, error) {
	return _MockStablecoin.Contract.Name(&_MockStablecoin.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockStablecoin *MockStablecoinCallerSession) Name() (string, error) {
	return _MockStablecoin.Contract.Name(&_MockStablecoin.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockStablecoin *MockStablecoinCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockStablecoin.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockStablecoin *MockStablecoinSession) Symbol() (string, error) {
	return _MockStablecoin.Contract.Symbol(&_MockStablecoin.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockStablecoin *MockStablecoinCallerSession) Symbol() (string, error) {
	return _MockStablecoin.Contract.Symbol(&_MockStablecoin.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockStablecoin *MockStablecoinCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockStablecoin.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockStablecoin *MockStablecoinSession) TotalSupply() (*big.Int, error) {
	return _MockStablecoin.Contract.TotalSupply(&_MockStablecoin.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockStablecoin *MockStablecoinCallerSession) TotalSupply() (*big.Int, error) {
	return _MockStablecoin.Contract.TotalSupply(&_MockStablecoin.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.Approve(&_MockStablecoin.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.Approve(&_MockStablecoin.TransactOpts, spender, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_MockStablecoin *MockStablecoinTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_MockStablecoin *MockStablecoinSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.Mint(&_MockStablecoin.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_MockStablecoin *MockStablecoinTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.Mint(&_MockStablecoin.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.Transfer(&_MockStablecoin.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.Transfer(&_MockStablecoin.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.TransferFrom(&_MockStablecoin.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockStablecoin *MockStablecoinTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockStablecoin.Contract.TransferFrom(&_MockStablecoin.TransactOpts, from, to, value)
}

// MockStablecoinApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the MockStablecoin contract.
type MockStablecoinApprovalIterator struct {
	Event *MockStablecoinApproval // Event containing the contract specifics and raw log

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
func (it *MockStablecoinApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockStablecoinApproval)
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
		it.Event = new(MockStablecoinApproval)
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
func (it *MockStablecoinApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockStablecoinApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockStablecoinApproval represents a Approval event raised by the MockStablecoin contract.
type MockStablecoinApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockStablecoin *MockStablecoinFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*MockStablecoinApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockStablecoin.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &MockStablecoinApprovalIterator{contract: _MockStablecoin.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockStablecoin *MockStablecoinFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MockStablecoinApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockStablecoin.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockStablecoinApproval)
				if err := _MockStablecoin.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_MockStablecoin *MockStablecoinFilterer) ParseApproval(log types.Log) (*MockStablecoinApproval, error) {
	event := new(MockStablecoinApproval)
	if err := _MockStablecoin.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockStablecoinTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the MockStablecoin contract.
type MockStablecoinTransferIterator struct {
	Event *MockStablecoinTransfer // Event containing the contract specifics and raw log

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
func (it *MockStablecoinTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockStablecoinTransfer)
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
		it.Event = new(MockStablecoinTransfer)
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
func (it *MockStablecoinTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockStablecoinTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockStablecoinTransfer represents a Transfer event raised by the MockStablecoin contract.
type MockStablecoinTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockStablecoin *MockStablecoinFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockStablecoinTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockStablecoin.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockStablecoinTransferIterator{contract: _MockStablecoin.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockStablecoin *MockStablecoinFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MockStablecoinTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockStablecoin.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockStablecoinTransfer)
				if err := _MockStablecoin.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_MockStablecoin *MockStablecoinFilterer) ParseTransfer(log types.Log) (*MockStablecoinTransfer, error) {
	event := new(MockStablecoinTransfer)
	if err := _MockStablecoin.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
