# TS2Go Desktop - Developer Guide

**Version:** 0.1.0  
**Last Updated:** November 4, 2025

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Project Structure](#project-structure)
3. [Development Setup](#development-setup)
4. [Building and Testing](#building-and-testing)
5. [Contributing](#contributing)
6. [Code Style](#code-style)

---

## Architecture Overview

TS2Go Desktop is built using the following technologies:

### Frontend Stack
- **Vue 3** - Progressive JavaScript framework (Composition API)
- **TypeScript** - Type-safe JavaScript
- **PrimeVue 4** - UI component library
- **Tailwind CSS 4** - Utility-first CSS framework
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Monaco Editor** - Code editor (VS Code engine)

### Backend Stack
- **Tauri 2** - Cross-platform desktop framework (Rust)
- **Rust** - Systems programming language for native backend

### Build Tools
- **Vite 5** - Fast build tool and dev server
- **Vitest** - Unit testing framework
- **Vue TSC** - TypeScript type checking for Vue

---

## Project Structure

```
desktop-ui/
├── src/
│   ├── components/          # Vue components
│   │   ├── CodeEditor.vue
│   │   ├── LogViewer.vue
│   │   ├── ExampleGallery.vue
│   │   ├── KeyboardShortcutsDialog.vue
│   │   └── __tests__/       # Component tests
│   ├── views/               # Page views
│   │   ├── HomeView.vue
│   │   ├── ProjectView.vue
│   │   ├── EditorView.vue
│   │   ├── ExamplesView.vue
│   │   └── SettingsView.vue
│   ├── stores/              # Pinia stores
│   │   ├── project.ts
│   │   ├── transpiler.ts
│   │   ├── settings.ts
│   │   ├── logs.ts
│   │   └── __tests__/       # Store tests
│   ├── composables/         # Vue composables
│   │   ├── useTheme.ts
│   │   └── useKeyboardShortcuts.ts
│   ├── router/              # Vue Router config
│   │   └── index.ts
│   ├── test/                # Test utilities
│   │   ├── setup.ts
│   │   └── monaco-mock.ts
│   ├── App.vue              # Root component
│   └── main.ts              # Application entry point
├── src-tauri/               # Rust backend
│   ├── src/
│   │   ├── commands.rs      # Tauri commands
│   │   └── main.rs          # Rust entry point
│   ├── tauri.conf.json      # Tauri configuration
│   └── Cargo.toml           # Rust dependencies
├── public/                  # Static assets
├── index.html               # HTML entry point
├── package.json             # Node.js dependencies
├── tsconfig.json            # TypeScript config
├── vite.config.ts           # Vite configuration
├── vitest.config.ts         # Vitest configuration
├── tailwind.config.js       # Tailwind CSS config
├── USER_GUIDE.md            # User documentation
├── DEVELOPER_GUIDE.md       # This file
└── README.md                # Project README
```

---

## Development Setup

### Prerequisites

- **Node.js** 18+ and npm
- **Rust** 1.70+ (for Tauri)
- **Operating System**:
  - Windows: Visual Studio C++ Build Tools
  - macOS: Xcode Command Line Tools
  - Linux: Development packages (see Tauri docs)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/el-j/ts2go.git
   cd ts2go/desktop-ui
   ```

2. **Install Node.js dependencies**
   ```bash
   npm install
   ```

3. **Install Rust dependencies** (Tauri will handle this)
   ```bash
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   ```

4. **Install Tauri CLI** (if not already installed)
   ```bash
   npm install -g @tauri-apps/cli
   ```

### Running Development Server

```bash
# Run Vite dev server only (frontend)
npm run dev

# Run Tauri dev mode (frontend + backend)
npm run tauri:dev
```

The dev server will:
- Hot-reload on file changes
- Open at `http://localhost:5173`
- Launch Tauri window (in tauri:dev mode)

---

## Building and Testing

### Running Tests

```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run tests with UI
npm run test:ui
```

### Test Structure

Tests are organized by component/store:
- **Component tests**: `src/components/__tests__/*.spec.ts`
- **Store tests**: `src/stores/__tests__/*.spec.ts`
- **Integration tests**: (Coming in Week 4)

Example test:
```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CodeEditor from '../CodeEditor.vue'

describe('CodeEditor', () => {
  it('should render', () => {
    const wrapper = mount(CodeEditor)
    expect(wrapper.exists()).toBe(true)
  })
})
```

### Building for Production

```bash
# Build Vite production bundle
npm run build

# Build Tauri application for current platform
npm run tauri:build

# Build for specific platform
npm run tauri:build -- --target x86_64-pc-windows-msvc  # Windows
npm run tauri:build -- --target x86_64-apple-darwin     # macOS
npm run tauri:build -- --target x86_64-unknown-linux-gnu # Linux
```

Build artifacts:
- **Windows**: `src-tauri/target/release/bundle/nsis/*.exe`
- **macOS**: `src-tauri/target/release/bundle/dmg/*.dmg`
- **Linux**: `src-tauri/target/release/bundle/deb/*.deb`

### Linting and Formatting

```bash
# Run ESLint
npm run lint

# Format code with Prettier
npm run format

# Type check with TypeScript
npm run type-check
```

---

## Contributing

### Development Workflow

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write code following style guide
   - Add tests for new features
   - Update documentation

3. **Run tests and linting**
   ```bash
   npm test
   npm run lint
   npm run type-check
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   Follow conventional commits:
   - `feat:` - New features
   - `fix:` - Bug fixes
   - `docs:` - Documentation changes
   - `style:` - Code style changes
   - `refactor:` - Code refactoring
   - `test:` - Test additions or changes
   - `chore:` - Build/tool changes

5. **Push and create Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```

### Adding New Components

1. **Create component file**
   ```bash
   # In src/components/
   touch MyComponent.vue
   ```

2. **Write component with TypeScript**
   ```vue
   <script setup lang="ts">
   import { ref } from 'vue'
   
   interface Props {
     title: string
   }
   
   const props = defineProps<Props>()
   const count = ref(0)
   </script>
   
   <template>
     <div>
       <h2>{{ props.title }}</h2>
       <p>Count: {{ count }}</p>
     </div>
   </template>
   ```

3. **Create test file**
   ```bash
   # In src/components/__tests__/
   touch MyComponent.spec.ts
   ```

4. **Write tests**
   ```typescript
   import { describe, it, expect } from 'vitest'
   import { mount } from '@vue/test-utils'
   import MyComponent from '../MyComponent.vue'
   
   describe('MyComponent', () => {
     it('renders title', () => {
       const wrapper = mount(MyComponent, {
         props: { title: 'Test Title' }
       })
       expect(wrapper.text()).toContain('Test Title')
     })
   })
   ```

### Adding New Stores

1. **Create store file**
   ```bash
   # In src/stores/
   touch myStore.ts
   ```

2. **Define store with Pinia**
   ```typescript
   import { defineStore } from 'pinia'
   import { ref } from 'vue'
   
   export const useMyStore = defineStore('myStore', () => {
     const data = ref<string[]>([])
     
     function addItem(item: string) {
       data.value.push(item)
     }
     
     return {
       data,
       addItem
     }
   })
   ```

3. **Create test file and write tests**
   ```typescript
   import { describe, it, expect, beforeEach } from 'vitest'
   import { setActivePinia, createPinia } from 'pinia'
   import { useMyStore } from '../myStore'
   
   describe('My Store', () => {
     beforeEach(() => {
       setActivePinia(createPinia())
     })
     
     it('adds items', () => {
       const store = useMyStore()
       store.addItem('test')
       expect(store.data).toContain('test')
     })
   })
   ```

### Adding Tauri Commands

1. **Define command in Rust** (`src-tauri/src/commands.rs`)
   ```rust
   #[tauri::command]
   pub async fn my_command(input: String) -> Result<String, String> {
       Ok(format!("Processed: {}", input))
   }
   ```

2. **Register command** (`src-tauri/src/main.rs`)
   ```rust
   fn main() {
       tauri::Builder::default()
           .invoke_handler(tauri::generate_handler![
               my_command
           ])
           .run(tauri::generate_context!())
           .expect("error while running tauri application");
   }
   ```

3. **Call from Vue**
   ```typescript
   import { invoke } from '@tauri-apps/api/core'
   
   async function callMyCommand() {
     const result = await invoke<string>('my_command', {
       input: 'Hello'
     })
     console.log(result) // "Processed: Hello"
   }
   ```

---

## Code Style

### TypeScript/Vue Style

- Use **Composition API** with `<script setup>`
- Prefer **`const`** over `let`
- Use **TypeScript interfaces** for props and data
- Follow **Vue 3 best practices**
- Use **camelCase** for variables and functions
- Use **PascalCase** for components and types

Example:
```vue
<script setup lang="ts">
import { ref, computed } from 'vue'

interface User {
  id: number
  name: string
}

const users = ref<User[]>([])

const userCount = computed(() => users.value.length)

function addUser(user: User) {
  users.value.push(user)
}
</script>
```

### CSS/Tailwind Style

- Use **Tailwind utility classes** first
- Create **custom classes** only when necessary
- Use **dark:** variants for dark mode
- Follow **mobile-first** responsive design

Example:
```vue
<template>
  <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md">
    <h2 class="text-xl font-semibold mb-2">Title</h2>
    <p class="text-gray-600 dark:text-gray-400">Content</p>
  </div>
</template>
```

### Rust Style

- Follow **Rust standard conventions**
- Use **snake_case** for functions and variables
- Use **PascalCase** for types and structs
- Add **error handling** with `Result<T, E>`
- Document **public functions**

Example:
```rust
/// Processes the input and returns a result
#[tauri::command]
pub async fn process_data(input: String) -> Result<String, String> {
    match validate_input(&input) {
        Ok(data) => Ok(format!("Processed: {}", data)),
        Err(e) => Err(format!("Validation failed: {}", e))
    }
}
```

---

## Debugging

### Frontend Debugging

**Chrome DevTools:**
1. Open dev server: `npm run dev`
2. Open `http://localhost:5173`
3. Press F12 to open DevTools
4. Use Vue DevTools extension

**VS Code:**
1. Install "Volar" extension
2. Use built-in debugger
3. Set breakpoints in `.vue` files

### Backend Debugging

**Rust debugging:**
1. Add `println!` or `dbg!` macros
2. Check logs in terminal
3. Use `cargo check` for compilation errors

**Tauri debugging:**
1. Enable dev console in Tauri app
2. Check stdout/stderr in terminal
3. Use `console.log` in Rust via `tauri::api::dialog`

---

## Performance Optimization

### Frontend

- **Lazy load** routes with `import()`
- **Debounce** expensive operations
- **Virtualize** long lists
- **Optimize** images and assets
- **Use** `v-memo` for static content

### Backend

- **Avoid** blocking operations
- **Use** async/await
- **Batch** file operations
- **Cache** repeated computations

---

## Resources

### Documentation
- [Vue 3 Docs](https://vuejs.org/)
- [PrimeVue Docs](https://primevue.org/)
- [Tauri Docs](https://tauri.app/)
- [Tailwind CSS Docs](https://tailwindcss.com/)
- [Vitest Docs](https://vitest.dev/)

### Tools
- [Vue DevTools](https://devtools.vuejs.org/)
- [Rust Analyzer](https://rust-analyzer.github.io/)
- [TypeScript Playground](https://www.typescriptlang.org/play)

---

## Support

- **Issues**: [GitHub Issues](https://github.com/el-j/ts2go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/el-j/ts2go/discussions)
- **Documentation**: [Project Docs](../docs/)

---

**Happy coding!** 🚀
