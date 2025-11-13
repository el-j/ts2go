# Roadmap to Alpha 2.0.1

**Target Release:** Alpha 2.0.1  
**Goal:** Add Go build & test integration to GUI and CLI  
**Timeline:** 2-3 weeks  
**Last Updated:** November 13, 2025

---

## Release Goal

Add **Go compilation** and **test execution** capabilities to ts2go, making it a complete TypeScript→Go development tool with:
- Build transpiled code into binaries
- Run tests on transpiled code
- View build/test results in GUI
- Manage build artifacts

**Current State:** Can transpile and **run** Go code ✅  
**Target State:** Can transpile, run, **build**, and **test** Go code ✅

---

## Phase 1: Go Build Integration (Week 1)

### 1.1 Backend Build Command (2-3 days)

**Add to Tauri Backend** (`desktop-ui/src-tauri/src/main.rs`):

```rust
#[tauri::command]
async fn build_go_file(code: String, output_path: String) -> Result<serde_json::Value, String> {
    // Create temp Go file
    // Execute: go build -o <output_path> <temp_file>
    // Return: success, errors, warnings, binary_path, size, duration_ms
}

#[tauri::command]
async fn build_go_project(
    source_dir: String,
    output_path: String,
    build_flags: Option<Vec<String>>
) -> Result<serde_json::Value, String> {
    // Execute: go build -o <output_path> [flags] .
    // Handle: main package validation, go.mod existence
    // Return: success, errors, binary_path, size, dependencies, duration_ms
}
```

**Add to CLI** (`pkg/cli/build.go` - NEW):

```go
func BuildCommand(args []string) error {
    // Parse flags: --source, --output, --flags
    // Validate transpiled Go code exists
    // Execute go build
    // Report progress and results
}
```

**Tests:**
- Unit tests for build functions
- Integration test: transpile + build simple project
- Verify binary is executable

### 1.2 UI Build Integration (2 days)

**Add Build Store** (`desktop-ui/src/stores/build.ts` - NEW):

```typescript
export const useBuildStore = defineStore('build', () => {
  const isBuilding = ref(false)
  const buildResult = ref<BuildResult | null>(null)
  const buildHistory = ref<BuildResult[]>([])
  
  interface BuildResult {
    success: boolean
    binary_path: string
    binary_size: number
    errors: string[]
    warnings: string[]
    duration_ms: number
    timestamp: number
  }
  
  async function buildFile(code: string, outputPath: string) {
    // Call build_go_file Tauri command
  }
  
  async function buildProject(sourceDir: string, outputPath: string) {
    // Call build_go_project Tauri command
  }
})
```

**Update ProjectView** (`desktop-ui/src/views/ProjectView.vue`):
- Add "🔨 Build" button next to "▶️ Run" button
- Add build output panel tab
- Show build errors/warnings
- Display binary location and size
- Add "Run Built Binary" button

**Update EditorView** (`desktop-ui/src/views/EditorView.vue`):
- Add build button for single files
- Show build results

**Success Criteria:**
- Can build single Go file from UI
- Can build full project from UI
- Build results display in output panel
- Binary path is clickable/copyable
- Build errors show clearly

---

## Phase 2: Go Test Integration (Week 2)

### 2.1 Backend Test Command (2-3 days)

**Add to Tauri Backend** (`desktop-ui/src-tauri/src/main.rs`):

```rust
#[tauri::command]
async fn test_go_file(file_path: String) -> Result<serde_json::Value, String> {
    // Execute: go test <file_path>
    // Parse: test results, pass/fail, timing
    // Return: success, passed, failed, skipped, duration_ms, test_results[]
}

#[tauri::command]
async fn test_go_project(
    source_dir: String,
    test_flags: Option<Vec<String>>
) -> Result<serde_json::Value, String> {
    // Execute: go test [flags] ./...
    // Parse: all test results
    // Return: summary + detailed results
}
```

**Add to CLI** (`pkg/cli/test.go` - NEW):

```go
func TestCommand(args []string) error {
    // Parse flags: --source, --verbose, --coverage
    // Execute go test
    // Parse and display results
}
```

**Test Result Parsing:**
- Parse `go test` JSON output (`-json` flag)
- Extract: test names, pass/fail, duration, output
- Handle: table-driven tests, subtests, benchmarks

**Tests:**
- Unit tests for test execution
- Integration test: transpile + test project with tests
- Verify test results are parsed correctly

### 2.2 UI Test Integration (2 days)

**Add Test Store** (`desktop-ui/src/stores/test.ts` - NEW):

```typescript
export const useTestStore = defineStore('test', () => {
  const isTesting = ref(false)
  const testResult = ref<TestResult | null>(null)
  const testHistory = ref<TestResult[]>([])
  
  interface TestResult {
    success: boolean
    passed: number
    failed: number
    skipped: number
    tests: TestCase[]
    coverage?: number
    duration_ms: number
    timestamp: number
  }
  
  interface TestCase {
    name: string
    passed: boolean
    duration_ms: number
    output: string
  }
  
  async function testFile(filePath: string) {
    // Call test_go_file Tauri command
  }
  
  async function testProject(sourceDir: string) {
    // Call test_go_project Tauri command
  }
})
```

**Create TestResultsPanel** (`desktop-ui/src/components/TestResultsPanel.vue` - NEW):
- Tree view of test results
- Pass/fail indicators (✅/❌)
- Test timing
- Expandable test output
- Summary stats (X passed, Y failed, Z skipped)
- Coverage percentage (if available)

**Update ProjectView:**
- Add "🧪 Test" button
- Add test results panel tab
- Show test progress
- Display test summary

**Success Criteria:**
- Can run tests from UI
- Test results display in tree format
- Pass/fail is visually clear
- Can view individual test output
- Test timing is shown

---

## Phase 3: Artifact Management (Week 3)

### 3.1 Build Artifacts Storage (1-2 days)

**Add Artifacts Store** (`desktop-ui/src/stores/artifacts.ts` - NEW):

```typescript
export const useArtifactsStore = defineStore('artifacts', () => {
  const artifacts = ref<Artifact[]>([])
  
  interface Artifact {
    id: string
    type: 'binary' | 'test_results'
    path: string
    size: number
    created_at: number
    project_path: string
    metadata: {
      go_version: string
      platform: string
      architecture: string
    }
  }
  
  function addArtifact(artifact: Artifact) {
    // Save to localStorage
  }
  
  function getProjectArtifacts(projectPath: string) {
    // Filter artifacts by project
  }
  
  async function openArtifact(artifactPath: string) {
    // Open file explorer / execute binary
  }
})
```

**Add Artifacts View** (`desktop-ui/src/views/ArtifactsView.vue` - NEW):
- List all build artifacts
- Filter by project
- Show size, date, platform
- Actions: open folder, run binary, delete
- Clean up old artifacts

### 3.2 Build History & Management (1-2 days)

**Enhance Build Store:**
- Persist build history to localStorage
- Track: builds over time, success rate, duration trends
- Associate builds with specific file versions (git hash if available)

**Create BuildHistoryPanel** (`desktop-ui/src/components/BuildHistoryPanel.vue` - NEW):
- Timeline of builds
- Success/failure indicators
- Build duration graph
- Click to view build details
- Rebuild from history

**Update HistoryView:**
- Add builds tab alongside transpilation history
- Show combined history: transpile → build → test

### 3.3 CLI Enhancements (1 day)

**Update CLI Commands:**

```bash
# Build command
ts2go build <source-dir> --output <binary-path> [--flags "..."]

# Test command
ts2go test <source-dir> [--verbose] [--coverage]

# Package command (future)
ts2go package <source-dir> --platform <target> --output <package-path>
```

**Add:** (`cmd/ts2go/main.go`)
```go
case "build", "b":
    if err := cli.BuildCommand(args); err != nil {
        // ...
    }
case "test", "t":
    if err := cli.TestCommand(args); err != nil {
        // ...
    }
```

---

## Phase 4: Polish & Testing (Ongoing)

### 4.1 Error Handling
- Graceful handling when Go is not installed
- Clear error messages for build failures
- Validation of Go environment

### 4.2 UI/UX Improvements
- Loading states for long builds
- Cancel build/test operations
- Keyboard shortcuts (Ctrl+B for build, Ctrl+T for test)
- Status bar indicators

### 4.3 Documentation
- Update user guide with build/test features
- Add examples of building projects
- Document binary management

### 4.4 Testing
- E2E tests for build workflow
- E2E tests for test workflow
- Integration tests for CLI commands

---

## Success Criteria for Alpha 2.0.1

### Must Have
- ✅ Build single Go files from GUI
- ✅ Build Go projects from GUI
- ✅ Run Go tests from GUI
- ✅ Display build results with errors/warnings
- ✅ Display test results with pass/fail
- ✅ CLI build command
- ✅ CLI test command
- ✅ Build artifacts stored and managed
- ✅ Build/test history tracking

### Nice to Have
- Coverage reports display
- Benchmark result visualization
- Cross-compilation support (build for different platforms)
- Binary size optimization suggestions
- Incremental builds

### Quality Gates
- All existing tests pass
- New features have unit tests
- Integration tests for build/test workflows
- No regressions in existing transpilation/run features
- Documentation updated

---

## Implementation Plan Summary

**Week 1: Build**
- Days 1-3: Backend build commands (Rust + Go CLI)
- Days 4-5: UI build integration

**Week 2: Test**
- Days 1-3: Backend test commands + parsing
- Days 4-5: UI test integration

**Week 3: Polish**
- Days 1-2: Artifact management
- Days 3-4: Build history & UI polish
- Day 5: CLI enhancements & testing

---

## Dependencies

### Required
- Go compiler installed on user machine
- Transpiled Go code is valid and compilable
- go.mod exists for projects (can be auto-generated)

### Optional
- Git (for version tracking in build history)
- Cross-compilation tools (for multi-platform builds)

---

## Risk Management

**Risk:** Build failures due to transpilation issues  
**Mitigation:** Validate transpiled code before building, show clear error messages

**Risk:** Go not installed on user machine  
**Mitigation:** Detect Go installation, show helpful error with install link

**Risk:** Test parsing fails on complex test output  
**Mitigation:** Use `go test -json` for structured output, graceful fallback

**Risk:** Large binaries in artifact storage  
**Mitigation:** Configurable artifact retention, cleanup utilities

---

## Post-Alpha 2.0.1 (Future)

- **Cross-compilation:** Build for Windows, macOS, Linux from any platform
- **Package command:** Create distributable packages (.deb, .rpm, .msi, .dmg)
- **CI/CD integration:** Export build/test configs for GitHub Actions, etc.
- **Plugin system:** Allow custom build steps, test frameworks
- **Profiling:** Integrate Go profiling tools (pprof)
- **Benchmarking:** Run and visualize benchmarks

---

*This roadmap focuses on completing the transpile → build → test → run cycle for alpha 2.0.1.*
