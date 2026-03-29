<script setup lang="ts">
import type { Log } from '../types/api'

const props = defineProps<{
  logs: Log[]
}>()

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString()
}

const getLogIcon = (type: string) => {
  return type === 'state_change' ? '🔄' : '📝'
}

const getLogTypeLabel = (type: string) => {
  return type === 'state_change' ? 'State Change' : 'Progress'
}
</script>

<template>
  <div class="log-list">
    <div v-if="logs.length === 0" class="log-empty">
      No logs yet
    </div>
    <div v-for="log in logs" :key="log.id" class="log-entry">
      <div class="log-header">
        <span class="log-icon">{{ getLogIcon(log.log_type) }}</span>
        <span class="log-type">{{ getLogTypeLabel(log.log_type) }}</span>
        <span class="log-date">{{ formatDate(log.created_at) }}</span>
      </div>
      <div class="log-content">{{ log.content }}</div>
    </div>
  </div>
</template>

<style scoped>
.log-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-empty {
  padding: 20px;
  text-align: center;
  color: var(--text);
  opacity: 0.7;
  font-size: 14px;
}

.log-entry {
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg);
}

.log-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.log-icon {
  font-size: 14px;
}

.log-type {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text);
  opacity: 0.7;
}

.log-date {
  margin-left: auto;
  font-size: 12px;
  color: var(--text);
  opacity: 0.5;
}

.log-content {
  font-size: 14px;
  color: var(--text-h);
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
