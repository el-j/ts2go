# TS2Go Desktop Application

**Phase 21 Implementation** - Modern desktop UI for the TS2Go TypeScript to Go transpiler.

## ✨ Features

- 🎨 **Modern UI** - Beautiful interface with dark/light themes
- ✨ **Real-time Transpilation** - Instant feedback as you code
- 📖 **Example Gallery** - 6 built-in examples to learn from
- ⌨️ **Keyboard Shortcuts** - Work faster with hotkeys
- ⚙️ **Customizable Settings** - Personalize your experience
- 📊 **Visual Feedback** - Progress tracking and detailed logs
- 🌓 **Theme Support** - Dark, light, and system themes
- 💾 **Auto-save** - Never lose your work

## 🚀 Quick Start

```bash
# From repository root:

# Build desktop app (recommended - handles all dependencies)
make build-desktop

# OR manually:

# 1. Install desktop dependencies
make install-desktop

# 2. Run development server
make dev-desktop

# 3. Build production version
make build-desktop
```

**Note:** The Makefile automatically builds the CLI binary and copies it to the correct location (`desktop-ui/src-tauri/bin/ts2go-cli`) before starting Tauri. This is required for the desktop app to function.

## 📚 Documentation

- [User Guide](USER_GUIDE.md) - How to use the desktop app
- [Developer Guide](DEVELOPER_GUIDE.md) - Contributing and development
- [Comprehensive Plan](../docs/PHASE21_COMPREHENSIVE_PLAN.md) - Full feature roadmap
- [TODO Checklist](../docs/PHASE21_TODO.md) - Implementation progress

## 🛠️ Technology Stack

- **Tauri 2.0** - Native desktop application framework
- **Vue 3** - Progressive JavaScript framework (Composition API)
- **TypeScript** - Type-safe development
- **PrimeVue 4** - Rich UI component library
- **Tailwind CSS 4** - Utility-first CSS framework
- **Pinia** - State management
- **Monaco Editor** - VS Code editor engine
- **Vite 5** - Fast build tool
- **Vitest** - Unit testing framework

## Development

### Prerequisites

- Node.js 20+ and npm
- Rust 1.70+
- Tauri CLI

### Setup

```bash
# Install dependencies
npm install

# Run development server (Vue only)
npm run dev

# Run Tauri development mode
npm run tauri:dev
```

### Build

```bash
# From repository root (recommended):
make build-desktop

# OR manually from desktop-ui directory:
# 1. First, build CLI from root: cd .. && make build-cli
# 2. Copy CLI to bin: mkdir -p src-tauri/bin && cp ../ts2go src-tauri/bin/ts2go-cli
# 3. Then build desktop: npm run tauri:build
```

**Important:** The desktop app bundles the CLI binary (`ts2go-cli`) as a resource. The Makefile handles this automatically, but if building manually, ensure the CLI is in `src-tauri/bin/` before building.

## Project Structure

```
desktop-ui/
├── src/                    # Vue source code
│   ├── components/         # Reusable Vue components
│   ├── views/              # Page views
│   ├── stores/             # Pinia stores
│   ├── services/           # Business logic
│   ├── router/             # Vue Router configuration
│   ├── assets/             # Static assets
│   ├── App.vue             # Root component
│   └── main.ts             # Application entry point
├── src-tauri/              # Rust backend
│   ├── src/                # Rust source code
│   │   └── main.rs         # Tauri entry point
│   ├── Cargo.toml          # Rust dependencies
│   └── tauri.conf.json     # Tauri configuration
├── public/                 # Public assets
├── index.html              # HTML entry point
├── package.json            # Node.js dependencies
├── vite.config.ts          # Vite configuration
├── tailwind.config.js      # Tailwind configuration
└── tsconfig.json           # TypeScript configuration
```

## 📦 Implementation Status

**Current Progress: 85%** (Phase 21 - Week 3)

### ✅ Completed Features

**Week 1 - Foundation (100%)**
- Tauri project setup with Rust backend
- Vue 3 + TypeScript + Vite configuration
- PrimeVue 4 & Tailwind CSS 4 integration
- Base layout with sidebar navigation
- Vue Router with multiple views
- Pinia stores (project, transpiler, settings, logs)
- 42 unit tests passing

**Week 2 - Editor & Transpilation (80%)**
- Monaco Editor integration
- Split-pane code editor (TypeScript → Go)
- Tauri CLI integration
- Real-time transpilation
- Progress tracking
- Log viewer with filtering
- Error handling

**Week 3 - Advanced Features (40%)**
- ✅ **Settings Panel** - Comprehensive customization
  - Application settings (theme, font size, auto-save)
  - Project settings (module name, patterns)
  - Editor settings (tab size, word wrap, line numbers)
  - Settings persistence with localStorage
  
- ✅ **Keyboard Shortcuts** - Work faster with hotkeys
  - Core shortcuts (Ctrl+S, Ctrl+L, Ctrl+/)
  - Shortcuts help dialog
  - Integration in Editor view

- ⏳ **Watch Mode** - Auto-transpile on file changes (Planned)
- ⏳ **Analyze Command** - Dependency analysis UI (Planned)
- ⏳ **Build History** - Track transpilation history (Planned)

**Week 4 - Polish & Release (30%)**
- ✅ **Example Gallery** - 6 built-in examples
  - Basic: Interface, Function, Enum, Control Flow
  - Intermediate: Class with Inheritance
  - Advanced: Async/Await
  - Category filtering & preview
  
- ✅ **Theme Customization** - Dark/Light/System themes
  - Theme composable with auto-detection
  - System theme listening
  - Tailwind dark mode support

- ✅ **Documentation**
  - User Guide (complete)
  - Developer Guide (complete)
  - Updated README

- ⏳ **Testing** - Expand coverage (In Progress)
- ⏳ **Build & Package** - Multi-platform builds (Planned)

### 🎯 Upcoming Features

- Watch Mode implementation
- Analyze Command UI
- Build History tracking
- E2E tests with Playwright
- Production builds for Windows, macOS, Linux
- Tutorial mode for first-time users

## 🧪 Testing

```bash
# Run all tests
npm test

# Run tests in watch mode
npm test -- --watch

# Run tests with coverage
npm run test:coverage

# Run tests with UI
npm run test:ui
```

**Current Test Coverage:**
- ✅ 42 tests passing
- ✅ Stores: 100% coverage
- ✅ Components: Core components tested
- ⏳ E2E: Planned for Week 4

## License

MIT - See parent repository for details
