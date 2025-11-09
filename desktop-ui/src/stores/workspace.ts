import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface FileNode {
  name: string
  path: string
  type: 'file' | 'folder'
  children?: FileNode[]
  content?: string
  isDirty?: boolean
}

export interface OpenFile {
  path: string
  name: string
  content: string
  isDirty: boolean
}

export const useWorkspaceStore = defineStore('workspace', () => {
  // State
  const projectPath = ref<string>('')
  const fileTree = ref<FileNode[]>([])
  const openFiles = ref<OpenFile[]>([])
  const activeFilePath = ref<string>('')
  
  // Computed
  const activeFile = computed(() => {
    return openFiles.value.find(f => f.path === activeFilePath.value)
  })
  
  const hasUnsavedChanges = computed(() => {
    return openFiles.value.some(f => f.isDirty)
  })
  
  // Actions
  function setProjectPath(path: string) {
    projectPath.value = path
  }
  
  function setFileTree(tree: FileNode[]) {
    fileTree.value = tree
  }
  
  function openFile(filePath: string, fileName: string, content: string = '') {
    const existing = openFiles.value.find(f => f.path === filePath)
    if (!existing) {
      openFiles.value.push({
        path: filePath,
        name: fileName,
        content,
        isDirty: false
      })
    }
    activeFilePath.value = filePath
  }
  
  function closeFile(filePath: string) {
    const index = openFiles.value.findIndex(f => f.path === filePath)
    if (index !== -1) {
      openFiles.value.splice(index, 1)
      
      // If closing active file, switch to another
      if (activeFilePath.value === filePath) {
        if (openFiles.value.length > 0) {
          activeFilePath.value = openFiles.value[Math.max(0, index - 1)].path
        } else {
          activeFilePath.value = ''
        }
      }
    }
  }
  
  function setActiveFile(filePath: string) {
    activeFilePath.value = filePath
  }
  
  function updateFileContent(filePath: string, content: string) {
    const file = openFiles.value.find(f => f.path === filePath)
    if (file) {
      file.content = content
      file.isDirty = true
    }
  }
  
  function markFileSaved(filePath: string) {
    const file = openFiles.value.find(f => f.path === filePath)
    if (file) {
      file.isDirty = false
    }
  }
  
  function createFile(parentPath: string, fileName: string) {
    const newPath = parentPath ? `${parentPath}/${fileName}` : fileName
    
    // Add to file tree
    const newNode: FileNode = {
      name: fileName,
      path: newPath,
      type: 'file',
      content: '',
      isDirty: false
    }
    
    if (parentPath) {
      addNodeToTree(fileTree.value, parentPath, newNode)
    } else {
      fileTree.value.push(newNode)
    }
    
    // Open the new file
    openFile(newPath, fileName, '')
  }
  
  function createFolder(parentPath: string, folderName: string) {
    const newPath = parentPath ? `${parentPath}/${folderName}` : folderName
    
    const newNode: FileNode = {
      name: folderName,
      path: newPath,
      type: 'folder',
      children: []
    }
    
    if (parentPath) {
      addNodeToTree(fileTree.value, parentPath, newNode)
    } else {
      fileTree.value.push(newNode)
    }
  }
  
  function renameFile(oldPath: string, newName: string) {
    const pathParts = oldPath.split('/')
    pathParts[pathParts.length - 1] = newName
    const newPath = pathParts.join('/')
    
    // Update in tree
    renameNodeInTree(fileTree.value, oldPath, newPath, newName)
    
    // Update open files
    const file = openFiles.value.find(f => f.path === oldPath)
    if (file) {
      file.path = newPath
      file.name = newName
    }
    
    if (activeFilePath.value === oldPath) {
      activeFilePath.value = newPath
    }
  }
  
  function deleteFile(filePath: string) {
    // Remove from tree
    removeNodeFromTree(fileTree.value, filePath)
    
    // Close if open
    closeFile(filePath)
  }
  
  // Helper functions
  function addNodeToTree(nodes: FileNode[], parentPath: string, newNode: FileNode): boolean {
    for (const node of nodes) {
      if (node.path === parentPath && node.type === 'folder') {
        if (!node.children) node.children = []
        node.children.push(newNode)
        return true
      }
      if (node.children && addNodeToTree(node.children, parentPath, newNode)) {
        return true
      }
    }
    return false
  }
  
  function renameNodeInTree(nodes: FileNode[], oldPath: string, newPath: string, newName: string): boolean {
    for (const node of nodes) {
      if (node.path === oldPath) {
        node.path = newPath
        node.name = newName
        return true
      }
      if (node.children && renameNodeInTree(node.children, oldPath, newPath, newName)) {
        return true
      }
    }
    return false
  }
  
  function removeNodeFromTree(nodes: FileNode[], targetPath: string): boolean {
    for (let i = 0; i < nodes.length; i++) {
      if (nodes[i].path === targetPath) {
        nodes.splice(i, 1)
        return true
      }
      if (nodes[i].children && removeNodeFromTree(nodes[i].children!, targetPath)) {
        return true
      }
    }
    return false
  }
  
  function closeAllFiles() {
    openFiles.value = []
    activeFilePath.value = ''
  }
  
  return {
    // State
    projectPath,
    fileTree,
    openFiles,
    activeFilePath,
    
    // Computed
    activeFile,
    hasUnsavedChanges,
    
    // Actions
    setProjectPath,
    setFileTree,
    openFile,
    closeFile,
    setActiveFile,
    updateFileContent,
    markFileSaved,
    createFile,
    createFolder,
    renameFile,
    deleteFile,
    closeAllFiles
  }
})
