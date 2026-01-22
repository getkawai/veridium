// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {KawaiToken} from "../contracts/KawaiToken.sol";
import {RevenueDistributor} from "../contracts/RevenueDistributor.sol";
import {MockStablecoin} from "../contracts/MockStablecoin.sol";
import {PaymentVault} from "../contracts/PaymentVault.sol";
import {OTCMarket} from "../contracts/OTCMarket.sol";

contract DeployKawai is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        console.log("=== Starting Full Suite Deployment ===");
        console.log("Deployer:", deployer);

        vm.startBroadcast(deployerPrivateKey);

        // 1. Get or Deploy Stablecoin
        // Mainnet: Use existing USDC (0x754704bc059f8c67012fed69bc8a327a5aafb603)
        // Testnet: Deploy MockStablecoin
        address stablecoinAddress;
        
        // Try to read USDC_ADDRESS from environment (for mainnet)
        try vm.envAddress("USDC_ADDRESS") returns (address existingUsdc) {
            stablecoinAddress = existingUsdc;
            console.log("Using existing stablecoin at:", stablecoinAddress);
        } catch {
            // If not set, deploy MockStablecoin (for testnet)
            MockStablecoin stablecoin = new MockStablecoin();
            stablecoinAddress = address(stablecoin);
            console.log("MockStablecoin deployed at:", stablecoinAddress);
        }

        // 2. Deploy KawaiToken
        // deployer is admin, deployer is initial minter
        KawaiToken token = new KawaiToken(deployer, deployer);
        console.log("KawaiToken deployed at:", address(token));

        // 3. Deploy RevenueDistributor (Revenue Sharing)
        // Transfers stablecoin from pre-funded balance
        RevenueDistributor revenueDistributor = new RevenueDistributor(
            stablecoinAddress
        );
        console.log(
            "RevenueDistributor (Revenue Sharing) deployed at:",
            address(revenueDistributor)
        );

        // 5. Deploy PaymentVault (Stablecoin Deposits for credits)
        PaymentVault vault = new PaymentVault(stablecoinAddress, deployer);
        console.log("PaymentVault deployed at:", address(vault));

        // 6. Deploy OTCMarket (Escrow)
        // deployer as fee recipient for now
        OTCMarket escrow = new OTCMarket(
            address(token),
            stablecoinAddress,
            deployer
        );
        console.log("OTCMarket (Escrow) deployed at:", address(escrow));

        // --- Setup & Permissions ---
        // No MINTER_ROLE needed for RevenueDistributor (transfer mode only)

        vm.stopBroadcast();

        console.log("\n=== Deployment Summary (SAVE THESE!) ===");
        console.log("Network:", vm.envOr("NETWORK", string("Unknown")));
        console.log("Stablecoin:", stablecoinAddress);
        console.log("KawaiToken:", address(token));
        console.log("RevenueDistributor:", address(revenueDistributor));
        console.log("PaymentVault:", address(vault));
        console.log("OTCMarket:", address(escrow));
        console.log("=========================================");
    }
}
