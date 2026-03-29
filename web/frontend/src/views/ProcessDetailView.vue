<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { processApi } from '../services/api'
import { useKeyboard } from '../composables/useKeyboard'
import type { Process, Log, ProcessStatus } from '../types/api'
import StatusBadge from '../components/StatusBadge.vue'
import LogList from '../components/LogList.vue'
import AddLogForm from '../components/AddLogForm.vue'

const route = useRoute()
const router = useRouter()
const { register, unregister } = useKeyboard()

const process = ref<Process | null>(null)
const logs = ref<Log[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const updatingStatus = ref(false)

const processId = computed(() => Number(route.params.id))

const statusOptions: Array<{ value: ProcessStatus; label: string; key: string }> = [
  { value: 'running', label: 'Running', key: '1' },
  { value: 'blocked', label: 'Blocked', key: '2' },
  { value: 'suspended', label: 'Suspended', key: '3' },
  { value: 'terminated', label: 'Terminated', key: '4' },
]

const fetchProcess = async () => {
  loading.value = true
  error.value = null
  try {
    const [processData, logsData] = await Promise.all([
      processApi.get(processId.value),
      processApi.getLogs(processId.value),
    ])
    process.value = processData
    logs.value = logsData
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load process'
  } finally {
    loading.value = false
  }
}

const changeStatus = async (status: ProcessStatus) => {
  if (!process.value || updatingStatus.value) return
  if (process.value.status === status) return

  updatingStatus.value = true
  try {
    await processApi.updateStatus(processId.value, { status })
    await fetchProcess()
  } catch (e) {
    console.error('Failed to update status:', e)
  } finally {
    updatingStatus.value = false
  }
}

const addLog = async (content: string) => {
  if (!process.value) return

  try {
    await processApi.addLog(processId.value, {
      log_type: 'progress',
      content,
    })
    await fetchProcess()
  } catch (e) {
    console.error('Failed to add log:', e)
  }
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString()
}

const handleKeydown = (e: KeyboardEvent) => {
  switch (e.key) {
    case 'l':
      e.preventDefault()
      // Focus add log form
      window.dispatchEvent(new CustomEvent('focus-log-input'))
      break
    case '1':
      e.preventDefault()
      changeStatus('running')
      break
    case '2':
      e.preventDefault()
      changeStatus('blocked')
      break
    case '3':
      e.preventDefault()
      changeStatus('suspended')
      break
    case '4':
      e.preventDefault()
      changeStatus('terminated')
      break
  }
}

onMounted(() => {
  fetchProcess()
  register('keydown', handleKeydown)

  window.addEventListener('focus-log-input', () => {
    logFormRef.value?.focus()
  })
})

onUnmounted(() => {
  unregister('keydown', handleKeydown)
})

const logFormRef = ref<InstanceType<typeof AddLogForm>>()
</script>

<template>
  <div class="process-detail-view">
    <div class="view-header">
      <button class="back-btn" @click="router.back()">
        ← Back
      </button>
      <div class="header-actions">
        <button
          class="btn-delete"
          :disabled="updatingStatus"
          @click="changeStatus('terminated')"
        >
          Terminate
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">
      Loading...
    </div>

    <div v-else-if="error" class="error">
      {{ error }}
    </div>

    <div v-else-if="process" class="detail-content">
      <!-- Process Info -->
      <div class="info-section">
        <div class="info-header">
          <h1 class="process-title">{{ process.title }}</h1>
          <StatusBadge :status="process.status" />
        </div>

        <div class="info-meta">
          <span class="meta-item">
            <span class="meta-label">Priority:</span>
            <span class="meta-value">{{ process.priority }}</span>
          </span>
          <span class="meta-item">
            <span class="meta-label">Created:</span>
            <span class="meta-value">{{ formatDate(process.created_at) }}</span>
          </span>
          <span class="meta-item">
            <span class="meta-label">Updated:</span>
            <span class="meta-value">{{ formatDate(process.updated_at) }}</span>
          </span>
        </div>

        <div v-if="process.description" class="description">
          <h3>Description</h3>
          <p>{{ process.description }}</p>
        </div>
      </div>

      <!-- Status Actions -->
      <div class="status-section">
        <h3>Status</h3>
        <div class="status-actions">
          <button
            v-for="option in statusOptions"
            :key="option.value"
            :class="['status-btn', {
              active: process.status === option.value,
              disabled: updatingStatus
            }]"
            :disabled="updatingStatus"
            @click="changeStatus(option.value)"
          >
            <span class="status-key">{{ option.key }}</span>
            {{ option.label }}
          </button>
        </div>
      </div>

      <!-- Logs -->
      <div class="logs-section">
        <h3>Logs</h3>
        <AddLogForm
          ref="logFormRef"
          :loading="updatingStatus"
          @submit="addLog"
        />
        <LogList :logs="logs" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.process-detail-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 1px solid var(--border);
}

.back-btn {
  padding: 6px 12px;
  background: transparent;
  color: var(--text-h);
  border: 1px solid var(--border);
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.15s;
}

.back-btn:hover {
  background: var(--code-bg);
}

.btn-delete {
  padding: 6px 12px;
  background: rgba(239, 68, 68, 0.15);
  color: rgb(239, 68, 68);
  border: 1px solid transparent;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-delete:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.25);
  border-color: rgb(239, 68, 68);
}

.btn-delete:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading,
.error {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text);
  font-size: 14px;
}

.error {
  color: rgb(239, 68, 68);
}

.detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-section {
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg);
}

.info-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.process-title {
  margin: 0;
  font-size: 22px;
  font-weight: 500;
  color: var(--text-h);
  line-height: 1.3;
}

.info-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border);
}

.meta-item {
  display: flex;
  gap: 6px;
  font-size: 13px;
}

.meta-label {
  color: var(--text);
  opacity: 0.7;
}

.meta-value {
  color: var(--text-h);
  font-weight: 500;
}

.description h3 {
  margin: 0 0 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-h);
}

.description p {
  margin: 0;
  font-size: 14px;
  color: var(--text);
  line-height: 1.6;
  white-space: pre-wrap;
}

.status-section h3,
.logs-section h3 {
  margin: 0 0 12px;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-h);
}

.status-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.status-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: var(--bg);
  color: var(--text-h);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s;
}

.status-btn:hover:not(.disabled) {
  border-color: var(--accent);
}

.status-btn.active {
  background: var(--accent);
  color: white;
  border-color: var(--accent);
}

.status-btn.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.status-key {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
  font-size: 11px;
  font-family: var(--mono);
}

.dark .status-key {
  background: rgba(255, 255, 255, 0.15);
}

.logs-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
