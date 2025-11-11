#!/usr/bin/env fish

echo "🧪 Testing URL Construction (Unit Test)"
echo "========================================"
echo ""

cd /Users/yuda/github.com/kawai-network/veridium

echo "Running URL construction tests..."
go test -v ./pkg/yzma/download -run TestURL

echo ""
if test $status -eq 0
    echo "✅ All URL tests PASSED!"
    echo ""
    echo "Now test real download:"
    echo "  chmod +x test-url-fix.fish"
    echo "  ./test-url-fix.fish"
else
    echo "❌ URL tests FAILED!"
end

