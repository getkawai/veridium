// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package miningdistributor

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

// MiningRewardDistributorMetaData contains all meta data concerning the MiningRewardDistributor contract.
var MiningRewardDistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"kawaiToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"advancePeriod\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimMultiplePeriods\",\"inputs\":[{\"name\":\"periods\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"contributorAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"developerAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"userAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"affiliatorAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"developers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"affiliators\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"merkleProofs\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimReward\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributorAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"developerAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"affiliatorAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"developer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"affiliator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStats\",\"inputs\":[],\"outputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributorRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"developerRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"affiliatorRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedPeriod\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kawaiToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"periodMerkleRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalAffiliatorRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalContributorRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalDeveloperRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalUserRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodAdvanced\",\"inputs\":[{\"name\":\"oldPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"contributorAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"developerAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"userAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"affiliatorAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"developer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"affiliator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a03461017057601f61167738819003918201601f19168301916001600160401b038311848410176101745780849260209460405283398101031261017057516001600160a01b0380821691829003610170573315610158575f543360018060a01b03198216175f55604051913391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00558115610116575060805260016002556040516114ee908161018982396080518181816101ae0152818161044c0152818161050c015281816105a901528181610659015281816111370152818161120f015281816112a901526113550152f35b62461bcd60e51b815260206004820152601560248201527f496e76616c6964204b41574149206164647265737300000000000000000000006044820152606490fd5b604051631e4fbdf760e01b81525f6004820152602490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6102e0806040526004361015610013575f80fd5b5f90813560e01c9081630604061814610d375750806306b7771714610d195780630ae6540314610c775780632eb4a7ab14610c595780633f08ccd014610753578063437f3d22146108a4578063715018a61461084a578063727a7c5a1461082057806377363251146108025780637cb647591461079c578063873f6f9e146107535780638da5cb5b1461072c57806395112df31461023b57806396e379f01461021d578063c59d4847146101dd578063cb56cd4f14610198578063f2fde38b146101055763fd8bfafc146100e5575f80fd5b346101025780600319360112610102576020600754604051908152f35b80fd5b5034610102576020366003190112610102576004356001600160a01b0381811691829003610194576101356113f6565b811561017b575f54826bffffffffffffffffffffffff60a01b8216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a380f35b604051631e4fbdf760e01b815260048101849052602490fd5b5f80fd5b50346101025780600319360112610102576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b503461010257806003193601126101025760a060025460055460065460075490600854926040519485526020850152604084015260608301526080820152f35b50346101025780600319360112610102576020600554604051908152f35b5034610102576101203660031901126101025760a4356001600160a01b03811690036101025760c4356001600160a01b038116908190036104e65760e4356001600160a01b03811690036104e657610104356001600160401b038111610728576102a9903690600401610d69565b6102b1611421565b6102c16002546004351115610d99565b600435845260036020526040842033855260205260ff6040852054166106e35761036591610360916102f4851515610dd6565b61035b6040516103318161032360e43560c43560a4356084356064356044356024353360043560208b01610e19565b03601f198101835282610e9d565b602081519101209260043588526004602052604088205492610354841515610ebe565b3691610eff565b611463565b610f55565b600435825260036020526040822033835260205260408220600160ff19825416179055602435610656575b604435151580610642575b6105a6575b606435610509575b6084351515806104f5575b610449575b60408051602435815260443560208201526064359181019190915260843560608201526001600160a01b0360a4358116608083015260e4351660a08201523390600435907f2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf89060c090a460017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b815260e4356001600160a01b0316600482015260843560248201529082908290604490829084905af180156104ea576104d2575b50506104ca608435600854610f91565b6008556103b8565b6104db90610e76565b6104e657815f6104ba565b5080fd5b6040513d84823e3d90fd5b5060e4356001600160a01b031615156103b3565b817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b815260c4356001600160a01b0316600482015260643560248201529082908290604490829084905af180156104ea57610592575b505061058a606435600754610f91565b6007556103a8565b61059b90610e76565b6104e657815f61057a565b817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b815260a4356001600160a01b03166004820152604480356024830152909183918391829084905af180156104ea5761062e575b5050610626604435600654610f91565b6006556103a0565b61063790610e76565b6104e657815f610616565b5060a4356001600160a01b0316151561039b565b907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b815233600482015260248035908201529082908290604490829084905af180156104ea576106d4575b50906106cc602435600554610f91565b600555610390565b6106dd90610e76565b5f6106bc565b60405162461bcd60e51b815260206004820152601f60248201527f416c726561647920636c61696d656420666f72207468697320706572696f64006044820152606490fd5b8280fd5b5034610102578060031936011261010257546040516001600160a01b039091168152602090f35b50346101025760403660031901126101025760ff6040602092610774610d53565b6004358252600385528282206001600160a01b03909116825284522054604051911615158152f35b5034610102576020366003190112610102576004356107b96113f6565b600254807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c60406001548151908152856020820152a28160015582526004602052604082205580f35b50346101025780600319360112610102576020600654604051908152f35b50346101025760203660031901126101025760406020916004358152600483522054604051908152f35b50346101025780600319360112610102576108636113f6565b80546001600160a01b03198116825581906001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b503461010257610120366003190112610102576004356001600160401b0381116104e6576108d6903690600401610d69565b610120526101e0526024356001600160401b0381116104e6576108fd903690600401610d69565b6102205260c0526044356001600160401b0381116104e657610923903690600401610d69565b6102405260a0526064356001600160401b0381116104e657610949903690600401610d69565b61020052610260526084356001600160401b0381116104e657610970903690600401610d69565b6102c0526101405260a4356001600160401b0381116104e657610997903690600401610d69565b61018052610280526001600160401b0360c43511610102576109be3660c435600401610d69565b610100526101c05260e4356001600160401b0381116104e6576109e5903690600401610d69565b6101a052608052610104356001600160401b0381116104e657610a0c903690600401610d69565b60e0526102a052610a1b611421565b61022051610120519081149081610c4c575b81610c3f575b81610c32575b81610c25575b81610c18575b81610c0b575b81610bff575b5015610bc25780610160525b610120516101605181811015610b9a57610a7a916101e051610fb2565b35610a8e610160516102205160c051610fb2565b35610aa2610160516102405160a051610fb2565b35610ab7610160516102005161026051610fb2565b35610acc610160516102c05161014051610fb2565b35610ae9610ae4610160516101805161028051610fb2565b610fd6565b610b00610ae461016051610100516101c051610fb2565b91610b17610ae4610160516101a051608051610fb2565b9360e051956101605196871015610b86576102a051600597881b81013590601e1981360301821215610b825701978835986001600160401b038a11610b825760208a9101981b36038813610b7e57610b6e99610fea565b6001610160510161016052610a5d565b8a80fd5b8b80fd5b634e487b7160e01b8a52603260045260248afd5b8260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b60405162461bcd60e51b8152602060048201526015602482015274082e4e4c2f240d8cadccee8d040dad2e6dac2e8c6d605b1b6044820152606490fd5b905060e051145f610a51565b6101a05181149150610a4b565b6101005181149150610a45565b6101805181149150610a3f565b6102c05181149150610a39565b6102005181149150610a33565b6102405181149150610a2d565b50346101025780600319360112610102576020600154604051908152f35b503461010257602036600319011261010257600435610c946113f6565b7f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c604060025492610cc4846113e8565b80600255816001558552600460205280828620557f5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc982600254958151908152866020820152a18151908582526020820152a280f35b50346101025780600319360112610102576020600854604051908152f35b9050346104e657816003193601126104e6576020906002548152f35b602435906001600160a01b038216820361019457565b9181601f84011215610194578235916001600160401b038311610194576020808501948460051b01011161019457565b15610da057565b60405162461bcd60e51b815260206004820152600e60248201526d125b9d985b1a59081c195c9a5bd960921b6044820152606490fd5b15610ddd57565b60405162461bcd60e51b8152602060048201526014602482015273496e76616c69642075736572206164647265737360601b6044820152606490fd5b97939691929460f099969189526bffffffffffffffffffffffff1997889687809660601b1660208c015260348b015260548a01526074890152609488015260601b1660b486015260601b1660c884015260601b1660dc8201520190565b6001600160401b038111610e8957604052565b634e487b7160e01b5f52604160045260245ffd5b90601f801991011681019081106001600160401b03821117610e8957604052565b15610ec557565b60405162461bcd60e51b815260206004820152601260248201527114195c9a5bd9081b9bdd081cd95d1d1b195960721b6044820152606490fd5b9092916001600160401b038411610e89578360051b6040519260208094610f2882850182610e9d565b809781520191810192831161019457905b828210610f465750505050565b81358152908301908301610f39565b15610f5c57565b60405162461bcd60e51b815260206004820152600d60248201526c24b73b30b634b210383937b7b360991b6044820152606490fd5b91908201809211610f9e57565b634e487b7160e01b5f52601160045260245ffd5b9190811015610fc25760051b0190565b634e487b7160e01b5f52603260045260245ffd5b356001600160a01b03811681036101945790565b969391989097959492955f9188835260036020526040832033845260205260ff6040842054166113db57611082916103609161035b8b8e6110658f6103238f8f8f918f908f9061103e6002548b1115610d99565b6110526001600160a01b0385161515610dd6565b604051988997602089019b33908d610e19565b519020928c87526004602052604087205492610354841515610ebe565b86815260036020526040812033825260205260408120600160ff1982541617905587611353575b88151580611341575b6112a7575b8561120d575b811515806111fb575b611135575b506040805197885260208801989098529686019390935260608501959095526001600160a01b03918216608085015293811660a08401529092169133917f2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf8908060c081015b0390a4565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b81526001600160a01b0386166004820152602481018490529082908290604490829084905af180156104ea576111e7575b50509161113091836111d97f2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf8979695600854610f91565b6008559193949550916110cb565b6111f18291610e76565b61010257806111a2565b506001600160a01b03841615156110c6565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b81526001600160a01b0387166004820152602481018890529082908290604490829084905af180156104ea57908291611293575b505061128b86600754610f91565b6007556110bd565b61129c90610e76565b61010257805f61127d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b81526001600160a01b0385166004820152602481018b90529082908290604490829084905af180156104ea5790829161132d575b505061132589600654610f91565b6006556110b7565b61133690610e76565b61010257805f611317565b506001600160a01b03831615156110b2565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b156104e6576040516340c10f1960e01b8152336004820152602481018a90529082908290604490829084905af180156104ea576113cc575b506113c488600554610f91565b6005556110a9565b6113d590610e76565b5f6113b7565b5050505050505050505050565b5f198114610f9e5760010190565b5f546001600160a01b0316330361140957565b60405163118cdaa760e01b8152336004820152602490fd5b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146114515760029055565b604051633ee5aeb560e01b8152600490fd5b9091905f915b81518310156114b1576020808460051b84010151915f8382105f146114a057505f525261149a60405f205b926113e8565b91611469565b9060409261149a9483525220611494565b915050149056fea264697066735822122079a314f2946b2067e3cfc8576446edcd53d8ba6d43608d8fd63d04a90558736064736f6c63430008140033",
}

// MiningRewardDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use MiningRewardDistributorMetaData.ABI instead.
var MiningRewardDistributorABI = MiningRewardDistributorMetaData.ABI

// MiningRewardDistributorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MiningRewardDistributorMetaData.Bin instead.
var MiningRewardDistributorBin = MiningRewardDistributorMetaData.Bin

// DeployMiningRewardDistributor deploys a new Ethereum contract, binding an instance of MiningRewardDistributor to it.
func DeployMiningRewardDistributor(auth *bind.TransactOpts, backend bind.ContractBackend, kawaiToken_ common.Address) (common.Address, *types.Transaction, *MiningRewardDistributor, error) {
	parsed, err := MiningRewardDistributorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MiningRewardDistributorBin), backend, kawaiToken_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MiningRewardDistributor{MiningRewardDistributorCaller: MiningRewardDistributorCaller{contract: contract}, MiningRewardDistributorTransactor: MiningRewardDistributorTransactor{contract: contract}, MiningRewardDistributorFilterer: MiningRewardDistributorFilterer{contract: contract}}, nil
}

// MiningRewardDistributor is an auto generated Go binding around an Ethereum contract.
type MiningRewardDistributor struct {
	MiningRewardDistributorCaller     // Read-only binding to the contract
	MiningRewardDistributorTransactor // Write-only binding to the contract
	MiningRewardDistributorFilterer   // Log filterer for contract events
}

// MiningRewardDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MiningRewardDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MiningRewardDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MiningRewardDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MiningRewardDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MiningRewardDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MiningRewardDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MiningRewardDistributorSession struct {
	Contract     *MiningRewardDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// MiningRewardDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MiningRewardDistributorCallerSession struct {
	Contract *MiningRewardDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// MiningRewardDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MiningRewardDistributorTransactorSession struct {
	Contract     *MiningRewardDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// MiningRewardDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MiningRewardDistributorRaw struct {
	Contract *MiningRewardDistributor // Generic contract binding to access the raw methods on
}

// MiningRewardDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MiningRewardDistributorCallerRaw struct {
	Contract *MiningRewardDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// MiningRewardDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MiningRewardDistributorTransactorRaw struct {
	Contract *MiningRewardDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMiningRewardDistributor creates a new instance of MiningRewardDistributor, bound to a specific deployed contract.
func NewMiningRewardDistributor(address common.Address, backend bind.ContractBackend) (*MiningRewardDistributor, error) {
	contract, err := bindMiningRewardDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributor{MiningRewardDistributorCaller: MiningRewardDistributorCaller{contract: contract}, MiningRewardDistributorTransactor: MiningRewardDistributorTransactor{contract: contract}, MiningRewardDistributorFilterer: MiningRewardDistributorFilterer{contract: contract}}, nil
}

// NewMiningRewardDistributorCaller creates a new read-only instance of MiningRewardDistributor, bound to a specific deployed contract.
func NewMiningRewardDistributorCaller(address common.Address, caller bind.ContractCaller) (*MiningRewardDistributorCaller, error) {
	contract, err := bindMiningRewardDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorCaller{contract: contract}, nil
}

// NewMiningRewardDistributorTransactor creates a new write-only instance of MiningRewardDistributor, bound to a specific deployed contract.
func NewMiningRewardDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*MiningRewardDistributorTransactor, error) {
	contract, err := bindMiningRewardDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorTransactor{contract: contract}, nil
}

// NewMiningRewardDistributorFilterer creates a new log filterer instance of MiningRewardDistributor, bound to a specific deployed contract.
func NewMiningRewardDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*MiningRewardDistributorFilterer, error) {
	contract, err := bindMiningRewardDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorFilterer{contract: contract}, nil
}

// bindMiningRewardDistributor binds a generic wrapper to an already deployed contract.
func bindMiningRewardDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MiningRewardDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MiningRewardDistributor *MiningRewardDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MiningRewardDistributor.Contract.MiningRewardDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MiningRewardDistributor *MiningRewardDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.MiningRewardDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MiningRewardDistributor *MiningRewardDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.MiningRewardDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MiningRewardDistributor *MiningRewardDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MiningRewardDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MiningRewardDistributor *MiningRewardDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MiningRewardDistributor *MiningRewardDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.contract.Transact(opts, method, params...)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) CurrentPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "currentPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorSession) CurrentPeriod() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.CurrentPeriod(&_MiningRewardDistributor.CallOpts)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) CurrentPeriod() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.CurrentPeriod(&_MiningRewardDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 contributorRewards, uint256 developerRewards, uint256 userRewards, uint256 affiliatorRewards)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) GetStats(opts *bind.CallOpts) (struct {
	Period             *big.Int
	ContributorRewards *big.Int
	DeveloperRewards   *big.Int
	UserRewards        *big.Int
	AffiliatorRewards  *big.Int
}, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "getStats")

	outstruct := new(struct {
		Period             *big.Int
		ContributorRewards *big.Int
		DeveloperRewards   *big.Int
		UserRewards        *big.Int
		AffiliatorRewards  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Period = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ContributorRewards = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.DeveloperRewards = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UserRewards = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AffiliatorRewards = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 contributorRewards, uint256 developerRewards, uint256 userRewards, uint256 affiliatorRewards)
func (_MiningRewardDistributor *MiningRewardDistributorSession) GetStats() (struct {
	Period             *big.Int
	ContributorRewards *big.Int
	DeveloperRewards   *big.Int
	UserRewards        *big.Int
	AffiliatorRewards  *big.Int
}, error) {
	return _MiningRewardDistributor.Contract.GetStats(&_MiningRewardDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 contributorRewards, uint256 developerRewards, uint256 userRewards, uint256 affiliatorRewards)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) GetStats() (struct {
	Period             *big.Int
	ContributorRewards *big.Int
	DeveloperRewards   *big.Int
	UserRewards        *big.Int
	AffiliatorRewards  *big.Int
}, error) {
	return _MiningRewardDistributor.Contract.GetStats(&_MiningRewardDistributor.CallOpts)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) HasClaimed(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "hasClaimed", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _MiningRewardDistributor.Contract.HasClaimed(&_MiningRewardDistributor.CallOpts, arg0, arg1)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _MiningRewardDistributor.Contract.HasClaimed(&_MiningRewardDistributor.CallOpts, arg0, arg1)
}

// HasClaimedPeriod is a free data retrieval call binding the contract method 0x3f08ccd0.
//
// Solidity: function hasClaimedPeriod(uint256 period, address contributor) view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) HasClaimedPeriod(opts *bind.CallOpts, period *big.Int, contributor common.Address) (bool, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "hasClaimedPeriod", period, contributor)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimedPeriod is a free data retrieval call binding the contract method 0x3f08ccd0.
//
// Solidity: function hasClaimedPeriod(uint256 period, address contributor) view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorSession) HasClaimedPeriod(period *big.Int, contributor common.Address) (bool, error) {
	return _MiningRewardDistributor.Contract.HasClaimedPeriod(&_MiningRewardDistributor.CallOpts, period, contributor)
}

// HasClaimedPeriod is a free data retrieval call binding the contract method 0x3f08ccd0.
//
// Solidity: function hasClaimedPeriod(uint256 period, address contributor) view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) HasClaimedPeriod(period *big.Int, contributor common.Address) (bool, error) {
	return _MiningRewardDistributor.Contract.HasClaimedPeriod(&_MiningRewardDistributor.CallOpts, period, contributor)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) KawaiToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "kawaiToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_MiningRewardDistributor *MiningRewardDistributorSession) KawaiToken() (common.Address, error) {
	return _MiningRewardDistributor.Contract.KawaiToken(&_MiningRewardDistributor.CallOpts)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) KawaiToken() (common.Address, error) {
	return _MiningRewardDistributor.Contract.KawaiToken(&_MiningRewardDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_MiningRewardDistributor *MiningRewardDistributorSession) MerkleRoot() ([32]byte, error) {
	return _MiningRewardDistributor.Contract.MerkleRoot(&_MiningRewardDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) MerkleRoot() ([32]byte, error) {
	return _MiningRewardDistributor.Contract.MerkleRoot(&_MiningRewardDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MiningRewardDistributor *MiningRewardDistributorSession) Owner() (common.Address, error) {
	return _MiningRewardDistributor.Contract.Owner(&_MiningRewardDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) Owner() (common.Address, error) {
	return _MiningRewardDistributor.Contract.Owner(&_MiningRewardDistributor.CallOpts)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) PeriodMerkleRoots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "periodMerkleRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_MiningRewardDistributor *MiningRewardDistributorSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _MiningRewardDistributor.Contract.PeriodMerkleRoots(&_MiningRewardDistributor.CallOpts, arg0)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _MiningRewardDistributor.Contract.PeriodMerkleRoots(&_MiningRewardDistributor.CallOpts, arg0)
}

// TotalAffiliatorRewards is a free data retrieval call binding the contract method 0x06b77717.
//
// Solidity: function totalAffiliatorRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) TotalAffiliatorRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "totalAffiliatorRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAffiliatorRewards is a free data retrieval call binding the contract method 0x06b77717.
//
// Solidity: function totalAffiliatorRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorSession) TotalAffiliatorRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalAffiliatorRewards(&_MiningRewardDistributor.CallOpts)
}

// TotalAffiliatorRewards is a free data retrieval call binding the contract method 0x06b77717.
//
// Solidity: function totalAffiliatorRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) TotalAffiliatorRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalAffiliatorRewards(&_MiningRewardDistributor.CallOpts)
}

// TotalContributorRewards is a free data retrieval call binding the contract method 0x96e379f0.
//
// Solidity: function totalContributorRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) TotalContributorRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "totalContributorRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalContributorRewards is a free data retrieval call binding the contract method 0x96e379f0.
//
// Solidity: function totalContributorRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorSession) TotalContributorRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalContributorRewards(&_MiningRewardDistributor.CallOpts)
}

// TotalContributorRewards is a free data retrieval call binding the contract method 0x96e379f0.
//
// Solidity: function totalContributorRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) TotalContributorRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalContributorRewards(&_MiningRewardDistributor.CallOpts)
}

// TotalDeveloperRewards is a free data retrieval call binding the contract method 0x77363251.
//
// Solidity: function totalDeveloperRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) TotalDeveloperRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "totalDeveloperRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDeveloperRewards is a free data retrieval call binding the contract method 0x77363251.
//
// Solidity: function totalDeveloperRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorSession) TotalDeveloperRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalDeveloperRewards(&_MiningRewardDistributor.CallOpts)
}

// TotalDeveloperRewards is a free data retrieval call binding the contract method 0x77363251.
//
// Solidity: function totalDeveloperRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) TotalDeveloperRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalDeveloperRewards(&_MiningRewardDistributor.CallOpts)
}

// TotalUserRewards is a free data retrieval call binding the contract method 0xfd8bfafc.
//
// Solidity: function totalUserRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) TotalUserRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "totalUserRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUserRewards is a free data retrieval call binding the contract method 0xfd8bfafc.
//
// Solidity: function totalUserRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorSession) TotalUserRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalUserRewards(&_MiningRewardDistributor.CallOpts)
}

// TotalUserRewards is a free data retrieval call binding the contract method 0xfd8bfafc.
//
// Solidity: function totalUserRewards() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) TotalUserRewards() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TotalUserRewards(&_MiningRewardDistributor.CallOpts)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) AdvancePeriod(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "advancePeriod", _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.AdvancePeriod(&_MiningRewardDistributor.TransactOpts, _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.AdvancePeriod(&_MiningRewardDistributor.TransactOpts, _merkleRoot)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0x437f3d22.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] contributorAmounts, uint256[] developerAmounts, uint256[] userAmounts, uint256[] affiliatorAmounts, address[] developers, address[] users, address[] affiliators, bytes32[][] merkleProofs) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) ClaimMultiplePeriods(opts *bind.TransactOpts, periods []*big.Int, contributorAmounts []*big.Int, developerAmounts []*big.Int, userAmounts []*big.Int, affiliatorAmounts []*big.Int, developers []common.Address, users []common.Address, affiliators []common.Address, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "claimMultiplePeriods", periods, contributorAmounts, developerAmounts, userAmounts, affiliatorAmounts, developers, users, affiliators, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0x437f3d22.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] contributorAmounts, uint256[] developerAmounts, uint256[] userAmounts, uint256[] affiliatorAmounts, address[] developers, address[] users, address[] affiliators, bytes32[][] merkleProofs) returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) ClaimMultiplePeriods(periods []*big.Int, contributorAmounts []*big.Int, developerAmounts []*big.Int, userAmounts []*big.Int, affiliatorAmounts []*big.Int, developers []common.Address, users []common.Address, affiliators []common.Address, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.ClaimMultiplePeriods(&_MiningRewardDistributor.TransactOpts, periods, contributorAmounts, developerAmounts, userAmounts, affiliatorAmounts, developers, users, affiliators, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0x437f3d22.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] contributorAmounts, uint256[] developerAmounts, uint256[] userAmounts, uint256[] affiliatorAmounts, address[] developers, address[] users, address[] affiliators, bytes32[][] merkleProofs) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) ClaimMultiplePeriods(periods []*big.Int, contributorAmounts []*big.Int, developerAmounts []*big.Int, userAmounts []*big.Int, affiliatorAmounts []*big.Int, developers []common.Address, users []common.Address, affiliators []common.Address, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.ClaimMultiplePeriods(&_MiningRewardDistributor.TransactOpts, periods, contributorAmounts, developerAmounts, userAmounts, affiliatorAmounts, developers, users, affiliators, merkleProofs)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x95112df3.
//
// Solidity: function claimReward(uint256 period, uint256 contributorAmount, uint256 developerAmount, uint256 userAmount, uint256 affiliatorAmount, address developer, address user, address affiliator, bytes32[] merkleProof) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) ClaimReward(opts *bind.TransactOpts, period *big.Int, contributorAmount *big.Int, developerAmount *big.Int, userAmount *big.Int, affiliatorAmount *big.Int, developer common.Address, user common.Address, affiliator common.Address, merkleProof [][32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "claimReward", period, contributorAmount, developerAmount, userAmount, affiliatorAmount, developer, user, affiliator, merkleProof)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x95112df3.
//
// Solidity: function claimReward(uint256 period, uint256 contributorAmount, uint256 developerAmount, uint256 userAmount, uint256 affiliatorAmount, address developer, address user, address affiliator, bytes32[] merkleProof) returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) ClaimReward(period *big.Int, contributorAmount *big.Int, developerAmount *big.Int, userAmount *big.Int, affiliatorAmount *big.Int, developer common.Address, user common.Address, affiliator common.Address, merkleProof [][32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.ClaimReward(&_MiningRewardDistributor.TransactOpts, period, contributorAmount, developerAmount, userAmount, affiliatorAmount, developer, user, affiliator, merkleProof)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x95112df3.
//
// Solidity: function claimReward(uint256 period, uint256 contributorAmount, uint256 developerAmount, uint256 userAmount, uint256 affiliatorAmount, address developer, address user, address affiliator, bytes32[] merkleProof) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) ClaimReward(period *big.Int, contributorAmount *big.Int, developerAmount *big.Int, userAmount *big.Int, affiliatorAmount *big.Int, developer common.Address, user common.Address, affiliator common.Address, merkleProof [][32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.ClaimReward(&_MiningRewardDistributor.TransactOpts, period, contributorAmount, developerAmount, userAmount, affiliatorAmount, developer, user, affiliator, merkleProof)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.RenounceOwnership(&_MiningRewardDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.RenounceOwnership(&_MiningRewardDistributor.TransactOpts)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) SetMerkleRoot(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "setMerkleRoot", _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.SetMerkleRoot(&_MiningRewardDistributor.TransactOpts, _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.SetMerkleRoot(&_MiningRewardDistributor.TransactOpts, _merkleRoot)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.TransferOwnership(&_MiningRewardDistributor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.TransferOwnership(&_MiningRewardDistributor.TransactOpts, newOwner)
}

// MiningRewardDistributorMerkleRootUpdatedIterator is returned from FilterMerkleRootUpdated and is used to iterate over the raw logs and unpacked data for MerkleRootUpdated events raised by the MiningRewardDistributor contract.
type MiningRewardDistributorMerkleRootUpdatedIterator struct {
	Event *MiningRewardDistributorMerkleRootUpdated // Event containing the contract specifics and raw log

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
func (it *MiningRewardDistributorMerkleRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MiningRewardDistributorMerkleRootUpdated)
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
		it.Event = new(MiningRewardDistributorMerkleRootUpdated)
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
func (it *MiningRewardDistributorMerkleRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MiningRewardDistributorMerkleRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MiningRewardDistributorMerkleRootUpdated represents a MerkleRootUpdated event raised by the MiningRewardDistributor contract.
type MiningRewardDistributorMerkleRootUpdated struct {
	Period  *big.Int
	OldRoot [32]byte
	NewRoot [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootUpdated is a free log retrieval operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) FilterMerkleRootUpdated(opts *bind.FilterOpts, period []*big.Int) (*MiningRewardDistributorMerkleRootUpdatedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _MiningRewardDistributor.contract.FilterLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorMerkleRootUpdatedIterator{contract: _MiningRewardDistributor.contract, event: "MerkleRootUpdated", logs: logs, sub: sub}, nil
}

// WatchMerkleRootUpdated is a free log subscription operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) WatchMerkleRootUpdated(opts *bind.WatchOpts, sink chan<- *MiningRewardDistributorMerkleRootUpdated, period []*big.Int) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _MiningRewardDistributor.contract.WatchLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MiningRewardDistributorMerkleRootUpdated)
				if err := _MiningRewardDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
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
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) ParseMerkleRootUpdated(log types.Log) (*MiningRewardDistributorMerkleRootUpdated, error) {
	event := new(MiningRewardDistributorMerkleRootUpdated)
	if err := _MiningRewardDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MiningRewardDistributorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MiningRewardDistributor contract.
type MiningRewardDistributorOwnershipTransferredIterator struct {
	Event *MiningRewardDistributorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MiningRewardDistributorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MiningRewardDistributorOwnershipTransferred)
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
		it.Event = new(MiningRewardDistributorOwnershipTransferred)
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
func (it *MiningRewardDistributorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MiningRewardDistributorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MiningRewardDistributorOwnershipTransferred represents a OwnershipTransferred event raised by the MiningRewardDistributor contract.
type MiningRewardDistributorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MiningRewardDistributorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MiningRewardDistributor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorOwnershipTransferredIterator{contract: _MiningRewardDistributor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MiningRewardDistributorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MiningRewardDistributor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MiningRewardDistributorOwnershipTransferred)
				if err := _MiningRewardDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) ParseOwnershipTransferred(log types.Log) (*MiningRewardDistributorOwnershipTransferred, error) {
	event := new(MiningRewardDistributorOwnershipTransferred)
	if err := _MiningRewardDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MiningRewardDistributorPeriodAdvancedIterator is returned from FilterPeriodAdvanced and is used to iterate over the raw logs and unpacked data for PeriodAdvanced events raised by the MiningRewardDistributor contract.
type MiningRewardDistributorPeriodAdvancedIterator struct {
	Event *MiningRewardDistributorPeriodAdvanced // Event containing the contract specifics and raw log

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
func (it *MiningRewardDistributorPeriodAdvancedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MiningRewardDistributorPeriodAdvanced)
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
		it.Event = new(MiningRewardDistributorPeriodAdvanced)
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
func (it *MiningRewardDistributorPeriodAdvancedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MiningRewardDistributorPeriodAdvancedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MiningRewardDistributorPeriodAdvanced represents a PeriodAdvanced event raised by the MiningRewardDistributor contract.
type MiningRewardDistributorPeriodAdvanced struct {
	OldPeriod *big.Int
	NewPeriod *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPeriodAdvanced is a free log retrieval operation binding the contract event 0x5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc9.
//
// Solidity: event PeriodAdvanced(uint256 oldPeriod, uint256 newPeriod)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) FilterPeriodAdvanced(opts *bind.FilterOpts) (*MiningRewardDistributorPeriodAdvancedIterator, error) {

	logs, sub, err := _MiningRewardDistributor.contract.FilterLogs(opts, "PeriodAdvanced")
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorPeriodAdvancedIterator{contract: _MiningRewardDistributor.contract, event: "PeriodAdvanced", logs: logs, sub: sub}, nil
}

// WatchPeriodAdvanced is a free log subscription operation binding the contract event 0x5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc9.
//
// Solidity: event PeriodAdvanced(uint256 oldPeriod, uint256 newPeriod)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) WatchPeriodAdvanced(opts *bind.WatchOpts, sink chan<- *MiningRewardDistributorPeriodAdvanced) (event.Subscription, error) {

	logs, sub, err := _MiningRewardDistributor.contract.WatchLogs(opts, "PeriodAdvanced")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MiningRewardDistributorPeriodAdvanced)
				if err := _MiningRewardDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
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
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) ParsePeriodAdvanced(log types.Log) (*MiningRewardDistributorPeriodAdvanced, error) {
	event := new(MiningRewardDistributorPeriodAdvanced)
	if err := _MiningRewardDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MiningRewardDistributorRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the MiningRewardDistributor contract.
type MiningRewardDistributorRewardClaimedIterator struct {
	Event *MiningRewardDistributorRewardClaimed // Event containing the contract specifics and raw log

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
func (it *MiningRewardDistributorRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MiningRewardDistributorRewardClaimed)
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
		it.Event = new(MiningRewardDistributorRewardClaimed)
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
func (it *MiningRewardDistributorRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MiningRewardDistributorRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MiningRewardDistributorRewardClaimed represents a RewardClaimed event raised by the MiningRewardDistributor contract.
type MiningRewardDistributorRewardClaimed struct {
	Period            *big.Int
	Contributor       common.Address
	User              common.Address
	ContributorAmount *big.Int
	DeveloperAmount   *big.Int
	UserAmount        *big.Int
	AffiliatorAmount  *big.Int
	Developer         common.Address
	Affiliator        common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf8.
//
// Solidity: event RewardClaimed(uint256 indexed period, address indexed contributor, address indexed user, uint256 contributorAmount, uint256 developerAmount, uint256 userAmount, uint256 affiliatorAmount, address developer, address affiliator)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) FilterRewardClaimed(opts *bind.FilterOpts, period []*big.Int, contributor []common.Address, user []common.Address) (*MiningRewardDistributorRewardClaimedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _MiningRewardDistributor.contract.FilterLogs(opts, "RewardClaimed", periodRule, contributorRule, userRule)
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorRewardClaimedIterator{contract: _MiningRewardDistributor.contract, event: "RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf8.
//
// Solidity: event RewardClaimed(uint256 indexed period, address indexed contributor, address indexed user, uint256 contributorAmount, uint256 developerAmount, uint256 userAmount, uint256 affiliatorAmount, address developer, address affiliator)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *MiningRewardDistributorRewardClaimed, period []*big.Int, contributor []common.Address, user []common.Address) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _MiningRewardDistributor.contract.WatchLogs(opts, "RewardClaimed", periodRule, contributorRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MiningRewardDistributorRewardClaimed)
				if err := _MiningRewardDistributor.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
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

// ParseRewardClaimed is a log parse operation binding the contract event 0x2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf8.
//
// Solidity: event RewardClaimed(uint256 indexed period, address indexed contributor, address indexed user, uint256 contributorAmount, uint256 developerAmount, uint256 userAmount, uint256 affiliatorAmount, address developer, address affiliator)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) ParseRewardClaimed(log types.Log) (*MiningRewardDistributorRewardClaimed, error) {
	event := new(MiningRewardDistributorRewardClaimed)
	if err := _MiningRewardDistributor.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
