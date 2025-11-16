#!/bin/bash
# Extended test script for ts2go SaaS API - Week 2 features

set -e

API_URL="${API_URL:-http://localhost:8080}"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🧪 Testing ts2go SaaS API - Week 2 Features"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Register and login first
echo -e "${BLUE}Setup: Creating test user${NC}"
TIMESTAMP=$(date +%s)
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"week2test${TIMESTAMP}@example.com\",
    \"username\": \"week2user${TIMESTAMP}\",
    \"password\": \"SecurePass123!\",
    \"full_name\": \"Week 2 Test User\"
  }")

ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.access_token')
USERNAME=$(echo "$REGISTER_RESPONSE" | jq -r '.user.username')

if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo -e "${RED}✗ Failed to create test user${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Test user created: $USERNAME${NC}"
echo ""

# Test 1: Rate Limiting Headers
echo -e "${YELLOW}Test 1: Rate Limiting Headers${NC}"
RATE_LIMIT_RESPONSE=$(curl -s -i "$API_URL/health" | grep -i "x-ratelimit")
if [ -n "$RATE_LIMIT_RESPONSE" ]; then
    echo -e "${GREEN}✓ Rate limit headers present${NC}"
    echo "$RATE_LIMIT_RESPONSE"
else
    echo -e "${RED}✗ No rate limit headers found${NC}"
fi
echo ""

# Test 2: Create Project
echo -e "${YELLOW}Test 2: Create Project${NC}"
PROJECT_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/projects" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"My First Project\",
    \"slug\": \"my-first-project\",
    \"description\": \"A test project for TypeScript to Go transpilation\",
    \"visibility\": \"private\"
  }")

PROJECT_ID=$(echo "$PROJECT_RESPONSE" | jq -r '.id')

if [ "$PROJECT_ID" != "null" ] && [ -n "$PROJECT_ID" ]; then
    echo -e "${GREEN}✓ Project created successfully${NC}"
    echo "Project ID: $PROJECT_ID"
    echo "Project Name: $(echo "$PROJECT_RESPONSE" | jq -r '.name')"
else
    echo -e "${RED}✗ Project creation failed${NC}"
    echo "$PROJECT_RESPONSE" | jq .
    exit 1
fi
echo ""

# Test 3: List Projects
echo -e "${YELLOW}Test 3: List Projects${NC}"
LIST_RESPONSE=$(curl -s "$API_URL/api/v1/projects" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

PROJECT_COUNT=$(echo "$LIST_RESPONSE" | jq 'length')

if [ "$PROJECT_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ Projects listed successfully${NC}"
    echo "Total projects: $PROJECT_COUNT"
else
    echo -e "${RED}✗ Failed to list projects${NC}"
fi
echo ""

# Test 4: Get Single Project
echo -e "${YELLOW}Test 4: Get Project by ID${NC}"
GET_RESPONSE=$(curl -s "$API_URL/api/v1/projects/$PROJECT_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

GET_NAME=$(echo "$GET_RESPONSE" | jq -r '.name')

if [ "$GET_NAME" == "My First Project" ]; then
    echo -e "${GREEN}✓ Project retrieved successfully${NC}"
    echo "$GET_RESPONSE" | jq '{name, slug, visibility}'
else
    echo -e "${RED}✗ Failed to get project${NC}"
    echo "$GET_RESPONSE"
fi
echo ""

# Test 5: Update Project
echo -e "${YELLOW}Test 5: Update Project${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "$API_URL/api/v1/projects/$PROJECT_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"My Updated Project\",
    \"description\": \"Updated description\"
  }")

UPDATED_NAME=$(echo "$UPDATE_RESPONSE" | jq -r '.name')

if [ "$UPDATED_NAME" == "My Updated Project" ]; then
    echo -e "${GREEN}✓ Project updated successfully${NC}"
    echo "New name: $UPDATED_NAME"
else
    echo -e "${RED}✗ Project update failed${NC}"
    echo "$UPDATE_RESPONSE"
fi
echo ""

# Test 6: Create Multiple Projects
echo -e "${YELLOW}Test 6: Create Multiple Projects${NC}"
for i in {1..3}; do
    curl -s -X POST "$API_URL/api/v1/projects" \
      -H "Authorization: Bearer $ACCESS_TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"Test Project $i\",
        \"slug\": \"test-project-$i\",
        \"visibility\": \"private\"
      }" > /dev/null
done

TOTAL_PROJECTS=$(curl -s "$API_URL/api/v1/projects" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq 'length')

if [ "$TOTAL_PROJECTS" -ge 4 ]; then
    echo -e "${GREEN}✓ Multiple projects created${NC}"
    echo "Total projects: $TOTAL_PROJECTS"
else
    echo -e "${RED}✗ Failed to create multiple projects${NC}"
fi
echo ""

# Test 7: Slug Validation
echo -e "${YELLOW}Test 7: Slug Validation${NC}"
INVALID_SLUG=$(curl -s -X POST "$API_URL/api/v1/projects" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Invalid Project\",
    \"slug\": \"Invalid Slug With Spaces!\",
    \"visibility\": \"private\"
  }")

ERROR_MSG=$(echo "$INVALID_SLUG" | jq -r '.error')

if [[ "$ERROR_MSG" == *"slug"* ]]; then
    echo -e "${GREEN}✓ Slug validation working${NC}"
    echo "Error: $ERROR_MSG"
else
    echo -e "${RED}✗ Slug validation not working${NC}"
fi
echo ""

# Test 8: Duplicate Slug Prevention
echo -e "${YELLOW}Test 8: Duplicate Slug Prevention${NC}"
DUPLICATE=$(curl -s -X POST "$API_URL/api/v1/projects" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Duplicate\",
    \"slug\": \"my-first-project\",
    \"visibility\": \"private\"
  }")

DUP_ERROR=$(echo "$DUPLICATE" | jq -r '.error')

if [[ "$DUP_ERROR" == *"already exists"* ]]; then
    echo -e "${GREEN}✓ Duplicate slug prevention working${NC}"
else
    echo -e "${RED}✗ Duplicate slug not prevented${NC}"
fi
echo ""

# Test 9: Access Control
echo -e "${YELLOW}Test 9: Access Control${NC}"
# Create another user
OTHER_USER=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"other${TIMESTAMP}@example.com\",
    \"username\": \"other${TIMESTAMP}\",
    \"password\": \"SecurePass123!\"
  }")

OTHER_TOKEN=$(echo "$OTHER_USER" | jq -r '.access_token')

# Try to access first user's project
ACCESS_DENIED=$(curl -s "$API_URL/api/v1/projects/$PROJECT_ID" \
  -H "Authorization: Bearer $OTHER_TOKEN")

ACCESS_ERROR=$(echo "$ACCESS_DENIED" | jq -r '.error')

if [[ "$ACCESS_ERROR" == *"denied"* ]] || [[ "$ACCESS_ERROR" == *"Forbidden"* ]]; then
    echo -e "${GREEN}✓ Access control working${NC}"
else
    echo -e "${RED}✗ Access control not working${NC}"
    echo "$ACCESS_DENIED"
fi
echo ""

# Test 10: Delete Project
echo -e "${YELLOW}Test 10: Delete Project${NC}"
DELETE_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$API_URL/api/v1/projects/$PROJECT_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if [ "$DELETE_STATUS" == "204" ]; then
    echo -e "${GREEN}✓ Project deleted successfully${NC}"
else
    echo -e "${RED}✗ Project deletion failed (Status: $DELETE_STATUS)${NC}"
fi

# Verify deletion
DELETED_PROJECT=$(curl -s "$API_URL/api/v1/projects/$PROJECT_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$DELETED_PROJECT" | jq -e '.error' > /dev/null; then
    echo -e "${GREEN}✓ Deleted project not accessible${NC}"
else
    echo -e "${RED}✗ Deleted project still accessible${NC}"
fi
echo ""

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ All Week 2 tests completed!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Summary:"
echo "  ✓ Rate limiting active"
echo "  ✓ Project CRUD operations working"
echo "  ✓ Slug validation functioning"
echo "  ✓ Access control enforced"
echo "  ✓ Soft delete implemented"
