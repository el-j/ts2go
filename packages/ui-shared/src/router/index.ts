import { createRouter, createWebHistory } from 'vue-router'
import { usePlatform } from '../composables/usePlatform'
import { useAuthStore } from '../stores/auth'
import HomeView from '../views/HomeView.vue'

const { isWeb } = usePlatform()

const router = createRouter({
  history: createWebHistory('/'),
  routes: [
    // Auth routes (web only)
    ...(isWeb ? [
      {
        path: '/login',
        name: 'login',
        component: () => import('../views/Login.vue'),
        meta: { requiresGuest: true }
      },
      {
        path: '/register',
        name: 'register',
        component: () => import('../views/Register.vue'),
        meta: { requiresGuest: true }
      }
    ] : []),
    // Main app routes
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: isWeb }
    },
    {
      path: '/project',
      name: 'project',
      component: () => import('../views/ProjectView.vue'),
      meta: { requiresAuth: isWeb }
    },
    {
      path: '/editor',
      name: 'editor',
      component: () => import('../views/EditorView.vue'),
      meta: { requiresAuth: isWeb }
    },
    {
      path: '/examples',
      name: 'examples',
      component: () => import('../views/ExamplesView.vue'),
      meta: { requiresAuth: isWeb }
    },
    {
      path: '/history',
      name: 'history',
      component: () => import('../views/HistoryView.vue'),
      meta: { requiresAuth: isWeb }
    },
    {
      path: '/analyze',
      name: 'analyze',
      component: () => import('../views/AnalyzeView.vue'),
      meta: { requiresAuth: isWeb }
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue'),
      meta: { requiresAuth: isWeb }
    }
  ]
})

// Navigation guard for web authentication
if (isWeb) {
  router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()
    
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      next({ name: 'login' })
    } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
      next({ name: 'home' })
    } else {
      next()
    }
  })
}

export default router

