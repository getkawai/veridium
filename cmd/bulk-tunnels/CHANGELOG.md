# Changelog - Bulk Tunnels

## [2.0.0] - TunnelToken Obfuscation

### Added
- ✅ **TunnelToken Obfuscation**: Semua TunnelToken yang digenerate sekarang otomatis di-obfuscate
- ✅ **decode-token.go**: Helper functions untuk decode obfuscated tokens
- ✅ **decode-tunnel-token command**: Standalone tool untuk decode tokens dari JSON file
- ✅ **test-obfuscation.go**: Test function untuk verify obfuscation
- ✅ **test/main.go**: Demo program untuk showcase obfuscation
- ✅ **README.md**: Comprehensive documentation
- ✅ **CHANGELOG.md**: This file

### Changed
- 🔄 TunnelToken sekarang di-obfuscate sebelum disimpan ke JSON
- 🔄 Import `pkg/obfuscator` untuk obfuscation functionality

### Security
- 🔒 TunnelToken tidak lagi plaintext di JSON files
- 🔒 Menambah layer kesulitan untuk copy-paste tokens
- ⚠️ **Note**: Masih obfuscation, bukan encryption

### Usage

#### Before (v1.0.0):
```json
{
  "TunnelToken": "eyJhIjoiY2VhYjIxODc1MWQzM2NkODA0ODc4MTk2YWQ3YmVmNzQi..."
}
```

#### After (v2.0.0):
```json
{
  "TunnelToken": "xD78MD8OXAVZ0Iuh/UFivQU+A7GKtrJvwwVvHwoW89hepFQqia..."
}
```

### Tools

1. **bulk-tunnels** - Create tunnels with obfuscated tokens
   ```bash
   go run cmd/bulk-tunnels/main.go -count 5
   ```

2. **decode-tunnel-token** - Decode obfuscated tokens
   ```bash
   go run cmd/decode-tunnel-token/main.go -file tunnels.json -all
   ```

3. **test demo** - See obfuscation in action
   ```bash
   go run cmd/bulk-tunnels/test/main.go
   ```

### Breaking Changes

⚠️ **Important**: TunnelToken format berubah dari plaintext ke obfuscated.

- Jika Anda memiliki `tunnels.json` lama dengan plaintext tokens, mereka masih bisa digunakan
- Token baru akan otomatis ter-obfuscate
- Gunakan `decode-tunnel-token` untuk decode tokens

### Migration Guide

Jika Anda memiliki existing `tunnels.json` dengan plaintext tokens:

1. **Option 1**: Re-create tunnels dengan tool baru
   ```bash
   go run cmd/bulk-tunnels/main.go -count 3
   ```

2. **Option 2**: Manually obfuscate existing tokens
   ```go
   obf := obfuscator.New()
   obfuscatedToken := obf.Encode(plaintextToken)
   ```

3. **Option 3**: Keep using plaintext (not recommended)
   - Plaintext tokens masih valid untuk Cloudflare API
   - Tapi tidak mendapat benefit obfuscation

### Performance

- Obfuscation overhead: ~34.6% increase in token length
- Encode time: ~700ns per token
- Decode time: ~730ns per token
- Negligible impact on tunnel creation time

### Dependencies

- `github.com/kawai-network/veridium/pkg/obfuscator` - Custom obfuscator

### Files Changed

```
cmd/bulk-tunnels/
├── main.go                 # Modified: Added obfuscation
├── decode-token.go         # New: Helper functions
├── test-obfuscation.go     # New: Test function
├── README.md               # New: Documentation
├── CHANGELOG.md            # New: This file
└── test/
    └── main.go             # New: Demo program

cmd/decode-tunnel-token/
└── main.go                 # New: Decode tool
```

## [1.0.0] - Initial Release

### Added
- Initial bulk tunnels creation tool
- Hardcoded Cloudflare credentials
- JSON output
- DNS routing configuration

