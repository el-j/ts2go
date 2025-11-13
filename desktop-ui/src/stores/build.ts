import { defineStore } from 'pinia'
import { ref } from 'vue'
import { invoke } from '@tauri-apps/api/core'

export interface BuildResult {
  success: boolean
  binary_path: string
  binary_size: number
  stdout: string
  stderr: string
  errors: string[]
  warnings: string[]
  exit_code: number
  duration_ms: number
  timestamp: number
}

export const useBuildStore = defineStore('build', () => {
  // State
  const isBuilding = ref(false)
  const buildResult = ref<BuildResult | null>(null)
  const buildHistory = ref<BuildResult[]>([])
  
  // Actions
  async function buildFile(code: string, outputPath: string): Promise<BuildResult> {
    isBuilding.value = true
    
    try {
      const result = await invoke<Omit<BuildResult, 'timestamp'>>('build_go_file', {
        code,
        outputPath
      })
      
      const buildRes: BuildResult = {
        ...result,
        timestamp: Date.now()
      }
      
      buildResult.value = buildRes
      
      // Add to history
      buildHistory.value.unshift(buildRes)
      
      // Keep only last 50 builds in history
      if (buildHistory.value.length > 50) {
        buildHistory.value = buildHistory.value.slice(0, 50)
      }
      
      return buildRes
    } finally {
      isBuilding.value = false
    }
  }
  
  async function buildProject(
    sourceDir: string,
    outputPath: string,
    buildFlags?: string[]
  ): Promise<BuildResult> {
    isBuilding.value = true
    
    try {
      const result = await invoke<Omit<BuildResult, 'timestamp'>>('build_go_project', {
        sourceDir,
        outputPath,
        buildFlags
      })
      
      const buildRes: BuildResult = {
        ...result,
        timestamp: Date.now()
      }
      
      buildResult.value = buildRes
      
      // Add to history
      buildHistory.value.unshift(buildRes)
      
      // Keep only last 50 builds in history
      if (buildHistory.value.length > 50) {
        buildHistory.value = buildHistory.value.slice(0, 50)
      }
      
      return buildRes
    } finally {
      isBuilding.value = false
    }
  }
  
  function clearResult() {
    buildResult.value = null
  }
  
  function clearHistory() {
    buildHistory.value = []
  }
  
  return {
    // State
    isBuilding,
    buildResult,
    buildHistory,
    
    // Actions
    buildFile,
    buildProject,
    clearResult,
    clearHistory
  }
})
