<script setup lang="ts">
import { computed } from 'vue'
import { invoke } from '@tauri-apps/api/core'
import { useProjectStore, type Project } from '@/stores/project'
import { useWorkspaceStore } from '@/stores/workspace'
import { useRouter } from 'vue-router'
import Card from 'primevue/card'
import Button from 'primevue/button'

const projectStore = useProjectStore()
const workspace = useWorkspaceStore()
const router = useRouter()

const recentProjects = computed(() => projectStore.recentProjects)

async function openProject(project: Project) {
  try {
    // Load project files from backend
    const projectInfo = await invoke('load_project_folder', { path: project.path })
    
    // Load project into workspace
    workspace.loadProject(projectInfo)
    
    // Update project store
    projectStore.setCurrentProject(project)
    
    // Navigate to project view
    router.push('/project')
  } catch (error) {
    console.error('Failed to load project:', error)
  }
}

function togglePin(project: Project) {
  projectStore.togglePinProject(project.path)
}

function removeProject(project: Project) {
  projectStore.removeFromRecentProjects(project.path)
}

function formatDate(dateString: string | undefined) {
  if (!dateString) return 'Never'
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60)
  
  if (diffInHours < 1) return 'Just now'
  if (diffInHours < 24) return `${Math.floor(diffInHours)} hours ago`
  if (diffInHours < 48) return 'Yesterday'
  if (diffInHours < 168) return `${Math.floor(diffInHours / 24)} days ago`
  return date.toLocaleDateString()
}

function getProjectInitials(name: string) {
  return name
    .split(' ')
    .map(word => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}
</script>

<template>
  <div class="recent-projects">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200">
        Recent Projects
      </h3>
      <Button 
        v-if="recentProjects.length > 0"
        label="Clear All" 
        icon="pi pi-trash" 
        severity="danger"
        text
        size="small"
        @click="projectStore.clearRecentProjects()"
      />
    </div>

    <!-- Empty State -->
    <div 
      v-if="recentProjects.length === 0" 
      class="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-700"
    >
      <i class="pi pi-folder-open text-4xl text-gray-400 dark:text-gray-600 mb-3"></i>
      <p class="text-gray-500 dark:text-gray-400">No recent projects</p>
      <p class="text-sm text-gray-400 dark:text-gray-500 mt-1">
        Open a project to see it here
      </p>
    </div>

    <!-- Projects Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <Card 
        v-for="project in recentProjects" 
        :key="project.id"
        class="hover:shadow-lg transition-shadow duration-200 cursor-pointer relative"
      >
        <!-- Pin Badge -->
        <div 
          v-if="project.isPinned" 
          class="absolute top-2 right-2 bg-yellow-100 dark:bg-yellow-900 text-yellow-700 dark:text-yellow-300 rounded-full p-1"
        >
          <i class="pi pi-thumbtack text-xs"></i>
        </div>

        <template #header>
          <div class="p-4 bg-gradient-to-r from-blue-500 to-cyan-500 text-white">
            <div class="text-2xl font-bold">
              {{ getProjectInitials(project.name) }}
            </div>
          </div>
        </template>

        <template #content>
          <div @click="openProject(project)">
            <h4 class="font-semibold text-gray-800 dark:text-gray-200 mb-1">
              {{ project.name }}
            </h4>
            <p class="text-sm text-gray-500 dark:text-gray-400 truncate mb-2">
              {{ project.path }}
            </p>
            <div class="flex items-center text-xs text-gray-400 dark:text-gray-500">
              <i class="pi pi-clock mr-1"></i>
              <span>{{ formatDate(project.lastOpened) }}</span>
              <span v-if="project.accessCount" class="ml-3">
                <i class="pi pi-eye mr-1"></i>
                {{ project.accessCount }} times
              </span>
            </div>
          </div>
        </template>

        <template #footer>
          <div class="flex gap-2">
            <Button 
              icon="pi pi-folder-open" 
              label="Open" 
              size="small"
              @click="openProject(project)"
              class="flex-1"
            />
            <Button 
              :icon="project.isPinned ? 'pi pi-thumbtack-fill' : 'pi pi-thumbtack'" 
              :severity="project.isPinned ? 'warning' : 'secondary'"
              outlined
              size="small"
              @click.stop="togglePin(project)"
              v-tooltip.top="project.isPinned ? 'Unpin' : 'Pin'"
            />
            <Button 
              icon="pi pi-times" 
              severity="danger"
              outlined
              size="small"
              @click.stop="removeProject(project)"
              v-tooltip.top="'Remove from recent'"
            />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.recent-projects {
  @apply w-full;
}
</style>
