// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import {Script, console} from "forge-std/Script.sol";
import {PaymentVault} from "../contracts/PaymentVault.sol";

contract DeployPaymentVault is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        
        // Read stablecoin address from environment
        address stablecoin = vm.envAddress("STABLECOIN_ADDRESS");

        console.log("=== Deploying PaymentVault ===");
        console.log("Deployer:", deployer);
        console.log("Stablecoin:", stablecoin);

        vm.startBroadcast(deployerPrivateKey);

        // Deploy PaymentVault (Stablecoin Deposits for credits)
        PaymentVault vault = new PaymentVault(stablecoin, deployer);
        console.log("PaymentVault deployed at:", address(vault));

        vm.stopBroadcast();

        console.log("\n=== Deployment Summary ===");
        console.log("Network:", vm.envOr("NETWORK", string("Unknown")));
        console.log("PaymentVault:", address(vault));
        console.log("==========================");
    }
}

