import { ref } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { useWorkspaceStore } from '@/stores/workspace'
import { useSettingsStore } from '@/stores/settings'
import { useToast } from 'primevue/usetoast'
import type { OpenFile } from '@/stores/workspace'
import { usePlatform } from './usePlatform'

export interface SaveOptions {
  showDialog?: boolean
  createBackup?: boolean
  validateContent?: boolean
}

export function useSaveFile() {
  const workspace = useWorkspaceStore()
  const settings = useSettingsStore()
  const toast = useToast()
  const isSaving = ref(false)
  const { isWeb } = usePlatform()

  /**
   * Save a single file
   */
  async function saveFile(file: OpenFile, options: SaveOptions = {}): Promise<boolean> {
    if (isSaving.value) return false

    try {
      isSaving.value = true

      // Create backup if requested or enabled in settings
      if (options.createBackup || settings.settings.enableBackups) {
        await createBackup(file)
      }

      // Validate content if requested
      if (options.validateContent) {
        const isValid = await validateSyntax(file)
        if (!isValid && options.showDialog) {
          const proceed = confirm('File has syntax errors. Save anyway?')
          if (!proceed) return false
        }
      }

      if (isWeb) {
        // In web mode, upload to SaaS storage API if configured
        const apiBase = import.meta.env.VITE_SAAS_API_URL || ''
        const projectId = (import.meta.env.VITE_PROJECT_ID as string) || ''

        if (!apiBase || !projectId) {
          throw new Error('Web save requires VITE_SAAS_API_URL and VITE_PROJECT_ID')
        }

        const form = new FormData()
        form.append('project_id', projectId)
        const blob = new Blob([file.content], { type: 'text/plain' })
        form.append('files', new File([blob], file.name, { type: 'text/plain' }))

        const resp = await fetch(`${apiBase}/api/v1/files/upload`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${localStorage.getItem('ts2go:token') || ''}`
          },
          body: form
        })
        if (!resp.ok) {
          const text = await resp.text()
          throw new Error(`Upload failed: ${resp.status} ${text}`)
        }
      } else {
        // Save file to disk via Tauri
        await invoke('write_file', {
          path: file.path,
          content: file.content
        })
      }

      // Mark file as saved in workspace
      workspace.markFileSaved(file.path)

      if (options.showDialog) {
        toast.add({
          severity: 'success',
          summary: 'File Saved',
          detail: `${file.name} saved successfully`,
          life: 3000
        })
      }

      return true
    } catch (error) {
      console.error('Failed to save file:', error)
      
      toast.add({
        severity: 'error',
        summary: 'Save Failed',
        detail: `Failed to save ${file.name}: ${error}`,
        life: 5000
      })

      return false
    } finally {
      isSaving.value = false
    }
  }

  /**
   * Save all dirty files
   */
  async function saveAll(): Promise<{ saved: number; failed: number }> {
    const dirtyFiles = workspace.openFiles.filter(f => f.isDirty)
    
    if (dirtyFiles.length === 0) {
      toast.add({
        severity: 'info',
        summary: 'Nothing to Save',
        detail: 'No unsaved changes',
        life: 3000
      })
      return { saved: 0, failed: 0 }
    }

    let saved = 0
    let failed = 0

    for (const file of dirtyFiles) {
      const success = await saveFile(file, { showDialog: false })
      if (success) {
        saved++
      } else {
        failed++
      }
    }

    toast.add({
      severity: saved > 0 && failed === 0 ? 'success' : 'warn',
      summary: 'Save All Complete',
      detail: `Saved ${saved} file(s)${failed > 0 ? `, ${failed} failed` : ''}`,
      life: 4000
    })

    return { saved, failed }
  }

  /**
   * Create a backup of the file before saving
   */
  async function createBackup(file: OpenFile): Promise<void> {
    // Use backup location from settings
    const backupDir = settings.settings.backupLocation || './.backups'
    const fileName = file.name
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-')
    const backupPath = `${backupDir}/${fileName}.${timestamp}.backup`
    
    try {
      await invoke('write_file', {
        path: backupPath,
        content: file.content
      })
    } catch (error) {
      console.error('Failed to create backup:', error)
      // Don't throw - backup failure shouldn't prevent save
    }
  }

  /**
   * Validate file syntax (basic check)
   */
  async function validateSyntax(file: OpenFile): Promise<boolean> {
    const ext = file.name.split('.').pop()?.toLowerCase()
    
    // Only validate JSON files for now
    if (ext === 'json') {
      try {
        JSON.parse(file.content)
        return true
      } catch {
        return false
      }
    }

    // Note: Syntax validation deferred to transpiler - validation happens during transpilation
    // to avoid duplicate parsing overhead. Users get feedback via transpile errors.
    return true
  }

  /**
   * Check if there are unsaved changes
   */
  function hasUnsavedChanges(): boolean {
    return workspace.hasUnsavedChanges
  }

  /**
   * Get list of dirty files
   */
  function getDirtyFiles(): OpenFile[] {
    return workspace.openFiles.filter(f => f.isDirty)
  }

  return {
    isSaving,
    saveFile,
    saveAll,
    hasUnsavedChanges,
    getDirtyFiles
  }
}
