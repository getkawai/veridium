#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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
echo -e "${BLUE}║   Kawai Release v${VERSION}${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if R2 credentials are set
if [ -z "$R2_ACCESS_KEY_ID" ] || [ -z "$R2_SECRET_ACCESS_KEY" ] || [ -z "$R2_ENDPOINT_URL" ]; then
  echo -e "${RED}❌ Error: R2 credentials not set${NC}"
  echo "Set these environment variables:"
  echo "  - R2_ACCESS_KEY_ID"
  echo "  - R2_SECRET_ACCESS_KEY"
  echo "  - R2_ENDPOINT_URL"
  exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
  echo -e "${YELLOW}⚠️  Warning: You have uncommitted changes${NC}"
  read -p "Continue anyway? (y/N) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
  fi
fi

# Update version in config
echo -e "${GREEN}📝 Updating version in build/config.yml${NC}"
sed -i '' "s/version: \".*\"/version: \"$VERSION\"/" build/config.yml

# Commit version change
echo -e "${GREEN}💾 Committing version change${NC}"
git add build/config.yml
git commit -m "chore: Bump version to $VERSION" || echo "No changes to commit"

# Create and push tag
echo -e "${GREEN}🏷️  Creating tag v${VERSION}${NC}"
git tag "v${VERSION}"

echo -e "${GREEN}⬆️  Pushing tag to GitHub${NC}"
git push origin "v${VERSION}"

echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   CircleCI Build Started               ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo -e "${YELLOW}🔄 Linux and Windows builds running on CircleCI...${NC}"
echo ""

# Build macOS in parallel
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Building macOS (Local)               ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Update version in config (in case it was reverted)
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
echo ""

# Upload to R2
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
echo -e "${YELLOW}   Uploading binary...${NC}"
aws s3 cp "Kawai-${VERSION}-macos-universal.tar.gz" \
  "s3://kawai/v${VERSION}/" \
  --endpoint-url "$R2_ENDPOINT_URL"

# Upload checksum
echo -e "${YELLOW}   Uploading checksum...${NC}"
aws s3 cp checksums-macos.txt \
  "s3://kawai/v${VERSION}/checksums-macos.txt" \
  --endpoint-url "$R2_ENDPOINT_URL"

cd ..

echo ""
echo -e "${GREEN}✅ macOS build uploaded!${NC}"
echo ""

# Summary
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Release Summary                      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✅ Tag pushed:${NC} v${VERSION}"
echo -e "${GREEN}✅ macOS build:${NC} Complete & Uploaded"
echo -e "${YELLOW}🔄 Linux build:${NC} Running on CircleCI"
echo -e "${YELLOW}🔄 Windows build:${NC} Running on CircleCI"
echo ""
echo -e "${BLUE}📦 Downloads will be available at:${NC}"
echo -e "   https://storage.getkawai.com/v${VERSION}/"
echo ""
echo -e "${BLUE}🔗 Monitor CircleCI:${NC}"
echo -e "   https://app.circleci.com/pipelines/circleci/DMNjjttDaMw1NbzvTJEDKw/RUy8Q4CrjMVqvmg31M9y3C"
echo ""
echo -e "${GREEN}🎉 Release process started!${NC}"
