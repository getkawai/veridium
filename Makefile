# ==============================================================================
# Veridium Makefile
# ==============================================================================

.PHONY: help dev dev-fast dev-hot dev-rebuild build clean test generate \
        db-generate bindings-generate db-dump db-restore \
        contracts-compile contracts-bindings contracts-update

# ------------------------------------------------------------------------------
# Configuration
# ------------------------------------------------------------------------------
CONTRACTS_DIR := contracts
ABIS_DIR      := ./internal/generate/abi
DATA_DIR      := data

# Foundry artifacts
TOKEN_ARTIFACT       := $(CONTRACTS_DIR)/out/KawaiToken.sol/KawaiToken.json
ESCROW_ARTIFACT      := $(CONTRACTS_DIR)/out/Escrow.sol/OTCMarket.json
VAULT_ARTIFACT       := $(CONTRACTS_DIR)/out/PaymentVault.sol/PaymentVault.json
DISTRIBUTOR_ARTIFACT := $(CONTRACTS_DIR)/out/MerkleDistributor.sol/MerkleDistributor.json
USDT_ARTIFACT        := $(CONTRACTS_DIR)/out/MockUSDT.sol/MockUSDT.json

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
	@echo "Code Generation:"
	@echo "  make generate         Run db-generate + bindings-generate"
	@echo "  make db-generate      Generate Go code from SQL (sqlc)"
	@echo "  make bindings-generate Generate TypeScript bindings (wails)"
	@echo ""
	@echo "Database:"
	@echo "  make db-dump          Dump database to seed file"
	@echo "  make db-restore       Restore database from seed file"
	@echo ""
	@echo "Smart Contracts:"
	@echo "  make contracts-update Compile + generate bindings"
	@echo "  make contracts-compile Compile with Foundry"
	@echo "  make contracts-bindings Generate Go bindings"
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
	@echo "� Rebuilding frontend (keep DB)..."
	killport 9245 || true
	rm -f backend-dev.log
	VERIDIUM_DEV=1 wails3 dev 2>&1 | tee backend-dev.log

build:
	@echo "� Building production binary..."
	wails3 build

# ==============================================================================
# Code Generation
# ==============================================================================
generate: db-generate bindings-generate
	@echo "✅ All code generated!"

db-generate:
	@echo "🔄 Generating Go code from SQL queries..."
	sqlc generate
	@echo "✅ Database code generated!"

bindings-generate:
	@echo "🔄 Generating TypeScript bindings..."
	wails3 generate bindings -clean=true -ts
	@echo "✅ TypeScript bindings generated!"

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

contracts-bindings: abi-token abi-escrow abi-vault abi-distributor abi-usdt generate-project-abis
	@echo "✅ Contract bindings generated!"

contracts-update: contracts-compile contracts-bindings
	@echo "🚀 Contracts updated!"

generate-project-abis:
	@echo "📦 Injecting project ABIs into Jarvis..."
	@go run ./cmd/generate-abis pkg/jarvis/common/project_abis.go $(ABIS_DIR)

abi-token:
	@mkdir -p $(ABIS_DIR)/kawaitoken
	@jq -r .abi $(TOKEN_ARTIFACT) > $(ABIS_DIR)/kawaitoken/KawaiToken.abi
	@jq -r .bytecode.object $(TOKEN_ARTIFACT) > $(ABIS_DIR)/kawaitoken/KawaiToken.bin
	@abigen --abi $(ABIS_DIR)/kawaitoken/KawaiToken.abi --bin $(ABIS_DIR)/kawaitoken/KawaiToken.bin \
		--pkg kawaitoken --type KawaiToken --out $(ABIS_DIR)/kawaitoken/kawaitoken.go

abi-escrow:
	@mkdir -p $(ABIS_DIR)/escrow
	@jq -r .abi $(ESCROW_ARTIFACT) > $(ABIS_DIR)/escrow/Escrow.abi
	@jq -r .bytecode.object $(ESCROW_ARTIFACT) > $(ABIS_DIR)/escrow/Escrow.bin
	@abigen --abi $(ABIS_DIR)/escrow/Escrow.abi --bin $(ABIS_DIR)/escrow/Escrow.bin \
		--pkg escrow --type OTCMarket --out $(ABIS_DIR)/escrow/escrow.go

abi-vault:
	@mkdir -p $(ABIS_DIR)/vault
	@jq -r .abi $(VAULT_ARTIFACT) > $(ABIS_DIR)/vault/PaymentVault.abi
	@jq -r .bytecode.object $(VAULT_ARTIFACT) > $(ABIS_DIR)/vault/PaymentVault.bin
	@abigen --abi $(ABIS_DIR)/vault/PaymentVault.abi --bin $(ABIS_DIR)/vault/PaymentVault.bin \
		--pkg vault --type PaymentVault --out $(ABIS_DIR)/vault/vault.go

abi-distributor:
	@mkdir -p $(ABIS_DIR)/distributor
	@jq -r .abi $(DISTRIBUTOR_ARTIFACT) > $(ABIS_DIR)/distributor/MerkleDistributor.abi
	@jq -r .bytecode.object $(DISTRIBUTOR_ARTIFACT) > $(ABIS_DIR)/distributor/MerkleDistributor.bin
	@abigen --abi $(ABIS_DIR)/distributor/MerkleDistributor.abi --bin $(ABIS_DIR)/distributor/MerkleDistributor.bin \
		--pkg distributor --type MerkleDistributor --out $(ABIS_DIR)/distributor/distributor.go

abi-usdt:
	@mkdir -p $(ABIS_DIR)/usdt
	@jq -r .abi $(USDT_ARTIFACT) > $(ABIS_DIR)/usdt/MockUSDT.abi
	@jq -r .bytecode.object $(USDT_ARTIFACT) > $(ABIS_DIR)/usdt/MockUSDT.bin
	@abigen --abi $(ABIS_DIR)/usdt/MockUSDT.abi --bin $(ABIS_DIR)/usdt/MockUSDT.bin \
		--pkg usdt --type MockUSDT --out $(ABIS_DIR)/usdt/usdt.go

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