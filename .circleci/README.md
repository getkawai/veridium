# CircleCI Release Pipeline

Automated multi-platform release pipeline for Kawai using CircleCI.

## Overview

The CI/CD pipeline builds Kawai for multiple platforms:
- **Linux (amd64)** - Built on CircleCI
- **Windows (amd64)** - Built on CircleCI  
- **macOS (Universal)** - Built locally, uploaded to R2

## Workflow

```
1. Build macOS locally (see below)
2. Push tag (e.g., v0.1.0)
3. CircleCI builds Linux + Windows
4. Finalize downloads macOS checksum from R2
5. Combine all checksums
6. Release published
```

## Building macOS Locally

### Prerequisites
```bash
# Install dependencies
brew install go bun awscli

# Install Wails
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# Set R2 credentials (add to ~/.zshrc or ~/.bashrc)
export R2_ACCESS_KEY_ID="a71e802dd7c1ab8cf407ffb937cdf6a8"
export R2_SECRET_ACCESS_KEY="0e3ce0d92faa9b337c83131efc7a4a64bb6f313171c309d5cb9a0fb76926d0ca"
export R2_ENDPOINT_URL="https://ceab218751d33cd804878196ad7bef74.r2.cloudflarestorage.com"
```

### Build and Upload
```bash
# Build for specific version
./scripts/build-macos-release.sh 0.1.0

# Or use current git tag
./scripts/build-macos-release.sh
```

The script will:
1. Update version in `build/config.yml`
2. Build macOS Universal Binary
3. Create tar.gz archive
4. Generate SHA256 checksum
5. Upload to Cloudflare R2

### Manual Upload (if script fails)
```bash
cd bin
aws s3 cp Kawai-0.1.0-macos-universal.tar.gz s3://kawai/v0.1.0/ --endpoint-url $R2_ENDPOINT_URL
aws s3 cp checksums-macos.txt s3://kawai/v0.1.0/checksums-macos.txt --endpoint-url $R2_ENDPOINT_URL
```

## Release Process

### Quick Release (Recommended)
```bash
# One command to do everything
./scripts/release.sh 0.1.0
```

This script will:
1. Update version in `build/config.yml`
2. Commit version change
3. Create and push tag `v0.1.0`
4. Trigger CircleCI builds (Linux + Windows)
5. Build macOS locally in parallel
6. Upload macOS binary + checksum to R2
7. CircleCI finalize will combine all checksums

### Manual Release

#### 1. Build macOS First
```bash
# Build and upload macOS binary
./scripts/build-macos-release.sh 0.1.0
```

#### 2. Create and Push Tag
```bash
git tag v0.1.0
git push origin v0.1.0
```

### 3. CircleCI Builds Automatically
CircleCI will:
- Build Linux binary
- Build Windows binary
- Download macOS checksum from R2
- Combine all checksums
- Upload to R2
- Publish GitHub release

## Artifacts

Each release includes:
- `Kawai-{version}-macos-universal.tar.gz` - macOS Universal Binary
- `Kawai-{version}-linux-amd64.tar.gz` - Linux binary
- `Kawai-{version}-windows-amd64.zip` - Windows executable
- `checksums.txt` - SHA256 checksums
- `update.json` - Auto-update manifest

All uploaded to:
- **Cloudflare R2**: `https://releases.kawai.network/v{version}/`
- **GitHub Releases**: Draft release (auto-published after finalize)

## Resource Classes

- **macOS**: `macos.m1.medium.gen1` (Apple Silicon)
- **Linux**: `large` (4 vCPU, 8GB RAM)
- **Windows**: `windows.large` (4 vCPU, 15GB RAM)

## Cost Estimation

CircleCI Free Tier: 6,000 credits/month

Estimated credits per release:
- macOS build: ~500 credits (10 min × 50 credits/min)
- Linux build: ~100 credits (10 min × 10 credits/min)
- Windows build: ~200 credits (10 min × 20 credits/min)
- **Total: ~800 credits per release**

Can do ~7 releases/month on free tier.

## Troubleshooting

### Build fails on macOS
- Check Xcode version compatibility
- Verify Wails v3 supports the macOS version

### Build fails on Linux
- Ensure GTK3/WebKit2GTK dependencies are available
- Check Go version compatibility

### Build fails on Windows
- Verify Windows Server 2022 compatibility
- Check Go and Node versions

### R2 upload fails
- Verify R2 credentials in environment variables
- Check bucket name and endpoint URL
- Ensure bucket has public access enabled

## Monitoring

View pipeline status:
```bash
# Via CircleCI CLI
circleci run list

# Via web
https://app.circleci.com/pipelines/github/kawai-network/veridium
```

## Migration from GitHub Actions

GitHub Actions workflows in `.github/workflows/` are kept for reference but not used. CircleCI config in `.circleci/config.yml` is now the primary CI/CD pipeline.

Key differences:
- CircleCI uses credits instead of minutes
- Different executor types (docker, machine, macos)
- Workspace persistence instead of artifacts between jobs
- Manual trigger via API instead of workflow_dispatch
