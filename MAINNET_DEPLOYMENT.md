# Mainnet Deployment Guide

**⚠️ CRITICAL: This is for PRODUCTION deployment on Monad Mainnet with real USDC**

**Current Status:** Ready for deployment (85% complete - pending contract deployment)

---

## 🎯 PREREQUISITES

### Required
- [ ] All code changes from `MAINNET_PRODUCTION_CHECKLIST.md` completed
- [ ] `.env.mainnet` configured with correct values
- [ ] Private keys secured (hardware wallet or KMS recommended)
- [ ] Sufficient MON for gas fees (~10 MON recommended)
- [ ] Backup of current production data (if applicable)
- [ ] Team notified of deployment window

### Verification
```bash
# Verify environment detection
go run -tags mainnet main.go --check-env

# Verify build passes
go build -o /dev/null .

# Verify tests pass
make test

# Verify frontend builds
cd frontend && npm run build
```

---

## 📋 PRE-DEPLOYMENT CHECKLIST

### Configuration
- [ ] `.env.mainnet` has `ENVIRONMENT=mainnet`
- [ ] USDC address: `0x754704bc059f8c67012fed69bc8a327a5aafb603`
- [ ] RPC URL points to Monad Mainnet
- [ ] Private keys are NOT in git
- [ ] Cloudflare KV configured for production

### Smart Contracts
- [ ] All contracts compiled: `make contracts-compile`
- [ ] Contract tests pass: `make contracts-test`
- [ ] Deployment scripts reviewed
- [ ] Gas price acceptable

### Code
- [ ] Latest code from master branch
- [ ] No debug logs in production code
- [ ] Environment detection working
- [ ] Constants regenerated: `make constants-generate`

### Documentation
- [ ] User guide published: https://getkawai.com/docs
- [ ] API documentation updated
- [ ] Changelog prepared

---

## 🚀 DEPLOYMENT SEQUENCE

### PHASE 1: Smart Contract Deployment (30 minutes)

**Choose Your Deployment Strategy:**

#### Option A: Full Suite Deployment (Recommended for Fresh Start)

Use `DeployKawai.s.sol` to deploy everything at once. This script now **auto-detects** the environment:

**How it works:**
- ✅ If `USDC_ADDRESS` is set in `.env` → Uses existing USDC (mainnet mode)
- ✅ If `USDC_ADDRESS` is NOT set → Deploys MockUSDT (testnet mode)

**For Mainnet:**
```bash
# 1. Set USDC_ADDRESS in contracts/.env.mainnet
echo "USDC_ADDRESS=0x754704bc059f8c67012fed69bc8a327a5aafb603" >> contracts/.env.mainnet

# 2. Deploy full suite
cd contracts
forge script script/DeployKawai.s.sol:DeployKawai \
  --rpc-url $MONAD_MAINNET_RPC \
  --private-key $DEPLOYER_PRIVATE_KEY \
  --broadcast \
  --verify
```

**What it deploys:**
- ✅ Uses existing USDC (no MockUSDT deployment)
- ✅ KawaiToken
- ✅ PaymentVault (with USDC)
- ✅ MerkleDistributor (KAWAI Mining)
- ✅ MerkleDistributor (USDT Dividends)
- ✅ OTCMarket (Escrow)
- ✅ Automatic MINTER_ROLE grant to KAWAI distributor

**Advantages:**
- One command deploys everything
- Automatic MINTER_ROLE grant
- Consistent deployment
- Auto-detects mainnet vs testnet

#### Option B: Modular Deployment (Recommended for Flexibility)

Deploy contracts individually for more control:

```bash
# 1. Deploy KawaiToken first (if needed)
forge script script/DeployKawai.s.sol:DeployKawai \
  --rpc-url $MONAD_MAINNET_RPC \
  --private-key $DEPLOYER_PRIVATE_KEY \
  --broadcast \
  --verify

# 2. Deploy PaymentVault separately
make contracts-deploy-vault

# 3. Deploy distributors
make contracts-deploy-mining-mainnet
make contracts-deploy-cashback-mainnet
make contracts-deploy-referral-mainnet

# 4. Grant permissions
make contracts-grant-minter-mainnet
```

**Advantages:**
- More control over each deployment
- Easier to test individually
- Can deploy to different networks separately
- Better for upgrades/redeployments

**⚠️ CRITICAL: Save all deployed addresses immediately!**

#### Post-Deployment Configuration

After deploying contracts, update your configuration files:

```bash
# 1. Update .env.mainnet with deployed addresses
cat >> .env.mainnet << EOF
KAWAI_TOKEN_ADDRESS=<deployed_kawai_address>
PAYMENT_VAULT_ADDRESS=<deployed_vault_address>
MINING_DISTRIBUTOR_ADDRESS=<deployed_mining_address>
CASHBACK_DISTRIBUTOR_ADDRESS=<deployed_cashback_address>
REFERRAL_DISTRIBUTOR_ADDRESS=<deployed_referral_address>
EOF

# 2. Regenerate backend constants
go run cmd/obfuscator-gen/main.go

# 3. Verify generated files
git diff internal/constant/blockchain.go
git diff pkg/jarvis/db/project_tokens.go

# 4. Verify all contracts on MonadScan
# Visit https://monadexplorer.com and check each contract
```

---

### PHASE 2: Grant Permissions (10 minutes)

#### 2.1 Grant MINTER_ROLE to Distributors
```bash
# Mining Distributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $MINING_DISTRIBUTOR_ADDRESS \
  --rpc-url $MONAD_MAINNET_RPC \
  --private-key $ADMIN_PRIVATE_KEY

# Cashback Distributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $CASHBACK_DISTRIBUTOR_ADDRESS \
  --rpc-url $MONAD_MAINNET_RPC \
  --private-key $ADMIN_PRIVATE_KEY

# Referral Distributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $REFERRAL_DISTRIBUTOR_ADDRESS \
  --rpc-url $MONAD_MAINNET_RPC \
  --private-key $ADMIN_PRIVATE_KEY
```

#### 2.2 Verify Permissions
```bash
# Check MINTER_ROLE for each distributor
go run cmd/dev/check-minter-role/main.go
```

**Expected output:**
```
✅ Mining Distributor has MINTER_ROLE
✅ Cashback Distributor has MINTER_ROLE
✅ Referral Distributor has MINTER_ROLE
```

---

### PHASE 3: Backend Deployment (15 minutes)

#### 3.1 Build Production Binary
```bash
# Build with mainnet configuration
CGO_ENABLED=1 go build -tags mainnet -o veridium-mainnet .

# Verify binary
./veridium-mainnet --version
./veridium-mainnet --check-env
```

#### 3.2 Deploy to Production Server
```bash
# Stop current service (if running)
systemctl stop veridium

# Backup current binary
cp /opt/veridium/veridium /opt/veridium/veridium.backup

# Deploy new binary
scp veridium-mainnet production-server:/opt/veridium/veridium

# Deploy .env.mainnet
scp .env.mainnet production-server:/opt/veridium/.env

# Set permissions
ssh production-server "chmod +x /opt/veridium/veridium"
```

#### 3.3 Start Service
```bash
# Start service
systemctl start veridium

# Check status
systemctl status veridium

# Check logs
journalctl -u veridium -f
```

**Expected logs:**
```
✅ Environment: mainnet
✅ USDC Address: 0x754704bc059f8c67012fed69bc8a327a5aafb603
✅ Kawai Token: 0x...
✅ All contracts initialized
✅ Server started on :8080
```

---

### PHASE 4: Frontend Deployment (10 minutes)

#### 4.1 Build Production Frontend
```bash
cd frontend

# Install dependencies
npm install

# Build for production
npm run build

# Verify build
ls -lh dist/
```

#### 4.2 Deploy to CDN/Hosting
```bash
# Example: Deploy to Cloudflare Pages
wrangler pages publish dist/

# Or deploy to your hosting provider
# rsync -avz dist/ production-server:/var/www/kawai/
```

#### 4.3 Verify Frontend
- Open https://app.getkawai.com
- Check network indicator shows "Monad Mainnet"
- Check deposit modal shows USDC warning
- Verify all links work

---

### PHASE 5: Initial Testing (20 minutes)

#### 5.1 Acquire Test USDC
```bash
# Option 1: Bridge from Ethereum
# Visit https://monadbridge.com
# Bridge small amount (e.g., $10 USDC)

# Option 2: Buy on DEX
# Swap MON → USDC on Monad DEX
```

#### 5.2 Test Deposit Flow
1. Open Kawai Desktop
2. Connect wallet with test USDC
3. Navigate to Deposit
4. Read warning message
5. Enter small amount (e.g., 5 USDC)
6. Approve USDC
7. Confirm deposit
8. Verify balance updated

#### 5.3 Monitor Transaction
```bash
# Check transaction on explorer
# Verify:
# - Transaction successful
# - USDC transferred to PaymentVault
# - Gas cost reasonable
# - No errors in logs
```

#### 5.4 Test Balance Checking
```bash
# Check USDC balance in PaymentVault
cast call $PAYMENT_VAULT_ADDRESS \
  "balanceOf(address)(uint256)" \
  $TEST_USER_ADDRESS \
  --rpc-url $MONAD_MAINNET_RPC

# Should match deposited amount
```

---

### PHASE 6: Monitoring Setup (15 minutes)

#### 6.1 Set Up Alerts
```bash
# Configure monitoring for:
# - Transaction failures
# - Low MON balance (for gas)
# - High error rates
# - Unusual activity
```

#### 6.2 Set Up Logging
```bash
# Configure log aggregation
# - Application logs
# - Transaction logs
# - Error logs
# - Performance metrics
```

#### 6.3 Set Up Dashboards
- Transaction volume
- User deposits
- USDC balance in vault
- Gas costs
- Error rates

---

## ✅ POST-DEPLOYMENT VERIFICATION

### Immediate Checks (First Hour)
- [ ] Backend is running without errors
- [ ] Frontend loads correctly
- [ ] Deposit flow works
- [ ] Balance checking works
- [ ] No critical errors in logs
- [ ] Monitoring alerts working

### First Day Checks
- [ ] Multiple users can deposit
- [ ] Transactions are fast (<5 seconds)
- [ ] Gas costs are reasonable
- [ ] No memory leaks
- [ ] No database issues
- [ ] User feedback is positive

### First Week Checks
- [ ] System is stable
- [ ] No security issues
- [ ] Performance is good
- [ ] User adoption growing
- [ ] No critical bugs reported

---

## 🆘 ROLLBACK PROCEDURE

### If Critical Issue Found

#### 1. Immediate Actions
```bash
# Stop accepting new deposits
# Option 1: Pause contracts (if pause mechanism available)
cast send $PAYMENT_VAULT_ADDRESS "pause()" \
  --rpc-url $MONAD_MAINNET_RPC \
  --private-key $ADMIN_PRIVATE_KEY

# Option 2: Stop backend service
systemctl stop veridium

# Option 3: Show maintenance page on frontend
```

#### 2. Assess Impact
- How many users affected?
- How much USDC at risk?
- What is the root cause?
- Can it be fixed quickly?

#### 3. Communication
- Notify users via Discord/Twitter
- Update status page
- Provide timeline for resolution

#### 4. Fix and Redeploy
- Fix issue in code
- Test thoroughly on testnet
- Deploy fix to mainnet
- Resume operations

---

## 🔒 SECURITY CONSIDERATIONS

### Private Key Management
- **NEVER** commit private keys to git
- Use hardware wallet for admin operations
- Use KMS for production keys
- Rotate keys regularly
- Limit key access to essential personnel

### Contract Security
- All contracts audited before deployment
- Emergency pause mechanism tested
- Multi-sig for critical operations
- Rate limiting on deposits
- Maximum deposit limits (if needed)

### Monitoring
- Alert on unusual transactions
- Monitor for reentrancy attacks
- Track failed transactions
- Monitor gas price spikes
- Alert on low balances

---

## 📞 EMERGENCY CONTACTS

### Technical Team
- **Smart Contracts:** [Contract Developer]
- **Backend:** [Backend Developer]
- **Frontend:** [Frontend Developer]
- **DevOps:** [DevOps Engineer]

### Emergency Procedures
- **Critical Bug:** Pause contracts immediately
- **Security Issue:** Contact security team
- **Infrastructure:** Contact DevOps
- **User Support:** Contact support team

---

## 📝 DEPLOYMENT CHECKLIST SUMMARY

### Pre-Deployment
- [x] Code ready (85% complete)
- [ ] Contracts deployed
- [ ] Permissions granted
- [ ] Configuration updated
- [ ] Tests passed

### Deployment
- [ ] Contracts deployed to mainnet
- [ ] Backend deployed
- [ ] Frontend deployed
- [ ] Monitoring configured

### Post-Deployment
- [ ] Initial testing complete
- [ ] No critical errors
- [ ] Users can deposit
- [ ] Monitoring active

---

## 🎯 SUCCESS CRITERIA

### Technical
✅ All contracts deployed and verified  
✅ Backend running without errors  
✅ Frontend loads correctly  
✅ Deposit flow works end-to-end  
✅ Monitoring and alerts active  

### Business
✅ Users can deposit USDC  
✅ Transactions are fast and cheap  
✅ No security issues  
✅ Positive user feedback  
✅ System is stable  

---

**Last Updated:** January 21, 2026  
**Status:** Ready for deployment  
**Estimated Deployment Time:** 2-3 hours  
**Risk Level:** LOW (all safety measures in place)
