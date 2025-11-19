import { ref, onUnmounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useSaveFile } from './useSaveFile'
import type { OpenFile } from '@/stores/workspace'

export function useAutoSave() {
  const settings = useSettingsStore()
  const { saveFile } = useSaveFile()
  
  const saveTimers = ref<Map<string, NodeJS.Timeout>>(new Map())

  /**
   * Schedule an auto-save for a file
   */
  function scheduleSave(file: OpenFile) {
    // Skip if auto-save is disabled
    if (!settings.settings.autoSave) return

    // Clear existing timer for this file
    const existingTimer = saveTimers.value.get(file.path)
    if (existingTimer) {
      clearTimeout(existingTimer)
    }

    // Schedule new save with configurable delay
    const delay = settings.settings.autoSaveDelay || 3000
    const timer = setTimeout(async () => {
      if (file.isDirty) {
        console.log(`Auto-saving: ${file.name}`)
        await saveFile(file, { showDialog: false })
      }
      saveTimers.value.delete(file.path)
    }, delay)

    saveTimers.value.set(file.path, timer)
  }

  /**
   * Cancel all pending auto-saves
   */
  function cancelAllSaves() {
    saveTimers.value.forEach(timer => clearTimeout(timer))
    saveTimers.value.clear()
  }

  /**
   * Cancel auto-save for a specific file
   */
  function cancelSave(filePath: string) {
    const timer = saveTimers.value.get(filePath)
    if (timer) {
      clearTimeout(timer)
      saveTimers.value.delete(filePath)
    }
  }

  // Cleanup on unmount
  onUnmounted(() => {
    cancelAllSaves()
  })

  return {
    scheduleSave,
    cancelSave,
    cancelAllSaves
  }
}
