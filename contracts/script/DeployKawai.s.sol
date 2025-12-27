// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {KawaiToken} from "../contracts/KawaiToken.sol";
import {MerkleDistributor} from "../contracts/MerkleDistributor.sol";
import {MockUSDT} from "../contracts/MockUSDT.sol";
import {PaymentVault} from "../contracts/PaymentVault.sol";
import {OTCMarket} from "../contracts/Escrow.sol";

contract DeployKawai is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        console.log("=== Starting Full Suite Deployment ===");
        console.log("Deployer:", deployer);

        vm.startBroadcast(deployerPrivateKey);

        // 1. Deploy MockUSDT (for testing purposes)
        MockUSDT usdt = new MockUSDT();
        console.log("MockUSDT deployed at:", address(usdt));

        // 2. Deploy KawaiToken
        // deployer is admin, deployer is initial minter
        KawaiToken token = new KawaiToken(deployer, deployer);
        console.log("KawaiToken deployed at:", address(token));

        // 3. Deploy KAWAI MerkleDistributor (Mining Rewards)
        // mintOnClaim=true: mints new KAWAI tokens when contributor claims (contributor pays gas)
        MerkleDistributor kawaiDistributor = new MerkleDistributor(
            address(token),
            true // mintOnClaim: mint tokens on claim
        );
        console.log(
            "KAWAI MerkleDistributor (Mining) deployed at:",
            address(kawaiDistributor)
        );

        // 4. Deploy USDT MerkleDistributor (Profit Sharing)
        // mintOnClaim=false: transfers USDT from pre-funded balance
        MerkleDistributor usdtDistributor = new MerkleDistributor(
            address(usdt),
            false // mintOnClaim: transfer from balance
        );
        console.log(
            "USDT MerkleDistributor (Dividends) deployed at:",
            address(usdtDistributor)
        );

        // 5. Deploy PaymentVault (USDT Deposits for credits)
        PaymentVault vault = new PaymentVault(address(usdt), deployer);
        console.log("PaymentVault deployed at:", address(vault));

        // 6. Deploy OTCMarket (Escrow)
        // deployer as fee recipient for now
        OTCMarket escrow = new OTCMarket(
            address(token),
            address(usdt),
            deployer
        );
        console.log("OTCMarket (Escrow) deployed at:", address(escrow));

        // --- Setup & Permissions ---

        // Grant MINTER_ROLE to KAWAI MerkleDistributor so it can mint mining rewards
        bytes32 MINTER_ROLE = keccak256("MINTER_ROLE");
        token.grantRole(MINTER_ROLE, address(kawaiDistributor));
        console.log(
            "Permission: Granted MINTER_ROLE to KAWAI MerkleDistributor"
        );

        vm.stopBroadcast();

        console.log("\n=== Deployment Summary (SAVE THESE!) ===");
        console.log("Network: Monad Testnet");
        console.log("MockUSDT:", address(usdt));
        console.log("KawaiToken:", address(token));
        console.log("KAWAI_Distributor:", address(kawaiDistributor));
        console.log("USDT_Distributor:", address(usdtDistributor));
        console.log("PaymentVault:", address(vault));
        console.log("OTCMarket:", address(escrow));
        console.log("=========================================");
    }
}
