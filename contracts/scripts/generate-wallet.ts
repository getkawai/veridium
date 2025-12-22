import { ethers } from "ethers";
import * as fs from "fs";
import * as path from "path";

// Generate a new wallet for deployment
const wallet = ethers.Wallet.createRandom();

console.log("=== New Deployment Wallet Generated ===\n");
console.log("Address:", wallet.address);
console.log("Private Key:", wallet.privateKey);
console.log("Mnemonic:", wallet.mnemonic?.phrase);

// Create .env file
const envContent = `PRIVATE_KEY=${wallet.privateKey.slice(2)}\n`;
const envPath = path.join(__dirname, "../.env");

fs.writeFileSync(envPath, envContent);
console.log("\n✓ .env file created at:", envPath);

console.log("\n=== Next Steps ===");
console.log("1. Save your mnemonic phrase in a safe place (for backup)");
console.log("2. Get testnet BNB for this address:", wallet.address);
console.log("   Visit: https://testnet.bnbchain.org/faucet-smart");
console.log("   Enter address:", wallet.address);
console.log("3. Wait for testnet BNB to arrive (check on BSCScan Testnet)");
console.log("4. Run deployment: npx hardhat run scripts/deploy.ts --network bsctestnet");

console.log("\n⚠️  IMPORTANT: Keep your private key and mnemonic safe!");
console.log("⚠️  Never share them or commit them to git!");
