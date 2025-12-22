.PHONY: help db-generate bindings-generate dev build clean test update-llama db-dump db-restore

# Default target
help:
	@echo "Veridium Development Commands"
	@echo ""
	@echo "Database:"
	@echo "  make db-generate          Generate Go code from SQL queries (sqlc)"
	@echo "  make db-dump              Dump SQLite database to seed file"
	@echo "  make db-restore           Restore database from seed file"
	@echo "  make bindings-generate    Generate TypeScript bindings from Go (wails)"
	@echo "  make generate             Run both db-generate and bindings-generate"
	@echo ""
	@echo "Development:"
	@echo "  make dev                  Start development server (full build)"
	@echo "  make devd                 Start dev server (keep existing DB)"
	@echo ""
	@echo "Quick Development (skip frontend build if possible):"
	@echo "  make dev-quick            Skip build if dist exists"
	@echo "  make devd-quick           Skip build if dist exists + keep DB"
	@echo "  make dev-smart            Reuse Vite OR skip build (smartest)"
	@echo "  make devd-smart           Reuse Vite OR skip build + keep DB"
	@echo ""
	@echo "Manual Control:"
	@echo "  make dev-skip-build       Skip 'bun run build:dev' (use dist)"
	@echo "  make devd-skip-build      Skip build + keep DB"
	@echo "  make dev-skip-frontend    Skip Vite dev server (manual start)"
	@echo "  make devd-skip-frontend   Skip Vite + keep DB"
	@echo ""
	@echo "Build:"
	@echo "  make build                Build production binary"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean                Clean generated files and build artifacts"
	@echo "  make test                 Run tests"
	@echo "  make update-llama         Update llama.cpp to latest version"

# Generate Go code from SQL queries using sqlc
db-generate:
	@echo "🔄 Generating Go code from SQL queries..."
	sqlc generate
	@echo "✅ Database code generated!"

# Dump SQLite database to seed file
db-dump:
	@echo "💾 Dumping database to seed file..."
	@mkdir -p internal/database/seed
	@sqlite3 data/veridium.db ".dump" > internal/database/seed/veridium_dump.sql
	@echo "✅ Database dumped to internal/database/seed/veridium_dump.sql"

# Restore database from seed file
db-restore:
	@echo "📥 Restoring database from seed file..."
	@if [ ! -f internal/database/seed/veridium_dump.sql ]; then \
		echo "❌ Error: Seed file not found at internal/database/seed/veridium_dump.sql"; \
		exit 1; \
	fi
	@mkdir -p data
	@rm -f data/veridium.db data/veridium.db-shm data/veridium.db-wal
	@sqlite3 data/veridium.db < internal/database/seed/veridium_dump.sql
	@echo "✅ Database restored from seed file"

# Generate TypeScript bindings from Go using Wails
bindings-generate:
	@echo "🔄 Generating TypeScript bindings..."
	wails3 generate bindings -clean=true -ts
	@echo "✅ TypeScript bindings generated!"

# Generate both database code and TypeScript bindings
generate: db-generate bindings-generate
	@echo "✅ All code generated successfully!"

# Start development server
dev:
	@echo "🚀 Starting development server..."
	killport 9245
	rm -rf data
	rm -rf files
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev 2>&1 | tee backend-dev.log

# Start development server without removing database
devd:
	@echo "🚀 Starting development server..."
	killport 9245
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev 2>&1 | tee backend-dev.log

# Start development server (skip frontend dev server, assume it's already running)
dev-skip-frontend:
	@echo "🚀 Starting development server (skipping frontend dev server)..."
	@echo "⚠️  Make sure Vite is already running on port 9245!"
	@echo "    Run in another terminal: cd frontend && bun run dev -- --port 9245"
	killport 9245 || true
	rm -rf data
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev -config ./build/config-skip-frontend.yml 2>&1 | tee backend-dev.log

# Start development server without removing database (skip frontend dev server)
devd-skip-frontend:
	@echo "🚀 Starting development server (skipping frontend dev server, keeping DB)..."
	@echo "⚠️  Make sure Vite is already running on port 9245!"
	@echo "    Run in another terminal: cd frontend && bun run dev -- --port 9245"
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev -config ./build/config-skip-frontend.yml 2>&1 | tee backend-dev.log

# Start development server (skip 'bun run build:dev', use existing dist)
dev-skip-build:
	@echo "🚀 Starting development server (skipping frontend build)..."
	@echo "⚠️  Requires existing frontend/dist from previous build!"
	@if [ ! -d frontend/dist ]; then \
		echo "❌ Error: frontend/dist not found!"; \
		echo "   Run 'make dev' first to build frontend."; \
		exit 1; \
	fi
	killport 9245
	rm -rf data
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev -config ./build/config-skip-frontend-build.yml 2>&1 | tee backend-dev.log

# Start development server without removing database (skip 'bun run build:dev')
devd-skip-build:
	@echo "🚀 Starting development server (skipping frontend build, keeping DB)..."
	@echo "⚠️  Requires existing frontend/dist from previous build!"
	@if [ ! -d frontend/dist ]; then \
		echo "❌ Error: frontend/dist not found!"; \
		echo "   Run 'make devd' first to build frontend."; \
		exit 1; \
	fi
	killport 9245
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev -config ./build/config-skip-frontend-build.yml 2>&1 | tee backend-dev.log

# Quick dev: skip 'bun run build:dev' if frontend/dist exists
dev-quick:
	@echo "⚡ Quick development mode..."
	@if [ -d frontend/dist ]; then \
		echo "✅ Found existing frontend/dist, skipping build..."; \
		$(MAKE) dev-skip-build; \
	else \
		echo "⚠️  No frontend/dist found, running full build..."; \
		$(MAKE) dev; \
	fi

# Quick dev without removing database
devd-quick:
	@echo "⚡ Quick development mode (keeping DB)..."
	@if [ -d frontend/dist ]; then \
		echo "✅ Found existing frontend/dist, skipping build..."; \
		$(MAKE) devd-skip-build; \
	else \
		echo "⚠️  No frontend/dist found, running full build..."; \
		$(MAKE) devd; \
	fi

# Smart dev: reuse Vite dev server if running, otherwise skip build if dist exists
dev-smart:
	@echo "🧠 Smart development mode..."
	@if lsof -Pi :9245 -sTCP:LISTEN -t >/dev/null 2>&1 ; then \
		echo "✅ Vite dev server already running on port 9245"; \
		echo "🚀 Reusing existing dev server..."; \
		rm -rf data; \
		rm -f backend-dev.log; \
		VERIDIUM_DEV=true wails3 dev -config ./build/config-skip-frontend.yml 2>&1 | tee backend-dev.log; \
	elif [ -d frontend/dist ]; then \
		echo "✅ Found existing frontend/dist"; \
		echo "🚀 Skipping build, using existing dist..."; \
		$(MAKE) dev-skip-build; \
	else \
		echo "⚠️  No Vite server or dist found, running full build..."; \
		$(MAKE) dev; \
	fi

# Smart dev without removing database
devd-smart:
	@echo "🧠 Smart development mode (keeping DB)..."
	@if lsof -Pi :9245 -sTCP:LISTEN -t >/dev/null 2>&1 ; then \
		echo "✅ Vite dev server already running on port 9245"; \
		echo "🚀 Reusing existing dev server..."; \
		rm -f backend-dev.log; \
		VERIDIUM_DEV=true wails3 dev -config ./build/config-skip-frontend.yml 2>&1 | tee backend-dev.log; \
	elif [ -d frontend/dist ]; then \
		echo "✅ Found existing frontend/dist"; \
		echo "🚀 Skipping build, using existing dist..."; \
		$(MAKE) devd-skip-build; \
	else \
		echo "⚠️  No Vite server or dist found, running full build..."; \
		$(MAKE) devd; \
	fi

# Build production binary
build:
	@echo "🔨 Building production binary..."
	wails3 build

# Clean generated files and build artifacts
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf frontend/bindings/
	rm -rf build/
	@echo "✅ Clean complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Watch for changes and regenerate (requires entr or similar)
watch:
	@echo "👀 Watching for changes in queries..."
	@echo "Install 'entr' first: brew install entr"
	find internal/database/queries -name "*.sql" | entr -r make generate

# Update llama.cpp to latest version
update-llama:
	@echo "🔧 Updating llama.cpp..."
	@./scripts/update-llama.sh

# Update llama.cpp (force reinstall)
update-llama-force:
	@echo "🔧 Force updating llama.cpp..."
	@./scripts/update-llama.sh --force

# Check llama.cpp version
llama-version:
	@echo "📦 Checking llama.cpp version..."
	@go run ./cmd/update-llama/main.go --version

# List available llama.cpp versions
llama-versions:
	@echo "📋 Listing available llama.cpp versions..."
	@go run ./cmd/update-llama/main.go --list

# Paths
CONTRACTS_DIR := contracts
ABIS_DIR := ./internal/generate/abi

# Artifacts
TOKEN_ARTIFACT := $(CONTRACTS_DIR)/artifacts/contracts/KawaiToken.sol/KawaiToken.json
ESCROW_ARTIFACT := $(CONTRACTS_DIR)/artifacts/contracts/Escrow.sol/OTCMarket.json
VAULT_ARTIFACT := $(CONTRACTS_DIR)/artifacts/contracts/PaymentVault.sol/PaymentVault.json

contracts-compile:
	@echo "Compiling smart contracts..."
	cd $(CONTRACTS_DIR) && npx hardhat compile

contracts-bindings: abi-token abi-escrow abi-vault generate-project-abis
	@echo "✅ Bindings generated."

contracts-update: contracts-compile contracts-bindings
	@echo "🚀 Contracts updated successfully (Compiled + Bindings Generated)"

generate-project-abis:
	@echo "Injecting project ABIs into Jarvis..."
	@go run ./cmd/generate-abis pkg/jarvis/common/project_abis.go $(ABIS_DIR)
	@echo "Project ABIs injected."

abi-token:
	@echo "Generating bindings for KawaiToken..."
	@mkdir -p $(ABIS_DIR)/kawaitoken
	@jq -r .abi $(TOKEN_ARTIFACT) > $(ABIS_DIR)/kawaitoken/KawaiToken.abi
	@jq -r .bytecode $(TOKEN_ARTIFACT) > $(ABIS_DIR)/kawaitoken/KawaiToken.bin
	@abigen --abi $(ABIS_DIR)/kawaitoken/KawaiToken.abi --bin $(ABIS_DIR)/kawaitoken/KawaiToken.bin --pkg kawaitoken --type KawaiToken --out $(ABIS_DIR)/kawaitoken/kawaitoken.go
	@echo "KawaiToken abi generated."

abi-escrow:
	@echo "Generating bindings for OTC Market (Escrow)..."
	@mkdir -p $(ABIS_DIR)/escrow
	@jq -r .abi $(ESCROW_ARTIFACT) > $(ABIS_DIR)/escrow/Escrow.abi
	@jq -r .bytecode $(ESCROW_ARTIFACT) > $(ABIS_DIR)/escrow/Escrow.bin
	@abigen --abi $(ABIS_DIR)/escrow/Escrow.abi --bin $(ABIS_DIR)/escrow/Escrow.bin --pkg escrow --type OTCMarket --out $(ABIS_DIR)/escrow/escrow.go
	@echo "OTC Market abi generated."

abi-vault:
	@echo "Generating bindings for PaymentVault..."
	@mkdir -p $(ABIS_DIR)/vault
	@jq -r .abi $(VAULT_ARTIFACT) > $(ABIS_DIR)/vault/PaymentVault.abi
	@jq -r .bytecode $(VAULT_ARTIFACT) > $(ABIS_DIR)/vault/PaymentVault.bin
	@abigen --abi $(ABIS_DIR)/vault/PaymentVault.abi --bin $(ABIS_DIR)/vault/PaymentVault.bin --pkg vault --type PaymentVault --out $(ABIS_DIR)/vault/vault.go
	@echo "PaymentVault abi generated."

clean:
	rm -rf $(ABIS_DIR)
	cd $(CONTRACTS_DIR) && npx hardhat clean