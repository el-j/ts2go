<template>
  <div class="app-container">
    <!-- Top Toolbar -->
    <Toolbar class="top-toolbar">
      <template #start>
        <div class="flex items-center gap-4">
          <i class="pi pi-code text-3xl text-primary"></i>
          <span class="text-2xl font-bold">TS2Go</span>
          <span class="text-sm text-gray-500">TypeScript to Go Transpiler</span>
        </div>
      </template>
      <template #end>
        <div class="flex gap-2">
          <Button icon="pi pi-cog" severity="secondary" text @click="showSettings = true" />
          <Button :icon="isDark ? 'pi pi-sun' : 'pi pi-moon'" severity="secondary" text @click="toggleTheme" />
        </div>
      </template>
    </Toolbar>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Left Sidebar - File Tree -->
      <div class="sidebar">
        <Panel header="Files" class="h-full">
          <div class="file-upload-area">
            <FileUpload
              mode="advanced"
              name="files[]"
              accept=".ts,.tsx"
              :multiple="true"
              :maxFileSize="10000000"
              :showUploadButton="false"
              :showCancelButton="false"
              @select="onFilesSelect"
            >
              <template #empty>
                <div class="text-center p-4">
                  <i class="pi pi-cloud-upload text-4xl text-gray-400 mb-2"></i>
                  <p class="text-gray-600">Drag & drop TypeScript files here</p>
                  <p class="text-sm text-gray-400 mt-2">or click to browse</p>
                </div>
              </template>
            </FileUpload>
          </div>

          <Divider />

          <div v-if="files.length > 0" class="file-list">
            <div
              v-for="file in files"
              :key="file.name"
              class="file-item"
              :class="{ active: selectedFile === file }"
              @click="selectFile(file)"
            >
              <i class="pi pi-file text-blue-500"></i>
              <span>{{ file.name }}</span>
            </div>
          </div>
          <div v-else class="text-center text-gray-400 py-4">
            No files loaded
          </div>
        </Panel>
      </div>

      <!-- Center - Editor Area -->
      <div class="editor-area">
        <Splitter>
          <SplitterPanel :size="50" :minSize="20">
            <Panel header="TypeScript" class="h-full">
              <div class="editor-container">
                <Textarea
                  v-model="typescriptCode"
                  class="editor-textarea"
                  :autoResize="false"
                  placeholder="TypeScript code will appear here..."
                />
              </div>
            </Panel>
          </SplitterPanel>
          <SplitterPanel :size="50" :minSize="20">
            <Panel header="Go Code" class="h-full">
              <div class="editor-container">
                <Textarea
                  v-model="goCode"
                  class="editor-textarea"
                  :autoResize="false"
                  readonly
                  placeholder="Transpiled Go code will appear here..."
                />
              </div>
            </Panel>
          </SplitterPanel>
        </Splitter>

        <!-- Action Buttons -->
        <div class="action-bar">
          <Button
            label="Transpile"
            icon="pi pi-play"
            :loading="isTranspiling"
            @click="transpileCode"
          />
          <Button
            label="Copy Go Code"
            icon="pi pi-copy"
            severity="secondary"
            :disabled="!goCode"
            @click="copyGoCode"
          />
          <Button
            label="Download"
            icon="pi pi-download"
            severity="secondary"
            :disabled="!goCode"
            @click="downloadGoCode"
          />
        </div>
      </div>
    </div>

    <!-- Status Bar -->
    <div class="status-bar">
      <div class="flex gap-4">
        <span>{{ files.length }} files</span>
        <span v-if="errorCount > 0" class="text-red-500">{{ errorCount }} errors</span>
        <span v-else class="text-green-500">Ready</span>
      </div>
      <div>
        <span v-if="lastTranspileTime" class="text-gray-500">
          Last transpiled: {{ lastTranspileTime }}
        </span>
      </div>
    </div>

    <!-- Settings Sidebar -->
    <Sidebar v-model:visible="showSettings" position="right" class="settings-sidebar">
      <template #header>
        <h3>Settings</h3>
      </template>
      <TabView>
        <TabPanel header="General">
          <div class="setting-item">
            <label>Theme</label>
            <SelectButton v-model="themeMode" :options="themeOptions" optionLabel="label" optionValue="value" />
          </div>
          <div class="setting-item">
            <label>Auto-transpile</label>
            <InputSwitch v-model="autoTranspile" />
          </div>
        </TabPanel>
        <TabPanel header="Transpiler">
          <div class="setting-item">
            <label>Module Name</label>
            <InputText v-model="moduleName" placeholder="github.com/example/project" />
          </div>
          <div class="setting-item">
            <label>Optimization</label>
            <Dropdown v-model="optimization" :options="optimizationOptions" placeholder="Select optimization level" />
          </div>
        </TabPanel>
      </TabView>
    </Sidebar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Toolbar from 'primevue/toolbar'
import Button from 'primevue/button'
import Panel from 'primevue/panel'
import FileUpload from 'primevue/fileupload'
import Textarea from 'primevue/textarea'
import Splitter from 'primevue/splitter'
import SplitterPanel from 'primevue/splitterpanel'
import Sidebar from 'primevue/sidebar'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import SelectButton from 'primevue/selectbutton'
import InputSwitch from 'primevue/inputswitch'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Divider from 'primevue/divider'

// State
const files = ref<File[]>([])
const selectedFile = ref<File | null>(null)
const typescriptCode = ref('')
const goCode = ref('')
const isTranspiling = ref(false)
const errorCount = ref(0)
const lastTranspileTime = ref('')
const showSettings = ref(false)
const isDark = ref(false)
const themeMode = ref('light')
const autoTranspile = ref(false)
const moduleName = ref('main')
const optimization = ref('basic')

const themeOptions = [
  { label: 'Light', value: 'light' },
  { label: 'Dark', value: 'dark' },
  { label: 'System', value: 'system' }
]

const optimizationOptions = ['none', 'basic', 'aggressive']

// Methods
const onFilesSelect = (event: any) => {
  const newFiles = event.files
  files.value = [...files.value, ...newFiles]
  if (newFiles.length > 0 && !selectedFile.value) {
    selectFile(newFiles[0])
  }
}

const selectFile = async (file: File) => {
  selectedFile.value = file
  const text = await file.text()
  typescriptCode.value = text
  if (autoTranspile.value) {
    transpileCode()
  }
}

const transpileCode = async () => {
  isTranspiling.value = true
  errorCount.value = 0
  
  try {
    // TODO: Implement actual transpilation via Tauri command
    // For now, just a placeholder
    await new Promise(resolve => setTimeout(resolve, 500))
    goCode.value = `package main\n\n// Transpiled from: ${selectedFile.value?.name || 'input.ts'}\n\nfunc main() {\n\t// Generated Go code will appear here\n}\n`
    lastTranspileTime.value = new Date().toLocaleTimeString()
  } catch (error) {
    errorCount.value = 1
    console.error('Transpilation error:', error)
  } finally {
    isTranspiling.value = false
  }
}

const copyGoCode = () => {
  navigator.clipboard.writeText(goCode.value)
}

const downloadGoCode = () => {
  const blob = new Blob([goCode.value], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = selectedFile.value ? selectedFile.value.name.replace('.ts', '.go') : 'output.go'
  a.click()
  URL.revokeObjectURL(url)
}

const toggleTheme = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark-mode')
}
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
}

.top-toolbar {
  border-bottom: 1px solid var(--surface-border);
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 300px;
  border-right: 1px solid var(--surface-border);
  overflow-y: auto;
}

.editor-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.editor-container {
  height: 100%;
  overflow: hidden;
}

.editor-textarea {
  width: 100%;
  height: 100%;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  resize: none;
  border: none;
  outline: none;
}

.file-upload-area {
  margin-bottom: 1rem;
}

.file-list {
  max-height: 400px;
  overflow-y: auto;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.file-item:hover {
  background-color: var(--surface-hover);
}

.file-item.active {
  background-color: var(--primary-color);
  color: white;
}

.action-bar {
  display: flex;
  gap: 0.5rem;
  padding: 1rem;
  border-top: 1px solid var(--surface-border);
  background: var(--surface-ground);
}

.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  border-top: 1px solid var(--surface-border);
  background: var(--surface-ground);
  font-size: 0.875rem;
}

.settings-sidebar {
  width: 400px;
}

.setting-item {
  margin-bottom: 1.5rem;
}

.setting-item label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}
</style>
