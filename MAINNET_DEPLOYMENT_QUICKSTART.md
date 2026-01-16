# Mainnet Deployment - Quick Start

**⚠️ CRITICAL: This is for PRODUCTION deployment. Double-check everything!**

## 📋 Pre-Deployment Checklist

### 1. Hardware Wallet Setup ✅
- [ ] Ledger/Trezor purchased and initialized
- [ ] Seed phrase backed up (fireproof safe)
- [ ] Test transaction on testnet successful
- [ ] Wallet funded with 1 MON minimum

### 2. Production Infrastructure ✅
- [ ] Production server provisioned
- [ ] Domain configured (api.kawai.network, app.kawai.network)
- [ ] SSL certificates installed
- [ ] Cloudflare KV production namespaces created (FRESH)
- [ ] Monitoring and alerts configured
- [ ] Backup strategy in place

### 3. Security Review ✅
- [ ] Emergency pause mechanism tested
- [ ] All critical bugs fixed
- [ ] Smart contracts reviewed (or risk accepted)
- [ ] Team trained on emergency procedures
- [ ] Rollback plan documented

## 🚀 Deployment Steps

### Step 1: Prepare Environment (5 minutes)

```bash
# Ensure you're on mainnet config
cat .env | grep MONAD_RPC_URL
# Should show: https://mainnet-rpc.monad.xyz

# If not, you're still on testnet - switch first!
# See ENV_SWITCHING_GUIDE.md
```

### Step 2: Create Production KV Namespaces (10 minutes)

```bash
# Login to Cloudflare dashboard
# Create 10 fresh KV namespaces:
# - kawai-contributors-prod
# - kawai-users-prod
# - kawai-proofs-prod
# - kawai-settlements-prod
# - kawai-apikey-prod
# - kawai-authz-prod
# - kawai-p2pmarketplace-prod
# - kawai-cashback-prod
# - kawai-revenue-prod
# - kawai-holder-prod

# Update .env with namespace IDs
```

### Step 3: Deploy Smart Contracts (30 minutes)

```bash
# Connect hardware wallet
# Ensure it's unlocked and ready

# Deploy main contracts
cd contracts
make contracts-deploy-mainnet

# ⚠️ IMPORTANT: Save all contract addresses!
# Update contracts/.env with addresses

# Deploy distributors
make contracts-deploy-mining-mainnet
make contracts-deploy-cashback-mainnet
make contracts-deploy-referral-mainnet

# Grant MINTER_ROLE
make contracts-grant-minter-mining
make contracts-grant-minter-cashback
make contracts-grant-minter-referral
```

### Step 4: Update Backend Configuration (5 minutes)

```bash
# Update .env with deployed contract addresses
# Copy from contracts/.env to .env

# Regenerate constants
go run cmd/obfuscator-gen/main.go

# Verify
cat internal/constant/blockchain.go
```

### Step 5: Test Pause Mechanism (5 minutes)

```bash
# Check status
make pause-status

# Test pause (dry-run first)
make pause-all-dry

# Actual pause
make pause-all

# Verify paused
make pause-status

# Unpause
make unpause-all

# Verify active
make pause-status
```

### Step 6: Deploy Backend (15 minutes)

```bash
# Build backend
go build -o main main.go

# Upload to production server
scp main user@server:/opt/kawai/

# Start service
ssh user@server
cd /opt/kawai
./main

# Verify health
curl https://api.kawai.network/health
```

### Step 7: Deploy Frontend (10 minutes)

```bash
# Update frontend config with mainnet addresses
cd frontend
# Edit config with mainnet contract addresses

# Build
npm run build

# Deploy to server
scp -r dist/* user@server:/var/www/kawai/

# Verify
curl https://app.kawai.network
```

### Step 8: Smoke Tests (15 minutes)

```bash
# Test basic operations:
# 1. Connect wallet
# 2. Check balances
# 3. Deposit USDT (small amount)
# 4. Use AI service (1 request)
# 5. Check rewards accumulation
# 6. Generate referral code
# 7. Check P2P marketplace

# If all pass → GO LIVE! 🎉
```

## 🔍 Post-Deployment Verification

### Immediate (First Hour)

```bash
# Monitor logs
tail -f backend-prod.log

# Check Telegram alerts
# Should receive deployment notification

# Verify contract interactions
make pause-status

# Check KV data
# Should see initial data being written
```

### First Day

- [ ] Monitor transaction volume
- [ ] Check error rates
- [ ] Verify reward calculations
- [ ] Test claiming flow
- [ ] Monitor gas usage
- [ ] Check user signups

### First Week

- [ ] Weekly settlement successful
- [ ] Users claiming rewards
- [ ] P2P trading active
- [ ] No critical bugs
- [ ] Performance acceptable
- [ ] Monitoring alerts working

## 🚨 Emergency Procedures

### If Critical Bug Found

```bash
# 1. PAUSE IMMEDIATELY
make pause-all

# 2. Notify team
# Send alert to Telegram

# 3. Investigate
# Check logs, transactions, KV data

# 4. Fix or rollback
# Deploy fix or revert to previous version

# 5. Test on testnet
# Verify fix works

# 6. Unpause
make unpause-all
```

### If Contracts Need Update

```bash
# 1. Deploy new contracts
# 2. Migrate data if needed
# 3. Update backend config
# 4. Regenerate constants
# 5. Restart services
# 6. Verify everything works
```

## 📊 Monitoring Checklist

### Metrics to Watch

- **Transaction success rate** (target: >99%)
- **API response time** (target: <500ms)
- **Error rate** (target: <0.1%)
- **Gas usage** (monitor for spikes)
- **KV read/write operations**
- **User signups** (growth rate)
- **Reward claims** (weekly)
- **P2P trading volume**

### Alerts to Configure

- Contract paused/unpaused
- High error rate (>1%)
- Slow API response (>2s)
- Failed transactions (>5 in 1 hour)
- Low MON balance (admin wallet)
- KV quota exceeded
- Server CPU/RAM high (>80%)

## 🔗 Quick Links

- **Full Guide**: [DEPLOYMENT_MAINNET.md](DEPLOYMENT_MAINNET.md)
- **Environment Switching**: [ENV_SWITCHING_GUIDE.md](ENV_SWITCHING_GUIDE.md)
- **Emergency Pause**: [EMERGENCY_PAUSE_GUIDE.md](EMERGENCY_PAUSE_GUIDE.md)
- **Testnet Guide**: [DEPLOYMENT_TESTNET.md](DEPLOYMENT_TESTNET.md)

## 📞 Emergency Contacts

- **Team Lead**: [Contact info]
- **DevOps**: [Contact info]
- **Security**: [Contact info]
- **Telegram**: [Group link]

## ⏱️ Estimated Timeline

| Phase | Duration | Critical? |
|-------|----------|-----------|
| Pre-deployment prep | 2-4 hours | ✅ Yes |
| Contract deployment | 30 minutes | ✅ Yes |
| Backend deployment | 15 minutes | ✅ Yes |
| Frontend deployment | 10 minutes | ✅ Yes |
| Testing & verification | 1 hour | ✅ Yes |
| **Total** | **4-6 hours** | |

## 💡 Pro Tips

1. **Deploy during low-traffic hours** (e.g., Sunday 2 AM UTC)
2. **Have rollback plan ready** before starting
3. **Test everything on testnet first** (same day)
4. **Keep team on standby** for first 24 hours
5. **Monitor closely** for first week
6. **Document everything** as you go
7. **Celebrate success** but stay vigilant! 🎉

---

**Ready to deploy?** Follow this guide step-by-step. Don't skip steps!

**Questions?** Review [DEPLOYMENT_MAINNET.md](DEPLOYMENT_MAINNET.md) for detailed instructions.

**Emergency?** Check [EMERGENCY_PAUSE_GUIDE.md](EMERGENCY_PAUSE_GUIDE.md) immediately.
