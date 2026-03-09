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
  isTranspiled?: boolean  // Whether this TS file has been transpiled
  transpilationStatus?: TranspilationStatus  // Current transpilation status
}

export interface OpenFile {
  path: string
  name: string
  content: string
  isDirty: boolean
}

export type TranspilationStatus = 'pending' | 'transpiling' | 'success' | 'error' | 'warning'

export interface TranspilationInfo {
  tsFilePath: string
  goFilePath: string
  goCode?: string
  timestamp: number
  success: boolean
  status: TranspilationStatus
  errorMessage?: string
  warningMessage?: string
  logs?: string[]
}

export interface TranspilationProgress {
  currentFile: string | null
  totalFiles: number
  completedFiles: number
  logs: Array<{ timestamp: number; message: string; level: 'info' | 'error' | 'warning' | 'success' }>
}

export const useWorkspaceStore = defineStore('workspace', () => {
  // State
  const projectPath = ref<string>('')
  const fileTree = ref<FileNode[]>([])
  const openFiles = ref<OpenFile[]>([])
  const activeFilePath = ref<string>('')
  const transpilationMap = ref<Map<string, TranspilationInfo>>(new Map())
  const transpilationProgress = ref<TranspilationProgress>({
    currentFile: null,
    totalFiles: 0,
    completedFiles: 0,
    logs: []
  })
  
  // Computed
  const activeFile = computed(() => {
    return openFiles.value.find(f => f.path === activeFilePath.value)
  })
  
  const activeFileTranspilation = computed(() => {
    if (!activeFilePath.value) return null
    return transpilationMap.value.get(activeFilePath.value) || null
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
  
  function setTranspilation(tsFilePath: string, goFilePath: string, goCode: string, success: boolean = true, errorMessage?: string, warningMessage?: string) {
    const status: TranspilationStatus = errorMessage ? 'error' : warningMessage ? 'warning' : success ? 'success' : 'pending'
    
    transpilationMap.value.set(tsFilePath, {
      tsFilePath,
      goFilePath,
      goCode,
      timestamp: Date.now(),
      success,
      status,
      errorMessage,
      warningMessage,
      logs: []
    })
    
    // Mark file as transpiled in tree with status
    updateFileTranspiledStatus(tsFilePath, success, status)
  }
  
  function updateTranspilationStatus(tsFilePath: string, status: TranspilationStatus, errorMessage?: string, warningMessage?: string) {
    const existing = transpilationMap.value.get(tsFilePath)
    if (existing) {
      existing.status = status
      existing.errorMessage = errorMessage
      existing.warningMessage = warningMessage
      transpilationMap.value.set(tsFilePath, existing)
    } else {
      // Create new entry if doesn't exist
      transpilationMap.value.set(tsFilePath, {
        tsFilePath,
        goFilePath: tsFilePath.replace(/\.tsx?$/, '.go'),
        timestamp: Date.now(),
        success: status === 'success',
        status,
        errorMessage,
        warningMessage,
        logs: []
      })
    }
    
    updateFileTranspiledStatus(tsFilePath, status === 'success', status)
  }
  
  function addTranspilationLog(message: string, level: 'info' | 'error' | 'warning' | 'success' = 'info') {
    transpilationProgress.value.logs.push({
      timestamp: Date.now(),
      message,
      level
    })
  }
  
  function startTranspilation(totalFiles: number) {
    transpilationProgress.value = {
      currentFile: null,
      totalFiles,
      completedFiles: 0,
      logs: []
    }
    addTranspilationLog(`Starting transpilation of ${totalFiles} file(s)...`, 'info')
  }
  
  function setCurrentTranspilingFile(filePath: string) {
    transpilationProgress.value.currentFile = filePath
    const fileName = filePath.split('/').pop() || filePath
    addTranspilationLog(`Transpiling ${fileName}...`, 'info')
  }
  
  function completeFileTranspilation(filePath: string, success: boolean, message?: string) {
    transpilationProgress.value.completedFiles++
    const fileName = filePath.split('/').pop() || filePath
    
    if (success) {
      addTranspilationLog(`✓ ${fileName} transpiled successfully`, 'success')
    } else {
      addTranspilationLog(`✗ ${fileName} failed: ${message}`, 'error')
    }
  }
  
  function finishTranspilation() {
    const { totalFiles, completedFiles } = transpilationProgress.value
    transpilationProgress.value.currentFile = null
    addTranspilationLog(`Transpilation complete: ${completedFiles}/${totalFiles} files processed`, 'info')
  }
  
  function getTranspilation(tsFilePath: string): TranspilationInfo | undefined {
    return transpilationMap.value.get(tsFilePath)
  }
  
  function clearTranspilation(tsFilePath: string) {
    transpilationMap.value.delete(tsFilePath)
    updateFileTranspiledStatus(tsFilePath, false)
  }
  
  function updateFileTranspiledStatus(filePath: string, isTranspiled: boolean, status?: TranspilationStatus) {
    const updateNode = (nodes: FileNode[]): boolean => {
      for (const node of nodes) {
        if (node.path === filePath && node.type === 'file') {
          node.isTranspiled = isTranspiled
          if (status) {
            (node as any).transpilationStatus = status
          }
          return true
        }
        if (node.children && updateNode(node.children)) {
          return true
        }
      }
      return false
    }
    
    updateNode(fileTree.value)
  }
  
  return {
    // State
    projectPath,
    fileTree,
    openFiles,
    activeFilePath,
    transpilationMap,
    transpilationProgress,
    
    // Computed
    activeFile,
    activeFileTranspilation,
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
    loadProject,
    setTranspilation,
    getTranspilation,
    clearTranspilation,
    updateTranspilationStatus,
    addTranspilationLog,
    startTranspilation,
    setCurrentTranspilingFile,
    completeFileTranspilation,
    finishTranspilation
  }
})
