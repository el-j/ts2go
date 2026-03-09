<script setup lang="ts">
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'

const visible = defineModel<boolean>('visible', { default: false })

interface ShortcutInfo {
  keys: string
  description: string
  context: string
}

const shortcuts = ref<ShortcutInfo[]>([
  {
    keys: 'Ctrl+S',
    description: 'Transpile current file',
    context: 'Editor'
  },
  {
    keys: 'Ctrl+O',
    description: 'Open file dialog',
    context: 'Global'
  },
  {
    keys: 'Ctrl+W',
    description: 'Close current tab',
    context: 'Editor'
  },
  {
    keys: 'Ctrl+Shift+W',
    description: 'Close all tabs',
    context: 'Editor'
  },
  {
    keys: 'Ctrl+F',
    description: 'Find in file',
    context: 'Editor'
  },
  {
    keys: 'Ctrl+H',
    description: 'Find and replace',
    context: 'Editor'
  },
  {
    keys: 'Ctrl+/',
    description: 'Show keyboard shortcuts',
    context: 'Global'
  },
  {
    keys: 'Ctrl+B',
    description: 'Toggle sidebar',
    context: 'Global'
  },
  {
    keys: 'Ctrl+`',
    description: 'Toggle terminal',
    context: 'Global'
  },
  {
    keys: 'Ctrl+Shift+P',
    description: 'Command palette',
    context: 'Global'
  },
  {
    keys: 'Escape',
    description: 'Close dialog',
    context: 'Global'
  }
])
</script>

<template>
  <Dialog 
    v-model:visible="visible" 
    modal 
    header="Keyboard Shortcuts"
    :style="{ width: '50rem' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
  >
    <div class="mb-4">
      <p class="text-gray-600 dark:text-gray-400 text-sm">
        Use these keyboard shortcuts to navigate and work faster in TS2Go Desktop.
      </p>
    </div>

    <DataTable 
      :value="shortcuts" 
      stripedRows
      class="text-sm"
    >
      <Column field="keys" header="Keys" style="width: 150px">
        <template #body="slotProps">
          <Tag :value="slotProps.data.keys" severity="secondary" />
        </template>
      </Column>
      <Column field="description" header="Description" />
      <Column field="context" header="Context" style="width: 120px">
        <template #body="slotProps">
          <Tag 
            :value="slotProps.data.context" 
            :severity="slotProps.data.context === 'Global' ? 'info' : 'success'"
          />
        </template>
      </Column>
    </DataTable>

    <template #footer>
      <div class="text-xs text-gray-500 dark:text-gray-400">
        <i class="pi pi-info-circle mr-1"></i>
        On macOS, use Cmd instead of Ctrl
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
/* Custom styles for shortcuts dialog */
</style>
