# Week 4: Unified Frontend Architecture

## 🎯 Refined Approach: Shared UI Codebase

**Strategy:** Extract and reuse the existing Vue 3 + Vite + PrimeVue desktop UI for both:
1. **Desktop App** - Tauri wrapper with IPC for local CLI
2. **SaaS Web** - Standalone web app with HTTP API

**Benefits:**
- ✅ Single codebase to maintain
- ✅ Consistent UX across platforms
- ✅ Reuse existing components (MonacoEditor, FileTree, etc.)
- ✅ Leverage Vue 3 + TypeScript + PrimeVue + TailwindCSS
- ✅ Already has advanced features (editor, diagnostics, history)

---

## 📦 New Architecture

```
ts2go/
├── desktop/
│   └── tauri-backend/        # Tauri Rust backend (extracted)
│       ├── src/
│       ├── Cargo.toml
│       └── tauri.conf.json
│
├── packages/
│   └── ui-shared/            # Shared Vue UI (moved from desktop/ui)
│       ├── src/
│       │   ├── components/   # Reusable components
│       │   ├── views/        # Page components
│       │   ├── composables/  # Vue composables
│       │   ├── stores/       # Pinia stores
│       │   ├── services/     # API abstraction layer
│       │   │   ├── platform.service.ts    # Platform detection
│       │   │   ├── transpile.service.ts   # Tauri IPC or HTTP
│       │   │   ├── project.service.ts
│       │   │   └── auth.service.ts
│       │   └── router/
│       ├── package.json
│       └── vite.config.ts
│
├── apps/
│   ├── desktop/              # Desktop app entry
│   │   ├── src/
│   │   │   └── main.ts       # Imports shared UI + Tauri config
│   │   ├── package.json
│   │   └── tauri.conf.json   # Points to tauri-backend
│   │
│   └── web/                  # SaaS web entry
│       ├── src/
│       │   └── main.ts       # Imports shared UI + web config
│       ├── package.json
│       └── vite.config.ts
│
└── saas/
    └── backend/              # Go API (existing)
```

---

## 🔧 Implementation Tasks

### Phase 1: Extract & Restructure (Tasks 1-2) ✅

**Task 1: Extract Tauri Backend**
```bash
# Move Tauri Rust code
mv desktop/tauri desktop/tauri-backend
# Update Cargo.toml, tauri.conf.json paths
```

**Task 2: Create Shared UI Package**
```bash
# Create package structure
mkdir -p packages/ui-shared
mv desktop/ui/* packages/ui-shared/
# Update package.json name to @ts2go/ui-shared
```

### Phase 2: Platform Abstraction (Tasks 3-4)

**Task 3: Platform Detection**
```typescript
// packages/ui-shared/src/composables/usePlatform.ts
export const usePlatform = () => {
  const isTauri = '__TAURI__' in window;
  const isWeb = !isTauri;
  
  return {
    isTauri,
    isWeb,
    platform: isTauri ? 'desktop' : 'web'
  };
};
```

**Task 4: API Abstraction Layer**
```typescript
// packages/ui-shared/src/services/transpile.service.ts
import { invoke } from '@tauri-apps/api/core';
import axios from 'axios';
import { usePlatform } from '@/composables/usePlatform';

export const transpileService = {
  async transpile(files: FileInput[], settings: Settings) {
    const { isTauri } = usePlatform();
    
    if (isTauri) {
      // Desktop: Use Tauri IPC
      return await invoke('transpile', { files, settings });
    } else {
      // Web: Use HTTP API
      const response = await axios.post('/api/v1/transpile', {
        project_id: currentProject.id,
        input_files: files.map(f => f.id),
        settings
      });
      return response.data;
    }
  }
};
```

### Phase 3: Add Web-Specific Features (Tasks 5-6)

**Task 5: Authentication for SaaS**
```vue
<!-- packages/ui-shared/src/views/LoginView.vue -->
<template v-if="isWeb">
  <div class="login-container">
    <form @submit.prevent="handleLogin">
      <input v-model="email" type="email" placeholder="Email" />
      <input v-model="password" type="password" placeholder="Password" />
      <button type="submit">Login</button>
    </form>
  </div>
</template>
```

**Task 6: WebSocket Real-time Updates**
```typescript
// packages/ui-shared/src/services/websocket.service.ts
export const useWebSocket = () => {
  const { isWeb } = usePlatform();
  
  if (!isWeb) return null; // Desktop doesn't need WS
  
  const ws = new WebSocket('ws://localhost:8080/api/v1/ws');
  
  ws.onmessage = (event) => {
    const update = JSON.parse(event.data);
    // Update job status in Pinia store
    jobStore.updateJobStatus(update);
  };
  
  return ws;
};
```

### Phase 4: Build Configuration (Task 7)

**Desktop Build:**
```typescript
// apps/desktop/vite.config.ts
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': '/packages/ui-shared/src'
    }
  },
  clearScreen: false,
  server: {
    strictPort: true,
  },
  envPrefix: ['VITE_', 'TAURI_'],
  build: {
    target: ['es2021', 'chrome100', 'safari13'],
    minify: !process.env.TAURI_DEBUG ? 'esbuild' : false,
    sourcemap: !!process.env.TAURI_DEBUG,
  },
});
```

**Web Build:**
```typescript
// apps/web/vite.config.ts
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': '/packages/ui-shared/src'
    }
  },
  build: {
    outDir: '../../saas/frontend/dist',
    target: 'es2020',
  },
});
```

### Phase 5: Entry Points (Tasks 8-9)

**Desktop Entry:**
```typescript
// apps/desktop/src/main.ts
import { createApp } from 'vue';
import App from '@ts2go/ui-shared/App.vue';
import { router } from '@ts2go/ui-shared/router';
import { createPinia } from 'pinia';

// Desktop-specific setup
import { setupTauri } from './tauri-setup';

const app = createApp(App);
app.use(createPinia());
app.use(router);

setupTauri(); // Initialize Tauri-specific features

app.mount('#app');
```

**Web Entry:**
```typescript
// apps/web/src/main.ts
import { createApp } from 'vue';
import App from '@ts2go/ui-shared/App.vue';
import { router } from '@ts2go/ui-shared/router';
import { createPinia } from 'pinia';

// Web-specific setup
import { setupAuth } from './auth-setup';
import { setupWebSocket } from './websocket-setup';

const app = createApp(App);
app.use(createPinia());
app.use(router);

setupAuth();      // JWT auth interceptors
setupWebSocket(); // Real-time updates

app.mount('#app');
```

---

## 🎨 Component Reusability Matrix

| Component | Desktop | Web | Notes |
|-----------|---------|-----|-------|
| MonacoEditor | ✅ | ✅ | Same editor experience |
| FileTree | ✅ | ✅ | Desktop: local files, Web: uploaded files |
| ProjectView | ✅ | ✅ | Desktop: local projects, Web: cloud projects |
| HistoryView | ✅ | ✅ | Desktop: local history, Web: job history |
| SettingsView | ✅ | ⚠️ | Web version has additional account settings |
| LoginView | ❌ | ✅ | Web-only |
| DiagnosticsPanel | ✅ | ✅ | Same diagnostics display |
| CodeEditor | ✅ | ✅ | Shared editor component |

---

## 🔀 Router Configuration

```typescript
// packages/ui-shared/src/router/index.ts
import { usePlatform } from '@/composables/usePlatform';

const routes = [
  {
    path: '/',
    component: HomeView,
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    component: LoginView,
    meta: { webOnly: true } // Only show on web
  },
  {
    path: '/projects',
    component: ProjectView,
    meta: { requiresAuth: true }
  },
  {
    path: '/editor',
    component: EditorView,
    meta: { requiresAuth: false }
  },
  {
    path: '/history',
    component: HistoryView,
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    component: SettingsView,
    meta: { requiresAuth: false }
  },
];

// Filter routes based on platform
router.beforeEach((to, from, next) => {
  const { isWeb, isTauri } = usePlatform();
  
  // Skip web-only routes on desktop
  if (to.meta.webOnly && isTauri) {
    return next('/');
  }
  
  // Auth check for web
  if (to.meta.requiresAuth && isWeb && !isAuthenticated()) {
    return next('/login');
  }
  
  next();
});
```

---

## 🗂️ Service Layer Pattern

```typescript
// packages/ui-shared/src/services/base.service.ts
export abstract class BaseService {
  protected platform: ReturnType<typeof usePlatform>;
  
  constructor() {
    this.platform = usePlatform();
  }
  
  protected async call<T>(
    tauriCommand: string,
    tauriArgs: any,
    httpMethod: string,
    httpUrl: string,
    httpData?: any
  ): Promise<T> {
    if (this.platform.isTauri) {
      return await invoke<T>(tauriCommand, tauriArgs);
    } else {
      const response = await axios({ method: httpMethod, url: httpUrl, data: httpData });
      return response.data;
    }
  }
}

// packages/ui-shared/src/services/project.service.ts
export class ProjectService extends BaseService {
  async listProjects() {
    return this.call<Project[]>(
      'list_projects',        // Tauri command
      {},                     // Tauri args
      'GET',                  // HTTP method
      '/api/v1/projects',     // HTTP URL
    );
  }
  
  async createProject(name: string, description: string) {
    return this.call<Project>(
      'create_project',
      { name, description },
      'POST',
      '/api/v1/projects',
      { name, description }
    );
  }
}
```

---

## 📋 Detailed Task Breakdown

### ✅ Task 1: Extract Tauri Backend (30 min)
- [x] Move `desktop/tauri` → `desktop/tauri-backend`
- [x] Update Cargo.toml workspace
- [x] Update tauri.conf.json paths
- [x] Test Tauri build

### ✅ Task 2: Create Shared UI Package (45 min)
- [x] Create `packages/ui-shared` structure
- [x] Move `desktop/ui/*` to `packages/ui-shared`
- [x] Update package.json name and dependencies
- [x] Setup monorepo workspace (pnpm/npm workspaces)
- [x] Update import paths

### 🔄 Task 3: Platform Detection Layer (30 min)
- [ ] Create `composables/usePlatform.ts`
- [ ] Add environment detection
- [ ] Export platform utilities

### 🔄 Task 4: API Abstraction Layer (2 hours)
- [ ] Create `services/base.service.ts`
- [ ] Implement `services/transpile.service.ts`
- [ ] Implement `services/project.service.ts`
- [ ] Implement `services/file.service.ts`
- [ ] Update stores to use new services

### 🔄 Task 5: Authentication for SaaS (1.5 hours)
- [ ] Create `views/LoginView.vue`
- [ ] Create `views/RegisterView.vue`
- [ ] Implement `services/auth.service.ts`
- [ ] Add auth store with JWT handling
- [ ] Add auth guards to router

### 🔄 Task 6: WebSocket Support (1 hour)
- [ ] Create `services/websocket.service.ts`
- [ ] Add job status updates
- [ ] Integrate with job store
- [ ] Add reconnection logic

### 🔄 Task 7: Build Configuration (1 hour)
- [ ] Create `apps/desktop` with Vite config
- [ ] Create `apps/web` with Vite config
- [ ] Setup workspace builds
- [ ] Configure environment variables

### 🔄 Task 8: Desktop App Entry (30 min)
- [ ] Create `apps/desktop/src/main.ts`
- [ ] Setup Tauri integration
- [ ] Test desktop build

### 🔄 Task 9: Web App Entry (30 min)
- [ ] Create `apps/web/src/main.ts`
- [ ] Setup HTTP API client
- [ ] Setup WebSocket
- [ ] Test web build

### 🔄 Task 10: Testing & Validation (1 hour)
- [ ] Test desktop app functionality
- [ ] Test web app functionality
- [ ] Verify component reuse
- [ ] Test file upload/download
- [ ] Test real-time updates

---

## 🚀 Implementation Order

**Phase 1: Foundation (Tasks 1-2)**
1. Extract Tauri backend
2. Create shared UI package

**Phase 2: Core Abstraction (Tasks 3-4)**
3. Add platform detection
4. Build API abstraction layer

**Phase 3: Web Features (Tasks 5-6)**
5. Add authentication
6. Add WebSocket support

**Phase 4: Build & Deploy (Tasks 7-9)**
7. Configure builds
8. Setup desktop entry
9. Setup web entry

**Phase 5: Validation (Task 10)**
10. Test both platforms

---

## 📊 Expected Outcomes

**Desktop App:**
- ✅ Uses Tauri IPC for file operations
- ✅ Local file system access
- ✅ Native OS integration
- ✅ No authentication required
- ✅ Fast local transpilation

**Web App:**
- ✅ Uses HTTP REST API
- ✅ Cloud file storage (MinIO)
- ✅ JWT authentication
- ✅ Real-time job updates (WebSocket)
- ✅ Multi-user support

**Shared:**
- ✅ Same UI/UX
- ✅ Same components
- ✅ Same editor experience
- ✅ Consistent behavior

---

## 🎯 Success Criteria

- [ ] Both apps build successfully
- [ ] Desktop app works with local files
- [ ] Web app works with HTTP API
- [ ] Authentication works on web
- [ ] WebSocket updates work on web
- [ ] Code reuse >90%
- [ ] No platform-specific bugs
- [ ] Performance acceptable on both

---

**Total Estimated Time:** 8-10 hours
**Complexity:** Medium-High
**Risk:** Low (incremental migration)
