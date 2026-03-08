#!/bin/bash

# CircleCI Setup Script
# Run this after following the project on CircleCI web UI

echo "Setting up CircleCI environment variables..."

# Check if project exists
if ! circleci project secret list github kawai-network veridium &>/dev/null; then
  echo "❌ Project not found. Please follow the project first:"
  echo "   1. Go to https://app.circleci.com/"
  echo "   2. Click 'Projects' → Find 'veridium' → Click 'Set Up Project'"
  echo "   3. Select 'Use Existing Config' (we already have .circleci/config.yml)"
  echo "   4. Run this script again"
  exit 1
fi

# Require values from local environment instead of hardcoded secrets.
required_vars=(
  R2_ACCOUNT_ID
  R2_ACCESS_KEY_ID
  R2_SECRET_ACCESS_KEY
  R2_ENDPOINT_URL
)

for var in "${required_vars[@]}"; do
  if [ -z "${!var}" ]; then
    echo "❌ Missing required environment variable: ${var}"
    echo "   Export values first, for example:"
    echo "   export ${var}=<your-value>"
    exit 1
  fi
done

# Add environment variables
echo "Adding R2_ACCOUNT_ID..."
circleci project secret create github kawai-network veridium R2_ACCOUNT_ID \
  --env-value "$R2_ACCOUNT_ID"

echo "Adding R2_ACCESS_KEY_ID..."
circleci project secret create github kawai-network veridium R2_ACCESS_KEY_ID \
  --env-value "$R2_ACCESS_KEY_ID"

echo "Adding R2_SECRET_ACCESS_KEY..."
circleci project secret create github kawai-network veridium R2_SECRET_ACCESS_KEY \
  --env-value "$R2_SECRET_ACCESS_KEY"

echo "Adding R2_ENDPOINT_URL..."
circleci project secret create github kawai-network veridium R2_ENDPOINT_URL \
  --env-value "$R2_ENDPOINT_URL"

echo ""
echo "⚠️  GITHUB_TOKEN needs to be added manually:"
echo "   1. Go to https://app.circleci.com/settings/project/github/kawai-network/veridium/environment-variables"
echo "   2. Click 'Add Environment Variable'"
echo "   3. Name: GITHUB_TOKEN"
echo "   4. Value: <your-github-personal-access-token>"
echo ""
echo "✅ Setup complete! You can now trigger a release:"
echo "   git tag v0.1.1"
echo "   git push origin v0.1.1"
