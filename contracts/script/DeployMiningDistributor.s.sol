// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "forge-std/Script.sol";
import "../contracts/MiningRewardDistributor.sol";

contract DeployMiningDistributor is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address kawaiToken = vm.envAddress("KAWAI_TOKEN_ADDRESS");
        
        vm.startBroadcast(deployerPrivateKey);
        
        MiningRewardDistributor distributor = new MiningRewardDistributor(kawaiToken);
        
        console.log("==============================================");
        console.log("MiningRewardDistributor deployed at:", address(distributor));
        console.log("==============================================");
        console.log("KAWAI Token:", kawaiToken);
        console.log("Current Period:", distributor.currentPeriod());
        console.log("==============================================");
        console.log("");
        console.log("IMPORTANT: Developer addresses specified per claim");
        console.log("Backend uses GetRandomTreasuryAddress() for distribution");
        console.log("==============================================");
        
        vm.stopBroadcast();

        console.log("\n=== Deployment Summary ===");
        console.log("Network:", vm.envOr("NETWORK", string("Unknown")));
        console.log("MiningRewardDistributor:", address(distributor));
        console.log("==========================");
    }
}

