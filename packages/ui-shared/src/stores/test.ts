import { defineStore } from 'pinia'
import { ref } from 'vue'
import { invoke } from '@tauri-apps/api/core'

export interface TestCase {
  Test: string
  Package: string
  Action: string
  Elapsed: number
  Output?: string
}

export interface TestResults {
  passed: number
  failed: number
  skipped: number
  total: number
  tests: TestCase[]
}

export interface TestResult {
  success: boolean
  stdout: string
  stderr: string
  exit_code: number
  duration_ms: number
  test_results: TestResults
  timestamp: number
}

export const useTestStore = defineStore('test', () => {
  // State
  const isTesting = ref(false)
  const testResult = ref<TestResult | null>(null)
  const testHistory = ref<TestResult[]>([])
  
  // Actions
  async function testFile(filePath: string): Promise<TestResult> {
    isTesting.value = true
    
    try {
      const result = await invoke<Omit<TestResult, 'timestamp'>>('test_go_file', {
        filePath
      })
      
      const testRes: TestResult = {
        ...result,
        timestamp: Date.now()
      }
      
      testResult.value = testRes
      
      // Add to history
      testHistory.value.unshift(testRes)
      
      // Keep only last 50 test runs in history
      if (testHistory.value.length > 50) {
        testHistory.value = testHistory.value.slice(0, 50)
      }
      
      return testRes
    } finally {
      isTesting.value = false
    }
  }
  
  async function testProject(
    sourceDir: string,
    testFlags?: string[]
  ): Promise<TestResult> {
    isTesting.value = true
    
    try {
      const result = await invoke<Omit<TestResult, 'timestamp'>>('test_go_project', {
        sourceDir,
        testFlags
      })
      
      const testRes: TestResult = {
        ...result,
        timestamp: Date.now()
      }
      
      testResult.value = testRes
      
      // Add to history
      testHistory.value.unshift(testRes)
      
      // Keep only last 50 test runs in history
      if (testHistory.value.length > 50) {
        testHistory.value = testHistory.value.slice(0, 50)
      }
      
      return testRes
    } finally {
      isTesting.value = false
    }
  }
  
  function clearResult() {
    testResult.value = null
  }
  
  function clearHistory() {
    testHistory.value = []
  }
  
  return {
    // State
    isTesting,
    testResult,
    testHistory,
    
    // Actions
    testFile,
    testProject,
    clearResult,
    clearHistory
  }
})
