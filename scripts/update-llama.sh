#!/bin/bash

# Script untuk update llama.cpp ke versi terbaru
# Usage: ./scripts/update-llama.sh [options]
#
# Options:
#   --force          Force update even if already on latest version
#   --version        Show currently installed version
#   --list           List available versions from GitHub
#   --processor TYPE Specify processor type: auto, cpu, cuda, vulkan, metal
#   --path PATH      Custom installation path

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}🔧 Veridium llama.cpp Update Script${NC}"
echo "===================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓${NC} Go is installed: $(go version)"
echo ""

# Change to project root
cd "$PROJECT_ROOT"

# Build the update tool
echo -e "${BLUE}📦 Building update tool...${NC}"
go build -o /tmp/update-llama ./cmd/update-llama

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to build update tool${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Build successful"
echo ""

# Run the update tool with all arguments passed to this script
echo -e "${BLUE}🚀 Running update tool...${NC}"
echo ""
/tmp/update-llama "$@"

# Cleanup
rm -f /tmp/update-llama

echo ""
echo -e "${GREEN}✨ Done!${NC}"

