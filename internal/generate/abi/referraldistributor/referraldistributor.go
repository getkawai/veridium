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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"kawaiToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"TOTAL_ALLOCATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"advancePeriod\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimMultiplePeriods\",\"inputs\":[{\"name\":\"periods\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"kawaiAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"merkleProofs\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimRewards\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStats\",\"inputs\":[],\"outputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingAllocation\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"referrers\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedAnyPeriod\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedPeriod\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kawaiToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"periodMerkleRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeriodMerkleRoot\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalKawaiDistributed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalReferrers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodAdvanced\",\"inputs\":[{\"name\":\"oldPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a03461014657601f6112a238819003918201601f19168301916001600160401b0383118484101761014a5780849260209460405283398101031261014657516001600160a01b038082169182900361014657331561012e575f543360018060a01b03198216175f55604051913391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005581156100ec57506080526001600255604051611143908161015f82396080518181816101e50152818161058901526109ee0152f35b62461bcd60e51b815260206004820152601560248201527f496e76616c6964204b41574149206164647265737300000000000000000000006044820152606490fd5b604051631e4fbdf760e01b81525f6004820152602490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6080806040526004361015610012575f80fd5b5f90813560e01c9081630604061814610d72575080630ae6540314610cd05780630dac7bea14610caa5780632eb4a7ab14610c8c5780633f08ccd014610c435780633f4ba83a14610bd25780635869bc5a14610bb45780635c975abb14610b8f57806362d03cb714610b7157806366351b94146108cb578063715018a614610872578063727a7c5a146108485780637cb64759146107e25780638456cb5914610781578063873f6f9e146107375780638da5cb5b14610710578063adeacbd3146106d1578063c40c91bd14610618578063c59d4847146105b8578063cb56cd4f14610573578063f2fde38b146104e65763f75cc2b914610110575f80fd5b346104e35760603660031901126104e35760043567ffffffffffffffff81116102cb57610141903690600401610dbe565b9060243567ffffffffffffffff81116104df57610162903690600401610dbe565b60443567ffffffffffffffff81116104db57610182903690600401610dbe565b929061018c611076565b610194611056565b828614806104d2575b15610495579493929190869587955b8087106102da5788886101c0811515610e2c565b6101e26af8277896582678ac0000006101db83600654610e6e565b1115610e8f565b817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156102cb576040516340c10f1960e01b8152336004820152602481018490529082908290604490829084905af180156102cf576102b7575b505061025590600654610e6e565b600655338152600460205260408120805460ff811615610298575b8260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b60ff191660011790556007546102ad90610ff9565b6007558180610270565b6102c090610edb565b6102cb578183610247565b5080fd5b6040513d84823e3d90fd5b909192939495966102ec888389611007565b3597888a52600360205260408a20335f5260205260ff60405f20541661048a5761031a6002548a1115610def565b8861036261037061032c848a8a611007565b60408051602081019586523360601b6bffffffffffffffffffffffff191691810191909152903560548201529182906074820190565b03601f198101835282610f03565b51902091898b52600560205260408b20549061038d821515610f25565b88831015610476578260051b860135601e198736030181121561047257860180359067ffffffffffffffff821161046e576020018160051b3603811361046e5761046395610424946103e76103ec936103f1953691610f66565b6110b8565b610fbd565b8a8c52600360205260408c20335f5260205260405f20600160ff1982541617905561041d838989611007565b3590610e6e565b98610430828888611007565b35906040519182527f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb79696860203393a3610ff9565b9594939291906101ac565b8d80fd5b8c80fd5b634e487b7160e01b8c52603260045260248cfd5b610463919850610ff9565b60405162461bcd60e51b8152602060048201526015602482015274082e4e4c2f240d8cadccee8d040dad2e6dac2e8c6d605b1b6044820152606490fd5b5083861461019d565b8580fd5b8380fd5b80fd5b50346104e35760203660031901126104e357610500610da8565b61050861102b565b6001600160a01b0390811690811561055a575f54826bffffffffffffffffffffffff60a01b8216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a380f35b604051631e4fbdf760e01b815260048101849052602490fd5b50346104e357806003193601126104e3576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b50346104e357806003193601126104e3576002546006546af8277896582678ac000000818103908111610604576080935060075491604051938452602084015260408301526060820152f35b634e487b7160e01b84526011600452602484fd5b50346104e35760403660031901126104e35760243560043561063861102b565b600254811161068c578083526005602052807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c6040808620548151908152856020820152a282526005602052604082205580f35b60405162461bcd60e51b815260206004820152601860248201527f43616e6e6f74207365742066757475726520706572696f6400000000000000006044820152606490fd5b50346104e35760203660031901126104e35760209060ff906040906001600160a01b036106fc610da8565b168152600484522054166040519015158152f35b50346104e357806003193601126104e357546040516001600160a01b039091168152602090f35b50346104e35760403660031901126104e3576040610753610d8e565b9160043581526003602052209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b50346104e357806003193601126104e35761079a61102b565b6107a2611056565b805460ff60a01b1916600160a01b1781556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602090a180f35b50346104e35760203660031901126104e3576004356107ff61102b565b600254807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c60406001548151908152856020820152a28160015582526005602052604082205580f35b50346104e35760203660031901126104e35760406020916004358152600583522054604051908152f35b50346104e357806003193601126104e35761088b61102b565b5f80546001600160a01b0319811682556001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b50346104e35760603660031901126104e35760243560043560443567ffffffffffffffff81116104df57610903903690600401610dbe565b9061090c611076565b610914611056565b610922600254841115610def565b8285526020916003835260408620338752835260ff604087205416610b2c576109cb916103ec91610954871515610e2c565b61096f6af8277896582678ac0000006101db89600654610e6e565b6040518581018781523360601b6bffffffffffffffffffffffff19166020820152603481018990526103e7916109a88160548401610362565b51902092878a526005875260408a2054926109c4841515610f25565b3691610f66565b81845260038152604080852033865282528420805460ff199081166001179091557f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104db576040516340c10f1960e01b8152336004820152602481018690529086908290604490829084905af18015610b2157610aec575b50907f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb79696891610a8085600654610e6e565b60065533865260048252604086209081549060ff821615610acf575b5050506040519384523393a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b60019116179055610ae1600754610ff9565b6007555f8080610a9c565b94610b197f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb796968939296610edb565b949091610a50565b6040513d88823e3d90fd5b60405162461bcd60e51b815260048101849052601f60248201527f416c726561647920636c61696d656420666f72207468697320706572696f64006044820152606490fd5b50346104e357806003193601126104e3576020600754604051908152f35b50346104e357806003193601126104e35760ff6020915460a01c166040519015158152f35b50346104e357806003193601126104e3576020600654604051908152f35b50346104e357806003193601126104e357610beb61102b565b805460ff8160a01c1615610c315760ff60a01b191681556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa90602090a180f35b604051638dfc202b60e01b8152600490fd5b50346104e35760403660031901126104e35760ff6040602092610c64610d8e565b6004358252600385528282206001600160a01b03909116825284522054604051911615158152f35b50346104e357806003193601126104e3576020600154604051908152f35b50346104e357806003193601126104e35760206040516af8277896582678ac0000008152f35b50346104e35760203660031901126104e357600435610ced61102b565b7f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c604060025492610d1d84610ff9565b80600255816001558552600560205280828620557f5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc982600254958151908152866020820152a18151908582526020820152a280f35b9050346102cb57816003193601126102cb576020906002548152f35b602435906001600160a01b0382168203610da457565b5f80fd5b600435906001600160a01b0382168203610da457565b9181601f84011215610da45782359167ffffffffffffffff8311610da4576020808501948460051b010111610da457565b15610df657565b60405162461bcd60e51b815260206004820152600e60248201526d125b9d985b1a59081c195c9a5bd960921b6044820152606490fd5b15610e3357565b60405162461bcd60e51b81526020600482015260136024820152724e6f207265776172647320746f20636c61696d60681b6044820152606490fd5b91908201809211610e7b57565b634e487b7160e01b5f52601160045260245ffd5b15610e9657565b60405162461bcd60e51b815260206004820152601860248201527f4578636565647320746f74616c20616c6c6f636174696f6e00000000000000006044820152606490fd5b67ffffffffffffffff8111610eef57604052565b634e487b7160e01b5f52604160045260245ffd5b90601f8019910116810190811067ffffffffffffffff821117610eef57604052565b15610f2c57565b60405162461bcd60e51b815260206004820152601260248201527114195c9a5bd9081b9bdd081cd95d1d1b195960721b6044820152606490fd5b90929167ffffffffffffffff8411610eef578360051b6040519260208094610f9082850182610f03565b8097815201918101928311610da457905b828210610fae5750505050565b81358152908301908301610fa1565b15610fc457565b60405162461bcd60e51b815260206004820152600d60248201526c24b73b30b634b210383937b7b360991b6044820152606490fd5b5f198114610e7b5760010190565b91908110156110175760051b0190565b634e487b7160e01b5f52603260045260245ffd5b5f546001600160a01b0316330361103e57565b60405163118cdaa760e01b8152336004820152602490fd5b60ff5f5460a01c1661106457565b60405163d93c066560e01b8152600490fd5b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146110a65760029055565b604051633ee5aeb560e01b8152600490fd5b9091905f915b8151831015611106576020808460051b84010151915f8382105f146110f557505f52526110ef60405f205b92610ff9565b916110be565b906040926110ef94835252206110e9565b915050149056fea26469706673582212209d867283b687ed0a86695a3fc00280ff3e5d864a73d3e056d13ef94fa285a12164736f6c63430008140033",
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

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) TOTALALLOCATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "TOTAL_ALLOCATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) TOTALALLOCATION() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.TOTALALLOCATION(&_ReferralRewardDistributor.CallOpts)
}

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) TOTALALLOCATION() (*big.Int, error) {
	return _ReferralRewardDistributor.Contract.TOTALALLOCATION(&_ReferralRewardDistributor.CallOpts)
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
// Solidity: function getStats() view returns(uint256 period, uint256 kawaiDistributed, uint256 remainingAllocation, uint256 referrers)
func (_ReferralRewardDistributor *ReferralRewardDistributorCaller) GetStats(opts *bind.CallOpts) (struct {
	Period              *big.Int
	KawaiDistributed    *big.Int
	RemainingAllocation *big.Int
	Referrers           *big.Int
}, error) {
	var out []interface{}
	err := _ReferralRewardDistributor.contract.Call(opts, &out, "getStats")

	outstruct := new(struct {
		Period              *big.Int
		KawaiDistributed    *big.Int
		RemainingAllocation *big.Int
		Referrers           *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Period = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.KawaiDistributed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RemainingAllocation = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Referrers = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 kawaiDistributed, uint256 remainingAllocation, uint256 referrers)
func (_ReferralRewardDistributor *ReferralRewardDistributorSession) GetStats() (struct {
	Period              *big.Int
	KawaiDistributed    *big.Int
	RemainingAllocation *big.Int
	Referrers           *big.Int
}, error) {
	return _ReferralRewardDistributor.Contract.GetStats(&_ReferralRewardDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 kawaiDistributed, uint256 remainingAllocation, uint256 referrers)
func (_ReferralRewardDistributor *ReferralRewardDistributorCallerSession) GetStats() (struct {
	Period              *big.Int
	KawaiDistributed    *big.Int
	RemainingAllocation *big.Int
	Referrers           *big.Int
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
