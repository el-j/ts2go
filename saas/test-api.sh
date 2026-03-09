#!/bin/bash
# Test script for ts2go SaaS API

set -e

API_URL="${API_URL:-http://localhost:8080}"
TEST_PASSWORD="${TEST_PASSWORD:-SecurePass123!}"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Testing ts2go SaaS API at $API_URL"
echo ""

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
HEALTH=$(curl -s "$API_URL/health")
if echo "$HEALTH" | jq -e '.status == "healthy"' > /dev/null; then
    echo -e "${GREEN}✓ Health check passed${NC}"
    echo "$HEALTH" | jq .
else
    echo -e "${RED}✗ Health check failed${NC}"
    echo "$HEALTH"
    exit 1
fi
echo ""

# Test 2: User Registration
echo -e "${YELLOW}Test 2: User Registration${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"testuser$(date +%s)@example.com\",
    \"username\": \"testuser$(date +%s)\",
    \"password\": \"$TEST_PASSWORD\",
    \"full_name\": \"Test User\"
  }")

ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.access_token')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.refresh_token')

if [ "$ACCESS_TOKEN" != "null" ] && [ -n "$ACCESS_TOKEN" ]; then
    echo -e "${GREEN}✓ Registration successful${NC}"
    echo "User: $(echo "$REGISTER_RESPONSE" | jq -r '.user.username')"
    echo "Email: $(echo "$REGISTER_RESPONSE" | jq -r '.user.email')"
else
    echo -e "${RED}✗ Registration failed${NC}"
    echo "$REGISTER_RESPONSE" | jq .
    exit 1
fi
echo ""

# Test 3: Get Current User (with JWT)
echo -e "${YELLOW}Test 3: Get Current User (Authenticated)${NC}"
ME_RESPONSE=$(curl -s "$API_URL/api/v1/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$ME_RESPONSE" | jq -e '.username' > /dev/null; then
    echo -e "${GREEN}✓ Get current user successful${NC}"
    echo "$ME_RESPONSE" | jq .
else
    echo -e "${RED}✗ Get current user failed${NC}"
    echo "$ME_RESPONSE"
    exit 1
fi
echo ""

# Test 4: Create API Key
echo -e "${YELLOW}Test 4: Create API Key${NC}"
API_KEY_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/api-keys" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test API Key",
    "scopes": ["read:projects", "write:transpile"]
  }')

API_KEY=$(echo "$API_KEY_RESPONSE" | jq -r '.api_key')

if [ "$API_KEY" != "null" ] && [ -n "$API_KEY" ]; then
    echo -e "${GREEN}✓ API key created${NC}"
    echo "Key: ${API_KEY:0:20}..."
    echo "Prefix: $(echo "$API_KEY_RESPONSE" | jq -r '.key_prefix')"
else
    echo -e "${RED}✗ API key creation failed${NC}"
    echo "$API_KEY_RESPONSE" | jq .
    exit 1
fi
echo ""

# Test 5: List API Keys
echo -e "${YELLOW}Test 5: List API Keys${NC}"
LIST_KEYS_RESPONSE=$(curl -s "$API_URL/api/v1/api-keys" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$LIST_KEYS_RESPONSE" | jq -e 'length > 0' > /dev/null; then
    echo -e "${GREEN}✓ List API keys successful${NC}"
    echo "Total keys: $(echo "$LIST_KEYS_RESPONSE" | jq 'length')"
else
    echo -e "${RED}✗ List API keys failed${NC}"
    echo "$LIST_KEYS_RESPONSE"
    exit 1
fi
echo ""

# Test 6: Refresh Token
echo -e "${YELLOW}Test 6: Refresh Token${NC}"
REFRESH_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.access_token')

if [ "$NEW_ACCESS_TOKEN" != "null" ] && [ -n "$NEW_ACCESS_TOKEN" ]; then
    echo -e "${GREEN}✓ Token refresh successful${NC}"
    echo "New token: ${NEW_ACCESS_TOKEN:0:20}..."
else
    echo -e "${RED}✗ Token refresh failed${NC}"
    echo "$REFRESH_RESPONSE" | jq .
    exit 1
fi
echo ""

# Test 7: Login with existing user
echo -e "${YELLOW}Test 7: Login (Existing User)${NC}"
EMAIL=$(echo "$REGISTER_RESPONSE" | jq -r '.user.email')
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email_or_username\": \"$EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

if echo "$LOGIN_RESPONSE" | jq -e '.access_token' > /dev/null; then
    echo -e "${GREEN}✓ Login successful${NC}"
    echo "User: $(echo "$LOGIN_RESPONSE" | jq -r '.user.username')"
else
    echo -e "${RED}✗ Login failed${NC}"
    echo "$LOGIN_RESPONSE" | jq .
    exit 1
fi
echo ""

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ All tests passed!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
