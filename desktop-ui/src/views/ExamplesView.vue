<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/AppLayout.vue'
import ExampleGallery from '@/components/ExampleGallery.vue'
import type { Example } from '@/data/examples'
import Splitter from 'primevue/splitter'
import SplitterPanel from 'primevue/splitterpanel'
import Card from 'primevue/card'
import Button from 'primevue/Button'

const router = useRouter()
const selectedExample = ref<Example | null>(null)

function loadExample(example: Example) {
  selectedExample.value = example
}

function openInEditor() {
  if (selectedExample.value) {
    // Navigate to editor with the example code
    router.push({
      name: 'editor',
      query: {
        code: selectedExample.value.typescript
      }
    })
  }
}
</script>

<template>
  <AppLayout>
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Example Gallery</h2>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Explore TypeScript to Go transpilation examples
            </p>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-hidden bg-gray-50 dark:bg-gray-900">
        <Splitter class="h-full">
          <!-- Left: Gallery -->
          <SplitterPanel :size="40" :minSize="30">
            <div class="h-full overflow-auto p-4">
              <ExampleGallery @load-example="loadExample" />
            </div>
          </SplitterPanel>

          <!-- Right: Preview -->
          <SplitterPanel :size="60" :minSize="30">
            <div class="h-full overflow-auto p-4">
              <div v-if="selectedExample">
                <Card class="mb-4">
                  <template #title>
                    <div class="flex items-center justify-between">
                      <span>{{ selectedExample.title }}</span>
                      <Button 
                        label="Open in Editor" 
                        icon="pi pi-external-link" 
                        size="small"
                        @click="openInEditor"
                      />
                    </div>
                  </template>
                  <template #content>
                    <p class="text-gray-600 dark:text-gray-400 mb-4">
                      {{ selectedExample.description }}
                    </p>

                    <!-- TypeScript Code -->
                    <div class="mb-4">
                      <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2 flex items-center gap-2">
                        <i class="pi pi-file text-blue-500"></i>
                        TypeScript Input
                      </h3>
                      <pre class="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-auto text-sm"><code>{{ selectedExample.typescript }}</code></pre>
                    </div>

                    <!-- Go Code -->
                    <div>
                      <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2 flex items-center gap-2">
                        <i class="pi pi-file text-cyan-500"></i>
                        Generated Go Output
                      </h3>
                      <pre class="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-auto text-sm"><code>{{ selectedExample.go }}</code></pre>
                    </div>
                  </template>
                </Card>
              </div>
              <div v-else class="flex items-center justify-center h-full">
                <div class="text-center text-gray-500 dark:text-gray-400">
                  <i class="pi pi-code text-4xl mb-4"></i>
                  <p>Select an example to view the code</p>
                </div>
              </div>
            </div>
          </SplitterPanel>
        </Splitter>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
pre {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  line-height: 1.5;
  max-height: 400px;
}
</style>
