.PHONY: help db-generate bindings-generate dev build clean test update-llama db-dump db-restore

# Default target
help:
	@echo "Veridium Development Commands"
	@echo ""
	@echo "Database:"
	@echo "  make db-generate          Generate Go code from SQL queries (sqlc)"
	@echo "  make db-seed              Fetch and seed agents data (runs cmd/seed-agents)"
	@echo "  make db-dump              Dump SQLite database to seed file"
	@echo "  make db-restore           Restore database from seed file"
	@echo "  make bindings-generate    Generate TypeScript bindings from Go (wails)"
	@echo "  make generate             Run both db-generate and bindings-generate"
	@echo ""
	@echo "Development:"
	@echo "  make dev                  Start development server"
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

# Fetch and seed agents data detailed
db-seed:
	@echo "🌱 Seeding agents..."
	@rm -f data/seed_agents.db
	@go run cmd/seed-agents/main.go
	@echo "✅ Agents seeded!"

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
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev 2>&1 | tee backend-dev.log

# Start development server without removing database
devd:
	@echo "🚀 Starting development server..."
	killport 9245
	rm -f backend-dev.log
	VERIDIUM_DEV=true wails3 dev 2>&1 | tee backend-dev.log

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
