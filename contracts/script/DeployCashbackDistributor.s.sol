// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "forge-std/Script.sol";
import "../contracts/DepositCashbackDistributor.sol";

contract DeployCashbackDistributor is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address kawaiTokenAddress = vm.envAddress("KAWAI_TOKEN_ADDRESS");
        
        vm.startBroadcast(deployerPrivateKey);
        
        DepositCashbackDistributor distributor = new DepositCashbackDistributor(kawaiTokenAddress);
        
        console.log("DepositCashbackDistributor deployed at:", address(distributor));
        
        vm.stopBroadcast();
    }
}

