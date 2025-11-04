<template>
  <div class="home-view flex h-screen">
    <!-- Sidebar -->
    <aside class="w-64 bg-gradient-to-br from-primary-600 to-primary-800 text-white flex flex-col">
      <div class="p-6">
        <h1 class="text-2xl font-bold">TS2Go</h1>
        <p class="text-sm text-primary-100 mt-1">TypeScript to Go Transpiler</p>
      </div>
      
      <nav class="flex-1 px-3">
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
          to="/settings" 
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-white/10 transition-colors mb-1"
          :class="{ 'bg-white/20': $route.path === '/settings' }"
        >
          <i class="pi pi-cog"></i>
          <span>Settings</span>
        </router-link>
      </nav>
      
      <div class="p-4 border-t border-white/10 text-xs text-primary-100">
        Phase 21: Desktop UI v0.1.0
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- Header -->
      <header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold">Welcome to TS2Go</h2>
          <div class="flex items-center gap-3">
            <button class="btn-secondary">
              <i class="pi pi-play mr-2"></i>
              Start Transpiling
            </button>
          </div>
        </div>
      </header>

      <!-- Content Area -->
      <div class="flex-1 overflow-auto p-6">
        <div class="max-w-4xl mx-auto">
          <div class="panel mb-6">
            <h3 class="text-2xl font-bold mb-4 text-primary-600">
              🚀 Get Started with TS2Go Desktop
            </h3>
            <p class="text-gray-600 dark:text-gray-400 mb-6">
              Transform your TypeScript projects into idiomatic Go code with our powerful desktop application.
            </p>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <router-link to="/project" class="panel hover:shadow-lg transition-shadow cursor-pointer border-2 border-transparent hover:border-primary-300">
                <div class="text-4xl mb-3">📁</div>
                <h4 class="font-semibold mb-2">Open Project</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Browse and select your TypeScript project to begin
                </p>
              </router-link>
              
              <router-link to="/editor" class="panel hover:shadow-lg transition-shadow cursor-pointer border-2 border-transparent hover:border-primary-300">
                <div class="text-4xl mb-3">📝</div>
                <h4 class="font-semibold mb-2">Quick Editor</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Try transpiling code snippets in real-time
                </p>
              </router-link>
              
              <div class="panel hover:shadow-lg transition-shadow cursor-pointer border-2 border-transparent hover:border-primary-300">
                <div class="text-4xl mb-3">💡</div>
                <h4 class="font-semibold mb-2">View Examples</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Explore pre-built examples and templates
                </p>
              </div>
            </div>
          </div>

          <div class="panel">
            <h4 class="font-semibold mb-4 flex items-center gap-2">
              <i class="pi pi-info-circle text-primary-600"></i>
              Quick Stats
            </h4>
            <div class="grid grid-cols-3 gap-4 text-center">
              <div>
                <div class="text-3xl font-bold text-primary-600">0</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Projects</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-secondary-600">0</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Files Transpiled</div>
              </div>
              <div>
                <div class="text-3xl font-bold text-green-600">100%</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Success Rate</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.home-view {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}
</style>
