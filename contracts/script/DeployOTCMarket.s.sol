// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import {Script, console} from "forge-std/Script.sol";
import {OTCMarket} from "../contracts/OTCMarket.sol";

contract DeployOTCMarket is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        
        // Read required addresses from environment
        address kawaiToken = vm.envAddress("KAWAI_TOKEN_ADDRESS");
        address stablecoin = vm.envAddress("STABLECOIN_ADDRESS");

        console.log("=== Deploying OTCMarket ===");
        console.log("Deployer:", deployer);
        console.log("KawaiToken:", kawaiToken);
        console.log("Stablecoin:", stablecoin);

        vm.startBroadcast(deployerPrivateKey);

        // Deploy OTCMarket
        // deployer as fee recipient
        OTCMarket otcMarket = new OTCMarket(
            kawaiToken,
            stablecoin,
            deployer
        );
        console.log("OTCMarket deployed at:", address(otcMarket));

        vm.stopBroadcast();

        console.log("\n=== Deployment Summary ===");
        console.log("Network:", vm.envOr("NETWORK", string("Unknown")));
        console.log("OTCMarket:", address(otcMarket));
        console.log("==========================");
    }
}
