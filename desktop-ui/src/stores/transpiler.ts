import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface TranspilationStatus {
  isRunning: boolean
  progress: number
  currentFile: string
  filesProcessed: number
  totalFiles: number
}

export const useTranspilerStore = defineStore('transpiler', () => {
  const status = ref<TranspilationStatus>({
    isRunning: false,
    progress: 0,
    currentFile: '',
    filesProcessed: 0,
    totalFiles: 0
  })

  function startTranspilation(totalFiles: number) {
    status.value = {
      isRunning: true,
      progress: 0,
      currentFile: '',
      filesProcessed: 0,
      totalFiles
    }
  }

  function updateProgress(filesProcessed: number, currentFile: string) {
    status.value.filesProcessed = filesProcessed
    status.value.currentFile = currentFile
    status.value.progress = (filesProcessed / status.value.totalFiles) * 100
  }

  function stopTranspilation() {
    status.value.isRunning = false
  }

  return {
    status,
    startTranspilation,
    updateProgress,
    stopTranspilation
  }
})
