# Technical FAQ

## What is a Merkle tree?

A Merkle tree is a cryptographic data structure that allows us to:
- Store thousands of rewards off-chain (gas-free)
- Prove individual rewards on-chain (secure)
- Verify claims without storing all data on blockchain

Think of it as a "receipt" that proves your rewards are legitimate.

## Why use off-chain + on-chain hybrid?

**Benefits:**
- ⚡ **Fast**: Instant reward tracking
- 💰 **Cheap**: No gas fees for accumulation
- 🔒 **Secure**: Blockchain verification when claiming
- 📊 **Transparent**: All data verifiable

## What is a smart contract?

A smart contract is code that runs on the blockchain. It:
- Executes automatically (no middleman)
- Cannot be changed (immutable)
- Is publicly verifiable (transparent)
- Handles token minting and claiming

## How do I get MON tokens?

MON is Monad's native token for gas fees. Get it from:
1. **Monad Faucet**: Free testnet MON
2. **Exchanges**: Buy on supported exchanges (mainnet)
3. **Bridge**: Bridge from other chains

The app will guide you when needed.

## What wallet does Kawai use?

Kawai has a **built-in self-custodial wallet**:
- Private keys stored on your device
- Compatible with standard Ethereum wallets
- Can export/import seed phrase
- No external wallet app needed

## Can I use MetaMask instead?

Currently, Kawai uses its built-in wallet for the best user experience. MetaMask integration may be added in the future.

## What is gas fee?

Gas fee is a small payment (in MON) to process blockchain transactions. You pay gas when:
- Claiming rewards (~0.001 MON)
- Trading on P2P marketplace
- Transferring tokens

## Why are there settlement delays?

Weekly settlement allows us to:
- Batch thousands of rewards efficiently
- Reduce blockchain congestion
- Lower overall gas costs
- Maintain system security

You still see rewards instantly - claiming just happens weekly.

## Is my data private?

Yes! We only store:
- **On-chain**: Wallet addresses, token balances, claims
- **Off-chain**: Reward calculations, usage stats
- **Never stored**: Private keys, personal info, AI conversations

## How can I verify my rewards?

1. Check your Rewards dashboard for pending amounts
2. After settlement, get your Merkle proof
3. Verify proof against on-chain Merkle root
4. All smart contracts are open source and auditable

## What happens if the server goes down?

Your tokens are safe! They're on the blockchain, not our servers. You can:
- Access tokens via any Ethereum wallet
- Claim rewards using Merkle proofs (stored redundantly)
- Trade on-chain without our interface

## Can I run my own node?

For contributors: Yes! The contributor client will be open source.

For users: Not required. The desktop app handles everything.

## What is the difference between testnet and mainnet?

- **Testnet** (current): Testing environment with free tokens
- **Mainnet** (future): Production environment with real value

Your testnet rewards won't transfer to mainnet, but you'll get early adopter benefits!

---

**Need more help?** Check [General FAQ](general.md) or [contact support](../support/contact.md).
