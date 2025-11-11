#!/bin/bash

echo "🧪 Testing Library Path Management"
echo "==================================="
echo ""

cd "$(dirname "$0")"

echo "Building test program..."
go build -o test-library-paths ./cmd/test-library-paths/main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"
echo ""
echo "Running tests..."
echo ""

./test-library-paths

exit_code=$?

# Cleanup
rm -f test-library-paths

exit $exit_code

