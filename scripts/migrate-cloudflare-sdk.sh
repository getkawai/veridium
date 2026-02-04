#!/bin/bash

# Migration script for Cloudflare Go SDK v0.114.0 -> v6.6.0
# This script updates all remaining files with old SDK calls

set -e

echo "🔄 Starting Cloudflare SDK migration..."

# List of files to migrate
FILES=(
    "pkg/store/contributor.go"
    "pkg/store/referral.go"
    "pkg/store/period_counter.go"
    "pkg/store/merkle.go"
    "pkg/store/settlement.go"
    "pkg/store/job_rewards.go"
    "pkg/store/cashback_kv.go"
)

# Backup files
echo "📦 Creating backups..."
for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        cp "$file" "$file.backup"
        echo "  ✓ Backed up $file"
    fi
done

# Remove old import - delete entire line instead of leaving empty string
echo "🔧 Removing old cloudflare-go imports..."
for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        # Use sed to delete the entire import line
        sed -i.tmp '/^[[:space:]]*"github\.com\/cloudflare\/cloudflare-go"[[:space:]]*$/d' "$file"
        rm -f "$file.tmp"
        echo "  ✓ Updated imports in $file"
    fi
done

echo "✅ Migration preparation complete!"
echo ""
echo "⚠️  Manual steps required:"
echo "1. Update all WriteWorkersKVEntry calls to s.client.SetValue()"
echo "2. Update all GetWorkersKV calls to s.client.GetValue()"
echo "3. Update all DeleteWorkersKVEntry calls to s.client.DeleteValue()"
echo "4. Update all ListWorkersKVKeys calls to s.client.ListKeys()"
echo "5. Run: go mod tidy"
echo "6. Run: make test"
echo ""
echo "To restore backups: for f in pkg/store/*.backup; do mv \$f \${f%.backup}; done"
