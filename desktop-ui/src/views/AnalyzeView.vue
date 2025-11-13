<script setup lang="ts">
import { ref } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { open } from '@tauri-apps/plugin-dialog'
import AppLayout from '../components/AppLayout.vue'

import { useLogsStore } from '@/stores/logs'

const logsStore = useLogsStore()

const projectPath = ref('')
const isAnalyzing = ref(false)
const analysisResult = ref<string>('')
const parsedAnalysis = ref<any>(null)
const errorMessage = ref('')

async function analyzeProject() {
  if (!projectPath.value.trim()) {
    errorMessage.value = 'Please enter a project path'
    return
  }

  isAnalyzing.value = true
  errorMessage.value = ''
  analysisResult.value = ''
  parsedAnalysis.value = null

  try {
    logsStore.addLog('info', `Analyzing project: ${projectPath.value}`)
    
    const result = await invoke<string>('analyze_project', {
      path: projectPath.value
    })
    
    analysisResult.value = result
    
    // Try to parse the result as JSON if possible
    try {
      parsedAnalysis.value = JSON.parse(result)
    } catch {
      // If not JSON, treat as plain text
      parsedAnalysis.value = { rawOutput: result }
    }
    
    logsStore.addLog('success', 'Project analysis completed')
  } catch (error: any) {
    errorMessage.value = error
    logsStore.addLog('error', `Analysis failed: ${error}`)
  } finally {
    isAnalyzing.value = false
  }
}

async function selectDirectory() {
  try {
    const selected = await open({
      directory: true,
      multiple: false,
      title: 'Select TypeScript Project Folder'
    })
    
    if (selected) {
      projectPath.value = selected as string
      errorMessage.value = ''
    }
  } catch (error: any) {
    errorMessage.value = `Failed to select directory: ${error}`
  }
}

function getImportCount(imports: any): number {
  if (!imports) return 0
  if (Array.isArray(imports)) return imports.length
  if (typeof imports === 'object') return Object.keys(imports).length
  return 0
}

function getExportCount(exports: any): number {
  if (!exports) return 0
  if (Array.isArray(exports)) return exports.length
  if (typeof exports === 'object') return Object.keys(exports).length
  return 0
}
</script>

<template>
  <AppLayout>
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">
              <i class="pi pi-search mr-2"></i>
              Project Analysis
            </h2>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Analyze TypeScript project structure, imports, and exports
            </p>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-auto p-6 bg-gray-50 dark:bg-gray-900">
        <div class="max-w-7xl mx-auto space-y-6">
          <!-- Input Section -->
          <Card>
            <template #title>
              <i class="pi pi-folder mr-2"></i>
              Select Project
            </template>
            <template #content>
              <div class="space-y-4">
                <div class="flex gap-2">
                  <div class="flex-1">
                    <input
                      v-model="projectPath"
                      type="text"
                      placeholder="Enter project path (e.g., /path/to/project)"
                      class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg 
                             bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100
                             focus:outline-none focus:ring-2 focus:ring-primary-500"
                      @keyup.enter="analyzeProject"
                    />
                  </div>
                  <Button
                    label="Browse"
                    icon="pi pi-folder-open"
                    @click="selectDirectory"
                    outlined
                  />
                  <Button
                    label="Analyze"
                    icon="pi pi-search"
                    @click="analyzeProject"
                    :loading="isAnalyzing"
                    :disabled="!projectPath.trim()"
                  />
                </div>

                <Message v-if="errorMessage" severity="error" :closable="false">
                  {{ errorMessage }}
                </Message>
              </div>
            </template>
          </Card>

          <!-- Loading -->
          <div v-if="isAnalyzing" class="flex flex-col items-center justify-center py-12">
            <ProgressSpinner />
            <p class="mt-4 text-gray-600 dark:text-gray-400">Analyzing project...</p>
          </div>

          <!-- Results -->
          <div v-if="analysisResult && !isAnalyzing" class="space-y-6">
            <!-- Summary -->
            <Card v-if="parsedAnalysis && !parsedAnalysis.rawOutput">
              <template #title>
                <i class="pi pi-chart-bar mr-2"></i>
                Analysis Summary
              </template>
              <template #content>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                    <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">
                      {{ parsedAnalysis.files?.length || 0 }}
                    </div>
                    <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Total Files</div>
                  </div>

                  <div class="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
                    <div class="text-2xl font-bold text-green-600 dark:text-green-400">
                      {{ getImportCount(parsedAnalysis.imports) }}
                    </div>
                    <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Total Imports</div>
                  </div>

                  <div class="p-4 bg-purple-50 dark:bg-purple-900/20 rounded-lg">
                    <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">
                      {{ getExportCount(parsedAnalysis.exports) }}
                    </div>
                    <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Total Exports</div>
                  </div>
                </div>
              </template>
            </Card>

            <!-- Detailed Results -->
            <Card>
              <template #title>
                <i class="pi pi-list mr-2"></i>
                Detailed Analysis
              </template>
              <template #content>
                <Accordion :multiple="true">
                  <!-- Files -->
                  <AccordionTab v-if="parsedAnalysis?.files">
                    <template #header>
                      <div class="flex items-center gap-2">
                        <i class="pi pi-file"></i>
                        <span class="font-medium">Files</span>
                        <Tag :value="`${parsedAnalysis.files.length}`" severity="info" />
                      </div>
                    </template>
                    
                    <DataTable 
                      :value="parsedAnalysis.files" 
                      stripedRows 
                      paginator 
                      :rows="10"
                    >
                      <Column field="path" header="File Path" sortable></Column>
                      <Column field="size" header="Size" sortable></Column>
                    </DataTable>
                  </AccordionTab>

                  <!-- Imports -->
                  <AccordionTab v-if="parsedAnalysis?.imports">
                    <template #header>
                      <div class="flex items-center gap-2">
                        <i class="pi pi-download"></i>
                        <span class="font-medium">Imports</span>
                        <Tag :value="`${getImportCount(parsedAnalysis.imports)}`" severity="success" />
                      </div>
                    </template>
                    
                    <pre class="bg-gray-900 text-gray-100 p-4 rounded overflow-auto text-sm">{{ JSON.stringify(parsedAnalysis.imports, null, 2) }}</pre>
                  </AccordionTab>

                  <!-- Exports -->
                  <AccordionTab v-if="parsedAnalysis?.exports">
                    <template #header>
                      <div class="flex items-center gap-2">
                        <i class="pi pi-upload"></i>
                        <span class="font-medium">Exports</span>
                        <Tag :value="`${getExportCount(parsedAnalysis.exports)}`" severity="warn" />
                      </div>
                    </template>
                    
                    <pre class="bg-gray-900 text-gray-100 p-4 rounded overflow-auto text-sm">{{ JSON.stringify(parsedAnalysis.exports, null, 2) }}</pre>
                  </AccordionTab>

                  <!-- Warnings -->
                  <AccordionTab v-if="parsedAnalysis?.warnings">
                    <template #header>
                      <div class="flex items-center gap-2">
                        <i class="pi pi-exclamation-triangle"></i>
                        <span class="font-medium">Warnings</span>
                        <Tag :value="`${parsedAnalysis.warnings.length}`" severity="danger" />
                      </div>
                    </template>
                    
                    <div class="space-y-2">
                      <Message 
                        v-for="(warning, index) in parsedAnalysis.warnings" 
                        :key="index"
                        severity="warn"
                        :closable="false"
                      >
                        {{ warning }}
                      </Message>
                    </div>
                  </AccordionTab>

                  <!-- Raw Output -->
                  <AccordionTab v-if="parsedAnalysis?.rawOutput">
                    <template #header>
                      <div class="flex items-center gap-2">
                        <i class="pi pi-code"></i>
                        <span class="font-medium">Raw Output</span>
                      </div>
                    </template>
                    
                    <pre class="bg-gray-900 text-gray-100 p-4 rounded overflow-auto text-sm">{{ parsedAnalysis.rawOutput }}</pre>
                  </AccordionTab>
                </Accordion>
              </template>
            </Card>
          </div>

          <!-- Empty State -->
          <Card v-if="!analysisResult && !isAnalyzing">
            <template #content>
              <div class="text-center py-12">
                <i class="pi pi-search text-6xl text-gray-400 mb-4"></i>
                <h3 class="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">
                  No Analysis Yet
                </h3>
                <p class="text-gray-600 dark:text-gray-400">
                  Enter a project path and click "Analyze" to see the analysis results
                </p>
              </div>
            </template>
          </Card>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
pre {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  line-height: 1.5;
}
</style>
