<template>
  <AppLayout>
    <div class="project-view">
      <!-- File Tabs -->
      <FileTabs />
      
      <!-- Main Content Area -->
      <div class="project-content">
        <!-- File Tree Sidebar -->
        <div class="file-tree-sidebar">
          <FileTree />
        </div>
        
        <!-- Editor Area -->
        <div class="editor-area">
          <div v-if="activeFile" class="editor-container">
            <div class="editor-header">
              <h3>{{ activeFile.name }}</h3>
              <span v-if="activeFile.isDirty" class="unsaved-indicator">(unsaved)</span>
            </div>
            <Textarea
              v-model="activeFile.content"
              @input="handleContentChange"
              class="code-editor"
              placeholder="Start typing..."
            ></Textarea>
          </div>
          <div v-else class="no-file-open">
            <div class="empty-state">
              <span class="empty-icon">📂</span>
              <h3>No File Open</h3>
              <p>Select a file from the tree or create a new one</p>
            </div>
          </div>
        </div>
        
        <!-- Output Panel -->
        <div class="output-panel">
          <div class="panel-header">
            <h4>Output</h4>
          </div>
          <div class="panel-content">
            <p class="placeholder">Transpilation output will appear here</p>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppLayout from '../components/AppLayout.vue'
import FileTree from '../components/FileTree.vue'
import FileTabs from '../components/FileTabs.vue'
import { useWorkspaceStore } from '../stores/workspace'

const workspace = useWorkspaceStore()

const activeFile = computed(() => workspace.activeFile)

function handleContentChange(event: Event) {
  const content = (event.target as HTMLTextAreaElement).value
  if (activeFile.value) {
    workspace.updateFileContent(activeFile.value.path, content)
  }
}
</script>

<style scoped>
.project-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.project-content {
  flex: 1;
  display: grid;
  grid-template-columns: 250px 1fr 300px;
  grid-template-rows: 1fr;
  overflow: hidden;
}

.file-tree-sidebar {
  border-right: 1px solid var(--color-border);
  overflow-y: auto;
}

.editor-area {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-soft);
}

.editor-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.unsaved-indicator {
  color: var(--color-warning);
  font-size: 12px;
}

.code-editor {
  flex: 1;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
  border: none;
  outline: none;
  resize: none;
  background: var(--color-background);
  color: var(--color-text);
}

.no-file-open {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 32px;
}

.empty-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

.empty-state h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
}

.empty-state p {
  color: var(--color-text-secondary);
  margin: 0;
}

.output-panel {
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
}

.panel-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-soft);
}

.panel-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.panel-content {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.placeholder {
  color: var(--color-text-secondary);
  font-size: 13px;
}

/* Responsive */
@media (max-width: 1024px) {
  .project-content {
    grid-template-columns: 200px 1fr;
  }
  
  .output-panel {
    display: none;
  }
}

@media (max-width: 768px) {
  .project-content {
    grid-template-columns: 1fr;
  }
  
  .file-tree-sidebar {
    display: none;
  }
}
</style>
