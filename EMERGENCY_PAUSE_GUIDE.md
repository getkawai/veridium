# Emergency Pause Mechanism Guide

## Overview

All distributor contracts now include an emergency pause mechanism using OpenZeppelin's `Pausable` pattern. This allows the contract owner to immediately stop all claim operations in case of:
- Critical bug discovery
- Security incident
- Suspicious activity
- Contract upgrade preparation

## Affected Contracts

1. **MiningRewardDistributor** - Mining reward claims
2. **DepositCashbackDistributor** - Cashback claims  
3. **ReferralRewardDistributor** - Referral reward claims

## How It Works

### Paused State
When paused, all claim functions are blocked:
- `claimReward()` / `claimCashback()` / `claimRewards()` - Single period claims
- `claimMultiplePeriods()` - Batch claims

### Admin Functions Still Work
Even when paused, owner can still:
- Set Merkle roots (`setMerkleRoot()`)
- Advance periods (`advancePeriod()`)
- View contract state (all view functions)
- Unpause the contract

## Usage

### Pause Contract (Emergency)

```bash
# Mining Distributor
cast send $MINING_DISTRIBUTOR_ADDRESS "pause()" \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY

# Cashback Distributor
cast send $CASHBACK_DISTRIBUTOR_ADDRESS "pause()" \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY

# Referral Distributor
cast send $REFERRAL_DISTRIBUTOR_ADDRESS "pause()" \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY
```

### Check Pause Status

```bash
# Returns true if paused, false if not
cast call $DISTRIBUTOR_ADDRESS "paused()" --rpc-url $RPC_URL
```

### Unpause Contract (After Fix)

```bash
cast send $DISTRIBUTOR_ADDRESS "unpause()" \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY
```

## Emergency Response Procedure

### 1. Detect Issue
- Monitor Telegram alerts for claim failures
- Check contract events for suspicious patterns
- Review user reports

### 2. Assess Severity
**CRITICAL** - Pause immediately if:
- Funds at risk (exploit possible)
- Merkle proof vulnerability
- Reentrancy attack detected
- Unauthorized minting

**HIGH** - Pause within 1 hour if:
- Incorrect reward calculations
- Gas optimization issues causing failures
- Frontend/backend sync issues

**MEDIUM** - No pause needed if:
- Individual user claim issues
- Proof generation errors (backend fix)
- UI display issues

### 3. Pause Contracts

```bash
# Quick pause all distributors
make pause-all-distributors
```

Or manually:
```bash
# Pause each contract
cast send $MINING_DISTRIBUTOR_ADDRESS "pause()" --rpc-url $RPC_URL --private-key $PRIVATE_KEY
cast send $CASHBACK_DISTRIBUTOR_ADDRESS "pause()" --rpc-url $RPC_URL --private-key $PRIVATE_KEY
cast send $REFERRAL_DISTRIBUTOR_ADDRESS "pause()" --rpc-url $RPC_URL --private-key $PRIVATE_KEY
```

### 4. Notify Users
Send Telegram alert:
```
🚨 EMERGENCY MAINTENANCE 🚨
All reward claims temporarily paused.
Investigating issue. Your rewards are safe.
ETA: [time estimate]
```

### 5. Investigate & Fix
- Review contract events and transactions
- Identify root cause
- Test fix on testnet
- Prepare deployment if contract upgrade needed

### 6. Test Unpause
On testnet:
```bash
# Unpause
cast send $DISTRIBUTOR_ADDRESS "unpause()" --rpc-url $TESTNET_RPC --private-key $TEST_KEY

# Test claim
cast send $DISTRIBUTOR_ADDRESS "claimReward(...)" --rpc-url $TESTNET_RPC --private-key $TEST_KEY
```

### 7. Unpause Production
```bash
# Unpause all distributors
make unpause-all-distributors
```

### 8. Monitor
- Watch first 10-20 claims closely
- Check Telegram alerts
- Verify reward amounts correct
- Monitor gas usage

## Testing

Run pause mechanism tests:
```bash
cd contracts
forge test --match-path test/PauseTest.t.sol -vv
```

Expected output:
```
✓ testMiningPause() - Paused contract blocks claims
✓ testMiningUnpause() - Unpaused contract allows claims
✓ testMiningOnlyOwnerCanPause() - Non-owner cannot pause
✓ testMiningOnlyOwnerCanUnpause() - Non-owner cannot unpause
(+ 8 more tests for Cashback and Referral distributors)
```

## Makefile Commands

Add to `Makefile`:

```makefile
# Pause all distributors (emergency)
pause-all-distributors:
	@echo "🚨 PAUSING ALL DISTRIBUTORS..."
	cast send $(MINING_DISTRIBUTOR_ADDRESS) "pause()" --rpc-url $(RPC_URL) --private-key $(PRIVATE_KEY)
	cast send $(CASHBACK_DISTRIBUTOR_ADDRESS) "pause()" --rpc-url $(RPC_URL) --private-key $(PRIVATE_KEY)
	cast send $(REFERRAL_DISTRIBUTOR_ADDRESS) "pause()" --rpc-url $(RPC_URL) --private-key $(PRIVATE_KEY)
	@echo "✅ All distributors paused"

# Unpause all distributors
unpause-all-distributors:
	@echo "🔓 UNPAUSING ALL DISTRIBUTORS..."
	cast send $(MINING_DISTRIBUTOR_ADDRESS) "unpause()" --rpc-url $(RPC_URL) --private-key $(PRIVATE_KEY)
	cast send $(CASHBACK_DISTRIBUTOR_ADDRESS) "unpause()" --rpc-url $(RPC_URL) --private-key $(PRIVATE_KEY)
	cast send $(REFERRAL_DISTRIBUTOR_ADDRESS) "unpause()" --rpc-url $(RPC_URL) --private-key $(PRIVATE_KEY)
	@echo "✅ All distributors unpaused"

# Check pause status
check-pause-status:
	@echo "Mining Distributor paused: $$(cast call $(MINING_DISTRIBUTOR_ADDRESS) 'paused()' --rpc-url $(RPC_URL))"
	@echo "Cashback Distributor paused: $$(cast call $(CASHBACK_DISTRIBUTOR_ADDRESS) 'paused()' --rpc-url $(RPC_URL))"
	@echo "Referral Distributor paused: $$(cast call $(REFERRAL_DISTRIBUTOR_ADDRESS) 'paused()' --rpc-url $(RPC_URL))"
```

## Security Considerations

### Access Control
- Only contract owner can pause/unpause
- Owner private key must be secured (hardware wallet/KMS)
- Consider multi-sig for production

### User Impact
- Users cannot claim while paused
- Rewards are NOT lost - just delayed
- No funds at risk during pause

### Gas Costs
- Pause: ~30,000 gas (~0.00003 MON)
- Unpause: ~30,000 gas (~0.00003 MON)
- Negligible cost for emergency response

### Limitations
- Cannot pause individual users (all or nothing)
- Cannot pause admin functions (by design)
- Cannot pause token transfers (only claims)

## Production Checklist

Before mainnet deployment:

- [x] Pausable imported in all distributors
- [x] `whenNotPaused` modifier on all claim functions
- [x] `pause()` and `unpause()` functions added
- [x] Only owner can pause/unpause
- [x] Tests written and passing
- [ ] Makefile commands added
- [ ] Emergency response procedure documented
- [ ] Team trained on pause procedure
- [ ] Hardware wallet/KMS configured for owner key
- [ ] Telegram alert integration tested
- [ ] Testnet pause/unpause tested
- [ ] Multi-sig considered (optional)

## Related Files

- `contracts/contracts/MiningRewardDistributor.sol` - Mining pause implementation
- `contracts/contracts/DepositCashbackDistributor.sol` - Cashback pause implementation
- `contracts/contracts/ReferralRewardDistributor.sol` - Referral pause implementation
- `contracts/test/PauseTest.t.sol` - Pause mechanism tests
- `PRODUCTION_CHECKLIST.md` - Updated with pause mechanism item

## Support

If emergency pause needed:
1. Contact: [Your emergency contact]
2. Telegram: [Your Telegram alert channel]
3. Backup owner key location: [Secure location]
