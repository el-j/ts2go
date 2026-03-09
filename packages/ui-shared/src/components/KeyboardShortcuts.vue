<template>
  <Dialog 
    v-model:visible="visible" 
    modal 
    header="Keyboard Shortcuts"
    :style="{ width: '650px' }"
  >
    <div class="shortcuts-panel">
      <div v-for="category in shortcuts" :key="category.name" class="shortcut-category">
        <h4 class="category-title">{{ category.name }}</h4>
        <div class="shortcuts-list">
          <div v-for="shortcut in category.items" :key="shortcut.key" class="shortcut-item">
            <span class="shortcut-action">{{ shortcut.action }}</span>
            <kbd class="shortcut-key">{{ shortcut.key }}</kbd>
          </div>
        </div>
      </div>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <p class="footer-note">
          <i class="pi pi-info-circle"></i>
          More shortcuts coming soon!
        </p>
        <Button label="Close" @click="visible = false" />
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'

interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const shortcuts = [
  {
    name: 'File Operations',
    items: [
      { action: 'Save File', key: 'Ctrl+S / ⌘S' },
      { action: 'Save All Files', key: 'Ctrl+Shift+S / ⌘⇧S' },
      { action: 'Open File', key: 'Ctrl+O / ⌘O' },
      { action: 'Close File', key: 'Ctrl+W / ⌘W' }
    ]
  },
  {
    name: 'Editor',
    items: [
      { action: 'Find', key: 'Ctrl+F / ⌘F' },
      { action: 'Replace', key: 'Ctrl+H / ⌘H' },
      { action: 'Go to Definition', key: 'F12' },
      { action: 'Format Document', key: 'Shift+Alt+F / ⇧⌥F' },
      { action: 'Toggle Comment', key: 'Ctrl+/ / ⌘/' }
    ]
  },
  {
    name: 'Transpilation',
    items: [
      { action: 'Transpile Current File', key: 'Ctrl+T / ⌘T' },
      { action: 'Transpile All Files', key: 'Ctrl+Shift+T / ⌘⇧T' },
      { action: 'Run Code', key: 'Ctrl+Enter / ⌘↵' }
    ]
  },
  {
    name: 'Navigation',
    items: [
      { action: 'Quick Open', key: 'Ctrl+P / ⌘P' },
      { action: 'Go to Line', key: 'Ctrl+G / ⌘G' },
      { action: 'Focus File Tree', key: 'Ctrl+B / ⌘B' },
      { action: 'Focus Editor', key: 'Ctrl+1 / ⌘1' }
    ]
  },
  {
    name: 'General',
    items: [
      { action: 'Show Shortcuts', key: 'Ctrl+K Ctrl+S / ⌘K ⌘S' },
      { action: 'Open Settings', key: 'Ctrl+, / ⌘,' },
      { action: 'Toggle Sidebar', key: 'Ctrl+Shift+B / ⌘⇧B' }
    ]
  }
]
</script>

<style scoped>
.shortcuts-panel {
  padding: 1rem 0;
  max-height: 500px;
  overflow-y: auto;
}

.shortcut-category {
  margin-bottom: 2rem;
}

.shortcut-category:last-child {
  margin-bottom: 0;
}

.category-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--primary-color);
  margin-bottom: 0.75rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--surface-border);
}

.shortcuts-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.shortcut-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0.75rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.shortcut-item:hover {
  background-color: var(--surface-50);
}

.shortcut-action {
  flex: 1;
  color: var(--text-color);
  font-size: 0.9rem;
}

.shortcut-key {
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaf0 100%);
  border: 1px solid var(--surface-border);
  border-radius: 4px;
  padding: 0.3rem 0.6rem;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 0.85rem;
  color: var(--text-color-secondary);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  white-space: nowrap;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.footer-note {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-color-secondary);
  font-size: 0.85rem;
  margin: 0;
}

.footer-note i {
  color: var(--primary-color);
}
</style>
