#!/usr/bin/env fish

# Test Download with Retry Logic for GitHub CDN Delays
# This tests the new retry mechanism for b7018

echo "🧪 Testing Download with Retry Logic"
echo "===================================="
echo ""

# Clean up any previous downloads
echo "🧹 Cleaning up previous downloads..."
rm -rf ~/.llama-cpp/bin ~/.llama-cpp/metadata

echo ""
echo "📥 Attempting to download llama.cpp b7018..."
echo "   (This will retry 3 times with 2s, 4s, 8s backoff if 404)"
echo ""

# Run the installer test
go run ./cmd/test-library-paths/main.go

echo ""
echo "✅ Test complete!"

