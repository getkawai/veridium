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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MINTER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f5ffd5b50604051611e97380380611e97833981810160405281019061003191906102c9565b6040518060400160405280600e81526020017f4b6177616920414920546f6b656e0000000000000000000000000000000000008152506040518060400160405280600581526020017f4b4157414900000000000000000000000000000000000000000000000000000081525081600390816100ac9190610544565b5080600490816100bc9190610544565b5050506100d15f5f1b8361010a60201b60201c565b506101027f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a68261010a60201b60201c565b505050610613565b5f61011b838361020060201b60201c565b6101f657600160055f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555061019361026460201b60201c565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4600190506101fa565b5f90505b92915050565b5f60055f8481526020019081526020015f205f015f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f33905090565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6102988261026f565b9050919050565b6102a88161028e565b81146102b2575f5ffd5b50565b5f815190506102c38161029f565b92915050565b5f5f604083850312156102df576102de61026b565b5b5f6102ec858286016102b5565b92505060206102fd858286016102b5565b9150509250929050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061038257607f821691505b6020821081036103955761039461033e565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026103f77fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826103bc565b61040186836103bc565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61044561044061043b84610419565b610422565b610419565b9050919050565b5f819050919050565b61045e8361042b565b61047261046a8261044c565b8484546103c8565b825550505050565b5f5f905090565b61048961047a565b610494818484610455565b505050565b5b818110156104b7576104ac5f82610481565b60018101905061049a565b5050565b601f8211156104fc576104cd8161039b565b6104d6846103ad565b810160208510156104e5578190505b6104f96104f1856103ad565b830182610499565b50505b505050565b5f82821c905092915050565b5f61051c5f1984600802610501565b1980831691505092915050565b5f610534838361050d565b9150826002028217905092915050565b61054d82610307565b67ffffffffffffffff81111561056657610565610311565b5b610570825461036b565b61057b8282856104bb565b5f60209050601f8311600181146105ac575f841561059a578287015190505b6105a48582610529565b86555061060b565b601f1984166105ba8661039b565b5f5b828110156105e1578489015182556001820191506020850194506020810190506105bc565b868310156105fe57848901516105fa601f89168261050d565b8355505b6001600288020188555050505b505050505050565b611877806106205f395ff3fe608060405234801561000f575f5ffd5b5060043610610135575f3560e01c806340c10f19116100b657806395d89b411161007a57806395d89b411461035d578063a217fddf1461037b578063a9059cbb14610399578063d5391393146103c9578063d547741f146103e7578063dd62ed3e1461040357610135565b806340c10f19146102a957806342966c68146102c557806370a08231146102e157806379cc67901461031157806391d148541461032d57610135565b8063248a9ca3116100fd578063248a9ca3146102055780632f2ff15d14610235578063313ce5671461025157806332cb6b0c1461026f57806336568abe1461028d57610135565b806301ffc9a71461013957806306fdde0314610169578063095ea7b31461018757806318160ddd146101b757806323b872dd146101d5575b5f5ffd5b610153600480360381019061014e91906112db565b610433565b6040516101609190611320565b60405180910390f35b6101716104ac565b60405161017e91906113a9565b60405180910390f35b6101a1600480360381019061019c9190611456565b61053c565b6040516101ae9190611320565b60405180910390f35b6101bf61055e565b6040516101cc91906114a3565b60405180910390f35b6101ef60048036038101906101ea91906114bc565b610567565b6040516101fc9190611320565b60405180910390f35b61021f600480360381019061021a919061153f565b610595565b60405161022c9190611579565b60405180910390f35b61024f600480360381019061024a9190611592565b6105b2565b005b6102596105d4565b60405161026691906115eb565b60405180910390f35b6102776105dc565b60405161028491906114a3565b60405180910390f35b6102a760048036038101906102a29190611592565b6105ec565b005b6102c360048036038101906102be9190611456565b610667565b005b6102df60048036038101906102da9190611604565b610701565b005b6102fb60048036038101906102f6919061162f565b610715565b60405161030891906114a3565b60405180910390f35b61032b60048036038101906103269190611456565b61075a565b005b61034760048036038101906103429190611592565b61077a565b6040516103549190611320565b60405180910390f35b6103656107de565b60405161037291906113a9565b60405180910390f35b61038361086e565b6040516103909190611579565b60405180910390f35b6103b360048036038101906103ae9190611456565b610874565b6040516103c09190611320565b60405180910390f35b6103d1610896565b6040516103de9190611579565b60405180910390f35b61040160048036038101906103fc9190611592565b6108ba565b005b61041d6004803603810190610418919061165a565b6108dc565b60405161042a91906114a3565b60405180910390f35b5f7f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806104a557506104a48261095e565b5b9050919050565b6060600380546104bb906116c5565b80601f01602080910402602001604051908101604052809291908181526020018280546104e7906116c5565b80156105325780601f1061050957610100808354040283529160200191610532565b820191905f5260205f20905b81548152906001019060200180831161051557829003601f168201915b5050505050905090565b5f5f6105466109c7565b90506105538185856109ce565b600191505092915050565b5f600254905090565b5f5f6105716109c7565b905061057e8582856109e0565b610589858585610a73565b60019150509392505050565b5f60055f8381526020019081526020015f20600101549050919050565b6105bb82610595565b6105c481610b63565b6105ce8383610b77565b50505050565b5f6012905090565b6b033b2e3c9fd0803ce800000081565b6105f46109c7565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610658576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106628282610c61565b505050565b7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a661069181610b63565b6b033b2e3c9fd0803ce8000000826106a761055e565b6106b19190611722565b11156106f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106e99061179f565b60405180910390fd5b6106fc8383610d4b565b505050565b61071261070c6109c7565b82610dca565b50565b5f5f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b61076c826107666109c7565b836109e0565b6107768282610dca565b5050565b5f60055f8481526020019081526020015f205f015f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b6060600480546107ed906116c5565b80601f0160208091040260200160405190810160405280929190818152602001828054610819906116c5565b80156108645780601f1061083b57610100808354040283529160200191610864565b820191905f5260205f20905b81548152906001019060200180831161084757829003601f168201915b5050505050905090565b5f5f1b81565b5f5f61087e6109c7565b905061088b818585610a73565b600191505092915050565b7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a681565b6108c382610595565b6108cc81610b63565b6108d68383610c61565b50505050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f33905090565b6109db8383836001610e49565b505050565b5f6109eb84846108dc565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811015610a6d5781811015610a5e578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401610a55939291906117cc565b60405180910390fd5b610a6c84848484035f610e49565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610ae3575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610ada9190611801565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610b53575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610b4a9190611801565b60405180910390fd5b610b5e838383611018565b505050565b610b7481610b6f6109c7565b611231565b50565b5f610b82838361077a565b610c5757600160055f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550610bf46109c7565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019050610c5b565b5f90505b92915050565b5f610c6c838361077a565b15610d41575f60055f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550610cde6109c7565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a460019050610d45565b5f90505b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610dbb575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610db29190611801565b60405180910390fd5b610dc65f8383611018565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610e3a575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610e319190611801565b60405180910390fd5b610e45825f83611018565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610eb9575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401610eb09190611801565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610f29575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610f209190611801565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015611012578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161100991906114a3565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603611068578060025f82825461105c9190611722565b92505081905550611136565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050818110156110f1578381836040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016110e8939291906117cc565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361117d578060025f82825403925050819055506111c7565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161122491906114a3565b60405180910390a3505050565b61123b828261077a565b61127e5780826040517fe2517d3f00000000000000000000000000000000000000000000000000000000815260040161127592919061181a565b60405180910390fd5b5050565b5f5ffd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b6112ba81611286565b81146112c4575f5ffd5b50565b5f813590506112d5816112b1565b92915050565b5f602082840312156112f0576112ef611282565b5b5f6112fd848285016112c7565b91505092915050565b5f8115159050919050565b61131a81611306565b82525050565b5f6020820190506113335f830184611311565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61137b82611339565b6113858185611343565b9350611395818560208601611353565b61139e81611361565b840191505092915050565b5f6020820190508181035f8301526113c18184611371565b905092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6113f2826113c9565b9050919050565b611402816113e8565b811461140c575f5ffd5b50565b5f8135905061141d816113f9565b92915050565b5f819050919050565b61143581611423565b811461143f575f5ffd5b50565b5f813590506114508161142c565b92915050565b5f5f6040838503121561146c5761146b611282565b5b5f6114798582860161140f565b925050602061148a85828601611442565b9150509250929050565b61149d81611423565b82525050565b5f6020820190506114b65f830184611494565b92915050565b5f5f5f606084860312156114d3576114d2611282565b5b5f6114e08682870161140f565b93505060206114f18682870161140f565b925050604061150286828701611442565b9150509250925092565b5f819050919050565b61151e8161150c565b8114611528575f5ffd5b50565b5f8135905061153981611515565b92915050565b5f6020828403121561155457611553611282565b5b5f6115618482850161152b565b91505092915050565b6115738161150c565b82525050565b5f60208201905061158c5f83018461156a565b92915050565b5f5f604083850312156115a8576115a7611282565b5b5f6115b58582860161152b565b92505060206115c68582860161140f565b9150509250929050565b5f60ff82169050919050565b6115e5816115d0565b82525050565b5f6020820190506115fe5f8301846115dc565b92915050565b5f6020828403121561161957611618611282565b5b5f61162684828501611442565b91505092915050565b5f6020828403121561164457611643611282565b5b5f6116518482850161140f565b91505092915050565b5f5f604083850312156116705761166f611282565b5b5f61167d8582860161140f565b925050602061168e8582860161140f565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806116dc57607f821691505b6020821081036116ef576116ee611698565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61172c82611423565b915061173783611423565b925082820190508082111561174f5761174e6116f5565b5b92915050565b7f4d617820737570706c79206578636565646564000000000000000000000000005f82015250565b5f611789601383611343565b915061179482611755565b602082019050919050565b5f6020820190508181035f8301526117b68161177d565b9050919050565b6117c6816113e8565b82525050565b5f6060820190506117df5f8301866117bd565b6117ec6020830185611494565b6117f96040830184611494565b949350505050565b5f6020820190506118145f8301846117bd565b92915050565b5f60408201905061182d5f8301856117bd565b61183a602083018461156a565b939250505056fea26469706673582212203ed23d0d1b7cc16822b496099434668a59ece674443a4bb9d8d482b6a9bb103764736f6c634300081e0033",
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

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_KawaiToken *KawaiTokenCaller) MAXSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KawaiToken.contract.Call(opts, &out, "MAX_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_KawaiToken *KawaiTokenSession) MAXSUPPLY() (*big.Int, error) {
	return _KawaiToken.Contract.MAXSUPPLY(&_KawaiToken.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_KawaiToken *KawaiTokenCallerSession) MAXSUPPLY() (*big.Int, error) {
	return _KawaiToken.Contract.MAXSUPPLY(&_KawaiToken.CallOpts)
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
