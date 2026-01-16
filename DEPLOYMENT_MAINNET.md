# Mainnet Deployment

**Status: ✅ LIVE (Promoted from Testnet Round 6)**

**Deployment Date:** 2026-01-16

---

## Current Mainnet Addresses

```
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
**RPC:** https://testnet-rpc.monad.xyz  
**Chain ID:** 10143

---

## System Status

✅ All 3 distributors deployed with pause mechanism  
✅ MINTER_ROLE granted to all distributors  
✅ Emergency pause tested and working  
✅ Monitoring active (Telegram alerts)  
✅ Production KV namespaces configured  

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

- Promoted from Testnet Round 6 (2026-01-15)
- All contracts tested and verified
- Ready for production use
- For new testnet, create fresh deployment
