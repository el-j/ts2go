<template>
  <Dialog 
    v-model:visible="visible" 
    modal 
    header="Unsaved Changes"
    :style="{ width: '450px' }"
    :closable="false"
  >
    <div class="unsaved-changes-dialog">
      <p class="dialog-message">
        You have unsaved changes in {{ dirtyFiles.length }} {{ dirtyFiles.length === 1 ? 'file' : 'files' }}:
      </p>
      
      <div class="file-list-container">
        <ul class="file-list">
          <li v-for="file in dirtyFiles" :key="file.path" class="file-item">
            <i class="pi pi-file text-primary"></i>
            <span class="file-name">{{ file.name }}</span>
            <span class="file-status">●</span>
          </li>
        </ul>
      </div>
      
      <p class="dialog-question">
        Do you want to save them before closing?
      </p>
    </div>
    
    <template #footer>
      <div class="dialog-actions">
        <Button 
          label="Don't Save" 
          severity="secondary"
          text
          @click="discardChanges"
          class="discard-btn"
        />
        <Button 
          label="Cancel" 
          severity="secondary"
          @click="cancel"
        />
        <Button 
          label="Save All" 
          @click="saveAll"
          autofocus
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import { useWorkspaceStore } from '@/stores/workspace'
import { useSaveFile } from '@/composables/useSaveFile'

interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'save': []
  'discard': []
  'cancel': []
}>()

const workspace = useWorkspaceStore()
const { saveAll: saveAllFiles } = useSaveFile()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const dirtyFiles = computed(() => {
  return workspace.openFiles.filter(f => f.isDirty)
})

async function saveAll() {
  await saveAllFiles()
  visible.value = false
  emit('save')
}

function discardChanges() {
  visible.value = false
  emit('discard')
}

function cancel() {
  visible.value = false
  emit('cancel')
}
</script>

<style scoped>
.unsaved-changes-dialog {
  padding: 1rem 0;
}

.dialog-message {
  margin-bottom: 1rem;
  color: var(--text-color);
  font-size: 0.95rem;
}

.file-list-container {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid var(--surface-border);
  border-radius: 4px;
  padding: 0.5rem;
  background: var(--surface-50);
  margin-bottom: 1rem;
}

.file-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.file-item:hover {
  background-color: var(--surface-100);
}

.file-name {
  flex: 1;
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 0.9rem;
}

.file-status {
  color: var(--primary-color);
  font-size: 0.7rem;
}

.dialog-question {
  margin-top: 1rem;
  margin-bottom: 0;
  color: var(--text-color);
  font-weight: 500;
}

.dialog-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.discard-btn {
  color: var(--red-500) !important;
}

.discard-btn:hover {
  background-color: var(--red-50) !important;
}
</style>
