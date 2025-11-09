<template>
  <div class="app-layout flex h-screen">
    <!-- Sidebar Navigation -->
    <aside class="w-64 bg-gradient-to-br from-primary-600 to-primary-800 text-white flex flex-col flex-shrink-0">
      <div class="p-6">
        <h1 class="text-2xl font-bold">TS2Go</h1>
        <p class="text-sm text-primary-100 mt-1">TypeScript to Go Transpiler</p>
      </div>
      
      <nav class="flex-1 px-3 overflow-y-auto">
        <router-link 
          to="/" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/' }"
        >
          <i class="pi pi-home"></i>
          <span>Home</span>
        </router-link>
        
        <router-link 
          to="/project" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/project' }"
        >
          <i class="pi pi-folder"></i>
          <span>Projects</span>
        </router-link>
        
        <router-link 
          to="/editor" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/editor' }"
        >
          <i class="pi pi-code"></i>
          <span>Editor</span>
        </router-link>
        
        <router-link 
          to="/examples" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/examples' }"
        >
          <i class="pi pi-book"></i>
          <span>Examples</span>
        </router-link>
        
        <router-link 
          to="/history" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/history' }"
        >
          <i class="pi pi-history"></i>
          <span>History</span>
        </router-link>
        
        <router-link 
          to="/analyze" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/analyze' }"
        >
          <i class="pi pi-search"></i>
          <span>Analyze</span>
        </router-link>
        
        <router-link 
          to="/settings" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/settings' }"
        >
          <i class="pi pi-cog"></i>
          <span>Settings</span>
        </router-link>
      </nav>
      
      <div class="p-4 border-t border-white/10 text-xs text-primary-100">
        v0.7.0-beta
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <slot></slot>
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.app-layout {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}
</style>
