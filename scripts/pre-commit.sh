#!/usr/bin/env bash
# TS2Go Pre-commit Hook
# Ensures formatting, linting, and type checking pass before every commit.

set -e

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Running TS2Go Pre-Commit Quality Checks...${NC}"

# 1. Check Go code formatting
echo -n "Checking Go code formatting (gofmt)... "
UNFORMATTED_GO=$(cd go && find . -name "*.go" -not -path "*/fixtures/*" -not -path "*/test-projects/*" | xargs gofmt -s -l)
if [ -n "$UNFORMATTED_GO" ]; then
    echo -e "${RED}FAILED${NC}"
    echo -e "${YELLOW}The following Go files are not formatted:${NC}"
    echo "$UNFORMATTED_GO"
    echo -e "${YELLOW}Run 'make fmt' to format them automatically.${NC}"
    exit 1
fi
echo -e "${GREEN}PASSED${NC}"

# 2. Check Go Vet
echo -n "Checking Go static analysis (go vet)... "
(cd go && go vet ./...)
echo -e "${GREEN}PASSED${NC}"

# 3. Check TypeScript types in packages/ui-shared
if [ -d "packages/ui-shared" ]; then
    echo -n "Checking TypeScript types (vue-tsc)... "
    (cd packages/ui-shared && npm run type-check --silent)
    echo -e "${GREEN}PASSED${NC}"

    echo -n "Checking Frontend formatting (prettier)... "
    (cd packages/ui-shared && npm run format:check --silent)
    echo -e "${GREEN}PASSED${NC}"
fi

echo -e "${GREEN}✅ All pre-commit quality checks passed!${NC}"
exit 0
