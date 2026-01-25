# Ephemeral Public Repo Release Strategy

## Concept

Create temporary public repo → Build via GitHub Actions → Download artifacts → Delete repo

This allows us to:
- ✅ Use free GitHub Actions (2000 min/month for public repos)
- ✅ Keep main repo private (secrets safe)
- ✅ Automated multi-platform builds (macOS, Linux, Windows)
- ✅ No permanent public exposure

## Workflow

```
┌─────────────────┐
│ Private Repo    │
│ (veridium)      │
└────────┬────────┘
         │
         │ 1. Copy & Clean
         ▼
┌─────────────────┐
│ Temp Directory  │
│ (kawai-0.1.0)   │
│ - Remove cmd/   │
│ - Remove .env   │
│ - Remove secrets│
└────────┬────────┘
         │
         │ 2. git init
         │ 3. gh repo create --public
         ▼
┌─────────────────┐
│ Public Repo     │
│ (ephemeral)     │
│ kawai-build-*   │
└────────┬────────┘
         │
         │ 4. Push tag
         │ 5. Trigger GitHub Actions
         ▼
┌─────────────────┐
│ GitHub Actions  │
│ - macOS build   │
│ - Linux build   │
│ - Windows build │
└────────┬────────┘
         │
         │ 6. Create Release
         ▼
┌─────────────────┐
│ Download        │
│ Artifacts       │
└────────┬────────┘
         │
         │ 7. Upload to R2
         ▼
┌─────────────────┐
│ Cloudflare R2   │
│ (permanent)     │
└────────┬────────┘
         │
         │ 8. gh repo delete
         ▼
┌─────────────────┐
│ Cleanup         │
│ (repo deleted)  │
└─────────────────┘
```

## Usage

### Prerequisites

```bash
# 1. Install GitHub CLI
brew install gh

# 2. Install AWS CLI (for R2)
brew install awscli

# 3. Add credentials to .env file
cat >> .env <<EOF
# GitHub Personal Access Token (for gh CLI)
GH_TOKEN=your_github_personal_access_token

# Cloudflare R2 credentials
R2_ACCESS_KEY_ID=your_key
R2_SECRET_ACCESS_KEY=your_secret
R2_ENDPOINT_URL=https://your-account.r2.cloudflarestorage.com
EOF
```

**Get GitHub Token**: https://github.com/settings/tokens (need `repo` and `delete_repo` scopes)

### Release Command

```bash
# Release version 0.1.0
./scripts/ephemeral-release.sh 0.1.0

# Or use current version from build/config.yml
./scripts/ephemeral-release.sh
```

### What Happens

1. **Copy & Clean** (5 seconds)
   - Copy entire repo to `../kawai-0.1.0/`
   - Remove `cmd/`, `.env*`, secrets
   - Remove `.git` directory

2. **Initialize Git** (5 seconds)
   - `git init`
   - `git add .`
   - `git commit -m "Release v0.1.0"`

3. **Create Public Repo** (10 seconds)
   - `gh repo create kawai-network/kawai-build-0.1.0 --public`
   - Repo name includes version for uniqueness

4. **Push & Tag** (10 seconds)
   - `git push origin master`
   - `git tag v0.1.0`
   - `git push origin v0.1.0`

5. **GitHub Actions Build** (15-20 minutes)
   - macOS: ~10 minutes
   - Linux: ~5 minutes
   - Windows: ~8 minutes
   - Parallel execution

6. **Download Artifacts** (1 minute)
   - `gh release download v0.1.0`
   - Downloads to `releases/v0.1.0/`

7. **Upload to R2** (2 minutes)
   - `aws s3 sync` to Cloudflare R2
   - Public URL: `https://storage.getkawai.com/v0.1.0/`

8. **Cleanup** (5 seconds)
   - `gh repo delete kawai-network/kawai-build-0.1.0 --yes`
   - `rm -rf ../kawai-0.1.0/`

**Total Time: ~20-25 minutes**

## Files Excluded from Public Repo

```
cmd/                      # Dev tools, admin tools
.env*                     # All environment files
contracts/.env*           # Contract environment files
data/                     # Local databases
.brv/                     # Kiro IDE metadata
.kiro/                    # Kiro settings
backend-dev.log           # Dev logs
contract-*.log            # Contract logs
verify.json               # Verification data
.circleci/                # CircleCI config (not needed)
```

## Files Included in Public Repo

```
main.go                   # Entry point
go.mod, go.sum           # Dependencies
internal/                # Core app logic
pkg/                     # Reusable packages
frontend/                # UI
contracts/contracts/     # Smart contracts (Solidity only)
build/                   # Build configs
scripts/                 # Build scripts (except ephemeral-release.sh)
.github/workflows/       # CI/CD
Makefile, Taskfile.yml   # Build automation
.gitignore               # Git ignore rules
README.md                # Documentation
```

## Security

### What's Safe to Publish

- ✅ Source code (Go, TypeScript, Solidity)
- ✅ Build scripts
- ✅ Documentation
- ✅ Smart contract addresses (already on blockchain)

### What's Never Published

- 🔒 API keys (OpenRouter, Gemini, etc)
- 🔒 Private keys (ADMIN_PRIVATE_KEY)
- 🔒 R2 credentials
- 🔒 Cloudflare tokens
- 🔒 Telegram bot tokens
- 🔒 Dev tools in `cmd/`

### GitHub Secrets

CI/CD uses GitHub Secrets (set via repo settings):
```
R2_ACCESS_KEY_ID
R2_SECRET_ACCESS_KEY
R2_ACCOUNT_ID
ADMIN_PRIVATE_KEY (for contract verification)
```

## Troubleshooting

### Build Fails

```bash
# Check GitHub Actions logs
gh run list --repo kawai-network/kawai-build-0.1.0
gh run view <run-id> --repo kawai-network/kawai-build-0.1.0 --log
```

### Repo Not Deleted

```bash
# Manual cleanup
gh repo delete kawai-network/kawai-build-0.1.0 --yes
```

### Artifacts Not Downloaded

```bash
# Manual download
gh release download v0.1.0 \
  --repo kawai-network/kawai-build-0.1.0 \
  --dir releases/v0.1.0
```

### R2 Upload Fails

```bash
# Check credentials
aws s3 ls s3://kawai/ --endpoint-url $R2_ENDPOINT_URL

# Manual upload
aws s3 sync releases/v0.1.0/ s3://kawai/v0.1.0/ \
  --endpoint-url $R2_ENDPOINT_URL
```

## Cost Analysis

### GitHub Actions (Public Repo)

- **Free Tier**: 2,000 minutes/month
- **macOS**: 10x multiplier = 200 minutes effective
- **Per Release**: ~20 minutes = 200 credits
- **Releases/Month**: ~10 releases

### Cloudflare R2

- **Storage**: $0.015/GB/month
- **Egress**: Free (Class A operations)
- **Per Release**: ~500 MB × 3 platforms = 1.5 GB
- **Cost**: ~$0.02/month per release

### Total Cost

- **GitHub Actions**: $0 (within free tier)
- **R2 Storage**: ~$0.20/month (10 releases)
- **Total**: ~$0.20/month

## Comparison with Alternatives

| Solution | Cost | Effort | Security | Automation |
|----------|------|--------|----------|------------|
| Ephemeral Public | $0 | Low | High | Full |
| Permanent Public | $0 | Medium | Medium | Full |
| Private + Self-hosted | $0 | High | High | Full |
| Private + GitHub Team | $4/mo | Low | High | Full |
| Manual Builds | $0 | High | High | None |

## Future Improvements

1. **Parallel Uploads**: Upload to R2 while build is running
2. **Caching**: Cache Go modules, npm packages between builds
3. **Notifications**: Telegram/Discord notification on completion
4. **Rollback**: Keep last N releases on R2
5. **Checksums**: Verify downloaded artifacts before upload

## References

- GitHub CLI: https://cli.github.com/
- GitHub Actions: https://docs.github.com/en/actions
- Cloudflare R2: https://developers.cloudflare.com/r2/
