# Quick Start - Bulk Tunnels with Obfuscated Tokens

## 🚀 Quick Commands

### Create Tunnels (with obfuscated tokens)
```bash
cd /Users/yuda/github.com/kawai-network/veridium
go run cmd/bulk-tunnels/main.go -count 5
```

### View Tunnels (without decoding)
```bash
go run cmd/decode-tunnel-token/main.go -file tunnels.json
```

### Decode All Tokens
```bash
go run cmd/decode-tunnel-token/main.go -file tunnels.json -all
```

### Decode Specific Token
```bash
go run cmd/decode-tunnel-token/main.go -file tunnels.json -id node-1.getkawai.com
```

### See Demo
```bash
go run cmd/bulk-tunnels/test/main.go
```

## 📝 What Changed?

### Before (Plaintext Token)
```json
{
  "TunnelID": "abc123...",
  "TunnelToken": "eyJhIjoiY2VhYjIxODc1MWQzM2NkODA0ODc4MTk2YWQ3YmVmNzQi...",
  "Hostname": "node-1.getkawai.com"
}
```

### After (Obfuscated Token)
```json
{
  "TunnelID": "abc123...",
  "TunnelToken": "xD78MD8OXAVZ0Iuh/UFivQU+A7GKtrJvwwVvHwoW89hepFQqia...",
  "Hostname": "node-1.getkawai.com"
}
```

## 🔐 How to Use Obfuscated Tokens

### In Your Code

```go
import (
    "github.com/kawai-network/veridium/pkg/obfuscator"
    "github.com/kawai-network/veridium/pkg/tunnelkit"
)

// Read tunnels.json
var tunnels []*tunnelkit.TunnelInfo
// ... load from JSON ...

// Decode token
obf := obfuscator.New()
decodedToken, err := obf.Decode(tunnels[0].TunnelToken)
if err != nil {
    log.Fatal(err)
}

// Use decoded token
err = tunnelkit.RunTunnel(ctx, decodedToken)
```

### From Command Line

```bash
# Get decoded token for specific tunnel
TOKEN=$(go run cmd/decode-tunnel-token/main.go -file tunnels.json -id node-1.getkawai.com | grep "Decoded Token:" | cut -d: -f2 | xargs)

# Use the token (example)
cloudflared tunnel run --token "$TOKEN"
```

## 🎯 Common Use Cases

### 1. Create Production Tunnels
```bash
go run cmd/bulk-tunnels/main.go \
  -domain production.example.com \
  -count 10 \
  -out production-tunnels.json
```

### 2. Create Development Tunnels
```bash
go run cmd/bulk-tunnels/main.go \
  -domain dev.example.com \
  -count 3 \
  -out dev-tunnels.json
```

### 3. Audit Existing Tunnels
```bash
# List all
go run cmd/decode-tunnel-token/main.go -file tunnels.json

# Decode specific one
go run cmd/decode-tunnel-token/main.go -file tunnels.json -id <tunnel-id>
```

### 4. Rotate Tokens
```bash
# 1. Delete old tunnels (via Cloudflare dashboard or API)
# 2. Create new ones
go run cmd/bulk-tunnels/main.go -count 5 -out new-tunnels.json
# 3. Tokens are automatically obfuscated
```

## 🛡️ Security Best Practices

### ✅ DO:
- Keep `tunnels.json` in `.gitignore`
- Set file permissions: `chmod 600 tunnels.json`
- Rotate tokens regularly
- Use decode tool only when needed
- Keep obfuscator code private

### ❌ DON'T:
- Commit `tunnels.json` to public repos
- Share obfuscated tokens publicly (still decodable)
- Rely on obfuscation for security compliance
- Use for highly sensitive production secrets

## 🔧 Troubleshooting

### Token Decode Fails
```bash
# Check if file is valid JSON
cat tunnels.json | jq .

# Try with verbose output
go run cmd/decode-tunnel-token/main.go -file tunnels.json -all
```

### Tunnel Not Found
```bash
# List all tunnels to see available IDs
go run cmd/decode-tunnel-token/main.go -file tunnels.json
```

### Build Issues
```bash
# Ensure dependencies are installed
go mod tidy

# Build manually
go build -o bulk-tunnels cmd/bulk-tunnels/main.go cmd/bulk-tunnels/decode-token.go
go build -o decode-tunnel-token cmd/decode-tunnel-token/main.go
```

## 📚 More Information

- Full documentation: [README.md](README.md)
- Changelog: [CHANGELOG.md](CHANGELOG.md)
- Obfuscator docs: [../../pkg/obfuscator/README.md](../../pkg/obfuscator/README.md)

## 💡 Tips

1. **Test First**: Run the demo to understand obfuscation
   ```bash
   go run cmd/bulk-tunnels/test/main.go
   ```

2. **Backup**: Keep backup of `tunnels.json` before rotating
   ```bash
   cp tunnels.json tunnels.json.backup
   ```

3. **Automation**: Use in scripts
   ```bash
   #!/bin/bash
   go run cmd/bulk-tunnels/main.go -count 5
   echo "Tunnels created and tokens obfuscated!"
   ```

4. **CI/CD**: Decode tokens in deployment
   ```bash
   TOKEN=$(go run cmd/decode-tunnel-token/main.go -file tunnels.json -id $TUNNEL_ID)
   # Use $TOKEN in deployment
   ```

## ⚠️ Important Reminders

- This is **obfuscation**, not **encryption**
- Anyone with the code can decode tokens
- Use proper encryption for production secrets
- Obfuscation adds obscurity, not security
- Perfect for hiding from casual viewing
- Not suitable for compliance requirements

