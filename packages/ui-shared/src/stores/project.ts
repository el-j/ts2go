import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useBackendState } from '../composables/useBackendState'
import type { LocalProject as Project } from '../types'

export type { LocalProject as Project } from '../types'

const MAX_RECENT_PROJECTS = 20

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const recentProjects = ref<Project[]>([])
  const backend = useBackendState<Project[]>('recent-projects', [])
  let isLoading = true

  // Load recent projects from backend
  const loadRecentProjects = async () => {
    try {
      const stored = await backend.load()
      recentProjects.value = stored.map((p: any) => ({
        ...p,
        lastModified: new Date(p.lastModified),
        lastOpened: p.lastOpened
      }))
    } catch (e) {
      console.error('Failed to load recent projects from backend:', e)
    } finally {
      isLoading = false
    }
  }

  // Save recent projects to backend
  const saveRecentProjects = async () => {
    if (isLoading) return
    
    try {
      await backend.save(recentProjects.value)
    } catch (e) {
      console.error('Failed to save recent projects to backend:', e)
    }
  }

  // Watch for changes and auto-save to backend
  watch(recentProjects, () => {
    saveRecentProjects()
  }, { deep: true })

  function addProject(project: Project) {
    projects.value.push(project)
  }

  function setCurrentProject(project: Project | null) {
    currentProject.value = project
    if (project) {
      addToRecentProjects(project)
    }
  }

  function addToRecentProjects(project: Project) {
    const now = new Date().toISOString()
    
    // Find if project already exists in recent
    const existingIndex = recentProjects.value.findIndex(p => p.path === project.path)
    
    if (existingIndex >= 0) {
      // Update existing project
      const existing = recentProjects.value[existingIndex]
      existing.lastOpened = now
      existing.accessCount = (existing.accessCount || 0) + 1
      existing.name = project.name
      existing.lastModified = project.lastModified
      
      // Move to front if not pinned
      if (!existing.isPinned) {
        recentProjects.value.splice(existingIndex, 1)
        recentProjects.value.unshift(existing)
      }
    } else {
      // Add new project
      const newProject: Project = {
        ...project,
        lastOpened: now,
        isPinned: false,
        accessCount: 1
      }
      
      recentProjects.value.unshift(newProject)
      
      // Remove oldest unpinned projects if we exceed max
      if (recentProjects.value.length > MAX_RECENT_PROJECTS) {
        const unpinnedProjects = recentProjects.value.filter(p => !p.isPinned)
        const pinnedProjects = recentProjects.value.filter(p => p.isPinned)
        const trimmedUnpinned = unpinnedProjects.slice(0, MAX_RECENT_PROJECTS - pinnedProjects.length)
        recentProjects.value = [...pinnedProjects, ...trimmedUnpinned]
      }
    }
  }

  function togglePinProject(projectPath: string) {
    const project = recentProjects.value.find(p => p.path === projectPath)
    if (project) {
      project.isPinned = !project.isPinned
      
      // Re-sort: pinned first, then by last opened
      recentProjects.value.sort((a, b) => {
        if (a.isPinned && !b.isPinned) return -1
        if (!a.isPinned && b.isPinned) return 1
        return (b.lastOpened || '').localeCompare(a.lastOpened || '')
      })
    }
  }

  function removeFromRecentProjects(projectPath: string) {
    const index = recentProjects.value.findIndex(p => p.path === projectPath)
    if (index >= 0) {
      recentProjects.value.splice(index, 1)
    }
  }

  function clearRecentProjects() {
    recentProjects.value = []
  }

  // Alias for convenience (matches usage in ProjectLoader)
  function addRecentProject(projectData: { name: string; path: string; lastOpened: number }) {
    const project: Project = {
      id: projectData.path, // Use path as ID
      name: projectData.name,
      path: projectData.path,
      lastModified: new Date(projectData.lastOpened),
      lastOpened: new Date(projectData.lastOpened).toISOString(),
      isPinned: false,
      accessCount: 1
    }
    addToRecentProjects(project)
  }

  // Load recent projects on initialization
  loadRecentProjects()

  return {
    projects,
    currentProject,
    recentProjects,
    addProject,
    setCurrentProject,
    addToRecentProjects,
    addRecentProject,
    togglePinProject,
    removeFromRecentProjects,
    clearRecentProjects
  }
})
