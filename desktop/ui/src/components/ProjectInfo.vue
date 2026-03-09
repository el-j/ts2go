<template>
  <div class="project-info" :class="{ collapsed: isCollapsed }">
    <div class="info-header" @click="toggleCollapse">
      <div class="header-content">
        <span class="icon">📂</span>
        <span class="project-name">{{ projectName }}</span>
      </div>
      <button class="collapse-btn" :title="isCollapsed ? 'Expand' : 'Collapse'">
        {{ isCollapsed ? '▶' : '▼' }}
      </button>
    </div>
    
    <div v-if="!isCollapsed" class="info-content">
      <div class="info-row">
        <span class="label">Path:</span>
        <span class="value" :title="projectPath">{{ projectPath }}</span>
      </div>
      <div class="info-row">
        <span class="label">Files:</span>
        <span class="value">{{ fileCount }} TypeScript files</span>
      </div>
      <div v-if="lastModified" class="info-row">
        <span class="label">Loaded:</span>
        <span class="value">{{ formatDate(lastModified) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useWorkspaceStore } from '../stores/workspace'

const workspace = useWorkspaceStore()
const isCollapsed = ref(false)

const projectPath = computed(() => workspace.projectPath)
const projectName = computed(() => {
  const path = workspace.projectPath
  return path.split('/').filter(Boolean).pop() || 'Project'
})

const fileCount = computed(() => {
  const countFiles = (nodes: any[]): number => {
    return nodes.reduce((count, node) => {
      if (node.type === 'file') {
        return count + 1
      }
      if (node.children) {
        return count + countFiles(node.children)
      }
      return count
    }, 0)
  }
  return countFiles(workspace.fileTree)
})

const lastModified = ref(new Date())

function toggleCollapse() {
  isCollapsed.value = !isCollapsed.value
}

function formatDate(date: Date): string {
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.project-info {
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-soft);
  transition: all 0.2s ease;
}

.info-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  cursor: pointer;
  user-select: none;
}

.info-header:hover {
  background: var(--color-background-mute);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.icon {
  font-size: 18px;
  flex-shrink: 0;
}

.project-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.collapse-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 12px;
  color: var(--color-text-secondary);
  padding: 4px;
  flex-shrink: 0;
}

.collapse-btn:hover {
  color: var(--color-text);
}

.info-content {
  padding: 8px 12px 12px;
  font-size: 12px;
  border-top: 1px solid var(--color-border);
}

.info-row {
  display: flex;
  gap: 8px;
  padding: 4px 0;
  align-items: flex-start;
}

.info-row .label {
  font-weight: 600;
  min-width: 50px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.info-row .value {
  color: var(--color-text);
  word-break: break-all;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 11px;
  line-height: 1.4;
}

.project-info.collapsed {
  border-bottom: 1px solid var(--color-border);
}

.project-info.collapsed .info-header {
  padding: 8px 12px;
}
</style>
