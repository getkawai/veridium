# Production Deployment Guide - Kawai Rewards System

**Target Network:** Monad Mainnet  
**Estimated Time:** 2-3 hours  
**Risk Level:** Medium (financial contracts)

---

## 📋 PRE-DEPLOYMENT CHECKLIST

### 1. Infrastructure Preparation

#### ✅ Server/Hosting
- [ ] Production server provisioned (AWS/GCP/DigitalOcean)
- [ ] Domain configured (e.g., api.kawai.network)
- [ ] SSL certificate installed (Let's Encrypt/Cloudflare)
- [ ] Firewall rules configured (allow 443, 8080)
- [ ] Monitoring setup (Datadog/New Relic/Grafana)
- [ ] Log aggregation configured (CloudWatch/Papertrail)
- [ ] Backup strategy defined (daily snapshots)

#### ✅ Database/Storage
- [ ] Cloudflare KV production namespaces created
- [ ] KV API tokens generated (production-only)
- [ ] Database backup schedule configured
- [ ] Data retention policy defined (90 days recommended)

#### ✅ Blockchain Access
- [ ] Monad Mainnet RPC endpoints configured
  - Primary: `https://mainnet-rpc.monad.xyz`
  - Backup: Multiple providers (Ankr, Infura equivalent)
- [ ] RPC rate limits verified (10k req/day minimum)
- [ ] Websocket connection for events (optional but recommended)

---

### 2. Wallet & Keys Preparation

#### ✅ Admin Wallet (Contract Owner)
- [ ] **CRITICAL:** Generate new wallet for production
  - Use hardware wallet (Ledger/Trezor) if possible
  - Or use secure key management service (AWS KMS/HashiCorp Vault)
- [ ] Fund with MON for gas (minimum 10 MON recommended)
- [ ] Backup private key securely (encrypted, offline)
- [ ] Document wallet address in secure location
- [ ] **NEVER** commit private key to git

#### ✅ Settlement Wallet (Backend)
- [ ] Generate dedicated wallet for settlement operations
- [ ] Fund with MON for gas (minimum 5 MON)
- [ ] Store private key in environment variable (not in code)
- [ ] Rotate key every 90 days (security best practice)

#### ✅ Treasury Wallets
- [ ] Verify treasury addresses for revenue distribution
- [ ] Test small transactions to each address
- [ ] Document all addresses in secure spreadsheet

---

### 3. Smart Contracts Preparation

#### ✅ Contract Deployment
- [ ] Review all contract code (security audit recommended)
- [ ] Test on testnet one final time
- [ ] Prepare deployment script with correct parameters
- [ ] Calculate total gas cost (estimate: 0.5-1 MON)
- [ ] Verify contract source code on explorer after deployment

#### ✅ Contract Configuration
- [ ] Set correct token allocations:
  - Mining: 500M KAWAI
  - Cashback: 200M KAWAI
  - Referral: 150M KAWAI
  - Revenue: Based on actual revenue
- [ ] Configure period duration (weekly = 604800 seconds)
- [ ] Set admin addresses correctly
- [ ] Grant MINTER_ROLE to distributors
- [ ] Test all admin functions

#### ✅ Contract Verification
- [ ] Verify on Monad Explorer (monadexplorer.com)
- [ ] Test read functions (currentPeriod, merkleRoots)
- [ ] Test write functions with small amounts first
- [ ] Document all contract addresses

---

### 4. Backend Configuration

#### ✅ Environment Variables (.env)
```bash
# Blockchain
MONAD_RPC_URL=https://mainnet-rpc.monad.xyz
MONAD_RPC_URL_BACKUP=https://backup-rpc.monad.xyz

# Contracts (FILL AFTER DEPLOYMENT)
TOKEN_ADDRESS=0x...
MINING_DISTRIBUTOR_ADDRESS=0x...
CASHBACK_DISTRIBUTOR_ADDRESS=0x...
REFERRAL_DISTRIBUTOR_ADDRESS=0x...
USDT_DISTRIBUTOR_ADDRESS=0x...
PAYMENT_VAULT_ADDRESS=0x...

# Admin (USE SECURE KEY MANAGEMENT)
ADMIN_ADDRESS=0x...
ADMIN_PRIVATE_KEY=  # Load from AWS Secrets Manager/Vault

# Cloudflare KV (PRODUCTION NAMESPACES)
CF_ACCOUNT_ID=...
CF_API_TOKEN=...  # Production token, not testnet
CF_KV_CONTRIBUTORS_NAMESPACE_ID=...
CF_KV_PROOFS_NAMESPACE_ID=...
CF_KV_SETTLEMENTS_NAMESPACE_ID=...
CF_KV_CASHBACK_NAMESPACE_ID=...
CF_KV_REVENUE_NAMESPACE_ID=...

# API Keys (PRODUCTION)
OPENROUTER_API_KEYS=...
GEMINI_API_KEYS=...

# Monitoring
SENTRY_DSN=...  # Error tracking
DATADOG_API_KEY=...  # Metrics

# Telegram Alerts
TELEGRAM_BOT_TOKEN=...
TELEGRAM_CHAT_ID=...  # For critical alerts
```

- [ ] All environment variables configured
- [ ] No hardcoded secrets in code
- [ ] Secrets stored in secure vault (AWS Secrets Manager/Vault)
- [ ] Environment variables validated on startup

#### ✅ Backend Deployment
- [ ] Build production binary: `go build -o kawai-backend main.go`
- [ ] Test binary locally with production config (dry-run mode)
- [ ] Deploy to production server
- [ ] Configure systemd service for auto-restart
- [ ] Set up health check endpoint (`/health`)
- [ ] Configure log rotation (logrotate)
- [ ] Test API endpoints (curl/Postman)

---

### 5. Frontend Configuration

#### ✅ Environment Variables
```typescript
// .env.production
NEXT_PUBLIC_API_URL=https://api.kawai.network
NEXT_PUBLIC_CHAIN_ID=10143  // Monad Mainnet
NEXT_PUBLIC_RPC_URL=https://mainnet-rpc.monad.xyz
NEXT_PUBLIC_EXPLORER_URL=https://monadexplorer.com

// Contract Addresses (FILL AFTER DEPLOYMENT)
NEXT_PUBLIC_TOKEN_ADDRESS=0x...
NEXT_PUBLIC_MINING_DISTRIBUTOR=0x...
NEXT_PUBLIC_CASHBACK_DISTRIBUTOR=0x...
```

- [ ] All contract addresses configured
- [ ] RPC URL points to mainnet
- [ ] API URL points to production backend
- [ ] No testnet addresses in production build

#### ✅ Frontend Deployment
- [ ] Build production bundle: `npm run build`
- [ ] Test build locally: `npm run start`
- [ ] Deploy to hosting (Vercel/Netlify/Cloudflare Pages)
- [ ] Configure custom domain
- [ ] Enable CDN caching
- [ ] Test all pages load correctly
- [ ] Test wallet connection (MetaMask/WalletConnect)

---

### 6. Security Hardening

#### ✅ Smart Contracts
- [ ] **CRITICAL:** Security audit completed (recommended: CertiK/OpenZeppelin)
- [ ] Timelock on admin functions (24-48 hours)
- [ ] Multi-sig wallet for admin operations (Gnosis Safe)
- [ ] Emergency pause mechanism tested
- [ ] Rate limiting on claims (prevent spam)

#### ✅ Backend
- [ ] Rate limiting on API endpoints (100 req/min per IP)
- [ ] Input validation on all endpoints
- [ ] SQL injection prevention (use parameterized queries)
- [ ] CORS configured correctly (whitelist frontend domain)
- [ ] API authentication (JWT/API keys)
- [ ] DDoS protection (Cloudflare)
- [ ] Regular security updates (weekly)

#### ✅ Infrastructure
- [ ] SSH key-only access (no password)
- [ ] Firewall rules (allow only necessary ports)
- [ ] Fail2ban configured (block brute force)
- [ ] Regular backups (daily, tested restore)
- [ ] Intrusion detection (OSSEC/Wazuh)

---

### 7. Monitoring & Alerts

#### ✅ Application Monitoring
- [ ] Error tracking (Sentry/Rollbar)
- [ ] Performance monitoring (New Relic/Datadog)
- [ ] Uptime monitoring (UptimeRobot/Pingdom)
- [ ] Log aggregation (CloudWatch/Papertrail)
- [ ] Custom metrics (claim success rate, gas costs)

#### ✅ Blockchain Monitoring
- [ ] Contract balance monitoring (alert if low)
- [ ] Failed transaction alerts
- [ ] Gas price monitoring (alert if spike)
- [ ] Merkle root upload verification
- [ ] Claim success rate tracking

#### ✅ Alert Configuration
```yaml
Alerts:
  Critical (Immediate):
    - Backend down (>5 min)
    - Contract out of tokens
    - Failed merkle root upload
    - Security breach detected
    
  Warning (1 hour):
    - High error rate (>5%)
    - Low MON balance (<1 MON)
    - Slow API response (>2s)
    - High gas prices (>500 gwei)
    
  Info (Daily):
    - Daily claim summary
    - Gas cost report
    - User activity stats
```

- [ ] Telegram bot configured for alerts
- [ ] Email alerts configured
- [ ] PagerDuty/OpsGenie for on-call (optional)

---

### 8. Testing Strategy

#### ✅ Pre-Launch Testing (Mainnet)
1. **Smoke Test (1-2 hours)**
   - [ ] Deploy contracts with minimal allocations (1000 KAWAI)
   - [ ] Create 1 test user with real wallet
   - [ ] Inject small mining reward (10 KAWAI)
   - [ ] Run settlement
   - [ ] Test claim via UI
   - [ ] Verify KAWAI balance increased
   - [ ] Check gas costs

2. **Beta Test (1 week)**
   - [ ] Invite 10-20 beta users
   - [ ] Monitor all claims closely
   - [ ] Collect feedback on UX
   - [ ] Fix any issues found
   - [ ] Verify gas costs acceptable

3. **Soft Launch (2 weeks)**
   - [ ] Open to 100-200 users
   - [ ] Monitor system load
   - [ ] Optimize performance if needed
   - [ ] Prepare for full launch

#### ✅ Load Testing
- [ ] Test with 100 concurrent users
- [ ] Test with 1000 claims in 1 hour
- [ ] Test settlement with 500 users
- [ ] Verify database performance
- [ ] Check API response times

---

### 9. Operational Procedures

#### ✅ Weekly Settlement Process
```bash
# Every Monday 00:00 UTC (automated via cron)
1. Generate mining settlement
   go run cmd/reward-settlement/main.go generate --type mining

2. Generate cashback settlement
   go run cmd/reward-settlement/main.go generate --type cashback

3. Upload merkle roots
   go run cmd/reward-settlement/main.go upload --type mining
   go run cmd/reward-settlement/main.go upload --type cashback

4. Verify on-chain
   - Check merkle roots uploaded
   - Verify period advanced
   - Test 1 claim manually

5. Monitor claims
   - Watch for failed claims
   - Check gas costs
   - Verify user satisfaction
```

- [ ] Cron job configured for automated settlement
- [ ] Manual verification checklist created
- [ ] Rollback procedure documented
- [ ] Emergency contact list prepared

#### ✅ Incident Response
```yaml
Incident Types:
  1. Backend Down:
     - Check server status
     - Check logs for errors
     - Restart service
     - Notify users if >30 min
     
  2. Failed Claims:
     - Check merkle root uploaded
     - Verify proof validity
     - Check contract has tokens
     - Check user hasn't claimed
     
  3. Contract Issues:
     - Pause contract if critical
     - Investigate root cause
     - Prepare fix
     - Test on testnet
     - Deploy fix with timelock
     
  4. Security Breach:
     - Pause all contracts immediately
     - Investigate attack vector
     - Notify users
     - Prepare incident report
     - Implement fixes
```

- [ ] Incident response playbook created
- [ ] Team roles defined (who does what)
- [ ] Communication templates prepared
- [ ] Post-mortem process defined

---

### 10. Documentation

#### ✅ User Documentation
- [ ] How to claim rewards (step-by-step guide)
- [ ] FAQ (common issues & solutions)
- [ ] Video tutorials (optional)
- [ ] Troubleshooting guide
- [ ] Contact support information

#### ✅ Developer Documentation
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Contract ABI documentation
- [ ] Settlement process documentation
- [ ] Deployment guide (this document)
- [ ] Architecture diagrams

#### ✅ Operations Documentation
- [ ] Runbook for common tasks
- [ ] Incident response procedures
- [ ] Monitoring dashboard guide
- [ ] Backup/restore procedures
- [ ] Key rotation procedures

---

## 🚀 DEPLOYMENT DAY CHECKLIST

### Phase 1: Contract Deployment (1 hour)
- [ ] **T-60min:** Final code review
- [ ] **T-45min:** Deploy contracts to mainnet
- [ ] **T-30min:** Verify contracts on explorer
- [ ] **T-20min:** Grant MINTER_ROLE
- [ ] **T-10min:** Test all contract functions
- [ ] **T-0min:** Contracts ready ✅

### Phase 2: Backend Deployment (30 min)
- [ ] **T+0min:** Update .env with contract addresses
- [ ] **T+5min:** Run obfuscator-gen
- [ ] **T+10min:** Build production binary
- [ ] **T+15min:** Deploy to production server
- [ ] **T+20min:** Start backend service
- [ ] **T+25min:** Test API endpoints
- [ ] **T+30min:** Backend ready ✅

### Phase 3: Frontend Deployment (30 min)
- [ ] **T+30min:** Update frontend .env
- [ ] **T+35min:** Build production bundle
- [ ] **T+40min:** Deploy to hosting
- [ ] **T+45min:** Test all pages
- [ ] **T+50min:** Test wallet connection
- [ ] **T+60min:** Frontend ready ✅

### Phase 4: Smoke Test (30 min)
- [ ] **T+60min:** Create test user
- [ ] **T+65min:** Inject test reward
- [ ] **T+70min:** Run settlement
- [ ] **T+75min:** Test claim via UI
- [ ] **T+85min:** Verify success
- [ ] **T+90min:** System verified ✅

### Phase 5: Go Live (Immediate)
- [ ] **T+90min:** Announce launch
- [ ] **T+90min:** Monitor closely for 24 hours
- [ ] **T+24h:** Review metrics
- [ ] **T+24h:** Collect feedback
- [ ] **T+1week:** Evaluate & optimize

---

## 📊 SUCCESS METRICS

### Week 1 Targets
- [ ] 100+ users claimed rewards
- [ ] <1% claim failure rate
- [ ] <2s average claim time
- [ ] <$0.50 average gas cost
- [ ] 99.9% uptime
- [ ] 0 critical incidents

### Month 1 Targets
- [ ] 1000+ users claimed rewards
- [ ] <0.5% claim failure rate
- [ ] <1s average claim time
- [ ] 99.95% uptime
- [ ] User satisfaction >4.5/5

---

## 🆘 EMERGENCY CONTACTS

```yaml
Team:
  Backend Lead: [Name] - [Phone] - [Email]
  Frontend Lead: [Name] - [Phone] - [Email]
  DevOps: [Name] - [Phone] - [Email]
  Security: [Name] - [Phone] - [Email]

External:
  Hosting Support: [Provider] - [Support URL]
  RPC Provider: [Provider] - [Support URL]
  Security Auditor: [Company] - [Contact]

Escalation:
  Level 1: Team Lead (respond within 1 hour)
  Level 2: CTO (respond within 30 min)
  Level 3: CEO (critical only)
```

---

## 📝 POST-DEPLOYMENT

### Day 1
- [ ] Monitor all metrics closely
- [ ] Respond to user feedback
- [ ] Fix any minor issues
- [ ] Document any problems

### Week 1
- [ ] Review all incidents
- [ ] Optimize performance
- [ ] Update documentation
- [ ] Plan improvements

### Month 1
- [ ] Comprehensive review
- [ ] User survey
- [ ] Performance report
- [ ] Plan next features

---

## ✅ FINAL CHECKLIST

Before going live, verify:
- [ ] All contracts deployed & verified
- [ ] All environment variables configured
- [ ] Backend deployed & tested
- [ ] Frontend deployed & tested
- [ ] Monitoring & alerts configured
- [ ] Documentation complete
- [ ] Team trained on procedures
- [ ] Emergency contacts updated
- [ ] Backup & restore tested
- [ ] Security audit completed (recommended)
- [ ] Legal compliance verified
- [ ] User communication prepared
- [ ] Support channels ready

---

**Ready to deploy?** Follow this guide step-by-step. Don't skip any steps!

**Questions?** Review DEPLOYMENT.md for technical details.

**Issues?** Check CHANGELOG_CASHBACK_FIX.md for known issues & solutions.

---

**Good luck! 🚀**
