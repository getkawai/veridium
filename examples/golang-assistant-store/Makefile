.PHONY: run build test clean fmt vet lint help

# Variables
BINARY_NAME=assistant-store
GO=go

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	$(GO) run .

run-advanced: ## Run advanced examples
	$(GO) run . -advanced

build: ## Build the application
	$(GO) build -o $(BINARY_NAME) .

test: ## Run tests
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

benchmark: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

fmt: ## Format code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

lint: ## Run golangci-lint (requires golangci-lint installed)
	golangci-lint run

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	$(GO) clean

deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

install: ## Install the application
	$(GO) install

# Development helpers
dev: fmt vet test ## Run fmt, vet, and test

all: clean fmt vet test build ## Run all checks and build

