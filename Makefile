# ==============================================================================
# Veridium Makefile
# ==============================================================================

.PHONY: help dev dev-fast dev-hot dev-rebuild build clean test generate \
        db-generate bindings-generate db-dump db-restore \
        contracts-compile contracts-bindings contracts-update \
        contracts-test contracts-test-gas contracts-coverage \
        contracts-deploy-local contracts-deploy-testnet contracts-verify \
        contracts-upgrade contracts-clean contracts-validate \
        contracts-gas-snapshot contracts-gas-compare \
        admin-register admin-register-dry \
        docs-install docs-serve docs-build docs-clean docs-deploy

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
REFERRAL_ARTIFACT    := $(CONTRACTS_DIR)/out/ReferralRewardDistributor.sol/ReferralRewardDistributor.json
CASHBACK_ARTIFACT    := $(CONTRACTS_DIR)/out/DepositCashbackDistributor.sol/DepositCashbackDistributor.json
MINING_ARTIFACT      := $(CONTRACTS_DIR)/out/MiningRewardDistributor.sol/MiningRewardDistributor.json
USDT_ARTIFACT        := $(CONTRACTS_DIR)/out/MockUSDT.sol/MockUSDT.json

# Load environment variables from contracts/.env if exists
-include $(CONTRACTS_DIR)/.env
export

# Deployment configuration (can be overridden via env vars or contracts/.env)
PRIVATE_KEY ?= 
RPC_URL ?= https://testnet-rpc.monad.xyz
CONTRACT_ADDRESS ?=
ETHERSCAN_API_KEY ?= MKB28KJN1TJKRPA4EYVXXBWYUYDX6P5ESF
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
	@echo "Reward Settlement (Unified):"
	@echo "  make reward-settlement-generate TYPE=<type>  Generate settlement (mining|cashback|referral)"
	@echo "  make reward-settlement-upload TYPE=<type>    Upload Merkle root to contract"
	@echo "  make reward-settlement-status                Show status for all reward types"
	@echo "  make reward-settlement-all                   Settle all reward types at once"
	@echo "  make settle-mining                           Shortcut for mining settlement"
	@echo "  make settle-cashback                         Shortcut for cashback settlement"
	@echo "  make settle-referral                         Shortcut for referral settlement"
	@echo "  make settle-all                              Shortcut for all settlements"
	@echo ""
	@echo "Mining Rewards Testing:"
	@echo "  make test-mining-rewards      Run all mining rewards tests"
	@echo "  make test-inject-mining-data  Inject test data to KV store"
	@echo "  make test-mining-settlement   Test full settlement flow"
	@echo ""
	@echo "KV Store Cleanup:"
	@echo "  make cleanup-kv-preview       Preview what will be deleted"
	@echo "  make cleanup-kv-jobs          Delete job reward records"
	@echo "  make cleanup-kv-proofs        Delete Merkle proofs"
	@echo "  make cleanup-kv-settlements   Delete settlement periods"
	@echo "  make cleanup-kv-all           Delete ALL mining data (⚠️  DANGEROUS)"
	@echo ""
	@echo "Documentation (MkDocs):"
	@echo "  make docs-install     Install MkDocs and dependencies"
	@echo "  make docs-serve       Start local documentation server (http://localhost:8000)"
	@echo "  make docs-build       Build static documentation site"
	@echo "  make docs-clean       Clean documentation build artifacts"
	@echo "  make docs-deploy      Deploy documentation to GitHub Pages"
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

contracts-bindings: abi-token abi-escrow abi-vault abi-distributor abi-referral abi-cashback abi-mining abi-usdt generate-project-abis
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
	@test -f $(ESCROW_ARTIFACT) || (echo "❌ Escrow artifact not found! Run: make contracts-compile" && exit 1)
	@test -f $(TOKEN_ARTIFACT) || (echo "❌ Token artifact not found! Run: make contracts-compile" && exit 1)
	@echo "✅ Contract artifacts found"
	@echo ""
	@echo "Checking if bindings generated..."
	@test -f $(ABIS_DIR)/escrow/escrow.go || (echo "❌ Escrow bindings not found! Run: make contracts-bindings" && exit 1)
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

settle-all:
	@make reward-settlement-all

# ==============================================================================
# Mining Rewards Testing (Legacy - use reward-settlement instead)
# ==============================================================================
test-mining-rewards:
	@echo "🧪 Running mining rewards tests..."
	@bash scripts/test-mining-rewards.sh

test-inject-mining-data:
	@echo "💉 Injecting test mining reward data..."
	@go run cmd/dev/test-inject-mining-data/main.go

test-mining-settlement:
	@echo "🌳 Testing full settlement flow..."
	@make test-inject-mining-data
	@echo ""
	@echo "📊 Generating settlement..."
	@go run cmd/mining-settlement/main.go generate --type kawai

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

docs-install:
	@echo "📚 Installing MkDocs and dependencies..."
	@pip install mkdocs-material
	@pip install pymdown-extensions
	@echo "✅ MkDocs installed successfully!"

docs-serve:
	@echo "📖 Starting MkDocs development server..."
	@echo "🌐 Open http://localhost:8000 in your browser"
	@mkdocs serve

docs-build:
	@echo "🔨 Building static documentation site..."
	@mkdocs build
	@echo "✅ Documentation built to ./site/"

docs-clean:
	@echo "🧹 Cleaning documentation build artifacts..."
	@rm -rf site/
	@echo "✅ Documentation cleaned!"

docs-deploy:
	@echo "🚀 Deploying documentation to GitHub Pages..."
	@mkdocs gh-deploy --force
	@echo "✅ Documentation deployed!"