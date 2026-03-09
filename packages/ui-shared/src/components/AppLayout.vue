<template>
  <div class="app-layout flex h-screen">
    <!-- Sidebar Navigation -->
    <aside :class="['sidebar', { 'collapsed': isCollapsed }]">
      <div class="sidebar-header">
        <div v-if="!isCollapsed">
          <h1 class="text-2xl font-bold">TS2Go</h1>
          <p class="text-sm mt-1">TypeScript to Go Transpiler</p>
        </div>
        <h1 v-else class="text-2xl font-bold">TS</h1>
      </div>
      
      <!-- Toggle Button -->
      <button @click="toggleSidebar" class="toggle-btn">
        <i :class="isCollapsed ? 'pi pi-angle-right' : 'pi pi-angle-left'"></i>
      </button>
      
      <nav class="flex-1 px-3 overflow-y-auto">
        <router-link 
          to="/" 
          class="nav-item"
          :class="{ 'active': $route.path === '/' }"
          v-tooltip.right="isCollapsed ? 'Home' : ''"
        >
          <i class="pi pi-home"></i>
          <span v-if="!isCollapsed">Home</span>
        </router-link>
        
        <router-link 
          to="/project" 
          class="nav-item"
          :class="{ 'active': $route.path === '/project' }"
          v-tooltip.right="isCollapsed ? 'Projects' : ''"
        >
          <i class="pi pi-folder"></i>
          <span v-if="!isCollapsed">Projects</span>
        </router-link>
        
        <router-link 
          to="/editor" 
          class="nav-item"
          :class="{ 'active': $route.path === '/editor' }"
          v-tooltip.right="isCollapsed ? 'Editor' : ''"
        >
          <i class="pi pi-code"></i>
          <span v-if="!isCollapsed">Editor</span>
        </router-link>
        
        <router-link 
          to="/examples" 
          class="nav-item"
          :class="{ 'active': $route.path === '/examples' }"
          v-tooltip.right="isCollapsed ? 'Examples' : ''"
        >
          <i class="pi pi-book"></i>
          <span v-if="!isCollapsed">Examples</span>
        </router-link>
        
        <router-link 
          to="/history" 
          class="nav-item"
          :class="{ 'active': $route.path === '/history' }"
          v-tooltip.right="isCollapsed ? 'History' : ''"
        >
          <i class="pi pi-history"></i>
          <span v-if="!isCollapsed">History</span>
        </router-link>
        
        <router-link 
          to="/analyze" 
          class="nav-item"
          :class="{ 'active': $route.path === '/analyze' }"
          v-tooltip.right="isCollapsed ? 'Analyze' : ''"
        >
          <i class="pi pi-search"></i>
          <span v-if="!isCollapsed">Analyze</span>
        </router-link>
        
        <router-link 
          to="/settings" 
          class="nav-item"
          :class="{ 'active': $route.path === '/settings' }"
          v-tooltip.right="isCollapsed ? 'Settings' : ''"
        >
          <i class="pi pi-cog"></i>
          <span v-if="!isCollapsed">Settings</span>
        </router-link>
      </nav>
      
      <div class="sidebar-footer">
        <span v-if="!isCollapsed">v0.7.0-beta</span>
        <span v-else>v0.7</span>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <slot></slot>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'

const isCollapsed = ref(false)

function toggleSidebar() {
  isCollapsed.value = !isCollapsed.value
}
</script>

<style scoped>
.app-layout {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.sidebar {
  width: 256px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: width 0.3s ease;
  position: relative;
  color: white;
}

.sidebar.collapsed {
  width: 70px;
}

.sidebar-header {
  padding: 24px;
  transition: padding 0.3s ease;
}

.sidebar.collapsed .sidebar-header {
  padding: 24px 12px;
  text-align: center;
}

.toggle-btn {
  position: absolute;
  top: 20px;
  right: -12px;
  width: 24px;
  height: 24px;
  background: white;
  border: 2px solid #667eea;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  transition: all 0.2s;
  color: #667eea;
  font-size: 12px;
}

.toggle-btn:hover {
  background: #667eea;
  color: white;
  transform: scale(1.1);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  transition: all 0.2s;
  margin-bottom: 4px;
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
}

.sidebar.collapsed .nav-item {
  justify-content: center;
  padding: 10px;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.nav-item.active {
  background: rgba(255, 255, 255, 0.2);
}

.nav-item i {
  font-size: 18px;
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.nav-item span {
  white-space: nowrap;
  overflow: hidden;
  opacity: 1;
  transition: opacity 0.2s ease;
}

.sidebar.collapsed .nav-item span {
  opacity: 0;
  width: 0;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  text-align: center;
  transition: padding 0.3s ease;
}

.sidebar.collapsed .sidebar-footer {
  padding: 16px 8px;
}
</style>
