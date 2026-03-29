<script setup lang="ts">
import { useTheme } from '../composables/useTheme'

const { mode, toggleTheme } = useTheme()

const getIcon = () => {
  switch (mode.value) {
    case 'light':
      return '☀️'
    case 'dark':
      return '🌙'
    case 'system':
      return '💻'
    default:
      return '💻'
  }
}
</script>

<template>
  <button
    @click="toggleTheme"
    :title="`Theme: ${mode} (click to toggle)`"
    class="theme-toggle"
  >
    <span class="icon">{{ getIcon() }}</span>
  </button>
</template>

<style scoped>
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg);
  color: var(--text-h);
  cursor: pointer;
  transition: all 0.2s ease;
}

.theme-toggle:hover {
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-bg);
}

.theme-toggle:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.icon {
  font-size: 18px;
  line-height: 1;
}
</style>
