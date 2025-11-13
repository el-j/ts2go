<template>
  <AppLayout>
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4 flex-1 min-w-0">
          <div>
            <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Project View</h2>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              {{ projectPath ? projectPath : 'See and edit your TypeScript project files' }}
            </p>
          </div>
          
          <!-- Inline Project Info (visible on larger screens) -->
          <div v-if="projectPath" class="project-info-inline hidden lg:flex">
            <div class="info-badge">
              <span class="icon">📂</span>
              <span class="name">{{ projectName }}</span>
            </div>
            <div class="info-details">
              <span class="detail">{{ fileCount }} files</span>
              <span class="detail-separator">•</span>
              <span class="detail">{{ formatDate(new Date()) }}</span>
            </div>
          </div>
          
          <!-- Info Icon with Popover (visible on smaller screens) -->
          <button 
            v-if="projectPath" 
            @click="toggleProjectInfo"
            class="info-icon-btn lg:hidden"
            type="button"
          >
            <i class="pi pi-info-circle"></i>
          </button>
          <Popover ref="projectInfoPopover">
            <div class="project-info-popover">
              <h4 class="popover-title">{{ projectName }}</h4>
              <div class="popover-details">
                <div class="popover-row">
                  <i class="pi pi-folder"></i>
                  <span>{{ fileCount }} files</span>
                </div>
                <div class="popover-row">
                  <i class="pi pi-map-marker"></i>
                  <span class="text-xs break-all">{{ projectPath }}</span>
                </div>
                <div class="popover-row">
                  <i class="pi pi-clock"></i>
                  <span>{{ formatDate(new Date()) }}</span>
                </div>
              </div>
            </div>
          </Popover>
        </div>
        
        <!-- Action Buttons -->
        <div v-if="projectPath" class="action-buttons">
          <button 
            @click="handleTranspile" 
            class="btn-action btn-primary"
            :disabled="transpileStore.isTranspiling"
          >
            {{ transpileStore.isTranspiling ? '⏳ Transpiling...' : '🚀 Transpile All' }}
          </button>
          <button 
            @click="handleTranspileFile" 
            class="btn-action btn-secondary"
            :disabled="!activeFile || transpileStore.isTranspiling"
          >
            📄 Transpile File
          </button>
          <button 
            v-if="transpileStore.currentResult?.output_dir"
            @click="openOutputFolder" 
            class="btn-action btn-secondary"
          >
            👁️ Preview
          </button>
        </div>
      </div>
    </div>
    
    <div class="project-view">
      <!-- Show ProjectLoader only if no project is loaded -->
      <div v-if="!projectPath" class="loader-wrapper">
        <ProjectLoader />
      </div>
      
      <!-- Show project workspace when project is loaded -->
      <div v-else class="workspace-container">
        <Splitter class="workspace-splitter">
          <!-- File Tree Sidebar -->
          <SplitterPanel :size="15" :min-size="10" class="file-tree-panel">
            <FileTree />
          </SplitterPanel>
          
          <!-- Editor Column with Tabs -->
          <SplitterPanel :size="55" :min-size="30" class="editor-panel">
            <div class="editor-column">
              <!-- File Tabs -->
              <FileTabs />
              
              <!-- Editor Area -->
              <div class="editor-area">
                <div v-if="activeFile" class="editor-container">
                  <MonacoEditor
                    v-model="activeFile.content"
                    :language="getEditorLanguage(activeFile.name)"
                    @update:modelValue="handleMonacoChange"
                  />
                </div>
                <div v-else class="no-file-open">
                  <div class="empty-state">
                    <span class="empty-icon">📂</span>
                    <h3>No File Open</h3>
                    <p>Select a file from the tree or create a new one</p>
                  </div>
                </div>
              </div>
            </div>
          </SplitterPanel>
          
          <!-- Output Panel -->
          <SplitterPanel :size="30" :min-size="20" class="output-panel-wrapper">
            <div class="output-panel">
              <div class="panel-header">
                <div class="header-left">
                  <h4>Output</h4>
                  <span v-if="transpileStore.isTranspiling" class="elapsed-badge">
                    {{ transpileStore.elapsedTime }}s
                  </span>
                </div>
                <button 
                  v-if="transpileStore.currentResult" 
                  @click="transpileStore.clearResult()"
                  class="clear-btn"
                >
                  ×
                </button>
              </div>
              <div class="panel-content">
                <!-- Loading State -->
                <div v-if="transpileStore.isTranspiling" class="output-loading">
                  <div class="spinner"></div>
                  <p>Transpiling project...</p>
                  <p class="elapsed">{{ transpileStore.elapsedTime }}s elapsed</p>
                </div>
                
                <!-- Success State with Go Code (Single File) -->
                <div v-else-if="transpileStore.currentResult?.goCode" class="go-code-output">
                  <MonacoEditor
                    :model-value="transpileStore.currentResult.goCode"
                    language="go"
                    :readonly="true"
                    theme="vs-dark"
                  />
                </div>
                
                <!-- Error State -->
                <div v-else-if="transpileStore.hasError" class="output-error">
                  <div class="error-header">
                    <span class="icon">❌</span>
                    <h3>Transpilation Failed</h3>
                  </div>
                  <pre class="error-message">{{ transpileStore.currentResult?.error }}</pre>
                </div>
                
                <!-- Empty State -->
                <p v-else class="placeholder">Transpilation output will appear here</p>
              </div>
            </div>
          </SplitterPanel>
        </Splitter>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { useToast } from 'primevue/usetoast'
import ProjectLoader from '../components/ProjectLoader.vue'
import AppLayout from '../components/AppLayout.vue'
import FileTree from '../components/FileTree.vue'
import FileTabs from '../components/FileTabs.vue'
import MonacoEditor from '../components/MonacoEditor.vue'

import { useWorkspaceStore } from '../stores/workspace'
import { useTranspileStore } from '../stores/transpile'

const workspace = useWorkspaceStore()
const transpileStore = useTranspileStore()
const toast = useToast()

const projectInfoPopover = ref()

function toggleProjectInfo(event: Event) {
  projectInfoPopover.value.toggle(event)
}

const activeFile = computed(() => workspace.activeFile)
const projectPath = computed(() => workspace.projectPath)

const projectName = computed(() => {
  const path = workspace.projectPath
  return path.split('/').filter(Boolean).pop() || 'Project'
})

const fileCount = computed(() => {
  const countFiles = (nodes: any[]): number => {
    return nodes.reduce((count, node) => {
      if (node.type === 'file') {
        return count + 1
      }
      if (node.children) {
        return count + countFiles(node.children)
      }
      return count
    }, 0)
  }
  return countFiles(workspace.fileTree)
})

function formatDate(date: Date): string {
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getEditorLanguage(filename: string): string {
  const ext = filename.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'ts':
    case 'tsx':
      return 'typescript'
    case 'js':
    case 'jsx':
      return 'javascript'
    case 'json':
      return 'json'
    case 'md':
      return 'markdown'
    case 'go':
      return 'go'
    case 'vue':
      return 'vue'
    default:
      return 'plaintext'
  }
}

function handleMonacoChange(content: string) {
  if (activeFile.value) {
    workspace.updateFileContent(activeFile.value.path, content)
  }
}

async function handleTranspile() {
  if (!projectPath.value) return
  
  try {
    await transpileStore.transpileProject(projectPath.value)
    
    // Show success toast
    toast.add({
      severity: 'success',
      summary: 'Transpilation Complete!',
      detail: `Transpiled ${transpileStore.currentResult?.files_transpiled} files successfully`,
      life: 4000
    })
  } catch (error) {
    console.error('Transpilation failed:', error)
    
    // Show error toast
    toast.add({
      severity: 'error',
      summary: 'Transpilation Failed',
      detail: String(error),
      life: 5000
    })
  }
}

async function handleTranspileFile() {
  if (!activeFile.value || !projectPath.value) return
  
  try {
    transpileStore.isTranspiling = true
    transpileStore.currentResult = null
    
    // Transpile the current file
    const goCode = await invoke<string>('transpile_code', {
      code: activeFile.value.content,
      filename: activeFile.value.name
    })
    
    // Store result in transpile store for display
    transpileStore.currentResult = {
      success: true,
      message: `Successfully transpiled ${activeFile.value.name}`,
      files_transpiled: 1,
      goCode: goCode // Add the Go code to result
    }
    
    // Show success toast
    toast.add({
      severity: 'success',
      summary: 'File Transpiled!',
      detail: `Successfully transpiled ${activeFile.value.name}`,
      life: 3000
    })
  } catch (error) {
    console.error('File transpilation failed:', error)
    transpileStore.currentResult = {
      success: false,
      error: `Failed to transpile file: ${error}`
    }
    
    // Show error toast
    toast.add({
      severity: 'error',
      summary: 'Transpilation Failed',
      detail: String(error),
      life: 5000
    })
  } finally {
    transpileStore.isTranspiling = false
  }
}

async function openOutputFolder() {
  if (!transpileStore.currentResult?.output_dir) return
  
  try {
    await invoke('open_in_explorer', {
      path: transpileStore.currentResult.output_dir
    })
  } catch (error) {
    console.error('Failed to open output folder:', error)
  }
}
</script>

<style scoped>
.project-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.loader-wrapper {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

.workspace-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.project-info-inline {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 16px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 8px;
  font-size: 13px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.info-badge {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-badge .icon {
  font-size: 18px;
}

.info-badge .name {
  font-weight: 700;
  color: #667eea;
  font-size: 14px;
}

.info-details {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #667eea;
  font-size: 12px;
  font-weight: 600;
}

.detail {
  font-family: 'Monaco', 'Menlo', monospace;
}

.detail-separator {
  color: rgba(102, 126, 234, 0.4);
  font-weight: 400;
}

.info-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(102, 126, 234, 0.1);
  border: 1px solid rgba(102, 126, 234, 0.2);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s;
}

.info-icon-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: scale(1.05);
}

.info-icon-btn i {
  font-size: 18px;
}

.project-info-popover {
  padding: 16px;
  min-width: 280px;
}

.popover-title {
  margin: 0 0 12px 0;
  font-size: 16px;
  font-weight: 700;
  color: #667eea;
}

.popover-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.popover-row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--color-text);
}

.popover-row i {
  color: #667eea;
  width: 16px;
  text-align: center;
}

.action-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
}

.btn-action {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  white-space: nowrap;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.btn-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
}

.btn-secondary:hover:not(:disabled) {
  background: #667eea;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(102, 126, 234, 0.3);
}

.workspace-splitter {
  flex: 1 !important;
  height: 100% !important;
}

.workspace-splitter :deep(.p-splitter-panel) {
  overflow: hidden;
}

.file-tree-panel {
  overflow-y: auto;
}

.editor-panel {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.output-panel-wrapper {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.file-tree-sidebar {
  border-right: 1px solid var(--color-border);
  overflow-y: auto;
}

.project-content {
  flex: 1;
  display: grid;
  grid-template-columns: 250px 1fr 400px;
  grid-template-rows: 1fr;
  overflow: hidden;
}

.file-tree-sidebar {
  border-right: 1px solid var(--color-border);
  overflow-y: auto;
}

.editor-column {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100%;
}

.editor-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.editor-container :deep(.monaco-editor-container) {
  flex: 1;
  height: 100%;
  min-height: 0;
}

.editor-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-soft);
}

.editor-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.unsaved-indicator {
  color: var(--color-warning);
  font-size: 12px;
}

.code-editor {
  flex: 1;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
  border: none;
  outline: none;
  resize: none;
  background: var(--color-background);
  color: var(--color-text);
}

.no-file-open {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 32px;
}

.empty-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

.empty-state h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
}

.empty-state p {
  color: var(--color-text-secondary);
  margin: 0;
}

.output-panel {
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
}

.panel-header {
  padding: 8px 12px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-soft);
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 40px;
  height: 40px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.panel-header h4 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--color-text-secondary);
}

.elapsed-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
  font-family: monospace;
}

.clear-btn {
  background: transparent;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--color-text-secondary);
  padding: 0 4px;
  line-height: 1;
}

.clear-btn:hover {
  color: var(--color-text);
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.panel-content > .placeholder,
.panel-content > .output-loading,
.panel-content > .output-error {
  padding: 16px;
}

.placeholder {
  color: var(--color-text-secondary);
  font-size: 13px;
}

.output-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 16px 0;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(102, 126, 234, 0.2);
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.output-loading p {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}

.elapsed {
  color: var(--color-text-secondary);
  font-size: 12px !important;
}

.output-error {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.success-header,
.error-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.success-header .icon {
  font-size: 24px;
}

.error-header .icon {
  font-size: 24px;
}

.success-header h3,
.error-header h3 {
  margin: 0;
  font-size: 16px;
}

.success-header h3 {
  color: #10b981;
}

.error-header h3 {
  color: #ef4444;
}

.go-code-output {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
  height: 100%;
}

.go-code-output :deep(.monaco-editor-container) {
  flex: 1;
  height: 100%;
  min-height: 0;
  width: 100%;
}

.code-header {
  padding: 8px 12px;
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid rgba(139, 92, 246, 0.2);
  border-radius: 4px;
  margin-bottom: 8px;
}

.code-header .label {
  color: #8b5cf6;
  font-weight: 600;
  font-size: 13px;
}

.go-code-output .monaco-editor-wrapper {
  flex: 1;
  min-height: 0;
}

.result-details {
  background: #d1fae5;
  padding: 12px;
  border-radius: 6px;
  font-size: 13px;
}

.detail-row {
  display: flex;
  gap: 8px;
  padding: 4px 0;
}

.detail-row .label {
  font-weight: 600;
  min-width: 130px;
}

.detail-row .value {
  color: #065f46;
  word-break: break-all;
}

.message {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #10b981;
  font-weight: 600;
  color: #065f46;
}

.error-message {
  background: #fee2e2;
  padding: 12px;
  border-radius: 6px;
  color: #991b1b;
  font-family: monospace;
  font-size: 12px;
  white-space: pre-wrap;
  overflow-x: auto;
}

/* Responsive */
@media (max-width: 1200px) {
  .project-content {
    grid-template-columns: 250px 1fr 350px;
  }
  
  .project-info-inline {
    display: none;
  }
}

@media (max-width: 1024px) {
  .project-content {
    grid-template-columns: 200px 1fr;
  }
  
  .output-panel {
    display: none;
  }
  
  .action-buttons .btn-secondary {
    display: none;
  }
}

@media (max-width: 768px) {
  .project-content {
    grid-template-columns: 1fr;
  }
  
  .file-tree-sidebar {
    display: none;
  }
}
</style>
