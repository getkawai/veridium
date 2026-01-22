# 🚀 Smart Contracts Guide

**Status:** Current Implementation (Jan 2026)  
**Network:** Monad Blockchain (Testnet)

---

## 📋 Current Contracts (8 Total)

### Deployed Contracts (Monad Testnet)

| Contract | Address | Purpose | Status |
|----------|---------|---------|--------|
| **KawaiToken** | `0xE32660b39D99988Df4bFdc7e4b68A4DC9D654722` | ERC20 utility token | ✅ Active |
| **MiningRewardDistributor** | `0x8117D77A219EeF5F7869897C3F0973Afb87d8427` | Mining rewards (85/5/5/5) | ✅ Active |
| **DepositCashbackDistributor** | `0xdE64f6F5bEe28762c91C76ff762365D553204e35` | Tiered cashback (1-5%) | ✅ Active |
| **USDT_Distributor** | `0x98a7590406a08Cc64dc074D8698B71e4D997a268` | USDT dividends | ✅ Active |
| **Escrow (OTCMarket)** | `0xd065F9DDb66aa90a1FF62c10868BeF921be2E103` | P2P OTC Market | ✅ Active |
| **PaymentVault** | `0x9a5A9e31977cB86cD502DC9E0B568d8F17977dAd` | USDT deposits | ✅ Active |
| **MockUSDT** | `0x2cBe796033377352158df11Ab388010ab3097F58` | Testnet USDT | ✅ Testnet Only |

**Network Info:**
- **Chain ID:** 10143
- **RPC:** `https://testnet-rpc.monad.xyz`
- **Explorer:** `https://testnet.monad.xyz`

---

## 🔧 Development Commands

### Prerequisites
```bash
# Install Foundry
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Setup environment
cd contracts
forge install
cp .env.example .env
# Edit .env with your values
```

### Essential Commands
```bash
# Compile contracts
make contracts-compile

# Run tests
make contracts-test

# Generate Go bindings
make contracts-bindings

# Full upgrade workflow
make contracts-upgrade

# Deploy to testnet
make contracts-deploy-testnet
```

### Development Workflow

#### Adding New Features
```bash
# 1. Edit contract
vim contracts/contracts/Escrow.sol

# 2. Run full upgrade workflow
make contracts-upgrade

# 3. Test locally
make dev-hot

# 4. Deploy when ready
make contracts-deploy-testnet
```

#### Bug Fixes
```bash
# 1. Write failing test
vim contracts/test/Escrow.t.sol
make contracts-test  # Should fail

# 2. Fix bug
vim contracts/contracts/Escrow.sol
make contracts-test  # Should pass

# 3. Deploy fix
make contracts-upgrade
make contracts-deploy-testnet
```

---

## 🔐 MINTER_ROLE Setup

### Why MINTER_ROLE is Required

**ALL reward distribution contracts need `MINTER_ROLE`** to mint KAWAI on-demand when users claim.

**Contracts Requiring MINTER_ROLE:**
- ✅ MiningRewardDistributor (1-4 mint calls per claim)
- ✅ DepositCashbackDistributor (1 mint call per claim)

### Grant MINTER_ROLE (Automated)
```bash
# Grant to all contracts at once
export PRIVATE_KEY=0x...
./GRANT_ALL_MINTER_ROLES.sh
```

### Verify MINTER_ROLE Status
```bash
make check-minter-role

# Expected output:
# ✅ MiningRewardDistributor: HAS MINTER_ROLE
# ✅ DepositCashbackDistributor: HAS MINTER_ROLE
```

---

## 🧪 Testing & Validation

### Run Tests
```bash
# All tests
make contracts-test

# With gas report
make contracts-test-gas

# Coverage report
make contracts-coverage
```

### Pre-Deployment Validation
```bash
# Validate everything
make contracts-validate

# Gas optimization
make contracts-gas-snapshot
make contracts-gas-compare
```

---

## 🎯 Reward Settlement

### Generate Settlements
```bash
# Mining rewards
make settle-mining

# Cashback rewards
make settle-cashback

# Referral rewards
make settle-referral

# All at once
make settle-all
```

### Upload Merkle Roots
```bash
# Upload mining settlement
make upload-merkle-root TYPE=mining ROOT=0x...

# Upload cashback settlement
make upload-merkle-root TYPE=cashback ROOT=0x...

# Upload referral settlement
make upload-merkle-root TYPE=referral ROOT=0x...
```

### Check Status
```bash
# Settlement status
make reward-settlement-status

# Check user balance
make check-balance ADDR=0x...

# Check claim status
make check-claim-status TYPE=mining PERIOD=123 ADDR=0x...
```

---

## 🐛 Troubleshooting

### Common Issues

#### "Artifact not found"
```bash
make contracts-clean
make contracts-compile
```

#### "Bindings out of sync"
```bash
make contracts-bindings
```

#### "Claim transaction fails"
```bash
# Check MINTER_ROLE
make check-minter-role

# Check Merkle root uploaded
cast call $MINING_DISTRIBUTOR "merkleRoot()" --rpc-url $RPC_URL

# Check user hasn't claimed already
make check-claim-status TYPE=mining PERIOD=<ID> ADDR=<USER>
```

#### "Deployment fails"
```bash
# Check RPC connection
cast block-number --rpc-url $RPC_URL

# Check balance
cast balance $YOUR_ADDRESS --rpc-url $RPC_URL
```

---

## 📚 Command Reference

| Command | Description | When to Use |
|---------|-------------|-------------|
| `make contracts-compile` | Compile contracts | After editing .sol files |
| `make contracts-test` | Run tests | Before deploying |
| `make contracts-bindings` | Generate Go bindings | After compile |
| `make contracts-upgrade` | Full workflow | Complete update |
| `make contracts-deploy-testnet` | Deploy to Monad | Production deploy |
| `make check-minter-role` | Check MINTER_ROLE | After deployment |
| `make settle-mining` | Mining settlement | Weekly |
| `make settle-cashback` | Cashback settlement | Weekly |
| `make settle-referral` | Referral settlement | Weekly |
| `make settle-all` | All settlements | Weekly automation |

---

## 🔗 Related Documentation

- **Main README:** Project overview and setup
- **Makefile:** All available commands
- **GRANT_ALL_MINTER_ROLES.sh:** Automated role setup
- **docs/development/:** Detailed technical guides
- **REWARD_SYSTEMS.md:** Tokenomics and reward mechanics

---

**Happy Coding! 🚀**