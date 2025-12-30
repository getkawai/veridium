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

// OTCMarketOrder is an auto generated low-level Go binding around an user-defined struct.
type OTCMarketOrder struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
}

// OTCMarketMetaData contains all meta data concerning the OTCMarket contract.
var OTCMarketMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_tokenDeAI\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_usdt\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"FEE_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buyOrder\",\"inputs\":[{\"name\":\"_orderId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"buyOrderPartial\",\"inputs\":[{\"name\":\"_orderId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOrder\",\"inputs\":[{\"name\":\"_orderId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createOrder\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveOrders\",\"inputs\":[{\"name\":\"_offset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOTCMarket.Order[]\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOrder\",\"inputs\":[{\"name\":\"_orderId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOrders\",\"inputs\":[{\"name\":\"_orderIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOTCMarket.Order[]\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOrdersBySeller\",\"inputs\":[{\"name\":\"_seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_offset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOTCMarket.Order[]\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOrdersCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"orders\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"priceInUSDT\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenDeAI\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usdt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OrderCancelled\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrderCreated\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrderFulfilled\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrderPartiallyFilled\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amountFilled\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"remainingAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"pricePaid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60c060405234801561000f575f5ffd5b506040516126fc3803806126fc83398181016040528101906100319190610270565b600161004f6100446101e060201b60201c565b61020960201b60201c565b5f01819055505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036100c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100ba9061031a565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610131576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161012890610382565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508173ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1681525050805f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050506103a0565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61023f82610216565b9050919050565b61024f81610235565b8114610259575f5ffd5b50565b5f8151905061026a81610246565b92915050565b5f5f5f6060848603121561028757610286610212565b5b5f6102948682870161025c565b93505060206102a58682870161025c565b92505060406102b68682870161025c565b9150509250925092565b5f82825260208201905092915050565b7f496e76616c696420746f6b656e206164647265737300000000000000000000005f82015250565b5f6103046015836102c0565b915061030f826102d0565b602082019050919050565b5f6020820190508181035f830152610331816102f8565b9050919050565b7f496e76616c6964205553445420616464726573730000000000000000000000005f82015250565b5f61036c6014836102c0565b915061037782610338565b602082019050919050565b5f6020820190508181035f83015261039981610360565b9050919050565b60805160a0516123186103e45f395f81816105a60152818161145001526114c101525f81816107720152818161089b01528181610cf2015261150d01526123185ff3fe608060405234801561000f575f5ffd5b50600436106100e8575f3560e01c806396d875dc1161008a578063bf333f2c11610064578063bf333f2c1461024d578063d09ef2411461026b578063f2e8553a146102a0578063f9567f5d146102d0576100e8565b806396d875dc146101dc578063a85c38ef146101fa578063b5b3b0511461022f576100e8565b806346904840116100c65780634690484014610156578063514fcac71461017457806379109baa146101905780637c95cdc6146101ac576100e8565b806303652027146100ec57806322f85eaa1461011c5780632f48ab7d14610138575b5f5ffd5b6101066004803603810190610101919061195b565b6102ec565b6040516101139190611b38565b60405180910390f35b61013660048036038101906101319190611b82565b6104c8565b005b6101406105a4565b60405161014d9190611c08565b60405180910390f35b61015e6105c8565b60405161016b9190611c30565b60405180910390f35b61018e60048036038101906101899190611b82565b6105ec565b005b6101aa60048036038101906101a59190611c49565b610807565b005b6101c660048036038101906101c19190611c49565b610a42565b6040516101d39190611b38565b60405180910390f35b6101e4610cf0565b6040516101f19190611c08565b60405180910390f35b610214600480360381019061020f9190611b82565b610d14565b60405161022696959493929190611ca5565b60405180910390f35b610237610d86565b6040516102449190611d04565b60405180910390f35b610255610d92565b6040516102629190611d04565b60405180910390f35b61028560048036038101906102809190611b82565b610d96565b60405161029796959493929190611ca5565b60405180910390f35b6102ba60048036038101906102b59190611d47565b610e64565b6040516102c79190611b38565b60405180910390f35b6102ea60048036038101906102e59190611c49565b611193565b005b60605f8383905067ffffffffffffffff81111561030c5761030b611d97565b5b60405190808252806020026020018201604052801561034557816020015b6103326118aa565b81526020019060019003908161032a5790505b5090505f5f90505b848490508110156104bd5760018054905085858381811061037157610370611dc4565b5b90506020020135106103b8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103af90611e4b565b60405180910390fd5b60018585838181106103cd576103cc611dc4565b5b90506020020135815481106103e5576103e4611dc4565b5b905f5260205f2090600602016040518060c00160405290815f8201548152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015481526020016003820154815260200160048201548152602001600582015f9054906101000a900460ff1615151515815250508282815181106104a5576104a4611dc4565b5b6020026020010181905250808060010191505061034d565b508091505092915050565b6104d06112f5565b6001805490508110610517576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161050e90611e4b565b60405180910390fd5b5f6001828154811061052c5761052b611dc4565b5b905f5260205f2090600602019050806005015f9054906101000a900460ff1661058a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058190611eb3565b60405180910390fd5b610598828260040154611317565b506105a1611686565b50565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6105f46112f5565b600180549050811061063b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161063290611e4b565b60405180910390fd5b5f600182815481106106505761064f611dc4565b5b905f5260205f20906006020190503373ffffffffffffffffffffffffffffffffffffffff16816001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146106ef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106e690611f1b565b60405180910390fd5b806005015f9054906101000a900460ff1661073f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073690611f83565b60405180910390fd5b5f816004015490505f826005015f6101000a81548160ff0219169083151502179055505f82600401819055506107b633827f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166116a09092919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff16837fc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d59660405160405180910390a35050610804611686565b50565b61080f6112f5565b5f8211610851576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161084890611feb565b60405180910390fd5b5f8111610893576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161088a90612053565b60405180910390fd5b6108e03330847f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166116f3909392919063ffffffff16565b5f600180549050905060016040518060c001604052808381526020013373ffffffffffffffffffffffffffffffffffffffff16815260200185815260200184815260200185815260200160011515815250908060018154018082558091505060019003905f5260205f2090600602015f909190919091505f820151815f01556020820151816001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160020155606082015181600301556080820151816004015560a0820151816005015f6101000a81548160ff02191690831515021790555050503373ffffffffffffffffffffffffffffffffffffffff16817ff7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d58585604051610a2d929190612071565b60405180910390a350610a3e611686565b5050565b60606064821115610a88576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a7f906120e2565b60405180910390fd5b5f5f90505f5f90505b600180549050811015610aef5760018181548110610ab257610ab1611dc4565b5b905f5260205f2090600602016005015f9054906101000a900460ff1615610ae2578180610ade9061212d565b9250505b8080600101915050610a91565b505f848211610afe575f610b0b565b8482610b0a9190612174565b5b905083811115610b19578390505b5f8167ffffffffffffffff811115610b3457610b33611d97565b5b604051908082528060200260200182016040528015610b6d57816020015b610b5a6118aa565b815260200190600190039081610b525790505b5090505f5f90505f5f90505f5f90505b60018054905081108015610b9057508483105b15610ce15760018181548110610ba957610ba8611dc4565b5b905f5260205f2090600602016005015f9054906101000a900460ff1615610cce57888210610cbf5760018181548110610be557610be4611dc4565b5b905f5260205f2090600602016040518060c00160405290815f8201548152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015481526020016003820154815260200160048201548152602001600582015f9054906101000a900460ff161515151581525050848481518110610ca557610ca4611dc4565b5b60200260200101819052508280610cbb9061212d565b9350505b8180610cca9061212d565b9250505b8080610cd99061212d565b915050610b7d565b50829550505050505092915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60018181548110610d23575f80fd5b905f5260205f2090600602015f91509050805f015490806001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806002015490806003015490806004015490806005015f9054906101000a900460ff16905086565b5f600180549050905090565b5f81565b5f5f5f5f5f5f6001805490508710610de3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dda90611e4b565b60405180910390fd5b5f60018881548110610df857610df7611dc4565b5b905f5260205f2090600602019050805f0154816001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16826002015483600301548460040154856005015f9054906101000a900460ff169650965096509650965096505091939550919395565b60606064821115610eaa576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ea1906120e2565b60405180910390fd5b5f5f90505f5f90505b600180549050811015610f51578573ffffffffffffffffffffffffffffffffffffffff1660018281548110610eeb57610eea611dc4565b5b905f5260205f2090600602016001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610f44578180610f409061212d565b9250505b8080600101915050610eb3565b505f848211610f60575f610f6d565b8482610f6c9190612174565b5b905083811115610f7b578390505b5f8167ffffffffffffffff811115610f9657610f95611d97565b5b604051908082528060200260200182016040528015610fcf57816020015b610fbc6118aa565b815260200190600190039081610fb45790505b5090505f5f90505f5f90505f5f90505b60018054905081108015610ff257508483105b15611183578973ffffffffffffffffffffffffffffffffffffffff166001828154811061102257611021611dc4565b5b905f5260205f2090600602016001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361117057888210611161576001818154811061108757611086611dc4565b5b905f5260205f2090600602016040518060c00160405290815f8201548152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015481526020016003820154815260200160048201548152602001600582015f9054906101000a900460ff16151515158152505084848151811061114757611146611dc4565b5b6020026020010181905250828061115d9061212d565b9350505b818061116c9061212d565b9250505b808061117b9061212d565b915050610fdf565b5082955050505050509392505050565b61119b6112f5565b60018054905082106111e2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111d990611e4b565b60405180910390fd5b5f600183815481106111f7576111f6611dc4565b5b905f5260205f2090600602019050806005015f9054906101000a900460ff16611255576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161124c90611eb3565b60405180910390fd5b5f8211611297576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161128e90611feb565b60405180910390fd5b80600401548211156112de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112d5906121f1565b60405180910390fd5b6112e88383611317565b506112f1611686565b5050565b6112fd611748565b600261130f61130a611789565b6117b2565b5f0181905550565b5f6001838154811061132c5761132b611dc4565b5b905f5260205f20906006020190505f8160020154838360030154611350919061220f565b61135a919061227d565b905082826004015f82825461136f9190612174565b925050819055505f82600401540361139d575f826005015f6101000a81548160ff0219169083151502179055505b5f6127105f836113ad919061220f565b6113b7919061227d565b90505f81836113c69190612174565b90505f8211801561142357505f73ffffffffffffffffffffffffffffffffffffffff165f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614155b1561149657611495335f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16847f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166116f3909392919063ffffffff16565b5b61150633856001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166116f3909392919063ffffffff16565b61155133867f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166116a09092919063ffffffff16565b5f8460040154036115ec57836001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16877fe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d4488876040516115df929190612071565b60405180910390a461167e565b836001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16877f17e85045a994095c87543a467b7b4b48abc118c93a651b291e2748e50f6e5b1488886004015488604051611675939291906122ad565b60405180910390a45b505050505050565b6001611698611693611789565b6117b2565b5f0181905550565b6116ad83838360016117bb565b6116ee57826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016116e59190611c30565b60405180910390fd5b505050565b61170184848484600161181d565b61174257836040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016117399190611c30565b60405180910390fd5b50505050565b61175061188e565b15611787576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b5f5f63a9059cbb60e01b9050604051815f525f1960601c86166004528460245260205f60445f5f8b5af1925060015f5114831661180f578383151615611803573d5f823e3d81fd5b5f873b113d1516831692505b806040525050949350505050565b5f5f6323b872dd60e01b9050604051815f525f1960601c87166004525f1960601c86166024528460445260205f60645f5f8c5af1925060015f5114831661187b57838315161561186f573d5f823e3d81fd5b5f883b113d1516831692505b806040525f606052505095945050505050565b5f60026118a161189c611789565b6117b2565b5f015414905090565b6040518060c001604052805f81526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f81526020015f81526020015f81526020015f151581525090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f84011261191b5761191a6118fa565b5b8235905067ffffffffffffffff811115611938576119376118fe565b5b60208301915083602082028301111561195457611953611902565b5b9250929050565b5f5f60208385031215611971576119706118f2565b5b5f83013567ffffffffffffffff81111561198e5761198d6118f6565b5b61199a85828601611906565b92509250509250929050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b5f819050919050565b6119e1816119cf565b82525050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f611a10826119e7565b9050919050565b611a2081611a06565b82525050565b5f8115159050919050565b611a3a81611a26565b82525050565b60c082015f820151611a545f8501826119d8565b506020820151611a676020850182611a17565b506040820151611a7a60408501826119d8565b506060820151611a8d60608501826119d8565b506080820151611aa060808501826119d8565b5060a0820151611ab360a0850182611a31565b50505050565b5f611ac48383611a40565b60c08301905092915050565b5f602082019050919050565b5f611ae6826119a6565b611af081856119b0565b9350611afb836119c0565b805f5b83811015611b2b578151611b128882611ab9565b9750611b1d83611ad0565b925050600181019050611afe565b5085935050505092915050565b5f6020820190508181035f830152611b508184611adc565b905092915050565b611b61816119cf565b8114611b6b575f5ffd5b50565b5f81359050611b7c81611b58565b92915050565b5f60208284031215611b9757611b966118f2565b5b5f611ba484828501611b6e565b91505092915050565b5f819050919050565b5f611bd0611bcb611bc6846119e7565b611bad565b6119e7565b9050919050565b5f611be182611bb6565b9050919050565b5f611bf282611bd7565b9050919050565b611c0281611be8565b82525050565b5f602082019050611c1b5f830184611bf9565b92915050565b611c2a81611a06565b82525050565b5f602082019050611c435f830184611c21565b92915050565b5f5f60408385031215611c5f57611c5e6118f2565b5b5f611c6c85828601611b6e565b9250506020611c7d85828601611b6e565b9150509250929050565b611c90816119cf565b82525050565b611c9f81611a26565b82525050565b5f60c082019050611cb85f830189611c87565b611cc56020830188611c21565b611cd26040830187611c87565b611cdf6060830186611c87565b611cec6080830185611c87565b611cf960a0830184611c96565b979650505050505050565b5f602082019050611d175f830184611c87565b92915050565b611d2681611a06565b8114611d30575f5ffd5b50565b5f81359050611d4181611d1d565b92915050565b5f5f5f60608486031215611d5e57611d5d6118f2565b5b5f611d6b86828701611d33565b9350506020611d7c86828701611b6e565b9250506040611d8d86828701611b6e565b9150509250925092565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f82825260208201905092915050565b7f496e76616c6964204f72646572204944000000000000000000000000000000005f82015250565b5f611e35601083611df1565b9150611e4082611e01565b602082019050919050565b5f6020820190508181035f830152611e6281611e29565b9050919050565b7f4f72646572206e6f7420616374697665000000000000000000000000000000005f82015250565b5f611e9d601083611df1565b9150611ea882611e69565b602082019050919050565b5f6020820190508181035f830152611eca81611e91565b9050919050565b7f4e6f7420796f7572206f726465720000000000000000000000000000000000005f82015250565b5f611f05600e83611df1565b9150611f1082611ed1565b602082019050919050565b5f6020820190508181035f830152611f3281611ef9565b9050919050565b7f4f7264657220616c726561647920736f6c642f63616e63656c6c6564000000005f82015250565b5f611f6d601c83611df1565b9150611f7882611f39565b602082019050919050565b5f6020820190508181035f830152611f9a81611f61565b9050919050565b7f416d6f756e74206d757374206265203e203000000000000000000000000000005f82015250565b5f611fd5601283611df1565b9150611fe082611fa1565b602082019050919050565b5f6020820190508181035f83015261200281611fc9565b9050919050565b7f5072696365206d757374206265203e20300000000000000000000000000000005f82015250565b5f61203d601183611df1565b915061204882612009565b602082019050919050565b5f6020820190508181035f83015261206a81612031565b9050919050565b5f6040820190506120845f830185611c87565b6120916020830184611c87565b9392505050565b7f4c696d697420746f6f20686967680000000000000000000000000000000000005f82015250565b5f6120cc600e83611df1565b91506120d782612098565b602082019050919050565b5f6020820190508181035f8301526120f9816120c0565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f612137826119cf565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361216957612168612100565b5b600182019050919050565b5f61217e826119cf565b9150612189836119cf565b92508282039050818111156121a1576121a0612100565b5b92915050565b7f416d6f756e7420657863656564732072656d61696e696e6700000000000000005f82015250565b5f6121db601883611df1565b91506121e6826121a7565b602082019050919050565b5f6020820190508181035f830152612208816121cf565b9050919050565b5f612219826119cf565b9150612224836119cf565b9250828202612232816119cf565b9150828204841483151761224957612248612100565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f612287826119cf565b9150612292836119cf565b9250826122a2576122a1612250565b5b828204905092915050565b5f6060820190506122c05f830186611c87565b6122cd6020830185611c87565b6122da6040830184611c87565b94935050505056fea2646970667358221220251d472006b8ffab061efcb0013871343fc23ade54e46a74cbb5c07eb645cced64736f6c634300081e0033",
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

// GetActiveOrders is a free data retrieval call binding the contract method 0x7c95cdc6.
//
// Solidity: function getActiveOrders(uint256 _offset, uint256 _limit) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketCaller) GetActiveOrders(opts *bind.CallOpts, _offset *big.Int, _limit *big.Int) ([]OTCMarketOrder, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "getActiveOrders", _offset, _limit)

	if err != nil {
		return *new([]OTCMarketOrder), err
	}

	out0 := *abi.ConvertType(out[0], new([]OTCMarketOrder)).(*[]OTCMarketOrder)

	return out0, err

}

// GetActiveOrders is a free data retrieval call binding the contract method 0x7c95cdc6.
//
// Solidity: function getActiveOrders(uint256 _offset, uint256 _limit) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketSession) GetActiveOrders(_offset *big.Int, _limit *big.Int) ([]OTCMarketOrder, error) {
	return _OTCMarket.Contract.GetActiveOrders(&_OTCMarket.CallOpts, _offset, _limit)
}

// GetActiveOrders is a free data retrieval call binding the contract method 0x7c95cdc6.
//
// Solidity: function getActiveOrders(uint256 _offset, uint256 _limit) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketCallerSession) GetActiveOrders(_offset *big.Int, _limit *big.Int) ([]OTCMarketOrder, error) {
	return _OTCMarket.Contract.GetActiveOrders(&_OTCMarket.CallOpts, _offset, _limit)
}

// GetOrder is a free data retrieval call binding the contract method 0xd09ef241.
//
// Solidity: function getOrder(uint256 _orderId) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, uint256 remainingAmount, bool isActive)
func (_OTCMarket *OTCMarketCaller) GetOrder(opts *bind.CallOpts, _orderId *big.Int) (struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
}, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "getOrder", _orderId)

	outstruct := new(struct {
		Id              *big.Int
		Seller          common.Address
		TokenAmount     *big.Int
		PriceInUSDT     *big.Int
		RemainingAmount *big.Int
		IsActive        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Seller = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokenAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PriceInUSDT = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RemainingAmount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.IsActive = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// GetOrder is a free data retrieval call binding the contract method 0xd09ef241.
//
// Solidity: function getOrder(uint256 _orderId) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, uint256 remainingAmount, bool isActive)
func (_OTCMarket *OTCMarketSession) GetOrder(_orderId *big.Int) (struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
}, error) {
	return _OTCMarket.Contract.GetOrder(&_OTCMarket.CallOpts, _orderId)
}

// GetOrder is a free data retrieval call binding the contract method 0xd09ef241.
//
// Solidity: function getOrder(uint256 _orderId) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, uint256 remainingAmount, bool isActive)
func (_OTCMarket *OTCMarketCallerSession) GetOrder(_orderId *big.Int) (struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
}, error) {
	return _OTCMarket.Contract.GetOrder(&_OTCMarket.CallOpts, _orderId)
}

// GetOrders is a free data retrieval call binding the contract method 0x03652027.
//
// Solidity: function getOrders(uint256[] _orderIds) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketCaller) GetOrders(opts *bind.CallOpts, _orderIds []*big.Int) ([]OTCMarketOrder, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "getOrders", _orderIds)

	if err != nil {
		return *new([]OTCMarketOrder), err
	}

	out0 := *abi.ConvertType(out[0], new([]OTCMarketOrder)).(*[]OTCMarketOrder)

	return out0, err

}

// GetOrders is a free data retrieval call binding the contract method 0x03652027.
//
// Solidity: function getOrders(uint256[] _orderIds) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketSession) GetOrders(_orderIds []*big.Int) ([]OTCMarketOrder, error) {
	return _OTCMarket.Contract.GetOrders(&_OTCMarket.CallOpts, _orderIds)
}

// GetOrders is a free data retrieval call binding the contract method 0x03652027.
//
// Solidity: function getOrders(uint256[] _orderIds) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketCallerSession) GetOrders(_orderIds []*big.Int) ([]OTCMarketOrder, error) {
	return _OTCMarket.Contract.GetOrders(&_OTCMarket.CallOpts, _orderIds)
}

// GetOrdersBySeller is a free data retrieval call binding the contract method 0xf2e8553a.
//
// Solidity: function getOrdersBySeller(address _seller, uint256 _offset, uint256 _limit) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketCaller) GetOrdersBySeller(opts *bind.CallOpts, _seller common.Address, _offset *big.Int, _limit *big.Int) ([]OTCMarketOrder, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "getOrdersBySeller", _seller, _offset, _limit)

	if err != nil {
		return *new([]OTCMarketOrder), err
	}

	out0 := *abi.ConvertType(out[0], new([]OTCMarketOrder)).(*[]OTCMarketOrder)

	return out0, err

}

// GetOrdersBySeller is a free data retrieval call binding the contract method 0xf2e8553a.
//
// Solidity: function getOrdersBySeller(address _seller, uint256 _offset, uint256 _limit) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketSession) GetOrdersBySeller(_seller common.Address, _offset *big.Int, _limit *big.Int) ([]OTCMarketOrder, error) {
	return _OTCMarket.Contract.GetOrdersBySeller(&_OTCMarket.CallOpts, _seller, _offset, _limit)
}

// GetOrdersBySeller is a free data retrieval call binding the contract method 0xf2e8553a.
//
// Solidity: function getOrdersBySeller(address _seller, uint256 _offset, uint256 _limit) view returns((uint256,address,uint256,uint256,uint256,bool)[])
func (_OTCMarket *OTCMarketCallerSession) GetOrdersBySeller(_seller common.Address, _offset *big.Int, _limit *big.Int) ([]OTCMarketOrder, error) {
	return _OTCMarket.Contract.GetOrdersBySeller(&_OTCMarket.CallOpts, _seller, _offset, _limit)
}

// GetOrdersCount is a free data retrieval call binding the contract method 0xb5b3b051.
//
// Solidity: function getOrdersCount() view returns(uint256)
func (_OTCMarket *OTCMarketCaller) GetOrdersCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "getOrdersCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOrdersCount is a free data retrieval call binding the contract method 0xb5b3b051.
//
// Solidity: function getOrdersCount() view returns(uint256)
func (_OTCMarket *OTCMarketSession) GetOrdersCount() (*big.Int, error) {
	return _OTCMarket.Contract.GetOrdersCount(&_OTCMarket.CallOpts)
}

// GetOrdersCount is a free data retrieval call binding the contract method 0xb5b3b051.
//
// Solidity: function getOrdersCount() view returns(uint256)
func (_OTCMarket *OTCMarketCallerSession) GetOrdersCount() (*big.Int, error) {
	return _OTCMarket.Contract.GetOrdersCount(&_OTCMarket.CallOpts)
}

// Orders is a free data retrieval call binding the contract method 0xa85c38ef.
//
// Solidity: function orders(uint256 ) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, uint256 remainingAmount, bool isActive)
func (_OTCMarket *OTCMarketCaller) Orders(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
}, error) {
	var out []interface{}
	err := _OTCMarket.contract.Call(opts, &out, "orders", arg0)

	outstruct := new(struct {
		Id              *big.Int
		Seller          common.Address
		TokenAmount     *big.Int
		PriceInUSDT     *big.Int
		RemainingAmount *big.Int
		IsActive        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Seller = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokenAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PriceInUSDT = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RemainingAmount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.IsActive = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// Orders is a free data retrieval call binding the contract method 0xa85c38ef.
//
// Solidity: function orders(uint256 ) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, uint256 remainingAmount, bool isActive)
func (_OTCMarket *OTCMarketSession) Orders(arg0 *big.Int) (struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
}, error) {
	return _OTCMarket.Contract.Orders(&_OTCMarket.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0xa85c38ef.
//
// Solidity: function orders(uint256 ) view returns(uint256 id, address seller, uint256 tokenAmount, uint256 priceInUSDT, uint256 remainingAmount, bool isActive)
func (_OTCMarket *OTCMarketCallerSession) Orders(arg0 *big.Int) (struct {
	Id              *big.Int
	Seller          common.Address
	TokenAmount     *big.Int
	PriceInUSDT     *big.Int
	RemainingAmount *big.Int
	IsActive        bool
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

// BuyOrderPartial is a paid mutator transaction binding the contract method 0xf9567f5d.
//
// Solidity: function buyOrderPartial(uint256 _orderId, uint256 _amount) returns()
func (_OTCMarket *OTCMarketTransactor) BuyOrderPartial(opts *bind.TransactOpts, _orderId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _OTCMarket.contract.Transact(opts, "buyOrderPartial", _orderId, _amount)
}

// BuyOrderPartial is a paid mutator transaction binding the contract method 0xf9567f5d.
//
// Solidity: function buyOrderPartial(uint256 _orderId, uint256 _amount) returns()
func (_OTCMarket *OTCMarketSession) BuyOrderPartial(_orderId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.BuyOrderPartial(&_OTCMarket.TransactOpts, _orderId, _amount)
}

// BuyOrderPartial is a paid mutator transaction binding the contract method 0xf9567f5d.
//
// Solidity: function buyOrderPartial(uint256 _orderId, uint256 _amount) returns()
func (_OTCMarket *OTCMarketTransactorSession) BuyOrderPartial(_orderId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _OTCMarket.Contract.BuyOrderPartial(&_OTCMarket.TransactOpts, _orderId, _amount)
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

// OTCMarketOrderPartiallyFilledIterator is returned from FilterOrderPartiallyFilled and is used to iterate over the raw logs and unpacked data for OrderPartiallyFilled events raised by the OTCMarket contract.
type OTCMarketOrderPartiallyFilledIterator struct {
	Event *OTCMarketOrderPartiallyFilled // Event containing the contract specifics and raw log

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
func (it *OTCMarketOrderPartiallyFilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OTCMarketOrderPartiallyFilled)
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
		it.Event = new(OTCMarketOrderPartiallyFilled)
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
func (it *OTCMarketOrderPartiallyFilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OTCMarketOrderPartiallyFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OTCMarketOrderPartiallyFilled represents a OrderPartiallyFilled event raised by the OTCMarket contract.
type OTCMarketOrderPartiallyFilled struct {
	OrderId         *big.Int
	Buyer           common.Address
	Seller          common.Address
	AmountFilled    *big.Int
	RemainingAmount *big.Int
	PricePaid       *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOrderPartiallyFilled is a free log retrieval operation binding the contract event 0x17e85045a994095c87543a467b7b4b48abc118c93a651b291e2748e50f6e5b14.
//
// Solidity: event OrderPartiallyFilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amountFilled, uint256 remainingAmount, uint256 pricePaid)
func (_OTCMarket *OTCMarketFilterer) FilterOrderPartiallyFilled(opts *bind.FilterOpts, orderId []*big.Int, buyer []common.Address, seller []common.Address) (*OTCMarketOrderPartiallyFilledIterator, error) {

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

	logs, sub, err := _OTCMarket.contract.FilterLogs(opts, "OrderPartiallyFilled", orderIdRule, buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &OTCMarketOrderPartiallyFilledIterator{contract: _OTCMarket.contract, event: "OrderPartiallyFilled", logs: logs, sub: sub}, nil
}

// WatchOrderPartiallyFilled is a free log subscription operation binding the contract event 0x17e85045a994095c87543a467b7b4b48abc118c93a651b291e2748e50f6e5b14.
//
// Solidity: event OrderPartiallyFilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amountFilled, uint256 remainingAmount, uint256 pricePaid)
func (_OTCMarket *OTCMarketFilterer) WatchOrderPartiallyFilled(opts *bind.WatchOpts, sink chan<- *OTCMarketOrderPartiallyFilled, orderId []*big.Int, buyer []common.Address, seller []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _OTCMarket.contract.WatchLogs(opts, "OrderPartiallyFilled", orderIdRule, buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OTCMarketOrderPartiallyFilled)
				if err := _OTCMarket.contract.UnpackLog(event, "OrderPartiallyFilled", log); err != nil {
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

// ParseOrderPartiallyFilled is a log parse operation binding the contract event 0x17e85045a994095c87543a467b7b4b48abc118c93a651b291e2748e50f6e5b14.
//
// Solidity: event OrderPartiallyFilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amountFilled, uint256 remainingAmount, uint256 pricePaid)
func (_OTCMarket *OTCMarketFilterer) ParseOrderPartiallyFilled(log types.Log) (*OTCMarketOrderPartiallyFilled, error) {
	event := new(OTCMarketOrderPartiallyFilled)
	if err := _OTCMarket.contract.UnpackLog(event, "OrderPartiallyFilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
