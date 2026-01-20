# Mainnet Production Checklist

## 🎯 Priority: Production-Ready for Mainnet

This document outlines all necessary steps to ensure the application is production-ready for Monad Mainnet with USDC.

---

## ✅ COMPLETED ITEMS

### 1. Smart Contract Updates
- [x] Updated `PaymentVault.sol` to use generic `stablecoin` variable
- [x] Compiled contracts with `forge build`
- [x] Generated Go bindings with `make contracts-bindings`

### 2. Environment Configuration
- [x] Updated `.env.mainnet` with USDC address: `0x754704bc059f8c67012fed69bc8a327a5aafb603`
- [x] Added clear documentation for testnet vs mainnet
- [x] Updated `DEPLOYMENT.md` with network-specific addresses

### 3. Backend Code Safety
- [x] Added environment check to `MintTestTokens()` to prevent mainnet usage
- [x] Updated all comments to use "stablecoin" terminology
- [x] Added compatibility notes in code comments

### 4. Documentation
- [x] Created `STABLECOIN_SUPPORT.md` with comprehensive guide
- [x] Created `USDT_TO_STABLECOIN_MIGRATION.md` tracking all changes
- [x] Documented functions that work vs don't work on mainnet

---

## 🔴 CRITICAL ITEMS TO COMPLETE

### 1. Environment Detection System

**Problem**: Current code uses hardcoded RPC URL check
**Solution**: Create proper environment detection

**Action Items**:
- [ ] Add `ENVIRONMENT` variable to `.env` (values: `testnet` | `mainnet`)
- [ ] Create `pkg/config/environment.go` with environment detection
- [ ] Update all environment checks to use centralized config
- [ ] Add runtime validation on startup

**Files to Update**:
- `.env` - Add `ENVIRONMENT=testnet`
- `.env.mainnet` - Add `ENVIRONMENT=mainnet`
- `internal/services/deai_service.go` - Use config instead of hardcoded check
- `pkg/blockchain/revenue_settlement.go` - Add environment validation

### 2. Frontend Safety Guards

**Problem**: Frontend can still call `MintTestTokens()` on mainnet
**Solution**: Hide/disable test functions in production UI

**Action Items**:
- [ ] Update frontend to check environment from backend config
- [ ] Hide "Mint Test Tokens" button on mainnet
- [ ] Show appropriate message: "On mainnet, acquire USDC via exchanges/bridges"
- [ ] Add environment indicator in UI (Testnet badge vs Mainnet badge)

**Files to Update**:
- `frontend/src/components/*` - Add environment checks
- `frontend/src/config/*` - Add environment detection
- Regenerate TypeScript bindings if needed

### 3. Deployment Scripts

**Problem**: Need separate deployment for testnet vs mainnet
**Solution**: Create environment-specific deployment scripts

**Action Items**:
- [ ] Create `scripts/deploy-testnet.sh`
- [ ] Create `scripts/deploy-mainnet.sh`
- [ ] Add pre-deployment validation checks
- [ ] Document deployment process for each environment

### 4. Contract Deployment

**Problem**: PaymentVault needs to be deployed with USDC address on mainnet
**Solution**: Deploy contracts to mainnet with correct addresses

**Action Items**:
- [ ] Review all contract constructor parameters
- [ ] Deploy `PaymentVault` to mainnet with USDC address
- [ ] Deploy all distributor contracts to mainnet
- [ ] Update `.env.mainnet` with deployed addresses
- [ ] Run `make constants-generate` to update Go constants
- [ ] Verify all contracts on MonadScan

**Contracts to Deploy**:
- [ ] `PaymentVault` (with USDC address)
- [ ] `MiningRewardDistributor`
- [ ] `CashbackDistributor`
- [ ] `ReferralDistributor`
- [ ] `KawaiToken` (if not already deployed)

### 5. Testing on Mainnet

**Problem**: Need to verify all functionality works with real USDC
**Solution**: Comprehensive mainnet testing

**Action Items**:
- [ ] Acquire test USDC on Monad Mainnet (small amount)
- [ ] Test deposit flow with real USDC
- [ ] Test approval flow
- [ ] Test balance checking
- [ ] Test reward claiming (if applicable)
- [ ] Monitor gas costs
- [ ] Test error handling

### 6. Monitoring & Alerts

**Problem**: Need to monitor mainnet operations
**Solution**: Set up monitoring and alerting

**Action Items**:
- [ ] Set up transaction monitoring
- [ ] Add error logging for failed transactions
- [ ] Monitor USDC balance in PaymentVault
- [ ] Set up alerts for low balance
- [ ] Monitor gas prices
- [ ] Track user deposits and withdrawals

---

## 🟡 RECOMMENDED IMPROVEMENTS

### 1. Graceful Degradation

**Add fallback mechanisms**:
- [ ] Handle RPC failures gracefully
- [ ] Add retry logic for failed transactions
- [ ] Cache balance data to reduce RPC calls
- [ ] Add circuit breaker for repeated failures

### 2. User Experience

**Improve UX for mainnet users**:
- [ ] Add clear instructions for acquiring USDC
- [ ] Show bridge links (e.g., to bridge USDC from Ethereum)
- [ ] Display current USDC balance prominently
- [ ] Add transaction history
- [ ] Show pending transactions

### 3. Security Hardening

**Additional security measures**:
- [ ] Add rate limiting for deposits
- [ ] Implement maximum deposit limits (if needed)
- [ ] Add multi-sig for admin operations
- [ ] Audit all smart contracts before mainnet launch
- [ ] Set up emergency pause mechanism

### 4. Documentation

**User-facing documentation**:
- [ ] Create user guide for mainnet deposits
- [ ] Document how to get USDC on Monad
- [ ] Create FAQ for common issues
- [ ] Add troubleshooting guide
- [ ] Document gas fee expectations

---

## 📋 PRE-LAUNCH CHECKLIST

Before launching on mainnet, verify:

### Configuration
- [ ] `.env.mainnet` has correct USDC address
- [ ] All contract addresses are verified on MonadScan
- [ ] RPC URL points to mainnet
- [ ] Private keys are secured (use hardware wallet/KMS)
- [ ] Environment variable is set to `mainnet`

### Smart Contracts
- [ ] All contracts deployed to mainnet
- [ ] Contract addresses updated in `.env.mainnet`
- [ ] Contracts verified on MonadScan
- [ ] Admin roles assigned correctly
- [ ] Emergency pause tested

### Backend
- [ ] Build passes: `go build -o /dev/null .`
- [ ] All tests pass: `make test`
- [ ] Constants regenerated: `make constants-generate`
- [ ] No hardcoded testnet addresses in code
- [ ] Environment detection working correctly

### Frontend
- [ ] TypeScript bindings regenerated: `make bindings-generate`
- [ ] Test functions hidden on mainnet
- [ ] Environment indicator visible
- [ ] USDC acquisition instructions visible
- [ ] Build passes: `cd frontend && npm run build`

### Testing
- [ ] Deposit flow tested with real USDC
- [ ] Withdrawal flow tested
- [ ] Balance checking works
- [ ] Error messages are user-friendly
- [ ] Gas estimation accurate

### Monitoring
- [ ] Logging configured
- [ ] Error tracking set up
- [ ] Transaction monitoring active
- [ ] Alerts configured
- [ ] Backup RPC endpoints configured

---

## 🚀 DEPLOYMENT SEQUENCE

### Phase 1: Preparation (Before Mainnet Launch)
1. Complete all critical items above
2. Deploy contracts to mainnet
3. Update configuration files
4. Regenerate all bindings and constants
5. Run full test suite

### Phase 2: Soft Launch (Limited Users)
1. Deploy backend to production
2. Deploy frontend to production
3. Test with small group of users
4. Monitor for issues
5. Gather feedback

### Phase 3: Full Launch
1. Address any issues from soft launch
2. Update documentation
3. Announce mainnet availability
4. Monitor closely for first 48 hours

---

## 🆘 ROLLBACK PLAN

If critical issues are found on mainnet:

1. **Immediate Actions**:
   - Pause contracts if emergency pause is available
   - Switch frontend to maintenance mode
   - Stop accepting new deposits

2. **Investigation**:
   - Review logs and error messages
   - Identify root cause
   - Assess impact on users

3. **Resolution**:
   - Fix issue in code
   - Test fix thoroughly on testnet
   - Deploy fix to mainnet
   - Resume operations

4. **Communication**:
   - Notify users of issue
   - Provide timeline for resolution
   - Update status page

---

## 📞 SUPPORT CONTACTS

- **Smart Contract Issues**: [Contract Developer]
- **Backend Issues**: [Backend Developer]
- **Frontend Issues**: [Frontend Developer]
- **Infrastructure**: [DevOps]
- **Emergency Contact**: [Project Lead]

---

## 📝 NOTES

- Always test on testnet first before deploying to mainnet
- Keep private keys secure - never commit to git
- Monitor gas prices - Monad should have low fees but verify
- USDC has 6 decimals, not 18 like most tokens
- Circle's USDC is the official USDC on Monad Mainnet
- Bridge from Ethereum may take time - inform users

---

**Last Updated**: January 21, 2026  
**Status**: 🟡 In Progress - Critical items pending  
**Target Launch**: [Set target date]
