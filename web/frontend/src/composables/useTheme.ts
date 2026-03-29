// Theme Composable
import { computed } from 'vue'
import { useThemeStore } from '../stores/theme'

export function useTheme() {
  const themeStore = useThemeStore()

  const mode = computed(() => themeStore.mode)
  const isDark = computed(() => themeStore.isDark)

  const setTheme = (newMode: 'light' | 'dark' | 'system') => {
    themeStore.setTheme(newMode)
  }

  const toggleTheme = () => {
    themeStore.toggle()
  }

  return {
    mode,
    isDark,
    setTheme,
    toggleTheme,
  }
}
