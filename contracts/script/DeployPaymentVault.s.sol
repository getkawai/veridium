// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {PaymentVault} from "../contracts/PaymentVault.sol";

/**
 * @title DeployPaymentVault
 * @notice Deploys PaymentVault contract for mainnet or testnet
 * @dev Usage:
 *   Testnet: forge script script/DeployPaymentVault.s.sol:DeployPaymentVault --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast
 *   Mainnet: Set USDC_ADDRESS=0x754704bc059f8c67012fed69bc8a327a5aafb603 in .env.mainnet
 */
contract DeployPaymentVault is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        // Get stablecoin address from environment
        // Testnet: MockUSDT address
        // Mainnet: USDC address (0x754704bc059f8c67012fed69bc8a327a5aafb603)
        address stablecoinAddress = vm.envAddress("USDC_ADDRESS");

        console.log("==============================================");
        console.log("Deploying PaymentVault...");
        console.log("==============================================");
        console.log("Deployer:", deployer);
        console.log("Stablecoin:", stablecoinAddress);
        console.log("==============================================");

        vm.startBroadcast(deployerPrivateKey);

        PaymentVault vault = new PaymentVault(stablecoinAddress, deployer);

        vm.stopBroadcast();

        console.log("");
        console.log("==============================================");
        console.log("PaymentVault deployed at:", address(vault));
        console.log("==============================================");
        console.log("Stablecoin:", stablecoinAddress);
        console.log("Owner:", deployer);
        console.log("==============================================");
        console.log("");
        console.log("Next steps:");
        console.log("1. Update PAYMENT_VAULT_ADDRESS in .env:");
        console.log("   PAYMENT_VAULT_ADDRESS=%s", address(vault));
        console.log("");
        console.log("2. Regenerate backend constants:");
        console.log("   go run cmd/obfuscator-gen/main.go");
        console.log("");
        console.log("3. Test deposit flow:");
        console.log("   - Approve stablecoin: cast send %s 'approve(address,uint256)' %s <amount> --rpc-url $RPC_URL --private-key $USER_KEY", stablecoinAddress, address(vault));
        console.log("   - Deposit: cast send %s 'deposit(uint256)' <amount> --rpc-url $RPC_URL --private-key $USER_KEY", address(vault));
        console.log("==============================================");
    }
}
