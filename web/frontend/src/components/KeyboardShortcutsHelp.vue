<script setup lang="ts">
const props = defineProps<{
  show?: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const shortcuts = [
  { key: 'c', description: 'Create new process' },
  { key: '/', description: 'Focus search' },
  { key: 'n', description: 'Navigate to process list' },
  { key: 'Escape', description: 'Close modal/dialog' },
  { key: '?', description: 'Show this help' },
  { section: 'In Process List' },
  { key: 'j / k', description: 'Navigate up/down' },
  { key: '↑ / ↓', description: 'Navigate up/down' },
  { key: 'Enter', description: 'Open process detail' },
  { key: 'd', description: 'Delete selected process' },
  { key: 'e', description: 'Edit selected process' },
  { section: 'In Process Detail' },
  { key: 'l', description: 'Focus add log input' },
  { key: '1', description: 'Set status: Running' },
  { key: '2', description: 'Set status: Blocked' },
  { key: '3', description: 'Set status: Suspended' },
  { key: '4', description: 'Set status: Terminated' },
]

const onBackdropClick = (e: MouseEvent) => {
  if (e.target === e.currentTarget) {
    emit('close')
  }
}
</script>

<template>
  <Transition name="fade">
    <div v-if="show" class="modal-backdrop" @click="onBackdropClick">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Keyboard Shortcuts</h2>
          <button class="close-btn" @click="emit('close')">×</button>
        </div>
        <div class="shortcuts-list">
          <template v-for="(item, index) in shortcuts" :key="index">
            <div v-if="item.section" class="shortcut-section">
              {{ item.section }}
            </div>
            <div v-else class="shortcut-item">
              <kbd>{{ item.key }}</kbd>
              <span>{{ item.description }}</span>
            </div>
          </template>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dark .modal-backdrop {
  background: rgba(0, 0, 0, 0.7);
}

.modal-content {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 8px;
  width: 90%;
  max-width: 450px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
  color: var(--text-h);
}

.close-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: var(--text);
  font-size: 24px;
  line-height: 1;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.15s;
}

.close-btn:hover {
  background: var(--code-bg);
}

.shortcuts-list {
  padding: 16px 20px;
  overflow-y: auto;
}

.shortcut-section {
  margin: 12px 0 8px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text);
  opacity: 0.7;
}

.shortcut-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 0;
}

kbd {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 80px;
  padding: 4px 8px;
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 4px;
  font-family: var(--mono);
  font-size: 12px;
  color: var(--text-h);
  text-align: center;
}

.shortcut-item span {
  font-size: 14px;
  color: var(--text);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
