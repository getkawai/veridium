# Installation & Security FAQ

Common questions about installing Kawai and security warnings.

## Windows Installation

### Why does Windows show a security warning?

**Short answer**: Kawai is a new application that hasn't built "reputation" with Microsoft yet.

**Long answer**:

Windows SmartScreen is a security feature that warns users about applications that:

1. **Don't have a code signing certificate** (costs $200-600/year)
2. **Haven't been downloaded by many users** (need thousands of downloads)
3. **Are new or recently updated** (each version needs reputation)

This is **normal for new open source projects**. Many legitimate applications show this warning initially, including:

- Early versions of VS Code
- Discord (in early days)
- Many indie games and tools
- Open source developer tools

### Is Kawai safe to install?

**Yes!** Here's why you can trust Kawai:

✅ **Open Source**: All code is public on [GitHub](https://github.com/kawai-network/veridium)

✅ **Verifiable**: You can review the code yourself or ask developers to audit it

✅ **Community**: Active Discord community with real users

✅ **Transparent**: We're honest about why the warning appears

✅ **No Malware**: Scan with any antivirus - it's clean

### How do I install despite the warning?

**Step-by-step guide:**

1. **Download** from official GitHub Releases only
   - ✅ `https://github.com/kawai-network/veridium/releases`
   - ❌ Don't download from random websites

2. **See the warning**: "Windows protected your PC"
   - This is expected - don't panic!

3. **Click "More info"**
   - Small text link at bottom of warning

4. **Click "Run anyway"**
   - Button appears after clicking "More info"

5. **Allow UAC prompt**
   - Click "Yes" when Windows asks for permission

6. **Install normally**
   - Follow the installation wizard

### Can I verify the download is safe?

**Yes! Verify the file hash:**

```powershell
# Open PowerShell in download folder
Get-FileHash Kawai-{VERSION}-windows-amd64.zip -Algorithm SHA256
```

Compare the output with the hash in `checksums.txt` from the release. If they match, the file is authentic and hasn't been tampered with.

### When will the warning go away?

We're working on it! Here's our plan:

**Phase 1** (Current): Unsigned builds with clear documentation

**Phase 2** (After 1000+ users): Purchase code signing certificate ($200-300/year)
- Warning will still appear for 3-6 months while building reputation
- But shows our verified identity

**Phase 3** (After validation): Upgrade to EV certificate ($400-600/year)
- Instant trust, no warning
- Requires registered company

**Phase 4** (Future): Microsoft Store version
- No warning at all
- Automatic updates
- Easier installation

### What if my antivirus blocks it?

Some antivirus software is extra cautious with unsigned applications.

**Solutions:**

1. **Temporarily disable** real-time protection during installation
2. **Add exception** for Kawai installer and installation folder
3. **Whitelist** the application after scanning
4. **Use Windows Defender** only (usually less aggressive)

**Popular antivirus compatibility:**

| Antivirus | Status | Notes |
|-----------|--------|-------|
| Windows Defender | ✅ Works | May show SmartScreen |
| Avast | ⚠️ May block | Add exception |
| AVG | ⚠️ May block | Add exception |
| Kaspersky | ⚠️ May block | Add exception |
| Norton | ⚠️ May block | Add exception |
| Bitdefender | ✅ Usually OK | Scan first |

### Can I use Microsoft Store instead?

We're working on a Microsoft Store version! Benefits:

- ✅ No SmartScreen warning
- ✅ Automatic updates
- ✅ Easier installation
- ✅ Trusted distribution

**Status**: In development. Join our [Discord](https://discord.gg/kawai) for updates.

## macOS Installation

### Why does macOS say "unidentified developer"?

Similar to Windows, macOS Gatekeeper protects users from unsigned applications.

**Solution:**

1. **Right-click** (or Control+click) on Kawai.app
2. **Select "Open"** from menu
3. **Click "Open"** in dialog
4. **Or**: System Preferences → Security & Privacy → Click "Open Anyway"

### When will macOS signing be available?

**Soon!** macOS signing is our top priority:

- **Cost**: $99/year (Apple Developer account)
- **Timeline**: Implementing in next release
- **Benefit**: No Gatekeeper warning

We're prioritizing macOS because:
- Gatekeeper is stricter than SmartScreen
- User experience impact is higher
- Cost is lower ($99 vs $200-600)

## Linux Installation

### What format is the Linux build?

We provide a `.tar.gz` archive containing the Kawai binary. This works on all Linux distributions.

**To install:**

```bash
# Extract the archive
tar -xzf Kawai-{VERSION}-linux-amd64.tar.gz

# Make executable (if needed)
chmod +x Kawai

# Run the application
./Kawai

# Optional: Install system-wide
sudo mv Kawai /usr/local/bin/
```

### Why isn't there a .deb or .rpm package?

We currently provide `.tar.gz` for maximum compatibility. Package formats coming soon:

**Planned:**
- ✅ .tar.gz (available now)
- 🔄 .deb package (in progress)
- 🔄 .rpm package (in progress)
- 🔄 AppImage (planned)
- 🔄 Flatpak (planned)
- 🔄 Snap (planned)

### Do I need to install dependencies?

Most modern Linux distributions include required dependencies. If you get errors:

**Ubuntu/Debian:**
```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
```

**Fedora/RHEL:**
```bash
sudo dnf install gtk3 webkit2gtk4.1
```

**Arch:**
```bash
sudo pacman -S gtk3 webkit2gtk-4.1
```

## General Security

### How do I know this isn't malware?

**Multiple ways to verify:**

1. **Check the source code**
   - All code is open source on GitHub
   - Review it yourself or ask a developer friend

2. **Scan with antivirus**
   - Use VirusTotal.com
   - Scan with your local antivirus
   - Should show 0 detections

3. **Verify file hash**
   - Compare SHA256 hash with GitHub Releases
   - Ensures file hasn't been tampered with

4. **Check community**
   - Join our Discord
   - Ask existing users
   - Read reviews and feedback

5. **Monitor behavior**
   - Use Process Monitor (Windows)
   - Check network connections
   - Review file system access

### What data does Kawai collect?

**We collect minimal data:**

✅ **Collected:**
- Wallet addresses (public blockchain data)
- Usage statistics (anonymous)
- Reward calculations (for claiming)
- Error logs (for debugging)

❌ **Never collected:**
- Private keys (stored locally only)
- Personal information
- AI conversation content
- Browsing history
- Location data

### Is my wallet safe?

**Yes!** Your wallet is self-custodial:

- 🔐 **Private keys** stored on your device only
- 🔒 **Encrypted** with your password
- 🚫 **Never sent** to our servers
- ✅ **You control** your funds completely

We **cannot** access your wallet or funds. You're in full control.

### Can I review the code?

**Absolutely!** We're fully open source:

- **Repository**: [github.com/kawai-network/veridium](https://github.com/kawai-network/veridium)
- **License**: Open source (check repo for details)
- **Contributions**: Welcome! Submit PRs
- **Audits**: Community audits encouraged

### What if I'm still concerned?

**We understand security concerns!** Here's what you can do:

1. **Wait for signed version**
   - Coming in next few months
   - Will remove security warnings

2. **Use in VM first**
   - Test in virtual machine
   - Verify behavior before main system

3. **Start with testnet**
   - No real money at risk
   - Test all features safely

4. **Ask the community**
   - Join our Discord
   - Ask questions
   - Get real user feedback

5. **Contact us**
   - Email: security@getkawai.com
   - Discord: [discord.gg/kawai](https://discord.gg/kawai)
   - GitHub: Open an issue

## Troubleshooting

### Installation fails completely

**Try these steps:**

1. **Run as Administrator** (Windows)
2. **Disable antivirus** temporarily
3. **Check disk space** (need 500MB+)
4. **Update Windows/macOS** to latest version
5. **Download again** (file might be corrupted)
6. **Check system requirements** in [Requirements](../getting-started/requirements.md)

### App crashes on startup

**Common causes:**

1. **Missing dependencies** (Linux)
   - Install required libraries
   - Check error messages

2. **Corrupted installation**
   - Uninstall completely
   - Delete app data folder
   - Reinstall fresh

3. **Conflicting software**
   - Close other apps
   - Disable VPN temporarily
   - Check firewall settings

4. **Outdated system**
   - Update OS to latest version
   - Update graphics drivers
   - Install system updates

### Can't connect to network

**Check these:**

1. **Internet connection** working
2. **Firewall** not blocking Kawai
3. **VPN** not interfering
4. **Proxy** settings correct
5. **Antivirus** not blocking network

**Firewall exceptions needed:**

- **Outbound**: Port 443 (HTTPS)
- **Outbound**: Port 8545 (Ethereum RPC)
- **Outbound**: Port 30303 (P2P, optional)

## Getting More Help

### Where can I get support?

**Multiple channels:**

1. **Discord** (fastest): [discord.gg/kawai](https://discord.gg/kawai)
   - Community help
   - Real-time chat
   - Voice support

2. **GitHub Issues**: [github.com/kawai-network/veridium/issues](https://github.com/kawai-network/veridium/issues)
   - Bug reports
   - Feature requests
   - Technical discussions

3. **Email**: support@getkawai.com
   - Detailed issues
   - Private concerns
   - Business inquiries

4. **Documentation**: [docs.getkawai.com](https://docs.getkawai.com)
   - Comprehensive guides
   - FAQ sections
   - Video tutorials

### What information should I provide?

**When asking for help, include:**

1. **Operating System**: Windows 10/11, macOS version, Linux distro
2. **App Version**: Check in About section
3. **Error Message**: Exact text or screenshot
4. **Steps to Reproduce**: What you did before error
5. **System Specs**: RAM, CPU, disk space
6. **Logs**: Check app data folder for log files

**Where to find logs:**

- **Windows**: `%APPDATA%\Kawai\logs`
- **macOS**: `~/Library/Application Support/Kawai/logs`
- **Linux**: `~/.config/Kawai/logs`

---

**Still have questions?** Check our other FAQs:

- [General FAQ](general.md) - About Kawai and features
- [Rewards FAQ](rewards.md) - About earning and claiming
- [Technical FAQ](technical.md) - About blockchain and tech

Or [contact support](../support/contact.md) directly!
