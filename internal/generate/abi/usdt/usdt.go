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
	Bin: "0x6080346200039c576040906001600160401b039080830182811182821017620003885783526009815260209068135bd8dac81554d11560ba1b828201528351928484018481108282111762000388578552600493848152631554d11560e21b84820152825182811162000375576003918254916001958684811c941680156200036a575b8885101462000357578190601f9485811162000304575b508890858311600114620002a1575f9262000295575b50505f1982861b1c191690861b1783555b8051938411620002825786548581811c9116801562000277575b8782101462000264578281116200021c575b5085918411600114620001b5579383949184925f95620001a9575b50501b925f19911b1c19161782555b3315620001935760025464e8d4a51000928382018092116200018057505f917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9160025533835282815284832084815401905584519384523393a3516107489081620003a18239f35b601190634e487b7160e01b5f525260245ffd5b505f602492519163ec442f0560e01b8352820152fd5b015193505f8062000108565b9190601f19841692875f5284875f20945f5b8989838310620002045750505010620001ea575b50505050811b01825562000117565b01519060f8845f19921b161c191690555f808080620001db565b868601518955909701969485019488935001620001c7565b875f52865f208380870160051c8201928988106200025a575b0160051c019086905b8281106200024e575050620000ed565b5f81550186906200023e565b9250819262000235565b602288634e487b7160e01b5f525260245ffd5b90607f1690620000db565b604187634e487b7160e01b5f525260245ffd5b015190505f80620000b0565b90889350601f19831691875f528a5f20925f5b8c828210620002ed5750508411620002d5575b505050811b018355620000c1565b01515f1983881b60f8161c191690555f8080620002c7565b8385015186558c97909501949384019301620002b4565b909150855f52885f208580850160051c8201928b86106200034d575b918a91869594930160051c01915b8281106200033e5750506200009a565b5f81558594508a91016200032e565b9250819262000320565b602289634e487b7160e01b5f525260245ffd5b93607f169362000083565b604186634e487b7160e01b5f525260245ffd5b634e487b7160e01b5f52604160045260245ffd5b5f80fdfe6080604081815260049182361015610015575f80fd5b5f92833560e01c91826306fdde03146104e857508163095ea7b31461043e57816318160ddd1461041f57816323b872dd14610329578163313ce5671461030d57816340c10f191461026157816370a082311461022a57816395d89b411461010b57508063a9059cbb146100db5763dd62ed3e14610090575f80fd5b346100d757806003193601126100d757806020926100ac610607565b6100b4610621565b6001600160a01b0391821683526001865283832091168252845220549051908152f35b5080fd5b50346100d757806003193601126100d7576020906101046100fa610607565b6024359033610637565b5160018152f35b8383346100d757816003193601126100d757805190828454600181811c90808316928315610220575b602093848410811461020d578388529081156101f1575060011461019c575b505050829003601f01601f191682019267ffffffffffffffff84118385101761018957508291826101859252826105c0565b0390f35b634e487b7160e01b815260418552602490fd5b8787529192508591837f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b5b8385106101dd5750505050830101858080610153565b8054888601830152930192849082016101c7565b60ff1916878501525050151560051b8401019050858080610153565b634e487b7160e01b895260228a52602489fd5b91607f1691610134565b5050346100d75760203660031901126100d75760209181906001600160a01b03610252610607565b16815280845220549051908152f35b9190503461030957806003193601126103095761027c610607565b6001600160a01b031691602435919083156102f457600254908382018092116102e1575084927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9260209260025585855284835280852082815401905551908152a380f35b634e487b7160e01b865260119052602485fd5b84602492519163ec442f0560e01b8352820152fd5b8280fd5b5050346100d757816003193601126100d7576020905160068152f35b9050823461041c57606036600319011261041c57610345610607565b61034d610621565b916044359360018060a01b038316808352600160205286832033845260205286832054915f198310610388575b602088610104898989610637565b8683106103f05781156103d95733156103c2575082526001602090815286832033845281529186902090859003905582906101048761037a565b8751634a1406b160e11b8152908101849052602490fd5b875163e602df0560e01b8152908101849052602490fd5b8751637dc7a0d960e11b8152339181019182526020820193909352604081018790528291506060010390fd5b80fd5b5050346100d757816003193601126100d7576020906002549051908152f35b905034610309578160031936011261030957610458610607565b6024359033156104d1576001600160a01b03169182156104ba57508083602095338152600187528181208582528752205582519081527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925843392a35160018152f35b8351634a1406b160e11b8152908101859052602490fd5b835163e602df0560e01b8152808401869052602490fd5b8490843461030957826003193601126103095782600354600181811c908083169283156105b6575b602093848410811461020d578388529081156101f1575060011461056057505050829003601f01601f191682019267ffffffffffffffff84118385101761018957508291826101859252826105c0565b600387529192508591837fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b5b8385106105a25750505050830101858080610153565b80548886018301529301928490820161058c565b91607f1691610510565b602080825282518183018190529093925f5b8281106105f357505060409293505f838284010152601f8019910116010190565b8181018601518482016040015285016105d2565b600435906001600160a01b038216820361061d57565b5f80fd5b602435906001600160a01b038216820361061d57565b916001600160a01b038084169283156106fa57169283156106e2575f90838252816020526040822054908382106106b0575091604082827fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef958760209652828652038282205586815220818154019055604051908152a3565b60405163391434e360e21b81526001600160a01b03919091166004820152602481019190915260448101839052606490fd5b60405163ec442f0560e01b81525f6004820152602490fd5b604051634b637e8f60e11b81525f6004820152602490fdfea2646970667358221220274940acf65fd3f82389c45dd500b2d8579591cf730c75df96df7fbca8aacef564736f6c63430008140033",
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
