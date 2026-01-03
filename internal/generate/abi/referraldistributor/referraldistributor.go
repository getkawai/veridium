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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"kawaiToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"advancePeriod\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimMultiplePeriods\",\"inputs\":[{\"name\":\"periods\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"kawaiAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"merkleProofs\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimRewards\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStats\",\"inputs\":[],\"outputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"referrers\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedAnyPeriod\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedPeriod\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kawaiToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"periodMerkleRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalKawaiDistributed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalReferrers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodAdvanced\",\"inputs\":[{\"name\":\"oldPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051611f0e380380611f0e833981810160405281019061003191906102d7565b335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100999190610311565b60405180910390fd5b6100b18161018660201b60201c565b5060016100d06100c561024760201b60201c565b61027060201b60201c565b5f01819055505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610144576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161013b90610384565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250506001600281905550506103a2565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6102a68261027d565b9050919050565b6102b68161029c565b81146102c0575f5ffd5b50565b5f815190506102d1816102ad565b92915050565b5f602082840312156102ec576102eb610279565b5b5f6102f9848285016102c3565b91505092915050565b61030b8161029c565b82525050565b5f6020820190506103245f830184610302565b92915050565b5f82825260208201905092915050565b7f496e76616c6964204b41574149206164647265737300000000000000000000005f82015250565b5f61036e60158361032a565b91506103798261033a565b602082019050919050565b5f6020820190508181035f83015261039b81610362565b9050919050565b608051611b466103c85f395f81816106ba015281816109880152610d180152611b465ff3fe608060405234801561000f575f5ffd5b5060043610610109575f3560e01c8063727a7c5a116100a0578063adeacbd31161006f578063adeacbd314610291578063c59d4847146102c1578063cb56cd4f146102e1578063f2fde38b146102ff578063f75cc2b91461031b57610109565b8063727a7c5a146101f75780637cb6475914610227578063873f6f9e146102435780638da5cb5b1461027357610109565b80635869bc5a116100dc5780635869bc5a1461019557806362d03cb7146101b357806366351b94146101d1578063715018a6146101ed57610109565b8063060406181461010d5780630ae654031461012b5780632eb4a7ab146101475780633f08ccd014610165575b5f5ffd5b610115610337565b6040516101229190611149565b60405180910390f35b6101456004803603810190610140919061119d565b61033d565b005b61014f6103fe565b60405161015c91906111d7565b60405180910390f35b61017f600480360381019061017a9190611274565b610404565b60405161018c91906112cc565b60405180910390f35b61019d610466565b6040516101aa9190611149565b60405180910390f35b6101bb61046c565b6040516101c89190611149565b60405180910390f35b6101eb60048036038101906101e69190611346565b610472565b005b6101f5610871565b005b610211600480360381019061020c91906113b7565b610884565b60405161021e91906111d7565b60405180910390f35b610241600480360381019061023c919061119d565b610899565b005b61025d60048036038101906102589190611274565b610901565b60405161026a91906112cc565b60405180910390f35b61027b61092b565b60405161028891906113f1565b60405180910390f35b6102ab60048036038101906102a6919061140a565b610952565b6040516102b891906112cc565b60405180910390f35b6102c961096f565b6040516102d893929190611435565b60405180910390f35b6102e9610986565b6040516102f691906114c5565b60405180910390f35b6103196004803603810190610314919061140a565b6109aa565b005b61033560048036038101906103309190611588565b610a2e565b005b60025481565b610345610e72565b5f600254905060025f81548092919061035d90611665565b9190505550816001819055508160055f60025481526020019081526020015f20819055507f5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc9816002546040516103b49291906116ac565b60405180910390a16002547f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c5f5f1b846040516103f29291906116d3565b60405180910390a25050565b60015481565b5f60035f8481526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b60065481565b60075481565b61047a610ef9565b6002548411156104bf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104b690611754565b60405180910390fd5b60035f8581526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1615610558576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161054f906117bc565b60405180910390fd5b5f831161059a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059190611824565b60405180910390fd5b5f8433856040516020016105b0939291906118a7565b6040516020818303038152906040528051906020012090506106158383808060200260200160405190810160405280939291908181526020018383602002808284375f81840152601f19601f8201169050808301925050505050505060015483610f1b565b610654576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064b9061192d565b60405180910390fd5b600160035f8781526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166340c10f1933866040518363ffffffff1660e01b815260040161071392919061194b565b5f604051808303815f87803b15801561072a575f5ffd5b505af115801561073c573d5f5f3e3d5ffd5b505050508360065f8282546107519190611972565b9250508190555060045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1661081357600160045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555060075f81548092919061080d90611665565b91905055505b3373ffffffffffffffffffffffffffffffffffffffff16857f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb7969688660405161085a9190611149565b60405180910390a35061086b610f31565b50505050565b610879610e72565b6108825f610f4b565b565b6005602052805f5260405f205f915090505481565b6108a1610e72565b6002547f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c600154836040516108d79291906116d3565b60405180910390a2806001819055508060055f60025481526020019081526020015f208190555050565b6003602052815f5260405f20602052805f5260405f205f915091509054906101000a900460ff1681565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6004602052805f5260405f205f915054906101000a900460ff1681565b5f5f5f600254600654600754925092509250909192565b7f000000000000000000000000000000000000000000000000000000000000000081565b6109b2610e72565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610a22575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610a1991906113f1565b60405180910390fd5b610a2b81610f4b565b50565b610a36610ef9565b8383905086869050148015610a5057508181905086869050145b610a8f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a86906119ef565b60405180910390fd5b5f5f90505b86869050811015610e61575f878783818110610ab357610ab2611a0d565b5b90506020020135905060035f8281526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1615610b205750610e54565b600254811115610b65576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b5c90611754565b60405180910390fd5b5f8133888886818110610b7b57610b7a611a0d565b5b90506020020135604051602001610b94939291906118a7565b6040516020818303038152906040528051906020012090505f60055f8481526020019081526020015f205490505f5f1b8103610c05576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bfc90611a84565b60405180910390fd5b610c73868686818110610c1b57610c1a611a0d565b5b9050602002810190610c2d9190611aae565b808060200260200160405190810160405280939291908181526020018383602002808284375f81840152601f19601f820116905080830192505050505050508284610f1b565b610cb2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ca99061192d565b60405180910390fd5b600160035f8581526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166340c10f19338a8a88818110610d6657610d65611a0d565b5b905060200201356040518363ffffffff1660e01b8152600401610d8a92919061194b565b5f604051808303815f87803b158015610da1575f5ffd5b505af1158015610db3573d5f5f3e3d5ffd5b50505050878785818110610dca57610dc9611a0d565b5b9050602002013560065f828254610de19190611972565b925050819055503373ffffffffffffffffffffffffffffffffffffffff16837f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb7969688a8a88818110610e3457610e33611a0d565b5b90506020020135604051610e489190611149565b60405180910390a35050505b8080600101915050610a94565b50610e6a610f31565b505050505050565b610e7a61100c565b73ffffffffffffffffffffffffffffffffffffffff16610e9861092b565b73ffffffffffffffffffffffffffffffffffffffff1614610ef757610ebb61100c565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401610eee91906113f1565b60405180910390fd5b565b610f01611013565b6002610f13610f0e611054565b61107d565b5f0181905550565b5f82610f278584611086565b1490509392505050565b6001610f43610f3e611054565b61107d565b5f0181905550565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f33905090565b61101b6110d7565b15611052576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5f8290505f5f90505b84518110156110cc576110bd828683815181106110b0576110af611a0d565b5b60200260200101516110f3565b91508080600101915050611090565b508091505092915050565b5f60026110ea6110e5611054565b61107d565b5f015414905090565b5f81831061110a57611105828461111d565b611115565b611114838361111d565b5b905092915050565b5f825f528160205260405f20905092915050565b5f819050919050565b61114381611131565b82525050565b5f60208201905061115c5f83018461113a565b92915050565b5f5ffd5b5f5ffd5b5f819050919050565b61117c8161116a565b8114611186575f5ffd5b50565b5f8135905061119781611173565b92915050565b5f602082840312156111b2576111b1611162565b5b5f6111bf84828501611189565b91505092915050565b6111d18161116a565b82525050565b5f6020820190506111ea5f8301846111c8565b92915050565b6111f981611131565b8114611203575f5ffd5b50565b5f81359050611214816111f0565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6112438261121a565b9050919050565b61125381611239565b811461125d575f5ffd5b50565b5f8135905061126e8161124a565b92915050565b5f5f6040838503121561128a57611289611162565b5b5f61129785828601611206565b92505060206112a885828601611260565b9150509250929050565b5f8115159050919050565b6112c6816112b2565b82525050565b5f6020820190506112df5f8301846112bd565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f840112611306576113056112e5565b5b8235905067ffffffffffffffff811115611323576113226112e9565b5b60208301915083602082028301111561133f5761133e6112ed565b5b9250929050565b5f5f5f5f6060858703121561135e5761135d611162565b5b5f61136b87828801611206565b945050602061137c87828801611206565b935050604085013567ffffffffffffffff81111561139d5761139c611166565b5b6113a9878288016112f1565b925092505092959194509250565b5f602082840312156113cc576113cb611162565b5b5f6113d984828501611206565b91505092915050565b6113eb81611239565b82525050565b5f6020820190506114045f8301846113e2565b92915050565b5f6020828403121561141f5761141e611162565b5b5f61142c84828501611260565b91505092915050565b5f6060820190506114485f83018661113a565b611455602083018561113a565b611462604083018461113a565b949350505050565b5f819050919050565b5f61148d6114886114838461121a565b61146a565b61121a565b9050919050565b5f61149e82611473565b9050919050565b5f6114af82611494565b9050919050565b6114bf816114a5565b82525050565b5f6020820190506114d85f8301846114b6565b92915050565b5f5f83601f8401126114f3576114f26112e5565b5b8235905067ffffffffffffffff8111156115105761150f6112e9565b5b60208301915083602082028301111561152c5761152b6112ed565b5b9250929050565b5f5f83601f840112611548576115476112e5565b5b8235905067ffffffffffffffff811115611565576115646112e9565b5b602083019150836020820283011115611581576115806112ed565b5b9250929050565b5f5f5f5f5f5f606087890312156115a2576115a1611162565b5b5f87013567ffffffffffffffff8111156115bf576115be611166565b5b6115cb89828a016114de565b9650965050602087013567ffffffffffffffff8111156115ee576115ed611166565b5b6115fa89828a016114de565b9450945050604087013567ffffffffffffffff81111561161d5761161c611166565b5b61162989828a01611533565b92509250509295509295509295565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61166f82611131565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036116a1576116a0611638565b5b600182019050919050565b5f6040820190506116bf5f83018561113a565b6116cc602083018461113a565b9392505050565b5f6040820190506116e65f8301856111c8565b6116f360208301846111c8565b9392505050565b5f82825260208201905092915050565b7f496e76616c696420706572696f640000000000000000000000000000000000005f82015250565b5f61173e600e836116fa565b91506117498261170a565b602082019050919050565b5f6020820190508181035f83015261176b81611732565b9050919050565b7f416c726561647920636c61696d656420666f72207468697320706572696f64005f82015250565b5f6117a6601f836116fa565b91506117b182611772565b602082019050919050565b5f6020820190508181035f8301526117d38161179a565b9050919050565b7f4e6f207265776172647320746f20636c61696d000000000000000000000000005f82015250565b5f61180e6013836116fa565b9150611819826117da565b602082019050919050565b5f6020820190508181035f83015261183b81611802565b9050919050565b5f819050919050565b61185c61185782611131565b611842565b82525050565b5f8160601b9050919050565b5f61187882611862565b9050919050565b5f6118898261186e565b9050919050565b6118a161189c82611239565b61187f565b82525050565b5f6118b2828661184b565b6020820191506118c28285611890565b6014820191506118d2828461184b565b602082019150819050949350505050565b7f496e76616c69642070726f6f66000000000000000000000000000000000000005f82015250565b5f611917600d836116fa565b9150611922826118e3565b602082019050919050565b5f6020820190508181035f8301526119448161190b565b9050919050565b5f60408201905061195e5f8301856113e2565b61196b602083018461113a565b9392505050565b5f61197c82611131565b915061198783611131565b925082820190508082111561199f5761199e611638565b5b92915050565b7f4172726179206c656e677468206d69736d6174636800000000000000000000005f82015250565b5f6119d96015836116fa565b91506119e4826119a5565b602082019050919050565b5f6020820190508181035f830152611a06816119cd565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f506572696f64206e6f7420736574746c656400000000000000000000000000005f82015250565b5f611a6e6012836116fa565b9150611a7982611a3a565b602082019050919050565b5f6020820190508181035f830152611a9b81611a62565b9050919050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83356001602003843603038112611aca57611ac9611aa2565b5b80840192508235915067ffffffffffffffff821115611aec57611aeb611aa6565b5b602083019250602082023603831315611b0857611b07611aaa565b5b50925092905056fea2646970667358221220426ddf404a97539be57ec1619f6a6ca436b249ab03a65437e6ff8937f7301a3b64736f6c634300081e0033",
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
