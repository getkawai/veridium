.PHONY: help db-generate bindings-generate dev build clean test

# Default target
help:
	@echo "Veridium Development Commands"
	@echo ""
	@echo "Database:"
	@echo "  make db-generate          Generate Go code from SQL queries (sqlc)"
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

# Generate Go code from SQL queries using sqlc
db-generate:
	@echo "🔄 Generating Go code from SQL queries..."
	sqlc generate
	@echo "✅ Database code generated!"

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
	rm -f backend-dev.log
	wails3 dev 2>&1 | tee backend-dev.log

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
