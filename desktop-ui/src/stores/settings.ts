import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface AppSettings {
  theme: 'light' | 'dark' | 'system'
  fontSize: number
  autoSave: boolean
  defaultOutputDir: string
}

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<AppSettings>({
    theme: 'system',
    fontSize: 14,
    autoSave: true,
    defaultOutputDir: './output'
  })

  function updateSettings(newSettings: Partial<AppSettings>) {
    settings.value = { ...settings.value, ...newSettings }
  }

  return {
    settings,
    updateSettings
  }
})
