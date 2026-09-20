#!/usr/bin/env bash
# TS2Go Mutation Testing Runner
# Executes Gremlins mutation testing on target Go packages to verify test efficacy.

set -e

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT/go"

# Check if gremlins is available
GREMLINS_BIN="$(which gremlins 2>/dev/null || echo "$HOME/go/bin/gremlins")"

if [ ! -x "$GREMLINS_BIN" ]; then
    echo -e "${YELLOW}⚠️  gremlins not found in PATH or ~/go/bin. Installing latest gremlins...${NC}"
    go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
    GREMLINS_BIN="$(which gremlins 2>/dev/null || echo "$HOME/go/bin/gremlins")"
    if [ ! -x "$GREMLINS_BIN" ]; then
        echo -e "${RED}❌ Failed to install gremlins.${NC}"
        exit 1
    fi
fi

echo -e "${BLUE}🧬 Running Mutation Testing Suite with Gremlins...${NC}"

# Target packages: defaults to runtime packages or user arguments
TARGETS="${*:-./runtime/path ./runtime/fs ./runtime/console}"

echo -e "${YELLOW}Target packages: ${TARGETS}${NC}"

for TARGET in $TARGETS; do
    echo -e "\n${BLUE}👉 Testing mutations in ${TARGET}...${NC}"
    "$GREMLINS_BIN" unleash "$TARGET" --threshold-efficacy=80 --timeout-coefficient=20
done

echo -e "\n${GREEN}✅ Mutation testing completed successfully! All mutant quality thresholds met.${NC}"
exit 0
