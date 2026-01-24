#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load .env file if exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | grep -E '^R2_' | xargs)
fi

# Get version from argument or git tag
VERSION=${1:-$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')}
if [ -z "$VERSION" ]; then
  echo -e "${RED}❌ Error: No version specified and no git tag found${NC}"
  echo "Usage: $0 [version]"
  echo "Example: $0 0.1.0"
  exit 1
fi

echo -e "${GREEN}🚀 Building Kawai macOS v${VERSION}${NC}"

# Check if R2 credentials are set
if [ -z "$R2_ACCESS_KEY_ID" ] || [ -z "$R2_SECRET_ACCESS_KEY" ] || [ -z "$R2_ENDPOINT_URL" ]; then
  echo -e "${YELLOW}⚠️  Warning: R2 credentials not set. Upload will be skipped.${NC}"
  echo "Add these to your .env file:"
  echo "  R2_ACCESS_KEY_ID=your-key"
  echo "  R2_SECRET_ACCESS_KEY=your-secret"
  echo "  R2_ENDPOINT_URL=your-endpoint"
  SKIP_UPLOAD=1
fi

# Update version in config
echo -e "${GREEN}📝 Updating version in build/config.yml${NC}"
sed -i '' "s/version: \".*\"/version: \"$VERSION\"/" build/config.yml

# Build
echo -e "${GREEN}🔨 Building macOS Universal Binary${NC}"
PRODUCTION=true make release-darwin

# Create distribution archive
echo -e "${GREEN}📦 Creating distribution archive${NC}"
cd bin
tar -czf "Kawai-${VERSION}-macos-universal.tar.gz" Kawai.app
shasum -a 256 "Kawai-${VERSION}-macos-universal.tar.gz" > checksums-macos.txt

echo -e "${GREEN}✅ Build complete!${NC}"
echo -e "   Archive: bin/Kawai-${VERSION}-macos-universal.tar.gz"
echo -e "   Checksum: bin/checksums-macos.txt"

# Upload to R2
if [ -z "$SKIP_UPLOAD" ]; then
  echo -e "${GREEN}☁️  Uploading to Cloudflare R2${NC}"
  
  # Check if aws cli is installed
  if ! command -v aws &> /dev/null; then
    echo -e "${YELLOW}⚠️  AWS CLI not found. Installing...${NC}"
    brew install awscli
  fi
  
  # Configure AWS CLI
  export AWS_ACCESS_KEY_ID=$R2_ACCESS_KEY_ID
  export AWS_SECRET_ACCESS_KEY=$R2_SECRET_ACCESS_KEY
  export AWS_DEFAULT_REGION=auto
  
  # Upload binary
  aws s3 cp "Kawai-${VERSION}-macos-universal.tar.gz" \
    "s3://kawai/v${VERSION}/" \
    --endpoint-url "$R2_ENDPOINT_URL"
  
  # Upload checksum
  aws s3 cp checksums-macos.txt \
    "s3://kawai/v${VERSION}/checksums-macos.txt" \
    --endpoint-url "$R2_ENDPOINT_URL"
  
  echo -e "${GREEN}✅ Upload complete!${NC}"
  echo -e "   URL: https://storage.getkawai.com/v${VERSION}/Kawai-${VERSION}-macos-universal.tar.gz"
else
  echo -e "${YELLOW}⚠️  Skipping upload (R2 credentials not set)${NC}"
  echo -e "${YELLOW}📋 Manual upload commands:${NC}"
  echo -e "   aws s3 cp Kawai-${VERSION}-macos-universal.tar.gz s3://kawai/v${VERSION}/ --endpoint-url \$R2_ENDPOINT_URL"
  echo -e "   aws s3 cp checksums-macos.txt s3://kawai/v${VERSION}/checksums-macos.txt --endpoint-url \$R2_ENDPOINT_URL"
fi

cd ..
echo -e "${GREEN}🎉 Done!${NC}"
