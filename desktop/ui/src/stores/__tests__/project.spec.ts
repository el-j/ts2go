import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useProjectStore, type Project } from '../project'

describe('Project Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with empty state', () => {
    const store = useProjectStore()
    
    expect(store.projects).toEqual([])
    expect(store.currentProject).toBeNull()
  })

  it('should add a project', () => {
    const store = useProjectStore()
    const project: Project = {
      id: '1',
      name: 'Test Project',
      path: '/path/to/project',
      lastModified: new Date()
    }

    store.addProject(project)

    expect(store.projects).toHaveLength(1)
    expect(store.projects[0]).toEqual(project)
  })

  it('should set current project', () => {
    const store = useProjectStore()
    const project: Project = {
      id: '1',
      name: 'Test Project',
      path: '/path/to/project',
      lastModified: new Date()
    }

    store.setCurrentProject(project)

    expect(store.currentProject).toEqual(project)
  })

  it('should clear current project', () => {
    const store = useProjectStore()
    const project: Project = {
      id: '1',
      name: 'Test Project',
      path: '/path/to/project',
      lastModified: new Date()
    }

    store.setCurrentProject(project)
    expect(store.currentProject).toEqual(project)

    store.setCurrentProject(null)
    expect(store.currentProject).toBeNull()
  })

  it('should handle multiple projects', () => {
    const store = useProjectStore()
    const project1: Project = {
      id: '1',
      name: 'Project 1',
      path: '/path/to/project1',
      lastModified: new Date()
    }
    const project2: Project = {
      id: '2',
      name: 'Project 2',
      path: '/path/to/project2',
      lastModified: new Date()
    }

    store.addProject(project1)
    store.addProject(project2)

    expect(store.projects).toHaveLength(2)
    expect(store.projects[0]).toEqual(project1)
    expect(store.projects[1]).toEqual(project2)
  })
})
