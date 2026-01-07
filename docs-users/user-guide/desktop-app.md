# Desktop App Overview

Welcome to the Kawai DeAI Network desktop application! This guide will help you understand all the features and how to use them effectively.

## 🎯 What is Kawai Desktop App?

The Kawai desktop app is your gateway to decentralized AI services. It's built with modern technology (Wails v3 + React) to provide a premium, native desktop experience for:

- **AI Chat**: Talk with powerful AI models locally or through the network
- **Image Generation**: Create stunning images using Google Gemini AI
- **Built-in Wallet**: Manage your USDT deposits and KAWAI tokens
- **Rewards Dashboard**: Track and claim your earnings
- **P2P Trading**: Buy and sell KAWAI tokens directly with other users

## 🏗️ App Structure

### Main Sections

The app is divided into several key sections:

#### 1. **Chat Interface**

Your main workspace for AI conversations.

- **New Chat**: Start fresh conversations anytime
- **Chat History**: All your conversations are saved locally
- **Model Selection**: Choose from different AI models
- **Knowledge Base**: Add documents to give AI context about your projects

**Features:**
- ✅ Fast local inference (no internet delays)
- ✅ OpenAI-compatible format
- ✅ Conversation history
- ✅ Multi-turn conversations
- ✅ Code highlighting
- ✅ Markdown support

#### 2. **Image Generation**

Create AI-generated images with advanced controls.

- **Two Models Available**:
  - `gemini-2.5-flash-image` (Nano Banana) - Fast, 1024px
  - `gemini-3-pro-image-preview` (Nano Banana Pro) - HD/4K quality

- **Aspect Ratios**: 10 options (1:1, 16:9, 9:16, 4:3, 3:4, etc.)
- **Quality Tiers**: Standard (1K), HD (2K), Ultra (4K)

**Usage:**
1. Enter your prompt (describe what you want)
2. Select aspect ratio and quality
3. Click "Generate"
4. Download your image

#### 3. **Wallet Dashboard**

Manage your funds and rewards.

**Tabs:**
- **Overview**: Balance summary, recent transactions
- **Deposit**: Add USDT to your account
- **Rewards**: View and claim mining, cashback, and referral rewards
- **Trading**: Access P2P marketplace

**Features:**
- ✅ Real-time balance updates
- ✅ Transaction history
- ✅ Multi-currency support (USDT, KAWAI)
- ✅ Secure Web3 integration

#### 4. **Marketplace**

Trade KAWAI ↔ USDT with other users.

- **Order Book**: See all buy/sell orders
- **Create Orders**: List your KAWAI for sale or buy from others
- **Market Stats**: Track price trends and volume
- **Order History**: Review your past trades

#### 5. **Settings**

Customize your experience.

- **Appearance**: Light/dark theme
- **Language**: Multiple language support
- **Network**: Switch between testnet and mainnet
- **API Keys**: Manage authentication

## 🔐 Security Features

### Local-First Design

- **Private Keys**: Never leave your device
- **Local Storage**: Chat history stored on your computer
- **Encrypted**: Sensitive data is encrypted at rest

### Web3 Integration

- **Secure Encryption**: Your wallet is password-protected and encrypted
- **Transaction Signing**: You approve every transaction
- **No Custody**: We never hold your funds

### API Key Security

- **Generated Locally**: API keys created on your device
- **Encrypted Storage**: Keys stored securely
- **Revocable**: Delete keys anytime

## 💡 Key Features

### 1. **AI Chat with Local Inference**

**What it means:**
- AI runs on your computer (no cloud needed)
- Faster responses (no network delay)
- Complete privacy (data stays local)

**How it works:**
1. App downloads AI model (one-time)
2. Model runs using `llama.cpp` (efficient C++ engine)
3. You get instant responses without internet

**Benefits:**
- 🚀 Faster than cloud AI
- 🔒 100% private
- 💰 Earn rewards while using AI

### 2. **Knowledge Base Integration**

**What it is:**
Upload documents so AI understands your context.

**Supported Formats:**
- PDF documents
- Word documents (DOCX)
- Markdown files
- Text files
- Code files

**How to use:**
1. Click "Knowledge Base" in sidebar
2. Upload your documents
3. AI automatically indexes them
4. Ask questions about your documents

**Example:**
```
You: "What's the revenue for Q3 in the financial report?"
AI: "Based on the uploaded Q3_Report.pdf, revenue was $2.5M..."
```

### 3. **Hybrid Reward System**

Earn KAWAI tokens in multiple ways:

**Use-to-Earn (5% cashback)**
- Get 5% back on every AI request
- Applies to both chat and image generation
- Instant accumulation

**Deposit Cashback (1-5% tiered)**
- Earn when depositing USDT
- Higher deposits = higher rate
- First-time bonus: 5%

**Referral Rewards (5 USDT + lifetime commission)**
- Invite friends, earn 5 USDT
- Get 5% of their mining rewards forever
- No limit on referrals

**Mining Rewards (85-90% of generated tokens)**
- Contribute GPU power (future feature)
- Earn majority of tokens
- Fair halving schedule

### 4. **P2P Marketplace**

Trade without intermediaries.

**How it works:**
1. Users create buy/sell orders
2. Orders stored in smart contracts
3. Others can fill orders (partial fills supported)
4. Funds held in escrow until complete

**Features:**
- ✅ No middleman fees
- ✅ Smart contract security
- ✅ Partial order filling
- ✅ Real-time order book

## 🎨 User Interface

### Modern Design

- **Material Design**: Clean, intuitive interface
- **Dark Mode**: Easy on the eyes
- **Responsive**: Adapts to your window size
- **Animations**: Smooth transitions

### Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| New Chat | `Cmd/Ctrl + N` |
| Search | `Cmd/Ctrl + K` |
| Toggle Sidebar | `Cmd/Ctrl + B` |
| Settings | `Cmd/Ctrl + ,` |

### Navigation

- **Sidebar**: Main navigation (Chat, Wallet, Marketplace)
- **Top Bar**: Current section and actions
- **Bottom Bar**: Status and notifications

## 📊 Performance

### Fast & Efficient

- **Local Inference**: 20-50 tokens/second (depending on hardware)
- **Image Generation**: 10-30 seconds (depending on quality)
- **UI Updates**: 60 FPS smooth animations
- **Memory Usage**: ~500MB-2GB (depending on model)

### System Requirements

See [System Requirements](../getting-started/requirements.md) for details.

## 🔄 Updates

### Auto-Update

- **Automatic**: App checks for updates on startup
- **Notifications**: You'll be notified when updates are available
- **One-Click**: Update with a single click
- **Safe**: Old version kept as backup

### Release Channels

- **Stable**: Tested, production-ready releases
- **Beta**: Early access to new features
- **Dev**: Latest development build (may be unstable)

## 🆘 Troubleshooting

### Common Issues

**App won't start:**
- Check system requirements
- Update to latest version
- Check logs: `~/.kawai/logs/`

**AI not responding:**
- Check if model is downloaded
- Restart the app
- Check disk space (models need 4-8GB)

**Wallet not connecting:**
- Check wallet is unlocked
- Switch to Monad network
- Check network connection

**Slow performance:**
- Close unused chats
- Clear cache (Settings → Advanced)
- Check GPU drivers are updated

### Getting Help

- **FAQ**: [View common questions](../faq/general.md)
- **Community**: Join our [Discord](https://discord.gg/kawai)
- **Support**: [Contact us](../support/contact.md)

## 🚀 Next Steps

Ready to dive deeper? Check out these guides:

- [Wallet Setup](wallet-setup.md) - Create your built-in wallet
- [Deposit USDT](deposit.md) - Add funds to your account
- [Using AI Chat](ai-chat.md) - Master AI conversations
- [Image Generation](image-generation.md) - Create stunning images
- [Free Trial](free-trial.md) - Claim your 5-10 USDT bonus

---

**Need more help?** Check our [FAQ](../faq/general.md) or [contact support](../support/contact.md).

