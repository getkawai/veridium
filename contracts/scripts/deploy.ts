import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("Deploying contracts with account:", deployer.address);
  console.log("Account balance:", (await ethers.provider.getBalance(deployer.address)).toString());

  const deployments: Record<string, string> = {};

  // 1. Deploy MockUSDT
  console.log("\n1. Deploying MockUSDT...");
  const MockUSDT = await ethers.getContractFactory("MockUSDT");
  const usdt = await MockUSDT.deploy();
  await usdt.waitForDeployment();
  const usdtAddress = await usdt.getAddress();
  deployments["MockUSDT"] = usdtAddress;
  console.log("✓ MockUSDT deployed to:", usdtAddress);

  // 2. Deploy KawaiToken
  console.log("\n2. Deploying KawaiToken...");
  const KawaiToken = await ethers.getContractFactory("KawaiToken");
  const kawaiToken = await KawaiToken.deploy(
    deployer.address, // defaultAdmin
    deployer.address  // minter
  );
  await kawaiToken.waitForDeployment();
  const kawaiTokenAddress = await kawaiToken.getAddress();
  deployments["KawaiToken"] = kawaiTokenAddress;
  console.log("✓ KawaiToken deployed to:", kawaiTokenAddress);

  // Check initial balance
  const balance = await kawaiToken.balanceOf(deployer.address);
  console.log("  Initial balance:", ethers.formatEther(balance), "KAWAI");

  // 3. Deploy Escrow (OTCMarket)
  console.log("\n3. Deploying Escrow (OTCMarket)...");
  const Escrow = await ethers.getContractFactory("OTCMarket");
  const escrow = await Escrow.deploy(
    kawaiTokenAddress,  // _tokenDeAI
    usdtAddress,        // _usdt
    deployer.address    // _feeRecipient
  );
  await escrow.waitForDeployment();
  const escrowAddress = await escrow.getAddress();
  deployments["Escrow"] = escrowAddress;
  console.log("✓ Escrow deployed to:", escrowAddress);

  // 4. Deploy PaymentVault
  console.log("\n4. Deploying PaymentVault...");
  const PaymentVault = await ethers.getContractFactory("PaymentVault");
  const vault = await PaymentVault.deploy(
    usdtAddress,        // _usdt
    deployer.address    // initialOwner
  );
  await vault.waitForDeployment();
  const vaultAddress = await vault.getAddress();
  deployments["PaymentVault"] = vaultAddress;
  console.log("✓ PaymentVault deployed to:", vaultAddress);

  // Save deployments to JSON
  const deploymentsPath = path.join(__dirname, "../deployments.json");
  fs.writeFileSync(deploymentsPath, JSON.stringify(deployments, null, 2));
  console.log("\n✓ Deployment addresses saved to:", deploymentsPath);

  // Summary
  console.log("\n=== Deployment Summary ===");
  console.log("Network:", (await ethers.provider.getNetwork()).name);
  console.log("Deployer:", deployer.address);
  console.log("\nContract Addresses:");
  Object.entries(deployments).forEach(([name, address]) => {
    console.log(`  ${name}: ${address}`);
  });

  console.log("\n=== Next Steps ===");
  console.log("1. Verify contracts on BSCScan:");
  console.log(`   npx hardhat verify --network bsctestnet ${usdtAddress}`);
  console.log(`   npx hardhat verify --network bsctestnet ${kawaiTokenAddress} "${deployer.address}" "${deployer.address}"`);
  console.log(`   npx hardhat verify --network bsctestnet ${escrowAddress} "${kawaiTokenAddress}" "${usdtAddress}" "${deployer.address}"`);
  console.log(`   npx hardhat verify --network bsctestnet ${vaultAddress} "${usdtAddress}" "${deployer.address}"`);
  console.log("\n2. Update pkg/jarvis address database with deployed addresses");
  console.log("3. Run 'make contracts-bindings' to regenerate Go bindings with new addresses");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
