# Bulk Tunnels Creator

Tool untuk membuat multiple Cloudflare Tunnels sekaligus dengan TunnelToken yang ter-obfuscate.

## Features

- ✅ Membuat multiple tunnels sekaligus
- ✅ Automatic DNS routing configuration
- ✅ **TunnelToken di-obfuscate** untuk keamanan tambahan
- ✅ Output ke JSON file
- ✅ Hardcoded credentials (dapat di-override via flags)

## TunnelToken Obfuscation

TunnelToken yang dihasilkan akan **otomatis di-obfuscate** menggunakan custom obfuscator tanpa secret key. Ini memberikan layer keamanan tambahan saat menyimpan tokens di file JSON.

### Mengapa Obfuscate?

- ✅ Menyembunyikan token dari casual inspection
- ✅ Menambah kesulitan untuk copy-paste token
- ✅ Tidak memerlukan key management
- ✅ Mudah di-decode kembali saat dibutuhkan

⚠️ **Note**: Ini adalah obfuscation, bukan encryption. Jangan mengandalkan ini sebagai security utama.

## Usage

### Create Tunnels

```bash
# Default: membuat 3 tunnels untuk getkawai.com
go run cmd/bulk-tunnels/main.go

# Custom jumlah tunnels
go run cmd/bulk-tunnels/main.go -count 5

# Custom domain
go run cmd/bulk-tunnels/main.go -domain example.com -count 10

# Custom output file
go run cmd/bulk-tunnels/main.go -out data/my-tunnels.json
```

### Decode Tokens

Gunakan tool `decode-tunnel-token` untuk decode obfuscated tokens:

```bash
# List semua tunnels (tanpa decode)
go run cmd/decode-tunnel-token/main.go -file tunnels.json

# Decode semua tokens
go run cmd/decode-tunnel-token/main.go -file tunnels.json -all

# Decode token spesifik by Tunnel ID
go run cmd/decode-tunnel-token/main.go -file tunnels.json -id <tunnel-id>

# Decode token spesifik by Hostname
go run cmd/decode-tunnel-token/main.go -file tunnels.json -id node-1.getkawai.com
```

## Output Format

File `tunnels.json` akan berisi:

```json
[
  {
    "TunnelID": "abc123...",
    "TunnelToken": "Q1nPvSVXvy2Ke5J7jT==...",  // OBFUSCATED
    "Hostname": "node-1.getkawai.com",
    "PublicURL": "https://node-1.getkawai.com"
  },
  {
    "TunnelID": "def456...",
    "TunnelToken": "MUYFV20jYHDbfJ8=...",  // OBFUSCATED
    "Hostname": "node-2.getkawai.com",
    "PublicURL": "https://node-2.getkawai.com"
  }
]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-acc-id` | `ceab218751d33cd804878196ad7bef74` | Cloudflare Account ID |
| `-token` | `OP8BZQhyeJxrovCPKt15eUOSC6i5LXTVECGRSMc1` | Cloudflare API Token |
| `-domain` | `getkawai.com` | Base domain untuk hostnames |
| `-count` | `3` | Jumlah tunnels yang akan dibuat |
| `-out` | `tunnels.json` | Output file path |

## Examples

### Example 1: Create 5 Tunnels

```bash
go run cmd/bulk-tunnels/main.go -count 5
```

Output:
```
Starting creation of 5 tunnels for domain getkawai.com...
[1/5] Creating tunnel 'node-1' (Host: node-1.getkawai.com)... SUCCESS (ID: abc123...)
[2/5] Creating tunnel 'node-2' (Host: node-2.getkawai.com)... SUCCESS (ID: def456...)
[3/5] Creating tunnel 'node-3' (Host: node-3.getkawai.com)... SUCCESS (ID: ghi789...)
[4/5] Creating tunnel 'node-4' (Host: node-4.getkawai.com)... SUCCESS (ID: jkl012...)
[5/5] Creating tunnel 'node-5' (Host: node-5.getkawai.com)... SUCCESS (ID: mno345...)

Done! Successfully created 5/5 tunnels.
Results saved to tunnels.json
```

### Example 2: Decode All Tokens

```bash
go run cmd/decode-tunnel-token/main.go -file tunnels.json -all
```

Output:
```
Decoding 5 tunnel token(s)...

[1] Tunnel ID: abc123...
    Hostname: node-1.getkawai.com
    Public URL: https://node-1.getkawai.com
    Token (Decoded): eyJhIjoiY2VhYjIxODc1MWQzM2NkODA0ODc4MTk2YWQ3YmVmNzQi...

[2] Tunnel ID: def456...
    Hostname: node-2.getkawai.com
    Public URL: https://node-2.getkawai.com
    Token (Decoded): eyJhIjoiY2VhYjIxODc1MWQzM2NkODA0ODc4MTk2YWQ3YmVmNzQi...
```

### Example 3: Decode Specific Token

```bash
go run cmd/decode-tunnel-token/main.go -file tunnels.json -id node-1.getkawai.com
```

Output:
```
Tunnel ID: abc123...
Hostname: node-1.getkawai.com
Public URL: https://node-1.getkawai.com

Obfuscated Token:
Q1nPvSVXvy2Ke5J7jT==...

Decoded Token:
eyJhIjoiY2VhYjIxODc1MWQzM2NkODA0ODc4MTk2YWQ3YmVmNzQi...
```

## Programmatic Usage

Untuk menggunakan decode function di code Anda:

```go
import (
    "github.com/kawai-network/veridium/pkg/obfuscator"
    "github.com/kawai-network/veridium/pkg/tunnelkit"
)

// Decode obfuscated token
obf := obfuscator.New()
decodedToken, err := obf.Decode(obfuscatedToken)
if err != nil {
    log.Fatal(err)
}

// Use the decoded token
err = tunnelkit.RunTunnel(ctx, decodedToken)
```

## Security Considerations

### ✅ What Obfuscation Provides:

- Hides tokens from casual viewing
- Makes copy-paste more difficult
- Adds a layer of obscurity
- No key management needed

### ❌ What Obfuscation Does NOT Provide:

- **NOT encryption** - Anyone with the code can decode
- **NOT secure** for sensitive production data
- **NOT a replacement** for proper secrets management
- **NOT compliant** with security standards (use proper encryption for that)

### Best Practices:

1. **Don't commit** `tunnels.json` to public repositories
2. **Use .gitignore** to exclude tunnel files
3. **Rotate tokens** regularly
4. **Use proper encryption** for production secrets
5. **Restrict file permissions** (chmod 600)

## Related Tools

- `cmd/decode-tunnel-token` - Decode obfuscated tunnel tokens
- `pkg/obfuscator` - String obfuscation library
- `pkg/tunnelkit` - Cloudflare Tunnel management

## Troubleshooting

### Token decode fails

```bash
# Verify file exists and is valid JSON
cat tunnels.json | jq .

# Try decoding with verbose error
go run cmd/decode-tunnel-token/main.go -file tunnels.json -all
```

### Tunnel creation fails

Check:
1. Cloudflare credentials are correct
2. Domain exists in your Cloudflare account
3. API token has required permissions
4. Network connectivity to Cloudflare API

## License

Part of Veridium project.

