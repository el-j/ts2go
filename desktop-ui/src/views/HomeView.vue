<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import AppLayout from '../components/AppLayout.vue'
import RecentProjects from '../components/RecentProjects.vue'
import Button from 'primevue/button'

const router = useRouter()
const projectStore = useProjectStore()

const projectCount = computed(() => projectStore.recentProjects.length)

function goToEditor() {
  router.push('/editor')
}

function goToProject() {
  router.push('/project')
}

function goToExamples() {
  router.push('/examples')
}
</script>

<template>
  <AppLayout>
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="bg-gradient-to-r from-blue-600 to-cyan-600 text-white px-6 py-8">
        <div class="max-w-6xl mx-auto">
          <h1 class="text-3xl font-bold mb-2">Welcome to TS2Go Desktop</h1>
          <p class="text-blue-100">
            Transform your TypeScript projects into idiomatic Go code with ease
          </p>
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-auto p-6 bg-gray-50 dark:bg-gray-900">
        <div class="max-w-6xl mx-auto space-y-6">
          <!-- Quick Actions -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
            <h3 class="text-xl font-semibold mb-4 text-gray-800 dark:text-gray-200">
              Quick Start
            </h3>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <button 
                @click="goToProject"
                class="p-6 border-2 border-gray-200 dark:border-gray-700 rounded-lg hover:border-blue-500 dark:hover:border-blue-500 hover:shadow-md transition-all duration-200 text-left"
              >
                <div class="text-4xl mb-3">📁</div>
                <h4 class="font-semibold mb-2 text-gray-800 dark:text-gray-200">Open Project</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Browse and select your TypeScript project to begin
                </p>
              </button>
              
              <button 
                @click="goToEditor"
                class="p-6 border-2 border-gray-200 dark:border-gray-700 rounded-lg hover:border-blue-500 dark:hover:border-blue-500 hover:shadow-md transition-all duration-200 text-left"
              >
                <div class="text-4xl mb-3">📝</div>
                <h4 class="font-semibold mb-2 text-gray-800 dark:text-gray-200">Quick Editor</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Try transpiling code snippets in real-time
                </p>
              </button>
              
              <button 
                @click="goToExamples"
                class="p-6 border-2 border-gray-200 dark:border-gray-700 rounded-lg hover:border-blue-500 dark:hover:border-blue-500 hover:shadow-md transition-all duration-200 text-left"
              >
                <div class="text-4xl mb-3">💡</div>
                <h4 class="font-semibold mb-2 text-gray-800 dark:text-gray-200">View Examples</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Explore pre-built examples and templates
                </p>
              </button>
            </div>
          </div>

          <!-- Recent Projects -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
            <RecentProjects />
          </div>

          <!-- Quick Stats -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
            <h4 class="text-lg font-semibold mb-4 flex items-center gap-2 text-gray-800 dark:text-gray-200">
              <i class="pi pi-chart-bar text-blue-600"></i>
              Quick Stats
            </h4>
            <div class="grid grid-cols-3 gap-4 text-center">
              <div>
                <div class="text-3xl font-bold text-blue-600">{{ projectCount }}</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Recent Projects</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-green-600">80%</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">TS Coverage</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-cyan-600">v0.7.0</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Version</div>
              </div>
            </div>
          </div>

          <!-- Features Highlight -->
          <div class="bg-gradient-to-r from-blue-50 to-cyan-50 dark:from-gray-800 dark:to-gray-800 rounded-lg p-6 border border-blue-200 dark:border-gray-700">
            <h4 class="text-lg font-semibold mb-3 text-gray-800 dark:text-gray-200">
              ✨ What's New in Phase 1 & 2
            </h4>
            <ul class="space-y-2 text-sm text-gray-700 dark:text-gray-300">
              <li class="flex items-start gap-2">
                <i class="pi pi-check-circle text-green-600 mt-1"></i>
                <span><strong>Array Runtime Library:</strong> 15+ JavaScript-like array methods (filter, map, reduce, etc.)</span>
              </li>
              <li class="flex items-start gap-2">
                <i class="pi pi-check-circle text-green-600 mt-1"></i>
                <span><strong>Modern JS Features:</strong> Spread operators, destructuring, rest parameters, defaults</span>
              </li>
              <li class="flex items-start gap-2">
                <i class="pi pi-check-circle text-green-600 mt-1"></i>
                <span><strong>Fixed Desktop UI:</strong> Navigation sidebar and CLI integration now working</span>
              </li>
              <li class="flex items-start gap-2">
                <i class="pi pi-check-circle text-green-600 mt-1"></i>
                <span><strong>Coverage Improvement:</strong> TypeScript pattern support increased from 40% to 80%</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
/* Custom styles if needed */
</style>
