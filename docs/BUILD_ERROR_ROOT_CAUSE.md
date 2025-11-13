# Build Pipeline Error - Root Cause Analysis & Fix Plan

## 🔴 CRITICAL ERROR IDENTIFIED

### The Problem
```
error: unexpected argument 'build' found
Usage: cargo build [OPTIONS]
```

### Root Cause
The release workflows (alpha, beta, release) use `tauri-apps/tauri-action@v0` which incorrectly constructs the build command:

**What it does:**
```bash
npm run tauri:build build -- --target universal-apple-darwin
```

**What package.json has:**
```json
"tauri:build": "tauri build"
```

**Result:**
```bash
tauri build build -- --target universal-apple-darwin
                ^^^^^ EXTRA "build" argument causes error!
```

This happens because `tauri-action` adds "build" when `tauriScript` is set to `npm run tauri:build`.

## ✅ THE SOLUTION

**Use the Makefile approach that WORKS in `build-desktop.yml`!**

The `build-desktop.yml` workflow uses `make build-desktop` which:
1. Builds CLI correctly
2. Bundles all resources (parser, mappings, node)
3. Runs `cd desktop-ui && npm run tauri:build` (correct command)
4. Creates proper release structure

## 📋 COMPREHENSIVE FIX CHECKLIST

### Phase 1: Understand Workflow Redundancy
- [ ] **Analyze which workflows trigger on which branches**
  - `build-desktop.yml` - Triggers on main, develop, tags, PRs
  - `ci-develop.yml` - CI for develop branch
  - `ci-feature.yml` - CI for feature branches  
  - `ci-main.yml` - CI for main branch
  - `release-alpha.yml` - Releases for feature/copilot branches ⚠️ BROKEN
  - `release-beta.yml` - Releases for develop branch ⚠️ LIKELY BROKEN
  - `release.yml` - Releases for version tags ⚠️ LIKELY BROKEN

- [ ] **Determine which workflows are duplicates**
  - CI workflows (ci-*) = Run tests, no builds
  - build-desktop.yml = Build artifacts, upload to Actions
  - release-* = Build + Create GitHub Releases

**FINDING:** Release workflows try to duplicate build logic but use broken `tauri-action`!

### Phase 2: Fix Strategy

**Option A: Replace tauri-action with Makefile (RECOMMENDED)**
- ✅ Proven to work
- ✅ Single source of truth
- ✅ Consistent with build-desktop.yml
- ✅ No magic/hidden behavior

**Option B: Fix tauri-action configuration**
- ❌ More complex
- ❌ Hides build logic in action
- ❌ Harder to debug
- ❌ Inconsistent with other workflows

**DECISION: Use Option A - Replace with Makefile**

### Phase 3: Implementation Plan

#### Step 1: Fix release-alpha.yml
- [x] Add workspace cleanup
- [x] Add cache improvements  
- [x] Add macOS universal targets
- [ ] **CRITICAL: Replace tauri-action with Makefile commands**
- [ ] Test on feature branch

#### Step 2: Fix release-beta.yml  
- [x] Add workspace cleanup
- [x] Add cache improvements
- [x] Add macOS universal targets
- [ ] **CRITICAL: Replace tauri-action with Makefile commands**
- [ ] Test on develop branch

#### Step 3: Fix release.yml
- [x] Add workspace cleanup
- [x] Add cache improvements
- [x] Add macOS universal targets
- [ ] **CRITICAL: Replace tauri-action with Makefile commands**
- [ ] Test with version tag

#### Step 4: Simplify Workflow Structure (OPTIONAL)
- [ ] Consider merging release workflows into build-desktop.yml
- [ ] Use workflow conditions instead of separate files
- [ ] Document workflow purposes clearly

### Phase 4: The Actual Fix

**Replace this (BROKEN):**
```yaml
- name: Build Tauri app
  uses: tauri-apps/tauri-action@v0
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  with:
    projectPath: desktop-ui
    tauriScript: npm run tauri:build
    tagName: ${{ needs.create-alpha-release.outputs.tag_name }}
    args: ${{ matrix.args }}
```

**With this (WORKS):**
```yaml
- name: Build Desktop App using Makefile
  run: make build-desktop
  env:
    CARGO_INCREMENTAL: 0
    MACOSX_DEPLOYMENT_TARGET: '10.13'

- name: Verify build artifacts
  shell: bash
  run: |
    if [ "${{ matrix.platform }}" == "macos-latest" ]; then
      ls -lah release/stable/macos/
    elif [ "${{ matrix.platform }}" == "windows-latest" ]; then
      ls -lah release/stable/windows/
    else
      ls -lah release/stable/linux/
    fi

- name: Upload to release
  uses: softprops/action-gh-release@v1
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  with:
    tag_name: ${{ needs.create-alpha-release.outputs.tag_name }}
    files: release/stable/${{ matrix.platform == 'macos-latest' && 'macos' || matrix.platform == 'windows-latest' && 'windows' || 'linux' }}/*
```

### Phase 5: Testing Strategy

1. **Local Testing:**
   ```bash
   # Clean everything
   rm -rf release/ desktop-ui/dist desktop-ui/src-tauri/target
   
   # Test Makefile
   make install-desktop
   make build-desktop
   
   # Verify artifacts
   ls -lah release/stable/macos/
   ```

2. **Branch Testing:**
   ```bash
   # Push to feature/copilot branch
   git add .github/workflows/
   git commit -m "Fix: Replace tauri-action with Makefile in release workflows"
   git push
   
   # Watch: https://github.com/el-j/ts2go/actions
   # Should trigger release-alpha.yml
   ```

3. **Validation Checklist:**
   - [ ] All 3 platforms build successfully
   - [ ] Artifacts uploaded to release
   - [ ] No "unexpected argument 'build'" error
   - [ ] Desktop apps work when downloaded

### Phase 6: Root Cause - Why tauri-action Fails

The `tauri-apps/tauri-action@v0` has this behavior:

```yaml
# When you set:
tauriScript: npm run tauri:build
args: --target universal-apple-darwin

# It constructs:
npm run tauri:build build -- --target universal-apple-darwin
#                     ^^^^^ ADDED BY ACTION

# Which expands to:
tauri build build -- --target universal-apple-darwin
#           ^^^^^ DUPLICATE!
```

The action assumes `tauriScript` is something like `npm run` (not `npm run tauri:build`), so it adds "build" itself.

**Correct usage would be:**
```yaml
tauriScript: npm run
args: tauri:build -- --target universal-apple-darwin
```

But this is confusing and fragile. **Better to use Makefile!**

## 📊 Workflow Redundancy Analysis

### Current Workflows
1. **build-desktop.yml** (WORKS ✅)
   - Triggers: main, develop, tags, PRs
   - Purpose: Build artifacts, upload to Actions
   - Method: Makefile (`make build-desktop`)
   - Status: ✅ Working perfectly

2. **ci-develop.yml**
   - Triggers: develop branch
   - Purpose: Run tests/linting
   - Method: Direct commands
   - Status: ✅ Different purpose (CI only)

3. **ci-feature.yml**
   - Triggers: feature/** branches
   - Purpose: Run tests/linting
   - Method: Direct commands
   - Status: ✅ Different purpose (CI only)

4. **ci-main.yml**
   - Triggers: main branch
   - Purpose: Run tests/linting
   - Method: Direct commands
   - Status: ✅ Different purpose (CI only)

5. **release-alpha.yml** (BROKEN ❌)
   - Triggers: feature/**, copilot/** branches
   - Purpose: Build + Create alpha releases
   - Method: tauri-action (BROKEN)
   - Status: ❌ Duplicates build-desktop logic incorrectly

6. **release-beta.yml** (BROKEN ❌)
   - Triggers: develop branch
   - Purpose: Build + Create beta releases
   - Method: tauri-action (BROKEN)
   - Status: ❌ Duplicates build-desktop logic incorrectly

7. **release.yml** (BROKEN ❌)
   - Triggers: version tags (v*)
   - Purpose: Build + Create stable releases
   - Method: tauri-action (BROKEN)
   - Status: ❌ Duplicates build-desktop logic incorrectly

### Redundancy Assessment

**FINDING:** 
- CI workflows (ci-*) are fine - different purpose
- build-desktop.yml is fine - works correctly
- release-* workflows duplicate build logic but use broken method

**RECOMMENDATION:**
Either:
1. Fix release-* to use Makefile (like build-desktop.yml) ✅ RECOMMENDED
2. Or merge functionality into build-desktop.yml with conditions

## 🎯 SUCCESS CRITERIA

After fixes applied:
- [ ] No "unexpected argument 'build'" errors
- [ ] All 3 platforms build successfully
- [ ] Artifacts attached to alpha release
- [ ] Desktop apps download and run correctly
- [ ] Branch can merge without errors
- [ ] Future releases work without manual intervention

## 🚀 NEXT STEPS

1. **Apply the fix** to all 3 release workflows
2. **Test** on feature branch (triggers alpha)
3. **Verify** artifacts created correctly
4. **Merge** branch with confidence
5. **Document** workflow architecture for team

## 📝 SUMMARY

**Problem:** `tauri-apps/tauri-action@v0` adds extra "build" argument  
**Solution:** Replace with `make build-desktop` (proven to work)  
**Benefit:** Consistent, maintainable, debuggable build process  
**Status:** Ready to implement  
**Confidence:** HIGH ✅ (Makefile works in build-desktop.yml)
