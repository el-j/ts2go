# CI/CD Guide for TS2Go

This document describes the Continuous Integration and Continuous Deployment (CI/CD) setup for the TS2Go project.

## Overview

TS2Go uses GitHub Actions for automated testing, building, and releasing. The workflows are designed to support a Git Flow branching strategy with feature branches, develop, and main branches.

## Workflows

### 1. CI - Feature Branch (`ci-feature.yml`)

**Triggers:**
- Push to `feature/**` or `copilot/**` branches
- Pull requests to `develop` or `main`

**Jobs:**
- **lint-and-test:** Runs on Ubuntu
  - Go formatting check (`gofmt`)
  - Go vet static analysis
  - Build CLI
  - Run tests with coverage
  - Upload coverage to Codecov

- **build-cross-platform:** Runs on Ubuntu, macOS, Windows
  - Build CLI for each platform
  - Verify binary works

**Purpose:** Quick feedback for feature development. Ensures code compiles and tests pass on all platforms.

### 2. CI - Develop Branch (`ci-develop.yml`)

**Triggers:**
- Push to `develop`
- Pull requests to `main`

**Jobs:**
- **test:** Runs on Ubuntu
  - All checks from feature CI
  - staticcheck static analysis
  - Full test suite with race detection
  
- **build-all-platforms:** Runs on Ubuntu, macOS, Windows
  - Build for all platforms
  - Test CLI commands

- **integration-tests:** Runs on Ubuntu
  - End-to-end transpilation tests
  - Example project analysis

**Purpose:** Comprehensive testing before merging to main. Catches integration issues and ensures stability.

### 3. CI - Main Branch (`ci-main.yml`)

**Triggers:**
- Push to `main`

**Jobs:**
- **test-and-quality:** Runs on Ubuntu
  - All checks from develop CI
  - Stricter formatting enforcement
  - Coverage reporting with token
  
- **build-all-platforms:** Runs on Ubuntu, macOS, Windows
  - Production-ready builds
  - Multi-architecture support

**Purpose:** Final verification before release. Only stable, tested code reaches main.

### 4. Release (`release.yml`)

**Triggers:**
- Push tags matching `v*.*.*` (e.g., `v1.0.0`)
- Manual workflow dispatch with version input

**Permissions:**
- `contents: write` - To create releases and upload assets

**Jobs:**
- **create-release:**
  - Extract version from tag
  - Generate changelog from CHANGELOG.md
  - Create GitHub release

- **build-and-upload:**
  - Build for multiple platforms:
    - Linux: amd64, arm64
    - macOS: amd64, arm64 (universal)
    - Windows: amd64
  - Create archives (tar.gz for Unix, zip for Windows)
  - Upload as release assets
  - Generate SHA256 checksums

- **publish-docker:**
  - Build multi-arch Docker images
  - Push to GitHub Container Registry (ghcr.io)
  - Tag with version and `latest`

**Purpose:** Automated release process with compiled binaries for all platforms.

### 5. Desktop UI (`desktop-ui.yml`)

**Triggers:**
- Push to `main` or `develop` with changes in `desktop-ui/`
- Pull requests affecting `desktop-ui/`
- Published releases

**Jobs:**
- **test:**
  - Run linter, type checker, and unit tests
  - Build frontend

- **build-tauri:**
  - Build Tauri app for Windows, macOS, Linux
  - Create platform-specific installers
  - Upload to release (if triggered by release)

- **build-web-only:**
  - Build web version (without Tauri)
  - Upload as artifact

**Purpose:** Build and test the desktop UI separately from the CLI.

## Branching Strategy

```
main (production)
  ↑
develop (integration)
  ↑
feature/* (development)
```

### Workflow:
1. Create feature branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/my-feature
   ```

2. Develop and push (triggers feature CI)
   ```bash
   git add .
   git commit -m "feat: add my feature"
   git push origin feature/my-feature
   ```

3. Create PR to `develop` (triggers feature CI)
   - CI runs automatically
   - Review and merge when green

4. Periodically merge `develop` to `main` (triggers main CI)
   ```bash
   git checkout main
   git pull origin main
   git merge develop
   git push origin main
   ```

5. Create release tag (triggers release workflow)
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

## Release Process

### Automated Release (Recommended)

1. **Update version and changelog**
   ```bash
   # Edit CHANGELOG.md - add section for new version
   # Edit cmd/ts2go/main.go - update version constant
   git add CHANGELOG.md cmd/ts2go/main.go
   git commit -m "chore: bump version to v1.0.0"
   git push origin main
   ```

2. **Create and push tag**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. **Wait for workflows**
   - GitHub Actions will automatically:
     - Create the release
     - Build binaries for all platforms
     - Upload release assets
     - Build and push Docker images
     - Build Desktop UI installers

4. **Verify release**
   - Check GitHub Releases page
   - Download and test binaries
   - Verify Docker images

### Manual Release (Emergency)

Use workflow dispatch:
1. Go to Actions → Release workflow
2. Click "Run workflow"
3. Enter version (e.g., `v1.0.0`)
4. Click "Run workflow"

## Docker Images

Docker images are automatically built and pushed on releases.

**Registry:** `ghcr.io/el-j/ts2go`

**Tags:**
- `latest` - Latest release
- `v1.0.0` - Specific version
- `v1.0` - Major.minor
- `v1` - Major only

**Usage:**
```bash
# Pull latest
docker pull ghcr.io/el-j/ts2go:latest

# Run
docker run -v $(pwd):/workspace ghcr.io/el-j/ts2go transpile /workspace/my-project

# Check version
docker run ghcr.io/el-j/ts2go version
```

## Environment Variables & Secrets

### Required Secrets:
- `GITHUB_TOKEN` - Automatically provided by GitHub
- `CODECOV_TOKEN` - (Optional) For coverage reporting on main
- `TAURI_PRIVATE_KEY` - (Optional) For code-signing Tauri apps
- `TAURI_KEY_PASSWORD` - (Optional) Password for Tauri signing key

### Setting up secrets:
1. Go to repository Settings → Secrets and variables → Actions
2. Add repository secrets as needed

## Status Badges

Add to README.md:
```markdown
[![CI - Main](https://github.com/el-j/ts2go/actions/workflows/ci-main.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/ci-main.yml)
[![CI - Develop](https://github.com/el-j/ts2go/actions/workflows/ci-develop.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/ci-develop.yml)
[![Release](https://github.com/el-j/ts2go/actions/workflows/release.yml/badge.svg)](https://github.com/el-j/ts2go/actions/workflows/release.yml)
```

## Troubleshooting

### Tests failing in CI but pass locally
- Check Node.js and Go versions match
- Ensure npm dependencies are installed
- Check for platform-specific issues

### Build fails on specific platform
- Review platform-specific build logs
- Check for platform-specific code paths
- Test locally with similar environment

### Release assets not uploading
- Check GitHub token permissions
- Verify tag format matches `v*.*.*`
- Check artifact paths in workflow

### Docker build fails
- Verify Dockerfile syntax
- Check base image availability
- Test Docker build locally

## Best Practices

1. **Always create PRs** - Don't push directly to develop or main
2. **Wait for CI** - Don't merge until all checks pass
3. **Keep branches updated** - Regularly sync with develop
4. **Write tests** - All new features should have tests
5. **Update docs** - Keep CHANGELOG.md and docs in sync
6. **Test releases** - Download and test release binaries before announcing

## Monitoring

- **Actions tab** - View all workflow runs
- **Releases page** - See all published releases
- **Insights → Actions** - View workflow statistics and trends

## Future Improvements

- [ ] Add performance benchmarks to CI
- [ ] Add automatic dependency updates (Dependabot)
- [ ] Add security scanning (CodeQL, Snyk)
- [ ] Add automatic changelog generation
- [ ] Add release notes generation from commits
- [ ] Add automatic PR labeling
- [ ] Add stale PR/issue management

## Support

For CI/CD issues:
1. Check workflow logs in Actions tab
2. Review this guide for common issues
3. Open an issue with workflow run link
4. Tag with `ci` label
