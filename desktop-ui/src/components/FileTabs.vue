<template>
  <div class="file-tabs">
    <div v-if="openFiles.length === 0" class="no-tabs">
      <!-- <span>No files open</span> -->
    </div>
    
    <div v-else class="tabs-container">
      <div
        v-for="file in openFiles"
        :key="file.path"
        class="tab"
        :class="{ active: file.path === activeFilePath }"
        @click="setActiveFile(file.path)"
      >
        <span class="tab-icon">{{ getFileIcon(file.name) }}</span>
        <span class="tab-name">{{ file.name }}</span>
        <span v-if="file.isDirty" class="dirty-indicator">●</span>
        <Button class="close-btn" @click.stop="closeFile(file.path)" title="Close">
          ×
        </Button>
      </div>
    </div>
    <!-- <ProjectLoader />  -->
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
// import ProjectLoader from './ProjectLoader.vue'
import { useWorkspaceStore } from '../stores/workspace'
const workspace = useWorkspaceStore()

const openFiles = computed(() => workspace.openFiles)
const activeFilePath = computed(() => workspace.activeFilePath)

function setActiveFile(path: string) {
  workspace.setActiveFile(path)
}

function closeFile(path: string) {
  workspace.closeFile(path)
}

function getFileIcon(fileName: string): string {
  const ext = fileName.split('.').pop()?.toLowerCase()
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
</script>

<style scoped>
.file-tabs {
  display: flex;
  background: var(--color-background-soft);
  border-bottom: 1px solid var(--color-border);
  overflow-x: auto;
  min-height: 40px;
}

.no-tabs {
  display: flex;
  align-items: center;
  padding: 0 16px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.tabs-container {
  display: flex;
  flex: 1;
}

.tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--color-background-soft);
  border-right: 1px solid var(--color-border);
  cursor: pointer;
  user-select: none;
  min-width: 120px;
  max-width: 200px;
  position: relative;
}

.tab:hover {
  background: var(--color-background-mute);
}

.tab.active {
  background: var(--color-background);
  border-bottom: 2px solid var(--color-primary);
}

.tab-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.tab-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
}

.dirty-indicator {
  color: var(--color-warning);
  font-size: 20px;
  line-height: 1;
  flex-shrink: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: 20px;
  line-height: 1;
  padding: 0 4px;
  cursor: pointer;
  flex-shrink: 0;
}

.close-btn:hover {
  color: var(--color-text);
  background: var(--color-background-mute);
  border-radius: 4px;
}

.tab.active .close-btn {
  color: var(--color-text);
}
</style>
