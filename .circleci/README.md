# CircleCI Release Pipeline

Automated multi-platform release pipeline for Kawai using CircleCI.

## Setup

### Quick Setup (Automated)

1. **Follow project on CircleCI**
   - Go to https://app.circleci.com/
   - Click "Projects" → Find "veridium" → Click "Set Up Project"
   - Select "Use Existing Config" (we already have `.circleci/config.yml`)

2. **Run setup script**
   ```bash
   ./.circleci/setup.sh
   ```

3. **Add GITHUB_TOKEN manually**
   - Go to https://app.circleci.com/settings/project/github/kawai-network/veridium/environment-variables
   - Click "Add Environment Variable"
   - Name: `GITHUB_TOKEN`
   - Value: Your GitHub Personal Access Token

### Manual Setup

If the script doesn't work, add environment variables manually:

```
R2_ACCOUNT_ID=ceab218751d33cd804878196ad7bef74
R2_ACCESS_KEY_ID=a71e802dd7c1ab8cf407ffb937cdf6a8
R2_SECRET_ACCESS_KEY=0e3ce0d92faa9b337c83131efc7a4a64bb6f313171c309d5cb9a0fb76926d0ca
R2_ENDPOINT_URL=https://ceab218751d33cd804878196ad7bef74.r2.cloudflarestorage.com
GITHUB_TOKEN=<your-github-token>
```

### 2. Enable Build Processing

CircleCI Project Settings → Advanced → Enable build processing for tags

### 3. Trigger Release

**Option 1: Push tag (auto-trigger)**
```bash
git tag v1.0.0
git push origin v1.0.0
```

**Option 2: Manual trigger via API**
```bash
curl -X POST \
  -H "Circle-Token: $CIRCLECI_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "parameters": {
      "version": "1.0.0"
    }
  }' \
  https://circleci.com/api/v2/project/github/kawai-network/veridium/pipeline
```

## Pipeline Jobs

1. **prepare-release** - Create GitHub draft release
2. **build-macos** - Build macOS Universal Binary (M1 runner)
3. **build-linux** - Build Linux amd64 binary
4. **build-windows** - Build Windows amd64 executable
5. **finalize-release** - Combine checksums, generate manifest, publish release

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
