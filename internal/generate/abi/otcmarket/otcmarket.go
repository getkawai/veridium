// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package otcmarket

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
	Bin: "0x60c03461019157601f6111a938819003918201601f19168301916001600160401b038311848410176101955780849260609460405283398101031261019157610047816101a9565b906100606040610059602084016101a9565b92016101a9565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055916001600160a01b031690811561014c576001600160a01b03169081156101075760805260a0525f80546001600160a01b0319166001600160a01b0392909216919091179055604051610feb90816101be82396080518181816103d101528181610565015281816107490152610d84015260a05181818161085d0152610d590152f35b60405162461bcd60e51b815260206004820152601460248201527f496e76616c6964205553445420616464726573730000000000000000000000006044820152606490fd5b60405162461bcd60e51b815260206004820152601560248201527f496e76616c696420746f6b656e206164647265737300000000000000000000006044820152606490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101915756fe60806040526004361015610011575f80fd5b5f3560e01c806303652027146108dc57806322f85eaa1461088c5780632f48ab7d146108485780634690484014610821578063514fcac7146106da57806379109baa146105365780637c95cdc61461040057806396d875dc146103bc578063a85c38ef14610395578063b5b3b05114610378578063bf333f2c1461035e578063d09ef241146102dc578063f2e8553a1461016a5763f9567f5d146100b3575f80fd5b34610166576100c136610a1c565b6100c9610cbd565b6100d66001548310610b2f565b60046100e183610a32565b506100f260ff600583015416610bd1565b6100fd831515610c10565b015481116101215761010e91610cf5565b60015f516020610f965f395f51905f5255005b60405162461bcd60e51b815260206004820152601860248201527f416d6f756e7420657863656564732072656d61696e696e6700000000000000006044820152606490fd5b5f80fd5b34610166576060366003190112610166576004356001600160a01b03811690819003610166576044356024356101a36064831115610c51565b600154915f805b8481106102a257508280821115610299576101c491610cb0565b905b808211610291575b506101d881610a9a565b5f5f925f5b86811080610288575b1561027657876101f582610a32565b50600101546001600160a01b031614610217575b61021290610c8e565b6101dd565b9385811015610235575b61022d61021291610c8e565b949050610209565b9161022d61026d6102129261025261024c89610a32565b50610b82565b61025c8289610b6e565b526102678188610b6e565b50610c8e565b93915050610221565b6040518061028486826109a5565b0390f35b508183106101e6565b9050846101ce565b50505f906101c6565b856102ac82610a32565b50600101546001600160a01b0316146102c8575b6001016101aa565b906102d4600191610c8e565b9190506102c0565b34610166576020366003190112610166576103056004356103006001548210610b2f565b610a32565b5080546001820154600283015460038401546004850154600590950154604080519586526001600160a01b039094166020860152928401919091526060830152608082019290925260ff909116151560a082015260c090f35b34610166575f3660031901126101665760206040515f8152f35b34610166575f366003190112610166576020600154604051908152f35b34610166576020366003190112610166576004356001548110156101665761030590610a32565b34610166575f366003190112610166576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346101665761040e36610a1c565b9061041c6064831115610c51565b600154905f805b838110610504575081808211156104fb5761043d91610cb0565b925b8084116104f3575b5061045183610a9a565b925f5f915f5b858110806104ea575b156104dc5760ff600561047283610a32565b50015416610489575b61048490610c8e565b610457565b92848110156104a7575b61049f61048491610c8e565b93905061047b565b9161049f6104d3610484926104be61024c88610a32565b6104c8828c610b6e565b52610267818b610b6e565b93915050610493565b6040518061028489826109a5565b50818310610460565b925083610447565b50505f9261043f565b60ff600561051183610a32565b50015416610522575b600101610423565b9061052e600191610c8e565b91905061051a565b346101665761054436610a1c565b61054c610cbd565b610557821515610c10565b80156106a1576105898230337f0000000000000000000000000000000000000000000000000000000000000000610f2c565b60015490610595610a62565b8281526020810190338252604081019185835260608201848152608083019087825260a0840194600186526801000000000000000088101561068d57600188016001556105e188610a32565b94909461067a57600595518555600185019060018060a01b039051166bffffffffffffffffffffffff60a01b825416179055516002840155516003830155516004820155019051151560ff8019835416911617905560405192835260208301527ff7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d560403393a360015f516020610f965f395f51905f5255005b634e487b7160e01b5f525f60045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b60405162461bcd60e51b815260206004820152601160248201527005072696365206d757374206265203e203607c1b6044820152606490fd5b34610166576020366003190112610166576004356106f6610cbd565b6107036001548210610b2f565b61070c81610a32565b5060018101546001600160a01b031633036107eb57600581019081549060ff8216156107a65761076d9260045f92019081549360ff1916905555337f0000000000000000000000000000000000000000000000000000000000000000610eab565b33907fc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d5965f80a360015f516020610f965f395f51905f5255005b60405162461bcd60e51b815260206004820152601c60248201527f4f7264657220616c726561647920736f6c642f63616e63656c6c6564000000006044820152606490fd5b60405162461bcd60e51b815260206004820152600e60248201526d2737ba103cb7bab91037b93232b960911b6044820152606490fd5b34610166575f366003190112610166575f546040516001600160a01b039091168152602090f35b34610166575f366003190112610166576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346101665760203660031901126101665761010e6004356108ab610cbd565b6108b86001548210610b2f565b60046108c382610a32565b506108d460ff600583015416610bd1565b015490610cf5565b346101665760203660031901126101665760043567ffffffffffffffff811161016657366023820112156101665780600401359067ffffffffffffffff8211610166576024810190602436918460051b0101116101665761093c82610a9a565b915f91600154925b828110610959576040518061028487826109a5565b806109728561096b6001948787610b1f565b3510610b2f565b61098961024c610983838787610b1f565b35610a32565b6109938288610b6e565b5261099e8187610b6e565b5001610944565b60206040818301928281528451809452019201905f5b8181106109c85750505090565b909192602060c060019260a08751805183528580831b038582015116858401526040810151604084015260608101516060840152608081015160808401520151151560a082015201940191019190916109bb565b6040906003190112610166576004359060243590565b600154811015610a4e5760015f52600660205f20910201905f90565b634e487b7160e01b5f52603260045260245ffd5b6040519060c0820182811067ffffffffffffffff82111761068d57604052565b67ffffffffffffffff811161068d5760051b60200190565b90610aa482610a82565b60405190601f01601f1916810167ffffffffffffffff81118282101761068d576040528281528092610ad8601f1991610a82565b01905f5b828110610ae857505050565b602090610af3610a62565b5f81525f838201525f60408201525f60608201525f60808201525f60a082015282828501015201610adc565b9190811015610a4e5760051b0190565b15610b3657565b60405162461bcd60e51b815260206004820152601060248201526f125b9d985b1a590813dc99195c88125160821b6044820152606490fd5b8051821015610a4e5760209160051b010190565b9060ff6005610b8f610a62565b8454815260018501546001600160a01b0316602082015260028501546040820152600385015460608201526004850154608082015293015416151560a0830152565b15610bd857565b60405162461bcd60e51b815260206004820152601060248201526f4f72646572206e6f742061637469766560801b6044820152606490fd5b15610c1757565b60405162461bcd60e51b81526020600482015260126024820152710416d6f756e74206d757374206265203e20360741b6044820152606490fd5b15610c5857565b60405162461bcd60e51b815260206004820152600e60248201526d098d2dad2e840e8dede40d0d2ced60931b6044820152606490fd5b5f198114610c9c5760010190565b634e487b7160e01b5f52601160045260245ffd5b91908203918211610c9c57565b60025f516020610f965f395f51905f525414610ce65760025f516020610f965f395f51905f5255565b633ee5aeb560e01b5f5260045ffd5b90610cff82610a32565b5091600383015482810290808204841490151715610c9c576002840154908115610e975704928315610e435760016004820191610d3d858454610cb0565b80845515610e32575b0190610d7d8560018060a01b03845416337f0000000000000000000000000000000000000000000000000000000000000000610f2c565b610da884337f0000000000000000000000000000000000000000000000000000000000000000610eab565b5480610ded575060018060a01b039054169260405192835260208301527fe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d4460403393a4565b939060018060a01b0390541693604051938452602084015260408301527f17e85045a994095c87543a467b7b4b48abc118c93a651b291e2748e50f6e5b1460603393a4565b60058101805460ff19169055610d46565b60405162461bcd60e51b815260206004820152602660248201527f416d6f756e7420746f6f20736d616c6c2c20707269636520726f756e647320746044820152656f207a65726f60d01b6064820152608490fd5b634e487b7160e01b5f52601260045260245ffd5b916040519163a9059cbb60e01b5f5260018060a01b031660045260245260205f60448180865af19060015f5114821615610f0b575b60405215610eeb5750565b635274afe760e01b5f9081526001600160a01b0391909116600452602490fd5b906001811516610f2357823b15153d15161690610ee0565b503d5f823e3d90fd5b6040516323b872dd60e01b5f9081526001600160a01b039384166004529290931660245260449390935260209060648180865af19060015f5114821615610f7d575b6040525f60605215610eeb5750565b906001811516610f2357823b15153d15161690610f6e56fe9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00a264697066735822122053baf9f68ceba6532f963660a2aaeb948b933a6d1dd67fe767f7b58166e8529e64736f6c63430008210033",
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
