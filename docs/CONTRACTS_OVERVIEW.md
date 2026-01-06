# Veridium Smart Contracts - Complete Overview

**Status:** Current Implementation (Jan 2026)  
**Network:** Monad Blockchain (Testnet)

## Overview

All smart contracts for the Veridium DePIN network, deployed on **Monad Blockchain** (high-throughput, low-cost, EVM-compatible).

## Current Contracts (8 Total)

### 1. **KawaiToken.sol** (ERC20 + AccessControl + Burnable)
- **Purpose**: Native utility token for the Veridium ecosystem
- **Features**:
  - Fair Launch: No initial mint, supply starts at 0
  - Max Supply: 1 Billion tokens (18 decimals)
  - Mining Rewards: Minted on-demand by backend (MINTER_ROLE)
  - Dynamic difficulty with halving mechanism
- **Tokenomics**:
  - 70% to Contributors (miners)
  - 30% to Admin (platform)
  - Post-mining: USDT profit sharing for KAWAI holders

### 2. **Escrow.sol** (P2P OTC Market v2)
- **Purpose**: Decentralized marketplace for KAWAI ↔ USDT trading
- **Features**:
  - Partial Fill: Buyers can fill any amount ≤ listing amount
  - Atomic Settlement: USDT ↔ KAWAI swaps in one transaction
  - Fee: 0.5% to platform
  - No order book, fully on-chain
- **Use Cases**:
  - Contributors sell mined KAWAI for USDT
  - Users buy KAWAI for staking/governance

### 3. **PaymentVault.sol**
- **Purpose**: Handle user deposits for AI service credits
- **Features**:
  - Users deposit USDT to get off-chain credits
  - Owner can withdraw for revenue distribution
  - ReentrancyGuard protection
  - Event emission for backend tracking
- **Flow**:
  1. User deposits USDT → emits `Deposited` event
  2. Backend listens to event → credits user balance in KV
  3. User consumes credits → backend deducts from KV
  4. Admin withdraws USDT → distributes to contributors + treasury

### 4. **MerkleDistributor.sol** (Generic Base)
- **Purpose**: Generic Merkle-based reward distribution
- **Modes**:
  - **Mint Mode**: For KAWAI rewards (mints new tokens on claim)
  - **Transfer Mode**: For USDT dividends (transfers from pre-funded balance)
- **Features**:
  - Bitmap-based claim tracking (gas-optimized)
  - Trustless verification with Merkle proofs
  - Scalable to millions of users
- **Use Cases**:
  - Base contract for specialized distributors
  - Monthly USDT profit sharing for KAWAI holders

### 5. **MiningRewardDistributor.sol** ⛏️
- **Purpose**: Specialized contract for mining rewards (85/5/5/5 split)
- **Features**:
  - **Multi-recipient Distribution**: Contributor, Developer, User (cashback), Affiliator
  - **Period-based Settlement**: Weekly batches
  - **Merkle Proof Verification**: Gas-efficient (~150k gas per claim)
  - **Mint-on-demand**: Requires MINTER_ROLE on KawaiToken
- **Rewards Split**:
  - Contributor: 85% (or 90% non-referral)
  - Developer: 5%
  - User: 5% (use-to-earn cashback)
  - Affiliator: 5% (if referral exists)
- **See**: Root `README.md` for mining tokenomics

### 6. **DepositCashbackDistributor.sol** 💰
- **Purpose**: Tiered KAWAI cashback on USDT deposits
- **Features**:
  - **5 Tiers**: 1% to 5% cashback based on total deposits
  - **Period-based Settlement**: Weekly batches
  - **Merkle Proof Verification**: Gas-efficient claiming
  - **Mint-on-demand**: Requires MINTER_ROLE on KawaiToken
- **Tier Structure**:
  - Tier 1: $100+ → 1% cashback
  - Tier 2: $500+ → 2% cashback
  - Tier 3: $2,000+ → 3% cashback
  - Tier 4: $10,000+ → 4% cashback
  - Tier 5: $50,000+ → 5% cashback
- **See**: `CASHBACK_SYSTEM.md` for complete implementation

### 7. **ReferralRewardDistributor.sol** 🎁
- **Purpose**: Dual rewards for referral program
- **Features**:
  - **Dual Rewards**: KAWAI tokens (minted) + USDT (transferred)
  - **Period-based Settlement**: Weekly/monthly batches
  - **Gas-efficient**: Merkle proof verification (~150k gas per claim)
  - **Batch Claiming**: Users can claim multiple periods at once
  - **Security**: ReentrancyGuard, double-claim prevention
- **Rewards Structure**:
  - New User with Referral: 10 USDT + 100 KAWAI
  - Referrer: 5 USDT + 50 KAWAI per successful referral
- **See**: `REFERRAL_SYSTEM.md` and `docs/REFERRAL_CONTRACT_GUIDE.md`

### 8. **MockUSDT.sol** (Testnet Only)
- **Purpose**: ERC20 mock for testing on Monad testnet
- **Features**: Standard ERC20 with public mint function

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Veridium Ecosystem                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐      ┌──────────────┐                    │
│  │ KawaiToken   │◄─────┤ Escrow       │                    │
│  │ (ERC20)      │      │ (P2P Market) │                    │
│  └──────┬───────┘      └──────────────┘                    │
│         │                                                    │
│         │ MINTER_ROLE                                       │
│         │                                                    │
│  ┌──────▼──────────────────────────────────────────┐       │
│  │      Reward Distributors (Merkle-based)          │       │
│  │                                                   │       │
│  │  • MiningRewardDistributor (85/5/5/5 split)     │       │
│  │  • DepositCashbackDistributor (1-5% tiers)      │       │
│  │  • ReferralRewardDistributor (KAWAI+USDT)       │       │
│  │  • MerkleDistributor (USDT profit sharing)      │       │
│  └──────────────────────────────────────────────────┘       │
│                                                              │
│  ┌──────────────┐                                           │
│  │ PaymentVault │                                           │
│  │ (USDT)       │                                           │
│  └──────────────┘                                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Deployment

### Prerequisites

1. Install Foundry:
```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

2. Set environment variables:
```bash
cp env.example .env
# Edit .env with your values:
# - PRIVATE_KEY
# - RPC_URL (Monad testnet/mainnet)
# - USDT_ADDRESS
```

### Deploy All Contracts

```bash
# 1. Deploy KawaiToken
forge script script/DeployKawai.s.sol:DeployKawai \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify

# 2. Deploy Escrow
forge script script/DeployEscrow.s.sol:DeployEscrow \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify

# 3. Deploy PaymentVault
forge script script/DeployPaymentVault.s.sol:DeployPaymentVault \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify

# 4. Deploy MerkleDistributor (for mining rewards)
forge script script/DeployMerkleDistributor.s.sol:DeployMerkleDistributor \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify

# 5. Deploy ReferralRewardDistributor
forge script script/DeployReferralDistributor.s.sol:DeployReferralDistributor \
  --rpc-url $RPC_URL \
  --broadcast \
  --verify
```

### Post-Deployment Setup

1. **Grant MINTER_ROLE**:
```bash
cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $MERKLE_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY

cast send $KAWAI_TOKEN_ADDRESS \
  "grantRole(bytes32,address)" \
  $(cast keccak "MINTER_ROLE") \
  $REFERRAL_DISTRIBUTOR_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY
```

2. **Fund ReferralRewardDistributor with USDT**:
```bash
# Approve USDT
cast send $USDT_ADDRESS \
  "approve(address,uint256)" \
  $REFERRAL_DISTRIBUTOR_ADDRESS \
  1000000000000 \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY

# Fund distributor
cast send $REFERRAL_DISTRIBUTOR_ADDRESS \
  "fundUSDT(uint256)" \
  1000000000000 \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY
```

## Testing

### Run All Tests

```bash
forge test -vvv
```

### Run Specific Test

```bash
forge test --match-contract ReferralRewardDistributorTest -vvv
```

### Gas Report

```bash
forge test --gas-report
```

### Coverage

```bash
forge coverage
```

## Integration with Backend

### 1. Listen to Events

```go
// Listen to PaymentVault deposits
logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
    Addresses: []common.Address{paymentVaultAddress},
    Topics:    [][]common.Hash{{depositEventTopic}},
})

for _, log := range logs {
    // Parse event
    event, err := vault.ParseDeposited(log)
    
    // Credit user balance in KV
    store.AddBalance(ctx, event.User.Hex(), event.Amount.String())
}
```

### 2. Generate Merkle Trees

```go
// Weekly settlement for mining rewards
rewards := collectMiningRewards()
tree := generateMerkleTree(rewards)
distributor.SetMerkleRoot(tree.Root())

// Monthly settlement for referral rewards
referralRewards := collectReferralRewards()
referralTree := generateMerkleTree(referralRewards)
referralDistributor.SetMerkleRoot(referralTree.Root())
```

### 3. Provide Proofs to Users

```go
// API endpoint: GET /v1/rewards/proof/:address/:period
func GetProof(c *gin.Context) {
    address := c.Param("address")
    period := c.Param("period")
    
    proof, amount := getProofFromKV(address, period)
    
    c.JSON(200, gin.H{
        "proof":  proof,
        "amount": amount,
    })
}
```

## Security Considerations

### Audits
- [ ] Internal audit completed
- [ ] External audit by reputable firm
- [ ] Bug bounty program launched

### Best Practices
- ✅ ReentrancyGuard on all state-changing functions
- ✅ AccessControl for privileged operations
- ✅ SafeERC20 for token transfers
- ✅ Comprehensive test coverage (>90%)
- ✅ Gas optimization (bitmap for claims)
- ✅ Event emission for off-chain tracking

### Upgrade Strategy
- Contracts are **immutable** (no proxy pattern)
- For upgrades, deploy new contracts and migrate users
- Use timelock for critical admin operations (future)

## Gas Costs (Monad Network)

| Operation | Gas | Cost @ 1 gwei | Cost @ 10 gwei |
|-----------|-----|---------------|----------------|
| Deposit USDT | 50k | $0.05 | $0.50 |
| Claim Mining Rewards | 150k | $0.15 | $1.50 |
| Claim Referral Rewards | 150k | $0.15 | $1.50 |
| Trade on Escrow | 100k | $0.10 | $1.00 |
| Set Merkle Root | 50k | $0.05 | $0.50 |

**Note**: Monad's high throughput (10k TPS) keeps gas prices consistently low.

## Monitoring

### On-chain Metrics
- Total KAWAI minted
- Total USDT deposited
- Total referral rewards claimed
- Active referrers count
- Trading volume on Escrow

### Alerts
- Low USDT balance in ReferralRewardDistributor
- Failed Merkle root updates
- Unusual claim patterns (potential exploit)

## Resources

- [Foundry Book](https://book.getfoundry.sh/)
- [Monad Documentation](https://docs.monad.xyz/)
- [OpenZeppelin Contracts](https://docs.openzeppelin.com/contracts/)
- [Referral Contract Guide](./REFERRAL_CONTRACT_GUIDE.md)

## Support

For questions or issues:
- GitHub Issues: [kawai-network/veridium](https://github.com/kawai-network/veridium)
- Discord: [Join our community](#)
- Email: dev@kawai.network

---

**Built with ❤️ by the Veridium Team**

