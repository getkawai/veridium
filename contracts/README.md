# Kawai AI Smart Contracts

Official smart contracts for the Kawai AI platform on Monad blockchain.

## 🌐 Links

- **Website**: <https://getkawai.com>
- **Documentation**: <https://getkawai.com/docs>
- **Telegram**: <https://t.me/getkawai>
- **Discord**: <https://discord.gg/SNf3ZEa8Eq>
- **Twitter**: <https://twitter.com/getkawai>

## 📋 Deployed Contracts

### Monad Mainnet (Chain ID: 143)

| Contract | Address | Description |
|----------|---------|-------------|
| KAWAI Token | TBD | ERC20 token with 1B max supply |
| Payment Vault | TBD | USDC deposit vault for AI credits |
| Mining Distributor | TBD | Mining rewards with referral splits |
| Cashback Distributor | TBD | Deposit cashback rewards |
| Referral Distributor | TBD | Referral commission rewards |
| OTCMarket | TBD | P2P OTC market for KAWAI-USDC trading |

### Monad Testnet (Chain ID: 10143)

| Contract | Address | Description |
|----------|---------|-------------|
| KAWAI Token | TBD | ERC20 token (testnet) |
| MockUSDT | `0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc` | Test stablecoin |
| Payment Vault | TBD | MockUSDT deposit vault |
| Mining Distributor | TBD | Mining rewards (testnet) |
| Cashback Distributor | TBD | Cashback rewards (testnet) |
| Referral Distributor | TBD | Referral rewards (testnet) |
| OTCMarket | TBD | P2P OTC market (testnet) |

## 🏗️ Architecture

### KawaiToken.sol
- ERC20 token with minting capability
- Max supply: 1,000,000,000 KAWAI (1 billion)
- Fair launch: No initial mint
- Access control: MINTER_ROLE for distributors

### PaymentVault.sol
- Accepts USDC deposits (mainnet) or MockUSDT (testnet)
- Converts deposits to AI service credits
- Owner can withdraw for revenue distribution

### MiningRewardDistributor.sol
- Distributes mining rewards via Merkle proofs
- Referral splits: 85% contributor, 5% developer, 5% user, 5% affiliator
- Non-referral: 90% contributor, 5% developer, 5% user
- Weekly settlement periods

### DepositCashbackDistributor.sol
- Distributes deposit cashback in KAWAI tokens
- Batch claiming support

### ReferralRewardDistributor.sol
- Distributes referral commissions in KAWAI
- Period-based Merkle roots
- Tracks unique referrers

## 🚀 Deployment

### Prerequisites
```bash
# Install Foundry
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Install dependencies
forge install
```

### Compile Contracts
```bash
forge build
```

### Run Tests
```bash
forge test -vvv
```

### Deploy to Mainnet
```bash
# Option A: Full suite deployment
cd contracts
forge script script/DeployKawai.s.sol:DeployKawai \
  --rpc-url $MONAD_MAINNET_RPC \
  --private-key $DEPLOYER_PRIVATE_KEY \
  --broadcast \
  --verify

# Option B: Modular deployment
make contracts-deploy-vault
make contracts-deploy-mining-mainnet
make contracts-deploy-cashback-mainnet
make contracts-deploy-referral-mainnet
```

## 🔒 Security

All contracts use OpenZeppelin libraries for battle-tested implementations.

**Security Features by Contract:**

| Contract | ReentrancyGuard | Access Control | Pausable | Merkle Proofs |
|----------|----------------|----------------|----------|---------------|
| KawaiToken | ❌ | ✅ AccessControl | ❌ | ❌ |
| PaymentVault | ✅ | ✅ Ownable | ❌ | ❌ |
| MiningRewardDistributor | ✅ | ✅ Ownable | ✅ | ✅ |
| DepositCashbackDistributor | ✅ | ✅ Ownable | ✅ | ✅ |
| ReferralRewardDistributor | ✅ | ✅ Ownable | ✅ | ✅ |
| OTCMarket (Escrow) | ✅ | ❌ | ❌ | ❌ |

**Key Security Measures:**
- ReentrancyGuard on all state-changing functions (where applicable)
- Role-based access control (Ownable/AccessControl)
- Emergency pause mechanisms (distributor contracts)
- Merkle proof verification for gas-efficient distributions
- SafeERC20 for all token transfers
- Input validation and zero-address checks

## 📄 License

MIT License - see [LICENSE](../LICENSE) for details

## 🤝 Contributing

See main repository [README](../README.md) for contribution guidelines.

## 📞 Support

- Discord: <https://discord.gg/SNf3ZEa8Eq>
- Email: <mailto:support@getkawai.com>
- Documentation: <https://getkawai.com/docs>

---

## Built with ❤️ by the Kawai AI team
