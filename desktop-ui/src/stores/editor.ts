import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface EditorTab {
  id: string
  name: string
  path: string
  content: string
  language: 'typescript' | 'go'
  isDirty: boolean
  isActive: boolean
}

export const useEditorStore = defineStore('editor', () => {
  const tabs = ref<EditorTab[]>([])
  const activeTabId = ref<string | null>(null)

  const activeTab = computed(() => {
    return tabs.value.find(tab => tab.id === activeTabId.value) || null
  })

  const dirtyTabs = computed(() => {
    return tabs.value.filter(tab => tab.isDirty)
  })

  function addTab(tab: Omit<EditorTab, 'id' | 'isDirty' | 'isActive'>) {
    // Check if tab with same path already exists
    const existingTab = tabs.value.find(t => t.path === tab.path)
    if (existingTab) {
      setActiveTab(existingTab.id)
      return existingTab
    }

    const newTab: EditorTab = {
      ...tab,
      id: `tab-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      isDirty: false,
      isActive: false
    }

    tabs.value.push(newTab)
    setActiveTab(newTab.id)
    return newTab
  }

  function removeTab(tabId: string) {
    const index = tabs.value.findIndex(tab => tab.id === tabId)
    if (index === -1) return

    tabs.value.splice(index, 1)

    // If removing active tab, activate another
    if (activeTabId.value === tabId) {
      if (tabs.value.length > 0) {
        const newIndex = Math.min(index, tabs.value.length - 1)
        setActiveTab(tabs.value[newIndex].id)
      } else {
        activeTabId.value = null
      }
    }
  }

  function setActiveTab(tabId: string) {
    tabs.value.forEach(tab => {
      tab.isActive = tab.id === tabId
    })
    activeTabId.value = tabId
  }

  function updateTabContent(tabId: string, content: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.content = content
      tab.isDirty = true
    }
  }

  function markTabAsSaved(tabId: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.isDirty = false
    }
  }

  function closeAllTabs() {
    tabs.value = []
    activeTabId.value = null
  }

  function closeOtherTabs(tabId: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tabs.value = [tab]
      setActiveTab(tabId)
    }
  }

  return {
    tabs,
    activeTabId,
    activeTab,
    dirtyTabs,
    addTab,
    removeTab,
    setActiveTab,
    updateTabContent,
    markTabAsSaved,
    closeAllTabs,
    closeOtherTabs
  }
})
