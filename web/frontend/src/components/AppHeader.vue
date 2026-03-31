<script setup lang="ts">
import { useRouter } from 'vue-router'
import ThemeToggle from './ThemeToggle.vue'

const router = useRouter()

const goToSearch = () => {
  router.push('/search')
}

const onCreateClick = () => {
  window.dispatchEvent(new CustomEvent('open-create-modal'))
}

const onHelpClick = () => {
  window.dispatchEvent(new CustomEvent('show-keyboard-help'))
}
</script>

<template>
  <header class="app-header">
    <div class="header-left">
      <h1 class="app-title" @click="router.push('/')">
        <span class="title-icon">⚡</span>
        Tasker
      </h1>
    </div>

    <div class="header-center">
      <button
        class="nav-link"
        :class="{ active: $route.name === 'list' }"
        @click="router.push('/')"
      >
        Processes
      </button>
    </div>

    <div class="header-right">
      <button class="icon-btn" title="Search ( / )" @click="goToSearch">
        🔍
      </button>
      <button class="icon-btn" title="Create process (c)" @click="onCreateClick">
        ➕
      </button>
      <button class="icon-btn" title="Keyboard shortcuts (?)" @click="onHelpClick">
        ⌨️
      </button>
      <ThemeToggle />
    </div>
  </header>
</template>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 50px;
  border-bottom: 1px solid var(--border);
  background: var(--bg);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.app-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-h);
  cursor: pointer;
  user-select: none;
}

.title-icon {
  font-size: 20px;
}

.header-center {
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav-link {
  padding: 6px 12px;
  border: none;
  background: transparent;
  color: var(--text);
  font-size: 14px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.15s;
}

.nav-link:hover {
  background: var(--code-bg);
  color: var(--text-h);
}

.nav-link.active {
  background: var(--accent-bg);
  color: var(--accent);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid transparent;
  background: transparent;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.15s;
}

.icon-btn:hover {
  background: var(--code-bg);
}

.icon-btn:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
</style>
