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
      <span class="file-icon">{{ getFileIcon(node) }}</span>
      <span class="node-name">{{ node.name }}</span>
      <span v-if="node.isDirty" class="dirty-indicator">●</span>
    </div>
    
    <div v-if="isExpanded && node.children && node.children.length > 0" class="node-children">
      <FileTreeNode
        v-for="child in node.children"
        :key="child.path"
        :node="child"
        :active-path="activePath"
        @select="$emit('select', $event.path, $event.name)"
        @rename="$emit('rename', $event)"
        @delete="$emit('delete', $event)"
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

function getFileIcon(node: FileNode): string {
  if (node.type === 'folder') return '📁'
  
  const ext = node.name.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'ts': return '🔷'
    case 'js': return '🟨'
    case 'json': return '📋'
    case 'md': return '📝'
    case 'go': return '🔵'
    case 'vue': return '💚'
    default: return '📄'
  }
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
}

.dirty-indicator {
  color: var(--color-warning);
  margin-left: 4px;
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
