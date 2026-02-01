#!/bin/sh
# Kawai Node Installer Script
# Usage: curl -fsSL getkawai.com/node | sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
R2_BASE_URL="https://storage.getkawai.com"
BUCKET="kawai"
INSTALL_DIR="/usr/local/bin"
VERSION="${VERSION:-latest}"

# Print functions
print_info() {
    echo "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo "${YELLOW}[WARNING]${NC} $1"
}

# Detect OS
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    echo "$OS"
}

# Detect architecture
detect_arch() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    echo "$ARCH"
}

# Get latest version from R2
get_latest_version() {
    print_info "Checking for latest version..."
    
    # List objects and get the latest version directory
    LATEST=$(curl -fsSL "${R2_BASE_URL}/node/?list-type=2&prefix=node/v&delimiter=/" 2>/dev/null | \
        grep -o 'node/v[^/]*/' | \
        sed 's|node/v||;s|/||' | \
        sort -V | \
        tail -1)
    
    if [ -z "$LATEST" ]; then
        # Fallback to a known version if we can't detect
        LATEST="69e4e5e"
    fi
    
    echo "$LATEST"
}

# Download and install
download_and_install() {
    OS=$1
    ARCH=$2
    VERSION=$3
    
    FILENAME="kawai-node-${VERSION}-${OS}-${ARCH}.tar.gz"
    URL="${R2_BASE_URL}/node/v${VERSION}/${FILENAME}"
    
    print_info "Downloading kawai-node ${VERSION} for ${OS}/${ARCH}..."
    print_info "URL: ${URL}"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT
    
    # Download
    if ! curl -fsSL "$URL" -o "${TMP_DIR}/${FILENAME}"; then
        print_error "Failed to download ${FILENAME}"
        print_info "Trying alternative URL..."
        
        # Try without the v prefix
        URL="${R2_BASE_URL}/node/${VERSION}/${FILENAME}"
        if ! curl -fsSL "$URL" -o "${TMP_DIR}/${FILENAME}"; then
            print_error "Failed to download from alternative URL"
            exit 1
        fi
    fi
    
    print_success "Download complete"
    
    # Extract
    print_info "Extracting..."
    tar -xzf "${TMP_DIR}/${FILENAME}" -C "$TMP_DIR"
    
    # Find the binary
    if [ -f "${TMP_DIR}/kawai-node" ]; then
        BINARY="${TMP_DIR}/kawai-node"
    else
        print_error "Could not find kawai-node binary in archive"
        exit 1
    fi
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        SUDO=""
    else
        print_warning "Installation directory requires elevated permissions"
        SUDO="sudo"
    fi
    
    # Install
    print_info "Installing to ${INSTALL_DIR}/kawai-node..."
    $SUDO mv "$BINARY" "${INSTALL_DIR}/kawai-node"
    $SUDO chmod +x "${INSTALL_DIR}/kawai-node"
    
    print_success "Installation complete!"
}

# Verify installation
verify_installation() {
    if command -v kawai-node >/dev/null 2>&1; then
        print_success "kawai-node is now available in PATH"
        print_info "Version: $(kawai-node --version 2>/dev/null || echo 'unknown')"
        print_info "Run 'kawai-node --help' to get started"
    else
        print_warning "kawai-node installed but not in PATH"
        print_info "Add ${INSTALL_DIR} to your PATH or run:"
        print_info "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

# Main
main() {
    echo "${GREEN}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║                                                            ║"
    echo "║   🌸 Kawai Node Installer                                  ║"
    echo "║                                                            ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo "${NC}"
    
    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)
    
    print_info "Detected platform: ${OS}/${ARCH}"
    
    # Get version
    if [ "$VERSION" = "latest" ]; then
        VERSION=$(get_latest_version)
    fi
    
    print_info "Installing version: ${VERSION}"
    
    # Download and install
    download_and_install "$OS" "$ARCH" "$VERSION"
    
    # Verify
    verify_installation
    
    echo ""
    echo "${GREEN}✨ Installation complete! Welcome to Kawai Network.${NC}"
    echo ""
    echo "To start contributing:"
    echo "  ${BLUE}kawai-node start${NC}"
    echo ""
    echo "For help:"
    echo "  ${BLUE}kawai-node --help${NC}"
    echo ""
}

# Run main
main
