import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface TranspilationStatus {
  isRunning: boolean
  progress: number
  currentFile: string
  filesProcessed: number
  totalFiles: number
  startTime: number | null
  estimatedTimeRemaining: number | null
  processingSpeed: number | null
}

export const useTranspilerStore = defineStore('transpiler', () => {
  const status = ref<TranspilationStatus>({
    isRunning: false,
    progress: 0,
    currentFile: '',
    filesProcessed: 0,
    totalFiles: 0,
    startTime: null,
    estimatedTimeRemaining: null,
    processingSpeed: null
  })

  const formattedTimeRemaining = computed(() => {
    if (!status.value.estimatedTimeRemaining) return null
    
    const seconds = Math.floor(status.value.estimatedTimeRemaining / 1000)
    if (seconds < 60) return `${seconds}s`
    
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
    return `${minutes}m ${remainingSeconds}s`
  })

  function startTranspilation(totalFiles: number) {
    status.value = {
      isRunning: true,
      progress: 0,
      currentFile: '',
      filesProcessed: 0,
      totalFiles,
      startTime: Date.now(),
      estimatedTimeRemaining: null,
      processingSpeed: null
    }
  }

  function updateProgress(filesProcessed: number, currentFile: string) {
    status.value.filesProcessed = filesProcessed
    status.value.currentFile = currentFile
    status.value.progress = (filesProcessed / status.value.totalFiles) * 100
    
    // Calculate processing speed and estimated time remaining
    if (status.value.startTime && filesProcessed > 0) {
      const elapsedTime = Date.now() - status.value.startTime
      status.value.processingSpeed = (filesProcessed / elapsedTime) * 1000 // files per second
      
      const remainingFiles = status.value.totalFiles - filesProcessed
      status.value.estimatedTimeRemaining = remainingFiles / status.value.processingSpeed * 1000
    }
  }

  function stopTranspilation() {
    status.value.isRunning = false
  }

  return {
    status,
    formattedTimeRemaining,
    startTranspilation,
    updateProgress,
    stopTranspilation
  }
})
