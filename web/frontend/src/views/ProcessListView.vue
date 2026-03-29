<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProcessesStore } from '../stores/processes'
import { useKeyboard } from '../composables/useKeyboard'
import type { ProcessStatus } from '../types/api'
import ProcessCard from '../components/ProcessCard.vue'

const router = useRouter()
const processesStore = useProcessesStore()
const { register, unregister } = useKeyboard()

const selectedId = ref<number | null>(null)
const deleteConfirm = ref<number | null>(null)

const statuses: Array<{ value: ProcessStatus | 'all'; label: string }> = [
  { value: 'all', label: 'All' },
  { value: 'running', label: 'Running' },
  { value: 'blocked', label: 'Blocked' },
  { value: 'suspended', label: 'Suspended' },
  { value: 'terminated', label: 'Terminated' },
]

const filteredProcesses = computed(() => {
  return processesStore.processes
})

const selectedIndex = computed(() => {
  return filteredProcesses.value.findIndex((p) => p.id === selectedId.value)
})

const selectNext = () => {
  if (filteredProcesses.value.length === 0) return
  const idx = selectedIndex.value
  const nextIdx = idx < filteredProcesses.value.length - 1 ? idx + 1 : 0
  selectedId.value = filteredProcesses.value[nextIdx].id
}

const selectPrev = () => {
  if (filteredProcesses.value.length === 0) return
  const idx = selectedIndex.value
  const prevIdx = idx > 0 ? idx - 1 : filteredProcesses.value.length - 1
  selectedId.value = filteredProcesses.value[prevIdx].id
}

const openSelected = () => {
  if (selectedId.value) {
    router.push(`/process/${selectedId.value}`)
  }
}

const deleteSelected = async () => {
  if (selectedId.value) {
    deleteConfirm.value = selectedId.value
  }
}

const confirmDelete = async () => {
  if (deleteConfirm.value) {
    try {
      await processesStore.deleteProcess(deleteConfirm.value)
      deleteConfirm.value = null
      selectedId.value = null
    } catch (e) {
      console.error('Failed to delete process:', e)
    }
  }
}

const editSelected = () => {
  // TODO: Implement edit
  console.log('Edit not implemented yet')
}

const openCreateModal = () => {
  window.dispatchEvent(new CustomEvent('open-create-modal'))
}

const handleKeydown = (e: KeyboardEvent) => {
  if (deleteConfirm.value) return

  switch (e.key) {
    case 'c':
      e.preventDefault()
      window.dispatchEvent(new CustomEvent('open-create-modal'))
      break
    case 'j':
    case 'ArrowDown':
      e.preventDefault()
      selectNext()
      break
    case 'k':
    case 'ArrowUp':
      e.preventDefault()
      selectPrev()
      break
    case 'Enter':
      e.preventDefault()
      openSelected()
      break
    case 'd':
      e.preventDefault()
      deleteSelected()
      break
    case 'e':
      e.preventDefault()
      editSelected()
      break
  }
}

onMounted(() => {
  processesStore.fetchProcesses()
  register('keydown', handleKeydown)
})

onUnmounted(() => {
  unregister('keydown', handleKeydown)
})

const setStatusFilter = (status: ProcessStatus | 'all') => {
  processesStore.setStatusFilter(status)
  selectedId.value = null
}

const openProcess = (processId: number) => {
  selectedId.value = processId
  router.push(`/process/${processId}`)
}
</script>

<template>
  <div class="process-list-view">
    <div class="view-header">
      <div class="header-left">
        <h2>Processes</h2>
        <span class="count">{{ filteredProcesses.length }}</span>
      </div>
      <button class="btn-primary" @click="openCreateModal">
        <span class="btn-icon">+</span>
        New Process
      </button>
    </div>

    <div class="filter-bar">
      <div class="filter-group">
        <span class="filter-label">Status:</span>
        <button
          v-for="status in statuses"
          :key="status.value"
          :class="['filter-btn', { active: processesStore.statusFilter === status.value }]"
          @click="setStatusFilter(status.value)"
        >
          {{ status.label }}
        </button>
      </div>
    </div>

    <div v-if="processesStore.loading" class="loading">
      Loading...
    </div>

    <div v-else-if="processesStore.error" class="error">
      {{ processesStore.error }}
    </div>

    <div v-else-if="filteredProcesses.length === 0" class="empty">
      <p>No processes found.</p>
      <p class="empty-hint">Press <kbd>c</kbd> to create a new process.</p>
    </div>

    <div v-else class="process-list">
      <ProcessCard
        v-for="process in filteredProcesses"
        :key="process.id"
        :process="process"
        :selected="process.id === selectedId"
        @click="openProcess(process.id)"
      />
    </div>

    <!-- Delete Confirmation -->
    <Transition name="fade">
      <div v-if="deleteConfirm" class="modal-backdrop" @click.self="deleteConfirm = null">
        <div class="confirm-dialog">
          <h3>Confirm Delete</h3>
          <p>Are you sure you want to delete this process? This will also delete all sub-processes.</p>
          <div class="dialog-actions">
            <button class="btn-secondary" @click="deleteConfirm = null">Cancel</button>
            <button class="btn-danger" @click="confirmDelete">Delete</button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.process-list-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.header-left h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 500;
  color: var(--text-h);
}

.count {
  font-size: 14px;
  color: var(--text);
  opacity: 0.7;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: var(--accent);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: filter 0.15s;
}

.btn-primary:hover {
  filter: brightness(1.1);
}

.btn-icon {
  font-size: 16px;
  line-height: 1;
}

.filter-bar {
  padding: 12px 20px;
  border-bottom: 1px solid var(--border);
  background: var(--code-bg);
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
}

.filter-btn {
  padding: 4px 10px;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: var(--bg);
  color: var(--text);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
}

.filter-btn:hover {
  border-color: var(--accent);
}

.filter-btn.active {
  background: var(--accent);
  color: white;
  border-color: var(--accent);
}

.process-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 12px;
  align-content: start;
}

.loading,
.error,
.empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text);
  font-size: 14px;
}

.error {
  color: rgb(239, 68, 68);
}

.empty-hint {
  margin-top: 8px;
  font-size: 13px;
  opacity: 0.7;
}

.empty-hint kbd {
  padding: 2px 6px;
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 3px;
  font-family: var(--mono);
  font-size: 11px;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.confirm-dialog {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 20px;
  width: 90%;
  max-width: 400px;
}

.confirm-dialog h3 {
  margin: 0 0 12px;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-h);
}

.confirm-dialog p {
  margin: 0 0 20px;
  font-size: 14px;
  color: var(--text);
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-secondary,
.btn-danger {
  padding: 8px 14px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  transition: filter 0.15s;
}

.btn-secondary {
  background: transparent;
  color: var(--text-h);
  border: 1px solid var(--border);
}

.btn-secondary:hover {
  background: var(--code-bg);
}

.btn-danger {
  background: rgb(239, 68, 68);
  color: white;
}

.btn-danger:hover {
  filter: brightness(1.1);
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
