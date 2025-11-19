import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService, type LoginRequest, type RegisterRequest } from '../services/auth'
import { usePlatform } from '../composables/usePlatform'

export const useAuthStore = defineStore('auth', () => {
  const { isWeb } = usePlatform()
  const token = ref<string | null>(null)
  const user = ref<any>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => {
    // For web, check if token exists
    if (isWeb) {
      return !!token.value
    }
    // For Tauri, always authenticated (no login required)
    return true
  })

  async function login(credentials: LoginRequest) {
    if (!isWeb) {
      // Tauri doesn't need authentication
      return
    }

    isLoading.value = true
    error.value = null
    try {
      const response = await authService.login(credentials)
      token.value = response.token
      user.value = response.user
    } catch (e: any) {
      error.value = e.message || 'Login failed'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function register(data: RegisterRequest) {
    if (!isWeb) {
      // Tauri doesn't need authentication
      return
    }

    isLoading.value = true
    error.value = null
    try {
      const response = await authService.register(data)
      token.value = response.token
      user.value = response.user
    } catch (e: any) {
      error.value = e.message || 'Registration failed'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    if (!isWeb) {
      // Tauri doesn't need authentication
      return
    }

    await authService.logout()
    token.value = null
    user.value = null
  }

  async function loadUser() {
    if (!isWeb) {
      // Tauri doesn't need authentication
      return
    }

    const savedToken = authService.getToken()
    if (!savedToken) return

    token.value = savedToken
    try {
      user.value = await authService.getCurrentUser()
    } catch (e) {
      // Token invalid, clear it
      token.value = null
      authService.logout()
    }
  }

  return {
    token,
    user,
    isLoading,
    error,
    isAuthenticated,
    login,
    register,
    logout,
    loadUser,
  }
})
