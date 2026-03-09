# COMPREHENSIVE FIX PLAN: DMG Build Issues

## Executive Summary

**Root Causes Identified:**

1. **Local Build (228MB .app):**
   - DUPLICATE Node.js binaries: `node` (89MB) + `node.exe` (89MB) = 178MB wasted
   - Stale `node.exe` from Nov 13 in `desktop-ui/src-tauri/bin/`
   - Tauri bundles both due to `resources: ["bin/node*"]` pattern
   - bundle_dmg.sh fails with "Not enough arguments" error (Tauri CLI bug)
   
2. **CI Build (46.9MB DMG):**
   - Likely missing Node.js runtime entirely (~90MB)
   - Or has corrupt/incomplete bundle
   - DMG shows "damaged and can't be opened" error

## Detailed Analysis

### Local .app Bundle Contents

```
Total: 228MB
├── ts2go-cli: 13MB ✓
├── parser/: 23MB ✓  
├── node: 89MB ✓ (correct)
├── node.exe: 89MB ✗ (DUPLICATE, wrong platform)
└── mappings/: 16KB ✓

Actual needed: 13 + 23 + 89 + 0.016 = 125MB
Wasted space: 89MB (node.exe)
```

### Why node.exe Exists on macOS

The Makefile only copies `node`:
```makefile
cp "$$(command -v node)" desktop-ui/src-tauri/bin/node
```

But `desktop-ui/src-tauri/bin/node.exe` exists from a previous build or manual copy (timestamp Nov 13 13:24:34).

Tauri configuration bundles everything matching the pattern:
```json
"resources": [
  "bin/node*"  // ← Matches BOTH node and node.exe
]
```

### bundle_dmg.sh Failure

The script exists and is executable but fails when run:
```bash
$ bash -x ./bundle_dmg.sh
+ [[ -z '' ]]
+ echo 'Not enough arguments. Run '\''create-dmg --help'\'' for help.'
Not enough arguments. Run 'create-dmg --help' for help.
+ exit 1
```

**Why:** Tauri CLI (v2.9.4) generates this script but doesn't invoke it correctly. This is a known issue with Tauri's DMG bundling on some macOS versions.

**Impact:** Build fails after .app creation, so no DMG is produced locally.

## Comprehensive Solution

### Phase 1: Fix Local Build (Immediate)

**Step 1.1: Clean bin/ directory before build**

Modify `Makefile` to remove all binaries before copying:

```makefile
build-desktop: build-cli
	@echo "Cleaning and preparing bundling directory..."
	@rm -rf desktop-ui/src-tauri/bin
	@mkdir -p desktop-ui/src-tauri/bin
	
	@echo "Copying CLI binary to desktop-ui/src-tauri/bin..."
	@cp $(BINARY_NAME) desktop-ui/src-tauri/bin/ts2go-cli$(if $(findstring .exe,$(BINARY_NAME)),.exe,)
	
	@echo "Copying parser dependencies to desktop-ui/src-tauri/bin/parser..."
	@mkdir -p desktop-ui/src-tauri/bin/parser
	@cp -R internal/transpiler/parser/* desktop-ui/src-tauri/bin/parser/
	
	@echo "Copying package mappings to desktop-ui/src-tauri/bin/mappings..."
	@mkdir -p desktop-ui/src-tauri/bin/mappings
	@cp -R mappings/* desktop-ui/src-tauri/bin/mappings/
	
	@echo "Bundling Node.js runtime for $(PLATFORM)..."
	@if command -v node >/dev/null 2>&1; then \
		if [ "$(PLATFORM)" = "windows" ]; then \
			cp "$$(command -v node)" desktop-ui/src-tauri/bin/node.exe; \
			chmod +x desktop-ui/src-tauri/bin/node.exe; \
			echo "  ✓ Node.js runtime bundled for Windows ($$(node --version))"; \
		else \
			cp "$$(command -v node)" desktop-ui/src-tauri/bin/node; \
			chmod +x desktop-ui/src-tauri/bin/node; \
			echo "  ✓ Node.js runtime bundled for $(PLATFORM) ($$(node --version))"; \
		fi; \
	else \
		echo "  ⚠️  Warning: Node.js not found in PATH. App will require system Node.js."; \
	fi
	
	# ... rest of build ...
```

**Step 1.2: Add verification after Tauri build**

```makefile
	@echo "🔍 Verifying build artifacts..."
	@if [ "$(PLATFORM)" = "macos" ]; then \
		APP_PATH="desktop-ui/src-tauri/target/release/bundle/macos/TS2Go Desktop.app"; \
		if [ ! -d "$$APP_PATH" ]; then \
			echo "❌ ERROR: .app bundle not found"; \
			exit 1; \
		fi; \
		\
		APP_SIZE=$$(du -sm "$$APP_PATH" | cut -f1); \
		echo "📦 .app size: $${APP_SIZE}MB"; \
		\
		if [ $$APP_SIZE -lt 100 ]; then \
			echo "❌ ERROR: .app too small ($${APP_SIZE}MB, expected >100MB)"; \
			exit 1; \
		fi; \
		\
		if [ $$APP_SIZE -gt 180 ]; then \
			echo "⚠️  WARNING: .app larger than expected ($${APP_SIZE}MB, expected ~125-150MB)"; \
			echo "  This might indicate duplicate binaries."; \
		fi; \
		\
		if [ ! -f "$$APP_PATH/Contents/Resources/bin/ts2go-cli" ]; then \
			echo "❌ ERROR: ts2go-cli missing from .app"; \
			exit 1; \
		fi; \
		\
		NODE_COUNT=$$(find "$$APP_PATH/Contents/Resources/bin" -name "node*" -type f | grep -v node_modules | wc -l); \
		if [ $$NODE_COUNT -eq 0 ]; then \
			echo "❌ ERROR: Node.js runtime missing from .app"; \
			exit 1; \
		elif [ $$NODE_COUNT -gt 1 ]; then \
			echo "⚠️  WARNING: Multiple Node.js binaries found ($$NODE_COUNT)"; \
			find "$$APP_PATH/Contents/Resources/bin" -name "node*" -type f | grep -v node_modules; \
		fi; \
		\
		echo "✅ .app bundle verified"; \
	fi
```

**Step 1.3: Implement DMG fallback**

```makefile
	@if [ "$(PLATFORM)" = "macos" ]; then \
		rm -rf "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).app"; \
		if [ -d "desktop-ui/src-tauri/target/release/bundle/macos" ]; then \
			cp -R "desktop-ui/src-tauri/target/release/bundle/macos/"*.app "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
		fi; \
		\
		# Try to copy DMG if it exists
		if [ -d "desktop-ui/src-tauri/target/release/bundle/dmg" ]; then \
			DMG_FILE=$$(find "desktop-ui/src-tauri/target/release/bundle/dmg" -name "*.dmg" -size +10M | head -1); \
			if [ -n "$$DMG_FILE" ]; then \
				cp "$$DMG_FILE" "$(RELEASE_DIR)/stable/$(PLATFORM)/" 2>/dev/null || true; \
				echo "✅ DMG copied from Tauri build"; \
			else \
				echo "⚠️  No valid DMG found from Tauri, creating manually..."; \
				hdiutil create -volname "TS2Go Desktop" \
					-srcfolder "$(RELEASE_DIR)/stable/$(PLATFORM)/"*.app \
					-ov -format UDZO "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).dmg"; \
				echo "✅ DMG created manually with hdiutil"; \
			fi; \
		else \
			echo "⚠️  DMG directory not found, creating manually..."; \
			hdiutil create -volname "TS2Go Desktop" \
				-srcfolder "$(RELEASE_DIR)/stable/$(PLATFORM)/"*.app \
				-ov -format UDZO "$(RELEASE_DIR)/stable/$(PLATFORM)/TS2Go-Desktop-v$(VERSION).dmg"; \
			echo "✅ DMG created manually with hdiutil"; \
		fi; \
		\
		echo "✅ Build complete!"; \
		echo "📍 Release location: $(RELEASE_DIR)/stable/$(PLATFORM)/"; \
	elif [ "$(PLATFORM)" = "windows" ]; then \
		# Windows build logic... \
	else \
		# Linux build logic... \
	fi
```

### Phase 2: Fix CI Builds

**Step 2.1: Add pre-build cleanup to workflows**

```yaml
- name: Clean workspace
  run: |
    echo "Cleaning workspace..."
    rm -rf release/ dist/ build/
    rm -rf desktop-ui/dist desktop-ui/src-tauri/target
    rm -rf desktop-ui/src-tauri/bin
    echo "✓ Workspace cleaned"
```

**Step 2.2: Add artifact verification**

```yaml
- name: Verify build artifacts
  shell: bash
  run: |
    echo "Verifying build artifacts..."
    if [ "${{ matrix.platform }}" == "macos-latest" ]; then
      APP_PATH="desktop-ui/src-tauri/target/release/bundle/macos/TS2Go Desktop.app"
      
      if [ ! -d "$APP_PATH" ]; then
        echo "❌ ERROR: .app not found"
        exit 1
      fi
      
      APP_SIZE=$(du -sm "$APP_PATH" | cut -f1)
      echo "📦 .app size: ${APP_SIZE}MB"
      
      if [ $APP_SIZE -lt 100 ]; then
        echo "❌ ERROR: .app too small (${APP_SIZE}MB)"
        exit 1
      fi
      
      if [ $APP_SIZE -gt 180 ]; then
        echo "⚠️  WARNING: .app unusually large (${APP_SIZE}MB, expected ~125-150MB)"
      fi
      
      # Check for Node.js
      if [ ! -f "$APP_PATH/Contents/Resources/bin/node" ]; then
        echo "❌ ERROR: Node.js runtime missing"
        exit 1
      fi
      
      # Check for duplicate Node.js
      NODE_COUNT=$(find "$APP_PATH/Contents/Resources/bin" -name "node*" -type f | grep -v node_modules | wc -l)
      if [ $NODE_COUNT -gt 1 ]; then
        echo "⚠️  WARNING: Multiple Node.js binaries found"
        find "$APP_PATH/Contents/Resources/bin" -name "node*" -type f | grep -v node_modules
      fi
      
      # Verify DMG
      if [ -d "release/stable/macos" ]; then
        DMG_PATH=$(find release/stable/macos -name "*.dmg" | head -1)
        if [ -n "$DMG_PATH" ]; then
          DMG_SIZE=$(du -sm "$DMG_PATH" | cut -f1)
          echo "💿 DMG size: ${DMG_SIZE}MB"
          
          if [ $DMG_SIZE -lt 80 ]; then
            echo "❌ ERROR: DMG too small (${DMG_SIZE}MB)"
            exit 1
          fi
        else
          echo "⚠️  No DMG found in release directory"
        fi
      fi
    fi
```

### Phase 3: Platform-Specific Node.js Bundling

**Option A: Conditional Makefile (Recommended)**

Already shown in Step 1.1 above - use `$(PLATFORM)` to determine which Node.js binary to copy.

**Option B: Platform-Specific Tauri Configs**

Create separate tauri configs per platform, but this is more complex and less maintainable.

## Testing Plan

### Local Testing
1. Clean build: `rm -rf desktop-ui/src-tauri/bin && make build-desktop`
2. Verify .app size: Should be ~125-150MB (not 228MB)
3. Verify only ONE node binary: `find "desktop-ui/src-tauri/target/release/bundle/macos/TS2Go Desktop.app/Contents/Resources/bin" -name "node*" -type f | grep -v node_modules`
4. Verify DMG created: Check `release/stable/macos/`
5. Test .app launches and works

### CI Testing
1. Push fixes to branch
2. Trigger alpha release
3. Verify all artifacts:
   - macOS .app: ~125-150MB
   - macOS DMG: ~100-130MB
   - Windows MSI: ~130-150MB
   - Linux AppImage: ~130-140MB
4. Download and test each artifact

## Implementation Checklist

- [ ] Update Makefile with bin/ cleanup
- [ ] Add platform-specific Node.js copying
- [ ] Add .app verification
- [ ] Implement DMG fallback with hdiutil
- [ ] Update all release workflows with cleanup
- [ ] Add artifact size verification to workflows
- [ ] Test local build
- [ ] Commit and push
- [ ] Trigger CI build
- [ ] Verify CI artifacts
- [ ] Test downloaded artifacts
- [ ] Document fixes

## Expected Outcomes

**Before Fix:**
- Local .app: 228MB (89MB wasted on node.exe)
- Local DMG: FAILS
- CI DMG: 46.9MB (broken/incomplete)

**After Fix:**
- Local .app: ~125-150MB (no duplicates)
- Local DMG: ~100-130MB (created successfully)
- CI DMG: ~100-130MB (functional)
- All artifacts work correctly

## Rollback Plan

If fixes cause issues:
1. Revert Makefile changes
2. Manually remove `desktop-ui/src-tauri/bin/node.exe`
3. Build with old Makefile
4. Create DMG manually with hdiutil

## Long-Term Improvements

1. **Update Tauri:** Upgrade to Tauri v2.1+ which may fix bundle_dmg.sh issues
2. **CI Matrix:** Add size assertions to prevent shipping broken artifacts
3. **E2E Tests:** Automated testing of all platform builds
4. **Documentation:** Add troubleshooting guide for build issues
