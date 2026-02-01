#!/bin/bash
# Script to check R2 Cloudflare bucket contents
# Usage: ./scripts/check-r2.sh [path]

set -e

# R2 Configuration from environment or .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

R2_ACCESS_KEY_ID="${R2_ACCESS_KEY_ID:-}"
R2_SECRET_ACCESS_KEY="${R2_SECRET_ACCESS_KEY:-}"
R2_ENDPOINT_URL="${R2_ENDPOINT_URL:-}"
BUCKET="${R2_BUCKET:-kawai}"
PREFIX="${1:-}"

if [ -z "$R2_ACCESS_KEY_ID" ] || [ -z "$R2_SECRET_ACCESS_KEY" ] || [ -z "$R2_ENDPOINT_URL" ]; then
    echo "Error: R2 credentials not found"
    echo "Please set R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, and R2_ENDPOINT_URL"
    exit 1
fi

# Export AWS credentials for AWS CLI to use R2 keys
export AWS_ACCESS_KEY_ID="$R2_ACCESS_KEY_ID"
export AWS_SECRET_ACCESS_KEY="$R2_SECRET_ACCESS_KEY"

echo "========================================"
echo "R2 Cloudflare Bucket Check"
echo "========================================"
echo "Endpoint: $R2_ENDPOINT_URL"
echo "Bucket: $BUCKET"
echo "Prefix: ${PREFIX:-(root)}"
echo ""

# Use AWS CLI with R2 endpoint
if command -v aws &> /dev/null; then
    echo "Listing objects..."
    aws s3 ls "s3://$BUCKET/$PREFIX" \
        --endpoint-url "$R2_ENDPOINT_URL" \
        --recursive \
        --human-readable \
        --summarize 2>/dev/null || {
        echo "Error: Failed to list objects"
        echo "Make sure AWS CLI is configured with the R2 credentials"
        exit 1
    }
else
    echo "AWS CLI not found. Trying with curl..."
    
    # Generate signature for R2 (S3 compatible)
    DATE=$(date -u +"%Y%m%dT%H%M%SZ")
    DATE_SHORT=$(date -u +"%Y%m%d")
    REGION="auto"
    SERVICE="s3"
    
    # Create canonical request
    HTTP_METHOD="GET"
    CANONICAL_URI="/"
    CANONICAL_QUERYSTRING="list-type=2&prefix=$PREFIX"
    CANONICAL_HEADERS="host:$(echo $R2_ENDPOINT_URL | sed 's|https://||')\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\nx-amz-date:$DATE\n"
    SIGNED_HEADERS="host;x-amz-content-sha256;x-amz-date"
    PAYLOAD_HASH="e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    
    CANONICAL_REQUEST="$HTTP_METHOD\n$CANONICAL_URI\n$CANONICAL_QUERYSTRING\n$CANONICAL_HEADERS\n$SIGNED_HEADERS\n$PAYLOAD_HASH"
    
    # Create string to sign
    CREDENTIAL_SCOPE="$DATE_SHORT/$REGION/$SERVICE/aws4_request"
    STRING_TO_SIGN="AWS4-HMAC-SHA256\n$DATE\n$CREDENTIAL_SCOPE\n$(echo -n "$CANONICAL_REQUEST" | sha256sum | cut -d' ' -f1)"
    
    # Calculate signature
    kDate=$(echo -n "$DATE_SHORT" | openssl dgst -sha256 -mac HMAC -macopt key:"AWS4$R2_SECRET_ACCESS_KEY" | cut -d' ' -f2)
    kRegion=$(echo -n "$REGION" | openssl dgst -sha256 -mac HMAC -macopt hexkey:"$kDate" | cut -d' ' -f2)
    kService=$(echo -n "$SERVICE" | openssl dgst -sha256 -mac HMAC -macopt hexkey:"$kRegion" | cut -d' ' -f2)
    kSigning=$(echo -n "aws4_request" | openssl dgst -sha256 -mac HMAC -macopt hexkey:"$kService" | cut -d' ' -f2)
    SIGNATURE=$(echo -n "$STRING_TO_SIGN" | openssl dgst -sha256 -mac HMAC -macopt hexkey:"$kSigning" | cut -d' ' -f2)
    
    # Make request
    HOST=$(echo $R2_ENDPOINT_URL | sed 's|https://||')
    curl -s "$R2_ENDPOINT_URL/?list-type=2&prefix=$PREFIX" \
        -H "Host: $HOST" \
        -H "X-Amz-Date: $DATE" \
        -H "X-Amz-Content-SHA256: $PAYLOAD_HASH" \
        -H "Authorization: AWS4-HMAC-SHA256 Credential=$R2_ACCESS_KEY_ID/$CREDENTIAL_SCOPE, SignedHeaders=$SIGNED_HEADERS, Signature=$SIGNATURE" \
        | xmllint --format - 2>/dev/null || cat
fi

echo ""
echo "========================================"
echo "Check complete!"
echo "========================================"
