# USDT to Stablecoin Terminology Migration

## Summary

Successfully migrated the project from USDT-specific terminology to generic "stablecoin" terminology to support both MockUSDT (testnet) and USDC (mainnet).

## Why This Change?

- **USDT not available on Monad Mainnet**: As of January 2026, USDT has not been deployed on Monad Mainnet
- **USDC is available**: Circle's USDC is deployed at `0x754704bc059f8c67012fed69bc8a327a5aafb603`
- **Avoid confusion**: Using generic "stablecoin" terminology prevents confusion about which token is used
- **Backward compatibility**: Variable names like `USDT_TOKEN_ADDRESS` and `UsdtTokenAddress` kept for compatibility

## Changes Made

### 1. Smart Contracts ✅

**File**: `contracts/contracts/PaymentVault.sol`
- Changed `IERC20 public usdt` → `IERC20 public stablecoin`
- Updated constructor parameter name
- Updated all internal references
- Compiled successfully with `forge build`
- Generated new ABI and Go bindings with `make contracts-bindings`

### 2. Environment Configuration ✅

**Files**: `.env.mainnet`, `.env.README`, `DEPLOYMENT.md`
- Updated `USDT_TOKEN_ADDRESS` to USDC address for mainnet: `0x754704bc059f8c67012fed69bc8a327a5aafb603`
- Added clear documentation distinguishing testnet (MockUSDT) vs mainnet (USDC)
- Updated deployment guide with separate addresses for each network

### 3. Documentation ✅

**Files**: `README.md`, `STABLECOIN_SUPPORT.md`
- Changed user-facing terminology from "USDT" to "stablecoin"
- Created comprehensive stablecoin support documentation
- Explained testnet vs mainnet token differences

### 4. Backend Services ✅

**File**: `pkg/blockchain/revenue_settlement.go`
- Changed struct field `usdtToken` → `stablecoinToken`
- Updated all comments and log messages
- Changed proof prefix from "usdt:" to "stablecoin:"
- Updated RewardType from "usdt" to "stablecoin"
- Updated function comments throughout

**File**: `internal/services/deai_service.go`
- Updated `GetVaultBalance()` comments to mention stablecoin
- Updated `DepositToVault()` comments and variable names
- Updated `GetUSDTAllowance()` with backward compatibility note
- Updated `ApproveUSDT()` with backward compatibility note
- Updated `TransferUSDT()` with backward compatibility note
- Updated `MintTestTokens()` comments
- Function names kept for backward compatibility

**File**: `internal/services/config_service.go`
- Added comment to `Usdt` field explaining it represents stablecoin address

**File**: `pkg/blockchain/client.go`
- Updated `MarketplaceCreateOrder()` parameter name and comments
- Updated `ValidateTradeBalance()` comments
- Updated `GetUSDTBalance()` with backward compatibility note

**File**: `cmd/dev/inject-test-usdt/main.go`
- Updated all comments and log messages
- Changed display text from "USDT" to "stablecoin"

**File**: `internal/constant/blockchain.go`
- Added comprehensive comment explaining `UsdtTokenAddress` usage
- Documented testnet vs mainnet addresses

### 5. Build Verification ✅

- Ran `go build -o /dev/null .` - **SUCCESS**
- No compilation errors
- All type changes compatible

## Backward Compatibility

The following names were **intentionally kept** for backward compatibility:

### Variable Names
- `USDT_TOKEN_ADDRESS` (env var)
- `UsdtTokenAddress` (Go constant)
- `Usdt` (struct field in ContractAddresses)

### Function Names
- `GetUSDTBalance()`
- `GetUSDTAllowance()`
- `ApproveUSDT()`
- `TransferUSDT()`

### Contract Resolution
- `contracts.ResolveAddress("MockUSDT")` - still works, points to stablecoin

All these names now have comments explaining they work with any stablecoin (MockUSDT on testnet, USDC on mainnet).

## Network-Specific Addresses

### Testnet (Monad Testnet)
- **Token**: MockUSDT
- **Address**: `0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc`
- **Decimals**: 6
- **Purpose**: Testing only
- **Special Features**: Has public `mint()` function for easy testing

### Mainnet (Monad Mainnet)
- **Token**: USDC (Circle)
- **Address**: `0x754704bc059f8c67012fed69bc8a327a5aafb603`
- **Decimals**: 6
- **Purpose**: Production deposits and payments
- **Special Features**: NO public `mint()` function - must acquire through exchanges/bridges

## ⚠️ Important Compatibility Notes

### Functions That Work on Both Networks
These standard ERC-20 functions work identically on both testnet (MockUSDT) and mainnet (USDC):
- ✅ `balanceOf(address)` - Check balance
- ✅ `transfer(address, uint256)` - Transfer tokens
- ✅ `approve(address, uint256)` - Approve spending
- ✅ `transferFrom(address, address, uint256)` - Transfer from approved address
- ✅ `allowance(address, address)` - Check allowance
- ✅ `decimals()` - Get decimals (both return 6)
- ✅ `totalSupply()` - Get total supply
- ✅ `name()` - Get token name
- ✅ `symbol()` - Get token symbol

### Functions That ONLY Work on Testnet
- ❌ `mint(address, uint256)` - **ONLY available on MockUSDT (testnet)**
  - This function does NOT exist on USDC mainnet
  - Code calling `mint()` will FAIL on mainnet
  - Protected by environment check in `MintTestTokens()`

### How to Get Stablecoin on Mainnet
Since USDC on mainnet doesn't have a public `mint()` function, users must:
1. **Bridge** USDC from Ethereum/other chains to Monad
2. **Buy** USDC on Monad DEXes
3. **Receive** USDC from other users via transfer

## Testing Checklist

- [x] Smart contracts compile successfully
- [x] Go bindings generated successfully
- [x] Backend builds without errors
- [ ] Frontend TypeScript bindings updated (if needed)
- [ ] Test deposit flow with USDC on mainnet
- [ ] Test cashback claims with stablecoin
- [ ] Verify all user-facing text shows "stablecoin" not "USDT"

## Next Steps

1. **Deploy Updated Contracts**: Deploy PaymentVault with new variable names to mainnet
2. **Update Frontend**: Check if frontend TypeScript bindings need regeneration
3. **Test on Mainnet**: Verify USDC deposits work correctly
4. **Update User Documentation**: Ensure all user-facing docs mention "stablecoin"

## References

- [STABLECOIN_SUPPORT.md](./STABLECOIN_SUPPORT.md) - Comprehensive stablecoin documentation
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment guide with network-specific addresses
- [MonadScan USDC](https://explorer.monad.xyz/address/0x754704bc059f8c67012fed69bc8a327a5aafb603) - USDC contract on Monad Mainnet

---

**Migration Date**: January 21, 2026  
**Status**: ✅ Complete  
**Build Status**: ✅ Passing
