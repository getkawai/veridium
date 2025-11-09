# Llama.cpp Embedding Quick Start Guide

## TL;DR

```bash
# 1. Download embedding model (one-time, ~550MB)
mkdir -p ~/models
cd ~/models
wget https://huggingface.co/nomic-ai/nomic-embed-text-v1.5-GGUF/resolve/main/nomic-embed-text-v1.5.Q4_K_M.gguf

# 2. Start llama-server with embedding model
llama-server \
  -m ~/models/nomic-embed-text-v1.5.Q4_K_M.gguf \
  --port 8080 \
  --embedding \
  --pooling mean \
  --embd-normalize 2

# 3. Verify it's working
curl -X POST http://localhost:8080/embedding \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello world"}' | jq '.embedding | length'
# Should output: 768

# 4. Run Veridium
# The app will automatically use llama.cpp for embeddings!
```

## Step-by-Step Guide

### 1. Check if llama-server is Installed

```bash
which llama-server
# If not found, install llama.cpp:
brew install llama.cpp
```

### 2. Download an Embedding Model

Choose one based on your needs:

#### Option A: nomic-embed-text-v1.5 (Recommended)
- **Best for**: Multilingual, general purpose
- **Size**: ~550MB
- **Dimensions**: 768

```bash
mkdir -p ~/models
cd ~/models
wget https://huggingface.co/nomic-ai/nomic-embed-text-v1.5-GGUF/resolve/main/nomic-embed-text-v1.5.Q4_K_M.gguf
```

#### Option B: all-MiniLM-L6-v2 (Fastest)
- **Best for**: English only, speed
- **Size**: ~90MB
- **Dimensions**: 384

```bash
cd ~/models
wget https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2-GGUF/resolve/main/all-MiniLM-L6-v2.Q4_K_M.gguf
```

### 3. Start llama-server

```bash
llama-server \
  -m ~/models/nomic-embed-text-v1.5.Q4_K_M.gguf \
  --port 8080 \
  --embedding \
  --pooling mean \
  --embd-normalize 2 \
  --ctx-size 2048 \
  --batch-size 512 \
  --threads 4
```

**Important flags:**
- `--embedding`: Enable embedding mode
- `--pooling mean`: Use mean pooling (recommended)
- `--embd-normalize 2`: L2 normalization for cosine similarity
- `--ctx-size 2048`: Context window size
- `--batch-size 512`: Batch size for processing

### 4. Test the Embedding Endpoint

```bash
# Test health
curl http://localhost:8080/health

# Test embedding generation
curl -X POST http://localhost:8080/embedding \
  -H "Content-Type: application/json" \
  -d '{"content":"The quick brown fox jumps over the lazy dog"}' \
  | jq '.embedding | length'

# Should output: 768 (for nomic-embed-text)
```

### 5. Run Veridium

The application is already configured to use llama.cpp by default!

```bash
# Just run the app
./veridium

# Or if building from source
go run .
```

You should see in the logs:

```
✅ Vector Search service initialized (chromem)
   Database path: /Users/you/Library/Application Support/veridium/vector-db
   Embedding provider: llama.cpp (llama-server)
   Embedding endpoint: http://localhost:8080
```

### 6. Upload a File and Test RAG

1. Open Veridium
2. Upload a document (PDF, DOCX, TXT, etc.)
3. The file will be automatically:
   - Parsed
   - Chunked
   - Embedded using llama.cpp
   - Stored in chromem vector database
4. Ask questions about the document!

## Troubleshooting

### Problem: "connection refused" error

**Solution**: Make sure llama-server is running on port 8080

```bash
# Check if llama-server is running
lsof -i :8080

# If not, start it
llama-server -m ~/models/nomic-embed-text-v1.5.Q4_K_M.gguf --port 8080 --embedding
```

### Problem: "no embeddings found in the response"

**Solution**: You forgot the `--embedding` flag

```bash
# Wrong (LLM mode)
llama-server -m model.gguf --port 8080

# Correct (embedding mode)
llama-server -m model.gguf --port 8080 --embedding
```

### Problem: Port 8080 already in use

**Solution**: Use a different port

```bash
# Start llama-server on port 8081
llama-server -m model.gguf --port 8081 --embedding

# Update main.go
vectorSearchService, err := services.NewVectorSearchService(
    vectorDBPath, 
    "llama", 
    "http://localhost:8081", // Changed port
)
```

### Problem: Model file not found

**Solution**: Check the path and download the model

```bash
# Check if model exists
ls -lh ~/models/nomic-embed-text-v1.5.Q4_K_M.gguf

# If not, download it
cd ~/models
wget https://huggingface.co/nomic-ai/nomic-embed-text-v1.5-GGUF/resolve/main/nomic-embed-text-v1.5.Q4_K_M.gguf
```

## Advanced: Run as Background Service

### macOS (launchd)

Create `~/Library/LaunchAgents/com.veridium.llama-embedding.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.veridium.llama-embedding</string>
    <key>ProgramArguments</key>
    <array>
        <string>/opt/homebrew/bin/llama-server</string>
        <string>-m</string>
        <string>/Users/YOUR_USERNAME/models/nomic-embed-text-v1.5.Q4_K_M.gguf</string>
        <string>--port</string>
        <string>8080</string>
        <string>--embedding</string>
        <string>--pooling</string>
        <string>mean</string>
        <string>--embd-normalize</string>
        <string>2</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/llama-embedding.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/llama-embedding.error.log</string>
</dict>
</plist>
```

Load the service:

```bash
launchctl load ~/Library/LaunchAgents/com.veridium.llama-embedding.plist
launchctl start com.veridium.llama-embedding
```

### Linux (systemd)

Create `/etc/systemd/system/llama-embedding.service`:

```ini
[Unit]
Description=Llama.cpp Embedding Server
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
ExecStart=/usr/local/bin/llama-server \
  -m /home/YOUR_USERNAME/models/nomic-embed-text-v1.5.Q4_K_M.gguf \
  --port 8080 \
  --embedding \
  --pooling mean \
  --embd-normalize 2
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable llama-embedding
sudo systemctl start llama-embedding
sudo systemctl status llama-embedding
```

## Performance Tips

### 1. Use GPU Acceleration (if available)

```bash
# Metal (macOS)
llama-server -m model.gguf --embedding --n-gpu-layers 99

# CUDA (Linux/Windows with NVIDIA GPU)
llama-server -m model.gguf --embedding --n-gpu-layers 99
```

### 2. Optimize Batch Size

```bash
# For M1/M2 Macs (16GB RAM)
llama-server -m model.gguf --embedding --batch-size 1024 --ubatch-size 512

# For M1/M2 Macs (8GB RAM)
llama-server -m model.gguf --embedding --batch-size 512 --ubatch-size 256
```

### 3. Adjust Thread Count

```bash
# Use all CPU cores
llama-server -m model.gguf --embedding --threads $(sysctl -n hw.ncpu)

# Or manually set (e.g., 4 threads)
llama-server -m model.gguf --embedding --threads 4
```

## Comparison: Ollama vs llama.cpp

Both are good choices, but here's why we chose llama.cpp as default:

| Feature | llama.cpp | Ollama |
|---------|-----------|--------|
| **Speed** | ✅ Faster (direct binary) | ⚠️ Slower (Docker overhead) |
| **Memory** | ✅ Lower (no Docker) | ⚠️ Higher (Docker + app) |
| **Setup** | ⚠️ Manual model download | ✅ Auto model download |
| **Integration** | ✅ Already used for LLM | ❌ Separate service |
| **Control** | ✅ Full control over flags | ⚠️ Limited configuration |

**Recommendation**: Use llama.cpp for production, Ollama for quick testing.

## Next Steps

1. **Try different models**: Experiment with different embedding models
2. **Monitor performance**: Check embedding generation speed
3. **Tune parameters**: Adjust batch size, threads, etc.
4. **Auto-start**: Set up as background service
5. **Scale up**: Run multiple instances for load balancing

---

**Need help?** Check the full documentation: [LLAMA_EMBEDDING_INTEGRATION.md](./LLAMA_EMBEDDING_INTEGRATION.md)

