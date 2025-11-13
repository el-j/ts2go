import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Artifact {
  id: string
  type: 'binary' | 'test_results'
  path: string
  size: number
  created_at: number
  project_path: string
  metadata: {
    go_version?: string
    platform?: string
    architecture?: string
    test_summary?: {
      passed: number
      failed: number
      skipped: number
      total: number
    }
  }
}

export const useArtifactsStore = defineStore('artifacts', () => {
  // State
  const artifacts = ref<Artifact[]>([])
  
  // Load from localStorage on init
  const loadArtifacts = () => {
    const stored = localStorage.getItem('ts2go_artifacts')
    if (stored) {
      try {
        artifacts.value = JSON.parse(stored)
      } catch (e) {
        console.error('Failed to load artifacts:', e)
        artifacts.value = []
      }
    }
  }
  
  // Save to localStorage
  const saveArtifacts = () => {
    try {
      localStorage.setItem('ts2go_artifacts', JSON.stringify(artifacts.value))
    } catch (e) {
      console.error('Failed to save artifacts:', e)
    }
  }
  
  // Actions
  function addArtifact(artifact: Artifact) {
    artifacts.value.unshift(artifact)
    
    // Keep only last 100 artifacts
    if (artifacts.value.length > 100) {
      artifacts.value = artifacts.value.slice(0, 100)
    }
    
    saveArtifacts()
  }
  
  function getProjectArtifacts(projectPath: string): Artifact[] {
    return artifacts.value.filter(a => a.project_path === projectPath)
  }
  
  function deleteArtifact(id: string) {
    artifacts.value = artifacts.value.filter(a => a.id !== id)
    saveArtifacts()
  }
  
  function clearOldArtifacts(daysOld: number = 30) {
    const cutoffTime = Date.now() - (daysOld * 24 * 60 * 60 * 1000)
    artifacts.value = artifacts.value.filter(a => a.created_at > cutoffTime)
    saveArtifacts()
  }
  
  function clearAllArtifacts() {
    artifacts.value = []
    saveArtifacts()
  }
  
  // Initialize
  loadArtifacts()
  
  return {
    // State
    artifacts,
    
    // Actions
    addArtifact,
    getProjectArtifacts,
    deleteArtifact,
    clearOldArtifacts,
    clearAllArtifacts
  }
})
