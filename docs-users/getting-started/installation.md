# Installation Guide

Get Kawai DeAI Network installed on your system in just a few minutes.

## 🖥️ Platform-Specific Installation

### macOS Installation

1. **Download** the `.tar.gz` file for macOS
2. **Extract** the archive (double-click or use `tar -xzf`)
3. **Drag** Kawai.app to Applications folder
4. **Launch** from Applications or Spotlight
5. **Allow** the app if macOS shows security warning:
   - Go to System Settings → Privacy & Security
   - Click "Open Anyway" for Kawai

!!! warning "macOS Security"
    First launch may show "unidentified developer" warning. This is normal for new applications. Click "Open" to proceed.

### Windows Installation

1. **Download** the `.zip` file for Windows
2. **Extract** the archive to a folder (e.g., `C:\Program Files\Kawai`)
3. **Run** `Kawai.exe` from the extracted folder
4. **Create** desktop shortcut (optional - right-click Kawai.exe → Send to → Desktop)
5. **Launch** from shortcut or Start Menu

!!! warning "Windows SmartScreen Warning"
    **You will see a security warning** when running Kawai. This is normal for new applications.
    
    **Why does this happen?**
    
    - Windows SmartScreen protects users from unknown software
    - New applications need to build "reputation" with Microsoft
    - This doesn't mean the software is unsafe - it's just new
    - We're working on getting our application signed to remove this warning
    
    **How to run safely:**
    
    === "Step-by-Step"
        1. **Download** from official source only (GitHub Releases or storage.getkawai.com)
        2. **Extract** the ZIP file
        3. **Verify** the file hash (optional but recommended)
        4. When you see "Windows protected your PC":
           - Click **"More info"**
           - Click **"Run anyway"**
        5. **Allow** User Account Control (UAC) prompt if shown
        6. **Run** Kawai.exe normally
    
    === "Visual Guide"
        **Step 1**: You'll see this warning
        ```
        Windows protected your PC
        Microsoft Defender SmartScreen prevented an unrecognized app from starting.
        Running this app might put your PC at risk.
        ```
        
        **Step 2**: Click "More info"
        
        **Step 3**: Click "Run anyway" button
        
        **Step 4**: Click "Yes" on UAC prompt if shown
    
    === "Verify Download"
        For extra security, verify the file hash:
        
        ```powershell
        # In PowerShell
        Get-FileHash Kawai-{VERSION}-windows-amd64.zip -Algorithm SHA256
        ```
        
        Compare with hash from checksums.txt in the release.
    
    !!! info "Why isn't Kawai signed yet?"
        Code signing certificates cost $200-600/year. As an early-stage open source project, we're prioritizing development over certificates. Once we have more users and funding, we'll get proper code signing to remove this warning.
    
    !!! tip "Alternative: Microsoft Store"
        We're working on a Microsoft Store version that won't show this warning. Join our Discord for updates!

### Linux Installation

1. **Download** the `.tar.gz` file for Linux
2. **Extract** the archive:
   ```bash
   tar -xzf Kawai-{VERSION}-linux-amd64.tar.gz
   ```
3. **Make executable** (if needed):
   ```bash
   chmod +x Kawai
   ```
4. **Run** the application:
   ```bash
   ./Kawai
   ```
5. **Optional**: Move to `/usr/local/bin` for system-wide access:
   ```bash
   sudo mv Kawai /usr/local/bin/
   ```

!!! note "Linux Dependencies"
    Most modern Linux distributions include required dependencies. If you encounter issues, install:
    ```bash
    # Ubuntu/Debian
    sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
    
    # Fedora
    sudo dnf install gtk3 webkit2gtk4.1
    
    # Arch
    sudo pacman -S gtk3 webkit2gtk-4.1
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
Follow the [Wallet Setup Guide](../user-guide/wallet-setup.md) to create your built-in wallet.

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

??? question "Windows SmartScreen blocks installation"
    **This is expected behavior for new applications.**
    
    **Solution:**
    1. Click "More info" on the warning
    2. Click "Run anyway"
    3. Allow UAC prompt
    
    **Still blocked?**
    - Check Windows Defender settings
    - Temporarily disable real-time protection
    - Add exception for Kawai installer
    - Download from official GitHub only
    
    **Security concerns?**
    - Verify file hash from GitHub Releases
    - Check our open source code
    - Scan with your antivirus
    - Join Discord to ask the community

??? question "App won't start"
    **Possible causes:**
    - Insufficient permissions
    - Missing dependencies
    - Corrupted download
    - Antivirus blocking
    
    **Solutions:**
    - Run as administrator (Windows)
    - Install missing libraries (Linux)
    - Re-download and reinstall
    - Add antivirus exception

??? question "Wallet connection fails"
    **Possible causes:**
    - MetaMask not installed
    - Wrong network configuration
    - Browser compatibility
    
    **Solutions:**
    - Create wallet in app
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

1. **Check** our [Installation & Security FAQ](../faq/installation.md) for detailed solutions
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