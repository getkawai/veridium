// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import {Script, console} from "forge-std/Script.sol";
import {KawaiToken} from "../contracts/KawaiToken.sol";

contract DeployKawai is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        console.log("=== Deploying KawaiToken ===");
        console.log("Deployer:", deployer);

        vm.startBroadcast(deployerPrivateKey);

        // Deploy KawaiToken
        // deployer is admin, deployer is initial minter
        KawaiToken token = new KawaiToken(deployer, deployer);
        console.log("KawaiToken deployed at:", address(token));

        vm.stopBroadcast();

        console.log("\n=== Deployment Summary ===");
        console.log("Network:", vm.envOr("NETWORK", string("Unknown")));
        console.log("KawaiToken:", address(token));
        console.log("==========================");
    }
}
