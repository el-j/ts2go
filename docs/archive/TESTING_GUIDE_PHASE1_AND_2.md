# Testing Guide - Phase 1 & 2 Implementation

**Purpose:** Comprehensive testing checklist for Go Configuration and State Persistence features  
**Status:** Tasks 1-8 Complete, Tasks 9-10 Pending  
**Last Updated:** November 2025

---

## Test Environment Setup

### Prerequisites
- ✅ Rust build passing: `cd desktop-ui/src-tauri && cargo check`
- ✅ TypeScript compilation clean
- ✅ Tauri app can be launched: `cd desktop-ui && npm run tauri dev`

---

## Test Suite 1: Go Configuration UI

### Test 1.1: Settings Page Go Configuration
**Location:** Settings → Project Tab → Go Configuration Section

**Steps:**
1. Launch app
2. Navigate to Settings (gear icon)
3. Click "Project" tab
4. Locate "Go Compiler Configuration" section

**Expected Results:**
- ✅ Dropdown shows "Go Binary Source" with options:
  - "System Go (from PATH)"
  - "Custom Path"
- ✅ "Detect Go Installation" button visible
- ✅ Auto-detection runs on mount
- ✅ Result banner shows Go version and path (if Go installed)

**Pass/Fail:** ___________

---

### Test 1.2: System Go Detection
**Precondition:** Go installed and in system PATH

**Steps:**
1. Open Settings → Project tab
2. Select "System Go (from PATH)" in dropdown
3. Click "Detect Go Installation"

**Expected Results:**
- ✅ Green banner appears with checkmark icon
- ✅ Shows "Go Compiler Detected"
- ✅ Displays Go version (e.g., "go version go1.21.0 darwin/arm64")
- ✅ Displays Go path (e.g., "/usr/local/go/bin/go")
- ✅ Detection completes in < 2 seconds

**Pass/Fail:** ___________

---

### Test 1.3: Go Not Found Error
**Precondition:** Temporarily rename Go binary or use custom path to non-existent file

**Steps:**
1. Open Settings → Project tab
2. Select "Custom Path" in dropdown
3. Enter invalid path: `/fake/path/to/go`
4. Click "Detect Go Installation"

**Expected Results:**
- ✅ Red banner appears with X icon
- ✅ Shows "Go Compiler Not Found"
- ✅ Error message: "Custom Go binary not found at: /fake/path/to/go"
- ✅ No version or path displayed

**Pass/Fail:** ___________

---

### Test 1.4: Custom Go Path Selection
**Precondition:** Know location of Go binary (run `which go`)

**Steps:**
1. Open Settings → Project tab
2. Select "Custom Path" in dropdown
3. Verify custom path input appears
4. Click folder icon button
5. Browse to Go binary location
6. Select Go binary
7. Auto-detection should run

**Expected Results:**
- ✅ File picker opens
- ✅ Custom path input shows selected path
- ✅ Auto-detection runs immediately
- ✅ Green banner appears if valid Go binary

**Pass/Fail:** ___________

---

## Test Suite 2: Go Command Smart Detection

### Test 2.1: Build with System Go
**Precondition:** Go in PATH, project transpiled

**Steps:**
1. Load TypeScript project
2. Transpile project
3. Click "🔨 Build Project"
4. Check output logs

**Expected Results:**
- ✅ Build succeeds
- ✅ Logs show: "Building Go project..."
- ✅ Binary created in output_dir
- ✅ No errors about Go not found

**Pass/Fail:** ___________

---

### Test 2.2: Test with System Go
**Precondition:** Go in PATH, project transpiled

**Steps:**
1. Transpile project with test files
2. Click "🧪 Test Project"
3. Check output logs

**Expected Results:**
- ✅ Tests run successfully
- ✅ Logs show test results
- ✅ Test summary displayed (passed/failed/skipped)
- ✅ No errors about Go not found

**Pass/Fail:** ___________

---

### Test 2.3: Run with System Go
**Precondition:** Go in PATH, project transpiled

**Steps:**
1. Transpile project
2. Click "▶️ Run Project"
3. Check output logs

**Expected Results:**
- ✅ Project runs successfully
- ✅ Logs show execution output
- ✅ Exit code displayed
- ✅ No errors about Go not found

**Pass/Fail:** ___________

---

### Test 2.4: Build without Go
**Precondition:** Temporarily remove Go from PATH or configure invalid custom path

**Steps:**
1. Settings → Custom Path → `/invalid/path`
2. Transpile project
3. Click "🔨 Build Project"
4. Check error message

**Expected Results:**
- ✅ Error toast appears
- ✅ Message: "Go not found. Please install Go or configure in Settings"
- ✅ Build fails gracefully
- ✅ Logs show clear error

**Pass/Fail:** ___________

---

## Test Suite 3: Transpilation State Persistence

### Test 3.1: State Saved After Transpilation
**Steps:**
1. Load TypeScript project
2. Transpile project successfully
3. Open browser DevTools → Application → Local Storage
4. Check `ts2go_transpilation_states` key

**Expected Results:**
- ✅ Key exists in localStorage
- ✅ Value is valid JSON
- ✅ Contains project path, output dir, files count, timestamp
- ✅ `success: true`

**Console Command:**
```javascript
JSON.parse(localStorage.getItem('ts2go_transpilation_states'))
```

**Pass/Fail:** ___________

---

### Test 3.2: State Restored on View Switch
**Steps:**
1. Load and transpile project
2. Navigate away from ProjectView (e.g., to Settings)
3. Return to ProjectView (click Home/Project icon)

**Expected Results:**
- ✅ "Previous Transpilation Restored" banner appears
- ✅ Shows file count and timestamp
- ✅ Build/Test/Run buttons are enabled
- ✅ No need to re-transpile

**Pass/Fail:** ___________

---

### Test 3.3: State Restored After App Restart
**Steps:**
1. Load and transpile project
2. Close Tauri app completely
3. Relaunch app
4. Navigate to ProjectView
5. Load the same project

**Expected Results:**
- ✅ State restoration banner appears automatically
- ✅ Correct file count and timestamp displayed
- ✅ Build/Test/Run buttons enabled immediately
- ✅ Output directory verified

**Pass/Fail:** ___________

---

### Test 3.4: State Cleared When Output Dir Missing
**Steps:**
1. Load and transpile project
2. Manually delete output directory (ts-go-project_output)
3. Close and reopen app
4. Load same project

**Expected Results:**
- ✅ No restoration banner appears
- ✅ Logs show: "Previous transpilation output not found, state cleared"
- ✅ Build/Test/Run buttons disabled
- ✅ State removed from localStorage

**Pass/Fail:** ___________

---

### Test 3.5: Clear State Button
**Steps:**
1. Load and transpile project
2. Verify restoration banner appears
3. Click "Clear State" button on banner

**Expected Results:**
- ✅ Banner disappears
- ✅ Toast notification: "State Cleared"
- ✅ Build/Test/Run buttons disabled
- ✅ Logs show: "Transpilation state cleared"
- ✅ State removed from localStorage

**Pass/Fail:** ___________

---

## Test Suite 4: Multi-Project State Management

### Test 4.1: Multiple Project States
**Steps:**
1. Load Project A, transpile, verify state saved
2. Load Project B, transpile, verify state saved
3. Load Project C, transpile, verify state saved
4. Check localStorage

**Expected Results:**
- ✅ All 3 projects in localStorage map
- ✅ Each has unique projectPath key
- ✅ Each has separate state data

**Console Command:**
```javascript
const states = JSON.parse(localStorage.getItem('ts2go_transpilation_states'))
Object.keys(states).forEach(key => console.log(key, states[key]))
```

**Pass/Fail:** ___________

---

### Test 4.2: Switch Between Projects with State
**Steps:**
1. Transpile Project A
2. Switch to Project B, transpile
3. Switch back to Project A

**Expected Results:**
- ✅ Project A state restored correctly
- ✅ Shows Project A's file count and timestamp
- ✅ Build buttons work for Project A

**Pass/Fail:** ___________

---

## Test Suite 5: Format Integration

### Test 5.1: Go Formatting After Transpile
**Steps:**
1. Transpile project
2. Check toast notifications
3. Check output logs

**Expected Results:**
- ✅ Toast shows: "Transpiled X files successfully (Y formatted)"
- ✅ Logs show format results
- ✅ Go files are properly formatted (check manually)

**Pass/Fail:** ___________

---

### Test 5.2: Format with Custom Go
**Steps:**
1. Configure custom Go path in Settings
2. Transpile project
3. Verify formatting works

**Expected Results:**
- ✅ Formatting uses custom Go binary
- ✅ No errors
- ✅ Files formatted correctly

**Pass/Fail:** ___________

---

## Test Suite 6: Error Handling

### Test 6.1: Graceful Go Not Found
**Steps:**
1. Remove Go from PATH
2. Try each command: Build, Test, Run

**Expected Results:**
- ✅ Each shows clear error
- ✅ No crashes or silent failures
- ✅ Error mentions "configure in Settings"

**Pass/Fail:** ___________

---

### Test 6.2: Invalid Custom Path Handling
**Steps:**
1. Set custom path to directory (not file)
2. Try to detect Go
3. Try to build project

**Expected Results:**
- ✅ Detection shows error
- ✅ Build shows error
- ✅ Clear message about invalid path

**Pass/Fail:** ___________

---

## Test Suite 7: UI/UX

### Test 7.1: Loading States
**Steps:**
1. Click "Detect Go Installation"
2. Observe button during detection

**Expected Results:**
- ✅ Button shows loading spinner
- ✅ Button disabled during detection
- ✅ Result appears after completion

**Pass/Fail:** ___________

---

### Test 7.2: Banner Responsiveness
**Steps:**
1. Restore state, view banner
2. Resize window to narrow width
3. Check banner layout

**Expected Results:**
- ✅ Banner remains readable
- ✅ No text overflow
- ✅ Clear State button accessible

**Pass/Fail:** ___________

---

## Test Suite 8: Edge Cases

### Test 8.1: Empty Output Directory
**Steps:**
1. Transpile project
2. Delete all files in output_dir (but keep directory)
3. Restart app, load project

**Expected Results:**
- ✅ State verification detects directory exists
- ✅ State restored (directory exists)
- ✅ Build may fail (no files), but state logic works

**Pass/Fail:** ___________

---

### Test 8.2: Corrupted localStorage
**Steps:**
1. Manually corrupt localStorage state:
   ```javascript
   localStorage.setItem('ts2go_transpilation_states', 'invalid json{{{')
   ```
2. Restart app
3. Load project

**Expected Results:**
- ✅ App doesn't crash
- ✅ Error logged to console
- ✅ Fresh state created
- ✅ User can transpile normally

**Pass/Fail:** ___________

---

## Test Suite 9: Integration Tests

### Test 9.1: Full Workflow - System Go
**Steps:**
1. Fresh app launch
2. Verify Settings shows detected Go
3. Load TS project
4. Transpile
5. Build
6. Test
7. Run
8. Close app
9. Reopen app
10. Load same project
11. Build without re-transpiling

**Expected Results:**
- ✅ All steps succeed
- ✅ State restored correctly
- ✅ Build works with restored state

**Pass/Fail:** ___________

---

### Test 9.2: Full Workflow - Custom Go
**Steps:**
1. Configure custom Go path
2. Load TS project
3. Transpile
4. Build
5. Verify binary uses custom Go

**Expected Results:**
- ✅ All commands use custom Go
- ✅ No fallback to system Go
- ✅ Success messages consistent

**Pass/Fail:** ___________

---

## Success Criteria

**Required to Pass:**
- [ ] All Test Suite 1 tests pass (Go Configuration UI)
- [ ] All Test Suite 2 tests pass (Smart Go Detection)
- [ ] All Test Suite 3 tests pass (State Persistence)
- [ ] At least 80% of Test Suite 4-9 pass (Multi-project, Format, Errors, UI, Edge Cases, Integration)
- [ ] No critical bugs or crashes
- [ ] No Rust compilation errors
- [ ] No TypeScript errors

**Optional (Nice to Have):**
- [ ] Performance tests (state load < 100ms)
- [ ] Accessibility tests (keyboard navigation)
- [ ] Cross-platform tests (macOS, Windows, Linux)

---

## Bug Reporting Template

```markdown
**Test ID:** [e.g., Test 3.2]
**Status:** FAIL
**Expected:** [What should happen]
**Actual:** [What actually happened]
**Steps to Reproduce:**
1. 
2. 
3. 

**Console Errors:**
```
[Paste console output]
```

**Screenshots:** [If applicable]
**Environment:** 
- OS: 
- Go Version: 
- App Version: 
```

---

## Testing Notes

### Common Issues to Watch For:
1. **localStorage quota exceeded** - Test with many projects
2. **Race conditions** - State save/load timing
3. **Path separators** - Windows vs Unix paths
4. **Permissions** - Go binary executable permissions
5. **Network delays** - File system operations

### Testing Tools:
- Browser DevTools → Application → Local Storage
- Console for localStorage inspection
- Network tab for Tauri command timing
- Terminal for Go command verification

---

**Last Updated:** November 2025  
**Version:** 1.0  
**Status:** Ready for Testing
