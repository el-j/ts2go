import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useHistoryStore } from '../history'

describe('History Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    // Clear localStorage before each test
    localStorage.clear()
  })

  it('should initialize with empty builds', () => {
    const store = useHistoryStore()
    
    expect(store.builds).toEqual([])
    expect(store.successfulBuilds).toEqual([])
    expect(store.failedBuilds).toEqual([])
    expect(store.successRate).toBe(0)
    expect(store.averageDuration).toBe(0)
  })

  it('should add a build', () => {
    const store = useHistoryStore()
    
    store.addBuild({
      projectPath: '/test/project',
      filesProcessed: 5,
      totalFiles: 5,
      duration: 1000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    expect(store.builds).toHaveLength(1)
    expect(store.builds[0].projectPath).toBe('/test/project')
    expect(store.builds[0].status).toBe('success')
    expect(store.builds[0]).toHaveProperty('id')
    expect(store.builds[0]).toHaveProperty('timestamp')
  })

  it('should calculate success rate correctly', () => {
    const store = useHistoryStore()
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 100,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 0,
      totalFiles: 1,
      duration: 100,
      status: 'failed',
      errors: 1,
      warnings: 0
    })
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 100,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    expect(store.builds).toHaveLength(3)
    expect(store.successfulBuilds).toHaveLength(2)
    expect(store.failedBuilds).toHaveLength(1)
    expect(store.successRate).toBeCloseTo(66.67, 1)
  })

  it('should calculate average duration correctly', () => {
    const store = useHistoryStore()
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 1000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 2000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 3000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    expect(store.averageDuration).toBe(2000)
  })

  it('should not include failed builds in average duration', () => {
    const store = useHistoryStore()
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 1000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 0,
      totalFiles: 1,
      duration: 5000,
      status: 'failed',
      errors: 1,
      warnings: 0
    })
    
    expect(store.averageDuration).toBe(1000)
  })

  it('should remove a build', () => {
    const store = useHistoryStore()
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 1000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    const buildId = store.builds[0].id
    
    store.removeBuild(buildId)
    
    expect(store.builds).toHaveLength(0)
  })

  it('should clear all history', () => {
    const store = useHistoryStore()
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 1000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    store.addBuild({
      projectPath: '',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 2000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    expect(store.builds).toHaveLength(2)
    
    store.clearHistory()
    
    expect(store.builds).toHaveLength(0)
  })

  it('should limit history to maxHistorySize', () => {
    const store = useHistoryStore()
    store.maxHistorySize = 3
    
    for (let i = 0; i < 5; i++) {
      store.addBuild({
        projectPath: '',
        filesProcessed: 1,
        totalFiles: 1,
        duration: 1000,
        status: 'success',
        errors: 0,
        warnings: 0
      })
    }
    
    expect(store.builds).toHaveLength(3)
  })

  it('should persist to localStorage', () => {
    const store = useHistoryStore()
    
    store.addBuild({
      projectPath: '/test',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 1000,
      status: 'success',
      errors: 0,
      warnings: 0
    })
    
    const stored = localStorage.getItem('ts2go-build-history')
    expect(stored).not.toBeNull()
    
    const parsed = JSON.parse(stored!)
    expect(parsed).toHaveLength(1)
    expect(parsed[0].projectPath).toBe('/test')
  })

  it('should load from localStorage on init', () => {
    // Set up data in localStorage
    const mockBuild = {
      id: 'test-id',
      timestamp: new Date().toISOString(),
      projectPath: '/test',
      filesProcessed: 1,
      totalFiles: 1,
      duration: 1000,
      status: 'success',
      errors: 0,
      warnings: 0
    }
    
    localStorage.setItem('ts2go-build-history', JSON.stringify([mockBuild]))
    
    // Create new store instance
    setActivePinia(createPinia())
    const store = useHistoryStore()
    
    expect(store.builds).toHaveLength(1)
    expect(store.builds[0].projectPath).toBe('/test')
    expect(store.builds[0].timestamp).toBeInstanceOf(Date)
  })
})
