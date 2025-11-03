# Monorepo Migration Plan

**Date:** November 2, 2025  
**Goal:** Restructure ts2go into a clean monorepo for better organization

---

## 🎯 Objectives

1. Separate concerns (core transpiler, CLI, desktop app)
2. Enable independent development and versioning
3. Prepare for desktop app integration
4. Improve maintainability and collaboration
5. Support future extensions (web app, VS Code extension)

---

## 📁 Current Structure

```
ts2go/
├── cmd/ts2go/              # CLI entry point
├── internal/               # Core transpiler
│   ├── transpiler/
│   ├── analyzer/
│   ├── mapper/
│   ├── module/
│   ├── optimizer/
│   ├── orchestrator/
│   └── project/
├── pkg/                    # Public packages
│   ├── cli/
│   └── optimizer/
├── runtime/                # Go runtime libraries
├── tests/                  # Integration tests
├── examples/
├── docs/
└── go.work
```

---

## 📁 Target Monorepo Structure

```
ts2go/
├── README.md
├── LICENSE
├── .gitignore
├── go.work                         # Go workspace for all packages
├── Makefile                        # Root-level build commands
│
├── docs/                           # Shared documentation
│   ├── COMPREHENSIVE_ROADMAP.md
│   ├── MONOREPO_GUIDE.md
│   └── ...
│
├── packages/
│   │
│   ├── core/                       # Core transpiler engine
│   │   ├── README.md
│   │   ├── go.mod
│   │   ├── internal/
│   │   │   ├── transpiler/        # Code generation
│   │   │   ├── analyzer/          # Import analysis
│   │   │   ├── mapper/            # npm → Go mapping
│   │   │   ├── module/            # Module system
│   │   │   ├── optimizer/         # Code optimization
│   │   │   ├── orchestrator/      # Multi-file orchestration
│   │   │   └── project/           # Project scanning
│   │   ├── runtime/               # Go runtime libraries
│   │   │   ├── fs/
│   │   │   ├── path/
│   │   │   ├── console/
│   │   │   ├── process/
│   │   │   ├── os/
│   │   │   ├── http/
│   │   │   ├── url/
│   │   │   └── buffer/
│   │   └── mappings/              # npm-to-go.yaml
│   │
│   ├── cli/                        # Command-line interface
│   │   ├── README.md
│   │   ├── go.mod
│   │   ├── cmd/ts2go/             # Main entry point
│   │   │   ├── main.go
│   │   │   └── go.mod
│   │   └── pkg/cli/               # CLI utilities
│   │       ├── analyze.go
│   │       ├── transpile.go
│   │       ├── progress.go
│   │       └── watch.go
│   │
│   ├── desktop/                    # Tauri desktop application
│   │   ├── README.md
│   │   ├── package.json
│   │   ├── tauri.conf.json
│   │   ├── vite.config.ts
│   │   ├── tailwind.config.js
│   │   ├── tsconfig.json
│   │   │
│   │   ├── src-tauri/             # Rust backend
│   │   │   ├── Cargo.toml
│   │   │   ├── tauri.conf.json
│   │   │   ├── build.rs
│   │   │   └── src/
│   │   │       ├── main.rs
│   │   │       ├── transpiler.rs  # Bridge to CLI
│   │   │       ├── commands.rs    # Tauri commands
│   │   │       └── lib.rs
│   │   │
│   │   └── src/                   # Vue 3 frontend
│   │       ├── main.ts
│   │       ├── App.vue
│   │       ├── assets/
│   │       ├── components/
│   │       │   ├── Editor/
│   │       │   │   ├── TypeScriptEditor.vue
│   │       │   │   ├── GoEditor.vue
│   │       │   │   └── SplitView.vue
│   │       │   ├── FileTree/
│   │       │   │   ├── FileTree.vue
│   │       │   │   └── FileNode.vue
│   │       │   ├── Visualization/
│   │       │   │   ├── ASTViewer.vue
│   │       │   │   └── DependencyGraph.vue
│   │       │   ├── Settings/
│   │       │   │   └── SettingsPanel.vue
│   │       │   └── Common/
│   │       │       ├── ErrorDisplay.vue
│   │       │       └── ProgressBar.vue
│   │       ├── composables/
│   │       │   ├── useTranspiler.ts
│   │       │   ├── useFileSystem.ts
│   │       │   └── useSettings.ts
│   │       ├── stores/
│   │       │   ├── transpiler.ts
│   │       │   ├── files.ts
│   │       │   └── settings.ts
│   │       ├── types/
│   │       │   └── index.ts
│   │       └── utils/
│   │
│   └── shared/                     # Shared utilities (future)
│       └── types/
│
├── examples/                       # Example projects
│   ├── simple-api/
│   ├── cli-tool/
│   └── data-processor/
│
├── tests/                          # Integration tests
│   ├── fixtures/
│   └── integration_test.go
│
└── tools/                          # Build scripts
    ├── setup.sh
    └── release.sh
```

---

## 🔄 Migration Steps

### Phase 1: Prepare Core Package (Week 1, Days 1-2)

**Step 1.1: Create packages structure**
```bash
mkdir -p packages/core
```

**Step 1.2: Move core transpiler**
```bash
# Move internal and runtime
mv internal packages/core/
mv runtime packages/core/
mv mappings packages/core/
```

**Step 1.3: Create core go.mod**
```bash
cd packages/core
go mod init github.com/el-j/ts2go/packages/core
```

**Step 1.4: Update imports**
- Update all imports from `github.com/yourusername/ts2go/internal/*`
- To `github.com/el-j/ts2go/packages/core/internal/*`

---

### Phase 2: Create CLI Package (Week 1, Days 3-4)

**Step 2.1: Create CLI structure**
```bash
mkdir -p packages/cli/cmd/ts2go
mkdir -p packages/cli/pkg/cli
```

**Step 2.2: Move CLI code**
```bash
mv cmd/ts2go/* packages/cli/cmd/ts2go/
mv pkg/cli packages/cli/pkg/
```

**Step 2.3: Create CLI go.mod**
```bash
cd packages/cli
go mod init github.com/el-j/ts2go/packages/cli
```

**Step 2.4: Update CLI dependencies**
- Add dependency on `packages/core`
- Update import paths

---

### Phase 3: Update Go Workspace (Week 1, Day 5)

**Step 3.1: Update go.work**
```go
go 1.22.5

use (
	./packages/core
	./packages/cli
	./packages/cli/cmd/ts2go
	./tests
)
```

**Step 3.2: Test build**
```bash
# Build core
cd packages/core && go build ./...

# Build CLI
cd packages/cli && go build ./cmd/ts2go
```

---

### Phase 4: Setup Desktop App Skeleton (Week 2)

**Step 4.1: Install Tauri CLI**
```bash
cargo install tauri-cli
npm install -g @tauri-apps/cli
```

**Step 4.2: Create Tauri project**
```bash
cd packages
npm create tauri-app desktop
# Choose:
# - Package manager: npm
# - UI template: Vue (TypeScript)
# - Add plugins: none
```

**Step 4.3: Install dependencies**
```bash
cd desktop
npm install
npm install -D tailwindcss@latest postcss autoprefixer
npm install primevue@^4.0.0 primeicons
npm install @vueuse/core
npm install monaco-editor
npm install @vue/composition-api
```

**Step 4.4: Configure Tailwind**
```bash
npx tailwindcss init -p
```

**Step 4.5: Configure PrimeVue**
Update `src/main.ts`:
```typescript
import PrimeVue from 'primevue/config';
import 'primevue/resources/themes/aura-light-blue/theme.css';
import 'primeicons/primeicons.css';

app.use(PrimeVue);
```

---

### Phase 5: Desktop App - Core Features (Week 3)

**Step 5.1: Create basic layout**
- Top toolbar
- Left file tree
- Center split editor (TS | Go)
- Bottom status bar

**Step 5.2: File upload**
- Drag & drop support
- File browser
- Recent files

**Step 5.3: Transpilation bridge**
Create Rust commands in `src-tauri/src/transpiler.rs`:
```rust
use std::process::Command;

#[tauri::command]
pub fn transpile_code(typescript: String) -> Result<String, String> {
    // Call ../cli/cmd/ts2go/ts2go convert
    // Return Go code or error
}
```

**Step 5.4: Real-time transpilation**
- Debounced transpilation on change
- Error highlighting
- Syntax highlighting with Monaco Editor

---

### Phase 6: Desktop App - Advanced Features (Week 4)

**Step 6.1: AST Viewer**
- Tree component showing AST structure
- Clickable nodes
- Sync with editor

**Step 6.2: Settings panel**
- Module name
- Go version target
- Optimization level
- Custom mappings

**Step 6.3: Export functionality**
- Download single file
- Download project zip
- Copy to clipboard

**Step 6.4: Themes**
- Dark/light mode toggle
- Syntax theme selection

---

### Phase 7: Testing & Documentation (Week 5)

**Step 7.1: Test all packages**
```bash
# Test core
cd packages/core && go test ./...

# Test CLI
cd packages/cli && go test ./...

# Test desktop app
cd packages/desktop && npm run test
```

**Step 7.2: Update documentation**
- Update README files for each package
- Create MONOREPO_GUIDE.md
- Update build instructions
- Update contribution guide

**Step 7.3: CI/CD updates**
- Update GitHub Actions workflows
- Separate workflows for each package
- Build and release automation

---

## 🔧 Updated Build Commands

### Root Makefile
```makefile
.PHONY: all build-core build-cli build-desktop test clean

all: build-core build-cli

build-core:
	cd packages/core && go build ./...

build-cli:
	cd packages/cli && go build -o ../../ts2go ./cmd/ts2go

build-desktop:
	cd packages/desktop && npm run tauri build

test:
	cd packages/core && go test ./...
	cd packages/cli && go test ./...
	cd tests && go test ./...

clean:
	rm -f ts2go
	cd packages/desktop && npm run clean
```

---

## 📦 Package Dependencies

```
packages/cli
    └── depends on packages/core

packages/desktop
    └── calls packages/cli (via Rust subprocess)

packages/core
    └── no dependencies (pure Go)
```

---

## 🎨 Desktop App Tech Stack Details

### Frontend
- **Vue 3.4+** - Composition API, `<script setup>`
- **TypeScript 5.x** - Full type safety
- **Vite 5.x** - Fast build tool
- **PrimeVue 4.x** - UI components (Button, Tree, Panel, etc.)
- **Tailwind CSS 4.x** - Utility-first styling
- **Monaco Editor** - Code editor (VS Code editor)
- **VueUse** - Composition utilities

### Backend (Tauri)
- **Rust 1.75+** - Native performance
- **Tauri 2.x** - Desktop app framework
- **tokio** - Async runtime
- **serde** - Serialization

### Features Enabled by Stack
- Native file system access (Tauri FS API)
- Cross-platform (Windows, macOS, Linux)
- Small bundle size (~3-5 MB)
- Native performance
- Auto-updates support
- System tray integration
- Rich UI components from PrimeVue
- Beautiful styling with Tailwind

---

## 📋 Migration Checklist

### Pre-Migration
- [ ] Backup current repository
- [ ] Document current import paths
- [ ] Freeze feature development

### Migration
- [ ] Create packages/ directory
- [ ] Move core to packages/core
- [ ] Move CLI to packages/cli
- [ ] Update go.work
- [ ] Update all import paths
- [ ] Test builds

### Desktop App
- [ ] Create packages/desktop
- [ ] Set up Tauri + Vue + PrimeVue + Tailwind
- [ ] Implement file upload
- [ ] Implement transpilation bridge
- [ ] Create basic UI layout
- [ ] Add advanced features

### Post-Migration
- [ ] Update documentation
- [ ] Update CI/CD
- [ ] Test all workflows
- [ ] Create release

---

## 🚀 Benefits After Migration

1. **Clear Organization**
   - Each package has clear purpose
   - Easy to find code
   - Logical structure

2. **Independent Development**
   - Core team works on core
   - CLI team works on CLI
   - UI team works on desktop
   - No conflicts

3. **Better Testing**
   - Test each package independently
   - Faster test cycles
   - Clear test organization

4. **Future-Proof**
   - Easy to add new packages
   - VS Code extension
   - Web app
   - Language server

5. **Professional Structure**
   - Industry-standard monorepo
   - Easy for contributors
   - Better documentation

---

**Timeline:** 5 weeks total  
**Next Step:** Start Phase 1 migration
