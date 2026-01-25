#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Load .env file if exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | grep -E '^R2_' | xargs)
fi

# Get version from argument
VERSION=$1
if [ -z "$VERSION" ]; then
  echo -e "${RED}❌ Error: Version required${NC}"
  echo "Usage: $0 <version>"
  echo "Example: $0 0.1.0"
  exit 1
fi

# Remove 'v' prefix if present
VERSION=${VERSION#v}

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Regenerate Checksums v${VERSION}${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if R2 credentials are set
if [ -z "$R2_ACCESS_KEY_ID" ] || [ -z "$R2_SECRET_ACCESS_KEY" ] || [ -z "$R2_ENDPOINT_URL" ]; then
  echo -e "${RED}❌ Error: R2 credentials not set${NC}"
  exit 1
fi

# Configure AWS CLI
export AWS_ACCESS_KEY_ID=$R2_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=$R2_SECRET_ACCESS_KEY
export AWS_DEFAULT_REGION=auto
unset AWS_SESSION_TOKEN

# Create temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo -e "${GREEN}📥 Downloading binaries from R2...${NC}"

# Download binaries
DOWNLOADED=0

if aws s3 cp "s3://kawai/v${VERSION}/Kawai-${VERSION}-macos-universal.tar.gz" . --endpoint-url "$R2_ENDPOINT_URL" 2>/dev/null; then
  echo "✅ Downloaded macOS binary"
  DOWNLOADED=$((DOWNLOADED + 1))
fi

if aws s3 cp "s3://kawai/v${VERSION}/Kawai-${VERSION}-linux-amd64.tar.gz" . --endpoint-url "$R2_ENDPOINT_URL" 2>/dev/null; then
  echo "✅ Downloaded Linux binary"
  DOWNLOADED=$((DOWNLOADED + 1))
fi

if aws s3 cp "s3://kawai/v${VERSION}/Kawai-${VERSION}-windows-amd64.zip" . --endpoint-url "$R2_ENDPOINT_URL" 2>/dev/null; then
  echo "✅ Downloaded Windows binary"
  DOWNLOADED=$((DOWNLOADED + 1))
fi

if [ $DOWNLOADED -eq 0 ]; then
  echo -e "${RED}❌ No binaries found in R2${NC}"
  cd -
  rm -rf "$TEMP_DIR"
  exit 1
fi

echo ""
echo -e "${GREEN}🔐 Generating checksums...${NC}"

# Generate checksums
shasum -a 256 Kawai-*.tar.gz Kawai-*.zip 2>/dev/null > checksums.txt || true

if [ -s checksums.txt ]; then
  echo -e "${GREEN}✅ Generated checksums:${NC}"
  cat checksums.txt
else
  echo -e "${RED}❌ Failed to generate checksums${NC}"
  cd -
  rm -rf "$TEMP_DIR"
  exit 1
fi

# Upload to R2
echo ""
echo -e "${GREEN}☁️  Uploading checksums to R2...${NC}"

aws s3 cp checksums.txt "s3://kawai/v${VERSION}/checksums.txt" \
  --endpoint-url "$R2_ENDPOINT_URL"

echo -e "${GREEN}✅ Checksums uploaded!${NC}"

# Cleanup
cd -
rm -rf "$TEMP_DIR"

echo ""
echo -e "${GREEN}🎉 Done!${NC}"
echo -e "   URL: https://storage.getkawai.com/v${VERSION}/checksums.txt"
