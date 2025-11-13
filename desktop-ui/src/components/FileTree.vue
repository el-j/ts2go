<template>
  <div class="file-tree">
    <div class="tree-header">
      <h3>Files</h3>
      <div class="tree-actions">
        <Button @click="createNewFile" title="New File" class="icon-btn">
          <span>📄</span>
        </Button>
        <Button @click="createNewFolder" title="New Folder" class="icon-btn">
          <span>📁</span>
        </Button>
      </div>
    </div>
    
    <div v-if="fileTree.length === 0" class="empty-state">
      <p>No files in project</p>
      <Button @click="createNewFile" class="btn-primary">Create File</Button>
    </div>
    
    <div v-else class="tree-content">
      <FileTreeNode
        v-for="node in fileTree"
        :key="node.path"
        :node="node"
        :active-path="activeFilePath"
        @select="handleSelect"
        @rename="handleRename"
        @delete="handleDelete"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { useWorkspaceStore } from '../stores/workspace'
import { useTranspileStore } from '../stores/transpile'
import FileTreeNode from './FileTreeNode.vue'

const workspace = useWorkspaceStore()
const transpileStore = useTranspileStore()

const fileTree = computed(() => workspace.fileTree)
const activeFilePath = computed(() => workspace.activeFilePath)

async function handleSelect({path, name}: {path: string, name: string}) {
  try {
    // Read file content from disk
    const content = await invoke<string>('read_file', { path })
    workspace.openFile(path, name, content)
    
    // Check if this file has been transpiled and auto-load Go output
    const transpilation = workspace.getTranspilation(path)
    if (transpilation && transpilation.success && transpilation.goCode) {
      // Auto-load the Go code into output panel
      transpileStore.currentResult = {
        success: true,
        message: `Showing transpiled output for ${name}`,
        files_transpiled: 1,
        goCode: transpilation.goCode
      }
    }
  } catch (error) {
    console.error('Failed to read file:', error)
    // Open with error message if read fails
    workspace.openFile(path, name, `Error loading file: ${error}`)
  }
}

function createNewFile() {
  const fileName = prompt('Enter file name:')
  if (fileName) {
    workspace.createFile('', fileName)
  }
}

function createNewFolder() {
  const folderName = prompt('Enter folder name:')
  if (folderName) {
    workspace.createFolder('', folderName)
  }
}

function handleRename(path: string) {
  const oldName = path.split('/').pop() || ''
  const newName = prompt('Enter new name:', oldName)
  if (newName && newName !== oldName) {
    workspace.renameFile(path, newName)
  }
}

function handleDelete(path: string) {
  if (confirm(`Delete ${path}?`)) {
    workspace.deleteFile(path)
  }
}
</script>

<style scoped>
.file-tree {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--color-background-soft);
  border-right: 1px solid var(--color-border);
}

.tree-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-soft);
}

.tree-header h3 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--color-text-secondary);
}

.tree-actions {
  display: flex;
  gap: 4px;
}

/* .icon-btn {
  padding: 4px 8px;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
} */

/* .icon-btn:hover {
  background: var(--color-background-mute);
} */

.tree-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.empty-state {
  padding: 24px 16px;
  text-align: center;
}

.empty-state p {
  color: var(--color-text-secondary);
  margin-bottom: 12px;
  font-size: 13px;
}

/* .btn-primary {
  padding: 8px 16px;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
} */

/* .btn-primary:hover {
  opacity: 0.9;
} */
</style>
