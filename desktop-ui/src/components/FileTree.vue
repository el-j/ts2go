<template>
  <div class="file-tree">
    <div class="tree-header">
      <h3>Files</h3>
      <div class="tree-actions">
        <button @click="createNewFile" title="New File" class="icon-btn">
          <span>📄</span>
        </button>
        <button @click="createNewFolder" title="New Folder" class="icon-btn">
          <span>📁</span>
        </button>
      </div>
    </div>
    
    <div v-if="fileTree.length === 0" class="empty-state">
      <p>No files in project</p>
      <button @click="createNewFile" class="btn-primary">Create File</button>
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
import { useWorkspaceStore } from '../stores/workspace'
import FileTreeNode from './FileTreeNode.vue'

const workspace = useWorkspaceStore()

const fileTree = computed(() => workspace.fileTree)
const activeFilePath = computed(() => workspace.activeFilePath)

function handleSelect(path: string, name: string) {
  workspace.openFile(path, name)
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
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
}

.tree-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.tree-actions {
  display: flex;
  gap: 4px;
}

.icon-btn {
  padding: 4px 8px;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.icon-btn:hover {
  background: var(--color-background-mute);
}

.tree-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.empty-state {
  padding: 32px 16px;
  text-align: center;
}

.empty-state p {
  color: var(--color-text-secondary);
  margin-bottom: 16px;
}

.btn-primary {
  padding: 8px 16px;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary:hover {
  opacity: 0.9;
}
</style>
