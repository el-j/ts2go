import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useEditorStore } from '../editor'

describe('Editor Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with empty tabs', () => {
    const store = useEditorStore()
    
    expect(store.tabs).toEqual([])
    expect(store.activeTabId).toBeNull()
    expect(store.activeTab).toBeNull()
  })

  it('should add a new tab', () => {
    const store = useEditorStore()
    
    const tab = store.addTab({
      name: 'test.ts',
      path: '/path/to/test.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    expect(store.tabs).toHaveLength(1)
    expect(tab.name).toBe('test.ts')
    expect(tab.isDirty).toBe(false)
    expect(store.activeTabId).toBe(tab.id)
  })

  it('should not add duplicate tabs', () => {
    const store = useEditorStore()
    
    const tab1 = store.addTab({
      name: 'test.ts',
      path: '/path/to/test.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    const tab2 = store.addTab({
      name: 'test.ts',
      path: '/path/to/test.ts',
      content: 'const y = 2;',
      language: 'typescript'
    })
    
    expect(store.tabs).toHaveLength(1)
    expect(tab1.id).toBe(tab2.id)
    expect(store.activeTabId).toBe(tab1.id)
  })

  it('should remove a tab', () => {
    const store = useEditorStore()
    
    const tab = store.addTab({
      name: 'test.ts',
      path: '/path/to/test.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    store.removeTab(tab.id)
    
    expect(store.tabs).toHaveLength(0)
    expect(store.activeTabId).toBeNull()
  })

  it('should switch active tab when removing current active', () => {
    const store = useEditorStore()
    
    const tab1 = store.addTab({
      name: 'test1.ts',
      path: '/path/to/test1.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    const tab2 = store.addTab({
      name: 'test2.ts',
      path: '/path/to/test2.ts',
      content: 'const y = 2;',
      language: 'typescript'
    })
    
    expect(store.activeTabId).toBe(tab2.id)
    
    store.removeTab(tab2.id)
    
    expect(store.tabs).toHaveLength(1)
    expect(store.activeTabId).toBe(tab1.id)
  })

  it('should set active tab', () => {
    const store = useEditorStore()
    
    const tab1 = store.addTab({
      name: 'test1.ts',
      path: '/path/to/test1.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    const tab2 = store.addTab({
      name: 'test2.ts',
      path: '/path/to/test2.ts',
      content: 'const y = 2;',
      language: 'typescript'
    })
    
    store.setActiveTab(tab1.id)
    
    expect(store.activeTabId).toBe(tab1.id)
    expect(store.activeTab?.id).toBe(tab1.id)
    expect(tab1.isActive).toBe(true)
    expect(tab2.isActive).toBe(false)
  })

  it('should update tab content and mark as dirty', () => {
    const store = useEditorStore()
    
    const tab = store.addTab({
      name: 'test.ts',
      path: '/path/to/test.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    expect(tab.isDirty).toBe(false)
    
    store.updateTabContent(tab.id, 'const x = 2;')
    
    expect(tab.content).toBe('const x = 2;')
    expect(tab.isDirty).toBe(true)
  })

  it('should mark tab as saved', () => {
    const store = useEditorStore()
    
    const tab = store.addTab({
      name: 'test.ts',
      path: '/path/to/test.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    store.updateTabContent(tab.id, 'const x = 2;')
    expect(tab.isDirty).toBe(true)
    
    store.markTabAsSaved(tab.id)
    expect(tab.isDirty).toBe(false)
  })

  it('should close all tabs', () => {
    const store = useEditorStore()
    
    store.addTab({
      name: 'test1.ts',
      path: '/path/to/test1.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    store.addTab({
      name: 'test2.ts',
      path: '/path/to/test2.ts',
      content: 'const y = 2;',
      language: 'typescript'
    })
    
    expect(store.tabs).toHaveLength(2)
    
    store.closeAllTabs()
    
    expect(store.tabs).toHaveLength(0)
    expect(store.activeTabId).toBeNull()
  })

  it('should close other tabs', () => {
    const store = useEditorStore()
    
    const tab1 = store.addTab({
      name: 'test1.ts',
      path: '/path/to/test1.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    store.addTab({
      name: 'test2.ts',
      path: '/path/to/test2.ts',
      content: 'const y = 2;',
      language: 'typescript'
    })
    
    store.addTab({
      name: 'test3.ts',
      path: '/path/to/test3.ts',
      content: 'const z = 3;',
      language: 'typescript'
    })
    
    expect(store.tabs).toHaveLength(3)
    
    store.closeOtherTabs(tab1.id)
    
    expect(store.tabs).toHaveLength(1)
    expect(store.tabs[0].id).toBe(tab1.id)
    expect(store.activeTabId).toBe(tab1.id)
  })

  it('should track dirty tabs', () => {
    const store = useEditorStore()
    
    const tab1 = store.addTab({
      name: 'test1.ts',
      path: '/path/to/test1.ts',
      content: 'const x = 1;',
      language: 'typescript'
    })
    
    const tab2 = store.addTab({
      name: 'test2.ts',
      path: '/path/to/test2.ts',
      content: 'const y = 2;',
      language: 'typescript'
    })
    
    expect(store.dirtyTabs).toHaveLength(0)
    
    store.updateTabContent(tab1.id, 'const x = 2;')
    expect(store.dirtyTabs).toHaveLength(1)
    
    store.updateTabContent(tab2.id, 'const y = 3;')
    expect(store.dirtyTabs).toHaveLength(2)
    
    store.markTabAsSaved(tab1.id)
    expect(store.dirtyTabs).toHaveLength(1)
  })
})
