<script setup lang="ts">
import { computed } from 'vue'
import Message from 'primevue/message'
import Accordion from 'primevue/accordion'
import AccordionTab from 'primevue/accordiontab'

interface TranspilationError {
  file?: string
  line?: number
  column?: number
  message: string
  context?: string
}

interface Props {
  errors: TranspilationError[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  jumpToError: [error: TranspilationError]
}>()

const errorCount = computed(() => props.errors.length)
const hasErrors = computed(() => errorCount.value > 0)
</script>

<template>
  <div v-if="hasErrors" class="error-display">
    <Message severity="error" :closable="false" class="mb-3">
      <div class="flex items-center justify-between w-full">
        <span class="font-semibold">
          <i class="pi pi-exclamation-triangle mr-2"></i>
          {{ errorCount }} {{ errorCount === 1 ? 'Error' : 'Errors' }} Found
        </span>
      </div>
    </Message>

    <Accordion :multiple="true">
      <AccordionTab 
        v-for="(error, index) in errors" 
        :key="index"
      >
        <template #header>
          <div class="flex items-center gap-2 flex-1">
            <i class="pi pi-times-circle text-red-500"></i>
            <span class="font-medium">
              {{ error.file || 'Unknown file' }}
              <span v-if="error.line" class="text-sm text-gray-600 dark:text-gray-400">
                (Line {{ error.line }}{{ error.column ? `:${error.column}` : '' }})
              </span>
            </span>
          </div>
        </template>

        <div class="space-y-3">
          <!-- Error Message -->
          <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3">
            <p class="text-sm text-red-800 dark:text-red-200">{{ error.message }}</p>
          </div>

          <!-- Error Context (code snippet) -->
          <div v-if="error.context" class="bg-gray-900 rounded p-3">
            <pre class="text-sm text-gray-100 overflow-auto"><code>{{ error.context }}</code></pre>
          </div>

          <!-- Actions -->
          <div class="flex gap-2">
            <button
              v-if="error.file && error.line"
              @click="emit('jumpToError', error)"
              class="px-3 py-1.5 text-sm bg-primary-600 hover:bg-primary-700 text-white rounded flex items-center gap-2"
            >
              <i class="pi pi-arrow-right"></i>
              Jump to Error
            </button>
          </div>
        </div>
      </AccordionTab>
    </Accordion>
  </div>

  <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
    <i class="pi pi-check-circle text-4xl mb-2 text-green-500"></i>
    <p>No errors found</p>
  </div>
</template>

<style scoped>
.error-display {
  padding: 1rem;
}

pre {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  line-height: 1.5;
  margin: 0;
}
</style>
