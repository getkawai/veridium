# Production Deployment - Quick Checklist

**Use this for quick reference. See PRODUCTION_DEPLOYMENT_GUIDE.md for details.**

---

## 🔴 CRITICAL (Must Do)

### Before Deployment
- [ ] Security audit completed (or accepted risk)
- [ ] Admin wallet secured (hardware wallet/KMS)
- [ ] Private keys NEVER in git
- [ ] Production RPC endpoints configured
- [ ] Cloudflare KV production namespaces created
- [ ] Monitoring & alerts configured (✅ Integrated - test with `make test-telegram-alert`)
- [ ] Backup strategy implemented

### During Deployment
- [ ] Deploy contracts to mainnet
- [ ] Verify contracts on explorer
- [ ] Grant MINTER_ROLE to distributors
- [ ] Test all contract functions
- [ ] Update backend .env with addresses
- [ ] Deploy backend to production
- [ ] Deploy frontend to production
- [ ] Run smoke test with real wallet

### After Deployment
- [ ] Monitor for 24 hours continuously
- [ ] Test claims with real users
- [ ] Verify gas costs acceptable
- [ ] Check all alerts working
- [ ] Document any issues

---

## 🟡 IMPORTANT (Should Do)

### Infrastructure
- [ ] SSL certificate installed
- [ ] Domain configured
- [ ] Firewall rules set
- [ ] Log aggregation configured
- [ ] Rate limiting enabled
- [ ] DDoS protection (Cloudflare)

### Security
- [ ] Multi-sig for admin (Gnosis Safe)
- [ ] Timelock on admin functions
- [x] Emergency pause mechanism tested (see EMERGENCY_PAUSE_GUIDE.md)
- [ ] Regular security updates
- [ ] Incident response plan

### Operations
- [ ] Automated settlement (cron)
- [ ] Health check endpoint
- [ ] Error tracking (Sentry)
- [ ] Performance monitoring
- [ ] User documentation

---

## 🟢 NICE TO HAVE (Optional)

- [ ] Load testing (100+ concurrent users)
- [ ] Beta testing (1-2 weeks)
- [ ] Video tutorials
- [ ] API documentation (Swagger)
- [ ] Architecture diagrams
- [ ] User satisfaction survey

---

## 📋 Deployment Day Timeline

| Time | Task | Duration |
|------|------|----------|
| T-60min | Final review | 15 min |
| T-45min | Deploy contracts | 30 min |
| T-15min | Deploy backend | 15 min |
| T-0min | Deploy frontend | 15 min |
| T+15min | Smoke test | 15 min |
| T+30min | Go live | - |
| T+30min to T+24h | Monitor closely | 24 hours |

---

## 🆘 Emergency Procedures

### Backend Down
1. Check server status
2. Check logs: `journalctl -u kawai-backend -n 100`
3. Restart: `systemctl restart kawai-backend`
4. If still down, rollback to previous version

### Failed Claims
1. Check merkle root: `cast call <CONTRACT> "periodMerkleRoots(uint256)" <PERIOD>`
2. Verify proof: `go run cmd/dev/verify-cashback-proof/ <USER>`
3. Check contract balance: `cast call <TOKEN> "balanceOf(address)" <CONTRACT>`
4. Check user claimed: `cast call <CONTRACT> "hasClaimed(uint256,address)" <PERIOD> <USER>`

### Contract Issues
1. **PAUSE IMMEDIATELY** if critical
2. Investigate root cause
3. Test fix on testnet
4. Deploy fix with timelock
5. Notify users

---

## 📞 Quick Contacts

- **Backend Issues:** [Name] - [Phone]
- **Contract Issues:** [Name] - [Phone]
- **Security Issues:** [Name] - [Phone]
- **Hosting Support:** [Provider Support URL]

---

## ✅ Pre-Launch Verification

Run these commands before going live:

```bash
# 1. Verify contracts deployed
cast call <TOKEN_ADDRESS> "name()(string)"
cast call <MINING_DISTRIBUTOR> "currentPeriod()(uint256)"
cast call <CASHBACK_DISTRIBUTOR> "currentPeriod()(uint256)"

# 2. Verify MINTER_ROLE granted
cast call <TOKEN> "hasRole(bytes32,address)(bool)" \
  $(cast keccak "MINTER_ROLE") <MINING_DISTRIBUTOR>

# 3. Test backend API
curl https://api.kawai.network/health
curl https://api.kawai.network/api/cashback/stats?address=0x...

# 4. Test frontend
open https://app.kawai.network
# - Connect wallet
# - Check rewards page loads
# - Verify contract addresses correct
```

---

## 📊 Success Criteria

### Day 1
- [ ] 10+ successful claims
- [ ] 0 critical errors
- [ ] <2s API response time
- [ ] 100% uptime

### Week 1
- [ ] 100+ successful claims
- [ ] <1% error rate
- [ ] <$0.50 avg gas cost
- [ ] 99.9% uptime

### Month 1
- [ ] 1000+ successful claims
- [ ] <0.5% error rate
- [ ] User satisfaction >4.5/5
- [ ] 99.95% uptime

---

**Ready?** Follow PRODUCTION_DEPLOYMENT_GUIDE.md step-by-step.

**Questions?** Check the detailed guide for explanations.

**Issues?** Use emergency procedures above.
