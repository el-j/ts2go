import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

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

  // Load settings from localStorage on init
  const loadSettings = () => {
    const stored = localStorage.getItem('ts2go-settings')
    if (stored) {
      try {
        const parsed = JSON.parse(stored)
        settings.value = { ...DEFAULT_SETTINGS, ...parsed }
      } catch (e) {
        console.error('Failed to load settings:', e)
      }
    }
  }

  // Save settings to localStorage
  const saveSettings = () => {
    try {
      localStorage.setItem('ts2go-settings', JSON.stringify(settings.value))
    } catch (e) {
      console.error('Failed to save settings:', e)
    }
  }

  // Watch for changes and auto-save
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
