import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();
  const targetAddress = process.env.TARGET_ADDRESS || deployer.address;

  console.log("Checking balance for:", targetAddress);

  const balance = await ethers.provider.getBalance(targetAddress);
  console.log("Balance:", ethers.formatEther(balance), "tBNB");

  if (balance === 0n) {
    console.log("\n⚠️  Balance is 0!");
  } else {
    console.log("\n✓ Wallet has sufficient balance!");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
