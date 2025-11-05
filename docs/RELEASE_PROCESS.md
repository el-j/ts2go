# Release Process Documentation

## Overview

This document describes the GitHub Actions-based release process for TS2Go, including how artifacts are collected and uploaded to releases.

## Release Types

TS2Go uses three types of releases, each with its own workflow:

### 1. Stable Release (release.yml)
- **Trigger**: Push to `main` branch or manual workflow dispatch with version tag
- **Tag Format**: `v1.0.0` (semantic versioning)
- **Artifacts**: CLI binaries (all platforms), Tauri desktop apps (all platforms), Docker images, checksums
- **Release Type**: Production release

### 2. Beta Release (release-beta.yml)
- **Trigger**: Push to `develop` branch
- **Tag Format**: `v0.1.0-beta.<num>.<sha>` (auto-generated)
- **Artifacts**: CLI binaries (all platforms), Tauri desktop apps (all platforms), checksums
- **Release Type**: Pre-release (beta)

### 3. Alpha Release (release-alpha.yml)
- **Trigger**: Push to `feature/**` or `copilot/**` branches
- **Tag Format**: `v0.0.0-alpha.<branch>.<sha>` (auto-generated)
- **Artifacts**: CLI binaries (3 platforms), Tauri desktop apps (all platforms), checksums
- **Release Type**: Pre-release (alpha) - testing only

## Release Workflow Steps

Each release workflow follows these steps:

### Step 1: Create Release
- Creates a GitHub release with the appropriate tag
- Generates release notes from CHANGELOG.md (stable) or commit history (beta/alpha)
- Sets release as draft/prerelease based on type

**Output**: `version` - the release tag name

### Step 2: Build CLI Binaries
- Builds CLI binaries for multiple platforms in parallel:
  - **Stable/Beta**: linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64
  - **Alpha**: linux-amd64, darwin-arm64, windows-amd64 (reduced set for faster feedback)
- Each build:
  - Compiles the Go binary with version information
  - Creates compressed archives (.tar.gz for Unix, .zip for Windows)
  - Includes README.md and LICENSE files
  - Uploads archive to the release created in Step 1

**Dependencies**: Requires Step 1 (create-release)

### Step 3: Build Desktop Applications
- Builds Tauri desktop applications for all platforms:
  - macOS (universal binary - both Intel and Apple Silicon)
  - Linux (Ubuntu 22.04 - AppImage, deb, rpm)
  - Windows (MSI installer, exe)
- Each build:
  - Builds the CLI binary for bundling
  - Installs frontend and parser dependencies
  - Runs Tauri build with platform-specific arguments
  - **Uploads to the SAME release** using `releaseId` parameter
- Uses Tauri's built-in code signing (when keys provided)

**Dependencies**: Requires Steps 1 and 2 (create-release, build-cli) to complete first

### Step 4: Generate Checksums
- Downloads all CLI artifacts from the release
- Generates SHA256 checksums for all .tar.gz and .zip files
- Uploads checksums.txt to the release

**Dependencies**: Requires all previous steps (create-release, build-cli, build-desktop)

### Step 5: Publish Docker Image (Stable Only)
- Builds and pushes Docker image to GitHub Container Registry
- Tags: version number, major.minor, major, and latest
- Multi-platform: linux/amd64, linux/arm64

**Dependencies**: Requires Steps 1-3

## Key Configuration Details

### Tauri Action Configuration
The Tauri action is configured to upload to the existing release:

```yaml
- uses: tauri-apps/tauri-action@v0
  with:
    projectPath: desktop-ui
    tagName: ${{ needs.create-release.outputs.version }}
    releaseId: ${{ needs.create-release.outputs.version }}  # Uses tag as release ID
    args: ${{ matrix.args }}
```

**Important**: The `releaseId` parameter ensures Tauri uploads to the existing release created in Step 1, rather than creating a new release.

### Job Dependencies
Jobs are structured to ensure proper sequencing:

```
create-release
    ↓
build-cli (parallel across platforms)
    ↓
build-desktop (parallel across platforms)
    ↓
generate-checksums
```

### Platform Matrix

**CLI Builds (Stable/Beta)**:
- linux-amd64
- linux-arm64
- darwin-amd64
- darwin-arm64
- windows-amd64

**CLI Builds (Alpha)** - reduced for speed:
- linux-amd64
- darwin-arm64
- windows-amd64

**Desktop Builds (All)**:
- macos-latest (universal)
- ubuntu-22.04
- windows-latest

## Expected Artifacts

For a complete stable release, you should see these artifacts:

### CLI Binaries
1. `ts2go-v1.0.0-linux-amd64.tar.gz`
2. `ts2go-v1.0.0-linux-arm64.tar.gz`
3. `ts2go-v1.0.0-darwin-amd64.tar.gz`
4. `ts2go-v1.0.0-darwin-arm64.tar.gz`
5. `ts2go-v1.0.0-windows-amd64.zip`

### Desktop Applications
6. `TS2Go-Desktop_<version>_universal.dmg` (macOS)
7. `TS2Go-Desktop_<version>_amd64.AppImage` (Linux)
8. `TS2Go-Desktop_<version>_amd64.deb` (Linux)
9. `TS2Go-Desktop_<version>_amd64.rpm` (Linux)
10. `TS2Go-Desktop_<version>_x64_en-US.msi` (Windows)
11. `TS2Go-Desktop_<version>_x64-setup.exe` (Windows)

### Other Files
12. `checksums.txt` - SHA256 checksums for CLI binaries

## Troubleshooting

### Missing Artifacts

**Problem**: Some artifacts missing from release
**Solution**: Check workflow run logs for each matrix job. Each platform builds independently.

**Problem**: Desktop apps not appearing
**Solution**: 
1. Check that CLI builds completed successfully (desktop builds depend on them)
2. Verify Tauri action is using `releaseId` parameter
3. Check for platform-specific build failures

### Checksum Issues

**Problem**: Checksums file missing or incomplete
**Solution**: 
1. Ensure generate-checksums job ran after all builds
2. Check that robinraju/release-downloader successfully downloaded artifacts
3. Verify artifacts have the expected naming pattern (`ts2go-*.tar.gz` or `ts2go-*.zip`)

### Release Not Created

**Problem**: Release not created at all
**Solution**:
1. Check if commit message contains `[skip-release]` (beta/alpha only)
2. Verify branch name matches trigger conditions
3. For stable releases, ensure VERSION file exists or manual version input provided

## Verification Checklist

Use this checklist to verify a successful release:

- [ ] Release created with correct tag name
- [ ] Release notes generated and attached
- [ ] Release marked as prerelease (beta/alpha) or stable correctly
- [ ] All 5 CLI binary archives present (stable/beta) or 3 (alpha)
- [ ] LICENSE file included in archives
- [ ] All desktop applications built and uploaded:
  - [ ] macOS .dmg
  - [ ] Linux .AppImage, .deb, .rpm
  - [ ] Windows .msi and .exe
- [ ] Checksums.txt file present with all CLI checksums
- [ ] Docker image published (stable only)
- [ ] All workflow jobs completed successfully

## Manual Release Trigger

To manually trigger a stable release:

1. Go to Actions → Release Stable (Main Branch)
2. Click "Run workflow"
3. Enter version tag (e.g., `v1.2.3`)
4. Click "Run workflow"

The workflow will create the release and build all artifacts.

## Security

### Code Signing
- Tauri builds can use code signing when `TAURI_PRIVATE_KEY` and `TAURI_KEY_PASSWORD` secrets are configured
- Stable releases use code signing, beta/alpha may not

### Secrets Required
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions
- `TAURI_PRIVATE_KEY`: (Optional) For code signing desktop apps
- `TAURI_KEY_PASSWORD`: (Optional) For code signing desktop apps

## Version Management

### Stable Releases
Version is determined by:
1. Manual input (workflow_dispatch)
2. VERSION file (auto-version workflow)
3. Git tag name (tag push)

### Beta Releases
Version auto-generated: `v0.1.0-beta.<number>.<short-sha>`

### Alpha Releases
Version auto-generated: `v0.0.0-alpha.<branch-name>.<short-sha>`

## Future Improvements

Potential enhancements to consider:

1. **Artifact attestation**: Add SLSA provenance for supply chain security
2. **Release notes automation**: Extract from conventional commits
3. **Homebrew formula**: Auto-update Homebrew tap on stable release
4. **APT/YUM repositories**: Publish to package repositories
5. **Update checks**: Built-in update notification in CLI/desktop apps
6. **Performance metrics**: Track build times and artifact sizes
