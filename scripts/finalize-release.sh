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
echo -e "${BLUE}║   Finalize Release v${VERSION}${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if R2 credentials are set
if [ -z "$R2_ACCESS_KEY_ID" ] || [ -z "$R2_SECRET_ACCESS_KEY" ] || [ -z "$R2_ENDPOINT_URL" ]; then
  echo -e "${RED}❌ Error: R2 credentials not set${NC}"
  echo "Add these to your .env file:"
  echo "  R2_ACCESS_KEY_ID=your-key"
  echo "  R2_SECRET_ACCESS_KEY=your-secret"
  echo "  R2_ENDPOINT_URL=your-endpoint"
  exit 1
fi

# Check if aws cli is installed
if ! command -v aws &> /dev/null; then
  echo -e "${YELLOW}⚠️  AWS CLI not found. Installing...${NC}"
  brew install awscli
fi

# Configure AWS CLI
export AWS_ACCESS_KEY_ID=$R2_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=$R2_SECRET_ACCESS_KEY
export AWS_DEFAULT_REGION=auto
unset AWS_SESSION_TOKEN

# Create temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo -e "${GREEN}📥 Downloading existing checksums from R2...${NC}"

# Download existing combined checksum if exists
aws s3 cp "s3://kawai/v${VERSION}/checksums.txt" existing-checksums.txt \
  --endpoint-url "$R2_ENDPOINT_URL" 2>/dev/null || touch existing-checksums.txt

# Download individual checksums
aws s3 cp "s3://kawai/v${VERSION}/checksums-macos.txt" . \
  --endpoint-url "$R2_ENDPOINT_URL" 2>/dev/null || echo "⚠️  No macOS checksum"

aws s3 cp "s3://kawai/v${VERSION}/checksums-linux.txt" . \
  --endpoint-url "$R2_ENDPOINT_URL" 2>/dev/null || echo "⚠️  No Linux checksum"

aws s3 cp "s3://kawai/v${VERSION}/checksums-windows.txt" . \
  --endpoint-url "$R2_ENDPOINT_URL" 2>/dev/null || echo "⚠️  No Windows checksum"

# Combine checksums (preserve existing + add new)
echo -e "${GREEN}🔗 Combining checksums...${NC}"

# Start with existing checksums
cat existing-checksums.txt > checksums.txt 2>/dev/null || true

# Add new checksums (avoid duplicates)
for file in checksums-*.txt; do
  if [ -f "$file" ]; then
    while IFS= read -r line; do
      # Extract filename from checksum line
      filename=$(echo "$line" | awk '{print $2}')
      # Check if this file is already in checksums.txt
      if ! grep -q "$filename" checksums.txt 2>/dev/null; then
        echo "$line" >> checksums.txt
      fi
    done < "$file"
  fi
done

if [ -s checksums.txt ]; then
  echo -e "${GREEN}✅ Combined checksums:${NC}"
  cat checksums.txt
else
  echo -e "${YELLOW}⚠️  No checksums found${NC}"
fi

# Generate update.json
echo -e "${GREEN}📝 Generating update manifest...${NC}"
cat > update.json <<EOF
{
  "version": "${VERSION}",
  "release_date": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "platforms": {
    "macos-universal": {
      "url": "https://storage.getkawai.com/v${VERSION}/Kawai-${VERSION}-macos-universal.tar.gz"
    },
    "linux-amd64": {
      "url": "https://storage.getkawai.com/v${VERSION}/Kawai-${VERSION}-linux-amd64.tar.gz"
    },
    "windows-amd64": {
      "url": "https://storage.getkawai.com/v${VERSION}/Kawai-${VERSION}-windows-amd64.zip"
    }
  }
}
EOF

echo -e "${GREEN}📄 Manifest:${NC}"
cat update.json

# Upload to R2
echo ""
echo -e "${GREEN}☁️  Uploading to R2...${NC}"

# Upload combined checksums
if [ -s checksums.txt ]; then
  echo -e "${YELLOW}   Uploading checksums.txt...${NC}"
  aws s3 cp checksums.txt "s3://kawai/v${VERSION}/checksums.txt" \
    --endpoint-url "$R2_ENDPOINT_URL"
fi

# Upload version-specific manifest
echo -e "${YELLOW}   Uploading v${VERSION}/update.json...${NC}"
aws s3 cp update.json "s3://kawai/v${VERSION}/update.json" \
  --endpoint-url "$R2_ENDPOINT_URL" \
  --content-type application/json

# Upload to latest/
echo -e "${YELLOW}   Uploading latest/update.json...${NC}"
aws s3 cp update.json "s3://kawai/latest/update.json" \
  --endpoint-url "$R2_ENDPOINT_URL" \
  --content-type application/json

# Cleanup
cd -
rm -rf "$TEMP_DIR"

echo ""
echo -e "${GREEN}✅ Release finalized!${NC}"
echo ""
echo -e "${BLUE}📦 Release structure:${NC}"
echo -e "   https://storage.getkawai.com/"
echo -e "   ├── latest/"
echo -e "   │   └── update.json"
echo -e "   └── v${VERSION}/"
echo -e "       ├── update.json"
echo -e "       ├── Kawai-${VERSION}-macos-universal.tar.gz"
echo -e "       ├── Kawai-${VERSION}-linux-amd64.tar.gz"
echo -e "       ├── Kawai-${VERSION}-windows-amd64.zip"
echo -e "       └── checksums.txt"
echo ""
echo -e "${GREEN}🎉 Done!${NC}"
