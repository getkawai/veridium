# Mainnet Deployment

**Simple guide. No bullshit.**

---

## Before Deploy

1. **Wallet ready?** Fund dengan 1 MON
2. **KV namespaces created?** 10 production namespaces di Cloudflare
3. **Tested on testnet?** Yes (Round 6)

---

## Deploy (2 hours)

## Deploy (2 hours)

### 1. Deploy Contracts (30 min)

```bash
cd contracts
# Update contracts/.env: MONAD_RPC_URL=https://mainnet-rpc.monad.xyz
make contracts-deploy-mainnet
make contracts-deploy-mining-mainnet
make contracts-deploy-cashback-mainnet
make contracts-deploy-referral-mainnet
make contracts-grant-minter-mining
make contracts-grant-minter-cashback
make contracts-grant-minter-referral
```

### 2. Update Backend (5 min)

```bash
# Copy addresses from contracts/.env to .env
go run cmd/obfuscator-gen/main.go
go build -o main main.go
```

### 3. Test Pause (5 min)

```bash
make pause-status
make pause-all
make unpause-all
```

### 4. Deploy (30 min)

- Upload backend to server
- Deploy frontend
- Test with real wallet

**Done!**

---

## Emergency

```bash
make pause-all  # Pause everything
# Fix issue
make unpause-all  # Resume
```

---

## Back to Testnet

```bash
cp .env.testnet .env
go run cmd/obfuscator-gen/main.go
```

### STEP 1: CREATE PRODUCTION KV NAMESPACES (10 minutes)

**⚠️ CRITICAL:** Do NOT use testnet namespaces! Create fresh production namespaces.

1. Login to Cloudflare Dashboard
2. Navigate to Workers & Pages → KV
3. Create 8 new namespaces:
   ```
   kawai-prod-contributors
   kawai-prod-proofs
   kawai-prod-settlements
   kawai-prod-cashback
   kawai-prod-holders
   kawai-prod-marketplace
   kawai-prod-users
   kawai-prod-revenue
   ```
4. Copy namespace IDs to `.env.production`

**Verify:** All namespaces are EMPTY (no data from testnet)

---

### STEP 2: PREPARE ENVIRONMENT FILES (5 minutes)

#### 2.1 Create `.env.production`
```bash
# Blockchain Configuration
MONAD_RPC_URL=https://mainnet-rpc.monad.xyz
MONAD_RPC_URL_BACKUP=<BACKUP_RPC_URL>

# Contract Addresses (WILL BE FILLED AFTER DEPLOYMENT)
USDT_TOKEN_ADDRESS=<REAL_USDT_ADDRESS_ON_MONAD>
TOKEN_ADDRESS=
ESCROW_ADDRESS=
PAYMENT_VAULT_ADDRESS=
KAWAI_DISTRIBUTOR_ADDRESS=
USDT_DISTRIBUTOR_ADDRESS=
MINING_DISTRIBUTOR_ADDRESS=
CASHBACK_DISTRIBUTOR_ADDRESS=

# Cloudflare KV (PRODUCTION NAMESPACES)
CF_ACCOUNT_ID=<YOUR_ACCOUNT_ID>
CF_API_TOKEN=<PRODUCTION_API_TOKEN>
CF_KV_CONTRIBUTORS_NAMESPACE_ID=<PROD_NAMESPACE_ID>
CF_KV_PROOFS_NAMESPACE_ID=<PROD_NAMESPACE_ID>
CF_KV_SETTLEMENTS_NAMESPACE_ID=<PROD_NAMESPACE_ID>
CF_KV_CASHBACK_NAMESPACE_ID=<PROD_NAMESPACE_ID>
CF_KV_HOLDER_NAMESPACE_ID=<PROD_NAMESPACE_ID>
CF_KV_P2PMARKETPLACE_NAMESPACE_ID=<PROD_NAMESPACE_ID>
CF_KV_USERS_NAMESPACE_ID=<PROD_NAMESPACE_ID>
CF_KV_REVENUE_NAMESPACE_ID=<PROD_NAMESPACE_ID>

# Admin Configuration
ADMIN_ADDRESS=<ADMIN_WALLET_ADDRESS>
ADMIN_PRIVATE_KEY=<FROM_HARDWARE_WALLET_OR_KMS>

# Treasury Addresses (VERIFY THESE!)
TREASURY_ADDRESSES=<COMMA_SEPARATED_ADDRESSES>

# Telegram Alerts (PRODUCTION BOT)
TELEGRAM_BOT_TOKEN=<PRODUCTION_BOT_TOKEN>
TELEGRAM_CHAT_ID=<PRODUCTION_CHAT_ID>

# API Keys (PRODUCTION KEYS)
ETHERSCAN_API_KEYS=<PRODUCTION_KEYS>
OPENROUTER_API_KEYS=<PRODUCTION_KEYS>
GEMINI_API_KEYS=<PRODUCTION_KEYS>
```

#### 2.2 Create `contracts/.env.production`
```bash
# Monad Mainnet
RPC_URL=https://mainnet-rpc.monad.xyz
PRIVATE_KEY=<ADMIN_WALLET_PRIVATE_KEY>
CHAIN_ID=<MONAD_MAINNET_CHAIN_ID>

# Will be filled after KawaiToken deployment
KAWAI_TOKEN_ADDRESS=
```

**⚠️ SECURITY CHECK:**
- [ ] `.env.production` is in `.gitignore`
- [ ] No private keys committed to git
- [ ] Private keys stored securely (KMS/hardware wallet)

---

### STEP 3: DEPLOY CONTRACTS TO MAINNET (30 minutes)

#### 3.1 Final Pre-Deployment Checks
```bash
# 1. Verify you're on mainnet RPC
echo $RPC_URL
# Should output: https://mainnet-rpc.monad.xyz

# 2. Check admin wallet balance
cast balance <ADMIN_WALLET_ADDRESS> --rpc-url $RPC_URL
# Should have at least 1 MON

# 3. Verify contract code one last time
cd contracts
forge test
# All tests must pass

# 4. Estimate gas costs
forge script script/DeployKawai.s.sol:DeployKawai --rpc-url $RPC_URL
# Review gas estimates
```

#### 3.2 Deploy Main Contracts
```bash
# Deploy: KawaiToken, KAWAI_Distributor, USDT_Distributor, PaymentVault, OTCMarket
# NOTE: NO MockUSDT on mainnet!
cd contracts
forge script script/DeployKawai.s.sol:DeployKawai \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --verify

# ⏳ Wait for deployment (5-10 minutes)
# ⏳ Wait for verification (5-10 minutes)
```

**Expected Output:**
```
✅ KawaiToken deployed to: 0x...
✅ KAWAI_Distributor deployed to: 0x...
✅ USDT_Distributor deployed to: 0x...
✅ PaymentVault deployed to: 0x...
✅ OTCMarket deployed to: 0x...
✅ All contracts verified on explorer
```

**⚠️ CRITICAL:** Save all contract addresses immediately!

#### 3.3 Update Environment Files
```bash
# Update contracts/.env.production
KAWAI_TOKEN_ADDRESS=<FROM_STEP_3.2>

# Update .env.production
TOKEN_ADDRESS=<KAWAI_TOKEN_ADDRESS>
ESCROW_ADDRESS=<OTC_MARKET_ADDRESS>
PAYMENT_VAULT_ADDRESS=<PAYMENT_VAULT_ADDRESS>
KAWAI_DISTRIBUTOR_ADDRESS=<KAWAI_DISTRIBUTOR_ADDRESS>
USDT_DISTRIBUTOR_ADDRESS=<USDT_DISTRIBUTOR_ADDRESS>
```

#### 3.4 Deploy MiningRewardDistributor
```bash
forge script script/DeployMiningDistributor.s.sol:DeployMiningDistributor \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --verify
```

**Save address:**
```bash
MINING_DISTRIBUTOR_ADDRESS=<FROM_OUTPUT>
```

#### 3.5 Deploy DepositCashbackDistributor
```bash
forge script script/DeployCashbackDistributor.s.sol:DeployCashbackDistributor \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --verify
```

**Save address:**
```bash
CASHBACK_DISTRIBUTOR_ADDRESS=<FROM_OUTPUT>
```

#### 3.6 Verify All Contracts on Explorer
```bash
# Visit Monad Explorer
open https://monadexplorer.com/address/<CONTRACT_ADDRESS>

# Verify for each contract:
# - Contract verified ✅
# - Source code visible ✅
# - Constructor arguments correct ✅
```

---

### STEP 4: GRANT MINTER_ROLE (5 minutes)

```bash
# Grant MINTER_ROLE to MiningRewardDistributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $MINING_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY

# Grant MINTER_ROLE to DepositCashbackDistributor
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $CASHBACK_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY

# Verify roles granted
cast call $KAWAI_TOKEN_ADDRESS \
  "hasRole(bytes32,address)(bool)" \
  $(cast keccak "MINTER_ROLE") \
  $MINING_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL
# Should return: true
```

**Telegram Alert:** You should receive alerts for each transaction

---

### STEP 5: UPDATE BACKEND CODE (5 minutes)

```bash
# Copy production env
cp .env.production .env

# Regenerate blockchain constants
go run cmd/obfuscator-gen/main.go

# Verify generated files
cat internal/constant/blockchain.go
# Should show mainnet addresses

# Build production binary
PRODUCTION=true go build -o kawai-backend \
  -tags production \
  -trimpath \
  -ldflags="-w -s" \
  main.go

# Test binary
./kawai-backend --version
```

---

### STEP 6: DEPLOY BACKEND TO PRODUCTION SERVER (15 minutes)

```bash
# 1. Upload binary to server
scp kawai-backend user@api.kawai.network:/opt/kawai/

# 2. Upload .env.production
scp .env.production user@api.kawai.network:/opt/kawai/.env

# 3. SSH to server
ssh user@api.kawai.network

# 4. Setup systemd service
sudo nano /etc/systemd/system/kawai-backend.service
```

**Service file:**
```ini
[Unit]
Description=Kawai Backend API
After=network.target

[Service]
Type=simple
User=kawai
WorkingDirectory=/opt/kawai
ExecStart=/opt/kawai/kawai-backend
Restart=always
RestartSec=10
Environment="PRODUCTION=true"

[Install]
WantedBy=multi-user.target
```

```bash
# 5. Start service
sudo systemctl daemon-reload
sudo systemctl enable kawai-backend
sudo systemctl start kawai-backend

# 6. Check status
sudo systemctl status kawai-backend
sudo journalctl -u kawai-backend -f

# 7. Test API
curl https://api.kawai.network/health
# Should return: {"status":"ok"}
```

---

### STEP 7: DEPLOY FRONTEND TO PRODUCTION (15 minutes)

```bash
# 1. Update frontend config
cd frontend
nano .env.production
```

**Frontend .env.production:**
```bash
VITE_API_URL=https://api.kawai.network
VITE_KAWAI_TOKEN_ADDRESS=<FROM_STEP_3>
VITE_MINING_DISTRIBUTOR_ADDRESS=<FROM_STEP_3>
VITE_CASHBACK_DISTRIBUTOR_ADDRESS=<FROM_STEP_3>
VITE_RPC_URL=https://mainnet-rpc.monad.xyz
VITE_CHAIN_ID=<MONAD_MAINNET_CHAIN_ID>
```

```bash
# 2. Build production frontend
npm run build

# 3. Deploy to hosting (Vercel/Netlify/Cloudflare Pages)
# Example for Vercel:
vercel --prod

# Or upload to server:
scp -r dist/* user@app.kawai.network:/var/www/kawai/

# 4. Test frontend
open https://app.kawai.network
```

---

### STEP 8: SMOKE TEST WITH REAL WALLET (15 minutes)

**⚠️ Use small amounts for testing!**

#### 8.1 Connect Wallet
1. Open https://app.kawai.network
2. Connect your personal wallet (NOT admin wallet!)
3. Verify network: Monad Mainnet
4. Verify contract addresses displayed correctly

#### 8.2 Test Mining Rewards (if you have real mining data)
1. Navigate to "Rewards" → "Mining Rewards"
2. Check if claimable rewards show correctly
3. If claimable, try claiming **small amount** first
4. Verify transaction on explorer
5. Verify KAWAI balance increased

#### 8.3 Test Cashback (if you have real deposits)
1. Navigate to "Rewards" → "Deposit Cashback"
2. Check if cashback shows correctly
3. If claimable, try claiming **small amount** first
4. Verify transaction on explorer
5. Verify KAWAI balance increased

#### 8.4 Verify Backend Logs
```bash
ssh user@api.kawai.network
sudo journalctl -u kawai-backend -n 100

# Should see:
# - No errors
# - Successful API calls
# - Telegram alerts sent
```

---

### STEP 9: FIRST SETTLEMENT (WHEN READY)

**⚠️ Only run when you have real user data!**

```bash
# 1. SSH to production server
ssh user@api.kawai.network

# 2. Generate mining settlement
cd /opt/kawai
./kawai-backend settlement generate --type mining

# 3. Review settlement
# - Check contributor count
# - Check total amount
# - Verify Merkle root generated

# 4. Upload Merkle root
./kawai-backend settlement upload --type mining

# 5. Verify on-chain
cast call $MINING_DISTRIBUTOR_ADDRESS \
  "periodMerkleRoots(uint256)(bytes32)" \
  <PERIOD_ID> \
  --rpc-url $RPC_URL

# 6. Check Telegram alerts
# Should receive:
# - Settlement generation success
# - Merkle root upload success
```

---

## 📊 POST-DEPLOYMENT MONITORING (24 HOURS)

### Hour 1-4: Critical Monitoring
- [ ] Check backend logs every 15 minutes
- [ ] Monitor Telegram alerts
- [ ] Check server resources (CPU, RAM, disk)
- [ ] Verify API response times (<2s)
- [ ] Test claiming with multiple users

### Hour 4-12: Active Monitoring
- [ ] Check logs every hour
- [ ] Monitor gas costs
- [ ] Check contract balances
- [ ] Verify settlement operations
- [ ] Monitor user feedback

### Hour 12-24: Passive Monitoring
- [ ] Check logs every 2 hours
- [ ] Monitor alerts
- [ ] Check success metrics
- [ ] Document any issues

---

## 🆘 EMERGENCY PROCEDURES

### Backend Down
```bash
# 1. Check status
ssh user@api.kawai.network
sudo systemctl status kawai-backend

# 2. Check logs
sudo journalctl -u kawai-backend -n 100

# 3. Restart
sudo systemctl restart kawai-backend

# 4. If still down, rollback
sudo systemctl stop kawai-backend
# Deploy previous version
sudo systemctl start kawai-backend
```

### Contract Issues
```bash
# 1. PAUSE IMMEDIATELY (if critical)
cast send $MINING_DISTRIBUTOR_ADDRESS \
  "pause()" \
  --rpc-url $RPC_URL \
  --private-key $ADMIN_PRIVATE_KEY

# 2. Investigate
# - Check explorer for failed transactions
# - Review contract events
# - Check Merkle roots

# 3. Fix and test on testnet first
# 4. Deploy fix to mainnet
# 5. Unpause
```

### High Gas Costs
```bash
# 1. Check current gas price
cast gas-price --rpc-url $RPC_URL

# 2. If too high, delay settlement
# 3. Monitor gas prices
# 4. Run settlement when gas is low
```

---

## ✅ SUCCESS CRITERIA

### Day 1
- [ ] 0 critical errors
- [ ] 10+ successful claims
- [ ] <2s API response time
- [ ] 100% uptime
- [ ] All alerts working

### Week 1
- [ ] 100+ successful claims
- [ ] <1% error rate
- [ ] <$0.50 avg gas cost
- [ ] 99.9% uptime
- [ ] No security incidents

### Month 1
- [ ] 1000+ successful claims
- [ ] <0.5% error rate
- [ ] User satisfaction >4.5/5
- [ ] 99.95% uptime
- [ ] Automated operations running smoothly

---

## 📝 DEPLOYMENT CHECKLIST SUMMARY

### Pre-Deployment
- [ ] Security audit completed
- [ ] Wallets secured (hardware wallet/KMS)
- [ ] Infrastructure ready (server, domain, SSL)
- [ ] Production KV namespaces created
- [ ] Real USDT address obtained
- [ ] Monitoring configured
- [ ] Backup strategy implemented

### Deployment
- [ ] Contracts deployed to mainnet
- [ ] Contracts verified on explorer
- [ ] MINTER_ROLE granted
- [ ] Backend deployed to production
- [ ] Frontend deployed to production
- [ ] Smoke test passed

### Post-Deployment
- [ ] 24-hour monitoring completed
- [ ] First settlement successful
- [ ] User claims working
- [ ] Alerts functioning
- [ ] Documentation updated

---

## 💰 ESTIMATED COSTS

### One-Time (Deployment)
- Contract deployment: ~0.1 MON
- Setup transactions: ~0.01 MON
- **Total:** ~0.11 MON (~$X USD)

### Monthly (Operations)
- Weekly settlements: ~0.02 MON/month
- Emergency operations: ~0.01 MON/month
- **Total:** ~0.03 MON/month (~$X USD)

### Infrastructure
- Server: $50-100/month
- Domain: $10/year
- SSL: Free (Let's Encrypt)
- Cloudflare KV: Free tier (up to 100k reads/day)
- **Total:** ~$50-100/month

---

## 📞 SUPPORT CONTACTS

- **Smart Contracts:** [Your Name] - [Contact]
- **Backend:** [Your Name] - [Contact]
- **Frontend:** [Your Name] - [Contact]
- **Infrastructure:** [Your Name] - [Contact]
- **Monad Support:** [Monad Team Contact]

---

## 🎯 FINAL CHECKLIST BEFORE GO-LIVE

- [ ] All tests passed
- [ ] Security audit reviewed
- [ ] Wallets secured
- [ ] Infrastructure ready
- [ ] Monitoring active
- [ ] Team briefed
- [ ] Emergency procedures documented
- [ ] Rollback plan ready
- [ ] User documentation complete
- [ ] Support channels ready

**Ready to deploy?** Review this checklist one more time, then execute step by step.

**Remember:** Take your time, double-check everything, and monitor closely for the first 24 hours.

**Good luck! 🚀**
