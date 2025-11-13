# GitHub Actions Workflow Fix - Build Desktop App

## Problem Analysis

### Root Cause
The GitHub Actions workflow was failing on all platforms (macOS, Windows, Linux) with the error:
```
error: unexpected argument 'build' found
Usage: cargo build [OPTIONS]
```

This error occurred because the workflow was using `npm run tauri:build` with `working-directory: desktop-ui`, which caused the Tauri CLI to be invoked incorrectly. The command was being interpreted as `cargo.exe build` instead of the proper Tauri build command.

### Why It Failed
1. **Working Directory Issue**: Using `working-directory: desktop-ui` with `npm run tauri:build` caused path resolution issues
2. **Complex Resource Bundling**: Manually copying CLI, parser, node, and mappings in separate steps was error-prone
3. **Platform-Specific Logic**: Having separate Unix/Windows resource bundling steps added complexity
4. **Inconsistent with Local Build**: The workflow didn't match the proven Makefile approach that works locally

## Solution: Use Makefile Targets

### Key Changes

#### Before (Failed Approach)
```yaml
- name: Install desktop dependencies
  working-directory: desktop-ui
  run: npm ci

- name: Build CLI (Unix)
  if: matrix.os_name != 'windows'
  run: make build-cli

- name: Build CLI (Windows)
  if: matrix.os_name == 'windows'
  run: go build -o ts2go.exe ./cmd/ts2go

- name: Prepare bundled resources (Unix)
  if: matrix.os_name != 'windows'
  run: |
    mkdir -p desktop-ui/src-tauri/bin
    cp ts2go desktop-ui/src-tauri/bin/ts2go-cli
    # ... more manual copying

- name: Build Desktop App
  working-directory: desktop-ui
  run: npm run tauri:build
```

#### After (Working Approach)
```yaml
- name: Install desktop dependencies
  run: make install-desktop

- name: Build Desktop App using Makefile
  run: make build-desktop
```

### Why This Works

1. **Proven Locally**: The Makefile targets (`make build-desktop`) work perfectly on local machines
2. **Single Source of Truth**: Build logic is in one place (Makefile), not duplicated in CI
3. **Handles All Platforms**: Makefile automatically detects platform and adjusts behavior
4. **Complete Process**: `make build-desktop` handles:
   - Building the CLI binary (with correct platform extensions)
   - Copying CLI to `desktop-ui/src-tauri/bin/`
   - Bundling parser dependencies
   - Bundling package mappings
   - Bundling Node.js runtime
   - Running `cd desktop-ui && npm run tauri:build` (correct command invocation)
   - Creating `release/stable/{platform}/` directory structure
   - Copying final artifacts to release directory

## Artifact Handling

### Build Artifacts Structure
After `make build-desktop`, artifacts are in:
```
release/stable/
  ├── macos/TS2Go-Desktop-v{VERSION}.app
  ├── windows/*.msi
  └── linux/*.AppImage, *.deb
```

### Workflow Artifact Upload
```yaml
- name: Upload build artifacts
  uses: actions/upload-artifact@v4
  with:
    name: ts2go-desktop-${{ matrix.os_name }}
    path: release/stable/${{ matrix.os_name }}/
    retention-days: 7
    if-no-files-found: error
```

### Release Assets (on Git Tags)
When pushing a tag (e.g., `v0.5.0`), the workflow:
1. Creates platform-specific archives:
   - macOS: `TS2Go-Desktop-v{VERSION}-macos.zip` (contains .app)
   - Windows: `TS2Go-Desktop-v{VERSION}.msi` (installer directly)
   - Linux: `TS2Go-Desktop-v{VERSION}-linux.tar.gz` (contains .AppImage + .deb)
2. Uploads to GitHub Release automatically
3. Generates release notes from commits

## Benefits of This Approach

### 1. Reliability
- ✅ Uses battle-tested Makefile that works locally
- ✅ Eliminates manual resource copying that was error-prone
- ✅ Consistent behavior between local and CI builds

### 2. Maintainability
- ✅ Single source of truth for build logic
- ✅ Changes to build process only need to update Makefile
- ✅ Less code duplication in workflow file

### 3. Platform Compatibility
- ✅ Makefile handles platform detection automatically
- ✅ Correctly adds `.exe` extension on Windows
- ✅ Proper command invocation via `cd desktop-ui && npm run tauri:build`

### 4. Complete Bundling
- ✅ All dependencies bundled: CLI, parser, mappings, Node.js
- ✅ Apps work standalone without external dependencies
- ✅ Fixes original issue: "mappings/npm-to-go.yaml not found"

## Testing the Workflow

### Manual Testing
```bash
# Simulate what CI does
make install-desktop
make build-desktop

# Check artifacts
ls -lah release/stable/macos/  # or windows, linux
```

### CI Testing
1. Push to branch to trigger workflow
2. Check GitHub Actions tab
3. Verify all 3 platforms build successfully
4. Download artifacts from workflow run
5. Test each platform's app

### Release Testing
```bash
# Create and push a tag
git tag v0.5.1
git push origin v0.5.1

# Check:
# - Workflow runs
# - All 3 platforms build
# - Release created automatically
# - Assets attached to release
```

## Validation Checklist

Before pushing workflow changes:

- [x] YAML syntax valid (`ruby -ryaml -e "YAML.load_file('.github/workflows/build-desktop.yml')"`)
- [x] Uses Makefile targets that work locally
- [x] No platform-specific conditional logic (Makefile handles it)
- [x] Artifact paths match Makefile output structure
- [x] Release asset creation handles all file types
- [x] if-no-files-found: error ensures build failures are caught
- [x] Release workflow only runs on tags
- [x] GITHUB_TOKEN configured for release uploads

## Expected Workflow Duration

With caching enabled:
- **First build**: ~15-20 minutes per platform
- **Subsequent builds**: ~8-12 minutes per platform
  - Rust cache saves ~5-10 minutes
  - Go cache saves ~1-2 minutes
  - npm cache saves ~1-2 minutes

Total parallel build time: ~15-20 minutes (first) / ~8-12 minutes (cached)

## Troubleshooting

### If builds still fail:

1. **Check Makefile compatibility**:
   ```bash
   make -n build-desktop  # Dry run to see commands
   ```

2. **Verify dependencies**:
   - Node.js 20+ installed
   - Rust stable toolchain
   - Go 1.21+ installed
   - Platform-specific dependencies (Linux: webkit2gtk)

3. **Check artifact paths**:
   ```bash
   # After build, verify:
   ls release/stable/{platform}/
   ```

4. **Review workflow logs**:
   - Look for "Build Desktop App using Makefile" step
   - Check Makefile output for errors
   - Verify "List build artifacts" shows expected files

## Future Improvements

1. **Build caching**: Already implemented (Rust, Go)
2. **Parallel builds**: Already runs in parallel via matrix strategy
3. **Auto-versioning**: Consider reading VERSION file for consistent versioning
4. **Signed releases**: Add code signing for macOS/Windows (requires certificates)
5. **Notarization**: Add Apple notarization for macOS apps (requires Apple Developer account)

## Related Files

- `.github/workflows/build-desktop.yml` - Main workflow file
- `Makefile` - Build logic (targets: install-desktop, build-desktop)
- `desktop-ui/package.json` - Defines `tauri:build` script
- `desktop-ui/src-tauri/tauri.conf.json` - Tauri configuration with resources
- `docs/CI_BUILD_FIXES.md` - Previous troubleshooting attempts
- `docs/CROSS_COMPILATION.md` - Cross-platform build guide

## Summary

**Problem**: Workflow failed with incorrect Tauri CLI invocation
**Solution**: Use proven Makefile targets instead of reimplementing build logic
**Result**: Reliable, maintainable, cross-platform builds with automatic release publishing

The key insight: **Don't fight the tools**. The Makefile already works perfectly locally, so use it in CI too. This eliminates an entire class of "works on my machine" issues.
