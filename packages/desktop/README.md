# TS2Go Desktop Application

A standalone desktop application for transpiling TypeScript to Go with a beautiful, intuitive interface.

## Technology Stack

- **Vue 3** - Progressive JavaScript framework (Composition API)
- **PrimeVue 4** - Rich UI component library (90+ components)
- **Tailwind CSS 4** - Utility-first CSS framework
- **Vite 5** - Next-generation build tool
- **TypeScript 5** - Type-safe development

## Features

### Current (MVP)
- ✅ Drag & drop TypeScript files
- ✅ Split-pane editor (TypeScript | Go)
- ✅ File tree view
- ✅ Basic transpilation UI
- ✅ Copy & download functionality
- ✅ Dark/light theme toggle
- ✅ Settings panel

### Planned
- 🔄 Real-time transpilation (< 500ms)
- 🔄 AST tree viewer
- 🔄 Dependency graph visualization
- 🔄 Error highlighting with line numbers
- 🔄 Monaco Editor integration (VS Code editor)
- 🔄 Tauri integration for native functionality
- 🔄 Cross-platform builds (Windows, macOS, Linux)

## Development

### Prerequisites
- Node.js 18+ and npm/pnpm
- (Optional) Rust and Tauri CLI for native builds

### Install Dependencies
```bash
npm install
```

### Run Development Server
```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### Build for Production
```bash
npm run build
```

## Project Structure

```
packages/desktop/
├── src/
│   ├── components/       # Vue components
│   ├── composables/      # Composition API utilities
│   ├── stores/           # State management
│   ├── App.vue           # Main app component
│   ├── main.ts           # App entry point
│   └── style.css         # Global styles
├── public/               # Static assets
├── index.html            # HTML template
├── vite.config.ts        # Vite configuration
├── tailwind.config.js    # Tailwind CSS configuration
├── tsconfig.json         # TypeScript configuration
└── package.json          # Dependencies and scripts
```

## UI Components

### Main Layout
- Top toolbar with branding and controls
- Left sidebar with file tree
- Center split-pane editor
- Bottom status bar

### Key Components
- **FileUpload** - Drag & drop file upload (PrimeVue)
- **Splitter** - Resizable split panes (PrimeVue)
- **Textarea** - Code editors (will be replaced with Monaco)
- **Toolbar** - Action buttons
- **Sidebar** - Settings panel
- **TabView** - Settings tabs

## Roadmap

### Phase 1: MVP (Current)
- [x] Basic UI layout
- [x] File upload
- [x] Split editor
- [x] Placeholder transpilation
- [x] Copy/download functionality

### Phase 2: Core Features (Week 1-2)
- [ ] Integrate with CLI transpiler
- [ ] Real-time transpilation
- [ ] Error display
- [ ] Monaco Editor integration
- [ ] Syntax highlighting

### Phase 3: Advanced Features (Week 3-4)
- [ ] AST tree viewer
- [ ] Dependency graph
- [ ] Settings persistence
- [ ] Recent files
- [ ] Batch transpilation

### Phase 4: Native Build (Week 4-5)
- [ ] Tauri integration
- [ ] Native file system access
- [ ] System tray
- [ ] Auto-updates
- [ ] Cross-platform builds

## Contributing

See the main project README for contributing guidelines.

## License

Same as the main ts2go project.
