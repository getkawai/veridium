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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"kawaiToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"TOTAL_ALLOCATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"advancePeriod\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimMultiplePeriods\",\"inputs\":[{\"name\":\"periods\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"contributorAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"developerAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"userAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"affiliatorAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"developers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"affiliators\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"merkleProofs\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimReward\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributorAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"developerAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"affiliatorAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"developer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"affiliator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStats\",\"inputs\":[],\"outputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributorRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"developerRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"affiliatorRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingAllocation\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedPeriod\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kawaiToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"periodMerkleRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRootForPeriod\",\"inputs\":[{\"name\":\"_period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalAffiliatorRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalContributorRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalDeveloperRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalUserRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodAdvanced\",\"inputs\":[{\"name\":\"oldPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"contributorAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"developerAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"userAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"affiliatorAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"developer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"affiliator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a03461017057601f6119a338819003918201601f19168301916001600160401b038311848410176101745780849260209460405283398101031261017057516001600160a01b0380821691829003610170573315610158575f543360018060a01b03198216175f55604051913391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005581156101165750608052600160025560405161181a908161018982396080518181816101e5015281816105f3015281816106b20152818161074f015281816107fe015281816114430152818161151b015281816115b501526116610152f35b62461bcd60e51b815260206004820152601560248201527f496e76616c6964204b41574149206164647265737300000000000000000000006044820152606490fd5b604051631e4fbdf760e01b81525f6004820152602490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6102e0806040526004361015610013575f80fd5b5f90813560e01c90816306040618146110435750806306b77717146110255780630ae6540314610f835780630dac7bea14610f5c5780632eb4a7ab14610f3e5780633f08ccd01461093c5780633f4ba83a14610ecd578063437f3d2214610b135780635c975abb14610aee578063715018a614610a94578063727a7c5a14610a6a5780637736325114610a4c5780637cb64759146109e65780638456cb5914610985578063873f6f9e1461093c5780638da5cb5b1461091557806395112df31461037257806396e379f014610354578063b24aa278146102a5578063c59d484714610214578063cb56cd4f146101cf578063f2fde38b1461013c5763fd8bfafc1461011c575f80fd5b346101395780600319360112610139576020600754604051908152f35b80fd5b5034610139576020366003190112610139576004356001600160a01b03818116918290036101cb5761016c611702565b81156101b2575f54826bffffffffffffffffffffffff60a01b8216175f55167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a380f35b604051631e4fbdf760e01b815260048101849052602490fd5b5f80fd5b50346101395780600319360112610139576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b50346101395780600319360112610139576005546006546102358183611125565b6102426007548092611125565b906102506008548093611125565b92600254946b019d971e4fe8401e740000009485039485116102915760c09650604051958652602086015260408501526060840152608083015260a0820152f35b634e487b7160e01b87526011600452602487fd5b5034610139576040366003190112610139576004356024356102c5611702565b6102d08215156110a5565b80156103195760407f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c91838552600460205281852090808254925582519182526020820152a280f35b60405162461bcd60e51b8152602060048201526013602482015272125b9d985b1a59081b595c9adb19481c9bdbdd606a1b6044820152606490fd5b50346101395780600319360112610139576020600554604051908152f35b5034610139576101203660031901126101395760a4356001600160a01b03811690036101395760c4356001600160a01b03811690036101395760e4356001600160a01b038116900361013957610104356001600160401b038111610887576103de903690600401611075565b6103e661174d565b6103ee61172d565b6103fe60025460043511156110a5565b600435835260036020526040832033845260205260ff6040842054166108d05761043460c4356001600160a01b031615156110e2565b6b019d971e4fe8401e7400000061048a61046360843561045e60643561045e604435602435611125565b611125565b61045e61048161047860055460065490611125565b60075490611125565b60085490611125565b1161088b5761050691610501916104fc60405160208101906104d6816104c860e43560c43560a435608435606435604435602435336004358d611146565b03601f1981018352826111ca565b51902092600435875260046020526040872054926104f58415156111eb565b369161122c565b61178f565b611282565b600435815260036020526040812033825260205260408120600160ff198254161790556024356107fc575b6044351515806107e8575b61074c575b6064356106af575b60843515158061069b575b6105f0575b60408051602435815260443560208201526064359181019190915260843560608201526001600160a01b0360a4358116608083015260e435811660a083015260c43516903390600435907f2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf89060c090a460017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610698576040516340c10f1960e01b815260e4356001600160a01b0316600482015260843560248201529082908290604490829084905af1801561068d57610679575b5050610671608435600854611125565b600855610559565b610682906111a3565b61013957805f610661565b6040513d84823e3d90fd5b50fd5b5060e4356001600160a01b03161515610554565b807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610698576040516340c10f1960e01b815260c4356001600160a01b0316600482015260643560248201529082908290604490829084905af1801561068d57610738575b5050610730606435600754611125565b600755610549565b610741906111a3565b61013957805f610720565b807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610698576040516340c10f1960e01b815260a4356001600160a01b03166004820152604480356024830152909183918391829084905af1801561068d576107d4575b50506107cc604435600654611125565b600655610541565b6107dd906111a3565b61013957805f6107bc565b5060a4356001600160a01b0316151561053c565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610887576040516340c10f1960e01b815233600482015260248035908201529082908290604490829084905af1801561068d57610878575b50610870602435600554611125565b600555610531565b610881906111a3565b5f610861565b5080fd5b60405162461bcd60e51b815260206004820152601860248201527f4578636565647320746f74616c20616c6c6f636174696f6e00000000000000006044820152606490fd5b60405162461bcd60e51b815260206004820152601f60248201527f416c726561647920636c61696d656420666f72207468697320706572696f64006044820152606490fd5b5034610139578060031936011261013957546040516001600160a01b039091168152602090f35b50346101395760403660031901126101395760ff604060209261095d61105f565b6004358252600385528282206001600160a01b03909116825284522054604051911615158152f35b503461013957806003193601126101395761099e611702565b6109a661172d565b805460ff60a01b1916600160a01b1781556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602090a180f35b503461013957602036600319011261013957600435610a03611702565b600254807f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c60406001548151908152856020820152a28160015582526004602052604082205580f35b50346101395780600319360112610139576020600654604051908152f35b50346101395760203660031901126101395760406020916004358152600483522054604051908152f35b5034610139578060031936011261013957610aad611702565b80546001600160a01b03198116825581906001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b503461013957806003193601126101395760ff6020915460a01c166040519015158152f35b503461013957610120366003190112610139576004356001600160401b03811161088757610b45903690600401611075565b61024052610200526024356001600160401b03811161088757610b6c903690600401611075565b6102c0526101e0526044356001600160401b03811161088757610b93903690600401611075565b61028052610100526064356001600160401b03811161088757610bba903690600401611075565b60a052610260526084356001600160401b03811161088757610be0903690600401611075565b60c0526101805260a4356001600160401b03811161088757610c06903690600401611075565b60e052610160526001600160401b0360c4351161013957610c2c3660c435600401611075565b6101a052610120526001600160401b0360e4351161013957610c533660e435600401611075565b610140526101c052610104356001600160401b03811161088757610c7b903690600401611075565b60805261022052610c8a61174d565b610c9261172d565b6102c051610240519081149081610ec0575b81610eb4575b81610ea8575b81610e9c575b81610e8f575b81610e82575b81610e76575b5015610e3957806102a0525b610240516102a05181811015610e1157610cf191610200516112be565b35610d066102a0516102c0516101e0516112be565b35610d1b6102a05161028051610100516112be565b35610d2f6102a05160a051610260516112be565b35610d436102a05160c051610180516112be565b35610d5f610d5a6102a05160e051610160516112be565b6112e2565b610d76610d5a6102a0516101a051610120516112be565b91610d8e610d5a6102a051610140516101c0516112be565b93608051956102a05196871015610dfd5761022051600597881b81013590601e1981360301821215610df95701978835986001600160401b038a11610df95760208a9101981b36038813610df557610de5996112f6565b60016102a051016102a052610cd4565b8a80fd5b8b80fd5b634e487b7160e01b8a52603260045260248afd5b8260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b60405162461bcd60e51b8152602060048201526015602482015274082e4e4c2f240d8cadccee8d040dad2e6dac2e8c6d605b1b6044820152606490fd5b9050608051145f610cc8565b6101405181149150610cc2565b6101a05181149150610cbc565b60e05181149150610cb6565b60c05181149150610cb0565b60a05181149150610caa565b6102805181149150610ca4565b5034610139578060031936011261013957610ee6611702565b805460ff8160a01c1615610f2c5760ff60a01b191681556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa90602090a180f35b604051638dfc202b60e01b8152600490fd5b50346101395780600319360112610139576020600154604051908152f35b503461013957806003193601126101395760206040516b019d971e4fe8401e740000008152f35b503461013957602036600319011261013957600435610fa0611702565b7f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c604060025492610fd0846116f4565b80600255816001558552600460205280828620557f5c12640e4659b07c515609d150d36890ae4b15c3d83514b006a6dfd16700cdc982600254958151908152866020820152a18151908582526020820152a280f35b50346101395780600319360112610139576020600854604051908152f35b9050346108875781600319360112610887576020906002548152f35b602435906001600160a01b03821682036101cb57565b9181601f840112156101cb578235916001600160401b0383116101cb576020808501948460051b0101116101cb57565b156110ac57565b60405162461bcd60e51b815260206004820152600e60248201526d125b9d985b1a59081c195c9a5bd960921b6044820152606490fd5b156110e957565b60405162461bcd60e51b8152602060048201526014602482015273496e76616c69642075736572206164647265737360601b6044820152606490fd5b9190820180921161113257565b634e487b7160e01b5f52601160045260245ffd5b97939691929460f099969189526bffffffffffffffffffffffff1997889687809660601b1660208c015260348b015260548a01526074890152609488015260601b1660b486015260601b1660c884015260601b1660dc8201520190565b6001600160401b0381116111b657604052565b634e487b7160e01b5f52604160045260245ffd5b90601f801991011681019081106001600160401b038211176111b657604052565b156111f257565b60405162461bcd60e51b815260206004820152601260248201527114195c9a5bd9081b9bdd081cd95d1d1b195960721b6044820152606490fd5b9092916001600160401b0384116111b6578360051b6040519260208094611255828501826111ca565b80978152019181019283116101cb57905b8282106112735750505050565b81358152908301908301611266565b1561128957565b60405162461bcd60e51b815260206004820152600d60248201526c24b73b30b634b210383937b7b360991b6044820152606490fd5b91908110156112ce5760051b0190565b634e487b7160e01b5f52603260045260245ffd5b356001600160a01b03811681036101cb5790565b969391989097959492955f9188835260036020526040832033845260205260ff6040842054166116e75761138e91610501916104fc8b8e6113718f6104c88f8f8f918f908f9061134a6002548b11156110a5565b61135e6001600160a01b03851615156110e2565b604051988997602089019b33908d611146565b519020928c875260046020526040872054926104f58415156111eb565b86815260036020526040812033825260205260408120600160ff198254161790558761165f575b8815158061164d575b6115b3575b85611519575b81151580611507575b611441575b506040805197885260208801989098529686019390935260608501959095526001600160a01b03918216608085015293811660a08401529092169133917f2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf8908060c081015b0390a4565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610887576040516340c10f1960e01b81526001600160a01b0386166004820152602481018490529082908290604490829084905af1801561068d576114f3575b50509161143c91836114e57f2d081fe3985c9f70b31e1b13fe5934e9007bb2283ea761c4e3ace7b222edcaf8979695600854611125565b6008559193949550916113d7565b6114fd82916111a3565b61013957806114ae565b506001600160a01b03841615156113d2565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610887576040516340c10f1960e01b81526001600160a01b0387166004820152602481018890529082908290604490829084905af1801561068d5790829161159f575b505061159786600754611125565b6007556113c9565b6115a8906111a3565b61013957805f611589565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610887576040516340c10f1960e01b81526001600160a01b0385166004820152602481018b90529082908290604490829084905af1801561068d57908291611639575b505061163189600654611125565b6006556113c3565b611642906111a3565b61013957805f611623565b506001600160a01b03831615156113be565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316803b15610887576040516340c10f1960e01b8152336004820152602481018a90529082908290604490829084905af1801561068d576116d8575b506116d088600554611125565b6005556113b5565b6116e1906111a3565b5f6116c3565b5050505050505050505050565b5f1981146111325760010190565b5f546001600160a01b0316330361171557565b60405163118cdaa760e01b8152336004820152602490fd5b60ff5f5460a01c1661173b57565b60405163d93c066560e01b8152600490fd5b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00600281541461177d5760029055565b604051633ee5aeb560e01b8152600490fd5b9091905f915b81518310156117dd576020808460051b84010151915f8382105f146117cc57505f52526117c660405f205b926116f4565b91611795565b906040926117c694835252206117c0565b915050149056fea26469706673582212209645ff3228730262d0dc0a5f640fda6678e8ea86bb2fc6e5fd0278b951813c3864736f6c63430008140033",
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

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) TOTALALLOCATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "TOTAL_ALLOCATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorSession) TOTALALLOCATION() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TOTALALLOCATION(&_MiningRewardDistributor.CallOpts)
}

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) TOTALALLOCATION() (*big.Int, error) {
	return _MiningRewardDistributor.Contract.TOTALALLOCATION(&_MiningRewardDistributor.CallOpts)
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
// Solidity: function getStats() view returns(uint256 period, uint256 contributorRewards, uint256 developerRewards, uint256 userRewards, uint256 affiliatorRewards, uint256 remainingAllocation)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) GetStats(opts *bind.CallOpts) (struct {
	Period              *big.Int
	ContributorRewards  *big.Int
	DeveloperRewards    *big.Int
	UserRewards         *big.Int
	AffiliatorRewards   *big.Int
	RemainingAllocation *big.Int
}, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "getStats")

	outstruct := new(struct {
		Period              *big.Int
		ContributorRewards  *big.Int
		DeveloperRewards    *big.Int
		UserRewards         *big.Int
		AffiliatorRewards   *big.Int
		RemainingAllocation *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Period = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ContributorRewards = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.DeveloperRewards = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UserRewards = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AffiliatorRewards = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.RemainingAllocation = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 contributorRewards, uint256 developerRewards, uint256 userRewards, uint256 affiliatorRewards, uint256 remainingAllocation)
func (_MiningRewardDistributor *MiningRewardDistributorSession) GetStats() (struct {
	Period              *big.Int
	ContributorRewards  *big.Int
	DeveloperRewards    *big.Int
	UserRewards         *big.Int
	AffiliatorRewards   *big.Int
	RemainingAllocation *big.Int
}, error) {
	return _MiningRewardDistributor.Contract.GetStats(&_MiningRewardDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 period, uint256 contributorRewards, uint256 developerRewards, uint256 userRewards, uint256 affiliatorRewards, uint256 remainingAllocation)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) GetStats() (struct {
	Period              *big.Int
	ContributorRewards  *big.Int
	DeveloperRewards    *big.Int
	UserRewards         *big.Int
	AffiliatorRewards   *big.Int
	RemainingAllocation *big.Int
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

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MiningRewardDistributor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorSession) Paused() (bool, error) {
	return _MiningRewardDistributor.Contract.Paused(&_MiningRewardDistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MiningRewardDistributor *MiningRewardDistributorCallerSession) Paused() (bool, error) {
	return _MiningRewardDistributor.Contract.Paused(&_MiningRewardDistributor.CallOpts)
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

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) Pause() (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.Pause(&_MiningRewardDistributor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) Pause() (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.Pause(&_MiningRewardDistributor.TransactOpts)
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

// SetMerkleRootForPeriod is a paid mutator transaction binding the contract method 0xb24aa278.
//
// Solidity: function setMerkleRootForPeriod(uint256 _period, bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) SetMerkleRootForPeriod(opts *bind.TransactOpts, _period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "setMerkleRootForPeriod", _period, _merkleRoot)
}

// SetMerkleRootForPeriod is a paid mutator transaction binding the contract method 0xb24aa278.
//
// Solidity: function setMerkleRootForPeriod(uint256 _period, bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) SetMerkleRootForPeriod(_period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.SetMerkleRootForPeriod(&_MiningRewardDistributor.TransactOpts, _period, _merkleRoot)
}

// SetMerkleRootForPeriod is a paid mutator transaction binding the contract method 0xb24aa278.
//
// Solidity: function setMerkleRootForPeriod(uint256 _period, bytes32 _merkleRoot) returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) SetMerkleRootForPeriod(_period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.SetMerkleRootForPeriod(&_MiningRewardDistributor.TransactOpts, _period, _merkleRoot)
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

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MiningRewardDistributor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MiningRewardDistributor *MiningRewardDistributorSession) Unpause() (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.Unpause(&_MiningRewardDistributor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MiningRewardDistributor *MiningRewardDistributorTransactorSession) Unpause() (*types.Transaction, error) {
	return _MiningRewardDistributor.Contract.Unpause(&_MiningRewardDistributor.TransactOpts)
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

// MiningRewardDistributorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the MiningRewardDistributor contract.
type MiningRewardDistributorPausedIterator struct {
	Event *MiningRewardDistributorPaused // Event containing the contract specifics and raw log

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
func (it *MiningRewardDistributorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MiningRewardDistributorPaused)
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
		it.Event = new(MiningRewardDistributorPaused)
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
func (it *MiningRewardDistributorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MiningRewardDistributorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MiningRewardDistributorPaused represents a Paused event raised by the MiningRewardDistributor contract.
type MiningRewardDistributorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) FilterPaused(opts *bind.FilterOpts) (*MiningRewardDistributorPausedIterator, error) {

	logs, sub, err := _MiningRewardDistributor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorPausedIterator{contract: _MiningRewardDistributor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MiningRewardDistributorPaused) (event.Subscription, error) {

	logs, sub, err := _MiningRewardDistributor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MiningRewardDistributorPaused)
				if err := _MiningRewardDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) ParsePaused(log types.Log) (*MiningRewardDistributorPaused, error) {
	event := new(MiningRewardDistributorPaused)
	if err := _MiningRewardDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
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

// MiningRewardDistributorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the MiningRewardDistributor contract.
type MiningRewardDistributorUnpausedIterator struct {
	Event *MiningRewardDistributorUnpaused // Event containing the contract specifics and raw log

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
func (it *MiningRewardDistributorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MiningRewardDistributorUnpaused)
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
		it.Event = new(MiningRewardDistributorUnpaused)
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
func (it *MiningRewardDistributorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MiningRewardDistributorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MiningRewardDistributorUnpaused represents a Unpaused event raised by the MiningRewardDistributor contract.
type MiningRewardDistributorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MiningRewardDistributorUnpausedIterator, error) {

	logs, sub, err := _MiningRewardDistributor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MiningRewardDistributorUnpausedIterator{contract: _MiningRewardDistributor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MiningRewardDistributorUnpaused) (event.Subscription, error) {

	logs, sub, err := _MiningRewardDistributor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MiningRewardDistributorUnpaused)
				if err := _MiningRewardDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_MiningRewardDistributor *MiningRewardDistributorFilterer) ParseUnpaused(log types.Log) (*MiningRewardDistributorUnpaused, error) {
	event := new(MiningRewardDistributorUnpaused)
	if err := _MiningRewardDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
