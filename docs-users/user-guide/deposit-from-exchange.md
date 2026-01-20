# How to Deposit USDC from Exchange (Binance, Coinbase, etc.)

## ⚠️ Important: Network Requirement

**Kawai Desktop only accepts USDC on Monad Network.** You cannot directly deposit USDC from exchanges like Binance or Coinbase because they don't support Monad Network withdrawals yet.

## Why Monad Network?

- All Kawai smart contracts are deployed on Monad Mainnet
- USDC address on Monad: `0x754704bc059f8c67012fed69bc8a327a5aafb603`
- This is a different USDC contract than on Ethereum, BSC, or other networks

## How to Get USDC on Monad Network

### Option 1: Bridge from Another Network (Recommended)

**Step 1: Withdraw USDC from Exchange**
1. Go to your exchange (Binance, Coinbase, etc.)
2. Withdraw USDC to your personal wallet (MetaMask, Trust Wallet, etc.)
3. Choose a network that has a bridge to Monad:
   - Ethereum (most common)
   - Arbitrum
   - Optimism
   - Base

**Step 2: Bridge to Monad**
1. Visit the official Monad Bridge: [https://bridge.monad.xyz](https://bridge.monad.xyz)
2. Connect your wallet
3. Select source network (e.g., Ethereum)
4. Select destination: Monad
5. Enter USDC amount
6. Confirm and wait for bridge completion (usually 5-30 minutes)

**Step 3: Deposit to Kawai**
1. Open Kawai Desktop
2. Go to Wallet → Deposit
3. Enter amount
4. Confirm transaction

**Total Cost:**
- Exchange withdrawal fee: ~$1-5 (varies by exchange)
- Bridge fee: ~$5-20 (depends on source network gas)
- Kawai deposit gas: ~$0.01 (very cheap on Monad!)

---

### Option 2: Buy MON Token First

If you can withdraw MON (Monad's native token) from your exchange:

**Step 1: Withdraw MON**
1. Withdraw MON from exchange to your wallet
2. Make sure to select Monad Network

**Step 2: Swap MON to USDC**
1. Visit a Monad DEX (Decentralized Exchange)
2. Swap MON → USDC
3. You now have USDC on Monad Network

**Step 3: Deposit to Kawai**
1. Open Kawai Desktop
2. Go to Wallet → Deposit
3. Enter amount
4. Confirm transaction

---

### Option 3: Direct Fiat On-Ramp (If Available)

Some on-ramp services may support Monad Network directly:

**Supported Services (Check Availability):**
- Moonpay
- Transak
- Ramp Network

**Steps:**
1. Visit the on-ramp service
2. Select "Buy USDC"
3. Choose network: **Monad**
4. Enter amount and payment method
5. Complete KYC if required
6. Receive USDC directly on Monad Network
7. Deposit to Kawai Desktop

---

## Common Mistakes to Avoid

### ❌ DON'T: Send USDC from Binance Directly
```
Binance (USDC on Ethereum) 
  → ❌ Kawai Wallet Address
  = FUNDS LOST! (Wrong network)
```

### ✅ DO: Bridge First
```
Binance (USDC on Ethereum)
  → Your Wallet (Ethereum)
  → Bridge to Monad
  → Kawai Wallet (Monad)
  = SUCCESS!
```

---

## Network Verification Checklist

Before depositing, verify:

1. ✅ Your wallet is connected to **Monad Network**
2. ✅ You have USDC on **Monad Network** (not Ethereum/BSC/etc.)
3. ✅ You have some MON for gas fees (~$0.50 worth is enough)
4. ✅ Kawai Desktop shows "Monad Mainnet" in the network selector

---

## Need Help?

If you're stuck or need assistance:

1. **Check Network Settings:**
   - Open Kawai Desktop
   - Go to Settings → Network
   - Verify it shows "Monad Mainnet"

2. **Verify Your USDC:**
   - Check your wallet on [MonadScan](https://monadexplorer.com)
   - Search your address
   - Confirm USDC balance shows up

3. **Contact Support:**
   - Discord: [Join our server](https://discord.gg/kawai)
   - Telegram: [@kawai_support](https://t.me/kawai_support)
   - Email: support@kawai.network

---

## FAQ

**Q: Can I use USDT instead of USDC?**
A: On mainnet, only USDC is supported. On testnet, we use MockUSDT for testing.

**Q: Why doesn't Binance support Monad withdrawals?**
A: Monad is a relatively new blockchain. Major exchanges will add support over time as adoption grows.

**Q: How long does bridging take?**
A: Usually 5-30 minutes depending on the source network and bridge congestion.

**Q: What if I sent USDC to the wrong network?**
A: Unfortunately, funds sent to the wrong network cannot be recovered by Kawai. Always verify the network before sending!

**Q: Is there a minimum deposit amount?**
A: No minimum, but consider gas fees. Depositing less than $10 may not be economical due to bridge costs.

---

## Summary

**The Golden Rule:** 
> Always ensure your USDC is on **Monad Network** before depositing to Kawai Desktop!

**Recommended Path:**
```
Exchange → Personal Wallet → Bridge to Monad → Kawai Desktop
```

This extra step ensures your funds are safe and arrive correctly! 🎯
