#!/bin/bash
set -e

# Ephemeral Public Repo Release Script
# Creates temporary public repo, builds via GitHub Actions, then deletes repo

VERSION=${1:-"0.1.0"}
TEMP_DIR="kawai-${VERSION}"
PUBLIC_REPO="kawai-network/kawai-build-${VERSION}"

# Load R2 credentials from .env
if [ -f .env ]; then
  export $(grep -v '^#' .env | grep -E '^R2_' | xargs)
fi

# Verify R2 credentials
if [ -z "$R2_ACCESS_KEY_ID" ] || [ -z "$R2_SECRET_ACCESS_KEY" ] || [ -z "$R2_ENDPOINT_URL" ]; then
  echo "❌ Error: R2 credentials not found in .env"
  echo "   Required: R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_ENDPOINT_URL"
  exit 1
fi

echo "🚀 Starting ephemeral release for v${VERSION}"

# ============================================================================
# Step 1: Create clean copy
# ============================================================================
echo "📦 Creating clean copy..."
rm -rf "../${TEMP_DIR}"
cp -r . "../${TEMP_DIR}"
cd "../${TEMP_DIR}"

# ============================================================================
# Step 2: Remove private files
# ============================================================================
echo "🧹 Removing private files..."
rm -rf \
  cmd/ \
  .env* \
  contracts/.env* \
  data/ \
  .brv/ \
  .kiro/ \
  backend-dev.log \
  contract-*.log \
  verify.json \
  .git

# Remove CircleCI (not needed for public)
rm -rf .circleci/

# ============================================================================
# Step 3: Update version
# ============================================================================
echo "📝 Updating version to ${VERSION}..."
# Portable sed (works on both macOS and Linux)
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' "s/version: \".*\"/version: \"${VERSION}\"/" build/config.yml
else
  sed -i "s/version: \".*\"/version: \"${VERSION}\"/" build/config.yml
fi

# ============================================================================
# Step 4: Initialize git
# ============================================================================
echo "🔧 Initializing git..."
git init
git add .
git commit -m "Release v${VERSION}"

# ============================================================================
# Step 5: Create public repo
# ============================================================================
echo "🌐 Creating public repo: ${PUBLIC_REPO}..."
gh repo create "${PUBLIC_REPO}" \
  --public \
  --source=. \
  --remote=origin \
  --description="Kawai v${VERSION} - Temporary build repo"

# ============================================================================
# Step 5.5: Setup GitHub Secrets
# ============================================================================
echo "🔐 Setting up GitHub Secrets..."

# Extract R2_ACCOUNT_ID from endpoint URL
R2_ACCOUNT_ID=$(echo "$R2_ENDPOINT_URL" | sed 's|https://||' | cut -d. -f1)

# Set secrets using gh CLI
gh secret set R2_ACCESS_KEY_ID --body "$R2_ACCESS_KEY_ID" --repo "${PUBLIC_REPO}"
gh secret set R2_SECRET_ACCESS_KEY --body "$R2_SECRET_ACCESS_KEY" --repo "${PUBLIC_REPO}"
gh secret set R2_ACCOUNT_ID --body "$R2_ACCOUNT_ID" --repo "${PUBLIC_REPO}"

echo "✅ Secrets configured"

# ============================================================================
# Step 6: Push and trigger build
# ============================================================================
echo "⬆️  Pushing code..."
git push -u origin master

echo "🏷️  Creating release tag..."
git tag "v${VERSION}"
git push origin "v${VERSION}"

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
cd ../veridium

# Download from GitHub Release
gh release download "v${VERSION}" \
  --repo "${PUBLIC_REPO}" \
  --dir "releases/v${VERSION}"

echo "✅ Artifacts downloaded to: releases/v${VERSION}/"
ls -lh "releases/v${VERSION}/"

# ============================================================================
# Step 9: Upload to R2
# ============================================================================
echo "☁️  Uploading to Cloudflare R2..."

# Configure AWS CLI for R2
export AWS_ACCESS_KEY_ID="$R2_ACCESS_KEY_ID"
export AWS_SECRET_ACCESS_KEY="$R2_SECRET_ACCESS_KEY"
export AWS_ENDPOINT_URL="$R2_ENDPOINT_URL"

aws s3 sync "releases/v${VERSION}/" "s3://kawai/v${VERSION}/" \
  --endpoint-url "$R2_ENDPOINT_URL" \
  --exclude "*" \
  --include "*.tar.gz" \
  --include "*.zip" \
  --include "*.deb" \
  --include "checksums.txt"

echo "✅ Uploaded to: https://storage.getkawai.com/v${VERSION}/"

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
echo "🎉 Release v${VERSION} completed!"
echo ""
echo "📦 Artifacts:"
echo "   - GitHub Release: Deleted (ephemeral)"
echo "   - Cloudflare R2: https://storage.getkawai.com/v${VERSION}/"
echo "   - Local: releases/v${VERSION}/"
echo ""
echo "🔗 Next steps:"
echo "   1. Test downloads from R2"
echo "   2. Update homebrew formula"
echo "   3. Announce release"
