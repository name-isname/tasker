<script setup lang="ts">
import type { Process } from '../types/api'
import StatusBadge from './StatusBadge.vue'

const props = defineProps<{
  process: Process
  selected?: boolean
}>()

const emit = defineEmits<{
  click: []
}>()

const priorityConfig: Record<string, { label: string; class: string }> = {
  low: { label: 'L', class: 'priority-low' },
  medium: { label: 'M', class: 'priority-medium' },
  high: { label: 'H', class: 'priority-high' },
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return 'Today'
  if (days === 1) return 'Yesterday'
  if (days < 7) return `${days}d ago`
  return date.toLocaleDateString()
}

const priority = priorityConfig[props.process.priority] || { label: '?', class: 'priority-unknown' }
</script>

<template>
  <div
    :class="['process-card', { selected }]"
    @click="emit('click')"
  >
    <div class="card-header">
      <h3 class="card-title">{{ process.title }}</h3>
      <div class="card-meta">
        <span :class="['priority-indicator', priority.class]" :title="`Priority: ${process.priority}`">
          {{ priority.label }}
        </span>
        <StatusBadge :status="process.status" />
      </div>
    </div>
    <p v-if="process.description" class="card-description">
      {{ process.description.slice(0, 100) }}{{ process.description.length > 100 ? '...' : '' }}
    </p>
    <div class="card-footer">
      <span class="card-date">{{ formatDate(process.created_at) }}</span>
      <span v-if="process.parent_id" class="card-parent">Sub-process</span>
    </div>
  </div>
</template>

<style scoped>
.process-card {
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg);
  cursor: pointer;
  transition: all 0.15s ease;
}

.process-card:hover {
  border-color: var(--accent);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.dark .process-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.process-card.selected {
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-bg);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 8px;
}

.card-title {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-h);
  line-height: 1.4;
  flex: 1;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.priority-indicator {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
}

.priority-low {
  background: rgba(107, 114, 128, 0.15);
  color: rgb(107, 114, 128);
}

.priority-medium {
  background: rgba(59, 130, 246, 0.15);
  color: rgb(59, 130, 246);
}

.priority-high {
  background: rgba(239, 68, 68, 0.15);
  color: rgb(239, 68, 68);
}

.dark .priority-low {
  background: rgba(107, 114, 128, 0.2);
  color: rgb(156, 163, 175);
}

.dark .priority-medium {
  background: rgba(59, 130, 246, 0.2);
  color: rgb(96, 165, 250);
}

.dark .priority-high {
  background: rgba(239, 68, 68, 0.2);
  color: rgb(248, 113, 113);
}

.priority-unknown {
  background: rgba(107, 114, 128, 0.15);
  color: rgb(107, 114, 128);
}

.card-description {
  margin: 0 0 10px;
  font-size: 13px;
  color: var(--text);
  line-height: 1.5;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--text);
}

.card-date {
  opacity: 0.7;
}

.card-parent {
  padding: 2px 6px;
  background: var(--code-bg);
  border-radius: 3px;
  font-size: 11px;
}
</style>
