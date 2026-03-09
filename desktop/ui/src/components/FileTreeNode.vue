<template>
  <div class="tree-node">
    <div
      class="node-label"
      :class="{ active: node.path === activePath, folder: node.type === 'folder' }"
      @click="handleClick"
      @contextmenu.prevent="showContextMenu"
    >
      <span class="expand-icon" v-if="node.type === 'folder'" @click.stop="toggleExpanded">
        {{ isExpanded ? '▼' : '▶' }}
      </span>
      <span class="file-icon">{{ getNodeIcon(node) }}</span>
      <span class="node-name" :title="node.displayName || node.name">{{ node.name }}</span>
      <span v-if="node.isDirty" class="dirty-indicator">●</span>
      <span 
        v-if="node.type === 'file' && node.transpilationStatus" 
        :class="['status-badge', `status-${node.transpilationStatus}`]"
        :title="getStatusTitle(node.transpilationStatus)"
      >
        {{ getStatusIcon(node.transpilationStatus) }}
      </span>
    </div>
    
    <div v-if="isExpanded && node.children && node.children.length > 0" class="node-children">
      <FileTreeNode
        v-for="child in node.children"
        :key="child.path"
        :node="child"
        :active-path="activePath"
        @select="(payload) => $emit('select', payload)"
        @rename="(path) => $emit('rename', path)"
        @delete="(path) => $emit('delete', path)"
      />
    </div>
    
    <!-- Context Menu -->
    <div
      v-if="contextMenuVisible"
      class="context-menu"
      :style="{ top: contextMenuY + 'px', left: contextMenuX + 'px' }"
      @click="hideContextMenu"
    >
      <div class="menu-item" @click="$emit('rename', node.path)">Rename</div>
      <div class="menu-item" @click="$emit('delete', node.path)">Delete</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FileNode } from '../stores/workspace'
import { getFileIcon } from '@/utils/fileUtils'

interface Props {
  node: FileNode
  activePath: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  select: [{ path: string; name: string }]
  rename: [path: string]
  delete: [path: string]
}>()

const isExpanded = ref(true)
const contextMenuVisible = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)

function handleClick() {
  if (props.node.type === 'file') {
    emit('select', { path: props.node.path, name: props.node.name })
  } else {
    toggleExpanded()
  }
}

function toggleExpanded() {
  isExpanded.value = !isExpanded.value
}

function getNodeIcon(node: FileNode): string {
  if (node.type === 'folder') return '📁'
  return getFileIcon(node.name)
}

function showContextMenu(event: MouseEvent) {
  contextMenuVisible.value = true
  contextMenuX.value = event.clientX
  contextMenuY.value = event.clientY
  
  // Hide on click outside
  setTimeout(() => {
    document.addEventListener('click', hideContextMenu, { once: true })
  }, 0)
}

function hideContextMenu() {
  contextMenuVisible.value = false
}

function getStatusIcon(status: string): string {
  switch (status) {
    case 'pending': return '⏸'
    case 'transpiling': return '⏳'
    case 'success': return '✓'
    case 'error': return '✗'
    case 'warning': return '⚠'
    default: return ''
  }
}

function getStatusTitle(status: string): string {
  switch (status) {
    case 'pending': return 'Pending transpilation'
    case 'transpiling': return 'Currently transpiling...'
    case 'success': return 'Transpiled successfully'
    case 'error': return 'Transpilation failed'
    case 'warning': return 'Transpiled with warnings'
    default: return ''
  }
}
</script>

<style scoped>
.tree-node {
  position: relative;
}

.node-label {
  display: flex;
  align-items: center;
  padding: 6px 16px;
  cursor: pointer;
  user-select: none;
  font-size: 13px;
}

.node-label:hover {
  background: var(--color-background-mute);
}

.node-label.active {
  background: var(--color-primary);
  color: white;
}

.expand-icon {
  width: 16px;
  margin-right: 4px;
  font-size: 10px;
}

.file-icon {
  margin-right: 8px;
  font-size: 16px;
}

.node-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-label.folder .node-name {
  font-weight: 500;
}

.dirty-indicator {
  color: var(--color-warning);
  margin-left: 4px;
}

.status-badge {
  margin-left: 4px;
  font-size: 11px;
  font-weight: bold;
  padding: 1px 5px;
  border-radius: 3px;
  display: inline-flex;
  align-items: center;
}

.status-pending {
  color: #6b7280;
  background: rgba(107, 114, 128, 0.15);
}

.status-transpiling {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.15);
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-success {
  color: #10b981;
  background: rgba(16, 185, 129, 0.15);
}

.status-error {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.15);
}

.status-warning {
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.15);
}

.transpiled-badge {
  margin-left: 4px;
  color: #10b981;
  font-size: 11px;
  font-weight: bold;
  background: rgba(16, 185, 129, 0.15);
  padding: 1px 5px;
  border-radius: 3px;
}

.node-children {
  padding-left: 20px;
}

.context-menu {
  position: fixed;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  min-width: 120px;
}

.menu-item {
  padding: 8px 16px;
  cursor: pointer;
  font-size: 13px;
}

.menu-item:hover {
  background: var(--color-background-mute);
}

.menu-item:first-child {
  border-radius: 6px 6px 0 0;
}

.menu-item:last-child {
  border-radius: 0 0 6px 6px;
}
</style>
