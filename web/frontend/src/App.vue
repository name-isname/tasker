<script setup lang="ts">
import { ref } from 'vue'
import AppHeader from './components/AppHeader.vue'
import CreateProcessModal from './components/CreateProcessModal.vue'
import KeyboardShortcutsHelp from './components/KeyboardShortcutsHelp.vue'
import { useGlobalKeyboard } from './composables/useKeyboard'

const showCreateModal = ref(false)
const showHelp = ref(false)

useGlobalKeyboard(showHelp)

// Listen for custom event to open modal
window.addEventListener('open-create-modal', () => {
  showCreateModal.value = true
})

// Listen for keyboard help event (fallback)
window.addEventListener('show-keyboard-help', () => {
  showHelp.value = true
})
</script>

<template>
  <div id="app" class="app">
    <AppHeader />
    <main class="main-content">
      <RouterView />
    </main>
    <CreateProcessModal v-model:show="showCreateModal" />
    <KeyboardShortcutsHelp v-model:show="showHelp" />
  </div>
</template>

<style>
#app {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  overflow: hidden;
}

/* Remove template styles from style.css that we don't need */
#app > *:not(main) {
  min-height: 0;
}
</style>
