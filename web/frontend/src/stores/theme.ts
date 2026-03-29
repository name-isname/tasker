// Theme Store
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ThemeMode } from '../types/api'

const THEME_STORAGE_KEY = 'tasker-theme'

export const useThemeStore = defineStore('theme', () => {
  // Initialize from localStorage or default to 'system'
  const storedTheme = localStorage.getItem(THEME_STORAGE_KEY) as ThemeMode | null
  const mode = ref<ThemeMode>(storedTheme || 'system')

  // Computed: actual theme being applied (light or dark)
  const isDark = ref<boolean>(false)

  // Apply theme to DOM
  const applyTheme = () => {
    const html = document.documentElement
    let dark = false

    if (mode.value === 'system') {
      dark = window.matchMedia('(prefers-color-scheme: dark)').matches
    } else {
      dark = mode.value === 'dark'
    }

    isDark.value = dark
    if (dark) {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  // Set theme mode
  const setTheme = (newMode: ThemeMode) => {
    mode.value = newMode
    localStorage.setItem(THEME_STORAGE_KEY, newMode)
    applyTheme()
  }

  // Toggle between light and dark
  const toggle = () => {
    if (mode.value === 'light') {
      setTheme('dark')
    } else if (mode.value === 'dark') {
      setTheme('system')
    } else {
      setTheme('light')
    }
  }

  // Watch for system preference changes
  const watchSystemPreference = () => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handler = () => {
      if (mode.value === 'system') {
        applyTheme()
      }
    }
    mediaQuery.addEventListener('change', handler)
    return () => mediaQuery.removeEventListener('change', handler)
  }

  // Initialize
  applyTheme()
  watchSystemPreference()

  return {
    mode,
    isDark,
    setTheme,
    toggle,
  }
})
