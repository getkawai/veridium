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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_tokenDeAI\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_usdt\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"FEE_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buyOrder\",\"inputs\":[{\"name\":\"_orderId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOrder\",\"inputs\":[{\"name\":\"_orderId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createOrder\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDocumentsCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"orders\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenDeAI\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usdt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OrderCancelled\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrderCreated\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrderFulfilled\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60c060405234801561000f575f5ffd5b5060405161164238038061164283398181016040528101906100319190610270565b600161004f6100446101e060201b60201c565b61020960201b60201c565b5f01819055505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036100c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100ba9061031a565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610131576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161012890610382565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508173ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1681525050805f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050506103a0565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61023f82610216565b9050919050565b61024f81610235565b8114610259575f5ffd5b50565b5f8151905061026a81610246565b92915050565b5f5f5f6060848603121561028757610286610212565b5b5f6102948682870161025c565b93505060206102a58682870161025c565b92505060406102b68682870161025c565b9150509250925092565b5f82825260208201905092915050565b7f496e76616c696420746f6b656e206164647265737300000000000000000000005f82015250565b5f6103046015836102c0565b915061030f826102d0565b602082019050919050565b5f6020820190508181035f830152610331816102f8565b9050919050565b7f496e76616c6964205553445420616464726573730000000000000000000000005f82015250565b5f61036c6014836102c0565b915061037782610338565b602082019050919050565b5f6020820190508181035f83015261039981610360565b9050919050565b60805160a05161125e6103e45f395f8181610357015281816103c801526104ff01525f8181610418015281816106be015281816107e6015261097f015261125e5ff3fe608060405234801561000f575f5ffd5b5060043610610091575f3560e01c8063514fcac711610064578063514fcac71461010b57806379109baa1461012757806396d875dc14610143578063a85c38ef14610161578063bf333f2c1461019557610091565b80630bcf963b1461009557806322f85eaa146100b35780632f48ab7d146100cf57806346904840146100ed575b5f5ffd5b61009d6101b3565b6040516100aa9190610c6f565b60405180910390f35b6100cd60048036038101906100c89190610cb6565b6101bf565b005b6100d76104fd565b6040516100e49190610d5b565b60405180910390f35b6100f5610521565b6040516101029190610d94565b60405180910390f35b61012560048036038101906101209190610cb6565b610545565b005b610141600480360381019061013c9190610dad565b610752565b005b61014b61097d565b6040516101589190610d5b565b60405180910390f35b61017b60048036038101906101769190610cb6565b6109a1565b60405161018c959493929190610e05565b60405180910390f35b61019d610a0d565b6040516101aa9190610c6f565b60405180910390f35b5f600180549050905090565b6101c7610a11565b600180549050811061020e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020590610eb0565b60405180910390fd5b5f6001828154811061022357610222610ece565b5b905f5260205f2090600502019050806004015f9054906101000a900460ff16610281576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161027890610f45565b60405180910390fd5b5f816004015f6101000a81548160ff0219169083151502179055505f6127105f83600301546102b09190610f90565b6102ba9190610ffe565b90505f8183600301546102cd919061102e565b90505f8211801561032a57505f73ffffffffffffffffffffffffffffffffffffffff165f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614155b1561039d5761039c335f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16847f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a33909392919063ffffffff16565b5b61040d33846001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a33909392919063ffffffff16565b61045c3384600201547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a889092919063ffffffff16565b826001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16857fe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d44866002015487600301546040516104e7929190611061565b60405180910390a45050506104fa610adb565b50565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b61054d610a11565b6001805490508110610594576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058b90610eb0565b60405180910390fd5b5f600182815481106105a9576105a8610ece565b5b905f5260205f20906005020190503373ffffffffffffffffffffffffffffffffffffffff16816001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610648576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161063f906110d2565b60405180910390fd5b806004015f9054906101000a900460ff16610698576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161068f9061113a565b60405180910390fd5b5f816004015f6101000a81548160ff0219169083151502179055506107023382600201547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a889092919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff16827fc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d59660405160405180910390a35061074f610adb565b50565b61075a610a11565b5f821161079c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610793906111a2565b60405180910390fd5b5f81116107de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107d59061120a565b60405180910390fd5b61082b3330847f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610a33909392919063ffffffff16565b5f600180549050905060016040518060a001604052808381526020013373ffffffffffffffffffffffffffffffffffffffff16815260200185815260200184815260200160011515815250908060018154018082558091505060019003905f5260205f2090600502015f909190919091505f820151815f01556020820151816001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160020155606082015181600301556080820151816004015f6101000a81548160ff02191690831515021790555050503373ffffffffffffffffffffffffffffffffffffffff16817ff7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d58585604051610968929190611061565b60405180910390a350610979610adb565b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600181815481106109b0575f80fd5b905f5260205f2090600502015f91509050805f015490806001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806002015490806003015490806004015f9054906101000a900460ff16905085565b5f81565b610a19610af5565b6002610a2b610a26610b36565b610b5f565b5f0181905550565b610a41848484846001610b68565b610a8257836040517f5274afe7000000000000000000000000000000000000000000000000000000008152600401610a799190610d94565b60405180910390fd5b50505050565b610a958383836001610bd9565b610ad657826040517f5274afe7000000000000000000000000000000000000000000000000000000008152600401610acd9190610d94565b60405180910390fd5b505050565b6001610aed610ae8610b36565b610b5f565b5f0181905550565b610afd610c3b565b15610b34576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5f6323b872dd60e01b9050604051815f525f1960601c87166004525f1960601c86166024528460445260205f60645f5f8c5af1925060015f51148316610bc6578383151615610bba573d5f823e3d81fd5b5f883b113d1516831692505b806040525f606052505095945050505050565b5f5f63a9059cbb60e01b9050604051815f525f1960601c86166004528460245260205f60445f5f8b5af1925060015f51148316610c2d578383151615610c21573d5f823e3d81fd5b5f873b113d1516831692505b806040525050949350505050565b5f6002610c4e610c49610b36565b610b5f565b5f015414905090565b5f819050919050565b610c6981610c57565b82525050565b5f602082019050610c825f830184610c60565b92915050565b5f5ffd5b610c9581610c57565b8114610c9f575f5ffd5b50565b5f81359050610cb081610c8c565b92915050565b5f60208284031215610ccb57610cca610c88565b5b5f610cd884828501610ca2565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f610d23610d1e610d1984610ce1565b610d00565b610ce1565b9050919050565b5f610d3482610d09565b9050919050565b5f610d4582610d2a565b9050919050565b610d5581610d3b565b82525050565b5f602082019050610d6e5f830184610d4c565b92915050565b5f610d7e82610ce1565b9050919050565b610d8e81610d74565b82525050565b5f602082019050610da75f830184610d85565b92915050565b5f5f60408385031215610dc357610dc2610c88565b5b5f610dd085828601610ca2565b9250506020610de185828601610ca2565b9150509250929050565b5f8115159050919050565b610dff81610deb565b82525050565b5f60a082019050610e185f830188610c60565b610e256020830187610d85565b610e326040830186610c60565b610e3f6060830185610c60565b610e4c6080830184610df6565b9695505050505050565b5f82825260208201905092915050565b7f496e76616c6964204f72646572204944000000000000000000000000000000005f82015250565b5f610e9a601083610e56565b9150610ea582610e66565b602082019050919050565b5f6020820190508181035f830152610ec781610e8e565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4f72646572206e6f7420616374697665000000000000000000000000000000005f82015250565b5f610f2f601083610e56565b9150610f3a82610efb565b602082019050919050565b5f6020820190508181035f830152610f5c81610f23565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610f9a82610c57565b9150610fa583610c57565b9250828202610fb381610c57565b91508282048414831517610fca57610fc9610f63565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61100882610c57565b915061101383610c57565b92508261102357611022610fd1565b5b828204905092915050565b5f61103882610c57565b915061104383610c57565b925082820390508181111561105b5761105a610f63565b5b92915050565b5f6040820190506110745f830185610c60565b6110816020830184610c60565b9392505050565b7f4e6f7420796f7572206f726465720000000000000000000000000000000000005f82015250565b5f6110bc600e83610e56565b91506110c782611088565b602082019050919050565b5f6020820190508181035f8301526110e9816110b0565b9050919050565b7f4f7264657220616c726561647920736f6c642f63616e63656c6c6564000000005f82015250565b5f611124601c83610e56565b915061112f826110f0565b602082019050919050565b5f6020820190508181035f83015261115181611118565b9050919050565b7f416d6f756e74206d757374206265203e203000000000000000000000000000005f82015250565b5f61118c601283610e56565b915061119782611158565b602082019050919050565b5f6020820190508181035f8301526111b981611180565b9050919050565b7f5072696365206d757374206265203e20300000000000000000000000000000005f82015250565b5f6111f4601183610e56565b91506111ff826111c0565b602082019050919050565b5f6020820190508181035f830152611221816111e8565b905091905056fea26469706673582212203ef65effbd15fdc47afe713f701c000c40ad31cefbcf5310476aec053eae536864736f6c634300081e0033",
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
