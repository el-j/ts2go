<template>
  <div class="log-viewer flex flex-col h-full">
    <!-- Header -->
    <div class="bg-gray-100 dark:bg-gray-900 px-4 py-2 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
      <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">
        <i class="pi pi-list mr-2"></i>
        Logs ({{ filteredLogs.length }})
      </span>
      
      <div class="flex items-center gap-2">
        <!-- Filter by level -->
        <select 
          v-model="filterLevel" 
          class="px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800"
        >
          <option value="all">All Levels</option>
          <option value="info">Info</option>
          <option value="success">Success</option>
          <option value="warning">Warning</option>
          <option value="error">Error</option>
        </select>
        
        <!-- Auto-scroll toggle -->
        <Button 
          @click="autoScroll = !autoScroll"
          :class="['px-2 py-1 text-sm rounded', autoScroll ? 'bg-primary-600 text-white' : 'bg-gray-200 dark:bg-gray-700']"
          title="Auto-scroll to latest"
        >
          <i class="pi pi-angle-double-down"></i>
        </Button>
        
        <!-- Clear logs -->
        <Button 
          @click="clearLogs"
          class="px-2 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600"
          title="Clear all logs"
        >
          <i class="pi pi-trash"></i>
        </Button>
        
        <!-- Export logs -->
        <Button 
          @click="exportLogs"
          class="px-2 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600"
          title="Export logs"
        >
          <i class="pi pi-download"></i>
        </Button>
      </div>
    </div>
    
    <!-- Search bar -->
    <div class="px-4 py-2 border-b border-gray-200 dark:border-gray-700">
      <input 
        v-model="searchQuery"
        type="text"
        placeholder="Search logs..."
        class="w-full px-3 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800"
      />
    </div>
    
    <!-- Logs list -->
    <div ref="logsContainer" class="flex-1 overflow-auto p-4 space-y-1 text-sm font-mono">
      <div 
        v-for="log in filteredLogs" 
        :key="log.id"
        :class="getLogClass(log.level)"
        class="px-3 py-2 rounded"
      >
        <span class="opacity-70">{{ formatTime(Number(log.timestamp)) }}</span>
        <span :class="getLevelClass(log.level)" class="mx-2 font-bold uppercase">{{ log.level }}</span>
        <span>{{ log.message }}</span>
      </div>
      
      <div v-if="filteredLogs.length === 0" class="text-center text-gray-500 py-8">
        No logs to display
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useLogsStore } from '../stores/logs'

const logsStore = useLogsStore()
const logsContainer = ref<HTMLElement>()

const filterLevel = ref<string>('all')
const searchQuery = ref<string>('')
const autoScroll = ref<boolean>(true)

const filteredLogs = computed(() => {
  let logs = logsStore.logs
  
  // Filter by level
  if (filterLevel.value !== 'all') {
    logs = logs.filter(log => log.level === filterLevel.value)
  }
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    logs = logs.filter(log => log.message.toLowerCase().includes(query))
  }
  
  return logs
})

function getLogClass(level: string): string {
  const baseClass = 'border-l-4 '
  switch (level) {
    case 'info':
      return baseClass + 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
    case 'success':
      return baseClass + 'border-green-500 bg-green-50 dark:bg-green-900/20'
    case 'warning':
      return baseClass + 'border-yellow-500 bg-yellow-50 dark:bg-yellow-900/20'
    case 'error':
      return baseClass + 'border-red-500 bg-red-50 dark:bg-red-900/20'
    default:
      return baseClass + 'border-gray-500 bg-gray-50 dark:bg-gray-900/20'
  }
}

function getLevelClass(level: string): string {
  switch (level) {
    case 'info':
      return 'text-blue-700 dark:text-blue-400'
    case 'success':
      return 'text-green-700 dark:text-green-400'
    case 'warning':
      return 'text-yellow-700 dark:text-yellow-400'
    case 'error':
      return 'text-red-700 dark:text-red-400'
    default:
      return 'text-gray-700 dark:text-gray-400'
  }
}

function formatTime(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', { 
    hour12: false, 
    hour: '2-digit', 
    minute: '2-digit', 
    second: '2-digit' 
  })
}

function clearLogs() {
  logsStore.clearLogs()
}

function exportLogs() {
  const logsText = filteredLogs.value
    .map(log => `[${formatTime(Number(log.timestamp))}] ${log.level.toUpperCase()}: ${log.message}`)
    .join('\n')
  
  const blob = new Blob([logsText], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ts2go-logs-${Date.now()}.txt`
  a.click()
  URL.revokeObjectURL(url)
}

// Auto-scroll to bottom when new logs arrive
watch(() => logsStore.logs.length, async () => {
  if (autoScroll.value) {
    await nextTick()
    if (logsContainer.value) {
      logsContainer.value.scrollTop = logsContainer.value.scrollHeight
    }
  }
})
</script>

<style scoped>
/* .log-viewer {
  background: white;
}

:deep(.dark) .log-viewer {
  background: #1f2937;
} */
</style>
