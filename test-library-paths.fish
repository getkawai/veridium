#!/usr/bin/env fish

echo "🧪 Testing Library Path Management"
echo "==================================="
echo ""

cd (dirname (status -f))

echo "Building test program..."
go build -o test-library-paths ./cmd/test-library-paths/main.go

if test $status -ne 0
    echo "❌ Build failed"
    exit 1
end

echo "✅ Build successful"
echo ""
echo "Running tests..."
echo ""

./test-library-paths

set exit_code $status

# Cleanup
rm -f test-library-paths

exit $exit_code

