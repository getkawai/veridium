// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package escrow

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

// OTCMarketMetaData contains all meta data concerning the OTCMarket contract.
var OTCMarketMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenDeAI\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_usdt\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_feeRecipient\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"}],\"name\":\"OrderCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"OrderCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"OrderFulfilled\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FEE_BPS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_orderId\",\"type\":\"uint256\"}],\"name\":\"buyOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_orderId\",\"type\":\"uint256\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_priceInUSDT\",\"type\":\"uint256\"}],\"name\":\"createOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeRecipient\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDocumentsCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"orders\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceInUSDT\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenDeAI\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdt\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620016af380380620016af83398181016040528101906200003791906200023f565b6001600081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603620000b1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000a890620002fc565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000123576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200011a906200036e565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508173ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff168152505080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505062000390565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200020782620001da565b9050919050565b6200021981620001fa565b81146200022557600080fd5b50565b60008151905062000239816200020e565b92915050565b6000806000606084860312156200025b576200025a620001d5565b5b60006200026b8682870162000228565b93505060206200027e8682870162000228565b9250506040620002918682870162000228565b9150509250925092565b600082825260208201905092915050565b7f496e76616c696420746f6b656e20616464726573730000000000000000000000600082015250565b6000620002e46015836200029b565b9150620002f182620002ac565b602082019050919050565b600060208201905081810360008301526200031781620002d5565b9050919050565b7f496e76616c696420555344542061646472657373000000000000000000000000600082015250565b6000620003566014836200029b565b915062000363826200031e565b602082019050919050565b60006020820190508181036000830152620003898162000347565b9050919050565b60805160a0516112d6620003d96000396000818161036a015281816103dc015261051401526000818161042c015281816106dc0152818161080601526109a701526112d66000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063514fcac711610066578063514fcac71461010e57806379109baa1461012a57806396d875dc14610146578063a85c38ef14610164578063bf333f2c1461019857610093565b80630bcf963b1461009857806322f85eaa146100b65780632f48ab7d146100d257806346904840146100f0575b600080fd5b6100a06101b6565b6040516100ad9190610c4d565b60405180910390f35b6100d060048036038101906100cb9190610c99565b6101c3565b005b6100da610512565b6040516100e79190610d45565b60405180910390f35b6100f8610536565b6040516101059190610d81565b60405180910390f35b61012860048036038101906101239190610c99565b61055c565b005b610144600480360381019061013f9190610d9c565b610770565b005b61014e6109a5565b60405161015b9190610d45565b60405180910390f35b61017e60048036038101906101799190610c99565b6109c9565b60405161018f959493929190610df7565b60405180910390f35b6101a0610a3c565b6040516101ad9190610c4d565b60405180910390f35b6000600280549050905090565b6101cb610a41565b6002805490508110610212576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020990610ea7565b60405180910390fd5b60006002828154811061022857610227610ec7565b5b906000526020600020906005020190508060040160009054906101000a900460ff16610289576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161028090610f42565b60405180910390fd5b60008160040160006101000a81548160ff0219169083151502179055506000612710600083600301546102bc9190610f91565b6102c69190611002565b905060008183600301546102da9190611033565b905060008211801561033b5750600073ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614155b156103b0576103af33600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16847f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a87909392919063ffffffff16565b5b610421338460010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a87909392919063ffffffff16565b6104703384600201547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610b099092919063ffffffff16565b8260010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16857fe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d44866002015487600301546040516104fc929190611067565b60405180910390a450505061050f610b88565b50565b7f000000000000000000000000000000000000000000000000000000000000000081565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610564610a41565b60028054905081106105ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105a290610ea7565b60405180910390fd5b6000600282815481106105c1576105c0610ec7565b5b906000526020600020906005020190503373ffffffffffffffffffffffffffffffffffffffff168160010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610663576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161065a906110dc565b60405180910390fd5b8060040160009054906101000a900460ff166106b4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106ab90611148565b60405180910390fd5b60008160040160006101000a81548160ff0219169083151502179055506107203382600201547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610b099092919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff16827fc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d59660405160405180910390a35061076d610b88565b50565b610778610a41565b600082116107bb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107b2906111b4565b60405180910390fd5b600081116107fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107f590611220565b60405180910390fd5b61084b3330847f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a87909392919063ffffffff16565b6000600280549050905060026040518060a001604052808381526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018481526020016001151581525090806001815401808255809150506001900390600052602060002090600502016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604082015181600201556060820151816003015560808201518160040160006101000a81548160ff02191690831515021790555050503373ffffffffffffffffffffffffffffffffffffffff16817ff7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d58585604051610990929190611067565b60405180910390a3506109a1610b88565b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600281815481106109d957600080fd5b90600052602060002090600502016000915090508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154908060040160009054906101000a900460ff16905085565b600081565b600260005403610a7d576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002600081905550565b610b03848573ffffffffffffffffffffffffffffffffffffffff166323b872dd868686604051602401610abc93929190611240565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050610b92565b50505050565b610b83838473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8585604051602401610b3c929190611277565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050610b92565b505050565b6001600081905550565b600080602060008451602086016000885af180610bb5576040513d6000823e3d81fd5b3d925060005191505060008214610bd0576001811415610bec565b60008473ffffffffffffffffffffffffffffffffffffffff163b145b15610c2e57836040517f5274afe7000000000000000000000000000000000000000000000000000000008152600401610c259190610d81565b60405180910390fd5b50505050565b6000819050919050565b610c4781610c34565b82525050565b6000602082019050610c626000830184610c3e565b92915050565b600080fd5b610c7681610c34565b8114610c8157600080fd5b50565b600081359050610c9381610c6d565b92915050565b600060208284031215610caf57610cae610c68565b5b6000610cbd84828501610c84565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6000610d0b610d06610d0184610cc6565b610ce6565b610cc6565b9050919050565b6000610d1d82610cf0565b9050919050565b6000610d2f82610d12565b9050919050565b610d3f81610d24565b82525050565b6000602082019050610d5a6000830184610d36565b92915050565b6000610d6b82610cc6565b9050919050565b610d7b81610d60565b82525050565b6000602082019050610d966000830184610d72565b92915050565b60008060408385031215610db357610db2610c68565b5b6000610dc185828601610c84565b9250506020610dd285828601610c84565b9150509250929050565b60008115159050919050565b610df181610ddc565b82525050565b600060a082019050610e0c6000830188610c3e565b610e196020830187610d72565b610e266040830186610c3e565b610e336060830185610c3e565b610e406080830184610de8565b9695505050505050565b600082825260208201905092915050565b7f496e76616c6964204f7264657220494400000000000000000000000000000000600082015250565b6000610e91601083610e4a565b9150610e9c82610e5b565b602082019050919050565b60006020820190508181036000830152610ec081610e84565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4f72646572206e6f742061637469766500000000000000000000000000000000600082015250565b6000610f2c601083610e4a565b9150610f3782610ef6565b602082019050919050565b60006020820190508181036000830152610f5b81610f1f565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610f9c82610c34565b9150610fa783610c34565b9250828202610fb581610c34565b91508282048414831517610fcc57610fcb610f62565b5b5092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061100d82610c34565b915061101883610c34565b92508261102857611027610fd3565b5b828204905092915050565b600061103e82610c34565b915061104983610c34565b925082820390508181111561106157611060610f62565b5b92915050565b600060408201905061107c6000830185610c3e565b6110896020830184610c3e565b9392505050565b7f4e6f7420796f7572206f72646572000000000000000000000000000000000000600082015250565b60006110c6600e83610e4a565b91506110d182611090565b602082019050919050565b600060208201905081810360008301526110f5816110b9565b9050919050565b7f4f7264657220616c726561647920736f6c642f63616e63656c6c656400000000600082015250565b6000611132601c83610e4a565b915061113d826110fc565b602082019050919050565b6000602082019050818103600083015261116181611125565b9050919050565b7f416d6f756e74206d757374206265203e20300000000000000000000000000000600082015250565b600061119e601283610e4a565b91506111a982611168565b602082019050919050565b600060208201905081810360008301526111cd81611191565b9050919050565b7f5072696365206d757374206265203e2030000000000000000000000000000000600082015250565b600061120a601183610e4a565b9150611215826111d4565b602082019050919050565b60006020820190508181036000830152611239816111fd565b9050919050565b60006060820190506112556000830186610d72565b6112626020830185610d72565b61126f6040830184610c3e565b949350505050565b600060408201905061128c6000830185610d72565b6112996020830184610c3e565b939250505056fea26469706673582212203cd2a1d0cb1cd033fc0f0f5e6cd9d248b2c6620ebb008ae535c3b6cfa3c949cc64736f6c63430008140033",
}

// OTCMarketABI is the input ABI used to generate the binding from.
// Deprecated: Use OTCMarketMetaData.ABI instead.
var OTCMarketABI = OTCMarketMetaData.ABI

// OTCMarketBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OTCMarketMetaData.Bin instead.
var OTCMarketBin = OTCMarketMetaData.Bin

// DeployOTCMarket deploys a new Ethereum contract, binding an instance of OTCMarket to it.
func DeployOTCMarket(auth *bind.TransactOpts, backend bind.ContractBackend, _tokenDeAI common.Address, _usdt common.Address, _feeRecipient common.Address) (common.Address, *types.Transaction, *OTCMarket, error) {
	parsed, err := OTCMarketMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OTCMarketBin), backend, _tokenDeAI, _usdt, _feeRecipient)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OTCMarket{OTCMarketCaller: OTCMarketCaller{contract: contract}, OTCMarketTransactor: OTCMarketTransactor{contract: contract}, OTCMarketFilterer: OTCMarketFilterer{contract: contract}}, nil
}

// OTCMarket is an auto generated Go binding around an Ethereum contract.
type OTCMarket struct {
	OTCMarketCaller     // Read-only binding to the contract
	OTCMarketTransactor // Write-only binding to the contract
	OTCMarketFilterer   // Log filterer for contract events
}

// OTCMarketCaller is an auto generated read-only Go binding around an Ethereum contract.
type OTCMarketCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OTCMarketTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OTCMarketTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OTCMarketFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OTCMarketFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OTCMarketSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OTCMarketSession struct {
	Contract     *OTCMarket        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OTCMarketCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OTCMarketCallerSession struct {
	Contract *OTCMarketCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// OTCMarketTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OTCMarketTransactorSession struct {
	Contract     *OTCMarketTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// OTCMarketRaw is an auto generated low-level Go binding around an Ethereum contract.
type OTCMarketRaw struct {
	Contract *OTCMarket // Generic contract binding to access the raw methods on
}

// OTCMarketCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OTCMarketCallerRaw struct {
	Contract *OTCMarketCaller // Generic read-only contract binding to access the raw methods on
}

// OTCMarketTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OTCMarketTransactorRaw struct {
	Contract *OTCMarketTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOTCMarket creates a new instance of OTCMarket, bound to a specific deployed contract.
func NewOTCMarket(address common.Address, backend bind.ContractBackend) (*OTCMarket, error) {
	contract, err := bindOTCMarket(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OTCMarket{OTCMarketCaller: OTCMarketCaller{contract: contract}, OTCMarketTransactor: OTCMarketTransactor{contract: contract}, OTCMarketFilterer: OTCMarketFilterer{contract: contract}}, nil
}

// NewOTCMarketCaller creates a new read-only instance of OTCMarket, bound to a specific deployed contract.
func NewOTCMarketCaller(address common.Address, caller bind.ContractCaller) (*OTCMarketCaller, error) {
	contract, err := bindOTCMarket(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OTCMarketCaller{contract: contract}, nil
}

// NewOTCMarketTransactor creates a new write-only instance of OTCMarket, bound to a specific deployed contract.
func NewOTCMarketTransactor(address common.Address, transactor bind.ContractTransactor) (*OTCMarketTransactor, error) {
	contract, err := bindOTCMarket(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OTCMarketTransactor{contract: contract}, nil
}

// NewOTCMarketFilterer creates a new log filterer instance of OTCMarket, bound to a specific deployed contract.
func NewOTCMarketFilterer(address common.Address, filterer bind.ContractFilterer) (*OTCMarketFilterer, error) {
	contract, err := bindOTCMarket(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OTCMarketFilterer{contract: contract}, nil
}

// bindOTCMarket binds a generic wrapper to an already deployed contract.
func bindOTCMarket(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OTCMarketMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OTCMarket *OTCMarketRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OTCMarket.Contract.OTCMarketCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OTCMarket *OTCMarketRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OTCMarket.Contract.OTCMarketTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OTCMarket *OTCMarketRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OTCMarket.Contract.OTCMarketTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OTCMarket *OTCMarketCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OTCMarket.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OTCMarket *OTCMarketTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OTCMarket.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OTCMarket *OTCMarketTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OTCMarket.Contract.contract.Transact(opts, method, params...)
}

// FEEBPS is a free data retrieval call binding the contract method 0xbf333f2c.
//
// Solidity: function FEE_BPS() view returns(uint256)
func (_OTCMarket *OTCMarketCaller) FEEBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "FEE_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEBPS is a free data retrieval call binding the contract method 0xbf333f2c.
//
// Solidity: function FEE_BPS() view returns(uint256)
func (_OTCMarket *OTCMarketSession) FEEBPS() (*big.Int, error) {
	return _OTCMarket.Contract.FEEBPS(&_OTCMarket.CallOpts)
}

// FEEBPS is a free data retrieval call binding the contract method 0xbf333f2c.
//
// Solidity: function FEE_BPS() view returns(uint256)
func (_OTCMarket *OTCMarketCallerSession) FEEBPS() (*big.Int, error) {
	return _OTCMarket.Contract.FEEBPS(&_OTCMarket.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_OTCMarket *OTCMarketCaller) FeeRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "feeRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_OTCMarket *OTCMarketSession) FeeRecipient() (common.Address, error) {
	return _OTCMarket.Contract.FeeRecipient(&_OTCMarket.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_OTCMarket *OTCMarketCallerSession) FeeRecipient() (common.Address, error) {
	return _OTCMarket.Contract.FeeRecipient(&_OTCMarket.CallOpts)
}

// GetDocumentsCount is a free data retrieval call binding the contract method 0x0bcf963b.
//
// Solidity: function getDocumentsCount() view returns(uint256)
func (_OTCMarket *OTCMarketCaller) GetDocumentsCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "getDocumentsCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDocumentsCount is a free data retrieval call binding the contract method 0x0bcf963b.
//
// Solidity: function getDocumentsCount() view returns(uint256)
func (_OTCMarket *OTCMarketSession) GetDocumentsCount() (*big.Int, error) {
	return _OTCMarket.Contract.GetDocumentsCount(&_OTCMarket.CallOpts)
}

// GetDocumentsCount is a free data retrieval call binding the contract method 0x0bcf963b.
//
// Solidity: function getDocumentsCount() view returns(uint256)
func (_OTCMarket *OTCMarketCallerSession) GetDocumentsCount() (*big.Int, error) {
	return _OTCMarket.Contract.GetDocumentsCount(&_OTCMarket.CallOpts)
}

// Orders is a free data retrieval call binding the contract method 0xa85c38ef.
//
// Solidity: function orders(uint256 ) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, bool isActive)
func (_OTCMarket *OTCMarketCaller) Orders(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id          *big.Int
	Seller      common.Address
	TokenAmount *big.Int
	PriceInUSDT *big.Int
	IsActive    bool
}, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "orders", arg0)

	outstruct := new(struct {
		Id          *big.Int
		Seller      common.Address
		TokenAmount *big.Int
		PriceInUSDT *big.Int
		IsActive    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Seller = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokenAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PriceInUSDT = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.IsActive = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// Orders is a free data retrieval call binding the contract method 0xa85c38ef.
//
// Solidity: function orders(uint256 ) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, bool isActive)
func (_OTCMarket *OTCMarketSession) Orders(arg0 *big.Int) (struct {
	Id          *big.Int
	Seller      common.Address
	TokenAmount *big.Int
	PriceInUSDT *big.Int
	IsActive    bool
}, error) {
	return _OTCMarket.Contract.Orders(&_OTCMarket.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0xa85c38ef.
//
// Solidity: function orders(uint256 ) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, bool isActive)
func (_OTCMarket *OTCMarketCallerSession) Orders(arg0 *big.Int) (struct {
	Id          *big.Int
	Seller      common.Address
	TokenAmount *big.Int
	PriceInUSDT *big.Int
	IsActive    bool
}, error) {
	return _OTCMarket.Contract.Orders(&_OTCMarket.CallOpts, arg0)
}

// TokenDeAI is a free data retrieval call binding the contract method 0x96d875dc.
//
// Solidity: function tokenDeAI() view returns(address)
func (_OTCMarket *OTCMarketCaller) TokenDeAI(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "tokenDeAI")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenDeAI is a free data retrieval call binding the contract method 0x96d875dc.
//
// Solidity: function tokenDeAI() view returns(address)
func (_OTCMarket *OTCMarketSession) TokenDeAI() (common.Address, error) {
	return _OTCMarket.Contract.TokenDeAI(&_OTCMarket.CallOpts)
}

// TokenDeAI is a free data retrieval call binding the contract method 0x96d875dc.
//
// Solidity: function tokenDeAI() view returns(address)
func (_OTCMarket *OTCMarketCallerSession) TokenDeAI() (common.Address, error) {
	return _OTCMarket.Contract.TokenDeAI(&_OTCMarket.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_OTCMarket *OTCMarketCaller) Usdt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "usdt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_OTCMarket *OTCMarketSession) Usdt() (common.Address, error) {
	return _OTCMarket.Contract.Usdt(&_OTCMarket.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_OTCMarket *OTCMarketCallerSession) Usdt() (common.Address, error) {
	return _OTCMarket.Contract.Usdt(&_OTCMarket.CallOpts)
}

// BuyOrder is a paid mutator transaction binding the contract method 0x22f85eaa.
//
// Solidity: function buyOrder(uint256 _orderId) returns()
func (_OTCMarket *OTCMarketTransactor) BuyOrder(opts *bind.TransactOpts, _orderId *big.Int) (*types.Transaction, error) {
	return _OTCMarket.contract.Transact(opts, "buyOrder", _orderId)
}

// BuyOrder is a paid mutator transaction binding the contract method 0x22f85eaa.
//
// Solidity: function buyOrder(uint256 _orderId) returns()
func (_OTCMarket *OTCMarketSession) BuyOrder(_orderId *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.BuyOrder(&_OTCMarket.TransactOpts, _orderId)
}

// BuyOrder is a paid mutator transaction binding the contract method 0x22f85eaa.
//
// Solidity: function buyOrder(uint256 _orderId) returns()
func (_OTCMarket *OTCMarketTransactorSession) BuyOrder(_orderId *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.BuyOrder(&_OTCMarket.TransactOpts, _orderId)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x514fcac7.
//
// Solidity: function cancelOrder(uint256 _orderId) returns()
func (_OTCMarket *OTCMarketTransactor) CancelOrder(opts *bind.TransactOpts, _orderId *big.Int) (*types.Transaction, error) {
	return _OTCMarket.contract.Transact(opts, "cancelOrder", _orderId)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x514fcac7.
//
// Solidity: function cancelOrder(uint256 _orderId) returns()
func (_OTCMarket *OTCMarketSession) CancelOrder(_orderId *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.CancelOrder(&_OTCMarket.TransactOpts, _orderId)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x514fcac7.
//
// Solidity: function cancelOrder(uint256 _orderId) returns()
func (_OTCMarket *OTCMarketTransactorSession) CancelOrder(_orderId *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.CancelOrder(&_OTCMarket.TransactOpts, _orderId)
}

// CreateOrder is a paid mutator transaction binding the contract method 0x79109baa.
//
// Solidity: function createOrder(uint256 _amount, uint256 _priceInUSDT) returns()
func (_OTCMarket *OTCMarketTransactor) CreateOrder(opts *bind.TransactOpts, _amount *big.Int, _priceInUSDT *big.Int) (*types.Transaction, error) {
	return _OTCMarket.contract.Transact(opts, "createOrder", _amount, _priceInUSDT)
}

// CreateOrder is a paid mutator transaction binding the contract method 0x79109baa.
//
// Solidity: function createOrder(uint256 _amount, uint256 _priceInUSDT) returns()
func (_OTCMarket *OTCMarketSession) CreateOrder(_amount *big.Int, _priceInUSDT *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.CreateOrder(&_OTCMarket.TransactOpts, _amount, _priceInUSDT)
}

// CreateOrder is a paid mutator transaction binding the contract method 0x79109baa.
//
// Solidity: function createOrder(uint256 _amount, uint256 _priceInUSDT) returns()
func (_OTCMarket *OTCMarketTransactorSession) CreateOrder(_amount *big.Int, _priceInUSDT *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.CreateOrder(&_OTCMarket.TransactOpts, _amount, _priceInUSDT)
}

// OTCMarketOrderCancelledIterator is returned from FilterOrderCancelled and is used to iterate over the raw logs and unpacked data for OrderCancelled events raised by the OTCMarket contract.
type OTCMarketOrderCancelledIterator struct {
	Event *OTCMarketOrderCancelled // Event containing the contract specifics and raw log

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
func (it *OTCMarketOrderCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OTCMarketOrderCancelled)
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
		it.Event = new(OTCMarketOrderCancelled)
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
func (it *OTCMarketOrderCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OTCMarketOrderCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OTCMarketOrderCancelled represents a OrderCancelled event raised by the OTCMarket contract.
type OTCMarketOrderCancelled struct {
	OrderId *big.Int
	Seller  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterOrderCancelled is a free log retrieval operation binding the contract event 0xc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d596.
//
// Solidity: event OrderCancelled(uint256 indexed orderId, address indexed seller)
func (_OTCMarket *OTCMarketFilterer) FilterOrderCancelled(opts *bind.FilterOpts, orderId []*big.Int, seller []common.Address) (*OTCMarketOrderCancelledIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _OTCMarket.contract.FilterLogs(opts, "OrderCancelled", orderIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &OTCMarketOrderCancelledIterator{contract: _OTCMarket.contract, event: "OrderCancelled", logs: logs, sub: sub}, nil
}

// WatchOrderCancelled is a free log subscription operation binding the contract event 0xc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d596.
//
// Solidity: event OrderCancelled(uint256 indexed orderId, address indexed seller)
func (_OTCMarket *OTCMarketFilterer) WatchOrderCancelled(opts *bind.WatchOpts, sink chan<- *OTCMarketOrderCancelled, orderId []*big.Int, seller []common.Address) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _OTCMarket.contract.WatchLogs(opts, "OrderCancelled", orderIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OTCMarketOrderCancelled)
				if err := _OTCMarket.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
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

// ParseOrderCancelled is a log parse operation binding the contract event 0xc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d596.
//
// Solidity: event OrderCancelled(uint256 indexed orderId, address indexed seller)
func (_OTCMarket *OTCMarketFilterer) ParseOrderCancelled(log types.Log) (*OTCMarketOrderCancelled, error) {
	event := new(OTCMarketOrderCancelled)
	if err := _OTCMarket.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OTCMarketOrderCreatedIterator is returned from FilterOrderCreated and is used to iterate over the raw logs and unpacked data for OrderCreated events raised by the OTCMarket contract.
type OTCMarketOrderCreatedIterator struct {
	Event *OTCMarketOrderCreated // Event containing the contract specifics and raw log

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
func (it *OTCMarketOrderCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OTCMarketOrderCreated)
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
		it.Event = new(OTCMarketOrderCreated)
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
func (it *OTCMarketOrderCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OTCMarketOrderCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OTCMarketOrderCreated represents a OrderCreated event raised by the OTCMarket contract.
type OTCMarketOrderCreated struct {
	OrderId *big.Int
	Seller  common.Address
	Amount  *big.Int
	Price   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterOrderCreated is a free log retrieval operation binding the contract event 0xf7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d5.
//
// Solidity: event OrderCreated(uint256 indexed orderId, address indexed seller, uint256 amount, uint256 price)
func (_OTCMarket *OTCMarketFilterer) FilterOrderCreated(opts *bind.FilterOpts, orderId []*big.Int, seller []common.Address) (*OTCMarketOrderCreatedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _OTCMarket.contract.FilterLogs(opts, "OrderCreated", orderIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &OTCMarketOrderCreatedIterator{contract: _OTCMarket.contract, event: "OrderCreated", logs: logs, sub: sub}, nil
}

// WatchOrderCreated is a free log subscription operation binding the contract event 0xf7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d5.
//
// Solidity: event OrderCreated(uint256 indexed orderId, address indexed seller, uint256 amount, uint256 price)
func (_OTCMarket *OTCMarketFilterer) WatchOrderCreated(opts *bind.WatchOpts, sink chan<- *OTCMarketOrderCreated, orderId []*big.Int, seller []common.Address) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _OTCMarket.contract.WatchLogs(opts, "OrderCreated", orderIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OTCMarketOrderCreated)
				if err := _OTCMarket.contract.UnpackLog(event, "OrderCreated", log); err != nil {
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

// ParseOrderCreated is a log parse operation binding the contract event 0xf7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d5.
//
// Solidity: event OrderCreated(uint256 indexed orderId, address indexed seller, uint256 amount, uint256 price)
func (_OTCMarket *OTCMarketFilterer) ParseOrderCreated(log types.Log) (*OTCMarketOrderCreated, error) {
	event := new(OTCMarketOrderCreated)
	if err := _OTCMarket.contract.UnpackLog(event, "OrderCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OTCMarketOrderFulfilledIterator is returned from FilterOrderFulfilled and is used to iterate over the raw logs and unpacked data for OrderFulfilled events raised by the OTCMarket contract.
type OTCMarketOrderFulfilledIterator struct {
	Event *OTCMarketOrderFulfilled // Event containing the contract specifics and raw log

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
func (it *OTCMarketOrderFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OTCMarketOrderFulfilled)
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
		it.Event = new(OTCMarketOrderFulfilled)
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
func (it *OTCMarketOrderFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OTCMarketOrderFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OTCMarketOrderFulfilled represents a OrderFulfilled event raised by the OTCMarket contract.
type OTCMarketOrderFulfilled struct {
	OrderId *big.Int
	Buyer   common.Address
	Seller  common.Address
	Amount  *big.Int
	Price   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterOrderFulfilled is a free log retrieval operation binding the contract event 0xe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d44.
//
// Solidity: event OrderFulfilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amount, uint256 price)
func (_OTCMarket *OTCMarketFilterer) FilterOrderFulfilled(opts *bind.FilterOpts, orderId []*big.Int, buyer []common.Address, seller []common.Address) (*OTCMarketOrderFulfilledIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _OTCMarket.contract.FilterLogs(opts, "OrderFulfilled", orderIdRule, buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &OTCMarketOrderFulfilledIterator{contract: _OTCMarket.contract, event: "OrderFulfilled", logs: logs, sub: sub}, nil
}

// WatchOrderFulfilled is a free log subscription operation binding the contract event 0xe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d44.
//
// Solidity: event OrderFulfilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amount, uint256 price)
func (_OTCMarket *OTCMarketFilterer) WatchOrderFulfilled(opts *bind.WatchOpts, sink chan<- *OTCMarketOrderFulfilled, orderId []*big.Int, buyer []common.Address, seller []common.Address) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _OTCMarket.contract.WatchLogs(opts, "OrderFulfilled", orderIdRule, buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OTCMarketOrderFulfilled)
				if err := _OTCMarket.contract.UnpackLog(event, "OrderFulfilled", log); err != nil {
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

// ParseOrderFulfilled is a log parse operation binding the contract event 0xe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d44.
//
// Solidity: event OrderFulfilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amount, uint256 price)
func (_OTCMarket *OTCMarketFilterer) ParseOrderFulfilled(log types.Log) (*OTCMarketOrderFulfilled, error) {
	event := new(OTCMarketOrderFulfilled)
	if err := _OTCMarket.contract.UnpackLog(event, "OrderFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
