import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { invoke } from '@tauri-apps/api/core'

export interface TranspileResult {
  success: boolean
  output_dir?: string
  files_transpiled?: number
  message?: string
  error?: string
  duration?: number
  goCode?: string  // For single file transpilation
}

export const useTranspileStore = defineStore('transpile', () => {
  const isTranspiling = ref(false)
  const currentResult = ref<TranspileResult | null>(null)
  const transpileHistory = ref<TranspileResult[]>([])
  const startTime = ref<number>(0)
  const elapsedTime = ref<number>(0)
  
  let timerInterval: number | null = null
  
  const hasError = computed(() => currentResult.value && !currentResult.value.success)
  const hasSuccess = computed(() => currentResult.value && currentResult.value.success)
  
  async function transpileProject(projectPath: string) {
    isTranspiling.value = true
    currentResult.value = null
    startTime.value = Date.now()
    elapsedTime.value = 0
    
    // Start timer
    timerInterval = window.setInterval(() => {
      elapsedTime.value = Math.floor((Date.now() - startTime.value) / 1000)
    }, 100)
    
    try {
      const result = await invoke<any>('auto_transpile_project', {
        path: projectPath
      })
      
      const duration = Math.floor((Date.now() - startTime.value) / 1000)
      
      currentResult.value = {
        ...result,
        duration
      }
      
      if (currentResult.value) {
        transpileHistory.value.unshift(currentResult.value)
      }
      
      // Keep only last 10 results
      if (transpileHistory.value.length > 10) {
        transpileHistory.value = transpileHistory.value.slice(0, 10)
      }
      
      return currentResult.value
    } catch (error: any) {
      const duration = Math.floor((Date.now() - startTime.value) / 1000)
      
      currentResult.value = {
        success: false,
        error: error.toString(),
        duration
      }
      
      if (currentResult.value) {
        transpileHistory.value.unshift(currentResult.value)
      }
      
      throw error
    } finally {
      isTranspiling.value = false
      
      if (timerInterval) {
        clearInterval(timerInterval)
        timerInterval = null
      }
    }
  }
  
  function clearResult() {
    currentResult.value = null
  }
  
  function clearHistory() {
    transpileHistory.value = []
  }
  
  return {
    // State
    isTranspiling,
    currentResult,
    transpileHistory,
    elapsedTime,
    
    // Computed
    hasError,
    hasSuccess,
    
    // Actions
    transpileProject,
    clearResult,
    clearHistory
  }
})
