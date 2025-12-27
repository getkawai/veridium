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
	Bin: "0x60c060405234801561000f575f5ffd5b506040516111ec3803806111ec8339818101604052810190610031919061024d565b335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610099919061029a565b60405180910390fd5b6100b1816100f960201b60201c565b508173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505080151560a08115158152505050506102b3565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101e7826101be565b9050919050565b6101f7816101dd565b8114610201575f5ffd5b50565b5f81519050610212816101ee565b92915050565b5f8115159050919050565b61022c81610218565b8114610236575f5ffd5b50565b5f8151905061024781610223565b92915050565b5f5f60408385031215610263576102626101ba565b5b5f61027085828601610204565b925050602061028185828601610239565b9150509250929050565b610294816101dd565b82525050565b5f6020820190506102ad5f83018461028b565b92915050565b60805160a051610f036102e95f395f81816102a901526103f101525f81816102cf0152818161036101526105700152610f035ff3fe608060405234801561000f575f5ffd5b5060043610610091575f3560e01c80637cb64759116100645780637cb64759146100f75780638da5cb5b146101135780639e34070f14610131578063f2fde38b14610161578063fc0c546a1461017d57610091565b80632e7ba6ef146100955780632eb4a7ab146100b157806359aae2fe146100cf578063715018a6146100ed575b5f5ffd5b6100af60048036038101906100aa9190610985565b61019b565b005b6100b96103e9565b6040516100c69190610a21565b60405180910390f35b6100d76103ef565b6040516100e49190610a54565b60405180910390f35b6100f5610413565b005b610111600480360381019061010c9190610a97565b610426565b005b61011b610473565b6040516101289190610ad1565b60405180910390f35b61014b60048036038101906101469190610aea565b61049a565b6040516101589190610a54565b60405180910390f35b61017b60048036038101906101769190610b15565b6104ea565b005b61018561056e565b6040516101929190610b9b565b60405180910390f35b6101a48561049a565b156101e4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101db90610c34565b60405180910390fd5b5f8585856040516020016101fa93929190610cb7565b60405160208183030381529060405280519060200120905061025f8383808060200260200160405190810160405280939291908181526020018383602002808284375f81840152601f19601f8201169050808301925050505050505060015483610592565b61029e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029590610d63565b60405180910390fd5b6102a7866105a8565b7f00000000000000000000000000000000000000000000000000000000000000001561035a577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166340c10f1986866040518363ffffffff1660e01b8152600401610328929190610d90565b5f604051808303815f87803b15801561033f575f5ffd5b505af1158015610351573d5f5f3e3d5ffd5b505050506103a6565b6103a585857f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166105fc9092919063ffffffff16565b5b7f4ec90e965519d92681267467f775ada5bd214aa92c0dc93d90a5e880ce9ed0268686866040516103d993929190610db7565b60405180910390a1505050505050565b60015481565b7f000000000000000000000000000000000000000000000000000000000000000081565b61041b61064f565b6104245f6106d6565b565b61042e61064f565b7ffd69edeceaf1d6832d935be1fba54ca93bf17e71520c6c9ffc08d6e9529f875760015482604051610461929190610dec565b60405180910390a18060018190555050565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f5f610100836104aa9190610e40565b90505f610100846104bb9190610e70565b90505f60025f8481526020019081526020015f205490505f826001901b90508081831614945050505050919050565b6104f261064f565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610562575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016105599190610ad1565b60405180910390fd5b61056b816106d6565b50565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f8261059e8584610797565b1490509392505050565b5f610100826105b79190610e40565b90505f610100836105c89190610e70565b9050806001901b60025f8481526020019081526020015f20541760025f8481526020019081526020015f2081905550505050565b61060983838360016107e8565b61064a57826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016106419190610ad1565b60405180910390fd5b505050565b61065761084a565b73ffffffffffffffffffffffffffffffffffffffff16610675610473565b73ffffffffffffffffffffffffffffffffffffffff16146106d45761069861084a565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016106cb9190610ad1565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5f8290505f5f90505b84518110156107dd576107ce828683815181106107c1576107c0610ea0565b5b6020026020010151610851565b915080806001019150506107a1565b508091505092915050565b5f5f63a9059cbb60e01b9050604051815f525f1960601c86166004528460245260205f60445f5f8b5af1925060015f5114831661083c578383151615610830573d5f823e3d81fd5b5f873b113d1516831692505b806040525050949350505050565b5f33905090565b5f81831061086857610863828461087b565b610873565b610872838361087b565b5b905092915050565b5f825f528160205260405f20905092915050565b5f5ffd5b5f5ffd5b5f819050919050565b6108a981610897565b81146108b3575f5ffd5b50565b5f813590506108c4816108a0565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6108f3826108ca565b9050919050565b610903816108e9565b811461090d575f5ffd5b50565b5f8135905061091e816108fa565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f84011261094557610944610924565b5b8235905067ffffffffffffffff81111561096257610961610928565b5b60208301915083602082028301111561097e5761097d61092c565b5b9250929050565b5f5f5f5f5f6080868803121561099e5761099d61088f565b5b5f6109ab888289016108b6565b95505060206109bc88828901610910565b94505060406109cd888289016108b6565b935050606086013567ffffffffffffffff8111156109ee576109ed610893565b5b6109fa88828901610930565b92509250509295509295909350565b5f819050919050565b610a1b81610a09565b82525050565b5f602082019050610a345f830184610a12565b92915050565b5f8115159050919050565b610a4e81610a3a565b82525050565b5f602082019050610a675f830184610a45565b92915050565b610a7681610a09565b8114610a80575f5ffd5b50565b5f81359050610a9181610a6d565b92915050565b5f60208284031215610aac57610aab61088f565b5b5f610ab984828501610a83565b91505092915050565b610acb816108e9565b82525050565b5f602082019050610ae45f830184610ac2565b92915050565b5f60208284031215610aff57610afe61088f565b5b5f610b0c848285016108b6565b91505092915050565b5f60208284031215610b2a57610b2961088f565b5b5f610b3784828501610910565b91505092915050565b5f819050919050565b5f610b63610b5e610b59846108ca565b610b40565b6108ca565b9050919050565b5f610b7482610b49565b9050919050565b5f610b8582610b6a565b9050919050565b610b9581610b7b565b82525050565b5f602082019050610bae5f830184610b8c565b92915050565b5f82825260208201905092915050565b7f4d65726b6c654469737472696275746f723a2044726f7020616c7265616479205f8201527f636c61696d65642e000000000000000000000000000000000000000000000000602082015250565b5f610c1e602883610bb4565b9150610c2982610bc4565b604082019050919050565b5f6020820190508181035f830152610c4b81610c12565b9050919050565b5f819050919050565b610c6c610c6782610897565b610c52565b82525050565b5f8160601b9050919050565b5f610c8882610c72565b9050919050565b5f610c9982610c7e565b9050919050565b610cb1610cac826108e9565b610c8f565b82525050565b5f610cc28286610c5b565b602082019150610cd28285610ca0565b601482019150610ce28284610c5b565b602082019150819050949350505050565b7f4d65726b6c654469737472696275746f723a20496e76616c69642070726f6f665f8201527f2e00000000000000000000000000000000000000000000000000000000000000602082015250565b5f610d4d602183610bb4565b9150610d5882610cf3565b604082019050919050565b5f6020820190508181035f830152610d7a81610d41565b9050919050565b610d8a81610897565b82525050565b5f604082019050610da35f830185610ac2565b610db06020830184610d81565b9392505050565b5f606082019050610dca5f830186610d81565b610dd76020830185610ac2565b610de46040830184610d81565b949350505050565b5f604082019050610dff5f830185610a12565b610e0c6020830184610a12565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f610e4a82610897565b9150610e5583610897565b925082610e6557610e64610e13565b5b828204905092915050565b5f610e7a82610897565b9150610e8583610897565b925082610e9557610e94610e13565b5b828206905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffdfea2646970667358221220130204181240cba989bc7e9b1cf7d690f549644743db513b28e85600652f479b64736f6c634300081e0033",
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
