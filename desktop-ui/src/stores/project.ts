import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export interface Project {
  id: string
  name: string
  path: string
  lastModified: Date
  lastOpened?: string
  isPinned?: boolean
  accessCount?: number
}

const MAX_RECENT_PROJECTS = 20
const STORAGE_KEY = 'ts2go-recent-projects'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const recentProjects = ref<Project[]>([])

  // Load recent projects from localStorage
  const loadRecentProjects = () => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        const parsed = JSON.parse(stored)
        recentProjects.value = parsed.map((p: any) => ({
          ...p,
          lastModified: new Date(p.lastModified),
          lastOpened: p.lastOpened
        }))
      }
    } catch (e) {
      console.error('Failed to load recent projects:', e)
    }
  }

  // Save recent projects to localStorage
  const saveRecentProjects = () => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(recentProjects.value))
    } catch (e) {
      console.error('Failed to save recent projects:', e)
    }
  }

  // Watch for changes and auto-save
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

  // Load recent projects on initialization
  loadRecentProjects()

  return {
    projects,
    currentProject,
    recentProjects,
    addProject,
    setCurrentProject,
    addToRecentProjects,
    togglePinProject,
    removeFromRecentProjects,
    clearRecentProjects
  }
})
