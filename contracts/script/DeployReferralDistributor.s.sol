// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "forge-std/Script.sol";
import "../contracts/ReferralRewardDistributor.sol";

contract DeployReferralDistributor is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address kawaiToken = vm.envAddress("KAWAI_TOKEN_ADDRESS");
        
        vm.startBroadcast(deployerPrivateKey);
        
        ReferralRewardDistributor distributor = new ReferralRewardDistributor(
            kawaiToken
        );
        
        console.log("ReferralRewardDistributor deployed at:", address(distributor));
        console.log("KAWAI Token:", kawaiToken);
        console.log("Note: USDT rewards are off-chain credits only");
        
        vm.stopBroadcast();
    }
}

