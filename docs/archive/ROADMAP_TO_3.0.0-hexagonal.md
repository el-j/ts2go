# TS2Go Hexagonal Architecture Refactoring Roadmap

**Status:** Phases 0-4 Complete ✅ (November 16, 2025)  
**Started:** November 16, 2025  
**Completion:** Phases 0-4 completed same day

This document outlines the comprehensive roadmap to refactor the `ts2go` application into a Hexagonal (Ports and Adapters) Architecture.

## 🎉 Progress Update - November 16, 2025

**Major Milestone Achieved:**

✅ **Phase 0: Foundation & Architecture Setup** - Complete  
✅ **Phase 1: Core Business Logic Migration** - Complete  
✅ **Phase 2: Infrastructure Adapters** - Complete  
✅ **Phase 3: CLI Refactoring** - Complete  
✅ **Phase 4: Tauri Backend Simplification** - Complete  
⏳ **Phase 5: Frontend State Migration** - Pending  
⏳ **Phase 6: Web API Proof of Concept** - Pending  
⏳ **Phase 7: Cleanup & Legacy Removal** - Pending  

**See `docs/HEXAGONAL_ARCHITECTURE_COMPLETE.md` for detailed completion report.**

---

## 🎯 Your Refactored Goal

The goal is to refactor the **`ts2go` Go application** (the core engine and CLI) into a Hexagonal Architecture. This will separate the core business logic (transpilation, analysis) from infrastructure concerns (file system, `go` compiler commands) and delivery mechanisms (the CLI, the Tauri UI).

This decoupling is essential for:

* **Team Scalability:** Allowing separate teams to work on the UI, the core engine, and future web services without conflict.
* **Maintainability:** Making the core logic pure, testable, and independent of any specific UI or database.
* **Future-Proofing:** Enabling you to add new "driving adapters" like a web service or a login system that reuse the same central core.

---

## 🏛️ The Hexagonal Architecture Plan

In this new structure:

* **The Hexagon (Core):** This is your pure, dependency-free business logic (e.g., `transpiler`, `analyzer`). It knows nothing about file systems, operating systems, or UIs.
* **Ports (Interfaces):** These are Go interfaces that define the contracts.
    * **Driving Ports (Inbound):** The API of your core logic (e.g., `TranspilationService`).
    * **Driven Ports (Outbound):** The dependencies the core logic needs (e.g., `FileSystem`, `GoCompiler`).
* **Adapters (Implementations):** These are the "pluggable" components from the outside world.
    * **Driving Adapters (Inbound):** They *call* the core logic. (e.g., your **CLI** and **Tauri UI**).
    * **Driven Adapters (Outbound):** They *implement* the core's dependencies. (e.g., a **DiskFileSystem** using the `os` package, an **OSGoCompiler** using `exec.Command`).



### High-Level Target Architecture



---

## 🚀 The Refactoring Roadmap

Here is a step-by-step plan to get this refactor done efficiently.

### Phase 0: Laying the Foundation (The Blueprint)

**Goal:** Create the new directory structure and define all the port interfaces. No logic is moved yet.

1.  **Create New Directory Structure:**
    Create a new top-level `pkg/` directory. Your existing `cmd/` and `internal/` can remain for now as you migrate.

    ```bash
    ts2go/
    ├── cmd/ts2go-cli/      # (Existing CLI)
    ├── desktop-ui/         # (Existing Tauri App)
    ├── internal/           # (Your old logic, will be moved from here)
    └── pkg/                # (The new architecture)
        ├── core/           # === THE HEXAGON ===
        │   ├── domain/     # Entities (e.g., Project, File, TranspilationResult)
        │   ├── ports/      # The Go interfaces (the most important part)
        │   └── services/   # Implementation of Driving Ports (the use cases)
        └── adapters/       # === THE OUTSIDE WORLD ===
            ├── driven/     # Implementations of Driven Ports
            │   ├── filesystem/ # Implements FileSystem port
            │   ├── gocompiler/ # Implements GoCompiler port
            │   └── persistence/ # Implements state/settings ports
            └── driving/    # (Code for driving adapters, e.g., web handlers)
                ├── cli/
                └── web/
    ```

2.  **Define Your Ports (Interfaces):**
    This is the most critical design step. In `pkg/core/ports/`, define the interfaces.

    > **File: `pkg/core/ports/ports.go`**
    ```go
    package ports
    
    // --- DRIVING / INBOUND PORTS (The Core API) ---
    
    // TranspilationService defines the core transpilation use case.
    type TranspilationService interface {
        TranspileProject(projectPath string) (*domain.TranspilationResult, error)
        AnalyzeProject(projectPath string) (*domain.AnalysisReport, error)
    }
    
    // GoRuntimeService defines the use cases for running Go code.
    type GoRuntimeService interface {
        DetectGoInstallation() (*domain.GoVersion, error)
        BuildProject(projectPath string) (*domain.BuildResult, error)
        RunProject(projectPath string) (*domain.RunResult, error)
        TestProject(projectPath string) (*domain.TestResult, error)
    }
    
    // StateService manages persistent state.
    type StateService interface {
        GetProjectState(projectPath string) (*domain.ProjectState, error)
        GetAllProjectStates() (map[string]*domain.ProjectState, error)
        GetSettings() (*domain.Settings, error)
        SaveSettings(settings *domain.Settings) error
    }
    
    // --- DRIVEN / OUTBOUND PORTS (The Dependencies) ---
    
    // FileSystem is what the core uses to access files.
    type FileSystem interface {
        ReadFile(path string) ([]byte, error)
        WriteFile(path string, data []byte) error
        ScanDirectory(path string) ([]domain.File, error)
        DirectoryExists(path string) (bool, error)
    }
    
    // GoCompiler is what the core uses to interact with the 'go' command.
    type GoCompiler interface {
        Detect() (*domain.GoVersion, error)
        Build(path string, outputPath string) (string, error)
        Run(path string) (string, error)
        Test(path string) (string, error)
        Format(path string) (string, error)
    }
    
    // StateRepository is what the core uses to persist state.
    type StateRepository interface {
        GetState(projectPath string) (*domain.ProjectState, error)
        SaveState(projectPath string, state *domain.ProjectState) error
        GetAllStates() (map[string]*domain.ProjectState, error)
    }
    
    // SettingsRepository is what the core uses to persist settings.
    type SettingsRepository interface {
        Get() (*domain.Settings, error)
        Save(settings *domain.Settings) error
    }
    ```

---

### Phase 1: Migrate the Core Domain Logic (The Hexagon)

**Goal:** Move your existing transpilation engine into `pkg/core/services/` and make it depend *only* on the new ports.

1.  **Copy, Paste, and Refactor:**
    * **Copy** your existing `internal/transpiler/`, `internal/analyzer/`, `internal/mapper/` logic into `pkg/core/services/`.
    * Create a struct for your service, e.g., `transpilation_service.go`.
    * This struct will hold its dependencies *as interfaces*.

    > **File: `pkg/core/services/transpilation_service.go`**
    ```go
    package services
    
    import "ts2go/pkg/core/ports"
    
    type transpilationService struct {
        fs       ports.FileSystem     // Dependency
        compiler ports.GoCompiler   // Dependency
        // ... other pure transpiler logic
    }
    
    func NewTranspilationService(fs ports.FileSystem, compiler ports.GoCompiler) ports.TranspilationService {
        return &transpilationService{
            fs:       fs,
            compiler: compiler,
        }
    }
    
    func (s *transpilationService) TranspileProject(projectPath string) (*domain.TranspilationResult, error) {
        // 1. Find all files using s.fs.ScanDirectory(projectPath)
        // 2. Read files using s.fs.ReadFile(...)
        // 3. Run analysis logic (your existing code)
        // 4. Run transpilation logic (your existing code)
        // 5. Write files using s.fs.WriteFile(...)
        // 6. Format files using s.compiler.Format(...)
        // ...
        return nil, nil //
    }
    
    // ... implement other TranspilationService methods
    ```

2.  **Replace All Infrastructure Calls:**
    * Go through your newly copied service code.
    * Find every `os.ReadFile`, `os.WriteFile`, `exec.Command("go", ...)`.
    * **Delete them.**
    * Replace them with calls to the port interfaces: `s.fs.ReadFile(...)`, `s.compiler.Build(...)`, etc.
    * Your core logic is now "pure." It has no direct link to the outside world.

---

### Phase 2: Create Infrastructure Adapters (The Plugs)

**Goal:** Implement the "driven" ports with real-world code. This is where you'll **copy** your *original* infrastructure logic.

1.  **Create FileSystem Adapter:**
    * Create `pkg/adapters/driven/filesystem/disk_fs.go`.
    * Implement the `ports.FileSystem` interface using the standard `os` package.
    * **This is easy:** You are just copying the `os` calls you *just deleted* from the core logic and putting them here.

2.  **Create GoCompiler Adapter:**
    * Create `pkg/adapters/driven/gocompiler/os_go.go`.
    * Implement the `ports.GoCompiler` interface using `exec.Command`.
    * **Batch Copy:** This is where you port the `go` command logic from your **Tauri `main.rs` file** (`detect_go_installation`, `run_go_code`, `build_go_file`, etc.) into Go. This D.R.Y.'s your logic, moving it from Rust into your core Go application.

3.  **Create Persistence Adapters:**
    * Create `pkg/adapters/driven/persistence/json_store.go`.
    * Implement the `ports.StateRepository` and `ports.SettingsRepository` interfaces.
    * This adapter will read/write state from a simple JSON file (e.g., `~/.ts2go/state.json`).
    * **This replaces the frontend `localStorage` logic** from `transpile.ts` and `settings.ts`, moving state management to the robust backend.

---

### Phase 3: Rewire the Driving Adapters (The Sockets)

**Goal:** Update your CLI and Tauri backend to use the new hexagonal core.

1.  **Refactor the CLI (`cmd/ts2go-cli/main.go`):**
    * This `main.go` file becomes your **Composition Root** (the place where you do Dependency Injection).
    * You will instantiate all your *concrete* adapters and *inject* them into your services.

    > **File: `cmd/ts2go-cli/main.go` (simplified)**
    ```go
    package main
    
    import (
        "ts2go/pkg/adapters/driven/filesystem"
        "ts2go/pkg/adapters/driven/gocompiler"
        "ts2go/pkg/adapters/driven/persistence"
        "ts2go/pkg/core/services"
        "ts2go/pkg/adapters/driving/cli" // Your new CLI handler
    )
    
    func main() {
        // --- 1. Create Concrete Adapters (The "Plugs") ---
        fsAdapter := filesystem.NewDiskFileSystem()
        goAdapter := gocompiler.NewOSGoCompiler()
        stateRepo := persistence.NewJSONStateRepository("~/.ts2go/state.json")
        settingsRepo := persistence.NewJSONSettingsRepository("~/.ts2go/settings.json")
        
        // --- 2. Create Core Services (The "Hexagon") ---
        // Inject the adapters into the services
        transpileSvc := services.NewTranspilationService(fsAdapter, goAdapter)
        runtimeSvc := services.NewGoRuntimeService(goAdapter)
        stateSvc := services.NewStateService(stateRepo, settingsRepo)
        
        // --- 3. Run the Driving Adapter (The "Socket") ---
        // Pass the core services to the CLI application
        cliApp := cli.NewCliApplication(transpileSvc, runtimeSvc, stateSvc)
        cliApp.Execute()
    }
    ```

2.  **Simplify CLI Commands:**
    * Your Cobra (or similar) command functions become simple one-liners.
    * The `transpile` command's `Run` function now just calls `cliApp.transpileSvc.TranspileProject(path)`.
    * The `run` command now just calls `cliApp.runtimeSvc.RunProject(path)`.

3.  **Gut the Tauri Backend (`desktop-ui/src-tauri/src/main.rs`):**
    * **Action:** Delete all the duplicated Go logic (`run_go_code`, `detect_go_installation`, etc.) from your Rust file.
    * **New Logic:** Your Tauri backend (`main.rs`) will *only* shell out to your newly-refactored `ts2go-cli`.
    * Add new CLI commands (e.g., `ts2go-cli detect-go`, `ts2go-cli state get-all`) to support your UI.
    * Your Tauri `#[tauri::command]` functions become simple wrappers.

    > **Example: `desktop-ui/src-tauri/src/main.rs` (After)**
    ```rust
    #[tauri::command]
    fn detect_go_installation() -> Result<String, String> {
        // Just call the CLI, which now holds all the logic
        let output = Command::new("ts2go-cli")
            .arg("detect-go")
            .output()
            .map_err(|e| e.to_string())?;
            
        Ok(String::from_utf8_lossy(&output.stdout).to_string())
    }
    
    #[tauri::command]
    fn run_go_project(path: String) -> Result<String, String> {
        let output = Command::new("ts2go-cli")
            .arg("run")
            .arg("--path")
            .arg(path)
            .output()
            .map_err(|e| e.to_string())?;
        
        Ok(String::from_utf8_lossy(&output.stdout).to_string())
    }
    ```

---

### Phase 4: Migrate Frontend State (The Final Polish)

**Goal:** Complete the state migration from the Vue frontend to the Go backend.

1.  **Refactor Vue Stores (`transpile.ts`, `settings.ts`):**
    * Delete all `localStorage` logic.
    * `loadTranspilationStates()` now calls the Tauri command `get_all_transpilation_states()`.
    * `saveSettings(settings)` now calls the Tauri command `save_settings(settings)`.
    * The frontend is now truly just a "view" of the backend state.

2.  **Add Supporting Tauri & CLI Commands:**
    * Add Tauri commands (`get_all_translilation_states`, `save_settings`) to `main.rs`.
    * Add the corresponding CLI commands (`ts2go-cli state list`, `ts2go-cli settings set`) that the Tauri commands will call.
    * These CLI commands will use the `stateSvc` you built in Phase 3, which uses the persistence adapter you built in Phase 2. The loop is complete.

---

### Phase 5: Prove the Architecture (The Web Service)

**Goal:** Prove the refactor was successful by adding a new "driving adapter" for your planned web service.

1.  **Create a New Main File:** `cmd/ts2go-web/main.go`.
2.  **Copy the Composition Root:** Copy the DI logic from the CLI's `main.go`. You will instantiate the *exact same services and adapters*.
3.  **Add a Web Server:**
    * Add a simple `net/http` or `gin-gonic` web server.
    * Create an endpoint `/transpile` that calls `transpileSvc.TranspileProject(...)`.
    * Create an endpoint `/run` that calls `runtimeSvc.RunProject(...)`.

You can now run `./ts2go-web` and you will have a fully functional web service sharing the *exact same* battle-tested core logic as your desktop application, fulfilling your primary goal.