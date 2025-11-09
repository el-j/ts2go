<script setup lang="ts">
import { ref } from 'vue'
import DataView from 'primevue/dataview'
import Card from 'primevue/card'
import Button from 'primevue/Button'
import Tag from 'primevue/tag'
import { EXAMPLES, type Example } from '@/data/examples'

const emit = defineEmits<{
  loadExample: [example: Example]
}>()

const examples = ref<Example[]>(EXAMPLES)
const selectedCategory = ref<string>('all')
const filteredExamples = ref(examples.value)

function filterByCategory(category: string) {
  selectedCategory.value = category
  if (category === 'all') {
    filteredExamples.value = examples.value
  } else {
    filteredExamples.value = examples.value.filter(ex => ex.category === category)
  }
}

function getCategoryColor(category: string): string {
  switch (category) {
    case 'basic': return 'success'
    case 'intermediate': return 'warn'
    case 'advanced': return 'danger'
    default: return 'info'
  }
}
</script>

<template>
  <div class="example-gallery">
    <!-- Category Filter -->
    <div class="mb-4 flex gap-2">
      <Button 
        label="All" 
        :outlined="selectedCategory !== 'all'"
        @click="filterByCategory('all')"
        size="small"
      />
      <Button 
        label="Basic" 
        :outlined="selectedCategory !== 'basic'"
        severity="success"
        @click="filterByCategory('basic')"
        size="small"
      />
      <Button 
        label="Intermediate" 
        :outlined="selectedCategory !== 'intermediate'"
        severity="warn"
        @click="filterByCategory('intermediate')"
        size="small"
      />
      <Button 
        label="Advanced" 
        :outlined="selectedCategory !== 'advanced'"
        severity="danger"
        @click="filterByCategory('advanced')"
        size="small"
      />
    </div>

    <!-- Examples Grid -->
    <DataView :value="filteredExamples" layout="grid">
      <template #grid="slotProps">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <Card 
            v-for="example in slotProps.items" 
            :key="example.id"
            class="hover:shadow-lg transition-shadow cursor-pointer"
          >
            <template #title>
              <div class="flex items-center justify-between">
                <span class="text-lg">{{ example.title }}</span>
                <Tag 
                  :value="example.category" 
                  :severity="getCategoryColor(example.category)"
                />
              </div>
            </template>
            <template #content>
              <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
                {{ example.description }}
              </p>
              <div class="flex gap-2">
                <Button 
                  label="Load Example" 
                  icon="pi pi-code" 
                  size="small"
                  @click="emit('loadExample', example)"
                />
              </div>
            </template>
          </Card>
        </div>
      </template>
    </DataView>
  </div>
</template>

<style scoped>
.example-gallery {
  padding: 1rem;
}
</style>
