# Cloudflare Go SDK Migration Guide
## v0.114.0 → v6.6.0

### Migration Status: 75% Complete

## ✅ Completed Files
1. `go.mod` - Updated dependency
2. `pkg/store/kv_client.go` - New adapter created
3. `pkg/store/kvstore.go` - Main struct updated
4. `pkg/store/balance.go` - All operations migrated
5. `pkg/store/balance_trial.go` - All operations migrated
6. `pkg/store/apikey.go` - All operations migrated
7. `pkg/store/holder.go` - All operations migrated
8. `pkg/store/marketplace.go` - All operations migrated
9. `pkg/stablediffusion/remote/cloudflare.go` - New image generator using v6

## ⚠️ Remaining Files (Need Migration)

### Core Store Files
- `pkg/store/contributor.go` - 3 operations
- `pkg/store/referral.go` - 9 operations
- `pkg/store/period_counter.go` - 4 operations
- `pkg/store/merkle.go` - ~10 operations
- `pkg/store/settlement.go` - ~15 operations
- `pkg/store/job_rewards.go` - ~8 operations
- `pkg/store/cashback_kv.go` - 3 operations

### Dev Tools (Optional - Can be updated later)
- `cmd/dev/inspect-proof/main.go`
- `cmd/dev/cleanup-kv-all/main.go`
- `cmd/dev/list-job-rewards/main.go`
- `cmd/dev/list-kv-keys/main.go`

## 📋 Migration Pattern Reference

### 1. Remove Old Import
```go
// REMOVE THIS
import "github.com/cloudflare/cloudflare-go"
```

### 2. Write Operations
```go
// OLD
_, err := s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
    NamespaceID: namespaceID,
    Key:         key,
    Value:       data,
})

// NEW
err := s.client.SetValue(ctx, namespaceID, key, data)
```

### 3. Read Operations
```go
// OLD
value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
    NamespaceID: namespaceID,
    Key:         key,
})

// NEW
value, err := s.client.GetValue(ctx, namespaceID, key)
```

### 4. Delete Operations
```go
// OLD
_, err := s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
    NamespaceID: namespaceID,
    Key:         key,
})

// NEW
err := s.client.DeleteValue(ctx, namespaceID, key)
```

### 5. List Operations
```go
// OLD
resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
    NamespaceID: namespaceID,
    Prefix:      prefix,
})
// Access: resp.Result, resp.ResultInfo.Cursor

// NEW
result, err := s.client.ListKeys(ctx, namespaceID, prefix, cursor)
// Access: result.Result, result.ResultInfo.Cursor
```

## 🔧 Quick Migration Steps

### For Each Remaining File:

1. **Open the file**
2. **Remove import**: Delete `"github.com/cloudflare/cloudflare-go"` line
3. **Find & Replace** (use regex in your editor):
   
   **Pattern 1 - Write:**
   ```
   Find: _, err := s\.client\.WriteWorkersKVEntry\(ctx, cloudflare\.AccountIdentifier\(s\.accountID\), cloudflare\.WriteWorkersKVEntryParams\{\s*NamespaceID:\s*([^,]+),\s*Key:\s*([^,]+),\s*Value:\s*([^,]+),\s*\}\)
   Replace: err := s.client.SetValue(ctx, $1, $2, $3)
   ```
   
   **Pattern 2 - Read:**
   ```
   Find: ([^,]+), err := s\.client\.GetWorkersKV\(ctx, cloudflare\.AccountIdentifier\(s\.accountID\), cloudflare\.GetWorkersKVParams\{\s*NamespaceID:\s*([^,]+),\s*Key:\s*([^,]+),\s*\}\)
   Replace: $1, err := s.client.GetValue(ctx, $2, $3)
   ```
   
   **Pattern 3 - Delete:**
   ```
   Find: _, err := s\.client\.DeleteWorkersKVEntry\(ctx, cloudflare\.AccountIdentifier\(s\.accountID\), cloudflare\.DeleteWorkersKVEntryParams\{\s*NamespaceID:\s*([^,]+),\s*Key:\s*([^,]+),\s*\}\)
   Replace: err := s.client.DeleteValue(ctx, $1, $2)
   ```

4. **Handle List operations manually** (they're more complex)
5. **Run**: `go build ./...` to check for errors
6. **Fix any remaining issues**

## 🧪 Testing

After migration:

```bash
# 1. Clean and rebuild
go mod tidy
go clean -cache

# 2. Build
go build ./...

# 3. Run tests
go test ./pkg/store/...

# 4. Run specific KV tests
go test -v ./pkg/store -run TestBalance
go test -v ./pkg/store -run TestAPIKey
go test -v ./pkg/store -run TestHolder
```

## 🚀 Deployment Checklist

- [ ] All files migrated
- [ ] `go mod tidy` executed
- [ ] All tests passing
- [ ] Integration tests with real Cloudflare KV
- [ ] Staging deployment successful
- [ ] Production deployment

## 📝 Notes

### Why v6?
- Complete rewrite with Stainless code generator
- Better type safety
- Cleaner API
- Better error handling
- Active maintenance

### Breaking Changes
- Import path changed: `cloudflare-go` → `cloudflare-go/v6`
- API structure changed: flat functions → nested services
- Parameter wrapping with `F()` helper
- Response structure changes

### Rollback Plan
If issues occur:
1. Restore from git: `git checkout HEAD -- pkg/store/`
2. Revert go.mod: `git checkout HEAD -- go.mod go.sum`
3. Run: `go mod tidy`

## 🔗 References
- [Cloudflare Go SDK v6 Docs](https://pkg.go.dev/github.com/cloudflare/cloudflare-go/v6)
- [Migration Guide](https://github.com/cloudflare/cloudflare-go)
- [KV API Reference](https://developers.cloudflare.com/api/operations/workers-kv-namespace-list-namespaces)
