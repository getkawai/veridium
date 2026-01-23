#!/bin/bash

# Script to switch between testnet and mainnet environments
# Usage: ./switch-env.sh [testnet|mainnet]

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 [testnet|mainnet]"
    echo ""
    echo "Current environment:"
    if grep -q "testnet-rpc.monad.xyz" contracts/.env 2>/dev/null; then
        echo "  TESTNET"
    elif grep -q "rpc.monad.xyz" contracts/.env 2>/dev/null; then
        echo "  MAINNET"
    else
        echo "  UNKNOWN"
    fi
    exit 1
fi

ENV=$1

case $ENV in
    testnet)
        echo "Switching to TESTNET..."
        cp contracts/.env.testnet contracts/.env
        echo "✓ Switched to TESTNET"
        echo ""
        echo "RPC: https://testnet-rpc.monad.xyz"
        ;;
    mainnet)
        echo "⚠️  WARNING: Switching to MAINNET"
        echo "This will use REAL funds on Monad Mainnet!"
        echo ""
        read -p "Are you sure? (type 'yes' to confirm): " confirm
        if [ "$confirm" != "yes" ]; then
            echo "Cancelled."
            exit 1
        fi
        cp contracts/.env.mainnet contracts/.env
        echo "✓ Switched to MAINNET"
        echo ""
        echo "RPC: https://rpc.monad.xyz"
        ;;
    *)
        echo "Invalid environment: $ENV"
        echo "Usage: $0 [testnet|mainnet]"
        exit 1
        ;;
esac