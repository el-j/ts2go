# Version Management System

## Overview

TS2Go uses an automated version management system that requires **zero manual intervention** for version bumping and releases.

## How It Works

### Automatic Versioning (Main Branch)

When you push to the `main` branch, the system automatically:

1. **Determines version bump type** from your commit message (following [Conventional Commits](https://www.conventionalcommits.org/))
2. **Updates the VERSION file** with the new version number
3. **Creates and pushes a git tag** (e.g., `v1.2.3`)
4. **Triggers the release workflow** which builds and publishes all artifacts

### Version Bump Rules

The system follows semantic versioning (MAJOR.MINOR.PATCH):

| Commit Message Pattern | Bump Type | Example |
|----------------------|-----------|---------|
| `BREAKING CHANGE:` or `feat!:` or `fix!:` | **Major** | 1.0.0 → 2.0.0 |
| `feat:` or `feature:` | **Minor** | 1.0.0 → 1.1.0 |
| Any other commit (fix, docs, etc.) | **Patch** | 1.0.0 → 1.0.1 |

### Examples

```bash
# Major version bump (breaking change)
git commit -m "feat!: redesign API with new architecture"
git push origin main
# Result: 1.5.3 → 2.0.0

# Minor version bump (new feature)
git commit -m "feat: add dark mode support"
git push origin main
# Result: 1.5.3 → 1.6.0

# Patch version bump (bug fix, docs, etc.)
git commit -m "fix: correct typo in settings panel"
git push origin main
# Result: 1.5.3 → 1.5.4
```

## Release Levels

### 1. Alpha Releases (Feature Branches)

- **Triggers:** Automatic on push to `feature/**` or `copilot/**` branches
- **Version format:** `v0.0.0-alpha.{branch-name}.{commit-sha}`
- **Example:** `v0.0.0-alpha.feature-new-ui.a1b2c3d`
- **Purpose:** Testing unreleased features
- **Marked as:** Prerelease

### 2. Beta Releases (Develop Branch)

- **Triggers:** Automatic on push to `develop` branch
- **Version format:** `v0.1.0-beta.{number}.{commit-sha}`
- **Example:** `v0.1.0-beta.5.a1b2c3d`
- **Purpose:** Integration testing before stable release
- **Marked as:** Prerelease
- **Auto-increments:** Beta number increases automatically

### 3. Stable Releases (Main Branch)

- **Triggers:** Automatic on push to `main` branch
- **Version format:** `v{major}.{minor}.{patch}`
- **Example:** `v1.2.3`
- **Purpose:** Production releases
- **Marked as:** Full release
- **Version source:** VERSION file

## Manual Version Bumping (Optional)

If you need to manually bump the version before committing:

```bash
# Bump patch version (1.0.0 → 1.0.1)
./scripts/bump-version.sh patch

# Bump minor version (1.0.0 → 1.1.0)
./scripts/bump-version.sh minor

# Bump major version (1.0.0 → 2.0.0)
./scripts/bump-version.sh major
```

The script will:
- Update the `VERSION` file
- Update `desktop-ui/package.json`
- Update `internal/transpiler/parser/package.json`
- Show you the next steps

## Skipping Releases

To push changes without triggering a release:

```bash
git commit -m "docs: update README [skip-release]"
git push origin main
```

Or use `[no-release]` in your commit message.

## Workflow Files

- **`auto-version.yml`** - Automatic version bumping and tagging on main
- **`release.yml`** - Stable releases triggered by tags
- **`release-beta.yml`** - Beta releases on develop
- **`release-alpha.yml`** - Alpha releases on feature branches

## Version File

The `VERSION` file at the repository root contains the current version:

```
0.5.0
```

This file is automatically updated by the CI system but can also be manually edited if needed.

## Best Practices

### For Regular Development

1. Work on feature branches (`feature/my-feature`)
2. Alpha releases are created automatically for testing
3. Merge to `develop` when ready
4. Beta releases are created automatically for integration testing
5. Merge to `main` when stable
6. Use conventional commit messages for automatic version bumping
7. Stable release is created automatically with proper version

### For Hotfixes

1. Create branch from `main`
2. Fix the issue
3. Commit with appropriate message: `fix: critical security patch`
4. Push to `main`
5. System automatically bumps patch version and releases

### For Major Releases

1. Merge all features to `develop`
2. Test with beta releases
3. When ready, commit with breaking change message:
   ```bash
   git commit -m "feat!: redesign core architecture

   BREAKING CHANGE: This release changes the API interface"
   ```
4. Merge to `main`
5. System automatically bumps major version (e.g., 1.5.3 → 2.0.0)

## Troubleshooting

### Tag Already Exists

If a tag already exists for the version, the workflow will skip creating it again. This prevents duplicate releases.

### VERSION File Out of Sync

If the VERSION file gets out of sync with git tags:

1. Check latest tag: `git describe --tags --abbrev=0`
2. Update VERSION file manually
3. Commit: `git commit -am "chore: sync VERSION file [skip-release]"`
4. Push: `git push origin main`

### Manual Tag Creation

If you prefer to create tags manually:

```bash
# After updating VERSION file
./scripts/bump-version.sh minor

# Create tag manually
git tag -a v1.6.0 -m "Release v1.6.0"
git push origin v1.6.0

# The release.yml workflow will still trigger
```

## Benefits

✅ **Zero manual intervention** - Just push to main  
✅ **No version conflicts** - System handles everything  
✅ **Semantic versioning** - Follows industry standards  
✅ **Conventional commits** - Encourages good commit messages  
✅ **Multi-level releases** - Alpha, Beta, Stable automatically  
✅ **Easy to override** - Manual script available if needed  
✅ **Skip when needed** - Use [skip-release] in commit message  

## Summary

**You don't need to worry about versions!** Just:

1. Write good commit messages (use `feat:` for features, `fix:` for fixes)
2. Push to `main` when ready to release
3. The system handles everything else automatically

That's it! 🎉
