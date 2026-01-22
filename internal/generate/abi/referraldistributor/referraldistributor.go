// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package referraldistributor

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

// ReferralRewardDistributorMetaData contains all meta data concerning the ReferralRewardDistributor contract.
var ReferralRewardDistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"kawaiToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"advancePeriod\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimMultiplePeriods\",\"inputs\":[{\"name\":\"periods\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"kawaiAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"merkleProofs\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimRewards\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStats\",\"inputs\":[],\"outputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"referrers\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedAnyPeriod\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedPeriod\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kawaiToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"periodMerkleRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeriodMerkleRoot\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalKawaiDistributed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalReferrers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodAdvanced\",\"inputs\":[{\"name\":\"oldPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a03461014d57601f61115838819003918201601f19168301916001600160401b038311848410176101515780849260209460405283398101031261014d57516001600160a01b038082169182900361014d573315610135575f543360018060a01b03198216175f55604051913391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005581156100f357506080526001600255604051610ff2908161016682396080518181816102f70152818161034f0152818161051c01526109220152f35b62461bcd60e51b815260206004820152601560248201527f496e76616c6964204b41574149206164647265737300000000000000000000006044820152606490fd5b604051631e4fbdf760e01b81525f6004820152602490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe608060408181526004361015610013575f80fd5b5f91823560e01c9081630604061814610caf575080630ae6540314610c0e5780632eb4a7ab14610bf05780633f08ccd014610baa5780633f4ba83a14610b3c5780635869bc5a14610b1e5780635c975abb14610afa57806362d03cb714610adc57806366351b941461081e578063715018a6146107c2578063727a7c5a1461079a5780637cb64759146107365780638456cb59146106d8578063873f6f9e146106935780638da5cb5b1461066c578063adeacbd31461062f578063c40c91bd1461057a578063c59d48471461054b578063cb56cd4f14610508578063f2fde38b1461047c5763f75cc2b914610106575f80fd5b346104015760603660031901126104015760043567ffffffffffffffff811161047857610137903690600401610cfb565b9160243567ffffffffffffffff811161047457610158903690600401610cfb565b939060443567ffffffffffffffff81116104705761017a903690600401610cfb565b92610183610f25565b61018b610f05565b86811480610467575b1561042b57875b8181106101ca578860017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b6101d5818389610eb6565b3590818a526003602052868a20335f5260205260ff875f20541661042157610201600254831115610d2c565b61020c818a87610eb6565b8751602081018481523360601b6bffffffffffffffffffffffff191660408301529135605482015261024b81607481015b03601f198101835282610d91565b519020828b526005602052878b2054610265811515610db3565b8783101561040d57601e19863603018360051b8701351215610409578260051b86013586019182359267ffffffffffffffff841161040557602001918360051b36038313610405576102bf6102c4936102c9953691610df4565b610f67565b610e4b565b818a526003602052868a20335f52602052865f20600160ff19825416179055896102f4828b88610eb6565b357f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03163b156104015788516340c10f1960e01b815233600482015260248101919091528181604481836001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af180156103f7576103df575b50506103da9161038d828b88610eb6565b3561039b6006918254610e87565b90556103a8828b88610eb6565b359088519182527f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb79696860203393a3610ea8565b61019b565b6103e890610d69565b6103f357895f61037c565b8980fd5b89513d84823e3d90fd5b5080fd5b8d80fd5b8b80fd5b634e487b7160e01b8c52603260045260248cfd5b6103da9150610ea8565b845162461bcd60e51b8152602060048201526015602482015274082e4e4c2f240d8cadccee8d040dad2e6dac2e8c6d605b1b6044820152606490fd5b50838114610194565b8680fd5b8480fd5b8280fd5b503461040157602036600319011261040157610496610ce5565b61049e610eda565b6001600160a01b039081169182156104f157505f54826bffffffffffffffffffffffff60a01b8216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a380f35b51631e4fbdf760e01b815260048101849052602490fd5b5034610401578160031936011261040157517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b503461040157816003193601126104015760609060025490600654906007549181519384526020840152820152f35b50346104015780600319360112610401576024359060043561059a610eda565b60025481116105eb578084526005602052807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c83808720548151908152866020820152a28352600560205282205580f35b815162461bcd60e51b815260206004820152601860248201527f43616e6e6f74207365742066757475726520706572696f6400000000000000006044820152606490fd5b50346104015760203660031901126104015760209160ff9082906001600160a01b03610659610ce5565b1681526004855220541690519015158152f35b5034610401578160031936011261040157905490516001600160a01b039091168152602090f35b5034610401578060031936011261040157602091816106b0610ccb565b91600435815260038552209060018060a01b03165f52825260ff815f20541690519015158152f35b503461040157816003193601126104015760207f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25891610715610eda565b61071d610f05565b835460ff60a01b1916600160a01b17845551338152a180f35b50346104015760203660031901126104015760043590610754610eda565b600254807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c836001548151908152866020820152a2826001558352600560205282205580f35b5034610401576020366003190112610401578060209260043581526005845220549051908152f35b823461081b578060031936011261081b576107db610eda565b5f80546001600160a01b0319811682556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b80fd5b5034610401576060366003190112610401576004356024359160443567ffffffffffffffff811161047457610857903690600401610cfb565b9190610861610f25565b610869610f05565b610877600254851115610d2c565b83865260209260038452828720338852845260ff8388205416610a98578515610a5e5782518481018681523360601b6bffffffffffffffffffffffff1916602082015260348101889052610900936102c493909290916102bf916108de816054840161023d565b51902092888b5260058852868b2054926108f9841515610db3565b3691610df4565b828552600382528085203386528252808520805460ff199081166001179091557f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104705782516340c10f1960e01b8152336004820152602481018790529087908290604490829084905af18015610a5457610a1d575b50907f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb79696892916109b486600654610e87565b600655338752600483528187209081549060ff821615610a00575b505050519384523393a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b60019116179055610a12600754610ea8565b6007555f80806109cf565b95610a4b7f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb79696894939297610d69565b95909192610983565b83513d89823e3d90fd5b825162461bcd60e51b81526004810185905260136024820152724e6f207265776172647320746f20636c61696d60681b6044820152606490fd5b825162461bcd60e51b815260048101859052601f60248201527f416c726561647920636c61696d656420666f72207468697320706572696f64006044820152606490fd5b50346104015781600319360112610401576020906007549051908152f35b503461040157816003193601126104015760ff6020925460a01c1690519015158152f35b50346104015781600319360112610401576020906006549051908152f35b5034610401578160031936011261040157610b55610eda565b815460ff8160a01c1615610b995760ff60a01b19168255513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa90602090a180f35b8151638dfc202b60e01b8152600490fd5b503461040157806003193601126104015760ff81602093610bc9610ccb565b6004358252600386528282206001600160a01b039091168252855220549151911615158152f35b50346104015781600319360112610401576020906001549051908152f35b5034610401576020366003190112610401577f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c60043591610c4d610eda565b60025492610c5a84610ea8565b80600255816001558552600560205280828620557f5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc982600254958151908152866020820152a18151908582526020820152a280f35b8390346104015781600319360112610401576020906002548152f35b602435906001600160a01b0382168203610ce157565b5f80fd5b600435906001600160a01b0382168203610ce157565b9181601f84011215610ce15782359167ffffffffffffffff8311610ce1576020808501948460051b010111610ce157565b15610d3357565b60405162461bcd60e51b815260206004820152600e60248201526d125b9d985b1a59081c195c9a5bd960921b6044820152606490fd5b67ffffffffffffffff8111610d7d57604052565b634e487b7160e01b5f52604160045260245ffd5b90601f8019910116810190811067ffffffffffffffff821117610d7d57604052565b15610dba57565b60405162461bcd60e51b815260206004820152601260248201527114195c9a5bd9081b9bdd081cd95d1d1b195960721b6044820152606490fd5b90929167ffffffffffffffff8411610d7d578360051b6040519260208094610e1e82850182610d91565b8097815201918101928311610ce157905b828210610e3c5750505050565b81358152908301908301610e2f565b15610e5257565b60405162461bcd60e51b815260206004820152600d60248201526c24b73b30b634b210383937b7b360991b6044820152606490fd5b91908201809211610e9457565b634e487b7160e01b5f52601160045260245ffd5b5f198114610e945760010190565b9190811015610ec65760051b0190565b634e487b7160e01b5f52603260045260245ffd5b5f546001600160a01b03163303610eed57565b60405163118cdaa760e01b8152336004820152602490fd5b60ff5f5460a01c16610f1357565b60405163d93c066560e01b8152600490fd5b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f006002815414610f555760029055565b604051633ee5aeb560e01b8152600490fd5b9091905f915b8151831015610fb5576020808460051b84010151915f8382105f14610fa457505f5252610f9e60405f205b92610ea8565b91610f6d565b90604092610f9e9483525220610f98565b915050149056fea2646970667358221220818b0740a5df9db9e77ab8501481405e2cea30c4e89dd4d37b336ea0e8b0a28e64736f6c63430008140033",
}

// ReferralRewardDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use ReferralRewardDistributorMetaData.ABI instead.
var ReferralRewardDistributorABI = ReferralRewardDistributorMetaData.ABI

// ReferralRewardDistributorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ReferralRewardDistributorMetaData.Bin instead.
var ReferralRewardDistributorBin = ReferralRewardDistributorMetaData.Bin

// DeployReferralRewardDistributor deploys a new Ethereum contract, binding an instance of ReferralRewardDistributor to it.
func DeployReferralRewardDistributor(auth *bind.TransactOpts, backend bind.ContractBackend, kawaiToken_ common.Address) (common.Address, *types.Transaction, *ReferralRewardDistributor, error) {
	parsed, err := ReferralRewardDistributorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReferralRewardDistributorBin), backend, kawaiToken_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReferralRewardDistributor{ReferralRewardDistributorCaller: ReferralRewardDistributorCaller{contract: contract}, ReferralRewardDistributorTransactor: ReferralRewardDistributorTransactor{contract: contract}, ReferralRewardDistributorFilterer: ReferralRewardDistributorFilterer{contract: contract}}, nil
}

// ReferralRewardDistributor is an auto generated Go binding around an Ethereum contract.
type ReferralRewardDistributor struct {
	ReferralRewardDistributorCaller     // Read-only binding to the contract
	ReferralRewardDistributorTransactor // Write-only binding to the contract
	ReferralRewardDistributorFilterer   // Log filterer for contract events
}

// ReferralRewardDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReferralRewardDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralRewardDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReferralRewardDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralRewardDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReferralRewardDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralRewardDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReferralRewardDistributorSession struct {
	Contract     *ReferralRewardDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ReferralRewardDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReferralRewardDistributorCallerSession struct {
	Contract *ReferralRewardDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ReferralRewardDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReferralRewardDistributorTransactorSession struct {
	Contract     *ReferralRewardDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ReferralRewardDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReferralRewardDistributorRaw struct {
	Contract *ReferralRewardDistributor // Generic contract binding to access the raw methods on
}

// ReferralRewardDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReferralRewardDistributorCallerRaw struct {
	Contract *ReferralRewardDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// ReferralRewardDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReferralRewardDistributorTransactorRaw struct {
	Contract *ReferralRewardDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReferralRewardDistributor creates a new instance of ReferralRewardDistributor, bound to a specific deployed contract.
func NewReferralRewardDistributor(address common.Address, backend bind.ContractBackend) (*ReferralRewardDistributor, error) {
	contract, err := bindReferralRewardDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributor{ReferralRewardDistributorCaller: ReferralRewardDistributorCaller{contract: contract}, ReferralRewardDistributorTransactor: ReferralRewardDistributorTransactor{contract: contract}, ReferralRewardDistributorFilterer: ReferralRewardDistributorFilterer{contract: contract}}, nil
}

// NewReferralRewardDistributorCaller creates a new read-only instance of ReferralRewardDistributor, bound to a specific deployed contract.
func NewReferralRewardDistributorCaller(address common.Address, caller bind.ContractCaller) (*ReferralRewardDistributorCaller, error) {
	contract, err := bindReferralRewardDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorCaller{contract: contract}, nil
}

// NewReferralRewardDistributorTransactor creates a new write-only instance of ReferralRewardDistributor, bound to a specific deployed contract.
func NewReferralRewardDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*ReferralRewardDistributorTransactor, error) {
	contract, err := bindReferralRewardDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorTransactor{contract: contract}, nil
}

// NewReferralRewardDistributorFilterer creates a new log filterer instance of ReferralRewardDistributor, bound to a specific deployed contract.
func NewReferralRewardDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*ReferralRewardDistributorFilterer, error) {
	contract, err := bindReferralRewardDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorFilterer{contract: contract}, nil
}

// bindReferralRewardDistributor binds a generic wrapper to an already deployed contract.
func bindReferralRewardDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReferralRewardDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReferralRewardDistributor *ReferralRewardDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReferralRewardDistributor.Contract.ReferralRewardDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReferralRewardDistributor *ReferralRewardDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.ReferralRewardDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReferralRewardDistributor *ReferralRewardDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.ReferralRewardDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReferralRewardDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.contract.Transact(opts, method, params...)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) CurrentPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "currentPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) CurrentPeriod() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.CurrentPeriod(&_ReferralRewardDistributor.CallOpts)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) CurrentPeriod() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.CurrentPeriod(&_ReferralRewardDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 kawaiDistributed, uint256 referrers)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) GetStats(opts *bind.CallOpts) (struct {
	Period           *big.Int
	KawaiDistributed *big.Int
	Referrers        *big.Int
}, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "getStats")

	outstruct := new(struct {
		Period           *big.Int
		KawaiDistributed *big.Int
		Referrers        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Period = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.KawaiDistributed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Referrers = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 kawaiDistributed, uint256 referrers)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) GetStats() (struct {
	Period           *big.Int
	KawaiDistributed *big.Int
	Referrers        *big.Int
}, error) {
	return _ReferralRewardDistributor.Contract.GetStats(&_ReferralRewardDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 kawaiDistributed, uint256 referrers)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) GetStats() (struct {
	Period           *big.Int
	KawaiDistributed *big.Int
	Referrers        *big.Int
}, error) {
	return _ReferralRewardDistributor.Contract.GetStats(&_ReferralRewardDistributor.CallOpts)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) HasClaimed(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "hasClaimed", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _ReferralRewardDistributor.Contract.HasClaimed(&_ReferralRewardDistributor.CallOpts, arg0, arg1)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _ReferralRewardDistributor.Contract.HasClaimed(&_ReferralRewardDistributor.CallOpts, arg0, arg1)
}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) HasClaimedAnyPeriod(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "hasClaimedAnyPeriod", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) HasClaimedAnyPeriod(arg0 common.Address) (bool, error) {
	return _ReferralRewardDistributor.Contract.HasClaimedAnyPeriod(&_ReferralRewardDistributor.CallOpts, arg0)
}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) HasClaimedAnyPeriod(arg0 common.Address) (bool, error) {
	return _ReferralRewardDistributor.Contract.HasClaimedAnyPeriod(&_ReferralRewardDistributor.CallOpts, arg0)
}

// HasClaimedPeriod is a free data retrieval call binding the contract method 0x3f08ccd0.
//
// Solidity: function hasClaimedPeriod(uint256 period, address user) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) HasClaimedPeriod(opts *bind.CallOpts, period *big.Int, user common.Address) (bool, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "hasClaimedPeriod", period, user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimedPeriod is a free data retrieval call binding the contract method 0x3f08ccd0.
//
// Solidity: function hasClaimedPeriod(uint256 period, address user) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) HasClaimedPeriod(period *big.Int, user common.Address) (bool, error) {
	return _ReferralRewardDistributor.Contract.HasClaimedPeriod(&_ReferralRewardDistributor.CallOpts, period, user)
}

// HasClaimedPeriod is a free data retrieval call binding the contract method 0x3f08ccd0.
//
// Solidity: function hasClaimedPeriod(uint256 period, address user) view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) HasClaimedPeriod(period *big.Int, user common.Address) (bool, error) {
	return _ReferralRewardDistributor.Contract.HasClaimedPeriod(&_ReferralRewardDistributor.CallOpts, period, user)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) KawaiToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "kawaiToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) KawaiToken() (common.Address, error) {
	return _ReferralRewardDistributor.Contract.KawaiToken(&_ReferralRewardDistributor.CallOpts)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) KawaiToken() (common.Address, error) {
	return _ReferralRewardDistributor.Contract.KawaiToken(&_ReferralRewardDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) MerkleRoot() ([32]byte, error) {
	return _ReferralRewardDistributor.Contract.MerkleRoot(&_ReferralRewardDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) MerkleRoot() ([32]byte, error) {
	return _ReferralRewardDistributor.Contract.MerkleRoot(&_ReferralRewardDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) Owner() (common.Address, error) {
	return _ReferralRewardDistributor.Contract.Owner(&_ReferralRewardDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) Owner() (common.Address, error) {
	return _ReferralRewardDistributor.Contract.Owner(&_ReferralRewardDistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) Paused() (bool, error) {
	return _ReferralRewardDistributor.Contract.Paused(&_ReferralRewardDistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) Paused() (bool, error) {
	return _ReferralRewardDistributor.Contract.Paused(&_ReferralRewardDistributor.CallOpts)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) PeriodMerkleRoots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "periodMerkleRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _ReferralRewardDistributor.Contract.PeriodMerkleRoots(&_ReferralRewardDistributor.CallOpts, arg0)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _ReferralRewardDistributor.Contract.PeriodMerkleRoots(&_ReferralRewardDistributor.CallOpts, arg0)
}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) TotalKawaiDistributed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "totalKawaiDistributed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) TotalKawaiDistributed() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.TotalKawaiDistributed(&_ReferralRewardDistributor.CallOpts)
}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) TotalKawaiDistributed() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.TotalKawaiDistributed(&_ReferralRewardDistributor.CallOpts)
}

// TotalReferrers is a free data retrieval call binding the contract method 0x62d03cb7.
//
// Solidity: function totalReferrers() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) TotalReferrers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "totalReferrers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalReferrers is a free data retrieval call binding the contract method 0x62d03cb7.
//
// Solidity: function totalReferrers() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) TotalReferrers() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.TotalReferrers(&_ReferralRewardDistributor.CallOpts)
}

// TotalReferrers is a free data retrieval call binding the contract method 0x62d03cb7.
//
// Solidity: function totalReferrers() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) TotalReferrers() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.TotalReferrers(&_ReferralRewardDistributor.CallOpts)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) AdvancePeriod(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "advancePeriod", _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.AdvancePeriod(&_ReferralRewardDistributor.TransactOpts, _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.AdvancePeriod(&_ReferralRewardDistributor.TransactOpts, _merkleRoot)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) ClaimMultiplePeriods(opts *bind.TransactOpts, periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "claimMultiplePeriods", periods, kawaiAmounts, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) ClaimMultiplePeriods(periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.ClaimMultiplePeriods(&_ReferralRewardDistributor.TransactOpts, periods, kawaiAmounts, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) ClaimMultiplePeriods(periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.ClaimMultiplePeriods(&_ReferralRewardDistributor.TransactOpts, periods, kawaiAmounts, merkleProofs)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x66351b94.
//
// Solidity: function claimRewards(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) ClaimRewards(opts *bind.TransactOpts, period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "claimRewards", period, kawaiAmount, merkleProof)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x66351b94.
//
// Solidity: function claimRewards(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) ClaimRewards(period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.ClaimRewards(&_ReferralRewardDistributor.TransactOpts, period, kawaiAmount, merkleProof)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x66351b94.
//
// Solidity: function claimRewards(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) ClaimRewards(period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.ClaimRewards(&_ReferralRewardDistributor.TransactOpts, period, kawaiAmount, merkleProof)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) Pause() (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.Pause(&_ReferralRewardDistributor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) Pause() (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.Pause(&_ReferralRewardDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.RenounceOwnership(&_ReferralRewardDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.RenounceOwnership(&_ReferralRewardDistributor.TransactOpts)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) SetMerkleRoot(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "setMerkleRoot", _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.SetMerkleRoot(&_ReferralRewardDistributor.TransactOpts, _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.SetMerkleRoot(&_ReferralRewardDistributor.TransactOpts, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) SetPeriodMerkleRoot(opts *bind.TransactOpts, period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "setPeriodMerkleRoot", period, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) SetPeriodMerkleRoot(period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.SetPeriodMerkleRoot(&_ReferralRewardDistributor.TransactOpts, period, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) SetPeriodMerkleRoot(period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.SetPeriodMerkleRoot(&_ReferralRewardDistributor.TransactOpts, period, _merkleRoot)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.TransferOwnership(&_ReferralRewardDistributor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.TransferOwnership(&_ReferralRewardDistributor.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardDistributor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) Unpause() (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.Unpause(&_ReferralRewardDistributor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ReferralRewardDistributor *ReferralRewardDistributorTransactorSession) Unpause() (*types.Transaction, error) {
	return _ReferralRewardDistributor.Contract.Unpause(&_ReferralRewardDistributor.TransactOpts)
}

// ReferralRewardDistributorMerkleRootUpdatedIterator is returned from FilterMerkleRootUpdated and is used to iterate over the raw logs and unpacked data for MerkleRootUpdated events raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorMerkleRootUpdatedIterator struct {
	Event *ReferralRewardDistributorMerkleRootUpdated // Event containing the contract specifics and raw log

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
func (it *ReferralRewardDistributorMerkleRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardDistributorMerkleRootUpdated)
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
		it.Event = new(ReferralRewardDistributorMerkleRootUpdated)
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
func (it *ReferralRewardDistributorMerkleRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardDistributorMerkleRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardDistributorMerkleRootUpdated represents a MerkleRootUpdated event raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorMerkleRootUpdated struct {
	Period  *big.Int
	OldRoot [32]byte
	NewRoot [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootUpdated is a free log retrieval operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) FilterMerkleRootUpdated(opts *bind.FilterOpts, period []*big.Int) (*ReferralRewardDistributorMerkleRootUpdatedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _ReferralRewardDistributor.contract.FilterLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorMerkleRootUpdatedIterator{contract: _ReferralRewardDistributor.contract, event: "MerkleRootUpdated", logs: logs, sub: sub}, nil
}

// WatchMerkleRootUpdated is a free log subscription operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) WatchMerkleRootUpdated(opts *bind.WatchOpts, sink chan<- *ReferralRewardDistributorMerkleRootUpdated, period []*big.Int) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _ReferralRewardDistributor.contract.WatchLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardDistributorMerkleRootUpdated)
				if err := _ReferralRewardDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
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
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) ParseMerkleRootUpdated(log types.Log) (*ReferralRewardDistributorMerkleRootUpdated, error) {
	event := new(ReferralRewardDistributorMerkleRootUpdated)
	if err := _ReferralRewardDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardDistributorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorOwnershipTransferredIterator struct {
	Event *ReferralRewardDistributorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ReferralRewardDistributorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardDistributorOwnershipTransferred)
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
		it.Event = new(ReferralRewardDistributorOwnershipTransferred)
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
func (it *ReferralRewardDistributorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardDistributorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardDistributorOwnershipTransferred represents a OwnershipTransferred event raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ReferralRewardDistributorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReferralRewardDistributor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorOwnershipTransferredIterator{contract: _ReferralRewardDistributor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ReferralRewardDistributorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReferralRewardDistributor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardDistributorOwnershipTransferred)
				if err := _ReferralRewardDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) ParseOwnershipTransferred(log types.Log) (*ReferralRewardDistributorOwnershipTransferred, error) {
	event := new(ReferralRewardDistributorOwnershipTransferred)
	if err := _ReferralRewardDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardDistributorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorPausedIterator struct {
	Event *ReferralRewardDistributorPaused // Event containing the contract specifics and raw log

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
func (it *ReferralRewardDistributorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardDistributorPaused)
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
		it.Event = new(ReferralRewardDistributorPaused)
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
func (it *ReferralRewardDistributorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardDistributorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardDistributorPaused represents a Paused event raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) FilterPaused(opts *bind.FilterOpts) (*ReferralRewardDistributorPausedIterator, error) {

	logs, sub, err := _ReferralRewardDistributor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorPausedIterator{contract: _ReferralRewardDistributor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ReferralRewardDistributorPaused) (event.Subscription, error) {

	logs, sub, err := _ReferralRewardDistributor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardDistributorPaused)
				if err := _ReferralRewardDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) ParsePaused(log types.Log) (*ReferralRewardDistributorPaused, error) {
	event := new(ReferralRewardDistributorPaused)
	if err := _ReferralRewardDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardDistributorPeriodAdvancedIterator is returned from FilterPeriodAdvanced and is used to iterate over the raw logs and unpacked data for PeriodAdvanced events raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorPeriodAdvancedIterator struct {
	Event *ReferralRewardDistributorPeriodAdvanced // Event containing the contract specifics and raw log

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
func (it *ReferralRewardDistributorPeriodAdvancedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardDistributorPeriodAdvanced)
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
		it.Event = new(ReferralRewardDistributorPeriodAdvanced)
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
func (it *ReferralRewardDistributorPeriodAdvancedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardDistributorPeriodAdvancedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardDistributorPeriodAdvanced represents a PeriodAdvanced event raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorPeriodAdvanced struct {
	OldPeriod *big.Int
	NewPeriod *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPeriodAdvanced is a free log retrieval operation binding the contract event 0x5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc9.
//
// Solidity: event PeriodAdvanced(uint256 oldPeriod, uint256 newPeriod)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) FilterPeriodAdvanced(opts *bind.FilterOpts) (*ReferralRewardDistributorPeriodAdvancedIterator, error) {

	logs, sub, err := _ReferralRewardDistributor.contract.FilterLogs(opts, "PeriodAdvanced")
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorPeriodAdvancedIterator{contract: _ReferralRewardDistributor.contract, event: "PeriodAdvanced", logs: logs, sub: sub}, nil
}

// WatchPeriodAdvanced is a free log subscription operation binding the contract event 0x5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc9.
//
// Solidity: event PeriodAdvanced(uint256 oldPeriod, uint256 newPeriod)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) WatchPeriodAdvanced(opts *bind.WatchOpts, sink chan<- *ReferralRewardDistributorPeriodAdvanced) (event.Subscription, error) {

	logs, sub, err := _ReferralRewardDistributor.contract.WatchLogs(opts, "PeriodAdvanced")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardDistributorPeriodAdvanced)
				if err := _ReferralRewardDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
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

// ParsePeriodAdvanced is a log parse operation binding the contract event 0x5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc9.
//
// Solidity: event PeriodAdvanced(uint256 oldPeriod, uint256 newPeriod)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) ParsePeriodAdvanced(log types.Log) (*ReferralRewardDistributorPeriodAdvanced, error) {
	event := new(ReferralRewardDistributorPeriodAdvanced)
	if err := _ReferralRewardDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardDistributorRewardsClaimedIterator is returned from FilterRewardsClaimed and is used to iterate over the raw logs and unpacked data for RewardsClaimed events raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorRewardsClaimedIterator struct {
	Event *ReferralRewardDistributorRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *ReferralRewardDistributorRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardDistributorRewardsClaimed)
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
		it.Event = new(ReferralRewardDistributorRewardsClaimed)
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
func (it *ReferralRewardDistributorRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardDistributorRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardDistributorRewardsClaimed represents a RewardsClaimed event raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorRewardsClaimed struct {
	Period      *big.Int
	User        common.Address
	KawaiAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRewardsClaimed is a free log retrieval operation binding the contract event 0x3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb796968.
//
// Solidity: event RewardsClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, period []*big.Int, user []common.Address) (*ReferralRewardDistributorRewardsClaimedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardDistributor.contract.FilterLogs(opts, "RewardsClaimed", periodRule, userRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorRewardsClaimedIterator{contract: _ReferralRewardDistributor.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardsClaimed is a free log subscription operation binding the contract event 0x3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb796968.
//
// Solidity: event RewardsClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *ReferralRewardDistributorRewardsClaimed, period []*big.Int, user []common.Address) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardDistributor.contract.WatchLogs(opts, "RewardsClaimed", periodRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardDistributorRewardsClaimed)
				if err := _ReferralRewardDistributor.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

// ParseRewardsClaimed is a log parse operation binding the contract event 0x3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb796968.
//
// Solidity: event RewardsClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) ParseRewardsClaimed(log types.Log) (*ReferralRewardDistributorRewardsClaimed, error) {
	event := new(ReferralRewardDistributorRewardsClaimed)
	if err := _ReferralRewardDistributor.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardDistributorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorUnpausedIterator struct {
	Event *ReferralRewardDistributorUnpaused // Event containing the contract specifics and raw log

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
func (it *ReferralRewardDistributorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardDistributorUnpaused)
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
		it.Event = new(ReferralRewardDistributorUnpaused)
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
func (it *ReferralRewardDistributorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardDistributorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardDistributorUnpaused represents a Unpaused event raised by the ReferralRewardDistributor contract.
type ReferralRewardDistributorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ReferralRewardDistributorUnpausedIterator, error) {

	logs, sub, err := _ReferralRewardDistributor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ReferralRewardDistributorUnpausedIterator{contract: _ReferralRewardDistributor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ReferralRewardDistributorUnpaused) (event.Subscription, error) {

	logs, sub, err := _ReferralRewardDistributor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardDistributorUnpaused)
				if err := _ReferralRewardDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_ReferralRewardDistributor *ReferralRewardDistributorFilterer) ParseUnpaused(log types.Log) (*ReferralRewardDistributorUnpaused, error) {
	event := new(ReferralRewardDistributorUnpaused)
	if err := _ReferralRewardDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
