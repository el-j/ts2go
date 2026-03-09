import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useSaveFile } from '../useSaveFile'
import { useWorkspaceStore } from '@/stores/workspace'
import { useSettingsStore } from '@/stores/settings'
import { createPinia, setActivePinia } from 'pinia'
import type { OpenFile } from '@/stores/workspace'

// Mock Tauri API
vi.mock('@tauri-apps/api/core', () => ({
  invoke: vi.fn()
}))

// Mock PrimeVue toast
vi.mock('primevue/usetoast', () => ({
  useToast: () => ({
    add: vi.fn()
  })
}))

describe('useSaveFile', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('should save a single file successfully', async () => {
    const { saveFile } = useSaveFile()
    const workspace = useWorkspaceStore()
    const { invoke } = await import('@tauri-apps/api/core')

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    workspace.openFiles.push(mockFile)
    vi.mocked(invoke).mockResolvedValue(undefined)

    const result = await saveFile(mockFile)

    expect(result).toBe(true)
    expect(invoke).toHaveBeenCalledWith('write_file', {
      path: '/test/file.ts',
      content: 'const x = 1;'
    })
  })

  it('should create backups when enabled in settings', async () => {
    const { saveFile } = useSaveFile()
    const settings = useSettingsStore()
    const { invoke } = await import('@tauri-apps/api/core')

    settings.settings.enableBackups = true
    settings.settings.backupLocation = './.backups'

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    vi.mocked(invoke).mockResolvedValue(undefined)

    await saveFile(mockFile)

    // Should be called twice: once for backup, once for actual save
    expect(invoke).toHaveBeenCalledTimes(2)
    
    // First call is backup
    const backupCall = vi.mocked(invoke).mock.calls[0]
    expect(backupCall[0]).toBe('write_file')
    const backupArgs = backupCall[1] as Record<string, string>
    expect(backupArgs.path).toContain('.backups/file.ts')
    expect(backupArgs.path).toContain('.backup')
  })

  it('should save all dirty files', async () => {
    const { saveAll } = useSaveFile()
    const workspace = useWorkspaceStore()
    const { invoke } = await import('@tauri-apps/api/core')

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

    workspace.openFiles.push(file1, file2)
    vi.mocked(invoke).mockResolvedValue(undefined)

    const result = await saveAll()

    expect(result.saved).toBe(2)
    expect(result.failed).toBe(0)
    expect(invoke).toHaveBeenCalledTimes(2)
  })

  it('should handle save failures gracefully', async () => {
    const { saveFile } = useSaveFile()
    const { invoke } = await import('@tauri-apps/api/core')

    const mockFile: OpenFile = {
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    vi.mocked(invoke).mockRejectedValue(new Error('Write failed'))

    const result = await saveFile(mockFile)

    expect(result).toBe(false)
  })

  it('should detect unsaved changes', () => {
    const { hasUnsavedChanges } = useSaveFile()
    const workspace = useWorkspaceStore()

    expect(hasUnsavedChanges()).toBe(false)

    workspace.openFiles.push({
      path: '/test/file.ts',
      name: 'file.ts',
      content: 'const x = 1;',
      isDirty: true
    })

    expect(hasUnsavedChanges()).toBe(true)
  })

  it('should return dirty files list', () => {
    const { getDirtyFiles } = useSaveFile()
    const workspace = useWorkspaceStore()

    const dirtyFile: OpenFile = {
      path: '/test/dirty.ts',
      name: 'dirty.ts',
      content: 'const x = 1;',
      isDirty: true
    }

    const cleanFile: OpenFile = {
      path: '/test/clean.ts',
      name: 'clean.ts',
      content: 'const y = 2;',
      isDirty: false
    }

    workspace.openFiles.push(dirtyFile, cleanFile)

    const dirtyFiles = getDirtyFiles()
    expect(dirtyFiles).toHaveLength(1)
    expect(dirtyFiles[0].path).toBe('/test/dirty.ts')
  })

  it('should validate JSON syntax', async () => {
    const { saveFile } = useSaveFile()
    const { invoke } = await import('@tauri-apps/api/core')

    const validJsonFile: OpenFile = {
      path: '/test/data.json',
      name: 'data.json',
      content: '{"valid": true}',
      isDirty: true
    }

    vi.mocked(invoke).mockResolvedValue(undefined)

    const result = await saveFile(validJsonFile, { validateContent: true })
    expect(result).toBe(true)
  })
})
