#!/bin/bash
set -e

# Ephemeral Public Repo Release Script for Kawai Node (cmd/server)
# Creates temporary public repo, builds cmd/server via GitHub Actions, then deletes repo
#
# Usage: ./scripts/ephemeral-release-node.sh [VERSION]
# Example: ./scripts/ephemeral-release-node.sh 1.0.0

VERSION=${1:-"0.1.0"}
TEMP_DIR="kawai-contributor-${VERSION}"
PUBLIC_REPO="kawai-network/kawai-contributor-build-${VERSION}"

# Save original directory name for later use
ORIGINAL_DIR=$(basename "$(pwd)")

# ============================================================================
# Step 0: Load credentials from .env
# ============================================================================
echo "🔑 Loading credentials from .env..."

if [ -f .env ]; then
  # Load GH_TOKEN
  export $(grep -v '^#' .env | grep -E '^GH_TOKEN=' | xargs)
  # Load R2 credentials
  export $(grep -v '^#' .env | grep -E '^R2_' | xargs)
fi

# Verify GH_TOKEN
if [ -z "$GH_TOKEN" ]; then
  echo "❌ Error: GH_TOKEN not found in .env"
  echo "   Add to .env: GH_TOKEN=your_github_token"
  echo "   Get token from: https://github.com/settings/tokens"
  exit 1
fi

# Verify R2 credentials
if [ -z "$R2_ACCESS_KEY_ID" ] || [ -z "$R2_SECRET_ACCESS_KEY" ] || [ -z "$R2_ENDPOINT_URL" ]; then
  echo "❌ Error: R2 credentials not found in .env"
  echo "   Required: R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_ENDPOINT_URL"
  exit 1
fi

echo "✅ Credentials loaded"
echo "🚀 Starting ephemeral release for Kawai Node v${VERSION}"

# ============================================================================
# Step 1: Create clean copy
# ============================================================================
echo "📦 Creating clean copy..."
rm -rf "../${TEMP_DIR}"
cp -r . "../${TEMP_DIR}"
cd "../${TEMP_DIR}"

# ============================================================================
# Step 2: Remove private files (keep only cmd/server and necessary files)
# ============================================================================
echo "🧹 Removing private files..."

# Keep only cmd/server, remove other cmd subdirectories
for dir in cmd/*/; do
  if [ "$dir" != "cmd/server/" ]; then
    rm -rf "$dir"
  fi
done

# Remove other private files
rm -rf \
  .env* \
  contracts/.env* \
  data/ \
  .brv/ \
  .kiro/ \
  backend-dev.log \
  contract-*.log \
  verify.json \
  .git

# Remove folders with files that have invalid characters on Windows
# (e.g., colons in filenames like "Alibaba Qwen 257 Model: A Deep Dive...")
rm -rf \
  kawai-website/ \
  docs/ \
  docs-users/ \
  logo_analysis/ \
  whitepaper/ \
  zz*/ \
  *.md \
  *.png \
  *.ico \
  code.html

# Remove CircleCI (not needed for public)
rm -rf .circleci/

# Remove unnecessary workflows (keep only release-node.yml)
# First check if source workflow exists in parent directory
if [ -f "../veridium/.github/workflows/release-node.yml" ]; then
  echo "📋 Copying release-node.yml from parent directory..."
  mkdir -p .github/workflows
  cp ../veridium/.github/workflows/release-node.yml .github/workflows/
fi

# ============================================================================
# Step 3: Update version in main.go
# ============================================================================
echo "📝 Updating version to ${VERSION}..."

# Update version in cmd/server/main.go or kronk.go
if [ -f "cmd/server/api/services/kronk/kronk.go" ]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/var tag = \".*\"/var tag = \"v${VERSION}\"/" cmd/server/api/services/kronk/kronk.go
  else
    sed -i "s/var tag = \".*\"/var tag = \"v${VERSION}\"/" cmd/server/api/services/kronk/kronk.go
  fi
fi

# ============================================================================
# Step 4: Initialize git
# ============================================================================
echo "🔧 Initializing git..."
git init
git add .
git commit -m "Release Kawai Node v${VERSION}"

# ============================================================================
# Step 5: Create public repo
# ============================================================================
echo "🌐 Creating public repo: ${PUBLIC_REPO}..."
gh repo create "${PUBLIC_REPO}" \
  --public \
  --source=. \
  --remote=origin \
  --description="Kawai Node v${VERSION} - Temporary build repo"

# ============================================================================
# Step 5.5: Setup GitHub Secrets
# ============================================================================
echo "🔐 Setting up GitHub Secrets..."

# Extract R2_ACCOUNT_ID from endpoint URL
R2_ACCOUNT_ID=$(echo "$R2_ENDPOINT_URL" | sed 's|https://||' | cut -d. -f1)

# Set secrets using gh CLI (via stdin to avoid exposing in process list)
echo "$GH_TOKEN" | gh secret set GH_TOKEN --repo "${PUBLIC_REPO}"
echo "$R2_ACCESS_KEY_ID" | gh secret set R2_ACCESS_KEY_ID --repo "${PUBLIC_REPO}"
echo "$R2_SECRET_ACCESS_KEY" | gh secret set R2_SECRET_ACCESS_KEY --repo "${PUBLIC_REPO}"
echo "$R2_ACCOUNT_ID" | gh secret set R2_ACCOUNT_ID --repo "${PUBLIC_REPO}"

echo "✅ Secrets configured"

# ============================================================================
# Step 6: Push and trigger build
# ============================================================================
echo "⬆️  Pushing code..."
git push -u origin master

echo "🏷️  Creating release tag..."
git tag "node-v${VERSION}"
git push origin "node-v${VERSION}"

# ============================================================================
# Step 7: Wait for build
# ============================================================================
echo "⏳ Waiting for GitHub Actions to complete..."
echo "   Monitor: https://github.com/${PUBLIC_REPO}/actions"

# Wait for workflow to start (max 2 minutes)
echo "   Waiting for workflow to start..."
sleep 30

# Poll for completion (max 30 minutes)
MAX_WAIT=1800  # 30 minutes
ELAPSED=0
INTERVAL=30

while [ $ELAPSED -lt $MAX_WAIT ]; do
  # Check workflow status
  STATUS=$(gh run list --repo "${PUBLIC_REPO}" --limit 1 --json status --jq '.[0].status')
  
  if [ "$STATUS" = "completed" ]; then
    CONCLUSION=$(gh run list --repo "${PUBLIC_REPO}" --limit 1 --json conclusion --jq '.[0].conclusion')
    
    if [ "$CONCLUSION" = "success" ]; then
      echo "✅ Build completed successfully!"
      break
    else
      echo "❌ Build failed with conclusion: ${CONCLUSION}"
      echo "   Check logs: https://github.com/${PUBLIC_REPO}/actions"
      exit 1
    fi
  fi
  
  echo "   Status: ${STATUS} (${ELAPSED}s elapsed)"
  sleep $INTERVAL
  ELAPSED=$((ELAPSED + INTERVAL))
done

if [ $ELAPSED -ge $MAX_WAIT ]; then
  echo "⏰ Timeout waiting for build (30 minutes)"
  echo "   Check manually: https://github.com/${PUBLIC_REPO}/actions"
  exit 1
fi

# ============================================================================
# Step 8: Download artifacts
# ============================================================================
echo "📥 Downloading release artifacts..."
cd "../${ORIGINAL_DIR}"

# Download from GitHub Release
gh release download "node-v${VERSION}" \
  --repo "${PUBLIC_REPO}" \
  --dir "releases/node-v${VERSION}"

echo "✅ Artifacts downloaded to: releases/node-v${VERSION}/"
ls -lh "releases/node-v${VERSION}/"

# ============================================================================
# Step 9: Upload to R2
# ============================================================================
echo "☁️  Uploading to Cloudflare R2..."

# Configure AWS CLI for R2
export AWS_ACCESS_KEY_ID="$R2_ACCESS_KEY_ID"
export AWS_SECRET_ACCESS_KEY="$R2_SECRET_ACCESS_KEY"
export AWS_ENDPOINT_URL="$R2_ENDPOINT_URL"

aws s3 sync "releases/node-v${VERSION}/" "s3://kawai/node/v${VERSION}/" \
  --endpoint-url "$R2_ENDPOINT_URL" \
  --exclude "*" \
  --include "*.tar.gz" \
  --include "*.zip" \
  --include "checksums.txt"

# Update 'latest' symlink
echo "🔗 Updating latest symlink..."
aws s3 sync "s3://kawai/node/v${VERSION}/" "s3://kawai/node/latest/" \
  --endpoint-url "$R2_ENDPOINT_URL"

echo "✅ Uploaded to: https://storage.getkawai.com/node/v${VERSION}/"

# ============================================================================
# Step 10: Cleanup - Delete public repo
# ============================================================================
echo "🗑️  Deleting ephemeral repo..."
gh repo delete "${PUBLIC_REPO}" --yes

echo "🧹 Cleaning up temp directory..."
cd ..
rm -rf "${TEMP_DIR}"

# ============================================================================
# Done!
# ============================================================================
echo ""
echo "🎉 Release Kawai Node v${VERSION} completed!"
echo ""
echo "📦 Artifacts:"
echo "   - GitHub Release: Deleted (ephemeral)"
echo "   - Cloudflare R2: https://storage.getkawai.com/node/v${VERSION}/"
echo "   - Local: releases/node-v${VERSION}/"
echo ""
echo "🔗 Install command:"
echo "   curl -fsSL https://getkawai.com/node | sh"
echo ""
echo "📚 Next steps:"
echo "   1. Test downloads from R2"
echo "   2. Update install script if needed"
echo "   3. Announce release"
