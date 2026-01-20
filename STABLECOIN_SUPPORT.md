# Stablecoin Support

**Last Updated:** January 21, 2026

## Overview

Kawai DeAI Network supports stablecoin deposits for AI service payments. The specific stablecoin used depends on the network:

## Supported Stablecoins

| Network | Stablecoin | Address | Type |
|---------|-----------|---------|------|
| **Monad Testnet** | MockUSDT | `0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc` | Test token (anyone can mint) |
| **Monad Mainnet** | USDC | `0x754704bc059f8c67012fed69bc8a327a5aafb603` | Circle: USDC Token |

## Why USDC on Mainnet?

**USDT is not yet available on Monad Mainnet.** As of January 2026:
- Monad Mainnet launched on November 24, 2025
- Tether (USDT issuer) has not deployed USDT to Monad yet
- USDC by Circle is the primary stablecoin available on Monad

**USDC is an excellent alternative:**
- ✅ Same $1 peg as USDT
- ✅ Issued by Circle (reputable company)
- ✅ Widely used in DeFi
- ✅ High liquidity on Monad
- ✅ Integrated with Uniswap and other protocols

## Technical Implementation

### Variable Naming

For backward compatibility, the environment variable remains `USDT_TOKEN_ADDRESS`:

```bash
# .env.testnet
USDT_TOKEN_ADDRESS=0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc  # MockUSDT

# .env.mainnet
USDT_TOKEN_ADDRESS=0x754704bc059f8c67012fed69bc8a327a5aafb603  # USDC
```

### Smart Contracts

All contracts are **stablecoin-agnostic**:
- `PaymentVault.sol` accepts any ERC-20 token
- `DepositCashbackDistributor.sol` works with any stablecoin
- No hardcoded assumptions about token name or symbol

### Frontend

The frontend should display the appropriate stablecoin symbol:
- Testnet: "MockUSDT" or "USDT (Test)"
- Mainnet: "USDC"

## User Experience

### For Users

**On Testnet:**
- Use MockUSDT (can mint for free)
- Perfect for testing without real money

**On Mainnet:**
- Use USDC (real money)
- Bridge USDC from other chains or buy on Monad DEXs
- Same deposit/withdrawal flow as USDT

### Deposit Flow

1. User connects wallet to Monad network
2. User approves USDC/MockUSDT to PaymentVault
3. User deposits amount
4. Backend tracks balance off-chain
5. User receives cashback in KAWAI tokens

## Future Plans

### If USDT Becomes Available

When Tether deploys USDT to Monad Mainnet:

**Option 1: Support Both**
- Allow users to deposit either USDC or USDT
- Track balances separately or convert to common unit

**Option 2: Migrate to USDT**
- Update `USDT_TOKEN_ADDRESS` to USDT contract
- Announce migration to users
- Provide conversion mechanism for existing USDC balances

**Option 3: Stay with USDC**
- USDC works perfectly fine
- No need to change if ecosystem is stable

## Developer Notes

### Testing

**Testnet (MockUSDT):**
```bash
# Mint test tokens
cast send 0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc \
  "mint(address,uint256)" \
  YOUR_ADDRESS \
  1000000000 \
  --rpc-url https://testnet-rpc.monad.xyz
```

**Mainnet (USDC):**
- Need real USDC
- Bridge from Ethereum/Base/Polygon
- Or buy on Monad DEXs

### Contract Deployment

When deploying to mainnet:
1. Update `USDT_TOKEN_ADDRESS` in `.env.mainnet`
2. Deploy PaymentVault with USDC address
3. Test deposit/withdrawal flow
4. Verify cashback calculations work correctly

### Documentation Updates

When referring to deposits in docs:
- Use "stablecoin" instead of "USDT" for generic references
- Specify "USDC (mainnet)" or "MockUSDT (testnet)" when being specific
- Update user guides to mention USDC on mainnet

## FAQ

**Q: Why not use USDT on mainnet?**
A: USDT is not yet deployed on Monad Mainnet. USDC is the primary stablecoin available.

**Q: Is USDC as stable as USDT?**
A: Yes, both maintain a $1 peg. USDC is issued by Circle and is widely trusted.

**Q: Can I use other stablecoins?**
A: Currently only USDC (mainnet) and MockUSDT (testnet) are supported. We may add more in the future.

**Q: Will you support USDT when it's available?**
A: We'll evaluate based on user demand and ecosystem adoption.

**Q: Do I need to do anything different with USDC?**
A: No, the deposit flow is identical. Just use USDC instead of USDT.

## References

- [Circle USDC Documentation](https://www.circle.com/en/usdc)
- [Monad Explorer - USDC Token](https://monadscan.com/address/0x754704bc059f8c67012fed69bc8a327a5aafb603)
- [Uniswap on Monad](https://docs.uniswap.org/contracts/v3/reference/deployments/monad-deployments)

---

**Status:** ✅ Production Ready  
**Network:** Monad Testnet (MockUSDT) + Monad Mainnet (USDC)  
**Last Verified:** January 21, 2026
