# Environment Switching Guide

Quick guide for switching between testnet and mainnet environments.

## 📁 Environment Files

### Backend (.env)
- `.env.testnet` - Testnet configuration (Round 6)
- `.env` - **Currently: Mainnet** (empty, ready for deployment)
- `.env.example` - Template for new environments

### Contracts (contracts/.env)
- `contracts/.env.testnet` - Testnet deployment addresses
- `contracts/.env` - **Currently: Mainnet** (empty, ready for deployment)

## 🔄 Switching Environments

### Switch to Testnet

```bash
# Backend
cp .env.testnet .env

# Contracts
cp contracts/.env.testnet contracts/.env

# Regenerate constants
go run cmd/obfuscator-gen/main.go

# Verify
make pause-status
```

### Switch to Mainnet

```bash
# Backend
# .env is already mainnet (current state)

# Contracts
# contracts/.env is already mainnet (current state)

# After deployment, regenerate constants
go run cmd/obfuscator-gen/main.go
```

## ⚠️ Important Notes

### Before Switching

1. **Commit your changes** - Don't lose work
2. **Check current environment** - Know what you're switching from
3. **Backup if needed** - Save important data

### After Switching

1. **Regenerate constants** - Run `obfuscator-gen`
2. **Restart services** - Backend, frontend, etc.
3. **Verify connection** - Check RPC and contract addresses
4. **Test basic operations** - Ensure everything works

## 🚨 Critical Warnings

### Mainnet

- ❌ **NEVER commit private keys**
- ❌ **NEVER use testnet keys on mainnet**
- ✅ **Always use hardware wallet for mainnet**
- ✅ **Double-check addresses before transactions**
- ✅ **Test on testnet first**

### Testnet

- ✅ Safe to experiment
- ✅ Free tokens from faucet
- ✅ Can reset anytime
- ⚠️ Data won't transfer to mainnet

## 📋 Environment Checklist

### Testnet Environment
- [ ] `.env` points to testnet RPC
- [ ] Contract addresses are testnet
- [ ] KV namespaces are testnet
- [ ] Using testnet private key (safe to commit obfuscated)
- [ ] Telegram alerts go to test channel

### Mainnet Environment
- [ ] `.env` points to mainnet RPC
- [ ] Contract addresses are mainnet (after deployment)
- [ ] KV namespaces are production (fresh)
- [ ] Using hardware wallet (NO private key in file)
- [ ] Telegram alerts go to production channel
- [ ] Monitoring enabled
- [ ] Backup strategy active

## 🔍 Verify Current Environment

```bash
# Check RPC URL
grep MONAD_RPC_URL .env

# Check contract addresses
grep TOKEN_ADDRESS .env

# Check KV namespaces
grep CF_KV_CONTRIBUTORS_NAMESPACE_ID .env

# Verify with pause status
make pause-status
```

## 🛠️ Quick Commands

### Check Environment
```bash
# Show current RPC
echo "Backend RPC: $(grep MONAD_RPC_URL .env | cut -d= -f2)"
echo "Contracts RPC: $(grep MONAD_RPC_URL contracts/.env | cut -d= -f2)"
```

### Backup Current Environment
```bash
# Backup current .env
cp .env .env.backup.$(date +%Y%m%d_%H%M%S)
cp contracts/.env contracts/.env.backup.$(date +%Y%m%d_%H%M%S)
```

### Restore Environment
```bash
# Restore from backup
cp .env.backup.YYYYMMDD_HHMMSS .env
cp contracts/.env.backup.YYYYMMDD_HHMMSS contracts/.env
go run cmd/obfuscator-gen/main.go
```

## 📊 Environment Comparison

| Feature | Testnet | Mainnet |
|---------|---------|---------|
| **RPC** | testnet-rpc.monad.xyz | mainnet-rpc.monad.xyz |
| **Tokens** | Free (faucet) | Real value |
| **Risk** | Zero | High |
| **Speed** | Fast | Production |
| **Data** | Temporary | Permanent |
| **Keys** | Test keys OK | Hardware wallet only |
| **KV** | Test namespaces | Production namespaces |
| **Monitoring** | Optional | Mandatory |

## 🎯 Best Practices

### Development Workflow

1. **Develop on testnet**
   - Test all features
   - Verify contracts
   - Run full E2E tests

2. **Prepare mainnet**
   - Create production KV namespaces
   - Generate hardware wallet
   - Update .env with mainnet config

3. **Deploy to mainnet**
   - Follow DEPLOYMENT_MAINNET.md
   - Deploy contracts
   - Update addresses
   - Verify everything

4. **Maintain both**
   - Keep testnet for testing
   - Use mainnet for production
   - Switch as needed

### Safety Rules

1. **Always know which environment you're in**
2. **Never mix testnet and mainnet keys**
3. **Test on testnet before mainnet**
4. **Backup before switching**
5. **Verify after switching**

## 🆘 Troubleshooting

### "Wrong network" errors
- Check RPC URL in .env
- Verify contract addresses match network
- Regenerate constants

### "Contract not found" errors
- Ensure contracts deployed on current network
- Check addresses in .env
- Verify RPC connection

### "Insufficient funds" errors
- Testnet: Get tokens from faucet
- Mainnet: Fund wallet with real MON

### "Transaction failed" errors
- Check gas settings
- Verify network status
- Ensure correct private key/wallet

---

**Questions?** Check [DEPLOYMENT_MAINNET.md](DEPLOYMENT_MAINNET.md) or [DEPLOYMENT_TESTNET.md](DEPLOYMENT_TESTNET.md).
