# 🔧 Smart Contracts Development Workflow

## 📋 Table of Contents
- [Quick Start](#quick-start)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)

---

## 🚀 Quick Start

### First Time Setup
```bash
# Install Foundry (if not installed)
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Install dependencies
cd contracts
forge install

# Setup environment variables
cp .env.example .env
# Edit .env with your actual values:
# - PRIVATE_KEY (without 0x prefix)
# - RPC_URL (default: https://testnet-rpc.monad.xyz)
# - CONTRACT_ADDRESS (after deployment)
# - ETHERSCAN_API_KEY (for verification)
```

### Basic Commands
```bash
# Compile contracts
make contracts-compile

# Run tests
make contracts-test

# Generate Go bindings
make contracts-bindings

# Full update (compile + bindings)
make contracts-update
```

---

## 🔄 Development Workflow

### Scenario 1: Adding New Features to Existing Contract

```bash
# 1. Edit contract
vim contracts/contracts/Escrow.sol

# 2. Run full upgrade workflow (test + compile + bindings)
make contracts-upgrade

# 3. Review generated bindings
ls -la internal/generate/abi/escrow/

# 4. Update backend code if needed
vim internal/services/marketplace_service.go

# 5. Test locally
make dev-hot

# 6. Deploy to testnet when ready
PRIVATE_KEY=0x... RPC_URL=https://... make contracts-deploy-testnet
```

### Scenario 2: Creating New Contract

```bash
# 1. Create new contract
vim contracts/contracts/MyNewContract.sol

# 2. Add to deployment script
vim contracts/script/DeployKawai.s.sol

# 3. Update Makefile to generate bindings
vim Makefile
# Add new abi-mynewcontract target

# 4. Run upgrade workflow
make contracts-upgrade

# 5. Deploy
make contracts-deploy-testnet
```

### Scenario 3: Fixing Bugs

```bash
# 1. Write test that reproduces bug
vim contracts/test/Escrow.t.sol

# 2. Verify test fails
make contracts-test

# 3. Fix the bug
vim contracts/contracts/Escrow.sol

# 4. Verify test passes
make contracts-test

# 5. Full upgrade
make contracts-upgrade
```

---

## 🧪 Testing

### Run All Tests
```bash
make contracts-test
```

### Run Tests with Gas Report
```bash
make contracts-test-gas
```

### Run Specific Test
```bash
cd contracts
forge test --match-test testBuyOrder -vvv
```

### Run Tests with Coverage
```bash
make contracts-coverage
```

### Create Gas Baseline
```bash
# Create baseline
make contracts-gas-snapshot

# Make changes to contract
vim contracts/contracts/Escrow.sol

# Compare gas usage
make contracts-gas-compare
```

---

## 🚀 Deployment

### Local Deployment (Anvil)

```bash
# Terminal 1: Start Anvil
anvil

# Terminal 2: Deploy
make contracts-deploy-local
```

### Testnet Deployment (Monad)

```bash
# Set environment variables
export PRIVATE_KEY="0x..."
export RPC_URL="https://testnet.monad.xyz"

# Validate before deploying
make contracts-validate

# Deploy
make contracts-deploy-testnet

# Save contract addresses from output!
```

### Verify on Block Explorer

```bash
export CONTRACT_ADDRESS="0x..."
export ETHERSCAN_API_KEY="..."

make contracts-verify
```

---

## 🔍 Validation & Quality Checks

### Before Deploying

```bash
# 1. Validate everything
make contracts-validate

# 2. Check gas usage
make contracts-test-gas

# 3. Review coverage
make contracts-coverage

# 4. Compare gas vs baseline
make contracts-gas-compare
```

---

## 🛠️ Troubleshooting

### Issue: "Artifact not found"

```bash
# Clean and rebuild
make contracts-clean
make contracts-compile
```

### Issue: "Bindings out of sync"

```bash
# Regenerate bindings
make contracts-bindings
```

### Issue: "Test fails after contract change"

```bash
# Run tests with verbose output
cd contracts
forge test -vvvv

# Check specific test
forge test --match-test testYourTest -vvvv
```

### Issue: "Deployment fails"

```bash
# Check RPC connection
cast block-number --rpc-url $RPC_URL

# Check balance
cast balance $YOUR_ADDRESS --rpc-url $RPC_URL

# Check gas price
cast gas-price --rpc-url $RPC_URL
```

### Issue: "Go bindings not working"

```bash
# Check if abigen is installed
which abigen

# Install if missing
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# Regenerate bindings
make contracts-clean
make contracts-update
```

---

## 📝 Best Practices

### 1. Always Test Before Deploying
```bash
make contracts-validate
```

### 2. Use Gas Snapshots
```bash
# Before making changes
make contracts-gas-snapshot

# After making changes
make contracts-gas-compare
```

### 3. Keep Bindings in Sync
```bash
# After every contract change
make contracts-upgrade
```

### 4. Document Contract Addresses
```bash
# Save deployment output to file
make contracts-deploy-testnet | tee deployment.log
```

### 5. Version Control
```bash
# Commit contract changes with bindings
git add contracts/ internal/generate/abi/
git commit -m "feat(contracts): add partial fill support"
```

---

## 🎯 Common Workflows

### Full Contract Update Workflow
```bash
# 1. Edit contract
vim contracts/contracts/Escrow.sol

# 2. Write/update tests
vim contracts/test/Escrow.t.sol

# 3. Run upgrade workflow
make contracts-upgrade

# 4. Update backend
vim internal/services/marketplace_service.go

# 5. Test locally
make dev-hot

# 6. Validate
make contracts-validate

# 7. Deploy
make contracts-deploy-testnet

# 8. Commit
git add .
git commit -m "feat(contracts): your feature"
```

### Emergency Bug Fix Workflow
```bash
# 1. Write failing test
vim contracts/test/Escrow.t.sol
make contracts-test  # Should fail

# 2. Fix bug
vim contracts/contracts/Escrow.sol
make contracts-test  # Should pass

# 3. Quick deploy
make contracts-upgrade
make contracts-validate
make contracts-deploy-testnet

# 4. Update app
make dev-hot
```

---

## 📚 Additional Resources

- [Foundry Book](https://book.getfoundry.sh/)
- [Solidity Docs](https://docs.soliditylang.org/)
- [OpenZeppelin Contracts](https://docs.openzeppelin.com/contracts/)
- [Monad Testnet](https://docs.monad.xyz/)

---

## 🆘 Getting Help

If you encounter issues:

1. Check this documentation
2. Run `make help` for available commands
3. Check Foundry logs in `contracts/`
4. Review test output with `-vvvv` flag
5. Ask in team chat

---

## 📊 Makefile Command Reference

| Command | Description | When to Use |
|---------|-------------|-------------|
| `make contracts-compile` | Compile contracts | After editing .sol files |
| `make contracts-test` | Run tests | Before deploying |
| `make contracts-bindings` | Generate Go bindings | After compile |
| `make contracts-update` | Compile + bindings | Quick update |
| `make contracts-upgrade` | Test + compile + bindings | Full workflow |
| `make contracts-validate` | Validate everything | Before deploy |
| `make contracts-deploy-local` | Deploy to Anvil | Local testing |
| `make contracts-deploy-testnet` | Deploy to Monad | Testnet deploy |
| `make contracts-verify` | Verify on explorer | After deploy |
| `make contracts-clean` | Clean artifacts | Fresh start |
| `make contracts-test-gas` | Gas report | Optimize gas |
| `make contracts-coverage` | Coverage report | Check tests |
| `make contracts-gas-snapshot` | Create baseline | Before changes |
| `make contracts-gas-compare` | Compare gas | After changes |

---

**Happy Coding! 🚀**

