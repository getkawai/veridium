# ==============================================================================
# Veridium Makefile
# ==============================================================================

.PHONY: help dev dev-fast dev-hot dev-rebuild build clean test generate \
        bindings-generate constants-generate \
        admin-register admin-register-dry \
        release-prepare release-version \
        release-darwin release-darwin-package release-darwin-archive \
        release-linux release-linux-deb release-linux-archive \
        release-windows release-windows-archive \
        release-all release-archives release-clean \
        contributor contributor-dev contributor-dev-fresh contributor-build

# ------------------------------------------------------------------------------
# Configuration
# ------------------------------------------------------------------------------
DATA_DIR := data

# ==============================================================================
# Help
# ==============================================================================
help:
	@echo "Veridium Development Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev              Start fresh (reset DB + full build)"
	@echo "  make dev-fast         Start fresh (reset DB, skip frontend build)"
	@echo "  make dev-hot          Hot reload (keep DB, skip build) - fastest"
	@echo "  make dev-rebuild      Rebuild frontend, keep DB"
	@echo "  make build            Build production binary"
	@echo ""
	@echo "Release Builds:"
	@echo "  make release-version         Show current version"
	@echo "  make release-darwin          Build for macOS (Universal Binary)"
	@echo "  make release-darwin-package  Build + package for macOS"
	@echo "  make release-darwin-archive  Build + create distribution archive"
	@echo "  make release-linux           Build for Linux (amd64)"
	@echo "  make release-linux-deb       Build + create .deb package"
	@echo "  make release-linux-archive   Build + create distribution archive"
	@echo "  make release-windows         Build for Windows (amd64)"
	@echo "  make release-windows-archive Build + create distribution archive"
	@echo "  make release-all             Build for all platforms"
	@echo "  make release-archives        Create distribution archives (with checksums)"
	@echo "  make release-clean           Clean release artifacts"
	@echo ""
	@echo "Code Generation:"
	@echo "  make generate         Run bindings-generate + constants-generate"
	@echo "  make bindings-generate Generate TypeScript bindings (wails)"
	@echo "  make constants-generate Generate constants from .env"
	@echo "  make api-docs-generate Generate API documentation (Markdown + TSX)"
	@echo ""
	@echo "Admin Operations:"
	@echo "  make admin-register       Register all treasury addresses as admin"
	@echo "  make admin-register-dry   Preview admin registration (dry-run)"
	@echo ""
	@echo "Contributor Node (cmd/server):"
	@echo "  make contributor            Run contributor node (builds if needed)"
	@echo "  make contributor-build      Build contributor binary"
	@echo "  make contributor-dev        Dev mode with hot reload"
	@echo "  make contributor-dev-fresh  Fresh dev mode (reset data)"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean            Clean generated files"
	@echo "  make test             Run tests"

# ==============================================================================
# Development
# ==============================================================================
dev:
	@echo "🚀 Starting fresh development server..."
	killport 9245 || true
	rm -rf $(DATA_DIR)
	rm -f backend-dev.log
	VERIDIUM_DEV=1 wails3 dev 2>&1 | tee backend-dev.log

dev-fast:
	@echo "⚡ Starting fresh development server (skip frontend build)..."
	killport 9245 || true
	rm -rf $(DATA_DIR)
	rm -f backend-dev.log
	VERIDIUM_DEV=1 wails3 dev -config ./build/config-skip-frontend-build.yml 2>&1 | tee backend-dev.log

dev-hot:
	@echo "🔥 Hot reload (keep DB, skip frontend build)..."
	killport 9245 || true
	rm -f backend-dev.log
	VERIDIUM_DEV=1 wails3 dev -config ./build/config-skip-frontend-build.yml 2>&1 | tee backend-dev.log

dev-rebuild:
	@echo "🔧 Rebuilding frontend (keep DB)..."
	killport 9245 || true
	rm -f backend-dev.log
	VERIDIUM_DEV=1 wails3 dev 2>&1 | tee backend-dev.log

build:
	@echo "🔨 Building production binary..."
	wails3 build

# ==============================================================================
# Release Builds (Platform-Specific)
# ==============================================================================
release-prepare:
	@echo "📋 Preparing for release build..."
	@echo "Step 1: Generating bindings..."
	@wails3 task common:generate:bindings
	@echo "✅ Release preparation complete!"
	@echo "Note: Constants should be pre-generated locally before release"

release-version:
	@echo "📌 Current version: $(shell grep 'version:' build/config.yml | tail -1 | awk '{print $$2}' | tr -d '"')"
	@echo ""
	@echo "To update version, edit build/config.yml and update the 'version' field"
	@echo "Then run: wails3 task common:update:build-assets"

release-darwin:
	@echo "🍎 Building for macOS (Universal Binary)..."
	@$(MAKE) release-prepare
	@echo "Building macOS application..."
	@wails3 task darwin:build:universal
	@echo "Creating macOS app bundle..."
	@wails3 task darwin:create:app:bundle
	@echo "✅ macOS build complete!"
	@echo "📦 Location: build/bin/Kawai.app"

release-darwin-package:
	@echo "📦 Packaging macOS application..."
	@$(MAKE) release-darwin
	@wails3 task darwin:package:universal
	@echo "✅ macOS package complete!"
	@echo "📦 Location: build/bin/"

release-darwin-archive:
	@echo "📦 Creating macOS distribution archive..."
	@$(MAKE) release-darwin
	@cd build/bin && tar -czf Kawai-$(shell grep 'version:' ../config.yml | tail -1 | awk '{print $$2}' | tr -d '"')-macos-universal.tar.gz Kawai.app
	@cd build/bin && shasum -a 256 Kawai-*.tar.gz > checksums.txt
	@echo "✅ macOS archive created with checksum!"
	@echo "📦 Location: build/bin/Kawai-*-macos-universal.tar.gz"

release-linux:
	@echo "🐧 Building for Linux (amd64)..."
	@$(MAKE) release-prepare
	@echo "Building Linux application..."
	@wails3 task linux:build
	@echo "✅ Linux build complete!"
	@echo "📦 Location: build/bin/Kawai"

release-linux-archive:
	@echo "📦 Creating Linux distribution archive..."
	@$(MAKE) release-linux
	@cd build/bin && tar -czf Kawai-$(shell grep 'version:' ../config.yml | tail -1 | awk '{print $$2}' | tr -d '"')-linux-amd64.tar.gz Kawai
	@cd build/bin && shasum -a 256 Kawai-*-linux-*.tar.gz >> checksums.txt
	@echo "✅ Linux archive created with checksum!"
	@echo "📦 Location: build/bin/Kawai-*-linux-amd64.tar.gz"

release-linux-deb:
	@echo "📦 Creating Linux .deb package..."
	@$(MAKE) release-linux
	@wails3 task linux:create:deb
	@echo "✅ Debian package complete!"
	@echo "📦 Location: build/bin/"

release-windows:
	@echo "🪟 Building for Windows (amd64)..."
	@$(MAKE) release-prepare
	@echo "Building Windows application..."
	@wails3 task windows:build
	@echo "✅ Windows build complete!"
	@echo "📦 Location: build/bin/Kawai.exe"

release-windows-archive:
	@echo "📦 Creating Windows distribution archive..."
	@$(MAKE) release-windows
	@cd build/bin && zip -r Kawai-$(shell grep 'version:' ../config.yml | tail -1 | awk '{print $$2}' | tr -d '"')-windows-amd64.zip Kawai.exe
	@cd build/bin && shasum -a 256 Kawai-*-windows-*.zip >> checksums.txt
	@echo "✅ Windows archive created with checksum!"
	@echo "📦 Location: build/bin/Kawai-*-windows-amd64.zip"

release-all:
	@echo "🚀 Building for all platforms..."
	@$(MAKE) release-darwin
	@$(MAKE) release-linux
	@$(MAKE) release-windows
	@echo ""
	@echo "✅ All platform builds complete!"
	@echo ""
	@echo "📦 Build artifacts:"
	@ls -lh build/bin/

release-archives:
	@echo "📦 Creating distribution archives for all platforms..."
	@rm -f build/bin/checksums.txt
	@$(MAKE) release-darwin-archive
	@$(MAKE) release-linux-archive
	@$(MAKE) release-windows-archive
	@echo ""
	@echo "✅ All archives created!"
	@echo ""
	@echo "📦 Distribution files:"
	@ls -lh build/bin/*.tar.gz build/bin/*.zip 2>/dev/null || true
	@echo ""
	@echo "🔐 Checksums:"
	@cat build/bin/checksums.txt

release-clean:
	@echo "🧹 Cleaning release artifacts..."
	@rm -rf build/bin/*
	@rm -rf build/darwin/*.app
	@echo "✅ Release artifacts cleaned!"

# ==============================================================================
# Code Generation
# ==============================================================================
generate: bindings-generate constants-generate
	@echo "✅ All code generated!"

bindings-generate:
	@echo "🔄 Generating TypeScript bindings..."
	wails3 generate bindings -clean=true -ts
	@echo "✅ TypeScript bindings generated!"

constants-generate:
	@echo "🔄 Generating constants from .env..."
	go run cmd/obfuscator-gen/main.go
	@echo "✅ Constants generated!"

api-docs-generate:
	@echo "🔄 Generating API documentation..."
	go run ./cmd/server/api/tooling/docs
	@echo "✅ API documentation generated!"

# ==============================================================================
# Maintenance
# ==============================================================================
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf frontend/bindings/
	rm -rf build/bin
	@echo "✅ Clean complete!"

test:
	@echo "🧪 Running tests..."
	go test ./...

# ==============================================================================
# Admin Operations
# ==============================================================================
admin-register:
	@echo "🔐 Registering admin addresses..."
	@go run cmd/register-admin/main.go
	@echo "✅ Admin registration complete!"

admin-register-dry:
	@echo "🔍 Preview admin registration (dry-run)..."
	@go run cmd/register-admin/main.go --dry-run

# ==============================================================================
# Contributor Node (cmd/server)
# ==============================================================================

CONTRIBUTOR_BIN := bin/kawai-contributor
CONTRIBUTOR_SRC := cmd/server

# Run contributor node (builds if needed)
contributor: contributor-build
	@echo "🚀 Starting contributor node..."
	@./$(CONTRIBUTOR_BIN)

# Build contributor binary
contributor-build:
	@echo "🔨 Building contributor binary..."
	@mkdir -p bin
	@go build -o $(CONTRIBUTOR_BIN) ./$(CONTRIBUTOR_SRC)
	@echo "✅ Built: $(CONTRIBUTOR_BIN)"

# Build optimized contributor binary (smaller size)
contributor-build-optimized:
	@echo "🔨 Building optimized contributor binary..."
	@mkdir -p bin
	@go build -ldflags="-s -w" -trimpath -o $(CONTRIBUTOR_BIN) ./$(CONTRIBUTOR_SRC)
	@echo "✅ Built (optimized): $(CONTRIBUTOR_BIN)"
	@ls -lh $(CONTRIBUTOR_BIN)

# Build and compress with UPX (smallest size, requires UPX installed)
contributor-build-compressed:
	@echo "🔨 Building and compressing contributor binary..."
	@mkdir -p bin
	@go build -ldflags="-s -w" -trimpath -o $(CONTRIBUTOR_BIN) ./$(CONTRIBUTOR_SRC)
	@echo "📦 Original size:"
	@ls -lh $(CONTRIBUTOR_BIN)
	@if command -v upx >/dev/null 2>&1; then \
		echo "🗜️  Compressing with UPX..."; \
		upx --best --lzma $(CONTRIBUTOR_BIN); \
		echo "✅ Compressed size:"; \
		ls -lh $(CONTRIBUTOR_BIN); \
	else \
		echo "⚠️  UPX not installed. Skipping compression."; \
		echo "   Install with: brew install upx (macOS) or apt install upx (Linux)"; \
	fi

# Dev mode with hot reload
contributor-dev:
	@echo "🔥 Starting contributor node (dev mode)..."
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@VERIDIUM_DEV=1 air -c .air.toml 2>&1 | tee contributor-dev.log || VERIDIUM_DEV=1 go run ./$(CONTRIBUTOR_SRC)

# Fresh dev mode (reset data)
contributor-dev-fresh:
	@echo "🚀 Fresh dev mode (resetting data)..."
	@rm -rf data/
	@mkdir -p data
	@$(MAKE) contributor-dev
