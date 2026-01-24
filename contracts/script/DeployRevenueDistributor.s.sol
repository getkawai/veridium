// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import {Script, console} from "forge-std/Script.sol";
import {RevenueDistributor} from "../contracts/RevenueDistributor.sol";

contract DeployRevenueDistributor is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        
        // Read stablecoin address from environment
        address stablecoin = vm.envAddress("STABLECOIN_ADDRESS");

        console.log("=== Deploying RevenueDistributor ===");
        console.log("Deployer:", deployer);
        console.log("Stablecoin:", stablecoin);

        vm.startBroadcast(deployerPrivateKey);

        // Deploy RevenueDistributor (Revenue Sharing)
        // Transfers stablecoin from pre-funded balance
        RevenueDistributor revenueDistributor = new RevenueDistributor(
            stablecoin
        );
        console.log("RevenueDistributor deployed at:", address(revenueDistributor));

        vm.stopBroadcast();

        console.log("\n=== Deployment Summary ===");
        console.log("Network:", vm.envOr("NETWORK", string("Unknown")));
        console.log("RevenueDistributor:", address(revenueDistributor));
        console.log("==========================");
    }
}
