# TS2Go Desktop Application Specification

**Version:** 1.0  
**Date:** November 2, 2025  
**Technology:** Tauri 2 + Vue 3 + PrimeVue 4 + Tailwind CSS 4

---

## 🎯 Overview

A standalone desktop application for transpiling TypeScript to Go with a beautiful, intuitive interface. Built with modern web technologies wrapped in a native desktop app.

---

## 🛠️ Technology Stack

### Desktop Framework
- **Tauri 2.x** - Rust-based desktop framework
  - Native performance
  - Small bundle size (~3-5MB)
  - Cross-platform (Windows, macOS, Linux)
  - Native file system access
  - System integration

### Frontend
- **Vue 3.4+** - Progressive JavaScript framework
  - Composition API with `<script setup>`
  - Reactive state management
  - Component-based architecture

- **TypeScript 5.x** - Type-safe development
  - Full IntelliSense support
  - Compile-time type checking
  - Better refactoring

- **Vite 5.x** - Next-generation build tool
  - Lightning-fast HMR
  - Optimized builds
  - Plugin ecosystem

- **PrimeVue 4.x** - Rich UI component library
  - 90+ components
  - Beautiful themes
  - Accessibility built-in
  - Tree, Panel, Button, Dialog, etc.

- **Tailwind CSS 4.x** - Utility-first CSS
  - Rapid styling
  - Consistent design system
  - Dark mode support
  - Custom theme

- **Monaco Editor** - Code editor
  - VS Code's editor engine
  - TypeScript/Go syntax highlighting
  - IntelliSense
  - Diff viewer

- **VueUse** - Composition utilities
  - File system
  - Storage
  - UI utilities

---

## 🎨 User Interface Design

### Main Layout

```
┌─────────────────────────────────────────────────────┐
│ [🏠 TS2Go]  File  Edit  View  Help    [⚙️] [◐] [✕] │ Top Bar
├───────────┬─────────────────────────────────────────┤
│           │                                         │
│  File     │         TypeScript Code                 │
│  Tree     │  ┌──────────────────────────────────┐  │
│           │  │ interface User {                 │  │
│  📁 src/  │  │   name: string;                  │  │ Main
│    📄 a.ts│  │   age: number;                   │  │ Content
│    📄 b.ts│  │ }                                │  │
│           │  └──────────────────────────────────┘  │
│  📂 test/ │                                         │
│    📄 *.ts│         Generated Go Code               │
│           │  ┌──────────────────────────────────┐  │
│           │  │ type User struct {               │  │
│  [Upload] │  │   Name string `json:"name"`     │  │
│  [Clear]  │  │   Age  float64 `json:"age"`     │  │
│           │  │ }                                │  │
│           │  └──────────────────────────────────┘  │
├───────────┴─────────────────────────────────────────┤
│ ⚡ Ready  |  0 files  |  0 errors  |  0 warnings    │ Status Bar
└─────────────────────────────────────────────────────┘
```

### Color Scheme

**Light Theme (Aura Light Blue)**
- Primary: Blue (#3B82F6)
- Background: White (#FFFFFF)
- Surface: Gray-50 (#F9FAFB)
- Text: Gray-900 (#111827)
- Border: Gray-200 (#E5E7EB)

**Dark Theme (Aura Dark Blue)**
- Primary: Blue (#60A5FA)
- Background: Gray-900 (#111827)
- Surface: Gray-800 (#1F2937)
- Text: Gray-100 (#F3F4F6)
- Border: Gray-700 (#374151)

---

## 📋 Features Specification

### 1. File Management

#### 1.1 File Upload
**Component:** FileUploader.vue
```vue
<FileUpload
  mode="advanced"
  :multiple="true"
  accept=".ts,.tsx"
  :maxFileSize="10000000"
  @upload="handleUpload"
>
  <template #empty>
    <p>Drag and drop TypeScript files here</p>
  </template>
</FileUpload>
```

**Features:**
- Drag & drop files
- Multi-file selection
- File size limit: 10MB per file
- Supported: `.ts`, `.tsx`, `.d.ts`
- Progress indicator
- Error handling

#### 1.2 File Tree
**Component:** FileTree.vue (PrimeVue Tree)
```vue
<Tree
  :value="files"
  selectionMode="single"
  @node-select="onFileSelect"
>
  <template #default="{ node }">
    <i :class="node.icon"></i>
    <span>{{ node.label }}</span>
  </template>
</Tree>
```

**Features:**
- Hierarchical file structure
- Folder expansion/collapse
- File selection
- Context menu (right-click)
  - Transpile file
  - Remove file
  - View in explorer
- Icons for file types
- Drag to reorder

#### 1.3 Recent Projects
**Component:** RecentProjects.vue
```vue
<DataTable :value="recentProjects">
  <Column field="name" header="Project"></Column>
  <Column field="date" header="Last Opened"></Column>
  <Column header="Actions">
    <template #body="{ data }">
      <Button @click="openProject(data)" />
    </template>
  </Column>
</DataTable>
```

**Features:**
- Last 10 projects
- Quick open
- Clear history
- Stored in local storage

---

### 2. Code Editors

#### 2.1 Split-Pane Editor
**Component:** SplitEditor.vue
```vue
<Splitter layout="horizontal">
  <SplitterPanel :size="50">
    <MonacoEditor
      v-model="typescript"
      language="typescript"
      theme="vs-dark"
      @change="onCodeChange"
    />
  </SplitterPanel>
  <SplitterPanel :size="50">
    <MonacoEditor
      v-model="goCode"
      language="go"
      theme="vs-dark"
      :readonly="true"
    />
  </SplitterPanel>
</Splitter>
```

**Features:**
- Resizable split panes
- TypeScript editor (left)
  - Syntax highlighting
  - Auto-completion
  - Error squiggles
  - Line numbers
  - Minimap
- Go editor (right)
  - Syntax highlighting
  - Read-only
  - Line numbers
  - Minimap
- Sync scroll (optional)
- Diff mode toggle

#### 2.2 Editor Toolbar
**Component:** EditorToolbar.vue
```vue
<Toolbar>
  <template #start>
    <Button icon="pi pi-play" label="Transpile" />
    <Button icon="pi pi-copy" label="Copy Go" />
  </template>
  <template #end>
    <Button icon="pi pi-download" label="Download" />
    <SelectButton :options="views" />
  </template>
</Toolbar>
```

**Actions:**
- Transpile now
- Copy Go code
- Download Go file
- View mode: Split | TypeScript | Go | Diff

---

### 3. Transpilation

#### 3.1 Real-Time Transpilation
**Composable:** useTranspiler.ts
```typescript
export function useTranspiler() {
  const typescript = ref('')
  const goCode = ref('')
  const errors = ref([])
  const isTranspiling = ref(false)
  
  const transpile = useDebounceFn(async () => {
    isTranspiling.value = true
    try {
      const result = await invoke('transpile_code', {
        typescript: typescript.value
      })
      goCode.value = result.go
      errors.value = result.errors
    } catch (e) {
      errors.value = [e.message]
    } finally {
      isTranspiling.value = false
    }
  }, 500)
  
  watch(typescript, transpile)
  
  return { typescript, goCode, errors, isTranspiling }
}
```

**Features:**
- Debounced (500ms) transpilation
- Loading indicator
- Error display
- Success feedback

#### 3.2 Batch Transpilation
**Component:** BatchTranspiler.vue

**Features:**
- Transpile all files
- Progress bar
- File-by-file status
- Parallel processing
- Error summary

#### 3.3 Transpilation Options
**Component:** TranspileOptions.vue (in Settings)
```vue
<Panel header="Transpilation Options">
  <div class="field">
    <label>Module Name</label>
    <InputText v-model="options.moduleName" />
  </div>
  <div class="field">
    <label>Optimization</label>
    <Dropdown v-model="options.optimization" :options="['none', 'basic', 'aggressive']" />
  </div>
  <div class="field">
    <label>Error Handling</label>
    <SelectButton v-model="options.errorHandling" :options="['panic', 'return']" />
  </div>
</Panel>
```

---

### 4. Visualization

#### 4.1 AST Viewer
**Component:** ASTViewer.vue
```vue
<Tree
  :value="astTree"
  :expandedKeys="expandedKeys"
  selectionMode="single"
  @node-select="onNodeSelect"
>
  <template #default="{ node }">
    <span :class="'ast-' + node.kind">
      {{ node.kind }}
      <span class="text-xs">{{ node.text }}</span>
    </span>
  </template>
</Tree>
```

**Features:**
- Collapsible tree structure
- Syntax-highlighted nodes
- Click to jump to code
- Filter by node type
- Export as JSON

#### 4.2 Dependency Graph
**Component:** DependencyGraph.vue

**Features:**
- Visual graph of imports
- D3.js or similar
- Interactive nodes
- Zoom/pan
- Circular dependency highlighting
- Export as SVG

#### 4.3 Type Mapping Preview
**Component:** TypeMappingPreview.vue
```vue
<DataTable :value="typeMappings">
  <Column field="typescript" header="TypeScript Type"></Column>
  <Column field="go" header="Go Type"></Column>
  <Column field="complexity" header="Complexity"></Column>
</DataTable>
```

**Shows:**
- Interface → Struct mappings
- Type aliases
- Union type handling
- Array mappings

---

### 5. Error Display

#### 5.1 Error Panel
**Component:** ErrorPanel.vue
```vue
<Panel header="Errors & Warnings" :collapsed="errors.length === 0">
  <Message
    v-for="error in errors"
    :key="error.id"
    :severity="error.severity"
    :closable="true"
  >
    <template #default>
      <div>
        <strong>{{ error.file }}:{{ error.line }}</strong>
        <p>{{ error.message }}</p>
        <Button label="Jump to code" size="small" @click="jumpToError(error)" />
      </div>
    </template>
  </Message>
</Panel>
```

**Features:**
- Error severity (error, warning, info)
- File name and line number
- Clear error message
- Jump to code location
- Suggested fixes
- Dismissible
- Copy error details

---

### 6. Settings

#### 6.1 Settings Panel
**Component:** SettingsPanel.vue (Sidebar)
```vue
<Sidebar v-model:visible="visible" position="right" :style="{width: '30rem'}">
  <TabView>
    <TabPanel header="General">
      <!-- Theme, language, auto-save -->
    </TabPanel>
    <TabPanel header="Transpiler">
      <!-- Module name, optimization, etc. -->
    </TabPanel>
    <TabPanel header="Editor">
      <!-- Font size, theme, keybindings -->
    </TabPanel>
    <TabPanel header="Advanced">
      <!-- Custom mappings, debug mode -->
    </TabPanel>
  </TabView>
</Sidebar>
```

**Settings:**

**General**
- Theme: Light / Dark / System
- Language: English (future: more)
- Auto-save: On / Off
- Auto-transpile: On / Off

**Transpiler**
- Module name
- Go version target
- Optimization level: None / Basic / Aggressive
- Error handling: Panic / Return errors
- Include comments

**Editor**
- Font family
- Font size: 12-24px
- Show minimap: On / Off
- Word wrap: On / Off
- Line numbers: On / Off

**Advanced**
- Custom package mappings
- Debug mode
- Log level
- Experimental features

---

### 7. Export & Download

#### 7.1 Single File Export
**Features:**
- Download as `.go` file
- Copy to clipboard
- Save to custom location

#### 7.2 Project Export
**Features:**
- Export entire project structure
- Generate `go.mod`
- Create directory structure
- ZIP download
- Direct save to folder

#### 7.3 Configuration Export
**Features:**
- Save transpilation config as JSON
- Import config
- Share configs

---

### 8. Status & Feedback

#### 8.1 Status Bar
**Component:** StatusBar.vue
```vue
<div class="status-bar">
  <span>⚡ {{ status }}</span>
  <span>📁 {{ fileCount }} files</span>
  <span>{{ errorCount }} errors</span>
  <span>{{ warningCount }} warnings</span>
  <span>🕐 Last transpiled: {{ lastTranspile }}</span>
</div>
```

**Shows:**
- Current status (Ready, Transpiling, Error)
- File count
- Error/warning count
- Last transpile time
- Memory usage (optional)

#### 8.2 Toast Notifications
```vue
<Toast position="top-right" />
```

**Notifications:**
- Transpilation success
- Errors occurred
- File saved
- Settings updated

---

## 🔌 Tauri Commands (Rust Backend)

### File Operations
```rust
#[tauri::command]
async fn read_file(path: String) -> Result<String, String>

#[tauri::command]
async fn write_file(path: String, content: String) -> Result<(), String>

#[tauri::command]
async fn pick_folder() -> Result<String, String>
```

### Transpilation
```rust
#[tauri::command]
async fn transpile_code(
    typescript: String,
    options: TranspileOptions
) -> Result<TranspileResult, String>

#[tauri::command]
async fn transpile_file(
    path: String,
    options: TranspileOptions
) -> Result<TranspileResult, String>

#[tauri::command]
async fn transpile_project(
    project_path: String,
    options: TranspileOptions
) -> Result<ProjectTranspileResult, String>
```

### Analysis
```rust
#[tauri::command]
async fn analyze_project(path: String) -> Result<ProjectAnalysis, String>

#[tauri::command]
async fn get_ast(typescript: String) -> Result<ASTNode, String>

#[tauri::command]
async fn get_dependencies(path: String) -> Result<Vec<Dependency>, String>
```

---

## 📱 User Workflows

### Workflow 1: Quick Transpilation
1. Launch app
2. Drag TypeScript file
3. See Go code instantly
4. Download or copy

### Workflow 2: Project Transpilation
1. Launch app
2. Click "Open Project"
3. Select folder
4. See file tree
5. Click "Transpile All"
6. Review results
7. Export project

### Workflow 3: Interactive Development
1. Open file
2. Edit TypeScript in left pane
3. See real-time Go updates in right pane
4. Fix errors as they appear
5. Download when satisfied

---

## 🚀 Performance Targets

- **Cold Start:** < 2 seconds
- **Small File Transpilation:** < 100ms
- **Medium File (1000 lines):** < 500ms
- **Large File (5000 lines):** < 2 seconds
- **Project (50 files):** < 10 seconds
- **Memory Usage:** < 200MB idle
- **Bundle Size:** < 5MB

---

## 🧪 Testing Strategy

### Unit Tests
- Vue components (Vitest)
- Composables
- Utilities

### Integration Tests
- File operations
- Transpilation pipeline
- Settings persistence

### E2E Tests
- Playwright for Tauri
- User workflows
- Cross-platform

---

## 📦 Distribution

### Platforms
- **Windows:** `.exe` installer (NSIS)
- **macOS:** `.dmg` and `.app` bundle
- **Linux:** `.AppImage` and `.deb`

### Auto-Updates
- GitHub Releases integration
- Silent updates
- Update notifications

### Installation Size
- Windows: ~4-6 MB
- macOS: ~3-5 MB
- Linux: ~4-6 MB

---

## 🗺️ Development Roadmap

### MVP (Week 1-2)
- [ ] Basic layout
- [ ] File upload
- [ ] Split editor
- [ ] Simple transpilation
- [ ] Download

### Beta (Week 3)
- [ ] File tree
- [ ] Real-time transpilation
- [ ] Error display
- [ ] Settings panel
- [ ] Dark/light theme

### Full Release (Week 4-5)
- [ ] AST viewer
- [ ] Dependency graph
- [ ] Batch transpilation
- [ ] Project export
- [ ] Polish & testing

---

**Target Launch:** 5 weeks from start  
**Maintenance:** Ongoing with core transpiler updates
