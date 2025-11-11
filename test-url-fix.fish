#!/usr/bin/env fish

echo "🧪 Testing URL Fix for Download"
echo "================================"
echo ""

# Clean up
echo "🧹 Cleaning up..."
rm -rf ~/.llama-cpp/bin ~/.llama-cpp/metadata

echo ""
echo "📥 Testing download with corrected URL..."
echo "   Expected URL: https://github.com/ggml-org/llama.cpp/releases/download/b7018/llama-b7018-bin-macos-arm64.zip"
echo ""

# Run quick test
go run ./cmd/test-library-paths/main.go

echo ""
if test $status -eq 0
    echo "✅ SUCCESS! Download worked!"
else
    echo "❌ FAILED! Check logs above"
end

