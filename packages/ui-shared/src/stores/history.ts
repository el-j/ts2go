import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useBackendState } from '../composables/useBackendState'

export interface BuildRecord {
  id: string
  timestamp: Date
  projectPath: string
  filesProcessed: number
  totalFiles: number
  duration: number
  status: 'success' | 'failed' | 'partial'
  errors: number
  warnings: number
}

export const useHistoryStore = defineStore('history', () => {
  const builds = ref<BuildRecord[]>([])
  const maxHistorySize = ref(50) // Keep last 50 builds
  const backend = useBackendState<BuildRecord[]>('build-history', [])
  let isLoading = true

  const successfulBuilds = computed(() => 
    builds.value.filter(b => b.status === 'success')
  )

  const failedBuilds = computed(() => 
    builds.value.filter(b => b.status === 'failed')
  )

  const successRate = computed(() => {
    if (builds.value.length === 0) return 0
    return (successfulBuilds.value.length / builds.value.length) * 100
  })

  const averageDuration = computed(() => {
    if (successfulBuilds.value.length === 0) return 0
    const total = successfulBuilds.value.reduce((sum, b) => sum + b.duration, 0)
    return total / successfulBuilds.value.length
  })

  function addBuild(build: Omit<BuildRecord, 'id' | 'timestamp'>) {
    const record: BuildRecord = {
      ...build,
      id: `build-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      timestamp: new Date()
    }

    builds.value.unshift(record) // Add to beginning

    // Keep only max history size
    if (builds.value.length > maxHistorySize.value) {
      builds.value = builds.value.slice(0, maxHistorySize.value)
    }

    // Persist to backend
    saveToBackend()
  }

  function removeBuild(buildId: string) {
    builds.value = builds.value.filter(b => b.id !== buildId)
    saveToBackend()
  }

  function clearHistory() {
    builds.value = []
    saveToBackend()
  }

  async function saveToBackend() {
    if (isLoading) return
    
    try {
      await backend.save(builds.value)
    } catch (e) {
      console.error('Failed to save build history to backend:', e)
    }
  }

  async function loadFromBackend() {
    try {
      const stored = await backend.load()
      // Convert timestamp strings back to Date objects
      builds.value = stored.map((b: any) => ({
        ...b,
        timestamp: new Date(b.timestamp)
      }))
    } catch (e) {
      console.error('Failed to load build history from backend:', e)
    } finally {
      isLoading = false
    }
  }

  // Load on init
  loadFromBackend()

  return {
    builds,
    maxHistorySize,
    successfulBuilds,
    failedBuilds,
    successRate,
    averageDuration,
    addBuild,
    removeBuild,
    clearHistory
  }
})
