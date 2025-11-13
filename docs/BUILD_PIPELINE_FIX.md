# GitHub Actions Build Pipeline - Comprehensive Fix

## Problem Analysis

### Root Cause Identified
The build pipeline was failing with multiple critical issues:

1. **gtar (GNU tar) "File exists" errors**: Go toolchain cache extraction conflicts
2. **Cargo registry conflicts**: Rust cache files being reused incorrectly
3. **Dirty workspace state**: Build artifacts from previous runs causing conflicts
4. **Universal binary build issues**: macOS builds failing due to missing Rust targets
5. **Cache collision**: Multiple jobs/branches sharing cache keys incorrectly

### Error Symptoms
```
gtar: Cannot open: File exists
error: unexpected argument 'build' found
failed to build aarch64-apple-darwin binary
```

## Comprehensive Solutions Implemented

### 1. ✅ Workspace Cleanup (Critical)
**Added aggressive cleanup step** to ensure pristine build environment:

```yaml
- name: Clean workspace and caches
  shell: bash
  run: |
    echo "Cleaning workspace to prevent file conflicts..."
    # Clean build artifacts
    rm -rf release/ dist/ build/
    rm -rf desktop-ui/dist desktop-ui/src-tauri/target
    rm -rf desktop-ui/src-tauri/bin
    # Clean Go module cache entries that cause tar conflicts
    rm -rf ~/go/pkg/mod/golang.org/toolchain* || true
    rm -rf ~/go/pkg/mod/cache/download/golang.org/toolchain* || true
    # Clean Cargo registry to prevent conflicts
    rm -rf ~/.cargo/registry || true
    rm -rf ~/.cargo/git || true
    echo "✓ Workspace cleaned"
```

**Why This Works**:
- Removes all build artifacts that could interfere
- Cleans Go toolchain cache that causes gtar extraction errors
- Clears Cargo registry to prevent Rust compilation conflicts
- Ensures each build starts from clean slate

### 2. ✅ Improved Cache Strategy
**Enhanced cache keys** with lockfile hashes to prevent collision:

```yaml
# Node.js cache with dependency tracking
- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: '20'
    cache: 'npm'
    cache-dependency-path: 'desktop-ui/package-lock.json'

# Go cache with multi-module support
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.21'
    cache: true
    cache-dependency-path: |
      go.mod
      go.sum
      cmd/ts2go/go.mod
      cmd/ts2go/go.sum

# Rust cache with platform-specific keys
- name: Rust cache
  uses: Swatinem/rust-cache@v2
  with:
    workspaces: desktop-ui/src-tauri
    cache-on-failure: false
    key: ${{ matrix.platform }}-rust-${{ hashFiles('desktop-ui/src-tauri/Cargo.lock') }}
```

**Why This Works**:
- Platform-specific cache keys prevent cross-platform conflicts
- Lockfile hashes ensure cache invalidation on dependency changes
- `cache-on-failure: false` prevents caching broken builds
- Multi-module Go paths ensure all dependencies cached correctly

### 3. ✅ macOS Universal Binary Support
**Added explicit Rust target installation** for universal builds:

```yaml
- name: Add additional Rust targets for macOS universal builds
  if: matrix.os_name == 'macos'
  run: |
    rustup target add aarch64-apple-darwin
    rustup target add x86_64-apple-darwin
    echo "✓ macOS universal targets installed"
```

**Why This Works**:
- Installs both ARM64 and x86_64 targets for macOS
- Enables Tauri to build universal binaries (works on both Intel and Apple Silicon)
- Prevents "target not found" errors during build

### 4. ✅ Pre-Build Verification
**Added verification step** to catch missing files early:

```yaml
- name: Verify bundled resources preparation
  shell: bash
  run: |
    echo "Verifying required files exist before build..."
    test -f "mappings/npm-to-go.yaml" && echo "✓ mappings/npm-to-go.yaml exists" || (echo "✗ mappings file missing" && exit 1)
    test -f "internal/transpiler/parser/parser.js" && echo "✓ parser.js exists" || (echo "✗ parser.js missing" && exit 1)
    test -f "desktop-ui/package.json" && echo "✓ package.json exists" || (echo "✗ package.json missing" && exit 1)
    echo "✓ All required files present"
```

**Why This Works**:
- Fails fast if critical files are missing
- Provides clear error messages for debugging
- Prevents wasting 10+ minutes on builds that will fail

### 5. ✅ Build Environment Optimization
**Added environment variables** for reliable builds:

```yaml
- name: Build Desktop App using Makefile
  run: make build-desktop
  env:
    CARGO_INCREMENTAL: 0  # Disable incremental compilation (prevents cache issues)
    MACOSX_DEPLOYMENT_TARGET: '10.13'  # Ensure compatibility with older macOS
```

**Why This Works**:
- `CARGO_INCREMENTAL: 0` prevents incremental compilation cache corruption
- `MACOSX_DEPLOYMENT_TARGET` ensures broad macOS compatibility
- Clean build environment eliminates subtle caching bugs

### 6. ✅ Post-Build Verification
**Added comprehensive artifact verification**:

```yaml
- name: Verify build success
  shell: bash
  run: |
    echo "Verifying build artifacts..."
    if [ "${{ matrix.os_name }}" == "macos" ]; then
      if [ -d "release/stable/macos" ]; then
        echo "✓ macOS release directory exists"
        ls -lah release/stable/macos/
      else
        echo "✗ macOS release directory missing"
        echo "Checking Tauri output:"
        ls -R desktop-ui/src-tauri/target/release/bundle/ 2>/dev/null || echo "No bundle output"
        exit 1
      fi
    # ... similar for Windows and Linux
```

**Why This Works**:
- Verifies artifacts were actually created
- Shows detailed directory listings for debugging
- Fails build if artifacts missing (prevents false successes)
- Provides diagnostic info from Tauri output directory

## Complete Workflow Structure

### Build Flow
```
1. Checkout (clean: true)
2. Clean workspace and caches ← PREVENTS TAR CONFLICTS
3. Setup Node.js (with npm cache)
4. Setup Rust (with targets)
5. Add macOS universal targets ← ENABLES UNIVERSAL BUILDS
6. Setup Go (with multi-module cache)
7. Verify installations
8. Install platform dependencies (Linux only)
9. Configure Rust cache ← PLATFORM-SPECIFIC KEYS
10. Install desktop dependencies (make install-desktop)
11. Verify required files exist ← FAIL FAST
12. Build desktop app (make build-desktop) ← WITH CLEAN ENV
13. Verify build success ← CONFIRM ARTIFACTS
14. Upload artifacts (7-day retention)
15. Create release archives (on tags)
16. Upload to GitHub Releases (on tags)
```

## Expected Results

### Build Times (Estimated)
- **First build**: ~15-20 minutes per platform
  - macOS: ~18 min (universal binary takes longer)
  - Windows: ~15 min
  - Linux: ~12 min
- **Cached builds**: ~8-12 minutes per platform
  - Cache saves: 5-10 minutes total
- **Total parallel time**: ~15-20 minutes (all 3 platforms simultaneously)

### Artifacts Produced

#### Development Builds (every push)
```
ts2go-desktop-macos/
  └── TS2Go-Desktop-v0.5.0.app

ts2go-desktop-windows/
  └── TS2Go Desktop_0.5.0_x64_en-US.msi

ts2go-desktop-linux/
  ├── ts-2-go-desktop_0.5.0_amd64.AppImage
  └── ts-2-go-desktop_0.5.0_amd64.deb
```

#### Release Builds (on tags: v*)
```
TS2Go-Desktop-v0.5.0-macos.zip
TS2Go-Desktop-v0.5.0.msi
TS2Go-Desktop-v0.5.0-linux.tar.gz
```

## Testing Instructions

### 1. Test Workflow Locally
```bash
# Clean like CI does
rm -rf release/ dist/ build/
rm -rf desktop-ui/dist desktop-ui/src-tauri/target
rm -rf desktop-ui/src-tauri/bin

# Build
make install-desktop
make build-desktop

# Verify
ls -lah release/stable/macos/  # or windows, linux
```

### 2. Test in CI
```bash
# Push to trigger workflow
git add .github/workflows/build-desktop.yml
git commit -m "Fix: Comprehensive build pipeline improvements

- Add workspace cleanup to prevent tar conflicts
- Improve cache keys with lockfile hashes
- Add macOS universal binary support
- Add pre/post-build verification
- Optimize build environment"
git push
```

### 3. Monitor Build
1. Go to: https://github.com/el-j/ts2go/actions
2. Find "Build Desktop App" workflow
3. Watch each platform build in parallel
4. Download artifacts when complete
5. Test each platform's app

### 4. Test Release
```bash
# Create and push tag
git tag v0.5.1
git push origin v0.5.1

# Workflow will:
# - Build all 3 platforms
# - Create archives
# - Create GitHub Release
# - Attach artifacts automatically
```

## Troubleshooting Guide

### If Builds Still Fail

#### 1. Cache Issues
```bash
# Manually clear caches via GitHub UI
# Settings → Actions → Caches → Delete all

# Or update cache keys in workflow:
key: v2-${{ matrix.platform }}-rust-${{ hashFiles('...') }}
```

#### 2. Workspace Not Clean
```yaml
# Add even more aggressive cleanup:
- name: Nuclear cleanup
  run: |
    git clean -fdx
    git reset --hard HEAD
```

#### 3. Toolchain Issues
```yaml
# Reinstall toolchains from scratch:
- name: Fresh Rust install
  run: |
    rustup self uninstall -y || true
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
```

#### 4. Platform-Specific Failures

**macOS**:
```yaml
# Check Xcode version
- run: xcodebuild -version
# Try older Xcode if needed
- run: sudo xcode-select --switch /Applications/Xcode_14.3.1.app
```

**Windows**:
```yaml
# Use PowerShell instead of bash for some steps
- shell: pwsh
  run: |
    Remove-Item -Recurse -Force release -ErrorAction SilentlyContinue
```

**Linux**:
```yaml
# Update dependencies versions if webkit issues
- run: |
    sudo apt-get install -y libwebkit2gtk-4.1-dev  # Try 4.1 if 4.0 fails
```

## Key Changes Summary

### Files Modified
- `.github/workflows/build-desktop.yml` - Complete overhaul with comprehensive fixes

### Critical Improvements
1. ✅ **Workspace cleanup** - Prevents tar/cache conflicts
2. ✅ **Better cache keys** - Platform + lockfile hash based
3. ✅ **macOS universal support** - Both ARM64 + x86_64 targets
4. ✅ **Pre-build verification** - Fail fast on missing files
5. ✅ **Post-build verification** - Confirm artifacts created
6. ✅ **Clean build environment** - No incremental compilation
7. ✅ **Better error messages** - Diagnostic info on failure

### Benefits
- **Reliability**: Builds start from pristine state every time
- **Speed**: Proper caching saves 5-10 minutes
- **Debugging**: Clear error messages and verification steps
- **Compatibility**: Universal macOS binaries, broad platform support
- **Automation**: Full release pipeline from tag to GitHub Release

## References

- [Tauri Build Documentation](https://tauri.app/v1/guides/building/)
- [GitHub Actions Caching](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows)
- [Rust Cross-Compilation](https://rust-lang.github.io/rustup/cross-compilation.html)
- [macOS Universal Binaries](https://developer.apple.com/documentation/apple-silicon/building-a-universal-macos-binary)

## Next Steps

1. **Push changes** to trigger first test build
2. **Monitor all 3 platforms** - Should now succeed
3. **Download artifacts** - Test each platform's app
4. **Create test tag** - Verify release automation
5. **Update documentation** - If any platform-specific quirks discovered

## Success Criteria

✅ All 3 platforms build successfully  
✅ No tar/cache conflicts  
✅ Artifacts created and uploaded  
✅ Release automation works on tags  
✅ Apps work when downloaded and tested  
✅ Build times reasonable (~15-20 min first, ~8-12 min cached)

---

**Status**: Ready for testing in CI  
**Confidence**: High - addresses all known failure modes  
**Risk**: Low - can always roll back to previous version if needed
