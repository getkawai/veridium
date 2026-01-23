#!/bin/bash
set -e

# ============================================================
# Grant MINTER_ROLE to ALL Reward Distributors
# ============================================================
# This script grants MINTER_ROLE on KawaiToken to all contracts
# that need to mint KAWAI rewards.
#
# Required contracts:
# 1. MiningRewardDistributor (mining rewards)
# 2. CashbackDistributor (deposit cashback)
# 3. ReferralRewardDistributor (referral rewards)
# ============================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}  Grant MINTER_ROLE to All Reward Distributors${NC}"
echo -e "${BLUE}============================================================${NC}"
echo ""

# ============================================================
# Configuration
# ============================================================

# Monad Testnet
RPC_URL="https://testnet.monad.xyz/"
CHAIN_ID=10143

# Source contract addresses from contracts/.env
if [ -f "contracts/.env" ]; then
    source contracts/.env
    KAWAI_TOKEN="$KAWAI_TOKEN_ADDRESS"
    MINING_DISTRIBUTOR="$MINING_DISTRIBUTOR_ADDRESS"
    CASHBACK_DISTRIBUTOR="$CASHBACK_DISTRIBUTOR_ADDRESS"
    REFERRAL_DISTRIBUTOR="$REFERRAL_DISTRIBUTOR_ADDRESS"
else
    echo -e "${RED}❌ Error: contracts/.env not found${NC}"
    echo ""
    echo "Please run deployment first:"
    echo "  make deploy-testnet"
    echo ""
    exit 1
fi

# MINTER_ROLE = keccak256("MINTER_ROLE") - OpenZeppelin standard
MINTER_ROLE="0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6"

# ============================================================
# Validation
# ============================================================

if [ -z "$PRIVATE_KEY" ]; then
    echo -e "${RED}❌ Error: PRIVATE_KEY environment variable not set${NC}"
    echo ""
    echo "Usage:"
    echo "  export PRIVATE_KEY=0x..."
    echo "  ./GRANT_ALL_MINTER_ROLES.sh"
    echo ""
    exit 1
fi

echo -e "${GREEN}✅ Configuration loaded${NC}"
echo "  RPC: $RPC_URL"
echo "  Chain ID: $CHAIN_ID"
echo ""
echo "  KawaiToken: $KAWAI_TOKEN"
echo "  MiningRewardDistributor: $MINING_DISTRIBUTOR"
echo "  CashbackDistributor: $CASHBACK_DISTRIBUTOR"
echo "  ReferralDistributor: $REFERRAL_DISTRIBUTOR"
echo ""

# ============================================================
# Helper Functions
# ============================================================

check_role() {
    local contract=$1
    local role=$2
    local account=$3
    
    echo -e "${YELLOW}🔍 Checking if $account has role...${NC}"
    
    # hasRole(bytes32 role, address account) returns (bool)
    local calldata="0x91d14854${role:2}000000000000000000000000${account:2}"
    
    local result=$(cast call "$contract" "$calldata" --rpc-url "$RPC_URL" 2>&1)
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Failed to check role: $result${NC}"
        return 1
    fi
    
    # Parse result (0x0...0 = false, 0x0...1 = true)
    if [[ "$result" == *"0000000000000000000000000000000000000000000000000000000000000001" ]]; then
        echo -e "${GREEN}✅ Already has MINTER_ROLE${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Does not have MINTER_ROLE${NC}"
        return 1
    fi
}

grant_role() {
    local contract=$1
    local role=$2
    local account=$3
    local name=$4
    
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  Granting MINTER_ROLE to: $name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo "  Contract: $account"
    echo ""
    
    # Check if already has role
    if check_role "$contract" "$role" "$account"; then
        echo -e "${GREEN}✅ Skipping (already granted)${NC}"
        return 0
    fi
    
    echo ""
    echo -e "${YELLOW}📝 Granting role...${NC}"
    
    # grantRole(bytes32 role, address account)
    local tx_hash=$(cast send "$contract" \
        "grantRole(bytes32,address)" \
        "$role" \
        "$account" \
        --private-key "$PRIVATE_KEY" \
        --rpc-url "$RPC_URL" \
        --chain-id "$CHAIN_ID" \
        --gas-limit 100000 \
        --json 2>&1 | jq -r '.transactionHash // empty')
    
    if [ -z "$tx_hash" ]; then
        echo -e "${RED}❌ Failed to grant role${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ Role granted!${NC}"
    echo "  TX: $tx_hash"
    echo "  Explorer: https://testnet.monad.xyz/tx/$tx_hash"
    
    # Wait for confirmation
    echo ""
    echo -e "${YELLOW}⏳ Waiting for confirmation...${NC}"
    sleep 3
    
    # Verify
    if check_role "$contract" "$role" "$account"; then
        echo -e "${GREEN}✅ Verification successful!${NC}"
        return 0
    else
        echo -e "${RED}❌ Verification failed${NC}"
        return 1
    fi
}

# ============================================================
# Main Execution
# ============================================================

echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}  Starting Role Grants${NC}"
echo -e "${BLUE}============================================================${NC}"

# Track success/failure
SUCCESS_COUNT=0
FAIL_COUNT=0

# 1. Grant to MiningRewardDistributor
if grant_role "$KAWAI_TOKEN" "$MINTER_ROLE" "$MINING_DISTRIBUTOR" "MiningRewardDistributor"; then
    ((SUCCESS_COUNT++))
else
    ((FAIL_COUNT++))
fi

# 2. Grant to CashbackDistributor
if grant_role "$KAWAI_TOKEN" "$MINTER_ROLE" "$CASHBACK_DISTRIBUTOR" "CashbackDistributor"; then
    ((SUCCESS_COUNT++))
else
    ((FAIL_COUNT++))
fi

# 3. Grant to ReferralDistributor
if grant_role "$KAWAI_TOKEN" "$MINTER_ROLE" "$REFERRAL_DISTRIBUTOR" "ReferralDistributor"; then
    ((SUCCESS_COUNT++))
else
    ((FAIL_COUNT++))
fi

# ============================================================
# Summary
# ============================================================

echo ""
echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}  Summary${NC}"
echo -e "${BLUE}============================================================${NC}"
echo ""
echo -e "  ${GREEN}✅ Successful: $SUCCESS_COUNT${NC}"
echo -e "  ${RED}❌ Failed: $FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}🎉 All roles granted successfully!${NC}"
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo "  1. Test mining reward claim"
    echo "  2. Test cashback reward claim"
    echo "  3. Test referral reward claim"
    echo ""
    exit 0
else
    echo -e "${RED}⚠️  Some role grants failed. Please check the errors above.${NC}"
    echo ""
    exit 1
fi

