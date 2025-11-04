import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSettingsStore } from '../settings'

describe('Settings Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    // Clear localStorage
    localStorage.clear()
  })

  it('should initialize with default settings', () => {
    const store = useSettingsStore()
    
    expect(store.settings.theme).toBe('system')
    expect(store.settings.fontSize).toBe(14)
    expect(store.settings.autoSave).toBe(true)
    expect(store.settings.defaultOutputDir).toBe('./output')
    expect(store.settings.tabSize).toBe(4)
    expect(store.settings.wordWrap).toBe(true)
    expect(store.settings.lineNumbers).toBe(true)
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

  it('should update editor settings', () => {
    const store = useSettingsStore()
    
    store.updateSettings({
      tabSize: 2,
      wordWrap: false,
      lineNumbers: false,
      minimap: false,
      autoFormatOnSave: false
    })
    
    expect(store.settings.tabSize).toBe(2)
    expect(store.settings.wordWrap).toBe(false)
    expect(store.settings.lineNumbers).toBe(false)
    expect(store.settings.minimap).toBe(false)
    expect(store.settings.autoFormatOnSave).toBe(false)
  })

  it('should update project settings', () => {
    const store = useSettingsStore()
    
    store.updateSettings({
      goModuleName: 'github.com/user/project',
      excludePatterns: 'node_modules, dist',
      includePatterns: '**/*.ts'
    })
    
    expect(store.settings.goModuleName).toBe('github.com/user/project')
    expect(store.settings.excludePatterns).toBe('node_modules, dist')
    expect(store.settings.includePatterns).toBe('**/*.ts')
  })

  it('should reset to defaults', () => {
    const store = useSettingsStore()
    
    store.updateSettings({
      theme: 'dark',
      fontSize: 20,
      autoSave: false,
      tabSize: 2
    })
    
    store.resetToDefaults()
    
    expect(store.settings.theme).toBe('system')
    expect(store.settings.fontSize).toBe(14)
    expect(store.settings.autoSave).toBe(true)
    expect(store.settings.tabSize).toBe(4)
  })
})
