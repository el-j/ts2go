#!/bin/bash
# Manual version bumping script
# Usage: ./scripts/bump-version.sh [major|minor|patch]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BUMP_TYPE="${1:-patch}"

if [ ! -f VERSION ]; then
    echo -e "${RED}ERROR: VERSION file not found${NC}"
    exit 1
fi

CURRENT_VERSION=$(cat VERSION)
echo -e "${YELLOW}Current version: ${CURRENT_VERSION}${NC}"

# Parse version
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Bump version based on type
case "$BUMP_TYPE" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo -e "${RED}ERROR: Invalid bump type. Use: major, minor, or patch${NC}"
        exit 1
        ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
echo -e "${GREEN}New version: ${NEW_VERSION}${NC}"

# Update VERSION file
echo "$NEW_VERSION" > VERSION

# Update package.json files
if [ -f desktop-ui/package.json ]; then
    sed -i.bak "s/\"version\": \".*\"/\"version\": \"$NEW_VERSION\"/" desktop-ui/package.json
    rm desktop-ui/package.json.bak
    echo -e "${GREEN}✓ Updated desktop-ui/package.json${NC}"
fi

if [ -f internal/transpiler/parser/package.json ]; then
    sed -i.bak "s/\"version\": \".*\"/\"version\": \"$NEW_VERSION\"/" internal/transpiler/parser/package.json
    rm internal/transpiler/parser/package.json.bak
    echo -e "${GREEN}✓ Updated internal/transpiler/parser/package.json${NC}"
fi

echo ""
echo -e "${GREEN}Version bumped to ${NEW_VERSION}${NC}"
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review the changes: git diff"
echo "2. Commit: git commit -am 'chore: bump version to ${NEW_VERSION}'"
echo "3. Push to main: git push origin main"
echo "4. Tag will be created automatically by CI"
echo ""
echo -e "${YELLOW}Or manually create tag:${NC}"
echo "git tag -a v${NEW_VERSION} -m 'Release v${NEW_VERSION}'"
echo "git push origin v${NEW_VERSION}"
