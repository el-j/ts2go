import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTranspilerStore } from '../transpiler'

describe('Transpiler Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default status', () => {
    const store = useTranspilerStore()
    
    expect(store.status.isRunning).toBe(false)
    expect(store.status.progress).toBe(0)
    expect(store.status.currentFile).toBe('')
    expect(store.status.filesProcessed).toBe(0)
    expect(store.status.totalFiles).toBe(0)
  })

  it('should start transpilation', () => {
    const store = useTranspilerStore()
    
    store.startTranspilation(10)

    expect(store.status.isRunning).toBe(true)
    expect(store.status.totalFiles).toBe(10)
    expect(store.status.progress).toBe(0)
    expect(store.status.filesProcessed).toBe(0)
  })

  it('should update progress', () => {
    const store = useTranspilerStore()
    
    store.startTranspilation(10)
    store.updateProgress(3, 'test.ts')

    expect(store.status.filesProcessed).toBe(3)
    expect(store.status.currentFile).toBe('test.ts')
    expect(store.status.progress).toBe(30)
  })

  it('should calculate progress correctly', () => {
    const store = useTranspilerStore()
    
    store.startTranspilation(4)
    
    store.updateProgress(1, 'file1.ts')
    expect(store.status.progress).toBe(25)
    
    store.updateProgress(2, 'file2.ts')
    expect(store.status.progress).toBe(50)
    
    store.updateProgress(4, 'file4.ts')
    expect(store.status.progress).toBe(100)
  })

  it('should stop transpilation', () => {
    const store = useTranspilerStore()
    
    store.startTranspilation(10)
    expect(store.status.isRunning).toBe(true)
    
    store.stopTranspilation()
    expect(store.status.isRunning).toBe(false)
  })

  it('should maintain progress after stopping', () => {
    const store = useTranspilerStore()
    
    store.startTranspilation(10)
    store.updateProgress(5, 'test.ts')
    
    expect(store.status.filesProcessed).toBe(5)
    expect(store.status.progress).toBe(50)
    
    store.stopTranspilation()
    
    expect(store.status.filesProcessed).toBe(5)
    expect(store.status.progress).toBe(50)
    expect(store.status.isRunning).toBe(false)
  })
})
