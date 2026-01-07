# Installation Guide

Get Kawai DeAI Network installed on your system in just a few minutes.

## 📥 Download Options

### Official Releases
Download the latest version from our official sources:

=== "GitHub Releases"
    1. Visit [GitHub Releases](https://github.com/kawai-network/veridium/releases)
    2. Find the latest version
    3. Download for your platform:
       - `Kawai-DeAI-macos.dmg` (macOS)
       - `Kawai-DeAI-windows.exe` (Windows)
       - `Kawai-DeAI-linux.AppImage` (Linux)

=== "Direct Download"
    - [Download for macOS](https://github.com/kawai-network/veridium/releases/latest/download/Kawai-DeAI-macos.dmg)
    - [Download for Windows](https://github.com/kawai-network/veridium/releases/latest/download/Kawai-DeAI-windows.exe)
    - [Download for Linux](https://github.com/kawai-network/veridium/releases/latest/download/Kawai-DeAI-linux.AppImage)

## 🖥️ Platform-Specific Installation

### macOS Installation

1. **Download** the `.dmg` file
2. **Open** the downloaded file
3. **Drag** Kawai DeAI to Applications folder
4. **Launch** from Applications or Spotlight
5. **Allow** the app if macOS shows security warning:
   - Go to System Preferences → Security & Privacy
   - Click "Open Anyway" for Kawai DeAI

!!! warning "macOS Security"
    First launch may show "unidentified developer" warning. This is normal for new applications. Click "Open" to proceed.

### Windows Installation

1. **Download** the `.exe` file
2. **Run** the installer as Administrator
3. **Follow** the installation wizard
4. **Choose** installation directory (default recommended)
5. **Create** desktop shortcut (optional)
6. **Launch** from Start Menu or desktop

!!! tip "Windows Defender"
    Windows may show SmartScreen warning. Click "More info" → "Run anyway" to proceed.

### Linux Installation

1. **Download** the `.AppImage` file
2. **Make executable**:
   ```bash
   chmod +x Kawai-DeAI-linux.AppImage
   ```
3. **Run** the application:
   ```bash
   ./Kawai-DeAI-linux.AppImage
   ```
4. **Optional**: Install AppImageLauncher for better integration

!!! note "Linux Dependencies"
    Most modern Linux distributions include required dependencies. If you encounter issues, install:
    ```bash
    sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
    ```

## 🔧 First Launch Setup

### 1. Initial Configuration
When you first launch Kawai DeAI:

1. **Accept** terms of service
2. **Choose** data directory (default recommended)
3. **Configure** basic settings
4. **Wait** for initial setup to complete

### 2. Network Configuration
The app will automatically:

- Connect to Monad Testnet
- Download necessary configurations
- Initialize local AI models
- Prepare reward systems

### 3. Wallet Connection
Follow the [Wallet Setup Guide](../user-guide/wallet-setup.md) to connect your Web3 wallet.

## 🔄 Updates & Maintenance

### Automatic Updates
Kawai DeAI includes automatic update checking:

- **Notification** when updates are available
- **One-click** update installation
- **Backup** of user data before updates
- **Rollback** option if needed

### Manual Updates
To manually check for updates:

1. Open app settings
2. Go to "About" section
3. Click "Check for Updates"
4. Follow update prompts

### Data Backup
Your data is automatically backed up:

- **Wallet connections** (encrypted)
- **Chat history** (local storage)
- **Reward data** (synced with blockchain)
- **Settings** (local preferences)

## 🛠️ Advanced Installation

### Development Build
For developers and testers:

```bash
# Clone repository
git clone https://github.com/kawai-network/veridium.git
cd veridium

# Install dependencies
make install

# Build and run
make dev
```

### Custom Configuration
Advanced users can customize:

- **Data directory** location
- **Network endpoints** (for testing)
- **AI model** preferences
- **Performance** settings

## 🔍 Verification

### Verify Installation
After installation, verify everything works:

1. **Launch** the application
2. **Check** version in About section
3. **Test** wallet connection
4. **Try** AI chat with free trial
5. **Confirm** reward tracking works

### File Integrity
For security-conscious users:

1. **Check** file hashes (provided in releases)
2. **Verify** digital signatures
3. **Scan** with antivirus software
4. **Monitor** network connections

## 🚨 Troubleshooting

### Common Installation Issues

??? question "App won't start"
    **Possible causes:**
    - Insufficient permissions
    - Missing dependencies
    - Corrupted download
    
    **Solutions:**
    - Run as administrator (Windows)
    - Install missing libraries (Linux)
    - Re-download and reinstall

??? question "Wallet connection fails"
    **Possible causes:**
    - MetaMask not installed
    - Wrong network configuration
    - Browser compatibility
    
    **Solutions:**
    - Install MetaMask extension
    - Add Monad Testnet manually
    - Try different browser

??? question "AI services not working"
    **Possible causes:**
    - No internet connection
    - Firewall blocking requests
    - Server maintenance
    
    **Solutions:**
    - Check internet connection
    - Configure firewall exceptions
    - Wait and retry later

### Getting Help

If you encounter installation issues:

1. **Check** our [troubleshooting FAQ](../faq/troubleshooting.md)
2. **Search** existing GitHub issues
3. **Join** our Discord for community help
4. **Contact** support with system details

## 🎯 Next Steps

After successful installation:

1. **[Set up your wallet](../user-guide/wallet-setup.md)** for Web3 features
2. **[Claim your free trial](../user-guide/free-trial.md)** to test AI services
3. **[Explore the app](../user-guide/desktop-app.md)** and its features
4. **[Learn about rewards](../rewards/overview.md)** to start earning

---

**Installation complete!** 🎉 You're ready to experience the future of decentralized AI.

[Next: Wallet Setup →](../user-guide/wallet-setup.md)