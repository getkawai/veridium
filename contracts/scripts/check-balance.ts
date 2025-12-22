import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("Checking balance for:", deployer.address);

  const balance = await ethers.provider.getBalance(deployer.address);
  console.log("Balance:", ethers.formatEther(balance), "BNB");

  if (balance === 0n) {
    console.log("\n⚠️  Balance is 0!");
    console.log("Please wait for testnet BNB to arrive from the faucet.");
    console.log("\nCheck transaction status:");
    console.log("https://testnet.bscscan.com/tx/0x9826c4927f5892e535fa142c170bc051593c2264cd8e1cc335b7a480d4a1e533");
    console.log("\nCheck wallet balance:");
    console.log(`https://testnet.bscscan.com/address/${deployer.address}`);
  } else {
    console.log("\n✓ Wallet has sufficient balance for deployment!");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
