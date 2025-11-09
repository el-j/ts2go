<script setup lang="ts">
import { computed } from 'vue'
import { useHistoryStore } from '@/stores/history'
import AppLayout from '../components/AppLayout.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Chart from 'primevue/chart'

const historyStore = useHistoryStore()

const chartData = computed(() => {
  const last10 = historyStore.builds.slice(0, 10).reverse()
  return {
    labels: last10.map((_b, i) => `Build ${i + 1}`),
    datasets: [
      {
        label: 'Duration (ms)',
        data: last10.map(b => b.duration),
        backgroundColor: 'rgba(102, 126, 234, 0.2)',
        borderColor: 'rgb(102, 126, 234)',
        borderWidth: 2
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      title: {
        display: true,
        text: 'Duration (ms)'
      }
    }
  }
}

function getStatusSeverity(status: string) {
  switch (status) {
    case 'success': return 'success'
    case 'failed': return 'danger'
    case 'partial': return 'warn'
    default: return 'info'
  }
}

function formatDuration(ms: number) {
  if (ms < 1000) return `${ms}ms`
  const seconds = (ms / 1000).toFixed(1)
  return `${seconds}s`
}

function formatTimestamp(date: Date) {
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }).format(date)
}

function confirmClearHistory() {
  if (confirm('Are you sure you want to clear all build history? This cannot be undone.')) {
    historyStore.clearHistory()
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
            <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Build History</h2>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Track your transpilation history and performance
            </p>
          </div>
          <Button 
            label="Clear History" 
            icon="pi pi-trash" 
            severity="danger" 
            size="small"
            outlined
            @click="confirmClearHistory"
            :disabled="historyStore.builds.length === 0"
          />
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-auto p-6 bg-gray-50 dark:bg-gray-900">
        <div class="max-w-7xl mx-auto space-y-6">
          <!-- Statistics Cards -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-3xl font-bold text-primary-600 dark:text-primary-400">
                    {{ historyStore.builds.length }}
                  </div>
                  <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Total Builds</div>
                </div>
              </template>
            </Card>

            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-3xl font-bold text-green-600 dark:text-green-400">
                    {{ historyStore.successRate.toFixed(1) }}%
                  </div>
                  <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Success Rate</div>
                </div>
              </template>
            </Card>

            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-3xl font-bold text-blue-600 dark:text-blue-400">
                    {{ formatDuration(historyStore.averageDuration) }}
                  </div>
                  <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Avg Duration</div>
                </div>
              </template>
            </Card>

            <Card>
              <template #content>
                <div class="text-center">
                  <div class="text-3xl font-bold text-red-600 dark:text-red-400">
                    {{ historyStore.failedBuilds.length }}
                  </div>
                  <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">Failed Builds</div>
                </div>
              </template>
            </Card>
          </div>

          <!-- Performance Chart -->
          <Card v-if="historyStore.builds.length > 0">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-chart-line"></i>
                Build Duration Trend (Last 10 Builds)
              </div>
            </template>
            <template #content>
              <div style="height: 250px">
                <Chart type="line" :data="chartData" :options="chartOptions" />
              </div>
            </template>
          </Card>

          <!-- Build History Table -->
          <Card>
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-history"></i>
                Build History
              </div>
            </template>
            <template #content>
              <DataTable 
                :value="historyStore.builds" 
                stripedRows
                paginator 
                :rows="10"
                :rowsPerPageOptions="[10, 20, 50]"
                :emptyMessage="'No build history yet. Start transpiling to see history.'"
              >
                <Column field="timestamp" header="Timestamp" sortable>
                  <template #body="slotProps">
                    {{ formatTimestamp(slotProps.data.timestamp) }}
                  </template>
                </Column>

                <Column field="projectPath" header="Project" sortable>
                  <template #body="slotProps">
                    <span class="text-sm font-mono">
                      {{ slotProps.data.projectPath || 'Single File' }}
                    </span>
                  </template>
                </Column>

                <Column field="filesProcessed" header="Files" sortable>
                  <template #body="slotProps">
                    {{ slotProps.data.filesProcessed }} / {{ slotProps.data.totalFiles }}
                  </template>
                </Column>

                <Column field="duration" header="Duration" sortable>
                  <template #body="slotProps">
                    {{ formatDuration(slotProps.data.duration) }}
                  </template>
                </Column>

                <Column field="status" header="Status" sortable>
                  <template #body="slotProps">
                    <Tag 
                      :value="slotProps.data.status" 
                      :severity="getStatusSeverity(slotProps.data.status)"
                    />
                  </template>
                </Column>

                <Column field="errors" header="Errors" sortable>
                  <template #body="slotProps">
                    <span v-if="slotProps.data.errors > 0" class="text-red-600 dark:text-red-400">
                      {{ slotProps.data.errors }}
                    </span>
                    <span v-else class="text-gray-400">0</span>
                  </template>
                </Column>

                <Column field="warnings" header="Warnings" sortable>
                  <template #body="slotProps">
                    <span v-if="slotProps.data.warnings > 0" class="text-yellow-600 dark:text-yellow-400">
                      {{ slotProps.data.warnings }}
                    </span>
                    <span v-else class="text-gray-400">0</span>
                  </template>
                </Column>

                <Column header="Actions">
                  <template #body="slotProps">
                    <Button 
                      icon="pi pi-trash" 
                      severity="danger" 
                      text
                      size="small"
                      @click="historyStore.removeBuild(slotProps.data.id)"
                    />
                  </template>
                </Column>
              </DataTable>
            </template>
          </Card>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
/* Custom styles for history view */
</style>
