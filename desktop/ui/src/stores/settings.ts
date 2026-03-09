import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useBackendSettings } from '../composables/useBackendState'

export interface AppSettings {
  // Application settings
  theme: 'light' | 'dark' | 'system'
  fontSize: number
  autoSave: boolean
  autoSaveDelay: number // in milliseconds
  enableBackups: boolean
  backupLocation: string
  defaultOutputDir: string
  
  // Project settings
  goModuleName: string
  excludePatterns: string
  includePatterns: string
  
  // Editor settings
  tabSize: number
  wordWrap: boolean
  lineNumbers: boolean
  minimap: boolean
  autoFormatOnSave: boolean
  
  // Go Configuration (Phase 2)
  goBinarySource: 'bundled' | 'system' | 'custom'
  customGoBinaryPath: string
}

const DEFAULT_SETTINGS: AppSettings = {
  // Application settings
  theme: 'system',
  fontSize: 14,
  autoSave: true,
  autoSaveDelay: 3000, // 3 seconds
  enableBackups: true,
  backupLocation: './.backups',
  defaultOutputDir: './output',
  
  // Project settings
  goModuleName: '',
  excludePatterns: 'node_modules, **/*.test.ts, dist',
  includePatterns: '**/*.ts, **/*.tsx',
  
  // Editor settings
  tabSize: 4,
  wordWrap: true,
  lineNumbers: true,
  minimap: true,
  autoFormatOnSave: true,
  
  // Go Configuration (Phase 2)
  goBinarySource: 'system',  // Default to system Go
  customGoBinaryPath: ''
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<AppSettings>({ ...DEFAULT_SETTINGS })
  const backend = useBackendSettings<AppSettings>(DEFAULT_SETTINGS)
  let isLoading = true

  // Load settings from backend on init
  const loadSettings = async () => {
    try {
      const stored = await backend.load()
      settings.value = { ...DEFAULT_SETTINGS, ...stored }
    } catch (e) {
      console.error('Failed to load settings from backend:', e)
      settings.value = { ...DEFAULT_SETTINGS }
    } finally {
      isLoading = false
    }
  }

  // Save settings to backend
  const saveSettings = async () => {
    if (isLoading) return // Don't save during initial load
    
    try {
      await backend.save(settings.value)
    } catch (e) {
      console.error('Failed to save settings to backend:', e)
    }
  }

  // Watch for changes and auto-save to backend
  watch(settings, () => {
    saveSettings()
  }, { deep: true })

  function updateSettings(newSettings: Partial<AppSettings>) {
    settings.value = { ...settings.value, ...newSettings }
  }

  function resetToDefaults() {
    settings.value = { ...DEFAULT_SETTINGS }
  }

  // Load settings on initialization
  loadSettings()

  return {
    settings,
    updateSettings,
    resetToDefaults
  }
})
