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

# Add environment variables
echo "Adding R2_ACCOUNT_ID..."
circleci project secret create github kawai-network veridium R2_ACCOUNT_ID \
  --env-value "ceab218751d33cd804878196ad7bef74"

echo "Adding R2_ACCESS_KEY_ID..."
circleci project secret create github kawai-network veridium R2_ACCESS_KEY_ID \
  --env-value "a71e802dd7c1ab8cf407ffb937cdf6a8"

echo "Adding R2_SECRET_ACCESS_KEY..."
circleci project secret create github kawai-network veridium R2_SECRET_ACCESS_KEY \
  --env-value "0e3ce0d92faa9b337c83131efc7a4a64bb6f313171c309d5cb9a0fb76926d0ca"

echo "Adding R2_ENDPOINT_URL..."
circleci project secret create github kawai-network veridium R2_ENDPOINT_URL \
  --env-value "https://ceab218751d33cd804878196ad7bef74.r2.cloudflarestorage.com"

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
