import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { useBackendState } from '../composables/useBackendState'
import type { TranspileResult, TranspilationState } from '../types'

export const useTranspileStore = defineStore('transpile', () => {
  const isTranspiling = ref(false)
  const currentResult = ref<TranspileResult | null>(null)
  const transpileHistory = ref<TranspileResult[]>([])
  const startTime = ref<number>(0)
  const elapsedTime = ref<number>(0)
  const transpilationStates = ref<Map<string, TranspilationState>>(new Map())
  const backend = useBackendState<Array<[string, TranspilationState]>>('transpilation-states', [])
  let isLoading = true
  
  let timerInterval: number | null = null
  
  const hasError = computed(() => currentResult.value && !currentResult.value.success)
  const hasSuccess = computed(() => currentResult.value && currentResult.value.success)
  
  // Load transpilation states from backend
  async function loadTranspilationStates() {
    try {
      const stored = await backend.load()
      transpilationStates.value = new Map(stored)
    } catch (e) {
      console.error('Failed to load transpilation states from backend:', e)
    } finally {
      isLoading = false
    }
  }
  
  // Save transpilation states to backend
  async function saveTranspilationStates() {
    if (isLoading) return
    
    try {
      const states = Array.from(transpilationStates.value.entries())
      await backend.save(states)
    } catch (e) {
      console.error('Failed to save transpilation states to backend:', e)
    }
  }
  
  // Save transpilation state for a project
  function saveTranspilationState(projectPath: string, result: TranspileResult) {
    if (result.success && result.output_dir) {
      transpilationStates.value.set(projectPath, {
        projectPath,
        outputDir: result.output_dir,
        filesTranspiled: result.files_transpiled || 0,
        timestamp: new Date().toISOString(),
        success: true
      })
      saveTranspilationStates()
    }
  }
  
  // Get transpilation state for a project
  function getTranspilationState(projectPath: string): TranspilationState | null {
    return transpilationStates.value.get(projectPath) || null
  }
  
  // Clear transpilation state for a project
  function clearTranspilationState(projectPath: string) {
    transpilationStates.value.delete(projectPath)
    saveTranspilationStates()
  }
  
  // Clear all transpilation states
  function clearAllTranspilationStates() {
    transpilationStates.value.clear()
    saveTranspilationStates()
  }
  
  // Verify if output directory still exists
  async function verifyTranspilationState(projectPath: string): Promise<boolean> {
    const state = getTranspilationState(projectPath)
    if (!state) return false
    
    try {
      const result = await invoke<{ exists: boolean }>('check_directory_exists', {
        path: state.outputDir
      })
      
      if (!result.exists) {
        clearTranspilationState(projectPath)
        return false
      }
      
      return true
    } catch {
      return false
    }
  }
  
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
      
      // Save transpilation state for this project
      if (currentResult.value && currentResult.value.success) {
        saveTranspilationState(projectPath, currentResult.value)
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
  
  // Initialize: load states from localStorage
  loadTranspilationStates()
  
  return {
    // State
    isTranspiling,
    currentResult,
    transpileHistory,
    elapsedTime,
    transpilationStates,
    
    // Computed
    hasError,
    hasSuccess,
    
    // Actions
    transpileProject,
    clearResult,
    clearHistory,
    getTranspilationState,
    saveTranspilationState,
    clearTranspilationState,
    clearAllTranspilationStates,
    verifyTranspilationState
  }
})
