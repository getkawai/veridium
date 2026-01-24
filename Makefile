# ==============================================================================
# Veridium Makefile
# ==============================================================================

.PHONY: help dev dev-fast dev-hot dev-rebuild build clean test generate \
        db-generate bindings-generate constants-generate db-dump db-restore \
        contracts-compile contracts-bindings contracts-update \
        contracts-test contracts-test-gas contracts-coverage \
        contracts-deploy-local contracts-deploy-testnet contracts-deploy-vault contracts-verify \
        contracts-upgrade contracts-clean contracts-validate \
        contracts-gas-snapshot contracts-gas-compare \
        admin-register admin-register-dry \
        docs-install docs-serve docs-build docs-clean docs-deploy \
        release-prepare release-version \
        release-darwin release-darwin-package release-darwin-archive \
        release-linux release-linux-deb release-linux-rpm release-linux-appimage release-linux-archive \
        release-windows release-windows-nsis release-windows-msix release-windows-archive \
        release-all release-archives release-packages release-clean

# ------------------------------------------------------------------------------
# Configuration
# ------------------------------------------------------------------------------
CONTRACTS_DIR := contracts
ABIS_DIR      := ./internal/generate/abi
DATA_DIR      := data

# Foundry artifacts
TOKEN_ARTIFACT       := $(CONTRACTS_DIR)/out/KawaiToken.sol/KawaiToken.json
OTCMARKET_ARTIFACT   := $(CONTRACTS_DIR)/out/OTCMarket.sol/OTCMarket.json
VAULT_ARTIFACT       := $(CONTRACTS_DIR)/out/PaymentVault.sol/PaymentVault.json
REVENUE_ARTIFACT     := $(CONTRACTS_DIR)/out/RevenueDistributor.sol/RevenueDistributor.json
REFERRAL_ARTIFACT    := $(CONTRACTS_DIR)/out/ReferralRewardDistributor.sol/ReferralRewardDistributor.json
CASHBACK_ARTIFACT    := $(CONTRACTS_DIR)/out/DepositCashbackDistributor.sol/DepositCashbackDistributor.json
MINING_ARTIFACT      := $(CONTRACTS_DIR)/out/MiningRewardDistributor.sol/MiningRewardDistributor.json
STABLECOIN_ARTIFACT  := $(CONTRACTS_DIR)/out/MockStablecoin.sol/MockStablecoin.json

# Load environment variables from contracts/.env if exists
-include $(CONTRACTS_DIR)/.env
export

# Deployment configuration (can be overridden via env vars or contracts/.env)
PRIVATE_KEY ?=
RPC_URL ?= https://testnet-rpc.monad.xyz
CONTRACT_ADDRESS ?=
ETHERSCAN_API_KEY ?=
CHAIN_ID ?= 10143

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
	@echo "  make release-linux-rpm       Build + create .rpm package"
	@echo "  make release-linux-appimage  Build + create AppImage"
	@echo "  make release-linux-archive   Build + create distribution archive"
	@echo "  make release-windows         Build for Windows (amd64)"
	@echo "  make release-windows-nsis    Build + create NSIS installer"
	@echo "  make release-windows-msix    Build + create MSIX package"
	@echo "  make release-windows-archive Build + create distribution archive"
	@echo "  make release-all             Build for all platforms"
	@echo "  make release-archives        Create distribution archives (with checksums)"
	@echo "  make release-packages        Create packages for all platforms"
	@echo "  make release-clean           Clean release artifacts"
	@echo ""
	@echo "Code Generation:"
	@echo "  make generate         Run db-generate + bindings-generate + constants-generate"
	@echo "  make db-generate      Generate Go code from SQL (sqlc)"
	@echo "  make bindings-generate Generate TypeScript bindings (wails)"
	@echo "  make constants-generate Generate constants from .env"
	@echo ""
	@echo "Database:"
	@echo "  make db-dump          Dump database to seed file"
	@echo "  make db-restore       Restore database from seed file"
	@echo ""
	@echo "Smart Contracts:"
	@echo "  make contracts-update       Compile + generate bindings"
	@echo "  make contracts-upgrade      Full upgrade workflow (test + compile + bindings)"
	@echo "  make contracts-compile      Compile with Foundry"
	@echo "  make contracts-bindings     Generate Go bindings"
	@echo "  make contracts-test         Run contract tests"
	@echo "  make contracts-test-gas     Run tests with gas report"
	@echo "  make contracts-coverage     Generate test coverage report"
	@echo "  make contracts-validate     Validate contracts before deployment"
	@echo "  make contracts-gas-snapshot Create gas usage baseline"
	@echo "  make contracts-gas-compare  Compare gas usage vs baseline"
	@echo "  make contracts-deploy-local Deploy to local Anvil"
	@echo "  make contracts-deploy-testnet Deploy to Monad Testnet"
	@echo "  make contracts-deploy-vault Deploy PaymentVault (mainnet/testnet)"
	@echo "  make contracts-deploy-referral-testnet Deploy ReferralRewardDistributor"
	@echo "  make contracts-grant-minter-referral Grant MINTER_ROLE to referral contract"
	@echo "  make contracts-deploy-cashback-testnet Deploy DepositCashbackDistributor"
	@echo "  make contracts-grant-minter-cashback Grant MINTER_ROLE to cashback contract"
	@echo "  make contracts-deploy-mining-testnet Deploy MiningRewardDistributor"
	@echo "  make contracts-grant-minter-mining Grant MINTER_ROLE to mining contract"
	@echo "  make contracts-verify       Verify contract on explorer"
	@echo "  make contracts-clean        Clean contract artifacts"
	@echo ""
	@echo "Admin Operations:"
	@echo "  make admin-register       Register all treasury addresses as admin"
	@echo "  make admin-register-dry   Preview admin registration (dry-run)"
	@echo ""
	@echo "Emergency Pause Operations:"
	@echo "  make pause-status         Check pause status of all distributors"
	@echo "  make pause-all            Pause all distributors (EMERGENCY)"
	@echo "  make unpause-all          Unpause all distributors"
	@echo "  make pause-mining         Pause mining distributor only"
	@echo "  make unpause-mining       Unpause mining distributor only"
	@echo "  make pause-cashback       Pause cashback distributor only"
	@echo "  make unpause-cashback     Unpause cashback distributor only"
	@echo "  make pause-referral       Pause referral distributor only"
	@echo "  make unpause-referral     Unpause referral distributor only"
	@echo ""
	@echo "Reward Settlement (Unified - All Types):"
	@echo "  make reward-settlement-generate TYPE=<type>  Generate settlement (mining|cashback|referral|revenue)"
	@echo "  make reward-settlement-upload TYPE=<type>    Upload Merkle root to contract"
	@echo "  make reward-settlement-status                Show status for all reward types"
	@echo "  make reward-settlement-all                   Settle all reward types at once"
	@echo ""
	@echo "Shortcuts:"
	@echo "  make settle-mining                           Mining settlement"
	@echo "  make settle-cashback                         Cashback settlement"
	@echo "  make settle-referral                         Referral settlement"
	@echo "  make settle-revenue                          Revenue settlement (USDT dividends)"
	@echo "  make settle-all                              All settlements at once"
	@echo ""
	@echo "Testing Helpers:"
	@echo "  make test-mining-rewards      Run all mining rewards tests"
	@echo "  make test-inject-mining-data  Inject test data to KV store"
	@echo "  make test-telegram-alert      Test Telegram alert system"
	@echo ""
	@echo "KV Store Cleanup:"
	@echo "  make cleanup-kv-preview       Preview what will be deleted"
	@echo "  make cleanup-kv-jobs          Delete job reward records"
	@echo "  make cleanup-kv-proofs        Delete Merkle proofs"
	@echo "  make cleanup-kv-settlements   Delete settlement periods"
	@echo "  make cleanup-kv-all           Delete ALL mining data (⚠️  DANGEROUS)"
	@echo ""
	@echo "Documentation (MkDocs):"
	@echo "  make docs-install        Install MkDocs in Python venv"
	@echo "  make docs-serve          Start local docs server (http://localhost:8000)"
	@echo "  make docs-build          Build static documentation site (./site/)"
	@echo "  make docs-build-website  Build docs to kawai-website/docs/"
	@echo "  make docs-clean          Clean documentation build artifacts"
	@echo "  make docs-deploy         Deploy documentation to GitHub Pages"
	@echo "  make docs-venv-clean     Remove Python virtual environment"
	@echo ""
	@echo "Website (kawai-website):"
	@echo "  make website-install  Install website dependencies (npm)"
	@echo "  make website-serve    Start local website server (http://localhost:3000)"
	@echo "  make website-dev      Start website with auto-reload"
	@echo "  make website-clean    Clean website node_modules"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean            Clean generated files"
	@echo "  make test             Run tests"
	@echo "  make update-llama     Update llama.cpp"

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

release-linux-deb:
	@echo "📦 Creating Linux .deb package..."
	@$(MAKE) release-linux
	@wails3 task linux:create:deb
	@echo "✅ Debian package complete!"
	@echo "📦 Location: build/bin/"

release-linux-rpm:
	@echo "📦 Creating Linux .rpm package..."
	@$(MAKE) release-linux
	@wails3 task linux:create:rpm
	@echo "✅ RPM package complete!"
	@echo "📦 Location: build/bin/"

release-linux-appimage:
	@echo "📦 Creating Linux AppImage..."
	@$(MAKE) release-linux
	@wails3 task linux:create:appimage
	@echo "✅ AppImage complete!"
	@echo "📦 Location: build/bin/"

release-linux-archive:
	@echo "📦 Creating Linux distribution archive..."
	@$(MAKE) release-linux
	@cd build/bin && tar -czf Kawai-$(shell grep 'version:' ../config.yml | tail -1 | awk '{print $$2}' | tr -d '"')-linux-amd64.tar.gz Kawai
	@cd build/bin && shasum -a 256 Kawai-*-linux-*.tar.gz >> checksums.txt
	@echo "✅ Linux archive created with checksum!"
	@echo "📦 Location: build/bin/Kawai-*-linux-amd64.tar.gz"

release-windows:
	@echo "🪟 Building for Windows (amd64)..."
	@$(MAKE) release-prepare
	@echo "Building Windows application..."
	@wails3 task windows:build
	@echo "✅ Windows build complete!"
	@echo "📦 Location: build/bin/Kawai.exe"

release-windows-nsis:
	@echo "📦 Creating Windows NSIS installer..."
	@$(MAKE) release-windows
	@wails3 task windows:create:nsis:installer
	@echo "✅ NSIS installer complete!"
	@echo "📦 Location: build/bin/"

release-windows-msix:
	@echo "📦 Creating Windows MSIX package..."
	@$(MAKE) release-windows
	@wails3 task windows:create:msix:package
	@echo "✅ MSIX package complete!"
	@echo "📦 Location: build/bin/"

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

release-packages:
	@echo "📦 Creating packages for all platforms..."
	@$(MAKE) release-darwin-package
	@$(MAKE) release-linux-deb
	@$(MAKE) release-linux-appimage
	@$(MAKE) release-windows-nsis
	@echo ""
	@echo "✅ All packages complete!"
	@echo ""
	@echo "📦 Package artifacts:"
	@ls -lh build/bin/

release-clean:
	@echo "🧹 Cleaning release artifacts..."
	@rm -rf build/bin/*
	@rm -rf build/darwin/*.app
	@echo "✅ Release artifacts cleaned!"

# ==============================================================================
# Code Generation
# ==============================================================================
generate: db-generate bindings-generate constants-generate
	@echo "✅ All code generated!"

db-generate:
	@echo "🔄 Generating Go code from SQL queries..."
	sqlc generate
	@echo "✅ Database code generated!"

bindings-generate:
	@echo "🔄 Generating TypeScript bindings..."
	wails3 generate bindings -clean=true -ts
	@echo "✅ TypeScript bindings generated!"

constants-generate:
	@echo "🔄 Generating constants from .env..."
	go run cmd/obfuscator-gen/main.go
	@echo "✅ Constants generated!"

# ==============================================================================
# Database
# ==============================================================================
db-dump:
	@echo "💾 Dumping database to seed file..."
	@mkdir -p internal/database/seed
	@sqlite3 $(DATA_DIR)/veridium.db ".dump" > internal/database/seed/veridium_dump.sql
	@echo "✅ Database dumped!"

db-restore:
	@echo "📥 Restoring database from seed file..."
	@test -f internal/database/seed/veridium_dump.sql || (echo "❌ Seed file not found!" && exit 1)
	@mkdir -p $(DATA_DIR)
	@rm -f $(DATA_DIR)/veridium.db $(DATA_DIR)/veridium.db-shm $(DATA_DIR)/veridium.db-wal
	@sqlite3 $(DATA_DIR)/veridium.db < internal/database/seed/veridium_dump.sql
	@echo "✅ Database restored!"

# ==============================================================================
# Smart Contracts
# ==============================================================================
contracts-compile:
	@echo "🔨 Compiling smart contracts..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge build

contracts-bindings: abi-token abi-otcmarket abi-vault abi-revenue abi-referral abi-cashback abi-mining abi-stablecoin generate-project-abis
	@echo "✅ Contract bindings generated!"

contracts-update: contracts-compile contracts-bindings
	@echo "🚀 Contracts updated!"

# ✅ NEW: Full upgrade workflow with validation
contracts-upgrade:
	@echo "🔄 Full contract upgrade workflow..."
	@echo ""
	@echo "Step 1: Running tests..."
	@$(MAKE) contracts-test
	@echo ""
	@echo "Step 2: Compiling contracts..."
	@$(MAKE) contracts-compile
	@echo ""
	@echo "Step 3: Generating bindings..."
	@$(MAKE) contracts-bindings
	@echo ""
	@echo "✅ Contract upgrade complete!"
	@echo ""
	@echo "⚠️  Next steps:"
	@echo "  1. Review changes in internal/generate/abi/"
	@echo "  2. Update backend code if needed"
	@echo "  3. Test locally: make dev-hot"
	@echo "  4. Deploy to testnet: make contracts-deploy-testnet"

# ✅ NEW: Contract testing
contracts-test:
	@echo "🧪 Running contract tests..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge test -vv

contracts-test-gas:
	@echo "⛽ Running contract tests with gas report..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge test --gas-report

contracts-coverage:
	@echo "📊 Running contract coverage..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge coverage

# ✅ NEW: Contract validation
contracts-validate:
	@echo "🔍 Validating contract changes..."
	@echo "Checking if contracts compiled..."
	@test -f $(OTCMARKET_ARTIFACT) || (echo "❌ OTCMarket artifact not found! Run: make contracts-compile" && exit 1)
	@test -f $(TOKEN_ARTIFACT) || (echo "❌ Token artifact not found! Run: make contracts-compile" && exit 1)
	@echo "✅ Contract artifacts found"
	@echo ""
	@echo "Checking if bindings generated..."
	@test -f $(ABIS_DIR)/otcmarket/otcmarket.go || (echo "❌ OTCMarket bindings not found! Run: make contracts-bindings" && exit 1)
	@test -f $(ABIS_DIR)/kawaitoken/kawaitoken.go || (echo "❌ Token bindings not found! Run: make contracts-bindings" && exit 1)
	@echo "✅ Contract bindings found"
	@echo ""
	@echo "Running contract tests..."
	@$(MAKE) contracts-test
	@echo ""
	@echo "✅ All validations passed!"

# ✅ NEW: Gas optimization
contracts-gas-snapshot:
	@echo "📸 Creating gas snapshot..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge snapshot
	@echo "✅ Gas snapshot saved to contracts/.gas-snapshot"

contracts-gas-compare:
	@echo "📊 Comparing gas usage..."
	@test -f $(CONTRACTS_DIR)/.gas-snapshot || (echo "❌ No baseline snapshot! Run: make contracts-gas-snapshot" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge snapshot --diff .gas-snapshot

# ✅ NEW: Deployment commands
contracts-deploy-local:
	@echo "🚀 Deploying to local Anvil..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployKawai.s.sol:DeployKawai \
		--rpc-url http://localhost:8545 \
		--private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
		--broadcast

contracts-deploy-testnet:
	@echo "🚀 Deploying to Monad Testnet..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployKawai.s.sol:DeployKawai \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify

contracts-deploy-vault:
	@echo "🚀 Deploying PaymentVault..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	@test -n "$(STABLECOIN_ADDRESS)" || (echo "❌ STABLECOIN_ADDRESS not set! Set it in contracts/.env" && exit 1)
	@if [ "$(CHAIN_ID)" = "143" ] || echo "$(RPC_URL)" | grep -q "mainnet"; then \
		echo "⚠️  WARNING: Deploying to MAINNET with USDC $(STABLECOIN_ADDRESS)"; \
		echo "⚠️  This will cost real MON for gas fees"; \
		read -p "Continue with mainnet deployment? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1; \
	fi
	@echo "ℹ️  Stablecoin: $(STABLECOIN_ADDRESS)"
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployPaymentVault.s.sol:DeployPaymentVault \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify
	@echo "✅ PaymentVault deployed!"
	@echo "⚠️  Update PAYMENT_VAULT_ADDRESS in .env"

contracts-deploy-referral-testnet:
	@echo "🚀 Deploying ReferralRewardDistributor to Monad Testnet..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set! Set it in contracts/.env" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployReferralDistributor.s.sol:DeployReferralDistributor \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify

contracts-grant-minter-referral:
	@echo "🔐 Granting MINTER_ROLE to ReferralRewardDistributor..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set!" && exit 1)
	@test -n "$(REFERRAL_DISTRIBUTOR_ADDRESS)" || (echo "❌ REFERRAL_DISTRIBUTOR_ADDRESS not set!" && exit 1)
	@echo "Granting MINTER_ROLE to $(REFERRAL_DISTRIBUTOR_ADDRESS)..."
	cd $(CONTRACTS_DIR) && cast send $(KAWAI_TOKEN_ADDRESS) \
		"grantRole(bytes32,address)" \
		$$(cast keccak "MINTER_ROLE") \
		$(REFERRAL_DISTRIBUTOR_ADDRESS) \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY)
	@echo "✅ MINTER_ROLE granted!"

contracts-deploy-cashback-testnet:
	@echo "🚀 Deploying DepositCashbackDistributor to Monad Testnet..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set! Set it in contracts/.env" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployCashbackDistributor.s.sol:DeployCashbackDistributor \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify

contracts-grant-minter-cashback:
	@echo "🔐 Granting MINTER_ROLE to DepositCashbackDistributor..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set!" && exit 1)
	@test -n "$(CASHBACK_DISTRIBUTOR_ADDRESS)" || (echo "❌ CASHBACK_DISTRIBUTOR_ADDRESS not set!" && exit 1)
	@echo "Granting MINTER_ROLE to $(CASHBACK_DISTRIBUTOR_ADDRESS)..."
	cd $(CONTRACTS_DIR) && cast send $(KAWAI_TOKEN_ADDRESS) \
		"grantRole(bytes32,address)" \
		$$(cast keccak "MINTER_ROLE") \
		$(CASHBACK_DISTRIBUTOR_ADDRESS) \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY)
	@echo "✅ MINTER_ROLE granted!"

contracts-deploy-mining-testnet:
	@echo "🚀 Deploying MiningRewardDistributor to Monad Testnet..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set! Set it in contracts/.env" && exit 1)
	@echo "ℹ️  Note: Developer addresses are specified per claim (flexible distribution)"
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployMiningDistributor.s.sol:DeployMiningDistributor \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify

contracts-grant-minter-mining:
	@echo "🔐 Granting MINTER_ROLE to MiningRewardDistributor..."
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(RPC_URL)" || (echo "❌ RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set!" && exit 1)
	@test -n "$(MINING_DISTRIBUTOR_ADDRESS)" || (echo "❌ MINING_DISTRIBUTOR_ADDRESS not set!" && exit 1)
	@echo "Granting MINTER_ROLE to $(MINING_DISTRIBUTOR_ADDRESS)..."
	cd $(CONTRACTS_DIR) && cast send $(KAWAI_TOKEN_ADDRESS) \
		"grantRole(bytes32,address)" \
		$$(cast keccak "MINTER_ROLE") \
		$(MINING_DISTRIBUTOR_ADDRESS) \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY)
	@echo "✅ MINTER_ROLE granted!"

contracts-verify:
	@echo "✅ Verifying contracts on explorer..."
	@test -n "$(CONTRACT_ADDRESS)" || (echo "❌ CONTRACT_ADDRESS not set!" && exit 1)
	@test -n "$(ETHERSCAN_API_KEY)" || (echo "❌ ETHERSCAN_API_KEY not set!" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge verify-contract \
		$(CONTRACT_ADDRESS) \
		contracts/Escrow.sol:OTCMarket \
		--chain-id $(CHAIN_ID) \
		--etherscan-api-key $(ETHERSCAN_API_KEY)

# ✅ NEW: Clean contract artifacts
contracts-clean:
	@echo "🧹 Cleaning contract artifacts..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge clean
	rm -rf $(ABIS_DIR)
	@echo "✅ Contract artifacts cleaned!"

generate-project-abis:
	@echo "📦 Injecting project ABIs into Jarvis..."
	@go run ./cmd/generate-abis pkg/jarvis/common/project_abis.go $(ABIS_DIR)

abi-token:
	@mkdir -p $(ABIS_DIR)/kawaitoken
	@jq -r .abi $(TOKEN_ARTIFACT) > $(ABIS_DIR)/kawaitoken/KawaiToken.abi
	@jq -r .bytecode.object $(TOKEN_ARTIFACT) > $(ABIS_DIR)/kawaitoken/KawaiToken.bin
	@abigen --abi $(ABIS_DIR)/kawaitoken/KawaiToken.abi --bin $(ABIS_DIR)/kawaitoken/KawaiToken.bin \
		--pkg kawaitoken --type KawaiToken --out $(ABIS_DIR)/kawaitoken/kawaitoken.go

abi-otcmarket:
	@mkdir -p $(ABIS_DIR)/otcmarket
	@jq -r .abi $(OTCMARKET_ARTIFACT) > $(ABIS_DIR)/otcmarket/OTCMarket.abi
	@jq -r .bytecode.object $(OTCMARKET_ARTIFACT) > $(ABIS_DIR)/otcmarket/OTCMarket.bin
	@abigen --abi $(ABIS_DIR)/otcmarket/OTCMarket.abi --bin $(ABIS_DIR)/otcmarket/OTCMarket.bin \
		--pkg otcmarket --type OTCMarket --out $(ABIS_DIR)/otcmarket/otcmarket.go

abi-vault:
	@mkdir -p $(ABIS_DIR)/vault
	@jq -r .abi $(VAULT_ARTIFACT) > $(ABIS_DIR)/vault/PaymentVault.abi
	@jq -r .bytecode.object $(VAULT_ARTIFACT) > $(ABIS_DIR)/vault/PaymentVault.bin
	@abigen --abi $(ABIS_DIR)/vault/PaymentVault.abi --bin $(ABIS_DIR)/vault/PaymentVault.bin \
		--pkg vault --type PaymentVault --out $(ABIS_DIR)/vault/vault.go

abi-revenue:
	@mkdir -p $(ABIS_DIR)/revenuedistributor
	@jq -r .abi $(REVENUE_ARTIFACT) > $(ABIS_DIR)/revenuedistributor/RevenueDistributor.abi
	@jq -r .bytecode.object $(REVENUE_ARTIFACT) > $(ABIS_DIR)/revenuedistributor/RevenueDistributor.bin
	@abigen --abi $(ABIS_DIR)/revenuedistributor/RevenueDistributor.abi --bin $(ABIS_DIR)/revenuedistributor/RevenueDistributor.bin \
		--pkg revenuedistributor --type RevenueDistributor --out $(ABIS_DIR)/revenuedistributor/revenuedistributor.go

abi-referral:
	@mkdir -p $(ABIS_DIR)/referraldistributor
	@jq -r .abi $(REFERRAL_ARTIFACT) > $(ABIS_DIR)/referraldistributor/ReferralRewardDistributor.abi
	@jq -r .bytecode.object $(REFERRAL_ARTIFACT) > $(ABIS_DIR)/referraldistributor/ReferralRewardDistributor.bin
	@abigen --abi $(ABIS_DIR)/referraldistributor/ReferralRewardDistributor.abi --bin $(ABIS_DIR)/referraldistributor/ReferralRewardDistributor.bin \
		--pkg referraldistributor --type ReferralRewardDistributor --out $(ABIS_DIR)/referraldistributor/referraldistributor.go

abi-cashback:
	@mkdir -p $(ABIS_DIR)/cashbackdistributor
	@jq -r .abi $(CASHBACK_ARTIFACT) > $(ABIS_DIR)/cashbackdistributor/DepositCashbackDistributor.abi
	@jq -r .bytecode.object $(CASHBACK_ARTIFACT) > $(ABIS_DIR)/cashbackdistributor/DepositCashbackDistributor.bin
	@abigen --abi $(ABIS_DIR)/cashbackdistributor/DepositCashbackDistributor.abi --bin $(ABIS_DIR)/cashbackdistributor/DepositCashbackDistributor.bin \
		--pkg cashbackdistributor --type DepositCashbackDistributor --out $(ABIS_DIR)/cashbackdistributor/cashbackdistributor.go

abi-mining:
	@mkdir -p $(ABIS_DIR)/miningdistributor
	@jq -r .abi $(MINING_ARTIFACT) > $(ABIS_DIR)/miningdistributor/MiningRewardDistributor.abi
	@jq -r .bytecode.object $(MINING_ARTIFACT) > $(ABIS_DIR)/miningdistributor/MiningRewardDistributor.bin
	@abigen --abi $(ABIS_DIR)/miningdistributor/MiningRewardDistributor.abi --bin $(ABIS_DIR)/miningdistributor/MiningRewardDistributor.bin \
		--pkg miningdistributor --type MiningRewardDistributor --out $(ABIS_DIR)/miningdistributor/miningdistributor.go

abi-stablecoin:
	@mkdir -p $(ABIS_DIR)/mockstablecoin
	@jq -r .abi $(STABLECOIN_ARTIFACT) > $(ABIS_DIR)/mockstablecoin/MockStablecoin.abi
	@jq -r .bytecode.object $(STABLECOIN_ARTIFACT) > $(ABIS_DIR)/mockstablecoin/MockStablecoin.bin
	@abigen --abi $(ABIS_DIR)/mockstablecoin/MockStablecoin.abi --bin $(ABIS_DIR)/mockstablecoin/MockStablecoin.bin \
		--pkg mockstablecoin --type MockStablecoin --out $(ABIS_DIR)/mockstablecoin/mockstablecoin.go

# ==============================================================================
# Maintenance
# ==============================================================================
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf frontend/bindings/
	rm -rf build/bin
	rm -rf $(ABIS_DIR)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge clean 2>/dev/null || true
	@echo "✅ Clean complete!"

test:
	@echo "🧪 Running tests..."
	go test ./...

update-llama:
	@echo "🔧 Updating llama.cpp..."
	@./scripts/update-llama.sh

update-llama-force:
	@echo "🔧 Force updating llama.cpp..."
	@./scripts/update-llama.sh --force

llama-version:
	@go run ./cmd/update-llama/main.go --version

llama-versions:
	@go run ./cmd/update-llama/main.go --list

watch:
	@echo "👀 Watching for changes in queries..."
	find internal/database/queries -name "*.sql" | entr -r make generate

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
# Emergency Pause Operations
# ==============================================================================
pause-status:
	@echo "📊 Checking pause status..."
	@go run cmd/adminops/pause/main.go -action status

pause-all:
	@echo "🚨 PAUSING ALL DISTRIBUTORS..."
	@go run cmd/adminops/pause/main.go -action pause -contract all

unpause-all:
	@echo "🔓 UNPAUSING ALL DISTRIBUTORS..."
	@go run cmd/adminops/pause/main.go -action unpause -contract all

pause-mining:
	@echo "🚨 Pausing Mining Distributor..."
	@go run cmd/adminops/pause/main.go -action pause -contract mining

unpause-mining:
	@echo "🔓 Unpausing Mining Distributor..."
	@go run cmd/adminops/pause/main.go -action unpause -contract mining

pause-cashback:
	@echo "🚨 Pausing Cashback Distributor..."
	@go run cmd/adminops/pause/main.go -action pause -contract cashback

unpause-cashback:
	@echo "🔓 Unpausing Cashback Distributor..."
	@go run cmd/adminops/pause/main.go -action unpause -contract cashback

pause-referral:
	@echo "🚨 Pausing Referral Distributor..."
	@go run cmd/adminops/pause/main.go -action pause -contract referral

unpause-referral:
	@echo "🔓 Unpausing Referral Distributor..."
	@go run cmd/adminops/pause/main.go -action unpause -contract referral

pause-all-dry:
	@echo "🔍 DRY RUN: Pausing all distributors..."
	@go run cmd/adminops/pause/main.go -action pause -contract all -dry-run

unpause-all-dry:
	@echo "🔍 DRY RUN: Unpausing all distributors..."
	@go run cmd/adminops/pause/main.go -action unpause -contract all -dry-run

# ==============================================================================
# Unified Reward Settlement
# ==============================================================================
reward-settlement-generate:
	@echo "🌳 Generating reward settlement..."
	@test -n "$(TYPE)" || (echo "❌ TYPE not set! Usage: make reward-settlement-generate TYPE=mining|cashback|referral" && exit 1)
	@go run cmd/reward-settlement/main.go generate --type $(TYPE)

reward-settlement-upload:
	@echo "🚀 Uploading Merkle root to contract..."
	@test -n "$(TYPE)" || (echo "❌ TYPE not set! Usage: make reward-settlement-upload TYPE=mining|cashback|referral" && exit 1)
	@go run cmd/reward-settlement/main.go upload --type $(TYPE)

reward-settlement-status:
	@echo "📊 Checking settlement status..."
	@go run cmd/reward-settlement/main.go status

reward-settlement-all:
	@echo "🚀 Settling all reward types..."
	@go run cmd/reward-settlement/main.go all

# Convenience shortcuts
settle-mining:
	@make reward-settlement-generate TYPE=mining

settle-cashback:
	@make reward-settlement-generate TYPE=cashback

settle-referral:
	@make reward-settlement-generate TYPE=referral

settle-revenue:
	@make reward-settlement-generate TYPE=revenue

settle-all:
	@make reward-settlement-all

# ==============================================================================
# Testing Helpers (Legacy - use reward-settlement instead)
# ==============================================================================
test-mining-rewards:
	@echo "🧪 Running mining rewards tests..."
	@bash scripts/test-mining-rewards.sh

test-inject-mining-data:
	@echo "💉 Injecting test mining reward data..."
	@go run cmd/dev/test-inject-mining-data/main.go

test-telegram-alert:
	@echo "📱 Testing Telegram alert system..."
	@go run cmd/dev/test-telegram-alert/main.go

inject-test-usdt:
	@echo "💵 Injecting test USDT to PaymentVault..."
	@go run cmd/dev/inject-test-usdt/main.go

mint-test-kawai:
	@echo "🪙 Minting test KAWAI tokens..."
	@go run cmd/dev/mint-test-kawai/main.go

send-test-mon:
	@echo "💰 Sending MON testnet tokens for gas fees..."
	@test -n "$(ADDR)" || (echo "❌ ADDR not set! Usage: make send-test-mon ADDR=0x... [AMOUNT=0.1]" && exit 1)
	@go run cmd/dev/send-test-mon/main.go $(ADDR) $(AMOUNT)

test-claiming-data:
	@echo "🎯 Testing claiming data for address..."
	@test -n "$(ADDR)" || (echo "❌ ADDR not set! Usage: make test-claiming-data ADDR=0x..." && exit 1)
	@go run cmd/dev/test-claiming-data/main.go $(ADDR)

# ==============================================================================
# Testing Helpers
# ==============================================================================
cleanup-test-data:
	@echo "🧹 Cleaning up test data from Cloudflare KV..."
	@go run cmd/dev/cleanup-test-data/main.go --confirm

check-minter-role:
	@echo "🔐 Checking MINTER_ROLE status..."
	@go run cmd/dev/check-minter-role/main.go

check-balance:
	@echo "💰 Checking KAWAI balance..."
	@test -n "$(ADDR)" || (echo "❌ ADDR not set! Usage: make check-balance ADDR=0x..." && exit 1)
	@go run cmd/dev/check-balance/main.go $(ADDR)

check-claim-status:
	@echo "🔍 Checking claim status..."
	@test -n "$(TYPE)" || (echo "❌ TYPE not set! Usage: make check-claim-status TYPE=mining|cashback|referral PERIOD=123 ADDR=0x..." && exit 1)
	@test -n "$(PERIOD)" || (echo "❌ PERIOD not set! Usage: make check-claim-status TYPE=mining|cashback|referral PERIOD=123 ADDR=0x..." && exit 1)
	@test -n "$(ADDR)" || (echo "❌ ADDR not set! Usage: make check-claim-status TYPE=mining|cashback|referral PERIOD=123 ADDR=0x..." && exit 1)
	@go run cmd/dev/check-claim-status/main.go $(TYPE) $(PERIOD) $(ADDR)

upload-merkle-root:
	@echo "🚀 Uploading Merkle root..."
	@test -n "$(TYPE)" || (echo "❌ TYPE not set! Usage: make upload-merkle-root TYPE=mining|cashback|referral ROOT=0x..." && exit 1)
	@test -n "$(ROOT)" || (echo "❌ ROOT not set! Usage: make upload-merkle-root TYPE=mining|cashback|referral ROOT=0x..." && exit 1)
	@go run cmd/dev/upload-merkle-root/main.go $(TYPE) $(ROOT)

# ==============================================================================
# KV Store Cleanup
# ==============================================================================
cleanup-kv-preview:
	@echo "🔍 Preview KV cleanup (dry-run)..."
	@go run cmd/cleanup-kv-mining-data/main.go --all --dry-run

cleanup-kv-jobs:
	@echo "🧹 Cleaning up job reward records..."
	@go run cmd/cleanup-kv-mining-data/main.go --jobs --confirm DELETE

cleanup-kv-proofs:
	@echo "🧹 Cleaning up Merkle proofs..."
	@go run cmd/cleanup-kv-mining-data/main.go --proofs --confirm DELETE

cleanup-kv-settlements:
	@echo "🧹 Cleaning up settlement periods..."
	@go run cmd/cleanup-kv-mining-data/main.go --settlements --confirm DELETE

cleanup-kv-all:
	@echo "🧹 Cleaning up ALL mining data from KV..."
	@echo "⚠️  This will delete all job records, proofs, and settlements!"
	@go run cmd/cleanup-kv-mining-data/main.go --all --confirm DELETE

# ==============================================================================
# Documentation (MkDocs)
# ==============================================================================

# Python virtual environment for docs
VENV_DIR := venv
PYTHON := python3
PIP := $(VENV_DIR)/bin/pip
MKDOCS := $(VENV_DIR)/bin/mkdocs

docs-install:
	@echo "📚 Setting up documentation environment..."
	@if [ ! -d "$(VENV_DIR)" ]; then \
		echo "Creating Python virtual environment..."; \
		$(PYTHON) -m venv $(VENV_DIR); \
	fi
	@echo "Installing MkDocs and dependencies..."
	@$(PIP) install --upgrade pip
	@$(PIP) install mkdocs-material pymdown-extensions
	@echo "✅ MkDocs installed successfully!"
	@echo "💡 Use 'make docs-serve' to start the documentation server"

docs-serve:
	@if [ ! -f "$(MKDOCS)" ]; then \
		echo "❌ MkDocs not installed. Run 'make docs-install' first."; \
		exit 1; \
	fi
	@echo "🔍 Checking for processes on port 8000..."
	@lsof -ti:8000 | xargs kill -9 2>/dev/null || true
	@echo "📖 Starting MkDocs development server..."
	@echo "🌐 Open http://localhost:8000 in your browser"
	@echo "⏹️  Press Ctrl+C to stop"
	@$(MKDOCS) serve

docs-build:
	@if [ ! -f "$(MKDOCS)" ]; then \
		echo "❌ MkDocs not installed. Run 'make docs-install' first."; \
		exit 1; \
	fi
	@echo "🔨 Building static documentation site..."
	@$(MKDOCS) build
	@echo "✅ Documentation built to ./site/"

docs-build-website:
	@if [ ! -f "$(MKDOCS)" ]; then \
		echo "❌ MkDocs not installed. Run 'make docs-install' first."; \
		exit 1; \
	fi
	@if [ ! -d "$(WEBSITE_DIR)" ]; then \
		echo "❌ kawai-website directory not found!"; \
		exit 1; \
	fi
	@echo "🗑️  Removing old docs..."
	@rm -rf $(WEBSITE_DIR)/docs
	@echo "🔨 Building documentation to kawai-website/docs/..."
	@$(MKDOCS) build -d $(WEBSITE_DIR)/docs
	@echo "✅ Documentation built to $(WEBSITE_DIR)/docs/"
	@echo "💡 Next: cd $(WEBSITE_DIR) && git add docs && git commit && git push"

docs-clean:
	@echo "🧹 Cleaning documentation build artifacts..."
	@rm -rf site/
	@echo "✅ Documentation cleaned!"

docs-deploy:
	@if [ ! -f "$(MKDOCS)" ]; then \
		echo "❌ MkDocs not installed. Run 'make docs-install' first."; \
		exit 1; \
	fi
	@echo "🚀 Deploying documentation to GitHub Pages..."
	@$(MKDOCS) gh-deploy --force
	@echo "✅ Documentation deployed!"

docs-venv-clean:
	@echo "🧹 Removing Python virtual environment..."
	@rm -rf $(VENV_DIR)
	@echo "✅ Virtual environment removed!"

# ==============================================================================
# Website (kawai-website)
# ==============================================================================

WEBSITE_DIR := kawai-website

website-install:
	@echo "📦 Installing website dependencies..."
	@if [ ! -d "$(WEBSITE_DIR)" ]; then \
		echo "❌ $(WEBSITE_DIR) directory not found!"; \
		echo "💡 Make sure kawai-website is cloned in the same parent directory"; \
		exit 1; \
	fi
	@cd $(WEBSITE_DIR) && npm install
	@echo "✅ Website dependencies installed!"

website-serve:
	@if [ ! -d "$(WEBSITE_DIR)" ]; then \
		echo "❌ $(WEBSITE_DIR) directory not found!"; \
		exit 1; \
	fi
	@if [ ! -d "$(WEBSITE_DIR)/node_modules" ]; then \
		echo "❌ Dependencies not installed. Run 'make website-install' first."; \
		exit 1; \
	fi
	@echo "🌐 Starting website server..."
	@echo "🔗 Open http://localhost:3000 in your browser"
	@echo "⏹️  Press Ctrl+C to stop"
	@cd $(WEBSITE_DIR) && npm run serve

website-dev:
	@if [ ! -d "$(WEBSITE_DIR)" ]; then \
		echo "❌ $(WEBSITE_DIR) directory not found!"; \
		exit 1; \
	fi
	@if [ ! -d "$(WEBSITE_DIR)/node_modules" ]; then \
		echo "❌ Dependencies not installed. Run 'make website-install' first."; \
		exit 1; \
	fi
	@echo "🔥 Starting website with auto-reload..."
	@echo "🔗 Open http://localhost:3000 in your browser"
	@echo "⏹️  Press Ctrl+C to stop"
	@cd $(WEBSITE_DIR) && npm run dev

website-clean:
	@echo "🧹 Cleaning website dependencies..."
	@if [ -d "$(WEBSITE_DIR)/node_modules" ]; then \
		cd $(WEBSITE_DIR) && rm -rf node_modules; \
		echo "✅ Website dependencies cleaned!"; \
	else \
		echo "⚠️  No node_modules to clean"; \
	fi


# ✅ MAINNET DEPLOYMENT TARGETS (with safety confirmations)
contracts-deploy-mainnet:
	@echo "🚀 Deploying to Monad Mainnet..."
	@echo "⚠️  WARNING: This will deploy to PRODUCTION mainnet!"
	@echo "⚠️  Make sure you are using the correct private key and RPC URL"
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(MAINNET_RPC_URL)" || (echo "❌ MAINNET_RPC_URL not set!" && exit 1)
	@echo "Deploying to mainnet..."
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployKawai.s.sol:DeployKawai \
		--rpc-url $(MAINNET_RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify
	@echo "✅ Deployment complete! Update .env.mainnet with deployed addresses"

contracts-deploy-mining-mainnet:
	@echo "🚀 Deploying MiningRewardDistributor to Monad Mainnet..."
	@echo "⚠️  WARNING: This will deploy to PRODUCTION mainnet!"
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(MAINNET_RPC_URL)" || (echo "❌ MAINNET_RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set! Set it in contracts/.env.mainnet" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployMiningDistributor.s.sol:DeployMiningDistributor \
		--rpc-url $(MAINNET_RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify

contracts-deploy-cashback-mainnet:
	@echo "🚀 Deploying DepositCashbackDistributor to Monad Mainnet..."
	@echo "⚠️  WARNING: This will deploy to PRODUCTION mainnet!"
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(MAINNET_RPC_URL)" || (echo "❌ MAINNET_RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set! Set it in contracts/.env.mainnet" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployCashbackDistributor.s.sol:DeployCashbackDistributor \
		--rpc-url $(MAINNET_RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify

contracts-deploy-referral-mainnet:
	@echo "🚀 Deploying ReferralRewardDistributor to Monad Mainnet..."
	@echo "⚠️  WARNING: This will deploy to PRODUCTION mainnet!"
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(MAINNET_RPC_URL)" || (echo "❌ MAINNET_RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set! Set it in contracts/.env.mainnet" && exit 1)
	cd $(CONTRACTS_DIR) && ~/.foundry/bin/forge script script/DeployReferralDistributor.s.sol:DeployReferralDistributor \
		--rpc-url $(MAINNET_RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		--verify

# Grant MINTER_ROLE on mainnet (with safety confirmation)
contracts-grant-minter-mainnet:
	@echo "🔐 Granting MINTER_ROLE on Monad Mainnet..."
	@echo "⚠️  WARNING: This will grant MINTER_ROLE on PRODUCTION mainnet!"
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@test -n "$(PRIVATE_KEY)" || (echo "❌ PRIVATE_KEY not set!" && exit 1)
	@test -n "$(MAINNET_RPC_URL)" || (echo "❌ MAINNET_RPC_URL not set!" && exit 1)
	@test -n "$(KAWAI_TOKEN_ADDRESS)" || (echo "❌ KAWAI_TOKEN_ADDRESS not set!" && exit 1)
	@test -n "$(MINING_DISTRIBUTOR_ADDRESS)" || (echo "❌ MINING_DISTRIBUTOR_ADDRESS not set!" && exit 1)
	@test -n "$(CASHBACK_DISTRIBUTOR_ADDRESS)" || (echo "❌ CASHBACK_DISTRIBUTOR_ADDRESS not set!" && exit 1)
	@test -n "$(REFERRAL_DISTRIBUTOR_ADDRESS)" || (echo "❌ REFERRAL_DISTRIBUTOR_ADDRESS not set!" && exit 1)
	@echo "Granting MINTER_ROLE to Mining Distributor..."
	cd $(CONTRACTS_DIR) && cast send $(KAWAI_TOKEN_ADDRESS) \
		"grantRole(bytes32,address)" \
		$$(cast keccak "MINTER_ROLE") \
		$(MINING_DISTRIBUTOR_ADDRESS) \
		--rpc-url $(MAINNET_RPC_URL) \
		--private-key $(PRIVATE_KEY)
	@echo "Granting MINTER_ROLE to Cashback Distributor..."
	cd $(CONTRACTS_DIR) && cast send $(KAWAI_TOKEN_ADDRESS) \
		"grantRole(bytes32,address)" \
		$$(cast keccak "MINTER_ROLE") \
		$(CASHBACK_DISTRIBUTOR_ADDRESS) \
		--rpc-url $(MAINNET_RPC_URL) \
		--private-key $(PRIVATE_KEY)
	@echo "Granting MINTER_ROLE to Referral Distributor..."
	cd $(CONTRACTS_DIR) && cast send $(KAWAI_TOKEN_ADDRESS) \
		"grantRole(bytes32,address)" \
		$$(cast keccak "MINTER_ROLE") \
		$(REFERRAL_DISTRIBUTOR_ADDRESS) \
		--rpc-url $(MAINNET_RPC_URL) \
		--private-key $(PRIVATE_KEY)
	@echo "✅ All MINTER_ROLEs granted on mainnet!"

# ==============================================================================
# Full Deployment with Logging
# ==============================================================================

.PHONY: deploy-testnet deploy-mainnet

deploy-testnet:
	@echo "🚀 Starting full testnet deployment..."
	@rm -f contract-deploy-testnet.log
	@./contracts/deploy-all.sh testnet 2>&1 | tee contract-deploy-testnet.log
	@echo "✅ Deployment complete! Log saved to contract-deploy-testnet.log"

deploy-mainnet:
	@echo "🚀 Starting full mainnet deployment..."
	@rm -f contract-deploy-mainnet.log
	@./contracts/deploy-all.sh mainnet 2>&1 | tee contract-deploy-mainnet.log
	@echo "✅ Deployment complete! Log saved to contract-deploy-mainnet.log"
