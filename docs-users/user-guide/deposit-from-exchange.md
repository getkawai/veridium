# How to Deposit USDC from Exchange (Binance, Coinbase, etc.)

## ⚠️ Important: Network Requirement

**Kawai Desktop only accepts USDC on Monad Network.** When withdrawing from exchanges:

1. **Check if Monad Network is available** in the withdrawal network options
2. **If Monad is available:** Select it and withdraw directly to your Kawai wallet ✅
3. **If Monad is NOT available:** You need to bridge from another network (see options below)

## Why Monad Network?

- All Kawai smart contracts are deployed on Monad Mainnet
- USDC address on Monad: `0x754704bc059f8c67012fed69bc8a327a5aafb603`
- This is a different USDC contract than on Ethereum, BSC, or other networks
- Each blockchain has its own USDC contract - they are NOT interchangeable

## How to Get USDC on Monad Network

### Option 1: Direct Withdrawal (If Available) ⚡ FASTEST

**Check your exchange withdrawal page:**
1. Go to Withdraw USDC
2. Look for "Monad" or "Monad Network" in the network selection
3. If available:
   - Select **Monad Network**
   - Enter your Kawai wallet address
   - Enter amount
   - Confirm withdrawal
   - Done! ✅

**Exchanges that may support Monad:**
- Check Binance, Coinbase, OKX, Bybit, Gate.io, etc.
- Support is being added gradually as Monad adoption grows

**Cost:** Only exchange withdrawal fee (~$1-5)  
**Time:** 5-15 minutes

---

### Option 2: Bridge from Another Network (If Monad Not Available)

**Step 1: Withdraw USDC from Exchange**
1. Go to your exchange (Binance, Coinbase, etc.)
2. Withdraw USDC to your personal wallet (MetaMask, Trust Wallet, etc.)
3. Choose a network that has a bridge to Monad:
   - Ethereum (most common)
   - Arbitrum (cheaper gas)
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

**Time:** 30-60 minutes total

---

### Option 3: Buy MON Token First

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

### Option 4: Direct Fiat On-Ramp (If Available)

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

### ❌ DON'T: Select Wrong Network
```
Binance → Withdraw USDC
  → Select "Ethereum Network" (WRONG!)
  → Send to Kawai Wallet Address
  = FUNDS LOST! (Wrong network)
```

### ✅ DO: Select Monad Network (If Available)
```
Binance → Withdraw USDC
  → Select "Monad Network" ✅
  → Send to Kawai Wallet Address
  = SUCCESS!
```

### ✅ DO: Bridge If Monad Not Available
```
Binance → Withdraw USDC (Ethereum)
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

**Q: Does my exchange support Monad Network withdrawals?**
A: Check your exchange's withdrawal page. Look for "Monad" in the network selection dropdown. If not available, you'll need to bridge from another network.

**Q: Which exchanges support Monad withdrawals?**
A: Support is being added gradually. Check: Binance, Coinbase, OKX, Bybit, Gate.io, KuCoin, Kraken. If your exchange doesn't support it yet, use the bridge method.

**Q: Can I use USDT instead of USDC?**
A: On mainnet, only USDC is supported. On testnet, we use MockUSDT for testing.

**Q: Why doesn't my exchange have Monad as an option?**
A: Monad is a relatively new blockchain (launched Nov 2024). Major exchanges add support over time as adoption grows. Use the bridge method in the meantime.

**Q: How long does bridging take?**
A: Usually 5-30 minutes depending on the source network and bridge congestion.

**Q: What if I sent USDC to the wrong network?**
A: Unfortunately, funds sent to the wrong network cannot be recovered by Kawai. Always verify you selected **Monad Network** before confirming the withdrawal!

**Q: How do I know if I selected the right network?**
A: Before confirming withdrawal:
1. Check the network name shows "Monad" or "Monad Network"
2. Verify the chain ID is 143 (Monad Mainnet)
3. Double-check your Kawai wallet address
4. Start with a small test amount first

**Q: Is there a minimum deposit amount?**
A: No minimum, but consider fees. If bridging, depositing less than $10 may not be economical due to bridge costs. Direct Monad withdrawals have lower fees.

---

## Summary

**The Golden Rule:** 
> Always ensure you select **Monad Network** when withdrawing USDC!

**Recommended Paths:**

**Path 1 (Fastest):** If your exchange supports Monad
```
Exchange → Select Monad Network → Withdraw to Kawai Wallet
```

**Path 2 (If Monad not available):** Bridge from another network
```
Exchange → Personal Wallet → Bridge to Monad → Kawai Desktop
```

**Before Every Withdrawal:**
1. ✅ Verify network is **Monad** (Chain ID: 143)
2. ✅ Double-check wallet address
3. ✅ Test with small amount first
4. ✅ Confirm you have MON for gas (~$0.50 worth)

This ensures your funds arrive safely! 🎯
