# Contract Deployment Guide

**Universal deployment guide for Monad Testnet and Mainnet**

---

## 🎯 QUICK START

### Testnet Deployment
```bash
# 1. Prepare environment
export NETWORK=testnet
export RPC_URL=https://testnet-rpc.monad.xyz
export ENV_FILE=.env.testnet

# 2. Deploy contracts (see PHASE 1)
# 3. Grant permissions (see PHASE 2)
# 4. Regenerate constants (see PHASE 3)
```

### Mainnet Deployment
```bash
# 1. Prepare environment
export NETWORK=mainnet
export RPC_URL=https://rpc.monad.xyz
export ENV_FILE=.env.mainnet

# 2. Deploy contracts (see PHASE 1)
# 3. Grant permissions (see PHASE 2)
# 4. Regenerate constants (see PHASE 3)
```

---

## 📋 PREREQUISITES

### Required
- [ ] Contracts compiled: `make contracts-compile`
- [ ] Contract tests pass: `make contracts-test`
- [ ] Private key ready (in `contracts/.env`)
- [ ] Sufficient MON for gas (~1 MON for testnet, ~1 MON for mainnet)
- [ ] Admin private key exported: `export ADMIN_PRIVATE_KEY=<your_private_key>`

### Environment Files
- **Testnet**: `.env.testnet` and `contracts/.env.testnet`
- **Mainnet**: `.env.mainnet` and `contracts/.env.mainnet`

### Admin Private Key
The `ADMIN_PRIVATE_KEY` is the private key of the deployer account that will:
- Deploy all contracts
- Grant MINTER_ROLE to distributors
- Have DEFAULT_ADMIN_ROLE on KawaiToken

**How to obtain:**
- Export from your wallet (MetaMask, etc.)
- Or generate via: `cast wallet new`

**Security:**
- Store in environment variable: `export ADMIN_PRIVATE_KEY=0x...`
- Never commit to source control
- Use secrets manager for production

**Note:** This is the same key used in `contracts/.env` as `PRIVATE_KEY`

---

## 🚀 PHASE 1: Deploy Contracts

### Step 1.1: Prepare Contracts Environment

**For Testnet:**
```bash
cd contracts
cp .env.testnet .env
```

**For Mainnet:**
```bash
cd contracts
cp .env.mainnet .env
# Ensure USDC_ADDRESS is set to: 0x754704bc059f8c67012fed69bc8a327a5aafb603
```

### Step 1.2: Deploy Base Contracts

```bash
# Deploy KawaiToken, PaymentVault, OTCMarket, KAWAI_Distributor, USDT_Distributor
forge script script/DeployKawai.s.sol:DeployKawai \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify

# Save the KawaiToken address from output
export KAWAI_TOKEN_ADDRESS=<deployed_address>
```

### Step 1.3: Update contracts/.env

Add the KawaiToken address to `contracts/.env`:
```bash
echo "KAWAI_TOKEN_ADDRESS=$KAWAI_TOKEN_ADDRESS" >> contracts/.env
```

### Step 1.4: Deploy Distributors

```bash
# Deploy Mining Distributor
forge script script/DeployMiningDistributor.s.sol:DeployMiningDistributor \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify

# Deploy Cashback Distributor
forge script script/DeployCashbackDistributor.s.sol:DeployCashbackDistributor \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify

# Deploy Referral Distributor
forge script script/DeployReferralDistributor.s.sol:DeployReferralDistributor \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify
```

**⚠️ CRITICAL: Save all deployed addresses!**

---

## 🔐 PHASE 2: Grant Permissions

### Step 2.1: Grant MINTER_ROLE

```bash
# Set addresses from deployment output
export MINING_DISTRIBUTOR_ADDRESS=<deployed_address>
export CASHBACK_DISTRIBUTOR_ADDRESS=<deployed_address>
export REFERRAL_DISTRIBUTOR_ADDRESS=<deployed_address>

# Grant MINTER_ROLE to Mining Distributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $MINING_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $ADMIN_PRIVATE_KEY

# Grant MINTER_ROLE to Cashback Distributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $CASHBACK_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $ADMIN_PRIVATE_KEY

# Grant MINTER_ROLE to Referral Distributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $REFERRAL_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $ADMIN_PRIVATE_KEY
```

### Step 2.2: Verify Permissions

```bash
go run cmd/dev/check-minter-role/main.go
```

**Expected output:**
```text
✅ Mining Distributor has MINTER_ROLE
✅ Cashback Distributor has MINTER_ROLE
✅ Referral Distributor has MINTER_ROLE
```

---

## 🔄 PHASE 3: Update Configuration

### Step 3.1: Update Environment File

Update the appropriate `.env` file with all deployed addresses:

**For Testnet:**
```bash
# Edit .env.testnet
TOKEN_ADDRESS=<kawai_token_address>
OTC_MARKET_ADDRESS=<otc_market_address>
PAYMENT_VAULT_ADDRESS=<payment_vault_address>
KAWAI_DISTRIBUTOR_ADDRESS=<kawai_distributor_address>
USDT_DISTRIBUTOR_ADDRESS=<usdt_distributor_address>
MINING_DISTRIBUTOR_ADDRESS=<mining_distributor_address>
CASHBACK_DISTRIBUTOR_ADDRESS=<cashback_distributor_address>
REFERRAL_DISTRIBUTOR_ADDRESS=<referral_distributor_address>
```

**For Mainnet:**
```bash
# Edit .env.mainnet (same format as above)
```

### Step 3.2: Regenerate Backend Constants

This step generates obfuscated constants for backend services from the environment file. The obfuscator converts sensitive values (private keys, API keys) and contract addresses into Go constants.

```bash
# For Testnet
go run cmd/obfuscator-gen/main.go .env.testnet

# For Mainnet
go run cmd/obfuscator-gen/main.go .env.mainnet
```

### Step 3.3: Verify Generated Files

```bash
git diff internal/constant/blockchain.go
git diff pkg/jarvis/db/project_tokens.go
```

**What to check:**
- ✅ RPC URL matches network (testnet or mainnet)
- ✅ All contract addresses updated
- ✅ Comment shows correct deployment date

### Step 3.4: Verify Build

```bash
go build -o /dev/null .
```

---

## ✅ PHASE 4: Commit Changes

```bash
# Add generated files
git add internal/constant/blockchain.go
git add pkg/jarvis/db/project_tokens.go
git add internal/constant/cloudflare.go  # If changed

# Commit with descriptive message
git commit -m "chore: deploy contracts to $NETWORK ($(date +%Y-%m-%d))" \
  -m "- Updated blockchain.go with new $NETWORK contract addresses" \
  -m "- Updated project_tokens.go with new $NETWORK addresses" \
  -m "- All contracts deployed with ERC20Capped" \
  -m "- Generated from $ENV_FILE"

# Push to remote
git push origin <branch_name>
```

---

## 📊 DEPLOYMENT HISTORY

### Mainnet (2026-01-22)
```text
USDC:                    0x754704bc059f8c67012fed69bc8a327a5aafb603
KawaiToken:              0x5B7408a485E3167c91e925e8701d35e71B80331C ✅ Verified
PaymentVault:            0xBDC7Ad4F9e911e2EdC1128809cBC0C870EddfD9a ✅ Verified
OTCMarket:               0x9CaE910e3faC30B9E85abB3053301B3fB5a8D9fb ✅ Verified
KAWAI_Distributor:       0xDDb60C1fdbeb670c522F3f859ba1EE57A5740a14 ✅ Verified
USDT_Distributor:        0x52f71a92D4e12f87F171D91c5134042A20893650 ✅ Verified
Mining_Distributor:      0xF447C701B43e4FC4A2a172D828268Eb1F0C092FB ✅ Verified
Cashback_Distributor:    0x3Fa14A2b33f95E590bDd57a812bE4012ea5d61FF ✅ Verified
Referral_Distributor:    0xBF4c7ae729223c5E6aDb85708D685855a6d9d077 ✅ Verified
Gas Used: ~0.93 MON
```

### Testnet (2026-01-22)
```text
MockUSDT:                0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc
KawaiToken:              0xb2C745FF051e681f7B7FE8645b160eDaAC350B85 ✅ Verified
PaymentVault:            0x587c1bcAD2eD250CE05AeAF056Ade892782f9E51 ✅ Verified
OTCMarket:               0xa5A154a3719099Ad80c9c8E98f7fF50068825B0b ✅ Verified
KAWAI_Distributor:       0x78760B1e3f3F80cd4bf6EC1442d2a81797a32289 ✅ Verified
USDT_Distributor:        0xC2cA5C30De4F9a9a1e1fdB491432FEe596702E6c ✅ Verified
Mining_Distributor:      0x3373C18A0be2774004DB11625eCc099A72218Af1 ✅ Verified
Cashback_Distributor:    0xfEEC7790Dc2079Bf1922761BeA1B09FD057F7CB1 ✅ Verified
Referral_Distributor:    0xe4d6D1E6926d2616d80060d70c10d0f6Def994ea ✅ Verified
Gas Used: ~0.93 MON
```

---

## 🔍 VERIFICATION

### Verify Contract on Explorer

**Testnet:**
- Visit: <https://testnet.monadvision.com/>
- Search for contract address
- Check: Source code verified, correct compiler version (v0.8.20)

**Mainnet:**
- Visit: <https://monadvision.com/>
- Search for contract address
- Check: Source code verified, correct compiler version (v0.8.20)

### Verify Contract Functionality

```bash
# Check KawaiToken cap (should be 1B)
cast call $KAWAI_TOKEN_ADDRESS "cap()(uint256)" --rpc-url $RPC_URL

# Expected: 1000000000000000000000000000 (1B with 18 decimals)

# Check MINTER_ROLE
go run cmd/dev/check-minter-role/main.go
```

---

## 🆘 TROUBLESHOOTING

### Deployment Fails

#### Error: Insufficient funds
```bash
# Check MON balance
cast balance $DEPLOYER_ADDRESS --rpc-url $RPC_URL

# Need at least 1 MON for testnet, 10 MON for mainnet
```

#### Error: Contract already deployed

```bash
# Check if address already exists in .env
# If redeploying, remove old address first
```

### Verification Fails

#### Error: Contract not verified
```bash
# Manual verification
# For Mainnet (chain ID 143)
forge verify-contract \
  --chain-id 143 \
  --compiler-version v0.8.20 \
  $CONTRACT_ADDRESS \
  contracts/contracts/KawaiToken.sol:KawaiToken \
  --constructor-args $(cast abi-encode "constructor(address,address)" $ADMIN $MINTER) \
  --rpc-url $RPC_URL

# For Testnet (chain ID 10143)
forge verify-contract \
  --chain-id 10143 \
  --compiler-version v0.8.20 \
  $CONTRACT_ADDRESS \
  contracts/contracts/KawaiToken.sol:KawaiToken \
  --constructor-args $(cast abi-encode "constructor(address,address)" $ADMIN $MINTER) \
  --rpc-url $RPC_URL
```

**Note:** Chain IDs:
- Mainnet: `143`
- Testnet: `10143`

### Permission Grant Fails

#### Error: Transaction reverted
```bash
# Most common causes:
# 1. Insufficient MON for gas
# 2. Wrong private key (not the deployer)
# 3. RPC connection issue

# Solution: Check MON balance
cast balance $DEPLOYER_ADDRESS --rpc-url $RPC_URL

# Solution: Verify you're using deployer's private key
cast wallet address --private-key $ADMIN_PRIVATE_KEY
# Should match the address that deployed KawaiToken
```

---

## 📝 CHECKLIST

### Pre-Deployment
- [ ] Contracts compiled
- [ ] Tests passing
- [ ] Environment file ready
- [ ] Private key secured
- [ ] Sufficient MON balance

### Deployment
- [ ] Base contracts deployed
- [ ] Distributors deployed
- [ ] All addresses saved
- [ ] Contracts verified on explorer

### Post-Deployment
- [ ] MINTER_ROLE granted
- [ ] Permissions verified
- [ ] Environment file updated
- [ ] Constants regenerated
- [ ] Build verified
- [ ] Changes committed

---

**Last Updated:** January 22, 2026  
**Status:** Production Ready  
**Estimated Time:** 30-45 minutes per network
