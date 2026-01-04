// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cashbackdistributor

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

// DepositCashbackDistributorMetaData contains all meta data concerning the DepositCashbackDistributor contract.
var DepositCashbackDistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"kawaiToken_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"TOTAL_ALLOCATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"advancePeriod\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimCashback\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimMultiplePeriods\",\"inputs\":[{\"name\":\"periods\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"kawaiAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"merkleProofs\",\"type\":\"bytes32[][]\",\"internalType\":\"bytes32[][]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeriodMerkleRoot\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStats\",\"inputs\":[],\"outputs\":[{\"name\":\"_currentPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalKawaiDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_remainingAllocation\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalUsers\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimed\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasClaimedAnyPeriod\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasUserClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kawaiToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"merkleRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"periodMerkleRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMerkleRoot\",\"inputs\":[{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeriodMerkleRoot\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalKawaiDistributed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalUsers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CashbackClaimed\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"kawaiAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PeriodAdvanced\",\"inputs\":[{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051612326380380612326833981810160405281019061003191906102d6565b335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100999190610310565b60405180910390fd5b6100b18161018560201b60201c565b5060016100d06100c561024660201b60201c565b61026f60201b60201c565b5f01819055505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610144576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161013b90610383565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505060018081905550506103a1565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6102a58261027c565b9050919050565b6102b58161029b565b81146102bf575f5ffd5b50565b5f815190506102d0816102ac565b92915050565b5f602082840312156102eb576102ea610278565b5b5f6102f8848285016102c2565b91505092915050565b61030a8161029b565b82525050565b5f6020820190506103235f830184610301565b92915050565b5f82825260208201905092915050565b7f496e76616c6964204b41574149206164647265737300000000000000000000005f82015250565b5f61036d601583610329565b915061037882610339565b602082019050919050565b5f6020820190508181035f83015261039a81610361565b9050919050565b608051611f5f6103c75f395f818161088f01528181610b7901526110610152611f5f5ff3fe608060405234801561000f575f5ffd5b506004361061012a575f3560e01c8063873f6f9e116100ab578063c40c91bd1161006f578063c40c91bd14610330578063c59d48471461034c578063cb56cd4f1461036d578063f2fde38b1461038b578063f75cc2b9146103a75761012a565b8063873f6f9e146102785780638a90e20f146102a85780638da5cb5b146102c4578063adeacbd3146102e2578063bff1f9e1146103125761012a565b80634a03da8a116100f25780634a03da8a146101d45780635869bc5a14610204578063715018a614610222578063727a7c5a1461022c5780637cb647591461025c5761012a565b8063060406181461012e57806307c7a72d1461014c5780630ae654031461017c5780630dac7bea146101985780632eb4a7ab146101b6575b5f5ffd5b6101366103c3565b60405161014391906114a2565b60405180910390f35b61016660048036038101906101619190611547565b6103c9565b604051610173919061159f565b60405180910390f35b610196600480360381019061019191906115eb565b61042b565b005b6101a06104a6565b6040516101ad91906114a2565b60405180910390f35b6101be6104b5565b6040516101cb9190611625565b60405180910390f35b6101ee60048036038101906101e9919061163e565b6104bb565b6040516101fb9190611625565b60405180910390f35b61020c6104d5565b60405161021991906114a2565b60405180910390f35b61022a6104db565b005b6102466004803603810190610241919061163e565b6104ee565b6040516102539190611625565b60405180910390f35b610276600480360381019061027191906115eb565b610503565b005b610292600480360381019061028d9190611547565b61056b565b60405161029f919061159f565b60405180910390f35b6102c260048036038101906102bd91906116ca565b610595565b005b6102cc610a47565b6040516102d9919061174a565b60405180910390f35b6102fc60048036038101906102f79190611763565b610a6e565b604051610309919061159f565b60405180910390f35b61031a610a8b565b60405161032791906114a2565b60405180910390f35b61034a6004803603810190610345919061178e565b610a91565b005b610354610b43565b60405161036494939291906117cc565b60405180910390f35b610375610b77565b604051610382919061186a565b60405180910390f35b6103a560048036038101906103a09190611763565b610b9b565b005b6103c160048036038101906103bc919061192d565b610c1f565b005b60015481565b5f60045f8481526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b6104336111cb565b60015f81548092919061044590611a0a565b9190505550806002819055508060035f60015481526020019081526020015f20819055506001547fa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf68260405161049b9190611625565b60405180910390a250565b6aa56fa5b99019a5c800000081565b60025481565b5f60035f8381526020019081526020015f20549050919050565b60065481565b6104e36111cb565b6104ec5f611252565b565b6003602052805f5260405f205f915090505481565b61050b6111cb565b6001547f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c60025483604051610541929190611a51565b60405180910390a2806002819055508060035f60015481526020019081526020015f208190555050565b6004602052815f5260405f20602052805f5260405f205f915091509054906101000a900460ff1681565b61059d611313565b6001548411156105e2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105d990611ad2565b60405180910390fd5b60045f8581526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff161561067b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067290611b3a565b60405180910390fd5b5f83116106bd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106b490611ba2565b60405180910390fd5b6aa56fa5b99019a5c8000000836006546106d79190611bc0565b1115610718576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161070f90611c3d565b60405180910390fd5b5f84338560405160200161072e93929190611cc0565b6040516020818303038152906040528051906020012090505f60035f8781526020019081526020015f205490505f5f1b810361079f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161079690611d46565b60405180910390fd5b6107ea8484808060200260200160405190810160405280939291908181526020018383602002808284375f81840152601f19601f820116905080830192505050505050508284611335565b610829576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082090611dae565b60405180910390fd5b600160045f8881526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166340c10f1933876040518363ffffffff1660e01b81526004016108e8929190611dcc565b5f604051808303815f87803b1580156108ff575f5ffd5b505af1158015610911573d5f5f3e3d5ffd5b505050508460065f8282546109269190611bc0565b9250508190555060055f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff166109e857600160055f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555060075f8154809291906109e290611a0a565b91905055505b3373ffffffffffffffffffffffffffffffffffffffff16867f81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f87604051610a2f91906114a2565b60405180910390a35050610a4161134b565b50505050565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6005602052805f5260405f205f915054906101000a900460ff1681565b60075481565b610a996111cb565b600154821115610ade576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ad590611ad2565b60405180910390fd5b817f1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c60035f8581526020019081526020015f205483604051610b21929190611a51565b60405180910390a28060035f8481526020019081526020015f20819055505050565b5f5f5f5f6001546006546006546aa56fa5b99019a5c8000000610b669190611df3565b600754935093509350935090919293565b7f000000000000000000000000000000000000000000000000000000000000000081565b610ba36111cb565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610c13575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610c0a919061174a565b60405180910390fd5b610c1c81611252565b50565b610c27611313565b8383905086869050148015610c4157508181905086869050145b610c80576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7790611e70565b60405180910390fd5b5f5f90505f5f90505b87879050811015610fc1575f888883818110610ca857610ca7611e8e565b5b9050602002013590505f878784818110610cc557610cc4611e8e565b5b90506020020135905060045f8381526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1615610d33575050610fb4565b600154821115610d78576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d6f90611ad2565b60405180910390fd5b5f8111610dba576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610db190611ba2565b60405180910390fd5b5f823383604051602001610dd093929190611cc0565b6040516020818303038152906040528051906020012090505f60035f8581526020019081526020015f205490505f5f1b8103610e41576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e3890611d46565b60405180910390fd5b610eaf888887818110610e5757610e56611e8e565b5b9050602002810190610e699190611ec7565b808060200260200160405190810160405280939291908181526020018383602002808284375f81840152601f19601f820116905080830192505050505050508284611335565b610eee576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ee590611dae565b60405180910390fd5b600160045f8681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508286610f5e9190611bc0565b95503373ffffffffffffffffffffffffffffffffffffffff16847f81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f85604051610fa791906114a2565b60405180910390a3505050505b8080600101915050610c89565b505f8111611004576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ffb90611ba2565b60405180910390fd5b6aa56fa5b99019a5c80000008160065461101e9190611bc0565b111561105f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161105690611c3d565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166340c10f1933836040518363ffffffff1660e01b81526004016110ba929190611dcc565b5f604051808303815f87803b1580156110d1575f5ffd5b505af11580156110e3573d5f5f3e3d5ffd5b505050508060065f8282546110f89190611bc0565b9250508190555060055f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff166111ba57600160055f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555060075f8154809291906111b490611a0a565b91905055505b506111c361134b565b505050505050565b6111d3611365565b73ffffffffffffffffffffffffffffffffffffffff166111f1610a47565b73ffffffffffffffffffffffffffffffffffffffff161461125057611214611365565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401611247919061174a565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b61131b61136c565b600261132d6113286113ad565b6113d6565b5f0181905550565b5f8261134185846113df565b1490509392505050565b600161135d6113586113ad565b6113d6565b5f0181905550565b5f33905090565b611374611430565b156113ab576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5f8290505f5f90505b8451811015611425576114168286838151811061140957611408611e8e565b5b602002602001015161144c565b915080806001019150506113e9565b508091505092915050565b5f600261144361143e6113ad565b6113d6565b5f015414905090565b5f8183106114635761145e8284611476565b61146e565b61146d8383611476565b5b905092915050565b5f825f528160205260405f20905092915050565b5f819050919050565b61149c8161148a565b82525050565b5f6020820190506114b55f830184611493565b92915050565b5f5ffd5b5f5ffd5b6114cc8161148a565b81146114d6575f5ffd5b50565b5f813590506114e7816114c3565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f611516826114ed565b9050919050565b6115268161150c565b8114611530575f5ffd5b50565b5f813590506115418161151d565b92915050565b5f5f6040838503121561155d5761155c6114bb565b5b5f61156a858286016114d9565b925050602061157b85828601611533565b9150509250929050565b5f8115159050919050565b61159981611585565b82525050565b5f6020820190506115b25f830184611590565b92915050565b5f819050919050565b6115ca816115b8565b81146115d4575f5ffd5b50565b5f813590506115e5816115c1565b92915050565b5f60208284031215611600576115ff6114bb565b5b5f61160d848285016115d7565b91505092915050565b61161f816115b8565b82525050565b5f6020820190506116385f830184611616565b92915050565b5f60208284031215611653576116526114bb565b5b5f611660848285016114d9565b91505092915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f84011261168a57611689611669565b5b8235905067ffffffffffffffff8111156116a7576116a661166d565b5b6020830191508360208202830111156116c3576116c2611671565b5b9250929050565b5f5f5f5f606085870312156116e2576116e16114bb565b5b5f6116ef878288016114d9565b9450506020611700878288016114d9565b935050604085013567ffffffffffffffff811115611721576117206114bf565b5b61172d87828801611675565b925092505092959194509250565b6117448161150c565b82525050565b5f60208201905061175d5f83018461173b565b92915050565b5f60208284031215611778576117776114bb565b5b5f61178584828501611533565b91505092915050565b5f5f604083850312156117a4576117a36114bb565b5b5f6117b1858286016114d9565b92505060206117c2858286016115d7565b9150509250929050565b5f6080820190506117df5f830187611493565b6117ec6020830186611493565b6117f96040830185611493565b6118066060830184611493565b95945050505050565b5f819050919050565b5f61183261182d611828846114ed565b61180f565b6114ed565b9050919050565b5f61184382611818565b9050919050565b5f61185482611839565b9050919050565b6118648161184a565b82525050565b5f60208201905061187d5f83018461185b565b92915050565b5f5f83601f84011261189857611897611669565b5b8235905067ffffffffffffffff8111156118b5576118b461166d565b5b6020830191508360208202830111156118d1576118d0611671565b5b9250929050565b5f5f83601f8401126118ed576118ec611669565b5b8235905067ffffffffffffffff81111561190a5761190961166d565b5b60208301915083602082028301111561192657611925611671565b5b9250929050565b5f5f5f5f5f5f60608789031215611947576119466114bb565b5b5f87013567ffffffffffffffff811115611964576119636114bf565b5b61197089828a01611883565b9650965050602087013567ffffffffffffffff811115611993576119926114bf565b5b61199f89828a01611883565b9450945050604087013567ffffffffffffffff8111156119c2576119c16114bf565b5b6119ce89828a016118d8565b92509250509295509295509295565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611a148261148a565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611a4657611a456119dd565b5b600182019050919050565b5f604082019050611a645f830185611616565b611a716020830184611616565b9392505050565b5f82825260208201905092915050565b7f496e76616c696420706572696f640000000000000000000000000000000000005f82015250565b5f611abc600e83611a78565b9150611ac782611a88565b602082019050919050565b5f6020820190508181035f830152611ae981611ab0565b9050919050565b7f416c726561647920636c61696d656420666f72207468697320706572696f64005f82015250565b5f611b24601f83611a78565b9150611b2f82611af0565b602082019050919050565b5f6020820190508181035f830152611b5181611b18565b9050919050565b7f4e6f20636173686261636b20746f20636c61696d0000000000000000000000005f82015250565b5f611b8c601483611a78565b9150611b9782611b58565b602082019050919050565b5f6020820190508181035f830152611bb981611b80565b9050919050565b5f611bca8261148a565b9150611bd58361148a565b9250828201905080821115611bed57611bec6119dd565b5b92915050565b7f4578636565647320746f74616c20616c6c6f636174696f6e00000000000000005f82015250565b5f611c27601883611a78565b9150611c3282611bf3565b602082019050919050565b5f6020820190508181035f830152611c5481611c1b565b9050919050565b5f819050919050565b611c75611c708261148a565b611c5b565b82525050565b5f8160601b9050919050565b5f611c9182611c7b565b9050919050565b5f611ca282611c87565b9050919050565b611cba611cb58261150c565b611c98565b82525050565b5f611ccb8286611c64565b602082019150611cdb8285611ca9565b601482019150611ceb8284611c64565b602082019150819050949350505050565b7f506572696f64206e6f7420736574746c656400000000000000000000000000005f82015250565b5f611d30601283611a78565b9150611d3b82611cfc565b602082019050919050565b5f6020820190508181035f830152611d5d81611d24565b9050919050565b7f496e76616c69642070726f6f66000000000000000000000000000000000000005f82015250565b5f611d98600d83611a78565b9150611da382611d64565b602082019050919050565b5f6020820190508181035f830152611dc581611d8c565b9050919050565b5f604082019050611ddf5f83018561173b565b611dec6020830184611493565b9392505050565b5f611dfd8261148a565b9150611e088361148a565b9250828203905081811115611e2057611e1f6119dd565b5b92915050565b7f4172726179206c656e677468206d69736d6174636800000000000000000000005f82015250565b5f611e5a601583611a78565b9150611e6582611e26565b602082019050919050565b5f6020820190508181035f830152611e8781611e4e565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83356001602003843603038112611ee357611ee2611ebb565b5b80840192508235915067ffffffffffffffff821115611f0557611f04611ebf565b5b602083019250602082023603831315611f2157611f20611ec3565b5b50925092905056fea264697066735822122083c22a8e29db1b036a3acaf6f30b26e5f33e1262fdf19ae9c62521498febbdfe64736f6c634300081e0033",
}

// DepositCashbackDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositCashbackDistributorMetaData.ABI instead.
var DepositCashbackDistributorABI = DepositCashbackDistributorMetaData.ABI

// DepositCashbackDistributorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepositCashbackDistributorMetaData.Bin instead.
var DepositCashbackDistributorBin = DepositCashbackDistributorMetaData.Bin

// DeployDepositCashbackDistributor deploys a new Ethereum contract, binding an instance of DepositCashbackDistributor to it.
func DeployDepositCashbackDistributor(auth *bind.TransactOpts, backend bind.ContractBackend, kawaiToken_ common.Address) (common.Address, *types.Transaction, *DepositCashbackDistributor, error) {
	parsed, err := DepositCashbackDistributorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepositCashbackDistributorBin), backend, kawaiToken_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DepositCashbackDistributor{DepositCashbackDistributorCaller: DepositCashbackDistributorCaller{contract: contract}, DepositCashbackDistributorTransactor: DepositCashbackDistributorTransactor{contract: contract}, DepositCashbackDistributorFilterer: DepositCashbackDistributorFilterer{contract: contract}}, nil
}

// DepositCashbackDistributor is an auto generated Go binding around an Ethereum contract.
type DepositCashbackDistributor struct {
	DepositCashbackDistributorCaller     // Read-only binding to the contract
	DepositCashbackDistributorTransactor // Write-only binding to the contract
	DepositCashbackDistributorFilterer   // Log filterer for contract events
}

// DepositCashbackDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositCashbackDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositCashbackDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositCashbackDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositCashbackDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositCashbackDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositCashbackDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositCashbackDistributorSession struct {
	Contract     *DepositCashbackDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DepositCashbackDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositCashbackDistributorCallerSession struct {
	Contract *DepositCashbackDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// DepositCashbackDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositCashbackDistributorTransactorSession struct {
	Contract     *DepositCashbackDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// DepositCashbackDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositCashbackDistributorRaw struct {
	Contract *DepositCashbackDistributor // Generic contract binding to access the raw methods on
}

// DepositCashbackDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositCashbackDistributorCallerRaw struct {
	Contract *DepositCashbackDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// DepositCashbackDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositCashbackDistributorTransactorRaw struct {
	Contract *DepositCashbackDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositCashbackDistributor creates a new instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributor(address common.Address, backend bind.ContractBackend) (*DepositCashbackDistributor, error) {
	contract, err := bindDepositCashbackDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributor{DepositCashbackDistributorCaller: DepositCashbackDistributorCaller{contract: contract}, DepositCashbackDistributorTransactor: DepositCashbackDistributorTransactor{contract: contract}, DepositCashbackDistributorFilterer: DepositCashbackDistributorFilterer{contract: contract}}, nil
}

// NewDepositCashbackDistributorCaller creates a new read-only instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributorCaller(address common.Address, caller bind.ContractCaller) (*DepositCashbackDistributorCaller, error) {
	contract, err := bindDepositCashbackDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorCaller{contract: contract}, nil
}

// NewDepositCashbackDistributorTransactor creates a new write-only instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositCashbackDistributorTransactor, error) {
	contract, err := bindDepositCashbackDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorTransactor{contract: contract}, nil
}

// NewDepositCashbackDistributorFilterer creates a new log filterer instance of DepositCashbackDistributor, bound to a specific deployed contract.
func NewDepositCashbackDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositCashbackDistributorFilterer, error) {
	contract, err := bindDepositCashbackDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorFilterer{contract: contract}, nil
}

// bindDepositCashbackDistributor binds a generic wrapper to an already deployed contract.
func bindDepositCashbackDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositCashbackDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositCashbackDistributor *DepositCashbackDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositCashbackDistributor.Contract.DepositCashbackDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositCashbackDistributor *DepositCashbackDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.DepositCashbackDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositCashbackDistributor *DepositCashbackDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.DepositCashbackDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositCashbackDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.contract.Transact(opts, method, params...)
}

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) TOTALALLOCATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "TOTAL_ALLOCATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) TOTALALLOCATION() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TOTALALLOCATION(&_DepositCashbackDistributor.CallOpts)
}

// TOTALALLOCATION is a free data retrieval call binding the contract method 0x0dac7bea.
//
// Solidity: function TOTAL_ALLOCATION() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) TOTALALLOCATION() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TOTALALLOCATION(&_DepositCashbackDistributor.CallOpts)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) CurrentPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "currentPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) CurrentPeriod() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.CurrentPeriod(&_DepositCashbackDistributor.CallOpts)
}

// CurrentPeriod is a free data retrieval call binding the contract method 0x06040618.
//
// Solidity: function currentPeriod() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) CurrentPeriod() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.CurrentPeriod(&_DepositCashbackDistributor.CallOpts)
}

// GetPeriodMerkleRoot is a free data retrieval call binding the contract method 0x4a03da8a.
//
// Solidity: function getPeriodMerkleRoot(uint256 period) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) GetPeriodMerkleRoot(opts *bind.CallOpts, period *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "getPeriodMerkleRoot", period)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPeriodMerkleRoot is a free data retrieval call binding the contract method 0x4a03da8a.
//
// Solidity: function getPeriodMerkleRoot(uint256 period) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) GetPeriodMerkleRoot(period *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.GetPeriodMerkleRoot(&_DepositCashbackDistributor.CallOpts, period)
}

// GetPeriodMerkleRoot is a free data retrieval call binding the contract method 0x4a03da8a.
//
// Solidity: function getPeriodMerkleRoot(uint256 period) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) GetPeriodMerkleRoot(period *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.GetPeriodMerkleRoot(&_DepositCashbackDistributor.CallOpts, period)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 _currentPeriod, uint256 _totalKawaiDistributed, uint256 _remainingAllocation, uint256 _totalUsers)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) GetStats(opts *bind.CallOpts) (struct {
	CurrentPeriod         *big.Int
	TotalKawaiDistributed *big.Int
	RemainingAllocation   *big.Int
	TotalUsers            *big.Int
}, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "getStats")

	outstruct := new(struct {
		CurrentPeriod         *big.Int
		TotalKawaiDistributed *big.Int
		RemainingAllocation   *big.Int
		TotalUsers            *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CurrentPeriod = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalKawaiDistributed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RemainingAllocation = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalUsers = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 _currentPeriod, uint256 _totalKawaiDistributed, uint256 _remainingAllocation, uint256 _totalUsers)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) GetStats() (struct {
	CurrentPeriod         *big.Int
	TotalKawaiDistributed *big.Int
	RemainingAllocation   *big.Int
	TotalUsers            *big.Int
}, error) {
	return _DepositCashbackDistributor.Contract.GetStats(&_DepositCashbackDistributor.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() view returns(uint256 _currentPeriod, uint256 _totalKawaiDistributed, uint256 _remainingAllocation, uint256 _totalUsers)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) GetStats() (struct {
	CurrentPeriod         *big.Int
	TotalKawaiDistributed *big.Int
	RemainingAllocation   *big.Int
	TotalUsers            *big.Int
}, error) {
	return _DepositCashbackDistributor.Contract.GetStats(&_DepositCashbackDistributor.CallOpts)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) HasClaimed(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "hasClaimed", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimed(&_DepositCashbackDistributor.CallOpts, arg0, arg1)
}

// HasClaimed is a free data retrieval call binding the contract method 0x873f6f9e.
//
// Solidity: function hasClaimed(uint256 , address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) HasClaimed(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimed(&_DepositCashbackDistributor.CallOpts, arg0, arg1)
}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) HasClaimedAnyPeriod(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "hasClaimedAnyPeriod", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) HasClaimedAnyPeriod(arg0 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimedAnyPeriod(&_DepositCashbackDistributor.CallOpts, arg0)
}

// HasClaimedAnyPeriod is a free data retrieval call binding the contract method 0xadeacbd3.
//
// Solidity: function hasClaimedAnyPeriod(address ) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) HasClaimedAnyPeriod(arg0 common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasClaimedAnyPeriod(&_DepositCashbackDistributor.CallOpts, arg0)
}

// HasUserClaimed is a free data retrieval call binding the contract method 0x07c7a72d.
//
// Solidity: function hasUserClaimed(uint256 period, address user) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) HasUserClaimed(opts *bind.CallOpts, period *big.Int, user common.Address) (bool, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "hasUserClaimed", period, user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasUserClaimed is a free data retrieval call binding the contract method 0x07c7a72d.
//
// Solidity: function hasUserClaimed(uint256 period, address user) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) HasUserClaimed(period *big.Int, user common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasUserClaimed(&_DepositCashbackDistributor.CallOpts, period, user)
}

// HasUserClaimed is a free data retrieval call binding the contract method 0x07c7a72d.
//
// Solidity: function hasUserClaimed(uint256 period, address user) view returns(bool)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) HasUserClaimed(period *big.Int, user common.Address) (bool, error) {
	return _DepositCashbackDistributor.Contract.HasUserClaimed(&_DepositCashbackDistributor.CallOpts, period, user)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) KawaiToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "kawaiToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) KawaiToken() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.KawaiToken(&_DepositCashbackDistributor.CallOpts)
}

// KawaiToken is a free data retrieval call binding the contract method 0xcb56cd4f.
//
// Solidity: function kawaiToken() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) KawaiToken() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.KawaiToken(&_DepositCashbackDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) MerkleRoot() ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.MerkleRoot(&_DepositCashbackDistributor.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) MerkleRoot() ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.MerkleRoot(&_DepositCashbackDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) Owner() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.Owner(&_DepositCashbackDistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) Owner() (common.Address, error) {
	return _DepositCashbackDistributor.Contract.Owner(&_DepositCashbackDistributor.CallOpts)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) PeriodMerkleRoots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "periodMerkleRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.PeriodMerkleRoots(&_DepositCashbackDistributor.CallOpts, arg0)
}

// PeriodMerkleRoots is a free data retrieval call binding the contract method 0x727a7c5a.
//
// Solidity: function periodMerkleRoots(uint256 ) view returns(bytes32)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) PeriodMerkleRoots(arg0 *big.Int) ([32]byte, error) {
	return _DepositCashbackDistributor.Contract.PeriodMerkleRoots(&_DepositCashbackDistributor.CallOpts, arg0)
}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) TotalKawaiDistributed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "totalKawaiDistributed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) TotalKawaiDistributed() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalKawaiDistributed(&_DepositCashbackDistributor.CallOpts)
}

// TotalKawaiDistributed is a free data retrieval call binding the contract method 0x5869bc5a.
//
// Solidity: function totalKawaiDistributed() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) TotalKawaiDistributed() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalKawaiDistributed(&_DepositCashbackDistributor.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCaller) TotalUsers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositCashbackDistributor.contract.Call(opts, &out, "totalUsers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) TotalUsers() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalUsers(&_DepositCashbackDistributor.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_DepositCashbackDistributor *DepositCashbackDistributorCallerSession) TotalUsers() (*big.Int, error) {
	return _DepositCashbackDistributor.Contract.TotalUsers(&_DepositCashbackDistributor.CallOpts)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) AdvancePeriod(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "advancePeriod", _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.AdvancePeriod(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// AdvancePeriod is a paid mutator transaction binding the contract method 0x0ae65403.
//
// Solidity: function advancePeriod(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) AdvancePeriod(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.AdvancePeriod(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// ClaimCashback is a paid mutator transaction binding the contract method 0x8a90e20f.
//
// Solidity: function claimCashback(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) ClaimCashback(opts *bind.TransactOpts, period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "claimCashback", period, kawaiAmount, merkleProof)
}

// ClaimCashback is a paid mutator transaction binding the contract method 0x8a90e20f.
//
// Solidity: function claimCashback(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) ClaimCashback(period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimCashback(&_DepositCashbackDistributor.TransactOpts, period, kawaiAmount, merkleProof)
}

// ClaimCashback is a paid mutator transaction binding the contract method 0x8a90e20f.
//
// Solidity: function claimCashback(uint256 period, uint256 kawaiAmount, bytes32[] merkleProof) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) ClaimCashback(period *big.Int, kawaiAmount *big.Int, merkleProof [][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimCashback(&_DepositCashbackDistributor.TransactOpts, period, kawaiAmount, merkleProof)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) ClaimMultiplePeriods(opts *bind.TransactOpts, periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "claimMultiplePeriods", periods, kawaiAmounts, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) ClaimMultiplePeriods(periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimMultiplePeriods(&_DepositCashbackDistributor.TransactOpts, periods, kawaiAmounts, merkleProofs)
}

// ClaimMultiplePeriods is a paid mutator transaction binding the contract method 0xf75cc2b9.
//
// Solidity: function claimMultiplePeriods(uint256[] periods, uint256[] kawaiAmounts, bytes32[][] merkleProofs) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) ClaimMultiplePeriods(periods []*big.Int, kawaiAmounts []*big.Int, merkleProofs [][][32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.ClaimMultiplePeriods(&_DepositCashbackDistributor.TransactOpts, periods, kawaiAmounts, merkleProofs)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.RenounceOwnership(&_DepositCashbackDistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.RenounceOwnership(&_DepositCashbackDistributor.TransactOpts)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) SetMerkleRoot(opts *bind.TransactOpts, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "setMerkleRoot", _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetMerkleRoot(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) SetMerkleRoot(_merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetMerkleRoot(&_DepositCashbackDistributor.TransactOpts, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) SetPeriodMerkleRoot(opts *bind.TransactOpts, period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "setPeriodMerkleRoot", period, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) SetPeriodMerkleRoot(period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetPeriodMerkleRoot(&_DepositCashbackDistributor.TransactOpts, period, _merkleRoot)
}

// SetPeriodMerkleRoot is a paid mutator transaction binding the contract method 0xc40c91bd.
//
// Solidity: function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) SetPeriodMerkleRoot(period *big.Int, _merkleRoot [32]byte) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.SetPeriodMerkleRoot(&_DepositCashbackDistributor.TransactOpts, period, _merkleRoot)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DepositCashbackDistributor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.TransferOwnership(&_DepositCashbackDistributor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DepositCashbackDistributor *DepositCashbackDistributorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DepositCashbackDistributor.Contract.TransferOwnership(&_DepositCashbackDistributor.TransactOpts, newOwner)
}

// DepositCashbackDistributorCashbackClaimedIterator is returned from FilterCashbackClaimed and is used to iterate over the raw logs and unpacked data for CashbackClaimed events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorCashbackClaimedIterator struct {
	Event *DepositCashbackDistributorCashbackClaimed // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorCashbackClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorCashbackClaimed)
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
		it.Event = new(DepositCashbackDistributorCashbackClaimed)
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
func (it *DepositCashbackDistributorCashbackClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorCashbackClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorCashbackClaimed represents a CashbackClaimed event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorCashbackClaimed struct {
	Period      *big.Int
	User        common.Address
	KawaiAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCashbackClaimed is a free log retrieval operation binding the contract event 0x81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f.
//
// Solidity: event CashbackClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterCashbackClaimed(opts *bind.FilterOpts, period []*big.Int, user []common.Address) (*DepositCashbackDistributorCashbackClaimedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "CashbackClaimed", periodRule, userRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorCashbackClaimedIterator{contract: _DepositCashbackDistributor.contract, event: "CashbackClaimed", logs: logs, sub: sub}, nil
}

// WatchCashbackClaimed is a free log subscription operation binding the contract event 0x81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f.
//
// Solidity: event CashbackClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchCashbackClaimed(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorCashbackClaimed, period []*big.Int, user []common.Address) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "CashbackClaimed", periodRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorCashbackClaimed)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "CashbackClaimed", log); err != nil {
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

// ParseCashbackClaimed is a log parse operation binding the contract event 0x81c5a7a76a0b67c33105d78bfd703d22da4934380ad2800d95ba6e5b87bd735f.
//
// Solidity: event CashbackClaimed(uint256 indexed period, address indexed user, uint256 kawaiAmount)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParseCashbackClaimed(log types.Log) (*DepositCashbackDistributorCashbackClaimed, error) {
	event := new(DepositCashbackDistributorCashbackClaimed)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "CashbackClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorMerkleRootUpdatedIterator is returned from FilterMerkleRootUpdated and is used to iterate over the raw logs and unpacked data for MerkleRootUpdated events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorMerkleRootUpdatedIterator struct {
	Event *DepositCashbackDistributorMerkleRootUpdated // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorMerkleRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorMerkleRootUpdated)
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
		it.Event = new(DepositCashbackDistributorMerkleRootUpdated)
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
func (it *DepositCashbackDistributorMerkleRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorMerkleRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorMerkleRootUpdated represents a MerkleRootUpdated event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorMerkleRootUpdated struct {
	Period  *big.Int
	OldRoot [32]byte
	NewRoot [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMerkleRootUpdated is a free log retrieval operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterMerkleRootUpdated(opts *bind.FilterOpts, period []*big.Int) (*DepositCashbackDistributorMerkleRootUpdatedIterator, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorMerkleRootUpdatedIterator{contract: _DepositCashbackDistributor.contract, event: "MerkleRootUpdated", logs: logs, sub: sub}, nil
}

// WatchMerkleRootUpdated is a free log subscription operation binding the contract event 0x1cb89f7d8697e1d5c6f3bcdfeb0272652e14939019b16dd05e212084b79d337c.
//
// Solidity: event MerkleRootUpdated(uint256 indexed period, bytes32 oldRoot, bytes32 newRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchMerkleRootUpdated(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorMerkleRootUpdated, period []*big.Int) (event.Subscription, error) {

	var periodRule []interface{}
	for _, periodItem := range period {
		periodRule = append(periodRule, periodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "MerkleRootUpdated", periodRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorMerkleRootUpdated)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
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
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParseMerkleRootUpdated(log types.Log) (*DepositCashbackDistributorMerkleRootUpdated, error) {
	event := new(DepositCashbackDistributorMerkleRootUpdated)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "MerkleRootUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorOwnershipTransferredIterator struct {
	Event *DepositCashbackDistributorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorOwnershipTransferred)
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
		it.Event = new(DepositCashbackDistributorOwnershipTransferred)
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
func (it *DepositCashbackDistributorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorOwnershipTransferred represents a OwnershipTransferred event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DepositCashbackDistributorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorOwnershipTransferredIterator{contract: _DepositCashbackDistributor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorOwnershipTransferred)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParseOwnershipTransferred(log types.Log) (*DepositCashbackDistributorOwnershipTransferred, error) {
	event := new(DepositCashbackDistributorOwnershipTransferred)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositCashbackDistributorPeriodAdvancedIterator is returned from FilterPeriodAdvanced and is used to iterate over the raw logs and unpacked data for PeriodAdvanced events raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorPeriodAdvancedIterator struct {
	Event *DepositCashbackDistributorPeriodAdvanced // Event containing the contract specifics and raw log

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
func (it *DepositCashbackDistributorPeriodAdvancedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositCashbackDistributorPeriodAdvanced)
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
		it.Event = new(DepositCashbackDistributorPeriodAdvanced)
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
func (it *DepositCashbackDistributorPeriodAdvancedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositCashbackDistributorPeriodAdvancedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositCashbackDistributorPeriodAdvanced represents a PeriodAdvanced event raised by the DepositCashbackDistributor contract.
type DepositCashbackDistributorPeriodAdvanced struct {
	NewPeriod  *big.Int
	MerkleRoot [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPeriodAdvanced is a free log retrieval operation binding the contract event 0xa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf6.
//
// Solidity: event PeriodAdvanced(uint256 indexed newPeriod, bytes32 merkleRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) FilterPeriodAdvanced(opts *bind.FilterOpts, newPeriod []*big.Int) (*DepositCashbackDistributorPeriodAdvancedIterator, error) {

	var newPeriodRule []interface{}
	for _, newPeriodItem := range newPeriod {
		newPeriodRule = append(newPeriodRule, newPeriodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.FilterLogs(opts, "PeriodAdvanced", newPeriodRule)
	if err != nil {
		return nil, err
	}
	return &DepositCashbackDistributorPeriodAdvancedIterator{contract: _DepositCashbackDistributor.contract, event: "PeriodAdvanced", logs: logs, sub: sub}, nil
}

// WatchPeriodAdvanced is a free log subscription operation binding the contract event 0xa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf6.
//
// Solidity: event PeriodAdvanced(uint256 indexed newPeriod, bytes32 merkleRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) WatchPeriodAdvanced(opts *bind.WatchOpts, sink chan<- *DepositCashbackDistributorPeriodAdvanced, newPeriod []*big.Int) (event.Subscription, error) {

	var newPeriodRule []interface{}
	for _, newPeriodItem := range newPeriod {
		newPeriodRule = append(newPeriodRule, newPeriodItem)
	}

	logs, sub, err := _DepositCashbackDistributor.contract.WatchLogs(opts, "PeriodAdvanced", newPeriodRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositCashbackDistributorPeriodAdvanced)
				if err := _DepositCashbackDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
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

// ParsePeriodAdvanced is a log parse operation binding the contract event 0xa3fa2e7f4d459160c2a2988ce319e83f8b535f0ab7dade6bc39c8786ca009cf6.
//
// Solidity: event PeriodAdvanced(uint256 indexed newPeriod, bytes32 merkleRoot)
func (_DepositCashbackDistributor *DepositCashbackDistributorFilterer) ParsePeriodAdvanced(log types.Log) (*DepositCashbackDistributorPeriodAdvanced, error) {
	event := new(DepositCashbackDistributorPeriodAdvanced)
	if err := _DepositCashbackDistributor.contract.UnpackLog(event, "PeriodAdvanced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
