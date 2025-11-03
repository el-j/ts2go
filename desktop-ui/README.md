# TS2Go Desktop Application

**Phase 21 Implementation** - Tauri-based desktop UI for the TS2Go TypeScript to Go transpiler.

## Technology Stack

- **Tauri 2.0** - Native desktop application framework
- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe development
- **PrimeVue 4** - Rich UI component library
- **Tailwind CSS 4** - Utility-first CSS framework
- **Pinia** - State management
- **Vite** - Fast build tool

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
# Build for production
npm run tauri:build
```

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

## Features

### Week 1 (Current) - Foundation
- ✅ Tauri project setup
- ✅ Vue 3 + TypeScript configuration
- ✅ PrimeVue 4 integration
- ✅ Tailwind CSS 4 setup
- ✅ Base layout with sidebar navigation
- ✅ Vue Router with multiple views
- ✅ Pinia stores (project, transpiler, settings, logs)
- ⏳ Project selection UI (in progress)

### Week 2 - Editor & Transpilation
- Code editor integration
- Split-pane editor
- CLI integration
- Progress tracking

### Week 3 - Advanced Features
- Watch mode
- Dependency visualization
- Settings panel

### Week 4 - Polish & Testing
- Examples & templates
- Testing & documentation
- Build & package

## Documentation

- [Comprehensive Plan](../docs/PHASE21_COMPREHENSIVE_PLAN.md)
- [TODO Checklist](../docs/PHASE21_TODO.md)
- [Project Status](../docs/STATUS.md)

## License

MIT - See parent repository for details
