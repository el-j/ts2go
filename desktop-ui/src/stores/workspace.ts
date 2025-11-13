import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface FileNode {
  name: string
  displayName?: string  // Relative path for display in UI
  path: string  // Absolute path for file operations
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
  
  function loadProject(projectInfo: any) {
    // Set project path
    setProjectPath(projectInfo.root)
    
    // Convert files array to file tree structure
    const tree: FileNode[] = []
    const fileMap = new Map<string, FileNode>()
    const root = projectInfo.root
    
    // Create nodes for all files
    projectInfo.files.forEach((absolutePath: string) => {
      // Get relative path from project root
      const relativePath = absolutePath.startsWith(root) 
        ? absolutePath.slice(root.length + 1)  // +1 to remove leading slash
        : absolutePath
      
      const parts = relativePath.split('/')
      const fileName = parts[parts.length - 1]
      
      const node: FileNode = {
        name: fileName,
        displayName: relativePath,  // Show relative path in UI
        path: absolutePath,  // Store absolute path for file operations
        type: 'file',
        content: '',
        isDirty: false
      }
      
      fileMap.set(relativePath, node)
    })
    
    // Build tree structure
    projectInfo.files.forEach((absolutePath: string) => {
      const relativePath = absolutePath.startsWith(root)
        ? absolutePath.slice(root.length + 1)
        : absolutePath
      
      const parts = relativePath.split('/')
      const node = fileMap.get(relativePath)!
      
      if (parts.length === 1) {
        // Root level file
        tree.push(node)
      } else {
        // Nested file - find or create parent folders
        let currentRelativePath = ''
        let currentLevel: FileNode[] = tree
        
        for (let i = 0; i < parts.length - 1; i++) {
          const part = parts[i]
          currentRelativePath = currentRelativePath ? `${currentRelativePath}/${part}` : part
          const currentAbsolutePath = `${root}/${currentRelativePath}`
          
          let folder = currentLevel.find(n => n.path === currentAbsolutePath && n.type === 'folder')
          if (!folder) {
            folder = {
              name: part,
              displayName: currentRelativePath,  // Show relative path for folders too
              path: currentAbsolutePath,  // Store absolute path for folders
              type: 'folder',
              children: []
            }
            currentLevel.push(folder)
          }
          
          currentLevel = folder.children!
        }
        
        currentLevel.push(node)
      }
    })
    
    setFileTree(tree)
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
    closeAllFiles,
    loadProject
  }
})
