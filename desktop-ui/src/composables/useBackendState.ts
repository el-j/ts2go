/**
 * Composable for backend-persisted state management
 * Uses Tauri commands to delegate state persistence to the hexagonal backend
 */

import { invoke } from '@tauri-apps/api/core'

export interface BackendStateOptions {
  key: string
  defaultValue?: any
}

/**
 * Use backend-persisted state instead of localStorage
 * Delegates to Go backend's StateRepository through Tauri commands
 */
export function useBackendState<T>(key: string, defaultValue: T) {
  const load = async (): Promise<T> => {
    try {
      const jsonString = await invoke<string>('get_app_state', { key })
      if (!jsonString || jsonString === '{}') {
        return defaultValue
      }
      return JSON.parse(jsonString) as T
    } catch (error) {
      console.error(`Failed to load state for key "${key}":`, error)
      return defaultValue
    }
  }

  const save = async (data: T): Promise<void> => {
    try {
      const jsonString = JSON.stringify(data)
      await invoke('save_app_state', { key, data: jsonString })
    } catch (error) {
      console.error(`Failed to save state for key "${key}":`, error)
      throw error
    }
  }

  return {
    load,
    save
  }
}

/**
 * Use backend-persisted settings (singleton)
 * Delegates to Go backend's SettingsRepository through Tauri commands
 */
export function useBackendSettings<T>(defaultValue: T) {
  const load = async (): Promise<T> => {
    try {
      const jsonString = await invoke<string>('get_app_settings')
      if (!jsonString || jsonString === '{}') {
        return defaultValue
      }
      return JSON.parse(jsonString) as T
    } catch (error) {
      console.error('Failed to load settings:', error)
      return defaultValue
    }
  }

  const save = async (data: T): Promise<void> => {
    try {
      const jsonString = JSON.stringify(data)
      await invoke('save_app_settings', { data: jsonString })
    } catch (error) {
      console.error('Failed to save settings:', error)
      throw error
    }
  }

  return {
    load,
    save
  }
}
