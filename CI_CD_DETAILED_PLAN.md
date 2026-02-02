# 🚀 CI/CD & DevOps - Detailed Implementation Plan

**Tanggal:** 2 Februari 2026  
**Status:** Implementation Ready  
**Target:** Complete CI/CD automation untuk Veridium project

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Current State Analysis](#current-state-analysis)
3. [GitHub Actions Workflows](#github-actions-workflows)
4. [Automated Release Process](#automated-release-process)
5. [Auto-update install.sh](#auto-update-installsh)
6. [Implementation Steps](#implementation-steps)
7. [Testing Strategy](#testing-strategy)
8. [Rollout Plan](#rollout-plan)

---

## 🎯 Overview

### Goals
- ✅ Automated testing pada setiap push/PR
- ✅ Automated linting dan code quality checks
- ✅ Automated contract testing dan deployment
- ✅ Automated release process untuk desktop app dan node server
- ✅ Auto-update installer script di kawai-website repo
- ✅ Zero-downtime deployment

### Success Metrics
| Metric | Current | Target |
|--------|---------|--------|
| Manual steps per release | 8-10 | 0 |
| Release time | 2-3 hours | 15 minutes |
| Test coverage | ~15% | 70%+ |
| Failed deployments | Unknown | <5% |
| Time to rollback | Manual | <5 minutes |

---

## 🔍 Current State Analysis

### Existing Workflows

#### 1. `.github/workflows/release-node.yml`
**Purpose:** Build dan release Kronk server (cmd/server)  
**Trigger:** Tag push `node-v*`  
**Status:** ✅ Working, needs enhancement

**Current Flow:**
```
Tag push (node-v1.2.3)
  ↓
Build for linux/darwin/windows
  ↓
Upload to R2 storage
  ↓
Create GitHub Release (draft)
  ↓
❌ Manual: Update install.sh
```

#### 2. `.github/workflows/release.yml`
**Purpose:** Build desktop app (Wails)  
**Status:** Needs review

### Gaps Identified
- ❌ No CI workflow untuk testing
- ❌ No linting workflow
- ❌ No contract testing workflow
- ❌ Manual installer update
- ❌ No automated changelog generation
- ❌ No deployment verification

---

## 🔧 GitHub Actions Workflows

### Workflow 1: CI Pipeline (New)

**File:** `.github/workflows/ci.yml`

```yaml
name: CI Pipeline

on:
  push:
    branches: [master, enhancement, 'feature/**']
  pull_request:
    branches: [master]

env:
  GO_VERSION: '1.25'
  NODE_VERSION: '20'

jobs:
  # ============================================================================
  # LINT JOBS
  # ============================================================================
  
  lint-go:
    name: Lint Go Code
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          args: --timeout=10m --config=.golangci.yml
          
      - name: Check go mod tidy
        run: |
          go mod tidy
          git diff --exit-code go.mod go.sum
  
  lint-frontend:
    name: Lint Frontend Code
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      
      - name: Install dependencies
        run: npm ci
      
      - name: ESLint
        run: npm run lint
      
      - name: TypeScript check
        run: npm run type-check
      
      - name: Prettier check
        run: npm run format:check

  # ============================================================================
  # TEST JOBS
  # ============================================================================
  
  test-go:
    name: Test Go Code
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: Run tests with race detector
        run: go test -race -coverprofile=coverage.out -covermode=atomic ./...
      
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
          flags: unittests
          name: codecov-umbrella
          token: ${{ secrets.CODECOV_TOKEN }}
      
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 70" | bc -l) )); then
            echo "❌ Coverage below 70% threshold"
            exit 1
          fi
  
  test-frontend:
    name: Test Frontend Code
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run tests
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./frontend/coverage/coverage-final.json
          flags: frontend
          token: ${{ secrets.CODECOV_TOKEN }}

  # ============================================================================
  # CONTRACT JOBS
  # ============================================================================
  
  test-contracts:
    name: Test Smart Contracts
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: contracts
    steps:
      - uses: actions/checkout@v4
        with:
          submodules: recursive
      
      - name: Install Foundry
        uses: foundry-rs/foundry-toolchain@v1
        with:
          version: nightly
      
      - name: Run Forge tests
        run: forge test -vvv
      
      - name: Run Forge coverage
        run: forge coverage --report lcov
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./contracts/lcov.info
          flags: contracts
          token: ${{ secrets.CODECOV_TOKEN }}
      
      - name: Gas snapshot
        run: forge snapshot --check
  
  lint-contracts:
    name: Lint Smart Contracts
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: contracts
    steps:
      - uses: actions/checkout@v4
        with:
          submodules: recursive
      
      - name: Install Foundry
        uses: foundry-rs/foundry-toolchain@v1
      
      - name: Check formatting
        run: forge fmt --check
      
      - name: Run Slither
        uses: crytic/slither-action@v0.3.0
        with:
          target: contracts/
          slither-args: --filter-paths "lib|test"
          fail-on: medium

  # ============================================================================
  # BUILD JOBS
  # ============================================================================
  
  build-backend:
    name: Build Backend
    needs: [lint-go, test-go]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: Build main app
        run: go build -v -o veridium .
      
      - name: Build server
        run: go build -v -o server ./cmd/server
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: binaries
          path: |
            veridium
            server
  
  build-frontend:
    name: Build Frontend
    needs: [lint-frontend, test-frontend]
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      
      - name: Install dependencies
        run: npm ci
      
      - name: Build
        run: npm run build
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: frontend-dist
          path: frontend/dist

  # ============================================================================
  # INTEGRATION TEST
  # ============================================================================
  
  integration-test:
    name: Integration Tests
    needs: [build-backend, build-frontend]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Download artifacts
        uses: actions/download-artifact@v4
      
      - name: Setup test environment
        run: |
          # Start local blockchain (anvil)
          # Deploy contracts
          # Start backend
          # Run integration tests
          echo "Integration tests placeholder"
      
      - name: Run integration tests
        run: go test -v ./tests/integration/...

  # ============================================================================
  # SECURITY SCAN
  # ============================================================================
  
  security-scan:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'
      
      - name: Upload Trivy results to GitHub Security
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'
      
      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out gosec-results.sarif ./...'
      
      - name: Upload gosec results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'gosec-results.sarif'
```

---

### Workflow 2: Release Desktop App (Enhanced)

**File:** `.github/workflows/release-desktop.yml`

```yaml
name: Release Desktop App

on:
  push:
    tags:
      - 'v*'  # v1.2.3

permissions:
  contents: write
  packages: write

env:
  GO_VERSION: '1.25'
  NODE_VERSION: '20'

jobs:
  prepare:
    name: Prepare Release
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.get_version.outputs.version }}
      changelog: ${{ steps.changelog.outputs.changelog }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for changelog
      
      - name: Get version
        id: get_version
        run: |
          VERSION="${GITHUB_REF#refs/tags/v}"
          echo "version=$VERSION" >> $GITHUB_OUTPUT
      
      - name: Generate changelog
        id: changelog
        uses: mikepenz/release-changelog-builder-action@v4
        with:
          configuration: ".github/changelog-config.json"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: v${{ steps.get_version.outputs.version }}
          name: Veridium v${{ steps.get_version.outputs.version }}
          body: ${{ steps.changelog.outputs.changelog }}
          draft: true
          prerelease: false

  build-macos:
    name: Build macOS
    needs: prepare
    runs-on: macos-latest
    strategy:
      matrix:
        arch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest
      
      - name: Build
        run: wails3 build -platform darwin/${{ matrix.arch }}
      
      - name: Create DMG
        run: |
          # Create DMG installer
          # Sign with Apple Developer ID
          # Notarize
          echo "DMG creation placeholder"
      
      - name: Upload to Release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: v${{ needs.prepare.outputs.version }}
          files: |
            build/bin/Veridium-darwin-${{ matrix.arch }}.dmg
            build/bin/Veridium-darwin-${{ matrix.arch }}.dmg.sha256

  build-windows:
    name: Build Windows
    needs: prepare
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest
      
      - name: Build
        run: wails3 build -platform windows/amd64
      
      - name: Create Installer
        run: |
          # Create NSIS installer
          # Sign with code signing certificate
          echo "Installer creation placeholder"
      
      - name: Upload to Release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: v${{ needs.prepare.outputs.version }}
          files: |
            build/bin/Veridium-windows-amd64.exe
            build/bin/Veridium-windows-amd64-installer.exe

  build-linux:
    name: Build Linux
    needs: prepare
    runs-on: ubuntu-latest
    strategy:
      matrix:
        arch: [amd64]
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev
      
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest
      
      - name: Build
        run: wails3 build -platform linux/${{ matrix.arch }}
      
      - name: Create packages
        run: |
          # Create .deb package
          # Create .rpm package
          # Create AppImage
          echo "Package creation placeholder"
      
      - name: Upload to Release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: v${{ needs.prepare.outputs.version }}
          files: |
            build/bin/veridium-linux-${{ matrix.arch }}.tar.gz
            build/bin/veridium_${{ needs.prepare.outputs.version }}_${{ matrix.arch }}.deb
            build/bin/veridium-${{ needs.prepare.outputs.version }}-1.${{ matrix.arch }}.rpm
            build/bin/Veridium-${{ needs.prepare.outputs.version }}-${{ matrix.arch }}.AppImage

  publish-release:
    name: Publish Release
    needs: [build-macos, build-windows, build-linux]
    runs-on: ubuntu-latest
    steps:
      - name: Publish Release
        uses: eregon/publish-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          release_id: ${{ needs.prepare.outputs.release_id }}
```

---

### Workflow 3: Release Node Server (Enhanced)

**File:** `.github/workflows/release-node.yml` (Enhanced version)

```yaml
name: Release Node Server

on:
  push:
    tags:
      - 'node-v*'

permissions:
  contents: write
  packages: write

env:
  GO_VERSION: '1.25'
  R2_BUCKET: 'kawai'
  R2_PUBLIC_URL: 'https://storage.getkawai.com'

jobs:
  prepare:
    name: Prepare Release
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.get_version.outputs.version }}
      short_sha: ${{ steps.get_sha.outputs.short_sha }}
    steps:
      - uses: actions/checkout@v4
      
      - name: Get version
        id: get_version
        run: |
          VERSION="${GITHUB_REF#refs/tags/node-v}"
          echo "version=$VERSION" >> $GITHUB_OUTPUT
      
      - name: Get short SHA
        id: get_sha
        run: |
          SHORT_SHA=$(git rev-parse --short=7 HEAD)
          echo "short_sha=$SHORT_SHA" >> $GITHUB_OUTPUT

  build:
    name: Build ${{ matrix.os }}-${{ matrix.arch }}
    needs: prepare
    runs-on: ${{ matrix.runner }}
    strategy:
      matrix:
        include:
          - os: linux
            arch: amd64
            runner: ubuntu-latest
          - os: linux
            arch: arm64
            runner: ubuntu-latest
          - os: darwin
            arch: amd64
            runner: macos-latest
          - os: darwin
            arch: arm64
            runner: macos-latest
          - os: windows
            arch: amd64
            runner: windows-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Build
        env:
          GOOS: ${{ matrix.os }}
          GOARCH: ${{ matrix.arch }}
          CGO_ENABLED: 0
        run: |
          go build -v -ldflags="-s -w" -o kawai-node${{ matrix.os == 'windows' && '.exe' || '' }} ./cmd/server
      
      - name: Create archive
        run: |
          if [ "${{ matrix.os }}" = "windows" ]; then
            7z a kawai-node-${{ needs.prepare.outputs.short_sha }}-${{ matrix.os }}-${{ matrix.arch }}.zip kawai-node.exe
          else
            tar czf kawai-node-${{ needs.prepare.outputs.short_sha }}-${{ matrix.os }}-${{ matrix.arch }}.tar.gz kawai-node
          fi
        shell: bash
      
      - name: Generate checksum
        run: |
          if [ "${{ matrix.os }}" = "windows" ]; then
            sha256sum kawai-node-${{ needs.prepare.outputs.short_sha }}-${{ matrix.os }}-${{ matrix.arch }}.zip > checksums-${{ matrix.os }}-${{ matrix.arch }}.txt
          else
            sha256sum kawai-node-${{ needs.prepare.outputs.short_sha }}-${{ matrix.os }}-${{ matrix.arch }}.tar.gz > checksums-${{ matrix.os }}-${{ matrix.arch }}.txt
          fi
        shell: bash
      
      - name: Upload to R2
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.R2_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.R2_SECRET_ACCESS_KEY }}
          AWS_ENDPOINT_URL: ${{ secrets.R2_ENDPOINT_URL }}
        run: |
          aws s3 cp kawai-node-${{ needs.prepare.outputs.short_sha }}-${{ matrix.os }}-${{ matrix.arch }}.* \
            s3://${{ env.R2_BUCKET }}/node/v${{ needs.prepare.outputs.short_sha }}/ \
            --endpoint-url $AWS_ENDPOINT_URL
          
          aws s3 cp checksums-${{ matrix.os }}-${{ matrix.arch }}.txt \
            s3://${{ env.R2_BUCKET }}/node/v${{ needs.prepare.outputs.short_sha }}/ \
            --endpoint-url $AWS_ENDPOINT_URL

  update-installer:
    name: Update Installer Script
    needs: [prepare, build]
    runs-on: ubuntu-latest
    steps:
      - name: Checkout veridium repo
        uses: actions/checkout@v4
      
      - name: Checkout kawai-website repo
        uses: actions/checkout@v4
        with:
          repository: kawai-network/kawai-website
          token: ${{ secrets.WEBSITE_REPO_TOKEN }}
          path: kawai-website
      
      - name: Update install.sh
        run: |
          cd kawai-website/node
          
          # Update DEFAULT_LATEST_VERSION
          sed -i "s/DEFAULT_LATEST_VERSION=.*/DEFAULT_LATEST_VERSION=\"v${{ needs.prepare.outputs.short_sha }}\"/" install.sh
          
          # Verify change
          grep "DEFAULT_LATEST_VERSION" install.sh
      
      - name: Commit and push
        run: |
          cd kawai-website
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add node/install.sh
          git commit -m "chore: update node installer to v${{ needs.prepare.outputs.short_sha }}"
          git push
      
      - name: Create versions.txt
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.R2_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.R2_SECRET_ACCESS_KEY }}
          AWS_ENDPOINT_URL: ${{ secrets.R2_ENDPOINT_URL }}
        run: |
          echo "v${{ needs.prepare.outputs.short_sha }}" > versions.txt
          aws s3 cp versions.txt s3://${{ env.R2_BUCKET }}/node/versions.txt --endpoint-url $AWS_ENDPOINT_URL

  create-release:
    name: Create GitHub Release
    needs: [prepare, build, update-installer]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: node-v${{ needs.prepare.outputs.version }}
          name: Kawai Node v${{ needs.prepare.outputs.version }}
          body: |
            ## 🎉 Kawai Node v${{ needs.prepare.outputs.version }}
            
            **Version:** v${{ needs.prepare.outputs.short_sha }}
            
            ### 🚀 Quick Install
            
            ```bash
            curl -fsSL https://getkawai.com/node | sh
            ```
            
            ### 📦 Downloads
            
            All binaries are available at:
            `${{ env.R2_PUBLIC_URL }}/node/v${{ needs.prepare.outputs.short_sha }}/`
            
            ### 🔐 Verification
            
            Verify with SHA256 checksums in `checksums-*.txt` files.
          draft: false
          prerelease: false
```

---

## 🔄 Automated Release Process

### Release Flow Diagram

```
Developer creates tag
        ↓
    CI Pipeline
        ↓
┌───────────────────────┐
│   Lint & Test         │
│   - Go tests          │
│   - Frontend tests    │
│   - Contract tests    │
└───────────────────────┘
        ↓
┌───────────────────────┐
│   Build Artifacts     │
│   - Desktop (3 OS)    │
│   - Node server       │
│   - Contracts         │
└───────────────────────┘
        ↓
┌───────────────────────┐
│   Upload & Deploy     │
│   - GitHub Releases   │
│   - R2 Storage        │
│   - Update installer  │
└───────────────────────┘
        ↓
┌───────────────────────┐
│   Post-Deploy         │
│   - Verify downloads  │
│   - Update docs       │
│   - Notify team       │
└───────────────────────┘
```

### Versioning Strategy

**Desktop App:**
- Format: `v1.2.3` (Semantic Versioning)
- Major: Breaking changes
- Minor: New features
- Patch: Bug fixes

**Node Server:**
- Format: `node-v1.2.3` or `v{short-sha}`
- Use commit SHA for rapid iteration
- Use semver for stable releases

**Contracts:**
- Format: `contracts-v1.2.3`
- Major: Breaking contract changes
- Minor: New features (backward compatible)
- Patch: Bug fixes

---

## 🔧 Auto-update install.sh Implementation

### Current install.sh Structure

```bash
#!/bin/bash
# Kawai Node Installer

DEFAULT_LATEST_VERSION="v51aec45"  # ← This needs auto-update
R2_BASE_URL="https://storage.getkawai.com/node"

# ... rest of installer logic
```

### Solution: GitHub Actions Cross-Repo Update

**Requirements:**
1. Personal Access Token (PAT) with `repo` scope
2. Access to kawai-website repository
3. Automated commit and push

**Implementation:**

```yaml
# In .github/workflows/release-node.yml

- name: Update installer script
  run: |
    # Clone kawai-website repo
    git clone https://${{ secrets.WEBSITE_REPO_TOKEN }}@github.com/kawai-network/kawai-website.git
    cd kawai-website/node
    
    # Update version
    sed -i "s/DEFAULT_LATEST_VERSION=.*/DEFAULT_LATEST_VERSION=\"v$SHORT_SHA\"/" install.sh
    
    # Commit and push
    git config user.name "github-actions[bot]"
    git config user.email "github-actions[bot]@users.noreply.github.com"
    git add install.sh
    git commit -m "chore: update node installer to v$SHORT_SHA"
    git push
```

### Alternative: Dynamic Version Fetching

**Modify install.sh to fetch version from R2:**

```bash
#!/bin/bash
# Kawai Node Installer

R2_BASE_URL="https://storage.getkawai.com/node"

# Fetch latest version from R2
LATEST_VERSION=$(curl -fsSL "$R2_BASE_URL/versions.txt" | head -n1)

if [ -z "$LATEST_VERSION" ]; then
  echo "❌ Failed to fetch latest version"
  exit 1
fi

echo "📦 Installing Kawai Node $LATEST_VERSION"

# ... rest of installer logic
```

**Pros:**
- No cross-repo access needed
- Always up-to-date
- Simple implementation

**Cons:**
- Extra HTTP request
- Dependency on R2 availability

---

## 📝 Implementation Steps

### Phase 1: Setup (Week 1)

#### Day 1-2: Repository Setup
- [ ] Create `.golangci.yml` for Go linting
- [ ] Create `.github/changelog-config.json` for changelog generation
- [ ] Setup Codecov account and get token
- [ ] Create GitHub secrets:
  - `CODECOV_TOKEN`
  - `WEBSITE_REPO_TOKEN`
  - `R2_ACCESS_KEY_ID`
  - `R2_SECRET_ACCESS_KEY`
  - `R2_ENDPOINT_URL`

#### Day 3-4: CI Workflow
- [ ] Create `.github/workflows/ci.yml`
- [ ] Test lint jobs locally
- [ ] Test test jobs locally
- [ ] Push and verify on GitHub Actions

#### Day 5: Contract Testing
- [ ] Add Foundry to CI
- [ ] Setup Slither for security scanning
- [ ] Test contract workflow

### Phase 2: Release Automation (Week 2)

#### Day 1-2: Desktop Release
- [ ] Enhance `.github/workflows/release-desktop.yml`
- [ ] Test on all platforms
- [ ] Setup code signing (macOS, Windows)

#### Day 3-4: Node Release
- [ ] Enhance `.github/workflows/release-node.yml`
- [ ] Implement cross-repo update
- [ ] Test installer update

#### Day 5: Integration
- [ ] Test full release flow
- [ ] Document release process
- [ ] Create runbook

### Phase 3: Monitoring & Optimization (Week 3)

#### Day 1-2: Monitoring
- [ ] Setup GitHub Actions monitoring
- [ ] Create Slack/Discord notifications
- [ ] Setup failure alerts

#### Day 3-4: Optimization
- [ ] Optimize build times (caching)
- [ ] Parallelize jobs
- [ ] Reduce artifact sizes

#### Day 5: Documentation
- [ ] Update README with CI badges
- [ ] Create CONTRIBUTING.md
- [ ] Document release process

---

## 🧪 Testing Strategy

### Local Testing

```bash
# Test Go linting
golangci-lint run --timeout=10m

# Test Go tests
go test -race -coverprofile=coverage.out ./...

# Test contract compilation
cd contracts && forge test

# Test frontend build
cd frontend && npm run build
```

### CI Testing

```bash
# Use act to test GitHub Actions locally
act -j lint-go
act -j test-go
act -j build-backend
```

### Release Testing

```bash
# Test release workflow (dry-run)
git tag -a v0.0.1-test -m "Test release"
git push origin v0.0.1-test

# Verify:
# 1. CI passes
# 2. Artifacts uploaded
# 3. Release created
# 4. Installer updated

# Cleanup
git tag -d v0.0.1-test
git push origin :refs/tags/v0.0.1-test
```

---

## 🚀 Rollout Plan

### Week 1: Soft Launch
- Deploy CI workflow to `enhancement` branch
- Test with feature branches
- Gather feedback from team

### Week 2: Beta Testing
- Merge to `master`
- Test with real releases (pre-release tags)
- Monitor for issues

### Week 3: Full Rollout
- Enable for all releases
- Update documentation
- Train team on new process

### Week 4: Optimization
- Analyze metrics
- Optimize slow jobs
- Implement improvements

---

## 📊 Success Criteria

### Metrics to Track

| Metric | Target | How to Measure |
|--------|--------|----------------|
| CI Success Rate | >95% | GitHub Actions dashboard |
| Average CI Time | <15 min | GitHub Actions insights |
| Release Time | <30 min | Manual tracking |
| Failed Releases | <5% | Release history |
| Test Coverage | >70% | Codecov dashboard |
| Security Issues | 0 critical | Trivy/Slither reports |

### Monitoring Dashboard

Create a dashboard to track:
- CI/CD pipeline health
- Build times trend
- Test coverage trend
- Release frequency
- Deployment success rate

---

## 🔐 Security Considerations

### Secrets Management
- Use GitHub Secrets for sensitive data
- Rotate tokens quarterly
- Use least-privilege access
- Audit secret usage

### Code Signing
- macOS: Apple Developer ID
- Windows: Code signing certificate
- Linux: GPG signing

### Supply Chain Security
- Pin action versions
- Use verified actions only
- Scan dependencies (Dependabot)
- SBOM generation

---

## 📚 Additional Resources

### GitHub Actions
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)
- [Best Practices](https://docs.github.com/en/actions/guides/security-hardening-for-github-actions)

### Tools
- [golangci-lint](https://golangci-lint.run/)
- [Foundry](https://book.getfoundry.sh/)
- [Codecov](https://docs.codecov.com/)
- [act](https://github.com/nektos/act) - Local GitHub Actions testing

---

## 🎯 Next Steps

1. **Review this plan** with the team
2. **Create GitHub issues** for each phase
3. **Assign owners** for each task
4. **Start with Phase 1** (Setup)
5. **Iterate and improve** based on feedback

---

**Prepared by:** AI Assistant  
**For:** Veridium/Kawai Network Team  
**Date:** 2 Februari 2026  
**Status:** Ready for Implementation
