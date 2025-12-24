// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package kawaitoken

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

// KawaiTokenMetaData contains all meta data concerning the KawaiToken contract.
var KawaiTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MINTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f5ffd5b5060405161231e38038061231e833981810160405281019061003191906105a7565b6040518060400160405280600e81526020017f4b6177616920414920546f6b656e0000000000000000000000000000000000008152506040518060400160405280600581526020017f4b4157414900000000000000000000000000000000000000000000000000000081525081600390816100ac9190610822565b5080600490816100bc9190610822565b5050506100d15f5f1b8361014260201b60201c565b506101027f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a68261014260201b60201c565b5061013b8261011561023860201b60201c565b600a6101219190610a59565b6305f5e1006101309190610aa3565b61024060201b60201c565b5050610b9c565b5f61015383836102c560201b60201c565b61022e57600160055f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055506101cb61032960201b60201c565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019050610232565b5f90505b92915050565b5f6012905090565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036102b0575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016102a79190610af3565b60405180910390fd5b6102c15f838361033060201b60201c565b5050565b5f60055f8481526020019081526020015f205f015f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f33905090565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610380578060025f8282546103749190610b0c565b9250508190555061044e565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610409578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161040093929190610b4e565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610495578060025f82825403925050819055506104df565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161053c9190610b83565b60405180910390a3505050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6105768261054d565b9050919050565b6105868161056c565b8114610590575f5ffd5b50565b5f815190506105a18161057d565b92915050565b5f5f604083850312156105bd576105bc610549565b5b5f6105ca85828601610593565b92505060206105db85828601610593565b9150509250929050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061066057607f821691505b6020821081036106735761067261061c565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026106d57fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261069a565b6106df868361069a565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61072361071e610719846106f7565b610700565b6106f7565b9050919050565b5f819050919050565b61073c83610709565b6107506107488261072a565b8484546106a6565b825550505050565b5f5f905090565b610767610758565b610772818484610733565b505050565b5b818110156107955761078a5f8261075f565b600181019050610778565b5050565b601f8211156107da576107ab81610679565b6107b48461068b565b810160208510156107c3578190505b6107d76107cf8561068b565b830182610777565b50505b505050565b5f82821c905092915050565b5f6107fa5f19846008026107df565b1980831691505092915050565b5f61081283836107eb565b9150826002028217905092915050565b61082b826105e5565b67ffffffffffffffff811115610844576108436105ef565b5b61084e8254610649565b610859828285610799565b5f60209050601f83116001811461088a575f8415610878578287015190505b6108828582610807565b8655506108e9565b601f19841661089886610679565b5f5b828110156108bf5784890151825560018201915060208501945060208101905061089a565b868310156108dc57848901516108d8601f8916826107eb565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f8160011c9050919050565b5f5f8291508390505b60018511156109735780860481111561094f5761094e6108f1565b5b600185161561095e5780820291505b808102905061096c8561091e565b9450610933565b94509492505050565b5f8261098b5760019050610a46565b81610998575f9050610a46565b81600181146109ae57600281146109b8576109e7565b6001915050610a46565b60ff8411156109ca576109c96108f1565b5b8360020a9150848211156109e1576109e06108f1565b5b50610a46565b5060208310610133831016604e8410600b8410161715610a1c5782820a905083811115610a1757610a166108f1565b5b610a46565b610a29848484600161092a565b92509050818404811115610a4057610a3f6108f1565b5b81810290505b9392505050565b5f60ff82169050919050565b5f610a63826106f7565b9150610a6e83610a4d565b9250610a9b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff848461097c565b905092915050565b5f610aad826106f7565b9150610ab8836106f7565b9250828202610ac6816106f7565b91508282048414831517610add57610adc6108f1565b5b5092915050565b610aed8161056c565b82525050565b5f602082019050610b065f830184610ae4565b92915050565b5f610b16826106f7565b9150610b21836106f7565b9250828201905080821115610b3957610b386108f1565b5b92915050565b610b48816106f7565b82525050565b5f606082019050610b615f830186610ae4565b610b6e6020830185610b3f565b610b7b6040830184610b3f565b949350505050565b5f602082019050610b965f830184610b3f565b92915050565b61177580610ba95f395ff3fe608060405234801561000f575f5ffd5b506004361061012a575f3560e01c806342966c68116100ab578063a217fddf1161006f578063a217fddf14610352578063a9059cbb14610370578063d5391393146103a0578063d547741f146103be578063dd62ed3e146103da5761012a565b806342966c681461029c57806370a08231146102b857806379cc6790146102e857806391d148541461030457806395d89b41146103345761012a565b8063248a9ca3116100f2578063248a9ca3146101fa5780632f2ff15d1461022a578063313ce5671461024657806336568abe1461026457806340c10f19146102805761012a565b806301ffc9a71461012e57806306fdde031461015e578063095ea7b31461017c57806318160ddd146101ac57806323b872dd146101ca575b5f5ffd5b61014860048036038101906101439190611241565b61040a565b6040516101559190611286565b60405180910390f35b610166610483565b604051610173919061130f565b60405180910390f35b610196600480360381019061019191906113bc565b610513565b6040516101a39190611286565b60405180910390f35b6101b4610535565b6040516101c19190611409565b60405180910390f35b6101e460048036038101906101df9190611422565b61053e565b6040516101f19190611286565b60405180910390f35b610214600480360381019061020f91906114a5565b61056c565b60405161022191906114df565b60405180910390f35b610244600480360381019061023f91906114f8565b610589565b005b61024e6105ab565b60405161025b9190611551565b60405180910390f35b61027e600480360381019061027991906114f8565b6105b3565b005b61029a600480360381019061029591906113bc565b61062e565b005b6102b660048036038101906102b1919061156a565b610667565b005b6102d260048036038101906102cd9190611595565b61067b565b6040516102df9190611409565b60405180910390f35b61030260048036038101906102fd91906113bc565b6106c0565b005b61031e600480360381019061031991906114f8565b6106e0565b60405161032b9190611286565b60405180910390f35b61033c610744565b604051610349919061130f565b60405180910390f35b61035a6107d4565b60405161036791906114df565b60405180910390f35b61038a600480360381019061038591906113bc565b6107da565b6040516103979190611286565b60405180910390f35b6103a86107fc565b6040516103b591906114df565b60405180910390f35b6103d860048036038101906103d391906114f8565b610820565b005b6103f460048036038101906103ef91906115c0565b610842565b6040516104019190611409565b60405180910390f35b5f7f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061047c575061047b826108c4565b5b9050919050565b6060600380546104929061162b565b80601f01602080910402602001604051908101604052809291908181526020018280546104be9061162b565b80156105095780601f106104e057610100808354040283529160200191610509565b820191905f5260205f20905b8154815290600101906020018083116104ec57829003601f168201915b5050505050905090565b5f5f61051d61092d565b905061052a818585610934565b600191505092915050565b5f600254905090565b5f5f61054861092d565b9050610555858285610946565b6105608585856109d9565b60019150509392505050565b5f60055f8381526020019081526020015f20600101549050919050565b6105928261056c565b61059b81610ac9565b6105a58383610add565b50505050565b5f6012905090565b6105bb61092d565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461061f576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106298282610bc7565b505050565b7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a661065881610ac9565b6106628383610cb1565b505050565b61067861067261092d565b82610d30565b50565b5f5f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b6106d2826106cc61092d565b83610946565b6106dc8282610d30565b5050565b5f60055f8481526020019081526020015f205f015f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b6060600480546107539061162b565b80601f016020809104026020016040519081016040528092919081815260200182805461077f9061162b565b80156107ca5780601f106107a1576101008083540402835291602001916107ca565b820191905f5260205f20905b8154815290600101906020018083116107ad57829003601f168201915b5050505050905090565b5f5f1b81565b5f5f6107e461092d565b90506107f18185856109d9565b600191505092915050565b7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a681565b6108298261056c565b61083281610ac9565b61083c8383610bc7565b50505050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f33905090565b6109418383836001610daf565b505050565b5f6109518484610842565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8110156109d357818110156109c4578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016109bb9392919061166a565b60405180910390fd5b6109d284848484035f610daf565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610a49575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610a40919061169f565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610ab9575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610ab0919061169f565b60405180910390fd5b610ac4838383610f7e565b505050565b610ada81610ad561092d565b611197565b50565b5f610ae883836106e0565b610bbd57600160055f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550610b5a61092d565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019050610bc1565b5f90505b92915050565b5f610bd283836106e0565b15610ca7575f60055f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550610c4461092d565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a460019050610cab565b5f90505b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610d21575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610d18919061169f565b60405180910390fd5b610d2c5f8383610f7e565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610da0575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610d97919061169f565b60405180910390fd5b610dab825f83610f7e565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610e1f575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401610e16919061169f565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610e8f575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610e86919061169f565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015610f78578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610f6f9190611409565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610fce578060025f828254610fc291906116e5565b9250508190555061109c565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015611057578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161104e9392919061166a565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036110e3578060025f828254039250508190555061112d565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161118a9190611409565b60405180910390a3505050565b6111a182826106e0565b6111e45780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016111db929190611718565b60405180910390fd5b5050565b5f5ffd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b611220816111ec565b811461122a575f5ffd5b50565b5f8135905061123b81611217565b92915050565b5f60208284031215611256576112556111e8565b5b5f6112638482850161122d565b91505092915050565b5f8115159050919050565b6112808161126c565b82525050565b5f6020820190506112995f830184611277565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6112e18261129f565b6112eb81856112a9565b93506112fb8185602086016112b9565b611304816112c7565b840191505092915050565b5f6020820190508181035f83015261132781846112d7565b905092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6113588261132f565b9050919050565b6113688161134e565b8114611372575f5ffd5b50565b5f813590506113838161135f565b92915050565b5f819050919050565b61139b81611389565b81146113a5575f5ffd5b50565b5f813590506113b681611392565b92915050565b5f5f604083850312156113d2576113d16111e8565b5b5f6113df85828601611375565b92505060206113f0858286016113a8565b9150509250929050565b61140381611389565b82525050565b5f60208201905061141c5f8301846113fa565b92915050565b5f5f5f60608486031215611439576114386111e8565b5b5f61144686828701611375565b935050602061145786828701611375565b9250506040611468868287016113a8565b9150509250925092565b5f819050919050565b61148481611472565b811461148e575f5ffd5b50565b5f8135905061149f8161147b565b92915050565b5f602082840312156114ba576114b96111e8565b5b5f6114c784828501611491565b91505092915050565b6114d981611472565b82525050565b5f6020820190506114f25f8301846114d0565b92915050565b5f5f6040838503121561150e5761150d6111e8565b5b5f61151b85828601611491565b925050602061152c85828601611375565b9150509250929050565b5f60ff82169050919050565b61154b81611536565b82525050565b5f6020820190506115645f830184611542565b92915050565b5f6020828403121561157f5761157e6111e8565b5b5f61158c848285016113a8565b91505092915050565b5f602082840312156115aa576115a96111e8565b5b5f6115b784828501611375565b91505092915050565b5f5f604083850312156115d6576115d56111e8565b5b5f6115e385828601611375565b92505060206115f485828601611375565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061164257607f821691505b602082108103611655576116546115fe565b5b50919050565b6116648161134e565b82525050565b5f60608201905061167d5f83018661165b565b61168a60208301856113fa565b61169760408301846113fa565b949350505050565b5f6020820190506116b25f83018461165b565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6116ef82611389565b91506116fa83611389565b9250828201905080821115611712576117116116b8565b5b92915050565b5f60408201905061172b5f83018561165b565b61173860208301846114d0565b939250505056fea26469706673582212201dededaa127aa0e00cff09c2c4f22fbd377af76ed9723d1f5b0dd82d3e002ee664736f6c634300081e0033",
}

// KawaiTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use KawaiTokenMetaData.ABI instead.
var KawaiTokenABI = KawaiTokenMetaData.ABI

// KawaiTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KawaiTokenMetaData.Bin instead.
var KawaiTokenBin = KawaiTokenMetaData.Bin

// DeployKawaiToken deploys a new Ethereum contract, binding an instance of KawaiToken to it.
func DeployKawaiToken(auth *bind.TransactOpts, backend bind.ContractBackend, defaultAdmin common.Address, minter common.Address) (common.Address, *types.Transaction, *KawaiToken, error) {
	parsed, err := KawaiTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KawaiTokenBin), backend, defaultAdmin, minter)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KawaiToken{KawaiTokenCaller: KawaiTokenCaller{contract: contract}, KawaiTokenTransactor: KawaiTokenTransactor{contract: contract}, KawaiTokenFilterer: KawaiTokenFilterer{contract: contract}}, nil
}

// KawaiToken is an auto generated Go binding around an Ethereum contract.
type KawaiToken struct {
	KawaiTokenCaller     // Read-only binding to the contract
	KawaiTokenTransactor // Write-only binding to the contract
	KawaiTokenFilterer   // Log filterer for contract events
}

// KawaiTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type KawaiTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KawaiTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KawaiTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KawaiTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KawaiTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KawaiTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KawaiTokenSession struct {
	Contract     *KawaiToken       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KawaiTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KawaiTokenCallerSession struct {
	Contract *KawaiTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// KawaiTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KawaiTokenTransactorSession struct {
	Contract     *KawaiTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// KawaiTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type KawaiTokenRaw struct {
	Contract *KawaiToken // Generic contract binding to access the raw methods on
}

// KawaiTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KawaiTokenCallerRaw struct {
	Contract *KawaiTokenCaller // Generic read-only contract binding to access the raw methods on
}

// KawaiTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KawaiTokenTransactorRaw struct {
	Contract *KawaiTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKawaiToken creates a new instance of KawaiToken, bound to a specific deployed contract.
func NewKawaiToken(address common.Address, backend bind.ContractBackend) (*KawaiToken, error) {
	contract, err := bindKawaiToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KawaiToken{KawaiTokenCaller: KawaiTokenCaller{contract: contract}, KawaiTokenTransactor: KawaiTokenTransactor{contract: contract}, KawaiTokenFilterer: KawaiTokenFilterer{contract: contract}}, nil
}

// NewKawaiTokenCaller creates a new read-only instance of KawaiToken, bound to a specific deployed contract.
func NewKawaiTokenCaller(address common.Address, caller bind.ContractCaller) (*KawaiTokenCaller, error) {
	contract, err := bindKawaiToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenCaller{contract: contract}, nil
}

// NewKawaiTokenTransactor creates a new write-only instance of KawaiToken, bound to a specific deployed contract.
func NewKawaiTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*KawaiTokenTransactor, error) {
	contract, err := bindKawaiToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenTransactor{contract: contract}, nil
}

// NewKawaiTokenFilterer creates a new log filterer instance of KawaiToken, bound to a specific deployed contract.
func NewKawaiTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*KawaiTokenFilterer, error) {
	contract, err := bindKawaiToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenFilterer{contract: contract}, nil
}

// bindKawaiToken binds a generic wrapper to an already deployed contract.
func bindKawaiToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KawaiTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KawaiToken *KawaiTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KawaiToken.Contract.KawaiTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KawaiToken *KawaiTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KawaiToken.Contract.KawaiTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KawaiToken *KawaiTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KawaiToken.Contract.KawaiTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KawaiToken *KawaiTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KawaiToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KawaiToken *KawaiTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KawaiToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KawaiToken *KawaiTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KawaiToken.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_KawaiToken *KawaiTokenCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_KawaiToken *KawaiTokenSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _KawaiToken.Contract.DEFAULTADMINROLE(&_KawaiToken.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_KawaiToken *KawaiTokenCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _KawaiToken.Contract.DEFAULTADMINROLE(&_KawaiToken.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_KawaiToken *KawaiTokenCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_KawaiToken *KawaiTokenSession) MINTERROLE() ([32]byte, error) {
	return _KawaiToken.Contract.MINTERROLE(&_KawaiToken.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_KawaiToken *KawaiTokenCallerSession) MINTERROLE() ([32]byte, error) {
	return _KawaiToken.Contract.MINTERROLE(&_KawaiToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_KawaiToken *KawaiTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_KawaiToken *KawaiTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _KawaiToken.Contract.Allowance(&_KawaiToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_KawaiToken *KawaiTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _KawaiToken.Contract.Allowance(&_KawaiToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_KawaiToken *KawaiTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_KawaiToken *KawaiTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _KawaiToken.Contract.BalanceOf(&_KawaiToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_KawaiToken *KawaiTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _KawaiToken.Contract.BalanceOf(&_KawaiToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_KawaiToken *KawaiTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_KawaiToken *KawaiTokenSession) Decimals() (uint8, error) {
	return _KawaiToken.Contract.Decimals(&_KawaiToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_KawaiToken *KawaiTokenCallerSession) Decimals() (uint8, error) {
	return _KawaiToken.Contract.Decimals(&_KawaiToken.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_KawaiToken *KawaiTokenCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_KawaiToken *KawaiTokenSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _KawaiToken.Contract.GetRoleAdmin(&_KawaiToken.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_KawaiToken *KawaiTokenCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _KawaiToken.Contract.GetRoleAdmin(&_KawaiToken.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_KawaiToken *KawaiTokenCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_KawaiToken *KawaiTokenSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _KawaiToken.Contract.HasRole(&_KawaiToken.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_KawaiToken *KawaiTokenCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _KawaiToken.Contract.HasRole(&_KawaiToken.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_KawaiToken *KawaiTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_KawaiToken *KawaiTokenSession) Name() (string, error) {
	return _KawaiToken.Contract.Name(&_KawaiToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_KawaiToken *KawaiTokenCallerSession) Name() (string, error) {
	return _KawaiToken.Contract.Name(&_KawaiToken.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_KawaiToken *KawaiTokenCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_KawaiToken *KawaiTokenSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _KawaiToken.Contract.SupportsInterface(&_KawaiToken.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_KawaiToken *KawaiTokenCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _KawaiToken.Contract.SupportsInterface(&_KawaiToken.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_KawaiToken *KawaiTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_KawaiToken *KawaiTokenSession) Symbol() (string, error) {
	return _KawaiToken.Contract.Symbol(&_KawaiToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_KawaiToken *KawaiTokenCallerSession) Symbol() (string, error) {
	return _KawaiToken.Contract.Symbol(&_KawaiToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_KawaiToken *KawaiTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_KawaiToken *KawaiTokenSession) TotalSupply() (*big.Int, error) {
	return _KawaiToken.Contract.TotalSupply(&_KawaiToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_KawaiToken *KawaiTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _KawaiToken.Contract.TotalSupply(&_KawaiToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Approve(&_KawaiToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Approve(&_KawaiToken.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_KawaiToken *KawaiTokenTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_KawaiToken *KawaiTokenSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Burn(&_KawaiToken.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_KawaiToken *KawaiTokenTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Burn(&_KawaiToken.TransactOpts, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_KawaiToken *KawaiTokenTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "burnFrom", account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_KawaiToken *KawaiTokenSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.BurnFrom(&_KawaiToken.TransactOpts, account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_KawaiToken *KawaiTokenTransactorSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.BurnFrom(&_KawaiToken.TransactOpts, account, value)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_KawaiToken *KawaiTokenTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_KawaiToken *KawaiTokenSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _KawaiToken.Contract.GrantRole(&_KawaiToken.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_KawaiToken *KawaiTokenTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _KawaiToken.Contract.GrantRole(&_KawaiToken.TransactOpts, role, account)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_KawaiToken *KawaiTokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_KawaiToken *KawaiTokenSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Mint(&_KawaiToken.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_KawaiToken *KawaiTokenTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Mint(&_KawaiToken.TransactOpts, to, amount)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_KawaiToken *KawaiTokenTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_KawaiToken *KawaiTokenSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _KawaiToken.Contract.RenounceRole(&_KawaiToken.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_KawaiToken *KawaiTokenTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _KawaiToken.Contract.RenounceRole(&_KawaiToken.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_KawaiToken *KawaiTokenTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_KawaiToken *KawaiTokenSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _KawaiToken.Contract.RevokeRole(&_KawaiToken.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_KawaiToken *KawaiTokenTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _KawaiToken.Contract.RevokeRole(&_KawaiToken.TransactOpts, role, account)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Transfer(&_KawaiToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.Transfer(&_KawaiToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.TransferFrom(&_KawaiToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_KawaiToken *KawaiTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _KawaiToken.Contract.TransferFrom(&_KawaiToken.TransactOpts, from, to, value)
}

// KawaiTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the KawaiToken contract.
type KawaiTokenApprovalIterator struct {
	Event *KawaiTokenApproval // Event containing the contract specifics and raw log

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
func (it *KawaiTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KawaiTokenApproval)
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
		it.Event = new(KawaiTokenApproval)
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
func (it *KawaiTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KawaiTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KawaiTokenApproval represents a Approval event raised by the KawaiToken contract.
type KawaiTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_KawaiToken *KawaiTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*KawaiTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _KawaiToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenApprovalIterator{contract: _KawaiToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_KawaiToken *KawaiTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *KawaiTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _KawaiToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KawaiTokenApproval)
				if err := _KawaiToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_KawaiToken *KawaiTokenFilterer) ParseApproval(log types.Log) (*KawaiTokenApproval, error) {
	event := new(KawaiTokenApproval)
	if err := _KawaiToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KawaiTokenRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the KawaiToken contract.
type KawaiTokenRoleAdminChangedIterator struct {
	Event *KawaiTokenRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *KawaiTokenRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KawaiTokenRoleAdminChanged)
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
		it.Event = new(KawaiTokenRoleAdminChanged)
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
func (it *KawaiTokenRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KawaiTokenRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KawaiTokenRoleAdminChanged represents a RoleAdminChanged event raised by the KawaiToken contract.
type KawaiTokenRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_KawaiToken *KawaiTokenFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*KawaiTokenRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _KawaiToken.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenRoleAdminChangedIterator{contract: _KawaiToken.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_KawaiToken *KawaiTokenFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *KawaiTokenRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _KawaiToken.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KawaiTokenRoleAdminChanged)
				if err := _KawaiToken.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_KawaiToken *KawaiTokenFilterer) ParseRoleAdminChanged(log types.Log) (*KawaiTokenRoleAdminChanged, error) {
	event := new(KawaiTokenRoleAdminChanged)
	if err := _KawaiToken.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KawaiTokenRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the KawaiToken contract.
type KawaiTokenRoleGrantedIterator struct {
	Event *KawaiTokenRoleGranted // Event containing the contract specifics and raw log

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
func (it *KawaiTokenRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KawaiTokenRoleGranted)
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
		it.Event = new(KawaiTokenRoleGranted)
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
func (it *KawaiTokenRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KawaiTokenRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KawaiTokenRoleGranted represents a RoleGranted event raised by the KawaiToken contract.
type KawaiTokenRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_KawaiToken *KawaiTokenFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*KawaiTokenRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _KawaiToken.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenRoleGrantedIterator{contract: _KawaiToken.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_KawaiToken *KawaiTokenFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *KawaiTokenRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _KawaiToken.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KawaiTokenRoleGranted)
				if err := _KawaiToken.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_KawaiToken *KawaiTokenFilterer) ParseRoleGranted(log types.Log) (*KawaiTokenRoleGranted, error) {
	event := new(KawaiTokenRoleGranted)
	if err := _KawaiToken.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KawaiTokenRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the KawaiToken contract.
type KawaiTokenRoleRevokedIterator struct {
	Event *KawaiTokenRoleRevoked // Event containing the contract specifics and raw log

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
func (it *KawaiTokenRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KawaiTokenRoleRevoked)
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
		it.Event = new(KawaiTokenRoleRevoked)
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
func (it *KawaiTokenRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KawaiTokenRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KawaiTokenRoleRevoked represents a RoleRevoked event raised by the KawaiToken contract.
type KawaiTokenRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_KawaiToken *KawaiTokenFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*KawaiTokenRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _KawaiToken.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenRoleRevokedIterator{contract: _KawaiToken.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_KawaiToken *KawaiTokenFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *KawaiTokenRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _KawaiToken.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KawaiTokenRoleRevoked)
				if err := _KawaiToken.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_KawaiToken *KawaiTokenFilterer) ParseRoleRevoked(log types.Log) (*KawaiTokenRoleRevoked, error) {
	event := new(KawaiTokenRoleRevoked)
	if err := _KawaiToken.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KawaiTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the KawaiToken contract.
type KawaiTokenTransferIterator struct {
	Event *KawaiTokenTransfer // Event containing the contract specifics and raw log

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
func (it *KawaiTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KawaiTokenTransfer)
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
		it.Event = new(KawaiTokenTransfer)
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
func (it *KawaiTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KawaiTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KawaiTokenTransfer represents a Transfer event raised by the KawaiToken contract.
type KawaiTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_KawaiToken *KawaiTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KawaiTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KawaiToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KawaiTokenTransferIterator{contract: _KawaiToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_KawaiToken *KawaiTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *KawaiTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KawaiToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KawaiTokenTransfer)
				if err := _KawaiToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_KawaiToken *KawaiTokenFilterer) ParseTransfer(log types.Log) (*KawaiTokenTransfer, error) {
	event := new(KawaiTokenTransfer)
	if err := _KawaiToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
