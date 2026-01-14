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
	Bin: "0x60c03461017e57601f6111ad38819003918201601f19168301916001600160401b038311848410176101825780849260609460405283398101031261017e5761004781610196565b90610060604061005960208401610196565b9201610196565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556001600160a01b039290919083169081156101395783169081156100f45760805260a0521660018060a01b03195f5416175f5560405161100290816101ab82396080518181816103ff015281816105ab015281816107910152610c77015260a0518181816108a60152610c4c0152f35b60405162461bcd60e51b815260206004820152601460248201527f496e76616c6964205553445420616464726573730000000000000000000000006044820152606490fd5b60405162461bcd60e51b815260206004820152601560248201527f496e76616c696420746f6b656e206164647265737300000000000000000000006044820152606490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b038216820361017e5756fe604060808152600480361015610013575f80fd5b5f90813560e01c8063036520271461092557806322f85eaa146108d55780632f48ab7d146108915780634690484014610869578063514fcac71461071d57806379109baa1461057a5780637c95cdc61461042e57806396d875dc146103ea578063a85c38ef146103bc578063b5b3b0511461039d578063bf333f2c14610382578063d09ef241146102fb578063f2e8553a1461016f5763f9567f5d146100b7575f80fd5b3461016b576100c536610a64565b916100ce610f7d565b6100db6001548310610b40565b806100e583610a7e565b506100f660ff600583015416610b7f565b610101851515610acb565b0154831161012857509061011491610bdf565b60015f80516020610fad8339815191525580f35b606490602086519162461bcd60e51b8352820152601860248201527f416d6f756e7420657863656564732072656d61696e696e6700000000000000006044820152fd5b5080fd5b503461016b57606036600319011261016b576001600160a01b0391903582811691908290036102f85760243590604435936101ad6064861115610e58565b819382916001958654935b8481106102bf575085808211156102b6576101d291610bbe565b965b8088116102ae575b506101e687610d53565b96849384885b610202575b8a51806101fe8c826109e7565b0390f35b818110806102a5575b156102a057908185858b610220819796610a7e565b500154161461023b575b61023390610dd7565b9091926101ec565b968a888a83101561025d575b505061025561023391610dd7565b97905061022a565b8861028f610233949a6102559461027f61027961029596610a7e565b50610e09565b6102898383610df5565b52610df5565b50610dd7565b979150508a88610247565b6101f1565b5082861061020b565b96505f6101dc565b505083966101d4565b8383896102cb84610a7e565b50015416146102e3575b6102de90610dd7565b6101b8565b906102f06102de91610dd7565b9190506102d5565b80fd5b509190346102f85760203660031901126102f857506101fe61032a83356103256001548210610b40565b610a7e565b5080546001820154600283015460038401549684015460059094015495519283526001600160a01b03909116602083015260408201526060810194909452608084015260ff909116151560a0830152819060c0820190565b82843461016b578160031936011261016b5751908152602090f35b82843461016b578160031936011261016b576020906001549051908152f35b509190346102f85760203660031901126102f8578235906001548210156102f8575061032a6101fe91610a7e565b82843461016b578160031936011261016b57517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b5050346102f85761043e36610a64565b919061044d6064841115610e58565b81806001938454915b828110610541575083808211156105385761047091610bbe565b945b808611610530575b5061048485610d53565b94829182865b61049c575b8851806101fe8a826109e7565b81811080610527575b156105225760ff60056104b783610a7e565b500154166104cf575b6104c990610dd7565b8661048a565b93858110156104ed575b6104e56104c991610dd7565b9490506104c0565b926104e56105196104c99261050461027989610a7e565b61050e828d610df5565b5261028f818c610df5565b949150506104d9565b61048f565b508284106104a5565b94505f61047a565b50508194610472565b60ff600561054e83610a7e565b50015416610565575b61056090610dd7565b610456565b9061057261056091610dd7565b919050610557565b50913461016b5761058a36610a64565b610592610f7d565b61059d821515610acb565b80156106e6576105cf8230337f0000000000000000000000000000000000000000000000000000000000000000610f13565b6001928354916105dd610b0c565b838152602081019033825283810186815260608201848152608083019188835260a08401948a8652680100000000000000008910156106d3578a89018b5561062489610a7e565b9590956106c057928d92600595927ff7c110a6973307f2bc91245c2c06344ada13add2c1741e83ac5c0bb332bc85d59a999897955186558d86019060018060a01b039051166bffffffffffffffffffffffff60a01b8254161790555160028501555160038401555190820155019051151560ff80198354169116179055815194855260208501523393a35f80516020610fad8339815191525580f35b50634e487b7160e01b8c528b8d5260248cfd5b634e487b7160e01b8c5260418d5260248cfd5b825162461bcd60e51b8152602081870152601160248201527005072696365206d757374206265203e203607c1b6044820152606490fd5b5082346108655760203660031901126108655781359161073b610f7d565b6107486001548410610b40565b61075183610a7e565b5060018101549091906001600160a01b0316330361083157600582019182549360ff8516156107ef575001805460ff199093169091558390556107b590337f0000000000000000000000000000000000000000000000000000000000000000610e95565b33907fc0362da6f2ff36b382b34aec0814f6b3cdf89f5ef282a1d1f114d0c0b036d5968380a360015f80516020610fad8339815191525580f35b5162461bcd60e51b8152602081840152601c60248201527f4f7264657220616c726561647920736f6c642f63616e63656c6c6564000000006044820152606490fd5b606490602084519162461bcd60e51b8352820152600e60248201526d2737ba103cb7bab91037b93232b960911b6044820152fd5b8280fd5b82843461016b578160031936011261016b57905490516001600160a01b039091168152602090f35b82843461016b578160031936011261016b57517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b503461016b57602036600319011261016b57806101149135906108f6610f7d565b6109036001548310610b40565b61090c82610a7e565b5061091d60ff600583015416610b7f565b015490610bdf565b50913461016b57602036600319011261016b57823567ffffffffffffffff938482116109e357366023830112156109e357810135938411610865576024810190602436918660051b0101116108655761097d84610d53565b9290600154915b858110610998578351806101fe87826109e7565b806109b2846109ab6109de948a87610de5565b3510610b40565b6109c96102796109c3838a87610de5565b35610a7e565b6109d38288610df5565b5261028f8187610df5565b610984565b8380fd5b60208082019080835283518092528060408094019401925f905b838210610a1057505050505090565b845180518752808401516001600160a01b0316878501528082015187830152606080820151908801526080808201519088015260a09081015115159087015260c09095019493820193600190910190610a01565b6040906003190112610a7a576004359060243590565b5f80fd5b600154811015610ab75760069060015f52027fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601905f90565b634e487b7160e01b5f52603260045260245ffd5b15610ad257565b60405162461bcd60e51b81526020600482015260126024820152710416d6f756e74206d757374206265203e20360741b6044820152606490fd5b6040519060c0820182811067ffffffffffffffff821117610b2c57604052565b634e487b7160e01b5f52604160045260245ffd5b15610b4757565b60405162461bcd60e51b815260206004820152601060248201526f125b9d985b1a590813dc99195c88125160821b6044820152606490fd5b15610b8657565b60405162461bcd60e51b815260206004820152601060248201526f4f72646572206e6f742061637469766560801b6044820152606490fd5b91908203918211610bcb57565b634e487b7160e01b5f52601160045260245ffd5b90610be982610a7e565b5091600383015482810290808204841490151715610bcb576002840154908115610d275704926004810190610c1f848354610bbe565b80835515610d16575b600185151715610bcb5760010180546001600160a01b039290610c709087908516337f0000000000000000000000000000000000000000000000000000000000000000610f13565b610c9b85337f0000000000000000000000000000000000000000000000000000000000000000610e95565b5480610cd8575054169260405192835260208301527fe847ce2e2eb43b46eebf1b3aa5cd5a85a80e2537dc01a5fe9e48038508ec0d4460403393a4565b949190541693604051938452602084015260408301527f17e85045a994095c87543a467b7b4b48abc118c93a651b291e2748e50f6e5b1460603393a4565b60058101805460ff19169055610c28565b634e487b7160e01b5f52601260045260245ffd5b67ffffffffffffffff8111610b2c5760051b60200190565b90610d5d82610d3b565b604080519091601f1990601f018116820167ffffffffffffffff811183821017610b2c578352848252610d908295610d3b565b01915f5b838110610da15750505050565b602090610dac610b0c565b5f8152825f818301525f858301525f60608301525f60808301525f60a0830152828601015201610d94565b5f198114610bcb5760010190565b9190811015610ab75760051b0190565b8051821015610ab75760209160051b010190565b9060ff6005610e16610b0c565b8454815260018501546001600160a01b0316602082015260028501546040820152600385015460608201526004850154608082015293015416151560a0830152565b15610e5f57565b60405162461bcd60e51b815260206004820152600e60248201526d098d2dad2e840e8dede40d0d2ced60931b6044820152606490fd5b60405163a9059cbb60e01b5f9081526001600160a01b039384166004526024949094529260209060448180855af160015f5114811615610ef4575b8360405215610ede57505050565b635274afe760e01b835216600482015260249150fd5b6001811516610f0a57813b15153d151616610ed0565b833d5f823e3d90fd5b6040516323b872dd60e01b5f9081526001600160a01b03938416600452938316602452604494909452909160209060648180855af160015f5114811615610f67575b836040525f60605215610ede57505050565b6001811516610f0a57813b15153d151616610f55565b5f80516020610fad8339815191526002815414610f9a5760029055565b604051633ee5aeb560e01b8152600490fdfe9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00a26469706673582212205b9264dcbe6c9cd8cc2c2ea783c7640ae172781dd67e0a3c41603b689d88ae0064736f6c63430008140033",
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
