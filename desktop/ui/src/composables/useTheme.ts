import { watch, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'

export function useTheme() {
  const settingsStore = useSettingsStore()

  const applyTheme = (theme: 'light' | 'dark' | 'system') => {
    const html = document.documentElement
    
    if (theme === 'system') {
      // Use system preference
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      if (prefersDark) {
        html.classList.add('dark')
      } else {
        html.classList.remove('dark')
      }
    } else if (theme === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  const initTheme = () => {
    // Apply initial theme
    applyTheme(settingsStore.settings.theme)

    // Watch for theme changes
    watch(
      () => settingsStore.settings.theme,
      (newTheme) => {
        applyTheme(newTheme)
      }
    )

    // Listen for system theme changes when in system mode
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handleSystemThemeChange = () => {
      if (settingsStore.settings.theme === 'system') {
        applyTheme('system')
      }
    }
    
    mediaQuery.addEventListener('change', handleSystemThemeChange)
  }

  onMounted(() => {
    initTheme()
  })

  return {
    applyTheme
  }
}
