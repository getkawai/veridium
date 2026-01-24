#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Veridium Contract Deployment Script ===${NC}"
echo ""

# Check parameter
if [ -z "$1" ]; then
    echo "Usage: $0 [testnet|mainnet] [--skip-clean]"
    echo ""
    echo "Options:"
    echo "  --skip-clean    Skip forge clean, reuse existing artifacts"
    exit 1
fi

NETWORK=$1
SKIP_CLEAN=false

# Check for --skip-clean flag
if [ "$2" == "--skip-clean" ]; then
    SKIP_CLEAN=true
fi

if [ "$NETWORK" == "testnet" ]; then
    echo -e "${YELLOW}Copying testnet config...${NC}"
    cp contracts/.env.testnet contracts/.env
    
    # Update foundry.toml for testnet
    echo -e "${YELLOW}Updating foundry.toml for testnet...${NC}"
    sed -i.bak 's|eth-rpc-url=".*"|eth-rpc-url="https://testnet-rpc.monad.xyz"|' contracts/foundry.toml
    sed -i.bak 's|chain_id = .*|chain_id = 10143|' contracts/foundry.toml
    rm -f contracts/foundry.toml.bak
    
elif [ "$NETWORK" == "mainnet" ]; then
    echo -e "${RED}Copying mainnet config...${NC}"
    cp contracts/.env.mainnet contracts/.env
    
    # Update foundry.toml for mainnet
    echo -e "${RED}Updating foundry.toml for mainnet...${NC}"
    sed -i.bak 's|eth-rpc-url=".*"|eth-rpc-url="https://rpc.monad.xyz"|' contracts/foundry.toml
    sed -i.bak 's|chain_id = .*|chain_id = 143|' contracts/foundry.toml
    rm -f contracts/foundry.toml.bak
    
    echo -e "${RED}WARNING: You are deploying to MAINNET!${NC}"
    read -p "Are you sure? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        echo "Deployment cancelled"
        exit 1
    fi
else
    echo -e "${RED}Invalid network: $NETWORK${NC}"
    echo "Usage: $0 [testnet|mainnet] [--skip-clean]"
    exit 1
fi

# Check if .env exists
if [ ! -f "contracts/.env" ]; then
    echo -e "${RED}Error: contracts/.env not found${NC}"
    exit 1
fi

# Source the .env file
source contracts/.env

echo -e "${GREEN}Network: ${NETWORK}${NC}"
echo "RPC_URL: $RPC_URL"
echo "CHAIN_ID: $CHAIN_ID"

# Validate CHAIN_ID is set
if [ -z "$CHAIN_ID" ]; then
    echo -e "${RED}Error: CHAIN_ID not found in .env${NC}"
    exit 1
fi

# Setup verification parameters based on network
if [ "$NETWORK" == "mainnet" ]; then
    # Monad Mainnet uses Sourcify via MonadVision
    VERIFY_FLAGS="--verify --verifier sourcify --verifier-url https://sourcify-api-monad.blockvision.org/"
    echo "Verifier: Sourcify (MonadVision Mainnet)"
elif [ "$NETWORK" == "testnet" ]; then
    # Monad Testnet uses Etherscan verifier via MonadVision
    if [ -z "$ETHERSCAN_API_KEY" ]; then
        echo -e "${RED}❌ Error: ETHERSCAN_API_KEY environment variable not set${NC}"
        exit 1
    fi
    VERIFY_FLAGS="--verify --verifier etherscan --etherscan-api-key $ETHERSCAN_API_KEY"
    echo "Verifier: Etherscan (MonadVision Testnet)"
else
    VERIFY_FLAGS=""
    echo -e "${YELLOW}⚠️  No verifier configured for this network${NC}"
fi

# Change to contracts directory for all operations
cd contracts

if [ "$SKIP_CLEAN" = true ]; then
    echo -e "${YELLOW}Skipping forge clean, reusing existing artifacts${NC}"
else
    echo ""
    echo -e "${GREEN}Step 1: Clean old artifacts${NC}"
    forge clean
fi

echo ""
echo -e "${GREEN}Step 2: Compile contracts${NC}"
forge build --build-info

echo ""
echo -e "${GREEN}Step 3: Run tests${NC}"
forge test

echo ""
echo -e "${GREEN}Step 4: Deploy KawaiToken${NC}"
forge script script/DeployKawai.s.sol:DeployKawai \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --legacy \
  $VERIFY_FLAGS

BROADCAST_FILE="broadcast/DeployKawai.s.sol/$CHAIN_ID/run-latest.json"
if [ -f "$BROADCAST_FILE" ]; then
    KAWAI_ADDRESS=$(jq -r '.transactions[] | select(.contractName == "KawaiToken") | .contractAddress' "$BROADCAST_FILE" | head -1)
    if [ -z "$KAWAI_ADDRESS" ] || [ "$KAWAI_ADDRESS" == "null" ]; then
        echo -e "${RED}Failed to extract KawaiToken address${NC}"
        exit 1
    fi
    # Convert to checksum address
    KAWAI_ADDRESS=$(cast to-check-sum-address "$KAWAI_ADDRESS")
    echo -e "${GREEN}KawaiToken deployed at: $KAWAI_ADDRESS${NC}"
    sed -i.bak "s/^KAWAI_TOKEN_ADDRESS=.*/KAWAI_TOKEN_ADDRESS=$KAWAI_ADDRESS/" .env
    echo -e "${GREEN}Updated KAWAI_TOKEN_ADDRESS in contracts/.env${NC}"
    source .env
else
    echo -e "${RED}Broadcast file not found: $BROADCAST_FILE${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Step 5: Deploy OTCMarket${NC}"
forge script script/DeployOTCMarket.s.sol:DeployOTCMarket \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --legacy \
  $VERIFY_FLAGS

BROADCAST_FILE="broadcast/DeployOTCMarket.s.sol/$CHAIN_ID/run-latest.json"
if [ -f "$BROADCAST_FILE" ]; then
    OTC_ADDRESS=$(jq -r '.transactions[] | select(.contractName == "OTCMarket") | .contractAddress' "$BROADCAST_FILE" | head -1)
    if [ -z "$OTC_ADDRESS" ] || [ "$OTC_ADDRESS" == "null" ]; then
        echo -e "${RED}Failed to extract OTCMarket address${NC}"
        exit 1
    fi
    # Convert to checksum address
    OTC_ADDRESS=$(cast to-check-sum-address "$OTC_ADDRESS")
    echo -e "${GREEN}OTCMarket deployed at: $OTC_ADDRESS${NC}"
    sed -i.bak "s/^OTC_MARKET_ADDRESS=.*/OTC_MARKET_ADDRESS=$OTC_ADDRESS/" .env
    echo -e "${GREEN}Updated OTC_MARKET_ADDRESS in contracts/.env${NC}"
    source .env
else
    echo -e "${RED}Broadcast file not found${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Step 6: Deploy PaymentVault${NC}"
forge script script/DeployPaymentVault.s.sol:DeployPaymentVault \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --legacy \
  $VERIFY_FLAGS

BROADCAST_FILE="broadcast/DeployPaymentVault.s.sol/$CHAIN_ID/run-latest.json"
if [ -f "$BROADCAST_FILE" ]; then
    VAULT_ADDRESS=$(jq -r '.transactions[] | select(.contractName == "PaymentVault") | .contractAddress' "$BROADCAST_FILE" | head -1)
    if [ -z "$VAULT_ADDRESS" ] || [ "$VAULT_ADDRESS" == "null" ]; then
        echo -e "${RED}Failed to extract PaymentVault address${NC}"
        exit 1
    fi
    # Convert to checksum address
    VAULT_ADDRESS=$(cast to-check-sum-address "$VAULT_ADDRESS")
    echo -e "${GREEN}PaymentVault deployed at: $VAULT_ADDRESS${NC}"
    sed -i.bak "s/^PAYMENT_VAULT_ADDRESS=.*/PAYMENT_VAULT_ADDRESS=$VAULT_ADDRESS/" .env
    echo -e "${GREEN}Updated PAYMENT_VAULT_ADDRESS in contracts/.env${NC}"
    source .env
else
    echo -e "${RED}Broadcast file not found${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Step 7: Deploy RevenueDistributor${NC}"
forge script script/DeployRevenueDistributor.s.sol:DeployRevenueDistributor \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --legacy \
  $VERIFY_FLAGS

BROADCAST_FILE="broadcast/DeployRevenueDistributor.s.sol/$CHAIN_ID/run-latest.json"
if [ -f "$BROADCAST_FILE" ]; then
    REVENUE_ADDRESS=$(jq -r '.transactions[] | select(.contractName == "RevenueDistributor") | .contractAddress' "$BROADCAST_FILE" | head -1)
    if [ -z "$REVENUE_ADDRESS" ] || [ "$REVENUE_ADDRESS" == "null" ]; then
        echo -e "${RED}Failed to extract RevenueDistributor address${NC}"
        exit 1
    fi
    # Convert to checksum address
    REVENUE_ADDRESS=$(cast to-check-sum-address "$REVENUE_ADDRESS")
    echo -e "${GREEN}RevenueDistributor deployed at: $REVENUE_ADDRESS${NC}"
    sed -i.bak "s/^REVENUE_DISTRIBUTOR_ADDRESS=.*/REVENUE_DISTRIBUTOR_ADDRESS=$REVENUE_ADDRESS/" .env
    echo -e "${GREEN}Updated REVENUE_DISTRIBUTOR_ADDRESS in contracts/.env${NC}"
    source .env
else
    echo -e "${RED}Broadcast file not found${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Step 8: Deploy MiningDistributor${NC}"
# Note: Verification may fail due to via_ir compilation
forge script script/DeployMiningDistributor.s.sol:DeployMiningDistributor \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --legacy \
  $VERIFY_FLAGS || echo -e "${YELLOW}⚠️  Verification failed (expected for via_ir contracts)${NC}"

BROADCAST_FILE="broadcast/DeployMiningDistributor.s.sol/$CHAIN_ID/run-latest.json"
if [ -f "$BROADCAST_FILE" ]; then
    MINING_ADDRESS=$(jq -r '.transactions[] | select(.contractName == "MiningRewardDistributor") | .contractAddress' "$BROADCAST_FILE" | head -1)
    if [ -z "$MINING_ADDRESS" ] || [ "$MINING_ADDRESS" == "null" ]; then
        echo -e "${RED}Failed to extract MiningDistributor address${NC}"
        exit 1
    fi
    # Convert to checksum address
    MINING_ADDRESS=$(cast to-check-sum-address "$MINING_ADDRESS")
    echo -e "${GREEN}MiningDistributor deployed at: $MINING_ADDRESS${NC}"
    sed -i.bak "s/^MINING_DISTRIBUTOR_ADDRESS=.*/MINING_DISTRIBUTOR_ADDRESS=$MINING_ADDRESS/" .env
    echo -e "${GREEN}Updated MINING_DISTRIBUTOR_ADDRESS in .env${NC}"
    source .env
else
    echo -e "${RED}Broadcast file not found${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Step 9: Deploy CashbackDistributor${NC}"
forge script script/DeployCashbackDistributor.s.sol:DeployCashbackDistributor \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --legacy \
  $VERIFY_FLAGS

BROADCAST_FILE="broadcast/DeployCashbackDistributor.s.sol/$CHAIN_ID/run-latest.json"
if [ -f "$BROADCAST_FILE" ]; then
    CASHBACK_ADDRESS=$(jq -r '.transactions[] | select(.contractName == "DepositCashbackDistributor") | .contractAddress' "$BROADCAST_FILE" | head -1)
    if [ -z "$CASHBACK_ADDRESS" ] || [ "$CASHBACK_ADDRESS" == "null" ]; then
        echo -e "${RED}Failed to extract CashbackDistributor address${NC}"
        exit 1
    fi
    # Convert to checksum address
    CASHBACK_ADDRESS=$(cast to-check-sum-address "$CASHBACK_ADDRESS")
    echo -e "${GREEN}CashbackDistributor deployed at: $CASHBACK_ADDRESS${NC}"
    sed -i.bak "s/^CASHBACK_DISTRIBUTOR_ADDRESS=.*/CASHBACK_DISTRIBUTOR_ADDRESS=$CASHBACK_ADDRESS/" .env
    echo -e "${GREEN}Updated CASHBACK_DISTRIBUTOR_ADDRESS in .env${NC}"
    source .env
else
    echo -e "${RED}Broadcast file not found${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Step 10: Deploy ReferralDistributor${NC}"
forge script script/DeployReferralDistributor.s.sol:DeployReferralDistributor \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --legacy \
  $VERIFY_FLAGS

BROADCAST_FILE="broadcast/DeployReferralDistributor.s.sol/$CHAIN_ID/run-latest.json"
if [ -f "$BROADCAST_FILE" ]; then
    REFERRAL_ADDRESS=$(jq -r '.transactions[] | select(.contractName == "ReferralRewardDistributor") | .contractAddress' "$BROADCAST_FILE" | head -1)
    if [ -z "$REFERRAL_ADDRESS" ] || [ "$REFERRAL_ADDRESS" == "null" ]; then
        echo -e "${RED}Failed to extract ReferralDistributor address${NC}"
        exit 1
    fi
    # Convert to checksum address
    REFERRAL_ADDRESS=$(cast to-check-sum-address "$REFERRAL_ADDRESS")
    echo -e "${GREEN}ReferralDistributor deployed at: $REFERRAL_ADDRESS${NC}"
    sed -i.bak "s/^REFERRAL_DISTRIBUTOR_ADDRESS=.*/REFERRAL_DISTRIBUTOR_ADDRESS=$REFERRAL_ADDRESS/" .env
    echo -e "${GREEN}Updated REFERRAL_DISTRIBUTOR_ADDRESS in .env${NC}"
    source .env
else
    echo -e "${RED}Broadcast file not found${NC}"
    exit 1
fi

# Go back to root directory
cd ..

# Clean up backup files
rm -f contracts/.env.bak

echo ""
echo -e "${GREEN}Step 11: Generate Go bindings${NC}"
make contracts-bindings

echo ""
echo -e "${GREEN}Step 12: Update root .env.$NETWORK file${NC}"

# Source the contract addresses
source contracts/.env

# Update root .env file
if [ -f ".env.$NETWORK" ]; then
    echo "Updating .env.$NETWORK with deployed addresses..."
    
    # Update or add contract addresses using extracted variables
    if grep -q "^KAWAI_TOKEN_ADDRESS=" ".env.$NETWORK"; then
        sed -i.bak "s|^KAWAI_TOKEN_ADDRESS=.*|KAWAI_TOKEN_ADDRESS=$KAWAI_ADDRESS|" ".env.$NETWORK"
    else
        echo "KAWAI_TOKEN_ADDRESS=$KAWAI_ADDRESS" >> ".env.$NETWORK"
    fi
    
    if grep -q "^OTC_MARKET_ADDRESS=" ".env.$NETWORK"; then
        sed -i.bak "s|^OTC_MARKET_ADDRESS=.*|OTC_MARKET_ADDRESS=$OTC_ADDRESS|" ".env.$NETWORK"
    else
        echo "OTC_MARKET_ADDRESS=$OTC_ADDRESS" >> ".env.$NETWORK"
    fi
    
    sed -i.bak "s|^PAYMENT_VAULT_ADDRESS=.*|PAYMENT_VAULT_ADDRESS=$VAULT_ADDRESS|" ".env.$NETWORK" || \
        echo "PAYMENT_VAULT_ADDRESS=$VAULT_ADDRESS" >> ".env.$NETWORK"
    
    if grep -q "^REVENUE_DISTRIBUTOR_ADDRESS=" ".env.$NETWORK"; then
        sed -i.bak "s|^REVENUE_DISTRIBUTOR_ADDRESS=.*|REVENUE_DISTRIBUTOR_ADDRESS=$REVENUE_ADDRESS|" ".env.$NETWORK"
    else
        echo "REVENUE_DISTRIBUTOR_ADDRESS=$REVENUE_ADDRESS" >> ".env.$NETWORK"
    fi
    
    sed -i.bak "s|^MINING_DISTRIBUTOR_ADDRESS=.*|MINING_DISTRIBUTOR_ADDRESS=$MINING_ADDRESS|" ".env.$NETWORK" || \
        echo "MINING_DISTRIBUTOR_ADDRESS=$MINING_ADDRESS" >> ".env.$NETWORK"
    
    sed -i.bak "s|^CASHBACK_DISTRIBUTOR_ADDRESS=.*|CASHBACK_DISTRIBUTOR_ADDRESS=$CASHBACK_ADDRESS|" ".env.$NETWORK" || \
        echo "CASHBACK_DISTRIBUTOR_ADDRESS=$CASHBACK_ADDRESS" >> ".env.$NETWORK"
    
    sed -i.bak "s|^REFERRAL_DISTRIBUTOR_ADDRESS=.*|REFERRAL_DISTRIBUTOR_ADDRESS=$REFERRAL_ADDRESS|" ".env.$NETWORK" || \
        echo "REFERRAL_DISTRIBUTOR_ADDRESS=$REFERRAL_ADDRESS" >> ".env.$NETWORK"
    
    # Remove backup file
    rm -f ".env.$NETWORK.bak"
    
    echo -e "${GREEN}✅ Updated .env.$NETWORK${NC}"
else
    echo -e "${YELLOW}⚠️  .env.$NETWORK not found, skipping root .env update${NC}"
fi

echo ""
echo -e "${GREEN}Step 13: Grant MINTER_ROLE to distributors${NC}"

# MINTER_ROLE = keccak256("MINTER_ROLE") - OpenZeppelin standard
MINTER_ROLE="0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6"

echo "MINTER_ROLE: $MINTER_ROLE"
echo "Granting MINTER_ROLE to:"
echo "  1. MiningRewardDistributor: $MINING_ADDRESS"
echo "  2. CashbackDistributor: $CASHBACK_ADDRESS"
echo "  3. ReferralDistributor: $REFERRAL_ADDRESS"
echo ""

# Track failures
GRANT_FAILED=0
FAILED_GRANTS=""

# Grant to MiningRewardDistributor
echo "Granting to MiningRewardDistributor..."
if cast send "$KAWAI_ADDRESS" \
    "grantRole(bytes32,address)" \
    "$MINTER_ROLE" \
    "$MINING_ADDRESS" \
    --private-key "$PRIVATE_KEY" \
    --rpc-url "$RPC_URL" \
    --gas-limit 100000 \
    > /dev/null 2>&1; then
    echo -e "${GREEN}✅ MiningRewardDistributor granted${NC}"
else
    echo -e "${RED}❌ Failed to grant to MiningRewardDistributor${NC}"
    GRANT_FAILED=1
    FAILED_GRANTS="$FAILED_GRANTS\n  - MiningRewardDistributor"
fi

# Grant to CashbackDistributor
echo "Granting to CashbackDistributor..."
if cast send "$KAWAI_ADDRESS" \
    "grantRole(bytes32,address)" \
    "$MINTER_ROLE" \
    "$CASHBACK_ADDRESS" \
    --private-key "$PRIVATE_KEY" \
    --rpc-url "$RPC_URL" \
    --gas-limit 100000 \
    > /dev/null 2>&1; then
    echo -e "${GREEN}✅ CashbackDistributor granted${NC}"
else
    echo -e "${RED}❌ Failed to grant to CashbackDistributor${NC}"
    GRANT_FAILED=1
    FAILED_GRANTS="$FAILED_GRANTS\n  - CashbackDistributor"
fi

# Grant to ReferralDistributor
echo "Granting to ReferralDistributor..."
if cast send "$KAWAI_ADDRESS" \
    "grantRole(bytes32,address)" \
    "$MINTER_ROLE" \
    "$REFERRAL_ADDRESS" \
    --private-key "$PRIVATE_KEY" \
    --rpc-url "$RPC_URL" \
    --gas-limit 100000 \
    > /dev/null 2>&1; then
    echo -e "${GREEN}✅ ReferralDistributor granted${NC}"
else
    echo -e "${RED}❌ Failed to grant to ReferralDistributor${NC}"
    GRANT_FAILED=1
    FAILED_GRANTS="$FAILED_GRANTS\n  - ReferralDistributor"
fi

echo ""
echo -e "${GREEN}Step 14: Verify MINTER_ROLE grants${NC}"

# Verify each distributor has MINTER_ROLE
VERIFY_FAILED=0
FAILED_VERIFICATIONS=""

echo "Verifying MiningRewardDistributor..."
HAS_ROLE=$(cast call "$KAWAI_ADDRESS" \
    "hasRole(bytes32,address)" \
    "$MINTER_ROLE" \
    "$MINING_ADDRESS" \
    --rpc-url "$RPC_URL" 2>/dev/null)

if [[ "$HAS_ROLE" == *"0000000000000000000000000000000000000000000000000000000000000001"* ]]; then
    echo -e "${GREEN}✅ MiningRewardDistributor verified${NC}"
else
    echo -e "${RED}❌ MiningRewardDistributor verification failed${NC}"
    VERIFY_FAILED=1
    FAILED_VERIFICATIONS="$FAILED_VERIFICATIONS\n  - MiningRewardDistributor"
fi

echo "Verifying CashbackDistributor..."
HAS_ROLE=$(cast call "$KAWAI_ADDRESS" \
    "hasRole(bytes32,address)" \
    "$MINTER_ROLE" \
    "$CASHBACK_ADDRESS" \
    --rpc-url "$RPC_URL" 2>/dev/null)

if [[ "$HAS_ROLE" == *"0000000000000000000000000000000000000000000000000000000000000001"* ]]; then
    echo -e "${GREEN}✅ CashbackDistributor verified${NC}"
else
    echo -e "${RED}❌ CashbackDistributor verification failed${NC}"
    VERIFY_FAILED=1
    FAILED_VERIFICATIONS="$FAILED_VERIFICATIONS\n  - CashbackDistributor"
fi

echo "Verifying ReferralDistributor..."
HAS_ROLE=$(cast call "$KAWAI_ADDRESS" \
    "hasRole(bytes32,address)" \
    "$MINTER_ROLE" \
    "$REFERRAL_ADDRESS" \
    --rpc-url "$RPC_URL" 2>/dev/null)

if [[ "$HAS_ROLE" == *"0000000000000000000000000000000000000000000000000000000000000001"* ]]; then
    echo -e "${GREEN}✅ ReferralDistributor verified${NC}"
else
    echo -e "${RED}❌ ReferralDistributor verification failed${NC}"
    VERIFY_FAILED=1
    FAILED_VERIFICATIONS="$FAILED_VERIFICATIONS\n  - ReferralDistributor"
fi

echo ""
echo -e "${GREEN}=== Deployment Complete ===${NC}"
echo ""
echo "Deployed addresses:"
echo "KAWAI_TOKEN_ADDRESS=$KAWAI_TOKEN_ADDRESS"
echo "OTC_MARKET_ADDRESS=$OTC_MARKET_ADDRESS"
echo "PAYMENT_VAULT_ADDRESS=$PAYMENT_VAULT_ADDRESS"
echo "REVENUE_DISTRIBUTOR_ADDRESS=$REVENUE_DISTRIBUTOR_ADDRESS"
echo "MINING_DISTRIBUTOR_ADDRESS=$MINING_DISTRIBUTOR_ADDRESS"
echo "CASHBACK_DISTRIBUTOR_ADDRESS=$CASHBACK_DISTRIBUTOR_ADDRESS"
echo "REFERRAL_DISTRIBUTOR_ADDRESS=$REFERRAL_DISTRIBUTOR_ADDRESS"
echo ""

# Check if any grants failed
if [ $GRANT_FAILED -eq 1 ]; then
    echo -e "${RED}⚠️  Warning: Some MINTER_ROLE grants failed:${NC}"
    echo -e "${RED}$FAILED_GRANTS${NC}"
    echo ""
    echo -e "${YELLOW}You can manually grant roles using:${NC}"
    echo "  ./GRANT_ALL_MINTER_ROLES.sh"
    echo ""
    exit 1
fi

# Check if any verifications failed
if [ $VERIFY_FAILED -eq 1 ]; then
    echo -e "${RED}⚠️  Warning: Some MINTER_ROLE verifications failed:${NC}"
    echo -e "${RED}$FAILED_VERIFICATIONS${NC}"
    echo ""
    echo -e "${YELLOW}Roles may not have been granted correctly. Please check manually.${NC}"
    echo ""
    exit 1
fi

echo ""
echo -e "${GREEN}Step 15: Generate constants for backend${NC}"

# Backup current .env (just in case)
if [ -f ".env" ]; then
    echo "Backing up current .env to .env.backup..."
    cp .env .env.backup
fi

# Copy network-specific .env to root
echo "Copying .env.$NETWORK to .env..."
cp ".env.$NETWORK" .env

# Generate constants
echo "Running make constants-generate..."
if make constants-generate > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Constants generated successfully${NC}"
    echo -e "${GREEN}✅ Root .env updated to $NETWORK config${NC}"
else
    echo -e "${RED}❌ Failed to generate constants${NC}"
    echo -e "${YELLOW}You can manually generate with: make constants-generate${NC}"
fi

echo -e "${GREEN}All setup complete! Ready to use.${NC}"
echo ""
echo -e "${YELLOW}Verify deployment:${NC}"
echo "1. Check contracts on explorer: https://$NETWORK.monad.xyz"
echo "2. Restart backend to use new contract addresses"
echo "3. Test reward claiming in the app"
echo ""
echo -e "${YELLOW}Note: Root .env is now using $NETWORK config${NC}"
echo "Backup saved to: .env.backup"
