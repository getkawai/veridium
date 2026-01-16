# Mainnet Deployment

**Current Status: 🟡 PRODUCTION ON TESTNET**

**What this means:**
- Running on Monad Testnet (not true mainnet yet)
- Production environment (real users, real data)
- Waiting for Monad Mainnet launch

**Deployment Date:** 2026-01-16

---

## Current Production Addresses (Monad Testnet)

```
Network: Monad Testnet
RPC: https://testnet-rpc.monad.xyz
Chain ID: 10143

MockUSDT:                    0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc
KawaiToken:                  0xf68910e8d19047A309f989FFB515E44FBca5D31A
KAWAI_Distributor:           0x2B11e8385A859Ea75C77E05Bc0D9756A87017E92
USDT_Distributor:            0x896fB97f81ECBEfDBe29DCc3550aC984704932bF
PaymentVault:                0xDA94C8ac2a61eafBd47853EE22702BDCd45B6d93
OTCMarket:                   0x3E597D76B40004c3fC517C404037fD6F16C8fc34
MiningRewardDistributor:     0x1f78c7c472205F1720aAb66a565981561b5EBac0 ⭐ WITH PAUSE
DepositCashbackDistributor:  0x3d5Bfe788782A90ac124096296B45eaFFc43C79B ⭐ WITH PAUSE
ReferralRewardDistributor:   0x1c218602218745B20CE201948CaE836f8E94E111 ⭐ WITH PAUSE
```

**Admin:** 0x94D5C06229811c4816107005ff05259f229Eb07b

---

## Migration to True Mainnet (When Available)

### ✅ What Can Be Reused (No Changes Needed)
- Backend code
- Frontend code
- Smart contract code (just redeploy)
- Database/KV structure
- Monitoring setup
- Emergency procedures

### ❌ What Must Change

**1. Network Configuration**
```bash
# Current (Testnet)
MONAD_RPC_URL=https://testnet-rpc.monad.xyz
Chain ID: 10143

# True Mainnet (Future)
MONAD_RPC_URL=https://mainnet-rpc.monad.xyz
Chain ID: TBD (will be different)
```

**2. Contract Addresses**
- Must deploy all contracts to mainnet
- All addresses will be different
- Update .env with new addresses

**3. USDT Token**
```bash
# Current: MockUSDT (our test token)
USDT_TOKEN_ADDRESS=0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc

# Mainnet: Official USDT
USDT_TOKEN_ADDRESS=<official USDT on Monad Mainnet>
```

**4. Wallet Security**
```bash
# Current: Test wallet (OK for testnet)
# Mainnet: MUST use hardware wallet or KMS
```

**5. KV Namespaces**
- Create fresh production namespaces
- Cannot reuse testnet namespaces

**6. User Data**
- Cannot migrate directly (different chain)
- Users must claim on new mainnet
- Consider airdrop for existing users

### Deployment Steps to True Mainnet

```bash
# 1. Wait for Monad Mainnet launch announcement

# 2. Update configuration
MONAD_RPC_URL=https://mainnet-rpc.monad.xyz

# 3. Create fresh KV namespaces (10 namespaces)

# 4. Use hardware wallet
# Fund with real MON

# 5. Deploy contracts
make contracts-deploy-mainnet
make contracts-deploy-mining-mainnet
make contracts-deploy-cashback-mainnet
make contracts-deploy-referral-mainnet
make contracts-grant-minter-mining
make contracts-grant-minter-cashback
make contracts-grant-minter-referral

# 6. Update .env with new addresses
# 7. Regenerate constants
go run cmd/obfuscator-gen/main.go

# 8. Deploy backend & frontend
# 9. Test with small amounts
# 10. Announce migration to users
```

---

## Current System Status

✅ All 3 distributors deployed with pause mechanism  
✅ MINTER_ROLE granted to all distributors  
✅ Emergency pause tested and working  
✅ Monitoring active (Telegram alerts)  
✅ Production KV namespaces configured  
✅ Ready for production use on testnet

---

## Emergency Procedures

**Pause everything:**
```bash
make pause-all
```

**Check status:**
```bash
make pause-status
```

**Resume:**
```bash
make unpause-all
```

---

## Notes

- **Current:** Production environment on Monad Testnet
- **Future:** Will migrate to Monad Mainnet when available
- **Users:** Can use now, but tokens have no real value until mainnet
- **Timeline:** Depends on Monad Mainnet launch date
