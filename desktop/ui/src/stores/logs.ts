import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface LogEntry {
  id: string
  timestamp: Date
  level: 'info' | 'warning' | 'error' | 'success'
  message: string
}

export const useLogsStore = defineStore('logs', () => {
  const logs = ref<LogEntry[]>([])

  function addLog(level: LogEntry['level'], message: string) {
    logs.value.push({
      id: Date.now().toString(),
      timestamp: new Date(),
      level,
      message
    })
  }

  function clearLogs() {
    logs.value = []
  }

  return {
    logs,
    addLog,
    clearLogs
  }
})
