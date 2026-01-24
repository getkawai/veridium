# GitHub Actions Workflows

This directory contains CI/CD workflows for automated testing, building, and releasing Kawai.

## Workflows

### 🚀 release.yml
**Trigger**: Version tags (v*.*.*) or manual dispatch

Automated release workflow that:
1. Creates a GitHub release draft
2. Builds for all platforms (macOS Universal, Linux amd64, Windows amd64)
3. Generates distribution archives with SHA256 checksums
4. Uploads all artifacts to the release
5. Publishes the release

**Platforms:**
- **macOS**: Universal Binary (Intel + Apple Silicon)
- **Linux**: amd64 binary + Debian package
- **Windows**: amd64 executable

### 📦 update-manifest.yml
**Trigger**: Release published

Generates `update.json` manifest for automatic updates:
- Extracts version and release info
- Downloads and parses checksums
- Creates update manifest with download URLs
- Uploads manifest to release

## Usage

### Creating a Release

#### Option 1: Tag-based (Recommended)
```bash
# Update version in build/config.yml first
git tag v1.0.0
git push origin v1.0.0
```

#### Option 2: Manual Dispatch
1. Go to Actions → Release workflow
2. Click "Run workflow"
3. Enter version (e.g., 1.0.0)
4. Click "Run workflow"

### Required Secrets

Add these secrets in GitHub repository settings:

- `R2_ACCOUNT_ID` - Cloudflare Account ID (from R2 dashboard)
- `R2_ACCESS_KEY_ID` - Cloudflare R2 Access Key ID
- `R2_SECRET_ACCESS_KEY` - Cloudflare R2 Secret Access Key

### Environment Variables

Configured in workflows:
- `GO_VERSION`: Go version (currently 1.23)
- `NODE_VERSION`: Node.js version (currently 20)
- `BUN_VERSION`: Bun version (currently 1.1.38)
- `R2_BUCKET`: R2 bucket name (default: `kawai-releases`)
- `R2_PUBLIC_URL`: R2 public URL (default: `https://releases.kawai.network`)

## Release Process

1. **Prepare Release**
   - Update version in `build/config.yml`
   - Update CHANGELOG.md
   - Commit changes

2. **Create Tag**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. **Automated Build**
   - GitHub Actions builds for all platforms
   - Creates draft release with artifacts
   - Generates checksums

4. **Review & Publish**
   - Review draft release on GitHub
   - Edit release notes if needed
   - Publish release

5. **Auto-Update Manifest**
   - `update.json` is automatically generated
   - Ready for automatic updates

## Artifacts

Each release includes:

**GitHub Release:**
- `Kawai-{version}-macos-universal.tar.gz` - macOS Universal Binary
- `Kawai-{version}-linux-amd64.tar.gz` - Linux binary
- `Kawai-{version}-windows-amd64.zip` - Windows executable
- `Kawai-{version}-linux-amd64.deb` - Debian package (optional)
- `checksums.txt` - SHA256 checksums for all files
- `update.json` - Auto-update manifest

**Cloudflare R2 (Public):**
- All artifacts uploaded to `https://releases.kawai.network/v{version}/`
- `update.json` also at `https://releases.kawai.network/latest/update.json`
- Public download URLs for end users

## Cloudflare R2 Setup

### 1. Create R2 Bucket
```bash
# Login to Cloudflare dashboard
# Navigate to R2 → Create bucket
# Bucket name: kawai
# Enable public access
```

### 2. Get R2 Credentials
```bash
# R2 → Manage R2 API Tokens → Create API Token
# Permissions: Object Read & Write
# Copy: Access Key ID, Secret Access Key, Endpoint URL
```

### 3. Configure Custom Domain (Optional)
```bash
# R2 bucket → Settings → Public Access
# Add custom domain: releases.kawai.network
# Update DNS CNAME record
```

### 4. Add GitHub Secrets
```bash
# Repository → Settings → Secrets and variables → Actions
R2_ACCOUNT_ID=ceab218751d33cd804878196ad7bef74
R2_ACCESS_KEY_ID=<your-access-key-id>
R2_SECRET_ACCESS_KEY=<your-secret-access-key>
```

### 5. Update Workflow Variables
Edit `.github/workflows/release.yml`:
```yaml
env:
  R2_BUCKET: 'your-bucket-name'
  R2_PUBLIC_URL: 'https://your-domain.com'
```

## Troubleshooting

### Build Fails on macOS
- Check Xcode Command Line Tools are installed
- Verify Wails v3 is compatible with macOS version

### Build Fails on Linux
- Ensure GTK3 and WebKit2GTK dependencies are available
- Check `libgtk-3-dev` and `libwebkit2gtk-4.0-dev` packages

### Build Fails on Windows
- Verify Go and Bun are properly installed
- Check Windows SDK is available

### Release Not Created
- Ensure tag follows `v*.*.*` format (e.g., v1.0.0)
- Check GitHub Actions permissions
- Verify GITHUB_TOKEN has release permissions

## Local Testing

Test workflows locally using [act](https://github.com/nektos/act):

```bash
# Install act
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Test release workflow
act -j build-macos --secret-file .secrets

# Test PR workflow
act pull_request
```

## Maintenance

### Updating Dependencies
Update versions in workflow files:
- Go version: `GO_VERSION`
- Node version: `NODE_VERSION`
- Bun version: `BUN_VERSION`
- Foundry: Uses `nightly` tag (auto-updates)

### Adding New Platforms
1. Add new job in `release.yml`
2. Configure platform-specific build steps
3. Update artifact upload steps
4. Update `update-manifest.yml` to include new platform

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Wails v3 Documentation](https://v3alpha.wails.io/)
- [Foundry Book](https://book.getfoundry.sh/)
