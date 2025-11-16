<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { invoke } from '@tauri-apps/api/core'
import { open } from '@tauri-apps/plugin-dialog'
import AppLayout from '../components/AppLayout.vue'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import InputSwitch from 'primevue/inputswitch'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Button from 'primevue/button'
import Message from 'primevue/message'

const settingsStore = useSettingsStore()
const saveMessage = ref(false)

// Go detection state
const goDetecting = ref(false)
const goDetectionResult = ref<{
  found: boolean
  version?: string
  path?: string
  message?: string
} | null>(null)

// Theme options
const themeOptions = [
  { label: 'System', value: 'system' },
  { label: 'Light', value: 'light' },
  { label: 'Dark', value: 'dark' }
]

// Font size options
const fontSizes = [
  { label: 'Small (12px)', value: 12 },
  { label: 'Medium (14px)', value: 14 },
  { label: 'Large (16px)', value: 16 },
  { label: 'Extra Large (18px)', value: 18 }
]

// Tab size options
const tabSizes = [
  { label: '2 spaces', value: 2 },
  { label: '4 spaces', value: 4 },
  { label: '8 spaces', value: 8 }
]

// Go binary source options
const goBinarySourceOptions = [
  { label: 'System Go (from PATH)', value: 'system' },
  { label: 'Custom Path', value: 'custom' }
]

// Watch for settings changes and show save message
watch(() => settingsStore.settings, () => {
  saveMessage.value = true
  setTimeout(() => {
    saveMessage.value = false
  }, 3000)
}, { deep: true })

function resetToDefaults() {
  settingsStore.resetToDefaults()
}

function selectOutputDirectory() {
  // This will be implemented with Tauri dialog
  console.log('Select output directory')
}

async function browseForGoBinary() {
  try {
    const selected = await open({
      multiple: false,
      directory: false,
      title: 'Select Go Binary'
    })
    
    if (selected && typeof selected === 'string') {
      settingsStore.settings.customGoBinaryPath = selected
      // Auto-detect after selecting
      await detectGo()
    }
  } catch (error) {
    console.error('Failed to browse for Go binary:', error)
  }
}

async function detectGo() {
  goDetecting.value = true
  goDetectionResult.value = null
  
  try {
    const customPath = settingsStore.settings.goBinarySource === 'custom' 
      ? settingsStore.settings.customGoBinaryPath 
      : null
    
    const result = await invoke<{
      found: boolean
      version?: string
      path?: string
      message?: string
    }>('detect_go_installation', { customPath })
    
    goDetectionResult.value = result
  } catch (error) {
    goDetectionResult.value = {
      found: false,
      message: `Error detecting Go: ${error}`
    }
  } finally {
    goDetecting.value = false
  }
}

// Auto-detect Go on mount
onMounted(() => {
  detectGo()
})
</script>

<template>
  <AppLayout>
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Settings</h2>
          <Button 
            label="Reset to Defaults" 
            icon="pi pi-refresh" 
            severity="secondary" 
            size="small"
            @click="resetToDefaults"
          />
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-auto p-6 bg-gray-50 dark:bg-gray-900">
        <div class="max-w-4xl mx-auto">
          <!-- Save Message -->
          <Message v-if="saveMessage" severity="success" :closable="false" class="mb-4">
            Settings saved automatically
          </Message>

          <!-- Settings Tabs -->
          <TabView class="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
            <!-- Application Settings -->
            <TabPanel header="Application" value="0">
              <div class="space-y-6">
                <!-- Theme -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Theme
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Choose your preferred color scheme
                    </p>
                  </div>
                  <Dropdown 
                    v-model="settingsStore.settings.theme" 
                    :options="themeOptions" 
                    optionLabel="label" 
                    optionValue="value"
                    class="w-48"
                  />
                </div>

                <!-- Font Size -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Font Size
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Adjust the editor font size
                    </p>
                  </div>
                  <Dropdown 
                    v-model="settingsStore.settings.fontSize" 
                    :options="fontSizes" 
                    optionLabel="label" 
                    optionValue="value"
                    class="w-48"
                  />
                </div>

                <!-- Auto Save -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Auto Save
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Automatically save changes
                    </p>
                  </div>
                  <InputSwitch v-model="settingsStore.settings.autoSave" />
                </div>

                <!-- Auto Save Delay -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div class="flex-1 mr-4">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Auto Save Delay
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
                      Time to wait before auto-saving (milliseconds)
                    </p>
                    <InputText 
                      v-model.number="settingsStore.settings.autoSaveDelay" 
                      type="number"
                      :min="1000"
                      :max="10000"
                      :step="500"
                      class="w-full"
                      :disabled="!settingsStore.settings.autoSave"
                    />
                  </div>
                  <span class="text-sm text-gray-500 dark:text-gray-400 mt-6">
                    {{ (settingsStore.settings.autoSaveDelay / 1000).toFixed(1) }}s
                  </span>
                </div>

                <!-- Enable Backups -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Enable Backups
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Create backup copies when saving files
                    </p>
                  </div>
                  <InputSwitch v-model="settingsStore.settings.enableBackups" />
                </div>

                <!-- Backup Location -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div class="flex-1 mr-4">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Backup Location
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
                      Directory where backup files are stored
                    </p>
                    <InputText 
                      v-model="settingsStore.settings.backupLocation" 
                      class="w-full"
                      placeholder="./.backups"
                      :disabled="!settingsStore.settings.enableBackups"
                    />
                  </div>
                  <Button 
                    icon="pi pi-folder-open" 
                    severity="secondary"
                    @click="selectOutputDirectory"
                    class="mt-6"
                    :disabled="!settingsStore.settings.enableBackups"
                  />
                </div>

                <!-- Default Output Directory -->
                <div class="flex items-center justify-between py-3">
                  <div class="flex-1 mr-4">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Default Output Directory
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
                      Where transpiled Go files will be saved
                    </p>
                    <InputText 
                      v-model="settingsStore.settings.defaultOutputDir" 
                      class="w-full"
                      placeholder="./output"
                    />
                  </div>
                  <Button 
                    icon="pi pi-folder-open" 
                    severity="secondary"
                    @click="selectOutputDirectory"
                    class="mt-6"
                  />
                </div>
              </div>
            </TabPanel>

            <!-- Project Settings -->
            <TabPanel header="Project" value="1">
              <div class="space-y-6">
                <!-- Go Module Name -->
                <div class="py-3 border-b border-gray-200 dark:border-gray-700">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Go Module Name
                  </label>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
                    Override the default Go module name
                  </p>
                  <InputText 
                    v-model="settingsStore.settings.goModuleName" 
                    class="w-full"
                    placeholder="github.com/username/project"
                  />
                </div>

                <!-- Exclude Patterns -->
                <div class="py-3 border-b border-gray-200 dark:border-gray-700">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Exclude Patterns
                  </label>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
                    Glob patterns for files to exclude (comma separated)
                  </p>
                  <InputText 
                    v-model="settingsStore.settings.excludePatterns" 
                    class="w-full"
                    placeholder="node_modules, **/*.test.ts, dist"
                  />
                </div>

                <!-- Include Patterns -->
                <div class="py-3 border-b border-gray-200 dark:border-gray-700">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Include Patterns
                  </label>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
                    Glob patterns for files to include (comma separated)
                  </p>
                  <InputText 
                    v-model="settingsStore.settings.includePatterns" 
                    class="w-full"
                    placeholder="**/*.ts, **/*.tsx"
                  />
                </div>

                <!-- Go Configuration -->
                <div class="py-3">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Go Compiler Configuration
                  </label>
                  
                  <!-- Go Binary Source -->
                  <div class="mb-4">
                    <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">
                      Go Binary Source
                    </label>
                    <Dropdown 
                      v-model="settingsStore.settings.goBinarySource" 
                      :options="goBinarySourceOptions" 
                      optionLabel="label" 
                      optionValue="value"
                      class="w-full"
                    />
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      Choose where to find the Go compiler
                    </p>
                  </div>

                  <!-- Custom Go Path (shown only when custom is selected) -->
                  <div v-if="settingsStore.settings.goBinarySource === 'custom'" class="mb-4">
                    <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">
                      Custom Go Binary Path
                    </label>
                    <div class="flex gap-2">
                      <InputText 
                        v-model="settingsStore.settings.customGoBinaryPath" 
                        class="flex-1"
                        placeholder="/usr/local/go/bin/go"
                      />
                      <Button 
                        icon="pi pi-folder-open" 
                        severity="secondary" 
                        @click="browseForGoBinary"
                        title="Browse for Go binary"
                      />
                    </div>
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      Full path to the Go executable
                    </p>
                  </div>

                  <!-- Detect Go Installation -->
                  <div class="mb-4">
                    <Button 
                      label="Detect Go Installation" 
                      icon="pi pi-search" 
                      @click="detectGo"
                      :loading="goDetecting"
                      class="w-full"
                    />
                  </div>

                  <!-- Go Detection Result -->
                  <div v-if="goDetectionResult" class="p-3 rounded-lg" :class="goDetectionResult.found ? 'bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800' : 'bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800'">
                    <div class="flex items-start gap-2">
                      <i :class="goDetectionResult.found ? 'pi pi-check-circle text-green-600 dark:text-green-400' : 'pi pi-times-circle text-red-600 dark:text-red-400'" class="mt-0.5"></i>
                      <div class="flex-1">
                        <p class="text-sm font-medium" :class="goDetectionResult.found ? 'text-green-800 dark:text-green-300' : 'text-red-800 dark:text-red-300'">
                          {{ goDetectionResult.found ? 'Go Compiler Detected' : 'Go Compiler Not Found' }}
                        </p>
                        <div v-if="goDetectionResult.found" class="mt-1 text-xs space-y-1">
                          <p class="text-gray-700 dark:text-gray-300">
                            <span class="font-medium">Version:</span> {{ goDetectionResult.version }}
                          </p>
                          <p class="text-gray-700 dark:text-gray-300 break-all">
                            <span class="font-medium">Path:</span> {{ goDetectionResult.path }}
                          </p>
                        </div>
                        <p v-else class="mt-1 text-xs text-red-700 dark:text-red-300">
                          {{ goDetectionResult.message }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </TabPanel>

            <!-- Editor Settings -->
            <TabPanel header="Editor" value="2">
              <div class="space-y-6">
                <!-- Tab Size -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Tab Size
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Number of spaces per tab
                    </p>
                  </div>
                  <Dropdown 
                    v-model="settingsStore.settings.tabSize" 
                    :options="tabSizes" 
                    optionLabel="label" 
                    optionValue="value"
                    class="w-48"
                  />
                </div>

                <!-- Word Wrap -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Word Wrap
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Wrap long lines in the editor
                    </p>
                  </div>
                  <InputSwitch v-model="settingsStore.settings.wordWrap" />
                </div>

                <!-- Line Numbers -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Line Numbers
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Show line numbers in the editor
                    </p>
                  </div>
                  <InputSwitch v-model="settingsStore.settings.lineNumbers" />
                </div>

                <!-- Minimap -->
                <div class="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Minimap
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Show code minimap overview
                    </p>
                  </div>
                  <InputSwitch v-model="settingsStore.settings.minimap" />
                </div>

                <!-- Auto Format on Save -->
                <div class="flex items-center justify-between py-3">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Auto Format on Save
                    </label>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      Automatically format code when saving
                    </p>
                  </div>
                  <InputSwitch v-model="settingsStore.settings.autoFormatOnSave" />
                </div>
              </div>
            </TabPanel>
          </TabView>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
/* Custom styles for settings view */
</style>
