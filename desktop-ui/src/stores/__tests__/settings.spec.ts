import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSettingsStore } from '../settings'

describe('Settings Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default settings', () => {
    const store = useSettingsStore()
    
    expect(store.settings.theme).toBe('system')
    expect(store.settings.fontSize).toBe(14)
    expect(store.settings.autoSave).toBe(true)
    expect(store.settings.defaultOutputDir).toBe('./output')
  })

  it('should update theme setting', () => {
    const store = useSettingsStore()
    
    store.updateSettings({ theme: 'dark' })
    
    expect(store.settings.theme).toBe('dark')
    expect(store.settings.fontSize).toBe(14)
  })

  it('should update fontSize setting', () => {
    const store = useSettingsStore()
    
    store.updateSettings({ fontSize: 16 })
    
    expect(store.settings.fontSize).toBe(16)
    expect(store.settings.theme).toBe('system')
  })

  it('should update multiple settings at once', () => {
    const store = useSettingsStore()
    
    store.updateSettings({
      theme: 'light',
      fontSize: 18,
      autoSave: false
    })
    
    expect(store.settings.theme).toBe('light')
    expect(store.settings.fontSize).toBe(18)
    expect(store.settings.autoSave).toBe(false)
    expect(store.settings.defaultOutputDir).toBe('./output')
  })

  it('should update defaultOutputDir', () => {
    const store = useSettingsStore()
    
    store.updateSettings({ defaultOutputDir: '/custom/path' })
    
    expect(store.settings.defaultOutputDir).toBe('/custom/path')
  })

  it('should preserve unchanged settings', () => {
    const store = useSettingsStore()
    
    store.updateSettings({ theme: 'dark' })
    expect(store.settings.fontSize).toBe(14)
    expect(store.settings.autoSave).toBe(true)
    
    store.updateSettings({ fontSize: 20 })
    expect(store.settings.theme).toBe('dark')
    expect(store.settings.autoSave).toBe(true)
  })
})
