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
            v-tooltip.bottom="'Transpile all TypeScript files to Go (Ctrl+Shift+T)'"
          >
            {{ transpileStore.isTranspiling ? '⏳ Transpiling...' : '🚀 Transpile All' }}
          </button>
          <button 
            @click="handleTranspileFile" 
            class="btn-action btn-secondary"
            :disabled="!activeFile || transpileStore.isTranspiling"
            v-tooltip.bottom="'Transpile current file only (Ctrl+T)'"
          >
            📄 Transpile File
          </button>
          <button 
            v-if="transpileStore.currentResult?.output_dir"
            @click="openOutputFolder" 
            class="btn-action btn-secondary"
            v-tooltip.bottom="'Open output folder in file explorer'"
          >
            👁️ Preview
          </button>
          <button 
            v-if="transpileStore.currentResult?.output_dir && !isRunningProject"
            @click="handleRunProject" 
            class="btn-action btn-success"
            v-tooltip.bottom="'Run the complete Go project'"
          >
            ▶️ Run Project
          </button>
          <span v-if="isRunningProject" class="running-project-indicator">⏳ Running Project...</span>
          <button 
            v-if="transpileStore.currentResult?.output_dir && !buildStore.isBuilding"
            @click="handleBuildProject" 
            class="btn-action btn-build"
            v-tooltip.bottom="'Build project into executable binary'"
          >
            🔨 Build Project
          </button>
          <span v-if="buildStore.isBuilding" class="building-indicator">⏳ Building...</span>
          <button 
            v-if="transpileStore.currentResult?.output_dir && !testStore.isTesting"
            @click="handleTestProject" 
            class="btn-action btn-test"
            v-tooltip.bottom="'Run tests on the Go project'"
          >
            🧪 Test Project
          </button>
          <span v-if="testStore.isTesting" class="testing-indicator">⏳ Testing...</span>
          
          <!-- Help Button -->
          <button 
            @click="showShortcutsDialog = true"
            class="btn-action btn-help"
            v-tooltip.bottom="'View keyboard shortcuts (?)'"
          >
            <i class="pi pi-question-circle"></i>
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
                  <span v-if="transpileStore.isTranspiling && workspace.transpilationProgress.currentFile" class="progress-badge">
                    {{ getProgressText() }}
                  </span>
                  <button 
                    v-if="transpileStore.currentResult?.goCode && !isRunning"
                    @click="handleRunGoCode"
                    class="run-btn"
                    title="Run & Test Go Code"
                  >
                    ▶️ Run
                  </button>
                  <span v-if="isRunning" class="running-indicator">⏳ Running...</span>
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
                <!-- State Restoration Banner -->
                <div 
                  v-if="hasValidTranspilationState && !transpileStore.isTranspiling && !transpileStore.currentResult?.goCode" 
                  class="state-restoration-banner"
                >
                  <div class="banner-content">
                    <div class="banner-icon">📦</div>
                    <div class="banner-info">
                      <h4>Previous Transpilation Restored</h4>
                      <p class="banner-details">
                        {{ lastTranspilationFiles }} files transpiled on {{ lastTranspilationTime }}
                      </p>
                      <p class="banner-hint">
                        Build, Test, and Run commands are available. Transpile again to update.
                      </p>
                    </div>
                    <button 
                      @click="clearTranspilationState"
                      class="banner-clear-btn"
                      title="Clear saved state"
                    >
                      Clear State
                    </button>
                  </div>
                </div>

                <!-- Restoring State -->
                <div v-if="restoringState" class="output-loading">
                  <div class="spinner"></div>
                  <p>Restoring previous transpilation state...</p>
                </div>

                <!-- Loading State -->
                <div v-else-if="transpileStore.isTranspiling" class="output-loading">
                  <div class="spinner"></div>
                  <p>Transpiling project...</p>
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
                
                <!-- Log Output (when no file selected and we have logs) -->
                <div v-else-if="workspace.transpilationProgress.logs.length > 0" class="log-output">
                  <div class="log-header">
                    <h4>Transpilation Log</h4>
                    <button @click="workspace.transpilationProgress.logs = []" class="clear-log-btn" title="Clear logs">
                      Clear
                    </button>
                  </div>
                  <div class="log-entries">
                    <div 
                      v-for="(log, index) in workspace.transpilationProgress.logs" 
                      :key="index"
                      :class="['log-entry', `log-${log.level}`]"
                    >
                      <span class="log-icon">{{ getLogIcon(log.level) }}</span>
                      <span class="log-time">{{ formatLogTime(log.timestamp) }}</span>
                      <span class="log-message">{{ log.message }}</span>
                    </div>
                  </div>
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
    
    <!-- Unsaved Changes Dialog -->
    <UnsavedChangesDialog 
      v-model="showUnsavedDialog"
      @save="handleSaveAll"
      @discard="handleDiscardChanges"
      @cancel="handleCancelClose"
    />
    
    <!-- Keyboard Shortcuts Dialog -->
    <KeyboardShortcuts v-model="showShortcutsDialog" />
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { useToast } from 'primevue/usetoast'
import ProjectLoader from '../components/ProjectLoader.vue'
import AppLayout from '../components/AppLayout.vue'
import FileTree from '../components/FileTree.vue'
import FileTabs from '../components/FileTabs.vue'
import MonacoEditor from '../components/MonacoEditor.vue'
import Popover from 'primevue/popover'
import Splitter from 'primevue/splitter'
import SplitterPanel from 'primevue/splitterpanel'
import UnsavedChangesDialog from '@/components/UnsavedChangesDialog.vue'
import KeyboardShortcuts from '@/components/KeyboardShortcuts.vue'

import { useWorkspaceStore } from '../stores/workspace'
import { useTranspileStore } from '../stores/transpile'
import { useBuildStore } from '../stores/build'
import { useTestStore } from '../stores/test'
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts'
import { useAutoSave } from '@/composables/useAutoSave'
import { useSaveFile } from '@/composables/useSaveFile'

const workspace = useWorkspaceStore()
const transpileStore = useTranspileStore()
const buildStore = useBuildStore()
const testStore = useTestStore()
const toast = useToast()

// Initialize keyboard shortcuts and auto-save
useKeyboardShortcuts()
const { scheduleSave } = useAutoSave()
const { saveAll } = useSaveFile()

// Dialog states
const showUnsavedDialog = ref(false)
const showShortcutsDialog = ref(false)

const projectInfoPopover = ref()
const isRunning = ref(false)
const isRunningProject = ref(false)

// State restoration
const hasValidTranspilationState = ref(false)
const lastTranspilationTime = ref<string | null>(null)
const lastTranspilationFiles = ref(0)
const restoringState = ref(false)

// Handler functions for dialogs
async function handleSaveAll() {
  await saveAll()
  showUnsavedDialog.value = false
}

function handleDiscardChanges() {
  showUnsavedDialog.value = false
  // Continue with close operation
}

function handleCancelClose() {
  showUnsavedDialog.value = false
  // Cancel close operation
}

// Before unload handler to prevent data loss
function handleBeforeUnload(e: BeforeUnloadEvent) {
  if (workspace.hasUnsavedChanges) {
    e.preventDefault()
    e.returnValue = ''
    showUnsavedDialog.value = true
  }
}

// Add beforeunload listener on mount
if (typeof window !== 'undefined') {
  window.addEventListener('beforeunload', handleBeforeUnload)
}

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

function toggleProjectInfo(event: Event) {
  projectInfoPopover.value.toggle(event)
}

function getProgressText(): string {
  const progress = workspace.transpilationProgress
  const fileName = progress.currentFile?.split('/').pop() || 'file'
  return `Transpiling: ${fileName} (${progress.completedFiles}/${progress.totalFiles})`
}

function getLogIcon(level: string): string {
  switch (level) {
    case 'success': return '✓'
    case 'error': return '✗'
    case 'warning': return '⚠'
    default: return 'ℹ'
  }
}

function formatLogTime(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', { 
    hour: '2-digit', 
    minute: '2-digit', 
    second: '2-digit',
    hour12: false
  })
}

// Watch for active file changes and auto-load transpiled output
watch(() => workspace.activeFilePath, async (newPath) => {
  if (!newPath) return
  
  // Check if this file has been transpiled
  const transpilation = workspace.getTranspilation(newPath)
  if (transpilation && transpilation.success && transpilation.goCode) {
    // Auto-load the Go code into output panel
    transpileStore.currentResult = {
      success: true,
      message: `Showing transpiled output for ${activeFile.value?.name}`,
      files_transpiled: 1,
      goCode: transpilation.goCode
    }
  }
})

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
    
    // Schedule auto-save after content change
    scheduleSave(activeFile.value)
  }
}

async function handleTranspile() {
  if (!projectPath.value) return
  
  try {
    // Count TypeScript files in the project
    const tsFiles = countTypeScriptFiles(workspace.fileTree)
    
    // Initialize transpilation progress
    workspace.startTranspilation(tsFiles)
    
    // Mark all TS files as pending
    markAllTsFilesStatus('pending')
    
    // Start actual transpilation
    await transpileStore.transpileProject(projectPath.value)
    
    // After successful transpilation, load all Go files into transpilation map
    if (transpileStore.currentResult?.success && transpileStore.currentResult?.output_dir) {
      await loadTranspiledGoFiles(transpileStore.currentResult.output_dir)
    }
    
    // Mark all files as success (we don't have per-file info from backend yet)
    // In a future update, the backend should provide per-file results
    markAllTsFilesStatus('success')
    workspace.finishTranspilation()
    
    // Show success toast with format info
    let detail = `Transpiled ${transpileStore.currentResult?.files_transpiled} files successfully`
    if (transpileStore.currentResult?.files_formatted !== undefined) {
      detail += ` (${transpileStore.currentResult.files_formatted} formatted)`
    }
    toast.add({
      severity: 'success',
      summary: 'Transpilation Complete!',
      detail: detail,
      life: 4000
    })
    
    // Show format warnings if any
    if (transpileStore.currentResult?.format_warnings && 
        transpileStore.currentResult.format_warnings.length > 0) {
      toast.add({
        severity: 'warn',
        summary: 'Formatting Warnings',
        detail: transpileStore.currentResult.format_warnings.join(', '),
        life: 5000
      })
    }
  } catch (error) {
    console.error('Transpilation failed:', error)
    
    // Mark files as error
    markAllTsFilesStatus('error', String(error))
    workspace.addTranspilationLog(`Transpilation failed: ${error}`, 'error')
    
    // Show error toast
    toast.add({
      severity: 'error',
      summary: 'Transpilation Failed',
      detail: String(error),
      life: 5000
    })
  }
}

async function loadTranspiledGoFiles(outputDir: string) {
  try {
    // Get list of all Go files in output directory
    const goFiles = await invoke<string[]>('get_go_files', { path: outputDir })
    
    workspace.addTranspilationLog(`Loading ${goFiles.length} transpiled Go files...`, 'info')
    
    // For each Go file, find corresponding TS file and populate map
    for (const goFilePath of goFiles) {
      try {
        // Read the Go file content
        const goCode = await invoke<string>('read_file', { path: goFilePath })
        
        // Determine corresponding TypeScript file path
        // Extract relative path from output dir
        const relativePath = goFilePath.replace(outputDir, '').replace(/^\//, '')
        // Replace .go with .ts/.tsx and construct original path
        const tsFileName = relativePath.replace(/\.go$/, '.ts')
        const tsxFileName = relativePath.replace(/\.go$/, '.tsx')
        
        // Find the actual TS file in the file tree
        const tsFilePath = findTsFileInTree(tsFileName, tsxFileName)
        
        if (tsFilePath) {
          // Store in transpilation map
          workspace.setTranspilation(tsFilePath, goFilePath, goCode, true)
          workspace.addTranspilationLog(`✓ Loaded ${relativePath}`, 'success')
        }
      } catch (error) {
        console.warn(`Failed to load Go file ${goFilePath}:`, error)
      }
    }
    
    workspace.addTranspilationLog(`Successfully loaded all transpiled files`, 'success')
  } catch (error) {
    console.error('Failed to load transpiled Go files:', error)
    workspace.addTranspilationLog(`Warning: Could not load Go files: ${error}`, 'warning')
  }
}

function findTsFileInTree(tsFileName: string, tsxFileName: string): string | null {
  const searchInNodes = (nodes: any[], targetName: string): string | null => {
    for (const node of nodes) {
      if (node.type === 'file' && node.name === targetName) {
        return node.path
      }
      if (node.children) {
        const found = searchInNodes(node.children, targetName)
        if (found) return found
      }
    }
    return null
  }
  
  // Try both .ts and .tsx extensions
  const tsName = tsFileName.split('/').pop() || ''
  const tsxName = tsxFileName.split('/').pop() || ''
  
  return searchInNodes(workspace.fileTree, tsName) || searchInNodes(workspace.fileTree, tsxName)
}

function countTypeScriptFiles(nodes: any[]): number {
  let count = 0
  for (const node of nodes) {
    if (node.type === 'file' && /\.tsx?$/.test(node.name)) {
      count++
    }
    if (node.children) {
      count += countTypeScriptFiles(node.children)
    }
  }
  return count
}

function markAllTsFilesStatus(status: any, errorMessage?: string) {
  const markNodes = (nodes: any[]) => {
    for (const node of nodes) {
      if (node.type === 'file' && /\.tsx?$/.test(node.name)) {
        workspace.updateTranspilationStatus(node.path, status, errorMessage)
      }
      if (node.children) {
        markNodes(node.children)
      }
    }
  }
  markNodes(workspace.fileTree)
}

async function handleTranspileFile() {
  if (!activeFile.value || !projectPath.value) return
  
  try {
    transpileStore.isTranspiling = true
    transpileStore.currentResult = null
    
    // Start single file transpilation
    workspace.startTranspilation(1)
    workspace.setCurrentTranspilingFile(activeFile.value.path)
    workspace.updateTranspilationStatus(activeFile.value.path, 'transpiling')
    
    // Transpile the current file
    const goCode = await invoke<string>('transpile_code', {
      code: activeFile.value.content,
      filename: activeFile.value.name
    })
    
    // Determine Go file path
    const goFilePath = activeFile.value.path.replace(/\.tsx?$/, '.go')
    
    // Store transpilation info in workspace
    workspace.setTranspilation(activeFile.value.path, goFilePath, goCode, true)
    workspace.completeFileTranspilation(activeFile.value.path, true)
    workspace.finishTranspilation()
    
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
    
    // Mark as error in workspace
    workspace.updateTranspilationStatus(activeFile.value.path, 'error', String(error))
    workspace.completeFileTranspilation(activeFile.value.path, false, String(error))
    workspace.finishTranspilation()
    
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

async function handleRunGoCode() {
  if (!transpileStore.currentResult?.goCode) return
  
  isRunning.value = true
  
  try {
    // Execute Go code via Tauri backend
    const result = await invoke<{
      success: boolean
      stdout: string
      stderr: string
      exit_code: number
      duration_ms: number
    }>('run_go_code', {
      code: transpileStore.currentResult.goCode
    })
    
    // Display results in output panel by switching to log view
    workspace.transpilationProgress.logs = []
    workspace.addTranspilationLog('=== Go Code Execution ===', 'info')
    workspace.addTranspilationLog(`Exit code: ${result.exit_code}`, result.success ? 'success' : 'error')
    workspace.addTranspilationLog(`Duration: ${result.duration_ms}ms`, 'info')
    
    if (result.stdout) {
      workspace.addTranspilationLog('--- Standard Output ---', 'info')
      result.stdout.split('\n').forEach(line => {
        if (line.trim()) {
          workspace.addTranspilationLog(line, 'success')
        }
      })
    }
    
    if (result.stderr) {
      workspace.addTranspilationLog('--- Standard Error ---', 'warning')
      result.stderr.split('\n').forEach(line => {
        if (line.trim()) {
          workspace.addTranspilationLog(line, 'error')
        }
      })
    }
    
    // Clear current result to show logs
    transpileStore.currentResult = null
    
    // Show toast notification
    toast.add({
      severity: result.success ? 'success' : 'error',
      summary: result.success ? 'Execution Complete' : 'Execution Failed',
      detail: `Finished in ${result.duration_ms}ms with exit code ${result.exit_code}`,
      life: 4000
    })
  } catch (error) {
    workspace.addTranspilationLog(`Execution error: ${error}`, 'error')
    transpileStore.currentResult = null
    
    // Check if error is Go-related
    const errorStr = String(error).toLowerCase()
    const isGoError = errorStr.includes('go not found') || 
                      errorStr.includes('go binary') || 
                      errorStr.includes('install go') ||
                      errorStr.includes('configure in settings')
    
    if (isGoError) {
      toast.add({
        severity: 'error',
        summary: 'Go Compiler Not Found',
        detail: 'Configure Go in Settings to use Build, Test, and Run features',
        life: 8000,
        group: 'go-error',
        closable: true
      })
    } else {
      toast.add({
        severity: 'error',
        summary: 'Execution Failed',
        detail: String(error),
        life: 5000
      })
    }
  } finally {
    isRunning.value = false
  }
}

async function handleRunProject() {
  if (!transpileStore.currentResult?.output_dir) return
  
  isRunningProject.value = true
  
  try {
    // Execute complete Go project via Tauri backend
    const result = await invoke<{
      success: boolean
      stdout: string
      stderr: string
      exit_code: number
      duration_ms: number
    }>('run_go_project', {
      outputDir: transpileStore.currentResult.output_dir
    })
    
    // Display results in output panel
    workspace.transpilationProgress.logs = []
    workspace.addTranspilationLog('=== Go Project Execution ===', 'info')
    workspace.addTranspilationLog(`Project: ${transpileStore.currentResult.output_dir}`, 'info')
    workspace.addTranspilationLog(`Exit code: ${result.exit_code}`, result.success ? 'success' : 'error')
    workspace.addTranspilationLog(`Duration: ${result.duration_ms}ms`, 'info')
    
    if (result.stdout) {
      workspace.addTranspilationLog('--- Standard Output ---', 'info')
      result.stdout.split('\n').forEach(line => {
        if (line.trim()) {
          workspace.addTranspilationLog(line, 'success')
        }
      })
    }
    
    if (result.stderr) {
      workspace.addTranspilationLog('--- Standard Error ---', 'warning')
      result.stderr.split('\n').forEach(line => {
        if (line.trim()) {
          workspace.addTranspilationLog(line, 'error')
        }
      })
    }
    
    // Clear current result to show logs
    transpileStore.currentResult = null
    
    // Show toast notification
    toast.add({
      severity: result.success ? 'success' : 'error',
      summary: result.success ? 'Project Execution Complete' : 'Project Execution Failed',
      detail: `Finished in ${result.duration_ms}ms with exit code ${result.exit_code}`,
      life: 4000
    })
  } catch (error) {
    workspace.addTranspilationLog(`Project execution error: ${error}`, 'error')
    transpileStore.currentResult = null
    
    // Check if error is Go-related
    const errorStr = String(error).toLowerCase()
    const isGoError = errorStr.includes('go not found') || 
                      errorStr.includes('go binary') || 
                      errorStr.includes('install go') ||
                      errorStr.includes('configure in settings')
    
    if (isGoError) {
      toast.add({
        severity: 'error',
        summary: 'Go Compiler Not Found',
        detail: 'Configure Go in Settings to use Build, Test, and Run features',
        life: 8000,
        group: 'go-error',
        closable: true
      })
    } else {
      toast.add({
        severity: 'error',
        summary: 'Project Execution Failed',
        detail: String(error),
        life: 5000
      })
    }
  } finally {
    isRunningProject.value = false
  }
}

async function handleBuildProject() {
  if (!transpileStore.currentResult?.output_dir) return
  
  try {
    const outputDir = transpileStore.currentResult.output_dir
    const projectName = workspace.projectPath.split('/').pop() || 'app'
    const binaryPath = `${outputDir}/${projectName}`
    
    // Build the project
    const result = await buildStore.buildProject(outputDir, binaryPath)
    
    // Display results in output panel
    workspace.transpilationProgress.logs = []
    workspace.addTranspilationLog('=== Go Project Build ===', 'info')
    workspace.addTranspilationLog(`Source: ${outputDir}`, 'info')
    workspace.addTranspilationLog(`Binary: ${result.binary_path}`, 'info')
    workspace.addTranspilationLog(`Size: ${(result.binary_size / (1024 * 1024)).toFixed(2)} MB`, 'info')
    workspace.addTranspilationLog(`Duration: ${result.duration_ms}ms`, 'info')
    workspace.addTranspilationLog(`Exit code: ${result.exit_code}`, result.success ? 'success' : 'error')
    
    if (result.errors.length > 0) {
      workspace.addTranspilationLog('--- Build Errors ---', 'error')
      result.errors.forEach(err => {
        workspace.addTranspilationLog(err, 'error')
      })
    }
    
    if (result.warnings.length > 0) {
      workspace.addTranspilationLog('--- Build Warnings ---', 'warning')
      result.warnings.forEach(warn => {
        workspace.addTranspilationLog(warn, 'warning')
      })
    }
    
    if (result.success) {
      workspace.addTranspilationLog(`✅ Build successful! Binary created at: ${result.binary_path}`, 'success')
    }
    
    // Clear current result to show logs
    transpileStore.currentResult = null
    
    // Show toast notification
    toast.add({
      severity: result.success ? 'success' : 'error',
      summary: result.success ? 'Build Complete' : 'Build Failed',
      detail: result.success 
        ? `Binary created in ${result.duration_ms}ms (${(result.binary_size / (1024 * 1024)).toFixed(2)} MB)`
        : `Build failed with ${result.errors.length} error(s)`,
      life: result.success ? 4000 : 8000
    })
  } catch (error) {
    workspace.addTranspilationLog(`Build error: ${error}`, 'error')
    transpileStore.currentResult = null
    
    // Check if error is Go-related
    const errorStr = String(error).toLowerCase()
    const isGoError = errorStr.includes('go not found') || 
                      errorStr.includes('go binary') || 
                      errorStr.includes('install go') ||
                      errorStr.includes('configure in settings')
    
    if (isGoError) {
      toast.add({
        severity: 'error',
        summary: 'Go Compiler Not Found',
        detail: 'Configure Go in Settings to use Build, Test, and Run features',
        life: 8000,
        group: 'go-error',
        closable: true
      })
    } else {
      toast.add({
        severity: 'error',
        summary: 'Build Failed',
        detail: String(error),
        life: 5000
      })
    }
  }
}

async function handleTestProject() {
  if (!transpileStore.currentResult?.output_dir) return
  
  try {
    const outputDir = transpileStore.currentResult.output_dir
    
    // Run tests
    const result = await testStore.testProject(outputDir)
    
    // Display results in output panel
    workspace.transpilationProgress.logs = []
    workspace.addTranspilationLog('=== Go Project Tests ===', 'info')
    workspace.addTranspilationLog(`Source: ${outputDir}`, 'info')
    workspace.addTranspilationLog(`Duration: ${result.duration_ms}ms`, 'info')
    workspace.addTranspilationLog(`Exit code: ${result.exit_code}`, result.success ? 'success' : 'error')
    workspace.addTranspilationLog('', 'info')
    
    // Display test summary
    const testResults = result.test_results
    workspace.addTranspilationLog(`📊 Test Summary:`, 'info')
    workspace.addTranspilationLog(`   Total: ${testResults.total}`, 'info')
    workspace.addTranspilationLog(`   ✅ Passed: ${testResults.passed}`, 'success')
    workspace.addTranspilationLog(`   ❌ Failed: ${testResults.failed}`, testResults.failed > 0 ? 'error' : 'info')
    workspace.addTranspilationLog(`   ⏭️  Skipped: ${testResults.skipped}`, 'warning')
    workspace.addTranspilationLog('', 'info')
    
    // Display individual test results
    if (testResults.tests && testResults.tests.length > 0) {
      workspace.addTranspilationLog('📝 Test Results:', 'info')
      for (const test of testResults.tests) {
        const icon = test.Action === 'pass' ? '✅' : test.Action === 'fail' ? '❌' : '⏭️'
        const level = test.Action === 'pass' ? 'success' : test.Action === 'fail' ? 'error' : 'warning'
        workspace.addTranspilationLog(
          `${icon} ${test.Test} (${test.Elapsed.toFixed(3)}s)`,
          level
        )
        if (test.Output && test.Action === 'fail') {
          workspace.addTranspilationLog(`   ${test.Output}`, 'error')
        }
      }
    }
    
    if (result.stderr && result.stderr.trim()) {
      workspace.addTranspilationLog('', 'info')
      workspace.addTranspilationLog('--- Additional Output ---', 'warning')
      result.stderr.split('\n').forEach(line => {
        if (line.trim()) {
          workspace.addTranspilationLog(line, 'warning')
        }
      })
    }
    
    // Clear current result to show logs
    transpileStore.currentResult = null
    
    // Show toast notification
    toast.add({
      severity: result.success ? 'success' : 'error',
      summary: result.success ? 'Tests Passed' : 'Tests Failed',
      detail: result.success
        ? `All ${testResults.passed} tests passed in ${result.duration_ms}ms`
        : `${testResults.failed} test(s) failed out of ${testResults.total}`,
      life: result.success ? 4000 : 8000
    })
  } catch (error) {
    workspace.addTranspilationLog(`Test error: ${error}`, 'error')
    transpileStore.currentResult = null
    
    // Check if error is Go-related
    const errorStr = String(error).toLowerCase()
    const isGoError = errorStr.includes('go not found') || 
                      errorStr.includes('go binary') || 
                      errorStr.includes('install go') ||
                      errorStr.includes('configure in settings')
    
    if (isGoError) {
      toast.add({
        severity: 'error',
        summary: 'Go Compiler Not Found',
        detail: 'Configure Go in Settings to use Build, Test, and Run features',
        life: 8000,
        group: 'go-error',
        closable: true
      })
    } else {
      toast.add({
        severity: 'error',
        summary: 'Tests Failed',
        detail: String(error),
        life: 5000
      })
    }
  }
}

// State restoration
async function restoreTranspilationState() {
  if (!projectPath.value) return
  
  restoringState.value = true
  hasValidTranspilationState.value = false
  
  try {
    // Check if there's a saved transpilation state for this project
    const state = transpileStore.getTranspilationState(projectPath.value)
    
    if (!state) {
      restoringState.value = false
      return
    }
    
    // Verify the output directory still exists
    const isValid = await transpileStore.verifyTranspilationState(projectPath.value)
    
    if (!isValid) {
      // State is stale, clear it
      workspace.addTranspilationLog('Previous transpilation output not found, state cleared', 'warning')
      transpileStore.clearTranspilationState(projectPath.value)
      restoringState.value = false
      return
    }
    
    // Restore state
    hasValidTranspilationState.value = true
    lastTranspilationTime.value = new Date(state.timestamp).toLocaleString()
    lastTranspilationFiles.value = state.filesTranspiled
    
    // Restore the currentResult so Build/Test/Run buttons are enabled
    transpileStore.currentResult = {
      success: true,
      output_dir: state.outputDir,
      files_transpiled: state.filesTranspiled,
      message: 'Previous transpilation state restored'
    }
    
    workspace.addTranspilationLog(`Restored previous transpilation from ${lastTranspilationTime.value}`, 'success')
    workspace.addTranspilationLog(`Output directory: ${state.outputDir}`, 'info')
  } catch (error) {
    console.error('Failed to restore transpilation state:', error)
    workspace.addTranspilationLog(`Failed to restore state: ${error}`, 'error')
  } finally {
    restoringState.value = false
  }
}

function clearTranspilationState() {
  if (!projectPath.value) return
  
  transpileStore.clearTranspilationState(projectPath.value)
  hasValidTranspilationState.value = false
  lastTranspilationTime.value = null
  lastTranspilationFiles.value = 0
  transpileStore.currentResult = null
  
  workspace.addTranspilationLog('Transpilation state cleared', 'info')
  
  toast.add({
    severity: 'info',
    summary: 'State Cleared',
    detail: 'Previous transpilation state has been cleared',
    life: 3000
  })
}

// Auto-restore state when project is loaded
watch(projectPath, async (newPath) => {
  if (newPath) {
    await restoreTranspilationState()
  }
}, { immediate: true })

onMounted(async () => {
  // Restore state on mount if project is already loaded
  if (projectPath.value) {
    await restoreTranspilationState()
  }
})
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

.btn-success {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.btn-success:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(16, 185, 129, 0.4);
}

.btn-action.btn-help {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333;
  padding: 0.5rem 0.75rem;
}

.btn-action.btn-help:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 237, 234, 0.4);
}

.running-project-indicator {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
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
  height: 100%;
  overflow: hidden;
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

.progress-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 10px;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
  font-family: monospace;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.run-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.run-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(16, 185, 129, 0.3);
}

.running-indicator {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: rgba(251, 191, 36, 0.1);
  color: #f59e0b;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
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

/* State Restoration Banner */
.state-restoration-banner {
  padding: 16px;
  border-bottom: 1px solid var(--color-border);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.banner-content {
  display: flex;
  align-items: start;
  gap: 12px;
}

.banner-icon {
  font-size: 24px;
  flex-shrink: 0;
  margin-top: 2px;
}

.banner-info {
  flex: 1;
  min-width: 0;
}

.banner-info h4 {
  margin: 0 0 6px 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
}

.banner-details {
  margin: 0 0 4px 0;
  font-size: 12px;
  color: var(--color-text-secondary);
  font-family: 'Monaco', 'Menlo', monospace;
}

.banner-hint {
  margin: 0;
  font-size: 11px;
  color: var(--color-text-secondary);
  font-style: italic;
}

.banner-clear-btn {
  padding: 6px 12px;
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.banner-clear-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.3);
  transform: translateY(-1px);
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

.log-output {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-soft);
}

.log-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
}

.clear-log-btn {
  padding: 4px 12px;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all 0.2s;
}

.clear-log-btn:hover {
  background: var(--color-background-soft);
  border-color: var(--color-text-secondary);
}

.log-entries {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
}

.log-entry {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 6px 8px;
  margin-bottom: 4px;
  border-radius: 4px;
  transition: background 0.2s;
}

.log-entry:hover {
  background: var(--color-background-soft);
}

.log-icon {
  font-size: 14px;
  width: 16px;
  text-align: center;
  flex-shrink: 0;
}

.log-time {
  color: var(--color-text-secondary);
  font-size: 11px;
  min-width: 70px;
  flex-shrink: 0;
}

.log-message {
  flex: 1;
  word-break: break-word;
}

.log-info .log-icon {
  color: #3b82f6;
}

.log-success .log-icon {
  color: #10b981;
}

.log-error .log-icon {
  color: #ef4444;
}

.log-warning .log-icon {
  color: #f59e0b;
}

.log-success {
  background: rgba(16, 185, 129, 0.05);
}

.log-error {
  background: rgba(239, 68, 68, 0.05);
}

.log-warning {
  background: rgba(245, 158, 11, 0.05);
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
