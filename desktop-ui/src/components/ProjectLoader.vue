<template>
  <div class="project-loader">
    <div class="loader-content">
      <h2>Load TypeScript Project</h2>
      
      <div class="folder-selector">
        <label>Project Folder:</label>
        <div class="input-group">
          <input 
            v-model="selectedPath" 
            type="text" 
            placeholder="Select a folder or enter path..."
            @keyup.enter="loadProject"
          >
          <button @click="browseFolder" class="btn-browse">
            📁 Browse
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading">
        <span class="spinner"></span>
        Loading project...
      </div>

      <div v-if="projectInfo" class="project-info">
        <h3>✅ Project Loaded</h3>
        <div class="info-row">
          <span class="label">Root:</span>
          <span class="value">{{ projectInfo.root }}</span>
        </div>
        <div class="info-row">
          <span class="label">Files Found:</span>
          <span class="value">{{ projectInfo.files.length }} TypeScript files</span>
        </div>
        
        <div class="actions">
          <button @click="autoTranspile" class="btn-primary" :disabled="transpiling">
            {{ transpiling ? '⏳ Transpiling...' : '🚀 Auto-Transpile Project' }}
          </button>
          <button @click="openInEditor" class="btn-secondary">
            📝 Open in Editor
          </button>
        </div>
      </div>

      <div v-if="transpileResult" class="transpile-result" :class="{ success: transpileResult.success, error: !transpileResult.success }">
        <h3>{{ transpileResult.success ? '✅ Transpilation Complete!' : '❌ Transpilation Failed' }}</h3>
        <div v-if="transpileResult.success">
          <p><strong>Output Directory:</strong> {{ transpileResult.output_dir }}</p>
          <p><strong>Files Transpiled:</strong> {{ transpileResult.files_transpiled }}</p>
          <p class="message">{{ transpileResult.message }}</p>
          
          <button @click="openOutputFolder" class="btn-secondary">
            📂 Open Output Folder
          </button>
        </div>
        <div v-else>
          <p class="error-message">{{ transpileResult.error }}</p>
        </div>
      </div>

      <div v-if="error" class="error-box">
        ⚠️ {{ error }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { open } from '@tauri-apps/plugin-dialog'
import { useWorkspaceStore } from '../stores/workspace'
import { useProjectStore } from '../stores/project'
import { useRouter } from 'vue-router'

const selectedPath = ref('')
const loading = ref(false)
const transpiling = ref(false)
const projectInfo = ref<any>(null)
const transpileResult = ref<any>(null)
const error = ref('')

const workspaceStore = useWorkspaceStore()
const projectStore = useProjectStore()
const router = useRouter()

async function browseFolder() {
  try {
    const selected = await open({
      directory: true,
      multiple: false,
      title: 'Select TypeScript Project Folder'
    })
    
    if (selected) {
      selectedPath.value = selected as string
      await loadProject()
    }
  } catch (err) {
    error.value = `Failed to browse folder: ${err}`
  }
}

async function loadProject() {
  if (!selectedPath.value) {
    error.value = 'Please select a folder'
    return
  }

  loading.value = true
  error.value = ''
  projectInfo.value = null
  transpileResult.value = null

  try {
    const result = await invoke<any>('load_project_folder', {
      path: selectedPath.value
    })
    
    projectInfo.value = result
    
    // Add to recent projects
    projectStore.addRecentProject({
      name: selectedPath.value.split('/').pop() || 'Project',
      path: selectedPath.value,
      lastOpened: Date.now()
    })
  } catch (err) {
    error.value = `Failed to load project: ${err}`
  } finally {
    loading.value = false
  }
}

async function autoTranspile() {
  if (!projectInfo.value) return

  transpiling.value = true
  error.value = ''
  transpileResult.value = null

  try {
    const result = await invoke<any>('auto_transpile_project', {
      path: projectInfo.value.root
    })
    
    transpileResult.value = result
  } catch (err) {
    error.value = `Failed to transpile project: ${err}`
  } finally {
    transpiling.value = false
  }
}

function openInEditor() {
  if (!projectInfo.value) return

  // Load files into workspace
  workspaceStore.loadProject(projectInfo.value)
  
  // Navigate to project view
  router.push('/project')
}

function openOutputFolder() {
  if (!transpileResult.value?.output_dir) return
  
  // Open the output folder in system file explorer
  invoke('open_in_explorer', {
    path: transpileResult.value.output_dir
  }).catch(err => {
    error.value = `Failed to open folder: ${err}`
  })
}
</script>

<style scoped>
.project-loader {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.loader-content {
  background: var(--color-background-soft);
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

h2 {
  margin-top: 0;
  color: var(--color-heading);
}

.folder-selector {
  margin: 1.5rem 0;
}

.folder-selector label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.input-group {
  display: flex;
  gap: 0.5rem;
}

.input-group input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 1rem;
}

.btn-browse {
  padding: 0.75rem 1.5rem;
  background: var(--color-background-mute);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
}

.btn-browse:hover {
  background: var(--color-background);
}

.loading {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--color-background-mute);
  border-radius: 4px;
  margin: 1rem 0;
}

.spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.project-info {
  margin: 1.5rem 0;
  padding: 1.5rem;
  background: var(--color-background);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.project-info h3 {
  margin-top: 0;
  color: #10b981;
}

.info-row {
  display: flex;
  gap: 1rem;
  padding: 0.5rem 0;
  border-bottom: 1px solid var(--color-border);
}

.info-row:last-of-type {
  border-bottom: none;
}

.info-row .label {
  font-weight: 600;
  min-width: 120px;
}

.info-row .value {
  color: var(--color-text-light);
}

.actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-primary {
  flex: 1;
  padding: 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  padding: 1rem 1.5rem;
  background: var(--color-background-mute);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
}

.btn-secondary:hover {
  background: var(--color-background);
}

.transpile-result {
  margin: 1.5rem 0;
  padding: 1.5rem;
  border-radius: 8px;
  border: 2px solid;
}

.transpile-result.success {
  background: #d1fae5;
  border-color: #10b981;
}

.transpile-result.error {
  background: #fee2e2;
  border-color: #ef4444;
}

.transpile-result h3 {
  margin-top: 0;
}

.transpile-result.success h3 {
  color: #065f46;
}

.transpile-result.error h3 {
  color: #991b1b;
}

.transpile-result p {
  margin: 0.5rem 0;
}

.message {
  font-weight: 600;
  color: #065f46;
}

.error-message {
  color: #991b1b;
  font-family: monospace;
  white-space: pre-wrap;
}

.error-box {
  padding: 1rem;
  background: #fee2e2;
  border: 1px solid #ef4444;
  border-radius: 4px;
  color: #991b1b;
  margin: 1rem 0;
}
</style>
