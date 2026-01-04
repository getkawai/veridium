// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

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
        console.log("");
        console.log("Next steps:");
        console.log("1. Grant MINTER_ROLE to distributor:");
        console.log("   make contracts-grant-minter-mining MINING_DISTRIBUTOR_ADDRESS=%s", address(distributor));
        console.log("");
        console.log("2. Set initial Merkle root (after first week):");
        console.log("   cast send %s 'setMerkleRoot(bytes32)' <merkleRoot> --rpc-url $RPC_URL --private-key $PRIVATE_KEY", address(distributor));
        console.log("==============================================");
        
        vm.stopBroadcast();
    }
}

