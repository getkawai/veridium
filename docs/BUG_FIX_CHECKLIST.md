# Bug Fix Implementation Checklist

## 🎯 Priority 1: Critical Fixes (Do First)

### ✅ Bug #1: Nil Pointer Dereference
- [x] Verify fix in `GetVaultBalance()`
- [x] Verify fix in `DepositToVault()`
- [x] Run test: `go test -v -run TestNilPointerBug ./internal/services/`
- [x] Status: **COMPLETE** ✅

### ✅ Bug #2: Balance Race Condition
- [x] **Step 1:** Review current implementation
  ```bash
  # Find all balance operations
  grep -rn "DeductBalance\|AddBalance" internal/services/
  ```

- [x] **Step 2:** Update `deposit_sync_service.go`
  - [x] Line 157: Replaced `AddBalance` with `AddBalanceAtomic`
  - [x] Added comment about atomic operation
  - [x] Status: **COMPLETE** ✅

- [x] **Step 3:** Update `gateway/handler.go`
  - [x] Line 222: Replaced `DeductBalance` with `DeductBalanceAtomic`
  - [x] Added comment about race condition prevention
  - [x] Status: **COMPLETE** ✅

- [x] **Step 4:** Update `payment_event_listener.go`
  - [x] Line 98: Replaced `AddBalance` with `AddBalanceAtomic`
  - [x] Added comment about atomic operation
  - [x] Status: **COMPLETE** ✅

- [ ] **Step 5:** Test atomic operations
  ```bash
  go test -v ./pkg/store/balance_atomic_test.go
  go test -race ./pkg/store/
  ```

- [ ] **Step 6:** Integration testing
  - [ ] Test with 10 concurrent API requests
  - [ ] Test with insufficient balance
  - [ ] Test with deposit during spending
  - [ ] Verify no negative balances

- [ ] **Step 7:** Deploy and monitor
  - [ ] Deploy to staging
  - [ ] Monitor for errors
  - [ ] Check balance consistency
  - [ ] Deploy to production

---

## 🎯 Priority 2: Medium Fixes

### ⚠️ Bug #3: Marketplace Silent Failures

- [ ] **Step 1:** Add retry logic to `storePartialTradeRecord()`
  ```go
  func (s *TradeService) storePartialTradeRecordWithRetry(...) error {
      maxRetries := 3
      for attempt := 0; attempt < maxRetries; attempt++ {
          err := s.storePartialTradeRecord(...)
          if err == nil {
              return nil
          }
          time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
      }
      return fmt.Errorf("failed after %d retries", maxRetries)
  }
  ```

- [ ] **Step 2:** Add circuit breaker for KV operations
  - [ ] Install circuit breaker library: `go get github.com/sony/gobreaker`
  - [ ] Wrap KV operations with circuit breaker
  - [ ] Configure thresholds (5 failures in 10 seconds)

- [ ] **Step 3:** Add admin alerts
  - [ ] Log critical failures to Sentry
  - [ ] Send webhook notification on persistent failures
  - [ ] Create dashboard for monitoring

- [ ] **Step 4:** Test failure scenarios
  - [ ] Simulate KV store timeout
  - [ ] Simulate network errors
  - [ ] Verify retry logic works
  - [ ] Verify circuit breaker opens

---

## 🎯 Priority 3: Race Condition Verification

### ⚠️ Bug #4: WalletService Race Conditions

- [ ] **Step 1:** Audit all WalletService methods
  ```bash
  # Check mutex usage
  grep -A 10 "func (s \*WalletService)" internal/services/wallet_service.go
  ```

- [ ] **Step 2:** Verify mutex protection
  - [ ] `GetStatus()` - Uses `RLock()` ✅
  - [ ] `LockWallet()` - Uses `Lock()` ✅
  - [ ] `UnlockWallet()` - Check mutex usage
  - [ ] `SwitchWallet()` - Check mutex usage
  - [ ] `GetCurrentAddress()` - Check mutex usage
  - [ ] `SignMessage()` - Check mutex usage
  - [ ] `getTransactOpts()` - Check mutex usage

- [ ] **Step 3:** Run race detector
  ```bash
  go test -race -run TestRaceCondition ./internal/services/
  ```

- [ ] **Step 4:** Fix any race conditions found
  - [ ] Add `mu.RLock()` for read operations
  - [ ] Add `mu.Lock()` for write operations
  - [ ] Ensure proper `defer mu.Unlock()`

---

## 🎯 Priority 4: Low Priority Fixes

### 🟡 Bug #5: Panic in MCP Gateway

- [ ] **Step 1:** Find all panic calls
  ```bash
  grep -rn "panic(" pkg/mcp-gateway/
  ```

- [ ] **Step 2:** Replace with error returns
  - [ ] `client/config.go:44`
  - [ ] `mcp/stdio.go:75`
  - [ ] `mcp/stdio.go:82`
  - [ ] `desktop/paths.go:17`

- [ ] **Step 3:** Update callers to handle errors
  - [ ] Add error checks
  - [ ] Add proper error messages
  - [ ] Test error paths

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] Run all tests: `go test ./...`
- [ ] Run with race detector: `go test -race ./...`
- [ ] Check coverage: `go test -cover ./...`

### Integration Tests
- [ ] Test concurrent balance operations
- [ ] Test KV store failures
- [ ] Test blockchain RPC failures
- [ ] Test wallet lock during transaction

### Load Tests
- [ ] 100 concurrent API requests
- [ ] 1000 concurrent balance operations
- [ ] Sustained load for 5 minutes
- [ ] Monitor memory usage

### Manual Tests
- [ ] Lock wallet → Try to get balance → Should show error
- [ ] Make 10 concurrent API requests with 500 USDT balance
- [ ] Deposit while spending simultaneously
- [ ] Cancel order while someone is buying

---

## 📊 Verification Checklist

### Before Deployment
- [ ] All tests passing
- [ ] No race conditions detected
- [ ] Code reviewed by team
- [ ] Documentation updated
- [ ] Changelog updated

### After Deployment (Staging)
- [ ] Monitor error rates
- [ ] Check balance consistency
- [ ] Verify no negative balances
- [ ] Test with real users
- [ ] Monitor for 24 hours

### After Deployment (Production)
- [ ] Gradual rollout (10% → 50% → 100%)
- [ ] Monitor Sentry for errors
- [ ] Check database for anomalies
- [ ] Verify financial transactions
- [ ] Monitor for 1 week

---

## 🚨 Rollback Plan

If issues are found after deployment:

1. **Immediate Actions:**
   - [ ] Revert to previous version
   - [ ] Notify team
   - [ ] Check for data corruption

2. **Investigation:**
   - [ ] Review error logs
   - [ ] Check Sentry reports
   - [ ] Analyze failed transactions

3. **Fix and Redeploy:**
   - [ ] Fix identified issues
   - [ ] Add more tests
   - [ ] Redeploy with caution

---

## 📝 Documentation Updates

- [ ] Update README.md with new testing commands
- [ ] Update API documentation
- [ ] Update deployment guide
- [ ] Add troubleshooting section
- [ ] Document known issues

---

## ✅ Sign-off

### Developer
- [ ] All fixes implemented
- [ ] All tests passing
- [ ] Code reviewed
- [ ] Documentation updated

**Signed:** _________________ **Date:** _________

### QA
- [ ] All tests executed
- [ ] No critical bugs found
- [ ] Performance acceptable
- [ ] Ready for deployment

**Signed:** _________________ **Date:** _________

### Tech Lead
- [ ] Code quality approved
- [ ] Architecture sound
- [ ] Security reviewed
- [ ] Approved for production

**Signed:** _________________ **Date:** _________

---

## 📞 Emergency Contacts

- **Tech Lead:** [Name] - [Contact]
- **DevOps:** [Name] - [Contact]
- **On-Call:** [Name] - [Contact]

---

## 📅 Timeline

| Task | Estimated Time | Actual Time | Status |
|------|---------------|-------------|--------|
| Atomic operations | 2-4 hours | ~1 hour | ✅ |
| Retry logic | 1-2 hours | | ⏳ |
| Testing | 2-3 hours | | ⏳ |
| Code review | 1 hour | | ⏳ |
| Deployment | 1 hour | | ⏳ |
| **Total** | **7-11 hours** | **~1 hour** | **In Progress** |

---

**Last Updated:** 2025-01-01  
**Next Review:** After implementation
