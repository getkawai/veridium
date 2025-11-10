#!/bin/bash

# Test script untuk llama-server chat completion endpoint
# Usage: ./test_llama_chat_completion.sh [port]
# Default port: 8080

PORT=${1:-8080}
BASE_URL="http://127.0.0.1:${PORT}"

echo "=========================================="
echo "Testing llama-server Chat Completion API"
echo "=========================================="
echo "Server URL: ${BASE_URL}"
echo ""

# Check if server is running
echo "1. Checking server health..."
HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/health" 2>/dev/null)

if [ "$HEALTH_RESPONSE" != "200" ]; then
    echo "❌ Server tidak berjalan di port ${PORT}"
    echo ""
    echo "Untuk start server, gunakan salah satu cara berikut:"
    echo ""
    echo "A. Menggunakan llama-server langsung:"
    echo "   llama-server -m ~/.veridium/models/qwen2.5-0.5b-instruct-q4_k_m.gguf --port ${PORT}"
    echo ""
    echo "B. Menggunakan aplikasi Veridium (akan auto-start jika ada model)"
    echo ""
    exit 1
fi

echo "✅ Server berjalan"
echo ""

# Test 1: Simple chat completion
echo "2. Test 1: Simple chat completion (non-streaming)"
echo "----------------------------------------"
curl -X POST "${BASE_URL}/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "kawai-auto",
    "messages": [
      {
        "role": "user",
        "content": "Hello! Can you tell me a short joke?"
      }
    ],
    "temperature": 0.7,
    "max_tokens": 100,
    "stream": false
  }' \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -s | jq '.' 2>/dev/null || cat

echo ""
echo ""

# Test 2: Streaming chat completion
echo "3. Test 2: Streaming chat completion"
echo "----------------------------------------"
echo "Request sent, streaming response:"
curl -X POST "${BASE_URL}/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "kawai-auto",
    "messages": [
      {
        "role": "user",
        "content": "Write a haiku about programming"
      }
    ],
    "temperature": 0.8,
    "max_tokens": 150,
    "stream": true
  }' \
  -s --no-buffer 2>&1 | head -20

echo ""
echo ""

# Test 3: Multi-turn conversation
echo "4. Test 3: Multi-turn conversation"
echo "----------------------------------------"
curl -X POST "${BASE_URL}/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "kawai-auto",
    "messages": [
      {
        "role": "system",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "What is 2+2?"
      },
      {
        "role": "assistant",
        "content": "2+2 equals 4."
      },
      {
        "role": "user",
        "content": "What about 3+3?"
      }
    ],
    "temperature": 0.7,
    "max_tokens": 50,
    "stream": false
  }' \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -s | jq '.choices[0].message.content' 2>/dev/null || cat

echo ""
echo ""
echo "=========================================="
echo "Test selesai!"
echo "=========================================="

