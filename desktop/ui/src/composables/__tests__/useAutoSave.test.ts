import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useAutoSave } from '../useAutoSave'
import { useSettingsStore } from '@/stores/settings'
import { createPinia, setActivePinia } from 'pinia'
import type { OpenFile } from '@/stores/workspace'

// Mock useSaveFile
vi.mock('../useSaveFile', () => ({
  useSaveFile: () => ({
    saveFile: vi.fn().mockResolvedValue(true)
  })
}))

describe('useAutoSave', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.clearAllMocks()
    vi.useRealTimers()
  })

  it('should schedule auto-save for a dirty file', async () => {
    const { scheduleSave } = useAutoSave()
    const settings = useSettingsStore()
    const { useSaveFile } = await import('../useSaveFile')
    
    settings.settings.autoSave = true
    settings.settings.autoSaveDelay = 3000

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    scheduleSave(mockFile)

    // Fast-forward time
    vi.advanceTimersByTime(3000)
    await vi.runAllTimersAsync()

    const { saveFile } = useSaveFile()
    expect(saveFile).toHaveBeenCalledWith(mockFile, { showDialog: false })
  })

  it('should respect auto-save delay setting', async () => {
    const { scheduleSave } = useAutoSave()
    const settings = useSettingsStore()
    const { useSaveFile } = await import('../useSaveFile')
    
    settings.settings.autoSave = true
    settings.settings.autoSaveDelay = 5000 // 5 seconds

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    scheduleSave(mockFile)

    // Should not save after 3 seconds
    vi.advanceTimersByTime(3000)
    await vi.runAllTimersAsync()

    const { saveFile } = useSaveFile()
    expect(saveFile).not.toHaveBeenCalled()

    // Should save after 5 seconds total
    vi.advanceTimersByTime(2000)
    await vi.runAllTimersAsync()

    expect(saveFile).toHaveBeenCalled()
  })

  it('should not schedule save when auto-save is disabled', () => {
    const { scheduleSave } = useAutoSave()
    const settings = useSettingsStore()
    
    settings.settings.autoSave = false

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    scheduleSave(mockFile)
    vi.advanceTimersByTime(3000)

    // No timers should be scheduled
    expect(vi.getTimerCount()).toBe(0)
  })

  it('should cancel existing timer when scheduling new save', async () => {
    const { scheduleSave } = useAutoSave()
    const settings = useSettingsStore()
    const { useSaveFile } = await import('../useSaveFile')
    
    settings.settings.autoSave = true
    settings.settings.autoSaveDelay = 3000

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    // Schedule first save
    scheduleSave(mockFile)
    vi.advanceTimersByTime(1000)

    // Schedule again (should cancel first)
    scheduleSave(mockFile)
    vi.advanceTimersByTime(3000)
    await vi.runAllTimersAsync()

    const { saveFile } = useSaveFile()
    // Should only be called once
    expect(saveFile).toHaveBeenCalledTimes(1)
  })

  it('should cancel all pending saves', () => {
    const { scheduleSave, cancelAllSaves } = useAutoSave()
    const settings = useSettingsStore()
    
    settings.settings.autoSave = true
    settings.settings.autoSaveDelay = 3000

    const file1: OpenFile = {
      path: '/test/file1.ts',
      name: 'file1.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    const file2: OpenFile = {
      path: '/test/file2.ts',
      name: 'file2.ts',
      content: 'const y = 2;',
      isDirty: true
    }

    scheduleSave(file1)
    scheduleSave(file2)

    expect(vi.getTimerCount()).toBe(2)

    cancelAllSaves()

    expect(vi.getTimerCount()).toBe(0)
  })

  it('should cancel save for specific file', () => {
    const { scheduleSave, cancelSave } = useAutoSave()
    const settings = useSettingsStore()
    
    settings.settings.autoSave = true
    settings.settings.autoSaveDelay = 3000

    const file1: OpenFile = {
      path: '/test/file1.ts',
      name: 'file1.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    const file2: OpenFile = {
      path: '/test/file2.ts',
      name: 'file2.ts',
      content: 'const y = 2;',
      isDirty: true
    }

    scheduleSave(file1)
    scheduleSave(file2)

    cancelSave('/test/file1.ts')

    expect(vi.getTimerCount()).toBe(1)
  })

  it('should not save file that is no longer dirty', async () => {
    const { scheduleSave } = useAutoSave()
    const settings = useSettingsStore()
    const { useSaveFile } = await import('../useSaveFile')
    
    settings.settings.autoSave = true
    settings.settings.autoSaveDelay = 3000

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    scheduleSave(mockFile)

    // Mark file as clean before timer fires
    mockFile.isDirty = false

    vi.advanceTimersByTime(3000)
    await vi.runAllTimersAsync()

    const { saveFile } = useSaveFile()
    expect(saveFile).not.toHaveBeenCalled()
  })
})
