# CI/CD Plan for KAWAI Node Installer

## Current R2 Storage State

### Bucket Info
- **Bucket Name**: kawai
- **Public URL**: https://storage.getkawai.com
- **R2 Endpoint**: https://ceab218751d33cd804878196ad7bef74.r2.cloudflarestorage.com

### Folder Structure
```
node/
├── v51aec45/
│   ├── checksums-darwin-arm64.txt
│   ├── checksums-linux-amd64.txt
│   ├── kawai-node-51aec45-darwin-arm64.tar.gz (36.4 MB)
│   ├── kawai-node-51aec45-linux-amd64.tar.gz (42.8 MB)
│   └── kawai-node (119.3 MB)
└── v69e4e5e/
    └── ...
```

### Naming Convention
- **Folder**: node/v{7-char-commit-sha}/
- **Archive**: kawai-node-{commit-sha}-{platform}-{arch}.tar.gz
- **Platforms**: linux/amd64, darwin/arm64
- **Version Format**: v{commit-sha} (e.g., v51aec45)

### Access Requirements
To check R2 contents locally:
```bash
# Set environment variables from .env
export $(grep -v '^#' .env | xargs)

# Unset AWS session tokens (if using AWS CLI)
unset AWS_SESSION_TOKEN AWS_SECURITY_TOKEN

# Check contents
export AWS_ACCESS_KEY_ID=$R2_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=$R2_SECRET_ACCESS_KEY
aws s3 ls s3://kawai/node/ --endpoint-url $R2_ENDPOINT_URL
```

### Notes
- R2 public URL does not support directory listing (returns 404)
- AWS CLI requires session token unset to work with R2 credentials
- Version format uses 7-character git commit SHA

## Current State

### Repositories
1. **veridium** (this repo) - Contains:
   - `.github/workflows/release-node.yml` - Build and upload to R2
   - `cmd/r2upload/main.go` - R2 upload utility
   - `scripts/check-r2.sh` - R2 checker

2. **kawai-website** - Contains:
   - `node/install.sh` - Installer script (hardcoded version)
   - `node/index.html` - Landing page

### Current Flow
```
Push to feature/release-node
    ↓
GitHub Actions builds binaries
    ↓
Upload to R2: node/v{commit-sha}/
    ↓
Manual update install.sh with new version
```

## Problems
1. **Manual version update** - Developer must manually update `install.sh` after each release
2. **No version tracking** - No single source of truth for "latest" version
3. **Two repos to manage** - Release in veridium, installer in kawai-website

## Proposed CI/CD Solutions

### Option 1: Auto-update installer via GitHub Actions (Recommended)

**Flow:**
```
Push to feature/release-node
    ↓
Build binaries → Upload to R2
    ↓
Auto-update kawai-website/install.sh
    ↓
Auto-commit and push to kawai-website
```

**Implementation:**
1. Add step in `release-node.yml` to:
   - Checkout kawai-website repo
   - Update `DEFAULT_LATEST_VERSION` in install.sh
   - Commit and push

**Pros:**
- Fully automated
- No manual steps
- Single push triggers everything

**Cons:**
- Need GitHub token with write access to kawai-website
- Cross-repo workflow complexity

---

### Option 2: versions.txt in R2

**Flow:**
```
Push to feature/release-node
    ↓
Build binaries → Upload to R2
    ↓
Upload versions.txt with latest version
    ↓
Installer reads versions.txt from R2
```

**Implementation:**
1. Workflow uploads `node/versions.txt` containing latest version
2. Installer fetches `https://storage.getkawai.com/node/versions.txt`

**Pros:**
- No hardcoded version
- No cross-repo access needed
- Simple to implement

**Cons:**
- Extra HTTP request on every install
- versions.txt could be out of sync

---

### Option 3: GitHub Releases as Source of Truth

**Flow:**
```
Push to feature/release-node
    ↓
Build binaries → Create GitHub Release
    ↓
Upload to R2 (mirror)
    ↓
Installer queries GitHub API for latest release
```

**Implementation:**
1. Create GitHub release with tag `node-v{version}`
2. Installer queries `api.github.com/repos/kawai-network/veridium/releases/latest`
3. Parse version from tag

**Pros:**
- GitHub is source of truth
- Can use release notes
- Standard practice

**Cons:**
- GitHub API rate limits
- Dependency on GitHub availability

---

### Option 4: Cloudflare Worker API

**Flow:**
```
Push to feature/release-node
    ↓
Build binaries → Upload to R2
    ↓
Cloudflare Worker serves latest version
    ↓
Installer queries worker
```

**Implementation:**
1. Create Cloudflare Worker at `api.getkawai.com/node/latest`
2. Worker reads R2 directory and returns latest version
3. Installer fetches from API

**Pros:**
- Fast (edge cached)
- Can add logic (stable/beta channels)
- No hardcoded version

**Cons:**
- Need to maintain worker
- Extra infrastructure

---

## Recommendation

**Option 1 (Auto-update via GitHub Actions)** for immediate simplicity
**Option 2 (versions.txt)** as fallback if cross-repo access is problematic

## Implementation Plan

### Phase 1: Auto-update installer (Option 1)

Update `.github/workflows/release-node.yml`:

```yaml
- name: Update installer version
  run: |
    VERSION="${{ needs.prepare.outputs.version }}"
    
    # Checkout kawai-website
    git clone https://x-access-token:${{ secrets.KAWAI_WEBSITE_TOKEN }}@github.com/kawai-network/kawai-website.git /tmp/kawai-website
    
    # Update version in install.sh
    sed -i "s/DEFAULT_LATEST_VERSION=\"v[0-9a-f]*/DEFAULT_LATEST_VERSION=\"v${VERSION}/" /tmp/kawai-website/node/install.sh
    
    # Commit and push
    cd /tmp/kawai-website
    git config user.name "GitHub Actions"
    git config user.email "actions@github.com"
    git add node/install.sh
    git commit -m "Update installer to v${VERSION}"
    git push
```

### Phase 2: versions.txt fallback (Option 2)

Add to workflow:
```yaml
- name: Update versions.txt
  run: |
    echo "v${{ needs.prepare.outputs.version }}" > versions.txt
    # Upload to R2
```

Update installer to try versions.txt first, fallback to hardcoded.

---

## Next Steps

1. Choose preferred option
2. Generate GitHub token with kawai-website write access (if Option 1)
3. Update workflow
4. Test end-to-end

Which option do you prefer?
