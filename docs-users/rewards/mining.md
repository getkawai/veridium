# Mining Rewards

Contribute GPU power to the Kawai network and earn the majority of generated tokens!

## 🎯 What is Mining?

Mining on Kawai Network means providing GPU compute power for AI inference. Your computer processes AI requests from users and you earn KAWAI tokens as reward.

!!! info "Current Status"
    Mining is in **bootstrap phase**. The contributor client is ready, but we're currently operating with treasury pool rewards while building the initial user base. GPU mining will open to public contributors when the network is ready for distributed operations.

## ⚡ How It Works

```
User requests AI inference
    ↓
Request routed to available contributor
    ↓
Your GPU processes the request
    ↓
Response sent back to user
    ↓
Reward calculated and split
    ↓
Your share recorded (85-90%)
    ↓
Weekly settlement to blockchain
    ↓
You claim your KAWAI tokens!
```

## 💰 Reward Distribution

### With Referral (85/5/5/5 Split)

If you were invited by a referrer:

| Recipient | Share | Example (100 KAWAI job) |
|-----------|-------|------------------------|
| **You (Contributor)** | 85% | 85 KAWAI |
| Developer (Protocol) | 5% | 5 KAWAI |
| User (Requester) | 5% | 5 KAWAI |
| Affiliator (Referrer) | 5% | 5 KAWAI |

### Without Referral (90/5/5 Split)

If you signed up without a referral code:

| Recipient | Share | Example (100 KAWAI job) |
|-----------|-------|------------------------|
| **You (Contributor)** | 90% | 90 KAWAI |
| Developer (Protocol) | 5% | 5 KAWAI |
| User (Requester) | 5% | 5 KAWAI |

**Tip:** Even if you didn't use a referral code initially, you still earn 90% - very fair!

## 📊 Halving Schedule

Rewards decrease as total supply increases to control inflation:

| Phase | Supply Range | Rate | Per 1M Tokens |
|-------|--------------|------|---------------|
| **1** | 0 - 250M | **100 KAWAI** | 2.5M jobs |
| **2** | 250M - 500M | **50 KAWAI** | 5M jobs |
| **3** | 500M - 750M | **25 KAWAI** | 10M jobs |
| **4** | 750M - 1B | **12 KAWAI** | 20M jobs |

### Example Calculation

**Phase 1 (Current):**
```
AI request generates: 100,000 tokens
Reward: 100 KAWAI per 1M tokens
Total reward pool: 100,000 / 1,000,000 × 100 = 10 KAWAI

Your share (with referral):
10 KAWAI × 85% = 8.5 KAWAI per job
```

**Phase 4 (Late stage):**
```
Same request: 100,000 tokens
Reward: 12 KAWAI per 1M tokens  
Total reward pool: 100,000 / 1,000,000 × 12 = 1.2 KAWAI

Your share (with referral):
1.2 KAWAI × 85% = 1.02 KAWAI per job
```

**Early mining advantage: 8.3x more tokens!**

## 🖥️ System Requirements

### Minimum Requirements

| Component | Requirement |
|-----------|-------------|
| **GPU** | NVIDIA GTX 1060 / AMD RX 580 |
| **VRAM** | 6GB |
| **RAM** | 8GB |
| **Storage** | 20GB free (for models) |
| **Internet** | 10 Mbps upload |
| **OS** | Windows 10+, macOS 11+, Linux |

### Recommended Requirements

| Component | Requirement |
|-----------|-------------|
| **GPU** | NVIDIA RTX 3060+ / AMD RX 6700+ |
| **VRAM** | 12GB+ |
| **RAM** | 16GB+ |
| **Storage** | 50GB+ SSD |
| **Internet** | 50+ Mbps upload |
| **OS** | Latest version |

### Performance Estimates

| GPU Tier | Jobs/Day | Daily Earnings | Monthly |
|----------|----------|----------------|---------|
| **Entry** (GTX 1060) | ~50 | ~25 KAWAI | ~750 KAWAI |
| **Mid** (RTX 3060) | ~150 | ~75 KAWAI | ~2,250 KAWAI |
| **High** (RTX 3080) | ~300 | ~150 KAWAI | ~4,500 KAWAI |
| **Pro** (RTX 4090) | ~500+ | ~250+ KAWAI | ~7,500+ KAWAI |

*Estimates based on Phase 1 rates and typical usage patterns*

## 🚀 Getting Started

### Step 1: Prepare Your System

1. **Check GPU compatibility**
   - NVIDIA: Install latest CUDA drivers
   - AMD: Install latest ROCm drivers

2. **Free up storage**
   - Need ~20GB for AI models
   - SSD recommended for speed

3. **Check internet speed**
   - Run speed test
   - Minimum 10 Mbps upload

### Step 2: Download Contributor Client

!!! warning "Coming Soon"
    Public contributor client will be available soon. Join our [Discord](https://discord.gg/kawai) to be notified!

```bash
# Future download command
curl -o kawai-miner https://kawai.network/download/contributor
chmod +x kawai-miner
```

### Step 3: Setup & Configuration

```bash
# Initialize contributor node
./kawai-miner init

# Enter your wallet address
# Choose models to serve
# Set resource limits
```

### Step 4: Start Mining

```bash
# Start the contributor node
./kawai-miner start

# Check status
./kawai-miner status

# View logs
./kawai-miner logs
```

## 📈 Maximizing Earnings

### Tip 1: Maximize Uptime

**More uptime = more jobs = more earnings**

```
90% uptime: ~90 jobs/day
95% uptime: ~95 jobs/day
99% uptime: ~99 jobs/day

5% difference = ~150 KAWAI/month extra!
```

**Strategies:**
- Run 24/7 if possible
- Use dedicated machine
- Set up automatic restart on failure
- Monitor remotely

### Tip 2: Optimize Hardware

**GPU optimization:**
- Keep drivers updated
- Maintain good cooling (prevents throttling)
- Clean dust regularly
- Consider undervolting (efficiency)

**Network optimization:**
- Use wired connection (not WiFi)
- QoS settings for mining traffic
- Monitor latency to network

### Tip 3: Serve Multiple Models

**More models = more job opportunities**

```
1 model: Limited job pool
3 models: Medium job pool
5+ models: Maximum jobs

But balance with your VRAM!
```

### Tip 4: Peak Hours Strategy

**Mining during peak usage:**
- More jobs available during peak hours
- Consider timezone differences
- Europe peak: 8AM-11PM CET
- US peak: 9AM-12AM EST
- Asia peak: 9AM-12AM JST

### Tip 5: Join Early

**Halving advantage:**

```
Mine in Phase 1:
  1,000 jobs × 100 KAWAI/1M × 0.1 = 10,000 KAWAI

Mine in Phase 4:
  1,000 jobs × 12 KAWAI/1M × 0.1 = 1,200 KAWAI

Same work, 8.3x difference!
```

## 🔐 Security & Safety

### Contributor Security

**What the node does:**
- ✅ Processes AI inference requests
- ✅ Reports completed jobs
- ✅ Maintains heartbeat connection

**What it does NOT do:**
- ❌ Access your files
- ❌ Control your wallet
- ❌ Steal private keys
- ❌ Mine cryptocurrency (no PoW)

### Wallet Security

**Best practices:**
- Use a separate wallet for mining
- Don't store large amounts in mining wallet
- Keep private keys offline
- Use hardware wallet for large holdings

### System Security

**Recommendations:**
- Keep OS updated
- Use firewall
- Monitor resource usage
- Check logs regularly for anomalies

## 📊 Monitoring & Analytics

### Dashboard Metrics

**Real-time stats:**
- Jobs completed (last hour/day/week)
- Earnings (pending/claimed)
- Success rate
- Average response time
- Uptime percentage

**Performance graphs:**
- Earnings over time
- Jobs per day
- GPU utilization
- Network latency

### CLI Commands

```bash
# Check earnings
./kawai-miner earnings

# View job history
./kawai-miner history

# Performance stats
./kawai-miner stats

# Network status
./kawai-miner network
```

## 💸 Claiming Your Rewards

### Settlement Schedule

**Weekly cycle:**
- **Sunday 11:59 PM**: Period closes
- **Monday 12:00 AM**: Settlement begins
- **Monday 12:30 AM**: Merkle tree generated
- **Monday 1:00 AM**: Proofs available
- **Monday onwards**: Claim anytime

### Claiming Process

1. **Go to Rewards Dashboard**
   - Wallet → Rewards → Mining tab

2. **View claimable periods**
   - See weeks with pending rewards
   - Check amounts and proofs

3. **Click "Claim"**
   - Select period(s) to claim
   - Review transaction
   - Confirm in wallet

4. **Receive tokens**
   - KAWAI minted to your address
   - Usually ~10-30 seconds

[Detailed claiming guide →](claiming.md)

### Batch Claiming

Save on gas fees by claiming multiple periods at once:

```
Single claim: ~0.001 MON per period
Batch claim 4 periods: ~0.002 MON total

Save ~50% on gas fees!
```

## 🆘 Troubleshooting

### "No Jobs Received"

**Check:**
- [ ] Is node running? (`./kawai-miner status`)
- [ ] Internet connection stable?
- [ ] Firewall blocking ports?
- [ ] GPU available (not used by other apps)?
- [ ] Models downloaded correctly?

**Solutions:**
- Restart node
- Check logs (`./kawai-miner logs`)
- Verify port forwarding
- Update to latest version

### "Low Success Rate"

**Possible causes:**
- GPU overheating (thermal throttling)
- Slow internet connection
- Insufficient VRAM
- Corrupted models

**Solutions:**
- Improve cooling
- Test internet speed
- Reduce served models
- Re-download models

### "Rewards Not Showing"

**Check:**
- Is it after Monday settlement?
- Did jobs complete successfully?
- Check dashboard: Wallet → Rewards → Mining
- Verify wallet address is correct

**If still not showing:**
- Wait 1 hour (settlement delay)
- Check settlement status on Discord
- Contact support with wallet address

## 💡 Advanced Topics

### Multi-GPU Setup

Run multiple GPUs for higher earnings:

```bash
# Configure multiple GPUs
./kawai-miner config --gpus 0,1,2

# Each GPU can serve different models
# 3x GPUs ≈ 3x earnings potential
```

### Remote Mining

Monitor and control mining remotely:

```bash
# Enable remote API
./kawai-miner config --api-enable

# Access from anywhere
curl http://your-ip:8080/status
```

### Automatic Claiming

Set up automatic reward claiming (coming soon):

```bash
# Configure auto-claim
./kawai-miner config --auto-claim-enabled
./kawai-miner config --auto-claim-threshold 100

# Claims automatically when > 100 KAWAI pending
```

## 📚 FAQs

**Q: Do I need to keep the app open?**  
A: The contributor node runs as a background service. You can close terminal but keep the service running.

**Q: Can I mine while gaming?**  
A: Not recommended. Mining uses GPU at ~80-100%. Will affect gaming performance.

**Q: What if my internet goes down?**  
A: Node will automatically reconnect. Missed jobs don't affect your reputation (yet).

**Q: Can I mine on laptop?**  
A: Possible but not recommended. Laptops have limited cooling and battery drain. Use desktop.

**Q: How much electricity does it use?**  
A: Depends on GPU. RTX 3060 uses ~170W. Calculate: 170W × 24h × $0.12/kWh = ~$0.49/day.

**Q: Is it profitable?**  
A: Depends on electricity costs and KAWAI price. Monitor earnings vs. costs.

## ✅ Checklist

Before starting to mine:

- [ ] GPU meets minimum requirements
- [ ] Latest drivers installed
- [ ] 20GB+ free storage
- [ ] Stable internet (10+ Mbps)
- [ ] Wallet address ready
- [ ] Understand reward splits
- [ ] Joined Discord for updates

## 🚀 Next Steps

1. **[Join Discord](https://discord.gg/kawai)** - Get notified when public mining launches
2. **[Prepare Hardware](#system-requirements)** - Check your GPU compatibility
3. **[Understand Rewards](#reward-distribution)** - Know what you'll earn
4. **[Learn Claiming](claiming.md)** - How to get your rewards

---

**Ready to mine?** Join our [Discord community](https://discord.gg/kawai) for latest updates!

