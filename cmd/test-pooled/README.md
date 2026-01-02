# Test Pooled Provider

CLI tool untuk testing pooled provider dengan CLIProxyAPI fallback mechanism.

## 🚀 Usage

```bash
# Build
go build -o test-pooled ./cmd/test-pooled/

# Run
./test-pooled
```

## 📊 What It Tests

### Test 1: OpenRouter Pooled
- Creates pooled provider with multiple OpenRouter API keys
- Tests with prompt: "What is recent news on AI?"
- Shows account status after request

### Test 2: ZAI Pooled
- Creates pooled provider with multiple ZAI API keys
- Tests with same prompt
- Shows account status after request

### Test 3: Full Chain with Fallback
- Creates chain: OpenRouter → Pollinations → ZAI
- Tests 3-level fallback mechanism
- Demonstrates circuit breaker in action

## 📈 Expected Output

```
🚀 Testing Pooled Provider with CLIProxyAPI Fallback
============================================================

📊 Test 1: OpenRouter with Multiple API Keys
------------------------------------------------------------
🔑 Found 2 OpenRouter API keys
PooledProvider[openrouter]: Registered 2 API keys
✅ Pooled provider created

🧪 Testing OpenRouter Pooled...
📝 Prompt: What is recent news on AI? Give me a brief summary in 2-3 sentences.

✅ Response received in 2.5s

📄 Response:
---
Recent AI news includes...
---

📊 Usage: 25 input tokens, 50 output tokens, 75 total

📈 Account Status:
  [1] openrouter-key-1
      Provider: openrouter
      Status: active
  [2] openrouter-key-2
      Provider: openrouter
      Status: active
```

## 🎯 Features Demonstrated

### ✅ Account Pooling
- Multiple API keys per provider
- Round-robin selection
- Load balancing

### ✅ Fallback Mechanism
- Level 1: Account failover
- Level 2: Provider rotation
- Level 3: Global retry

### ✅ Circuit Breaker
- Skips failed models
- Prevents repeated failures
- Automatic recovery

### ✅ State Management
- Tracks account status
- Monitors quota
- Shows cooldown state

## 🐛 Troubleshooting

### All tests fail
**Check:**
- API keys are valid
- Network connectivity
- Rate limits

### Timeout errors
**Solution:**
- Increase timeout in code
- Check network latency
- Verify API endpoints

### Quota exceeded
**Solution:**
- Wait for cooldown
- Add more API keys
- Check account status

## 📚 Related Documentation

- `docs/INTEGRATION_GUIDE.md` - How to use pooled providers
- `docs/IMPLEMENTATION_COMPLETE.md` - What was implemented
- `examples/pooled_provider_example.go` - More examples

## ✨ Status

✅ **Working** - Demonstrates full fallback mechanism
⚠️ **Note**: Some providers may have conversion issues (being fixed)


