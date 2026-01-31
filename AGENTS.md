# Veridium Developer Guide for Agentic Coding

## Build, Lint, and Test Commands

### Building
```bash
make dev           # Full build (reset DB)
make dev-fast      # Fast build (reset DB, no frontend)
make dev-hot       # Hot reload (keep DB, no build)
make build         # Production build
```

### Testing
```bash
make test              # Run all Go tests
go test ./...          # Alternative to make test
go test -v ./internal/services -run TestCreateWallet  # Specific test
go test -v ./internal/services/wallet_service_test.go    # Single file
go test -cover ./...    # With coverage
```

### Code Generation
```bash
make generate         # Generate all (db + bindings + constants)
make db-generate     # SQL to Go via sqlc
make bindings-generate # TypeScript bindings via wails
make constants-generate # Constants from .env
```

### Smart Contracts
```bash
make contracts-compile    # Foundry compilation
make contracts-bindings   # Go bindings from ABIs
make contracts-test       # Contract tests
make contracts-upgrade    # Full workflow (test + compile + bindings)
```

### Database
```bash
make db-dump     # Dump to seed file
make db-restore   # Restore from seed file
```

### Clean
```bash
make clean              # Clean generated files
make contracts-clean     # Clean contract artifacts
```

## Code Style Guidelines

### Imports
- No enforced order, but group: stdlib, third-party, internal
- Remove unused imports
- Use `goimports` to organize

### Formatting
- Use `gofmt` (standard Go formatter)
- No explicit .golangci.yml at root

### Types
- Use `context.Context` for async operations
- Prefer interfaces (see `pkg/store/kvstore.go`)
- Structs with exported fields for API responses
- Pointers for optional struct fields

### Naming Conventions

#### Packages
- Lowercase, single word: `services`, `store`, `blockchain`
- Avoid underscores in package names

#### Types
- Exported: CamelCase, e.g., `ReferralService`, `WalletService`
- Unexported: camelCase, e.g., `agentSession`

#### Functions/Methods
- Exported: PascalCase, e.g., `CreateReferralCode`, `GetWallets`
- Unexported: camelCase, e.g., `formatTimestamp`, `setupForList`
- Constructors: `New<Type>`, e.g., `NewReferralService`

#### Constants
- Exported: PascalCase, e.g., `UIMessageRoleUser`
- Unexported: camelCase, e.g., `totalReferrals`

#### Variables
- CamelCase: `userAddress`, `referralCode`
- Loop vars: `i`, `k`, `v`

#### JSON Tags
- snake_case: `total_earnings_usdt`
- Use `omitempty` for optional fields
- Match frontend contracts

#### Database Fields
- snake_case columns: `total_earnings_usdt`
- Named params via sqlc

### Error Handling
- Always check errors
- Use `fmt.Errorf` with `%w`: `fmt.Errorf("failed: %w", err)`
- Return errors as last return value
- Use package-level errors for known types
- Log errors with context

### Comments
- Exported functions must have comments
- Use `//` for single-line, block comments for packages
- Explain "why", not "what"

### Struct Organization
- Group related fields
- JSON tags on same line as field
- Pointers for nullables: `FirstDepositAt *time.Time`
- Order: basic types, pointers, nested structs

### Testing
- Use `github.com/stretchr/testify`: `require.NoError`, `assert.Equal`
- Table-driven tests for multiple cases
- `t.Helper()` in helpers, `t.TempDir()` for temp dirs, `t.Cleanup()` for cleanup
- Test files end with `_test.go`
- Descriptive names: `TestCreateWallet_WithValidMnemonic`

### Concurrency
- Use `sync.RWMutex` for shared state
- Always lock before read/write
- Defer unlock: `defer s.mu.RUnlock()`

### Context Usage
- First param: `ctx context.Context`
- Pass through call chains
- Use for cancellation/timeouts
- Background context: `context.Background()`

### Logging
- Use `log` or `slog` for structured logging
- Include context in messages

### File Organization
- `main.go` at root
- `internal/` for private app code
- `pkg/` for reusable libraries
- `cmd/` for CLI tools
- `contracts/` for Solidity contracts

### Path Conventions
- All application data storage paths must be resolved through `internal/paths`.
- Do not use hardcoded paths or `os.UserHomeDir()` directly for app data.
- Use `paths.Base()`, `paths.Database()`, etc., to ensure cross-platform compatibility and proper development/packaged behavior.

### Smart Contracts
- Bindings in `internal/generate/abi/`
- Handle transaction errors gracefully
- Verify addresses and ABIs
- Use `bind` package for calls

### Database Queries
- Queries in `internal/database/queries/`
- Generated code in `internal/database/generated/`
- Named queries: `-- name: QueryName :one/:many/:exec`
- Transactions via `*sql.DB` or `*sqlx.DB`

### Wails Integration
- Expose via `application.NewService()`
- Frontend calls via RPC
- Return structured types with JSON tags
- Handle errors as return value

### Constants & Config
- Env vars in `.env`
- Generate constants with `make constants-generate`
- Use obfuscated constants for secrets
- Never commit `.env`
