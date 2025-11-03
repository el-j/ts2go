<template>
  <div class="h-screen flex flex-col">
    <!-- Toolbar -->
    <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-3 flex items-center gap-4">
      <h2 class="text-xl font-semibold">Code Editor</h2>
      
      <div class="flex-1"></div>
      
      <button 
        @click="transpileCode" 
        class="btn-secondary flex items-center gap-2"
        :disabled="transpilerStore.status.isRunning"
      >
        <i class="pi pi-play"></i>
        <span>{{ transpilerStore.status.isRunning ? 'Transpiling...' : 'Transpile' }}</span>
      </button>
      
      <button @click="clearCode" class="px-3 py-2 text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200">
        <i class="pi pi-trash"></i>
      </button>
    </div>

    <!-- Split Pane Editor -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Left Pane: TypeScript Input -->
      <div class="flex-1 flex flex-col border-r border-gray-200 dark:border-gray-700">
        <div class="bg-gray-100 dark:bg-gray-900 px-4 py-2 border-b border-gray-200 dark:border-gray-700">
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">
            <i class="pi pi-code mr-2"></i>
            TypeScript Input
          </span>
        </div>
        <div class="flex-1">
          <CodeEditor 
            v-model="typescriptCode" 
            language="typescript"
            theme="vs-dark"
          />
        </div>
      </div>

      <!-- Right Pane: Go Output -->
      <div class="flex-1 flex flex-col">
        <div class="bg-gray-100 dark:bg-gray-900 px-4 py-2 border-b border-gray-200 dark:border-gray-700">
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">
            <i class="pi pi-file-code mr-2"></i>
            Generated Go Code
          </span>
        </div>
        <div class="flex-1">
          <CodeEditor 
            v-model="goCode" 
            language="go"
            :readonly="true"
            theme="vs-dark"
          />
        </div>
      </div>
    </div>

    <!-- Progress Bar (shown when transpiling) -->
    <div 
      v-if="transpilerStore.status.isRunning" 
      class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 px-6 py-3"
    >
      <div class="flex items-center gap-4">
        <div class="flex-1">
          <div class="flex justify-between text-sm mb-1">
            <span>{{ transpilerStore.status.currentFile || 'Initializing...' }}</span>
            <span>{{ transpilerStore.status.filesProcessed }} / {{ transpilerStore.status.totalFiles }} files</span>
          </div>
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
            <div 
              class="bg-primary-600 h-2 rounded-full transition-all duration-300"
              :style="{ width: `${transpilerStore.status.progress}%` }"
            ></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import CodeEditor from '../components/CodeEditor.vue'
import { useTranspilerStore } from '../stores/transpiler'
import { useLogsStore } from '../stores/logs'
import { invoke } from '@tauri-apps/api/core'

const transpilerStore = useTranspilerStore()
const logsStore = useLogsStore()

const typescriptCode = ref(`// TypeScript Example
interface Person {
  name: string;
  age: number;
  email?: string;
}

function greet(person: Person): string {
  return \`Hello, \${person.name}! You are \${person.age} years old.\`;
}

const user: Person = {
  name: "Alice",
  age: 30,
  email: "alice@example.com"
};

console.log(greet(user));
`)

const goCode = ref('// Click "Transpile" to generate Go code')

async function transpileCode() {
  if (!typescriptCode.value.trim()) {
    logsStore.addLog('error', 'No TypeScript code to transpile')
    return
  }

  try {
    transpilerStore.startTranspilation(1)
    logsStore.addLog('info', 'Starting transpilation...')
    
    transpilerStore.updateProgress(0, 'input.ts')
    
    // Call Tauri command to transpile
    const result = await invoke<string>('transpile_code', {
      code: typescriptCode.value,
      filename: 'input.ts'
    })
    
    goCode.value = result
    transpilerStore.updateProgress(1, 'input.ts')
    logsStore.addLog('success', 'Transpilation completed successfully')
    
  } catch (error: any) {
    logsStore.addLog('error', `Transpilation failed: ${error}`)
    goCode.value = `// Error during transpilation:\n// ${error}`
  } finally {
    transpilerStore.stopTranspilation()
  }
}

function clearCode() {
  typescriptCode.value = ''
  goCode.value = '// Click "Transpile" to generate Go code'
  logsStore.addLog('info', 'Editor cleared')
}
</script>
