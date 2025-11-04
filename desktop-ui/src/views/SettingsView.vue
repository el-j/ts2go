<script setup lang="ts">
import { ref, watch } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import InputNumber from 'primevue/inputnumber'
import InputSwitch from 'primevue/inputswitch'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Button from 'primevue/button'
import Message from 'primevue/message'

const settingsStore = useSettingsStore()
const saveMessage = ref(false)

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
</script>

<template>
  <div class="h-screen flex">
    <div class="flex-1 flex flex-col">
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
            <TabPanel header="Application">
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
            <TabPanel header="Project">
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
                <div class="py-3">
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
              </div>
            </TabPanel>

            <!-- Editor Settings -->
            <TabPanel header="Editor">
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
  </div>
</template>

<style scoped>
/* Custom styles for settings view */
</style>
