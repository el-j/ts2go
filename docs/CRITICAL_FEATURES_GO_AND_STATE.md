# CRITICAL FEATURES: Go Binary Bundling & Transpilation State Persistence

## Problem Statement

### Issue 0: Transpiled Go Code Not Formatted (NEW - CRITICAL)

**Problem:**
> Transpiled Go code is not automatically formatted with `go fmt` after transpilation completes. This violates Go community standards and produces poorly formatted, potentially unreadable code.

**Current Behavior:**
- `auto_transpile_project()` in main.rs calls CLI transpiler
- CLI generates Go code successfully
- Function counts transpiled files and returns success
- **No formatting step** - code written directly to disk unformatted

**Why This Matters:**
1. **Go Standard Practice:** All Go code should be `gofmt`'d - it's not optional
2. **Code Readability:** Unformatted Go code is harder to read and review
3. **Syntax Validation:** `go fmt` catches basic syntax errors early
4. **Version Control:** Inconsistent formatting creates noisy diffs
5. **Professional Quality:** Properly formatted code demonstrates quality tooling

**Affected Code:**
- `desktop-ui/src-tauri/src/main.rs` - Line 172: `auto_transpile_project()` function
- `pkg/cli/transpile.go` - TranspileCommand doesn't format output
- All transpiled Go files written without formatting pass

**Expected Behavior:**
After successful transpilation:
1. Count transpiled .go files
2. Run `go fmt ./...` in output directory
3. Report formatting results (files formatted, warnings if any)
4. Handle errors gracefully (warnings, not failures)
5. Return enhanced JSON with formatting status

**Implementation Notes:**
- Can use system Go initially (Phase 0a)
- Will use bundled Go once available (Phase 0b - after Phase 3)
- Should NOT fail entire transpilation if formatting fails
- Should capture and report formatting errors/warnings
- Requires Go in PATH or bundled Go binary

**Code Location:**
```rust
// desktop-ui/src-tauri/src/main.rs:172
#[tauri::command]
async fn auto_transpile_project(path: String) -> Result<serde_json::Value, String> {
    // ... transpilation code ...
    
    if result.status.success() {
        let mut go_files = Vec::new();
        find_go_files(Path::new(&output_dir), &mut go_files)?;
        
        // ❌ MISSING: Format all Go files
        // Should add:
        // let format_result = format_go_files(&output_dir)?;
        
        Ok(serde_json::json!({
            "success": true,
            "output_dir": output_dir,
            "files_transpiled": go_files.len(),
            // Missing: "files_formatted": N, "format_warnings": []
        }))
    }
}
```

### Issue 1: Missing Go Binary for Build/Test/Run Operations

**User Report:**
> "When I transpile a complete project in the UI, the buttons for build and test and run are not able to work as they miss the go binary to execute. I looked into the settings tab, and does not find any option to either set a path for my local go version, or use an internal go version."

**Current Behavior:**
- After transpiling a project, Build/Test/Run buttons appear in ProjectView
- Clicking these buttons calls Rust commands that execute `Command::new("go")`
- If user doesn't have Go installed or Go is not in PATH, operations fail silently or with errors
- No way to configure Go binary path in settings
- No bundled Go compiler with the app

**Affected Commands in main.rs:**
1. `run_go_file()` - Line 277: `Command::new("go").arg("run")`
2. `run_go_project()` - Line 341: `Command::new("go")`
3. `build_go_file()` - Line 384: `Command::new("go").arg("build")`  
4. `build_go_project()` - Line 478: `Command::new("go")`
5. `test_go_file()` - Line 551: `Command::new("go").arg("test")`
6. `test_go_project()` - Line 597: `Command::new("go")`

**Impact:**
- Users cannot use Build/Test/Run features without manually installing Go
- No clear error message explaining what's wrong
- Desktop app should be self-contained and not require external dependencies

### Issue 2: Lost Transpilation State When Switching Views

**User Report:**
> "When I transpiled a complete project and switch the main view and come back to the project-view, it seems to be missing the compilation state and so the buttons for test and build and run. The app always should know the state of what the user did and never forget about an important thing, like 'transpile a project' and it's outputs."

**Current Behavior:**
- User transpiles project in ProjectView
- `transpileStore.currentResult` is set with `output_dir`, `files_transpiled`, etc.
- Build/Test/Run buttons appear conditionally: `v-if="transpileStore.currentResult?.output_dir"`
- User navigates to Settings or Examples view
- User returns to ProjectView
- `transpileStore.currentResult` is still in memory (within same session)
- **BUT** if app is closed or page refreshed, state is lost
- Buttons disappear even though transpiled Go code still exists on disk

**Impact:**
- Poor user experience - app "forgets" important work
- User must re-transpile just to access Build/Test/Run buttons
- No persistence across app restarts
- No way to track which projects have been transpiled

## Proposed Solutions

### Solution 1: Bundle Go Compiler with Desktop App

**Approach A: Full Go Installation (Recommended for Beta)**

Bundle complete Go distribution with desktop app, similar to how Node.js is bundled.

**Implementation:**

1. **Download Official Go Binaries**
   - Go 1.22.9 or 1.23.x (latest stable)
   - Platforms:
     - `darwin-arm64` (macOS Apple Silicon) - ~430MB extracted
     - `darwin-amd64` (macOS Intel) - ~450MB extracted
     - `windows-amd64` (Windows x64) - ~400MB extracted
     - `linux-amd64` (Linux x64) - ~430MB extracted

2. **Makefile Integration**
   ```makefile
   # Add to build-desktop target
   @echo "Bundling Go compiler for $(PLATFORM)..."
   @if [ -n "$(GO_BUNDLE_PATH)" ]; then \
       mkdir -p desktop-ui/src-tauri/bin/go; \
       cp -R "$(GO_BUNDLE_PATH)"/* desktop-ui/src-tauri/bin/go/; \
       chmod +x desktop-ui/src-tauri/bin/go/bin/go; \
       echo "  ✓ Go $(shell $(GO_BUNDLE_PATH)/bin/go version) bundled"; \
   else \
       echo "  ⚠️  GO_BUNDLE_PATH not set, skipping Go bundling"; \
   fi
   ```

3. **Tauri Configuration**
   ```json
   // tauri.conf.json
   "resources": [
     "bin/ts2go-cli*",
     "bin/parser",
     "bin/node*",
     "bin/mappings",
     "bin/go/**"  // NEW: Bundle entire Go installation
   ]
   ```

4. **Rust Helper Function**
   ```rust
   // desktop-ui/src-tauri/src/main.rs
   
   fn get_go_binary_path(custom_path: Option<String>) -> Result<PathBuf, String> {
       // Check for custom path from settings
       if let Some(path) = custom_path {
           let go_path = PathBuf::from(path);
           if go_path.exists() {
               return Ok(go_path);
           }
           return Err(format!("Custom Go path not found: {}", path));
       }
       
       // In release, try bundled Go first
       #[cfg(not(debug_assertions))]
       {
           if let Ok(exe_path) = std::env::current_exe() {
               if let Some(contents_dir) = exe_path.parent().and_then(|p| p.parent()) {
                   #[cfg(target_os = "macos")]
                   {
                       let go_path = contents_dir.join("Resources/bin/go/bin/go");
                       if go_path.exists() {
                           return Ok(go_path);
                       }
                   }
                   
                   #[cfg(target_os = "windows")]
                   {
                       let go_path = contents_dir.join("bin/go/bin/go.exe");
                       if go_path.exists() {
                           return Ok(go_path);
                       }
                   }
                   
                   #[cfg(target_os = "linux")]
                   {
                       let go_path = contents_dir.join("bin/go/bin/go");
                       if go_path.exists() {
                           return Ok(go_path);
                       }
                   }
               }
           }
       }
       
       // Fallback to system Go
       #[cfg(target_os = "windows")]
       { Ok(PathBuf::from("go.exe")) }
       #[cfg(not(target_os = "windows"))]
       { Ok(PathBuf::from("go")) }
   }
   
   // New command to detect Go installation
   #[tauri::command]
   async fn detect_go_installation(custom_path: Option<String>) -> Result<serde_json::Value, String> {
       let go_path = get_go_binary_path(custom_path)?;
       
       let output = Command::new(&go_path)
           .arg("version")
           .output()
           .map_err(|e| format!("Failed to execute Go: {}. Make sure Go is installed.", e))?;
       
       if output.status.success() {
           let version_str = String::from_utf8_lossy(&output.stdout);
           Ok(serde_json::json!({
               "found": true,
               "version": version_str.trim(),
               "path": go_path.to_string_lossy(),
               "valid": true
           }))
       } else {
           Ok(serde_json::json!({
               "found": false,
               "valid": false,
               "error": "Go binary found but version check failed"
           }))
       }
   }
   ```

5. **Update All Command::new("go") Calls**
   ```rust
   // Before:
   let output = Command::new("go")
       .arg("run")
       .arg(&temp_file)
       .output()
   
   // After:
   let go_path = get_go_binary_path(None)?;  // Get from settings if needed
   let output = Command::new(&go_path)
       .arg("run")
       .arg(&temp_file)
       .output()
   ```

**App Size Impact:**
- Current .app: ~139MB
- With bundled Go: ~139MB + ~430MB = **~570MB**
- DMG: ~54MB → **~400MB** (compressed)

**Approach B: Minimal Go Runtime (Alternative)**

Bundle only essential Go binaries without full stdlib source:
- `go` binary
- `compile`, `link`, `asm` tools from `pkg/tool/`
- Compiled standard library from `pkg/`
- Skip: `src/`, `doc/`, `test/`, `misc/`

**Size:** ~200MB instead of 430MB  
**Risk:** May not work for all Go operations

**Recommended:** Start with Approach A (full Go), optimize later if size is issue.

### Solution 2: Add Go Configuration to Settings

**UI Changes (SettingsView.vue):**

Add new "Go Configuration" section under "Project" tab:

```vue
<!-- Go Configuration -->
<div class="py-4 border-b border-gray-200 dark:border-gray-700">
  <h3 class="text-lg font-medium text-gray-800 dark:text-gray-200 mb-4">
    Go Configuration
  </h3>
  
  <!-- Go Binary Source -->
  <div class="mb-4">
    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
      Go Binary Source
    </label>
    <Dropdown 
      v-model="settingsStore.settings.goBinarySource" 
      :options="goBinarySourceOptions" 
      optionLabel="label" 
      optionValue="value"
      class="w-full"
    />
  </div>
  
  <!-- Custom Go Path (only if 'custom' selected) -->
  <div v-if="settingsStore.settings.goBinarySource === 'custom'" class="mb-4">
    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
      Custom Go Binary Path
    </label>
    <div class="flex gap-2">
      <InputText 
        v-model="settingsStore.settings.customGoBinaryPath" 
        class="flex-1"
        placeholder="/usr/local/go/bin/go"
      />
      <Button 
        icon="pi pi-folder-open" 
        severity="secondary"
        @click="selectGoBinary"
      />
    </div>
  </div>
  
  <!-- Go Version Detection -->
  <div class="bg-gray-50 dark:bg-gray-900 p-4 rounded">
    <div class="flex items-center justify-between mb-2">
      <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
        Go Installation Status
      </span>
      <Button 
        label="Detect" 
        size="small" 
        @click="detectGo"
        :loading="detectingGo"
      />
    </div>
    
    <div v-if="goDetectionResult">
      <div v-if="goDetectionResult.found" class="flex items-start gap-2">
        <i class="pi pi-check-circle text-green-500 mt-1"></i>
        <div>
          <p class="text-sm text-gray-700 dark:text-gray-300">
            {{ goDetectionResult.version }}
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-400">
            {{ goDetectionResult.path }}
          </p>
        </div>
      </div>
      <div v-else class="flex items-start gap-2">
        <i class="pi pi-times-circle text-red-500 mt-1"></i>
        <div>
          <p class="text-sm text-red-600 dark:text-red-400">
            Go not found or invalid
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-400">
            {{ goDetectionResult.error }}
          </p>
        </div>
      </div>
    </div>
  </div>
</div>
```

**Settings Store Changes (settings.ts):**

```typescript
export interface AppSettings {
  // ... existing settings ...
  
  // Go Configuration
  goBinarySource: 'bundled' | 'system' | 'custom'
  customGoBinaryPath: string
}

const DEFAULT_SETTINGS: AppSettings = {
  // ... existing defaults ...
  
  goBinarySource: 'bundled',  // Default to bundled Go
  customGoBinaryPath: ''
}
```

**SettingsView Script:**

```typescript
const goBinarySourceOptions = [
  { label: 'Bundled Go (Recommended)', value: 'bundled' },
  { label: 'System Go (from PATH)', value: 'system' },
  { label: 'Custom Go Path', value: 'custom' }
]

const detectingGo = ref(false)
const goDetectionResult = ref<any>(null)

async function detectGo() {
  detectingGo.value = true
  try {
    const customPath = settingsStore.settings.goBinarySource === 'custom' 
      ? settingsStore.settings.customGoBinaryPath 
      : undefined
    
    goDetectionResult.value = await invoke('detect_go_installation', { 
      customPath 
    })
  } catch (error) {
    goDetectionResult.value = {
      found: false,
      error: String(error)
    }
  } finally {
    detectingGo.value = false
  }
}

async function selectGoBinary() {
  // Use Tauri dialog API to select Go binary
  const { open } = await import('@tauri-apps/plugin-dialog')
  const selected = await open({
    title: 'Select Go Binary',
    filters: [{
      name: 'Go Binary',
      extensions: process.platform === 'win32' ? ['exe'] : ['*']
    }]
  })
  
  if (selected) {
    settingsStore.updateSettings({ 
      customGoBinaryPath: selected 
    })
    await detectGo()  // Verify selection
  }
}

// Auto-detect on mount
onMounted(() => {
  detectGo()
})
```

### Solution 3: Persist Transpilation State

**Transpile Store Enhancement (transpile.ts):**

```typescript
export interface TranspilationState {
  projectPath: string
  outputDir: string
  filesTranspiled: number
  timestamp: Date
  success: boolean
}

export const useTranspileStore = defineStore('transpile', () => {
  // ... existing state ...
  
  const transpilationStates = ref<Map<string, TranspilationState>>(new Map())
  
  // Load states from localStorage
  function loadTranspilationStates() {
    try {
      const stored = localStorage.getItem('ts2go-transpilation-states')
      if (stored) {
        const parsed = JSON.parse(stored)
        transpilationStates.value = new Map(
          parsed.map((state: any) => [
            state.projectPath,
            {
              ...state,
              timestamp: new Date(state.timestamp)
            }
          ])
        )
      }
    } catch (e) {
      console.error('Failed to load transpilation states:', e)
    }
  }
  
  // Save states to localStorage
  function saveTranspilationStates() {
    try {
      const states = Array.from(transpilationStates.value.entries()).map(
        ([path, state]) => ({ ...state, timestamp: state.timestamp.toISOString() })
      )
      localStorage.setItem('ts2go-transpilation-states', JSON.stringify(states))
    } catch (e) {
      console.error('Failed to save transpilation states:', e)
    }
  }
  
  // Save state after successful transpilation
  function saveTranspilationState(projectPath: string, result: TranspileResult) {
    if (result.success && result.output_dir) {
      transpilationStates.value.set(projectPath, {
        projectPath,
        outputDir: result.output_dir,
        filesTranspiled: result.files_transpiled || 0,
        timestamp: new Date(),
        success: true
      })
      saveTranspilationStates()
    }
  }
  
  // Get state for specific project
  function getTranspilationState(projectPath: string): TranspilationState | null {
    return transpilationStates.value.get(projectPath) || null
  }
  
  // Clear state for specific project
  function clearTranspilationState(projectPath: string) {
    transpilationStates.value.delete(projectPath)
    saveTranspilationStates()
  }
  
  // Clear all states
  function clearAllTranspilationStates() {
    transpilationStates.value.clear()
    saveTranspilationStates()
  }
  
  // Verify if output directory still exists
  async function verifyTranspilationState(projectPath: string): Promise<boolean> {
    const state = getTranspilationState(projectPath)
    if (!state) return false
    
    try {
      const { exists } = await invoke('check_directory_exists', {
        path: state.outputDir
      })
      
      if (!exists) {
        clearTranspilationState(projectPath)
        return false
      }
      
      return true
    } catch {
      return false
    }
  }
  
  // Update transpileProject to save state
  async function transpileProject(projectPath: string) {
    // ... existing transpilation logic ...
    
    if (currentResult.value && currentResult.value.success) {
      saveTranspilationState(projectPath, currentResult.value)
      transpileHistory.value.push(currentResult.value)
    }
    
    // ... rest of function ...
  }
  
  // Initialize
  loadTranspilationStates()
  
  return {
    // ... existing exports ...
    transpilationStates,
    getTranspilationState,
    clearTranspilationState,
    clearAllTranspilationStates,
    verifyTranspilationState
  }
})
```

**ProjectView Enhancement:**

```vue
<script setup lang="ts">
// ... existing imports ...

const hasValidTranspilationState = ref(false)
const lastTranspilationTime = ref<string>('')

onMounted(async () => {
  // Check if current project has transpilation state
  if (projectPath.value) {
    const state = transpileStore.getTranspilationState(projectPath.value)
    if (state) {
      // Verify output directory still exists
      const isValid = await transpileStore.verifyTranspilationState(projectPath.value)
      
      if (isValid) {
        hasValidTranspilationState.value = true
        
        // Restore currentResult so buttons appear
        transpileStore.currentResult = {
          success: true,
          output_dir: state.outputDir,
          files_transpiled: state.filesTranspiled
        }
        
        // Calculate time ago
        const minutesAgo = Math.floor(
          (Date.now() - state.timestamp.getTime()) / 60000
        )
        lastTranspilationTime.value = minutesAgo < 1 
          ? 'Just now' 
          : `${minutesAgo} minute${minutesAgo > 1 ? 's' : ''} ago`
      }
    }
  }
})

async function verifyAndRestoreState() {
  if (!projectPath.value) return
  
  const isValid = await transpileStore.verifyTranspilationState(projectPath.value)
  
  if (isValid) {
    hasValidTranspilationState.value = true
  } else {
    hasValidTranspilationState.value = false
    toast.add({
      severity: 'warn',
      summary: 'Output Missing',
      detail: 'Previously transpiled output directory no longer exists.',
      life: 5000
    })
  }
}
</script>

<template>
  <!-- Add this section above action buttons -->
  <div v-if="hasValidTranspilationState && !transpileStore.currentResult" 
       class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 mb-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <i class="pi pi-check-circle text-blue-600 dark:text-blue-400 text-xl"></i>
        <div>
          <p class="text-sm font-medium text-blue-900 dark:text-blue-100">
            Previously Transpiled
          </p>
          <p class="text-xs text-blue-700 dark:text-blue-300">
            {{ lastTranspilationTime }} • {{ transpileStore.getTranspilationState(projectPath)?.filesTranspiled }} files
          </p>
        </div>
      </div>
      <div class="flex gap-2">
        <Button 
          label="Verify Output" 
          icon="pi pi-refresh" 
          size="small" 
          severity="secondary"
          @click="verifyAndRestoreState"
        />
        <Button 
          label="Clear State" 
          icon="pi pi-times" 
          size="small" 
          severity="danger"
          text
          @click="transpileStore.clearTranspilationState(projectPath)"
        />
      </div>
    </div>
  </div>
  
  <!-- Existing buttons - now also check hasValidTranspilationState -->
  <button 
    v-if="hasValidTranspilationState || transpileStore.currentResult?.output_dir"
    @click="handleBuildProject" 
    class="btn-action btn-build"
  >
    🔨 Build Project
  </button>
  <!-- ... other buttons ... -->
</template>
```

**New Tauri Command:**

```rust
#[tauri::command]
async fn check_directory_exists(path: String) -> Result<serde_json::Value, String> {
    let dir_path = Path::new(&path);
    Ok(serde_json::json!({
        "exists": dir_path.exists() && dir_path.is_dir()
    }))
}
```

## Implementation Priority

### Phase 0: Go Formatting After Transpilation (CRITICAL - Quick Win)
**Effort:** ~2-4 hours  
**Impact:** Essential code quality improvement  
**Risk:** Low  
**Dependency:** Can use system Go initially, bundled Go later

**Why First:**
- Simplest of all critical features
- Required regardless of other features
- Demonstrates proper Go toolchain integration
- Sets pattern for other Go commands (build, test, run)
- Can be implemented immediately with system Go

**Implementation Steps:**

1. **Add format_go_files() helper in main.rs** (1 hour)
   ```rust
   fn format_go_files(output_dir: &str) -> Result<FormatResult, String> {
       let go_path = get_go_binary_path(None)?;  // Use bundled or system
       
       let output = Command::new(&go_path)
           .arg("fmt")
           .arg("./...")
           .current_dir(output_dir)
           .output()
           .map_err(|e| format!("Failed to format Go code: {}", e))?;
       
       let stdout = String::from_utf8_lossy(&output.stdout);
       let stderr = String::from_utf8_lossy(&output.stderr);
       
       // go fmt prints formatted file paths to stdout
       let formatted_files: Vec<&str> = stdout.lines()
           .filter(|l| !l.is_empty())
           .collect();
       
       Ok(FormatResult {
           success: output.status.success(),
           files_formatted: formatted_files.len(),
           warnings: if !stderr.is_empty() { 
               vec![stderr.to_string()] 
           } else { 
               vec![] 
           }
       })
   }
   ```

2. **Update auto_transpile_project()** (30 min)
   ```rust
   if result.status.success() {
       let mut go_files = Vec::new();
       find_go_files(Path::new(&output_dir), &mut go_files)?;
       
       // NEW: Format all Go files
       let format_result = match format_go_files(&output_dir) {
           Ok(res) => res,
           Err(e) => {
               // Don't fail transpilation on format errors
               eprintln!("Warning: Failed to format Go code: {}", e);
               FormatResult {
                   success: false,
                   files_formatted: 0,
                   warnings: vec![e]
               }
           }
       };
       
       Ok(serde_json::json!({
           "success": true,
           "output_dir": output_dir,
           "files_transpiled": go_files.len(),
           "files_formatted": format_result.files_formatted,
           "format_warnings": format_result.warnings
       }))
   }
   ```

3. **Add FormatResult struct** (15 min)
   ```rust
   struct FormatResult {
       success: bool,
       files_formatted: usize,
       warnings: Vec<String>
   }
   ```

4. **Update transpileStore types** (15 min)
   ```typescript
   export interface TranspileResult {
     success: boolean
     output_dir?: string
     files_transpiled?: number
     files_formatted?: number        // NEW
     format_warnings?: string[]      // NEW
     goCode?: string
     error?: string
     message?: string
   }
   ```

5. **Update ProjectView to show format status** (30 min)
   ```vue
   <div v-if="transpileStore.currentResult?.files_formatted" 
        class="text-sm text-gray-600 dark:text-gray-400">
     ✓ {{ transpileStore.currentResult.files_formatted }} files formatted
   </div>
   
   <div v-if="transpileStore.currentResult?.format_warnings?.length" 
        class="text-sm text-yellow-600 dark:text-yellow-400">
     ⚠️ Formatting warnings: {{ transpileStore.currentResult.format_warnings.join(', ') }}
   </div>
   ```

6. **Test with valid and invalid Go syntax** (1 hour)
   - Transpile valid TypeScript → verify Go code formatted
   - Manually break Go syntax → verify warning (not failure)
   - Test with system Go
   - Test after Phase 3 with bundled Go

**Testing:**
- [ ] Valid TypeScript transpiles and formats successfully
- [ ] Format errors produce warnings, not failures
- [ ] UI shows formatting status
- [ ] Works with system Go (Phase 0a)
- [ ] Works with bundled Go after Phase 3 (Phase 0b)

### Phase 1: Transpilation State Persistence (High Priority - Quick Win)
**Effort:** ~4-6 hours  
**Impact:** Immediate UX improvement  
**Risk:** Low  

1. Update transpile.ts with state persistence (2 hours)
2. Add new Tauri command for directory check (30 min)
3. Update ProjectView.vue with state restoration (2 hours)
4. Test across views and app restarts (1 hour)

### Phase 2: Go Configuration UI (Medium Priority)
**Effort:** ~4-6 hours  
**Impact:** Enables custom Go paths  
**Risk:** Low  

1. Add settings interface and UI (2 hours)
2. Implement detect_go_installation command (1 hour)
3. Add file picker integration (1 hour)
4. Update all Command::new("go") calls (1 hour)
5. Test with system Go and custom paths (1 hour)

### Phase 3: Bundle Go Compiler (High Priority - Complex)
**Effort:** ~16-24 hours  
**Impact:** Makes app fully self-contained  
**Risk:** Medium (app size, platform testing)  

1. Research and download Go binaries (2 hours)
2. Update Makefile for Go bundling (3 hours)
3. Update Tauri config and resources (1 hour)
4. Implement get_go_binary_path() helper (2 hours)
5. Update all Command calls to use helper (2 hours)
6. Test bundled Go on all platforms (4 hours)
7. Optimize size if needed (4 hours)
8. Update documentation (2 hours)

## Testing Checklist

**Phase 0: Go Formatting**
- [ ] Transpiled Go code is automatically formatted
- [ ] Format errors produce warnings (not failures)
- [ ] UI shows "N files formatted" status
- [ ] Format warnings displayed in UI
- [ ] Works with system Go
- [ ] Works with bundled Go (after Phase 3)

**Phase 1: Transpilation State**
- [ ] Transpilation state persists across view changes
- [ ] Transpilation state persists across app restarts
- [ ] Transpilation state cleared when output deleted
- [ ] Build/Test/Run buttons appear with valid state

**Phase 2: Go Configuration**
- [ ] Go detection shows correct version and path
- [ ] Settings persist Go configuration
- [ ] Build/Test/Run work with system Go
- [ ] Build/Test/Run work with custom Go path

**Phase 3: Bundled Go**
- [ ] Build/Test/Run work with bundled Go
- [ ] App size acceptable with bundled Go
- [ ] Works on macOS (ARM64 and Intel)
- [ ] Works on Windows (x64)
- [ ] Works on Linux (x64)
- [ ] All formatting operations use bundled Go

## Documentation Updates

- [ ] **Phase 0:** Document automatic Go formatting in USER_GUIDE.md
- [ ] **Phase 0:** Add troubleshooting for format errors
- [ ] **Phase 1:** Document transpilation state persistence
- [ ] **Phase 2:** Add Go Configuration section to USER_GUIDE.md
- [ ] **Phase 2:** Add troubleshooting for Go-related errors
- [ ] **Phase 3:** Update build instructions for Go bundling
- [ ] **Phase 3:** Document app size with/without bundled Go
- [ ] **Phase 2:** Add screenshots of new settings UI

## Alternative Approaches Considered

### 1. Download Go on Demand
**Pros:** Smaller initial app size  
**Cons:** Requires internet, complex implementation, platform-specific installers  
**Decision:** Rejected - too complex

### 2. Require User to Install Go
**Pros:** No app size increase  
**Cons:** Poor UX, defeats purpose of self-contained desktop app  
**Decision:** Rejected - unacceptable UX

### 3. Use Go WASM
**Pros:** Browser-compatible  
**Cons:** Performance issues, limited stdlib, experimental  
**Decision:** Rejected - not production-ready

## Success Criteria

1. ✅ **Phase 0:** Transpiled Go code is automatically formatted with `go fmt`
2. ✅ **Phase 0:** Format errors produce warnings but don't fail transpilation
3. ✅ **Phase 1:** App remembers transpilation state across sessions
4. ✅ **Phase 2:** Users can configure Go binary source in settings
5. ✅ **Phase 3:** Users can Build/Test/Run without installing Go
6. ✅ **Phase 2:** Clear error messages if Go not available
7. ✅ **Phase 3:** App size increase acceptable (<600MB total)
8. ✅ **All Phases:** Works reliably on all platforms

## Next Steps

1. **Immediate (Phase 0):** Implement automatic `go fmt` formatting (~2-4 hours)
   - Simple, high-impact, prerequisite for quality
   - Can use system Go initially
   - Will automatically use bundled Go once Phase 3 complete

2. **Short-term (Phase 1):** Implement transpilation state persistence (~4-6 hours)
   - Critical UX improvement
   - Low risk, immediate value

3. **Short-term (Phase 2):** Implement Go configuration UI (~4-6 hours)
   - Enables flexibility
   - Prepares for Phase 3

4. **Medium-term (Phase 3):** Bundle Go compiler (~16-24 hours)
   - Makes app fully self-contained
   - Largest effort but highest long-term value

5. **Validate:** Test with real users, gather feedback on app size vs. convenience trade-off

**Total Estimated Effort:** ~26-40 hours (was 24-36 before Phase 0)

**Priority Rationale:**
- Phase 0 is prerequisite for code quality - do first
- Phase 1 provides immediate UX wins
- Phase 2 enables user flexibility
- Phase 3 achieves self-contained desktop app goal
